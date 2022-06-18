package loader_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ksahli/baal/pkg/loader"
	"github.com/ksahli/baal/pkg/monitor"
)

type Reader struct {
	fail        bool
	definitions []loader.Definition
}

func (r Reader) Read(p []byte) (int, error) {
	if r.fail {
		err := errors.New(`test reader error`)
		return 0, err
	}
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	if err := encoder.Encode(&r.definitions); err != nil {
		err := fmt.Errorf(`test reader error: %w`, err)
		return 0, err
	}
	return buffer.Read(p)
}

func (r Reader) Close() error {
	return nil
}

func location(t *testing.T, l string) *url.URL {
	URL, err := url.Parse(l)
	if err != nil {
		t.Fatalf("unwanted error %v", err)
	}
	return URL
}

func TestLoad(t *testing.T) {
	reader := Reader{
		definitions: []loader.Definition{
			{
				Location:  "http:/domain-1.com",
				Method:    "GET",
				Frequency: "5m",
			},
			{
				Location:  "http:/domain-2.com",
				Method:    "GET",
				Frequency: "5m",
			},
			{
				Location:  "http:/domain-3.com",
				Method:    "GET",
				Frequency: "15m",
			},
			{
				Location:  "http:/domain-4.com",
				Method:    "GET",
				Frequency: "15m",
			},
			{
				Location:  "http:/domain-5.com",
				Method:    "GET",
				Frequency: "1h",
			},
			{
				Location:  "http:/domain-6.com",
				Method:    "GET",
				Frequency: "1h",
			},
		},
	}
	logger := log.New(os.Stderr, " [loader] ", log.Ldate)
	sut := loader.New(&reader, logger)
	got, err := sut.Load()
	if err != nil {
		msg := "unwanted error: %v"
		t.Fatalf(msg, err)
	}
	want := map[time.Duration][]monitor.Job{
		5 * time.Minute: []monitor.Job{
			{
				Location: location(t, "http:/domain-1.com"),
				Method:   "GET",
			},
			{
				Location: location(t, "http:/domain-2.com"),
				Method:   "GET",
			},
		},
		15 * time.Minute: []monitor.Job{
			{
				Location: location(t, "http:/domain-3.com"),
				Method:   "GET",
			},
			{
				Location: location(t, "http:/domain-4.com"),
				Method:   "GET",
			},
		},
		time.Hour: []monitor.Job{
			{
				Location: location(t, "http:/domain-5.com"),
				Method:   "GET",
			},
			{
				Location: location(t, "http:/domain-6.com"),
				Method:   "GET",
			},
		},
	}
	if !reflect.DeepEqual(want, got) {
		msg := "\n want %v\n got  %v"
		t.Fatalf(msg, want, got)
	}
}

func TestLoadDecoderError(t *testing.T) {
	reader := Reader{fail: true}
	logger := log.New(os.Stderr, " [loader] ", log.Ldate)
	sut := loader.New(&reader, logger)
	got, err := sut.Load()
	if err == nil {
		t.Fatal("want an error, got nothing")
	}
	if len(got) != 0 {
		msg := "want no definitions, got %d"
		t.Fatalf(msg, len(got))
	}
}

func TestLoadParseLocationError(t *testing.T) {
	reader := Reader{
		definitions: []loader.Definition{
			{
				Location:  "://invalid domain",
				Method:    "GET",
				Frequency: "5m",
			},
		},
	}
	logger := log.New(os.Stderr, " [loader] ", log.Ldate)
	sut := loader.New(&reader, logger)
	got, err := sut.Load()
	if err == nil {
		t.Fatal("want an error, got nothing")
	}
	if len(got) != 0 {
		msg := "want no definitions, got %d"
		t.Fatalf(msg, len(got))
	}
}

func TestLoadParseDurationError(t *testing.T) {
	reader := Reader{
		definitions: []loader.Definition{
			{
				Location:  "http:/domain-1.com",
				Method:    "GET",
				Frequency: "invalid duration",
			},
		},
	}
	logger := log.New(os.Stderr, " [loader] ", log.Ldate)
	sut := loader.New(&reader, logger)
	got, err := sut.Load()
	if err == nil {
		t.Fatal("want an error, got nothing")
	}
	if len(got) != 0 {
		msg := "want no definitions, got %d"
		t.Fatalf(msg, len(got))
	}
}

func TestFile(t *testing.T) {
	directory := t.TempDir()
	path := fmt.Sprintf("%s/definitions.json", directory)
	if _, err := os.Create(path); err != nil {
		msg := "unwanted error %v"
		t.Fatalf(msg, err)
	}
	logger := log.New(os.Stderr, " [loader] ", log.Ldate)
	loader, err := loader.File(path, logger)
	if err != nil {
		msg := "unwanted error %v"
		t.Fatalf(msg, err)
	}
	if loader == nil {
		t.Fatal("want a loader, got nothing")
	}
}

func TestFileError(t *testing.T) {
	logger := log.New(os.Stderr, " [loader] ", log.Ldate)
	loader, err := loader.File("invalid path", logger)
	if err == nil {
		t.Fatal("want an error, got nothing")
	}
	if loader != nil {
		msg := "want nothing, got %v"
		t.Fatalf(msg, loader)
	}
}
