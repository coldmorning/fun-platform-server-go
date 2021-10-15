package model

type AccessDetails struct {
	AccessUuid string
	UserId     string
}

type RefreshDetails struct {
	RefreshUuid string
	UserId      string
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}
