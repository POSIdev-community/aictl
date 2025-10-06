package regexfilter

import "regexp"

type RegexFilter struct {
	value regexp.Regexp
	empty bool

	createdFromConstructor bool
}

func NewRegexFilter(value string) (RegexFilter, error) {
	if value == "" {
		return RegexFilter{
			empty:                  true,
			createdFromConstructor: true,
		}, nil
	}

	regex, err := regexp.Compile(value)
	if err != nil {
		return RegexFilter{}, err
	}

	return RegexFilter{*regex, false, true}, nil
}

func (r RegexFilter) Empty() bool {
	return r.empty
}

func (r RegexFilter) Execute(name string) (matched bool) {
	return r.value.MatchString(name)
}
