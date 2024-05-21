package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/swarnimcodes/apigo/utils"
)

func FetchAllUsers(w http.ResponseWriter, r *http.Request) {
	dbLocation := os.Getenv("DATABASE_LOCATION")
	if dbLocation == "" {
		message := "No database location specifed"
		statusCode := http.StatusInternalServerError
		utils.SendErrorResponse(w, message, statusCode)
		return
	}
}

func FetchUser(w http.ResponseWriter, r *http.Request) {
	// Get id to look for
	// id := strings.TrimPrefix(r.URL.Path, "/user/")
	id := r.URL.Query().Get("id")
	if id == "" {
		message := "No id provided"
		statusCode := http.StatusBadRequest
		utils.SendErrorResponse(w, message, statusCode)
		return
	}

	// Check if database location is specified
	dbLocation := os.Getenv("DATABASE_LOCATION")
	if dbLocation == "" {
		message := "No database location specifed"
		statusCode := http.StatusInternalServerError
		utils.SendErrorResponse(w, message, statusCode)
		return
	}

	db, err := sql.Open("sqlite3", dbLocation)
	if err != nil {
		message := fmt.Sprintf("Failed to connect to the database: %v", err)
		statusCode := http.StatusInternalServerError
		utils.SendErrorResponse(w, message, statusCode)
		return
	}
	defer db.Close()

	var user User
	query := `SELECT firstName, lastName, password, location
	FROM user WHERE id = ?`

	err = db.QueryRow(query, id).Scan(&user.FirstName, &user.LastName, &user.Password, &user.Location)
	if err != nil {
		if err == sql.ErrNoRows {
			message := "No such user found"
			statusCode := http.StatusNotFound
			utils.SendErrorResponse(w, message, statusCode)

		} else {
			message := fmt.Sprintf("Failed to fetch user details: %v", err)
			statusCode := http.StatusInternalServerError
			utils.SendErrorResponse(w, message, statusCode)

		}
		return
	}

	// message := ""
	// statusCode := http.StatusOK
	// utils.SendMessageResponse()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
