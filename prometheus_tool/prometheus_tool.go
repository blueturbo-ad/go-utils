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
	c prometheus.Gauge
}

func GetSingleton() *PrometheusTool {
	once.Do(func() {
		instance = &PrometheusTool{}
	})
	return instance
}

func (p *PrometheusTool) NewPrometheusGauge(option *prometheus.GaugeOpts) *PrometheusTool {
	p.c = prometheus.NewGauge(*option)

	return p
}

func (p *PrometheusTool) SetVal(val any) {
	switch val.(type) {
	case *int:
		val := val.(*int)
		if val == nil {
			p.c.Set(float64(0))
		} else {
			p.c.Set(float64(*val))
		}
	case *int32:
		val := val.(*int32)
		if val == nil {
			p.c.Set(float64(0))
		} else {
			p.c.Set(float64(*val))
		}
	case *int64:
		val := val.(*int64)
		if val == nil {
			p.c.Set(float64(0))
		} else {
			p.c.Set(float64(*val))
		}
	case *float32:
		val := val.(*float32)
		if val == nil {
			p.c.Set(float64(0))
		} else {
			p.c.Set(float64(*val))
		}
	case *float64:
		val := val.(*float64)
		if val == nil {
			p.c.Set(float64(0))
		} else {
			p.c.Set(float64(*val))
		}
	case *string:
		val := val.(*string)
		if val == nil {
			p.c.Set(float64(0))
		} else {
			p.c.Set(float64(0))
		}
	case *bool:
		val := val.(*bool)
		if val == nil {
			p.c.Set(float64(0))
		} else {
			if *val {
				p.c.Set(float64(1))
			} else {
				p.c.Set(float64(0))
			}
		}
	default:
		p.c.Set(float64(0))
	}
}
