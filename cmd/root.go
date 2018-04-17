package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var broker string
var rootCA string
var certPEM string
var keyPEM string
var verbose bool
var saramaLog bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "kafcat",
	Short: "The pretended Swiss army knife for Apache Kafka",
	Long:  `The pretended Swiss army knife for Apache Kafka`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !verbose {
			log.SetOutput(ioutil.Discard)
		}
	},
	SilenceUsage: true,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&saramaLog, "log-client", "w", false, "enable sarama's (underlying kafka client) log to stderr")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "be verbose")
	RootCmd.PersistentFlags().StringVarP(&broker, "broker-list", "b", "localhost:9092", "brokers")
	RootCmd.PersistentFlags().StringVar(&rootCA, "root-ca", "", "filename of the root certificate in PEM format")
	RootCmd.PersistentFlags().StringVar(&certPEM, "client-cert", "", "filename of the client certificate in PEM format")
	RootCmd.PersistentFlags().StringVar(&keyPEM, "client-key", "", "filename of the client's private key in PEM format")
}
