package main

import (
	"context"
	"flag"
	"os"

	"github.com/ksahli/baal/cmd/observe"
)

type Command interface {
	Execute(ctx context.Context) error
}

var ctx = context.Background()

func main() {
	if err := run(ctx); err != nil {
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	var command Command

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	switch os.Args[1] {
	case "observe":
		flags := flag.NewFlagSet("observe", flag.ExitOnError)
		var (
			definitions = flags.String("definitions", "", "domains definitions file")
			results     = flags.String("results", "", "monitoring results file")
		)
		if err := flags.Parse(os.Args[2:]); err != nil {
			return err
		}
		command = observe.Command{
			Definitions: *definitions,
			Results:     *results,
		}
	}

	if err := command.Execute(ctx); err != nil {
		return err
	}

	return nil
}
