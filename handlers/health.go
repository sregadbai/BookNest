package handlers

import (
	"net/http"

	"github.com/sregadbai/BookNest/db"
)

// HealthCheck is a handler that checks the health of the service
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Check if the database client is connected
	if db.Client == nil {
		http.Error(w, "Database not connected", http.StatusInternalServerError)
		return
	}

	// If everything is okay, return 200 OK
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
