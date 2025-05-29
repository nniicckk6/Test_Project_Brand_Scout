package quotes

import (
	"errors"
	"math/rand"
	"sync"
)

// MemoryRepo реализует хранение цитат в памяти с поддержкой конкурентного доступа.
type MemoryRepo struct {
	mu     sync.RWMutex
	data   []Quote
	nextID int
}

// NewMemoryRepository создаёт и возвращает новый in-memory репозиторий цитат.
func NewMemoryRepository() Repository {
	return &MemoryRepo{data: make([]Quote, 0), nextID: 1}
}

// Create сохраняет новую цитату в памяти и возвращает её ID или ошибку.
func (r *MemoryRepo) Create(q Quote) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	q.ID = r.nextID
	r.nextID++
	r.data = append(r.data, q)
	return q.ID, nil
}

// GetAll возвращает копию списка всех сохранённых цитат.
func (r *MemoryRepo) GetAll() ([]Quote, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]Quote(nil), r.data...), nil
}

// GetRandom возвращает случайную цитату или ошибку, если список пуст.
func (r *MemoryRepo) GetRandom() (Quote, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if len(r.data) == 0 {
		return Quote{}, errors.New("нету доступных цитат")
	}
	return r.data[rand.Intn(len(r.data))], nil
}

// FilterByAuthor возвращает все цитаты указанного автора.
func (r *MemoryRepo) FilterByAuthor(author string) ([]Quote, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []Quote
	for _, q := range r.data {
		if q.Author == author {
			res = append(res, q)
		}
	}
	return res, nil
}

// Delete удаляет цитату по ID, возвращает ошибку, если цитата не найдена.
func (r *MemoryRepo) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, q := range r.data {
		if q.ID == id {
			r.data = append(r.data[:i], r.data[i+1:]...)
			return nil
		}
	}
	return errors.New("нет такой цитаты")
}
