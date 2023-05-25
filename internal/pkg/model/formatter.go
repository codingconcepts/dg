package model

// Formatter determines the behaviour for anything that can take a format
// string and return another.
type Formatter interface {
	Format(string) string
}

// FormatterProcessor can be called to get the Format string out of a struct
// that implements this interface.
type FormatterProcessor interface {
	GetFormat() string
}
