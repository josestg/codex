package model

type CodexSubmission struct {
	QuestionID string          `json:"question_id"`
	SourceCode string          `json:"source_code"`
	Language   string          `json:"language"`
	TimeLimit  string          `json:"time_limit"`
	TestCases  []*SandboxStdin `json:"test_cases"`
}

type SandboxStdin struct {
	ID             string `json:"id"`
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
	IsPrivate      bool   `json:"is_private"`
}

type SandboxStdout struct {
	TestCaseID string      `json:"test_case_id"`
	Status     string      `json:"status"`
	Log        *SandboxLog `json:"log"`
	Error      error       `json:"error"`
}

type SandboxTestStdout struct {
	TestCaseID string                  `json:"test_case_id"`
	Status     string                  `json:"status"`
	Log        *SandboxLog             `json:"log"`
	Error      *SandboxTestStdoutError `json:"error"`
}

type SandboxTestStdoutError struct {
	Message string `json:"message"`
}

type SandboxLog struct {
	ID             int    `json:"id"`
	Stdin          string `json:"stdin"`
	Stdout         string `json:"stdout"`
	Stderr         string `json:"stderr"`
	RunningTime    string `json:"running_time"`
	ExpectedStdout string `json:"expected_stdout"`
}
