package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
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

// verifyTurnstile validates a Cloudflare Turnstile token server-side.
// Returns true if verification succeeds or if Turnstile is not configured.
func verifyTurnstile(token, remoteIP string) bool {
	secret := os.Getenv("TURNSTILE_SECRET_KEY")
	if secret == "" {
		log.Printf("TURNSTILE_SECRET_KEY not set — skipping verification")
		return true
	}

	form := fmt.Sprintf("secret=%s&response=%s&remoteip=%s", secret, token, remoteIP)
	resp, err := http.Post(
		"https://challenges.cloudflare.com/turnstile/v0/siteverify",
		"application/x-www-form-urlencoded",
		strings.NewReader(form),
	)
	if err != nil {
		log.Printf("turnstile verification request failed: %v", err)
		return false
	}
	defer resp.Body.Close()

	var result struct {
		Success bool `json:"success"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("turnstile response decode error: %v", err)
		return false
	}

	if !result.Success {
		log.Printf("turnstile verification failed for IP %s", remoteIP)
	}
	return result.Success
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

// pageData holds template data for pages that need dynamic values.
type pageData struct {
	TurnstileSiteKey string
}

// nearbyArea represents a linked service area for cross-navigation.
type nearbyArea struct {
	City string
	Slug string
}

// serviceAreaData holds template data for service area pages.
type serviceAreaData struct {
	City            string
	Slug            string
	MetaTitle       string
	MetaDescription string
	Intro           string
	Context         string
	HeroImage       string
	NearbyAreas     []nearbyArea
}

func serveServiceArea(tmpl *template.Template, data serviceAreaData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
			log.Printf("template error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// sitemapURL represents a single <url> entry in the sitemap.
type sitemapURL struct {
	XMLName    xml.Name `xml:"url"`
	Loc        string   `xml:"loc"`
	ChangeFreq string   `xml:"changefreq,omitempty"`
	Priority   string   `xml:"priority,omitempty"`
}

// sitemapIndex is the root <urlset> element.
type sitemapIndex struct {
	XMLName xml.Name     `xml:"urlset"`
	XMLNS   string       `xml:"xmlns,attr"`
	URLs    []sitemapURL `xml:"url"`
}

func serveSitemap(baseURL string, serviceAreas []serviceAreaData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urls := []sitemapURL{
			{Loc: baseURL + "/", ChangeFreq: "weekly", Priority: "1.0"},
			{Loc: baseURL + "/services/gutter-installation/", ChangeFreq: "monthly", Priority: "0.9"},
			{Loc: baseURL + "/services/gutter-cleaning/", ChangeFreq: "monthly", Priority: "0.9"},
			{Loc: baseURL + "/services/gutter-repair/", ChangeFreq: "monthly", Priority: "0.9"},
			{Loc: baseURL + "/services/gutter-guards/", ChangeFreq: "monthly", Priority: "0.9"},
			{Loc: baseURL + "/services/fascia-soffit-repair/", ChangeFreq: "monthly", Priority: "0.9"},
			{Loc: baseURL + "/gallery/", ChangeFreq: "monthly", Priority: "0.7"},
			{Loc: baseURL + "/about/", ChangeFreq: "monthly", Priority: "0.7"},
			{Loc: baseURL + "/contact/", ChangeFreq: "monthly", Priority: "0.8"},
			{Loc: baseURL + "/faq/", ChangeFreq: "monthly", Priority: "0.7"},
		}

		for _, area := range serviceAreas {
			urls = append(urls, sitemapURL{
				Loc:        baseURL + "/service-areas/" + area.Slug + "/",
				ChangeFreq: "monthly",
				Priority:   "0.8",
			})
		}

		sitemap := sitemapIndex{
			XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
			URLs:  urls,
		}

		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		w.Write([]byte(xml.Header))
		enc := xml.NewEncoder(w)
		enc.Indent("", "  ")
		enc.Encode(sitemap)
	}
}

func serveRobotsTxt(baseURL string) http.HandlerFunc {
	body := "User-agent: *\nAllow: /\n\nSitemap: " + baseURL + "/sitemap.xml\n"
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(body))
	}
}

func serveContact(tmpl *template.Template) http.HandlerFunc {
	siteKey := os.Getenv("TURNSTILE_SITE_KEY")
	return func(w http.ResponseWriter, r *http.Request) {
		data := pageData{TurnstileSiteKey: siteKey}
		if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
			log.Printf("template error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
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

		// Turnstile verification
		turnstileToken := r.FormValue("cf-turnstile-response")
		remoteIP := r.RemoteAddr
		if idx := strings.LastIndex(remoteIP, ":"); idx != -1 {
			remoteIP = remoteIP[:idx]
		}
		if !verifyTurnstile(turnstileToken, remoteIP) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Header().Set("HX-Retarget", "#form-errors")
			w.Header().Set("HX-Reswap", "innerHTML")
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprint(w, `<div class="form-error-banner" role="alert">Verification failed — please try again.</div>`)
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
	mux.HandleFunc("GET /contact/", serveContact(pages["contact"]))
	mux.HandleFunc("POST /contact/", handleContactSubmit(pages["contact"]))
	mux.HandleFunc("GET /faq/", servePage(pages["faq"]))

	// Service area pages
	serviceAreaTmpl := loadPage("service-area.html")
	serviceAreas := []serviceAreaData{
		{
			City:            "Bonney Lake",
			Slug:            "bonney-lake",
			MetaTitle:       "Gutter Installation & Repair in Bonney Lake, WA | A-Team Gutters",
			MetaDescription: "A-Team Gutters is based in Bonney Lake, WA — seamless gutter installation, cleaning, repair, and screen systems for Pierce County homes. Licensed contractor TEAMGGL760KN. Free estimates.",
			Intro:           "A-Team Gutters is based right here in Bonney Lake. We know the neighborhoods, the tree coverage, the rainfall patterns along Lake Tapps, and the roofline styles common to homes throughout the 98391 zip code. When you call A-Team, you're calling a neighbor.",
			Context:         "Bonney Lake sits at the foothills of the Cascades, where mature Douglas fir and cedar canopy is dense and needle drop is year-round. Most homes in the area deal with conifer debris volumes that standard gutter systems and generic guard products aren't designed for. We've been working in Bonney Lake long enough to know what holds up here and what doesn't.",
			HeroImage:       "hero-gutters-home.jpeg",
			NearbyAreas:     []nearbyArea{{City: "Sumner", Slug: "sumner"}, {City: "Buckley", Slug: "buckley"}, {City: "Puyallup", Slug: "puyallup"}},
		},
		{
			City:            "Sumner",
			Slug:            "sumner",
			MetaTitle:       "Gutter Installation & Repair in Sumner, WA | A-Team Gutters",
			MetaDescription: "Professional gutter installation, cleaning, and repair in Sumner, WA. A-Team Gutters serves the Sumner area with seamless aluminum systems built for PNW weather. Free estimates.",
			Intro:           "A-Team Gutters serves Sumner and the surrounding valley communities from our base in neighboring Bonney Lake. Whether you're in the older neighborhoods near downtown Sumner or the newer developments along the valley edge, we provide free on-site estimates and same-standard installation on every job.",
			Context:         "Sumner sits at the confluence of the Puyallup and White Rivers, where the valley transitions to the foothills. The area's mix of mature deciduous trees and conifers means gutters here deal with both heavy leaf loads in fall and needle buildup year-round. Valley-floor properties can also see faster water runoff during heavy rain events — downspout sizing and placement matters here.",
			HeroImage:       "hero-rain-gutters.jpeg",
			NearbyAreas:     []nearbyArea{{City: "Bonney Lake", Slug: "bonney-lake"}, {City: "Puyallup", Slug: "puyallup"}, {City: "Auburn", Slug: "auburn"}},
		},
		{
			City:            "Puyallup",
			Slug:            "puyallup",
			MetaTitle:       "Gutter Installation & Repair in Puyallup, WA | A-Team Gutters",
			MetaDescription: "Seamless gutter installation, cleaning, repair, and screens in Puyallup, WA. A-Team Gutters serves Puyallup and Pierce County. Licensed, insured, free estimates.",
			Intro:           "A-Team Gutters serves Puyallup and the surrounding South Hill communities with the same standard of installation we bring to every job across Pierce County. Free estimates, on-site forming, no subcontractors.",
			Context:         "Puyallup is one of the largest communities in our service area, and the variation in housing stock is significant — from valley-floor homes near the fairgrounds to the elevated neighborhoods of South Hill with steeper rooflines and heavier tree coverage. We assess each property individually. A South Hill home with a 9/12 pitch and mature firs overhead needs a different approach than a flat-lot ranch in the valley.",
			HeroImage:       "hero-green-house.jpg",
			NearbyAreas:     []nearbyArea{{City: "Sumner", Slug: "sumner"}, {City: "Bonney Lake", Slug: "bonney-lake"}, {City: "Auburn", Slug: "auburn"}},
		},
		{
			City:            "Auburn",
			Slug:            "auburn",
			MetaTitle:       "Gutter Installation & Repair in Auburn, WA | A-Team Gutters",
			MetaDescription: "Gutter installation and repair in Auburn, WA. A-Team Gutters serves Auburn and southern King County with seamless systems built for Pacific Northwest conditions. Free estimates.",
			Intro:           "A-Team Gutters serves Auburn and the southern King County communities along the valley corridor. From the older neighborhoods near downtown Auburn to the developments along Highway 18, we bring the same seamless installation standard to every job.",
			Context:         "Auburn sits at the boundary of Pierce and King Counties — a location that puts it within easy reach of our Bonney Lake base. The Green River Valley's mix of agricultural land and residential development means properties here often deal with a combination of valley fog, heavy seasonal rainfall, and the maintenance demands of mature trees planted decades ago.",
			HeroImage:       "hero-copper-gutters.jpg",
			NearbyAreas:     []nearbyArea{{City: "Sumner", Slug: "sumner"}, {City: "Puyallup", Slug: "puyallup"}, {City: "Black Diamond", Slug: "black-diamond"}},
		},
		{
			City:            "Enumclaw",
			Slug:            "enumclaw",
			MetaTitle:       "Gutter Installation & Repair in Enumclaw, WA | A-Team Gutters",
			MetaDescription: "Gutter installation, repair, and cleaning in Enumclaw, WA. A-Team Gutters serves the Enumclaw foothills — licensed contractor with 30+ years experience. Free estimates.",
			Intro:           "A-Team Gutters serves Enumclaw and the surrounding foothills communities at the base of the Cascades. If you're in Enumclaw, Buckley, or the communities between them and the mountains, you're in our primary service area.",
			Context:         "Enumclaw sits higher than most of our service area — close to 700 feet elevation — and experiences conditions the valley floor doesn't. Freeze events are more frequent, ice loading on gutters is a real seasonal concern, and the proximity to the Cascades means wind and precipitation events hit harder. We account for these conditions in how we hang and size systems for Enumclaw properties.",
			HeroImage:       "hero-gutter-install.jpeg",
			NearbyAreas:     []nearbyArea{{City: "Buckley", Slug: "buckley"}, {City: "Bonney Lake", Slug: "bonney-lake"}, {City: "Black Diamond", Slug: "black-diamond"}},
		},
		{
			City:            "Buckley",
			Slug:            "buckley",
			MetaTitle:       "Gutter Installation & Repair in Buckley, WA | A-Team Gutters",
			MetaDescription: "Gutter contractor serving Buckley, WA. Seamless installation, repair, cleaning, and screen systems for Pierce County foothills homes. A-Team Gutters. Free estimates.",
			Intro:           "A-Team Gutters serves Buckley and the Upper White River valley communities. Buckley is one of the closer foothills towns to our Bonney Lake base, and we work in the area regularly.",
			Context:         "Buckley's location at the base of the Carbon River drainage means it sees meaningful rainfall and the kind of tree cover — dense second-growth fir and cedar — that puts real demands on gutter systems. Homes here tend to be older, which means more spike-and-ferrule installations that have reached the end of their service life. Replacing aging hardware with hidden hanger systems is one of the most common jobs we do in the Buckley area.",
			HeroImage:       "hero-rain-gutters.jpeg",
			NearbyAreas:     []nearbyArea{{City: "Enumclaw", Slug: "enumclaw"}, {City: "Bonney Lake", Slug: "bonney-lake"}, {City: "Black Diamond", Slug: "black-diamond"}},
		},
		{
			City:            "Black Diamond",
			Slug:            "black-diamond",
			MetaTitle:       "Gutter Installation & Repair in Black Diamond, WA | A-Team Gutters",
			MetaDescription: "Gutter installation and repair in Black Diamond, WA. A-Team Gutters serves the Black Diamond area with seamless aluminum systems. Licensed, insured, free estimates.",
			Intro:           "A-Team Gutters serves Black Diamond and the communities in the upper Green River valley. The drive from Bonney Lake takes us through some of the most densely canopied residential areas in our service territory — and the gutter work here reflects it.",
			Context:         "Black Diamond is a small community surrounded by second-growth forest, and the gutter maintenance demands are among the most significant in our service area. Properties here deal with heavy needle loads, persistent moss, and the kind of shade that keeps gutters damp between rain events. If you're in Black Diamond and cleaning gutters more than twice a year, it's worth a conversation about screen systems.",
			HeroImage:       "hero-green-house.jpg",
			NearbyAreas:     []nearbyArea{{City: "Auburn", Slug: "auburn"}, {City: "Enumclaw", Slug: "enumclaw"}, {City: "Buckley", Slug: "buckley"}},
		},
	}

	for _, area := range serviceAreas {
		mux.HandleFunc("GET /service-areas/"+area.Slug+"/", serveServiceArea(serviceAreaTmpl, area))
	}

	// Sitemap and robots.txt
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "https://ateamgutter.com"
	}
	mux.HandleFunc("GET /sitemap.xml", serveSitemap(baseURL, serviceAreas))
	mux.HandleFunc("GET /robots.txt", serveRobotsTxt(baseURL))

	// Static files
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
