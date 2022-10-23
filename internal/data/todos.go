// Filname: internal/data/todos.go

package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"alysianorales.net/TODO/internal/validator"
	"github.com/lib/pq"
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
	v.Check(len(todo.Mode) >= 1, "mode", "must contain at least one todos")
	v.Check(len(todo.Mode) <= 5, "mode", "must contain at most 5 todos")
	v.Check(validator.Unique(todo.Mode), "mode", "must not contain duplicate todos")
}

// Define a ListModel which wraps a sql.DB connection pool
type TodoModel struct {
	DB *sql.DB
}

// Insert () allows us to create a new List
func (m TodoModel) Insert(todo *Todo) error {
	query := `
	INSERT INTO todo (name, level, contact, phone, email, website, address, mode)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id, created_at, version
	`
	//Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//Cleanup to prevent memory leaks
	defer cancel()
	//Collect data fields into a slice
	args := []interface{}{
		todo.Name, todo.Level, todo.Contact,
		todo.Phone, todo.Email, todo.Website,
		todo.Address, pq.Array(todo.Mode),
	}
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&todo.ID, &todo.CreatedAt, &todo.Version)
}

//Get() allows us to retrieve a specific List
func (m TodoModel) Get(id int64) (*Todo, error) {
	//Ensure that there is a valid id
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	//Create the query
	query := `
		SELECT id, created_at, name, level, contact, phone, email, website, address, mode, version
		FROM todo
		WHERE id = $1
	`
	//Declare a Todo variable to hold the returned data
	var todo Todo
	//Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//Cleanup to prevent memory leaks
	defer cancel()
	//Exexcute the query using QueryRow()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&todo.ID,
		&todo.CreatedAt,
		&todo.Name,
		&todo.Level,
		&todo.Contact,
		&todo.Phone,
		&todo.Email,
		&todo.Website,
		&todo.Address,
		pq.Array(&todo.Mode),
		&todo.Version,
	)
	//Handle any errors
	if err != nil {
		//Check the type of error
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	//Success
	return &todo, nil
}

//Update() allows us to edit/alter a specific List
//Optimistic locking (version number)
func (m TodoModel) Update(todo *Todo) error {
	//Create a query
	query := `
		UPDATE todo
		SET name = $1, level = $2, contact = $3, 
			phone = $4, email = $5, website = $6, 
			address = $7, mode = $8, version = version + 1
		WHERE id = $9
		AND version = $10		
		RETURNING version
	`
	//Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//Cleanup to prevent memory leaks
	defer cancel()

	args := []interface{}{
		todo.Name,
		todo.Level,
		todo.Contact,
		todo.Phone,
		todo.Email,
		todo.Website,
		todo.Address,
		pq.Array(todo.Mode),
		todo.ID,
		todo.Version,
	}
	//Check for edit conflicts
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&todo.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

//Delete() removes a specific List
func (m TodoModel) Delete(id int64) error {
	//Ensure that there is a valid id
	if id < 1 {
		return ErrRecordNotFound
	}
	//Create the delete query
	query := `
		DELETE FROM todo
		WHERE id = $1
	`
	//Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//Cleanup to prevent memory leaks
	defer cancel()
	//Execute the query.
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	//Check how many rows were affected by the delete operation.
	//Call the RowsAffected() method on the result variable
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	//Check if no rows were affected
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil

}

//The GetAll() method returns a list of all the todos sorted by id
func (m TodoModel) GetAll(name string, level string, mode []string, filters Filters) ([]*Todo, Metadata, error) {
	//Construct the query
	query := fmt.Sprintf(`
		SELECT COUNT (*) OVER(), id, created_at, name, 
				level, contact,phone, 
				email, website, address, mode, version
		FROM todo
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (to_tsvector('simple', level) @@ plainto_tsquery('simple', $2) OR $2 = '')
		AND (mode @> $3 OR $3 = '{}' )
		ORDER BY %s %s, id ASC
		LIMIT $4 OFFSET $5`, filters.sortColumn(), filters.sortOrder())
	//Create a 3-second-timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	//Execute the query
	args := []interface{}{name, level, pq.Array(mode), filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	//Close the result set
	defer rows.Close()
	totalRecords := 0
	//Initialize an empty slice to hold the Todo data
	todos := []*Todo{}
	//Iterate over the rows in the resultset
	for rows.Next() {
		var todo Todo
		//Scan the values from the row into Todo
		err := rows.Scan(
			&totalRecords,
			&todo.ID,
			&todo.CreatedAt,
			&todo.Name,
			&todo.Level,
			&todo.Contact,
			&todo.Phone,
			&todo.Email,
			&todo.Website,
			&todo.Address,
			pq.Array(&todo.Mode),
			&todo.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		//Add the List to our slice
		todos = append(todos, &todo)
	}
	// Check for errors after looping through the results set
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	//Return the slice of Lists
	return todos, metadata, nil
}
