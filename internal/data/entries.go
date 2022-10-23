// Filname: internal/data/entries.go

package data

import (
	"database/sql"
	"time"

	"alysianorales.net/TODO/internal/validator"
)

type Todo struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Name      string    `json:"name"`
	Level     string    `json:"level"`
	Contact   string    `json:"contact"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email,omitempty"`
	Website   string    `json:"website,omitempty"`
	Address   string    `json:"address"`
	Mode      []string  `json:"mode"`
	Version   int32     `json:"version"`
}

func ValidateList(v *validator.Validator, todo *Todo) {
	// Check() method to execute
	v.Check(todo.Name != "", "name", "must be provided")
	v.Check(len(todo.Name) <= 200, "name", "must not be more than 200 bytes long")

	v.Check(todo.Level != "", "level", "must be provided")
	v.Check(len(todo.Level) <= 200, "level", "must not be more than 200 bytes long")

	v.Check(todo.Contact != "", "contact", "must be provided")
	v.Check(len(todo.Contact) <= 200, "contact", "must not be more than 200 bytes long")

	v.Check(todo.Phone != "", "phone", "must be provided")
	v.Check(validator.Matches(todo.Phone, validator.PhoneRX), "phone", "must be a valid phone number")

	v.Check(todo.Email != "", "email", "must be provided")
	v.Check(validator.Matches(todo.Email, validator.EmailRX), "email", "must be a valid email")

	v.Check(todo.Website != "", "website", "must be provided")
	v.Check(validator.ValidWebsite(todo.Website), "website", "must be a valid url")

	v.Check(todo.Address != "", "address", "must be provided")
	v.Check(len(todo.Address) <= 500, "address", "must not be more than 500 bytes long")

	v.Check(todo.Mode != nil, "mode", "must be provided")
	v.Check(len(todo.Mode) >= 1, "mode", "must contain at least one entries")
	v.Check(len(todo.Mode) <= 5, "mode", "must contain at most 5 entries")
	v.Check(validator.Unique(todo.Mode), "mode", "must not contain duplicate entries")
}

// Define a ListModel which wraps a sql.DB connection pool
type TodoModel struct {
	DB *sql.DB
}

// Insert () allows us to create a new List
func (m TodoModel) Insert(todo *Todo) error {
	return nil
}

//Get() allows us to retrieve a specific List
func (m TodoModel) Get(id int64) (*Todo, error) {
	return nil, nil
}

//Update() allows us to edit/alter a specific List
func (m TodoModel) Update(todo *Todo) error {
	return nil
}

//Delete() removes a specific List
func (m TodoModel) Delete(id int64) error {
	return nil
}
