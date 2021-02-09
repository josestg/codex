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
		name    string
		eval    funcEval
		wantErr bool
	}{
		{
			name:    "golang",
			wantErr: false,
			eval: func(b engine.Builder) bool {
				_, ok := b.(*code.Golang)
				return ok
			},
		},
		{
			name:    "javascript",
			wantErr: false,
			eval: func(b engine.Builder) bool {
				_, ok := b.(*code.JavaScript)
				return ok
			},
		}, {
			name:    "unknown",
			wantErr: true,
			eval: func(b engine.Builder) bool {
				return true
			},
		},
	}

	for _, c := range tc {
		name := fmt.Sprintf("Get %s builder", c.name)
		t.Run(name, func(tt *testing.T) {
			builder, err := code.GetBuilder(c.name)
			if c.wantErr {
				if err != code.ErrUnknownProgrammingLanguage {
					t.Errorf("Expected error %v nil; got %v", code.ErrUnknownProgrammingLanguage, err)
				}
				if builder != nil {
					t.Errorf("Expected builder nil; got %v", builder)
				}
			} else {
				if err != nil {
					t.Errorf("Expected error nil; got %v", err)
				}

				if !c.eval(builder) {
					t.Errorf("Expected %s builder but got %v", c.name, builder)
				}
			}
		})
	}

}
