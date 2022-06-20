package ticker_test

import (
	"context"
	"fmt"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/ksahli/baal/pkg/monitor"
	"github.com/ksahli/baal/pkg/ticker"
)

var ctx = context.Background()

func TestTick(t *testing.T) {
	want := make([]monitor.Job, 0, 10)
	for i := 0; i < 10; i++ {
		f := fmt.Sprintf("https://domain-%d.com", i)
		location, err := url.Parse(f)
		if err != nil {
			msg := "unwanted error: %v"
			t.Fatalf(msg, err)
		}
		job := monitor.Job{
			Location: location,
			Method:   "GET",
		}
		want = append(want, job)
	}
	ticker := ticker.New(2*time.Second, want)

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	go ticker.Tick(ctx)

	got := make([]monitor.Job, 0, 10)
	for job := range ticker.Jobsc() {
		got = append(got, job)
	}

	if !reflect.DeepEqual(want, got) {
		msg := "\nwant %v\n  got %v"
		t.Fatalf(msg, want, got)
	}
}
