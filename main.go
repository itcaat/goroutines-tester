package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/itcaat/goroutines-tester/internal/benchmark"
	"github.com/itcaat/goroutines-tester/internal/metrics"
	"github.com/itcaat/goroutines-tester/internal/profiler"
)

// Version information injected by GoReleaser
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func initConfig() {
	// Set default values
	viper.SetDefault("tasks", 200)
	viper.SetDefault("blockKB", 1024)
	viper.SetDefault("mode", "single")
	viper.SetDefault("workers", 0)
	viper.SetDefault("debug", false)
	viper.SetDefault("metrics", false)
	viper.SetDefault("metrics-port", "8888")

	// Bind environment variables
	viper.BindEnv("tasks", "TASKS")
	viper.BindEnv("blockKB", "BLOCK_KB")
	viper.BindEnv("mode", "MODE")
	viper.BindEnv("workers", "WORKERS")
	viper.BindEnv("debug", "DEBUG")
	viper.BindEnv("metrics", "METRICS")
	viper.BindEnv("metrics-port", "METRICS_PORT")

	// Configure command line flags with short versions
	pflag.IntP("tasks", "t", viper.GetInt("tasks"), "number of tasks to execute")
	pflag.IntP("blockKB", "b", viper.GetInt("blockKB"), "data block size per task in KB")
	pflag.StringP("mode", "m", viper.GetString("mode"), "execution mode: single | pool")
	pflag.IntP("workers", "w", viper.GetInt("workers"), "number of workers for pool mode (0 = auto-detect)")
	pflag.BoolP("version", "v", false, "show version information")
	pflag.BoolP("debug", "d", viper.GetBool("debug"), "enable CPU and trace profiling")
	pflag.Bool("metrics", viper.GetBool("metrics"), "enable HTTP metrics server")
	pflag.StringP("metrics-port", "p", viper.GetString("metrics-port"), "port for HTTP metrics server")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
}

func main() {
	initConfig()

	fmt.Println("NumCPU =", runtime.NumCPU())
	fmt.Println("GOMAXPROCS =", runtime.GOMAXPROCS(0))

	// Show version and exit
	if viper.GetBool("version") {
		fmt.Printf("CPU Benchmarking Tool\n")
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Commit: %s\n", commit)
		fmt.Printf("Date: %s\n", date)
		fmt.Printf("Built by: %s\n", builtBy)
		return
	}

	// Get configuration values
	tasks := viper.GetInt("tasks")
	blockKB := viper.GetInt("blockKB")
	mode := viper.GetString("mode")
	workers := viper.GetInt("workers")
	debug := viper.GetBool("debug")
	enableMetrics := viper.GetBool("metrics")
	metricsPort := viper.GetString("metrics-port")

	// Initialize profiling
	var prof *profiler.Profiler
	if debug {
		prof = profiler.New()
		if err := prof.Start(); err != nil {
			fmt.Printf("Failed to start profiling: %v\n", err)
			os.Exit(1)
		}
		defer prof.Stop()
	}

	// Initialize metrics server
	var metricsServer *metrics.Server
	if enableMetrics {
		metricsServer = metrics.NewServer(version, commit, date)
		metricsServer.Start(metricsPort)
		time.Sleep(100 * time.Millisecond) // give server time to start
	}

	if workers < 1 {
		workers = runtime.GOMAXPROCS(0) // взять то, что Go сам посчитал по cgroups
	} else {
		runtime.GOMAXPROCS(workers) // если явно указано через env/флаг
	}

	// Configure benchmark
	config := benchmark.Config{
		Tasks:   tasks,
		BlockKB: blockKB,
		Mode:    mode,
		Workers: workers,
	}

	// Validate mode
	if mode != "single" && mode != "pool" {
		fmt.Println("unknown mode")
		return
	}

	fmt.Printf("mode=%s tasks=%d block=%dKB workers=%d\n", mode, tasks, blockKB, workers)

	// Execute benchmark
	runner := benchmark.NewRunner()
	start := time.Now()
	result := runner.Run(config)
	elapsed := time.Since(start)

	fmt.Printf("done in %v (sink=%d)\n", elapsed, result.Sink)

	// Update metrics
	if enableMetrics {
		metricsServer.UpdateMetrics(tasks, mode, workers, blockKB, elapsed)
		metricsServer.ShowInfo(metricsPort)
	}
}
