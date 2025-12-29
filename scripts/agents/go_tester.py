#!/usr/bin/env python3
"""
Go Tester Agent - Analyzes Go code and suggests comprehensive tests.

Usage:
    python scripts/agents/go_tester.py [--files FILE1 FILE2 ...]

Requires:
    pip install anthropic
    export ANTHROPIC_API_KEY=your-key
"""

import argparse
import os
import sys
from pathlib import Path

try:
    import anthropic
except ImportError:
    print("Error: anthropic package not installed. Run: pip install anthropic")
    sys.exit(1)

SYSTEM_PROMPT = """You are a Go testing expert focused on writing comprehensive, maintainable tests.

Your expertise includes:
- Idiomatic Go unit tests using `testing.T`
- Table-driven test patterns
- Integration tests with proper setup/teardown
- HTTP handler testing with `httptest`
- Test fixtures, mocks, and test utilities
- Benchmarking performance-critical code
- Testing concurrent code safely

## Test Writing Standards

### Table-Driven Tests
```go
func TestOperation(t *testing.T) {
    tests := []struct {
        name    string
        input   InputType
        want    OutputType
        wantErr bool
    }{
        {"valid input", validInput, expectedOutput, false},
        {"invalid input", badInput, nil, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Operation(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Operation() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("Operation() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### HTTP Handler Tests
```go
func TestHandler(t *testing.T) {
    req := httptest.NewRequest(http.MethodPost, "/endpoint", strings.NewReader(`{"a":1,"b":2}`))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    handler(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
    }
}
```

## Output Format

Provide complete, runnable test code that:
1. Covers all public functions and handlers
2. Tests happy path scenarios
3. Tests error conditions
4. Tests edge cases and boundary conditions
5. Uses table-driven tests where appropriate
6. Includes clear test names that describe the scenario

Output the complete test file(s) ready to be saved and run."""


def read_file(filepath: str) -> str:
    """Read a file and return its contents with the filename header."""
    try:
        with open(filepath, "r") as f:
            content = f.read()
        return f"### File: {filepath}\n```go\n{content}\n```\n"
    except Exception as e:
        return f"### File: {filepath}\nError reading file: {e}\n"


def find_go_files(directory: str = ".") -> list[str]:
    """Find all Go source files (excluding tests)."""
    go_files = []
    for path in Path(directory).rglob("*.go"):
        if "_test.go" not in str(path):
            go_files.append(str(path))
    return sorted(go_files)


def find_test_files(directory: str = ".") -> list[str]:
    """Find existing test files."""
    return sorted([str(p) for p in Path(directory).rglob("*_test.go")])


def run_test_generation(files: list[str] | None = None) -> str:
    """Generate tests for the given files."""
    client = anthropic.Anthropic()

    # Gather code to test
    if files:
        go_files = [f for f in files if f.endswith(".go") and "_test.go" not in f]
    else:
        go_files = find_go_files()

    if not go_files:
        return "No Go files found to generate tests for."

    # Build the code context
    code_context = "# Source Code to Test\n\n"
    for filepath in go_files:
        code_context += read_file(filepath)

    # Include existing tests for context
    existing_tests = find_test_files()
    if existing_tests:
        code_context += "\n# Existing Tests (for reference)\n\n"
        for filepath in existing_tests:
            code_context += read_file(filepath)

    # Create the test generation request
    message = client.messages.create(
        model="claude-sonnet-4-20250514",
        max_tokens=4096,
        system=SYSTEM_PROMPT,
        messages=[
            {
                "role": "user",
                "content": f"Please generate comprehensive tests for the following Go code:\n\n{code_context}\n\nProvide complete, runnable test files.",
            }
        ],
    )

    return message.content[0].text


def main():
    parser = argparse.ArgumentParser(
        description="Go Tester Agent - Generates comprehensive tests for Go code"
    )
    parser.add_argument(
        "--files",
        nargs="*",
        help="Specific files to generate tests for (defaults to all Go files)",
    )
    parser.add_argument(
        "--output",
        "-o",
        help="Output file for the generated tests (defaults to stdout)",
    )

    args = parser.parse_args()

    # Check for API key
    if not os.environ.get("ANTHROPIC_API_KEY"):
        print("Error: ANTHROPIC_API_KEY environment variable not set")
        sys.exit(1)

    print("Generating Go Tests...\n", file=sys.stderr)

    tests = run_test_generation(args.files)

    if args.output:
        with open(args.output, "w") as f:
            f.write(tests)
        print(f"Tests written to {args.output}", file=sys.stderr)
    else:
        print(tests)


if __name__ == "__main__":
    main()
