// Filename: cmd/api/routes
package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	// Create a new HttpRouter router instance
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/todo", app.requirePermission("todo:read", app.listTodoHandler))
	router.HandlerFunc(http.MethodPost, "/v1/todo", app.requirePermission("todo:write", app.createTodoHandler))
	router.HandlerFunc(http.MethodGet, "/v1/todo/:id", app.requirePermission("todo:read", app.showTodoHandler))
	//router.HandlerFunc(http.MethodGet, "/v1/stringrandom/:id", app.showRandomString)
	router.HandlerFunc(http.MethodPatch, "/v1/todo/:id", app.requirePermission("todo:write", app.updateTodoHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/todo/:id", app.requirePermission("todo:write", app.deleteTodoHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	return app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router))))

}
