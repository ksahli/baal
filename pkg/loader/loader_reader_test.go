package loader_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ksahli/baal/pkg/loader"
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
