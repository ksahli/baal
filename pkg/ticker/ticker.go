package ticker

import (
	"context"
	"time"

	"github.com/ksahli/baal/pkg/monitor"
)

type Ticker struct {
	frequency time.Duration
	jobs      []monitor.Job
	jobsc     chan monitor.Job
}

func (t *Ticker) Tick(ctx context.Context) {
	defer close(t.jobsc)

	for {
		select {
		case <-time.Tick(t.frequency):
			for _, job := range t.jobs {
				t.jobsc <- job
			}
		case <-ctx.Done():
			return
		}
	}
}

func (t *Ticker) Jobsc() <-chan monitor.Job {
	return t.jobsc
}

func New(frequency time.Duration, jobs []monitor.Job) *Ticker {
	jobsc := make(chan monitor.Job, 100)
	ticker := Ticker{
		frequency: frequency,
		jobs:      jobs,
		jobsc:     jobsc,
	}
	return &ticker
}
