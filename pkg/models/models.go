package models

// model for database
type User struct {
	ID           int
	Company      string
	Email        string
	Password     string
	Access_level string
}
