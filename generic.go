package validator

import (
	"fmt"
)

func Required() Rule {
	return Rule{
		Name: "required",
		ErrorMessage: func(set Rule) string {
			return fmt.Sprintf("%s is a required field", set.Target)
		},
		Validator: func(rule Rule) (bool, ErrorList) {
			str, ok := rule.GivenValue.(string)
			if !ok {
				return false, nil
			}
			return len(str) > 0, nil
		},
	}
}
