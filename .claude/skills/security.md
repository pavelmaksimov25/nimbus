---
name: security
description: Run Go security analysis — gosec, govulncheck, go vet, and staticcheck
user_invocable: true
---

# Go Security Scan

Run a comprehensive security audit on this Go project using native Go tooling. No npm/npx — only Go binaries and bash.

## Instructions

Execute the following security checks **sequentially**. For each tool, install it if not already present, then run the scan. Collect all output and present a unified security report at the end.

### Step 1: go vet (built-in static analysis)

```bash
go vet ./...
```

### Step 2: govulncheck (known vulnerability detection)

Install if missing, then run:

```bash
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

This checks Go dependencies against the Go vulnerability database.

### Step 3: gosec (security-focused linter)

Install if missing, then run:

```bash
go install github.com/securego/gosec/v2/cmd/gosec@latest
gosec -fmt text -confidence medium -severity medium ./...
```

This detects common security issues: SQL injection, hardcoded credentials, insecure crypto, command injection, path traversal, weak random numbers, etc.

### Step 4: staticcheck (advanced static analysis)

Install if missing, then run:

```bash
go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck ./...
```

This catches bugs, performance issues, and code simplifications that `go vet` misses.

## Output Format

After running all four tools, present a unified report:

```
## Security Scan Results

### go vet
<status: clean or list of findings>

### govulncheck
<status: no known vulnerabilities or list of CVEs with affected packages>

### gosec
<status: clean or list of findings grouped by severity (HIGH/MEDIUM/LOW)>

### staticcheck
<status: clean or list of findings>

### Summary
- Total issues: <count>
- Critical/High: <count>
- Medium: <count>
- Low/Info: <count>
- Recommendations: <actionable next steps if any issues found>
```

If a tool fails to install or run, note the error and continue with the remaining tools. Do not stop the entire scan because one tool failed.
