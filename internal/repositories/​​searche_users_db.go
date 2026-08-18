package repositories

import (
	"HOTA/internal/models"
	"fmt"
)

// ищет пользавателей по запросу юзера с главного меню
func SearcheUsersBD(query string) ([]models.User, error) {

	rules := "%" + query + "%"

	quer := `SELECT DISTINCT users.id, users.Nickname, users.Rolle FROM users 
	LEFT JOIN user_stacks
	ON users.id = user_stacks.user_id
	LEFT JOIN stacks 
	ON user_stacks.stack_id = stacks.id
	WHERE Nickname LIKE ?
	OR users.Rolle LIKE ?
	OR stacks.name LIKE ?
	  ORDER BY
	       ц
	
	
	`

	rows, err := models.UserDB.Query(quer, rules, rules, rules)
	if err != nil {
		return nil, fmt.Errorf("Ошибка запроса в БД: SearcheUsersBD %w", err)
	}

	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Nickname, &user.Rolle)
		if err != nil {
			return nil, fmt.Errorf("Ошибка сканирования пользователя из БД %w", err)
		}

		users = append(users, user)
	}

	return users, nil

}
