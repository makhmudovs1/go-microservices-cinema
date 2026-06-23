package main

import (
	"encoding/json"
	"net/http"

	"cinema/bookings/pkg/models"
	"github.com/gorilla/mux"
)

func (app *application) all(w http.ResponseWriter, r *http.Request) {
	// Получаем все сохраненные бронирования
	bookings, err := app.bookings.All()
	if err != nil {
		app.serverError(w, err)
	}

	// Преобразуем список бронирований в JSON
	b, err := json.Marshal(bookings)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Println("Bookings have been listed")

	// Отправляем ответ клиенту
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) findByID(w http.ResponseWriter, r *http.Request) {
	// Получаем id из входящего URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Ищем бронирование по id
	m, err := app.bookings.FindByID(id)
	if err != nil {
		if err.Error() == "ErrNoDocuments" {
			app.infoLog.Println("Booking not found")
			return
		}
		// Любая другая ошибка вернет внутреннюю ошибку сервера
		app.serverError(w, err)
	}

	// Преобразуем бронирование в JSON
	b, err := json.Marshal(m)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Println("Have been found a booking")

	// Отправляем ответ клиенту
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) insert(w http.ResponseWriter, r *http.Request) {
	// Определяем модель бронирования
	var m models.Booking
	// Получаем данные запроса
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		app.serverError(w, err)
	}

	// Добавляем новое бронирование
	insertResult, err := app.bookings.Insert(m)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Printf("New booking have been created, id=%s", insertResult.InsertedID)
}

func (app *application) delete(w http.ResponseWriter, r *http.Request) {
	// Получаем id из входящего URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Удаляем бронирование по id
	deleteResult, err := app.bookings.Delete(id)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Printf("Have been eliminated %d booking(s)", deleteResult.DeletedCount)
}
