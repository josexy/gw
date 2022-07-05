package validator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type ValidateFunc = func(validator.FieldLevel) bool

type Validator struct {
	v *validator.Validate
}

func newValidator() *Validator {
	va := new(Validator)
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		va.v = v
	} else {
		panic("can not bind validator")
	}
	return va
}

func (v *Validator) register(name string, fn ValidateFunc) {
	_ = v.v.RegisterValidation(name, fn)
}

func InitValidator() {
	v := newValidator()
	// 校验服务名
	v.register("valid_service_name", validServiceName)
	// 校验规则
	v.register("valid_rule", validRule)
	// 校验URL重写
	v.register("valid_url_rewrite", validUrlRewrite)
	// 校验header头
	v.register("valid_header_transfer", validHeaderTransfer)
	// 校验ip+端口列表
	v.register("valid_ipportlist", validIPPortList)
	// 校验ip权重列表
	v.register("valid_weightlist", validWeightList)
}
