package mysql

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/rlr524/snippetbox/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	// Retrieve the id and hashed password of the user. If no matching email exists, or the user is not
	// active, return the ErrInvalidCredentials error.
	var id int
	var hashedPassword []byte
	stmt := "SELECT id, hashed_password FROM users WHERE email = ? AND active = TRUE"
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// Check whether the hashed password and plain-text password provided match. If they don't,
	// return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
}

// TODO: Does it make more sense to check the user email using a simple lookup with a UserModel.EmailTaken() function
// and not worry about a race condition, the likelihood of which is near zero? In this case, we're relying on
// MySQL to not change error codes or messages in future versions if we upgrade or on updating this function. Neither
// option is perfect, however it seems like a better option to not couple so tightly to this specific version of MySQL.

func (m *UserModel) Insert(name, email, password string) error {
	// Create a bcrypt hash of the plain-text password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created) VALUES(?, ?, ?, UTC_TIMESTAMP())`

	// Use the Exec() method to insert the user details and hashed password into the users table
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		// If this returns an error, we use the errors.As() function to check whether the error has the
		// type *mysql.MySQLError. If it does, the error will be assigned to the mySQLError variable. We
		// can then check if the error relates to our users_uc_email key by checking the contents
		// of the message string. If it does, we return an ErrDuplicateEmail error.
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
