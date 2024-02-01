package validator

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var (
	emailRegex = regexp.MustCompile("^(?:(?:(?:[a-zA-Z]|\\d|[!#$%&'*+\\-/=?^_`{|}~]" +
		"|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#$%&'*+\\-/=?\\^_`{|}~]" +
		"|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)" +
		"|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))" +
		"?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]" +
		"|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])" +
		"|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]" +
		"|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))" +
		"*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22)))@(?:(?:(?:[a-zA-Z]" +
		"|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]" +
		"|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])" +
		"(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])" +
		"*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]" +
		"|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]" +
		"|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])" +
		"(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])" +
		"*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$")
)

func Url() Rule {
	return Rule{
		Name: "url",
		ErrorMessage: func(r Rule) string {
			return "not a valid url"
		},
		Validator: func(r Rule) (bool, ErrorList) {
			v, ok := r.GivenValue.(string)
			if !ok {
				return false, nil
			}

			if strings.HasPrefix(v, "http://") == false && strings.HasPrefix(v, "https://") == false {
				v = fmt.Sprintf("//%s", v)
			}

			u, err := url.Parse(v)
			if err != nil {
				return false, nil
			}

			return u.Host != "", nil
		},
	}
}

func Email() Rule {
	return Rule{
		Name: "email",
		ErrorMessage: func(r Rule) string {
			return "email address is invalid"
		},
		Validator: func(r Rule) (bool, ErrorList) {
			email, ok := r.GivenValue.(string)
			if !ok {
				return false, nil
			}
			return emailRegex.MatchString(email), nil
		},
	}
}

func MaxLength(n int) RuleFunc {
	return func() Rule {
		return Rule{
			Name:          "maxLength",
			ExpectedValue: n,
			Validator: func(r Rule) (bool, ErrorList) {
				str, ok := r.GivenValue.(string)
				if !ok {
					return false, nil
				}
				return len(str) <= n, nil
			},
			ErrorMessage: func(r Rule) string {
				return fmt.Sprintf("%s should be maximum %d characters long", r.Target, n)
			},
		}
	}
}

func MinLength(n int) RuleFunc {
	return func() Rule {
		return Rule{
			Name:          "minLength",
			ExpectedValue: n,
			Validator: func(r Rule) (bool, ErrorList) {
				str, ok := r.GivenValue.(string)
				if !ok {
					return false, nil
				}
				return len(str) >= n, nil
			},
			ErrorMessage: func(r Rule) string {
				return fmt.Sprintf("%s should be at least %d characters long", r.Target, n)
			},
		}
	}
}
