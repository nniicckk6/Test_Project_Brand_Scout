package quotes

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// SQLiteRepo реализует хранение цитат в SQLite базе данных.
type SQLiteRepo struct {
	db *sql.DB
}

// NewSQLiteRepository открывает SQLite базу по указанному пути, создаёт таблицу при необходимости и возвращает репозиторий.
func NewSQLiteRepository(path string) Repository {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		panic(err)
	}
	stmt := `
	CREATE TABLE IF NOT EXISTS quotes (
	    id INTEGER PRIMARY KEY,
	    author TEXT,
	    quote TEXT
	);`
	if _, err := db.Exec(stmt); err != nil {
		panic(err)
	}
	return &SQLiteRepo{db: db}
}

// Create сохраняет новую цитату в базе, возвращает её ID или ошибку.
func (r *SQLiteRepo) Create(q Quote) (int, error) {
	res, err := r.db.Exec("INSERT INTO quotes(author, quote) VALUES(?, ?)", q.Author, q.Text)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

// GetAll возвращает все цитаты из таблицы quotes.
func (r *SQLiteRepo) GetAll() ([]Quote, error) {
	rows, err := r.db.Query("SELECT id, author, quote FROM quotes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Quote
	for rows.Next() {
		var q Quote
		if err := rows.Scan(&q.ID, &q.Author, &q.Text); err != nil {
			return nil, err
		}
		list = append(list, q)
	}
	return list, nil
}

// GetRandom возвращает случайную цитату из базы или ошибку, если записей нет.
func (r *SQLiteRepo) GetRandom() (Quote, error) {
	var q Quote
	err := r.db.QueryRow("SELECT id, author, quote FROM quotes ORDER BY RANDOM() LIMIT 1").
		Scan(&q.ID, &q.Author, &q.Text)
	return q, err
}

// FilterByAuthor возвращает цитаты указанного автора из базы.
func (r *SQLiteRepo) FilterByAuthor(author string) ([]Quote, error) {
	rows, err := r.db.Query("SELECT id, author, quote FROM quotes WHERE author = ?", author)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Quote
	for rows.Next() {
		var q Quote
		if err := rows.Scan(&q.ID, &q.Author, &q.Text); err != nil {
			return nil, err
		}
		list = append(list, q)
	}
	return list, nil
}

// Delete удаляет цитату по ID, возвращает ошибку при неудаче.
func (r *SQLiteRepo) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM quotes WHERE id = ?", id)
	return err
}
