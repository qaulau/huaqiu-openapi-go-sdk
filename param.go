package openapi

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
)

type Params map[string]interface{}

// 合并新的参数
func (p Params) Merge(params url.Values) (Params)  {
	c := Params{}
	for k, v := range p{
		c[k] = v
	}
	if params == nil{
		return c
	}
	for k, v := range params{
		c[k] = v
	}
	return c
}

// 更新原有的数据
func (p Params) Update (params Params) (Params) {
	if params == nil {
		return p
	}
	for k, v := range params{
		p[k] = v
	}
	return p
}

// 转换为标准库 url 数值类型
func (p Params) UrlValues() (url.Values) {
	var u = url.Values{}
	for k, v := range p {
		switch v.(type) {
		case string:
			s := v.(string)
			u.Add(k, s)
		default:
			s := fmt.Sprintf("%s", v)
			u.Add(k, s)
		}
	}
	return u
}

// 转换为请求body 数值类型
func (p Params) Body () (io.ReadCloser) {
	v := p.UrlValues()
	b := []byte(v.Encode())
	return ioutil.NopCloser(bytes.NewReader(b))
}