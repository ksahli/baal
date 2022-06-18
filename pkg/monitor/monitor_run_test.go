package monitor_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"sync"
	"testing"

	"github.com/ksahli/baal/pkg/monitor"
)

func TestRun(t *testing.T) {

	client := new(http.Client)
	sut := monitor.New(client, stamper)
	defer sut.Stop()

	jobs := make(chan monitor.Job, 10)
	wg := new(sync.WaitGroup)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go sut.Run(wg, jobs)
	}

	got := []monitor.Result{}
	go func() {
		for result := range sut.Results() {
			got = append(got, result)
		}
	}()

	want := map[*url.URL]monitor.Result{}
	for i := 0; i < 10; i++ {
		handle := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
		}
		handler := http.HandlerFunc(handle)
		server := httptest.NewServer(handler)
		defer server.Close()

		location, err := url.Parse(server.URL)
		if err != nil {
			msg := "unwanted error %v"
			t.Fatalf(msg, err)
		}

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

		want[location] = result
		jobs <- job
	}

	close(jobs)
	wg.Wait()

	for _, g := range got {
		w, ok := want[g.Location]
		if !ok {
			t.Fatalf("no entry was found for %s", g.Location)
		}
		if !reflect.DeepEqual(w, g) {
			msg := "\n want: %v \n got:  %v"
			t.Fatalf(msg, w, g)
		}
	}
}
