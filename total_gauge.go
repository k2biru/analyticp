package analyticp

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type TotalGauge interface {
	Set(tags []string, value float64)
	Reset()
	Inc(tags []string)
}
type TotalGaugeConfig struct {
	App          string
	Name         string
	Tags         []string
	MustRegister bool
}

func NewTotalGauge(config *TotalGaugeConfig) TotalGauge {
	opt := prometheus.GaugeOpts{
		Name: "total_gauge",
		Help: "Total gauge",
		ConstLabels: prometheus.Labels{
			"app": config.App,
		},
	}

	if config.Name != "" {
		name := strings.ReplaceAll(config.Name, " ", "_")
		opt.Name = "total_" + name
		opt.Help = "Total " + name
	}

	m := &gauge{
		prometheusTotal: *prometheus.NewGaugeVec(
			opt,
			config.Tags,
		),
	}
	if config.MustRegister {
		prometheus.MustRegister(m.prometheusTotal)
	} else {
		_ = prometheus.Register(m.prometheusTotal)
	}
	return m
}

type gauge struct {
	prometheusTotal prometheus.GaugeVec
}

func (m *gauge) Set(tags []string, value float64) {
	m.prometheusTotal.WithLabelValues(tags...).Set(value)
}

func (m *gauge) Reset() {
	m.prometheusTotal.Reset()
}

func (m *gauge) Inc(tags []string) {
	m.prometheusTotal.WithLabelValues(tags...).Inc()
}
