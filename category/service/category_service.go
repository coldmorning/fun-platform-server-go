package categoryservice
import (
	"time"
	
	"github.com/coldmorning/fun-platform/category/dao/postgresql"
)


func CreateCategoryRequest(request *CreateCategoryRequest) error{
	err := categorypostresql.CreateCategory(request)
	return err
}