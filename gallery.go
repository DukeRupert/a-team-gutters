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
	// September 2026 batch. Specs (color, profile, size) come from the client's
	// job notes. Cities appear only where the client named one. The photo-to-job
	// mapping for the before/after pair, the white home with dark trim (side and
	// rear), the farmhouse and the board-and-batten home is inferred, not confirmed.
	// "Tuxedo" is the manufacturer's color name; keep it verbatim.
	{"gallery/cedar-lodge-home-gutters", "Cedar-sided lodge-style home with 5-inch black seamless aluminum gutters and downspout under an overcast Pacific Northwest sky", "Cedar Home — 5\" Black Seamless", "", false},
	{"gallery/gutter-replacement-before", "Before: two-story Bonney Lake home with its old gutters stripped off and stacked on the lawn beside the work truck and ladder", "Before — Bonney Lake Tear-Off", "", false},
	{"gallery/gutter-replacement-after", "After: the same Bonney Lake home with new 6-inch seamless K-style gutters and downspouts in Tuxedo", "After — 6\" Tuxedo K-Style", "", false},
	{"gallery/navy-two-story-gutters", "Navy two-story home in Puyallup with 5-inch low-gloss white gutters above a covered porch and flower beds", "Puyallup — 5\" Low-Gloss White", "", false},
	{"gallery/board-batten-home-gutters", "Board-and-batten home in Bonney Lake with stone accents and 6-inch black seamless gutters and downspout, with mature firs behind", "Bonney Lake — 6\" Black Seamless", "", false},
	{"gallery/white-farmhouse-wraparound", "White farmhouse with a wraparound covered porch and 5-inch black gutters along the full roofline", "Farmhouse — 5\" Black Gutters", "", false},
	{"gallery/white-home-dark-downspout", "White Bonney Lake home with 6-inch seamless K-style gutters in Tuxedo and a corner downspout beside an arched entry window", "Bonney Lake — 6\" Tuxedo Seamless", "", false},
	{"gallery/farmhouse-dark-trim-gutters", "Side of a white home with dark trim and new 6-inch seamless K-style gutters in Tuxedo along the eave", "White Home — 6\" Tuxedo K-Style", "", false},
	{"gallery/cream-home-rear-gutters", "Rear elevation of a cream two-story home with 6-inch Tuxedo gutters running the full fascia line", "Rear Elevation — 6\" Tuxedo", "", false},
	{"gallery/stucco-home-gutters", "Stucco home in Bonney Lake with a tile roof and 6-inch wicker-color gutters along the garage and entry rooflines", "Bonney Lake — 6\" Wicker", "", false},

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
