#!/bin/bash

# Base Linux Setup - GitHub Repository Setup Script
# This script helps configure the repository for public release on GitHub

set -e

REPO_URL="https://github.com/GuilhermeVozniak/base-linux-setup"
REPO_NAME="base-linux-setup"

echo "🚀 Base Linux Setup - GitHub Repository Setup"
echo "=============================================="
echo

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
print_step() {
    echo -e "${BLUE}📋 $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if we're in the right directory
check_directory() {
    print_step "Checking project directory..."
    
    if [ ! -f "main.go" ] || [ ! -f "go.mod" ]; then
        print_error "This doesn't appear to be the base-linux-setup project directory"
        print_error "Please run this script from the project root"
        exit 1
    fi
    
    if [ ! -d ".git" ]; then
        print_error "This is not a git repository"
        print_error "Please initialize git first: git init"
        exit 1
    fi
    
    print_success "Project directory confirmed"
}

# Check prerequisites
check_prerequisites() {
    print_step "Checking prerequisites..."
    
    # Check for required tools
    for tool in git go make; do
        if ! command -v $tool &> /dev/null; then
            print_error "$tool is not installed"
            exit 1
        fi
    done
    
    # Check Go version
    GO_VERSION=$(go version | grep -oE 'go[0-9]+\.[0-9]+' | sed 's/go//')
    if [ "$(printf '%s\n' "1.21" "$GO_VERSION" | sort -V | head -n1)" != "1.21" ]; then
        print_warning "Go version $GO_VERSION detected. Recommended: 1.21+"
    fi
    
    print_success "Prerequisites verified"
}

# Setup git configuration
setup_git() {
    print_step "Configuring git repository..."
    
    # Check if remote origin exists
    if ! git remote get-url origin &> /dev/null; then
        read -p "Enter your GitHub repository URL (e.g., https://github.com/username/repo.git): " REPO_URL
        git remote add origin "$REPO_URL"
        print_success "Added remote origin: $REPO_URL"
    else
        print_success "Remote origin already configured"
    fi
    
    # Set up git hooks directory
    mkdir -p .git/hooks
    
    # Create pre-commit hook
    cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash
# Pre-commit hook for Base Linux Setup

echo "Running pre-commit checks..."

# Check Go formatting
if ! gofmt -l . | grep -q '^$'; then
    echo "Code is not properly formatted. Run 'make fmt'"
    exit 1
fi

# Run tests
if ! make test; then
    echo "Tests failed"
    exit 1
fi

echo "Pre-commit checks passed ✅"
EOF
    
    chmod +x .git/hooks/pre-commit
    print_success "Git hooks configured"
}

# Build the project
build_project() {
    print_step "Building project..."
    
    if ! make build; then
        print_error "Build failed"
        exit 1
    fi
    
    if ! make test; then
        print_error "Tests failed"
        exit 1
    fi
    
    print_success "Project built and tested successfully"
}

# Create initial release
create_release() {
    print_step "Preparing for initial release..."
    
    # Check if we have any tags
    if git tag -l | grep -q .; then
        print_warning "Tags already exist. Skipping tag creation."
        return
    fi
    
    # Get current version from main.go or set default
    VERSION="v0.1.0"
    read -p "Enter version for initial release [$VERSION]: " input_version
    if [ -n "$input_version" ]; then
        VERSION="$input_version"
    fi
    
    # Ensure version starts with 'v'
    if [[ ! $VERSION == v* ]]; then
        VERSION="v$VERSION"
    fi
    
    # Create and push tag
    git tag -a "$VERSION" -m "Initial release $VERSION"
    print_success "Created tag: $VERSION"
    
    echo
    print_warning "To trigger the automated release:"
    print_warning "1. Commit and push your changes: git push origin main"
    print_warning "2. Push the tag: git push origin $VERSION"
    print_warning "3. GitHub Actions will automatically create the release"
}

# Setup GitHub Actions
setup_github_actions() {
    print_step "Verifying GitHub Actions setup..."
    
    if [ ! -f ".github/workflows/release.yml" ]; then
        print_error "Release workflow not found"
        print_error "Make sure .github/workflows/release.yml exists"
        exit 1
    fi
    
    if [ ! -f ".github/workflows/ci.yml" ]; then
        print_error "CI workflow not found"
        print_error "Make sure .github/workflows/ci.yml exists"
        exit 1
    fi
    
    print_success "GitHub Actions workflows configured"
}

# Setup GitHub repository settings
setup_github_repository() {
    print_step "GitHub repository setup checklist..."
    
    echo
    echo "Please ensure the following settings in your GitHub repository:"
    echo "1. 📋 Repository Settings:"
    echo "   - Description: 'CLI tool for automated Linux environment setup'"
    echo "   - Website: (optional)"
    echo "   - Topics: linux, automation, cli, golang, raspberry-pi, kali-linux"
    echo
    echo "2. 🔧 Features:"
    echo "   - ✅ Wikis (enable for documentation)"
    echo "   - ✅ Issues (enable for bug reports)"
    echo "   - ✅ Discussions (enable for community)"
    echo
    echo "3. 🏷️ Releases:"
    echo "   - Will be created automatically by GitHub Actions"
    echo "   - Make sure repository has push access to create releases"
    echo
    echo "4. 🔒 Security:"
    echo "   - Enable 'Require status checks' for main branch"
    echo "   - Enable 'Require pull request reviews'"
    echo
}

# Setup wiki
setup_wiki() {
    print_step "Wiki setup instructions..."
    
    echo
    echo "To set up the project wiki:"
    echo "1. Enable Wiki in repository settings"
    echo "2. Go to the Wiki tab in your repository"
    echo "3. Import the documentation from wiki/ directory:"
    echo "   - Copy content from wiki/Home.md to the Home page"
    echo "   - Create 'Installation' page from wiki/Installation.md"
    echo "   - Create 'Usage Guide' page from wiki/Usage-Guide.md"
    echo "   - Create 'Preset Development' page from wiki/Preset-Development.md"
    echo
    echo "Or use the git method:"
    echo "   git clone https://github.com/YOUR_USERNAME/$REPO_NAME.wiki.git"
    echo "   cp wiki/*.md $REPO_NAME.wiki/"
    echo "   cd $REPO_NAME.wiki && git add . && git commit -m 'Initial wiki' && git push"
    echo
}

# Main execution
main() {
    echo "This script will help you set up the Base Linux Setup repository for GitHub."
    echo "It will configure git, build the project, and prepare for release."
    echo
    read -p "Continue? (y/N): " -n 1 -r
    echo
    
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Setup cancelled."
        exit 0
    fi
    
    check_directory
    check_prerequisites
    setup_git
    build_project
    setup_github_actions
    create_release
    setup_github_repository
    setup_wiki
    
    echo
    print_success "GitHub setup completed!"
    echo
    echo "🎉 Next steps:"
    echo "1. Review and commit any remaining changes"
    echo "2. Push to GitHub: git push origin main"
    echo "3. Push the tag: git push origin \$(git describe --tags)"
    echo "4. Set up repository settings as described above"
    echo "5. Import wiki documentation"
    echo "6. Configure branch protection rules"
    echo
    echo "🚀 Your project is ready for the public!"
}

# Run main function
main "$@" 