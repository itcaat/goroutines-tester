# CI/CD Quick Reference

## Docker Build & Push Strategy

### 🔄 Continuous Build (Every Commit)

**Triggers:**
- Push to `main` or `develop` branches
- Pull requests to `main`

**Actions:**
- ✅ Build Docker image (multi-arch: amd64, arm64)
- ✅ Test build process
- ✅ Cache layers for faster builds
- ✅ Validate multi-architecture compatibility
- ❌ NO push to Docker Hub
- ❌ NO artifacts stored (just validation)

**Benefit:** Fast feedback on build issues without cluttering Docker Hub

### 🚀 Release Publishing (Tags Only)

**Triggers:**
- Version tags matching `v*` (e.g., `v1.2.3`)

**Actions:**
- ✅ Build Docker image (multi-arch)
- ✅ Push to Docker Hub `itcaat/goroutines-tester`
- ✅ Create semantic version tags:
  - `latest`
  - `v1.2.3` (exact version)
  - `v1.2` (minor version)
  - `v1` (major version)
- ✅ Run security scan with Trivy
- ✅ Upload security results to GitHub

## Release Workflow

### 1. Development
```bash
# Regular development - triggers build only
git add .
git commit -m "feat: add new feature"
git push origin main
# ✅ Builds Docker image, ❌ doesn't publish
```

### 2. Release
```bash
# Create release - triggers build AND publish
git tag -a v1.2.3 -m "Release v1.2.3: Description"
git push origin v1.2.3
# ✅ Builds and publishes to Docker Hub
```

### 3. Using Published Images
```bash
# Pull latest version
docker pull itcaat/goroutines-tester:latest

# Pull specific version
docker pull itcaat/goroutines-tester:v1.2.3

# Pull minor version (auto-updates patch versions)
docker pull itcaat/goroutines-tester:v1.2
```

## GitHub Actions Jobs

### Job 1: `build` (Always runs)
- Validates Docker build process
- Tests multi-architecture builds (amd64, arm64)
- Verifies Dockerfile syntax and dependencies
- Caches layers for performance
- Fast feedback without publishing

### Job 2: `push` (Tag-triggered only)
- Requires `build` job success
- Authenticates with Docker Hub
- Publishes tagged images
- Runs security scanning

## Benefits

✅ **No unnecessary publishes** - Clean Docker Hub repository  
✅ **Fast feedback** - Build issues caught immediately  
✅ **Semantic versioning** - Proper version management  
✅ **Security scanning** - Only on published images  
✅ **Multi-architecture** - Supports AMD64 and ARM64  
✅ **Caching** - Faster subsequent builds  

## Troubleshooting

### Build Fails on Commit
```bash
# Check GitHub Actions logs
# Fix issues and push again
git add .
git commit -m "fix: resolve build issue"
git push origin main
```

### Release Fails
```bash
# Delete problematic tag
git tag -d v1.2.3
git push origin :refs/tags/v1.2.3

# Fix issues and retag
git tag -a v1.2.3 -m "Release v1.2.3: Fixed"
git push origin v1.2.3
```

### Docker Hub Issues
- Verify `DOCKER_USERNAME` and `DOCKER_TOKEN` secrets
- Check Docker Hub repository exists: `itcaat/goroutines-tester`
- Ensure Docker Hub account has push permissions
