package validationutil

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type ValidationErr struct {
	fieldErrs map[string]string
	msg       string
}

func NewValidationErr(msg string, fieldErrs map[string]string) *ValidationErr {
	return &ValidationErr{
		fieldErrs: fieldErrs,
		msg:       msg,
	}
}

func NewSingleValidationErr(fieldName, msg string) error {
	return &ValidationErr{
		fieldErrs: map[string]string{
			fieldName: msg,
		},
		msg: msg,
	}
}

type Validation interface {
	ValidateStruct(structPtr interface{}, fields ...*validation.FieldRules) error
}

func (v *ValidationErr) Error() string {
	return v.msg
}

func (v *ValidationErr) GetFieldErrs() map[string]string {
	return v.fieldErrs
}

func ValidateStruct(structPtr interface{}, fields ...*validation.FieldRules) error {
	err := validation.ValidateStruct(structPtr, fields...)
	if e, ok := err.(validation.InternalError); ok {
		// an internal error happened
		panic(e.InternalError())
	} else if errs, ok := err.(validation.Errors); ok {
		fieldErrs := map[string]string{}
		for fieldName, fieldErr := range errs {
			fieldErrs[fieldName] = fieldErr.Error()
		}

		return &ValidationErr{
			msg:       err.Error(),
			fieldErrs: fieldErrs,
		}
	}

	return err
}
