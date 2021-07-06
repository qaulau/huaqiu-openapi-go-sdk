package hqpcb

import (
	"github.com/qaulau/huaqiu-openapi-go-sdk"
)

// 实例化 OpenAPI 客户端对象
func New(appid string, secret string, opts ... openapi.OptionFunc) (*openapi.Client){
	opts = append(opts, openapi.WithAppIdName("appid"))
	opts = append(opts, openapi.WithSignatureName("signature"))
	client := openapi.New(appid, secret, opts...)
	return client
}

