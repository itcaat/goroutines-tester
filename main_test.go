package main

import (
	"testing"

	"github.com/itcaat/goroutines-tester/internal/benchmark"
)

func TestMainIntegration(t *testing.T) {
	// Интеграционный тест для проверки совместной работы компонентов
	runner := benchmark.NewRunner()

	testCases := []struct {
		name   string
		config benchmark.Config
	}{
		{
			name: "single mode small",
			config: benchmark.Config{
				Tasks:   3,
				BlockKB: 1,
				Mode:    "single",
				Workers: 1,
			},
		},
		{
			name: "pool mode small",
			config: benchmark.Config{
				Tasks:   5,
				BlockKB: 1,
				Mode:    "pool",
				Workers: 2,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := runner.Run(tc.config)
			// Проверяем, что получили какой-то результат
			t.Logf("Config: %+v, Result sink: %d", tc.config, result.Sink)
		})
	}
}
