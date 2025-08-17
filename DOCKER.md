# Goroutines Tester Docker Image

A containerized CPU benchmarking tool for comparing single-threaded vs multi-threaded performance in Go.

## Quick Start

```bash
# Run with metrics enabled
docker run -p 8888:8888 itcaat/goroutines-tester:latest

# One-shot execution without metrics
docker run itcaat/goroutines-tester:latest \
  ./goroutines-tester -t 100 -m single
```

## Configuration

### Environment Variables

- `TASKS` - Number of tasks to execute (default: 200)
- `BLOCK_KB` - Size of data block per task in KB (default: 1024)
- `MODE` - Execution mode: `single` or `pool` (default: single)
- `WORKERS` - Number of workers for pool mode (default: 0 = auto)
- `METRICS_PORT` - Port for metrics server (default: 8888)

### Example with Custom Configuration

```bash
docker run -p 8888:8888 \
  -e TASKS=500 \
  -e BLOCK_KB=2048 \
  -e MODE=pool \
  -e WORKERS=8 \
  itcaat/goroutines-tester:latest
```

## Metrics

The container exposes Prometheus-compatible metrics at `/metrics` endpoint:

- Application info and version
- Task execution metrics
- Performance timing
- Resource utilization

```bash
# View metrics
curl http://localhost:8888/metrics
```

## Health Check

The image includes a built-in health check that monitors the metrics endpoint:

```bash
# Check container health
docker inspect --format='{{.State.Health.Status}}' <container-id>
```

## Multi-Architecture Support

This image supports multiple architectures:
- `linux/amd64`
- `linux/arm64`

## Tags

- `latest` - Latest stable release
- `v1.x.x` - Specific version tags
- `main` - Latest development build

## Docker Compose

See the full docker-compose.yml with monitoring stack in the [GitHub repository](https://github.com/itcaat/goroutines-tester).

## Source Code

Full source code and documentation: https://github.com/itcaat/goroutines-tester
