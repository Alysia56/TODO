//Filename: internal/data/models.go

package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFounc = errors.New("record not found")
)

//A wrapper for our data models
type Models struct {
	Todo TodoModel
}

//New Models() allows us to create a new Model
func NewModels(db *sql.DB) Models {
	return Models{
		Todo: TodoModel{DB: db},
	}
}
