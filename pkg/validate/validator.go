package validate

import (
	"fmt"
	"github.com/zktnotify/zktnotify/viewmodel"
	"gopkg.in/go-playground/validator.v9"
)

var validate = validator.New()

func VerifyStruct(param interface{}) error {
	err := validate.Struct(param)
	if err != nil {
		//只有在您的代码可以生成时才需要进行此检查
		//验证的无效值，例如与nil的接口
		//大多数包括我自己的值通常不会有这样的代码。
		if _, ok := err.(*validator.InvalidValidationError); ok {
			//TODO log
			return err
		}
		return err
	}
	return nil
}

func HandleVerifyErrorResult(errors validator.ValidationErrors) (result []*viewmodel.VerifyResult) {
	for _, err := range errors {
		result = append(result, &viewmodel.VerifyResult{
			Namespace:       err.Namespace(),
			Field:           err.Field(),
			StructNamespace: err.StructNamespace(),
			StructField:     err.StructField(),
			Tag:             err.Tag(),
			ActualTag:       err.ActualTag(),
			Kind:            err.Kind().String(),
			Type:            err.Type().String(),
			Value:           fmt.Sprintf("%v", err.Value()),
			Param:           err.Param(),
		})
	}
	return
}
