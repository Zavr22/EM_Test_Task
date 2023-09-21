package main

import (
	"EMTestTask/pkg/model"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
)

func main() {
	topic := "FIO"
	brokers := []string{"localhost:9092"}

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})

	fioMessages := []*model.FIO{
		{Name: "Dmitriy", Surname: "Ushakov", Patronymic: "Vasilevich"},
		{Name: "Mihail", Surname: "Alexandrovich", Patronymic: "Kulik"},
	}

	for _, fio := range fioMessages {
		messageBytes, err := json.Marshal(fio)
		if err != nil {
			log.Printf("Error marshalling JSON: %v", err)
			continue
		}

		kafkaMessage := kafka.Message{
			Key:   nil,
			Value: messageBytes,
		}
		if err := w.WriteMessages(context.Background(), kafkaMessage); err != nil {
			log.Printf("Error sending JSON message to Kafka: %v", err)
		}
	}

	w.Close()
}
