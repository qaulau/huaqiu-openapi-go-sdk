package openapi

import (
	"fmt"
	"strconv"
)

const (
	_ = iota
	StringType
	BytesType
	Int64Type
	Int32Type
	IntType
	Float64Type
	Float32Type
	BoolType
	ValType
)

const (
	StringFlag int32 = 1
	IntFlag int32 = 2
	FloatFlag int32 = 4
	BoolFlag int32 = 8
)

type object interface {}

type any struct {
	object
	otype int32
	flag int32
	text string
	ival int64
	fval float64
	bval bool
}

func NewAny(val object) any {
	a := any{}
	a.object = val
	return a
}

// 分析数据
func (this *any) analysis(flag int32) *any{
	if this.flag == 0 {
		switch this.object.(type) {
		case string:
			this.text, _ = this.object.(string)
			this.otype = StringType
			this.flag |= StringFlag
		case int64:
			this.ival, _ = this.object.(int64)
			this.text = strconv.FormatInt(this.ival, 10)
			this.otype = Int64Type
			this.flag |= IntFlag
		case int32:
			ival, _ := this.object.(int32)
			this.ival = int64(ival)
			this.text = strconv.FormatInt(this.ival, 10)
			this.otype = Int32Type
			this.flag |= IntFlag
		case int:
			ival := this.object.(int)
			this.ival = int64(ival)
			this.text = strconv.FormatInt(this.ival, 10)
			this.otype = IntType
			this.flag |= IntFlag
		case float64:
			this.fval, _ = this.object.(float64)
			this.text = strconv.FormatFloat(this.fval, 'E', -1, 64)
			this.otype = Float64Type
			this.flag |= FloatFlag
		case float32:
			fval := this.object.(float32)
			this.fval = float64(fval)
			this.text = strconv.FormatFloat(this.fval, 'E', -1, 32)
			this.otype = Float32Type
			this.flag |= FloatFlag
		case []byte:
			bval := this.object.([]byte)
			this.text = string(bval)
			this.otype = BytesType
		case bool:
			this.bval, _ = this.object.(bool)
			this.text = strconv.FormatBool(this.bval)
			this.otype = BoolType
			this.flag |= BoolFlag
		default:
			this.text = fmt.Sprintf("%v", this.object)
			this.otype = ValType
		}
		this.flag |= StringFlag
	}
	if flag == IntFlag && this.flag & IntFlag == 0 {
		this.ival, _ = strconv.ParseInt(this.text, 10, 64)
	}
	if flag == FloatFlag && this.flag & FloatFlag == 0 {
		this.fval, _ = strconv.ParseFloat(this.text, 64)
	}
	if flag == BoolFlag && this.flag & BoolFlag == 0 {
		this.bval, _ = strconv.ParseBool(this.text)
	}
	return this
}

// 转字符串
func (this *any) String() string {
	this.analysis(StringFlag)
	return this.text
}

// 转字节
func (this *any) Bytes() []byte {
	return []byte(this.text)
}

// 转 bool
func (this *any) Bool() bool {
	this.analysis(BoolFlag)
	return this.bval
}

// 转 int64
func (this *any) Int64() int64 {
	this.analysis(IntFlag)
	return this.ival
}

func (this *any) Int32() int32 {
	return int32(this.Int64())
}

func (this *any) Int() int {
	return int(this.Int64())
}

// 转 float64
func (this *any) Float64() float64 {
	this.analysis(FloatFlag)
	return this.fval
}

// 转 float32
func (this *any) Float32() float32 {
	return float32(this.Float64())
}

// 原始值
func (this *any) Value() interface{} {
	return this.object
}