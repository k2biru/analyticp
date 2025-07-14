package analyticp

import (
	"errors"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type TotalConsume interface {
	Inc(topic string, partition int, err error)
}

type TotalConsumeConfig struct {
	App          string
	Client       string
	ErrorLookup  ErrorLookup
	MustRegister bool
}

type totalConsume struct {
	prometheusTotal prometheus.CounterVec
	errorLookup     ErrorLookup
}

func NewTotalConsume(config *TotalConsumeConfig) TotalConsume {
	m := &totalConsume{
		errorLookup: config.ErrorLookup,
		prometheusTotal: *prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "total_consume",
				Help: "Total consume",
				ConstLabels: prometheus.Labels{
					"app":    config.App,
					"client": config.Client,
				},
			},
			[]string{"topic", "partition", "status"},
		),
	}
	if config.MustRegister {
		prometheus.MustRegister(m.prometheusTotal)
	} else {
		_ = prometheus.Register(m.prometheusTotal)
	}
	return m
}

func (m *totalConsume) Inc(topic string, partition int, err error) {
	status := "ok"
	if err != nil {
		status = "error"
		for v := range m.errorLookup {
			if errors.Is(err, v) {
				status = m.errorLookup[v]
			}
		}
	}
	m.prometheusTotal.WithLabelValues(topic, strconv.Itoa(partition), status).Inc()
}
