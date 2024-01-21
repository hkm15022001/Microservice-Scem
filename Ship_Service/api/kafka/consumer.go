package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/Shopify/sarama"
	CommonService "github.com/hkm15022001/Supply-Chain-Event-Management/internal/service/common"
	CommonMessage "github.com/hkm15022001/Supply-Chain-Event-Management/internal/service/common_message"
)

var consumer sarama.Consumer
var stopConsumer = make(chan struct{})

func ConnectConsumer(brokerList []string) error {
	config := sarama.NewConfig()
	var err error

	// Tạo consumer
	consumer, err = sarama.NewConsumer(brokerList, config)
	if err != nil {
		return err
	}

	return nil
}

func CloseConsumer() {

	// Đóng consumer
	if err := consumer.Close(); err != nil {
		log.Fatalf("Error closing the consumer: %s", err.Error())
	}
}

func ConsumeMessages(topic string, wg *sync.WaitGroup) {
	type MessageData map[string][]int
	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log.Fatalf("Error getting partitions: %s", err.Error())
		return
	}

	for _, partition := range partitions {
		partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Fatalf("Error creating partition consumer: %s", err.Error())
			return
		}

		go func(pc sarama.PartitionConsumer) {
			defer wg.Done()
		ConsumerLoop:
			for {
				select {
				case msg := <-pc.Messages():
					// Xử lý message ở đây
					fmt.Printf("Received message: %s\n", msg.Value)
					var data MessageData
					err := json.Unmarshal([]byte(msg.Value), &data)
					if err != nil {
						fmt.Println("Lỗi chuyển đổi JSON:", err)
						return
					}

					// In ra các giá trị từ struct
					for key, values := range data {
						fmt.Printf("Key: %s, Values: %v\n", key, values)
						longShipId, err := CommonService.CreateLongShipHandlerGRPC()
						if err != nil {
							log.Fatal("Error when create longship from plan service: ", err)
						}
						for _, orderId := range values {
							CommonService.UpdateOrderLongShipGRPC(uint(orderId), longShipId)
							err = CommonMessage.PublishPaymentConfirmedMessage(uint(orderId))
							if err != nil {
								log.Fatal("Error when update payment confirm from consumer: ", err)
							}
						}
					}

				case err := <-pc.Errors():
					// Xử lý lỗi nếu có
					fmt.Printf("Error: %v\n", err)

				case <-stopConsumer:
					// Dừng goroutine khi được báo hiệu
					break ConsumerLoop
				}
			}
		}(partitionConsumer)
	}

	// Chờ tất cả các goroutine hoàn thành
	wg.Wait()
}

func StartConsumer(brokerList []string, topic string) {
	err := ConnectConsumer(brokerList)
	if err != nil {
		log.Fatalf("Error creating the consumer: %s", err.Error())
	}
	defer CloseConsumer()

	var wg sync.WaitGroup

	// Tạo goroutine cho consumer
	wg.Add(1)
	go ConsumeMessages(topic, &wg)

	log.Println("Connected to Kafka-consumer")

	// Chờ tín hiệu tắt
	signals := make(chan os.Signal, 1)

	select {
	case <-signals:
		fmt.Println("Exiting consumer...")
		close(stopConsumer) // Signal goroutines to stop
	case <-stopConsumer:
		fmt.Println("Exiting consumer...")
	}
	wg.Wait()
}
