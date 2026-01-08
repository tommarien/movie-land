package validator

import (
	"fmt"
	"regexp"
)

var slugRegex = regexp.MustCompile(`^[a-z-]+$`)

type Validator struct {
	errors map[string]string
}

func New() *Validator {
	return &Validator{
		errors: make(map[string]string),
	}
}

func (v *Validator) Required(name, value string) {
	if value == "" {
		v.errors[name] = fmt.Sprintf("%s is required", name)
	}
}

func (v *Validator) Slug(name, value string) {
	if value != "" && !slugRegex.MatchString(value) {
		v.errors[name] = fmt.Sprintf("%s must contain only lowercase letters and hyphens", name)
	}
}

func (v *Validator) MaxLength(name, value string, maxLength int) {
	if value != "" && len(value) > maxLength {
		v.errors[name] = fmt.Sprintf("%s must not exceed %d characters", name, maxLength)
	}
}

func (v *Validator) IsValid() bool {
	return len(v.errors) == 0
}

func (v *Validator) GetErrors() []string {
	fieldErrors := make([]string, 0, len(v.errors))
	for _, message := range v.errors {
		fieldErrors = append(fieldErrors, message)
	}
	return fieldErrors
}
