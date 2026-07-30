package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
)

// Static assets get far-future caching, which is only safe if the URL changes
// when the bytes do. Two mechanisms cover the two kinds of asset:
//
//   - CSS/JS are referenced as {{staticHash "…"}}, which appends ?v=<content
//     hash>. A rebuild changes the query string and browsers refetch.
//   - Images and fonts are referenced by name. Their names encode what they
//     are (hero-copper-gutters-1280.avif, barlow-400.woff2), so replacing the
//     underlying photo at the same name would go stale — bump the ladder or
//     rename in that case rather than reusing a name.

var (
	hashOnce  sync.Once
	hashCache map[string]string
	hashMu    sync.Mutex
)

// staticHash returns urlPath with a ?v=<short content hash> suffix.
func staticHash(urlPath string) string {
	hashOnce.Do(func() { hashCache = map[string]string{} })

	hashMu.Lock()
	defer hashMu.Unlock()
	if v, ok := hashCache[urlPath]; ok {
		return v
	}

	versioned := urlPath
	if fp := strings.TrimPrefix(urlPath, "/"); strings.HasPrefix(fp, "static/") {
		if f, err := os.Open(fp); err == nil {
			h := sha256.New()
			if _, err := io.Copy(h, f); err == nil {
				versioned = urlPath + "?v=" + hex.EncodeToString(h.Sum(nil))[:12]
			}
			f.Close()
		} else {
			log.Printf("staticHash: %v", err)
		}
	}
	hashCache[urlPath] = versioned
	return versioned
}

// cacheControl sets Cache-Control by file type before delegating to h.
func cacheControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch path.Ext(r.URL.Path) {
		case ".woff2", ".woff", ".avif", ".webp", ".jpg", ".jpeg", ".png", ".svg", ".ico":
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		case ".css", ".js":
			// Always requested with ?v=<hash>; the hash is the cache key.
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		case ".webmanifest", ".json", ".xml":
			w.Header().Set("Cache-Control", "public, max-age=3600")
		default:
			w.Header().Set("Cache-Control", "public, max-age=86400")
		}
		h.ServeHTTP(w, r)
	})
}
