package prometheustool

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestPrometheus(t *testing.T) {

	t.Run("TestPrometheus", func(t *testing.T) {
		a, err := GetSingleton().NewPrometheusGauge(&prometheus.GaugeOpts{
			Name:        "test",
			Help:        "test",
			ConstLabels: map[string]string{"test": "test"},
		}, "test")
		if err != nil {
			t.Errorf("Error creating Prometheus gauge: %v", err)
		}

		c, err := a.GetPrometheusGauge("test")
		if err != nil {
			t.Errorf("Error getting Prometheus gauge: %v", err)
		}
		c.(prometheus.Gauge).Set(1)

	})

}
