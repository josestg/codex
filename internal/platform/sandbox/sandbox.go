package sandbox

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/josestg/codex/internal/model"
	"github.com/josestg/codex/internal/platform/sandbox/internal/code"
	"github.com/josestg/codex/internal/platform/sandbox/internal/engine"
)

var (
	ErrTimeLimitExceeded = errors.New("TIME LIMIT EXCEEDED")
)

// Sandbox knows how to run User's code remotely.
type Sandbox struct {
	executable engine.Executable
	logger     *log.Logger
}

// Create creates a new sandbox.
func Create(logger *log.Logger, wm engine.WorkdirManager, language string, sourceCode []byte) (*Sandbox, error) {
	builder, err := code.GetBuilder(language)
	if err != nil {
		return nil, errors.Wrap(err, "Sandbox: Create: Could not get executable builder.")
	}

	executable, err := builder.Build(wm, sourceCode)
	if err != nil {
		return nil, errors.Wrap(err, "Sandbox: Create: Could not executable the executable.")
	}

	sb := Sandbox{executable: executable, logger: logger}
	return &sb, nil
}

// RunTest runs User's code with the given test case inside the sandbox.
func (s *Sandbox) RunTest(sandboxStdin *model.SandboxStdin, timeLimit time.Duration) *model.SandboxStdout {

	testResult := model.SandboxStdout{
		TestCaseID: sandboxStdin.ID,
		Status:     "FAILED",
		Error:      nil,
		Log:        nil,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeLimit)
	defer cancel()

	process := s.executable.Exec(ctx)

	// Writes test case input into process stdin.
	stdinWriter, err := process.StdinPipe()
	if err != nil {
		testResult.Error = errors.Wrap(err, "Sandbox: run: Could not create stdin writer.")
		return &testResult
	}
	go func() {
		defer stdinWriter.Close()
		io.WriteString(stdinWriter, sandboxStdin.Input)
	}()

	// Reads process stderr and stores into stderr variable.
	var stderr strings.Builder
	stderrReader, err := process.StderrPipe()

	if err != nil {
		testResult.Error = errors.Wrap(err, "Sandbox: run: Could not create stderr reader.")
		return &testResult
	}
	go func() {
		defer stderrReader.Close()
		io.Copy(&stderr, stderrReader)
	}()

	// Reads process stderr and stores into stdout variable.
	var stdout strings.Builder
	stdoutReader, err := process.StdoutPipe()
	if err != nil {
		testResult.Error = errors.Wrap(err, "Sandbox: run: Could not create stdout reader.")
		return &testResult
	}
	go func() {
		defer stdoutReader.Close()
		io.Copy(&stdout, stdoutReader)
	}()

	// Runs a timed process.
	if err := process.Run(); err != nil {
		testResult.Error = ErrTimeLimitExceeded
		return &testResult
	}

	runningTime := fmt.Sprintf("%dms", process.ProcessState.UserTime().Milliseconds())
	processStdout := stdout.String()

	if sandboxStdin.ExpectedOutput == processStdout {
		testResult.Status = "PASSED"
	}

	if !sandboxStdin.IsPrivate {
		testResult.Log = &model.SandboxLog{
			ID:             process.ProcessState.Pid(),
			Stdin:          sandboxStdin.Input,
			Stdout:         processStdout,
			Stderr:         stderr.String(),
			RunningTime:    runningTime,
			ExpectedStdout: sandboxStdin.ExpectedOutput,
		}
	}

	return &testResult
}

// RunTestBulk runs User's code with the given test cases concurrently inside the sandbox.
func (s *Sandbox) RunTestBulk(stdinList []*model.SandboxStdin, timeLimit time.Duration) []*model.SandboxTestStdout {
	wg := new(sync.WaitGroup)
	// Create a buffer channel to store the output of each test case running on each goroutine.
	runTestResultChan := make(chan *model.SandboxStdout, len(stdinList))
	for _, sandboxStdin := range stdinList {
		wg.Add(1)
		go func(tc *model.SandboxStdin) {
			defer wg.Done()
			runTestResult := s.RunTest(tc, timeLimit)
			runTestResultChan <- runTestResult
		}(sandboxStdin)
	}

	// To make sure the channels are closed after all the goroutines are finished.
	go func() {
		wg.Wait()
		close(runTestResultChan)
	}()

	// Reads data from runTestResultChan channel and stores into result variable.
	var result []*model.SandboxTestStdout
	for runTestResult := range runTestResultChan {
		sandboxTestStdout := model.SandboxTestStdout{
			TestCaseID: runTestResult.TestCaseID,
			Status:     runTestResult.Status,
			Log:        runTestResult.Log,
		}

		if runTestResult.Error != nil {
			if errors.Cause(runTestResult.Error) == ErrTimeLimitExceeded {
				sandboxTestStdout.Error = &model.SandboxTestStdoutError{
					Message: runTestResult.Error.Error(),
				}
			} else {
				s.logger.Println(runTestResult.Error)
				sandboxTestStdout.Error = &model.SandboxTestStdoutError{
					Message: "INTERNAL CODEX ERROR",
				}
			}
		}

		result = append(result, &sandboxTestStdout)
	}
	return result
}
