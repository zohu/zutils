package zutils

import (
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhtrans "github.com/go-playground/validator/v10/translations/zh"
)

var ValiTrans ut.Translator
var ValiObj *validator.Validate

func Transfer() {
	ValiObj = validator.New()
	english := en.New()
	chinese := zh.New()
	uni := ut.New(chinese, chinese, english)
	ValiTrans, _ = uni.GetTranslator("zh")
	_ = zhtrans.RegisterDefaultTranslations(ValiObj, ValiTrans)
}
