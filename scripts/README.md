# Asset pipeline

## Images

Originals live in `assets-src/images/` and are **not** served. `build-images.sh`
generates everything under `static/images/`:

| output | purpose |
| --- | --- |
| `<base>.jpg` | recompressed fallback at native size — the `<img src>` |
| `<base>-<w>.avif` | one per ladder width ≤ source width |
| `<base>-<w>.webp` | likewise |

```bash
./scripts/build-images.sh      # needs avifenc, cwebp, ImageMagick
```

Re-run after adding or replacing a source. Generated files are committed so the
Docker build stays dependency-free.

`main.go` scans `static/images/**` at startup and builds a manifest keyed by
path-without-extension (`hero-copper-gutters`, `gallery/dark-gutters-gable`).
Templates then call:

```gotemplate
{{picture "hero-copper-gutters" "alt" "…" "sizes" "100vw" "priority" "true"}}
```

The helper emits the AVIF and WebP `<source>` sets plus a fallback `<img>`
carrying the real intrinsic `width`/`height`. **No template hardcodes a width, a
height, or a ladder step** — three pages previously declared `1920x1280` for
images that were actually 1000x545, which is where the layout shift came from.

Options: `alt`, `sizes` (default `100vw`), `class`, `style`, `priority`
(`"true"` ⇒ `fetchpriority=high loading=eager`), `loading`. An unknown option is
a render-time error rather than a silent no-op.

### Quality tiers

- **Photos** — AVIF q60 / WebP q78.
- **Backdrops** (`BACKDROPS` in the script) — AVIF q35 / WebP q55. These render
  only inside `.page-hero-image`, at `opacity: 0.15` beneath a near-opaque dark
  gradient. Their templates also pass `sizes="50vw"` so HiDPI screens don't pull
  the top rung for an image nobody can resolve detail in. `hero-copper-gutters`
  is deliberately excluded: the homepage paints it at 0.55 and it is the LCP
  element.
- **Icons** — `logo.png` keeps its alpha, so the PNG stays the fallback and the
  ladder is 48/96/192 instead of the photographic one.

## Fonts

Self-hosted latin subsets in `static/fonts/`, declared via `@font-face` in
`styles.css` with `font-display: swap`. Only the four weights the design uses
are shipped: Barlow 400/500, Barlow Condensed 600/700. DM Sans was dropped — it
was a whole extra family for two 10px lines in the footer signature.

To refresh, fetch the Google Fonts stylesheet with a modern browser UA, take the
`latin` `unicode-range` block for each face, and download those `woff2` URLs:

```bash
curl -A 'Mozilla/5.0 … Chrome/120.0 …' \
  'https://fonts.googleapis.com/css2?family=Barlow+Condensed:wght@600;700&family=Barlow:wght@400;500&display=swap'
```

## JavaScript

Alpine, the Alpine collapse plugin, and htmx are vendored into `static/js/` at
pinned versions. The previous `@3.x.x` CDN range meant an upstream release could
break the site unannounced. All are deferred. htmx loads **only on `/contact/`**
(via the `js` template block) since that is the only page using it.

## Cache busting

CSS and JS are referenced as `{{staticHash "/static/…"}}`, which appends
`?v=<content hash>`; `static.go` then serves everything under `/static/` with a
one-year `immutable` `Cache-Control`. Images and fonts are versioned by name, so
replace-in-place is not safe for them — add a new name instead.
