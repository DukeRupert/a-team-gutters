package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
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

// postmarkEmail represents the Postmark API email payload.
type postmarkEmail struct {
	From     string `json:"From"`
	To       string `json:"To"`
	Subject  string `json:"Subject"`
	TextBody string `json:"TextBody"`
	ReplyTo  string `json:"ReplyTo"`
}

// sendPostmarkEmail sends an email via the Postmark API.
// Returns nil if Postmark is not configured (missing env vars).
func sendPostmarkEmail(to, from, replyTo, subject, body string) error {
	token := os.Getenv("POSTMARK_SERVER_TOKEN")
	if token == "" {
		log.Printf("POSTMARK_SERVER_TOKEN not set — skipping email send")
		return nil
	}

	payload := postmarkEmail{
		From:     from,
		To:       to,
		Subject:  subject,
		TextBody: body,
		ReplyTo:  replyTo,
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal email payload: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.postmarkapp.com/email", bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Postmark-Server-Token", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("postmark request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("postmark returned %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func handleContactSubmit(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Honeypot check — if filled, silently accept (don't reveal to bots)
		if r.FormValue("website") != "" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			tmpl.ExecuteTemplate(w, "contact_success", nil)
			return
		}

		// Collect and trim fields
		name := strings.TrimSpace(r.FormValue("name"))
		phone := strings.TrimSpace(r.FormValue("phone"))
		email := strings.TrimSpace(r.FormValue("email"))
		address := strings.TrimSpace(r.FormValue("address"))
		service := strings.TrimSpace(r.FormValue("service"))
		message := strings.TrimSpace(r.FormValue("message"))

		// Server-side validation (defense in depth — Alpine validates client-side first)
		if name == "" || phone == "" || email == "" || !strings.Contains(email, "@") || address == "" || service == "" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Header().Set("HX-Retarget", "#form-errors")
			w.Header().Set("HX-Reswap", "innerHTML")
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprint(w, `<div class="form-error-banner" role="alert">Please fill out all required fields and try again.</div>`)
			return
		}

		// Log the submission
		timestamp := time.Now().Format("2006-01-02 15:04:05 MST")
		log.Printf("=== New Estimate Request ===")
		log.Printf("Name:    %s", name)
		log.Printf("Phone:   %s", phone)
		log.Printf("Email:   %s", email)
		log.Printf("Address: %s", address)
		log.Printf("Service: %s", service)
		log.Printf("Message: %s", message)
		log.Printf("Time:    %s", timestamp)
		log.Printf("============================")

		// Send email via Postmark
		postmarkTo := os.Getenv("POSTMARK_TO")
		postmarkFrom := os.Getenv("POSTMARK_FROM")

		if postmarkTo != "" && postmarkFrom != "" {
			emailBody := fmt.Sprintf(`New Estimate Request — A-Team Gutters

Name:    %s
Phone:   %s
Email:   %s
Address: %s
Service: %s
Message: %s

Submitted: %s`, name, phone, email, address, service, message, timestamp)

			subject := fmt.Sprintf("New Estimate Request — %s (%s)", name, service)

			if err := sendPostmarkEmail(postmarkTo, postmarkFrom, email, subject, emailBody); err != nil {
				log.Printf("ERROR sending email: %v", err)
				// Don't fail the submission — the log has the data
			}
		}

		// Return success partial
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.ExecuteTemplate(w, "contact_success", nil)
	}
}

func main() {
	// Load each page template with base layout + partials
	pages := map[string]*template.Template{
		"home":                 loadPage("home.html"),
		"about":               loadPage("about.html"),
		"contact":             loadPage("contact.html"),
		"faq":                 loadPage("faq.html"),
		"gutter-installation": loadPage("gutter-installation.html"),
		"gutter-cleaning":     loadPage("gutter-cleaning.html"),
		"gutter-repair":       loadPage("gutter-repair.html"),
		"gutter-guards":       loadPage("gutter-guards.html"),
		"fascia-soffit-repair": loadPage("fascia-soffit-repair.html"),
		"gallery":              loadPage("gallery.html"),
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
	mux.HandleFunc("GET /gallery/", servePage(pages["gallery"]))
	mux.HandleFunc("GET /about/", servePage(pages["about"]))
	mux.HandleFunc("GET /contact/", servePage(pages["contact"]))
	mux.HandleFunc("POST /contact/", handleContactSubmit(pages["contact"]))
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
