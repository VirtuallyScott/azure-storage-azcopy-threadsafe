# Git Flow Configuration for azcopy Thread-Safe Edition

## Branch Structure

This repository follows **Git Flow** methodology for organized development and releases:

### Main Branches
- **`main`** - Production-ready code, tagged releases
- **`develop`** - Integration branch for ongoing development

### Supporting Branches
- **`feature/*`** - Feature development branches
- **`release/*`** - Release preparation branches  
- **`hotfix/*`** - Emergency fixes to production
- **`support/*`** - Long-term support branches

## Workflow Commands

### Starting New Work

```bash
# Start a new feature
git flow feature start feature-name
# Example: git flow feature start process-locking-improvements

# Start a new release
git flow release start 1.2.0

# Start a hotfix
git flow hotfix start critical-bug-fix
```

### Finishing Work

```bash
# Finish a feature (merges to develop)
git flow feature finish feature-name

# Finish a release (merges to main and develop, creates tag)
git flow release finish 1.2.0

# Finish a hotfix (merges to main and develop, creates tag)
git flow hotfix finish critical-bug-fix
```

### Publishing Branches

```bash
# Publish feature branch for collaboration
git flow feature publish feature-name

# Track someone else's feature
git flow feature track feature-name
```

## Development Workflow

### For New Features

1. **Start from develop**: `git flow feature start my-feature`
2. **Develop**: Make your changes, commit regularly
3. **Test**: Ensure all tests pass and no regressions
4. **Finish**: `git flow feature finish my-feature`
5. **Push**: `git push origin develop`

### For Releases

1. **Start release**: `git flow release start 1.x.x`
2. **Prepare**: Update version numbers, documentation
3. **Test**: Final integration testing
4. **Finish**: `git flow release finish 1.x.x`
5. **Push**: `git push origin main && git push origin develop && git push --tags`

### For Hotfixes

1. **Start from main**: `git flow hotfix start fix-name`
2. **Fix**: Make minimal changes to address the urgent issue
3. **Test**: Thorough testing of the fix
4. **Finish**: `git flow hotfix finish fix-name`
5. **Push**: `git push origin main && git push origin develop && git push --tags`

## GitHub Integration

### Automated Builds
- **All branches** trigger GitHub Actions builds
- **Pull requests** to `main` and `develop` are tested
- **Tagged releases** create GitHub releases with binaries

### Branch Protection Rules (Recommended)

Set up the following protections on GitHub:

#### Main Branch
- Require pull request reviews
- Require status checks (GitHub Actions)
- Require up-to-date branches
- Restrict force pushes
- Restrict deletions

#### Develop Branch  
- Require status checks (GitHub Actions)
- Require up-to-date branches
- Allow administrators to bypass

### Pull Request Guidelines

#### To `develop`:
- All feature work
- Documentation updates
- Non-critical bug fixes
- Process improvements

#### To `main`:
- Only through git flow release/hotfix process
- Emergency hotfixes (with proper review)
- Manual merges discouraged

## Thread-Safe Development

### Feature Branch Testing
Always test process-locking features on Unix systems:

```bash
# Test concurrent operations
./azcopy copy source1/ dest1/ &
./azcopy copy source2/ dest2/ &
wait

# Check for lock-related errors
ls -la ~/.azcopy/.locks/
```

### Integration Testing
Before finishing features:
- Test on Linux, macOS, and Windows
- Verify no lock file leaks
- Check error handling and fallback behavior

## Release Process

### Version Numbering
Follow semantic versioning: `MAJOR.MINOR.PATCH`

- **MAJOR**: Breaking changes
- **MINOR**: New features, backward compatible  
- **PATCH**: Bug fixes, improvements

### Release Checklist

1. **Start Release**: `git flow release start X.Y.Z`
2. **Update Documentation**: BUILD.md, INSTALL.md, README.md
3. **Version Bumping**: Update version strings in code
4. **Testing**: Full regression testing on all platforms
5. **Changelog**: Update with release notes
6. **Finish Release**: `git flow release finish X.Y.Z`
7. **GitHub Release**: Tag automatically creates GitHub release
8. **Announcement**: Update README with new version info

## Configuration Files

Git flow configuration is stored in `.git/config`:

```ini
[gitflow "branch"]
    master = main
    develop = develop
[gitflow "prefix"]
    feature = feature/
    release = release/
    hotfix = hotfix/
    support = support/
    versiontag = 
```

## Best Practices

### Commit Messages
```
type(scope): description

feat(locking): add flock-based process synchronization
fix(mmf): resolve memory mapping race condition
docs(build): update installation instructions
test(integration): add concurrent operation tests
```

### Branch Naming
- `feature/process-locking-improvements`
- `feature/github-actions-optimization`
- `release/1.2.0`
- `hotfix/memory-leak-fix`

### Code Reviews
- All feature branches should be reviewed
- Focus on thread-safety implications
- Test build artifacts on multiple platforms
- Verify documentation updates

## Quick Reference

```bash
# Setup (already done)
git flow init

# Daily workflow
git flow feature start my-feature
# ... make changes ...
git flow feature finish my-feature

# Release workflow  
git flow release start 1.2.0
# ... prepare release ...
git flow release finish 1.2.0

# Emergency fixes
git flow hotfix start urgent-fix
# ... fix issue ...
git flow hotfix finish urgent-fix
```

## Troubleshooting

### Sync Issues
```bash
# Sync develop with origin
git checkout develop
git pull origin develop

# Sync main with origin
git checkout main  
git pull origin main
```

### Clean State
```bash
# Abort feature in progress
git flow feature rebase

# Reset to clean state
git checkout develop
git reset --hard origin/develop
```