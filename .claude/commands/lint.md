Run linting and fix any issues found:

1. Run `golangci-lint run` and capture the output
2. If there are errors, fix each one:
   - gocritic issues: fix the flagged patterns
   - gofmt issues: run `gofmt -s -w` on the affected files
   - other lint errors: apply the appropriate fix
3. Re-run `golangci-lint run` to confirm all issues are resolved
4. Report what was fixed
