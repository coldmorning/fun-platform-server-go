package model

import (
	"time"
)

type board struct{
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


type BoardRequest struct{
	Uuid uint32 `form:"uuid" binding:"required"`

}



type BoardListRequest struct {
	
	
}

type CreateBoardRequest struct {
	Uuid string 
	Name string `form:"name" binding:"required",max=150`
	Picture_path string `form:"picture_path" binding:"required",max=300`
	Managers_master string `form:"Managers_master" binding:"required",max=40`
	Managers string  `form:"managers" binding:"required",max=400`
	
	Create_time time.Time `form:"create_time" binding:"required",min=1`
	Create_by string `form:"create_by" binding:"required",min=1`
}



type UpdateBoardRequest struct {
	Uuid string `form:"uuid" binding:"required",max=40`
	Name string `form:"name" binding:"required",max=150`
	Picture_path string `form:"picture_path" binding:"required",max=300`
	Managers_master string `form:"Managers_master" binding:"required",max=40`
	Managers string  `form:"managers" binding:"required",max=400`
	
	Update_time time.Time `form:"update_time" binding:"required",min=1`
	Update_by string `form:"update_by" binding:"required",min=1`
	
}

type DeleteBoardRequest struct {
	Uuid uint32 `form:"uuid" binding:"required,gte=1"`
	Delete_time time.Time `form:"delete_time" binding:"required",min=1`
	Delete_by string `form:"delete_by" binding:"required",min=1`
}