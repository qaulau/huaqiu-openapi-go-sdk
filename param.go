package openapi

import (
	"bytes"
	"net/url"
)

type Params map[string]any

// 获取转换参数类型
func NewParams(dict *map[string]interface{}) *Params {
	params := &Params{}
	for k, v := range *dict{
		params.Add(k, v)
	}
	return params
}

// 合并新的参数
func (p *Params) Merge(params *url.Values) (*Params)  {
	c := Params{}
	for k, v := range *p{
		c[k] = v
	}
	if params == nil{
		return &c
	}
	for k, v := range *params{
		c[k] = NewAny(v[0])
	}
	return &c
}

// 更新原有的数据
func (p *Params) Update (params *Params) (Params) {
	if params == nil {
		return *p
	}
	for k, v := range *params{
		(*p)[k] = v
	}
	return *p
}

// 转换为标准库 url 数值类型
func (p *Params) UrlValues() (*url.Values) {
	var u = url.Values{}
	for k, v := range *p {
		u.Add(k, v.String())
	}
	return &u
}

// 转换为请求body 数值类型
func (p *Params) Buffer () (*bytes.Buffer) {
	v := p.UrlValues()
	b := []byte(v.Encode())
	buf := &bytes.Buffer{}
	buf.Write(b)
	return buf
}

// 获取指定键值
func (p *Params) Get(name string) string  {
	v, ok := (*p)[name]
	if !ok {
		return ""
	}
	return v.String()
}

// 设置值
func (p *Params) Set(name string, value interface{})  {
	(*p)[name] = NewAny(value)
}

// 设置值
func (p *Params) Add(name string, value interface{})  {
	(*p)[name] = NewAny(value)
}

// 检测是否存在减值
func (p *Params) Existed(name string) bool {
	_, ok := (*p)[name]
	return ok
}

// 检测键值大小
func (p *Params) Size() (size int)  {
	size = len(*p)
	return
}