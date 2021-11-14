package boardservice
import (
	"github.com/coldmorning/fun-platform/model"
	"github.com/coldmorning/fun-platform/board/dao/postgresql"
)


func CreateBoardRequest(board *model.CreateBoardRequest) error{
	err := boardpostresql.CreateBoard(board)
	return err
}