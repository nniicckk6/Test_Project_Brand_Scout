package quotes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// RegisterHandlers регистрирует HTTP-эндпоинты для операций с цитатами: создание, получение списка, случайная цитата, фильтрация и удаление.
func RegisterHandlers(r *mux.Router, svc *Service) {
	r.HandleFunc("/quotes", createQuote(svc)).Methods("POST")
	r.HandleFunc("/quotes", filterQuotes(svc)).Queries("author", "{author}").Methods("GET")
	r.HandleFunc("/quotes", listQuotes(svc)).Methods("GET")
	r.HandleFunc("/quotes/random", randomQuote(svc)).Methods("GET")
	r.HandleFunc("/quotes/{id}", deleteQuote(svc)).Methods("DELETE")
}

// createQuote возвращает HandlerFunc для создания новой цитаты.
// Ожидает JSON с полями author и quote, возвращает id созданной цитаты.
func createQuote(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var q Quote
		if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := svc.Create(q.Author, q.Text)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]int{"id": id})
	}
}

// listQuotes возвращает HandlerFunc для получения списка всех цитат.
func listQuotes(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list, _ := svc.GetAll()
		if list == nil {
			list = []Quote{}
		}
		json.NewEncoder(w).Encode(list)
	}
}

// randomQuote возвращает HandlerFunc для получения случайной цитаты.
func randomQuote(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := svc.GetRandom()
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(q)
	}
}

// filterQuotes возвращает HandlerFunc для фильтрации цитат по автору.
// Принимает параметр URL author и возвращает все цитаты указанного автора.
func filterQuotes(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		author := mux.Vars(r)["author"]
		list, _ := svc.FilterByAuthor(author)
		json.NewEncoder(w).Encode(list)
	}
}

// deleteQuote возвращает HandlerFunc для удаления цитаты по заданному id.
// Возвращает 204 No Content при успешном удалении или 404, если цитата не найдена.
func deleteQuote(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := svc.Delete(id); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
