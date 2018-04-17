package cmd

import (
	"io"
	"log"
	"os"

	"github.com/Shopify/sarama"

	"github.com/rollulus/kafcat/pkg/chisel"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

func init() {
	RootCmd.AddCommand(produceCmd)
}

var produceCmd = &cobra.Command{
	Use:   "produce [<topic>]",
	Short: "Produce messages to Kafka",
	Long: `Produce messages to Kafka from stdin. The input format is identical to the 
output of kafcat's "cat" command. The key, value and topic fields are used;
partition, offset and timestamp are ignored. Topic can be overridden by
specifying it.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var overrideTopic string
		if len(args) == 1 {
			overrideTopic = args[0]
		}
		cfg, err := getConfig()
		if err != nil {
			return err
		}
		cfg.Producer.Return.Successes = true

		client, err := sarama.NewClient([]string{broker}, cfg)
		if err != nil {
			return err
		}
		defer client.Close()
		prod, err := sarama.NewSyncProducerFromClient(client)
		if err != nil {
			return err
		}
		defer prod.Close()

		dec := yaml.NewDecoder(os.Stdin)
		for {
			var x ConsumerMessage
			err := dec.Decode(&x)
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			// TODO: validate input
			keyBs, err := chisel.ParseHexDump(x.Key)
			if err != nil {
				return err
			}
			valueBs, err := chisel.ParseHexDump(x.Value)
			if err != nil {
				return err
			}
			topic := x.Topic
			if overrideTopic != "" {
				topic = overrideTopic
			}
			part, offset, err := prod.SendMessage(&sarama.ProducerMessage{
				Key:   sarama.ByteEncoder(keyBs),
				Value: sarama.ByteEncoder(valueBs),
				Topic: topic,
			})
			if err != nil {
				return err
			}
			log.Printf("produced message; part=%d offset=%d", part, offset)
		}
		log.Printf("fini")
		return nil
	},
}
