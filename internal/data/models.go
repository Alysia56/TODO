//Filename: internal/data/models.go

package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

//A wrapper for our data models
type Models struct {
	Todo   TodoModel
	Tokens TokenModel
	Users  UserModel
}

//New Models() allows us to create a new Model
func NewModels(db *sql.DB) Models {
	return Models{
		Todo:   TodoModel{DB: db},
		Tokens: TokenModel{DB: db},
		Users:  UserModel{DB: db},
	}
}
