package kafka

import (
	"context"
	"encoding/json"
	"github.com/Zavr22/EMTestTask/pkg/model"
	"github.com/Zavr22/EMTestTask/web/rest/service"
	"github.com/segmentio/kafka-go"
	"log"
)

type KafkaConsumer struct {
	userServ *service.UserService
}

func NewKafkaConsumer(userServ *service.UserService) *KafkaConsumer {
	return &KafkaConsumer{userServ: userServ}
}

func (k *KafkaConsumer) ListenToKafkaTopic() {
	topic := "FIO"
	brokers := []string{"kafka:9092"}
	groupID := "my-group"

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	for {
		message, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message from Kafka: %v", err)
			continue
		}

		var fio model.FIO
		if err := json.Unmarshal(message.Value, &fio); err != nil {
			log.Printf("Error unmarshalling JSON message: %v", err)
			if err := k.sendToFailedTopic(message.Value); err != nil {
				log.Printf("Error sending message to FIO_FAILED: %v", err)
			}
			continue
		}

		if userID, err := k.userServ.EnrichAndSaveToDB(context.Background(), fio.Name, fio.Surname, fio.Patronymic); err != nil {
			log.Printf("Error processing message: %v", err)
			if err := k.sendToFailedTopic(message.Value); err != nil {
				log.Printf("Error sending message to FIO_FAILED: %v", err)
			}
			log.Printf("userID: %v", userID)
			continue
		}
	}
}

func (k *KafkaConsumer) sendToFailedTopic(message []byte) error {
	topic := "FIO_FAILED"
	brokers := []string{"kafka:9092"}

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})

	kafkaMessage := kafka.Message{
		Key:   nil,
		Value: message,
	}

	if err := w.WriteMessages(context.Background(), kafkaMessage); err != nil {
		return err
	}

	w.Close()
	return nil
}
