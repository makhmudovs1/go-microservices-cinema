package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
)

type movieTemplateData struct {
	Movie  websiteMovie
	Movies []websiteMovie
}

func posterURL(title string) string {
	posters := map[string]string{
		"avatar: the way of water": "/static/img/posters/avatar-way-of-water.jpg",
		"аватар: путь воды":        "/static/img/posters/avatar-way-of-water.jpg",
		"аватар: огонь и вода":     "/static/img/posters/avatar-way-of-water.jpg",
		"человек-паук":             "/static/img/posters/spider-man-new-day.webp",
		"человек паук":             "/static/img/posters/spider-man-new-day.webp",
		"spider-man":               "/static/img/posters/spider-man-new-day.webp",
		"interstellar":             "/static/img/posters/interstellar.jpg",
		"dune: part two":           "/static/img/posters/dune-part-two.jpg",
		"дюна: часть вторая":       "/static/img/posters/dune-part-two.jpg",
		"obsession":                "/static/img/posters/obsession.webp",
		"мстители: думсдей":        "/static/img/posters/avengers-doomsday.svg",
		"avengers: doomsday":       "/static/img/posters/avengers-doomsday.svg",
		"the lighthouse":           "/static/img/posters/the-lighthouse.jpg",
		"a hidden life":            "/static/img/posters/a-hidden-life.jpg",
		"soul":                     "/static/img/posters/soul.jpg",
		"little joe":               "/static/img/posters/little-joe.jpg",
		"onward":                   "/static/img/posters/onward.jpg",
	}

	if poster, ok := posters[strings.ToLower(strings.TrimSpace(title))]; ok {
		return poster
	}
	return "/static/img/posters/cinema-placeholder.svg"
}

func movieTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"posterURL": posterURL,
	}
}

func (app *application) moviesList(w http.ResponseWriter, r *http.Request) {

	// Получаем список фильмов из API
	var mtd movieTemplateData
	app.infoLog.Println("Calling movies API...")
	app.getAPIContent(app.apis.movies, &mtd.Movies)
	app.infoLog.Println(mtd.Movies)

	// Загружаем файлы шаблонов
	files := []string{
		"./ui/html/movies/list.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.New("list.page.tmpl").Funcs(movieTemplateFuncs()).ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, mtd)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func (app *application) moviesCreate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Redirect(w, r, "/movies/list", http.StatusSeeOther)
		return
	}

	rating, err := strconv.ParseFloat(r.PostForm.Get("rating"), 32)
	if err != nil {
		rating = 0
	}

	movie := websiteMovie{
		Title:    r.PostForm.Get("title"),
		Director: r.PostForm.Get("director"),
		Rating:   float32(rating),
	}

	err = app.postAPIContent(app.apis.movies, movie)
	if err != nil {
		app.errorLog.Println(err.Error())
	}

	http.Redirect(w, r, "/movies/list", http.StatusSeeOther)
}

func (app *application) moviesDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	movieID := vars["id"]

	err := app.deleteAPIContent(app.apiURL(app.apis.movies, movieID))
	if err != nil {
		app.errorLog.Println(err.Error())
	}

	http.Redirect(w, r, "/movies/list", http.StatusSeeOther)
}

func (app *application) moviesView(w http.ResponseWriter, r *http.Request) {
	// Получаем id из входящего URL
	vars := mux.Vars(r)
	movieID := vars["id"]

	// Получаем список фильмов из API
	app.infoLog.Println("Calling movies API...")
	url := fmt.Sprintf("%s/%s", app.apis.movies, movieID)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	var td movieTemplateData
	json.Unmarshal(bodyBytes, &td.Movie)
	app.infoLog.Println(td.Movie)

	// Загружаем файлы шаблонов
	files := []string{
		"./ui/html/movies/view.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.New("view.page.tmpl").Funcs(movieTemplateFuncs()).ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, td.Movie)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}
