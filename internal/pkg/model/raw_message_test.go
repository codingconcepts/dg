package model

import (
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestRawMessageUnmarshal(t *testing.T) {
	type test struct {
		R RawMessage `yaml:"r"`
	}

	y := `r: hello raw message`

	var tst test
	if err := yaml.NewDecoder(strings.NewReader(y)).Decode(&tst); err != nil {
		t.Fatalf("error decoding yaml: %v", err)
	}

	var s string
	if err := tst.R.UnmarshalFunc(&s); err != nil {
		log.Fatalf("error decoding yaml: %v", err)
	}

	assert.Equal(t, "hello raw message", s)
}
