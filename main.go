package main

import (
	"fmt"
	"os"

	"github.com/blueturbo-ad/go-utils/environment"
	gcpcloudstorage "github.com/blueturbo-ad/go-utils/gcp_cloud_tool/gcp_cloud_storage"
	gcpcloudtool "github.com/blueturbo-ad/go-utils/gcp_cloud_tool/gcp_svc_acc_token"
)

func main() {
	os.Setenv("POD_NAME", "test")
	os.Setenv("POD_NAMESPACE", "dsp-ns")
	environment.Init()
	workPath, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
	}
	cs := workPath + "/config/cloud_storage.yaml"
	err = gcpcloudstorage.GetSingleton().UpdateFromFile(cs, "Dev")
	if err != nil {
		fmt.Println(err.Error())
	}
	ac := workPath + "/config/gcp_acc_token_conf.yaml"
	err = gcpcloudtool.GetSingleton().UpdateFromFile(ac, "Dev")
	if err != nil {
		fmt.Println(err.Error())
	}
	// var e = redisclient.GetSingleton()

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
