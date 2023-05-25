package model

// RawMessage does what json.RawMessage does but for YAML.
type RawMessage struct {
	UnmarshalFunc func(interface{}) error
}

func (msg *RawMessage) UnmarshalYAML(unmarshal func(interface{}) error) error {
	msg.UnmarshalFunc = unmarshal
	return nil
}

func (msg *RawMessage) Unmarshal(v any) error {
	return msg.UnmarshalFunc(v)
}
