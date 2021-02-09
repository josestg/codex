package handlers

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/josestg/codex/internal/model"
	"github.com/josestg/codex/internal/platform/sandbox"
	"github.com/josestg/codex/internal/platform/web"
)

type CodeExecution struct {
	log *log.Logger
}

// Exec handle remote code execution request.
func (ce *CodeExecution) Exec(w http.ResponseWriter, r *http.Request) error {

	codexSubmission := model.CodexSubmission{}
	err := json.NewDecoder(r.Body).Decode(&codexSubmission)
	if err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	// Decode request body source_code from base64-encoded string into plain text.
	plainTextSourceCode, err := base64.StdEncoding.DecodeString(codexSubmission.SourceCode)
	if err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	// Creates a workspace for the sandbox.
	ws, err := sandbox.CreateWorkspace(os.TempDir())
	if err != nil {
		return errors.Wrap(err, "Exec: Could not create sandbox workspace")
	}

	// Removes workspace.
	defer func() {
		err := ws.RemoveAll()
		if err != nil {
			ce.log.Println("Exec:", err)
		}
	}()

	// Creates new sandbox
	sb, err := sandbox.Create(ce.log, ws, codexSubmission.Language, plainTextSourceCode)
	if err != nil {
		if errors.Cause(err) == sandbox.ErrUnknownProgrammingLanguage {
			return web.NewRequestError(err, http.StatusBadRequest)
		}
		return errors.Wrap(err, "Exec: Could not create sandbox")
	}

	timeLimit, err := time.ParseDuration(codexSubmission.TimeLimit)
	if err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	testLogs := sb.RunTestBulk(codexSubmission.TestCases, timeLimit)

	if err != web.Respond(w, testLogs, http.StatusOK) {
		ce.log.Printf("Could not write response: err %v", err)
	}

	return nil
}
