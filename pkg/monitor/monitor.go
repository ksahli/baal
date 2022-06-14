package monitor

import (
	"net/http"
	"net/url"
	"time"
)

type Job struct {
	Location *url.URL
	Method   string
}

type Result struct {
	Location  *url.URL
	Status    int
	Reachable bool
	Time      time.Time
}

type Monitor struct {
	stamper func() time.Time
	client  *http.Client
	results chan Result
}

func (m *Monitor) Do(job Job) Result {
	request := http.Request{
		URL:    job.Location,
		Method: job.Method,
	}
	response, err := m.client.Do(&request)
	result := Result{
		Location: job.Location,
		Time:     m.stamper(),
	}
	if err == nil {
		result.Reachable = true
		result.Status = response.StatusCode
	}
	return result
}

func (m *Monitor) Run(jobs <-chan Job) {
	defer close(m.results)
	for job := range jobs {
		result := m.Do(job)
		m.results <- result
	}
}

func (m Monitor) Results() <-chan Result {
	return m.results
}

func New(client *http.Client, stamper func() time.Time) *Monitor {
	results := make(chan Result, 100)
	monitor := Monitor{
		stamper: stamper,
		client:  client,
		results: results,
	}
	return &monitor
}
