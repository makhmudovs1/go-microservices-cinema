package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type websiteUser struct {
	ID       primitive.ObjectID
	Name     string
	LastName string
	Email    string
}

type userPayload struct {
	Name     string
	LastName string
	Email    string
}

type userTemplateData struct {
	User  websiteUser
	Users []websiteUser
}

func (app *application) usersList(w http.ResponseWriter, r *http.Request) {

	// Получаем список пользователей из API
	var utd userTemplateData
	err := app.getAPIContent(app.apis.users, &utd.Users)
	if err != nil {
		app.errorLog.Println(err.Error())
	}
	app.infoLog.Println(utd.Users)

	// Загружаем файлы шаблонов
	files := []string{
		"./ui/html/users/list.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, utd)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func (app *application) usersView(w http.ResponseWriter, r *http.Request) {
	// Получаем id из входящего URL
	vars := mux.Vars(r)
	userID := vars["id"]

	// Получаем список пользователей из API
	app.infoLog.Println("Calling users API...")
	url := fmt.Sprintf("%s/%s", app.apis.users, userID)

	var utd userTemplateData
	app.getAPIContent(url, &utd.User)
	app.infoLog.Println(utd.User)

	// Загружаем файлы шаблонов
	files := []string{
		"./ui/html/users/view.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, utd.User)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func (app *application) usersCreate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Redirect(w, r, "/users/list", http.StatusSeeOther)
		return
	}

	user := userPayload{
		Name:     r.PostForm.Get("name"),
		LastName: r.PostForm.Get("lastname"),
		Email:    r.PostForm.Get("email"),
	}

	err = app.postAPIContent(app.apis.users, user)
	if err != nil {
		app.errorLog.Println(err.Error())
	}

	http.Redirect(w, r, "/users/list", http.StatusSeeOther)
}

func (app *application) register(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/users/register.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func (app *application) registerCreate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	user := userPayload{
		Name:     r.PostForm.Get("name"),
		LastName: r.PostForm.Get("lastname"),
		Email:    r.PostForm.Get("email"),
	}

	err = app.postAPIContent(app.apis.users, user)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/bookings/list", http.StatusSeeOther)
}

func (app *application) usersDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	err := app.deleteAPIContent(app.apiURL(app.apis.users, userID))
	if err != nil {
		app.errorLog.Println(err.Error())
	}

	http.Redirect(w, r, "/users/list", http.StatusSeeOther)
}
