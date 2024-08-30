package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
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

func initTableUsers() error {
	err := LoadEnv()

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL,
        email TEXT NOT NULL UNIQUE,
        first_name TEXT,
        last_name TEXT,
        url TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`)
	if err != nil {
		return err
	}

	// Check if the table is empty
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return err
	}

	// If the table is empty, insert the default user
	if count == 0 {
		defaultUsername := os.Getenv("DEFAULT_USERNAME")
		defaultPassword := os.Getenv("DEFAULT_PASSWORD")
		defaultEmail := os.Getenv("DEFAULT_EMAIL")
		defaultFirstName := os.Getenv("DEFAULT_FIRST_NAME")
		defaultLastName := os.Getenv("DEFAULT_LAST_NAME")
		defaultURL := os.Getenv("DEFAULT_URL")

		log.Printf("Inserting default user with username: %s, email: %s", defaultUsername, defaultEmail)

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			return err
		}

		// Log the dan hash password
		log.Printf("Hashed password: %s", string(hashedPassword))
		CreateUser(defaultUsername, string(hashedPassword), defaultEmail, defaultFirstName, defaultLastName, defaultURL)
	}

	return nil
}

func initTableUserDevices() error {
	err := LoadEnv()
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS user_devices (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		device_jid TEXT NOT NULL
    )`)
	if err != nil {
		fmt.Println(err)
		return err
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

	initTableUsers()
	initTableUserDevices()

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
	stmt, err := db.Prepare("INSERT INTO users (username, password, email, first_name, last_name, url, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, password, email, firstName, lastName, url)
	if err != nil {
		return err
	}

	log.Println("User created successfully")
	return nil
}

func GetUserByUserId(userid int) (model.User, error) {
	var user model.User

	// Prepare the SQL statement.
	stmt, err := db.Prepare("SELECT id, username, password, email, first_name, last_name, url FROM users WHERE id = ?")
	if err != nil {
		return user, err
	}
	defer stmt.Close()

	// Execute the SQL statement and scan the result into the user struct.
	row := stmt.QueryRow(userid)
	err = row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.FirstName, &user.LastName, &user.Url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Return a custom error when the user is not found
			return user, errors.New("user not found")
		}
		return user, err
	}

	return user, nil
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

func InsertUserDevice(userDevice model.UserDevice) error {

	// err := LoadEnv()
	// if err != nil {
	// 	return err
	// }

	// databasePath := os.Getenv("DB_PATH")
	// if databasePath == "" {
	// 	return errors.New("DB_PATH environment variable not set")
	// }

	// db, err = sql.Open("sqlite3", databasePath)
	// if err != nil {
	// 	return err
	// }

	// err = db.Ping()
	// if err != nil {
	// 	return err
	// }

    // Prepare the SQL statement for inserting data
    stmt, err := db.Prepare("INSERT INTO user_devices (user_id, device_jid) VALUES (?, ?)")
    if err != nil {
        return err
    }
    defer stmt.Close()

    // Execute the SQL statement with the values
    _, err = stmt.Exec(userDevice.UserId, userDevice.DeviceJid)
    if err != nil {
        return err
    }

    return nil
}

func DeleteUserDevice(phone string) error {
	fmt.Println("deleting ", phone)

    // Prepare the SQL statement for inserting data
    stmt, err := db.Prepare("DELETE FROM user_devices WHERE device_jid = ?")
    if err != nil {
        return err
    }
    defer stmt.Close()

    // Execute the SQL statement with the values
    _, err = stmt.Exec(phone)
    if err != nil {
        return err
    }

    return nil
}

func CheckDatabase() error {

	// err := LoadEnv()
	// if err != nil {
	// 	return err
	// }

	// databasePath := os.Getenv("DB_PATH")
	// if databasePath == "" {
	// 	return errors.New("DB_PATH environment variable not set")
	// }

	// db, err = sql.Open("sqlite3", databasePath)
	// if err != nil {
	// 	return err
	// }

	err := db.Ping()
	if err != nil {
		fmt.Println("------- errrr db ------")
		return err
	}

	fmt.Println("------- db ok ------")

    // Prepare the SQL statement for inserting data
    stmt, err := db.Prepare("INSERT INTO user_devices (user_id, device_jid) VALUES (?, ?)")
    if err != nil {
        return err
    }
    defer stmt.Close()

	userDevice := model.UserDevice{
		UserId: 1,
		DeviceJid: "halo",
	}

    // Execute the SQL statement with the values
    _, err = stmt.Exec(userDevice.UserId, userDevice.DeviceJid)
    if err != nil {
		fmt.Println("------- hasil ------", err)
        return err
    }

    return nil
}

func GetUserByClientJID(jid string) (model.UserDevice, error) {
	var userDevice model.UserDevice

	stmt, err := db.Prepare("SELECT id, user_id, device_jid FROM user_devices WHERE device_jid = ?")
	if err != nil {
		return userDevice, err
	}
	defer stmt.Close()

	// Execute the SQL statement and scan the result into the user struct.
	row := stmt.QueryRow(jid)
	err = row.Scan(&userDevice.ID, &userDevice.UserId, &userDevice.DeviceJid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Return a custom error when the user is not found
			return userDevice, errors.New("user not found")
		}
		return userDevice, err
	}

	return userDevice, err
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

func GetUserByID(username string) (model.User, error) {
	var user model.User
	query := `SELECT id, username, email, first_name, last_name, url FROM users WHERE username = ?`
	err := db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName, &user.Url)
	return user, err
}

func UpdateUserProfile(username, firstName, lastName, email, password string) error {
	// Prepare the SQL statement for updating user profile
	var stmt *sql.Stmt
	var err error

	if password != "" {
		// Update profile including password
		stmt, err = db.Prepare("UPDATE users SET first_name=?, last_name=?, email=?, password=?, updated_at=CURRENT_TIMESTAMP WHERE username=?")
		if err != nil {
			return fmt.Errorf("failed to prepare update statement: %v", err)
		}
		defer stmt.Close()

		//hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		//if err != nil {
		//	return fmt.Errorf("failed to hash password: %v", err)
		//}

		result, err := stmt.Exec(firstName, lastName, email, password, username)
		if err != nil {
			return fmt.Errorf("failed to execute update statement: %v", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected: %v", err)
		}
		if rowsAffected == 0 {
			return sql.ErrNoRows
		}
	} else {
		// Update profile without changing the password
		stmt, err = db.Prepare("UPDATE users SET first_name=?, last_name=?, email=?, updated_at=CURRENT_TIMESTAMP WHERE username=?")
		if err != nil {
			return fmt.Errorf("failed to prepare update statement: %v", err)
		}
		defer stmt.Close()

		result, err := stmt.Exec(firstName, lastName, email, username)
		if err != nil {
			return fmt.Errorf("failed to execute update statement: %v", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected: %v", err)
		}
		if rowsAffected == 0 {
			return sql.ErrNoRows
		}
	}

	return nil // Update successful
}
