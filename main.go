package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"sync"
	"time"
)

// имитация тяжёлой CPU-задачи: генерим блок и хэшируем его
func doTask(id, blockBytes int) [32]byte {
	r := rand.New(rand.NewSource(int64(id)))
	buf := make([]byte, blockBytes)
	for i := range buf {
		buf[i] = byte(r.Intn(256))
	}
	return sha256.Sum256(buf)
}

func main() {
	var tasks int
	var blockKB int
	var mode string
	var workers int

	f, _ := os.Create("cpu.out")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	tr, _ := os.Create("trace.out")
	trace.Start(tr)
	defer trace.Stop()

	flag.IntVar(&tasks, "tasks", 200, "сколько задач выполнить")
	flag.IntVar(&blockKB, "blockKB", 1024, "размер блока данных на задачу (KB)")
	flag.StringVar(&mode, "mode", "single", "режим: single | pool")
	flag.IntVar(&workers, "workers", runtime.NumCPU(), "кол-во воркеров для режима pool")
	flag.Parse()

	if workers < 1 {
		workers = 1
	}
	blockBytes := blockKB * 1024
	runtime.GOMAXPROCS(workers)

	fmt.Printf("mode=%s tasks=%d block=%dKB workers=%d\n", mode, tasks, blockKB, workers)
	start := time.Now()

	var sink byte

	switch mode {
	case "single":
		for i := 0; i < tasks; i++ {
			sum := doTask(i, blockBytes)
			sink ^= sum[0]
		}

	case "pool":
		jobs := make(chan int, workers*2)
		results := make(chan [32]byte, workers*2)
		var wg sync.WaitGroup

		wg.Add(workers)
		for w := 0; w < workers; w++ {
			go func() {
				defer wg.Done()
				for id := range jobs {
					results <- doTask(id, blockBytes)
				}
			}()
		}

		go func() {
			for i := 0; i < tasks; i++ {
				jobs <- i
			}
			close(jobs)
		}()

		go func() {
			wg.Wait()
			close(results)
		}()

		for s := range results {
			sink ^= s[0]
		}

	default:
		fmt.Println("unknown mode")
		return
	}

	elapsed := time.Since(start)
	fmt.Printf("done in %v (sink=%d)\n", elapsed, sink)
}
