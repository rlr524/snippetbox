package forms

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

// The Form struct anonymously embeds a url.Values object (to hold the form data) and an
// FormErrors field (of type Errors) to hold any validation errors for the form data.
type Form struct {
	url.Values
	FormErrors Errors
}

// The New function initializes a custom Form struct, takes the form data as a param and returns a pointer to
// the Form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		Errors(map[string][]string{}),
	}
}

// The Required method checks that specific fields in the form data are present and not blank. If any fields fail
// this check, add the appropriate message to the form errors. Not this is a function pointer receiver as it is
// modifying the Form.
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.FormErrors.Add(field, "This field cannot be blank")
		}
	}
}

// The MaxLength method checks that a specific field in the form contains less than a maximum number of characters.
// If the check fails then add the appropriate message to the form errors.
func (f *Form) MaxLength(field string, max int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > max {
		f.FormErrors.Add(field, fmt.Sprintf("This field is too long, (maximum is %d characters)", max))
	}
}

// The PermittedValues method checks that a specific field in the form matches one of a set of specific permitted
// values. If the check fails, then add the appropriate message to the form errors.
func (f *Form) PermittedValues(field string, permitted ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	for _, p := range permitted {
		if value == p {
			return
		}
	}
	f.FormErrors.Add(field, "This field is invalid")
}

// The Valid method returns true if there are no errors
func (f *Form) Valid() bool {
	if len(f.FormErrors) == 0 {
		return true
	}
	return false
}

func main() {

}
