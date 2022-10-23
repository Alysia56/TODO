// Filename: cmd/api/routes
package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	// Create a new HttpRouter router instance
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/todo", app.createTodoHandler)
	router.HandlerFunc(http.MethodGet, "/v1/todo/:id", app.showTodoHandler)
	router.HandlerFunc(http.MethodGet, "/v1/stringrandom/:id", app.showRandomString)

	return router
}
