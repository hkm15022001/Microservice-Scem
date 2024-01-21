package kafka

import (
	"encoding/json"
	"fmt"
	"log"

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

func ProduceMessage(topic string, orderInfo interface{}) {
	orderInfoJson, err := json.Marshal(orderInfo)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(orderInfoJson),
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

	log.Println("Connected to Kafka-producer")
	select {}
}

// Usage:
// func main() {
//     brokerList := []string{"localhost:9092"}
//     topic := "order-topic"
//     StartProducer(brokerList, topic)
// }
