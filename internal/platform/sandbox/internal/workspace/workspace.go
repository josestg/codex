package workspace

import (
	"fmt"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const SandboxWorkspacePrefix = "CODEX_SANDBOX_"

// Workspace knows how to manage the sandbox workspace.
type Workspace struct {
	root      string
	namespace string
}

// Create creates a new workspace.
// Workspace automatically generates a unique namespace which is used as a workdir.
func Create(root string) (*Workspace, error) {
	// Generates a new namespace.
	uniqueDirName := uuid.New().String()
	namespace := fmt.Sprintf("%s_%s", SandboxWorkspacePrefix, uniqueDirName)

	fullPath := path.Join(root, namespace)
	err := os.MkdirAll(fullPath, os.ModePerm)

	if err != nil && os.IsNotExist(err) {
		return nil, errors.Wrapf(err, "Workspace: Create: Could not create workspace: %s", fullPath)
	}

	ws := Workspace{root: root, namespace: namespace}
	return &ws, nil
}

func (w *Workspace) Location() string {
	return path.Join(w.root, w.namespace)
}

func (w *Workspace) AddFile(name string, content []byte) (string, error) {
	filename := path.Join(w.Location(), name)
	file, err := os.Create(filename)
	if err != nil {
		return "", errors.Wrapf(err, "Workspace: AddFile: Could not create file: %v", filename)
	}

	_, err = file.Write(content)
	if err != nil {
		return "", errors.Wrap(err, "Workspace: AddFile: Could not writes content into file")
	}

	return filename, err
}

func (w *Workspace) Chdir() error {
	return os.Chdir(w.Location())
}

// RemoveALl removes all file and directory inside the workdir.
func (w *Workspace) RemoveAll() error {
	workspacePath := w.Location()
	if err := os.RemoveAll(workspacePath); err != nil {
		return errors.Wrapf(err, "Workspace: RemoveAll: Could not remove workspace: %s", workspacePath)
	}
	return nil
}
