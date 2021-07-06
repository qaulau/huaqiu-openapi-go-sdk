package openapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Result struct {
	code 			int64 `json:"response_code"`
	info            string `json:"response_info"`
	msg             string `json:"response_message"`
	data            map[string]interface{} `json:"response_data"`
}

type Response struct {
	http.Response
	Content         string
	Data            Result
}

func (this *Response) parse(){
	result := Result{}
	body, _ := ioutil.ReadAll(this.Body)
	error := json.Unmarshal(body, result)
	if error != nil {
		result.code = 3001
		result.msg  = error.Error()
		result.info = "SYSTEM_ERROR"
	}
	this.Data = result
	this.Content = string(body)
}
