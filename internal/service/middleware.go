package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	_ "modernc.org/sqlite"
)

// Поэтому на бэкенде обязательно нужно делать проверку прав (авторизацию): имеет ли право текущий вошедший пользователь смотреть данные профиля с этим ID.
// 1. Создаем функцию Middleware для CORS
func CORSMiddleware(next http.Handler) http.Handler {
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

func POSTMiddleware(next http.Handler) http.Handler {
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

func GETMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем фронтенду доступ
		w.Header().Set("Access-Control-Allow-Origin", "*") //* - разрешен доступ со всех доменов
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method != http.MethodGet {
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

func PutPatchMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем фронтенду доступ
		w.Header().Set("Access-Control-Allow-Origin", "*") //* - разрешен доступ со всех доменов
		w.Header().Set("Access-Control-Allow-Methods", "PUT, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method != http.MethodPut && r.Method != http.MethodPatch {
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

func DeleteMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем фронтенду доступ
		w.Header().Set("Access-Control-Allow-Origin", "*") //* - разрешен доступ со всех доменов
		w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method != http.MethodDelete {
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

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			http.Error(w, "Отсутствует токен", http.StatusForbidden)
			return
		}

		token := strings.TrimPrefix(tokenHeader, "Bearer ")

		userID, err := ValidateToken(token)
		if err != nil {
			fmt.Printf("ошибка валидации токена %v", err)
			http.Error(w, "ошибка валидации токена", http.StatusForbidden)
			return
		}

		// 2. КЛАДЕМ ID В КОНТЕКСТ ЗАПРОСА
		ctx := context.WithValue(r.Context(), "userID", userID)
		// Создаем новый запрос с этим контекстом
		rWithContext := r.WithContext(ctx)
		
		next.ServeHTTP(w, rWithContext)
	})
}
