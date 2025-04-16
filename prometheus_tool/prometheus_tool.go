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
	c        map[string]any
	reigster *prometheus.Registry
}

func GetSingleton() *PrometheusTool {
	once.Do(func() {
		instance = &PrometheusTool{
			c:        make(map[string]any),
			reigster: prometheus.NewRegistry(),
		}
	})
	return instance
}

func (p *PrometheusTool) GetPrometheus(name string) (any, error) {
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
		c := prometheus.NewGauge(*option)
		p.reigster.MustRegister(c)
		p.c[name] = c

	}

	return p, nil
}

// 统计qps的
func (p *PrometheusTool) NewPrometheusCounter(option *prometheus.CounterOpts, name string) (*PrometheusTool, error) {
	if option == nil {
		return nil, fmt.Errorf("option is nil")
	}
	if name == "" {
		name = option.Name
	}

	if _, ok := p.c[name]; ok {
		return p, nil
	} else {
		c := prometheus.NewCounter(*option)
		p.reigster.MustRegister(c)
		p.c[name] = c

	}

	return p, nil
}

func (p *PrometheusTool) NewPrometheusHistogram(option *prometheus.HistogramOpts, name string) (*PrometheusTool, error) {
	if option == nil {
		return nil, fmt.Errorf("option is nil")
	}
	if name == "" {
		name = option.Name
	}

	if _, ok := p.c[name]; ok {
		return p, nil
	} else {
		c := prometheus.NewHistogram(*option)
		p.reigster.MustRegister(c)
		p.c[name] = c

	}

	return p, nil
}

// 计算通过率的
func (p *PrometheusTool) NewPrometheusCounterVec(option *prometheus.CounterOpts, labelNames []string, name string) (*PrometheusTool, error) {
	if option == nil {
		return nil, fmt.Errorf("option is nil")
	}
	if name == "" {
		name = option.Name
	}

	if _, ok := p.c[name]; ok {
		return p, nil
	} else {
		c := prometheus.NewCounterVec(*option, labelNames)
		p.reigster.MustRegister(c)
		p.c[name] = c

	}

	return p, nil
}
func (P PrometheusTool) GetPrometheusRegister() *prometheus.Registry {
	return P.reigster
}
