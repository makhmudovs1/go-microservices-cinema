package main

import (
	"github.com/gorilla/mux"
)

func (app *application) routes() *mux.Router {
	// Регистрируем функции-обработчики.
	r := mux.NewRouter()
	r.HandleFunc("/", app.home)
	r.HandleFunc("/register", app.register).Methods("GET")
	r.HandleFunc("/register", app.registerCreate).Methods("POST")
	r.HandleFunc("/users/list", app.usersList)
	r.HandleFunc("/users/view/{id}", app.usersView)
	r.HandleFunc("/users/create", app.usersCreate).Methods("POST")
	r.HandleFunc("/users/delete/{id}", app.usersDelete).Methods("POST")
	r.HandleFunc("/movies/list", app.moviesList)
	r.HandleFunc("/movies/view/{id}", app.moviesView)
	r.HandleFunc("/movies/create", app.moviesCreate).Methods("POST")
	r.HandleFunc("/movies/delete/{id}", app.moviesDelete).Methods("POST")
	r.HandleFunc("/showtimes/list", app.showtimesList)
	r.HandleFunc("/showtimes/view/{id}", app.showtimesView)
	r.HandleFunc("/showtimes/create", app.showtimesCreate).Methods("POST")
	r.HandleFunc("/showtimes/delete/{id}", app.showtimesDelete).Methods("POST")
	r.HandleFunc("/bookings/list", app.bookingsList)
	r.HandleFunc("/bookings/view/{id}", app.bookingsView)
	r.HandleFunc("/bookings/create", app.bookingsCreate).Methods("POST")
	r.HandleFunc("/bookings/delete/{id}", app.bookingsDelete).Methods("POST")

	// Это будет отдавать файлы по адресу http://localhost:8000/static/<filename>
	r.PathPrefix("/static/").Handler(app.static("./ui/static/"))
	return r
}
