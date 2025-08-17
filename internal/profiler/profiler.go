package profiler

import (
	"fmt"
	"os"
	"runtime/pprof"
	"runtime/trace"
)

// Profiler manages application profiling
type Profiler struct {
	cpuFile   *os.File
	traceFile *os.File
	enabled   bool
}

// New creates a new Profiler
func New() *Profiler {
	return &Profiler{}
}

// Start begins profiling (CPU and trace)
func (p *Profiler) Start() error {
	if p.enabled {
		return fmt.Errorf("profiling already started")
	}

	var err error

	// Create file for CPU profile
	p.cpuFile, err = os.Create("cpu.out")
	if err != nil {
		return fmt.Errorf("failed to create cpu.out: %w", err)
	}

	// Start CPU profiling
	err = pprof.StartCPUProfile(p.cpuFile)
	if err != nil {
		p.cpuFile.Close()
		return fmt.Errorf("failed to start CPU profiling: %w", err)
	}

	// Create file for trace
	p.traceFile, err = os.Create("trace.out")
	if err != nil {
		pprof.StopCPUProfile()
		p.cpuFile.Close()
		return fmt.Errorf("failed to create trace.out: %w", err)
	}

	// Start trace profiling
	err = trace.Start(p.traceFile)
	if err != nil {
		pprof.StopCPUProfile()
		p.cpuFile.Close()
		p.traceFile.Close()
		return fmt.Errorf("failed to start trace profiling: %w", err)
	}

	p.enabled = true
	fmt.Println("Debug mode enabled: profiles will be saved to cpu.out and trace.out")
	return nil
}

// Stop stops profiling and closes files
func (p *Profiler) Stop() {
	if !p.enabled {
		return
	}

	// Stop profiling
	trace.Stop()
	pprof.StopCPUProfile()

	// Close files
	if p.traceFile != nil {
		p.traceFile.Close()
	}
	if p.cpuFile != nil {
		p.cpuFile.Close()
	}

	p.enabled = false
}
