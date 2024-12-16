package forms

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/http"
	"net/url"
	"strings"
)

// Form contains all form information
type Form struct {
	url.Values
	Errors errors
}

// Valid returns true if there is no errors, otherwise false.
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New initialises a form type and passes back a pointer to the newly created form struct.
func New(data url.Values) *Form {
	return &Form{
		Values: data,
		Errors: make(errors),
	}
}

// Required checks if form fields (text type) is, when trimmed, empty.
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be empty")
		}
	}
}

// Has checks if form field is in post any not empty.
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.FormValue(field)
	if x == "" {
		return false
	}

	return true
}

// MinLength checks if the form field, has a required length.
func (f *Form) MinLength(field string, length int, r *http.Request) bool {
	x := r.FormValue(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}

	return true
}

// IsEmail checks for valid email address in form field.
func (f *Form) IsEmail(field string) bool {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Please enter a valid email address")
		return false
	}

	return true
}
