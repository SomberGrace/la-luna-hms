# La Luna HMS · Prototype Style Guide

The prototype mirrors the running SvelteKit app: a dark "night sky" hospitality interface with a warm amber moonlight accent. This guide is the visual contract for the Axure rebuild — use it when you replace a capture with native Axure widgets, or when a new screen is added later.

---

## 1. Canvas & layout

| Token | Value |
|---|---|
| Logical screen | **390 × 844 px** (iPhone 13/14 class) |
| Capture scale | **2×** — PNGs are 780 px wide, place them at 390 px |
| Side gutter | **14 px** (`viewport` padding) |
| Screen rhythm | **12 px** between blocks |
| Grid | 2-column cards, 3-column quick actions |
| Status bar | 36 px (hidden below 520 px in the browser prototype) |
| App bar | 56 px |
| Tab bar | 64 px |
| Device frame (browser only) | 410 × 864 px with 10 px bezel |

Vertical screens are captured **full-content**, so a page image can be taller than 844 px. In Axure set the page height to the value in `pages.csv` and keep the image at `x = 0, y = 0`, width **390 px** (never 780).

---

## 2. Colour

### Core surfaces

| Token | Hex / value | Use |
|---|---|---|
| Night | `#05070F` | App background, page background |
| Surface | `#0B1120` | Cards, sheets |
| Surface 2 | `#0E1628` | Raised rows, tiles |
| Input | `#070C16` | Text fields, selects |
| Hairline | `rgba(148,163,184,.13)` | Panel borders |
| Hairline soft | `rgba(148,163,184,.07)` | Row dividers |
| Text | `#E8EEFB` | Primary text |
| Muted | `#8B9AB3` | Secondary text |
| Dim | `#5B6B85` | Captions, hints |

### Accent & tones

| Tone | Hex | Soft fill | Departments |
|---|---|---|---|
| Amber (brand) | `#F5B942` | `rgba(245,185,66,.14)` | Portal, Front Desk, Guest, Booking |
| Orange | `#FB923C` | `rgba(251,146,60,.14)` | Restaurant |
| Emerald | `#34D399` | `rgba(52,211,153,.13)` | Kitchen, success states |
| Indigo | `#818CF8` | `rgba(129,140,248,.14)` | Manager |
| Purple | `#A78BFA` | `rgba(167,139,250,.14)` | Bar |
| Rose | `#FB7185` | `rgba(251,113,133,.13)` | Errors, out of service |
| Slate | `#94A3B8` | `rgba(148,163,184,.11)` | Neutral notes |

**Primary gradient:** `linear-gradient(150deg, #F8C65A, #E0A221)` with text `#150F02`.
**Success/green text:** `#34D399`. **Negative amounts:** `#34D399` for payments, rose for warnings.

> Use the department tone for icons, left card rules, badges and the active tab. Never use more than one non-neutral tone on a single card.

---

## 3. Typography

| Role | Family | Size / weight |
|---|---|---|
| Brand / page display | **Oxanium** | 23–24 px / 700, letter-spacing `.02em` |
| Stat values | Oxanium | 21 px / 700 |
| Modal titles | Oxanium | 14.5 px / 700 |
| Room numbers | Oxanium | 17 px / 700 |
| Body / row titles | **Inter** | 10.5–11.5 px / 500–650 |
| Secondary text | Inter | 8.5–9 px / 400–500 |
| Section labels | Inter | 9.5 px / 700, uppercase, tracking `.14em` |
| Micro labels / badges | Inter | 7–8.5 px / 700, uppercase |

The UI is intentionally dense — it is a hotel operations tool designed to fit a phone. Do not scale type up when rebuilding; scale the whole 390 px canvas instead.

Fonts in the browser prototype are bundled locally: `assets/fonts/inter-*.woff2` and `assets/fonts/oxanium-*.woff2`. Axure exports do **not** need the fonts because the screen captures already contain the rendered text.

---

## 4. Shape, border & depth

| Element | Radius | Border | Shadow |
|---|---|---|---|
| Hero | 22 px | 1 px tone @ 26% | `0 18px 36px -22px rgba(0,0,0,.9)` |
| Card / list | 17 px | 1 px hairline | `0 14px 30px -24px rgba(0,0,0,.9)` |
| Tile / row icon | 15–16 px / 10 px | 1 px hairline | — |
| Button | 11 px | 1 px hairline | primary: amber glow |
| Chip | 999 px | tone-tinted | — |
| Modal sheet | 24 px top corners | 1 px hairline | `0 -30px 70px -20px rgba(0,0,0,.9)` |
| Center dialog | 22 px | 1 px hairline | same |

Cards use a subtle vertical gradient (`#0D1626 → #0A101B`) rather than a flat fill. A 2 px tone-coloured rule runs down the left edge of stat, room, dish, ticket and timeline cards.

---

## 5. Components

- **Hero** — tone-tinted panel with kicker, Oxanium headline, 10.5 px body, pill chips and a large low-opacity icon watermark at the bottom-right.
- **Stats** — 2×2 grid; uppercase micro label, Oxanium value, dim note, 2 px tone rule.
- **Section head** — uppercase label + optional text link (`View all`, `New docket`) or a muted note.
- **List row** — 34 px tinted icon, title + subtitle, right value/CTA, chevron when clickable. Rows are a single tappable target.
- **Tile grid** — 2 columns, icon, title, subtitle; whole tile is tappable.
- **Quick actions** — 3 columns, 74 px minimum height, icon over label.
- **Filter chips** — horizontal scroll, single-select, amber active state.
- **Room card** — room number, status badge, type · bed, nightly rate; tone rule on the left.
- **Dish row** — name, price, description, tag, circular `+` button.
- **Form card** — 2-column field grid, full-width submit, helper note. Inputs are 9 px radius, focus ring is tone at 55% + 3 px soft halo.
- **Statement** — label/sub/amount rows, total rule in the accent, positive amounts in emerald.
- **Ticket (KDS)** — docket id, meta, quantity lines, note per line, footer total + bump button.
- **Bars** — label + value, 5 px track, tone gradient fill with a soft glow.
- **Timeline** — 8 px tone dot, hairline spine, time/title/sub.
- **Note** — tinted 30 px icon, title + body, optional chevron when tappable.
- **Tray bar** — sticky amber bar with count badge, label and total (in-room dining).
- **App bar** — back button or moon mark, brand + department tag, notification bell with amber dot.
- **Tab bar** — 4–5 items, 19 px icon, 8.5 px label, amber active icon plate + 2 px top indicator.
- **Modal** — grab handle, icon, Oxanium title + subtitle, scrollable body, full-width action row.
- **Toast** — bottom pill, emerald check, 2.6 s auto-hide.
- **Screen map drawer** — right-hand 350 px panel listing all 34 screens and 19 popups.

---

## 6. Motion

| Interaction | Timing |
|---|---|
| Screen enter | 340 ms, `cubic-bezier(.22,.8,.2,1)`, 10 px rise + 0.995 scale |
| Modal enter | 260 ms rise from bottom, 200 ms scrim fade |
| Drawer | 280 ms slide |
| Toast | 220 ms in, auto-hide after 2.6 s |

The browser prototype honours `prefers-reduced-motion`. Capture mode disables all animation.

---

## 7. Iconography

- **Phosphor Icons (Regular)** — bundled as a webfont in `assets/icons/`.
- Sizes: 15–17 px inside 30–34 px tinted containers; 19 px in the tab bar; 30 px in success dialogs.
- Common glyphs: `ph-moon`, `ph-house`, `ph-door-open`, `ph-receipt`, `ph-credit-card`, `ph-fork-knife`, `ph-wine`, `ph-cooking-pot`, `ph-briefcase`, `ph-bell`, `ph-bed`, `ph-key`, `ph-shield-check`, `ph-chart-line-up`, `ph-package`, `ph-broom`, `ph-wifi-high`, `ph-towel`.
- Do not mix icon families. Keep the icon in the department tone.

---

## 8. Accessibility

- Body text meets AA on the dark surfaces; captions are decorative/secondary only.
- Interactive rows are ≥ 34 px tall; primary buttons are 40 px; quick actions are 74 px.
- Focus is always visible: tone border + 3 px soft ring.
- Every icon-only control has an `aria-label`.
- The prototype supports keyboard: **M** screen map, **H** hotspots, **Esc** close, and native tab/enter on buttons.
- Colour is never the only signal — labels and icons accompany every tone.

---

## 9. Asset specification for Axure

| Asset | Folder | Size | Notes |
|---|---|---|---|
| Screens | `axure/screens/` | 780 × variable | Place at 390 px width, preserve aspect |
| Popup context | `axure/popups/` | 780 × 1688 | Full 390 × 844 frame with scrim + sheet |
| Hotspots | `axure/hotspot-map.csv`, `axure/hotspots/*.csv` | — | Logical px, origin top-left |
| Page table | `axure/pages.csv` | — | Route, title, image, height, hotspot count |
| Manifest | `axure/manifest.json` | — | Machine-readable version of all of the above |

Use **PNG** (not JPEG) to keep the hairlines and amber glows crisp. The page background behind each image must be `#05070F` so the rounded panels blend.
