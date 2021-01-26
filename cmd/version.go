package cmd

import (
	"fmt"
	version "kafka-producer/version"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	log "github.com/sirupsen/logrus"
	cobra "github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			err := showVersion(cmd, args)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		},
	}
)

func showVersion(cmd *cobra.Command, args []string) error {
	log.Tracef("version(...) called")

	fmt.Printf("Version:            %s\n", version.Version)
	fmt.Printf("Built:              %s\n", version.BuildDate)
	fmt.Printf("LibrdkafkaLinkInfo: %s\n", kafka.LibrdkafkaLinkInfo)

	return nil
}
