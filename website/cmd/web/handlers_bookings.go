package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

type bookingTemplateData struct {
	Booking        websiteBooking
	Bookings       []websiteBooking
	BookingData    bookingData
	BookingsData   []bookingData
	Users          []websiteUser
	ShowTimes      []websiteShowTime
	Movies         []websiteMovie
	SelectedMovies []websiteMovie
}

type bookingData struct {
	ID           string
	UserFullName string
	ShowTimeDate string
}

func (app *application) loadBookingData(btd *bookingTemplateData, isList bool) {
	// Очищаем данные бронирований
	btd.BookingsData = []bookingData{}
	btd.BookingData = bookingData{}

	// Загружаем данные бронирований
	if isList {
		for _, b := range btd.Bookings {
			// Загружаем данные пользователя
			userURL := fmt.Sprintf("%s/%s", app.apis.users, b.UserID)
			var user websiteUser
			err := app.getAPIContent(userURL, &user)
			if err != nil {
				app.errorLog.Println(err.Error())
			}

			// Загружаем данные сеанса
			showtimeURL := fmt.Sprintf("%s/%s", app.apis.showtimes, b.ShowtimeID)
			var showtime websiteShowTime
			err = app.getAPIContent(showtimeURL, &showtime)
			if err != nil {
				app.errorLog.Println(err.Error())
			}

			bookingData := bookingData{
				ID:           b.ID.Hex(),
				UserFullName: fmt.Sprintf("%s %s", user.Name, user.LastName),
				ShowTimeDate: showtime.Date,
			}
			btd.BookingsData = append(btd.BookingsData, bookingData)
			app.infoLog.Println(b.UserID)
		}
	} else {
		b := btd.Booking

		// Загружаем данные пользователя
		userURL := fmt.Sprintf("%s/%s", app.apis.users, b.UserID)
		var user websiteUser
		err := app.getAPIContent(userURL, &user)
		if err != nil {
			app.errorLog.Println(err.Error())
		}

		// Загружаем данные сеанса
		showtimeURL := fmt.Sprintf("%s/%s", app.apis.showtimes, b.ShowtimeID)
		var showtime websiteShowTime

		err = app.getAPIContent(showtimeURL, &showtime)
		if err != nil {
			app.errorLog.Println(err.Error())
		}

		btd.BookingData = bookingData{
			ID:           b.ID.Hex(),
			UserFullName: fmt.Sprintf("%s %s", user.Name, user.LastName),
			ShowTimeDate: showtime.Date,
		}

		btd.SelectedMovies = []websiteMovie{}
		for _, movieID := range b.Movies {
			movieURL := fmt.Sprintf("%s/%s", app.apis.movies, movieID)
			var movie websiteMovie
			err = app.getAPIContent(movieURL, &movie)
			if err != nil {
				app.errorLog.Println(err.Error())
				continue
			}
			btd.SelectedMovies = append(btd.SelectedMovies, movie)
		}
	}
}

func (app *application) bookingsList(w http.ResponseWriter, r *http.Request) {

	// Получаем список бронирований из API
	var td bookingTemplateData
	app.infoLog.Println("Calling bookings API...")

	err := app.getAPIContent(app.apis.bookings, &td.Bookings)
	if err != nil {
		app.errorLog.Println(err.Error())
	}
	app.infoLog.Println(td.Bookings)
	app.infoLog.Println(td)

	app.loadBookingData(&td, true)
	app.loadBookingFormData(&td)

	// Загружаем файлы шаблонов
	files := []string{
		"./ui/html/bookings/list.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.New("list.page.tmpl").Funcs(movieTemplateFuncs()).ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, td)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func (app *application) loadBookingFormData(td *bookingTemplateData) {
	err := app.getAPIContent(app.apis.users, &td.Users)
	if err != nil {
		app.errorLog.Println(err.Error())
	}

	err = app.getAPIContent(app.apis.showtimes, &td.ShowTimes)
	if err != nil {
		app.errorLog.Println(err.Error())
	}

	err = app.getAPIContent(app.apis.movies, &td.Movies)
	if err != nil {
		app.errorLog.Println(err.Error())
	}
}

func (app *application) bookingsCreate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Redirect(w, r, "/bookings/list", http.StatusSeeOther)
		return
	}
	if len(r.PostForm["movies"]) == 0 {
		http.Redirect(w, r, "/bookings/list", http.StatusSeeOther)
		return
	}

	booking := websiteBooking{
		UserID:     r.PostForm.Get("userid"),
		ShowtimeID: r.PostForm.Get("showtimeid"),
		Movies:     r.PostForm["movies"],
	}

	err = app.postAPIContent(app.apis.bookings, booking)
	if err != nil {
		app.errorLog.Println(err.Error())
	}

	http.Redirect(w, r, "/bookings/list", http.StatusSeeOther)
}

func (app *application) bookingsDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookingID := vars["id"]

	err := app.deleteAPIContent(app.apiURL(app.apis.bookings, bookingID))
	if err != nil {
		app.errorLog.Println(err.Error())
	}

	http.Redirect(w, r, "/bookings/list", http.StatusSeeOther)
}

func (app *application) bookingsView(w http.ResponseWriter, r *http.Request) {
	// Получаем id из входящего URL
	vars := mux.Vars(r)
	bookingID := vars["id"]

	// Получаем список бронирований из API
	var td bookingTemplateData
	app.infoLog.Println("Calling bookings API...")
	url := fmt.Sprintf("%s/%s", app.apis.bookings, bookingID)

	err := app.getAPIContent(url, &td.Booking)
	if err != nil {
		app.errorLog.Println(err.Error())
	}
	app.infoLog.Println(td.Booking)
	app.infoLog.Println(url)

	app.loadBookingData(&td, false)

	// Загружаем файлы шаблонов
	files := []string{
		"./ui/html/bookings/view.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.New("view.page.tmpl").Funcs(movieTemplateFuncs()).ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, td)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}
