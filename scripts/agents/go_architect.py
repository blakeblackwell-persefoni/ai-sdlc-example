#!/usr/bin/env python3
"""
Go Architect Agent - Reviews Go API code for architectural best practices.

Usage:
    python scripts/agents/go_architect.py [--files FILE1 FILE2 ...]

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

SYSTEM_PROMPT = """You are a senior Go architect specializing in API design and distributed systems.

Your expertise includes:
- REST API design (resource modeling, HTTP semantics, status codes)
- gRPC and protobuf patterns
- Error handling strategies (wrapping, sentinel errors, custom error types)
- Concurrency patterns (goroutines, channels, context, sync primitives)
- Package structure and interface design
- Performance optimization for APIs
- Security best practices

## Review Process

1. **API Contract**: Examine endpoints, HTTP methods, request/response schemas
2. **Error Handling**: Verify comprehensive error wrapping and informative responses
3. **Concurrency**: Review goroutine safety, channel usage, context propagation
4. **Package Structure**: Assess modularity, interface design, dependency management
5. **Performance**: Identify bottlenecks, unnecessary allocations, N+1 patterns
6. **Security**: Check input validation, authentication/authorization patterns

## Output Format

Organize findings by priority:

### Critical
Security issues, data loss risks, race conditions

### Warnings
Design issues, performance problems, maintainability concerns

### Suggestions
Code quality improvements, alternative patterns

For each finding include:
- Clear explanation of the issue
- Why it matters (impact)
- How to fix it with code examples
- Reference to Go best practices when applicable

Be thorough but constructive. Focus on actionable improvements."""


def read_file(filepath: str) -> str:
    """Read a file and return its contents with the filename header."""
    try:
        with open(filepath, "r") as f:
            content = f.read()
        return f"### File: {filepath}\n```go\n{content}\n```\n"
    except Exception as e:
        return f"### File: {filepath}\nError reading file: {e}\n"


def find_go_files(directory: str = ".") -> list[str]:
    """Find all Go source files in the directory."""
    go_files = []
    for path in Path(directory).rglob("*.go"):
        if "_test.go" not in str(path):
            go_files.append(str(path))
    return sorted(go_files)


def find_api_specs(directory: str = ".") -> list[str]:
    """Find OpenAPI/Swagger specs."""
    specs = []
    for pattern in ["*.yaml", "*.yml", "*.json"]:
        for path in Path(directory).rglob(pattern):
            if "openapi" in str(path).lower() or "swagger" in str(path).lower():
                specs.append(str(path))
    # Also check root for openapi.yaml
    root_spec = Path(directory) / "openapi.yaml"
    if root_spec.exists() and str(root_spec) not in specs:
        specs.append(str(root_spec))
    return specs


def run_review(files: list[str] | None = None) -> str:
    """Run the architecture review."""
    client = anthropic.Anthropic()

    # Gather code to review
    if files:
        go_files = [f for f in files if f.endswith(".go")]
    else:
        go_files = find_go_files()

    if not go_files:
        return "No Go files found to review."

    # Build the code context
    code_context = "# Code to Review\n\n"
    for filepath in go_files:
        code_context += read_file(filepath)

    # Add API specs if available
    specs = find_api_specs()
    if specs:
        code_context += "\n# API Specifications\n\n"
        for spec in specs:
            try:
                with open(spec, "r") as f:
                    content = f.read()
                ext = Path(spec).suffix
                code_context += f"### File: {spec}\n```{ext[1:]}\n{content}\n```\n"
            except Exception as e:
                code_context += f"### File: {spec}\nError reading file: {e}\n"

    # Create the review request
    message = client.messages.create(
        model="claude-sonnet-4-20250514",
        max_tokens=4096,
        system=SYSTEM_PROMPT,
        messages=[
            {
                "role": "user",
                "content": f"Please review the following Go API code for architectural best practices:\n\n{code_context}",
            }
        ],
    )

    return message.content[0].text


def main():
    parser = argparse.ArgumentParser(
        description="Go Architect Agent - Reviews Go API code for architectural best practices"
    )
    parser.add_argument(
        "--files",
        nargs="*",
        help="Specific files to review (defaults to all Go files)",
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

    print("Running Go Architecture Review...\n", file=sys.stderr)

    review = run_review(args.files)

    if args.output:
        with open(args.output, "w") as f:
            f.write(review)
        print(f"Review written to {args.output}", file=sys.stderr)
    else:
        print(review)


if __name__ == "__main__":
    main()
