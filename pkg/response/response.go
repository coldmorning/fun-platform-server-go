package cusresponse

import (
	"encoding/json"
)

type err struct {
    BusinessCode  int `json:"business_code"`        
    Msg  string       `json:"msg"`  
    Data interface{}  `json:"data"`   
	HttpCode int `json:"http_status"` 
	Err error
  
}

func NewError(code int, msg string, httpCode int)  *err {
    return &err{
        BusinessCode: code,
        Msg:  msg,
        Data: nil,
		HttpCode: httpCode,
    }
}

func (e *err) i() {}




func (e *err) GetHttpCode() int {
	return e.HttpCode
}

func (e *err) GetBusinessCode() int {
	return e.BusinessCode
}

func (e *err) GetMsg() string {
	return e.Msg
}

func (e *err) GetErr() error {
	return e.Err
}

func (e *err) ToString() string {
	err := &err{
		HttpCode:     e.HttpCode,
		BusinessCode: e.BusinessCode,
		Msg:      e.Msg,
	}

	raw, _ := json.Marshal(err)
	return string(raw)
}