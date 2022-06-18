package collector_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/ksahli/baal/pkg/collector"
	"github.com/ksahli/baal/pkg/monitor"
)

type Writer struct {
	cfail, wfail bool
	closed       bool
	results      []monitor.Result
}

func (w *Writer) Write(p []byte) (int, error) {
	if w.wfail {
		err := errors.New("writer failure")
		return 0, err
	}
	b, result := bytes.NewBuffer(p), monitor.Result{}
	decoder := json.NewDecoder(b)
	if err := decoder.Decode(&result); err != nil {
		err := fmt.Errorf("writer error: %w", err)
		return 0, err
	}
	w.results = append(w.results, result)
	return len(p), nil
}

func (w *Writer) Close() error {
	if w.cfail {
		err := errors.New("closer error")
		w.closed = false
		return err
	}
	w.closed = true
	return nil
}

type Out struct {
	messages []string
}

func (w *Out) Write(p []byte) (int, error) {
	message := string(p)
	w.messages = append(w.messages, message)
	return len(p), nil
}

func TestWrite(t *testing.T) {
	writer := Writer{
		wfail: false,
		cfail: false,
	}

	logger := log.New(os.Stderr, "collector", log.Ldate)

	sut := collector.New(&writer, logger)
	defer sut.Close()

	location, err := url.Parse("https://localhost")
	if err != nil {
		msg := "unwanted error %w"
		t.Fatalf(msg, err)
	}

	want := monitor.Result{
		Location:  location,
		Reachable: false,
		Status:    200,
		Time:      time.Now().UTC(),
	}
	if err := sut.Write(want); err != nil {
		msg := "unwanted error %w"
		t.Fatalf(msg, err)
	}

	got := writer.results[0]
	if !reflect.DeepEqual(want, got) {
		msg := "want %v, got %v"
		t.Fatalf(msg, want, got)
	}
}

func TestWriteError(t *testing.T) {
	writer := Writer{
		wfail: true,
		cfail: false,
	}

	logger := log.New(os.Stderr, "collector", log.Ldate)

	sut := collector.New(&writer, logger)
	defer sut.Close()

	location, err := url.Parse("https://localhost")
	if err != nil {
		msg := "unwanted error %w"
		t.Fatalf(msg, err)
	}

	want := monitor.Result{
		Location:  location,
		Reachable: false,
		Status:    200,
		Time:      time.Now().UTC(),
	}
	if err := sut.Write(want); err == nil {
		t.Fatal("want an error, got nothing")
	}
}

func TestRun(t *testing.T) {
	writer := Writer{
		wfail: false,
		cfail: false,
	}

	out := new(Out)
	logger := log.New(out, "collector", log.Ldate)

	sut := collector.New(&writer, logger)
	defer sut.Close()

	wg := new(sync.WaitGroup)
	results := make(chan monitor.Result, 10)

	wg.Add(1)
	go sut.Run(wg, results)

	want := []monitor.Result{}
	for i := 0; i < 10; i++ {
		location, err := url.Parse("https://localhost")
		if err != nil {
			msg := "unwanted error %w"
			t.Fatalf(msg, err)
		}
		result := monitor.Result{
			Location:  location,
			Reachable: false,
			Status:    200,
			Time:      time.Now().UTC(),
		}
		want = append(want, result)
		results <- result
	}
	close(results)
	wg.Wait()

	got := writer.results
	if !reflect.DeepEqual(want, got) {
		msg := "want %v, got %v"
		t.Fatalf(msg, want, got)
	}
}

func TestRunWriteError(t *testing.T) {
	writer := Writer{
		wfail: true,
		cfail: false,
	}

	out := new(Out)
	logger := log.New(out, " [collector] ", log.Ldate)

	sut := collector.New(&writer, logger)
	defer sut.Close()

	wg := new(sync.WaitGroup)
	results := make(chan monitor.Result, 10)

	wg.Add(1)
	go sut.Run(wg, results)

	for i := 0; i < 10; i++ {
		location, err := url.Parse("https://localhost")
		if err != nil {
			msg := "unwanted error %w"
			t.Fatalf(msg, err)
		}
		result := monitor.Result{
			Location:  location,
			Reachable: false,
			Status:    200,
			Time:      time.Now().UTC(),
		}
		results <- result
	}
	close(results)
	wg.Wait()

	got := writer.results
	if len(got) != 0 {
		msg := "want nothing, got %d"
		t.Fatalf(msg, len(got))
	}

	if len(out.messages) != 10 {
		msg := "want 10 messages, got %d"
		t.Fatalf(msg, len(out.messages))
	}
}

func TestRunCloseError(t *testing.T) {
	writer := Writer{
		wfail: false,
		cfail: true,
	}

	out := new(Out)
	logger := log.New(out, " [collector] ", log.Ldate)

	sut := collector.New(&writer, logger)
	sut.Close()

	if len(out.messages) != 1 {
		msg := "want 1 messages, got %d"
		t.Fatalf(msg, len(out.messages))
	}
}

func TestFile(t *testing.T) {
	logger := log.New(os.Stderr, " [collector] ", log.Ldate)
	directory := t.TempDir()
	path := fmt.Sprintf("%s/definitions.json", directory)
	if _, err := os.Create(path); err != nil {
		msg := "unwanted error %v"
		t.Fatalf(msg, err)
	}
	collector, err := collector.File(path, logger)
	if err != nil {
		msg := "unwanted error %v"
		t.Fatalf(msg, err)
	}
	if collector == nil {
		t.Fatal("want a collector, got nothing")
	}
}

func TestFileError(t *testing.T) {
	logger := log.New(os.Stderr, " [collector] ", log.Ldate)
	collector, err := collector.File("/invalid path", logger)
	if err == nil {
		t.Fatal("want an error, got nothing")
	}
	if collector != nil {
		msg := "want nothing, got %v"
		t.Fatalf(msg, collector)

	}
}
