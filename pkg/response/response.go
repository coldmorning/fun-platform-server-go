package cusresponse

import (
	"encoding/json"
)

type err struct {
    BusinessCode  int `json:"business_code"`        
    Msg  string       `json:"msg"`  
    Data interface{}  `json:"data"`   
	HttpCode string `json:"http_status"` 
	Err error
  
}

func NewError(code int, msg string, httpCode string) Error {
    return &err{
        Code: code,
        Msg:  msg,
        Data: nil,
		HttpStatus: httpstatus,
    }
}

func (e *err) i() {}

func (e *err) WithData(data interface{}) Error {
    e.Data = data
    return e
}


func (e *err) GetHttpCode() int {
	return e.HttpCode
}

func (e *err) GetBusinessCode() int {
	return e.BusinessCode
}

func (e *err) GetMsg() string {
	return e.Message
}

func (e *err) GetErr() error {
	return e.Err
}

func (e *err) ToString() string {
	err := &error{
		HttpCode:     e.HttpCode,
		BusinessCode: e.BusinessCode,
		Message:      e.Message,
	}

	raw, _ := json.Marshal(err)
	return string(raw)
}