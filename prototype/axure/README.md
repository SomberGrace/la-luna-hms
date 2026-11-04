# La Luna HMS · Axure Kit

Everything needed to open, review or rebuild the clickable Axure prototype.

**Ready to open:** `LaLuna-HMS-prototype_RP11.rp` — Axure-saved RP 11 project with 53 pages (34 screens + 19 popups), 500 widgets, 288 interactions. It uses a centered iPhone 17 shell with a phone status bar, Dynamic Island, gesture bar, and a scrollable viewport for long screens. Installed for delivery at `C:\Users\Lanterns\Documents\Axure\LALUNA\Laluna HMSv3.rp`.

`LaLuna-HMS-prototype.rp` is the generator output (RP 9 container) kept as the reproducible build artifact. Open either file in Axure RP 11; both render identically.

| | |
|---|---|
| Screens | **34** @2x PNG captures (`screens/`, 780 px wide) |
| Popups | **19** @2x PNG captures (`popups/`, 780 × 1688) |
| Phone chrome | iPhone 17 frame, status icons/time, Dynamic Island, gesture bar, scroll panels |
| Hotspots | **378** with exact coordinates (`hotspot-map.csv`, `hotspots/`) |
| Departments | Portal · Booking · Guest · Front Desk · Restaurant · Bar · Kitchen · Manager |

## Start here

1. **[AXURE-BUILD-GUIDE.md](AXURE-BUILD-GUIDE.md)** — step-by-step build, interactions, lightboxes, QA.
2. **[SITEMAP.md](SITEMAP.md)** — all pages, routes, popups and navigation model.
3. **[INTERACTION-MATRIX.md](INTERACTION-MATRIX.md)** — every hotspot grouped by screen.
4. **[STYLE-GUIDE.md](STYLE-GUIDE.md)** — palette, typography and components.
5. **[pages.csv](pages.csv)** — page name, route, image, height, hotspot count.
6. **[hotspot-map.csv](hotspot-map.csv)** — global coordinate map (and `hotspots/` per screen).

## Native file

`LaLuna-HMS-prototype_RP11.rp` is the Axure-saved RP 11 project: 53 pages, 500 widgets, 288 interactions, 0 problems. It opens directly in Axure RP 11 with every screen image and hotspot already placed. The installed delivery copy is `C:\Users\Lanterns\Documents\Axure\LALUNA\Laluna HMSv3.rp`.

`LaLuna-HMS-prototype.rp` is the generator's RP 9 output kept for traceability. Axure RP 11 produced the RP 11 copy from it with **Save As**; the two files have identical page, widget and interaction counts.

The file was generated from [`rp-build/spec.json`](rp-build/spec.json), which is rebuilt from `pages.csv` + `hotspot-map.csv` by `../tools/build-axure-rp.mjs`. See [AXURE-BUILD-GUIDE.md §14](AXURE-BUILD-GUIDE.md#14-rebuilding-the-native-rp) for the regeneration command.

`toast` hotspots (90 rows with no target page) are intentionally skipped in the native build. The 288 real navigation, modal and tray hotspots are all wired.

The browser prototype in the parent folder (`../index.html`) remains the zero-dependency way to demo the same flows.
