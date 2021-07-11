package openapi

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"sort"
	"time"
)

// 获取字符串 MD5 16进制值
func Md5digest(data string) string{
	m := md5.New()
	io.WriteString(m, data)
	return fmt.Sprintf("%x", m.Sum(nil))
}

// 生成签名
func GenSign(secret string, params Params) string{
	var keys []string
	buf := bytes.Buffer{}
	buf.WriteString(secret)
	for k, _ := range params {
		if k != "sign" && k != "signature" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	len := len(keys) - 1
	for i, v := range keys {
		buf.WriteString(v)
		buf.WriteString("=")
		buf.WriteString(params.Get(v))
		if i != len {
			buf.WriteString("&")
		}
	}
	buf.WriteString(secret)
	data := buf.String()
	return Md5digest(data)
}

// 获取GM+8时区(中国官方标准时区)时间戳
func UnixTime()  int64{
	return UnixTimeByLocation("Asia/Shanghai")
}

// 获取指定时区的时间戳
func UnixTimeByLocation(location string) int64 {
	loc, err := time.LoadLocation(location)
	if err != nil{
		return 0
	}
	return time.Now().In(loc).Unix()
}

// 转换字符串
func String(v interface{}) string{
	switch v.(type) {
	case string:
		s := v.(string)
		return s
	case int,int32,int64:
		s := fmt.Sprintf("%d", v)
		return s
	default:
		s := fmt.Sprintf("%s", v)
		return s
	}
}
