Prepare a new release for this project.

The argument should be a semver version like `v1.2.0`. If no argument is provided, ask the user what version to release.

Steps:
1. Ensure working tree is clean: `git status`
2. Run full verification: `go test ./... && go vet ./... && go build -o /dev/null .`
3. Validate JSON presets: `python3 -m json.tool scripts/*.json > /dev/null`
4. Show what will be released: `git log --oneline $(git describe --tags --abbrev=0 2>/dev/null || echo HEAD~10)..HEAD`
5. Ask user to confirm the release
6. Create and push the tag:
   ```
   git tag $ARGUMENTS
   git push origin $ARGUMENTS
   ```

The GitHub Actions release workflow will automatically build cross-platform binaries and create the GitHub Release.
