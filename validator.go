package util

import (
	"bytes"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidatorHelperErrors []string

func (ve ValidatorHelperErrors) Error() string {
	if len(ve) == 0 {
		return ""
	}
	buff := bytes.NewBufferString("")
	for i := 0; i < len(ve); i++ {
		buff.WriteString(ve[i])
		buff.WriteString("\n")
	}
	return strings.TrimSpace(buff.String())
}

func ValidatorHelper(s interface{}, errs error) ValidatorHelperErrors {
	validationErrors, ok := errs.(validator.ValidationErrors)
	if !ok {
		return ValidatorHelperErrors{"errors type error"}
	}
	validatorHelperErrors := make(ValidatorHelperErrors, 0, len(validationErrors))
	for _, err := range validationErrors {
		ns := err.Namespace()
		index := strings.Index(ns, ".")

		msg := findMsg(s, ns[index+1:])
		if msg != "" {
			validatorHelperErrors = append(validatorHelperErrors, msg)
		}
	}
	return validatorHelperErrors
}

func findMsg(src interface{}, name string) (msg string) {
	val := reflect.ValueOf(src)
	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}
	// src 只能是结构体
	// 判断name . [ 的索引
	indexDoc := strings.Index(name, ".")
	indexStr := strings.Index(name, "[")
	// A
	if indexDoc < 0 && indexStr < 0 {
		// 直接看本字段的错误
		if field, ok := val.Type().FieldByName(name); ok {
			return field.Tag.Get("msg")
		}
		// A.B  A.B[0]
	} else if indexDoc > 0 && (indexDoc < indexStr || indexStr < 0) {
		// 取结构体字段，进入下层解析
		field := val.FieldByName(name[:indexDoc])
		return findMsg(field.Interface(), name[indexDoc+1:])
		// A[0].B
	} else if indexStr > 0 && indexDoc > indexStr {
		// 数组或切片内部错误
		field := val.FieldByName(name[:indexStr])
		return findMsg(field.Index(0).Interface(), name[indexDoc+1:])
		// A[0]
	} else {
		if field, ok := val.Type().FieldByName(name[:indexStr]); ok {
			return field.Tag.Get("msg")
		}
	}
	return
}
