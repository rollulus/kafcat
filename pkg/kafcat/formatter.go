package kafcat

import (
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/Shopify/sarama"
	yaml "gopkg.in/yaml.v2"
)

type Formatter struct {
	Out      io.Writer
	ForceHex bool
}

type ConsumerMessage struct {
	Key, Value string
	Topic      string
	Partition  int32
	Offset     int64
	Timestamp  time.Time // only set if kafka is version 0.10+
}

func (f *Formatter) bytesToString(bs []byte) string {
	if f.ForceHex {
		return hex.Dump(bs)
	}
	// TODO: if heuristic ...
	return string(bs)
}

func (f *Formatter) FormatConsumerMessage(m *sarama.ConsumerMessage) error {
	fmt.Fprintln(f.Out, "---")
	cm := &ConsumerMessage{f.bytesToString(m.Key), f.bytesToString(m.Value), m.Topic, m.Partition, m.Offset, m.Timestamp}
	return yaml.NewEncoder(f.Out).Encode(cm)
}
