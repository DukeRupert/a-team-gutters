package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// loadPage parses the base layout, partials, and the given page template.
func loadPage(page string) *template.Template {
	return template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/partials/nav.html",
		"templates/partials/footer.html",
		filepath.Join("templates/pages", page),
	))
}

func servePage(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.ExecuteTemplate(w, "base", nil); err != nil {
			log.Printf("template error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func main() {
	// Load each page template with base layout + partials
	pages := map[string]*template.Template{
		"home":                loadPage("home.html"),
		"about":              loadPage("about.html"),
		"contact":            loadPage("contact.html"),
		"faq":                loadPage("faq.html"),
		"gutter-installation": loadPage("gutter-installation.html"),
		"gutter-cleaning":    loadPage("gutter-cleaning.html"),
		"gutter-repair":      loadPage("gutter-repair.html"),
		"gutter-guards":      loadPage("gutter-guards.html"),
		"fascia-soffit-repair": loadPage("fascia-soffit-repair.html"),
	}

	mux := http.NewServeMux()

	// Home — exact match only
	mux.HandleFunc("GET /{$}", servePage(pages["home"]))

	// Services
	mux.HandleFunc("GET /services/gutter-installation/", servePage(pages["gutter-installation"]))
	mux.HandleFunc("GET /services/gutter-cleaning/", servePage(pages["gutter-cleaning"]))
	mux.HandleFunc("GET /services/gutter-repair/", servePage(pages["gutter-repair"]))
	mux.HandleFunc("GET /services/gutter-guards/", servePage(pages["gutter-guards"]))
	mux.HandleFunc("GET /services/fascia-soffit-repair/", servePage(pages["fascia-soffit-repair"]))

	// Core pages
	mux.HandleFunc("GET /about/", servePage(pages["about"]))
	mux.HandleFunc("GET /contact/", servePage(pages["contact"]))
	mux.HandleFunc("GET /faq/", servePage(pages["faq"]))

	// Static files
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
