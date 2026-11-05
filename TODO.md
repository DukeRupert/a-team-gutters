# A-Team Gutters — Feature Checklist

Jessie Patterson · Brochure Site · Signed 03/23/2026 · Two payments ($375 at signing, $375 at launch)

---

## Website

- [x] Up to 5 pages (services, story, contact info)
- [x] Mobile-first responsive design

## Integrations & Setup

- [x] Working contact form connected to client email
- [x] Google Analytics installed and configured
- [ ] Facebook Analytics installed and configured
- [x] Google Business Profile set up and verified (business info, service area, hours, site link)

## Hosting

- [x] Site live on Firefly Software hosting ($35/month, first 3 months free)

---

## Photo batch — 09/23/2026

Client sent job notes on 11/04 (color, profile, size, city for some). Captions
and alt text now carry them. Puyallup and Bonney Lake service-area heroes now
use `navy-two-story-gutters` and `stucco-home-gutters`.

- [x] City per photo, where the client gave one (Bonney Lake, Puyallup)
- [x] Gutter color / profile / size per job
- [ ] Confirm the inferred photo-to-job mapping. Client numbered 7 jobs for
      10 photos; #1 Cedar, #3 Navy, #7 Stucco are certain. Assumed:
      #2 White House (Bonney Lake) = before/after pair + `white-home-dark-downspout`;
      #6 (same specs as #2, no city) = `farmhouse-dark-trim-gutters` +
      `cream-home-rear-gutters`; #5 Farm House = `white-farmhouse-wraparound`;
      #4 (Bonney Lake) = `board-batten-home-gutters`
- [ ] City for #1 Cedar, #5 Farm House, #6
- [ ] Service performed per photo (install / repair / guards / fascia) — lets
      the strongest ones move onto the matching service pages
- [ ] Confirm homeowner consent to publish, esp. the before/after pair

Held back pending a decision — all four are in
`~/Pictures/a-team-gutter/09_23_2026/`, none copied into the repo:

- [ ] `IMG_0996` / `IMG_1217` — mid-construction (housewrap, exposed sheathing,
      tools on the deck). Reads as unfinished work in a "completed jobs" gallery,
      but would be strong on an About / process section
- [ ] `IMG_1373` — crew member cropped mid-frame at the edge; needs a crop, and
      a person in frame needs their OK
- [ ] `IMG_1367` — near-duplicate of `farmhouse-dark-trim-gutters` with lens
      flare and a ladder in shot

Notes:

- Street number above the garage in the "after" shot is blurred in
  `assets-src/images/gallery/gutter-replacement-after.jpg`
- Re-run `scripts/build-images.sh` after adding or replacing any source
