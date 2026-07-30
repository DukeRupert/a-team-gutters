package main

import (
	"fmt"
	"html/template"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// The responsive-image manifest. scripts/build-images.sh generates, for every
// source in assets-src/images:
//
//	<base>.jpg        recompressed fallback at native size
//	<base>-<w>.avif   one per ladder width <= source width
//	<base>-<w>.webp   likewise
//
// At startup we scan static/images for those files and read the fallback's
// intrinsic dimensions off the JPEG header. Templates then call {{picture ...}}
// and never hardcode a width, a height, or a ladder step — which is how the
// three wrong-aspect-ratio CLS bugs got in here in the first place.

const imageRoot = "static/images"

type imageVariants struct {
	Fallback      string // URL path of the <img src> fallback
	Width, Height int    // intrinsic dimensions of the fallback
	AVIF, WebP    []int  // available widths, ascending
}

// imageManifest is keyed by the image's path under static/images without an
// extension — "hero-copper-gutters", "gallery/dark-gutters-gable".
var imageManifest = map[string]*imageVariants{}

func loadImageManifest() {
	err := filepath.WalkDir(imageRoot, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		ext := strings.ToLower(filepath.Ext(p))
		if ext != ".avif" && ext != ".webp" {
			return nil
		}
		// static/images/gallery/foo-640.avif -> "gallery/foo", 640
		rel := filepath.ToSlash(strings.TrimPrefix(p, imageRoot+string(os.PathSeparator)))
		stem := strings.TrimSuffix(rel, filepath.Ext(rel))
		dash := strings.LastIndex(stem, "-")
		if dash < 0 {
			return nil
		}
		w, convErr := strconv.Atoi(stem[dash+1:])
		if convErr != nil {
			return nil
		}
		key := stem[:dash]

		v := imageManifest[key]
		if v == nil {
			v = &imageVariants{}
			imageManifest[key] = v
		}
		if ext == ".avif" {
			v.AVIF = append(v.AVIF, w)
		} else {
			v.WebP = append(v.WebP, w)
		}
		return nil
	})
	if err != nil {
		log.Printf("image manifest: %v", err)
	}

	for key, v := range imageManifest {
		sort.Ints(v.AVIF)
		sort.Ints(v.WebP)

		// The fallback carries the intrinsic aspect ratio, so it must exist.
		found := false
		for _, ext := range []string{".jpg", ".jpeg", ".png"} {
			fp := filepath.Join(imageRoot, filepath.FromSlash(key)+ext)
			f, err := os.Open(fp)
			if err != nil {
				continue
			}
			cfg, _, decErr := image.DecodeConfig(f)
			f.Close()
			if decErr != nil {
				log.Printf("image manifest: %s: %v", fp, decErr)
				continue
			}
			v.Fallback = "/" + imageRoot + "/" + key + ext
			v.Width, v.Height = cfg.Width, cfg.Height
			found = true
			break
		}
		if !found {
			log.Printf("image manifest: no fallback image for %q — run scripts/build-images.sh", key)
			delete(imageManifest, key)
		}
	}
	log.Printf("image manifest: %d images", len(imageManifest))
}

func webpURL(key string, w int) string {
	return fmt.Sprintf("/%s/%s-%d.webp", imageRoot, key, w)
}

func srcset(key string, widths []int, ext string) string {
	parts := make([]string, 0, len(widths))
	for _, w := range widths {
		parts = append(parts, fmt.Sprintf("/%s/%s-%d.%s %dw", imageRoot, key, w, ext, w))
	}
	return strings.Join(parts, ", ")
}

// picture renders a <picture> with AVIF and WebP sources plus a JPEG fallback.
// Options are trailing key/value pairs:
//
//	{{picture "hero-copper-gutters" "alt" "Copper gutters…" "sizes" "100vw" "priority" "true"}}
//
// Recognised keys: alt, sizes (default 100vw), class, priority ("true" =>
// eager + fetchpriority=high, otherwise lazy), style. Unknown keys are an
// error so a typo fails loudly at render time instead of silently dropping.
//
// A global `picture { display: contents }` rule keeps the wrapper transparent
// to layout, so existing `.hero-image img { … }` style rules still apply.
func picture(key string, opts ...string) (template.HTML, error) {
	if len(opts)%2 != 0 {
		return "", fmt.Errorf("picture %q: options must be key/value pairs, got %d", key, len(opts))
	}
	v := imageManifest[key]
	if v == nil {
		return "", fmt.Errorf("picture %q: not in image manifest", key)
	}

	var alt, class, style string
	sizes := "100vw"
	priority := false
	loading := ""
	for i := 0; i < len(opts); i += 2 {
		switch opts[i] {
		case "alt":
			alt = opts[i+1]
		case "sizes":
			sizes = opts[i+1]
		case "class":
			class = opts[i+1]
		case "style":
			style = opts[i+1]
		case "priority":
			priority = opts[i+1] == "true"
		case "loading":
			// For above-the-fold images that are not the LCP candidate (the nav
			// logo): eager, but without fetchpriority stealing bandwidth from
			// the hero.
			loading = opts[i+1]
		default:
			return "", fmt.Errorf("picture %q: unknown option %q", key, opts[i])
		}
	}

	esc := template.HTMLEscapeString
	var b strings.Builder
	b.WriteString("<picture>")
	if len(v.AVIF) > 0 {
		fmt.Fprintf(&b, `<source type="image/avif" srcset="%s" sizes="%s">`,
			esc(srcset(key, v.AVIF, "avif")), esc(sizes))
	}
	if len(v.WebP) > 0 {
		fmt.Fprintf(&b, `<source type="image/webp" srcset="%s" sizes="%s">`,
			esc(srcset(key, v.WebP, "webp")), esc(sizes))
	}
	fmt.Fprintf(&b, `<img src="%s" width="%d" height="%d" alt="%s"`,
		esc(v.Fallback), v.Width, v.Height, esc(alt))
	if class != "" {
		fmt.Fprintf(&b, ` class="%s"`, esc(class))
	}
	if style != "" {
		fmt.Fprintf(&b, ` style="%s"`, esc(style))
	}
	switch {
	case priority:
		b.WriteString(` fetchpriority="high" loading="eager" decoding="async"`)
	case loading != "":
		fmt.Fprintf(&b, ` loading="%s" decoding="async"`, esc(loading))
	default:
		b.WriteString(` loading="lazy" decoding="async"`)
	}
	b.WriteString("></picture>")
	return template.HTML(b.String()), nil
}

// preloadAVIF emits a <link rel=preload> for the LCP image's AVIF srcset so the
// browser starts fetching it before it has parsed the hero markup.
func preloadAVIF(key, sizes string) (template.HTML, error) {
	v := imageManifest[key]
	if v == nil {
		return "", fmt.Errorf("preloadAVIF %q: not in image manifest", key)
	}
	if len(v.AVIF) == 0 {
		return "", nil
	}
	esc := template.HTMLEscapeString
	return template.HTML(fmt.Sprintf(
		`<link rel="preload" as="image" type="image/avif" imagesrcset="%s" imagesizes="%s" fetchpriority="high">`,
		esc(srcset(key, v.AVIF, "avif")), esc(sizes))), nil
}

// stripImageExt lets templates pass a filename from Go data ("hero-rain-gutters.jpeg")
// where a manifest key is expected.
func stripImageExt(name string) string {
	return strings.TrimSuffix(name, path.Ext(name))
}
