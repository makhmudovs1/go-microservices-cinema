package main

import (
	"encoding/json"
	"net/http"

	"cinema/users/pkg/models"
	"github.com/gorilla/mux"
)

func (app *application) all(w http.ResponseWriter, r *http.Request) {
	// Получаем всех сохраненных пользователей
	users, err := app.users.All()
	if err != nil {
		app.serverError(w, err)
	}

	// Преобразуем список пользователей в JSON
	b, err := json.Marshal(users)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Println("Users have been listed")

	// Отправляем ответ клиенту
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) findByID(w http.ResponseWriter, r *http.Request) {
	// Получаем id из входящего URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Ищем пользователя по id
	m, err := app.users.FindByID(id)
	if err != nil {
		if err.Error() == "ErrNoDocuments" {
			app.infoLog.Println("User not found")
			return
		}
		// Любая другая ошибка вернет внутреннюю ошибку сервера
		app.serverError(w, err)
	}

	// Преобразуем пользователя в JSON
	b, err := json.Marshal(m)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Println("Have been found a user")

	// Отправляем ответ клиенту
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) insert(w http.ResponseWriter, r *http.Request) {
	// Определяем модель пользователя
	var u models.User
	// Получаем данные запроса
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		app.serverError(w, err)
	}

	// Добавляем нового пользователя
	insertResult, err := app.users.Insert(u)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Printf("New user have been created, id=%s", insertResult.InsertedID)
}

func (app *application) delete(w http.ResponseWriter, r *http.Request) {
	// Получаем id из входящего URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Удаляем пользователя по id
	deleteResult, err := app.users.Delete(id)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Printf("Have been eliminated %d user(s)", deleteResult.DeletedCount)
}
