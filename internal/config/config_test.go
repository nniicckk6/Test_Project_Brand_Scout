package config

import (
	"os"
	"testing"
)

// withEnv очищает и восстанавливает переменную окружения (ключ, значение)
func withEnv(key, value string, fn func()) {
	orig, had := os.LookupEnv(key)
	_ = os.Setenv(key, value)
	fn()
	if had {
		_ = os.Setenv(key, orig)
	} else {
		_ = os.Unsetenv(key)
	}
}

// TestLoadDefaults проверяет значения по умолчанию при отсутствии переменных окружения
func TestLoadDefaults(t *testing.T) {
	// очищаем переменные окружения
	_ = os.Unsetenv("PORT")
	_ = os.Unsetenv("DB_MODE")

	cfg := Load()
	if cfg.Port != "8080" {
		t.Errorf("ожидаемый порт по умолчанию 8080, получен %s", cfg.Port)
	}
	if cfg.DBMode != "memory" {
		t.Errorf("ожидаемый режим хранения 'memory', получен %s", cfg.DBMode)
	}
}

// TestLoadFromEnv проверяет корректную загрузку из переменных окружения
func TestLoadFromEnv(t *testing.T) {
	withEnv("PORT", "9090", func() {
		withEnv("DB_MODE", "sqlite", func() {
			cfg := Load()
			if cfg.Port != "9090" {
				t.Errorf("ожидаемый порт 9090, получен %s", cfg.Port)
			}
			if cfg.DBMode != "sqlite" {
				t.Errorf("ожидаемый режим хранения 'sqlite', получен %s", cfg.DBMode)
			}
		})
	})
}
