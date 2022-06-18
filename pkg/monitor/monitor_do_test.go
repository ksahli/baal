package monitor_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/ksahli/baal/pkg/monitor"
)

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
