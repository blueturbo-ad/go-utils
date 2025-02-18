package gcpcloudtool

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/blueturbo-ad/go-utils/config_manage"
	gcpcloudstorage "github.com/blueturbo-ad/go-utils/gcp_cloud_tool/gcp_cloud_storage"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/api/option"
	gtransport "google.golang.org/api/transport/grpc"
	"google.golang.org/genproto/googleapis/iam/credentials/v1"
)

var (
	instance    *GcpSvcAccountToken
	once        sync.Once
	EmptyString = ""
)

type GcpSvcAccountToken struct {
	Tokens map[string]string
	confs  []config_manage.CloudAcc
}

func GetSingleton() *GcpSvcAccountToken {
	once.Do(func() {
		instance = &GcpSvcAccountToken{}
	})
	return instance
}

func (g *GcpSvcAccountToken) GetToken(accounName string) string {
	return g.Tokens[accounName]
}

func (g *GcpSvcAccountToken) UpdateLoadK8sConfigMap(configMapName, env string) error {
	var e = new(config_manage.GcpSvcAccountTokenConfig)
	err := e.LoadK8sConfigMap(configMapName, env)
	if err != nil {
		return fmt.Errorf("redis client  LoadK8sConfigMap is error %s", err.Error())
	}
	return g.retrieveToken(e.CloudAcc)
}

// 函数用于内存更新etcd配置
func (r *GcpSvcAccountToken) UpdateFromEtcd(env string, eventType string, key string, value string) error {
	fmt.Printf("Event Type: %s, Key: %s, Value: %s\n", eventType, key, value)

	var err error
	switch key {
	case "redis-conf":
		var e = new(config_manage.GcpSvcAccountTokenConfig)
		err = e.LoadMemoryConfig([]byte(value), env)
		if err != nil {
			return err
		}
		if err := r.retrieveToken(e.CloudAcc); err != nil {
			return err
		}
	default:
		return nil
	}
	return nil
}

func (f *GcpSvcAccountToken) UpdateFromFile(confPath string, env string) error {
	var err error
	var e = new(config_manage.GcpSvcAccountTokenConfig)
	err = e.LoadConfig(confPath, env)
	if err != nil {
		return err
	}
	return f.retrieveToken(e.CloudAcc)
}

func (g *GcpSvcAccountToken) retrieveToken(confs []config_manage.CloudAcc) error {
	g.confs = confs
	for _, v := range confs {
		ctx := context.Background()
		conn, err := gtransport.Dial(ctx,
			option.WithEndpoint("iamcredentials.googleapis.com:443"),
			option.WithScopes("https://www.googleapis.com/auth/cloud-platform"))

		if err != nil {
			return err
		}
		client := credentials.NewIAMCredentialsClient(conn)

		svcAccount := flag.String("a", "projects/-/serviceAccounts/blueturbo-dsp-deploy@blueturbo.iam.gserviceaccount.com", v.AccountEmail)
		req := credentials.GenerateAccessTokenRequest{
			Name:     *svcAccount,
			Scope:    []string{"https://www.googleapis.com/auth/cloud-platform"},
			Lifetime: ptypes.DurationProto(*flag.Duration("d", time.Hour, "lifetime of token")),
		}
		rsp, err := client.GenerateAccessToken(ctx, &req)
		if err != nil {
			return err
		}
		if g.Tokens == nil {
			g.Tokens = make(map[string]string)
		}
		g.Tokens[v.Name] = rsp.AccessToken

	}
	t, err := json.Marshal(g.Tokens)
	if err != nil {
		return err
	}
	ctx := context.Background()
	client := gcpcloudstorage.GetSingleton().GetClient("dsp_bucket")
	if client == nil {
		return fmt.Errorf("failed to get GCP cloud storage client")
	}
	wc := client.Object("access_token.json").NewWriter(ctx)
	if _, err := wc.Write(t); err != nil {
		return err
	}
	if err := wc.Close(); err != nil { // 检查 wc.Close() 的返回值
		fmt.Printf("Error closing writer: %+v\n", err)
		return err
	}

	return nil
}

func (g *GcpSvcAccountToken) RefreshToken() error {
	return g.retrieveToken(g.confs)
}
