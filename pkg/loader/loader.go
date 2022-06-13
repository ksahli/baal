package loader

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Definition struct {
	Location string `json:"location"`
}

type Loader struct {
	reader io.ReadCloser
}

func (l Loader) Load() ([]Definition, error) {
	defer l.reader.Close()
	decoder := json.NewDecoder(l.reader)
	definitions := []Definition{}
	if err := decoder.Decode(&definitions); err != nil {
		err := fmt.Errorf("loader error: %w", err)
		return nil, err
	}
	return definitions, nil
}

func New(reader io.ReadCloser) *Loader {
	loader := Loader{reader: reader}
	return &loader
}

func File(path string) (*Loader, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		err := fmt.Errorf("loader error: %w", err)
		return nil, err
	}
	loader := New(file)
	return loader, nil

}
