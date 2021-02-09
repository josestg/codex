package code

import (
	"github.com/josestg/codex/internal/platform/sandbox/internal/engine"
)

// GetBuilder returns engine builder based-on programming language.
func GetBuilder(language string) engine.Builder {
	switch language {
	default:
		return nil
	case "golang":
		return new(Golang)
	case "javascript":
		return new(JavaScript)
	}
}
