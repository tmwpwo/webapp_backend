package renders

import (
	"bytes"
	"go_server/pkg/config"
	"go_server/pkg/models"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/justinas/nosurf"
)

var app *config.AppConfig

// NewTemplates sets the config for template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, req *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(req)
	return td
}

// RenderTmpl renders templates using render html/template
func RenderTmpl(w http.ResponseWriter, req *http.Request, tmpl string, td *models.TemplateData) {

	var templateCache map[string]*template.Template
	if app.UseCache {
		//get the tempate cache from the app config
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}

	//get requested template from cache
	t, ok := templateCache[tmpl]
	if !ok {
		log.Fatal("could not get template form template cache")
	}

	buffer := new(bytes.Buffer)

	td = AddDefaultData(td, req)

	_ = t.Execute(buffer, td)

	//render the template
	_, err := buffer.WriteTo(w)
	if err != nil {
		log.Println(err)
	}

}

func CreateTemplateCache() (map[string]*template.Template, error) {
	// myCache := make(map[string]*template.Template)
	myCache := map[string]*template.Template{}

	// get all of the files named *.page.html from ./templates

	pages, err := filepath.Glob("./templates/*page.html")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		templateSet, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob("./templates/*.layout.html")
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = templateSet
	}
	return myCache, nil
}
