package cmd

import (
	"encoding/json"
	"fmt"
	"kafka-producer/config"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	processCmd = &cobra.Command{
		Use: "process",
		Run: func(cmd *cobra.Command, args []string) {
			err := process(cmd, args)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		},
	}
)

// RecordValue represents the struct of the value in a Kafka message
type RecordValue struct {
	Count int
}

func process(cmd *cobra.Command, args []string) error {
	log.Tracef("process(...) called")

	// Create Kafka config map
	kafkaConfigMap := make(kafka.ConfigMap)
	for k, v := range config.ParsedConfig.GetKafkaConfigMap() {
		kafkaConfigMap.SetKey(k, v)
	}

	// Enable the log channel of Confluent Kafka Go library
	kafkaConfigMap.SetKey("go.logs.channel.enable", true)

	// Create producer instance
	p, err := kafka.NewProducer(&kafkaConfigMap)
	if err != nil {
		log.Errorf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	// Go-routine to handle log messages
	go func() {
		for {
			select {
			case logEntry, ok := <-p.Logs():
				if !ok {
					return
				}

				switch logEntry.Level {
				case 0:
					// KERN_EMERG
					log.Panicf("%s|%s|%s", logEntry.Tag, logEntry.Name, logEntry.Message)
				case 1:
					// KERN_ALERT
					log.Fatalf("%s|%s|%s", logEntry.Tag, logEntry.Name, logEntry.Message)
				case 2:
					// KERN_CRIT
					log.Fatalf("%s|%s|%s", logEntry.Tag, logEntry.Name, logEntry.Message)
				case 3:
					// KERN_ERR
					log.Errorf("%s|%s|%s", logEntry.Tag, logEntry.Name, logEntry.Message)
				case 4:
					// KERN_WARNING
					log.Warnf("%s|%s|%s", logEntry.Tag, logEntry.Name, logEntry.Message)
				case 5:
					// KERN_NOTICE
					log.Infof("%s|%s|%s", logEntry.Tag, logEntry.Name, logEntry.Message)
				case 6:
					// KERN_INFO
					log.Infof("%s|%s|%s", logEntry.Tag, logEntry.Name, logEntry.Message)
				case 7:
					// KERN_DEBUG
					log.Debugf("%s|%s|%s", logEntry.Tag, logEntry.Name, logEntry.Message)
				}
			}
		}
	}()

	// Go-routine to handle message delivery reports and
	// possibly other event types (errors, stats, etc)
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					log.Printf("Successfully produced record to topic %s partition [%d] @ offset %v\n",
						*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				}
			}
		}
	}()

	for n := 0; n < 10; n++ {
		recordKey := "alice"
		data := &RecordValue{
			Count: n,
		}
		recordValue, _ := json.Marshal(&data)

		log.Printf("Preparing to produce record: %s\t%s\n", recordKey, recordValue)

		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &config.ParsedConfig.Topic,
				Partition: kafka.PartitionAny,
			},
			Key:   []byte(recordKey),
			Value: []byte(recordValue),
		}, nil)
	}

	// Wait for all messages to be delivered
	p.Flush(15 * 1000)

	log.Printf("10 messages were produced to topic %s!\n", config.ParsedConfig.Topic)

	p.Close()

	return nil
}
