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
	Title     string    `json:"title"`
	Label     string    `json:"label"`
	Task      string    `json:"task"`
	Priority  string    `json:"priority"`
	Status    string    `json:"status,omitempty"`
	Website   string    `json:"website,omitempty"`
	Address   string    `json:"address"`
	Mode      []string  `json:"mode"`
	Version   int32     `json:"version"`
}

func ValidateList(v *validator.Validator, todo *Todo) {
	// Check() method to execute
	v.Check(todo.Title != "", "title", "must be provided")
	v.Check(len(todo.Title) <= 200, "title", "must not be more than 200 bytes long")

	v.Check(todo.Label != "", "label", "must be provided")
	v.Check(len(todo.Label) <= 200, "label", "must not be more than 200 bytes long")

	v.Check(todo.Task != "", "task", "must be provided")
	v.Check(len(todo.Task) <= 200, "task", "must not be more than 200 bytes long")

	v.Check(todo.Priority != "", "priority", "must be provided")
	//v.Check(validator.Matches(todo.Priority, validator.PriorityRX), "priority", "must be a valid priority number")

	v.Check(todo.Status != "", "status", "must be provided")
	//v.Check(validator.Matches(todo.Status, validator.StatusRX), "status", "must be a valid status")

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
	INSERT INTO todo (title, label, task, priority, status, website, address, mode)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id, created_at, version
	`
	//Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//Cleanup to prevent memory leaks
	defer cancel()
	//Collect data fields into a slice
	args := []interface{}{
		todo.Title, todo.Label, todo.Task,
		todo.Priority, todo.Status, todo.Website,
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
		SELECT id, created_at, title, label, task, priority, status, website, address, mode, version
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
		&todo.Title,
		&todo.Label,
		&todo.Task,
		&todo.Priority,
		&todo.Status,
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
		SET title = $1, label = $2, task = $3, 
			priority = $4, status = $5, website = $6, 
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
		todo.Title,
		todo.Label,
		todo.Task,
		todo.Priority,
		todo.Status,
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
func (m TodoModel) GetAll(title string, label string, mode []string, filters Filters) ([]*Todo, Metadata, error) {
	//Construct the query
	query := fmt.Sprintf(`
		SELECT COUNT (*) OVER(), id, created_at, title, 
				label, task,priority, 
				status, website, address, mode, version
		FROM todo
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (to_tsvector('simple', label) @@ plainto_tsquery('simple', $2) OR $2 = '')
		AND (mode @> $3 OR $3 = '{}' )
		ORDER BY %s %s, id ASC
		LIMIT $4 OFFSET $5`, filters.sortColumn(), filters.sortOrder())
	//Create a 3-second-timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	//Execute the query
	args := []interface{}{title, label, pq.Array(mode), filters.limit(), filters.offset()}
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
			&todo.Title,
			&todo.Label,
			&todo.Task,
			&todo.Priority,
			&todo.Status,
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
