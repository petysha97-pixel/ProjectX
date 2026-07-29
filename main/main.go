package main

import (
	"HOTA/internal/handlers"
	"HOTA/internal/models"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

// Поэтому на бэкенде обязательно нужно делать проверку прав (авторизацию): имеет ли право текущий вошедший пользователь смотреть данные профиля с этим ID.
// 1. Создаем функцию Middleware для CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем фронтенду доступ
		w.Header().Set("Access-Control-Allow-Origin", "*") //* - разрешен доступ со всех доменов
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Если браузер делает предварительную проверку (Preflight request)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Передаем запрос дальше в наши http.HandleFunc
		next.ServeHTTP(w, r)
	})
}

func SaveMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем фронтенду доступ
		w.Header().Set("Access-Control-Allow-Origin", "*") //* - разрешен доступ со всех доменов
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method != http.MethodPost {
			w.WriteHeader(405)
			return
		}

		// Если браузер делает предварительную проверку (Preflight request)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Передаем запрос дальше в наши http.HandleFunc
		next.ServeHTTP(w, r)
	})
}

func main() {

	db, err := sql.Open("sqlite", "../db.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	models.UserDB = db
	fmt.Println("Connected to SQLite")

	mux := http.NewServeMux()

	mux.Handle("/user", SaveMiddleware(http.HandlerFunc(handlers.NewUser)))
	mux.HandleFunc("/stack", handlers.GetStacks)
	mux.HandleFunc("/profile/{id}", handlers.GetUser)
	mux.HandleFunc("/getusers", handlers.GetAllUser)
	mux.HandleFunc("/users/searche", handlers.SearcheUsers)
	// http.HandleFunc("/profile", profile.Profile)

	// 3. Оборачиваем весь роутер в наше CORS Middleware
	fmt.Println("Сервер запущен: 8080")
	http.ListenAndServe(":8080", corsMiddleware(mux))

}

// ЗАДАЧИ
// 1. написать 2 функции для уникальности логина и никнейма
// 2. Прописать логику сохранения стека в БД + будем брать из БД и отображать на главной странице

// sql.Open - db, err := sql.Open("sqlite", "app.db") Открывает БД.

// db.Exec Выполняет запрос без результата.
// _, err := db.Exec(`
// CREATE TABLE users(
//     id INTEGER PRIMARY KEY,
//     name TEXT
// )
// `)

// db.Query  Возвращает много строк.
// rows, err := db.Query(
//     "SELECT id, name FROM users",
// )

// db.QueryRow Возвращает одну строку.
// var name string

// err := db.QueryRow(
//     "SELECT name FROM users WHERE id = ?",
//     1,
// ).Scan(&name)

// rows.Next - Переход к следующей записи.
// for rows.Next() {
// }

// rows.Scan - Чтение данных.
// var id int
// var name string

// rows.Scan(&id, &name)

// db.Prepare - Подготовленный запрос.
// stmt, err := db.Prepare(
//     "INSERT INTO users(name) VALUES(?)",
// )

// stmt.Exec("John")
// stmt.Exec("Kate")
// stmt.Exec("Alex")
