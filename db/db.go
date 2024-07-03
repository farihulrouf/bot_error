package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"wagobot.com/model"
)

var db *sql.DB // Define db as a global variable

func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}
	return nil
}

func InitDB() error {
	err := LoadEnv()
	if err != nil {
		return err
	}

	databasePath := os.Getenv("DB_PATH")
	if databasePath == "" {
		return errors.New("DB_PATH environment variable not set")
	}

	db, err = sql.Open("sqlite3", databasePath)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		first_name TEXT,
		last_name TEXT,
		url TEXT
	)`)
	if err != nil {
		return err
	}

	log.Println("Connected to the database")
	return nil
}

func CloseDB() {
	if db != nil {
		db.Close()
		log.Println("Database connection closed")
	}
}

func CreateUser(username, password, email, firstName, lastName, url string) error {
	// Prepare the SQL statement.
	stmt, err := db.Prepare("INSERT INTO users (username, password, email, first_name, last_name, url) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement with user data.
	_, err = stmt.Exec(username, password, email, firstName, lastName, url)
	if err != nil {
		return err
	}

	log.Println("User created successfully")
	return nil
}

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

func OpenDatabase() (*sql.DB, error) {
	err := LoadEnv()
	if err != nil {
		return nil, err
	}

	databasePath := os.Getenv("DB_PATH")
	if databasePath == "" {
		return nil, errors.New("DB_PATH environment variable not set")
	}

	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	// Ensure the database connection is successful.
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping the database: %w", err)
	}

	return db, nil
}

func UpdateUserURLWebhook(Username string, url string) error {
	query := `UPDATE users SET url = ? WHERE username = ?`
	_, err := db.Exec(query, url, Username)
	return err
}
