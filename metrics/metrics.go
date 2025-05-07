package metrics

import (
	"sort"
	"strings"
	"sync"
	"time"

	"maps"

	"github.com/blueturbo-ad/go-utils/environment"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricsInstance *Metrics
	metricsOnce     sync.Once
)

const (
	MetricTypeBusiness = "business"
	MetricTypeSystem   = "system"
)

const (
	TimeDurSuffix = "milliseconds"
	CountSuffix   = "total"
)

type Metrics struct {
	gauges     sync.Map
	counters   sync.Map
	histograms sync.Map
	summaries  sync.Map

	keyPrefix string
	labels    map[string]string
	registry  *prometheus.Registry
}

func GetInstance() *Metrics {
	metricsOnce.Do(func() {
		metricsInstance = newMetrics()
	})

	return metricsInstance
}

func newMetrics() *Metrics {
	return &Metrics{
		gauges:     sync.Map{},
		counters:   sync.Map{},
		histograms: sync.Map{},
		summaries:  sync.Map{},
		keyPrefix:  "default",
		labels:     make(map[string]string),
		registry:   prometheus.NewRegistry(),
	}
}

func (m *Metrics) Init(keyPrefix string, labels map[string]string) error {
	maps.Copy(m.labels, labels)

	m.keyPrefix = keyPrefix
	m.labels["environment"] = environment.GetEnv()
	m.labels["type"] = MetricTypeSystem

	return nil
}

func (m *Metrics) SetGauge(name string, value float64) {
	m.SetGaugeWithLabels(name, value, nil)
}

func (m *Metrics) SetGaugeWithLabels(name string, value float64, labels map[string]string) {
	mergedLabels := m.mergeLabels(labels)
	key := strings.Join([]string{m.keyPrefix, name}, "_")
	gaugeVec, ok := m.gauges.Load(key)

	if !ok {
		labelNames := getSortedKeys(mergedLabels)
		newGaugeVec := prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: key,
				Help: "Gauge for " + name,
			},
			labelNames,
		)

		actualGaugeVec, loaded := m.gauges.LoadOrStore(key, newGaugeVec)
		if loaded {
			gaugeVec = actualGaugeVec
		} else {
			m.registry.MustRegister(newGaugeVec)
			gaugeVec = newGaugeVec
		}
	}

	gaugeVec.(*prometheus.GaugeVec).WithLabelValues(getLabelValues(mergedLabels)...).Set(value)
}

func (m *Metrics) IncrCounter(name string) {
	m.IncrCounterWithLabels(name, nil)
}

func (m *Metrics) IncrCounterWithLabels(name string, labels map[string]string) {
	mergedLabels := m.mergeLabels(labels)
	key := strings.Join([]string{m.keyPrefix, name, CountSuffix}, "_")
	counterVec, ok := m.counters.Load(key)

	if !ok {
		labelNames := getSortedKeys(mergedLabels)
		newCounterVec := prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: key,
				Help: "Counter for " + name,
			},
			labelNames,
		)

		actualCounterVec, loaded := m.counters.LoadOrStore(key, newCounterVec)
		if loaded {
			counterVec = actualCounterVec
		} else {
			m.registry.MustRegister(newCounterVec)
			counterVec = newCounterVec
		}
	}

	counterVec.(*prometheus.CounterVec).WithLabelValues(getLabelValues(mergedLabels)...).Inc()
}

func (m *Metrics) MeasureSince(name string, begTs time.Time) {
	m.MeasureSinceWithLabels(name, begTs, nil)
}

func (m *Metrics) MeasureSinceWithLabels(name string, begTs time.Time, labels map[string]string) {
	duration := float64(time.Since(begTs).Milliseconds())
	m.ObserveHistogramWithLabels(name, duration, labels)
}

func (m *Metrics) ObserveHistogram(name string, value float64) {
	m.ObserveHistogramWithLabels(name, value, nil)
}

func (m *Metrics) ObserveHistogramWithLabels(name string, value float64, labels map[string]string) {
	mergedLabels := m.mergeLabels(labels)
	key := strings.Join([]string{m.keyPrefix, name, TimeDurSuffix}, "_")
	histogramVec, ok := m.histograms.Load(key)

	if !ok {
		labelNames := getSortedKeys(mergedLabels)
		newHistogramVec := prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: key,
				Help: "Histogram for " + name,
			},
			labelNames,
		)

		actualHistogramVec, loaded := m.histograms.LoadOrStore(key, newHistogramVec)
		if loaded {
			histogramVec = actualHistogramVec
		} else {
			m.registry.MustRegister(newHistogramVec)
			histogramVec = newHistogramVec
		}
	}

	histogramVec.(*prometheus.HistogramVec).WithLabelValues(getLabelValues(mergedLabels)...).Observe(value)
}

func (m *Metrics) ObserveSummary(name string, value float64) {
	m.ObserveSummaryWithLabels(name, value, nil)
}

func (m *Metrics) ObserveSummaryWithLabels(name string, value float64, labels map[string]string) {
	mergedLabels := m.mergeLabels(labels)
	key := strings.Join([]string{m.keyPrefix, name, TimeDurSuffix}, "_")
	summaryVec, ok := m.summaries.Load(key)

	if !ok {
		labelNames := getSortedKeys(mergedLabels)
		newSummaryVec := prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name: key,
				Help: "Summary for " + name,
			},
			labelNames,
		)

		actualSummaryVec, loaded := m.summaries.LoadOrStore(key, newSummaryVec)
		if loaded {
			summaryVec = actualSummaryVec
		} else {
			m.registry.MustRegister(newSummaryVec)
			summaryVec = newSummaryVec
		}
	}

	summaryVec.(*prometheus.SummaryVec).WithLabelValues(getLabelValues(mergedLabels)...).Observe(value)
}

func (m *Metrics) GetRegistry() *prometheus.Registry {
	return m.registry
}

func getSortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}

func (m *Metrics) mergeLabels(labels map[string]string) map[string]string {
	if labels == nil {
		return m.labels
	}

	mergedLabels := make(map[string]string, len(m.labels)+len(labels))

	maps.Copy(mergedLabels, m.labels)
	maps.Copy(mergedLabels, labels)

	return mergedLabels
}

func getLabelValues(labels map[string]string) []string {
	labelValues := make([]string, 0, len(labels))
	for _, k := range getSortedKeys(labels) {
		labelValues = append(labelValues, labels[k])
	}

	return labelValues
}
