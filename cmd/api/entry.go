// Filename/cmd/api/entry.go

package main

import (
	"fmt"
	"net/http"
	"time"

	"alysianorales.net/TODO/internal/data"
	"alysianorales.net/TODO/internal/validator"
)

//createEntryHandler for the "POST /v1/entry" endpoint
func (app *application) createEntryHandler(w http.ResponseWriter, r *http.Request) {
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
	// Display the Request
	fmt.Fprintf(w, "%+v\n", input)
}

//createEntryHandler for the "GET /v1/entry/:id" endpoint
func (app *application) showEntryHandler(w http.ResponseWriter, r *http.Request) {
	// Get the value of the "id" parameter
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Create a new instance of the Entries struct containing the ID we extracted
	//from our URL and some sample data
	entry := data.Todo{
		ID:        id,
		CreatedAt: time.Now(),
		Name:      "Yo Mama",
		Level:     "High School",
		Contact:   "Inita Lyfe",
		Phone:     "666-7777",
		Address:   "14 Upyoaph Street",
		Mode:      []string{"blended", "online"},
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"entry": entry}, nil)
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
