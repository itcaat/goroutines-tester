package metrics

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics holds all Prometheus metrics
type Metrics struct {
	appInfo          *prometheus.GaugeVec
	uptime           prometheus.Gauge
	tasksTotal       prometheus.Gauge
	tasksCompleted   prometheus.Gauge
	executionTime    prometheus.Gauge
	modeInfo         *prometheus.GaugeVec
	workers          prometheus.Gauge
	blockSizeKB      prometheus.Gauge
	totalRuns        prometheus.Counter
	lastRunTimestamp prometheus.Gauge
	startTime        time.Time
}

// Server represents HTTP server for metrics
type Server struct {
	metrics  *Metrics
	registry *prometheus.Registry
	version  string
	commit   string
	date     string
}

// NewServer creates a new metrics server
func NewServer(version, commit, date string) *Server {
	registry := prometheus.NewRegistry()

	metrics := &Metrics{
		appInfo: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "goroutines_tester_info",
				Help: "Application information",
			},
			[]string{"version", "commit", "date"},
		),
		uptime: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "goroutines_tester_uptime_seconds",
			Help: "Application uptime in seconds",
		}),
		tasksTotal: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "goroutines_tester_tasks_total",
			Help: "Total number of tasks configured",
		}),
		tasksCompleted: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "goroutines_tester_tasks_completed",
			Help: "Number of completed tasks",
		}),
		executionTime: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "goroutines_tester_execution_time_seconds",
			Help: "Last execution time in seconds",
		}),
		modeInfo: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "goroutines_tester_mode_info",
				Help: "Current execution mode",
			},
			[]string{"mode"},
		),
		workers: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "goroutines_tester_workers",
			Help: "Number of workers",
		}),
		blockSizeKB: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "goroutines_tester_block_size_kb",
			Help: "Block size in KB",
		}),
		totalRuns: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "goroutines_tester_total_runs_total",
			Help: "Total number of benchmark runs",
		}),
		lastRunTimestamp: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "goroutines_tester_last_run_timestamp_seconds",
			Help: "Timestamp of last run",
		}),
		startTime: time.Now(),
	}

	// Register all metrics
	registry.MustRegister(
		metrics.appInfo,
		metrics.uptime,
		metrics.tasksTotal,
		metrics.tasksCompleted,
		metrics.executionTime,
		metrics.modeInfo,
		metrics.workers,
		metrics.blockSizeKB,
		metrics.totalRuns,
		metrics.lastRunTimestamp,
	)

	// Register standard Go runtime metrics
	registry.MustRegister(prometheus.NewGoCollector())
	registry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))

	// Set static metrics
	metrics.appInfo.WithLabelValues(version, commit, date).Set(1)

	return &Server{
		metrics:  metrics,
		registry: registry,
		version:  version,
		commit:   commit,
		date:     date,
	}
}

// Start launches HTTP server for metrics
func (s *Server) Start(port string) {
	// Create a new ServeMux for this server
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(s.registry, promhttp.HandlerOpts{}))
	mux.HandleFunc("/", s.indexHandler)

	fmt.Printf("Metrics server started on http://localhost:%s\n", port)
	fmt.Printf("Metrics endpoint: http://localhost:%s/metrics\n", port)

	go func() {
		if err := http.ListenAndServe(":"+port, mux); err != nil {
			fmt.Printf("Failed to start metrics server: %v\n", err)
		}
	}()

	// Start background goroutine to update uptime
	go s.updateUptime()
}

// UpdateMetrics updates metrics after benchmark execution
func (s *Server) UpdateMetrics(tasks int, mode string, workers int, blockKB int, duration time.Duration) {
	s.metrics.tasksTotal.Set(float64(tasks))
	s.metrics.tasksCompleted.Set(float64(tasks))
	s.metrics.executionTime.Set(duration.Seconds())
	s.metrics.workers.Set(float64(workers))
	s.metrics.blockSizeKB.Set(float64(blockKB))
	s.metrics.lastRunTimestamp.Set(float64(time.Now().Unix()))

	// Reset mode info and set current mode
	s.metrics.modeInfo.Reset()
	s.metrics.modeInfo.WithLabelValues(mode).Set(1)

	// Increment total runs counter
	s.metrics.totalRuns.Inc()
}

// updateUptime runs in background to update uptime metric
func (s *Server) updateUptime() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		uptime := time.Since(s.metrics.startTime).Seconds()
		s.metrics.uptime.Set(uptime)
	}
}

// indexHandler returns the main page
func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
    <title>Goroutines Tester Metrics</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .metric-link { background: #f0f0f0; padding: 10px; border-radius: 5px; display: inline-block; margin: 10px 0; }
        .info { background: #e8f4f8; padding: 15px; border-radius: 5px; margin: 20px 0; }
    </style>
</head>
<body>
    <h1> Goroutines Tester Metrics</h1>
    <div class="metric-link">
        <strong>Metrics endpoint:</strong> <a href="/metrics">/metrics</a>
    </div>
    <div class="info">
        <p><strong>Version:</strong> %s</p>
        <p><strong>Commit:</strong> %s</p>
        <p><strong>Started:</strong> %s</p>
    </div>
    <h3>Available Metrics:</h3>
    
    <h4> Application Metrics:</h4>
    <ul>
        <li><strong>goroutines_tester_info</strong> - Application information</li>
        <li><strong>goroutines_tester_uptime_seconds</strong> - Application uptime</li>
        <li><strong>goroutines_tester_tasks_total</strong> - Total number of tasks</li>
        <li><strong>goroutines_tester_execution_time_seconds</strong> - Last execution time</li>
        <li><strong>goroutines_tester_total_runs_total</strong> - Total benchmark runs</li>
    </ul>
    
    <h4> Go Runtime Metrics:</h4>
    <ul>
        <li><strong>go_memstats_*</strong> - Memory statistics (heap, stack, GC)</li>
        <li><strong>go_gc_*</strong> - Garbage collector metrics</li>
        <li><strong>go_goroutines</strong> - Number of goroutines</li>
        <li><strong>go_threads</strong> - Number of OS threads</li>
        <li><strong>process_*</strong> - Process metrics (CPU, memory, file descriptors)</li>
    </ul>
</body>
</html>`, s.version, s.commit, s.metrics.startTime.Format(time.RFC3339))
}

// ShowInfo displays metrics information and waits for completion
func (s *Server) ShowInfo(port string) {
	fmt.Printf("\nMetrics available at: http://localhost:%s/metrics\n", port)
	fmt.Printf("Press Ctrl+C to stop\n")

	// Keep program running for metrics access
	select {}
}
