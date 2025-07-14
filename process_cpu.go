package analyticp

import (
	"github.com/prometheus/client_golang/prometheus"
)

type ProcessCPU interface {
	Set(cpu float64, name, act string)
}

type ProcessCPUConfig struct {
	App          string
	Process      string
	MustRegister bool
}

type processCPU struct {
	gauge prometheus.GaugeVec
}

func NewProcessCPU(config *ProcessCPUConfig) ProcessCPU {
	m := &processCPU{
		gauge: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "process_cpu",
				Help: "percen of CPU",
				ConstLabels: prometheus.Labels{
					"app":     config.App,
					"process": config.Process,
				},
			},
			[]string{"name", "act"},
		),
	}
	if config.MustRegister {
		prometheus.MustRegister(m.gauge)
	} else {
		_ = prometheus.Register(m.gauge)
	}
	return m
}

func (m *processCPU) Set(cpu float64, name, act string) {
	m.gauge.WithLabelValues(name, act).Set(cpu)
}
