package collector

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/ksahli/baal/pkg/monitor"
)

type Results = <-chan monitor.Result

type Collector struct {
	wlock, clock *sync.Mutex
	encoder      *json.Encoder
	logger       *log.Logger

	closer io.Closer
}

func (c *Collector) Write(result monitor.Result) error {
	c.wlock.Lock()
	defer c.wlock.Unlock()

	if err := c.encoder.Encode(&result); err != nil {
		return err
	}
	return nil
}

func (c *Collector) Run(wg *sync.WaitGroup, results Results) {
	defer c.Close()
	defer wg.Done()

	for result := range results {
		if err := c.Write(result); err != nil {
			c.logger.Print(err)
		}
	}
}

func (c *Collector) Close() {
	c.clock.Lock()
	defer c.clock.Unlock()

	if err := c.closer.Close(); err != nil {
		c.logger.Print(err)
	}
}

func New(writer io.WriteCloser, logger *log.Logger) *Collector {
	var (
		wlock, clock = new(sync.Mutex), new(sync.Mutex)
		encoder      = json.NewEncoder(writer)
	)
	collector := Collector{
		wlock:   wlock,
		clock:   clock,
		logger:  logger,
		encoder: encoder,
		closer:  writer,
	}
	return &collector
}

func File(path string, logger *log.Logger) (*Collector, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		err := fmt.Errorf("collector error: %w", err)
		return nil, err
	}
	return New(file, logger), nil
}
