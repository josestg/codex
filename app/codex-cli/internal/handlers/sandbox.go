package handlers

import (
	"encoding/base64"
	"log"
	"os"
	"time"

	"github.com/josestg/codex/internal/model"
	"github.com/josestg/codex/internal/platform/sandbox"
)

func Eval(logger *log.Logger, submission *model.CodexSubmission) ([]*model.SandboxTestStdout, error) {

	// Decode request body source_code from base64-encoded string into plain text.
	plainTextSourceCode, err := base64.StdEncoding.DecodeString(submission.SourceCode)
	if err != nil {
		return nil, err
	}

	// Creates a workspace for the sandbox.
	ws, err := sandbox.CreateWorkspace(os.TempDir())
	if err != nil {
		return nil, err
	}

	// Removes workspace.
	defer func() {
		err := ws.RemoveAll()
		if err != nil {
			logger.Println("Exec:", err)
		}
	}()

	// Creates new sandbox
	sb, err := sandbox.Create(logger, ws, submission.Language, plainTextSourceCode)
	if err != nil {
		return nil, err
	}

	timeLimit, err := time.ParseDuration(submission.TimeLimit)
	if err != nil {
		return nil, err
	}

	testLogs := sb.RunTestBulk(submission.TestCases, timeLimit)
	return testLogs, nil
}
