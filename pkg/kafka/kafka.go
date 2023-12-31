package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Zavr22/EMTestTask/pkg/models"
	"github.com/Zavr22/EMTestTask/pkg/rest/service"
	"github.com/segmentio/kafka-go"
	"log"
)

type Kafka struct {
	userServ   *service.UserService
	writerConn *kafka.Conn
	readerConn *kafka.Conn
}

func NewKafkaConsumer(userServ *service.UserService, conn *kafka.Conn) *Kafka {
	return &Kafka{userServ: userServ, writerConn: conn, readerConn: conn}
}

func (k *Kafka) ListenToKafkaTopic() {
	batch := k.readerConn.ReadBatch(0, 1e6)
	for {
		message, err := batch.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from Kafka: %v", err)
			return
		}

		var fio models.FIO
		if err := json.Unmarshal(message.Value, &fio); err != nil {
			log.Printf("Error unmarshalling JSON message: %v", err)
			if err := k.sendToFailedTopic(message.Value); err != nil {
				log.Printf("Error sending message to FIO_FAILED: %v", err)
			}
			continue
		}

		log.Println("Received message from Kafka:", string(message.Value))

		if userID, err := k.userServ.EnrichAndSaveToDB(context.Background(), fio.Name, fio.Surname, fio.Patronymic); err != nil {
			log.Printf("Error processing message: %v", err)
			if err := k.sendToFailedTopic(message.Value); err != nil {
				log.Printf("Error sending message to FIO_FAILED: %v", err)
			}
			log.Printf("userID: %v", userID)
			continue
		}

		log.Println("Message processed successfully")
	}
}

func (k *Kafka) sendToFailedTopic(message []byte) error {

	kafkaMessage := kafka.Message{
		Key:   nil,
		Value: message,
	}

	fmt.Println("Sending message to FIO_FAILED topic:", string(message))

	if _, err := k.writerConn.WriteMessages(kafkaMessage); err != nil {
		return err
	} else {
		log.Println("Message sent to FIO_FAILED topic successfully")
	}

	return nil
}

func (k *Kafka) ProduceMessage() {
	fioMessages := &models.FIO{
		Name:       "Dmitriy",
		Surname:    "Ushakov",
		Patronymic: "Vasilevich",
	}
	messageBytes, err := json.Marshal(fioMessages)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
	}
	kafkaMessage := kafka.Message{
		Key:   nil,
		Value: messageBytes,
	}
	if _, err := k.writerConn.WriteMessages(kafkaMessage); err != nil {
		log.Printf("Error sending JSON message to Kafka: %v", err)
	} else {
		log.Println("Message sent to Kafka successfully")
	}
}
