# La Luna HMS · Axure RP Build Guide

This kit turns the working HTML prototype into an Axure RP project with the same 34 screens, 19 popups and 378 clickable targets. Everything is captured at 2× and every coordinate is recorded, so the Axure build is a mechanical assembly pass rather than a redesign.

---

## 0. Read this first — about the `.rp` file

The kit ships a **real, clickable Axure project**: `LaLuna-HMS-prototype_RP11.rp` (53 pages, 500 widgets, 288 interactions, 0 problems), saved by Axure RP 11 itself. Open it directly — every screen capture, popup and hotspot is already in place. A copy is installed at `C:\Users\Lanterns\Documents\Axure\LALUNA\Laluna HMSv3.rp`.

`LaLuna-HMS-prototype.rp` is the generator's RP 9 output, kept so the build stays reproducible. It opens in Axure RP 11 as well, and **Save As** produces the RP 11 file. All content, page names, geometry and interactions survive the upgrade.

The rest of this guide is retained for two purposes:

1. **Manual assembly** — if you prefer to rebuild any screen by hand in Axure.
2. **Regeneration** — §14 documents how to rebuild the `.rp` from the CSV source data after a screen, hotspot or popup changes.

> **Do not overwrite** `C:\Users\Lanterns\Documents\Axure\LALUNA\Laluna HMSv2.rp`. `Laluna HMSv3.rp` is the intended v3 destination.

---

## 1. What is in the kit

```
prototype/
├── index.html                     ← landing page (open this first)
├── prototype.html                 ← the clickable prototype
├── assets/                        ← fonts, icons, CSS, data + runtime
├── axure/
│   ├── AXURE-BUILD-GUIDE.md       ← this file
│   ├── SITEMAP.md                 ← routes, page list, popup list
│   ├── INTERACTION-MATRIX.md      ← every hotspot, grouped by screen
│   ├── STYLE-GUIDE.md             ← colours, type, components
│   ├── pages.csv                  ← page name, route, image, height
│   ├── hotspot-map.csv            ← all 378 hotspots (global)
│   ├── hotspot-map.json           ← same data, machine readable
│   ├── hotspots/                  ← one CSV per screen (use these)
│   ├── manifest.json              ← capture metadata
│   ├── screens/                   ← 34 × @2x screen PNGs (780 px wide)
│   └── popups/                    ← 19 × @2x popup PNGs (780 × 1688)
└── export/                        ← standalone copy of the HTML prototype
```

There are **34 screens**, **19 popups** and **378 hotspots**. Department counts: Front Desk 11, Guest 5, Restaurant 4, Bar 4, Kitchen 4, Manager 3, Booking 2, Portal 1.

---

## 2. Before you start

- **Axure RP 11** (RP 10 also works; action names differ slightly and are noted below).
- The kit folder copied somewhere local. Keep the PNG paths unchanged so the CSVs stay valid.
- About 2–4 hours for the full build, or ~45 minutes for a vertical slice (Portal → Front Desk Hub → Rooms → a popup).
- Basic familiarity with Axure pages, hot spots and lightboxes.

---

## 3. Quick start (15 minutes)

1. Open Axure RP 11 → **File → New** → save as `LaLuna-HMS-Prototype.rp`.
2. Use the **iPhone 17** device canvas (**402 × 874 px**); the generated project adds a centered **442 × 914 px** phone frame, 9:41 status row, signal/Wi-Fi/100% battery icons, Dynamic Island and bottom gesture bar.
3. Create the folder structure from §4.
4. Put `axure/screens/portal-01-station-sign-in.png` in a **402 × 780 px** vertical scrolling panel at local **X = 6, Y = 0**, width **390**. The panel sits at page **X = 0, Y = 70** below the status row.
5. Drag six hot spots over the six role cards using the first six rows of `axure/hotspots/portal-01-station-sign-in.csv`.
6. Wire each hot spot to `Open Link → Page` and pick the matching page (create the target pages as you go).
7. Press **Preview** (F5). You now have a clickable Portal → stations flow.

---

## 4. Page structure

Create these folders (right-click the Pages panel → **Add Folder**) so the Axure sitemap matches the product:

| Folder | Pages |
|---|---|
| `00 · Portal` | Station Sign In |
| `01 · Booking` | Book Your Stay, Guest Account |
| `02 · Guest` | Your Stay, In-Room Dining, My Room Statement, Concierge & Requests, Resort Guide & Directory |
| `03 · Front Desk` | Front Desk Hub, Room Inventory, Rapid Walk-In, Reservation Check-In, Guest Billing Folio, Payment Terminal, Check-Out Terminal, Stay Extension, Housekeeping, Offline Sync & Queue, Reception Reports |
| `04 · Restaurant` | Restaurant Hub, Menu Items, Dockets, Restaurant Report |
| `05 · Bar` | Bar Hub, Stock Audit, Record Sales, Bar Daily Report |
| `06 · Kitchen` | Kitchen Orders, Kitchen Stock, Recipe BOM Costing, Kitchen Shift Report |
| `07 · Manager` | Executive Desk, Shift Audits & Approvals, Consolidated Report |
| `08 · Popups` | 19 popup pages (see §7) |

Name pages with the screen number from `pages.csv`, for example `03.01 Front Desk Hub`. Put the route in **Page Notes** (`/frontdesk`) so the HTML prototype and the `.rp` stay traceable.

### Page style (apply to all pages)

| Setting | Value |
|---|---|
| Device screen | **402 × 874 px** (iPhone 17) |
| Phone frame | **442 × 914 px**, centered around the device screen |
| Scrolling viewport | **402 × 780 px**, at page X=0/Y=70 |
| Screenshot | **390 px** wide, centered at panel X=6; height from `pages.csv` |
| Background colour | **`#05070F`** |
| Page alignment | centered device preview |
| Scrollbars | vertical inside the content panel; status bar, island and gesture bar stay fixed |

The status asset is `phone-status-bar.png`. The phone shell supplies the Dynamic Island and home indicator; the content panel clips full-height captures to the phone viewport.

---

## 5. Import the screens

For each row in `pages.csv`:

1. Create a vertical **402 × 780 px** dynamic panel at page **X = 0, Y = 70**.
2. Place `axure/screens/<image>` inside its content state at **X = 6, Y = 0**.
3. Set image **W = 390** and keep its aspect ratio (the PNG is 780 px wide, so it lands at 390 × `height`). Do not leave it at 780 px.
4. Lock the image so hot spots stay above it.
5. Keep the phone page at **402 × 874 px**; use the `height` column in `pages.csv` for the scrollable content height.
6. Add the route and a one-line purpose to **Page Notes**.

The captures include the app bar, content and tab bar. Their full height remains in the panel state, so pages taller than the viewport scroll while the phone frame and system chrome remain stationary.

---

## 6. Hotspots and interactions

### 6.1 Draw the hotspots

For each screen, open `axure/hotspots/<same-name>.csv`. Columns:

| Column | Meaning |
|---|---|
| `label` | Human name of the target (use it as the widget name) |
| `action` | `nav`, `modal`, `close`, `toast`, `back`, `tray` |
| `target` | Route (`/frontdesk`) or popup id (`room-status`) |
| `x`, `y`, `w`, `h` | Logical pixels from the image's top-left |

Draw a **Hot Spot** widget for each row at **panel-local `x + 6, y`** with size `w × h`; the 6 px X offset centers the 390 px capture in the 402 px screen. Keep the hotspot inside the same scrolling panel state as its image so it moves with the content.

During the build, style hot spots with a 10% amber fill and a dashed amber border (`#F5B942`) so you can see coverage. Clear the fill and border before publishing, or keep a separate "hotspot reveal" build.

### 6.2 Wire the actions

| `action` | Axure interaction |
|---|---|
| `nav` | **Click or Tap → Open Link → Page** → choose the page for `target` |
| `back` | **Click or Tap → Open Link → Back** (or the explicit back route from `pages.csv`) |
| `modal` / `tray` | **Click or Tap → Show Lightbox** → pick the popup page (see §7) |
| `close` | **Click or Tap → Hide Lightbox** (use on the scrim, the X and the Cancel/Done buttons) |
| `toast` | **Click or Tap → Show** the `Toast` master → **Wait 2600 ms** → **Hide** it (or **Set Text** then show) |
| `toast` with `close` | Hide the lightbox first, then show the toast |

Use **Axure RP 11** wording: *Click or Tap → Open Link → Page*. In RP 10 the same action is *On Click → Open Page*. Both produce the same prototype.

### 6.3 Navigation targets

The 211 page-opening hotspots are all listed in `INTERACTION-MATRIX.md`. Every `target` is a route that exists in `SITEMAP.md`, so there are no dead links. The bottom tab bar changes per department:

- Front Desk: Hub · Rooms · Walk-In · Folio · Checkout
- Restaurant: Hub · Menu · Dockets · Reports
- Bar: Hub · Stock · Sales · Reports
- Kitchen: Orders · Stock · Recipes · Reports
- Manager: Overview · Audits · Consolidated
- Guest: Stay · Dining · My Bill · Services

Each tab is already a hotspot in the capture — wire it to its route. Sub-pages that are not tabs (Housekeeping, Payments, Offline…) are reached from the Hub quick actions, exactly as in the browser prototype.

### 6.4 The toast master

Create a master named `Toast`:

1. 350 × 44 px dark pill (`#0A101C`, 12 px radius, hairline border).
2. 15 px emerald check icon + 10 px text.
3. Position it 78 px from the bottom, centred.
4. Hide it on page load.
5. On a `toast` hotspot: **Show Toast → Set Text → Wait 2600 → Hide**.

You only need one toast master for the whole file. The text values are in `INTERACTION-MATRIX.md` (the "Toast / confirm" rows) and in the browser prototype.

### 6.5 Optional — hotspot reveal

Add a small fixed **Reveal hotspots** button on a master, then use a dynamic panel or a global interaction to show/hide a translucent amber overlay behind the hot spots. This mirrors the **H** shortcut in the HTML prototype and is useful during reviews.

---

## 7. Popups & lightboxes

The 19 popups are defined in `SITEMAP.md` (with the screen that opens each one) and captured at full phone size in `axure/popups/`.

1. Create the 19 pages in the `08 · Popups` folder.
2. Place each `popups/popup-XX-*.png` at **W = 390**, **X = 6, Y = 0** inside the same **402 × 780 px** scrolling panel; popup page canvas is **402 × 874 px**.
3. On the trigger screen, add a hotspot over the control that opens the popup.
4. **Click or Tap → Show Lightbox → Page** → choose the popup page.
5. Set **Custom size: 442 × 914**, placement **center** (the iPhone 17 frame size).
6. Set the lightbox **background to transparent** (0% opacity). The popup capture already contains the dark scrim and the sheet, so an opaque Axure background would double-darken the screen.

For a centre dialog (`booking-confirmed`, `walkin-success`, `checkin-success`, `checkout-success`, `extension-success`, `guest-request-success`, `sale-recorded`) the capture already shows the modal centred at 362 px wide. The surrounding scrim is part of the image.

Inside a popup:

- **X / Cancel / Done / scrim** → `Hide Lightbox`.
- **A button that navigates** (e.g. *Open folio*, *Proceed to checkout*) → `Hide Lightbox`, then `Open Link → Page` to the target route.
- **Chips** → `Set Selected` style or a small dynamic panel; a toast is enough for the prototype.
- **Form submit** → show the confirmation popup or close with a toast.

Popup → popup chains also exist: `book-room` → `booking-confirmed`, and `guest-request` → `guest-request-success` (the browser prototype opens the success dialog after sending). Wire those as **Show Lightbox** on top of the current one, then **Hide Lightbox** to return.

> Deep link reference: every popup is directly viewable in the browser at `prototype.html?screen=/frontdesk/rooms&modal=room-status`.

---

## 8. Optional masters for a more native build

The captures are pixel-accurate, so native masters are optional. If you want live states for a demo:

| Master | Contents |
|---|---|
| `App Bar` | 390 × 56, back button / moon mark, brand + department tag, bell |
| `Tab Bar` | 390 × 64, 4–5 items, amber active state |
| `Toast` | see §6.4 |
| `Sheet Header` | grab handle, icon, title + subtitle, close |
| `Primary Button` | 40 px amber gradient, 11 px radius |

Use masters only where a state changes (active tab, toast, selected chip). Static screen content should stay as the captured image — it is faster and visually exact.

---

## 9. Preview, QA and publish

1. **Preview (F5)** using the iPhone 17 mobile shell. The phone frame is centered; long pages scroll inside the screen without moving the system chrome.
2. Walk every department from the Portal: Front Desk → all 11 screens; Restaurant, Bar, Kitchen, Manager, Guest, Booking.
3. Click every popup trigger and close it with X, scrim and Cancel/Done.
4. Check that no hotspot links to a missing page — `pages.csv` and `SITEMAP.md` are the source of truth.
5. Spot-check counts: `pages.csv` lists the expected hotspot count per page. If a page is short, you missed rows in its CSV.
6. **Publish → Preview** for a shareable HTML build, or **Publish → Axure Cloud** for a link. If you want the published output inside this deliverable, publish to `prototype/axure/preview/`.

### Acceptance checklist

- [ ] 34 pages, named and grouped exactly as `pages.csv`
- [ ] 19 popup pages in `08 · Popups`
- [ ] 402 × 874 px iPhone screen, centered 442 × 914 px frame, and `#05070F` background
- [ ] 390 px captures centered at X=6 inside the 402 × 780 px scrolling panel
- [ ] 378 hotspots placed (211 open page, 60 show lightbox, 107 toast — see the action summary in `INTERACTION-MATRIX.md`)
- [ ] Every popup closes from X, scrim and its Cancel/Done action
- [ ] No dead links in Preview
- [ ] Centered iPhone 17 shell shows time, cellular, WiFi, battery percentage, Dynamic Island, and gesture bar
- [ ] Pages taller than the phone viewport scroll without moving the device shell or status/gesture bars
- [ ] Status row shows 9:41, cellular, Wi-Fi and battery percentage; Dynamic Island and gesture bar are visible
- [ ] Preview works at phone dimensions and on a phone browser

---

## 10. Worked example — Front Desk Hub (`/frontdesk`)

Image: `screens/frontdesk-09-front-desk-hub.png` (390 × 1156). CSV: `hotspots/frontdesk-09-front-desk-hub.csv` (22 rows).

| Row | Label | Action | Target | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Back to / | nav | `/` | 14 | 11 | 34 | 34 |
| 2 | Notifications | toast | — | 342 | 11 | 34 | 34 |
| 3 | Check in Alice Banda | nav | `/frontdesk/checkin` | 14 | 391 | 362 | 55 |
| 4 | Rooms | nav | `/frontdesk/rooms` | 14 | 811 | 112 | 74 |
| 5 | Folio | nav | `/frontdesk/folio` | 14 | 892 | 112 | 74 |
| 6 | Review room 102 folio | nav | `/frontdesk/folio` | 14 | 975 | 362 | 55 |
| … | … | … | … | … | … | … | … |

Work top-to-bottom: image → rows 1–22 as hot spots → interactions → preview. The same pattern applies to all 34 screens.

---

## 11. Automation tips

- **Copy/paste interactions:** build one `nav` hotspot, then select it, copy, and paste onto the next hotspot; only the target page changes.
- **Multi-edit:** select several hot spots and set size/position together where they share a grid (tab bars are evenly spaced).
- **Tab bar reuse:** the five tab positions are identical on every page in a department. Draw them once per page from the top or bottom of the CSV; the coordinates repeat.
- **Two-monitor workflow:** keep `hotspots/<screen>.csv` open in a spreadsheet and type coordinates into the Inspector. Sorting by `y` matches the visual order.
- **Naming:** use the CSV `label` as the widget name. The Outline panel then reads like the interaction matrix.
- **Batch page creation:** create all 34 blank pages first from `pages.csv`, then import images; Axure keeps the folder order.

---

## 12. Troubleshooting

| Symptom | Cause / fix |
|---|---|
| Hotspots are offset | Keep the image and hotspots inside the same panel; use image X=6 and add 6 to each hotspot X coordinate. |
| Blurry screen | The PNG is @2x; 390 px display width keeps it sharp. Do not upscale beyond 390. |
| Lightbox looks too dark | Set the lightbox background opacity to 0% — the scrim is already in the PNG. |
| Popup opens at the wrong size | Use the centered **442 × 914** iPhone 17 frame; do not use "Fit to browser". |
| Page is cut off | Make sure the full capture is inside the 402 × 780 scrolling panel; use the `height` column in `pages.csv` for content height. |
| Tab bar is far down the page | Expected: captures are full-content and scroll inside the phone viewport. |
| Toast never hides | Add the **Wait 2600 ms → Hide** actions after **Show**. |
| Axure warns about missing pages | A `target` route was mapped to the wrong page name; re-check `SITEMAP.md`. |

---

## 13. File map

| Need | File |
|---|---|
| Axure RP 11 project (deliverable) | `LaLuna-HMS-prototype_RP11.rp` |
| Generator output (RP 9) | `LaLuna-HMS-prototype.rp` |
| Copy installed for delivery | `C:\Users\Lanterns\Documents\Axure\LALUNA\Laluna HMSv3.rp` |
| RP/1 source for regeneration | `rp-build/spec.json` |
| Spec builder | `../tools/build-axure-rp.mjs` |
| All routes & popups | `SITEMAP.md` |
| Every hotspot + coordinates | `hotspot-map.csv` / `INTERACTION-MATRIX.md` |
| Per-screen hotspot list | `hotspots/<screen>.csv` |
| Page size & image height | `pages.csv` |
| Colours, type, components | `STYLE-GUIDE.md` |
| Capture metadata / counts | `manifest.json` |
| Screen images | `screens/*.png` (34) |
| Popup images | `popups/*.png` (19) |

---

## 14. Rebuilding the native `.rp`

The checked-in file is reproducible. Regenerate it whenever `pages.csv`, `hotspot-map.csv`, a screen capture or a popup capture changes.

### 14.1 Prerequisites

- Node.js 18 or newer.
- Python 3.10+ with `Pillow`, `numpy`, `scipy` and `jsonschema` (`pip install -r requirements.txt` inside the zee-rp clone).
- A clone of [zee-rp](https://github.com/LilZeeCN/zee-rp).
- **A genuine Axure RP 9 blank project** used as the container template.

> **Template rule.** The generator only accepts a real RP 9 container (first bytes `AC EF 09 00`). Renaming or byte-flipping an RP 11 file to look like RP 9 is rejected by Axure — it opens with an empty page and zero widgets. Derive the template from a real RP 9 project: strip all pages, panel states and global variables, then verify it reads as `1 page, 0 widgets, 0 panelStates, 0 problems`. Keep it outside the repo (for example `%TEMP%\laluna-rp\blank9-real.rp`) and point `AXURE_BLANK_TEMPLATE_RP9` at it.

### 14.2 Build the spec

```powershell
node prototype/tools/build-axure-rp.mjs
```

This reads `pages.csv` + `hotspot-map.csv` and writes `prototype/axure/rp-build/spec.json`. Each page uses the iPhone 17 device canvas and a phone shell; the 390 px capture and its hotspots are paired inside a fixed 402 × 780 dynamic panel below the status row. The generated RP adds the 442 × 914 frame, status asset, Dynamic Island and gesture bar. Current output: 53 pages, 34 screens, 19 popups, 288 hotspots.

### 14.3 Generate the `.rp`

```powershell
$env:PYTHON3 = 'C:\Python315\python.exe'   # any Python 3 executable
$env:AXURE_BLANK_TEMPLATE_RP9 = "$env:TEMP\laluna-rp\blank9-real.rp"

node <zee-rp>/scripts/proto.mjs check --spec prototype/axure/rp-build/spec.json
node <zee-rp>/scripts/proto.mjs generate `
  --spec prototype/axure/rp-build/spec.json `
  --out-dir prototype/axure/rp-build/out `
  --out-rp prototype/axure/LaLuna-HMS-prototype.rp
```

`generate` refuses to overwrite an existing `.rp`; delete the old file or use a new name. On Windows the zee-rp `scripts/proto.mjs` calls `python3`, which may not exist — set `PYTHON3` and patch the four `spawnSync('python3', …)` calls to `spawnSync(process.env.PYTHON3 || 'python3', …)` if needed.

### 14.4 Verify

```powershell
node <zee-rp>/scripts/rp-read.mjs prototype/axure/LaLuna-HMS-prototype.rp
```

Expected summary: `packages=57 pages=53 widgets=500 interactions=288 actions=288 problems=0`.

`proto.mjs check` reports `structural_ok: true` with whitespace warnings for screenshot-only pages — those are expected and do not affect the native file.

### 14.5 Save an RP 11 copy

1. Open `LaLuna-HMS-prototype.rp` in Axure RP 11.
2. **File → Save As**. Axure writes `LaLuna-HMS-prototype_RP11.rp`; the editor adds the `_RP11` suffix when upgrading a legacy container.
3. Confirm the header is `AC EF 0B 00` and repeat §14.4 against the new file — counts must be identical (`pages=53`, `widgets=500`, `interactions=288`, `problems=0`).
4. Copy the RP 11 file to `C:\Users\Lanterns\Documents\Axure\LALUNA\Laluna HMSv3.rp`.

### 14.6 What the generator wires

| Source | Pages | Interaction |
|---|---|---|
| `nav` rows | 34 screen pages | `Open Link → Page` to the target route |
| `modal` / `tray` rows | 19 popup pages | `Open Link → Page` to the popup page (Axure preview shows it as a page, not a lightbox) |
| popup close hotspot | all 19 popups | `Open Link → Page` back to the parent screen |
| `book-room` confirm | `p02` | `Open Link → Page` to `booking-confirmed` |
| `guest-request` send | `p16` | `Open Link → Page` to `guest-request-success` |
| `toast` rows (90) | — | skipped — no target page |

If you want true Axure **lightboxes** instead of page navigation, open the generated file and convert the 19 popup links to `Show Lightbox → Page` (§7). The images and coordinates are already correct; only the action type changes.

---

*La Luna Hotel & Resort · HMS degree final project · Prototype 1.0*
