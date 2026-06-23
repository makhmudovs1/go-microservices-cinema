package main

import (
	"encoding/json"
	"net/http"
	"time"

	"cinema/showtimes/pkg/models"
	"github.com/gorilla/mux"
)

func (app *application) all(w http.ResponseWriter, r *http.Request) {
	// Получаем все сохраненные сеансы
	showtimes, err := app.showtimes.All()
	if err != nil {
		app.serverError(w, err)
	}

	// Преобразуем список сеансов в JSON
	b, err := json.Marshal(showtimes)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Println("Showtimes have been listed")

	// Отправляем ответ клиенту
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) findByID(w http.ResponseWriter, r *http.Request) {
	// Получаем id из входящего URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Ищем сеанс по id
	m, err := app.showtimes.FindByID(id)
	if err != nil {
		if err.Error() == "ErrNoDocuments" {
			app.infoLog.Println("Showtime not found")
			return
		}
		// Любая другая ошибка вернет внутреннюю ошибку сервера
		app.serverError(w, err)
	}

	// Преобразуем сеанс в JSON
	b, err := json.Marshal(m)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Println("Have been found a showtime")

	// Отправляем ответ клиенту
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) findByDate(w http.ResponseWriter, r *http.Request) {
	// Получаем id из входящего URL
	vars := mux.Vars(r)
	date := vars["date"]

	// Ищем сеанс по дате
	m, err := app.showtimes.FindByDate(date)
	if err != nil {
		if err.Error() == "ErrNoDocuments" {
			app.infoLog.Println("Showtime not found")
			return
		}
		// Любая другая ошибка вернет внутреннюю ошибку сервера
		app.serverError(w, err)
	}

	// Преобразуем сеанс в JSON
	b, err := json.Marshal(m)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Println("Have been found a showtime")

	// Отправляем ответ клиенту
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) insert(w http.ResponseWriter, r *http.Request) {
	// Определяем модель сеанса
	var m models.ShowTime
	// Получаем данные запроса
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		app.serverError(w, err)
	}

	// Добавляем новый сеанс
	m.CreatedAt = time.Now()
	insertResult, err := app.showtimes.Insert(m)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Printf("New showtime have been created, id=%s", insertResult.InsertedID)
}

func (app *application) delete(w http.ResponseWriter, r *http.Request) {
	// Получаем id из входящего URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Удаляем сеанс по id
	deleteResult, err := app.showtimes.Delete(id)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Printf("Have been eliminated %d showtime(s)", deleteResult.DeletedCount)
}
