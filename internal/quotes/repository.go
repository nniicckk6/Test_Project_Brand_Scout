package quotes

// Quote представляет цитату с ID, автором и текстом.
// Поле ID содержит уникальный идентификатор, Author — автора цитаты, Text — сам текст цитаты.
type Quote struct {
	ID     int    `json:"id"`
	Author string `json:"author"`
	Text   string `json:"quote"`
}

// Repository описывает набор операций для хранения и получения цитат.
// Реализации могут хранить данные в памяти, базе данных и т.д.
type Repository interface {
	// Create сохраняет новую цитату и возвращает её ID или ошибку.
	Create(q Quote) (int, error)
	// GetAll возвращает все сохранённые цитаты или ошибку.
	GetAll() ([]Quote, error)
	// GetRandom возвращает случайную цитату или ошибку, если нечего вернуть.
	GetRandom() (Quote, error)
	// FilterByAuthor возвращает цитаты указанного автора.
	FilterByAuthor(author string) ([]Quote, error)
	// Delete удаляет цитату по ID, возвращает ошибку, если не найдена.
	Delete(id int) error
}
