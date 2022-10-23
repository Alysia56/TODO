// Filename/cmd/api/todo.go

package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"alysianorales.net/TODO/internal/data"
	"alysianorales.net/TODO/internal/validator"
)

//createTodoHandler for the "POST /v1/todo" endpoint
func (app *application) createTodoHandler(w http.ResponseWriter, r *http.Request) {
	// Our Target Decode destination
	var input struct {
		Name    string   `json:"name"`
		Level   string   `json:"level"`
		Contact string   `json:"contact"`
		Phone   string   `json:"phone"`
		Email   string   `json:"email"`
		Website string   `json:"website"`
		Address string   `json:"address"`
		Mode    []string `json:"mode"`
	}
	// Initialize a new json.
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	//Copy the values from the input struct to a new Todo Struct
	todo := &data.Todo{
		ID:        0,
		CreatedAt: time.Time{},
		Name:      input.Name,
		Level:     input.Level,
		Contact:   input.Contact,
		Phone:     input.Phone,
		Email:     input.Email,
		Website:   input.Website,
		Address:   input.Address,
		Mode:      input.Mode,
		Version:   0,
	}

	// Initialize a new Validator instance
	v := validator.New()

	// check the map to see if there were validation errors
	if data.ValidateList(v, todo); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Create a Todo
	err = app.models.Todo.Insert(todo)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	// Create a Location header for the newly created resource/List
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/todo/%d", todo.ID))
	//Write the JSON response with 201 - Created status code with the body
	//being the List data and the header being the headers map
	err = app.writeJSON(w, http.StatusCreated, envelope{"todo": todo}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

//showTodoHandler for the "GET /v1/todo/:id" endpoint
func (app *application) showTodoHandler(w http.ResponseWriter, r *http.Request) {
	// Get the value of the "id" parameter
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Fetch the specific List
	todo, err := app.models.Todo.Get(id)
	//Handle errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	//Write the data returned by Get()
	err = app.writeJSON(w, http.StatusOK, envelope{"todo": todo}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showRandomString(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	integer := int(id)
	tools := &data.Tools{}

	random := tools.GenerateRandomString(integer)
	data := envelope{
		"Here is your randomize string": random,
		"Your :id is ":                  integer,
	}
	err = app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	//This method does a complete replacement
	//Get the id for the List that needs updating
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	//Fetch the original record from the database
	todo, err := app.models.Todo.Get(id)
	//Handle errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Create an input struct to hold data read in from the client
	var input struct {
		Name    string   `json:"name"`
		Level   string   `json:"level"`
		Contact string   `json:"contact"`
		Phone   string   `json:"phone"`
		Email   string   `json:"email"`
		Website string   `json:"website"`
		Address string   `json:"address"`
		Mode    []string `json:"mode"`
	}
	// Initialize a new json.
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	//Copy / Update the fields / values in the todo variable using the fields in the input struct
	todo.Name = input.Name
	todo.Level = input.Level
	todo.Contact = input.Contact
	todo.Phone = input.Phone
	todo.Email = input.Email
	todo.Website = input.Website
	todo.Address = input.Address
	todo.Mode = input.Mode
	//Perform validation on the updated Todo.
	//If validation fails, then send a 422 - Unprocessable Entity response to the client
	// Initialize a new Validator instance
	v := validator.New()

	// check the map to see if there were validation errors
	if data.ValidateList(v, todo); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Pass the updates Todo record to the Update() method
	err = app.models.Todo.Update(todo)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	//Write the data returned by Get()
	err = app.writeJSON(w, http.StatusOK, envelope{"todo": todo}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	//Delete the School from the database. Send a 404 Not Found status code
	//to the client if there is no matching record
	err = app.models.Todo.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Return 200 Status OK to the client with a success message
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "List successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
