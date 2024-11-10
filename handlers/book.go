package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/sregadbai/BookNest/db"
	"github.com/sregadbai/BookNest/logger"
	"github.com/sregadbai/BookNest/models"
)

func CreateBook(w http.ResponseWriter, r *http.Request) {
	var book *models.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if book.ID == "" {
		book.ID = uuid.NewString()
	}

	av, err := attributevalue.MarshalMap(book)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to marshal book data")
		http.Error(w, "Failed to marshal book data", http.StatusInternalServerError)
		return
	}

	_, err = db.DynamoDBClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("Books"),
		Item:      av,
	})
	if err != nil {
		logger.Log.WithError(err).Error("Failed to put item in DynamoDB")
		http.Error(w, "Failed to put item", http.StatusInternalServerError)
		return
	}
	logger.Log.WithFields(logrus.Fields{
		"book_id": book.ID,
	}).Info("Book created successfully")
	json.NewEncoder(w).Encode(book)
}

func GetBooks(w http.ResponseWriter, r *http.Request) {
	result, err := db.DynamoDBClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String("Books"),
	})
	if err != nil {
		logger.Log.WithError(err).Error("Failed to retrieve books")
		http.Error(w, "Failed to retrieve books", http.StatusInternalServerError)
		return
	}

	var books []models.Book
	err = attributevalue.UnmarshalListOfMaps(result.Items, &books)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to unmarshal book data")
		http.Error(w, "Failed to unmarshal book data", http.StatusInternalServerError)
		return
	}
	logger.Log.Info("Books retrieved successfully")
	json.NewEncoder(w).Encode(books)
}

func GetBook(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	result, err := db.DynamoDBClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("Books"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil || result.Item == nil {
		logger.Log.WithError(err).Error("Book not found")
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	var book models.Book
	err = attributevalue.UnmarshalMap(result.Item, &book)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to unmarshal book data")
		http.Error(w, "Failed to unmarshal book data", http.StatusInternalServerError)
		return
	}

	logger.Log.WithField("book_id", id).Info("Book retrieved successfully")
	json.NewEncoder(w).Encode(book)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	result, err := db.DynamoDBClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("Books"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil || result.Item == nil {
		logger.Log.WithError(err).Error("Book not found")
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	var updatedBook models.Book
	err = json.NewDecoder(r.Body).Decode(&updatedBook)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	updatedBook.ID = id
	av, err := attributevalue.MarshalMap(updatedBook)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to marshal updated book data")
		http.Error(w, "Failed to marshal updated book data", http.StatusInternalServerError)
		return
	}
	_, err = db.DynamoDBClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("Books"),
		Item:      av,
	})
	if err != nil {
		logger.Log.WithError(err).Error("Failed to update book in DynamoDB")
		http.Error(w, "Failed to update book in DynamoDB", http.StatusInternalServerError)
		return
	}
	logger.Log.WithFields(logrus.Fields{"book_id": id}).Info("Book updated successfully")
	json.NewEncoder(w).Encode(updatedBook)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	result, err := db.DynamoDBClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("Books"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		logger.Log.WithError(err).Error("Failed to retrieve book for deletion")
		http.Error(w, "Failed to retrieve book", http.StatusInternalServerError)
		return
	}

	// If the item does not exist, return 404
	if result.Item == nil {
		logger.Log.WithField("book_id", id).Warn("Book not found for deletion")
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	_, err = db.DynamoDBClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String("Books"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		logger.Log.WithError(err).Error("Failed to delete book")
		http.Error(w, "Failed to delete book", http.StatusInternalServerError)
		return
	}
	logger.Log.WithField("book_id", id).Info("Book deleted successfully")
	w.WriteHeader(http.StatusNoContent)
}

func DeleteAllBooks(w http.ResponseWriter, r *http.Request) {
	result, err := db.DynamoDBClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String("Books"),
	})
	if err != nil {
		logger.Log.WithError(err).Error("Failed to scan books")
		http.Error(w, "Failed to retrive books", http.StatusInternalServerError)
		return
	}

	if len(result.Items) == 0 {
		logger.Log.Info("No books to delete")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	writeRequests := []types.WriteRequest{}
	for _, item := range result.Items {
		writeRequests = append(writeRequests, types.WriteRequest{
			DeleteRequest: &types.DeleteRequest{
				Key: map[string]types.AttributeValue{
					"id": item["id"],
				},
			},
		})
	}

	for i := 0; i < len(writeRequests); i += 25 { // DynamoDB limit per BatchEriteItem
		end := i + 25
		if end > len(writeRequests) {
			end = len(writeRequests)
		}

		_, err = db.DynamoDBClient.BatchWriteItem(context.TODO(), &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				"Books": writeRequests[i:end],
			},
		})
		if err != nil {
			logger.Log.WithError(err).Error("Failed to delete books in batch")
			http.Error(w, "Failed to delete books in batch", http.StatusInternalServerError)
			return
		}
		logger.Log.Info("All books deleted successfully")
		w.WriteHeader(http.StatusNoContent)
	}
}
