package loader_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/ksahli/baal/pkg/loader"
)

func TestLoadSuccess(t *testing.T) {
	want := []loader.Definition{
		{Location: "domain.1"},
		{Location: "domain.2"},
		{Location: "domain.3"},
	}
	reader := Reader{definitions: want}
	sut := loader.New(reader)
	got, err := sut.Load()
	if err != nil {
		msg := "unwanted error: %v"
		t.Fatalf(msg, err)
	}
	if !reflect.DeepEqual(want, got) {
		msg := "want %v, got %v"
		t.Fatalf(msg, want, got)
	}
}

func TestLoadError(t *testing.T) {
	reader := Reader{fail: true}
	sut := loader.New(reader)
	got, err := sut.Load()
	if err == nil {
		t.Fatal("want an error got nothing")
	}
	if len(got) != 0 {
		msg := "want no defintions, got %d"
		t.Fatalf(msg, len(got))
	}
}

func TestFile(t *testing.T) {
	directory := t.TempDir()
	path := fmt.Sprintf("%s/domains.json", directory)
	file, err := os.Create(path)
	if err != nil {
		msg := "unwanted error: %v"
		t.Fatalf(msg, err)
	}
	defer file.Close()
	sut, err := loader.File(path)
	if err != nil {
		msg := `unwanted error for path: %s`
		t.Fatalf(msg, path)
	}
	if sut == nil {
		t.Fatal("want a loader, got nothing")
	}
}

func TestFileError(t *testing.T) {
	sut, err := loader.File("")
	if err == nil {
		t.Fatalf("want an error, got nothing")
	}
	if sut != nil {
		msg := "want no loader, %v"
		t.Fatalf(msg, sut)
	}
}
