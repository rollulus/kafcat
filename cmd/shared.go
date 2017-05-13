package cmd

import (
	"encoding/hex"
	"fmt"
	"log"
	"time"

	yaml "gopkg.in/yaml.v1"

	"github.com/Shopify/sarama"
)

type ConsumerMessage struct {
	Key, Value string
	Topic      string
	Partition  int32
	Offset     int64
	Timestamp  time.Time // only set if kafka is version 0.10+
}

func getClient() (sarama.Client, error) {
	log.Printf("broker: %s\n", broker)
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V0_10_1_0
	return sarama.NewClient([]string{broker}, cfg)
}

func format(m *sarama.ConsumerMessage) error {
	fmt.Printf("---\n")
	cm := ConsumerMessage{hex.Dump(m.Key), hex.Dump(m.Value), m.Topic, m.Partition, m.Offset, m.Timestamp}
	bs, err := yaml.Marshal(cm)
	fmt.Printf("%s", string(bs))
	return err
}
