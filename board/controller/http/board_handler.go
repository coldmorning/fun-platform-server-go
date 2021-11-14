package boardhttp

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	
	"github.com/coldmorning/fun-platform/model"
    "github.com/coldmorning/fun-platform/board/service"

)



func List(ctx *gin.Context){
	var form model.BoardRequest
	if err := ctx.ShouldBind(&form); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
	}
}

func Create(ctx *gin.Context){
	var form *model.CreateBoardRequest
	if err := ctx.ShouldBind(&form); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
	}
	err := boardservice.CreateBoardRequest(form)
	if err !=nil {
		//422
		ctx.JSON(http.StatusUnprocessableEntity, err.Error())	
	}else{
		//201
		ctx.JSON(http.StatusCreated, "ok")
	}

}

func Delete(ctx *gin.Context){
}



func Update(ctx *gin.Context){
}
