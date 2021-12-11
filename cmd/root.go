package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/gobwas/glob"
	"github.com/gol4ng/signal"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/nathanhack/aof"
	"github.com/nathanhack/redcon/v2"
	"github.com/nathanhack/redisbadger/commands"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const aofAppendOnlyFile = "appendonly.aof"

var cfgFile string
var addr string
var databasePathname string
var debug bool
var backup bool
var load string
var backupChan chan *aof.Command
var backupPathname string

type scannerState struct {
	txn    *badger.Txn
	it     *badger.Iterator
	offset uint64
	match  string
	glob   glob.Glob
}

func (ss *scannerState) Close() {
	if ss.it != nil {
		ss.it.Close()
		ss.it = nil
	}
	if ss.txn != nil {
		ss.txn.Discard()
		ss.txn = nil
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "redisbadger",
	Short: "Starts up a redis compatible server backed by badger",
	Long:  `Starts up a redis compatible server backed by badger`,
	RunE: func(cmd *cobra.Command, args []string) error {
		wg := sync.WaitGroup{}
		ctx, cancel := context.WithCancel(context.Background())

		defer signal.Subscribe(func(signal os.Signal) {
			logrus.Warnf("Ctrl-c pressed closing gracefully")
			cancel()
		}, os.Interrupt, syscall.SIGTERM)()

		var ps redcon.PubSub
		activeScans := map[string]*scannerState{}
		activeScansMux := sync.Mutex{}

		if debug {
			logrus.SetLevel(logrus.DebugLevel)
			logrus.Info("Debug enabled")
		}

		db, err := badger.Open(badger.DefaultOptions(databasePathname))
		if err != nil {
			err = fmt.Errorf("error while opening badgerDB at %v:%v", databasePathname, err)
			logrus.Error(err)
			return err
		}
		defer func() {
			cancel()
			wg.Wait()
			fmt.Println("database closed")
			db.Close()
		}()

		backupPathname, err = filepath.Abs(filepath.Join(databasePathname, aofAppendOnlyFile))
		if err != nil {
			logrus.Error(err)
			return err
		}
		backupChan = make(chan *aof.Command, 100)
		defer close(backupChan)

		wg.Add(1)
		go func() {
			backupRoutine(ctx, backupPathname, backupChan)
			wg.Done()
		}()

		if load != "" {
			if _, err := os.Stat(load); errors.Is(err, os.ErrNotExist) {
				err = fmt.Errorf("load file does not exist: %v", err)
				logrus.Error(err)
				return err
			}

			logrus.Infof("loading data from %v", load)
			absLoadFile, err := filepath.Abs(load)
			if err != nil {
				logrus.Error(err)
				return err
			}
			err = aofLoad(ctx, absLoadFile, db)
			if err != nil {
				err = fmt.Errorf("error while loading AOF: %v", err)
				logrus.Error(err)
				return err
			}
		}

		aofBackupActive := false // should really be mux protected but this shouldn't be called that frequently
		logrus.Printf("started server at %s", addr)
		err = redcon.ListenAndServe(ctx, addr,
			func(conn redcon.Conn, cmd redcon.Command) {
				command := strings.ToUpper(string(cmd.Args[0]))
				switch command {
				default:
					conn.WriteError("ERR unknown command '" + string(cmd.Args[0]) + "'")
				case commands.BgRewriteAOF:
					if aofBackupActive {
						conn.WriteError("ERR already running")
						return
					}
					aofBackupActive = true
					go func() {
						aofBackup(ctx, filepath.Join(databasePathname, aofAppendOnlyFile+".bak"), db)
						aofBackupActive = false
					}()
					conn.WriteString("ASAP")
				case commands.Ping:
					//PING [message]
					switch len(cmd.Args) {
					case 1:
						conn.WriteString("PONG")
					case 2:
						conn.WriteBulk(cmd.Args[1])
					default:
						conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", cmd.Args[0]))
					}
				case commands.Quit:
					//QUIT
					conn.WriteString("OK")
					conn.Close()
				case commands.Set:
					//SET key value [EX seconds|PX milliseconds|EXAT timestamp|PXAT milliseconds-timestamp|KEEPTTL] [NX|XX] [GET]
					if len(cmd.Args) != 3 {
						conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
						return
					}
					err := db.Update(func(txn *badger.Txn) error {
						logrus.Debugf("SET %s %s", cmd.Args[1], cmd.Args[2])
						err := txn.Set(cmd.Args[1], cmd.Args[2])
						return err
					})

					if err != nil {
						logrus.Error(err)
						conn.WriteError(fmt.Sprintf("Error %v", err))
						return
					}

					if backup {
						backupChan <- &aof.Command{
							Name:      commands.Set,
							Arguments: []string{string(cmd.Args[1]), string(cmd.Args[2])},
						}
					}
					conn.WriteString("OK")

				case commands.Get:
					//GET key
					if len(cmd.Args) != 2 {
						conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
						return
					}

					key := cmd.Args[1]
					var valCopy []byte
					err := db.View(func(txn *badger.Txn) error {
						item, err := txn.Get(key)
						if err == nil {
							valCopy, err = item.ValueCopy(nil)
						}
						return err
					})

					if err != nil {
						conn.WriteNull()
					} else {
						conn.WriteBulk(valCopy)
					}

				case commands.Del:
					//DEL key [key ...]
					if len(cmd.Args) < 2 {
						conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
						return
					}

					deleted := 0
					for _, key := range cmd.Args[1:] {
						err := db.Update(func(txn *badger.Txn) error {
							return txn.Delete(key)
						})
						if err == nil {
							deleted++
						}
						if backup {
							backupChan <- &aof.Command{
								Name:      commands.Del,
								Arguments: []string{string(key)},
							}
						}
					}
					conn.WriteInt(deleted)

				case commands.Scan:
					//SCAN cursor [MATCH pattern] [COUNT count] [TYPE type]

					if len(cmd.Args) < 2 || 8 < len(cmd.Args) {
						conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
						return
					}

					argStrings := argsToStrings(cmd.Args)
					_, hasType := scanFlags(argStrings, "TYPE")
					match, hasMatch := scanFlags(argStrings, "MATCH")
					count, hasCount := scanFlags(argStrings, "COUNT")

					if hasType {
						conn.WriteError(fmt.Sprintf("ERR TYPE not supported for %s", cmd.Args[0]))
						return
					}

					if logrus.GetLevel() == logrus.DebugLevel {
						sb := strings.Builder{}
						for _, arg := range cmd.Args {
							sb.WriteString(string(arg))
							sb.WriteString(" ")
						}
						logrus.Debugln(sb.String())
					}

					cursorValue, err := strconv.ParseUint(string(cmd.Args[1]), 10, 64)
					if err != nil {
						conn.WriteError(fmt.Sprintf("ERR cursor (required >=0) was not parsable: %v", err))
						return
					}

					// if no match is found use default
					if !hasMatch {
						match = "*"
					}

					maxCount := badger.DefaultIteratorOptions.PrefetchSize
					if hasCount {
						num, err := strconv.ParseInt(count, 10, 64)
						if err != nil {
							conn.WriteError(fmt.Sprintf("ERR COUNT value (required >=0) was not parsable: %v", err))
							return
						}
						if int64(maxCount) > num {
							maxCount = int(num)
						}
					}

					activeScansMux.Lock()

					scan, has := activeScans[conn.RemoteAddr()]

					if !has ||
						(scan != nil && scan.offset != cursorValue) ||
						(scan != nil && scan.match != match) {

						if scan != nil {
							scan.Close()

							//if we switch the match then
							// we need to restart from zero
							if scan.match != match {
								cursorValue = 0
							}
						}
						txn := db.NewTransaction(false)
						opts := badger.DefaultIteratorOptions
						opts.PrefetchValues = false

						pattern, err := glob.Compile(match)
						if err != nil {
							conn.WriteError(fmt.Sprintf("ERR MATCH string was not vaild glob syntax: %v", err))
							return
						}

						scan = &scannerState{
							txn:    txn,
							it:     txn.NewIterator(opts),
							offset: cursorValue,
							match:  match,
							glob:   pattern,
						}

						activeScans[conn.RemoteAddr()] = scan

						scan.it.Rewind()
						for i := uint64(0); i < cursorValue; i++ {
							scan.it.Next()
						}
					}

					keys := make([][]byte, 0)
					for ; scan.it.Valid() && len(keys) < maxCount; scan.it.Next() {
						tmp := scan.it.Item().Key()
						if scan.glob.Match(string(tmp)) {
							keys = append(keys, scan.it.Item().KeyCopy(nil))
						}
						scan.offset++
					}

					nextOffset := scan.offset
					if !scan.it.Valid() {
						//clean up since we made it to the end
						scan.Close()
						delete(activeScans, conn.RemoteAddr())
						nextOffset = 0
					}

					activeScansMux.Unlock()

					//well now we have data or possibly not, in either case we write it out
					conn.WriteArray(2)
					conn.WriteBulkString(fmt.Sprint(nextOffset))
					conn.WriteArray(len(keys))
					for _, key := range keys {
						conn.WriteBulkString(string(key))
					}
				case commands.Publish:
					if len(cmd.Args) != 3 {
						conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
						return
					}
					conn.WriteInt(ps.Publish(string(cmd.Args[1]), string(cmd.Args[2])))
				case commands.Subscribe, commands.PSubscribe:
					if len(cmd.Args) < 2 {
						conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
						return
					}

					for i := 1; i < len(cmd.Args); i++ {
						if command == commands.PSubscribe {
							ps.Psubscribe(conn, string(cmd.Args[i]))
						} else {
							ps.Subscribe(conn, string(cmd.Args[i]))
						}
					}
				}
			},
			func(conn redcon.Conn) bool {
				// Use this function to accept or deny the connection.
				// log.Printf("accept: %s", conn.RemoteAddr())
				return true
			},
			func(conn redcon.Conn, err error) {
				// This is called when the connection has been closed
				logrus.Printf("closed: %s, err: %v", conn.RemoteAddr(), err)
			},
		)
		if err != nil {
			logrus.Error(err)
		}

		return nil
	},
}

func backupRoutine(ctx context.Context, backupPathname string, backupChan chan *aof.Command) {
	logrus.Infof("starting backup file: %v", backupPathname)
	var f *os.File
	_, err := os.Stat(load)
	if errors.Is(err, os.ErrNotExist) {
		f, err = os.OpenFile(backupPathname, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	} else {
		f, err = os.OpenFile(backupPathname, os.O_RDWR|os.O_APPEND, 0666)
	}

	if err != nil {
		panic(fmt.Sprintf("error while opening AOF file: %v", err))
	}
	defer f.Close()
	writer := bufio.NewWriter(f)

	ctx2, cancel := context.WithCancel(ctx)
	wg := sync.WaitGroup{}
	wg.Add(1)

	// we start a go routine to call f.Sync() every 1 second until done
	// so we ensure data is pushed to disk
	go func() {
		defer f.Sync()
		defer wg.Done()

		for {
			select {
			case <-ctx2.Done():
				return
			case <-time.After(time.Second):
			}
			f.Sync()
		}
	}()

	for c := range backupChan {
		if err := aof.WriteCommand(c, writer); err != nil {
			logrus.Errorf("backupRoutine: failed to write command: %v", err)
			continue
		}
		writer.Flush()
	}

	cancel()
	wg.Wait()
	logrus.Infof("backupRoutine finished")

}

func aofLoad(ctx context.Context, aofFilepath string, db *badger.DB) error {
	//here we blindly just load the AOF file
	f, err := os.Open(aofFilepath)
	if err != nil {
		err = fmt.Errorf("error while opening AOF file: %v", err)
		logrus.Fatal(err)
		return err
	}

	reader := bufio.NewReader(f)

	count := 0
	added := 0
	removed := 0
	logrus.Infof("AOF ingest started")
	for command, bs, err := aof.ReadCommand(reader); !errors.Is(err, io.EOF); command, bs, err = aof.ReadCommand(reader) {
		logrus.Debugf("aofLoad:ReadCommand: %v %v %v", command, bs, err)
		count++
		if count%1_000_000 == 0 {
			logrus.Infof("AOF ingest index: %v added: %v removed: %v", count, added, removed)
		}

		if err != nil {
			logrus.Errorf("error when loading AOF: %v during these bytes %v", err, bs)
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("context done")
		default:
		}

		//if we're not a Set or Del we skip
		if command.Name != commands.Set && command.Name != commands.Del {
			continue
		}

		err := db.Update(func(txn *badger.Txn) error {
			_, err := txn.Get([]byte(command.Arguments[0]))
			if errors.Is(err, badger.ErrKeyNotFound) {
				if command.Name == commands.Set {
					added++
				} else {
					return nil
				}
			}

			//older commands should be applied
			// so we run anyway
			if command.Name == commands.Set {
				return txn.Set([]byte(command.Arguments[0]), []byte(command.Arguments[1]))
			}

			removed++
			return txn.Delete([]byte(command.Arguments[0]))
		})

		if err != nil {
			err = fmt.Errorf("error while adding data to DB from AOF: %v", err)
			logrus.Error(err)
			return err
		}

		if backup {
			if aofFilepath == backupPathname {
				continue
			}

			if command.Name == commands.Set {
				backupChan <- &aof.Command{
					Name:      commands.Set,
					Arguments: []string{string(command.Arguments[0]), string(command.Arguments[1])},
				}
			} else {
				backupChan <- &aof.Command{
					Name:      commands.Del,
					Arguments: []string{string(command.Arguments[0])},
				}
			}
		}
	}
	logrus.Infof("finished loading %v command (keys added %v, removed %v) items from %v", count, added, removed, load)

	return nil
}

func aofBackup(ctx context.Context, aofFilepath string, db *badger.DB) error {
	// here the point is to back up what's in the badgerDB
	// we don't do anything efficient we simply
	// over right the AOF file with what's in the db at this moment
	logrus.Infof("creating backup AOF %v", aofFilepath)
	f, err := os.Create(aofFilepath)
	if err != nil {
		err = fmt.Errorf("error while opening AOF file: %v", err)
		logrus.Fatal(err)
		return err
	}
	defer f.Close()
	writer := bufio.NewWriter(f)
	count := 0
	err = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			select {
			case <-ctx.Done():
				return fmt.Errorf("context done")
			default:
			}

			count++
			item, err := txn.Get(it.Item().Key())
			if err != nil {
				return fmt.Errorf("error while getting value item for key(%v) :%v", string(it.Item().Key()), err)
			}
			value, err := item.ValueCopy(nil)
			if err != nil {
				return fmt.Errorf("error while getting value for key(%v) :%v", string(it.Item().Key()), err)
			}
			opt := &aof.Command{
				Name:      commands.Set,
				Arguments: []string{string(it.Item().KeyCopy(nil)), string(value)},
			}

			err = aof.WriteCommand(opt, writer)
			if err != nil {
				return fmt.Errorf("error while getting value for key(%v) :%v", string(it.Item().Key()), err)
			}
			err = writer.Flush()
			if err != nil {
				logrus.Warnf("error during flush: %v", err)
			}
		}
		logrus.Infof("AOF backup created with %v keys", count)
		return nil
	})
	if err != nil {
		return fmt.Errorf("error while determining badgerDB item count:%v", err)
	}
	return f.Sync()
}

func argsToStrings(args [][]byte) (results []string) {
	results = make([]string, len(args))
	for i, arg := range args {
		results[i] = string(arg)
	}
	return results
}

func scanFlags(args []string, flag string) (flagValue string, flagFound bool) {
	for i, arg := range args {
		if strings.ToUpper(flag) == strings.ToUpper(arg) && i < len(args)-1 {
			return args[i+1], true
		}
	}
	return "", false
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.redisbadger.yaml)")
	rootCmd.PersistentFlags().StringVar(&addr, "address", ":6379", "the address and port to listen for redis commands")
	rootCmd.PersistentFlags().StringVar(&databasePathname, "database", "./badger", "the directory that will store the badger database")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enables debug logging")
	rootCmd.PersistentFlags().BoolVar(&backup, "backup", false, "enables an AOF backup file in the database directory")
	rootCmd.PersistentFlags().StringVar(&load, "load", "", "before starting the server load this specific AOF")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".redisbadger" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".redisbadger")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
