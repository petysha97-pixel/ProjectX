package main



// nota-profile.dev
// nota-dev.dev
// n-o-t-a.online
// n-o-t-a.dev
// END-NOTA
// end-nota.dev
// var Secret = "DF#!_ASDNBUJ@)_JSHDF2"

import (
	"HOTA/internal/handlers"
	"HOTA/internal/models"
	"HOTA/internal/service"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Ошибка в загрузке файла .env")
	}
	fmt.Println("Connected файла .env")

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

	//регистрация
	mux.Handle("/user", service.POSTMiddleware(http.HandlerFunc(handlers.NewUser)))
	mux.HandleFunc("/stack", handlers.GetStacks)

	//Авторизация 
	mux.Handle("/user/auth", service.POSTMiddleware(http.HandlerFunc(handlers.Auth)))

	//Главный профиль
	mux.Handle("/profile", service.JWTMiddleware(http.HandlerFunc((handlers.GetUser))))

	//CRUD
	mux.Handle("/user/update", service.JWTMiddleware(http.HandlerFunc(handlers.UpdateUser)))
    mux.Handle("/user/password", service.JWTMiddleware(http.HandlerFunc(handlers.UpdatePassword)))
	mux.Handle("/user/email", service.JWTMiddleware(http.HandlerFunc(handlers.UpdateEmail)))
	mux.Handle("/user/delete", service.JWTMiddleware(http.HandlerFunc(handlers.DeleteUser)))

	//"Умный поиск"
	mux.Handle("/users/searche", service.GETMiddleware(http.HandlerFunc(handlers.SearcheUsers)))

	// 3. Оборачиваем весь роутер в наше CORS Middleware
	fmt.Println("Сервер запущен: 8080")
	http.ListenAndServe(":8080", service.CORSMiddleware(mux))

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
