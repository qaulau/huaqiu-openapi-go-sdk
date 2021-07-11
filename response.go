package openapi

import (
	jsoniter "github.com/json-iterator/go"
	"net/http"
)

type Result struct {
	Code   int64 `json:"response_code"`
	Info   string `json:"response_info"`
	Msg    string `json:"error_message"`
	Data   jsoniter.Any
}

type Response struct {
	http.Response
	client          *Client
	Text         	string
	Content 		[]byte
	json            Result
}

func (this *Response) parse(){
	result := Result{}
	error := jsoniter.Unmarshal(this.Content, &result)
	if error != nil {
		result.Code = 3001
		result.Msg  = error.Error()
		result.Info = "SYSTEM_ERROR"
	}else{
		result.Data = jsoniter.Get(this.Content, this.client.dataName)
	}
	this.json = result
}

func (this *Response) Result() (Result) {
	if this.json.Code == 0 {
		this.parse()
	}
	return this.json
}
