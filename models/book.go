package models

type Book struct {
	ID     string `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string `json:"name" bson:"name"`
	Author string `json:"author" bson:"author"`
	ISBN   string `json:"isbn" bson:"isbn"`
	Genre  string `json:"genre" bson:"genre"`
}
