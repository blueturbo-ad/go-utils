package prometheustool

import (
	"flag"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestPrometheus(t *testing.T) {

	t.Run("TestPrometheus", func(t *testing.T) {
		a := GetSingleton().NewPrometheusGauge(&prometheus.GaugeOpts{
			Name:        "test",
			Help:        "test",
			ConstLabels: map[string]string{"test": "test"},
		}, "test")
		a.GetPrometheusGauge("test").Set(*flag.Float64("test", 0, "test"))

	})

}
