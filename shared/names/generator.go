package names

// Generator - generator interface
type Generator interface {
	Next() (string, error)
}

// GenerationError - generation error instance
type GenerationError struct {
	message string
}

func (e *GenerationError) Error() string {
	return e.message
}

// MissError - shoud be thrown on too many misses
var MissError = &GenerationError{"to many miss times, increase max size"}

// NewGenerationError - generation error generator
func NewGenerationError(message string) error {
	return &GenerationError{message}
}
