package handlers

import (
	"log"
	"net/http"

	"github.com/josestg/codex/internal/platform/web"
)

func API(log *log.Logger) http.Handler {
	app := web.NewApp(log)

	codex := &CodeExecution{log: log}

	app.Handle(http.MethodPost, "/api/v1/exec", codex.Exec)

	return app
}
