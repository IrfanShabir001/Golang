package models

type Student struct {
	Id   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	Meta Meta   `json:"meta" bson:"meta"`
}
