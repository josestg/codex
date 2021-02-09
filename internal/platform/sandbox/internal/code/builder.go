package code

import (
	"errors"

	"github.com/josestg/codex/internal/platform/sandbox/internal/engine"
)

var (
	ErrUnknownProgrammingLanguage = errors.New("UNKNOWN PROGRAMMING LANGUAGE")
)

// GetBuilder returns engine builder based-on programming language.
func GetBuilder(language string) (engine.Builder, error) {
	switch language {
	default:
		return nil, ErrUnknownProgrammingLanguage
	case "golang":
		return new(Golang), nil
	case "javascript":
		return new(JavaScript), nil
	}
}
