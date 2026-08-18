package service

import (
	"HOTA/internal/models"
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func ValidateNewEmailOnly(newEmail string) error {
	data := models.UpdateEmail{
		Email: newEmail,
	}

	return validation.ValidateStruct(&data, validation.Field(&data.Email, validation.Required,
		is.Email,
		validation.Length(1, 30),
		validation.By(UNIK_email),
	),
	)
}

// Обновление email
func ValidateUpdateEmail(id int, newEmail, password string) error {
	if newEmail == "" {
		return errors.New("email не может быть пустым")
	}
	if password == "" {
		return errors.New("пароль не может быть пустым")
	}

	// 1. Получаем текущий email и хеш пароля из БД
	var currentEmail, storedHash string
	query := "SELECT Email, Password FROM users WHERE id = ?"
	err := models.UserDB.QueryRow(query, id).Scan(&currentEmail, &storedHash)
	if err != nil {
		return fmt.Errorf("ошибка получения данных пользователя: %w", err)
	}

	// 2. Проверяем пароль
	if err := CheckPassword(storedHash, password); err != nil {
		return errors.New("неверный пароль")
	}

	// 3. Проверяем что новый email не совпадает с текущим
	if newEmail == currentEmail {
		return errors.New("новый email совпадает с текущим")
	}

	// 4. Валидация нового email (формат + уникальность)
	if err := ValidateNewEmailOnly(newEmail); 
	err != nil {
		return fmt.Errorf("Ошибка валидации при обновлении нового Email: %w", err)
	}

	//Обновление email в БД
	updateQuery := "UPDATE users SET Email = ? WHERE id = ?"
	_, err = models.UserDB.Exec(updateQuery, newEmail, id)
	if err != nil {
		return fmt.Errorf("ошибка обновления email: %w", err)
	}

	return nil
}
