package domain

import (
	"time"
)
 
type User struct {

	Uuid string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email string `json:"email"`
	Delete_time time.Time
	Create_time time.Time
	Update_time time.Time
	Create_by string
	Update_by string
}