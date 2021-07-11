package hqchip

import (
	"github.com/qaulau/huaqiu-openapi-go-sdk"
)

// 实例化 OpenAPI 客户端对象
func New(appid string, secret string, opts ... openapi.OptionFunc) (*openapi.Client){
	defs := []openapi.OptionFunc{}
	defs = append(defs, openapi.WithDomain("api.hqchip.com"))
	defs = append(defs, opts...)
	defs = append(defs, openapi.WithAppIdName("app_key"))
	defs = append(defs, openapi.WithSignatureName("sign"))
	defs = append(defs, openapi.WithDataName("data"))
	client := openapi.New(appid, secret, defs...)
	return client
}