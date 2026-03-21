package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	tmpl := template.Must(template.ParseGlob("templates/*.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		tmpl.ExecuteTemplate(w, "selector.html", nil)
	})

	http.HandleFunc("/forest-rain", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "forest-rain.html", nil)
	})

	http.HandleFunc("/rain-slate", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "rain-slate.html", nil)
	})

	http.HandleFunc("/classic", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "classic.html", nil)
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
