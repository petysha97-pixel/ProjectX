package repositories

import (
	"HOTA/internal/models"
)


//отправляем запрос в БД для взятия всего стека на фронт-регистрации
func GetStacks() ([]models.UserStackID, error) {

	rows, err := models.UserDB.Query(`SELECT id, name FROM stacks`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stacks []models.UserStackID

	for rows.Next() {

		var stak models.UserStackID
		err := rows.Scan(&stak.ID, &stak.Name)
		if err != nil {
			return nil, err
		}
		stacks = append(stacks, stak)
	}

	return stacks, nil

}
