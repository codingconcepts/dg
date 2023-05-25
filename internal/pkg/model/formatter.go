package model

// Formatter determines the behaviour for anything that can take a format
// string and return another.
type Formatter interface {
	Format(string) string
}
