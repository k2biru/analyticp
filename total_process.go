package analyticp

import (
	"errors"

	"github.com/prometheus/client_golang/prometheus"
)

type TotalProcess interface {
	Inc(subProcess string, err error)
}

type TotalProcessConfig struct {
	App          string
	Process      string
	ErrorLookup  ErrorLookup
	MustRegister bool
}

type totalProcess struct {
	prometheusTotal prometheus.CounterVec
	errorLookup     ErrorLookup
}

func NewTotalProcess(config *TotalProcessConfig) TotalProcess {
	m := &totalProcess{
		errorLookup: config.ErrorLookup,
		prometheusTotal: *prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "total_process",
				Help: "Total process",
				ConstLabels: prometheus.Labels{
					"app":     config.App,
					"process": config.Process,
				},
			},
			[]string{"sub_process", "status"},
		),
	}
	if config.MustRegister {
		prometheus.MustRegister(m.prometheusTotal)
	} else {
		_ = prometheus.Register(m.prometheusTotal)
	}
	return m
}

func (m *totalProcess) Inc(subProcess string, err error) {
	status := "ok"
	if err != nil {
		status = "error"
		for v := range m.errorLookup {
			if errors.Is(err, v) {
				status = m.errorLookup[v]
			}
		}
	}
	m.prometheusTotal.WithLabelValues(subProcess, status).Inc()
}
