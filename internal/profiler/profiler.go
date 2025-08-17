package profiler

import (
	"fmt"
	"os"
	"runtime/pprof"
	"runtime/trace"
)

// Profiler управляет профилированием приложения
type Profiler struct {
	cpuFile   *os.File
	traceFile *os.File
	enabled   bool
}

// New создает новый Profiler
func New() *Profiler {
	return &Profiler{}
}

// Start запускает профилирование (CPU и trace)
func (p *Profiler) Start() error {
	if p.enabled {
		return fmt.Errorf("профилирование уже запущено")
	}

	var err error

	// Создаем файл для CPU профиля
	p.cpuFile, err = os.Create("cpu.out")
	if err != nil {
		return fmt.Errorf("ошибка создания cpu.out: %w", err)
	}

	// Запускаем CPU профилирование
	err = pprof.StartCPUProfile(p.cpuFile)
	if err != nil {
		p.cpuFile.Close()
		return fmt.Errorf("ошибка запуска CPU профилирования: %w", err)
	}

	// Создаем файл для trace
	p.traceFile, err = os.Create("trace.out")
	if err != nil {
		pprof.StopCPUProfile()
		p.cpuFile.Close()
		return fmt.Errorf("ошибка создания trace.out: %w", err)
	}

	// Запускаем trace
	err = trace.Start(p.traceFile)
	if err != nil {
		pprof.StopCPUProfile()
		p.cpuFile.Close()
		p.traceFile.Close()
		return fmt.Errorf("ошибка запуска trace: %w", err)
	}

	p.enabled = true
	fmt.Println("Debug режим включен: профили будут сохранены в cpu.out и trace.out")
	return nil
}

// Stop останавливает профилирование и закрывает файлы
func (p *Profiler) Stop() {
	if !p.enabled {
		return
	}

	// Останавливаем профилирование
	trace.Stop()
	pprof.StopCPUProfile()

	// Закрываем файлы
	if p.traceFile != nil {
		p.traceFile.Close()
	}
	if p.cpuFile != nil {
		p.cpuFile.Close()
	}

	p.enabled = false
}
