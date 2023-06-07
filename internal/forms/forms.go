package forms

import (
	"net/http"
	"net/url"
	"strings"
)

// creates custom form struct, embebeds a url.Values
type Form struct {
	url.Values
	Errors errors
}

// Valid returns true if there are no errors, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// checks if form field is in post and not empty
func (f *Form) Has(field string, req *http.Request) bool {
	if req.Form.Get(field) == "" {
		f.Errors.Add(field, "this field cannnot be blank")
		return false
	}
	return true
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "this field cannot be blank")
		}
	}
}
