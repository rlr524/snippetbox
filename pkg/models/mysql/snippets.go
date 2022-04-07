package mysql

import (
	"database/sql"
	"errors"

	"github.com/rlr524/snippetbox/pkg/models"
)

// SnippetModel type which wraps a sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// Insert function inserts a new snippet into the database
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Get function returns a specific snippet based on its ID
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	// SQL statement to execute
	stmt := `SELECT id, title, content, created, expires FROM snippets
WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// Use QueryRow() method on the connection pool to execute the statement, passing in the untrusted id
	// variable as the value for the placeholder parameter. This returns a pointer to a sql.Row object which
	// holds the result from the database.
	row := m.DB.QueryRow(stmt, id)

	// Initialize a pointer to a new zeroed Snippet struct
	s := &models.Snippet{}

	// Use row.Scan() to copy the values from each field in sql.Row to the corresponding field in the Snippet
	// struct. Notice that the arguments to row.Scan are &pointers to the place we want to copy the data into, and
	// the number of arguments must be exactly the same as the number of columns returned by the statement.
	err := row.Scan(&s.ID, &s.Content, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// If the query returns no rows, then row.Scan() will return a sql.ErrNoRows error. We use the errors.Is()
		// function to check for that error specifically, and return our own models.ErrNoRecord instead
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	// If everything went ok, then return the Snippet object
	return s, nil
}

// Latest function returns the 10 most recently created snippets
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}
