package kafka

import (
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
)

var producer sarama.AsyncProducer

func Connect(brokerList []string) error {
	config := sarama.NewConfig()
	var err error
	producer, err = sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		return err
	}
	return nil
}

func CloseProducer() {
	if err := producer.Close(); err != nil {
		log.Fatalf("Error closing the producer: %s", err.Error())
	}
}

func ProduceMessage(topic string, message string) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	producer.Input() <- msg
	fmt.Println("Message sent")
}

func StartProducer(brokerList []string, topic string) {
	err := Connect(brokerList)
	if err != nil {
		log.Fatalf("Error creating the producer: %s", err.Error())
	}
	defer CloseProducer()

	log.Println("Connected to Kafka")
	go func() {
		for {
			ProduceMessage(topic, "Hello from Golang!")
			time.Sleep(10 * time.Second)
		}
	}()
	select {}
}

// Usage:
// func main() {
//     brokerList := []string{"localhost:9092"}
//     topic := "order-topic"
//     StartProducer(brokerList, topic)
// }
