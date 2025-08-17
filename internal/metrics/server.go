package metrics

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

// Metrics структура для хранения метрик приложения
type Metrics struct {
	TasksTotal     int64         // общее количество выполненных задач
	TasksCompleted int64         // количество завершенных задач
	ExecutionTime  time.Duration // время выполнения последнего запуска
	Mode           string        // режим работы (single/pool)
	Workers        int           // количество воркеров
	BlockSizeKB    int           // размер блока в KB
	StartTime      time.Time     // время запуска приложения
	LastRunTime    time.Time     // время последнего запуска
	TotalRuns      int64         // общее количество запусков
}

// Server представляет HTTP сервер для метрик
type Server struct {
	metrics *Metrics
	version string
	commit  string
	date    string
}

// NewServer создает новый сервер метрик
func NewServer(version, commit, date string) *Server {
	return &Server{
		metrics: &Metrics{
			StartTime: time.Now(),
		},
		version: version,
		commit:  commit,
		date:    date,
	}
}

// Start запускает HTTP сервер для метрик
func (s *Server) Start(port string) {
	http.HandleFunc("/metrics", s.metricsHandler)
	http.HandleFunc("/", s.indexHandler)

	fmt.Printf("Metrics server started on http://localhost:%s\n", port)
	fmt.Printf("Metrics endpoint: http://localhost:%s/metrics\n", port)

	go func() {
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			fmt.Printf("Ошибка запуска metrics сервера: %v\n", err)
		}
	}()
}

// UpdateMetrics обновляет метрики после выполнения бенчмарка
func (s *Server) UpdateMetrics(tasks int, mode string, workers int, blockKB int, duration time.Duration) {
	atomic.StoreInt64(&s.metrics.TasksTotal, int64(tasks))
	atomic.StoreInt64(&s.metrics.TasksCompleted, int64(tasks))
	atomic.AddInt64(&s.metrics.TotalRuns, 1)

	s.metrics.ExecutionTime = duration
	s.metrics.Mode = mode
	s.metrics.Workers = workers
	s.metrics.BlockSizeKB = blockKB
	s.metrics.LastRunTime = time.Now()
}

// metricsHandler возвращает метрики в формате Prometheus
func (s *Server) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	uptime := time.Since(s.metrics.StartTime).Seconds()
	tasksTotal := atomic.LoadInt64(&s.metrics.TasksTotal)
	tasksCompleted := atomic.LoadInt64(&s.metrics.TasksCompleted)
	totalRuns := atomic.LoadInt64(&s.metrics.TotalRuns)

	fmt.Fprintf(w, "# HELP goroutines_tester_info Application information\n")
	fmt.Fprintf(w, "# TYPE goroutines_tester_info gauge\n")
	fmt.Fprintf(w, "goroutines_tester_info{version=\"%s\",commit=\"%s\",date=\"%s\"} 1\n", s.version, s.commit, s.date)

	fmt.Fprintf(w, "# HELP goroutines_tester_uptime_seconds Application uptime in seconds\n")
	fmt.Fprintf(w, "# TYPE goroutines_tester_uptime_seconds gauge\n")
	fmt.Fprintf(w, "goroutines_tester_uptime_seconds %.2f\n", uptime)

	fmt.Fprintf(w, "# HELP goroutines_tester_tasks_total Total number of tasks configured\n")
	fmt.Fprintf(w, "# TYPE goroutines_tester_tasks_total gauge\n")
	fmt.Fprintf(w, "goroutines_tester_tasks_total %d\n", tasksTotal)

	fmt.Fprintf(w, "# HELP goroutines_tester_tasks_completed Number of completed tasks\n")
	fmt.Fprintf(w, "# TYPE goroutines_tester_tasks_completed gauge\n")
	fmt.Fprintf(w, "goroutines_tester_tasks_completed %d\n", tasksCompleted)

	fmt.Fprintf(w, "# HELP goroutines_tester_execution_time_seconds Last execution time in seconds\n")
	fmt.Fprintf(w, "# TYPE goroutines_tester_execution_time_seconds gauge\n")
	fmt.Fprintf(w, "goroutines_tester_execution_time_seconds %.6f\n", s.metrics.ExecutionTime.Seconds())

	fmt.Fprintf(w, "# HELP goroutines_tester_mode_info Current execution mode\n")
	fmt.Fprintf(w, "# TYPE goroutines_tester_mode_info gauge\n")
	fmt.Fprintf(w, "goroutines_tester_mode_info{mode=\"%s\"} 1\n", s.metrics.Mode)

	fmt.Fprintf(w, "# HELP goroutines_tester_workers Number of workers\n")
	fmt.Fprintf(w, "# TYPE goroutines_tester_workers gauge\n")
	fmt.Fprintf(w, "goroutines_tester_workers %d\n", s.metrics.Workers)

	fmt.Fprintf(w, "# HELP goroutines_tester_block_size_kb Block size in KB\n")
	fmt.Fprintf(w, "# TYPE goroutines_tester_block_size_kb gauge\n")
	fmt.Fprintf(w, "goroutines_tester_block_size_kb %d\n", s.metrics.BlockSizeKB)

	fmt.Fprintf(w, "# HELP goroutines_tester_total_runs_total Total number of benchmark runs\n")
	fmt.Fprintf(w, "# TYPE goroutines_tester_total_runs_total counter\n")
	fmt.Fprintf(w, "goroutines_tester_total_runs_total %d\n", totalRuns)

	if !s.metrics.LastRunTime.IsZero() {
		fmt.Fprintf(w, "# HELP goroutines_tester_last_run_timestamp_seconds Timestamp of last run\n")
		fmt.Fprintf(w, "# TYPE goroutines_tester_last_run_timestamp_seconds gauge\n")
		fmt.Fprintf(w, "goroutines_tester_last_run_timestamp_seconds %d\n", s.metrics.LastRunTime.Unix())
	}
}

// indexHandler возвращает главную страницу
func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
    <title>Goroutines Tester Metrics</title>
</head>
<body>
    <h1>Goroutines Tester</h1>
    <p>Metrics endpoint: <a href="/metrics">/metrics</a></p>
    <p>Application version: %s</p>
    <p>Started: %s</p>
</body>
</html>`, s.version, s.metrics.StartTime.Format(time.RFC3339))
}

// ShowInfo показывает информацию о метриках и ожидает завершения
func (s *Server) ShowInfo(port string) {
	fmt.Printf("\nMetrics доступны по адресу: http://localhost:%s/metrics\n", port)
	fmt.Printf("Для остановки нажмите Ctrl+C\n")

	// Держим программу запущенной для доступа к метрикам
	select {}
}
