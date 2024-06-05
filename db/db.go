package db

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"wagobot.com/model"
)

var db *sql.DB // Define db as a global variable

// InitDB initializes the database connection.
func InitDB(databasePath string) error {
	var err error
	db, err = sql.Open("sqlite3", databasePath)
	if err != nil {
		return err
	}

	// Ensure the database connection is successful.
	err = db.Ping()
	if err != nil {
		return err
	}

	log.Println("Connected to the database")
	return nil
}

// CloseDB closes the database connection.
func CloseDB() {
	if db != nil {
		db.Close()
		log.Println("Database connection closed")
	}
}

// CreateUser inserts a new user into the database.
func CreateUser(username, password, email, firstName, lastName string) error {
	// Prepare the SQL statement.
	stmt, err := db.Prepare("INSERT INTO users (username, password, email, first_name, last_name) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement with user data.
	_, err = stmt.Exec(username, password, email, firstName, lastName)
	if err != nil {
		return err
	}

	log.Println("User created successfully")
	return nil
}

// GetUserByUsername retrieves a user from the database by username.
func GetUserByUsername(username string) (model.User, error) {
	var user model.User

	// Prepare the SQL statement.
	stmt, err := db.Prepare("SELECT id, username, password, email, first_name, last_name FROM users WHERE username = ?")
	if err != nil {
		return user, err
	}
	defer stmt.Close()

	// Execute the SQL statement and scan the result into the user struct.
	row := stmt.QueryRow(username)
	err = row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.FirstName, &user.LastName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Return a custom error when the user is not found
			return user, errors.New("user not found")
		}
		return user, err
	}

	return user, nil
}
