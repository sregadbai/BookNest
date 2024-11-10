package models

type Book struct {
	ID     string `json:"id,omitempty" dynamodbav:"id,omitempty"`
	Name   string `json:"name" dynamodbav:"name"`
	Author string `json:"author" dynamodbav:"author"`
	ISBN   string `json:"isbn" dynamodbav:"isbn"`
	Genre  string `json:"genre" dynamodbav:"genre"`
}
