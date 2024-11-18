package main

import (
	"log"
	"net/http"
	"os"

	corsHandler "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sregadbai/BookNest/db"
	"github.com/sregadbai/BookNest/handlers"
	"github.com/sregadbai/BookNest/logger"
)

func loadEnv() {
	// Determine environment (default to "local")
	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}

	// Skip loading .env in CI/CD environments
	if os.Getenv("CI") == "true" {
		log.Println("Running in CI/CD environment, relying on system environment variables only")
		return
	}

	// Load the appropriate .env file
	envFile := env + ".env"
	err := godotenv.Load(envFile)
	if err != nil {
		log.Printf("Error loading %s file, relying on system environment variables\n", envFile)
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func main() {
	// Initialize logger
	logger.InitLogger()

	// Load environment variables
	loadEnv()

	// Initialize MongoDB connection
	dbURI := os.Getenv("MONGO_URI")
	if dbURI == "" {
		dbURI = "mongodb://localhost:27017"
	}
	db.ConnectMongoDB(dbURI)

	// Create Gorilla Mux router
	r := mux.NewRouter()

	// Define API routes
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/books", handlers.CreateBook).Methods("POST")             // Create a book
	api.HandleFunc("/books", handlers.GetBooks).Methods("GET")                // Get all books
	api.HandleFunc("/books/book/{id}", handlers.GetBook).Methods("GET")       // Get a book by ID
	api.HandleFunc("/books/book/{id}", handlers.UpdateBook).Methods("PUT")    // Update a book by ID
	api.HandleFunc("/books/book/{id}", handlers.DeleteBook).Methods("DELETE") // Delete a book by ID
	api.HandleFunc("/books", handlers.DeleteAllBooks).Methods("DELETE")       // Delete all books

	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Apply CORS
	cors := corsHandler.CORS(
		corsHandler.AllowedOrigins([]string{"http://localhost:3000"}), // Frontend origin
		corsHandler.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		corsHandler.AllowedHeaders([]string{"Content-Type"}),
	)

	// Start the server
	port := getEnv("PORT", "8080")
	logger.Log.Infof("Server is running on port %s", port)
	logger.Log.Fatal(http.ListenAndServe(":"+port, cors(r)))
}
