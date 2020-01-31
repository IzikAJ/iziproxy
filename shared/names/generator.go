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

// names.GenerationError{S: "no fallback, sorry"}

// MissError - shoud be thrown on too many misses
var MissError = &GenerationError{"to many miss times, increase max size"}

// NewGenerationError - generation error generator
func NewGenerationError(message string) error {
	return &GenerationError{S: message}
}
