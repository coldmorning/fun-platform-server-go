package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Uuid uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email string `json:"email"`
	Is_delete string 
}