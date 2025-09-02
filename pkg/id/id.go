package id

import "github.com/google/uuid"

func New() string {
	uuidv7, _ := uuid.NewV7()
	return uuidv7.String()
}
