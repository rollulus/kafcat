package cmd

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"time"

	"github.com/spf13/cobra"
)

var follow bool
var since string

func init() {
	dumpCmd.Flags().BoolVarP(&follow, "follow", "f", false, "don't quit on end of log; keep following the topic")
	dumpCmd.Flags().StringVarP(&since, "since", "s", "", "only return logs newer than a relative duration like 5s, 2m, or 3h, or shorthands 0 and now")
	RootCmd.AddCommand(dumpCmd)
}

func parseSince() (ConsumerOption, error) {
	switch since {
	case "":
		return FromOldestOffset(), nil
	case "0", "now":
		return FromCurrentOffset(), nil
	default:
		du, err := time.ParseDuration(since)
		if err != nil {
			return nil, err
		}
		return FromTime(time.Now().Add(-du)), nil
	}
}

func parseFollow() (ConsumerOption, error) {
	return BeyondHighWaterMark(follow), nil
}

var dumpCmd = &cobra.Command{
	Use: "cat <topic>[:partition] [flags]",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("expect a single argument: topic")
		}
		topic := args[0]

		client, err := getClient()
		if err != nil {
			return err
		}

		fOpt, err := parseFollow()
		if err != nil {
			return err
		}

		sOpt, err := parseSince()
		if err != nil {
			return err
		}

		con, err := NewConsumer(client, topic, AllPartitions(), fOpt, sOpt)

		if err != nil {
			return err
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancelOnSigterm(cancel)

		msgs, errs, err := con.Messages(ctx)
		if err != nil {
			cancel()
			return err
		}

		for {
			select {
			case m := <-msgs:
				if m == nil {
					log.Printf("got nil")
					cancel()
					return nil
				}
				format(m)

			case err := <-errs:
				log.Printf("%s", err)
				cancel()
				return nil
			}
		}

		return nil

		// return dump(c, "tweets")
	},
}

func cancelOnSigterm(cancel context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Printf("received signal %s; cancel", sig)
		cancel()
	}()
}
