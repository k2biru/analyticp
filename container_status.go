package analyticp

import (
	"github.com/prometheus/client_golang/prometheus"
)

type ContainerStatus interface {
	Set(err error)
}

type StatusConfig struct {
	App          string
	MustRegister bool
}

type containerStatus struct {
	prometheusStatus prometheus.GaugeVec
}

func NewContainerStatus(config *StatusConfig) ContainerStatus {
	m := &containerStatus{
		prometheusStatus: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "container_status",
				Help: "Container status",
				ConstLabels: prometheus.Labels{
					"app": config.App,
				},
			},
			[]string{"status"},
		),
	}
	if config.MustRegister {
		prometheus.MustRegister(m.prometheusStatus)
	} else {
		_ = prometheus.Register(m.prometheusStatus)
	}
	return m
}

func (m *containerStatus) Set(err error) {
	if err == nil {
		m.prometheusStatus.WithLabelValues("ok").Set(1)
		m.prometheusStatus.WithLabelValues("error").Set(0)
		return
	}
	m.prometheusStatus.WithLabelValues("ok").Set(0)
	m.prometheusStatus.WithLabelValues("error").Set(1)
}
