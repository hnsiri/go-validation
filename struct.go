package validator

import (
	"fmt"
	"reflect"
)

func Struct(data any, fields Fields) RuleFunc {
	return func() Rule {
		return Rule{
			Name: "struct",
			ErrorMessage: func(r Rule) string {
				return fmt.Sprintf("%s is not a struct", r.Target)
			},
			Validator: func(r Rule) (bool, ErrorList) {
				val := reflect.ValueOf(r.GivenValue)
				if val.Kind() != reflect.Struct {
					return false, nil
				}

				sv := New(data, fields)
				ok := sv.Validate()
				return ok, sv.Errors()
			},
		}
	}
}

func ArrayOfStructs(data any, fields Fields) RuleFunc {
	return func() Rule {
		return Rule{
			Name: "arrayOfStructs",
			ErrorMessage: func(r Rule) string {
				return fmt.Sprintf("%s is not an array of structs", r.Target)
			},
			Validator: func(r Rule) (bool, ErrorList) {
				val := reflect.ValueOf(r.GivenValue)
				switch val.Kind() {
				case reflect.Slice, reflect.Array:
					if val.Index(0).Kind() != reflect.Struct {
						return false, nil
					}

					errs := make(ErrorList)
					var res bool
					for i, s := range r.GivenValue.([]any) {
						sv := New(s, fields)
						ok := sv.Validate()
						if !ok {
							res = false
							for k, e := range sv.Errors() {
								errs[fmt.Sprintf("%d.%s", i, k)] = e
							}
						}
					}
					return res, errs

				default:
					return false, nil
				}
			},
		}
	}
}
