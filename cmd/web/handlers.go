package main

import (
	"errors"
	"fmt"
	"net/http"
	"snippetbox.kirill.ru/internal/models"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	snippets, err := app.snippets.Latest()

	if err != nil {
		app.serverError(w, err)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(w, "%+v\n\n", snippet)
	}

	//files := []string{
	//	"./ui/html/pages/base.tmpl",
	//	"./ui/html/pages/nav.tmpl",
	//	"./ui/html/pages/home.tmpl",
	//}
	//
	//ts, err := template.ParseFiles(files...)
	//if err != nil {
	//	app.serverError(w, err)
	//	return
	//}
	//err = ts.ExecuteTemplate(w, "base", nil)
	//
	//if err != nil {
	//	app.serverError(w, err)
	//	return
	//}
}

// Add a snippetView handler function.
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	snippet, err := app.snippets.Get(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRows) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	//w.Write([]byte("Display a specific snippet..."))
	fmt.Fprintf(w, "%+v", snippet)
}

// Add a snippetCreate handler function.
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		//w.WriteHeader(405)
		//w.Write([]byte("Method Not Allowed"))
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	title := "Testing title"
	content := "Lorem ipsum \ndolor sit amet, consectetur adipisicing elit. \n\nBeatae cumque cupiditate doloribus"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)

	w.Write([]byte("Create a new snippet..."))
}
