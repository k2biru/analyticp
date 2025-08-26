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
	App          string
	Tags         []string
	MustRegister bool
}

func NewTotalOnline(config *TotalOnlineConfig) TotalOnline {
	m := &gauge{
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
	if config.MustRegister {
		prometheus.MustRegister(m.prometheusTotal)
	} else {
		_ = prometheus.Register(m.prometheusTotal)
	}
	return m
}
