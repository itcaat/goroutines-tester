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
	// Устанавливаем дефолтные значения
	viper.SetDefault("tasks", 200)
	viper.SetDefault("blockKB", 1024)
	viper.SetDefault("mode", "single")
	viper.SetDefault("workers", 0)
	viper.SetDefault("debug", false)
	viper.SetDefault("metrics", false)
	viper.SetDefault("metrics-port", "8888")

	// Связываем ENV переменные
	viper.BindEnv("tasks", "TASKS")
	viper.BindEnv("blockKB", "BLOCK_KB")
	viper.BindEnv("mode", "MODE")
	viper.BindEnv("workers", "WORKERS")
	viper.BindEnv("debug", "DEBUG")
	viper.BindEnv("metrics", "METRICS")
	viper.BindEnv("metrics-port", "METRICS_PORT")

	// Настраиваем флаги командной строки с короткими версиями
	pflag.IntP("tasks", "t", viper.GetInt("tasks"), "сколько задач выполнить")
	pflag.IntP("blockKB", "b", viper.GetInt("blockKB"), "размер блока данных на задачу (KB)")
	pflag.StringP("mode", "m", viper.GetString("mode"), "режим: single | pool")
	pflag.IntP("workers", "w", viper.GetInt("workers"), "кол-во воркеров для режима pool (0 = автоопределение)")
	pflag.BoolP("version", "v", false, "показать версию")
	pflag.BoolP("debug", "d", viper.GetBool("debug"), "включить сбор профилей trace.out и cpu.out")
	pflag.Bool("metrics", viper.GetBool("metrics"), "включить HTTP сервер для экспорта метрик")
	pflag.StringP("metrics-port", "p", viper.GetString("metrics-port"), "порт для HTTP сервера метрик")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
}

func main() {
	initConfig()

	// Показать версию и выйти
	if viper.GetBool("version") {
		fmt.Printf("CPU Benchmarking Tool\n")
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Commit: %s\n", commit)
		fmt.Printf("Date: %s\n", date)
		fmt.Printf("Built by: %s\n", builtBy)
		return
	}

	// Получаем значения конфигурации
	tasks := viper.GetInt("tasks")
	blockKB := viper.GetInt("blockKB")
	mode := viper.GetString("mode")
	workers := viper.GetInt("workers")
	debug := viper.GetBool("debug")
	enableMetrics := viper.GetBool("metrics")
	metricsPort := viper.GetString("metrics-port")

	// Инициализация профилирования
	var prof *profiler.Profiler
	if debug {
		prof = profiler.New()
		if err := prof.Start(); err != nil {
			fmt.Printf("Ошибка запуска профилирования: %v\n", err)
			os.Exit(1)
		}
		defer prof.Stop()
	}

	// Инициализация метрик сервера
	var metricsServer *metrics.Server
	if enableMetrics {
		metricsServer = metrics.NewServer(version, commit, date)
		metricsServer.Start(metricsPort)
		time.Sleep(100 * time.Millisecond) // даем серверу время на запуск
	}

	// Валидация параметров
	if workers < 1 {
		workers = runtime.NumCPU() // если workers не задан или 0, используем все ядра
	}
	runtime.GOMAXPROCS(workers)

	// Конфигурация бенчмарка
	config := benchmark.Config{
		Tasks:   tasks,
		BlockKB: blockKB,
		Mode:    mode,
		Workers: workers,
	}

	// Проверка режима
	if mode != "single" && mode != "pool" {
		fmt.Println("unknown mode")
		return
	}

	fmt.Printf("mode=%s tasks=%d block=%dKB workers=%d\n", mode, tasks, blockKB, workers)

	// Выполнение бенчмарка
	runner := benchmark.NewRunner()
	start := time.Now()
	result := runner.Run(config)
	elapsed := time.Since(start)

	fmt.Printf("done in %v (sink=%d)\n", elapsed, result.Sink)

	// Обновление метрик
	if enableMetrics {
		metricsServer.UpdateMetrics(tasks, mode, workers, blockKB, elapsed)
		metricsServer.ShowInfo(metricsPort)
	}
}
