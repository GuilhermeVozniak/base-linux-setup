Help create a new JSON preset for this project.

Ask the user for:
1. Target distribution/OS name
2. What tasks they want included

Then:
1. Create the JSON file in `scripts/<distro-name>.json` following the schema in CLAUDE.md
2. Include a proper `match` field for auto-detection
3. Validate the JSON: `python3 -m json.tool scripts/<filename>.json`
4. Verify it loads: `go build -o /tmp/bls-test . && /tmp/bls-test list-presets`
5. Clean up: `rm /tmp/bls-test`

Use existing presets in `scripts/` as reference for style and structure.
