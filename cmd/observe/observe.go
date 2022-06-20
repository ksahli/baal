package observe

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/ksahli/baal/pkg/collector"
	"github.com/ksahli/baal/pkg/loader"
	"github.com/ksahli/baal/pkg/monitor"
	"github.com/ksahli/baal/pkg/ticker"
)

type Command struct {
	Definitions string
	Results     string
}

func (c Command) Execute(ctx context.Context) error {
	logger := log.New(os.Stderr, " [baal] ", log.Ldate)

	loader, err := loader.File(c.Definitions, logger)
	if err != nil {
		err := fmt.Errorf("observe: %w", err)
		return err
	}

	client := new(http.Client)
	stamper := time.Now
	monitor := monitor.New(client, stamper)

	collector, err := collector.File(c.Results, logger)
	if err != nil {
		err := fmt.Errorf("observe failed: %w", err)
		return err
	}

	cwg, mwg := new(sync.WaitGroup), new(sync.WaitGroup)

	cwg.Add(1)
	go collector.Run(cwg, monitor.Results())

	frequencies, err := loader.Load()
	if err != nil {
		err := fmt.Errorf("observe: %w", err)
		return err
	}

	for frequency, jobs := range frequencies {
		ticker := ticker.New(frequency, jobs)

		for i := 0; i <= 10; i++ {
			mwg.Add(1)
			go monitor.Run(mwg, ticker.Jobsc())
		}

		go ticker.Tick(ctx)
	}

	mwg.Wait()
	monitor.Stop()

	cwg.Wait()
	collector.Stop()

	return nil
}
