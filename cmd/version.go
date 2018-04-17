package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	GitCommit string
	GitTag    string
	SemVer    string
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays version information",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("kafcat v%s; tag: %s; commit: %s\n", SemVer, GitTag, GitCommit)
		return nil
	},
}
