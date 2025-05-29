// internal/config/config.go
package config

import (
	"bufio"
	"os"
	"strings"
)

// Config содержит параметры конфигурации приложения.
type Config struct {
	Port   string // порт HTTP-сервера
	DBMode string // режим хранения (sqlite или memory)
	DBPath string // путь к SQLite базе (файл или :memory:)
}

// Load загружает конфигурацию из .env и переменных окружения.
// Если значение не задано, используются значения по умолчанию: PORT=8080, DB_MODE=memory.
func Load() Config {
	loadDotEnv()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	mode := os.Getenv("DB_MODE")
	if mode == "" {
		mode = "memory"
	}
	// Путь к SQLite базе из переменной окружения или quotes.db по умолчанию
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "quotes.db"
	}
	return Config{Port: port, DBMode: mode, DBPath: dbPath}
}

// loadDotEnv читает файл .env в корне проекта и устанавливает переменные окружения.
func loadDotEnv() {
	file, err := os.Open(".env")
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		os.Setenv(key, val)
	}
}
