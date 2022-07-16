package validation

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	vLib "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type (
	Validator struct {
		validation *vLib.Validate
		trans      ut.Translator
	}

	Field vLib.FieldLevel
)

func NewValidator() *Validator {
	uni := ut.New(en.New())
	trans, _ := uni.GetTranslator("en")
	v := vLib.New()
	en_translations.RegisterDefaultTranslations(v, trans)

	return &Validator{
		validation: v,
		trans:      trans,
	}
}

func (v *Validator) RegisterValidation(tag string, check func(Field) bool, errMsg string) error {
	err := v.validation.RegisterTranslation(tag, v.trans,
		func(ut ut.Translator) error {
			return ut.Add(tag, errMsg, true)
		},
		func(ut ut.Translator, fe vLib.FieldError) string {
			t, err := ut.T(tag, fe.Field())
			if err != nil {
				// TODO: change this to something more error tolerant
				panic("could not register validation")
			}
			return t
		},
	)
	if err != nil {
		return err
	}
	return v.validation.RegisterValidation(tag, func(fl vLib.FieldLevel) bool {
		return check(fl)
	})
}

func (v *Validator) ValidateStruct(value any) error {
	return v.validation.Struct(value)
}

func (v *Validator) UnpackErrors(e error) []string {
	values, ok := e.(vLib.ValidationErrors)
	if !ok {
		return nil
	}
	errs := make([]string, 0, len(values))
	for _, vv := range values {
		errs = append(errs, vv.Translate(v.trans))
	}
	return errs
}
