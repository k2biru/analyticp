package analyticp

import (
	"github.com/prometheus/client_golang/prometheus"
)

type TotalOnline interface {
	Set(tags []string, value float64)
	Reset()
	Inc(tags []string)
}

type TotalOnlineConfig struct {
	App  string
	Tags []string
}

type totalOnline struct {
	prometheusTotal prometheus.GaugeVec
}

func NewTotalOnline(config *TotalOnlineConfig) TotalOnline {
	m := &totalOnline{
		prometheusTotal: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "total_online",
				Help: "Total online",
				ConstLabels: prometheus.Labels{
					"app": config.App,
				},
			},
			config.Tags,
		),
	}
	_ = prometheus.Register(m.prometheusTotal)
	return m
}

func (m *totalOnline) Set(tags []string, value float64) {
	m.prometheusTotal.WithLabelValues(tags...).Set(value)
}

func (m *totalOnline) Reset() {
	m.prometheusTotal.Reset()
}

func (m *totalOnline) Inc(tags []string) {
	m.prometheusTotal.WithLabelValues(tags...).Inc()
}
