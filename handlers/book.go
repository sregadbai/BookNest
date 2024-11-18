package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sregadbai/BookNest/db"
	"github.com/sregadbai/BookNest/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetBooks fetches all books
func GetBooks(w http.ResponseWriter, r *http.Request) {
	collection := db.GetCollection("books")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var books []models.Book
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch books", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var book models.Book
		if err := cursor.Decode(&book); err != nil {
			http.Error(w, "Failed to decode book", http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// GetBook fetches a single book by ID
func GetBook(w http.ResponseWriter, r *http.Request) {
	collection := db.GetCollection("books")
	vars := mux.Vars(r)
	id := vars["id"]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var book models.Book
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&book)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Book not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// CreateBook adds a new book
func CreateBook(w http.ResponseWriter, r *http.Request) {
	collection := db.GetCollection("books")

	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Generate a UUID if the id field is empty
	if book.ID == "" {
		book.ID = uuid.New().String()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, book)
	if err != nil {
		http.Error(w, "Failed to add book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

// UpdateBook updates a book by ID
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	collection := db.GetCollection("books")
	vars := mux.Vars(r)
	id := vars["id"]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if the book exists
	var existingBook models.Book
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&existingBook)
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	// Decode the new book data
	var updatedBook models.Book
	if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update the book
	update := bson.M{"$set": updatedBook}
	_, err = collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		http.Error(w, "Failed to update book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteBook deletes a single book by ID
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	collection := db.GetCollection("books")
	vars := mux.Vars(r)
	id := vars["id"]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if the book exists
	var existingBook models.Book
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&existingBook)
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	// Delete the book
	_, err = collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		http.Error(w, "Failed to delete book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteAllBooks deletes all books
func DeleteAllBooks(w http.ResponseWriter, r *http.Request) {
	collection := db.GetCollection("books")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to delete books", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
