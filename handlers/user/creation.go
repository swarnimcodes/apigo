package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/swarnimcodes/apigo/utils"
)

func createUserTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS user (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		firstName TEXT NOT NULL,
		lastName TEXT NOT NULL,
		password TEXT NOT NULL,
		location TEXT NOT NULL
	)
	`)
	return err
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	dbLocation := os.Getenv("DATABASE_LOCATION")
	if dbLocation == "" {
		message := "No database location specified"
		statusCode := http.StatusInternalServerError
		utils.SendErrorResponse(w, message, statusCode)
		return
	}

	// Ensure database location directory exists
	if err := os.MkdirAll(filepath.Dir(dbLocation), 0755); err != nil {
		message := fmt.Sprintf("Failed to create directories for database creation: %v", err)
		statusCode := http.StatusInternalServerError
		utils.SendErrorResponse(w, message, statusCode)
		return
	}

	// Open or create the .db file
	// Establish connection with the database
	db, err := sql.Open("sqlite3", dbLocation)
	if err != nil {
		message := fmt.Sprintf("Failed to connect to the database: %v", err)
		statusCode := http.StatusInternalServerError
		utils.SendErrorResponse(w, message, statusCode)
		return
	}
	defer db.Close()

	// Create user table if it doesn't exist
	if err := createUserTable(db); err != nil {
		message := fmt.Sprintf("Failed to create user table in the database: %v", err)
		statusCode := http.StatusInternalServerError
		utils.SendErrorResponse(w, message, statusCode)
		return
	}

	// parse request body into User struct
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		message := fmt.Sprintf("Invalid JSON request input payload: %v", err)
		statusCode := http.StatusBadRequest
		utils.SendErrorResponse(w, message, statusCode)
		return
	}

	// Prepare SQL statement
	stmt, err := db.Prepare("INSERT INTO user (firstName, lastName, password, location) VALUES (?, ?, ?, ?)")
	if err != nil {
		message := fmt.Sprintf("Failed to prepare insert statement: %v", err)
		statusCode := http.StatusInternalServerError
		utils.SendErrorResponse(w, message, statusCode)
		return
	}
	defer stmt.Close()

	// Execute SQL statement to insert new user
	_, err = stmt.Exec(user.FirstName, user.LastName, user.Password, user.Location)
	if err != nil {
		message := fmt.Sprintf("Failed to insert user into database: %v", err)
		statusCode := http.StatusInternalServerError
		utils.SendErrorResponse(w, message, statusCode)
		return
	}
	message := "User saved successfully in database"
	statusCode := http.StatusOK
	utils.SendMessageResponse(w, message, statusCode)
}
