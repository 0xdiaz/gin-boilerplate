# 📦 Documentation Porting Guide

**How to use this boilerplate's documentation in your Go + Gin projects**

---

## 🎯 TL;DR

```bash
# One command to port all documentation:
./scripts/port-docs.sh /path/to/your-project

# Done! 98% ready in ~5 seconds
```

---

## 📊 What Gets Ported?

### ✅ All Documentation Files

```
docs/
├── 00_AI_CRITICAL_RULES.md       ✅ 100% portable
├── AI_AGENT_RULES.md             ✅ 100% portable
├── AI_QUICK_REFERENCE.md         ✅ 100% portable
├── CODING_STANDARDS.md           ✅ 98% portable (auto-updated)
├── DESIGN_PATTERNS.md            ✅ 100% portable
├── DOCS_INDEX.md                 ✅ 90% portable (may need file list update)
├── README.md                     ✅ 100% portable
├── AUTHENTICATION.md             ✅ 95% portable (may need env var update)
├── CONFIGURATION.md              ✅ 90% portable (update for your project)
└── OBSERVABILITY.md              ✅ 90% portable (update for your project)
```

**Total:** 10 comprehensive documentation files (~8,500+ lines)

---

## 🚀 Quick Start

### Step 1: Run the Script

```bash
# Navigate to boilerplate directory
cd /path/to/go-gin-boilerplate

# Run port script
./scripts/port-docs.sh /path/to/your-project
```

### Step 2: Follow the Prompts

The script will:
1. ✅ Auto-detect your project's module path from `go.mod`
2. ✅ Copy all documentation files
3. ✅ Update module paths automatically
4. ✅ Update project names automatically
5. ✅ Ask if you want to update README.md

### Step 3: Quick Review (10 minutes)

Review these files for project-specific updates:

```bash
cd /path/to/your-project/docs

# Check these files (may need minor edits):
# - CONFIGURATION.md (env variables)
# - AUTHENTICATION.md (if auth differs)
# - OBSERVABILITY.md (metrics endpoints)
```

### Step 4: Commit

```bash
git add docs/
git commit -m "docs: add comprehensive project documentation"
```

**Done! 🎉**

---

## 🔧 Using this repo as a template (changing the Go module path)

When you **fork or clone this boilerplate** to start a new project under your own org/repo, you must change the Go module path so imports and `go build` use your module name. Do this **once** when creating the new project.

### Steps

1. **Clone or fork the repo** into your project directory.

2. **Update the module path in go.mod:**
   ```bash
   go mod edit -module=github.com/your-org/your-repo
   ```
   Replace `github.com/your-org/your-repo` with your actual module path.

3. **Replace the old import path** across the codebase:
   - Old path: `github.com/0xdiaz/gin-boilerplate`
   - New path: your module path from step 2  
   Use your editor’s find-replace (or a script) over all `.go` files (and any config/docs that reference the path).

4. **Tidy and build:**
   ```bash
   go mod tidy
   go build ./...
   ```

5. **Commit** the updated `go.mod`, `go.sum`, and all modified files.

After this, the project builds and runs under your module path. You only need to do this when creating a new project from this boilerplate.

---

## 📖 Detailed Usage

### Scenario 1: New Project from Scratch

```bash
# Create new project
mkdir my-awesome-api
cd my-awesome-api
go mod init github.com/myorg/my-awesome-api

# Create basic structure
mkdir -p internal/app/{controllers,services,dto}
mkdir -p internal/domain/{models,repositories}
mkdir -p pkg/{config,logger,utils}

# Port documentation
cd /path/to/boilerplate
./scripts/port-docs.sh ../my-awesome-api

# Result: Full docs in ~/my-awesome-api/docs/
```

**Time:** 5 seconds (script) + 10 minutes (review) = **~10 minutes total**

---

### Scenario 2: Existing Project (Add Documentation)

```bash
# Your existing project
cd /path/to/existing-project

# Port docs from boilerplate
/path/to/boilerplate/scripts/port-docs.sh .

# Script will:
# - Backup existing docs (if any)
# - Copy new documentation
# - Update paths to match your project
# - Offer to update README.md
```

**Time:** 5 seconds (script) + 15 minutes (review) = **~15 minutes total**

---

### Scenario 3: Company-Wide Template

```bash
# Create company template
mkdir /company/templates/go-gin-standard
cd /company/templates/go-gin-standard

# Copy boilerplate structure
cp -r /path/to/boilerplate/* .

# Port docs
./scripts/port-docs.sh . github.com/mycompany/service-template

# All teams can now use this template
```

**Result:** Consistent documentation across all company projects! 🏢

---

## 🔄 What Gets Updated Automatically

### Module Paths

```go
// BEFORE (in boilerplate)
import "github.com/0xdiaz/gin-boilerplate/internal/app/dto"

// AFTER (in your project)
import "github.com/myorg/my-api/internal/app/dto"
```

### Project Names

```markdown
<!-- BEFORE -->
# Go Gin Enterprise Boilerplate

<!-- AFTER -->
# My Awesome API
```

### Code Examples

```go
// BEFORE
package main

import (
    "github.com/0xdiaz/gin-boilerplate/pkg/config"
)

// AFTER
package main

import (
    "github.com/myorg/my-api/pkg/config"
)
```

---

## ✏️ What You Need to Update Manually (Optional)

### 1. CONFIGURATION.md (~5 minutes)

Update environment variables section if your project uses different config:

```bash
# Your project uses different env vars?
vim docs/CONFIGURATION.md

# Update:
# - Environment variable names
# - Default values
# - Configuration examples
```

### 2. AUTHENTICATION.md (~3 minutes)

Update if you customize auth:

```bash
# Different JWT expiry? Different tokens?
vim docs/AUTHENTICATION.md

# Update:
# - Token expiry times
# - Authentication endpoints (if different)
# - Security features (if added/removed)
```

### 3. OBSERVABILITY.md (~2 minutes)

Update metrics/health check endpoints:

```bash
vim docs/OBSERVABILITY.md

# Update:
# - Metrics endpoint paths
# - Health check indicators
# - Monitoring setup
```

---

## 📋 Verification Checklist

After porting, verify:

```bash
cd /path/to/your-project/docs

# ✅ Check critical files exist
[ -f "00_AI_CRITICAL_RULES.md" ] && echo "✓ Critical rules"
[ -f "CODING_STANDARDS.md" ] && echo "✓ Coding standards"
[ -f "DESIGN_PATTERNS.md" ] && echo "✓ Design patterns"
[ -f "DOCS_INDEX.md" ] && echo "✓ Docs index"

# ✅ Check no old references remain
grep -r "gin-boilerplate" . && echo "⚠️  Found old references" || echo "✓ No old references"
grep -r "0xdiaz" . && echo "⚠️  Found old author" || echo "✓ No old author refs"

# ✅ Check module path updated
grep -r "github.com/myorg/my-api" . && echo "✓ Module path updated" || echo "⚠️  Module path not found"
```

---

## 🎯 Benefits of Porting

### For Your Project

- 📚 **8,500+ lines** of comprehensive documentation
- 🤖 **AI-friendly** - Optimized for AI-assisted development
- 📏 **Standards enforced** - File size limits, function limits, testing requirements
- 🏗️ **Architecture patterns** - Clean architecture, repository pattern, DI
- 🧪 **Testing guides** - Unit test examples and patterns
- 🔒 **Security practices** - SQL injection prevention, password hashing, token security

### For Your Team

- ✅ **Onboarding speed** - New devs productive in hours, not days
- ✅ **Code consistency** - Everyone follows same patterns
- ✅ **Quality assurance** - Built-in standards prevent common mistakes
- ✅ **Knowledge sharing** - Documentation serves as training material
- ✅ **Reduced reviews** - Standards documented, less back-and-forth

### ROI

**Time Investment:**
- Initial port: 10-15 minutes
- Occasional updates: 5 minutes

**Time Saved:**
- Per new developer: 2-3 days (no need to learn codebase from scratch)
- Per feature: 30-60 minutes (patterns already documented)
- Code reviews: 50% reduction (standards clear)

**Annual savings for a 5-person team:** ~40-60 developer days! 🚀

---

## 💡 Pro Tips

### Tip 1: Keep Docs in Sync

When boilerplate docs get updated:

```bash
# Pull latest docs from boilerplate
cd /path/to/boilerplate
git pull

# Re-run port script
./scripts/port-docs.sh /path/to/your-project

# Review changes
cd /path/to/your-project
git diff docs/
```

### Tip 2: Customize for Your Domain

Add domain-specific docs:

```bash
cd /path/to/your-project/docs

# Add your own documentation
echo "# Payment Processing" > PAYMENTS.md
echo "# Order Management" > ORDERS.md

# Update DOCS_INDEX.md to include them
vim DOCS_INDEX.md
```

### Tip 3: Share Across Microservices

For microservices architecture:

```bash
# Port to each microservice
for service in auth-service payment-service order-service; do
    echo "Porting docs to $service..."
    ./scripts/port-docs.sh /services/$service
done

# All services now have same standards!
```

### Tip 4: Create Pull Request Template

Add docs reference to PR template:

```markdown
## Checklist

- [ ] Code follows [docs/CODING_STANDARDS.md](docs/CODING_STANDARDS.md)
- [ ] Architecture follows [docs/DESIGN_PATTERNS.md](docs/DESIGN_PATTERNS.md)
- [ ] Tests added (see [docs/CODING_STANDARDS.md](docs/CODING_STANDARDS.md#testing))
- [ ] Response utilities used (see [docs/00_AI_CRITICAL_RULES.md](docs/00_AI_CRITICAL_RULES.md))
```

---

## 🆘 Troubleshooting

### Issue: "Permission denied"

```bash
# Solution: Make script executable
chmod +x scripts/port-docs.sh
```

### Issue: "Old references still present"

```bash
# Solution: Re-run with correct module path
./scripts/port-docs.sh /path/to/project github.com/correct/module/path
```

### Issue: "Docs don't match project structure"

This is normal! The docs describe **patterns and principles**, not specific implementation.

Your project structure can differ (e.g., different folder names), but patterns remain same:
- Controller → Service → Repository flow
- Struct-based approach
- DTO usage
- Error handling

---

## 📊 Success Stories

### Before Porting Docs

```
❌ Inconsistent code styles across team
❌ New devs take 1-2 weeks to be productive
❌ Code reviews take 2-3 days
❌ Frequent back-and-forth on patterns
❌ No clear architecture guidelines
```

### After Porting Docs

```
✅ Consistent code following standards
✅ New devs productive in 1-2 days
✅ Code reviews take hours, not days
✅ Clear reference for all patterns
✅ Architecture clearly documented
```

---

## 🎓 Learning Path for New Devs

### Day 1: Quick Start (1 hour)

```bash
# Read critical rules
docs/00_AI_CRITICAL_RULES.md          # 5 minutes

# Read quick reference
docs/AI_QUICK_REFERENCE.md            # 10 minutes

# Scan coding standards TOC
docs/CODING_STANDARDS.md              # 5 minutes (just TOC)

# Review one complete example
docs/DESIGN_PATTERNS.md               # 30 minutes (Step-by-step guide)
```

**Result:** Ready to write first controller/service! ✅

### Week 1: Deep Dive (3-4 hours)

```bash
# Read sections as needed:
docs/CODING_STANDARDS.md              # Read specific sections when implementing
docs/DESIGN_PATTERNS.md               # Reference patterns as you encounter them
docs/AUTHENTICATION.md                # When working on auth features
```

**Result:** Confident in all patterns! ✅

### Month 1: Expert (Reference as needed)

```bash
# Docs become reference material
# Look up specific topics when needed
# Contribute improvements to docs
```

**Result:** Team expert, mentoring others! ✅

---

## 📈 Metrics

Track documentation effectiveness:

```markdown
## Documentation Metrics

- Time to first commit (new dev): _____ hours
- Code review turnaround: _____ hours
- Standards violations per PR: _____
- Developer satisfaction (1-10): _____

**Target:**
- Time to first commit: < 4 hours
- Code review turnaround: < 24 hours
- Standards violations: < 2 per PR
- Developer satisfaction: > 8/10
```

---

## 🔗 Additional Resources

- [Main Boilerplate README](README.md)
- [Contributing Guidelines](CONTRIBUTING.md) _(if exists)_
- [Docs Index](docs/DOCS_INDEX.md)
- [Script Documentation](scripts/README.md)

---

## 🎬 Video Tutorial (Coming Soon)

Watch a 5-minute video showing:
- How to run the port script
- What gets updated automatically
- How to customize for your project
- Real-world examples

---

## ❓ FAQ

**Q: Can I customize the docs after porting?**
A: Yes! They're yours now. Add, modify, remove as needed.

**Q: Will updates from boilerplate break my customizations?**
A: No. After porting, docs are independent. Only re-port if you want updates.

**Q: Can I use this for non-Gin Go projects?**
A: Yes! 90% of docs are framework-agnostic. Just update framework-specific examples.

**Q: What license are the docs under?**
A: Same as boilerplate (MIT). Use freely!

---

## 🤝 Contributing

Found improvements to the porting process?

1. Fork the boilerplate repo
2. Improve `scripts/port-docs.sh`
3. Update this guide
4. Submit PR

Help make porting even easier! 🙌

---

**Happy porting! 🚀**

Questions? Issues? Open a GitHub issue on the boilerplate repo.
