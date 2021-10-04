package domain

import (
	"time"
)

type category struct{
	Uuid string `json:"uuid"`
	Name string `json:"name"`
	Picture_path string `json:"picture_path"`
	Managers_master string  `json:"managers_master"`
	Managers string `json:"managers"`
	
	Delete_time time.Time
	Create_time time.Time
	Update_time time.Time
	Create_by string
	Update_by string
	Delete_by string
}