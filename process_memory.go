package analyticp

import (
	"github.com/prometheus/client_golang/prometheus"
)

type ProcessMemory interface {
	Set(memory float64, name, act string)
}
type ProcessMemoryConfig struct {
	App          string
	Process      string
	MustRegister bool
}

type processMemory struct {
	gauge prometheus.GaugeVec
}

func NewProcessMemory(config *ProcessMemoryConfig) ProcessMemory {
	m := &processMemory{
		gauge: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "process_memory",
				Help: "percen of memory",
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

func (m *processMemory) Set(cpu float64, name, act string) {
	m.gauge.WithLabelValues(name, act).Set(cpu)
}
