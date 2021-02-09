package sandbox

import "github.com/josestg/codex/internal/platform/sandbox/internal/workspace"

// CreateWorkspace creates a new workspace.
func CreateWorkspace(root string) (*workspace.Workspace, error) {
	return workspace.Create(root)
}
