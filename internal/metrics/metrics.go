package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Reporter interface {
	IncCreate(v uint, handler string)
	IncRemove(v uint, handler string)
	IncUpdate(v uint, handler string)
	IncRead(v uint, handler string)
	IncList(v uint, handler string)
}

type promReporter struct {
	createCounter *prometheus.CounterVec
	readCounter   *prometheus.CounterVec
	updateCounter *prometheus.CounterVec
	removeCounter *prometheus.CounterVec
	listCounter   *prometheus.CounterVec
}

func NewReporter() Reporter {
	return &promReporter{
		createCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "experiences_objects_create",
			Help: "The total number of create experiences",
		}, []string{"handler"}),
		readCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "experiences_objects_read",
			Help: "The total number of reads experiences",
		}, []string{"handler"}),
		updateCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "experiences_objects_update",
			Help: "The total number of update experiences",
		}, []string{"handler"}),
		removeCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "experiences_objects_remove",
			Help: "The total number of remove experiences",
		}, []string{"handler"}),
		listCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "experiences_objects_list",
			Help: "The total number of list experiences",
		}, []string{"handler"}),
	}
}

func (p *promReporter) IncCreate(v uint, handler string) {
	p.inc(p.createCounter, v, handler)
}

func (p *promReporter) IncRemove(v uint, handler string) {
	p.inc(p.removeCounter, v, handler)
}

func (p *promReporter) IncUpdate(v uint, handler string) {
	p.inc(p.updateCounter, v, handler)
}

func (p *promReporter) IncRead(v uint, handler string) {
	p.inc(p.readCounter, v, handler)
}

func (p *promReporter) IncList(v uint, handler string) {
	p.inc(p.listCounter, v, handler)
}

func (p *promReporter) inc(counter *prometheus.CounterVec, v uint, handler string) {
	counter.With(prometheus.Labels{"handler": handler}).Add(float64(v))
}
