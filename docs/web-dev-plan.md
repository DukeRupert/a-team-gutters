# A-Team Gutters, LLC — Website Development Plan

> Prepared by Firefly Software
> Brand Direction: Option E — A-Team Classic (Black / Silver / Red)

---

## Project Overview

A brochure website for A-Team Gutters, LLC — a family-owned gutter contractor based in Bonney Lake, WA. The site's primary goals are:

1. Establish digital legitimacy and trust for a brand-new LLC backed by 30+ years of trade experience
2. Rank in local search for gutter services across Pierce and King Counties
3. Convert visitors into estimate requests via a clean, friction-free contact form

**Business details baked into the build:**
- Business name: A-Team Gutters, LLC
- Phone: (253) 293-7732
- Email: info@ateamgutter.com
- Address: 13010 211th Ave Ct E, Bonney Lake, WA 98391
- Contractor license: TEAMGGL760KN
- Business license: 605469918
- LLC formed: 2024
- Trade experience: 30+ years
- Service area: Pierce County and King County, WA

---

## Brand System

Per Brand Option E — A-Team Classic. Full spec in `brand-option-e.md`.

**Palette:**
- Background: `#FFFFFF`
- Dark sections: `#0A0A0A`
- Surface alt: `#F5F5F5`
- Silver accent: `#9BAAB8`
- Text: `#111111`
- Muted: `#777777`
- Border: `#DDDDDD`
- Action red: `#B01C1C` — CTA buttons and eyebrow/tagline text only

**Typography:** Barlow Condensed (headings, buttons) + Barlow (body)

**The red rule:** Red appears only on the primary CTA button and eyebrow/tagline text. Nowhere else. See brand spec for full enforcement notes.

---

## Site Architecture

### Core Pages

| Page | URL | Primary SEO Target |
|---|---|---|
| Home | `/` | "gutter company Bonney Lake WA" / "gutter contractor Pierce County" |
| Gutter Installation | `/services/gutter-installation/` | "gutter installation Bonney Lake WA" |
| Gutter Cleaning | `/services/gutter-cleaning/` | "gutter cleaning Bonney Lake WA" |
| Gutter Repair | `/services/gutter-repair/` | "gutter repair Bonney Lake WA" |
| Gutter Guards | `/services/gutter-guards/` | "gutter guards Pierce County WA" |
| Fascia & Soffit Repair | `/services/fascia-soffit-repair/` | "fascia soffit repair Bonney Lake" |
| Gallery | `/gallery/` | Supporting trust / photo SEO |
| About | `/about/` | Brand trust, E-E-A-T signals |
| Service Areas | `/service-areas/` | Index page linking to city pages |
| Bonney Lake | `/service-areas/bonney-lake/` | "gutters Bonney Lake WA" |
| Sumner | `/service-areas/sumner/` | "gutter company Sumner WA" |
| Puyallup | `/service-areas/puyallup/` | "gutter installation Puyallup WA" |
| Auburn | `/service-areas/auburn/` | "gutter repair Auburn WA" |
| Enumclaw | `/service-areas/enumclaw/` | "gutters Enumclaw WA" |
| Buckley | `/service-areas/buckley/` | "gutter company Buckley WA" |
| Black Diamond | `/service-areas/black-diamond/` | "gutters Black Diamond WA" |
| FAQ | `/faq/` | Long-tail question queries |
| Contact | `/contact/` | Conversion endpoint |
| Blog | `/blog/` | Ongoing SEO content |

---

## Page Specifications

### Home (`/`)

**Structure:**
1. **Hero** — Dark (`#0A0A0A`) full-width section. Eyebrow in red: *"Gutter Excellence, A-Team Style"*. H1 in white Barlow Condensed. Single red primary CTA button. Phone number visible and clickable.
2. **Trust bar** — Immediately below hero. Five items: ⭐ 5-Star Rated on Angi | 30+ Years Experience | Licensed & Insured (TEAMGGL760KN) | Free Estimates | Satisfaction Guaranteed.
3. **Services grid** — Six service cards on white background. Silver left-border accent. Brief description + "Learn more →" link per card.
4. **Why A-Team** — Three-column section on `#F5F5F5`. Family-owned, 30 years of hands-on experience, Bonney Lake-based. No subcontractors.
5. **Service area callout** — Named city list with a simple map visual or county reference.
6. **CTA band** — Full-width `#0A0A0A` section. White headline, red primary button.
7. **Footer** — Logo, phone, email, address, license numbers, nav links, copyright.

**Meta:**
- Title: `Gutter Installation & Repair | Bonney Lake, WA | A-Team Gutters`
- Description: `A-Team Gutters LLC — family-owned gutter installation, cleaning, repair, and guard services in Bonney Lake, Sumner, Puyallup, and Pierce County. Free estimates. Licensed contractor TEAMGGL760KN.`

---

### Service Pages (all five follow this template)

**Structure:**
1. **Page hero** — Dark section, H1 targeting the service + location keyword, red eyebrow, short subhead, red CTA button.
2. **What it is** — Two-column layout: left is explanatory prose (PNW-specific), right is a bullet list of what's included.
3. **Why it matters in the PNW** — Short section specific to Pacific Northwest conditions: heavy rainfall, fir/pine debris, moss, steep rooflines, ice in winter.
4. **Our process** — 3–4 numbered steps showing what happens from estimate to completion.
5. **Pricing transparency** — A general range or "most residential jobs run between $X and $Y" — enough to pre-qualify leads without committing to firm quotes.
6. **FAQ accordion** — 4–6 questions specific to this service.
7. **Related services** — Cards linking to the other four service pages.
8. **Estimate CTA** — Red button + phone number.

**Service-specific SEO notes:**

- **Gutter Installation:** Target "seamless gutter installation Bonney Lake," "new gutters Pierce County." Mention .027 and .032 aluminum thickness, hidden hangers, ceramic-coated screws, on-site forming — these specifics build E-E-A-T and match how Seth described his work.
- **Gutter Cleaning:** Target "gutter cleaning Bonney Lake," "gutter cleaning Puyallup WA." Lean into seasonal demand — fall needle drop, spring moss, winter prep. Mention downspout flushing and flow testing.
- **Gutter Repair:** Target "gutter repair Bonney Lake," "leaking gutters Pierce County." Speak to the cost-of-waiting angle — small leaks become fascia rot, foundation issues, cracked driveways.
- **Gutter Guards:** Target "gutter guards Bonney Lake WA," "leaf guard installation Pierce County." Address the pine needle problem specifically — most guard products are reviewed for leaf loads, not needle loads. Position A-Team as knowing the difference.
- **Fascia & Soffit Repair:** Target "fascia repair Bonney Lake," "soffit replacement Pierce County." Position as an add-on to gutter installs — discovered during assessment, handled in the same visit.

---

### Service Area Pages

Each city page follows a tight template — unique enough for SEO, consistent enough to build efficiently:

1. Short intro paragraph naming the city and A-Team's presence there
2. Services offered (same five, linked)
3. One paragraph of city-specific context (nearest landmarks, neighborhood types, common gutter issues in that area)
4. Trust signals restatement
5. Estimate CTA

These pages must not be copy-paste clones with only the city name swapped — Google penalizes thin duplicate content. Each needs at least one paragraph of genuinely city-specific information.

---

### About (`/about/`)

- Seth and Jessica's story — 30 years of experience, why they started A-Team in Bonney Lake
- Photo of crew, truck, or job site (real photos, not stock)
- License and insurance details: Contractor TEAMGGL760KN, Business License 605469918
- LLC formed 2024, trade experience since mid-1990s
- Community connection — family business, serving the Pierce County foothills

**E-E-A-T note:** The About page is a significant Google trust signal for local service businesses. Named individuals, license numbers, verifiable credentials, and real photos all contribute to Experience, Expertise, Authoritativeness, and Trustworthiness scoring.

---

### FAQ (`/faq/`)

Seed questions to build out — each becomes a target for long-tail search:

- How much does gutter installation cost in Washington State?
- How often should I clean my gutters in the Pacific Northwest?
- What size gutters do I need — 5 inch or 6 inch?
- Do I need gutter guards if I have pine trees?
- How long does a gutter installation take?
- What is fascia and why does it matter?
- Are you licensed and insured in Washington State?
- Do you offer free estimates?
- What areas do you serve?
- What's the difference between seamless and sectional gutters?

---

### Contact (`/contact/`)

**Form fields:**
- Full name (required)
- Phone number (required)
- Email address (required)
- Service address (required)
- Service type — dropdown: Installation / Cleaning / Repair / Guards / Fascia & Soffit / Not sure
- Message / project details (optional)
- Submit button — red primary, "Request a Free Estimate"

**Form behavior:**
- Submits to `info@ateamgutter.com`
- Success message: "Thanks — we'll be in touch within 24 hours."
- Phone number displayed prominently alongside the form for call-in preference
- No Captcha friction if possible — use honeypot field instead

**Additional contact info displayed:**
- (253) 293-7732
- info@ateamgutter.com
- Bonney Lake, WA 98391
- Hours: Mon–Fri 9:00 AM – 5:30 PM
- Licensed contractor: TEAMGGL760KN

---

## SEO Foundation

### On-Page Essentials (every page)

- Unique `<title>` tag — format: `[Service] [City] WA | A-Team Gutters`
- Unique meta description — 150–160 characters, include city and primary keyword
- One H1 per page containing the primary keyword
- H2s for supporting keywords and section structure
- Image alt text on every photo — descriptive, location-referenced where natural
- Internal linking — every service page links to related services and relevant city pages
- Canonical tags to prevent duplicate content across city page variants

### Schema Markup

Implement the following structured data:

**LocalBusiness schema (site-wide):**
```json
{
  "@type": "LocalBusiness",
  "name": "A-Team Gutters LLC",
  "telephone": "253-293-7732",
  "email": "info@ateamgutter.com",
  "address": {
    "streetAddress": "13010 211th Ave Ct E",
    "addressLocality": "Bonney Lake",
    "addressRegion": "WA",
    "postalCode": "98391"
  },
  "areaServed": ["Bonney Lake", "Sumner", "Puyallup", "Auburn", "Enumclaw", "Buckley", "Black Diamond"],
  "licenseNumber": "TEAMGGL760KN"
}
```

**Service schema** — on each service page, implement `Service` type with name, description, areaServed, and provider reference.

**FAQPage schema** — on the FAQ page, implement `FAQPage` with `Question` and `Answer` pairs. This enables rich results (expanded FAQ display) in Google SERPs.

**BreadcrumbList schema** — on all interior pages for navigational rich results.

### Google Business Profile

Before or at launch, ensure the Google Business Profile is:
- Claimed and verified at the Bonney Lake address
- Category set to "Gutter Cleaning Service" (primary) + "Gutter Installation Service" (secondary)
- Phone, hours, and website URL matching the site exactly
- License number added to the profile
- Service area set to Pierce County + King County
- Photos uploaded (truck, crew, job photos)

This is the single highest-leverage SEO action for a local contractor — GBP drives the map pack results that appear above organic listings.

### Content Velocity

At launch, the site should have:
- 5 service pages (fully written, PNW-specific)
- 7 city pages (unique, non-duplicate)
- About page with real credentials
- FAQ page with 10+ questions
- 1–2 blog posts (seasonal: fall gutter prep, spring cleaning)

Post-launch, target 1 blog post per month minimum. Priority topics:
- "How to prepare your gutters for a Pacific Northwest winter"
- "Pine needles vs. leaves: why PNW gutter guards need to be different"
- "5" vs 6" gutters: what's right for your Bonney Lake home"
- "How to spot fascia rot before it becomes a $5,000 problem"

---

## Contact Form — Technical Notes

Build as a native Go + htmx form (no third-party embed). On submit:

1. Validate fields server-side
2. Send formatted email to `info@ateamgutter.com` via Postmark (consistent with Firefly stack)
3. Return htmx partial swapping the form for a success message
4. Log submission to database for backup (in case email delivery fails)
5. Honeypot field for spam prevention — hidden field that legitimate users never fill; bots do

**Email format to Seth/Jessica:**
```
New Estimate Request — A-Team Gutters

Name:    [name]
Phone:   [phone]
Email:   [email]
Address: [address]
Service: [service type]
Message: [message]

Submitted: [timestamp]
```

---

## Footer Content

Every page footer includes:

- A-Team Gutters LLC logo
- (253) 293-7732 — click-to-call
- info@ateamgutter.com
- 13010 211th Ave Ct E, Bonney Lake, WA 98391
- WA Contractor License: TEAMGGL760KN
- WA Business License: 605469918
- Navigation: Home | Services | Gallery | About | FAQ | Contact
- © 2024 A-Team Gutters, LLC. All rights reserved.

Displaying license numbers in the footer is a local SEO and trust signal — it makes them indexable and verifiable.

---

## Development Phases

### Phase 1 — Foundation (launch-ready)
- Home
- 5 service pages
- About
- Contact with working form
- Footer with license info
- Basic schema markup
- Google Search Console + Analytics connected

### Phase 2 — SEO Expansion (30–60 days post-launch)
- 7 city/service area pages
- FAQ page with schema
- Gallery page
- First 2 blog posts
- GBP fully configured and linked

### Phase 3 — Content Growth (ongoing)
- Monthly blog posts
- Additional city pages as service area expands
- Review integration (Google reviews widget or embed)
- Performance reporting setup

---

## Pre-Launch Checklist

- [ ] Domain registered (`ateamgutter.com` or `ateamgutters.com` — confirm with client)
- [ ] Email `info@ateamgutter.com` live and tested
- [ ] SSL certificate active
- [ ] All pages have unique title tags and meta descriptions
- [ ] Schema markup validated via Google Rich Results Test
- [ ] Contact form tested end-to-end — submissions arriving at correct email
- [ ] All phone numbers click-to-call on mobile
- [ ] Google Business Profile claimed and website URL set
- [ ] Google Search Console verified
- [ ] Google Analytics 4 connected
- [ ] Sitemap.xml generated and submitted
- [ ] robots.txt configured
- [ ] Page speed baseline captured (Lighthouse / PageSpeed Insights)
- [ ] License numbers visible in footer
- [ ] Real photos uploaded (not placeholder stock)

---

*Prepared by Firefly Software · fireflysoftware.dev*
