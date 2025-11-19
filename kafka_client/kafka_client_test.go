package kafkaclient

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/blueturbo-ad/go-utils/environment"
	"github.com/blueturbo-ad/go-utils/feishu"
	k8sclient "github.com/blueturbo-ad/go-utils/k8s_tool/k8s_client"
	"github.com/blueturbo-ad/go-utils/zap_loggerex"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
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
		topic := "topic-us-mis-bid_server-ad_library"
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
		for {
			topic := "topic-us-mis-bid_server-ad_library"
			err = c.Subscribe(topic, nil)
			if err != nil {
				fmt.Printf("failed to subscribe topic: %s\n", err)
				return
			}

			metadata, err := c.GetMetadata(&topic, false, int(time.Duration(20*time.Millisecond)))
			if err != nil {
				fmt.Printf("failed to get metadata: %s\n", err.Error())
				zap_loggerex.GetSingleton().Error("bid_stdout_logger", "failed to get metadata, %+v", err)

			}
			topicMetadata, ok := metadata.Topics[topic]
			if !ok {
				fmt.Printf("failed to get topic metadata for topic: %s\n", topic)
				zap_loggerex.GetSingleton().Error("bid_stdout_logger", "failed to get topic metadata")
				return
			}
			queryPartitions := []kafka.TopicPartition{}

			for _, partition := range topicMetadata.Partitions {
				queryPartitions = append(queryPartitions, kafka.TopicPartition{
					Topic:     &topic,
					Partition: partition.ID,
					Offset:    kafka.Offset(1761150160107712058),
				})
			}
			resultPartitions, err := c.OffsetsForTimes(queryPartitions, int(time.Duration(20*time.Millisecond)))
			if err != nil {
				fmt.Printf("failed to get offsets for times: %s\n", err.Error())
				zap_loggerex.GetSingleton().Error("bid_stdout_logger", "failed to get offsets for times, %+v", err)
			}
			assignPartitions := []kafka.TopicPartition{}

			for _, partition := range resultPartitions {
				offset := kafka.OffsetEnd
				if partition.Error != nil {
					offset = kafka.OffsetEnd
				} else {
					if partition.Offset < 0 && partition.Offset != kafka.OffsetEnd {
						zap_loggerex.GetSingleton().Info("bid_stdout_logger", "invalid offset for partition %d, set to the newest offset", partition.Partition)
						offset = kafka.OffsetEnd
					} else {
						offset = partition.Offset
					}
				}
				fmt.Printf("partition offset info, %+v\n", partition)
				assignPartitions = append(assignPartitions, kafka.TopicPartition{
					Topic:     &topic,
					Partition: partition.Partition,
					Offset:    offset,
				})
			}
			if err := c.Assign(assignPartitions); err != nil {
				msg := fmt.Sprintf("failed to assign partitions, %s", err.Error())
				zap_loggerex.GetSingleton().Error("bid_stdout_logger", "%s", msg)
				if err := feishu.GetInstance().Send("error", string(msg)); err != nil {
					zap_loggerex.GetSingleton().Error("bid_stdout_logger", "failed to send feishu alert, %+v", err)
					fmt.Println(err.Error())
				}

			}
			// func() {
			// 	for {
			// 		msg, err := c.ReadMessage(-1)
			// 		if err == nil {
			// 			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			// 			// c.CommitMessage(msg)

			// 		} else {
			// 			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			// 		}
			// 	}
			// }()
		}
	})
}
