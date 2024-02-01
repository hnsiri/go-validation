package validator

import (
	"errors"
	"fmt"
	"reflect"
	"unicode"

	"github.com/iancoleman/strcase"
)

const (
	SnakeCase = "snake"
	CamelCase = "camel"
)

var (
	Lang       = "en"
	NamingCase = SnakeCase
)

type RuleFunc func() Rule

type ErrorList map[string][]error

type Rule struct {
	Name          string
	Target        any
	ExpectedValue any
	GivenValue    any
	ErrorMessage  func(Rule) string
	Validator     func(Rule) (bool, ErrorList)
}

type Fields map[string][]Rule

func Rules(rules ...RuleFunc) []Rule {
	ruleSets := make([]Rule, len(rules))
	for i := 0; i < len(ruleSets); i++ {
		ruleSets[i] = rules[i]()
	}
	return ruleSets
}

type Validator struct {
	data      any
	fields    Fields
	validated bool
	errors    ErrorList
}

func New(data any, fields Fields) *Validator {
	errs := make(ErrorList, len(fields))
	return &Validator{
		fields: fields,
		data:   data,
		errors: errs,
	}
}

func (v *Validator) Validate() bool {
	ok := true
	for fieldName, rules := range v.fields {
		if !unicode.IsUpper(rune(fieldName[0])) {
			continue
		}
		fieldValue := getFieldValueByName(v.data, fieldName)
		for _, rule := range rules {
			rule.GivenValue = fieldValue
			rule.Target = fieldName
			isValid, errs := rule.Validator(rule)
			if !isValid {
				v.addError(fieldName, errors.New(rule.ErrorMessage(rule)))
				ok = false
				if errs != nil {
					for nestedFieldName, nestedErrors := range errs {
						for _, err := range nestedErrors {
							v.addError(fmt.Sprintf("%s.%s", fieldName, nestedFieldName), err)
						}
					}
				}
			}
		}
	}

	v.validated = true
	return ok
}

func (v *Validator) Errors() map[string][]error {
	return v.errors
}

func (v *Validator) IsValid() bool {
	return len(v.errors) == 0 && v.validated
}

func (v *Validator) addError(fieldName string, err error) {
	fieldName = convertToCase(fieldName)
	v.errors[fieldName] = append(v.errors[fieldName], err)
}

func getFieldValueByName(v any, name string) any {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}
	fieldVal := val.FieldByName(name)
	if !fieldVal.IsValid() {
		return nil
	}
	return fieldVal.Interface()
}

func convertToCase(f string) string {
	switch NamingCase {
	case SnakeCase:
		return strcase.ToSnakeWithIgnore(f, ".")
	case CamelCase:
		return strcase.ToLowerCamel(f)
	default:
		return f
	}
}
