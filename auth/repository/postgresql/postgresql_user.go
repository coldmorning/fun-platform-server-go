package postgresql

import (
	"fun-platform-server/domain"
)


func FetchUser(u domain.User)(*domain.User,error){

	//get user from database.
	
	var user = domain.User{
		ID: 1,
		Username: "username@gmail.com",
		Password: "password",
	}
	return &user,nil;

}