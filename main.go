package main

import (
	"fmt"

	"github.com/blueturbo-ad/dsp_data_interface/gen/go/dmp_pb"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"google.golang.org/protobuf/proto"
)

// func main() {
// 	os.Setenv("POD_NAME", "test")
// 	os.Setenv("POD_NAMESPACE", "dsp-ns")
// 	environment.Init()
// 	workPath, err := os.Getwd()
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	cs := workPath + "/config/cloud_storage.yaml"
// 	err = gcpcloudstorage.GetSingleton().UpdateFromFile(cs, "Dev")
// 	if err != nil {
// 		fmt.Println("gcp cloud storage", err.Error())
// 	}
// 	ac := workPath + "/config/gcp_acc_token_conf.yaml"
// 	err = gcpcloudtool.GetSingleton().UpdateFromFile(ac, "Dev")
// 	if err != nil {
// 		fmt.Println("gcp access token", err.Error())
// 	}
// 	var e = redisclient.GetSingleton()

// 	p := workPath + "/config/redis_conf.yaml"
// 	err = e.UpdateFromFile(p, "Dev")
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	client := e.GetReadClient("event_redis")
// 	if client == nil {
// 		fmt.Println("redis get client", err.Error())
// 	}
// 	ctx := context.Background()
// 	pong, err := client.Ping(ctx).Result()
// 	if err != nil {
// 		log.Fatalf("Error connecting to Redis: %v", err)
// 	}
// 	fmt.Println(pong)
// }

func main() {

	// cli := redis.NewClusterClient(&redis.ClusterOptions{
	// 	Addrs: []string{
	// 		"10.0.0.215:6390",
	// 		"10.0.0.215:6391",
	// 		"10.0.0.215:6392",
	// 		"10.0.0.215:6393",
	// 		"10.0.0.215:6394",
	// 		"10.0.0.215:6395",
	// 	},
	// 	ReadTimeout:     10 * time.Second,
	// 	WriteTimeout:    10 * time.Second,
	// 	PoolSize:        10,
	// 	MaxIdleConns:    10,
	// 	ConnMaxIdleTime: 30 * time.Second,
	// 	NewClient: func(opt *redis.Options) *redis.Client {
	// 		return redis.NewClient(opt)
	// 	},
	// })
	// ctx := context.Background()
	// pong, err := cli.Ping(ctx).Result()
	// if err != nil {
	// 	fmt.Printf("Error connecting to Redis: %v", err)
	// 	log.Fatalf("Error connecting to Redis: %v", err)
	// }
	// fmt.Println("Connected to Redis:", pong)
	// res, err := cli.Set(context.Background(), "key", "value", time.Duration(0)).Result()
	// if err != nil {
	// 	panic(err)
	// }
	// println(res)
	// val, err := cli.Get(context.Background(), "key").Result()
	// if err != nil {
	// 	panic(err)
	// }
	// println(val)
	SendKafkaConverUser()
}

func SendKafkaConverUser() {
	config := &kafka.ConfigMap{
		"bootstrap.servers": "kafka01-dev.domob-inc.cn:9093,kafka02-dev.domob-inc.cn:9093,kafka03-dev.domob-inc.cn:9093",
		"client.id":         "test",
		"acks":              "all",
	}
	p, err := kafka.NewProducer(config)
	if err != nil {
		fmt.Println("Failed to create producer:", err)
		return
	}
	defer p.Close()
	topic := "test_dmp_conver_user_topic"
	crowdMap := map[string]*dmp_pb.Crowd{
		"interest-crowds": {
			CrowdIds: []int32{101, 105, 203},
		},
		"purchase-crowds": {
			CrowdIds: []int32{501, 510},
		},
	}
	msg := dmp_pb.UserProfileDataMsg{
		Type: dmp_pb.UserProfileDataMsgType_RESET,
		Msg: &dmp_pb.UserProfileData{
			Crowd: crowdMap,
		},
	}
	// proto to json
	msgbyte, err := proto.Marshal(&msg)
	if err != nil {
		fmt.Println("Failed to marshal msg:", err)
		return
	}
	// send to kafka
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
}
