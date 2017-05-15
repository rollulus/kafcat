package consumer

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

type Consumer struct {
	Client          sarama.Client
	Topic           string
	Partitions      map[int32]bool
	TMax            *time.Time
	BeyondHighWater bool
	StartOffset     map[int32]int64 // per partition
}

type ConsumerOption func(*Consumer) error

func SpecificPartitions(ps map[int32]bool) ConsumerOption {
	return func(c *Consumer) error {
		log.Printf("consopt: consume specific partitions %#v", ps)
		c.Partitions = ps
		return nil
	}
}

func AllPartitions() ConsumerOption {
	return func(c *Consumer) error {
		log.Printf("consopt: consume all partitions")
		parts, err := c.Client.Partitions(c.Topic)
		if err != nil {
			return err
		}
		log.Printf("consopt: consume all partitions %#v", parts)
		c.Partitions = map[int32]bool{}
		for _, p := range parts {
			c.Partitions[p] = true
		}
		return nil
	}
}

func UntilTime(t time.Time) ConsumerOption {
	return func(c *Consumer) error {
		log.Printf("consopt: consume tmax %s", t.Format(time.RFC3339Nano))
		c.TMax = &t
		return nil
	}
}

func BeyondHighWaterMark(b bool) ConsumerOption {
	return func(c *Consumer) error {
		log.Printf("consopt: beyond hwm %v", b)
		c.BeyondHighWater = b
		return nil
	}
}

func FromTime(t time.Time) ConsumerOption {
	return func(c *Consumer) error {
		tMs := t.UnixNano() / 1000000
		log.Printf("consopt: consume tmin %s (%d ms)", t.Format(time.RFC3339Nano), tMs)
		return c.bootstrapStartOffset(func(p int32) (int64, error) {
			log.Printf("obtain offset for %s:%d", c.Topic, p)
			return c.Client.GetOffset(c.Topic, p, tMs)
		})
	}
}

func FromCurrentOffset() ConsumerOption {
	return func(c *Consumer) error {
		log.Printf("consopt: offset newest")
		return c.bootstrapStartOffset(func(p int32) (int64, error) {
			return sarama.OffsetNewest, nil
		})
	}
}

func FromOldestOffset() ConsumerOption {
	return func(c *Consumer) error {
		log.Printf("consopt: offset oldest")
		return c.bootstrapStartOffset(func(p int32) (int64, error) {
			return sarama.OffsetOldest, nil
		})
	}
}

func Nop() ConsumerOption {
	return func(c *Consumer) error {
		log.Printf("consopt: nop")
		return nil
	}
}

func New(client sarama.Client, topic string, opts ...ConsumerOption) (*Consumer, error) {
	c := &Consumer{Client: client, Topic: topic}

	for _, o := range opts {
		if err := o(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// for all existing partitions
func (c *Consumer) bootstrapStartOffset(f func(p int32) (int64, error)) error {
	if c.StartOffset != nil {
		return fmt.Errorf("StartOffset already set")
	}
	c.StartOffset = map[int32]int64{}
	parts, err := c.Client.Partitions(c.Topic)
	if err != nil {
		return err
	}
	for _, p := range parts {
		o, err := f(p)
		if err != nil {
			return err
		}
		c.StartOffset[p] = o
	}
	return nil
}

func (c *Consumer) Messages(ctx context.Context) (chan *sarama.ConsumerMessage, chan error, error) {
	if c.Partitions == nil {
		return nil, nil, fmt.Errorf("no partitions defined")
	}

	cons, err := sarama.NewConsumerFromClient(c.Client)
	if err != nil {
		return nil, nil, err
	}

	fanIn := make(chan *sarama.ConsumerMessage)
	errors := make(chan error)

	var wg sync.WaitGroup

	for p := range c.Partitions {
		off := c.StartOffset[p]
		wg.Add(1)
		log.Printf("consume %d %d", p, off)
		go func(p int32, off int64) {
			pc, err := cons.ConsumePartition(c.Topic, p, off)
			if err != nil {
				wg.Done()
				errors <- err
				return
			}
			// waits for context and async closes PartitionConsumer
			go func(pc sarama.PartitionConsumer) {
				<-ctx.Done()
				log.Printf("async close gorout %d %d", p, off)
				pc.AsyncClose()
			}(pc)
			// fans PartitionConsumer's messages into fanIn
			go func(pc sarama.PartitionConsumer) {
				for m := range pc.Messages() {
					if c.TMax != nil && m.Timestamp.After(*c.TMax) {
						log.Printf("timestamp %s > tmax %s; quit", m.Timestamp, c.TMax)
						break
					}
					fanIn <- m
					if !c.BeyondHighWater && m.Offset+1 >= pc.HighWaterMarkOffset() {
						log.Printf("o %d >= hwmo", m.Offset)
						break
					}
				}
				log.Printf("exit infanner gorout %d %d", p, off)
				wg.Done()
			}(pc)
		}(p, off)
	}

	go func() {
		wg.Wait()
		log.Printf("wg waited")
		close(fanIn)
	}()
	return fanIn, errors, nil
}
