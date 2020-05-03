package ui

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/ImageWare/TLSential/certificate"
	"github.com/ImageWare/TLSential/model"
	"github.com/gorilla/mux"
)

// Handler provides an interface for all ui/calls.
type Handler interface {
	Route() *mux.Router
}

type uiHandler struct {
	Title              string
	certificateService certificate.Service
}

func NewHandler(title string, cs certificate.Service) Handler {
	return &uiHandler{title, cs}
}

func (h *uiHandler) Route() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/ui/dashboard", h.Dashboard())
	r.HandleFunc("/ui/cert/{id}", h.Certificate())
	r.HandleFunc("/ui/login", h.Login())
	return r
}

type headTemplate struct {
	Title         string
	CustomCSSFile string
}

type loginTemplate struct {
	Head headTemplate
}

func (h *uiHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseGlob("ui/templates/*.html")
		if err != nil {
			log.Print(err.Error())
		}
		head := headTemplate{"Login", "login.css"}
		p := loginTemplate{head}
		t.ExecuteTemplate(w, "login", p)
	}
}

type dashboardTemplate struct {
	Head              headTemplate
	TotalCerts        int
	TotalRenewedCerts int
	TotalDomains      int
}

// Serve /ui/dashboard page.
func (h *uiHandler) Dashboard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseGlob("ui/templates/*.html")
		if err != nil {
			log.Print(err.Error())
		}
		head := headTemplate{"Dashboard", "dashboard.css"}

		// TODO: Fill out appropriate data for cert, renewed cert, and domain counts.
		d := dashboardTemplate{head, 4, 20, 69}
		t, _ = template.ParseFiles("layout.html", "dashboard.html")
		t.ExecuteTemplate(w, "dashboard", d)
	}
}

type certTemplate struct {
	Head headTemplate
	Cert *model.Certificate
}

// Serve /ui/certificate page.
func (h *uiHandler) Certificate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseGlob("ui/templates/*.html")
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "drats", http.StatusInternalServerError)
			return
		}

		id := mux.Vars(r)["id"]

		cert, err := h.certificateService.Cert(id)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "whoops", http.StatusInternalServerError)
			return
		}

		head := headTemplate{
			fmt.Sprintf("Certificate - %s", cert.CommonName),
			"certificate.css",
		}
		p := certTemplate{
			head,
			cert,
		}
		t.ExecuteTemplate(w, "certificate", p)
	}
}