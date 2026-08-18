package handlers

import (
	"HOTA/internal/models"
	"HOTA/internal/service"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Auth(w http.ResponseWriter, r *http.Request) {
	var authDTO models.AuotIn

	body, err := io.ReadAll(r.Body)
	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Не верные данные 1 ")
		return
	}

	err = json.Unmarshal(body, &authDTO)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Не верные данные 2")
		return
	}

	//Получаем текущий email и хеш пароля из БД
	var Email string
	var Hashpass string
	var id int

	query := "SELECT id, Email, Password FROM users WHERE Email = ?"
	err = models.UserDB.QueryRow(query, authDTO.Email).Scan(&id, &Email, &Hashpass)
	if err != nil {
		http.Error(w, "Данный email не зарегистрирован", http.StatusBadRequest)
	}

	//сравниванием пароли
	err = service.CheckPassword(Hashpass, authDTO.Password)
	if err != nil {
		w.WriteHeader(401)
		fmt.Fprintf(w, "Пароли не совпадают")
		return
	}

	jwtDTO, err := service.CreatJWT(id)
	if err != nil {
		w.WriteHeader(401)
		fmt.Printf("Ошибка создания токена %v", err)
		http.Error(w, "Ошибка создания токена", 400)
		return
	}

	fmt.Println(jwtDTO)
	w.Header().Set("Context-Type", "application/json")
	json.NewEncoder(w).Encode(jwtDTO)

}
