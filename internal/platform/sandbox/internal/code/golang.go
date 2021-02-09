package code

import (
	"context"
	"path"
	"time"

	"github.com/pkg/errors"

	"github.com/josestg/codex/internal/platform/sandbox/internal/engine"
)

type Golang struct {
}

var ErrGoBuild = errors.New("ERROR BUILD GO BINARY")

// Build builds Go code into a binary file and returns the executable.
func (g Golang) Build(wm engine.WorkdirManager, sourceCode []byte) (engine.Executable, error) {
	filename := "main.go"
	binaryFilename := path.Join(wm.Location(), "binary")

	// Adds main.go into workdir.
	_, err := wm.AddFile(filename, sourceCode)
	if err != nil {
		return nil, errors.Wrap(err, "Golang: Build: Could not add a new file to the workdir.")
	}

	// Goto workdir.
	err = wm.Chdir()
	if err != nil {
		return nil, errors.Wrap(err, "Golang: Build: Could not changes the active directory to the workdir.")
	}

	// Create go builder.
	builder := engine.New("go", []string{"build", "-o", binaryFilename, filename})

	// Sets limit builds time maximum to 3 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Builds binary file.
	if err := builder.Exec(ctx).Run(); err != nil {
		return nil, ErrGoBuild
	}

	// Creates the executable engine.
	executor := engine.New(binaryFilename, nil)
	return executor, nil
}
