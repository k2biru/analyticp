package analyticp

import (
	"errors"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type TotalAccess interface {
	Inc(api string, err error)
	IncHTTPStatus(api string, status int)
}

type TotalAccessConfig struct {
	App          string
	Client       string
	ErrorLookup  ErrorLookup
	MustRegister bool
}

type totalAccess struct {
	prometheusTotal prometheus.CounterVec
	errorLookup     ErrorLookup
}

func NewTotalAccess(config *TotalAccessConfig) TotalAccess {
	m := &totalAccess{
		errorLookup: config.ErrorLookup,
		prometheusTotal: *prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "total_access",
				Help: "Total access",
				ConstLabels: prometheus.Labels{
					"app":    config.App,
					"client": config.Client,
				},
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

func (m *totalAccess) Inc(api string, err error) {
	status := "ok"
	if err != nil {
		status = "error"
		for v := range m.errorLookup {
			if errors.Is(err, v) {
				status = m.errorLookup[v]
			}
		}
	}
	m.prometheusTotal.WithLabelValues(api, status).Inc()
}

func (m *totalAccess) IncHTTPStatus(api string, status int) {
	statusStr := strconv.Itoa(status)
	m.prometheusTotal.WithLabelValues(api, statusStr).Inc()
}
