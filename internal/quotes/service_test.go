// Тесты для бизнес-логики сервиса цитат
package quotes

import (
	"errors"
	"testing"
)

// TestServiceCreateInvalid проверяет, что при пустом авторе или тексте возвращается ошибка валидации
func TestServiceCreateInvalid(t *testing.T) {
	svc := NewService(NewMemoryRepository())
	// Пустой автор
	_, err := svc.Create("", "текст")
	if !errors.Is(err, ErrInvalidInput) {
		t.Errorf("ожидается ErrInvalidInput для пустого автора, получено %v", err)
	}
	// Пустой текст цитаты
	_, err = svc.Create("автор", "")
	if !errors.Is(err, ErrInvalidInput) {
		t.Errorf("ожидается ErrInvalidInput для пустого текста, получено %v", err)
	}
}

// TestServiceCRUD проверяет основные CRUD-операции через сервис
func TestServiceCRUD(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo)

	// --- Создание цитат ---
	id1, err := svc.Create("АвторA", "Первая цитата")
	if err != nil || id1 != 1 {
		t.Fatalf("ошибка создания первой цитаты: id=%d, err=%v", id1, err)
	}

	id2, err := svc.Create("АвторB", "Вторая цитата")
	if err != nil || id2 != 2 {
		t.Fatalf("ошибка создания второй цитаты: id=%d, err=%v", id2, err)
	}

	// --- Получение всех цитат ---
	all, err := svc.GetAll()
	if err != nil {
		t.Fatalf("ошибка GetAll: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("ожидается 2 цитаты, получено %d", len(all))
	}

	// --- Фильтрация по автору ---
	filtered, err := svc.FilterByAuthor("АвторA")
	if err != nil {
		t.Fatalf("ошибка FilterByAuthor: %v", err)
	}
	if len(filtered) != 1 || filtered[0].Author != "АвторA" {
		t.Errorf("ожидается один результат по АвторA, получено %v", filtered)
	}

	// Пустой результат фильтрации
	filteredEmpty, err := svc.FilterByAuthor("АвторX")
	if err != nil {
		t.Fatalf("ошибка FilterByAuthor для несуществующего автора: %v", err)
	}
	if len(filteredEmpty) != 0 {
		t.Errorf("ожидается 0 результатов, получено %v", filteredEmpty)
	}

	// --- Случайная цитата ---
	randQ, err := svc.GetRandom()
	if err != nil {
		t.Fatalf("ошибка GetRandom: %v", err)
	}
	if randQ.ID != 1 && randQ.ID != 2 {
		t.Errorf("неожиданный ID случайной цитаты: %d", randQ.ID)
	}

	// --- Удаление несуществующей цитаты ---
	err = svc.Delete(100)
	if err == nil {
		t.Error("ожидается ошибка при удалении несуществующей цитаты")
	}

	// --- Удаление существующей цитаты ---
	err = svc.Delete(id1)
	if err != nil {
		t.Errorf("не удалось удалить цитату id=%d: %v", id1, err)
	}
	after, _ := svc.GetAll()
	if len(after) != 1 || after[0].ID != 2 {
		t.Errorf("ожидается, что останется цитата id=2, получено %v", after)
	}
}
