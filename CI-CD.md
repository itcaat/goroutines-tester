# CI/CD Quick Reference

## Workflow Organization

### 📋 CI Workflow (`ci.yml`)

**Triggers:**
- Push to `main` branch
- Pull requests to `main`

**Jobs:**
- ✅ **Go Testing** - Multiple Go versions (1.23.x, 1.24.x)
- ✅ **Cross-platform builds** - Linux/macOS for amd64/arm64
- ✅ **Binary testing** - Functional tests
- ✅ **GoReleaser validation** - Config check
- ✅ **Docker build** - Multi-arch build test (no push)

**Benefit:** Fast comprehensive validation on every commit

### 🚀 Release Workflow (`release.yml`)

**Triggers:**
- Version tags matching `v*` (e.g., `v1.2.3`)

**Jobs:**
- ✅ **GoReleaser** - Binary releases for GitHub
- ✅ **Docker publish** - Push to Docker Hub `itcaat/goroutines-tester`
- ✅ **Security scanning** - Trivy scans on published images

**Docker Tags Created:**
- `latest`
- `v1.2.3` (exact version)
- `v1.2` (minor version)
- `v1` (major version)

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

## GitHub Actions Workflows

### CI Workflow Jobs

**`test`** - Go testing matrix (1.23.x, 1.24.x)
**`build-matrix`** - Cross-platform builds (linux/darwin × amd64/arm64)  
**`test-binary`** - Binary functionality testing
**`goreleaser-check`** - GoReleaser config validation
**`docker-build`** - Docker multi-arch build test

### Release Workflow Jobs

**`goreleaser`** - Creates GitHub releases with binaries
**`docker-publish`** - Publishes Docker images to Hub with security scanning

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
