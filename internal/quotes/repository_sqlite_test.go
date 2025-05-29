package quotes

import (
	"os"
	"testing"
)

// Тесты для SQLite-репозитория цитат
func TestSQLiteRepoCRUD(t *testing.T) {
	// Создание временного файла для SQLite-базы данных
	tmpfile, err := os.CreateTemp("", "quotes-*.db")
	if err != nil {
		t.Fatalf("не удалось создать временный файл: %v", err)
	}
	path := tmpfile.Name()
	tmpfile.Close()
	defer os.Remove(path)

	// Инициализация репозитория с созданной базой
	repo := NewSQLiteRepository(path)

	// Проверяем, что база изначально пуста
	all, err := repo.GetAll()
	if err != nil {
		t.Fatalf("ошибка получения цитат из SQLite: %v", err)
	}
	if len(all) != 0 {
		t.Errorf("ожидается отсутствие цитат в базе, получено %d", len(all))
	}

	// Создание двух цитат в базе
	id1, err := repo.Create(Quote{Author: "A", Text: "first"})
	if err != nil || id1 <= 0 {
		t.Fatalf("ошибка создания первой цитаты: id=%d, err=%v", id1, err)
	}
	id2, err := repo.Create(Quote{Author: "B", Text: "second"})
	if err != nil || id2 <= 0 {
		t.Fatalf("ошибка создания второй цитаты: id=%d, err=%v", id2, err)
	}

	// Проверяем, что GetAll возвращает обе цитаты
	all, err = repo.GetAll()
	if err != nil {
		t.Fatalf("ошибка получения цитат после создания: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("ожидается 2 цитаты, получено %d", len(all))
	}

	// Получение случайной цитаты из базы
	randQ, err := repo.GetRandom()
	if err != nil {
		t.Fatalf("ошибка получения случайной цитаты: %v", err)
	}
	if randQ.ID != id1 && randQ.ID != id2 {
		t.Errorf("GetRandom вернул некорректный ID %d", randQ.ID)
	}

	// Проверка фильтрации по автору
	listA, err := repo.FilterByAuthor("A")
	if err != nil {
		t.Fatalf("ошибка фильтрации по автору: %v", err)
	}
	if len(listA) != 1 || listA[0].Author != "A" {
		t.Errorf("ожидается одна цитата от автора A, получено %v", listA)
	}

	// Удаление несуществующей цитаты не должно приводить к ошибке
	err = repo.Delete(100)
	if err != nil {
		t.Errorf("ожидается отсутствие ошибки при удалении несуществующего ID, получено %v", err)
	}

	// Удаление существующей цитаты и проверка оставшегося
	err = repo.Delete(id1)
	if err != nil {
		t.Fatalf("ошибка удаления существующей цитаты: %v", err)
	}
	after, _ := repo.GetAll()
	if len(after) != 1 || after[0].ID != id2 {
		t.Errorf("ожидается оставшаяся цитата с ID=%d, получено %v", id2, after)
	}
}

// TestSQLiteRepoRandomEmpty проверяет, что GetRandom возвращает ошибку при отсутствии записей
func TestSQLiteRepoRandomEmpty(t *testing.T) {
	repo := NewSQLiteRepository(":memory:")
	// Попытка получения случайной цитаты из пустой базы должна привести к ошибке
	_, err := repo.GetRandom()
	if err == nil {
		t.Errorf("ожидается ошибка при получении случайной цитаты из пустой базы, получен nil")
	}
}
