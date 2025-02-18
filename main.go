package main

import (
	"os"

	"github.com/blueturbo-ad/go-utils/environment"
	gcpcloudstorage "github.com/blueturbo-ad/go-utils/gcp_cloud_tool/gcp_cloud_storage"
	k8sclient "github.com/blueturbo-ad/go-utils/k8s_tool/k8s_client"
)

func main() {
	os.Setenv("POD_NAME", "test")
	os.Setenv("POD_NAMESPACE", "dsp-ns")
	environment.Init()
	k8sclient.GetSingleton().SetUp()
	gcpcloudstorage.GetSingleton().UpdateLoadK8sConfigMap("gcp-cloud-storage-config", "Dev")

	// var e = redisclient.GetSingleton()
	// workPath, err := os.Getwd()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// p := workPath + "/config/redis_conf.yaml"
	// err = e.UpdateFromFile(p, "Dev")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// client := e.GetReadClient("event_redis")
	// if client == nil {
	// 	fmt.Println(err.Error())
	// }
	// ctx := context.Background()
	// res, err := client.Get(ctx, "test").Result()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// fmt.Println(res)
}
