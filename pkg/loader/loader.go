package loader

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/ksahli/baal/pkg/monitor"
)

type Definition struct {
	Location  string `json:"location"`
	Frequency string `json:"frequency"`
	Method    string `josn:"method"`
}

type Loader struct {
	logger  *log.Logger
	decoder *json.Decoder

	closer io.Closer
}

func (l Loader) Load() (map[time.Duration][]monitor.Job, error) {
	definitions := []Definition{}
	defer l.closer.Close()
	if err := l.decoder.Decode(&definitions); err != nil {
		err := fmt.Errorf("loader error: %w", err)
		return nil, err
	}
	jobs := map[time.Duration][]monitor.Job{}
	for _, definition := range definitions {
		duration, err := time.ParseDuration(definition.Frequency)
		if err != nil {
			err := fmt.Errorf("loader error: %w", err)
			return nil, err
		}
		location, err := url.Parse(definition.Location)
		if err != nil {
			err := fmt.Errorf("loader error: %w", err)
			return nil, err
		}
		job := monitor.Job{
			Location: location,
			Method:   definition.Method,
		}
		jobs[duration] = append(jobs[duration], job)
	}
	return jobs, nil
}

func New(reader io.ReadCloser, logger *log.Logger) *Loader {
	decoder := json.NewDecoder(reader)
	loader := Loader{
		logger:  logger,
		decoder: decoder,
		closer:  reader,
	}
	return &loader
}

func File(path string, logger *log.Logger) (*Loader, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		err := fmt.Errorf("loader: %w", err)
		return nil, err
	}
	loader := New(file, logger)
	return loader, nil
}
