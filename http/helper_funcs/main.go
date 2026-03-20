package helperfuncs

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalPages int `json:"total_pages"`
	Total      int `json:"total"`
}

func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func RespondSuccess(w http.ResponseWriter, data interface{}) {
	RespondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

func RespondError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, Response{
		Success: false,
		Error:   message,
	})
}

func RespondPaginated(w http.ResponseWriter, data interface{}, meta *Meta) {
	RespondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

// Usage
func getUsersFromDB() (any, error)

func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := getUsersFromDB()
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	meta := &Meta{
		Page:       1,
		PerPage:    10,
		TotalPages: 5,
		Total:      50,
	}

	RespondPaginated(w, users, meta)
}
