package models

import "go_server/internal/forms"

type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string //cross-site request forgery token
	Flash     string
	Warning   string
	Error     string
	Form      *forms.Form
}
