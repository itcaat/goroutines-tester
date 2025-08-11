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
- **Built-in profiling:** Automatically generates CPU and trace profiles
- **Configurable parameters:** Number of tasks, block size, and worker count
- **Performance measurement:** Reports execution time for comparison

## Installation

Make sure you have Go 1.24.3 or later installed.

```bash
git clone https://github.com/itcaat/prime-numbers
cd prime-numbers
go mod tidy
```

## Usage

### Basic Usage

```bash
# Run with default settings (200 tasks, 1024KB blocks, single mode)
go run main.go

# Run in parallel mode with worker pool
go run main.go -mode=pool

# Custom configuration
go run main.go -tasks=500 -blockKB=2048 -mode=pool -workers=8
```

### Command Line Options

- `-tasks`: Number of tasks to execute (default: 200)
- `-blockKB`: Size of data block per task in KB (default: 1024)
- `-mode`: Execution mode - `single` or `pool` (default: "single")
- `-workers`: Number of workers for pool mode (default: number of CPU cores)

### Examples

```bash
# Compare single-threaded vs multi-threaded performance
go run main.go -mode=single -tasks=100
go run main.go -mode=pool -tasks=100

# Stress test with larger workload
go run main.go -mode=pool -tasks=1000 -blockKB=4096 -workers=16

# Minimal workload for quick testing
go run main.go -tasks=50 -blockKB=512
```

## Performance Analysis

The program automatically generates profiling files:
- `cpu.out`: CPU profile for analyzing performance bottlenecks
- `trace.out`: Execution trace for understanding goroutine behavior

### Analyzing Profiles

```bash
# View CPU profile
go tool pprof cpu.out

# View execution trace
go tool trace trace.out
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

## Contributing

This is an educational project. Feel free to fork and experiment with different optimizations or add new execution modes for learning purposes.