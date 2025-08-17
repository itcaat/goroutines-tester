package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"

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

// getEnvInt возвращает значение из ENV или дефолт если ENV не установлена
func getEnvInt(key string, defaultValue int) int {
	if val := os.Getenv(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// getEnvString возвращает значение из ENV или дефолт если ENV не установлена
func getEnvString(key string, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

// getEnvBool возвращает значение из ENV или дефолт если ENV не установлена
func getEnvBool(key string, defaultValue bool) bool {
	if val := os.Getenv(key); val != "" {
		if boolVal, err := strconv.ParseBool(val); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func main() {
	// Настройка флагов командной строки с учетом ENV переменных
	var (
		tasks         = flag.Int("tasks", getEnvInt("TASKS", 200), "сколько задач выполнить")
		blockKB       = flag.Int("blockKB", getEnvInt("BLOCK_KB", 1024), "размер блока данных на задачу (KB)")
		mode          = flag.String("mode", getEnvString("MODE", "single"), "режим: single | pool")
		workers       = flag.Int("workers", getEnvInt("WORKERS", runtime.NumCPU()), "кол-во воркеров для режима pool")
		showVersion   = flag.Bool("version", false, "показать версию")
		debug         = flag.Bool("debug", getEnvBool("DEBUG", false), "включить сбор профилей trace.out и cpu.out")
		enableMetrics = flag.Bool("metrics", getEnvBool("METRICS", false), "включить HTTP сервер для экспорта метрик")
		metricsPort   = flag.String("metrics-port", getEnvString("METRICS_PORT", "8888"), "порт для HTTP сервера метрик")
	)
	flag.Parse()

	// Показать версию и выйти
	if *showVersion {
		fmt.Printf("CPU Benchmarking Tool\n")
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Commit: %s\n", commit)
		fmt.Printf("Date: %s\n", date)
		fmt.Printf("Built by: %s\n", builtBy)
		return
	}

	// Инициализация профилирования
	var prof *profiler.Profiler
	if *debug {
		prof = profiler.New()
		if err := prof.Start(); err != nil {
			fmt.Printf("Ошибка запуска профилирования: %v\n", err)
			os.Exit(1)
		}
		defer prof.Stop()
	}

	// Инициализация метрик сервера
	var metricsServer *metrics.Server
	if *enableMetrics {
		metricsServer = metrics.NewServer(version, commit, date)
		metricsServer.Start(*metricsPort)
		time.Sleep(100 * time.Millisecond) // даем серверу время на запуск
	}

	// Валидация параметров
	if *workers < 1 {
		*workers = 1
	}
	runtime.GOMAXPROCS(*workers)

	// Конфигурация бенчмарка
	config := benchmark.Config{
		Tasks:   *tasks,
		BlockKB: *blockKB,
		Mode:    *mode,
		Workers: *workers,
	}

	// Проверка режима
	if *mode != "single" && *mode != "pool" {
		fmt.Println("unknown mode")
		return
	}

	fmt.Printf("mode=%s tasks=%d block=%dKB workers=%d\n", *mode, *tasks, *blockKB, *workers)

	// Выполнение бенчмарка
	runner := benchmark.NewRunner()
	start := time.Now()
	result := runner.Run(config)
	elapsed := time.Since(start)

	fmt.Printf("done in %v (sink=%d)\n", elapsed, result.Sink)

	// Обновление метрик
	if *enableMetrics {
		metricsServer.UpdateMetrics(*tasks, *mode, *workers, *blockKB, elapsed)
		metricsServer.ShowInfo(*metricsPort)
	}
}
