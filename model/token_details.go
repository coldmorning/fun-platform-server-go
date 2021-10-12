package model

type TokenDetails struct{
	AccessToken string
	RefreshToken string
	AcessUuid string
	RefreshUuid string
	AtExpires int64
	RtExpires int64
}