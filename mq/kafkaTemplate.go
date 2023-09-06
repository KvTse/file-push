package mq

import (
	"context"
	ftpClient "file-push/ftp"
	"file-push/global"
	"file-push/limit"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)
import "github.com/segmentio/kafka-go"

func SendMessage2Kafka(ctx context.Context, message string) {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(global.GVA_CONFIG.KafkaConfig.Brokers),
		Topic:                  global.GVA_CONFIG.KafkaConfig.Topic,
		Balancer:               &kafka.Hash{},
		WriteTimeout:           3 * time.Second,
		RequiredAcks:           kafka.RequireAll,
		AllowAutoTopicCreation: true,
	}
	defer func() {
		if err := writer.Close(); err != nil {
			log.Fatal("failed to close writer:", err)
		}
	}()
	const retries = 3
	for i := 0; i < retries; i++ {
		if err := writer.WriteMessages(ctx, kafka.Message{Value: []byte(message)}); err != nil {
			log.Printf("send message %s error %v", message, err)
			if err == kafka.LeaderNotAvailable {
				log.Printf("kafka leader is not available ...")
				time.Sleep(time.Second)
			} else {
				log.Printf("kafka write error %v", err)
			}
		} else {
			log.Printf("send message %s to topic %s", message, writer.Topic)
			break
		}
	}

}

var reader *kafka.Reader

func ConsumeMessageFromKafka(gLimit *limit.GLimit) {
	reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{global.GVA_CONFIG.KafkaConfig.Brokers},
		GroupID:        global.GVA_CONFIG.KafkaConfig.ConsumerGroupId,
		Topic:          global.GVA_CONFIG.KafkaConfig.Topic,
		MaxBytes:       10e6,
		StartOffset:    kafka.LastOffset,
		CommitInterval: 0,
	})
	log.Println("start a kafka consumer ......")
	for {
		goFunc := func() {
			ctx := context.Background()
			message, err := reader.FetchMessage(ctx)
			if err != nil {
				log.Fatalf("read message from kafka error %v", err)
			}
			log.Printf("read message from topic %s, value is %s", message.Topic, string(message.Value))
			business := ftpClient.DoConsumeFtpPushBusiness(message.Value)
			if business {
				log.Printf("====================commit message...")
				if err := reader.CommitMessages(ctx, message); err != nil {
					log.Printf("commit message error %v", err)
				}
			}
		}
		gLimit.Run(goFunc)
	}
}
func ListenerSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	sig := <-c
	log.Printf("receive signal %s", sig)
	defer func() {
		if reader != nil {
			if err := reader.Close(); err != nil {
				log.Fatalf("kafka reader close occur error %v", err)
			}
		}
	}()
	os.Exit(0)
}
