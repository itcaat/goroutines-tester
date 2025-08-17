package benchmark

import (
	"crypto/sha256"
	"math/rand"
	"sync"
)

// Config конфигурация для бенчмарка
type Config struct {
	Tasks   int    // количество задач
	BlockKB int    // размер блока в KB
	Mode    string // режим выполнения: single | pool
	Workers int    // количество воркеров для режима pool
}

// Result результат выполнения бенчмарка
type Result struct {
	Sink byte // результат вычислений (для предотвращения оптимизации)
}

// Runner выполняет бенчмарки
type Runner struct{}

// NewRunner создает новый Runner
func NewRunner() *Runner {
	return &Runner{}
}

// Run выполняет бенчмарк согласно конфигурации
func (r *Runner) Run(config Config) Result {
	blockBytes := config.BlockKB * 1024
	var sink byte

	switch config.Mode {
	case "single":
		sink = r.runSingle(config.Tasks, blockBytes)
	case "pool":
		sink = r.runPool(config.Tasks, blockBytes, config.Workers)
	default:
		// для неизвестного режима возвращаем нулевой результат
		return Result{Sink: 0}
	}

	return Result{Sink: sink}
}

// runSingle выполняет задачи последовательно в одной горутине
func (r *Runner) runSingle(tasks int, blockBytes int) byte {
	var sink byte
	for i := 0; i < tasks; i++ {
		sum := doTask(i, blockBytes)
		sink ^= sum[0]
	}
	return sink
}

// runPool выполняет задачи с использованием пула воркеров
func (r *Runner) runPool(tasks int, blockBytes int, workers int) byte {
	jobs := make(chan int, workers*2)
	results := make(chan [32]byte, workers*2)
	var wg sync.WaitGroup

	// Запускаем воркеров
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for id := range jobs {
				results <- doTask(id, blockBytes)
			}
		}()
	}

	// Отправляем задачи
	go func() {
		for i := 0; i < tasks; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	// Ожидаем завершения воркеров и закрываем канал результатов
	go func() {
		wg.Wait()
		close(results)
	}()

	// Собираем результаты
	var sink byte
	for s := range results {
		sink ^= s[0]
	}

	return sink
}

// doTask имитация тяжёлой CPU-задачи: генерим блок и хэшируем его
func doTask(id, blockBytes int) [32]byte {
	r := rand.New(rand.NewSource(int64(id)))
	buf := make([]byte, blockBytes)
	for i := range buf {
		buf[i] = byte(r.Intn(256))
	}
	return sha256.Sum256(buf)
}
