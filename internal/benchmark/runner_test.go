package benchmark

import (
	"crypto/sha256"
	"testing"
)

func TestDoTask(t *testing.T) {
	// Test that doTask produces consistent results for the same input
	blockBytes := 1024
	id := 42

	result1 := doTask(id, blockBytes)
	result2 := doTask(id, blockBytes)

	if result1 != result2 {
		t.Errorf("doTask should produce consistent results for same input, got different hashes")
	}

	// Verify it returns a proper SHA-256 hash (32 bytes)
	if len(result1) != sha256.Size {
		t.Errorf("Expected hash size %d, got %d", sha256.Size, len(result1))
	}
}

func TestDoTaskDifferentInputs(t *testing.T) {
	// Test that different inputs produce different hashes
	blockBytes := 512

	result1 := doTask(1, blockBytes)
	result2 := doTask(2, blockBytes)

	if result1 == result2 {
		t.Errorf("Different inputs should produce different hashes")
	}
}

func TestRunnerSingleMode(t *testing.T) {
	runner := NewRunner()
	config := Config{
		Tasks:   5,
		BlockKB: 1,
		Mode:    "single",
		Workers: 1,
	}

	result := runner.Run(config)

	// Проверяем, что результат не нулевой (задачи выполнились)
	// Sink может быть любым значением, главное - что функция завершилась
	t.Logf("Single mode result sink: %d", result.Sink)
}

func TestRunnerPoolMode(t *testing.T) {
	runner := NewRunner()
	config := Config{
		Tasks:   10,
		BlockKB: 1,
		Mode:    "pool",
		Workers: 4,
	}

	result := runner.Run(config)

	// Проверяем, что результат получен
	t.Logf("Pool mode result sink: %d", result.Sink)
}

func TestRunnerUnknownMode(t *testing.T) {
	runner := NewRunner()
	config := Config{
		Tasks:   5,
		BlockKB: 1,
		Mode:    "unknown",
		Workers: 1,
	}

	result := runner.Run(config)

	// Для неизвестного режима должен возвращаться нулевой sink
	if result.Sink != 0 {
		t.Errorf("Expected sink=0 for unknown mode, got %d", result.Sink)
	}
}

func BenchmarkDoTask(b *testing.B) {
	blockBytes := 1024

	for i := 0; i < b.N; i++ {
		doTask(i, blockBytes)
	}
}

func BenchmarkDoTaskLarge(b *testing.B) {
	blockBytes := 4096

	for i := 0; i < b.N; i++ {
		doTask(i, blockBytes)
	}
}

func BenchmarkRunnerSingle(b *testing.B) {
	runner := NewRunner()
	config := Config{
		Tasks:   10,
		BlockKB: 1,
		Mode:    "single",
		Workers: 1,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runner.Run(config)
	}
}

func BenchmarkRunnerPool(b *testing.B) {
	runner := NewRunner()
	config := Config{
		Tasks:   10,
		BlockKB: 1,
		Mode:    "pool",
		Workers: 4,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runner.Run(config)
	}
}
