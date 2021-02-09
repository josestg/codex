package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"

	"github.com/josestg/codex/app/codex-api/internal/handlers"
)

func main() {
	logger := log.New(os.Stdout, "CODEX_API_", log.LstdFlags|log.Lshortfile)
	if err := run(logger); err != nil {
		logger.Fatal(err)
	}
}

func run(logger *log.Logger) error {

	// =========================================================================
	// Configuration

	var cfg struct {
		Web struct {
			Address         string        `conf:"default:localhost:8000"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
	}

	if err := conf.Parse(os.Args[1:], "CODEX", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("CODEX", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	// =========================================================================
	// Start API Service

	log.Printf("main : Started")
	defer log.Println("main : Completed")

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Printf("main : Config :\n%v\n", out)

	server := http.Server{
		Addr:         cfg.Web.Address,
		Handler:      handlers.API(logger),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}

	serverErr := make(chan error, 1)
	go func() {
		logger.Printf("main : API listening on %s", server.Addr)
		serverErr <- server.ListenAndServe()
	}()

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return errors.Wrap(err, "Could not start server.")
	case <-shutdownSignal:
		log.Println("main : Start shutdown")

		shutdownTimeout := cfg.Web.ShutdownTimeout
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", shutdownTimeout, err)
			err = server.Close()
			if err != nil {
				return errors.Wrap(err, "could not stop server gracefully")
			}
		}
	}
	return nil
}
