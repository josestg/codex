package code_test

import (
	"fmt"
	"testing"

	"github.com/josestg/codex/internal/platform/sandbox/internal/code"
	"github.com/josestg/codex/internal/platform/sandbox/internal/engine"
)

func TestGetBuilder(t *testing.T) {
	type funcEval func(b engine.Builder) bool
	tc := []struct {
		name string
		eval funcEval
	}{
		{
			name: "golang",
			eval: func(b engine.Builder) bool {
				_, ok := b.(*code.Golang)
				return ok
			},
		},
		{
			name: "javascript",
			eval: func(b engine.Builder) bool {
				_, ok := b.(*code.JavaScript)
				return ok
			},
		}, {
			name: "unknown",
			eval: func(b engine.Builder) bool {
				return b == nil
			},
		},
	}

	for _, c := range tc {
		name := fmt.Sprintf("Get %s builder", c.name)
		t.Run(name, func(tt *testing.T) {
			builder := code.GetBuilder(c.name)
			if !c.eval(builder) {
				t.Errorf("Expected %s builder but got %v", c.name, builder)
			}
		})
	}

}
