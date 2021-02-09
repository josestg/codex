
SHELL := /bin/bash

codex-api:
	 go run app/codex-api/main.go

codex-cli:
	go run app/codex-cli/main.go --data $(data)