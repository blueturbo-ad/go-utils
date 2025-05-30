package kafkaclient

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/blueturbo-ad/go-utils/environment"
	k8sclient "github.com/blueturbo-ad/go-utils/k8s_tool/k8s_client"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaMsg struct {
	Msg string `json:"msg"`
}

func TestKafkaClient(t *testing.T) {
	os.Setenv("POD_NAMESPACE", "dsp-ns")
	environment.Init()
	k8sclient.GetSingleton().SetUp()
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	p := dir + "/config/kafka_config.yaml"
	// GetSingleton().UpdateLoadK8sConfigMap("kafka-config", "Pro")
	GetSingleton().UpdateFromFile(p, "Pro")
	t.Run("kafka client producer", func(t *testing.T) {
		p, err := GetSingleton().GetProducerClient("win_kafka")
		if err != nil {
			t.Errorf("GetProducerClient() = %v; want nil", err)
		}
		kmsg := &KafkaMsg{
			Msg: "{'msg_type': 6, 'strategy_id': 554, 'strategy_data': 'CKoEEL0CGgtrYWZrYea1i+ivlSCAiXoogIl6Ohz///8H////B////wf///8H////B////wf///8HQgNBTELoAQHwAcCEPQ=='}",
		}

		msgbyte, err := json.Marshal(kmsg)
		if err != nil {
			fmt.Printf("failed to marshal message: %s\n", err)
			return
		}
		topic := "topic-sg-mis-bid_server-ad_library"
		err = p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          msgbyte,
		}, nil)
		if err != nil {
			fmt.Printf("failed to produce message: %s\n", err)
			return
		}

		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					msg := fmt.Sprintf("failed to deliver message: %s\n", ev.TopicPartition.Error)
					fmt.Println(msg)
					return
				} else {
					fmt.Println("send success")
					// msg := fmt.Sprintf("message delivered to %s [%d] at offset %v\n",
					// 	*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
					// fmt.Println(msg)
					return
				}
			case kafka.Error:
				fmt.Printf("kafka error: %v\n", ev)
				msg := fmt.Sprintf("kafka error: %s\n", ev)
				fmt.Println(msg)
				return
			default:
				fmt.Printf("ignored event: %v\n", ev)
				msg := fmt.Sprintf("ignored event: %s\n", ev)
				fmt.Println(msg)
				return
			}
		}

		time.Sleep(1 * time.Second)
	})
	t.Run("kafka client consumer", func(t *testing.T) {
		c, err := GetSingleton().GetConsumerClient("win_kafka", "test")
		if err != nil {
			t.Errorf("GetConsumerClient() = %v; want nil", err)
		}
		topic := "topic-sg-mis-event_server-budget_pacing"
		err = c.Subscribe(topic, nil)
		if err != nil {
			fmt.Printf("failed to subscribe topic: %s\n", err)
			return
		}
		func() {
			for {
				msg, err := c.ReadMessage(-1)
				if err == nil {
					fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
					// c.CommitMessage(msg)

				} else {
					fmt.Printf("Consumer error: %v (%v)\n", err, msg)
				}
			}
		}()
	})
}
