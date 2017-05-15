package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"time"

	"strings"

	"github.com/rollulus/kafcat/pkg/consumer"
	"github.com/spf13/cobra"
)

var follow bool
var since string

func init() {
	dumpCmd.Flags().BoolVarP(&follow, "follow", "f", false, "don't quit on end of log; keep following the topic")
	dumpCmd.Flags().StringVarP(&since, "since", "s", "", "only return logs newer than a relative duration like 5s, 2m, or 3h, or shorthands 0 and now")
	RootCmd.AddCommand(dumpCmd)
}

// topic[:partitions]
func parseTopicPartition(tp string) (string, consumer.ConsumerOption, error) {
	ss := strings.Split(tp, ":")
	if len(ss) == 1 {
		return ss[0], consumer.AllPartitions(), nil
	}
	if len(ss) != 2 {
		return "", nil, fmt.Errorf("topic `%s` not understood", tp)
	}
	ps := map[int32]bool{}
	for _, s := range strings.Split(ss[1], ",") {
		i, err := strconv.ParseInt(s, 0, 32)
		if err != nil {
			return "", nil, err
		}
		ps[int32(i)] = true
	}
	return ss[0], consumer.SpecificPartitions(ps), nil
}

func parseSince() (consumer.ConsumerOption, error) {
	switch since {
	case "":
		return consumer.FromOldestOffset(), nil
	case "0", "now":
		return consumer.FromCurrentOffset(), nil
	default:
		du, err := time.ParseDuration(since)
		if err != nil {
			return nil, err
		}
		return consumer.FromTime(time.Now().Add(-du)), nil
	}
}

func parseFollow() (consumer.ConsumerOption, error) {
	return consumer.BeyondHighWaterMark(follow), nil
}

var dumpCmd = &cobra.Command{
	Use: "cat <topic>[:partition[,partition]*] [flags]",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("expect a single argument: topic")
		}
		topicPartition := args[0]

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

		topic, partitionsOpts, err := parseTopicPartition(topicPartition)
		if err != nil {
			return err
		}

		con, err := consumer.New(client, topic, partitionsOpts, fOpt, sOpt)

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
				return err
			}
		}
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
