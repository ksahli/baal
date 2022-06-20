package observe_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/ksahli/baal/cmd/observe"
)

var ctx = context.Background()

func TestExecute(t *testing.T) {
	directory := t.TempDir()
	path := fmt.Sprintf("%s/results.json", directory)
	cmd := observe.Command{
		Definitions: "testdata/definitions.json",
		Results:     path,
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer func() {
		cancel()
		file, err := os.OpenFile(path, os.O_RDONLY, 0644)
		if err != nil {
			msg := "unwanted error: %w"
			t.Fatalf(msg, err)
		}

		content, err := ioutil.ReadAll(file)
		if err != nil {
			msg := "unwanted error: %w"
			t.Fatalf(msg, err)
		}

		if content == nil {
			t.Fatalf("want content written to results file, got nothing")
		}
	}()

	if err := cmd.Execute(ctx); err != nil {
		msg := "unwanted error: %w"
		t.Fatalf(msg, err)
	}

}

func TestExecuteInvalidDefinitionsPath(t *testing.T) {
	directory := t.TempDir()
	path := fmt.Sprintf("%s/results.json", directory)
	cmd := observe.Command{
		Definitions: "/invalid_path",
		Results:     path,
	}

	if err := cmd.Execute(ctx); err == nil {
		t.Fatal("want an error, got nothing")
	}
}

func TestExecuteInvalidDefinitions(t *testing.T) {
	directory := t.TempDir()
	path := fmt.Sprintf("%s/results.json", directory)
	cmd := observe.Command{
		Definitions: "testdata/invalid.json",
		Results:     path,
	}

	if err := cmd.Execute(ctx); err == nil {
		t.Fatal("want an error, got nothing")
	}
}

func TestExecuteInvaidResultsPath(t *testing.T) {
	cmd := observe.Command{
		Definitions: "testdata/definitions.json",
		Results:     "/invalid",
	}

	if err := cmd.Execute(ctx); err == nil {
		t.Fatal("want an error, got nothing")
	}
}
