# CODEX

CODEX or Remote Code Execution. (_Work in process_)

## Usage

```bazaar
git clone https://github.com/josestg/codex.git

cd codex

make codex-api 
# or
go run un app/codex-api/main.go
```

CODEX API listening at http://localhost:8000

## API

| Endpoint | Params | Header | Body | Response |
|--- | --- | --- | --- | ---  |
| __POST__ /api/v1/exec | `none` | `none` | [CodexSubmission](https://github.com/josestg/codex/blob/b797c9c8cc7a743e14d3f0861a4586fcab3a55aa/internal/model/sandbox.go#L3) | Slice of [SandboxTestStdout](https://github.com/josestg/codex/blob/b797c9c8cc7a743e14d3f0861a4586fcab3a55aa/internal/model/sandbox.go#L25)

### Example

1. Create a plain source code.

```go
package main

import (
	"fmt"
	"strings"
)

func Solve(n int) string {
	str := make([]string, n)
	content := byte('*')

	for i := 0; i < n; i++ {

		line := make([]byte, n-i)
		for j := 0; j < len(line); j++ {
			line[j] = content
		}

		str[i] = string(line)
	}

	return strings.Join(str, "\n")
}

func main() {
	var n int
	fmt.Scanf("%d", &n)
	out := Solve(n)
	fmt.Println(out)
}
```

2. Encode the plain source code into a base64 string.

```bazaar
cGFja2FnZSBtYWluCgppbXBvcnQgKAoJImZtdCIKCSJzdHJpbmdzIgopCgpmdW5jIFNvbHZlKG4gaW50KSBzdHJpbmcgewoJc3RyIDo9IG1ha2UoW11zdHJpbmcsIG4pCgljb250ZW50IDo9IGJ5dGUoJyonKQoKCWZvciBpIDo9IDA7IGkgPCBuOyBpKysgewoKCQlsaW5lIDo9IG1ha2UoW11ieXRlLCBuLWkpCgkJZm9yIGogOj0gMDsgaiA8IGxlbihsaW5lKTsgaisrIHsKCQkJbGluZVtqXSA9IGNvbnRlbnQKCQl9CgoJCXN0cltpXSA9IHN0cmluZyhsaW5lKQoJfQoKCXJldHVybiBzdHJpbmdzLkpvaW4oc3RyLCAiXG4iKQp9CgpmdW5jIG1haW4oKSB7Cgl2YXIgbiBpbnQKCWZtdC5TY2FuZigiJWQiLCAmbikKCW91dCA6PSBTb2x2ZShuKQoJZm10LlByaW50bG4ob3V0KQp9Cg==
```

3. Create a JSON file as a request body. (_request.json_)

```json
{
  "source_code": "cGFja2FnZSBtYWluCgppbXBvcnQgKAoJImZtdCIKCSJzdHJpbmdzIgopCgpmdW5jIFNvbHZlKG4gaW50KSBzdHJpbmcgewoJc3RyIDo9IG1ha2UoW11zdHJpbmcsIG4pCgljb250ZW50IDo9IGJ5dGUoJyonKQoKCWZvciBpIDo9IDA7IGkgPCBuOyBpKysgewoKCQlsaW5lIDo9IG1ha2UoW11ieXRlLCBuLWkpCgkJZm9yIGogOj0gMDsgaiA8IGxlbihsaW5lKTsgaisrIHsKCQkJbGluZVtqXSA9IGNvbnRlbnQKCQl9CgoJCXN0cltpXSA9IHN0cmluZyhsaW5lKQoJfQoKCXJldHVybiBzdHJpbmdzLkpvaW4oc3RyLCAiXG4iKQp9CgpmdW5jIG1haW4oKSB7Cgl2YXIgbiBpbnQKCWZtdC5TY2FuZigiJWQiLCAmbikKCW91dCA6PSBTb2x2ZShuKQoJZm10LlByaW50bG4ob3V0KQp9Cg==",
  "language": "golang",
  "time_limit": "1s",

  "test_cases": [
    {
      "id": "1",
      "input": "1",
			"is_private": false,
      "expected_output": "*\n"
    },
    {
      "id": "2",
      "input": "2",
			"is_private": true,
      "expected_output": "**\n*\n"
    },{
      "id": "3",
      "input": "3",
			"is_private": false,
      "expected_output": "*\n**\n**\n***"
    }
  ]
}
```

4. Make an HTTP request.

```shell
curl -XPOST -d @request.json http://localhost:8000/api/v1/exec
```

5. Response body.

```shell
[
  {
    "test_case_id": "3",
    "status": "FAILED",
    "log": {
      "id": 706991,
      "stdin": "3",
      "stdout": "***\n**\n*\n",
      "stderr": "",
      "running_time": "2ms",
      "expected_stdout": "*\n**\n**\n***"
    },
    "error": null
  },
  {
    "test_case_id": "1",
    "status": "PASSED",
    "log": {
      "id": 706990,
      "stdin": "1",
      "stdout": "*\n",
      "stderr": "",
      "running_time": "0ms",
      "expected_stdout": "*\n"
    },
    "error": null
  },
  {
    "test_case_id": "2",
    "status": "PASSED",
    "log": null,
    "error": null
  }
]
```