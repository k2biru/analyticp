package analyticp

import (
	"errors"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type TotalDurration interface {
	Inc(api string, durr time.Duration, err error)
	IncHTTPStatus(api string, durr time.Duration, status int)
}

type TotalDurrationConfig struct {
	App          string
	Client       string
	ErrorLookup  ErrorLookup
	MustRegister bool
	Bucket       []float64
}

type totalDurration struct {
	prometheusTotal prometheus.HistogramVec
	errorLookup     ErrorLookup
}

func NewTotalDurration(config *TotalDurrationConfig) TotalDurration {
	if len(config.Bucket) == 0 {
		config.Bucket = []float64{
			0.01, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0,
		}
	}
	m := &totalDurration{
		errorLookup: config.ErrorLookup,
		prometheusTotal: *prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "total_durration",
				Help: "Total durration second",
				ConstLabels: prometheus.Labels{
					"app":    config.App,
					"client": config.Client,
				},
				Buckets: config.Bucket,
			},
			[]string{"api", "status"},
		),
	}
	if config.MustRegister {
		prometheus.MustRegister(m.prometheusTotal)
	} else {
		_ = prometheus.Register(m.prometheusTotal)
	}
	return m
}

func (m *totalDurration) Inc(api string, durr time.Duration, err error) {
	status := "ok"
	if err != nil {
		status = "error"
		for v := range m.errorLookup {
			if errors.Is(err, v) {
				status = m.errorLookup[v]
			}
		}
	}
	m.prometheusTotal.WithLabelValues(api, status).Observe(durr.Seconds())
}

func (m *totalDurration) IncHTTPStatus(api string, durr time.Duration, status int) {
	statusStr := strconv.Itoa(status)
	m.prometheusTotal.WithLabelValues(api, statusStr).Observe(durr.Seconds())
}
