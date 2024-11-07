package models

type User struct{
	ID   string `bson:"_id,omitempty"`
    Name string `json:"name" bson:"name"`
}