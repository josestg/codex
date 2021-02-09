package code

import (
	"github.com/pkg/errors"

	"github.com/josestg/codex/internal/platform/sandbox/internal/engine"
)

type JavaScript struct {
}

// Build creates a command to execute javascript code.
func (js *JavaScript) Build(wm engine.WorkdirManager, sourceCode []byte) (engine.Executable, error) {
	filename := "main.js"

	// Adds main.js into workdir.
	fullPath, err := wm.AddFile(filename, sourceCode)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not adds %s to workspace", filename)
	}

	// Creates the executable command.
	executor := engine.New("node", []string{fullPath})
	return executor, err
}
