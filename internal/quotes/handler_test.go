// Тесты HTTP-хендлеров для работы с цитатами (создание, получение, фильтрация, удаление)
package quotes

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

// setupRouter создаёт маршрутизатор с in-memory сервисом цитат и регистрирует хендлеры
func setupRouter() *mux.Router {
	repo := NewMemoryRepository()
	svc := NewService(repo)
	r := mux.NewRouter()
	RegisterHandlers(r, svc)
	return r
}

// TestCreateQuoteHandler проверяет создание новой цитаты через POST /quotes
func TestCreateQuoteHandler(t *testing.T) {
	r := setupRouter()
	// valid request
	body := bytes.NewBufferString(`{"author":"A","quote":"Q"}`)
	req := httptest.NewRequest(http.MethodPost, "/quotes", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("ожидаемый статус 201 Created, получен %d", resp.StatusCode)
	}
	data, _ := ioutil.ReadAll(resp.Body)
	var res map[string]int
	if err := json.Unmarshal(data, &res); err != nil {
		t.Fatalf("ошибка разбора JSON ответа: %v", err)
	}
	if res["id"] != 1 {
		t.Errorf("ожидаемый id=1, получен %d", res["id"])
	}
}

// TestListAndRandomAndFilterDeleteHandlers проверяет последовательную логику работы эндпоинтов:
// получение списка, случайной цитаты, фильтрацию по автору и удаление
func TestListAndRandomAndFilterDeleteHandlers(t *testing.T) {
	r := setupRouter()
	// create two quotes
	for _, q := range []Quote{{Author: "A", Text: "one"}, {Author: "B", Text: "two"}} {
		b, _ := json.Marshal(q)
		req := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}

	// list all
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/quotes", nil))
	if w.Code != http.StatusOK {
		t.Errorf("ожидаемый статус 200 OK для списка, получен %d", w.Code)
	}
	var list []Quote
	if err := json.Unmarshal(w.Body.Bytes(), &list); err != nil {
		t.Errorf("ошибка разбора JSON списка: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("ожидалось 2 цитаты, получено %d", len(list))
	}

	// random
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/quotes/random", nil))
	if w.Code != http.StatusOK {
		t.Errorf("ожидаемый статус 200 OK для случайной цитаты, получен %d", w.Code)
	}
	var rnd Quote
	if err := json.Unmarshal(w.Body.Bytes(), &rnd); err != nil {
		t.Errorf("ошибка разбора JSON случайной цитаты: %v", err)
	}
	if rnd.ID < 1 || rnd.ID > 2 {
		t.Errorf("ID случайной цитаты вне диапазона: %d", rnd.ID)
	}

	// filter by author
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/quotes?author=A", nil))
	if w.Code != http.StatusOK {
		t.Errorf("ожидаемый статус 200 OK для фильтрации, получен %d", w.Code)
	}
	var fa []Quote
	if err := json.Unmarshal(w.Body.Bytes(), &fa); err != nil {
		t.Errorf("ошибка разбора JSON фильтрации: %v", err)
	}
	if len(fa) != 1 || fa[0].Author != "A" {
		t.Errorf("фильтрация по автору вернула %v", fa)
	}

	// delete one
	w = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/quotes/1", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Errorf("ожидаемый статус 204 No Content для удаления, получен %d", w.Code)
	}

	// confirm deletion
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/quotes", nil))
	var after []Quote
	if err := json.Unmarshal(w.Body.Bytes(), &after); err != nil {
		t.Errorf("ошибка разбора JSON после удаления: %v", err)
	}
	if len(after) != 1 || after[0].ID != 2 {
		t.Errorf("ожидался остаток цитаты с id=2, получено %v", after)
	}
}
