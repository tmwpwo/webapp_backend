package handlers

import (
	"go_server/internal/forms"
	"go_server/pkg/config"
	"go_server/pkg/driver"
	"go_server/pkg/models"
	"go_server/pkg/renders"
	"go_server/pkg/repository"
	"go_server/pkg/repository/dbrepo"
	"log"
	"net/http"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repositiory is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHAndlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, req *http.Request) {

	remoteIP := req.RemoteAddr
	m.App.Session.Put(req.Context(), "remote_ip", remoteIP)

	renders.RenderTmpl(w, req, "home.page.html", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, req *http.Request) {

	stringMap := make(map[string]string)
	stringMap["test"] = "hellooooooooo"

	remoteIP := m.App.Session.GetString(req.Context(), "remote_ip")

	stringMap["remote_ip"] = remoteIP

	renders.RenderTmpl(w, req, "about.page.html", &models.TemplateData{StringMap: stringMap})
}

func (m *Repository) ShowLogin(w http.ResponseWriter, req *http.Request) {
	renders.RenderTmpl(w, req, "login.page.html", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// logs user in
func (m *Repository) PostShowLogin(w http.ResponseWriter, req *http.Request) {
	_ = m.App.Session.RenewToken(req.Context())

	err := req.ParseForm()
	if err != nil {
		log.Println(err, "XD")

	}

	email := req.Form.Get("email")
	password := req.Form.Get("password")

	form := forms.New(req.PostForm)
	form.Required("email", "password")
	if !form.Valid() {
		//takes user back to page
		renders.RenderTmpl(w, req, "login.page.html", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, lol, err := m.DB.Authentication(email, password)
	log.Println(id, lol)
	if err != nil {
		log.Println(err, "XDDD")
		m.App.Session.Put(req.Context(), "error", "invalid credentials")
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(req.Context(), "user_id", id)
	m.App.Session.Put(req.Context(), "flash", "logged successfully")
	http.Redirect(w, req, "/home", http.StatusSeeOther)

}
