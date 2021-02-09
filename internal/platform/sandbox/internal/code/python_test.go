package code_test

import (
	"testing"

	"github.com/pkg/errors"

	"github.com/josestg/codex/internal/platform/sandbox/internal/code"
	"github.com/josestg/codex/internal/platform/sandbox/internal/engine"
)

func TestPython_Build_AddFile(t *testing.T) {
	wm := &WorkdirManagerMock{
		methodFail: "AddFile",
	}
	builder := code.Python{}
	_, err := builder.Build(wm, []byte("print('Hello')"))
	if errors.Cause(err) != errMockAddFile {
		t.Errorf("Expected %v; but got %v", errMockAddFile, err)
	}
}

func TestPython_Build(t *testing.T) {
	wm := &WorkdirManagerMock{}
	builder := code.Python{}
	e, err := builder.Build(wm, []byte("print('Hello')"))
	if err != nil {
		t.Errorf("Expected %v; but got %v", nil, err)
	}

	if _, ok := e.(*engine.Engine); !ok {
		t.Errorf("Expected e == *engine.Engine")
	}
}
