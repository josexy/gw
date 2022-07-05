package validator

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

func validServiceName(fl validator.FieldLevel) bool {
	matched, _ := regexp.Match("^[a-zA-Z0-9_]{6,128}$", []byte(fl.Field().String()))
	return matched
}

func validRule(fl validator.FieldLevel) bool {
	matched, _ := regexp.Match("^\\S+$", []byte(fl.Field().String()))
	return matched
}

// ^/test_http_service/abb/(.*) /test_http_service/bba/$1
func validUrlRewrite(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if len(value) == 0 {
		return true
	}
	for _, ms := range strings.Split(value, ",") {
		if len(strings.Split(ms, " ")) != 2 {
			return false
		}
	}
	return true
}

// add header_name header_value
func validHeaderTransfer(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if len(value) == 0 {
		return true
	}
	for _, ms := range strings.Split(value, ",") {
		if len(strings.Split(ms, " ")) != 3 {
			return false
		}
	}
	return true
}

func validIPPortList(fl validator.FieldLevel) bool {
	for _, ms := range strings.Split(fl.Field().String(), ",") {
		if match, _ := regexp.Match(`^\S+\:\d+$`, []byte(ms)); !match {
			return false
		}
	}
	return true
}

func validWeightList(fl validator.FieldLevel) bool {
	for _, ms := range strings.Split(fl.Field().String(), ",") {
		if match, _ := regexp.Match(`^\d+$`, []byte(ms)); !match {
			return false
		}
	}
	return true
}
