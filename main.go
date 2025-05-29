package main

import (
	"Test_Project_Brand_Scout/internal/config"
	"Test_Project_Brand_Scout/internal/quotes"
	"bytes"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

// newRouter создаёт маршрутизатор с middleware и хендлерами на основе конфигурации
func newRouter(cfg config.Config) *mux.Router {
	var repo quotes.Repository
	if cfg.DBMode == "sqlite" {
		repo = quotes.NewSQLiteRepository(cfg.DBPath)
	} else {
		repo = quotes.NewMemoryRepository()
	}
	svc := quotes.NewService(repo)
	router := mux.NewRouter()
	router.Use(loggingMiddleware)
	quotes.RegisterHandlers(router, svc)
	return router
}

func main() {
	cfg := config.Load()
	log.Printf("Используется режим хранения: %s", cfg.DBMode)

	router := newRouter(cfg)

	log.Println("Сервер запущен на порту", cfg.Port)
	log.Fatal(http.ListenAndServe(
		":"+cfg.Port,
		router,
	))
}

// loggingResponseWriter оборачивает http.ResponseWriter для перехвата статуса и тела ответа и сохранения их для последующего логирования
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

// WriteHeader сохраняет код статуса и передаёт его оригинальному ResponseWriter
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Write записывает тело ответа в буфер и отправляет его через оригинальный ResponseWriter
func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	lrw.body.Write(b)
	return lrw.ResponseWriter.Write(b)
}

// loggingMiddleware создаёт middleware для логирования каждого HTTP-запроса и ответа
// Лог содержит метод, URI, тело запроса, статус и тело ответа
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody []byte
		if r.Body != nil {
			reqBody, _ = io.ReadAll(r.Body)
		}
		r.Body = io.NopCloser(bytes.NewBuffer(reqBody))

		log.Printf("Получен запрос: %s %s, тело: %s", r.Method, r.RequestURI, string(reqBody))
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lrw, r)
		log.Printf("Отправлен ответ: статус %d, тело: %s", lrw.statusCode, lrw.body.String())
	})
}
