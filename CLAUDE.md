# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Static brochure website for A-Team Gutters, LLC — a family-owned gutter company in Bonney Lake, WA serving Pierce and King County. The site's primary goal is converting homeowner visitors into estimate requests.

## Tech Stack

- **Go 1.25** (net/http) — serves HTML templates, handles routing and form submission
- **Go HTML templates** — `templates/*.html` parsed via `template.ParseGlob`
- **HTMX** — server round-trip interactivity (contact form submission, content swaps)
- **Alpine.js** — client-side UI state (mobile menu, accordions, toggles)
- **Tailwind CSS** — utility-first styling, purged at build time

## Build & Run Commands

```bash
# Build and run locally
go build -o server . && ./server
# Server listens on :8080 (override with PORT env var)

# Tailwind CSS build (watch mode for development)
npx tailwindcss -i ./static/css/input.css -o ./static/css/output.css --watch

# Docker
docker build -t dukerupert/a-team-gutters:latest .
```

## Project Structure

Currently a single `main.go` with inline route handlers and one template (`classic.html`). The planned architecture as the site grows:

- `main.go` — entry point, route registration
- `templates/` — Go HTML templates. Target: `base.html` layout with page templates extending via `{{template "base" .}}`, partials (nav, footer, cards) in `templates/partials/`
- `static/` — CSS, images, and other assets served at `/static/`
- `docs/` — development plan and brand specs
- `deploy/` — docker-compose configs for dev and prod environments

## Deployment

Docker containers on a Hetzner VPS. Caddy reverse proxy handles TLS via Let's Encrypt. DNS through Cloudflare.

- **Dev:** Push to `dev` branch triggers GitHub Actions → builds Docker image → deploys to VPS via SSH (port 3000)
- **Prod:** Push to `main` branch triggers prod deploy workflow
- Docker images pushed to DockerHub under `dukerupert/a-team-gutters`

## Brand & Design

**Option E — A-Team Classic** (Black / Silver / Red). Full development plan in `docs/web-dev-plan.md`.

**Palette:** Background `#FFFFFF` · Dark sections `#0A0A0A` · Surface alt `#F5F5F5` · Silver accent `#9BAAB8` · Text `#111111` · Muted `#777777` · Border `#DDDDDD` · Action red `#B01C1C`

**The red rule:** Red appears only on the primary CTA button and eyebrow/tagline text. Nowhere else.

**Typography:** Barlow Condensed (headings, buttons, uppercase, 700 weight) + Barlow (body, 400/500 weight).

Key design constraints:
- Every page ends with a single CTA: "Get a Free Estimate"
- 5-star Angi rating and social proof must appear above the fold
- PNW-specific climate copy on every service page (rainfall, fir/pine debris, moss)
- Warranty and licensing language must be visible, not buried in footer
- Near-zero border-radius (2px) — intentional, do not round corners further

## Content & Voice

Write like a trusted neighbor, not a corporate contractor. Be geographically specific (name cities, landmarks, local weather patterns). Avoid generic contractor-speak. See `docs/web-dev-plan.md` for page-by-page content priorities and the voice/tone guide.

## Site Pages

Core pages: Home, Gutter Installation/Cleaning/Repair/Guards, Gallery, About, FAQ, Contact. Service area pages for Bonney Lake, Sumner, Puyallup, Auburn, Enumclaw, Buckley, Black Diamond. Blog for seasonal SEO content. Full URL structure in `docs/web-dev-plan.md`.

## Design Context

### Users
Homeowners in Pierce and King County, WA — typically searching for gutter services after noticing a problem (overflowing gutters, fascia rot, moss buildup) or preparing for PNW rain season. They're comparing 2–3 local contractors and deciding fast. They want to see credentials, location, and a clear way to request an estimate. Mobile-heavy traffic — many will find the site via Google Maps or local search on their phone.

### Brand Personality
**Tough, reliable, straightforward.** Blue-collar confidence backed by 30+ years of trade experience. No corporate polish, no marketing fluff — just a crew that shows up, does the work right, and stands behind it. The brand says: "We've been doing this longer than most companies have existed."

### Emotional Goals
The site should make homeowners feel **confidence and relief** — "These guys clearly know what they're doing, I can stop worrying about my gutters." Every design choice should reduce anxiety and build trust: visible license numbers, real experience claims, clear pricing signals, and a frictionless path to requesting an estimate.

### Aesthetic Direction
- **Visual tone:** Dark, industrial, high-contrast. The `#0A0A0A` dark sections anchor the identity — they feel serious and professional without being corporate. White sections breathe. Silver (`#9BAAB8`) adds sophistication without softness.
- **The red rule is sacred:** `#B01C1C` appears ONLY on primary CTA buttons and eyebrow/tagline text. This constraint gives red its power — it means "act now" every time it appears. Never dilute it with decorative use.
- **Typography carries the personality:** Barlow Condensed uppercase headings feel industrial and commanding. Barlow body text is clean and readable. No decorative fonts, no cursive, no playfulness.
- **Near-zero border-radius (2px):** Sharp edges are intentional — they reinforce the tough, no-nonsense brand. Do not round corners further.
- **Texture and depth:** Subtle noise overlays on dark sections, gradient overlays on hero images, and carefully controlled opacity layers create visual richness without clutter.
- **Motion:** Restrained. Fade-up entrance animations, smooth hover transitions on CTAs (skewed wipe effect), and subtle scale/transform feedback on interaction. No bouncing, no parallax, no decorative animation.
- **Anti-references:** Generic contractor sites with stock photos of smiling families, bright green/blue palettes, rounded bubbly UI, or templated WordPress themes. Also avoid overly minimalist tech-startup aesthetics — this is a trades business, not a SaaS product.

### Accessibility
- **Target: WCAG 2.1 AA compliance**
- Semantic HTML throughout — proper heading hierarchy, landmark regions, form labels
- All interactive elements keyboard-accessible with visible focus states (2px solid `#B01C1C` outline)
- Color contrast ratios meeting AA minimums (4.5:1 for body text, 3:1 for large text)
- Alt text on all images — descriptive and location-referenced where natural
- Reduced-motion support via `prefers-reduced-motion` media query
- Click-to-call phone links on mobile
- Form validation with clear, accessible error messages

### Design Principles

1. **Credentials before creativity.** License numbers, years of experience, and verifiable trust signals take priority over visual flair. If it builds trust, make it visible.
2. **Red means action.** The CTA color is reserved exclusively for moments where we want the visitor to do something. Every other element earns attention through contrast, typography, and layout — not color.
3. **Sharp and direct.** Squared-off corners, uppercase condensed headings, high-contrast dark/light sections. The design should feel like a firm handshake — confident, no hesitation.
4. **Local beats generic.** Name real cities, reference PNW weather patterns, use language a Bonney Lake homeowner would use. The design and copy should feel like it belongs to this specific place.
5. **One clear path.** Every page funnels toward a single action: requesting a free estimate. Remove friction, reduce choices, make the next step obvious.
