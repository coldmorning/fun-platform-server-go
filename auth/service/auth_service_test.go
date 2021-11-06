package authservice

import (
    "testing"
	"github.com/coldmorning/fun-platform/model"
	"github.com/stretchr/testify/assert"
)






func TestCreateToken(t *testing.T){
	td := &model.TokenDetails{}
	userId := "15e3eae4-20a1-40a7-bb7e-380dbf091357"
	td,_ = CreateAccessToken(userId,td)
	td,_ = CreateRefreshToken(userId,td)
	assert.NotNil(t, td)
	assert.NotNil(t, td.AccessToken)
	assert.NotNil(t, td.RefreshToken)


}

func TestCreateAccessToken(t *testing.T){
	td := &model.TokenDetails{}
	userId := "15e3eae4-20a1-40a7-bb7e-380dbf091357"
	td,_ = CreateAccessToken(userId,td)
	assert.NotNil(t, td)
	assert.NotNil(t, td.AccessToken)


}

func TestCreateRefreshToken(t *testing.T){
	td := &model.TokenDetails{}
	userId := "15e3eae4-20a1-40a7-bb7e-380dbf091357"
	td,_ = CreateRefreshToken(userId,td)
	assert.NotNil(t, td)
	assert.NotNil(t, td.RefreshToken)

}