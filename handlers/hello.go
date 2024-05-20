package handlers

import (
	"net/http"

	"github.com/swarnimcodes/apigo/utils"
)

type Response struct {
	Message string `json:"message"`
}

func Hello(w http.ResponseWriter, r *http.Request) {
	message := "Welcome to APIGo"
	statusCode := http.StatusOK
	utils.SendMessageResponse(w, message, statusCode)
}
