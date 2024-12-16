package forms

import (
	"net/http"
	"net/url"
)

// Form contains all form information
type Form struct {
	url.Values
	Errors errors
}

// New initialises a form type and passes back a pointer to the newly created form struct.
func New(data url.Values) *Form {
	return &Form{
		Values: data,
		Errors: make(errors),
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
