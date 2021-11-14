package cmd

import (
	"fmt"
	badger "github.com/dgraph-io/badger/v3"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/nathanhack/redisbadger/commands"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tidwall/redcon"
	"os"
	"strconv"
	"strings"
	"sync"
)

var cfgFile string
var addr string
var databasePathname string
var debug bool

type scannerState struct {
	txn    *badger.Txn
	it     *badger.Iterator
	offset uint64
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
		var ps redcon.PubSub
		activeScans := map[string]*scannerState{}
		activeScansMux := sync.Mutex{}

		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}

		logrus.Printf("started server at %s", addr)
		db, err := badger.Open(badger.DefaultOptions(databasePathname))
		if err != nil {
			logrus.Fatal(err)
		}
		defer db.Close()
		err = redcon.ListenAndServe(addr,
			func(conn redcon.Conn, cmd redcon.Command) {

				command := strings.ToUpper(string(cmd.Args[0]))
				switch command {
				default:
					conn.WriteError("ERR unknoWn command '" + string(cmd.Args[0]) + "'")
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
					}
					conn.WriteInt(deleted)
				case commands.Scan:
					//SCAN cursor [MATCH pattern] [COUNT count] [TYPE type]

					if len(cmd.Args) != 2 {
						conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
						return
					}

					cursorValue, err := strconv.ParseUint(string(cmd.Args[1]), 10, 64)
					if err != nil {
						conn.WriteError(fmt.Sprintf("ERR cursor was not parsable %v", err))
						return
					}
					activeScansMux.Lock()
					scan, has := activeScans[conn.RemoteAddr()]
					if !has || (scan != nil && scan.offset != cursorValue) {
						if scan != nil {
							scan.Close()
							scan.offset = cursorValue
						}
						txn := db.NewTransaction(false)
						opts := badger.DefaultIteratorOptions
						opts.PrefetchValues = false
						scan = &scannerState{
							txn: txn,
							it:  txn.NewIterator(opts),
						}

						activeScans[conn.RemoteAddr()] = scan

						scan.it.Rewind()
						for i := uint64(0); i < uint64(badger.DefaultIteratorOptions.PrefetchSize)*cursorValue; i++ {
							scan.it.Next()
						}
					}

					keys := make([][]byte, 0)
					for ; scan.it.Valid() && len(keys) < badger.DefaultIteratorOptions.PrefetchSize; scan.it.Next() {
						keys = append(keys, scan.it.Item().KeyCopy(nil))
					}

					scan.offset++

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
			logrus.Fatal(err)
		}

		return nil
	},
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
