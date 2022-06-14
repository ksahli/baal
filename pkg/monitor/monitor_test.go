package monitor_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/ksahli/baal/pkg/monitor"
)

var timestamp = time.Now()

func stamper() time.Time {
	return timestamp
}

func TestDo(t *testing.T) {
	handle := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}
	handler := http.HandlerFunc(handle)
	server := httptest.NewServer(handler)
	defer server.Close()
	client := new(http.Client)
	sut := monitor.New(client, stamper)
	location, err := url.Parse(server.URL)
	if err != nil {
		msg := "unwanted error %v"
		t.Fatalf(msg, err)
	}
	job := monitor.Job{
		Location: location,
		Method:   "GET",
	}
	got := sut.Do(job)
	want := monitor.Result{
		Location:  location,
		Status:    200,
		Reachable: true,
		Time:      timestamp,
	}
	if !reflect.DeepEqual(want, got) {
		msg := "want %v, got %v"
		t.Fatalf(msg, want, got)
	}
}

func TestDoError(t *testing.T) {
	client := new(http.Client)
	sut := monitor.New(client, stamper)
	location, err := url.Parse("0.0.0.0")
	if err != nil {
		msg := "unwanted error %v"
		t.Fatalf(msg, err)
	}
	job := monitor.Job{
		Location: location,
		Method:   "GET",
	}
	got := sut.Do(job)
	want := monitor.Result{
		Location:  location,
		Status:    0,
		Reachable: false,
		Time:      timestamp,
	}
	if !reflect.DeepEqual(want, got) {
		msg := "want %v, got %v"
		t.Fatalf(msg, want, got)
	}
}

func TestRun(t *testing.T) {
	handle := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}
	handler := http.HandlerFunc(handle)
	server := httptest.NewServer(handler)
	defer server.Close()
	client := new(http.Client)
	sut := monitor.New(client, stamper)
	location, err := url.Parse(server.URL)
	if err != nil {
		msg := "unwanted error %v"
		t.Fatalf(msg, err)
	}
	want := []monitor.Result{}
	jobs := make(chan monitor.Job, 10)
	go sut.Run(jobs)
	for i := 0; i < 10; i++ {
		job := monitor.Job{
			Location: location,
			Method:   "GET",
		}
		result := monitor.Result{
			Location:  location,
			Status:    200,
			Reachable: true,
			Time:      timestamp,
		}
		want = append(want, result)
		jobs <- job
	}
	close(jobs)
	got := []monitor.Result{}
	for result := range sut.Results() {
		got = append(got, result)
	}
	if !reflect.DeepEqual(want, got) {
		msg := "want %v, got %v"
		t.Fatalf(msg, want, got)
	}
}
