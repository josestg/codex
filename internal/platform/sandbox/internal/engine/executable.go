package engine

import (
	"context"
	"os/exec"
)

// Executable knows how to execute a binary file or command.
type Executable interface {
	Exec(ctx context.Context) *exec.Cmd
}

// WorkdirManager knows how to manage workdir.
type WorkdirManager interface {
	// Location returns the path to the workdir.
	Location() string
	// AddFile adds a new file into the workdir.
	AddFile(name string, content []byte) (string, error)
	// Chdir changes the active directory to the workdir.
	Chdir() error
}

// Builder knows how to build executable in a given workdir.
type Builder interface {
	Build(wm WorkdirManager, sourceCode []byte) (Executable, error)
}
