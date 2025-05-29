package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"Test_Project_Brand_Scout/internal/config"
)

// TestNewRouterMemory проверяет эндпоинты при хранении в памяти
func TestNewRouterMemory(t *testing.T) {
	cfg := config.Config{Port: "8080", DBMode: "memory"}
	r := newRouter(cfg)

	// GET /quotes должен вернуть пустой список
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/quotes", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("ожидается статус 200 OK, получен %d", w.Code)
	}
	body := w.Body.String()
	if body != "[]" && body != "[]\n" {
		t.Errorf("ожидается пустой JSON-массив, получено %s", body)
	}

	// GET /quotes/random на пустом должен вернуть 404 Not Found
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/quotes/random", nil))
	if w.Code != http.StatusNotFound {
		t.Errorf("ожидается 404 Not Found, получен %d", w.Code)
	}

	// POST /quotes создаёт новую цитату
	input := map[string]string{"author": "M", "quote": "Q"}
	buf, _ := json.Marshal(input)
	w = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewBuffer(buf))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("ожидается статус 201 Created, получен %d", w.Code)
	}
	var res map[string]int
	data, _ := io.ReadAll(w.Body)
	if err := json.Unmarshal(data, &res); err != nil {
		t.Fatalf("ошибка разбора JSON ответа: %v", err)
	}
	id, ok := res["id"]
	if !ok || id != 1 {
		t.Errorf("ожидается id=1, получен %d", id)
	}

	// GET /quotes должен вернуть массив с одним элементом
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/quotes", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("ожидается статус 200 OK, получен %d", w.Code)
	}
	var list []map[string]interface{}
	data, _ = io.ReadAll(w.Body)
	if err := json.Unmarshal(data, &list); err != nil {
		t.Fatalf("ошибка разбора JSON списка: %v", err)
	}
	if len(list) != 1 || list[0]["author"] != "M" {
		t.Errorf("ожидается одна цитата {author:M}, получено %v", list)
	}
}

// TestNewRouterSQLite проверяет эндпоинты при хранении в SQLite режиме с in-memory базой
func TestNewRouterSQLite(t *testing.T) {
	// Используем :memory: чтобы не трогать файл на диске
	cfg := config.Config{Port: "8080", DBMode: "sqlite", DBPath: ":memory:"}
	r := newRouter(cfg)

	// GET /quotes должен вернуть пустой список
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/quotes", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("SQLite: ожидается статус 200 OK, получен %d", w.Code)
	}
	body := w.Body.String()
	if body != "[]" && body != "[]\n" {
		t.Errorf("SQLite: ожидается пустой JSON-массив, получено %s", body)
	}

	// POST /quotes создаёт новую цитату
	input := map[string]string{"author": "S", "quote": "SqliteQ"}
	buf, _ := json.Marshal(input)
	w = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewBuffer(buf))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("SQLite: ожидается статус 201 Created, получен %d", w.Code)
	}

	// GET /quotes должен вернуть массив с созданной цитатой
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/quotes", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("SQLite: ожидается статус 200 OK, получен %d", w.Code)
	}
	var list []map[string]interface{}
	data, _ := io.ReadAll(w.Body)
	if err := json.Unmarshal(data, &list); err != nil {
		t.Fatalf("SQLite: ошибка разбора JSON списка: %v", err)
	}
	if len(list) != 1 || list[0]["author"] != "S" {
		t.Errorf("SQLite: ожидается одна цитата от S, получено %v", list)
	}
}
