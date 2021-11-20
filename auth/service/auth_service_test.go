package authservice

import (

    "testing"
	"time"
	"github.com/coldmorning/fun-platform/model"
	"github.com/stretchr/testify/assert"
	"github.com/dgrijalva/jwt-go"
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

func TestVerifyTokenOk(t *testing.T) {
	td := &model.TokenDetails{}
	userId := "15e3eae4-20a1-40a7-bb7e-380dbf091357"
	td,_ = CreateAccessToken(userId,td)
	token, _:= VerifyToken(Access_secret,td.AccessToken)
	assert.NotNil(t, token)
}


func TestVerifyTokenErrorExpired(t *testing.T) {
	//Create expect token
	var err error
	td := &model.TokenDetails{}
	td.AtExpires = time.Now().Add( -1 * time.Second).Unix()
	atClaims := jwt.MapClaims{}
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)
	td.AccessToken, err = at.SignedString(Access_secret)
	if err == nil {
		_, errToken := VerifyToken(Access_secret,td.AccessToken)
		  // assert for not nil (good when you expect something)
		  
		  if assert.NotNil(t,errToken) {
			// now we know that object isn't nil, we are safe to make
			// further assertions without causing any errors
			assert.Equal(t, errToken.Error(), ErrExpiredToken.Error())
		  }else{
			  t.Errorf("errToken is nil")
		  }
	}else{
		t.Errorf("AccessToken can not create correctly")
	}
}

func TestVerifyTokenErrorMalfFormat(t *testing.T) {
	var err error
	td := &model.TokenDetails{}
	td.AtExpires = time.Now().Add( 1 * time.Second).Unix()
	atClaims := jwt.MapClaims{}
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)
	td.AccessToken, err = at.SignedString(Access_secret)
	errorTokenFormat := "bearer "+td.AccessToken
	
	if err == nil {
		_, errToken := VerifyToken(Access_secret,errorTokenFormat)
		  // assert for not nil (good when you expect something)
	
		  if assert.NotNil(t,errToken) {
			// now we know that object isn't nil, we are safe to make
			// further assertions without causing any errors
			assert.Equal(t,ErrorMalformedToken.Error(), errToken.Error())
		  }else{
			  t.Errorf("errToken is nil")
		  }
	}else{
		t.Errorf("AccessToken can not create correctly")
	}
}

func TestVerifyTokenErrorOther(t *testing.T) {
	var err error
	td := &model.TokenDetails{}
	td.AtExpires = time.Now().Add( 1 * time.Second).Unix()
	atClaims := jwt.MapClaims{}
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)
	td.AccessToken, err = at.SignedString(Access_secret)
	errorTokenFormat := td.AccessToken +"123"
	
	if err == nil {
		_, errToken := VerifyToken(Access_secret,errorTokenFormat)
		  // assert for not nil (good when you expect something)
		  
		  if assert.NotNil(t,errToken) {
			// now we know that object isn't nil, we are safe to make
			// further assertions without causing any errors
			assert.Equal(t,ErrorOtherToken.Error(), errToken.Error())
		  }else{
			  t.Errorf("errToken is nil")
		  }
	}else{
		t.Errorf("AccessToken can not create correctly")
	}
}