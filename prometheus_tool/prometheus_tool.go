package prometheustool

import (
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type Message struct {
	Val any
}

var (
	instance    *PrometheusTool
	once        sync.Once
	EmptyString = ""
)

type PrometheusTool struct {
	c map[string]prometheus.Gauge
}

func GetSingleton() *PrometheusTool {
	once.Do(func() {
		instance = &PrometheusTool{
			c: make(map[string]prometheus.Gauge),
		}
	})
	return instance
}

func (p *PrometheusTool) GetPrometheusGauge(name string) (prometheus.Gauge, error) {
	if name == "" {
		return nil, fmt.Errorf("name is empty")
	}
	if c, ok := p.c[name]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("prometheus gauge not found")
}

func (p *PrometheusTool) NewPrometheusGauge(option *prometheus.GaugeOpts, name string) (*PrometheusTool, error) {
	if option == nil {
		return nil, fmt.Errorf("option is nil")
	}
	if name == "" {
		name = option.Name
	}

	if _, ok := p.c[name]; ok {
		return p, nil
	} else {
		reigster := prometheus.NewRegistry()
		c := prometheus.NewGauge(*option)
		if err := reigster.Register(c); err != nil {
			return nil, fmt.Errorf("register prometheus gauge failed: %v", err)
		}
		p.c[name] = c

	}

	return p, nil
}
