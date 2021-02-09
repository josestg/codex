package code

import (
	"github.com/pkg/errors"

	"github.com/josestg/codex/internal/platform/sandbox/internal/engine"
)

type Python struct {
}

func (p *Python) Build(wm engine.WorkdirManager, sourceCode []byte) (engine.Executable, error) {
	filename := "main.py"

	// Adds main.py into workdir.
	fullPath, err := wm.AddFile(filename, sourceCode)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not adds %s to workspace", filename)
	}

	// Creates the executable command.
	executor := engine.New("python3", []string{fullPath})
	return executor, err
}
