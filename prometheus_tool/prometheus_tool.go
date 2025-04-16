package prometheustool

import (
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

func (p *PrometheusTool) GetPrometheusGauge(name string) prometheus.Gauge {
	if name == "" {
		return nil
	}
	if c, ok := p.c[name]; ok {
		return c
	}
	return nil
}

func (p *PrometheusTool) NewPrometheusGauge(option *prometheus.GaugeOpts, name string) *PrometheusTool {
	if option == nil {
		return nil
	}
	if name == "" {
		name = option.Name
	}
	c := prometheus.NewGauge(*option)
	if _, ok := p.c[name]; ok {
		return p
	}
	p.c[name] = c
	return p
}
