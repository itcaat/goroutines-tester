# CPU Benchmarking Tool

An educational Go program that demonstrates CPU-intensive task processing using different execution modes to compare single-threaded vs multi-threaded performance.

## Overview

This program simulates heavy CPU work by generating random data blocks and computing their SHA-256 hashes. It's designed to help understand:
- Single-threaded vs multi-threaded processing
- Worker pool patterns in Go
- CPU profiling and tracing
- Performance optimization techniques

## Features

- **Two execution modes:**
  - `single`: Sequential processing on a single thread
  - `pool`: Parallel processing using a worker pool pattern
- **Built-in profiling:** Optional CPU and trace profiles generation
- **Metrics export:** Prometheus-compatible metrics via HTTP endpoint
- **Docker support:** Containerized deployment with multi-architecture builds
- **Configurable parameters:** Number of tasks, block size, and worker count
- **Performance measurement:** Reports execution time for comparison
- **Monitoring stack:** Optional Prometheus + Grafana integration

## Installation

### From Source

Make sure you have Go 1.24.3 or later installed.

```bash
git clone https://github.com/itcaat/goroutines-tester
cd goroutines-tester
go mod tidy
```

### Using Docker

```bash
# Pull from Docker Hub
docker pull itcaat/goroutines-tester:latest

# Or build locally
docker build -t goroutines-tester .
```

### Using Docker Compose

```bash
# Run application with metrics
docker-compose up -d

# Run with full monitoring stack (Prometheus + Grafana)
docker-compose --profile monitoring up -d
```

## Usage

### Basic Usage

```bash
# Run with default settings (200 tasks, 1024KB blocks, single mode)
go run main.go

# Run in parallel mode with worker pool
go run main.go -m pool

# Custom configuration
go run main.go -t 500 -b 2048 -m pool -w 8
```

### Command Line Options

- `-t, --tasks`: Number of tasks to execute (default: 200)
- `-b, --blockKB`: Size of data block per task in KB (default: 1024)
- `-m, --mode`: Execution mode - `single` or `pool` (default: "single")
- `-w, --workers`: Number of workers for pool mode (default: number of CPU cores)
- `-d, --debug`: Enable CPU and trace profiling (default: false)
- `--metrics`: Enable HTTP metrics server (default: false)
- `-p, --metrics-port`: Port for metrics server (default: "8888")
- `-v, --version`: Show version information

### Examples

```bash
# Compare single-threaded vs multi-threaded performance
go run main.go -m single -t 100
go run main.go -m pool -t 100

# Stress test with larger workload
go run main.go -m pool -t 1000 -b 4096 -w 16

# Minimal workload for quick testing
go run main.go -t 50 -b 512
```

## Docker Usage

### Quick Start

```bash
# Run with default settings
docker run -p 8888:8888 itcaat/goroutines-tester:latest

# Custom configuration via environment variables
docker run -p 8888:8888 \
  -e TASKS=500 \
  -e BLOCK_KB=2048 \
  -e MODE=pool \
  -e WORKERS=8 \
  itcaat/goroutines-tester:latest

# Run without metrics (one-shot execution)
docker run itcaat/goroutines-tester:latest \
  ./goroutines-tester -t 100 -m single
```

### Using Makefile

```bash
# Build Docker image
make docker-build

# Run with Docker Compose
make docker-run

# Start full monitoring stack
make monitoring

# View logs
make docker-logs

# Stop containers
make docker-stop
```

### Monitoring with Prometheus & Grafana

When running with the monitoring profile:

- **Application**: http://localhost:8888
- **Metrics**: http://localhost:8888/metrics
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000 (admin/admin)

```bash
# Start monitoring stack
docker-compose --profile monitoring up -d

# Stop monitoring stack
docker-compose --profile monitoring down
```

## Performance Analysis

### Profiling (Debug Mode)

When using `-d` or `--debug` flag, the program generates profiling files:
- `cpu.out`: CPU profile for analyzing performance bottlenecks
- `trace.out`: Execution trace for understanding goroutine behavior

```bash
# View CPU profile
go tool pprof cpu.out

# View execution trace
go tool trace trace.out
```

### Metrics (Metrics Mode)

When using `-metrics` flag, the program exposes Prometheus metrics at `/metrics` endpoint:

```bash
# View metrics
curl http://localhost:8888/metrics

# Key metrics available:
# - goroutines_tester_tasks_total
# - goroutines_tester_execution_time_seconds
# - goroutines_tester_workers
# - goroutines_tester_uptime_seconds
```

## Educational Value

This program demonstrates several important concepts:

1. **Concurrency Patterns**: Compare sequential vs parallel execution
2. **Worker Pools**: Efficient task distribution among goroutines
3. **Performance Profiling**: Built-in tools for performance analysis
4. **CPU-bound Workloads**: Understanding CPU utilization patterns
5. **Go Best Practices**: Proper use of channels, waitgroups, and goroutines

## How It Works

1. **Task Generation**: Each task generates a random data block of specified size
2. **Hash Computation**: SHA-256 hash is computed for each block (CPU-intensive)
3. **Result Aggregation**: Results are XORed together to prevent compiler optimizations
4. **Performance Measurement**: Execution time is measured and reported

### Single Mode
Tasks are processed sequentially in a single goroutine.

### Pool Mode
Tasks are distributed among multiple worker goroutines using channels for coordination.

## Expected Results

You should observe:
- **Pool mode** significantly faster than single mode on multi-core systems
- **Performance scaling** with the number of workers (up to CPU core count)
- **CPU utilization** differences between modes when analyzing profiles

## License

This project is licensed under the terms specified in the LICENSE file.

## Development & CI

### Continuous Integration

The project includes automated CI that runs on every commit to the main branch and on pull requests:

- **Testing**: Runs tests across multiple Go versions (1.23.x, 1.24.x)
- **Cross-platform builds**: Verifies the code builds on Linux, Windows, and macOS
- **Code quality**: Runs `go vet` and `staticcheck` for code analysis
- **Docker builds**: Builds and tests Docker images on every commit
- **Docker publishing**: Pushes to Docker Hub only when creating version tags
- **Security scanning**: Runs Trivy security scans on published images
- **GoReleaser validation**: Ensures release configuration is valid

### Running Tests Locally

```bash
# Run all tests
go test -v ./...

# Run tests with race detection
go test -race ./...

# Run benchmarks
go test -bench=. ./...

# Run static analysis
go vet ./...
staticcheck ./...
```

## Release Process

This project uses [GoReleaser](https://goreleaser.com/) with GitHub Actions for automated releases.

### Creating a Release

1. **Tag a new version:**
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

2. **GitHub Actions will automatically:**
   - Build binaries for Linux, Windows, and macOS (both amd64 and arm64)
   - Create release archives
   - Generate checksums
   - Create a GitHub release with changelog

### Local Testing

Test the release configuration locally:

```bash
# Install GoReleaser (if not already installed)
go install github.com/goreleaser/goreleaser@latest

# Test the build without publishing
goreleaser build --snapshot --clean

# Test the full release process without publishing
goreleaser release --snapshot --clean
```

### Version Information

The binary includes version information that can be displayed:

```bash
./cpu-benchmarking-tool -v
```

## Contributing

This is an educational project. Feel free to fork and experiment with different optimizations or add new execution modes for learning purposes.