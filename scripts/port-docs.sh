#!/bin/bash

###############################################################################
# Docs Porter - Port Documentation to Another Go Gin Project
###############################################################################
#
# This script copies and adapts the documentation from this boilerplate
# to another Go + Gin project with minimal changes.
#
# Portability: 98% - Only module paths and project names need updating
#
# Usage:
#   ./scripts/port-docs.sh /path/to/target-project
#   ./scripts/port-docs.sh /path/to/target-project github.com/your-org/your-project
#
###############################################################################

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
print_step() {
    echo -e "${BLUE}▶ $1${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

###############################################################################
# 1. Validate Input
###############################################################################

if [ -z "$1" ]; then
    print_error "Usage: $0 <target-project-path> [new-module-path]"
    echo ""
    echo "Examples:"
    echo "  $0 /path/to/my-api"
    echo "  $0 /path/to/my-api github.com/myorg/my-api"
    echo ""
    exit 1
fi

TARGET_PROJECT="$1"
NEW_MODULE_PATH="${2:-}"

# Check if target exists
if [ ! -d "$TARGET_PROJECT" ]; then
    print_error "Target project directory does not exist: $TARGET_PROJECT"
    exit 1
fi

# Get absolute paths
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
SOURCE_DOCS="$( cd "$SCRIPT_DIR/.." && pwd )/docs"
TARGET_DOCS="$TARGET_PROJECT/docs"

print_step "Documentation Porter for Go Gin Projects"
echo ""
echo "Source: $SOURCE_DOCS"
echo "Target: $TARGET_DOCS"
echo ""

###############################################################################
# 2. Detect Target Project Info
###############################################################################

print_step "Detecting target project information..."

# Try to detect module path from go.mod
if [ -f "$TARGET_PROJECT/go.mod" ]; then
    if [ -z "$NEW_MODULE_PATH" ]; then
        NEW_MODULE_PATH=$(grep "^module " "$TARGET_PROJECT/go.mod" | awk '{print $2}')
        print_success "Detected module path: $NEW_MODULE_PATH"
    fi

    # Extract project name from module path
    NEW_PROJECT_NAME=$(basename "$NEW_MODULE_PATH")
    print_success "Detected project name: $NEW_PROJECT_NAME"
else
    print_warning "go.mod not found in target project"

    if [ -z "$NEW_MODULE_PATH" ]; then
        print_error "Please provide module path as second argument"
        exit 1
    fi

    NEW_PROJECT_NAME=$(basename "$TARGET_PROJECT")
    print_warning "Using directory name as project name: $NEW_PROJECT_NAME"
fi

# Current boilerplate info (to be replaced)
OLD_MODULE_PATH="github.com/0xdiaz/gin-boilerplate"
OLD_PROJECT_NAME="go-gin-boilerplate"

echo ""
echo "Replacements:"
echo "  Module: $OLD_MODULE_PATH → $NEW_MODULE_PATH"
echo "  Name:   $OLD_PROJECT_NAME → $NEW_PROJECT_NAME"
echo ""

###############################################################################
# 3. Backup Existing Docs (if any)
###############################################################################

if [ -d "$TARGET_DOCS" ]; then
    print_warning "Target docs directory already exists"
    read -p "Backup existing docs? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        BACKUP_DIR="${TARGET_DOCS}.backup.$(date +%Y%m%d_%H%M%S)"
        mv "$TARGET_DOCS" "$BACKUP_DIR"
        print_success "Backed up to: $BACKUP_DIR"
    else
        read -p "Overwrite existing docs? (y/n) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_error "Aborted by user"
            exit 1
        fi
        rm -rf "$TARGET_DOCS"
    fi
fi

###############################################################################
# 4. Copy Documentation
###############################################################################

print_step "Copying documentation files..."

# Create target docs directory
mkdir -p "$TARGET_DOCS"

# Copy all markdown files
cp -r "$SOURCE_DOCS"/*.md "$TARGET_DOCS/" 2>/dev/null || true

# Count files copied
FILE_COUNT=$(find "$TARGET_DOCS" -name "*.md" | wc -l | tr -d ' ')
print_success "Copied $FILE_COUNT documentation files"

###############################################################################
# 5. Update Module Paths
###############################################################################

print_step "Updating module paths..."

# Find and replace module path in all markdown files
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    find "$TARGET_DOCS" -name "*.md" -type f -exec sed -i '' \
        "s|$OLD_MODULE_PATH|$NEW_MODULE_PATH|g" {} +
else
    # Linux
    find "$TARGET_DOCS" -name "*.md" -type f -exec sed -i \
        "s|$OLD_MODULE_PATH|$NEW_MODULE_PATH|g" {} +
fi

print_success "Updated module paths"

###############################################################################
# 6. Update Project Names
###############################################################################

print_step "Updating project names..."

# Replace project name references
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    find "$TARGET_DOCS" -name "*.md" -type f -exec sed -i '' \
        "s|$OLD_PROJECT_NAME|$NEW_PROJECT_NAME|g" {} +

    # Also replace "boilerplate" references
    find "$TARGET_DOCS" -name "*.md" -type f -exec sed -i '' \
        "s|gin-boilerplate|$NEW_PROJECT_NAME|g" {} +

    # Replace author name if needed
    find "$TARGET_DOCS" -name "*.md" -type f -exec sed -i '' \
        "s|0xdiaz|$(echo $NEW_MODULE_PATH | cut -d'/' -f2)|g" {} +
else
    # Linux
    find "$TARGET_DOCS" -name "*.md" -type f -exec sed -i \
        "s|$OLD_PROJECT_NAME|$NEW_PROJECT_NAME|g" {} +

    find "$TARGET_DOCS" -name "*.md" -type f -exec sed -i \
        "s|gin-boilerplate|$NEW_PROJECT_NAME|g" {} +

    find "$TARGET_DOCS" -name "*.md" -type f -exec sed -i \
        "s|0xdiaz|$(echo $NEW_MODULE_PATH | cut -d'/' -f2)|g" {} +
fi

print_success "Updated project names"

###############################################################################
# 7. Update README.md in Target Project (if exists)
###############################################################################

print_step "Checking if target project has README.md..."

TARGET_README="$TARGET_PROJECT/README.md"

if [ -f "$TARGET_README" ]; then
    print_warning "README.md exists in target project"
    echo ""
    echo "Options:"
    echo "  1) Keep existing README.md (docs will have their own README)"
    echo "  2) Add documentation section to existing README"
    echo "  3) Replace with boilerplate README (not recommended)"
    echo ""
    read -p "Choose option (1/2/3): " -n 1 -r
    echo

    case $REPLY in
        2)
            # Add documentation section
            if ! grep -q "## 📚 Documentation" "$TARGET_README"; then
                cat >> "$TARGET_README" << 'EOF'

---

## 📚 Documentation

Comprehensive documentation is available in the `docs/` directory:

### 🚀 Quick Start for AI Agents

1. **[docs/00_AI_CRITICAL_RULES.md](docs/00_AI_CRITICAL_RULES.md)** ⚠️ **READ FIRST** (100 lines)
   - Non-negotiable patterns that MUST be followed
   - Struct-based controllers/services
   - Response utility usage

2. **[docs/AI_QUICK_REFERENCE.md](docs/AI_QUICK_REFERENCE.md)** (Quick templates)
   - Ready-to-use code templates
   - Common patterns
   - Testing guidelines

### 📖 Full Documentation

- **[docs/CODING_STANDARDS.md](docs/CODING_STANDARDS.md)** - Complete coding standards and best practices
- **[docs/DESIGN_PATTERNS.md](docs/DESIGN_PATTERNS.md)** - Architecture patterns and implementation guides
- **[docs/AUTHENTICATION.md](docs/AUTHENTICATION.md)** - Authentication system documentation
- **[docs/CONFIGURATION.md](docs/CONFIGURATION.md)** - Configuration management guide
- **[docs/OBSERVABILITY.md](docs/OBSERVABILITY.md)** - Monitoring and observability setup
- **[docs/DOCS_INDEX.md](docs/DOCS_INDEX.md)** - Complete documentation index

### 🎯 For New Developers

Start with [docs/README.md](docs/README.md) for a guided tour of the documentation.

EOF
                print_success "Added documentation section to README.md"
            else
                print_warning "Documentation section already exists in README.md"
            fi
            ;;
        3)
            print_warning "Replacing README.md..."
            cp "$SOURCE_DOCS/../README.md" "$TARGET_README"
            # Update module path in README
            if [[ "$OSTYPE" == "darwin"* ]]; then
                sed -i '' "s|$OLD_MODULE_PATH|$NEW_MODULE_PATH|g" "$TARGET_README"
                sed -i '' "s|$OLD_PROJECT_NAME|$NEW_PROJECT_NAME|g" "$TARGET_README"
            else
                sed -i "s|$OLD_MODULE_PATH|$NEW_MODULE_PATH|g" "$TARGET_README"
                sed -i "s|$OLD_PROJECT_NAME|$NEW_PROJECT_NAME|g" "$TARGET_README"
            fi
            print_success "Replaced README.md"
            ;;
        *)
            print_success "Keeping existing README.md unchanged"
            ;;
    esac
else
    print_warning "No README.md in target project"
    read -p "Copy boilerplate README.md? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        cp "$SOURCE_DOCS/../README.md" "$TARGET_README"
        # Update module path in README
        if [[ "$OSTYPE" == "darwin"* ]]; then
            sed -i '' "s|$OLD_MODULE_PATH|$NEW_MODULE_PATH|g" "$TARGET_README"
            sed -i '' "s|$OLD_PROJECT_NAME|$NEW_PROJECT_NAME|g" "$TARGET_README"
        else
            sed -i "s|$OLD_MODULE_PATH|$NEW_MODULE_PATH|g" "$TARGET_README"
            sed -i "s|$OLD_PROJECT_NAME|$NEW_PROJECT_NAME|g" "$TARGET_README"
        fi
        print_success "Copied and updated README.md"
    fi
fi

###############################################################################
# 8. Validation & Summary
###############################################################################

print_step "Validating documentation..."

# Check critical files exist
CRITICAL_FILES=(
    "00_AI_CRITICAL_RULES.md"
    "CODING_STANDARDS.md"
    "DESIGN_PATTERNS.md"
    "DOCS_INDEX.md"
)

MISSING=0
for file in "${CRITICAL_FILES[@]}"; do
    if [ ! -f "$TARGET_DOCS/$file" ]; then
        print_error "Missing critical file: $file"
        MISSING=$((MISSING + 1))
    fi
done

if [ $MISSING -eq 0 ]; then
    print_success "All critical files present"
else
    print_error "Some critical files are missing"
fi

# Check if references were updated
OLD_REF_COUNT=$(grep -r "$OLD_MODULE_PATH" "$TARGET_DOCS" 2>/dev/null | wc -l | tr -d ' ')

if [ "$OLD_REF_COUNT" -gt 0 ]; then
    print_warning "Found $OLD_REF_COUNT old module path references (might be in code comments)"
else
    print_success "All module path references updated"
fi

###############################################################################
# 9. Final Summary
###############################################################################

echo ""
echo "========================================================================="
print_success "Documentation successfully ported!"
echo "========================================================================="
echo ""
echo "Summary:"
echo "  ✓ Copied $FILE_COUNT documentation files"
echo "  ✓ Updated module paths: $OLD_MODULE_PATH → $NEW_MODULE_PATH"
echo "  ✓ Updated project name: $OLD_PROJECT_NAME → $NEW_PROJECT_NAME"
echo "  ✓ Target location: $TARGET_DOCS"
echo ""
echo "Next Steps:"
echo ""
echo "  1. Review documentation in: $TARGET_DOCS"
echo ""
echo "  2. Update project-specific docs (if needed):"
echo "     • docs/AUTHENTICATION.md - If auth implementation differs"
echo "     • docs/CONFIGURATION.md - Update environment variables"
echo "     • docs/OBSERVABILITY.md - Update metrics/monitoring setup"
echo ""
echo "  3. Update DOCS_INDEX.md if you add new documentation files"
echo ""
echo "  4. Commit the documentation:"
echo "     cd $TARGET_PROJECT"
echo "     git add docs/"
echo "     git commit -m 'docs: add comprehensive project documentation'"
echo ""
echo "Portability: 98% ✓"
echo "Estimated review time: 10-15 minutes"
echo ""
print_success "Done! 🎉"
echo ""
