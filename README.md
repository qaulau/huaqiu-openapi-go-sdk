HuaQiu OpenAPI Go SDK
=================

# Install
```$ go get github.com/qaulau/huaqiu-openapi-go-sdk```


### Usage

- HuaQiu HQCHIP OpenAPI
```
package main

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	openapi "github.com/qaulau/huaqiu-openapi-go-sdk"
	"github.com/qaulau/huaqiu-openapi-go-sdk/hqchip"
)

func main(){
	client := hqchip.New("Your HQCHIP AppKey", "Your HQCHIP AppSecret")
	params := openapi.Params{}
	params.Add("order_id", 257194)
	resp, err := client.Get("/order/detail/", params)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Result().Code, resp.Result().Msg, resp.Result().Data.ToString())
	goods := []map[string]jsoniter.Any{}
	resp.Result().Data.Get("goods_list").ToVal(&goods)
	fmt.Println(goods[0]["avg_unit_price"].ToFloat64())
}
```

- HuaQiu HQPCB OpenAPI
```
package main

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	openapi "github.com/qaulau/huaqiu-openapi-go-sdk"
	"github.com/qaulau/huaqiu-openapi-go-sdk/hqpcb"
)

func main(){
	client := hqpcb.New("Your HQPCB AppID", "Your HQPCB AppSecret")
	params := openapi.Params{}
	params.Add("order_id", 1574974)
	resp, err := client.Get("/order/", params)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Result().Code, resp.Result().Msg, resp.Result().Data.ToString())
	detail := make(map[string]string)
	resp.Result().Data.Get("order_detail").ToVal(&detail)
	fmt.Println(detail["bwidth"])
}
```
