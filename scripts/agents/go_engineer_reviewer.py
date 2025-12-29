#!/usr/bin/env python3
"""
Go Engineer Reviewer Agent - Conducts thorough code reviews.

Usage:
    python scripts/agents/go_engineer_reviewer.py [--files FILE1 FILE2 ...]
    python scripts/agents/go_engineer_reviewer.py --diff  # Review git diff

Requires:
    pip install anthropic
    export ANTHROPIC_API_KEY=your-key
"""

import argparse
import os
import subprocess
import sys
from pathlib import Path

try:
    import anthropic
except ImportError:
    print("Error: anthropic package not installed. Run: pip install anthropic")
    sys.exit(1)

SYSTEM_PROMPT = """You are a senior Go engineer conducting a thorough code review.

## Review Dimensions

### 1. Go Idioms
- Idiomatic error handling (`if err != nil`)
- Error wrapping with context (`fmt.Errorf("context: %w", err)`)
- Interface design (accept interfaces, return structs)
- Effective use of standard library

### 2. Code Quality
- Clear, descriptive naming (packages, functions, variables)
- Single responsibility principle
- Appropriate abstraction levels
- DRY without over-abstraction

### 3. Concurrency Safety
- No data races (shared state properly synchronized)
- Goroutine lifecycle management (no leaks)
- Proper context usage and cancellation
- Channel patterns (buffered vs unbuffered)

### 4. Performance
- Unnecessary allocations (preallocate slices when size known)
- String concatenation in loops (use strings.Builder)
- Efficient data structures
- Database query patterns (N+1 problems)

### 5. Security
- Input validation at boundaries
- SQL injection prevention (parameterized queries)
- Path traversal protection
- No secrets in code
- Proper authentication/authorization

### 6. Error Handling
- All errors checked
- Errors wrapped with context
- Appropriate error types (sentinel vs wrapped vs custom)
- User-facing vs internal error messages

### 7. Testing
- Adequate test coverage
- Edge cases covered
- Tests are deterministic
- Clear test naming

## Output Format

For each finding:

**[Category]** - Severity: Critical|Warning|Suggestion

**Finding:** Brief description of the issue

**Location:** file:line or code snippet

**Impact:** Why this matters

**Recommendation:** How to fix it with code example

Be thorough but constructive. Prioritize findings by severity."""


def read_file(filepath: str) -> str:
    """Read a file and return its contents with the filename header."""
    try:
        with open(filepath, "r") as f:
            content = f.read()
        return f"### File: {filepath}\n```go\n{content}\n```\n"
    except Exception as e:
        return f"### File: {filepath}\nError reading file: {e}\n"


def get_git_diff(base: str = "main") -> str:
    """Get the git diff against a base branch."""
    try:
        result = subprocess.run(
            ["git", "diff", base, "--", "*.go"],
            capture_output=True,
            text=True,
            check=True,
        )
        return result.stdout
    except subprocess.CalledProcessError as e:
        return f"Error getting git diff: {e.stderr}"


def get_changed_files(base: str = "main") -> list[str]:
    """Get list of changed Go files."""
    try:
        result = subprocess.run(
            ["git", "diff", "--name-only", base, "--", "*.go"],
            capture_output=True,
            text=True,
            check=True,
        )
        return [f.strip() for f in result.stdout.strip().split("\n") if f.strip()]
    except subprocess.CalledProcessError:
        return []


def find_go_files(directory: str = ".") -> list[str]:
    """Find all Go source files."""
    return sorted([str(p) for p in Path(directory).rglob("*.go")])


def run_review(
    files: list[str] | None = None, use_diff: bool = False, base: str = "main"
) -> str:
    """Run the code review."""
    client = anthropic.Anthropic()

    code_context = ""

    if use_diff:
        # Review only the diff
        diff = get_git_diff(base)
        if not diff or diff.startswith("Error"):
            return f"No changes to review or error: {diff}"

        code_context = f"# Git Diff (against {base})\n\n```diff\n{diff}\n```\n"

        # Also include full files for context
        changed_files = get_changed_files(base)
        if changed_files:
            code_context += "\n# Full Files (for context)\n\n"
            for filepath in changed_files:
                if os.path.exists(filepath):
                    code_context += read_file(filepath)
    else:
        # Review specified files or all Go files
        if files:
            go_files = [f for f in files if f.endswith(".go")]
        else:
            go_files = find_go_files()

        if not go_files:
            return "No Go files found to review."

        code_context = "# Code to Review\n\n"
        for filepath in go_files:
            code_context += read_file(filepath)

    # Create the review request
    message = client.messages.create(
        model="claude-sonnet-4-20250514",
        max_tokens=4096,
        system=SYSTEM_PROMPT,
        messages=[
            {
                "role": "user",
                "content": f"Please conduct a thorough code review:\n\n{code_context}",
            }
        ],
    )

    return message.content[0].text


def main():
    parser = argparse.ArgumentParser(
        description="Go Engineer Reviewer Agent - Conducts thorough code reviews"
    )
    parser.add_argument(
        "--files",
        nargs="*",
        help="Specific files to review (defaults to all Go files)",
    )
    parser.add_argument(
        "--diff",
        action="store_true",
        help="Review only the git diff instead of full files",
    )
    parser.add_argument(
        "--base",
        default="main",
        help="Base branch for diff comparison (default: main)",
    )
    parser.add_argument(
        "--output",
        "-o",
        help="Output file for the review (defaults to stdout)",
    )

    args = parser.parse_args()

    # Check for API key
    if not os.environ.get("ANTHROPIC_API_KEY"):
        print("Error: ANTHROPIC_API_KEY environment variable not set")
        sys.exit(1)

    print("Running Go Code Review...\n", file=sys.stderr)

    review = run_review(args.files, args.diff, args.base)

    if args.output:
        with open(args.output, "w") as f:
            f.write(review)
        print(f"Review written to {args.output}", file=sys.stderr)
    else:
        print(review)


if __name__ == "__main__":
    main()
