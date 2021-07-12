HuaQiu OpenAPI Go SDK
=================

# Install
```$ go get github.com/qaulau/huaqiu-openapi-go-sdk```


### Usage

- HuaQiu HQCHIP OpenAPI

order query
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
	params := &openapi.Params{}
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

order make
```
package main

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	openapi "github.com/qaulau/huaqiu-openapi-go-sdk"
	"github.com/qaulau/huaqiu-openapi-go-sdk/hqchip"
)

type dict map[string]interface{}

func main(){
	client := hqchip.New("Your HQCHIP AppKey", "Your HQCHIP AppSecret")
	goods := []dict{
		dict{
			"out_goods_name":"ULN2803ADWR",
			"qty": 10,
			"goods_id": 2500217756,
		},
		dict{
			"out_goods_name":"CL10A106KP8NNNC",
			"qty": 50,
			"goods_id": 2500014137,
		},
	}
	invoice := dict{
		"type": 1,
		"inv_title": "刘权",
	}
	recive := dict{
		"consignee": "刘权",
		"province": 6,
		"city": 77,
		"district": 705,
		"address": "深圳市福田区新一代产业园1栋5层",
		"mobile": "13600000001",
		"tel": "",
	}
	goods_json, _ := jsoniter.Marshal(goods)
	invoice_json, _ := jsoniter.Marshal(invoice)
	recive_json, _ := jsoniter.Marshal(recive)
	data := dict{
		"goods_list": string(goods_json),
		"invoice": string(invoice_json),
		"receive": string(recive_json),
		"shipping_type": 1,
		"goods_type": 1,
		"out_order_no": "T202020930001",
		"product_num": "600",
	}
	params := openapi.NewParams(data)
	resp, err := client.Post("/order/make/", params)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Text)
	fmt.Printf("success make order: %s", resp.Result().Data.Get("order_sn").ToString())
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
	params := &openapi.Params{}
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
