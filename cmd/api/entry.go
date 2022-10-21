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
	// Initialize a new Validator instance
	v := validator.New()

	// Check() method to execute
	v.Check(input.Name != "", "name", "must be provided")
	v.Check(len(input.Name) <= 200, "name", "must not be more than 200 bytes long")

	v.Check(input.Level != "", "level", "must be provided")
	v.Check(len(input.Level) <= 200, "level", "must not be more than 200 bytes long")

	v.Check(input.Contact != "", "contact", "must be provided")
	v.Check(len(input.Contact) <= 200, "contact", "must not be more than 200 bytes long")

	v.Check(input.Phone != "", "phone", "must be provided")
	v.Check(validator.Matches(input.Phone, validator.PhoneRX), "phone", "must be a valid phone number")

	v.Check(input.Email != "", "email", "must be provided")
	v.Check(validator.Matches(input.Email, validator.EmailRX), "email", "must be a valid email")

	v.Check(input.Website != "", "website", "must be provided")
	v.Check(validator.ValidWebsite(input.Website), "website", "must be a valid url")

	v.Check(input.Address != "", "address", "must be provided")
	v.Check(len(input.Address) <= 500, "address", "must not be more than 500 bytes long")

	v.Check(input.Mode != nil, "mode", "must be provided")
	v.Check(len(input.Mode) >= 1, "mode", "must contain at least one entries")
	v.Check(len(input.Mode) <= 5, "mode", "must contain at most 5 entries")
	v.Check(validator.Unique(input.Mode), "mode", "must not contain duplicate entries")
	// check the map to see if there were validation errors
	if !v.Valid() {
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
	entry := data.Entry{
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
