Run the full verification pipeline for this project:

1. Format all Go files: `gofmt -s -w .`
2. Run static analysis: `go vet ./...`
3. Run all tests: `go test ./... -v`
4. Validate all JSON presets: `python3 -m json.tool scripts/*.json > /dev/null`
5. Quick compile check: `go build -o /dev/null .`

Report any failures with details. If everything passes, confirm with a summary.
