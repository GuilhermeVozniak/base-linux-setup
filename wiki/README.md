# Wiki Documentation

This directory contains the documentation pages for the Base Linux Setup project wiki. These markdown files are designed to be imported into GitHub Wiki.

## Files Overview

| File | Description | Wiki Page |
|------|-------------|-----------|
| `Home.md` | Main wiki landing page with overview and navigation | **Home** |
| `Installation.md` | Complete installation guide for all platforms | **Installation** |
| `Usage-Guide.md` | Comprehensive usage guide and examples | **Usage Guide** |
| `Preset-Development.md` | Guide for creating custom presets | **Preset Development** |

## How to Import to GitHub Wiki

### Method 1: Manual Import (Recommended)

1. **Enable Wiki** for your GitHub repository:
   - Go to repository Settings
   - Scroll to "Features" section
   - Check "Wikis"

2. **Create Wiki Pages**:
   - Navigate to the Wiki tab in your repository
   - Click "Create the first page"
   - For each file in this directory:
     - Create a new page with the corresponding name
     - Copy the content from the `.md` file
     - Save the page

3. **Page Mapping**:
   ```
   Home.md → Home (default page)
   Installation.md → Installation
   Usage-Guide.md → Usage Guide
   Preset-Development.md → Preset Development
   ```

### Method 2: Git Clone (Advanced)

GitHub Wiki is actually a git repository that you can clone and manage:

```bash
# Clone the wiki repository
git clone https://github.com/GuilhermeVozniak/base-linux-setup.wiki.git

# Copy files from this directory
cp wiki/*.md base-linux-setup.wiki/

# Rename files to match GitHub Wiki format
cd base-linux-setup.wiki
mv Usage-Guide.md "Usage-Guide.md"
mv Preset-Development.md "Preset-Development.md"

# Commit and push
git add .
git commit -m "Import initial wiki documentation"
git push origin master
```

### Method 3: Automated Script

Use the provided `setup-github.sh` script which will help set up the repository and wiki.

## Wiki Structure

The wiki follows this navigation structure:

```
Home
├── Getting Started
│   ├── Installation
│   ├── Usage Guide
│   └── Configuration
├── Development
│   ├── Development Guide
│   ├── Preset Development
│   └── API Reference
└── Support
    ├── Troubleshooting
    ├── FAQ
    └── Contributing
```

## Linking Between Pages

Use GitHub Wiki's double bracket syntax for internal links:

```markdown
[[Installation]]           # Links to Installation page
[[Usage Guide]]           # Links to Usage Guide page
[[Preset Development]]    # Links to Preset Development page
```

## External Links

For links to GitHub issues, PRs, etc., use full URLs:

```markdown
[Create an issue](https://github.com/GuilhermeVozniak/base-linux-setup/issues/new)
[Latest release](https://github.com/GuilhermeVozniak/base-linux-setup/releases/latest)
```

## Maintenance

### Updating Wiki Content

1. **Update source files** in this `wiki/` directory
2. **Copy changes** to GitHub Wiki pages
3. **Keep both in sync** for consistency

### Adding New Pages

1. **Create `.md` file** in this directory
2. **Add to this README** in the files overview table
3. **Import to GitHub Wiki** using one of the methods above
4. **Update navigation** in Home.md if needed

## Content Guidelines

### Writing Style
- Use clear, concise language
- Include code examples for technical concepts
- Provide step-by-step instructions
- Use consistent formatting

### Structure
- Start with a brief overview
- Use hierarchical headings (H1, H2, H3)
- Include navigation links at the end
- Add troubleshooting sections where appropriate

### Code Examples
- Use appropriate syntax highlighting
- Include full command examples
- Show expected output where helpful
- Provide multiple options when applicable

## Contributing to Wiki

1. **Edit source files** in this directory
2. **Submit PR** with your changes
3. **Update GitHub Wiki** after PR is merged
4. **Link to related issues/PRs** in commit messages

## Templates

### New Page Template

```markdown
# Page Title

Brief description of what this page covers.

## Section 1

Content here...

## Section 2

More content...

## Troubleshooting

Common issues and solutions...

## See Also

- [[Related Page 1]]
- [[Related Page 2]]
- [External Link](https://example.com)

---

*Last updated: [Date]*
```

### Code Block Template

```markdown
### Task Description

```bash
# Comments explaining the command
command-to-run --with-options

# Expected output
Output example here
``` 