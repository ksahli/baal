package monitor

import (
	"net/http"
	"net/url"
	"sync"
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
	lock    *sync.Mutex
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

func (m *Monitor) Run(wg *sync.WaitGroup, jobs <-chan Job) {
	defer wg.Done()
	for job := range jobs {
		result := m.Do(job)
		m.results <- result
	}
}

func (m *Monitor) Stop() {
	close(m.results)
}

func (m Monitor) Results() <-chan Result {
	return m.results
}

func New(client *http.Client, stamper func() time.Time) *Monitor {
	var (
		lock    = new(sync.Mutex)
		results = make(chan Result, 100)
	)
	monitor := Monitor{
		lock:    lock,
		stamper: stamper,
		client:  client,
		results: results,
	}
	return &monitor
}
