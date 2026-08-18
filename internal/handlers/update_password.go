package handlers

import (
	"HOTA/internal/models"
	"HOTA/internal/repositories"
	"HOTA/internal/service"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"strconv"
)

func UpdatePassword(w http.ResponseWriter, r *http.Request) {

	// Достаем userID из контекста запроса
	// r.Context().Value возвращает тип interface{}, поэтому приводим его к string через .(string)
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Ошибка сервера: ID пользователя не найден", http.StatusInternalServerError)
		return
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Некорректный id пользователя для обновления пароля", http.StatusNotFound)
		fmt.Printf("Не получилось конвертировать строку в число: %s", userID)
		return
	}

	//Проверка пользователя в БД по айди
	usr, err := repositories.GetUsersByID(id)
	if err != nil {
		fmt.Printf("ошибка поиск пользователя в БД: %v", err)
		http.Error(w, "Ошибка поиска пользователя в БД", http.StatusBadRequest)
		return
	}
	if usr == nil {
		fmt.Printf("ошибка поиск пользователя: %v", err)
		http.Error(w, "Пользователь не найден", http.StatusBadRequest)
		return
	}

	//
	var pass models.UpdatePasswors

	body, err := io.ReadAll(r.Body)
	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Не верные данные 1 ")
		return
	}

	err = json.Unmarshal(body, &pass)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Не верные данные 2")
		return
	}

	err = service.ValidateUpdatePass(id, pass.OldPassword, pass.NewPassword)
	if err != nil {
		fmt.Printf("Ошибка обновления пароля: %v", err)
		http.Error(w, "ошибка обновления пароля", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Пароль успешно обновлен")
}
