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
	"github.com/nathanhack/redisbadger/command"
	"github.com/nathanhack/redisbadger/command/commandname"
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
	txn       *badger.Txn
	it        *badger.Iterator
	offset    uint64
	match     string
	glob      glob.Glob
	mux       sync.Mutex
	ctx       context.Context
	ctxCancel context.CancelFunc
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
		dbMux := sync.Mutex{}
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

		if backup {
			backupChan = make(chan *aof.Command, 100)
			defer close(backupChan)

			wg.Add(1)
			go func() {
				backupRoutine(ctx, backupPathname, backupChan)
				wg.Done()
			}()
		}

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
			func(conn redcon.Conn, redisCmd redcon.Command) {

				cmd, err := command.ParseCommand(redisCmd.Args)
				if err != nil {
					logrus.Debugf("ParseCommand: %v", err)
					conn.WriteError(fmt.Sprintf("ERR :%v", err))
					return
				}

				switch cmd.Cmd {
				default:
					conn.WriteError("ERR unimplemented command '" + string(cmd.Cmd) + "'")
				case commandname.BgRewriteAOF:
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
				case commandname.Ping:
					//PING [message]
					switch len(cmd.Args) {
					case 0:
						conn.WriteString("PONG")
					case 1:
						conn.WriteBulk([]byte(cmd.Args[0]))
					default:
						conn.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command found: %v", cmd.Cmd, cmd.Args))
					}
				case commandname.Quit:
					//QUIT
					conn.WriteString("OK")
					conn.Close()
				case commandname.Set:
					//SET key value [EX seconds|PX milliseconds|EXAT timestamp|PXAT milliseconds-timestamp|KEEPTTL] [NX|XX] [GET]

					//first we will go through the non supported options/flags
					for _, k := range []string{"EX", "PX", "EXAT", "PXAT"} {
						_, has := scanOptions(k, cmd.Options)
						if has {
							conn.WriteError(fmt.Sprintf("Error the option '%v' is not implemented for command '%v'", k, cmd.Cmd))
							return
						}
					}

					for _, k := range []string{"KEEPTTL", "NX", "XX"} {
						if scanFlags(k, cmd.Options) {
							conn.WriteError(fmt.Sprintf("Error the option flagn '%v' is not implemented for command '%v'", k, cmd.Cmd))
							return
						}
					}

					_, hasGet := scanOptions("GET", cmd.Options)
					dbMux.Lock()
					var valCopy []byte
					err := db.Update(func(txn *badger.Txn) error {
						logrus.Debugf("cmd: %#v", cmd)
						if hasGet {
							item, err := txn.Get([]byte(cmd.Args[0]))
							if err == nil {
								valCopy, err = item.ValueCopy(nil)
								if err != nil {
									return err
								}
							}
						}

						return txn.Set([]byte(cmd.Args[0]), []byte(cmd.Args[1]))
					})
					dbMux.Unlock()
					if err != nil {
						logrus.Error(err)
						conn.WriteError(fmt.Sprintf("Error %v", err))
						return
					}

					if backup {
						backupChan <- &aof.Command{
							Name:      string(commandname.Set),
							Arguments: cmd.Args,
						}
					}

					if hasGet {
						if valCopy == nil {
							conn.WriteNull()
						} else {
							conn.WriteBulk(valCopy)
						}
					} else {
						conn.WriteString("OK")
					}
				case commandname.Get:
					//GET key
					var valCopy []byte
					err := db.View(func(txn *badger.Txn) error {
						item, err := txn.Get([]byte(cmd.Args[0]))
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

				case commandname.Del:
					//DEL key [key ...]
					deleted := 0
					for _, key := range cmd.Args {
						dbMux.Lock()
						err := db.Update(func(txn *badger.Txn) error {
							return txn.Delete([]byte(key))
						})
						dbMux.Unlock()
						if err == nil {
							deleted++
						}

						if backup {
							backupChan <- &aof.Command{
								Name:      string(commandname.Del),
								Arguments: []string{key},
							}
						}
					}
					conn.WriteInt(deleted)

				case commandname.Scan:
					//SCAN cursor [MATCH pattern] [COUNT count] [TYPE type]

					_, hasType := scanOptions("TYPE", cmd.Options)
					if hasType {
						conn.WriteError(fmt.Sprintf("ERR TYPE not supported for %s", cmd.Cmd))
						return
					}

					match, hasMatch := scanOptions("MATCH", cmd.Options)
					count, hasCount := scanOptions("COUNT", cmd.Options)

					if logrus.GetLevel() == logrus.DebugLevel {
						logrus.Debugln(strings.Join(cmd.Args, " "))
					}

					cursorValue, err := strconv.ParseUint(string(cmd.Args[0]), 10, 64)
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
						num, err := strconv.Atoi(count)
						if err != nil {
							conn.WriteError(fmt.Sprintf("ERR COUNT value (required >=0) was not parsable: %v", err))
							return
						}
						if maxCount > num {
							maxCount = num
						}
					}

					activeScansMux.Lock()
					scan, has := activeScans[conn.RemoteAddr()]
					activeScansMux.Unlock()

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
							logrus.Debugf("command:SCAN: glob.Compile('%v') -> %v", match, err)
							conn.WriteError(fmt.Sprintf("ERR MATCH string was not vaild glob syntax: %v", err))
							return
						}

						scan = &scannerState{
							txn:    txn,
							it:     txn.NewIterator(opts),
							offset: cursorValue,
							match:  match,
							glob:   pattern,
							mux:    sync.Mutex{},
						}

						scan.ctx, scan.ctxCancel = context.WithCancel(context.Background())

						activeScansMux.Lock()
						activeScans[conn.RemoteAddr()] = scan
						activeScansMux.Unlock()

						scan.it.Rewind()
						for i := uint64(0); i < cursorValue; i++ {
							scan.it.Next()

							// this loop could take a while we need to check the context
							select {
							case <-ctx.Done():
								return
							case <-scan.ctx.Done():
								return
							default:
							}
						}
					}

					scan.mux.Lock()

					keys := make([][]byte, 0)
					for ; scan.it.Valid() && len(keys) < maxCount; scan.it.Next() {
						tmp := scan.it.Item().Key()
						if scan.glob.Match(string(tmp)) {
							keys = append(keys, scan.it.Item().KeyCopy(nil))
						}
						scan.offset++

						// this loop could take a while we need to check the context
						select {
						case <-ctx.Done():
							return
						case <-scan.ctx.Done():
							return
						default:
						}
					}

					nextOffset := scan.offset
					if !scan.it.Valid() {
						//clean up since we made it to the end
						scan.Close()
						activeScansMux.Lock()
						delete(activeScans, conn.RemoteAddr())
						activeScansMux.Unlock()
						nextOffset = 0
					}
					scan.mux.Unlock()

					//well now we have data or possibly not, in either case we write it out
					conn.WriteArray(2)
					conn.WriteBulkString(fmt.Sprint(nextOffset))
					conn.WriteArray(len(keys))
					for _, key := range keys {
						conn.WriteBulkString(string(key))
					}
				case commandname.Publish:
					conn.WriteInt(ps.Publish(string(cmd.Args[0]), string(cmd.Args[1])))
				case commandname.Subscribe, commandname.PSubscribe:
					for _, arg := range cmd.Args {
						if cmd.Cmd == commandname.PSubscribe {
							ps.Psubscribe(conn, arg)
						} else {
							ps.Subscribe(conn, arg)
						}
					}
				}
			},
			func(conn redcon.Conn) bool {
				// Use this function to accept or deny the connection.
				return true
			},
			func(conn redcon.Conn, err error) {
				// This is called when the connection has been closed
				logrus.Printf("closed: %s, err: %v", conn.RemoteAddr(), err)
				activeScansMux.Lock()
				scan, has := activeScans[conn.RemoteAddr()]
				activeScansMux.Unlock()
				if has {
					scan.ctxCancel()
					logrus.Infof("active scan canceled")
				}
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
		if !strings.EqualFold(command.Name, string(commandname.Set)) && !strings.EqualFold(command.Name, string(commandname.Del)) {
			continue
		}

		err := db.Update(func(txn *badger.Txn) error {
			_, err := txn.Get([]byte(command.Arguments[0]))
			if errors.Is(err, badger.ErrKeyNotFound) {
				if strings.EqualFold(command.Name, string(commandname.Set)) {
					added++
				} else {
					return nil
				}
			}

			//older commands should be applied
			// so we run anyway
			if strings.EqualFold(command.Name, string(commandname.Set)) {
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

			if strings.EqualFold(command.Name, string(commandname.Set)) {
				backupChan <- &aof.Command{
					Name:      string(commandname.Set),
					Arguments: []string{string(command.Arguments[0]), string(command.Arguments[1])},
				}
			} else {
				backupChan <- &aof.Command{
					Name:      string(commandname.Del),
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
				Name:      string(commandname.Set),
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

func scanOptions(target string, options []command.Option) (optionValue string, found bool) {
	for _, o := range options {
		if strings.EqualFold(target, o.Name) {
			return o.Value, true
		}
	}
	return "", false
}

func scanFlags(target string, options []command.Option) (found bool) {
	for _, o := range options {
		if strings.EqualFold(target, o.Name) {
			return true
		}
	}
	return false
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.redisbadger.yaml)")
	rootCmd.Flags().StringVar(&addr, "address", ":6379", "the address and port to listen for redis commands")
	rootCmd.Flags().StringVar(&databasePathname, "database", "./badger", "the directory that will store the badger database")
	rootCmd.Flags().BoolVar(&backup, "backup", false, "enables an AOF backup file in the database directory")
	rootCmd.Flags().StringVar(&load, "load", "", "before starting the server load this specific AOF")

	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enables debug logging")
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
