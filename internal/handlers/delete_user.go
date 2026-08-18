package handlers

import (
	"HOTA/internal/models"
	"HOTA/internal/repositories"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	// Достаем userID из контекста запроса
	// r.Context().Value возвращает тип interface{}, поэтому приводим его к string через .(string)
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Ошибка сервера: ID пользователя не найден", http.StatusInternalServerError)
		return
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Некорректный id пользователя для удаления", http.StatusNotFound)
		fmt.Printf("Не получилось конвертировать строку в число: %s", userID)
		return
	}

	//Провеока пользователя в БД по айди
	usr, err := repositories.GetUsersByID(id)
	if err != nil {
		fmt.Printf("ошибка поиск пользователя в БД: %v", err)
		http.Error(w, "Ошибка поиска пользователя в БД для удалания", http.StatusBadRequest)
		return
	}
	if usr == nil {
		fmt.Printf("ошибка поиск пользователя: %v", err)
		http.Error(w, "Пользователь не найден", http.StatusBadRequest)
		return
	}

	err = repositories.DeleteUser(id)
	if err != nil {
		fmt.Printf("ошибка удаления пользователя в БД: %v", err)
		http.Error(w, "Ошибка удаления пользователя", http.StatusBadRequest)
	}

	otvet := models.RegistrationResponseError{
		ID:           id,
		StatusGlobal: "Пользователь удален успешно",
	}

	w.Header().Set("Context-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(otvet)

}
