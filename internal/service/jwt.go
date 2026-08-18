package service

import (
	"HOTA/internal/models"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
)

// создаем токен при авторизации
func CreatJWT(id int) (*models.AuotOut, error) {
	// Читаем секрет из переменных окружения
	secretStr := os.Getenv("SECRET_JWT")
	if secretStr == "" {
		return nil, fmt.Errorf("критическая ошибка: SECRET_JWT не задан в .env")
	}

	// Переводим строковый секрет в массив байт
	var secretJWT = []byte(secretStr)

	exp := time.Now().Add(time.Hour * 24)

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(exp),
		Issuer:    "balance",
		Subject:   strconv.Itoa(id),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(secretJWT)
	if err != nil {
		return nil, fmt.Errorf("Ошибка создания токена %w", err)
	}

	DTO := models.AuotOut{
		Token: ss,
	}

	return &DTO, nil
}

// валидация токена
func ValidateToken(tokenString string) (string, error) {
	// Читаем секрет из переменных окружения
	secretStr := os.Getenv("SECRET_JWT")
	if secretStr == "" {
		return "", fmt.Errorf("критическая ошибка: SECRET_JWT не задан в .env")
	}

	// Переводим строковый секрет в массив байт
	var secretJWT = []byte(secretStr)

	// Парсим токен, используя структуру jwt.RegisteredClaims (так как вы её указали при создании)
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		// ОБЯЗАТЕЛЬНАЯ ПРОВЕРКА: убеждаемся, что алгоритм именно HMAC (HS256)
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", t.Header["alg"])
		}
		// Возвращаем секрет для проверки подписи
		return secretJWT, nil
	})

	if err != nil {
		return "", fmt.Errorf("токен невалиден или просрочен: %w", err)
	}

	// Проверяем валидность и достаем Claims
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		// Так как вы записали ID пользователя в поле Subject: Subject: strconv.Itoa(id),
		// мы достаем его обратно из поля Subject (оно будет строкой)
		userID := claims.Subject
		return userID, nil
	}

	return "", fmt.Errorf("не удалось прочитать данные из токена")
}
