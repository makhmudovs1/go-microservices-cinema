package main

import (
	"encoding/json"
	"net/http"

	"cinema/movies/pkg/models"
	"github.com/gorilla/mux"
)

func (app *application) all(w http.ResponseWriter, r *http.Request) {
	// Получаем все сохраненные фильмы
	movies, err := app.movies.All()
	if err != nil {
		app.serverError(w, err)
	}

	// Преобразуем список фильмов в JSON
	b, err := json.Marshal(movies)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Println("Movies have been listed")

	// Отправляем ответ клиенту
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) findByID(w http.ResponseWriter, r *http.Request) {
	// Получаем id из входящего URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Ищем фильм по id
	m, err := app.movies.FindByID(id)
	if err != nil {
		if err.Error() == "ErrNoDocuments" {
			app.infoLog.Println("Movie not found")
			return
		}
		// Любая другая ошибка вернет внутреннюю ошибку сервера
		app.serverError(w, err)
	}

	// Преобразуем фильм в JSON
	b, err := json.Marshal(m)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Println("Have been found a movie")

	// Отправляем ответ клиенту
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) insert(w http.ResponseWriter, r *http.Request) {
	// Определяем модель фильма
	var m models.Movie
	// Получаем данные запроса
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		app.serverError(w, err)
	}

	// Добавляем новый фильм
	insertResult, err := app.movies.Insert(m)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Printf("New movie have been created, id=%s", insertResult.InsertedID)
}

func (app *application) delete(w http.ResponseWriter, r *http.Request) {
	// Получаем id из входящего URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Удаляем фильм по id
	deleteResult, err := app.movies.Delete(id)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Printf("Have been eliminated %d movie(s)", deleteResult.DeletedCount)
}
