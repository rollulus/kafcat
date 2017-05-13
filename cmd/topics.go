package cmd

import (
	"github.com/rollulus/kafcat/pkg/kafcat"
	"github.com/spf13/cobra"
)

var offsets bool
var leader bool
var replicas bool
var isr bool
var writable bool

var all bool

var partitions bool

func init() {
	RootCmd.AddCommand(topicsCmd)
	topicsCmd.Flags().BoolVarP(&all, "all", "a", false, "show all possible info; equals -olwrsp")
	topicsCmd.Flags().BoolVarP(&offsets, "offsets", "o", false, "show oldest and newest offsets for each partition; implies --partitions")
	topicsCmd.Flags().BoolVarP(&leader, "leader", "l", false, "show leader information for each partition; implies --partitions")
	topicsCmd.Flags().BoolVarP(&writable, "writable", "w", false, "show writable state (valid leader accepting writes) for each partition; implies --partitions")
	topicsCmd.Flags().BoolVarP(&replicas, "replicas", "r", false, "show broker ids of replicas; implies --partitions")
	topicsCmd.Flags().BoolVarP(&isr, "isr", "s", false, "show broker ids of in-sync replicas; implies --partitions")
	topicsCmd.Flags().BoolVarP(&partitions, "partitions", "p", false, "show partitions")
}

var topicsCmd = &cobra.Command{
	Use: "topics [flags]",
	RunE: func(cmd *cobra.Command, args []string) error {
		if all {
			offsets = true
			leader = true
			replicas = true
			isr = true
			writable = true
			partitions = true
		}
		partitions = partitions || offsets || leader || replicas || isr || writable
		client, err := getClient()
		if err != nil {
			return err
		}

		infos, err := kafcat.GetTopicInfos(client, partitions, leader, offsets, writable, isr, replicas)
		if err != nil {
			return err
		}

		formatTopics(infos)

		return nil
	},
}
