package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/josestg/codex/app/codex-cli/internal/handlers"
	"github.com/josestg/codex/internal/model"
)

func main() {
	logger := log.New(os.Stderr, "CODEX_CLI_", log.LstdFlags|log.Lshortfile)

	var str string
	flag.StringVar(&str, "data", "", "base64-encoded")
	flag.Parse()

	if str == "" {
		logger.Fatal("data is required")
	}

	b, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		logger.Fatal(err)
	}

	var sub model.CodexSubmission
	if err := json.Unmarshal(b, &sub); err != nil {
		logger.Fatal(err)
	}

	out, err := handlers.Eval(logger, &sub)
	if err != nil {
		logger.Fatal(err)
	}

	d, err := json.Marshal(out)
	if err != nil {
		logger.Fatal(err)
	}

	fmt.Println(string(d))
}
