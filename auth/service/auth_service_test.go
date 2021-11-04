package authservice

import (

    "testing"
	"github.com/coldmorning/fun-platform/model"
	"github.com/stretchr/testify/assert"

)


func TestCreateAccessToken(t *testing.T){
	td := &model.TokenDetails{}
	userId := "15e3eae4-20a1-40f7-bb7e-380dbf091356"
	td, err := authservice.CreateAccessToken(userId,td)
	assert.NotNil(t, td)
	assert.Equal(t, 36,len( td.AccessToken))
	
}
