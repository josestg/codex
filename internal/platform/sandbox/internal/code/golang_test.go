package code_test

import (
	"os"
	"testing"

	"github.com/pkg/errors"

	"github.com/josestg/codex/internal/platform/sandbox/internal/code"
	"github.com/josestg/codex/internal/platform/sandbox/internal/engine"
	"github.com/josestg/codex/internal/platform/sandbox/internal/workspace"
)

type WorkdirManagerMock struct {
	methodFail string
}

func (w WorkdirManagerMock) Location() string {
	if w.methodFail == "Location" {
		return "invalid_location"
	}
	return os.TempDir()
}

var (
	errMockAddFile = errors.New("MOCK ADD FILE")
	errMockChdir   = errors.New("MOCK CHDIR")
)

func (w WorkdirManagerMock) AddFile(name string, content []byte) (string, error) {
	if w.methodFail == "AddFile" {
		return "", errMockAddFile
	}
	return "", nil
}

func (w WorkdirManagerMock) Chdir() error {
	if w.methodFail == "Chdir" {
		return errMockChdir
	}
	return nil
}

const sourceCode = `
package main

func main() {
	
}
`

func TestGolang_Build_AddFile(t *testing.T) {
	wm := &WorkdirManagerMock{
		methodFail: "AddFile",
	}
	builder := code.Golang{}
	_, err := builder.Build(wm, []byte(sourceCode))
	if errors.Cause(err) != errMockAddFile {
		t.Errorf("Expected %v; but got %v", errMockAddFile, err)
	}
}

func TestGolang_Build_Chdir(t *testing.T) {
	wm := &WorkdirManagerMock{
		methodFail: "Chdir",
	}
	builder := code.Golang{}
	_, err := builder.Build(wm, []byte(sourceCode))
	if errors.Cause(err) != errMockChdir {
		t.Errorf("Expected %v; but got %v", errMockChdir, err)
	}
}

func TestGolang_Build_Exec(t *testing.T) {
	wm := &WorkdirManagerMock{}
	builder := code.Golang{}
	_, err := builder.Build(wm, []byte("Invalid Golang code"))
	if errors.Cause(err) != code.ErrGoBuild {
		t.Errorf("Expected %v; but got %v", code.ErrGoBuild, err)
	}
}

func TestGolang_Build(t *testing.T) {
	wm, err := workspace.Create(os.TempDir())
	if err != nil {
		t.Errorf("Exected Error nil")
	}
	defer wm.RemoveAll()
	builder := code.Golang{}
	e, err := builder.Build(wm, []byte(sourceCode))
	if err != nil {
		t.Errorf("Expected %v; but got %v", nil, err)
	}

	if _, ok := e.(*engine.Engine); !ok {
		t.Errorf("Expected e == *engine.Engine")
	}
}
