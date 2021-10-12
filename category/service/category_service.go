package categoryservice
import (
	"github.com/coldmorning/fun-platform/model"
	"github.com/coldmorning/fun-platform/category/dao/postgresql"
)


func CreateCategoryRequest(category *model.CreateCategoryRequest) error{
	err := categorypostresql.CreateCategory(category)
	return err
}