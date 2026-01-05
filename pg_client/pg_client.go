package pgclient

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/blueturbo-ad/go-utils/config_manage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgClientAttr struct {
	Nodes    []string
	Port     int
	Database string
	User     string
	Password string
	PoolSize int
	Timeout  int
	Env      string
}

type PgClientManager struct {
	ReadClient  [2]map[string]*pgxpool.Pool
	WriteClient [2]map[string]*pgxpool.Pool
	index       int
	rwMutex     sync.RWMutex
}

var (
	instance    *PgClientManager
	once        sync.Once
	EmptyString = ""
)

func GetSingleton() *PgClientManager {
	once.Do(func() {
		instance = &PgClientManager{
			index: -1,
		}

	})
	return instance
}

func (p *PgClientManager) GetReadClient(name string) *pgxpool.Pool {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()
	redisClient := p.ReadClient[p.index][name]
	if redisClient != nil {
		return redisClient
	}
	return nil
}

func (p *PgClientManager) GetWriteClient(name string) *pgxpool.Pool {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()
	if p.WriteClient[p.index][name] != nil {
		return p.WriteClient[p.index][name]
	}
	return nil
}

func (p *PgClientManager) UpdateLoadK8sConfigMap(configMapName, env string, hookName string) error {
	var e = new(config_manage.PostgreSQLConfigManager)
	err := e.LoadK8sConfigMap(configMapName, env)
	if err != nil {
		return fmt.Errorf("redis client  LoadK8sConfigMap is error %s", err.Error())
	}
	return p.refreshPgClient(e, env)
}

func (p *PgClientManager) UpdateFromFile(confPath string, env string) error {
	var err error
	var e = new(config_manage.PostgreSQLConfigManager)
	err = e.LoadConfig(confPath, env)
	if err != nil {
		return err
	}
	return p.refreshPgClient(e, env)
}
func (p *PgClientManager) refreshPgClient(confs *config_manage.PostgreSQLConfigManager, env string) error {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()
	newIndex := (p.index + 1) % 2
	p.ReadClient[newIndex] = make(map[string]*pgxpool.Pool)
	p.WriteClient[newIndex] = make(map[string]*pgxpool.Pool)
	for _, v := range *confs.Config {
		p.ReadClient[newIndex][v.Name] = p.BuildReadPgClient(&v, env)
		p.WriteClient[newIndex][v.Name] = p.BuildWritePgClient(&v, env)
	}
	p.index = newIndex
	return nil
}

func (p *PgClientManager) BuildReadPgClient(conf *config_manage.PostgreSQLConfig, env string) *pgxpool.Pool {
	pgAttr := &PgClientAttr{
		Nodes:    conf.WritePool.Nodes,
		Port:     conf.WritePool.Port,
		Database: conf.WritePool.Database,
		User:     conf.WritePool.User,
		Password: conf.WritePool.Password,
		PoolSize: conf.WritePool.PoolSize,
		Timeout:  conf.WritePool.Timeout,
		Env:      env,
	}
	if pool, err := p.Initialize(pgAttr); err != nil {
		return nil
	} else {
		return pool
	}
}

func (p *PgClientManager) BuildWritePgClient(conf *config_manage.PostgreSQLConfig, env string) *pgxpool.Pool {
	pgAttr := &PgClientAttr{
		Nodes:    conf.WritePool.Nodes,
		Port:     conf.WritePool.Port,
		Database: conf.WritePool.Database,
		User:     conf.WritePool.User,
		Password: conf.WritePool.Password,
		PoolSize: conf.WritePool.PoolSize,
		Timeout:  conf.WritePool.Timeout,
		Env:      env,
	}
	if pool, err := p.Initialize(pgAttr); err != nil {
		return nil
	} else {
		return pool
	}
}

// Initialize 初始化连接池
func (p *PgClientManager) Initialize(pgAttr *PgClientAttr) (*pgxpool.Pool, error) {
	if len(pgAttr.Nodes) == 0 {
		return nil, fmt.Errorf("no nodes configured")
	}

	// 构建连接字符串
	encodedPassword := url.QueryEscape(pgAttr.Password)

	// 构建连接字符串
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?pool_max_conns=%d&pool_min_conns=%d&pool_max_conn_lifetime=%ds",
		pgAttr.User,
		encodedPassword,
		pgAttr.Nodes[0],
		pgAttr.Port,
		pgAttr.Database,
		pgAttr.PoolSize,
		pgAttr.PoolSize/4, // 最小连接数设为最大的1/4
		3600,              // 连接最大生命周期1小时
	)

	// 创建连接池配置
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// 设置连接池参数
	config.MaxConns = int32(pgAttr.PoolSize)
	config.MinConns = int32(pgAttr.PoolSize / 4)
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute
	config.ConnConfig.ConnectTimeout = time.Duration(pgAttr.Timeout) * time.Millisecond

	// 创建连接池
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pgAttr.Timeout)*time.Millisecond)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	// 测试连接
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return pool, nil
}
