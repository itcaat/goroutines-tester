# Setup Instructions

## GitHub Actions & Docker Hub Integration

### 1. Create Docker Hub Account

1. Go to [Docker Hub](https://hub.docker.com/)
2. Create an account or sign in
3. Create a repository named `goroutines-tester`

### 2. Generate Docker Hub Access Token

1. Go to Docker Hub Settings → Security
2. Click "New Access Token"
3. Give it a name (e.g., "github-actions")
4. Set permissions to "Read, Write, Delete"
5. Copy the generated token

### 3. Configure GitHub Repository Secrets

Go to your GitHub repository → Settings → Secrets and variables → Actions

Add the following secrets:

- **`DOCKER_USERNAME`**: Your Docker Hub username
- **`DOCKER_TOKEN`**: The access token from step 2

### 4. Update Docker References

Update the Docker Hub username references to `itcaat` in:

- `README.md` ✅ Already updated
- `DOCKER.md` ✅ Already updated  
- `Makefile` ✅ Already updated

### 5. Test GitHub Actions

1. Push changes to main branch
2. Check Actions tab in GitHub to see the Docker build (image is built but not pushed)
3. Create a version tag to trigger Docker Hub publishing

## Local Development Setup

### Prerequisites

- Go 1.24.3+
- Docker and Docker Compose (optional)
- Make (optional, for convenience)

### Quick Start

```bash
# Clone and setup
git clone https://github.com/itcaat/goroutines-tester
cd goroutines-tester
make dev-setup

# Build and test
make build
make test

# Run locally
make run

# Build Docker image
make docker-build
```

## CI/CD Workflow

### Docker Build Strategy

The project uses a two-stage Docker workflow:

**1. Build Stage (Always runs):**
- Triggers on: commits to `main`, `develop`, and pull requests
- Actions: Builds Docker images for testing, stores as artifacts
- Does NOT push to Docker Hub
- Runs security scans locally
- Verifies image builds correctly on multiple architectures

**2. Push Stage (Only for tags):**
- Triggers on: version tags (e.g., `v1.2.3`)
- Actions: Builds and pushes images to Docker Hub
- Creates tags: `latest`, `v1.2.3`, `v1.2`, `v1`
- Runs security scans on published images
- Uploads scan results to GitHub Security

### Benefits:
- ✅ Fast feedback on build issues
- ✅ No unnecessary Docker Hub pushes
- ✅ Clean version management
- ✅ Security scanning on published images only

## Release Process

### Create a New Release

1. **Update version in code** (if needed)
2. **Tag the release:**
   ```bash
   git tag -a v1.1.0 -m "Release v1.1.0: Added Docker support"
   git push origin v1.1.0
   ```
3. **GitHub Actions will automatically:**
   - Build multi-architecture Docker images (on every commit)
   - Push to Docker Hub with version tags (only for version tags)
   - Run security scans (only when publishing)
   - Create GitHub release (via separate GoReleaser workflow)

### Manual Docker Push

```bash
# Set your Docker Hub username
export DOCKER_USERNAME=itcaat

# Build and push
make docker-push
```

## Monitoring Stack

### Local Development with Monitoring

```bash
# Start application + Prometheus + Grafana
make monitoring

# Access services:
# - Application: http://localhost:8080
# - Prometheus: http://localhost:9090  
# - Grafana: http://localhost:3000 (admin/admin)

# Stop everything
make monitoring-stop
```

### Production Deployment

For production use, consider:

1. **Resource limits** in docker-compose.yml
2. **Persistent storage** for Grafana dashboards
3. **Security** - change default passwords
4. **SSL/TLS** termination with reverse proxy
5. **Alerting** configuration in Prometheus

## Troubleshooting

### Docker Build Issues

```bash
# Clean Docker cache
make docker-clean

# Rebuild without cache
docker build --no-cache -t goroutines-tester .
```

### GitHub Actions Issues

1. Check secrets are set correctly
2. Verify Docker Hub repository exists
3. Check Actions logs for detailed errors
4. Ensure main branch protection rules allow Actions

### Local Issues

```bash
# Clean build artifacts
make clean

# Reset dependencies
go mod tidy
go mod download

# Check linting
make lint
```
