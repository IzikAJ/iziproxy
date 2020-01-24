package names

// Generator - generator interface
type Generator interface {
	Next() (string, error)
}

// GenerationError - generation error instance
type GenerationError struct {
	s string
}

func (e *GenerationError) Error() string {
	return e.s
}
