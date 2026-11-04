# La Luna HMS · Axure RP Build Guide

This kit turns the working HTML prototype into an Axure RP project with the same 34 screens, 19 popups and 378 clickable targets. Everything is captured at 2× and every coordinate is recorded, so the Axure build is a mechanical assembly pass rather than a redesign.

---

## 0. Read this first — about the `.rp` file

Axure RP project files (`.rp`) are a **proprietary packed binary format** written by the Axure editor. They cannot be authored reliably by a script or by hand, and fabricating one would produce a file Axure refuses to open. This kit therefore ships:

1. **A fully clickable browser prototype** — the working deliverable you can present today.
2. **Every asset needed to assemble the `.rp`** — 34 screen captures, 19 popup captures, a hotspot map with exact coordinates, and this guide.

Assembling the `.rp` in Axure RP 11 is a repeatable 2–4 hour pass. Nothing has to be designed from scratch.

> **Do not overwrite** `C:\Users\Lanterns\Documents\Axure\LALUNA\Laluna HMSv2.rp`. Create a new project file next to it, for example `LaLuna-HMS-Prototype.rp`.

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
2. Set the page width to **390 px** and the background to **`#05070F`**.
3. Create the folder structure from §4.
4. Drag `axure/screens/portal-01-station-sign-in.png` onto the Home page, set width **390**, position **0, 0**, lock it.
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
| Page width | **390 px** |
| Page height | the `height` column in `pages.csv` (the image is a full-content capture) |
| Background colour | **`#05070F`** |
| Page alignment | left / top |
| Scrollbars | as needed (most screens are taller than 844 px) |

Tip: select all pages in the Pages panel → right-click → **Page Style** to set the background in one go. Set the height per page afterwards.

---

## 5. Import the screens

For each row in `pages.csv`:

1. Drag `axure/screens/<image>` onto the page.
2. Set **W = 390** and keep the aspect ratio (the PNG is 780 px wide, so it lands at 390 × `height`). If you drop the image and it is 780 wide, set W to 390 — **never leave it at 780**.
3. Set **X = 0, Y = 0**.
4. Lock the image (right-click → **Lock** or the lock icon in the Inspector) so hot spots stay on top.
5. Set the page height to the image height (`pages.csv`).
6. Add the route and a one-line purpose to **Page Notes**.

The captures include the app bar, content and tab bar. The tab bar sits at the end of a long page (the capture is full-content, not a clipped 844 viewport). That is intentional: every control on the screen is visible and clickable in the Axure preview.

> Optional: for screens taller than 844 px you can instead use the capture as a scrolling image inside a 390 × 844 dynamic panel. Most builds do not need this — the full-height page is simpler and keeps every hotspot reachable.

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

Draw a **Hot Spot** widget for each row at exactly `x, y` with size `w × h`. Because the image is placed at 390 px wide, the coordinates map 1:1. Name the widget after the label (e.g. `Open Front Desk`) so the Outline panel is readable.

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
2. Import each `popups/popup-XX-*.png` at **W = 390**, **X = 0, Y = 0**; set page height **844**.
3. On the trigger screen, add a hotspot over the control that opens the popup.
4. **Click or Tap → Show Lightbox → Page** → choose the popup page.
5. Set **Custom size: 390 × 844**, placement **center**.
6. Set the lightbox **background to transparent** (0% opacity). The popup capture already contains the dark scrim and the sheet, so an opaque Axure background would double-darken the screen.

For a centre dialog (`booking-confirmed`, `walkin-success`, `checkin-success`, `checkout-success`, `extension-success`, `guest-request-success`, `sale-recorded`) the capture already shows the modal centred at 362 px wide. The lightbox page stays 390 × 844; the surrounding scrim is part of the image.

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

1. **Preview (F5)** at 390 px width (Axure's mobile preset).
2. Walk every department from the Portal: Front Desk → all 11 screens; Restaurant, Bar, Kitchen, Manager, Guest, Booking.
3. Click every popup trigger and close it with X, scrim and Cancel/Done.
4. Check that no hotspot links to a missing page — `pages.csv` and `SITEMAP.md` are the source of truth.
5. Spot-check counts: `pages.csv` lists the expected hotspot count per page. If a page is short, you missed rows in its CSV.
6. **Publish → Preview** for a shareable HTML build, or **Publish → Axure Cloud** for a link. If you want the published output inside this deliverable, publish to `prototype/axure/preview/`.

### Acceptance checklist

- [ ] 34 pages, named and grouped exactly as `pages.csv`
- [ ] 19 popup pages in `08 · Popups`
- [ ] 390 px page width, `#05070F` background, images at x=0/y=0
- [ ] 378 hotspots placed (211 open page, 60 show lightbox, 107 toast — see the action summary in `INTERACTION-MATRIX.md`)
- [ ] Every popup closes from X, scrim and its Cancel/Done action
- [ ] No dead links in Preview
- [ ] Preview works at 390 × 844 and on a phone browser

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
| Hotspots are offset | The image was placed at 780 px wide or not at 0,0. Set W=390, X=0, Y=0. |
| Blurry screen | The PNG is @2x; 390 px display width keeps it sharp. Do not upscale beyond 390. |
| Lightbox looks too dark | Set the lightbox background opacity to 0% — the scrim is already in the PNG. |
| Popup opens at the wrong size | Use **Custom size 390 × 844**; do not use "Fit to browser". |
| Page is cut off | Page height is smaller than the image; use the `height` column in `pages.csv`. |
| Tab bar is far down the page | Expected: captures are full-content. Use the hotspot coordinates, or wrap the image in a 390 × 844 scrolling panel. |
| Toast never hides | Add the **Wait 2600 ms → Hide** actions after **Show**. |
| Axure warns about missing pages | A `target` route was mapped to the wrong page name; re-check `SITEMAP.md`. |

---

## 13. File map

| Need | File |
|---|---|
| All routes & popups | `SITEMAP.md` |
| Every hotspot + coordinates | `hotspot-map.csv` / `INTERACTION-MATRIX.md` |
| Per-screen hotspot list | `hotspots/<screen>.csv` |
| Page size & image height | `pages.csv` |
| Colours, type, components | `STYLE-GUIDE.md` |
| Capture metadata / counts | `manifest.json` |
| Screen images | `screens/*.png` (34) |
| Popup images | `popups/*.png` (19) |

---

*La Luna Hotel & Resort · HMS degree final project · Prototype 1.0*
