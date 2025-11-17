# **CONTRIBUTING.md**

## Contributing to Net-Nudge

Thank you for your interest in contributing to **Net-Nudge**, a modular Go-based penetration testing scan scheduler and notifier.
This document describes how to contribute code, tests, documentation, and improvements to the project.

---

## **Project Overview**

Net-Nudge is a Go application that:

* Schedules and runs security scans (domains, subdomains, services)
* Supports multiple scanning tools (Nmap, Amass, custom scanners)
* Loads configuration from YAML
* Uses a plug-in style job registry
* Provides both direct CLI usage and an interactive menu
* Includes automated linting and testing via GitHub Actions

---

# 1. Getting Started

## **Prerequisites**

* Go **1.22+**
* Make (optional but recommended)
* GitHub account if submitting PRs

## **Install Dependencies**

```sh
go mod tidy
```

## **Project Structure**

```
/cmd/netnudge      → main application
/internal/jobs     → job definitions (Nmap, Amass, Subdomain scans, etc.)
/internal/registry → job registry and plugin lookup
/internal/scheduler→ scan scheduler + notifier logic
/config            → example YAML configuration files
/tests             → Go test suite
```

---

# 2. How to Contribute

## **Step 1: Fork and Clone**

```sh
git clone https://github.com/<yourname>/net-nudge.git
cd net-nudge
```

## **Step 2: Create a Feature Branch**

```sh
git checkout -b feature/<short-description>
```

## **Step 3: Write Your Code**

* Follow Go best practices.
* Keep functions small and focused.
* Use interfaces when adding new scanning tools.
* Place new jobs in `internal/jobs/` and register them in `internal/registry/`.

---

# 3. Coding Guidelines

### **General Go Style**

* Follow standard `gofmt` formatting.
* Do not leave unused imports, variables, or commented-out code.
* Prefer dependency injection over global state.

### **Adding a New Scan Job**

Each scanning tool should implement:

```go
type Job interface {
    Name() string
    Run(ctx context.Context, cfg any) error
}
```

Place job code in:

```
internal/jobs/<toolname>.go
```

Register it in:

```
internal/registry/registry.go
```

---

# 4. Configuration Guidelines (YAML)

Net-Nudge uses YAML config files under `/config/`.
When adding new job types, update:

* The YAML schema
* Struct bindings in the scheduler
* Example configurations

YAML files **must** be validated before running scans.

---

# 5. Testing

All contributions must include tests when possible.

### **Run the full test suite**

```sh
go test ./... -v
```

### **Coverage report**

```sh
go test ./... -cover
```

### **Where to Add Tests**

* Unit tests → `/tests/`
* Package-specific tests → alongside source files
* Integration tests → `/tests/integration/` (optional)

Tests should cover:

* Registry behavior
* Job execution
* Scheduler logic
* YAML config parsing
* Menu interactions (non-interactive mode)

---

# 6. Linting & CI Requirements

This project uses GitHub Actions for automated linting and testing.

### **Local linting**

Install:

```sh
golangci-lint run ./...
```

### **GitHub Workflows Included**

* **Go Lint** (`lint.yml`)
* **Go Test with Coverage** (`test.yml`)

Before submitting a PR:

* `golangci-lint` must pass
* `go test ./...` must pass
* No failing coverage checks
* CI must pass completely

---

# 7. Commit and PR Standards

### **Commit Messages**

Use clear, descriptive commit messages:

```
feat: add new amass domain discovery job
fix: resolve YAML config unmarshalling issue
refactor: reorganize scheduler interface
test: add coverage for registry lookups
docs: update configuration examples
```

### **Pull Requests**

All PRs must include:

* Summary of changes
* Test coverage (if applicable)
* Link to related issue (if any)

PRs should be **small, focused, and reviewable**.

---

# 8. Security Contributions

Because this is a penetration testing tool:

* Do not submit malicious payloads.
* Avoid adding tools or features that violate legal guidelines.
* Ensure all scanners require explicit user configuration and consent.

If unsure, open an issue before submitting a PR.

---

# 9. Reporting Issues

When reporting bugs, include:

* OS and Go version
* Steps to reproduce
* Full error output
* Relevant YAML configuration

Example:

```
go version go1.22 linux/amd64
nmap job fails when using custom ports
```

---

# 10. License

By contributing, you agree that your code will be licensed under the project’s existing license (MIT unless modified).

---

# 11. Thank You

Your contributions help make Net-Nudge a more capable and reliable penetration testing tool.
Whether you fix a bug or add an entire scanning module—thank you.
