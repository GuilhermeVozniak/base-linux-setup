# Contributing to Base Linux Setup

Thank you for your interest in contributing to Base Linux Setup! This document provides guidelines and information for contributors.

## 🚀 Quick Start

1. **Fork the repository** on GitHub

2. **Clone your fork** locally:

   ```bash
   git clone https://github.com/YOUR_USERNAME/base-linux-setup.git
   cd base-linux-setup
   ```

3. **Create a feature branch**:

   ```bash
   git checkout -b feature/your-feature-name
   ```

4. **Make your changes** and test them

5. **Commit your changes**:

   ```bash
   git commit -m "Add: your feature description"
   ```

6. **Push to your fork**:

   ```bash
   git push origin feature/your-feature-name
   ```

7. **Create a Pull Request** on GitHub

## 📋 Types of Contributions

### 🐛 Bug Reports

- Use the bug report template
- Include system information (`./base-linux-setup detect`)
- Provide steps to reproduce
- Include error output if applicable

### ✨ Feature Requests

- Use the feature request template
- Clearly describe the use case
- Consider implementation complexity
- Discuss in issues before large changes

### 🛠️ New Presets

- Use the preset request template
- Test on target systems
- Follow JSON format in `scripts/README.md`
- Include comprehensive task lists

### 📚 Documentation

- Fix typos and unclear sections
- Add examples and use cases
- Update README for new features
- Contribute to the Wiki

### 🔧 Code Improvements

- Bug fixes
- Performance improvements
- Code refactoring
- Test additions

## 🏗️ Development Setup

### Prerequisites

- Go 1.21 or higher
- Make (for build automation)
- Git

### Building

```bash
# Install dependencies
go mod download

# Build the application
make build

# Run tests
make test

# Run linting
make lint

# Format code
make fmt
```

### Testing

```bash
# Run all tests
go test ./...

# Test specific package
go test ./internal/presets

# Run with coverage
go test -cover ./...

# Test the built binary
./build/base-linux-setup --help
./build/base-linux-setup detect
./build/base-linux-setup list-presets
```

## 📝 Code Guidelines

### Go Code Style

- Follow standard Go formatting (`go fmt`)
- Use meaningful variable and function names
- Add comments for exported functions
- Handle errors properly
- Use context where appropriate

### Commit Messages

Use conventional commit format:

- `feat:` new feature
- `fix:` bug fix
- `docs:` documentation changes
- `style:` formatting, missing semicolons, etc.
- `refactor:` code change that neither fixes a bug nor adds a feature
- `test:` adding missing tests
- `chore:` changes to build process, auxiliary tools, libraries

Examples:

```
feat: add support for CentOS preset
fix: resolve JSON parsing error for special characters
docs: update installation instructions
```

### File Organization

```
base-linux-setup/
├── cmd/                 # CLI commands
├── internal/            # Internal packages
│   ├── detector/        # Environment detection
│   ├── presets/         # Preset management
│   ├── ui/             # User interface
│   └── executor/        # Task execution
├── scripts/            # JSON presets
└── .github/            # GitHub workflows and templates
```

## 🧪 Adding New Presets

All presets are JSON files in the `scripts/` directory. No Go code changes are needed.

1. **Create JSON file** in `scripts/` directory with a `match` field:

   ```json
   {
     "name": "My OS Setup",
     "environment": "My OS",
     "description": "Setup for My OS",
     "match": {
       "distribution": "myos"
     },
     "tasks": [...]
   }
   ```

2. **Rebuild** with `make build` — the new file is automatically embedded

3. **Test** with `./build/base-linux-setup list-presets`

The `match` field controls auto-detection. Available criteria:
- `distribution` — case-insensitive substring match against detected distro
- `os` — match against detected OS name
- `architecture` — match against detected architecture
- `is_raspberry_pi` — `true`/`false` for exact Raspberry Pi match

The most specific match (most fields) wins. A preset without `match` is the default fallback.

### Testing Presets

1. **Syntax validation**:

   ```bash
   # For JSON presets
   python3 -m json.tool scripts/my-preset.json
   ```

2. **Application testing**:

   ```bash
   make build
   ./build/base-linux-setup list-presets
   ./build/base-linux-setup detect
   ```

3. **Target system testing**:
   - Test on actual target OS
   - Verify all tasks execute correctly
   - Check for permission issues
   - Validate installed software

## 🔍 Pull Request Guidelines

### Before Submitting

- [ ] Code follows project style guidelines
- [ ] Tests pass (`make test`)
- [ ] Build succeeds (`make build`)
- [ ] Documentation updated if needed
- [ ] Commit messages follow convention
- [ ] No debugging code or temporary files

### PR Description

Include:

- **What**: Brief description of changes
- **Why**: Motivation for the change
- **How**: Technical approach if complex
- **Testing**: How you tested the changes
- **Screenshots**: If UI changes

### Review Process

1. Automated CI checks must pass
2. Code review by maintainers
3. Testing on relevant systems
4. Documentation review
5. Merge after approval

## 🏷️ Release Process

Releases are automated via GitHub Actions:

1. **Version tagging**:

   ```bash
   git tag v1.2.3
   git push origin v1.2.3
   ```

2. **Automatic builds** for all platforms
3. **Release creation** with binaries and checksums
4. **Changelog generation** from commits

### Version Scheme

We use Semantic Versioning (semver):

- `v1.0.0` - Major release (breaking changes)
- `v1.1.0` - Minor release (new features)
- `v1.1.1` - Patch release (bug fixes)

## 🤝 Community Guidelines

### Code of Conduct

- Be respectful and inclusive
- Welcome newcomers
- Focus on constructive feedback
- Help others learn and grow

### Communication

- **Issues**: Bug reports, feature requests
- **Discussions**: General questions, ideas
- **Pull Requests**: Code contributions
- **Wiki**: Collaborative documentation

### Getting Help

- Check existing issues and documentation
- Ask questions in discussions
- Join community conversations
- Contribute to help others

## 📚 Resources

- [Go Documentation](https://golang.org/doc/)
- [GitHub Flow](https://guides.github.com/introduction/flow/)
- [Semantic Versioning](https://semver.org/)
- [Conventional Commits](https://conventionalcommits.org/)

## 🎉 Recognition

Contributors are recognized in:

- Release notes
- README.md acknowledgments
- GitHub contributors graph
- Project wiki

Thank you for contributing to Base Linux Setup! 🚀
