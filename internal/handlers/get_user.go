package handlers

import (
	"HOTA/internal/models"
	"HOTA/internal/repositories"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// получаем нашего пользователя для отображения данных в главном меню
func GetUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}


	// Достаем userID из контекста запроса
	// r.Context().Value возвращает тип interface{}, поэтому приводим его к string через .(string)
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Ошибка сервера: ID пользователя не найден", http.StatusInternalServerError)
		return
	}

	

	id, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Некорректный id", http.StatusNotFound)
		fmt.Printf("Не получилось конвертировать строку в число: %s", userID)
		return
	}

	// Получаем юзера
	user := repositories.Get_userdb(id)
	// Получаем его стеки
	stack, err := repositories.GetStacksByUserID(id)
	if err != nil {
		fmt.Printf("Ошибка получения стеков пользователя %v", err)
		http.Error(w, "Ошибка получения стеков пользователя", http.StatusInternalServerError)
		return
	}
	// Собираем респонс
	reposonse := models.UserResponse{
		ID:       user.ID,
		Nickname: user.Nickname,
		Rolle:    user.Rolle,
		Stack:    stack,
	}

	w.Header().Set("Context-Type", "application/json")
	json.NewEncoder(w).Encode(reposonse)
}
