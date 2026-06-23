package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
)

type showtimeTemplateData struct {
	ShowTime        websiteShowTime
	ShowTimes       []websiteShowTime
	Movies          string
	AvailableMovies []websiteMovie
}

func (app *application) showtimesList(w http.ResponseWriter, r *http.Request) {

	// Получаем список сеансов из API
	app.infoLog.Println("Calling showtimes API...")
	var td showtimeTemplateData
	err := app.getAPIContent(app.apis.showtimes, &td.ShowTimes)
	if err != nil {
		app.errorLog.Println(err.Error())
	}
	err = app.getAPIContent(app.apis.movies, &td.AvailableMovies)
	if err != nil {
		app.errorLog.Println(err.Error())
	}
	app.infoLog.Println(td.ShowTimes)

	// Загружаем файлы шаблонов
	files := []string{
		"./ui/html/showtimes/list.page.tmpl",
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

func (app *application) showtimesCreate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Redirect(w, r, "/showtimes/list", http.StatusSeeOther)
		return
	}
	if len(r.PostForm["movies"]) == 0 {
		http.Redirect(w, r, "/showtimes/list", http.StatusSeeOther)
		return
	}

	showtime := websiteShowTime{
		Date:   r.PostForm.Get("date"),
		Movies: r.PostForm["movies"],
	}

	err = app.postAPIContent(app.apis.showtimes, showtime)
	if err != nil {
		app.errorLog.Println(err.Error())
	}

	http.Redirect(w, r, "/showtimes/list", http.StatusSeeOther)
}

func (app *application) showtimesDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	showtimeID := vars["id"]

	err := app.deleteAPIContent(app.apiURL(app.apis.showtimes, showtimeID))
	if err != nil {
		app.errorLog.Println(err.Error())
	}

	http.Redirect(w, r, "/showtimes/list", http.StatusSeeOther)
}

func (app *application) showtimesView(w http.ResponseWriter, r *http.Request) {
	// Получаем id из входящего URL
	vars := mux.Vars(r)
	showtimeID := vars["id"]

	// Получаем список сеансов из API
	app.infoLog.Println("Calling showtimes API...")
	url := fmt.Sprintf("%s/%s", app.apis.showtimes, showtimeID)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	var td showtimeTemplateData
	json.Unmarshal(bodyBytes, &td.ShowTime)
	app.infoLog.Println(td.ShowTime)

	// Загружаем названия фильмов
	var movies []string
	for _, m := range td.ShowTime.Movies {
		url := fmt.Sprintf("%s/%s", app.apis.movies, m)

		resp, err := http.Get(url)
		if err != nil {
			fmt.Print(err.Error())
		}
		defer resp.Body.Close()

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Print(err.Error())
		}

		var movie websiteMovie
		json.Unmarshal(bodyBytes, &movie)
		movies = append(movies, movie.Title)
	}
	td.Movies = strings.Join(movies, ", ")

	// Загружаем файлы шаблонов
	files := []string{
		"./ui/html/showtimes/view.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
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
