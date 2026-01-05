package config_manage

import (
	"log"
	"os"
	"testing"

	"github.com/blueturbo-ad/go-utils/environment"
	k8sclient "github.com/blueturbo-ad/go-utils/k8s_tool/k8s_client"
)

func TestPostgreSQLConfig(t *testing.T) {
	os.Setenv("POD_NAME", "test")
	os.Setenv("POD_NAMESPACE", "dsp-ns")
	environment.Init()
	k8sclient.GetSingleton().SetUp()

	t.Run("test postgresql config file", func(t *testing.T) {
		var p = new(PostgreSQLConfigManager)
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		filePath := dir + "/config/postgresql_conf.yaml"
		err = p.LoadConfig(filePath, "Dev")
		if err != nil {
			t.Errorf("LoadConfig() error = %v; want nil", err)
		}
		if p.Config == nil {
			t.Error("Config is nil")
			return
		}
		if len(*p.Config) == 0 {
			t.Error("Config is empty")
			return
		}
		// 验证配置内容
		config := (*p.Config)[0]
		if config.Name != "mis_postgresql" {
			t.Errorf("config.Name = %v; want mis_postgresql", config.Name)
		}
		if config.WritePool.Database != "blueturbo_business" {
			t.Errorf("config.WritePool.Database = %v; want blueturbo_business", config.WritePool.Database)
		}
		if config.WritePool.PoolSize != 20 {
			t.Errorf("config.WritePool.PoolSize = %v; want 20", config.WritePool.PoolSize)
		}
		if config.WritePool.Port != 5432 {
			t.Errorf("config.WritePool.Port = %v; want 5432", config.WritePool.Port)
		}
	})

	t.Run("test postgresql configmap", func(t *testing.T) {
		var p = new(PostgreSQLConfigManager)
		err := p.LoadK8sConfigMap("postgresql-conf", "Dev")
		if err != nil {
			t.Errorf("LoadK8sConfigMap() error = %v; want nil", err)
		}
	})

	t.Run("test postgresql memory config", func(t *testing.T) {
		yamlData := []byte(`
postgresql_conf:
  - name: mis_postgresql
    write_pool:
      database: blueturbo_business
      pool_size: 20
      timeout: 5000
      port: 5432
      nodes:
        - 34.142.238.89
      user: blueturbo
      password: 7n<{qqu2pgh97ns~
    read_pool:
      database: blueturbo_business
      pool_size: 20
      timeout: 5000
      port: 5432
      nodes:
        - 34.142.238.89
      user: blueturbo_read
      password: Cu5B3%l3M1pRXa0q
`)
		var p = new(PostgreSQLConfigManager)
		err := p.LoadMemoryConfig(yamlData, "Dev")
		if err != nil {
			t.Errorf("LoadMemoryConfig() error = %v; want nil", err)
		}
		if p.Config == nil {
			t.Error("Config is nil")
			return
		}
		if len(*p.Config) == 0 {
			t.Error("Config is empty")
			return
		}
		// 验证配置内容
		config := (*p.Config)[0]
		if config.Name != "mis_postgresql" {
			t.Errorf("config.Name = %v; want mis_postgresql", config.Name)
		}
		if config.ReadPool.User != "blueturbo_read" {
			t.Errorf("config.ReadPool.User = %v; want blueturbo_read", config.ReadPool.User)
		}
	})

	t.Run("test postgresql config with Pro env", func(t *testing.T) {
		var p = new(PostgreSQLConfigManager)
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		filePath := dir + "/config/postgresql_conf.yaml"
		err = p.LoadConfig(filePath, "Pro")
		if err != nil {
			t.Errorf("LoadConfig() with Pro env error = %v; want nil", err)
		}
		if p.Config == nil {
			t.Error("Config is nil")
			return
		}
	})

	t.Run("test postgresql config with invalid file", func(t *testing.T) {
		var p = new(PostgreSQLConfigManager)
		err := p.LoadConfig("/invalid/path/config.yaml", "Dev")
		if err == nil {
			t.Error("LoadConfig() with invalid file should return error")
		}
	})

	t.Run("test postgresql config with invalid env", func(t *testing.T) {
		var p = new(PostgreSQLConfigManager)
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		filePath := dir + "/config/postgresql_conf.yaml"
		err = p.LoadConfig(filePath, "InvalidEnv")
		if err == nil {
			t.Error("LoadConfig() with invalid env should return error")
		}
	})
}
