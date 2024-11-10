package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sregadbai/BookNest/db"
	"github.com/sregadbai/BookNest/handlers"
	"github.com/sregadbai/BookNest/logger"
)

func loadEnv() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}
	envFile := ".env." + env
	err := godotenv.Load(envFile)
	if err != nil {
		log.Printf("Error loading %s file, relying on system environment variables\n", envFile)
	}
}

func main() {
	logger.InitLogger()
	// Load environment variables
	loadEnv()

	// Initialize DynamoDB connection
	db.Connect()

	// Set up routes
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/books", handlers.CreateBook).Methods("POST")
	r.HandleFunc("/api/v1/books", handlers.GetBooks).Methods("GET")
	r.HandleFunc("/api/v1/books/book/{id}", handlers.GetBook).Methods("GET")
	r.HandleFunc("/api/v1/books/book/{id}", handlers.UpdateBook).Methods("PUT")
	r.HandleFunc("/api/v1/books/book/{id}", handlers.DeleteBook).Methods("DELETE")
	r.HandleFunc("/api/v1/books", handlers.DeleteAllBooks).Methods("DELETE")

	// Start the server
	logger.Log.Info("Server is running..")
	logger.Log.Fatal(http.ListenAndServe(":8080", r))
}
