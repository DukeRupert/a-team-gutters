package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

// The gallery grid used to be rendered entirely by Alpine from a JS array,
// which meant 13 full-size JPEGs (1.7MB) downloaded to fill thumbnail slots,
// with no width/height and so no CLS reservation. The photo list now lives here
// so the grid can be rendered server-side through {{picture}} — responsive,
// correctly sized, and working without JS. Alpine keeps only the lightbox,
// which needs plain URLs in a JS array.

type galleryPhoto struct {
	// Key is the manifest key (path under static/images, no extension).
	Key     string `json:"-"`
	Alt     string `json:"alt"`
	Caption string `json:"caption"`

	// Filled in from the manifest at startup.
	Src      string `json:"src"` // largest WebP — what the lightbox loads
	Portrait bool   `json:"-"`
}

var galleryPhotos = []galleryPhoto{
	{"gallery/copper-gutter-detail", "Copper gutter installation detail with Seattle skyline and overcast Pacific Northwest sky", "Copper Gutter Detail — Seattle", "", false},
	{"hero-copper-gutters", "Copper gutter on residential roofline with Seattle Space Needle visible", "Copper Gutters — Roofline View", "", false},
	{"gallery/green-house-gutters", "Green residential home with new dark seamless gutter installation and PNW fir trees", "Seamless Install — Pacific Northwest", "", false},
	{"gallery/dark-gutters-bay-window", "Dark seamless gutters on a bay window with mature fir trees in background", "Dark Gutters — Bay Window", "", false},
	{"gallery/dark-gutters-gable", "Gable end of residential home with dark seamless gutters and PNW conifers", "Gable Detail — Dark Gutters", "", false},
	{"gallery/craftsman-home-gutters-1", "Craftsman-style home with new gutter system on an overcast day in Pierce County", "Craftsman Home — Pierce County", "", false},
	{"gallery/craftsman-home-gutters-2", "Wide view of craftsman home with complete gutter system and lush PNW landscaping", "Craftsman Home — Full View", "", false},
	{"gallery/white-home-gutters-1", "Multi-story white home with seamless gutter installation", "Multi-Story Install — White Trim", "", false},
	{"gallery/white-home-gutters-2", "White residential home with gutter system and PNW blue sky", "White Home — Wide Angle", "", false},
	{"gallery/gutter-guard-closeup", "Close-up of micro-mesh gutter screen installed along residential roofline", "Gutter Screen Close-Up", "", false},
	{"gallery/copper-gutters-golden-hour", "Home with copper gutter system catching golden hour sunlight", "Copper Gutters — Golden Hour", "", false},
	{"gallery/brick-home-copper-gutters", "Brick home with copper gutter installation and clear blue sky", "Brick Home — Copper Install", "", false},
	{"gallery/dark-roof-home-gutters", "Home with dark roof and seamless gutter system surrounded by mature landscaping", "Dark Roof — Seamless Gutters", "", false},
}

type galleryData struct {
	Photos []galleryPhoto
	// PhotosJSON is a plain string, not template.JS, so html/template escapes
	// the quotes for the x-data attribute; the browser decodes them before
	// Alpine evaluates the expression.
	PhotosJSON string
}

// resolveGalleryPhotos fills Src and Portrait from the image manifest. Call
// after loadImageManifest.
func resolveGalleryPhotos() *galleryData {
	out := make([]galleryPhoto, 0, len(galleryPhotos))
	for _, p := range galleryPhotos {
		v := imageManifest[p.Key]
		if v == nil {
			log.Printf("gallery: %q not in image manifest — skipping", p.Key)
			continue
		}
		p.Portrait = v.Height > v.Width
		// WebP rather than the JPEG fallback: ~40% smaller, and every browser
		// that can run Alpine 3 supports it. Lightbox loads are click-driven,
		// so this never touches LCP.
		p.Src = v.Fallback
		if n := len(v.WebP); n > 0 {
			p.Src = webpURL(p.Key, v.WebP[n-1])
		}
		out = append(out, p)
	}
	j, err := json.Marshal(out)
	if err != nil {
		log.Printf("gallery: %v", err)
	}
	return &galleryData{Photos: out, PhotosJSON: string(j)}
}

func serveGallery(tmpl *template.Template, data *galleryData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
			log.Printf("template error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
