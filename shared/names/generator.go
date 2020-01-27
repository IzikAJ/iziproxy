package names

// Generator - generator interface
type Generator interface {
	Next() (string, error)
}

// GenerationError - generation error instance
type GenerationError struct {
	S string
}

func (e *GenerationError) Error() string {
	return e.S
}
