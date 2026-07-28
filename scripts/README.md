# Scripts

Utility scripts for the Go Gin Boilerplate project.

---

## 📦 port-docs.sh

**Purpose:** Port documentation to another Go + Gin project

### Quick Start

```bash
# Basic usage (auto-detect module path from target's go.mod)
./scripts/port-docs.sh /path/to/target-project

# Specify module path explicitly
./scripts/port-docs.sh /path/to/target-project github.com/myorg/my-api

# Example
./scripts/port-docs.sh ../my-new-api
```

### What It Does

1. ✅ Copies all documentation files from `docs/` to target project
2. ✅ Auto-detects target module path from `go.mod`
3. ✅ Updates all module path references
4. ✅ Updates project name references
5. ✅ Optionally updates target README.md with docs section
6. ✅ Validates all critical files copied correctly

### Portability

**98% portable** - Only module paths and project names are updated.

All architecture patterns, coding standards, and best practices remain unchanged.

### Time Required

- **Script execution:** ~5 seconds
- **Manual review:** ~10-15 minutes

### Files Updated

The script automatically updates:
- Module imports: `github.com/old/path` → `github.com/new/path`
- Project names: `go-gin-boilerplate` → `your-project-name`
- Author references: `0xdiaz` → `your-github-username`

### Safety Features

- Backs up existing docs before overwriting
- Confirms before replacing README.md
- Validates all critical files copied
- Shows summary of changes made

### Example Output

```
▶ Documentation Porter for Go Gin Projects

Source: /path/to/boilerplate/docs
Target: /path/to/my-api/docs

▶ Detecting target project information...
✓ Detected module path: github.com/myorg/my-api
✓ Detected project name: my-api

Replacements:
  Module: github.com/0xdiaz/gin-boilerplate → github.com/myorg/my-api
  Name:   go-gin-boilerplate → my-api

▶ Copying documentation files...
✓ Copied 10 documentation files

▶ Updating module paths...
✓ Updated module paths

▶ Updating project names...
✓ Updated project names

▶ Validating documentation...
✓ All critical files present
✓ All module path references updated

=========================================================================
✓ Documentation successfully ported!
=========================================================================

Portability: 98% ✓
Estimated review time: 10-15 minutes

Done! 🎉
```

### Post-Port Checklist

After running the script, review these files:

1. **`docs/AUTHENTICATION.md`** - Update if auth implementation differs
2. **`docs/CONFIGURATION.md`** - Update environment variables for your project
3. **`docs/OBSERVABILITY.md`** - Update metrics/monitoring endpoints
4. **`docs/DOCS_INDEX.md`** - Update if you add new documentation files

### Common Use Cases

#### Use Case 1: New Microservice

```bash
# Create new microservice project
mkdir my-payment-service
cd my-payment-service
go mod init github.com/myorg/payment-service

# Port documentation
cd /path/to/boilerplate
./scripts/port-docs.sh ../my-payment-service

# Done! Documentation ready
```

#### Use Case 2: Existing Project (Add Docs)

```bash
# Port to existing project
./scripts/port-docs.sh /path/to/existing-project

# Backs up existing docs automatically
# Updates README.md with docs section
```

#### Use Case 3: Internal Template

```bash
# Create company template
./scripts/port-docs.sh /company/templates/go-gin-template

# All devs can use this template
# Consistent docs across all projects
```

### Troubleshooting

**Issue:** "Target project directory does not exist"
```bash
# Solution: Create directory first
mkdir -p /path/to/target-project
```

**Issue:** "go.mod not found"
```bash
# Solution: Provide module path explicitly
./scripts/port-docs.sh /path/to/target github.com/myorg/myproject
```

**Issue:** "Permission denied"
```bash
# Solution: Make script executable
chmod +x scripts/port-docs.sh
```

### Advanced Usage

#### Dry Run (Preview Changes)

```bash
# Copy to temp location first
./scripts/port-docs.sh /tmp/preview-docs github.com/myorg/myproject

# Review changes
cd /tmp/preview-docs
cat *.md

# If satisfied, run on actual project
```

#### Batch Port (Multiple Projects)

```bash
# Port to multiple projects
for project in ../project-a ../project-b ../project-c; do
    echo "Porting to $project..."
    ./scripts/port-docs.sh "$project"
done
```

#### Custom Replacements

After running the script, you can add custom find/replace:

```bash
# Additional customizations
cd /path/to/target-project/docs
find . -name "*.md" -exec sed -i 's/old-value/new-value/g' {} +
```

---

## 🛠️ Future Scripts (Coming Soon)

- `setup-project.sh` - Complete project setup from boilerplate
- `update-docs.sh` - Pull latest docs updates from boilerplate
- `validate-standards.sh` - Check project compliance with coding standards
- `generate-changelog.sh` - Auto-generate changelog from commits

---

## 📚 Related Documentation

- [Main README](../README.md) - Project overview
- [Coding Standards](../docs/CODING_STANDARDS.md) - Code quality guidelines
- [Design Patterns](../docs/DESIGN_PATTERNS.md) - Architecture patterns

---

**Built with ❤️ for the Go community**
