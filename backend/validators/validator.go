package validators

import (
	"github.com/go-playground/validator/v10"
)

type Valid interface {
	Validate(v *validator.Validate) error
}

var Validate *FinalTestinationValidator

type FinalTestinationValidator struct {
	inner validator.Validate
}

func (v *FinalTestinationValidator) GetValidator() *validator.Validate {
	return &v.inner
}

func (v *FinalTestinationValidator) addValidator(tag string, validator func(validator.FieldLevel) bool) {
	_ = v.inner.RegisterValidation(tag, validator)
}

func (v *FinalTestinationValidator) addAlias(alias, validatorString string) {
	v.inner.RegisterAlias(alias, validatorString)
}

const (
	USERNAME_VALIDATOR = "alphanum,min=3,max=30"
	PASSWORD_VALIDATOR = "min=8,max=30"
)

func init() {
	Validate = &FinalTestinationValidator{inner: *validator.New()}

	Validate.addAlias("testination-username", USERNAME_VALIDATOR)
	Validate.addAlias("testination-password", PASSWORD_VALIDATOR)

	// Note: cannot just do an alias to "testination-username|email" because there is a know issue
	// what will panic as soon as the validator is used.
	// Source: https://github.com/go-playground/validator/issues/766
	// There also seems to be a bug with the notation `(alphanum,max=30,min=3)|email`, so this is
	// the only way to do it.
	Validate.addValidator("testination-credential", func(fl validator.FieldLevel) bool {
		usernameErr := Validate.GetValidator().Var(fl.Field().String(), "testination-username")
		emailErr := Validate.GetValidator().Var(fl.Field().String(), "email")

		// If none of the validators match, then the credential is invalid
		if usernameErr != nil && emailErr != nil {
			return false
		}

		return true
	})
}
