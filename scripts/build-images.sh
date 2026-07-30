#!/usr/bin/env bash
# Generate responsive image derivatives from assets-src/images into static/images.
#
# For every photographic source it emits:
#   <base>.jpg            recompressed fallback at native size (also what the
#                         Alpine-driven gallery lightbox loads)
#   <base>-<w>.avif       AVIF at each ladder width <= source width
#   <base>-<w>.webp       WebP at each ladder width <= source width
#
# main.go scans the generated -<w>.avif files at startup to build the srcset
# manifest, so the ladder here is the single source of truth — templates never
# hardcode widths or dimensions.
#
# Requires: avifenc, cwebp, ImageMagick (magick). Re-run after adding a source.
set -euo pipefail

SRC="assets-src/images"
OUT="static/images"
LADDER=(400 640 960 1280 1920)

# Quality: tuned so AVIF lands well under the old JPEGs without visible loss on
# photographic content. AVIF q60 ~= JPEG q85 perceptually for these images.
AVIF_Q=60
AVIF_SPEED=4
WEBP_Q=78
JPEG_Q=82

# Backdrop sources render only inside .page-hero-image, which paints them at
# opacity 0.15 under a near-opaque dark gradient. Detail that survives that is
# not detail anyone can see, so they get a much lower quality tier. Anything
# ever shown at full strength (hero-copper-gutters is the homepage LCP at
# opacity 0.55) must NOT be listed here.
BACKDROPS=(hero-green-house hero-gutter-install hero-gutters-home hero-rain-gutters)
BACKDROP_AVIF_Q=35
BACKDROP_WEBP_Q=55

command -v avifenc >/dev/null || { echo "avifenc not found" >&2; exit 1; }
command -v cwebp   >/dev/null || { echo "cwebp not found" >&2; exit 1; }
command -v magick  >/dev/null || { echo "magick not found" >&2; exit 1; }

total_src=0
total_out=0

while IFS= read -r -d '' src; do
    rel="${src#"$SRC"/}"          # e.g. gallery/foo.jpg
    dir=$(dirname "$rel")
    base=$(basename "$rel"); base="${base%.*}"
    outdir="$OUT/$dir"; [ "$dir" = "." ] && outdir="$OUT"
    mkdir -p "$outdir"

    # Note the trailing \n: `magick identify -format` emits no newline of its
    # own, and a `read` that hits EOF before a delimiter exits non-zero, which
    # under `set -e` kills the script silently on the first file.
    read -r sw sh < <(magick identify -format '%w %h\n' "$src")

    aq=$AVIF_Q; wq=$WEBP_Q; tier=""
    if [[ " ${BACKDROPS[*]} " == *" $base "* ]]; then
        aq=$BACKDROP_AVIF_Q; wq=$BACKDROP_WEBP_Q; tier=" [backdrop]"
    fi
    echo "  $rel (${sw}x${sh})$tier"

    # Recompressed fallback at native size, EXIF stripped, progressive.
    magick "$src" -auto-orient -strip -quality "$JPEG_Q" \
        -sampling-factor 4:2:0 -interlace JPEG "$outdir/$base.jpg"

    for w in "${LADDER[@]}"; do
        [ "$w" -gt "$sw" ] && continue
        magick "$src" -auto-orient -strip -resize "${w}x" -quality 95 "/tmp/_ri.png"
        avifenc -q "$aq" -s "$AVIF_SPEED" "/tmp/_ri.png" "$outdir/$base-$w.avif" >/dev/null
        cwebp -q "$wq" -quiet "/tmp/_ri.png" -o "$outdir/$base-$w.webp"
    done

    # Always make the native width available when it falls between ladder steps,
    # so full-bleed heroes are never upscaled by the browser.
    if [[ ! " ${LADDER[*]} " =~ " ${sw} " ]] && [ "$sw" -lt 1920 ]; then
        magick "$src" -auto-orient -strip -quality 95 "/tmp/_ri.png"
        avifenc -q "$aq" -s "$AVIF_SPEED" "/tmp/_ri.png" "$outdir/$base-$sw.avif" >/dev/null
        cwebp -q "$wq" -quiet "/tmp/_ri.png" -o "$outdir/$base-$sw.webp"
    fi

    total_src=$((total_src + $(stat -c%s "$src")))
done < <(find "$SRC" -type f \( -iname '*.jpg' -o -iname '*.jpeg' -o -iname '*.png' \) -print0)

# ── Icons ──
# The brand logo needs different handling from the photos: it must keep its
# alpha channel (so no JPEG fallback — the source PNG stays the fallback) and it
# renders at 48px, so the photographic ladder is far too wide. Sources live in
# static/images already because the PNG is referenced directly by schema.org
# markup and the web manifest.
ICON_LADDER=(48 96 192)
for icon in logo; do
    isrc="$OUT/$icon.png"
    [ -f "$isrc" ] || { echo "  (no $isrc, skipping)"; continue; }
    iw=$(magick identify -format '%w' "$isrc")
    echo "  $icon.png (icon, ${iw}px)"
    for w in "${ICON_LADDER[@]}"; do
        [ "$w" -gt "$iw" ] && continue
        magick "$isrc" -strip -resize "${w}x" "/tmp/_ri.png"
        avifenc -q "$AVIF_Q" -s "$AVIF_SPEED" "/tmp/_ri.png" "$OUT/$icon-$w.avif" >/dev/null
        cwebp -q 90 -quiet -alpha_q 100 "/tmp/_ri.png" -o "$OUT/$icon-$w.webp"
    done
    total_src=$((total_src + $(stat -c%s "$isrc")))
done

rm -f /tmp/_ri.png
total_out=$(find "$OUT" -type f \( -name '*.avif' -o -name '*.webp' -o -name '*.jpg' \) \
    -newermt '-1 hour' -printf '%s\n' | awk '{s+=$1} END {print s+0}')

echo
echo "sources: $(numfmt --to=iec "$total_src")   generated: $(numfmt --to=iec "$total_out")"
echo "done."
