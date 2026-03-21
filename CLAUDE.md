# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Static brochure website for A-Team Gutters, LLC — a family-owned gutter company in Bonney Lake, WA serving Pierce and King County. The site's primary goal is converting homeowner visitors into estimate requests.

## Tech Stack

- **Go** (net/http or Chi) — serves HTML templates, handles routing
- **HTMX** — server round-trip interactivity (contact forms, content swaps)
- **Alpine.js** — client-side UI state (menus, accordions, toggles)
- **Tailwind CSS** — utility-first styling, purged at build time
- **Go HTML templates** — `base.html` layout with page templates and partials

## Build & Run Commands

```bash
# Build
go build -o server .

# Run
./server

# Tailwind CSS build (watch)
npx tailwindcss -i ./static/css/input.css -o ./static/css/output.css --watch

# Docker build & deploy
docker build -t dockerhub-user/a-team-gutters:latest .
docker push dockerhub-user/a-team-gutters:latest
```

## Project Structure

Pages map 1:1 to handlers — if a page exists, there's a handler for it in `handlers/`. Templates use a `base.html` layout that pages extend via `{{template "base" .}}`. Partials (nav, footer, cards) live in `templates/partials/`.

## Deployment

Docker containers on a Hetzner VPS. Caddy reverse proxy handles TLS via Let's Encrypt. DNS through Cloudflare.

## Brand & Design

Three brand options are documented in `docs/` (Option C: Forest+Rain, Option D: Rain+Slate, Option E: A-Team Classic). All share the same typography system: **Barlow Condensed** (headlines, uppercase, 700 weight) and **Barlow** (body, 400/500 weight). See the chosen brand option doc for exact color palette, component specs, and constraints.

Key design constraints across all options:
- Every page ends with a single CTA: "Get a Free Estimate"
- 5-star Angi rating and social proof must appear above the fold
- PNW-specific climate copy on every service page (rainfall, fir/pine debris, moss)
- Warranty and licensing language must be visible, not buried in footer
- Near-zero border-radius (2px) — intentional, do not round corners further

## Content & Voice

Write like a trusted neighbor, not a corporate contractor. Be geographically specific (name cities, landmarks, local weather patterns). Avoid generic contractor-speak. See `docs/site-plan.md` for page-by-page content priorities and the voice/tone guide.

## Site Pages

Core pages: Home, Gutter Installation/Cleaning/Repair/Guards, Gallery, About, FAQ, Contact. Service area pages for Bonney Lake, Sumner, Puyallup, Auburn, Enumclaw, Buckley, Black Diamond. Blog for seasonal SEO content. Full URL structure in `docs/site-plan.md`.
