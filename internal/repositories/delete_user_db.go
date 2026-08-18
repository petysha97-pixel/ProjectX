package repositories

import (
	"HOTA/internal/models"
	"fmt"
)

func DeleteUser(id int) error {

	deletestack := "DELETE FROM user_stacks WHERE user_id = ?"
	_, err := models.UserDB.Exec(deletestack, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления стеков при удалении пользователя: %w", err)
	}

	deleteuser := "DELETE FROM users WHERE id = ?"
	_, err = models.UserDB.Exec(deleteuser, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении пользователя: %w", err)
	}

	return nil
}
