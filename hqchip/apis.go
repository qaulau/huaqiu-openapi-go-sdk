package hqchip

import (
	"github.com/qaulau/huaqiu-openapi-go-sdk"
)

// 实例化 OpenAPI 客户端对象
func New(appid string, secret string, opts ... openapi.OptionFunc) (*openapi.Client){
	opts = append(opts, openapi.WithSignatureName("sign"))
	opts = append(opts, openapi.WithAppIdName("app_key"))
	client := openapi.New(appid, secret, opts...)
	return client
}