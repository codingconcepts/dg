package model

import (
	"bytes"
	"testing"

	"gopkg.in/yaml.v3"
)

// RawMessage does what json.RawMessage does but for YAML.
type RawMessage struct {
	UnmarshalFunc func(interface{}) error
}

func (msg *RawMessage) UnmarshalYAML(unmarshal func(interface{}) error) error {
	msg.UnmarshalFunc = unmarshal
	return nil
}

// ToRawMessage converts an object into a model.RawMessage for testing purposes.
func ToRawMessage(t *testing.T, v any) RawMessage {
	buf := &bytes.Buffer{}
	if err := yaml.NewEncoder(buf).Encode(v); err != nil {
		t.Fatalf("error encoding to yaml: %v", err)
	}

	var rawMessage RawMessage
	if err := yaml.NewDecoder(buf).Decode(&rawMessage); err != nil {
		t.Fatalf("error decoding from yaml: %v", err)
	}

	return rawMessage
}
