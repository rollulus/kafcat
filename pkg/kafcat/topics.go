package kafcat

import (
	"github.com/Shopify/sarama"
)

type TopicInfo struct {
	Name         string
	Partitions   map[int32]PartitionInfo `yaml:",omitempty"`
	PartitionIds []int32                 `yaml:",omitempty,flow"`
}

type LeaderInfo struct {
	ID      int32
	Address string
}

type OffsetInfo struct {
	Oldest int64
	Newest int64
}

type PartitionInfo struct {
	Leader         *LeaderInfo `yaml:",omitempty"`
	Offsets        *OffsetInfo `yaml:",omitempty"`
	Writable       *bool       `yaml:",omitempty"`
	Replicas       []int32     `yaml:",omitempty,flow"`
	InSyncReplicas []int32     `yaml:",omitempty,flow"`
}

func GetTopicInfos(client sarama.Client, partitions, leader, offsets, writable, isr, replicas bool) ([]TopicInfo, error) {
	infos := []TopicInfo{}

	ts, err := client.Topics()
	if err != nil {
		return nil, err
	}

	for _, topic := range ts {
		infos = append(infos, TopicInfo{Name: topic})

		if !partitions {
			continue
		}

		info := &infos[len(infos)-1]

		info.Partitions = map[int32]PartitionInfo{}

		ps, err := client.Partitions(topic)
		if err != nil {
			return nil, err
		}

		if !(leader || offsets || writable || isr || replicas) {
			info.PartitionIds = ps
			continue
		}

		// construct writable-partition map for quick lookups
		wps, err := client.WritablePartitions(topic)
		if err != nil {
			return nil, err
		}
		wpsm := map[int32]bool{}
		for _, wp := range wps {
			wpsm[wp] = true
		}

		for _, p := range ps {
			pi := PartitionInfo{}

			if leader {
				leadBroker, err := client.Leader(topic, p)
				if err != nil {
					return nil, err
				}
				pi.Leader = &LeaderInfo{ID: leadBroker.ID(), Address: leadBroker.Addr()}
			}

			if offsets {
				pi.Offsets = &OffsetInfo{}
				if pi.Offsets.Oldest, err = client.GetOffset(topic, p, sarama.OffsetOldest); err != nil {
					return nil, err
				}
				if pi.Offsets.Newest, err = client.GetOffset(topic, p, sarama.OffsetNewest); err != nil {
					return nil, err
				}
			}

			if writable {
				True := true
				False := false
				if wpsm[p] {
					pi.Writable = &True
				} else {
					pi.Writable = &False
				}
			}

			if isr {
				if pi.InSyncReplicas, err = client.InSyncReplicas(topic, p); err != nil {
					return nil, err
				}
			}

			if replicas {
				if pi.Replicas, err = client.Replicas(topic, p); err != nil {
					return nil, err
				}

			}

			info.Partitions[p] = pi
		}
	}

	return infos, nil
}
