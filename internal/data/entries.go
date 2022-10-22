// Filname: internal/data/entries.go

package data

import (
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
