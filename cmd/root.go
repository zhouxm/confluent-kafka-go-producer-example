package cmd

import (
	"kafka-producer/config"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use: "kafka-producer",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			os.Exit(1)
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			log.Tracef("PersistentPreRunE(...) called")
			log.Info(args)

			err := config.ParsedConfig.Parse()
			if err != nil {
				log.Errorf("Cannot parse config: %s", err.Error())
				os.Exit(1)
			}

			logLevel, err := log.ParseLevel(config.ParsedConfig.LogLevel)
			if err != nil {
				log.Errorf(err.Error())
				os.Exit(1)
			}
			log.SetLevel(logLevel)
			log.Infof("Log level set to \"%s\"", logLevel)

			return nil
		},
	}
)

func init() {
	log.Tracef("init() called")

	rootCmd.PersistentFlags().StringP("config", "c", "", "config file")
	rootCmd.AddCommand(processCmd)
	rootCmd.AddCommand(versionCmd)

	viper.BindPFlags(rootCmd.PersistentFlags())
}

// Execute ...
func Execute() error {
	return rootCmd.Execute()
}
