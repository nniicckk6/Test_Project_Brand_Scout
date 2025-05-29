// Тесты для репозитория цитат в памяти
package quotes

import (
	"testing"
)

func TestMemoryRepoCRUD(t *testing.T) {
	repo := NewMemoryRepository()

	// Проверяем, что изначально нет цитат в репозитории
	all, err := repo.GetAll()
	if err != nil {
		t.Fatalf("ошибка получения списка цитат: %v", err)
	}
	if len(all) != 0 {
		t.Errorf("ожидается отсутствие цитат, получено %d", len(all))
	}

	// Создаём две цитаты в репозитории
	id1, err := repo.Create(Quote{Author: "A", Text: "first"})
	if err != nil || id1 != 1 {
		t.Fatalf("ошибка при создании первой цитаты: id=%d, err=%v", id1, err)
	}
	id2, err := repo.Create(Quote{Author: "B", Text: "second"})
	if err != nil || id2 != 2 {
		t.Fatalf("ошибка при создании второй цитаты: id=%d, err=%v", id2, err)
	}

	// Проверяем, что GetAll возвращает обе цитаты
	all, err = repo.GetAll()
	if err != nil {
		t.Fatalf("ошибка получения списка цитат после создания: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("ожидается 2 цитаты, получено %d", len(all))
	}

	// GetRandom должен вернуть одну из созданных цитат
	randQ, err := repo.GetRandom()
	if err != nil {
		t.Fatalf("ошибка получения случайной цитаты: %v", err)
	}
	if randQ.ID != 1 && randQ.ID != 2 {
		t.Errorf("GetRandom вернул некорректный ID %d", randQ.ID)
	}

	// Проверяем фильтрацию по автору
	listA, err := repo.FilterByAuthor("A")
	if err != nil {
		t.Fatalf("ошибка фильтрации по автору: %v", err)
	}
	if len(listA) != 1 || listA[0].Author != "A" {
		t.Errorf("ожидается одна цитата от автора A, получено %v", listA)
	}

	// Попытка удалить несуществующую цитату должна вернуть ошибку
	err = repo.Delete(100)
	if err == nil {
		t.Error("ожидается ошибка при удалении несуществующего ID")
	}

	// Удаляем существующую цитату и проверяем результат
	err = repo.Delete(1)
	if err != nil {
		t.Fatalf("ошибка удаления существующей цитаты: %v", err)
	}
	after, _ := repo.GetAll()
	if len(after) != 1 || after[0].ID != 2 {
		t.Errorf("ожидается оставшаяся цитата с ID=2, получено %v", after)
	}
}

// TestMemoryRepoRandomEmpty проверяет, что GetRandom возвращает ошибку при отсутствии цитат
func TestMemoryRepoRandomEmpty(t *testing.T) {
	repo := NewMemoryRepository()
	// Попытка получения случайной цитаты из пустого репозитория должна вернуть ошибку
	_, err := repo.GetRandom()
	if err == nil {
		t.Error("ожидается ошибка при получении случайной цитаты из пустого репозитория")
	}
}
