# La Luna HMS · Interactive Prototype

A clickable, Axure-style walkthrough of the La Luna Hotel & Resort management system. It covers every screen of the real SvelteKit app — front desk, restaurant, bar, kitchen, management, guest portal and the public booking flow — in a dark amber interface that matches the running product.

**Open `index.html` in Chrome, Edge or Firefox.** No server, build step or internet connection is required.

---

## What is inside

| File / folder | Purpose |
|---|---|
| `index.html` | Landing page with a live phone preview |
| `prototype.html` | The prototype itself (hash routing, modals, drawer) |
| `assets/css/` | Prototype and landing styles |
| `assets/js/` | Screen data, popup data and the runtime |
| `assets/fonts/` | Bundled Inter + Oxanium webfonts |
| `assets/icons/` | Bundled Phosphor icon font |
| `axure/` | 34 screen captures, 19 popup captures, hotspot map and build guide |
| `export/` | Clean standalone copy for submission |
| `tools/` | Playwright verification, capture and documentation scripts |

---

## Using the prototype

- **Portal** — choose a workstation: Front Desk, Restaurant, Bar, Kitchen, Manager, Guest or Booking.
- **Navigate** — tap rows, tiles, quick actions, tabs and the back arrow. Every route has a stable hash URL such as `#/frontdesk/rooms`.
- **Popups** — 19 modals (bottom sheets and centre dialogs) open from the screens. Close with X, the scrim, Cancel/Done or `Esc`.
- **Screen map** — the toolbar **Screens** button (or `M`) opens a drawer listing all 34 screens and 19 popups.
- **Hotspots** — **Hotspots** (or `H`) outlines every clickable target on the current screen.
- **Focus mode** — hides the toolbar for a clean demo.
- **Reset** — returns to the portal.
- **Deep links** — `prototype.html#/guest/dining` opens a screen directly; add `?modal=wifi` to open a popup, or `?capture=1&screen=/frontdesk` for a clean capture view.

### Keyboard

| Key | Action |
|---|---|
| `M` | Screen map |
| `H` | Reveal hotspots |
| `Esc` | Close popup / drawer / focus mode |

---

## Coverage

- **34 screens** across 8 areas
- **19 popups**
- **378 hotspots**
- **0 broken links** — verified by `tools/verify.mjs`

Departments: Front Desk (11 screens), Guest (5), Restaurant (4), Bar (4), Kitchen (4), Manager (3), Booking (2), Portal (1).

The data is deliberately realistic: Malawi kwacha pricing, Lake Malawi dishes (nsima, chambo, mandasi), local staff names, occupancy and folio maths, bar variance, kitchen BOM costing and manager approvals. All screens work offline.

---

## Browser support

Chrome/Edge 120+, Firefox 120+, Safari 17+. The layout is responsive: below 520 px the device frame becomes a full-screen mobile app; above it the prototype sits in a phone frame on a presentation stage.

`prefers-reduced-motion` is respected. The prototype has no analytics, network calls or storage.

---

## Rebuilding the kit

From `tools/` (requires Node 20+ and system Chrome):

```powershell
npm install            # installs playwright-core
node verify.mjs        # route/action/modal regression + QA screenshots
node capture.mjs       # regenerates axure/screens, axure/popups, hotspot map
node docs.mjs          # regenerates SITEMAP.md, INTERACTION-MATRIX.md, CSVs
```

`capture.mjs` uses `C:/Program Files/Google/Chrome/Application/chrome.exe`. Change the `exe` constant if Chrome lives elsewhere.

---

## Axure

See **[axure/README.md](axure/README.md)** and **[axure/AXURE-BUILD-GUIDE.md](axure/AXURE-BUILD-GUIDE.md)**. The kit contains high-resolution captures, the complete hotspot coordinate map and a build checklist for Axure RP 11.

A `.rp` project file is a proprietary Axure binary and cannot be generated programmatically; the kit is designed to be assembled in Axure in one short, mechanical pass. The working HTML prototype is the clickable deliverable.

---

*La Luna Hotel & Resort · HMS degree final project · Prototype 1.0*
