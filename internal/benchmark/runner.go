package benchmark

import (
	"crypto/sha256"
	"math/rand"
	"sync"
)

// Config configuration for benchmark
type Config struct {
	Tasks   int    // number of tasks
	BlockKB int    // block size in KB
	Mode    string // execution mode: single | pool
	Workers int    // number of workers for pool mode
}

// Result benchmark execution result
type Result struct {
	Sink byte // computation result (to prevent optimization)
}

// Runner executes benchmarks
type Runner struct{}

// NewRunner creates a new Runner
func NewRunner() *Runner {
	return &Runner{}
}

// Run executes benchmark according to configuration
func (r *Runner) Run(config Config) Result {
	blockBytes := config.BlockKB * 1024
	var sink byte

	switch config.Mode {
	case "single":
		sink = r.runSingle(config.Tasks, blockBytes)
	case "pool":
		sink = r.runPool(config.Tasks, blockBytes, config.Workers)
	default:
		// for unknown mode return zero result
		return Result{Sink: 0}
	}

	return Result{Sink: sink}
}

// runSingle executes tasks sequentially in a single goroutine
func (r *Runner) runSingle(tasks int, blockBytes int) byte {
	var sink byte
	for i := 0; i < tasks; i++ {
		sum := doTask(i, blockBytes)
		sink ^= sum[0]
	}
	return sink
}

// runPool executes tasks using a worker pool
func (r *Runner) runPool(tasks int, blockBytes int, workers int) byte {
	jobs := make(chan int, workers*2)
	results := make(chan [32]byte, workers*2)
	var wg sync.WaitGroup

	// Start workers
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for id := range jobs {
				results <- doTask(id, blockBytes)
			}
		}()
	}

	// Send tasks
	go func() {
		for i := 0; i < tasks; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	// Wait for workers to finish and close results channel
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var sink byte
	for s := range results {
		sink ^= s[0]
	}

	return sink
}

// doTask simulates heavy CPU task: generate block and hash it
func doTask(id, blockBytes int) [32]byte {
	r := rand.New(rand.NewSource(int64(id)))
	buf := make([]byte, blockBytes)
	for i := range buf {
		buf[i] = byte(r.Intn(256))
	}
	return sha256.Sum256(buf)
}
