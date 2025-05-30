package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/blueturbo-ad/dsp_data_interface/gen/go/dmp_pb"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/redis/go-redis/v9"
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
	// SendKafkaConverUser()
	// ReadRedis()
	// ReadBundleRedis()
	// SendModelKafka()
	// SendKafkaPromoConverUser()
	// ReadKafkaCrowds()
	ReadModelRedis()
}
func ReadModelRedis() {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{"10.152.0.27:6379"},
		NewClient: func(opt *redis.Options) *redis.Client {
			return redis.NewClient(opt)
		},
	})
	res, err := client.HGetAll(context.Background(), "123123131").Result()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	interest := res["model_fea"]
	interestCrowd := &dmp_pb.ModelFeature{}
	fmt.Println("interest:", []byte(interest))
	err = proto.Unmarshal([]byte(interest), interestCrowd)
	if err != nil {
		fmt.Println("Error unmarshalling interestCrowd:", err)
		return
	}
	interestCrowdJson, err := json.Marshal(interestCrowd)
	if err != nil {
		fmt.Println("Error marshalling interestCrowd to JSON:", err)
		return
	}
	fmt.Println("interestCrowdJson:", string(interestCrowdJson))
}
func ReadBundleRedis() {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{"10.152.0.27:6379"},
		NewClient: func(opt *redis.Options) *redis.Client {
			return redis.NewClient(opt)
		},
	})
	res, err := client.HGetAll(context.Background(), "123123131").Result()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Result:", res)
	interest := res["com.dubox.drive"]
	interestCrowd := &dmp_pb.PromoBundleBehavior{}
	fmt.Println("interest:", []byte(interest))
	err = proto.Unmarshal([]byte(interest), interestCrowd)
	if err != nil {
		fmt.Println("Error unmarshalling interestCrowd:", err)
		return
	}
	interestCrowdJson, err := json.Marshal(interestCrowd)
	if err != nil {
		fmt.Println("Error marshalling interestCrowd to JSON:", err)
		return
	}
	fmt.Println("interestCrowdJson:", string(interestCrowdJson))
}

func ReadRedis() {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{"10.152.0.27:6379"},
		NewClient: func(opt *redis.Options) *redis.Client {
			return redis.NewClient(opt)
		},
	})
	res, err := client.HGetAll(context.Background(), "123123131").Result()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Result:", res)
	interest := res["interest-crowds"]
	//base64解码
	// interestByte, err := base64.RawStdEncoding.DecodeString(interest)
	// if err != nil {
	// 	fmt.Println("Error decoding base64:", err)
	// 	return
	// }
	// fmt.Println("interest:", string(interestByte))
	// grpc to json
	interestCrowd := &dmp_pb.Crowd{}
	fmt.Println("interest:", []byte(interest))
	err = proto.Unmarshal([]byte(interest), interestCrowd)
	if err != nil {
		fmt.Println("Error unmarshalling interestCrowd:", err)
		return
	}
	interestCrowdJson, err := json.Marshal(interestCrowd)
	if err != nil {
		fmt.Println("Error marshalling interestCrowd to JSON:", err)
		return
	}
	fmt.Println("interestCrowdJson:", string(interestCrowdJson))

}

func SendKafkaConverUser() {
	config := &kafka.ConfigMap{
		"bootstrap.servers": "10.152.3.194:9092,10.152.3.195:9092,10.152.3.246:9092",
		"client.id":         "test",
		"acks":              "all",
	}
	p, err := kafka.NewProducer(config)
	if err != nil {
		fmt.Println("Failed to create producer:", err)
		return
	}
	defer p.Close()
	topic := "test_news_dmp_conver_user_topic"
	msg := dmp_pb.UserProfileDataMsg{
		Type:     *dmp_pb.UserProfileDataMsgType_RESET.Enum(),
		RedisTtl: 1000,
		Msg: &dmp_pb.UserProfileData{
			BtUid: "123123131",
			Crowd: map[string]*dmp_pb.Crowd{
				"interest-crowds": {
					CrowdIds: []int32{101, 105, 203, 301, 401},
				},
				"purchase-crowds": {
					CrowdIds: []int32{501, 510},
				},
			},
		},
	}

	// proto to json
	msgbyte, err := json.Marshal(&msg)
	if err != nil {
		fmt.Println("Failed to marshal msg:", err)
		return
	}
	fmt.Println("msgbyte:", string(msgbyte))
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

func ReadKafkaCrowds() {
	config := &kafka.ConfigMap{
		"bootstrap.servers":  "10.152.3.194:9092,10.152.3.195:9092,10.152.3.246:9092",
		"group.id":           "win_group3",
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	}
	c, err := kafka.NewConsumer(config)
	if err != nil {
		fmt.Println("Failed to create consumer:", err)
		return
	}
	defer c.Close()
	topic := "test_news_dmp_conver_user_topic"
	err = c.Subscribe(topic, nil)
	if err != nil {
		fmt.Printf("failed to subscribe topic: %s\n", err)
		return
	}
	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			// c.CommitMessage(msg)
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
		fmt.Println("msg:", msg)
		val := make(map[string]interface{})
		err = json.Unmarshal(msg.Value, &val)
		if err != nil {
			fmt.Println("Error unmarshalling message:", err)
			continue
		}
		fmt.Println("val:", val)
		msgVal := val["msg"]
		crowd, ok := msgVal.(map[string]interface{})["crowd"]
		if !ok {
			fmt.Println("interest-crowds not found in msg")
			continue
		}
		crowdMap, ok := crowd.(map[string]interface{})
		if !ok {
			fmt.Println("interest-crowds is not a map")
			continue
		}
		interest, ok := crowdMap["interest-crowds"]
		if !ok {
			fmt.Println("interest-crowds not found in msg")
			continue
		}
		interestMap, ok := interest.(map[string]interface{})
		if !ok {
			fmt.Println("interest-crowds is not a map")
			continue
		}
		var crowdIds []int32
		crowdIdsInterface, ok := interestMap["crowd_ids"]
		if !ok {
			fmt.Println("crowd_ids not found in interest-crowds")
			continue
		}
		for _, id := range crowdIdsInterface.([]interface{}) {
			// JSON 中的数字默认解析为 float64
			if idFloat, ok := id.(float64); ok {
				crowdIds = append(crowdIds, int32(idFloat))
			} else {
				fmt.Printf("非数字元素: %T %v\n", id, id)
			}
		}

		protoCrowd := &dmp_pb.Crowd{
			CrowdIds: crowdIds,
		}

		protoCrowdJson, err := proto.Marshal(protoCrowd)
		if err != nil {
			fmt.Println("Error marshalling interestCrowd to JSON:", err)
			continue
		}
		// fmt.Println("protoCrowdJson:", string(protoCrowdJson))
		// // base64编码
		// interestBase64 := base64.RawStdEncoding.EncodeToString(protoCrowdJson)
		// fmt.Println("interestBase64:", interestBase64)
		// redis set
		client := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: []string{"10.152.0.27:6379"},
			NewClient: func(opt *redis.Options) *redis.Client {
				return redis.NewClient(opt)
			},
		})
		client.HSet(context.Background(), "123123131", "interest-crowds", protoCrowdJson).Err()

	}
}

func SendKafkaPromoConverUser() {
	config := &kafka.ConfigMap{
		"bootstrap.servers": "10.152.3.194:9092,10.152.3.195:9092,10.152.3.246:9092",
		"client.id":         "test",
		"acks":              "all",
	}
	p, err := kafka.NewProducer(config)
	if err != nil {
		fmt.Println("Failed to create producer:", err)
		return
	}
	defer p.Close()
	topic := "test_news_dmp_conver_user_topic"
	msg := dmp_pb.UserProfileDataMsg{
		Type: *dmp_pb.UserProfileDataMsgType_UPDATE.Enum(),
		Msg: &dmp_pb.UserProfileData{
			BtUid: "123123131",
			// Crowd: map[string]*dmp_pb.Crowd{
			// 	"interest-crowds": {
			// 		CrowdIds: []int32{101, 105, 203, 301, 401},
			// 	},
			// 	"purchase-crowds": {
			// 		CrowdIds: []int32{501, 510},
			// 	},
			// },
			PromoBundleBehavior: map[string]*dmp_pb.PromoBundleBehavior{
				"com.dubox.drive": {
					ImpTs:     []int32{1, 2, 3, 4, 6, 7},
					ClkTs:     []int32{1, 2, 5, 7, 8},
					InstallTs: []int32{1, 2},
					EventTs:   []int32{3, 4},
				},
			},
		},
	}

	// proto to json
	msgbyte, err := json.Marshal(&msg)
	if err != nil {
		fmt.Println("Failed to marshal msg:", err)
		return
	}
	fmt.Println("msgbyte:", string(msgbyte))
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

func SendModelKafka() {
	config := &kafka.ConfigMap{
		"bootstrap.servers": "10.152.3.194:9092,10.152.3.195:9092,10.152.3.246:9092",
		"client.id":         "test",
		"acks":              "all",
	}
	p, err := kafka.NewProducer(config)
	if err != nil {
		fmt.Println("Failed to create producer:", err)
		return
	}
	defer p.Close()
	topic := "test_news_dmp_conver_user_topic"
	msg := dmp_pb.UserProfileDataMsg{
		Type: *dmp_pb.UserProfileDataMsgType_RESET.Enum(),
		Msg: &dmp_pb.UserProfileData{
			BtUid: "123123131",
			// Crowd: map[string]*dmp_pb.Crowd{
			// 	"interest-crowds": {
			// 		CrowdIds: []int32{101, 105, 203, 301, 401},
			// 	},
			// 	"purchase-crowds": {
			// 		CrowdIds: []int32{501, 510},
			// 	},
			// },
			// PromoBundleBehavior: map[string]*dmp_pb.PromoBundleBehavior{
			// 	"com.dubox.drive": {
			// 		ImpTs:     []int32{1, 2, 3, 4, 6, 7},
			// 		ClkTs:     []int32{1, 2, 5, 7, 8},
			// 		InstallTs: []int32{1, 2},
			// 		EventTs:   []int32{3, 4},
			// 	},
			// },
			ModelFea: &dmp_pb.ModelFeature{
				CreativeClkSeq: "44444",
				BundleDateSeq:  "123123131",
			},
		},
	}

	// proto to json
	msgbyte, err := json.Marshal(&msg)
	if err != nil {
		fmt.Println("Failed to marshal msg:", err)
		return
	}
	fmt.Println("msgbyte:", string(msgbyte))
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
