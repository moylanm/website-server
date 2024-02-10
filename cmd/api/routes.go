package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"mylesmoylan.net/ui"
)

func (app *application) routes() http.Handler {
	router := mux.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	fileServer := http.FileServer(http.FS(ui.Files))
	router.PathPrefix("/static/").Handler(fileServer)

	excerptsPath := "/excerpts/{id:[0-9]+}"
	router.HandleFunc(excerptsPath, app.showExcerpt).Methods(http.MethodGet)
	router.HandleFunc(excerptsPath, app.authenticate(app.deleteExcerpt)).Methods(http.MethodDelete)
	router.HandleFunc(excerptsPath, app.authenticate(app.updateExcerpt)).Methods(http.MethodPatch)
	router.HandleFunc("/excerpts/json", app.listExcerptsJson).Methods(http.MethodGet)
	router.HandleFunc("/excerpts", app.authenticate(app.createExcerpt)).Methods(http.MethodPost)
	router.HandleFunc("/excerpts", app.listExcerpts).Methods(http.MethodGet)

	router.HandleFunc("/about", app.about).Methods(http.MethodGet)
	router.HandleFunc("/", app.home).Methods(http.MethodGet)

	standard := alice.New(app.recoverPanic, app.rateLimit, secureHeaders)

	return standard.Then(router)
}
