package workpool

import (
	"errors"
	"fmt"
	"sync"

	"github.com/blueturbo-ad/go-utils/config_manage"
	"github.com/panjf2000/ants"
)

var (
	oncePool sync.Once
	instance *WorkPool
)

type WorkPool struct {
	Pools map[string]*ants.Pool
}

func GetSingleton() *WorkPool {
	oncePool.Do(func() {
		instance = NewAntsPool()

	})
	return instance
}

func NewAntsPool() *WorkPool {

	return &WorkPool{
		Pools: make(map[string]*ants.Pool),
	}
}

func (w *WorkPool) UpdateLoadK8sConfigMap(configMapName, env string, hookName string) error {
	var e = new(config_manage.WorkPoolConfigManager)
	err := e.LoadK8sConfigMap(configMapName, env)
	if err != nil {
		return fmt.Errorf("kafka client  LoadK8sConfigMap is error %s", err.Error())
	}
	return w.WorkPoolTune(e)
}
func (w *WorkPool) InitLoadK8sConfigMap(configMapName, env string, hookName string) error {
	var e = new(config_manage.WorkPoolConfigManager)
	err := e.LoadK8sConfigMap(configMapName, env)
	if err != nil {
		return fmt.Errorf("kafka client  LoadK8sConfigMap is error %s", err.Error())
	}
	return w.BuildWorkPool(e)
}

func (w *WorkPool) UpdateFromFile(confPath string, env string) error {
	var err error
	var e = new(config_manage.WorkPoolConfigManager)
	err = e.LoadConfig(confPath, env)
	if err != nil {
		return err
	}
	return w.BuildWorkPool(e)
}

func (w *WorkPool) BuildWorkPool(e *config_manage.WorkPoolConfigManager) error {
	var Pool *ants.Pool
	var err error
	for _, conf := range *e.Config {
		Pool, err = ants.NewPool(
			conf.PoolSize,
			ants.WithPreAlloc(true),
			ants.WithNonblocking(true),
		)
		if err != nil {
			return errors.New("failed to create ants pool: " + err.Error())
		}
		w.Pools[conf.Name] = Pool
	}
	return nil
}

func (p *WorkPool) GetGinCtxPool(key string) (*ants.Pool, error) {
	if p.Pools == nil {
		return nil, errors.New("ants pool is not initialized")
	}
	pool, exists := p.Pools[key]
	if !exists {
		return nil, fmt.Errorf("ants pool with key %s not found", key)
	}

	return pool, nil
}

func (p *WorkPool) WorkPoolTune(e *config_manage.WorkPoolConfigManager) error {
	if p.Pools == nil {
		p.Pools = make(map[string]*ants.Pool)
	}
	for _, conf := range *e.Config {
		if pool, exists := p.Pools[conf.Name]; exists {
			if pool.Cap() != conf.PoolSize {
				pool.Tune(uint(conf.PoolSize))
			}
		}
	}
	return nil
}

func (p *WorkPool) Release() {
	if p.Pools != nil {
		for key, pool := range p.Pools {
			pool.Release()
			delete(p.Pools, key)
		}
	}
}
