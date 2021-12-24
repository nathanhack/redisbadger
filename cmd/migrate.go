package cmd

import (
	"context"
	"fmt"
	"os"
	"syscall"

	"github.com/cheggaaa/pb"
	"github.com/go-redis/redis/v8"
	"github.com/gol4ng/signal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var move bool
var match string
var from string
var to string
var ctx, cancel = context.WithCancel(context.Background())

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrates data from one Redisbadger to another",
	Long:  `Migrates data from one Redisbadger to another.  Data can be copied or moved.`,
	Run: func(cmd *cobra.Command, args []string) {
		defer signal.Subscribe(func(signal os.Signal) {
			logrus.Warnf("Ctrl-c pressed closing gracefully")
			cancel()
		}, os.Interrupt, syscall.SIGTERM)()

		if move {
			logrus.Infof("MOVING data that matches '%v' from '%v' to '%v'", match, from, to)
		} else {
			logrus.Infof("COPYING data that matches '%v' from '%v' to '%v'", match, from, to)
		}

		rdb1 := redis.NewClient(&redis.Options{
			Addr:        from,
			Password:    "",
			DB:          0,
			ReadTimeout: -1,
		})

		rdb2 := redis.NewClient(&redis.Options{
			Addr:        to,
			Password:    "",
			DB:          0,
			ReadTimeout: -1,
		})

		var cursor uint64
		var keys []string
		var err error
		bar := pb.StartNew(-1)
		defer bar.Finish()

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			var nextcursor uint64
			keys, nextcursor, err = rdb1.Scan(ctx, cursor, match, 2000).Result()
			if err != nil {
				panic(err)
			}
			cursor = nextcursor

			for _, key := range keys {
				bar.Increment()

				select {
				case <-ctx.Done():
					logrus.Warn("Context done: all data may not have been migrated")
					return
				default:
				}

				if key == "" {
					logrus.Warnf("An empty string was found when getting keys from %v", from)
					continue
				}

				var value string
				value, err = rdb1.Get(ctx, key).Result()
				if err != nil {
					panic(err)
				}

				_, err = rdb2.Set(ctx, key, value, 0).Result()
				if err != nil {
					panic(err)
				}
			}

			if move {
				//if it's a move then we delete keys from the first one
				n, err := rdb1.Del(ctx, keys...).Result()
				if err != nil {
					panic(err)
				}
				if n != int64(len(keys)) {
					panic(fmt.Sprintf("failed to remove keys from '%v' during move", from))
				}
			}

			if cursor == 0 {
				break
			}
		}
	},
}

func init() {
	redisCmd.AddCommand(migrateCmd)

	migrateCmd.Flags().BoolVarP(&move, "move", "", false, "If set, data will be moved instead of copied.")
	migrateCmd.Flags().StringVarP(&match, "match", "", "*", "The glob filter to match key to be migrated.")

	migrateCmd.Flags().StringVarP(&from, "from", "f", "", "The address to query for data. (Ex. '0.0.0.0:6379')")
	migrateCmd.MarkFlagRequired("from")
	migrateCmd.Flags().StringVarP(&to, "to", "t", "", "The address to move or copy data to. (Ex. '0.0.0.0:6379')")
	migrateCmd.MarkFlagRequired("to")
}
