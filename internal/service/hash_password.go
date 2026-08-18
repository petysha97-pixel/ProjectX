package service

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// хешируем пароль + солим
func Hash_password(password string) (string, error) {
	hashpassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("неудалось захешировать пароль %w", err)
	}

	return string(hashpassword), nil

}

// Проверка пароля с хешем из БД
func CheckPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return errors.New("неверный пароль")
	}
	return nil
}
