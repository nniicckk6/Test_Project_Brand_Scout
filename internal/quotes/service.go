package quotes

import "errors"

// ErrInvalidInput возвращается, если автор или текст цитаты пустые.
var ErrInvalidInput = errors.New("Поля не должны быть пустыми (Автор и цитата)")

// Service предоставляет бизнес-логику работы с цитатами, используя репозиторий для хранения данных.
type Service struct {
	repo Repository
}

// NewService создаёт новый Service с указанным репозиторием.
func NewService(r Repository) *Service {
	return &Service{repo: r}
}

// Create добавляет новую цитату с указанным автором и текстом.
// Возвращает идентификатор созданной цитаты или ошибку при невалидном вводе.
func (s *Service) Create(author, text string) (int, error) {
	if author == "" || text == "" {
		return 0, ErrInvalidInput
	}
	return s.repo.Create(Quote{Author: author, Text: text})
}

// GetAll возвращает список всех сохранённых цитат или ошибку.
func (s *Service) GetAll() ([]Quote, error) {
	return s.repo.GetAll()
}

// GetRandom возвращает случайную цитату или ошибку, если цитаты отсутствуют.
func (s *Service) GetRandom() (Quote, error) {
	return s.repo.GetRandom()
}

// FilterByAuthor возвращает все цитаты указанного автора или ошибку.
func (s *Service) FilterByAuthor(author string) ([]Quote, error) {
	return s.repo.FilterByAuthor(author)
}

// Delete удаляет цитату по заданному идентификатору.
// Возвращает ошибку, если цитата не найдена.
func (s *Service) Delete(id int) error {
	return s.repo.Delete(id)
}
