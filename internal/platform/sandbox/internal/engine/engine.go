package engine

import (
	"context"
	"os/exec"
)

// Engine knows how to execute the command.
type Engine struct {
	name string
	args []string
}

// New creates a new engine.
func New(name string, args []string) *Engine {
	return &Engine{name: name, args: args}
}

// Exec creates a new cmd.
func (e *Engine) Exec(ctx context.Context) *exec.Cmd {
	return exec.CommandContext(ctx, e.name, e.args...)
}
