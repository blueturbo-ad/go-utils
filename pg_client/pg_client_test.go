package pgclient

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/blueturbo-ad/go-utils/config_manage"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TestGetSingleton 测试单例模式
func TestGetSingleton(t *testing.T) {
	t.Run("test singleton instance", func(t *testing.T) {
		instance1 := GetSingleton()
		instance2 := GetSingleton()
		if instance1 != instance2 {
			t.Error("GetSingleton() should return the same instance")
		}
	})

	t.Run("test singleton initial state", func(t *testing.T) {
		instance := GetSingleton()
		if instance == nil {
			t.Error("GetSingleton() should not return nil")
		}
	})

	t.Run("test concurrent singleton access", func(t *testing.T) {
		var wg sync.WaitGroup
		instances := make([]*PgClientManager, 100)
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				instances[idx] = GetSingleton()
			}(i)
		}
		wg.Wait()
		for i := 1; i < 100; i++ {
			if instances[i] != instances[0] {
				t.Error("GetSingleton() should return the same instance in concurrent access")
			}
		}
	})
}

// TestPgClientAttr 测试 PgClientAttr 结构体
func TestPgClientAttr(t *testing.T) {
	t.Run("test PgClientAttr creation", func(t *testing.T) {
		attr := &PgClientAttr{
			Nodes:    []string{"127.0.0.1"},
			Port:     5432,
			Database: "test_db",
			User:     "test_user",
			Password: "test_password",
			PoolSize: 20,
			Timeout:  5000,
			Env:      "Dev",
		}
		if attr.Nodes[0] != "127.0.0.1" {
			t.Errorf("Nodes[0] = %v; want 127.0.0.1", attr.Nodes[0])
		}
		if attr.Port != 5432 {
			t.Errorf("Port = %v; want 5432", attr.Port)
		}
		if attr.Database != "test_db" {
			t.Errorf("Database = %v; want test_db", attr.Database)
		}
		if attr.PoolSize != 20 {
			t.Errorf("PoolSize = %v; want 20", attr.PoolSize)
		}
	})
}

// TestInitialize 测试连接池初始化
func TestInitialize(t *testing.T) {
	manager := &PgClientManager{
		ReadClient:  [2]map[string]*pgxpool.Pool{},
		WriteClient: [2]map[string]*pgxpool.Pool{},
		index:       0,
	}

	t.Run("test initialize with empty nodes", func(t *testing.T) {
		attr := &PgClientAttr{
			Nodes:    []string{},
			Port:     5432,
			Database: "test_db",
			User:     "test_user",
			Password: "test_password",
			PoolSize: 20,
			Timeout:  5000,
			Env:      "Dev",
		}
		_, err := manager.Initialize(attr)
		if err == nil {
			t.Error("Initialize() with empty nodes should return error")
		}
		if err.Error() != "no nodes configured" {
			t.Errorf("Initialize() error = %v; want 'no nodes configured'", err)
		}
	})

	t.Run("test initialize with nil nodes", func(t *testing.T) {
		attr := &PgClientAttr{
			Nodes:    nil,
			Port:     5432,
			Database: "test_db",
			User:     "test_user",
			Password: "test_password",
			PoolSize: 20,
			Timeout:  5000,
			Env:      "Dev",
		}
		_, err := manager.Initialize(attr)
		if err == nil {
			t.Error("Initialize() with nil nodes should return error")
		}
	})

	t.Run("test initialize with invalid connection", func(t *testing.T) {
		attr := &PgClientAttr{
			Nodes:    []string{"invalid_host"},
			Port:     5432,
			Database: "test_db",
			User:     "test_user",
			Password: "test_password",
			PoolSize: 20,
			Timeout:  1000,
			Env:      "Dev",
		}
		_, err := manager.Initialize(attr)
		// 预期会失败，因为无法连接到无效主机
		if err == nil {
			t.Log("Initialize() succeeded unexpectedly - a real database might be running")
		}
	})
}

// TestBuildReadPgClient 测试构建读客户端
func TestBuildReadPgClient(t *testing.T) {
	manager := &PgClientManager{
		ReadClient:  [2]map[string]*pgxpool.Pool{},
		WriteClient: [2]map[string]*pgxpool.Pool{},
		index:       0,
	}

	t.Run("test build read client with empty nodes", func(t *testing.T) {
		conf := &config_manage.PostgreSQLConfig{
			Name: "test_pg",
		}
		conf.WritePool.Nodes = []string{}
		conf.WritePool.Port = 5432
		conf.WritePool.Database = "test_db"
		conf.WritePool.User = "test_user"
		conf.WritePool.Password = "test_password"
		conf.WritePool.PoolSize = 20
		conf.WritePool.Timeout = 5000

		result := manager.BuildReadPgClient(conf, "Dev")
		if result != nil {
			t.Error("BuildReadPgClient() with empty nodes should return nil")
		}
	})
}

// TestBuildWritePgClient 测试构建写客户端
func TestBuildWritePgClient(t *testing.T) {
	manager := &PgClientManager{
		ReadClient:  [2]map[string]*pgxpool.Pool{},
		WriteClient: [2]map[string]*pgxpool.Pool{},
		index:       0,
	}

	t.Run("test build write client with empty nodes", func(t *testing.T) {
		conf := &config_manage.PostgreSQLConfig{
			Name: "test_pg",
		}
		conf.WritePool.Nodes = []string{}
		conf.WritePool.Port = 5432
		conf.WritePool.Database = "test_db"
		conf.WritePool.User = "test_user"
		conf.WritePool.Password = "test_password"
		conf.WritePool.PoolSize = 20
		conf.WritePool.Timeout = 5000

		result := manager.BuildWritePgClient(conf, "Dev")
		if result != nil {
			t.Error("BuildWritePgClient() with empty nodes should return nil")
		}
	})
}

// TestPgClientManager 测试 PgClientManager 结构体
func TestPgClientManager(t *testing.T) {
	t.Run("test manager creation", func(t *testing.T) {
		manager := &PgClientManager{
			ReadClient:  [2]map[string]*pgxpool.Pool{},
			WriteClient: [2]map[string]*pgxpool.Pool{},
			index:       0,
		}
		manager.ReadClient[0] = make(map[string]*pgxpool.Pool)
		manager.WriteClient[0] = make(map[string]*pgxpool.Pool)

		if manager.ReadClient[0] == nil {
			t.Error("ReadClient[0] should not be nil")
		}
		if manager.WriteClient[0] == nil {
			t.Error("WriteClient[0] should not be nil")
		}
	})
}

// TestGetReadClient 测试获取读客户端
func TestGetReadClient(t *testing.T) {
	t.Run("test get non-existing read client", func(t *testing.T) {
		manager := &PgClientManager{
			ReadClient:  [2]map[string]*pgxpool.Pool{},
			WriteClient: [2]map[string]*pgxpool.Pool{},
			index:       0,
		}
		manager.ReadClient[0] = make(map[string]*pgxpool.Pool)
		manager.WriteClient[0] = make(map[string]*pgxpool.Pool)

		result := manager.GetReadClient("non_existing")
		if result != nil {
			t.Errorf("GetReadClient() = %v; want nil for non-existing client", result)
		}
	})
}

// TestGetWriteClient 测试获取写客户端
func TestGetWriteClient(t *testing.T) {
	t.Run("test get non-existing write client", func(t *testing.T) {
		manager := &PgClientManager{
			ReadClient:  [2]map[string]*pgxpool.Pool{},
			WriteClient: [2]map[string]*pgxpool.Pool{},
			index:       0,
		}
		manager.ReadClient[0] = make(map[string]*pgxpool.Pool)
		manager.WriteClient[0] = make(map[string]*pgxpool.Pool)

		result := manager.GetWriteClient("non_existing")
		if result != nil {
			t.Errorf("GetWriteClient() = %v; want nil for non-existing client", result)
		}
	})
}

// TestUpdateFromFile 测试从文件更新配置
func TestUpdateFromFile(t *testing.T) {
	t.Run("test update from invalid file", func(t *testing.T) {

		workPath, err := os.Getwd()
		if err != nil {
			fmt.Println(err.Error())
		}
		// 去workpath 的上一级目录
		works := strings.Split(workPath, "/")
		workPath = strings.Join(works[:len(works)-1], "/")

		cs := workPath + "/config/postgresql_conf.yaml"
		err = GetSingleton().UpdateFromFile(cs, "Dev")
		if err != nil {
			fmt.Println(err.Error())
		}
		var e = GetSingleton()
		client := e.GetWriteClient("mis_postgresql")
		if client == nil {
			t.Errorf("os.Getwd() = %v; want nil", client)
		}
		rows, err := client.Query(context.Background(), "SELECT id, name FROM advertiser ORDER BY name")
		if err != nil {
			t.Errorf("client.Query() error = %v; want nil", err)
		}
		defer rows.Close()
		for rows.Next() {
			var id int
			var name string
			err := rows.Scan(&id, &name)
			if err != nil {
				t.Errorf("rows.Scan() error = %v; want nil", err)
			}
			fmt.Printf("id: %d, name: %s\n", id, name)
		}
	})
}

// TestConcurrentAccess 测试并发访问
func TestConcurrentAccess(t *testing.T) {
	t.Run("test concurrent read and write access", func(t *testing.T) {
		manager := &PgClientManager{
			ReadClient:  [2]map[string]*pgxpool.Pool{},
			WriteClient: [2]map[string]*pgxpool.Pool{},
			index:       0,
		}
		manager.ReadClient[0] = make(map[string]*pgxpool.Pool)
		manager.WriteClient[0] = make(map[string]*pgxpool.Pool)

		var wg sync.WaitGroup
		concurrency := 50

		// 并发读取
		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = manager.GetReadClient("test_pg")
			}()
		}

		// 并发写入
		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = manager.GetWriteClient("test_pg")
			}()
		}

		wg.Wait()
	})
}

// TestEmptyString 测试空字符串常量
func TestEmptyString(t *testing.T) {
	t.Run("test empty string constant", func(t *testing.T) {
		if EmptyString != "" {
			t.Errorf("EmptyString = %v; want empty string", EmptyString)
		}
	})
}
