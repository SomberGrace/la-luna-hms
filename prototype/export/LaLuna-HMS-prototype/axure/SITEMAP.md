# La Luna HMS · Axure Sitemap

Generated from the live prototype. Every screen below has a matching high-resolution capture in `screens/` and a searchable hotspot list in `hotspot-map.csv`.

- **Screens:** 34
- **Popups:** 19
- **Hotspots:** 378
- **Capture scale:** 2× (390 px logical width → 780 px image width)

## Navigation model

```
Portal (station sign in)
├── Booking web  /book ──► /book/login ──► Guest portal
├── Guest portal /guest ──► dining · folio · services · guide
├── Front desk   /frontdesk ──► rooms · walk-in · check-in · folio · payments
│                              checkout · extension · housekeeping · offline · reports
├── Restaurant   /restaurant ──► menu · dockets · reports
├── Bar          /bar ──► stock · sales · reports
├── Kitchen      /kitchen ──► stock · recipes · reports
└── Manager      /manager ──► audits · reports
```

## Departments

| # | Department | Tag | Screens | Bottom navigation |
|---|---|---|---|---|
| 01 | Portal | Access | 1 | — |
| 02 | Front Desk | Mobile FD | 11 | `/frontdesk` · `/frontdesk/rooms` · `/frontdesk/walkin` · `/frontdesk/folio` · `/frontdesk/checkout` |
| 03 | Restaurant | Restaurant | 4 | `/restaurant` · `/restaurant/menu` · `/restaurant/dockets` · `/restaurant/reports` |
| 04 | Bar | Bar | 4 | `/bar` · `/bar/stock` · `/bar/sales` · `/bar/reports` |
| 05 | Kitchen | Kitchen KDS | 4 | `/kitchen` · `/kitchen/stock` · `/kitchen/recipes` · `/kitchen/reports` |
| 06 | Manager | Executive | 3 | `/manager` · `/manager/audits` · `/manager/reports` |
| 07 | Guest | Resident | 5 | `/guest` · `/guest/dining` · `/guest/folio` · `/guest/services` |
| 08 | Booking | Guest Web | 2 | — |

## Screens

| No | Route | Screen | Dept | Back | Image | Size (px) | Hotspots |
|---|---|---|---|---|---|---|---|
| 1 | `/` | Station Sign In | Portal | — | `screens/portal-01-station-sign-in.png` | 390 × 654 | 7 |
| 2 | `/book` | Book Your Stay | Booking | — | `screens/book-02-book-your-stay.png` | 390 × 959 | 8 |
| 3 | `/book/login` | Guest Account | Booking | `/book` | `screens/book-03-guest-account.png` | 390 × 580 | 5 |
| 4 | `/guest` | Your Stay | Guest | — | `screens/guest-04-your-stay.png` | 390 × 858 | 11 |
| 5 | `/guest/dining` | In-Room Dining | Guest | `/guest` | `screens/guest-05-in-room-dining.png` | 390 × 1032 | 19 |
| 6 | `/guest/folio` | My Room Statement | Guest | `/guest` | `screens/guest-06-my-room-statement.png` | 390 × 688 | 8 |
| 7 | `/guest/services` | Concierge & Requests | Guest | `/guest` | `screens/guest-07-concierge-requests.png` | 390 × 819 | 15 |
| 8 | `/guest/guide` | Resort Guide & Directory | Guest | `/guest` | `screens/guest-08-resort-guide-directory.png` | 390 × 1178 | 19 |
| 9 | `/frontdesk` | Front Desk Hub | Front Desk | — | `screens/frontdesk-09-front-desk-hub.png` | 390 × 1156 | 22 |
| 10 | `/frontdesk/rooms` | Room Inventory | Front Desk | `/frontdesk` | `screens/frontdesk-10-room-inventory.png` | 390 × 859 | 20 |
| 11 | `/frontdesk/walkin` | Rapid Walk-In | Front Desk | `/frontdesk` | `screens/frontdesk-11-rapid-walk-in.png` | 390 × 816 | 8 |
| 12 | `/frontdesk/checkin` | Reservation Check-In | Front Desk | `/frontdesk` | `screens/frontdesk-12-reservation-check-in.png` | 390 × 732 | 11 |
| 13 | `/frontdesk/folio` | Guest Billing Folio | Front Desk | `/frontdesk` | `screens/frontdesk-13-guest-billing-folio.png` | 390 × 818 | 10 |
| 14 | `/frontdesk/payments` | Payment Terminal | Front Desk | `/frontdesk` | `screens/frontdesk-14-payment-terminal.png` | 390 × 715 | 10 |
| 15 | `/frontdesk/checkout` | Check-Out Terminal | Front Desk | `/frontdesk` | `screens/frontdesk-15-check-out-terminal.png` | 390 × 812 | 8 |
| 16 | `/frontdesk/extension` | Stay Extension | Front Desk | `/frontdesk` | `screens/frontdesk-16-stay-extension.png` | 390 × 618 | 8 |
| 17 | `/frontdesk/housekeeping` | Housekeeping | Front Desk | `/frontdesk` | `screens/frontdesk-17-housekeeping.png` | 390 × 784 | 10 |
| 18 | `/frontdesk/offline` | Offline Sync & Queue | Front Desk | `/frontdesk` | `screens/frontdesk-18-offline-sync-queue.png` | 390 × 932 | 11 |
| 19 | `/frontdesk/reports` | Reception Reports | Front Desk | `/frontdesk` | `screens/frontdesk-19-reception-reports.png` | 390 × 907 | 8 |
| 20 | `/restaurant` | Restaurant Hub | Restaurant | — | `screens/restaurant-20-restaurant-hub.png` | 390 × 942 | 14 |
| 21 | `/restaurant/menu` | Menu Items | Restaurant | `/restaurant` | `screens/restaurant-21-menu-items.png` | 390 × 974 | 18 |
| 22 | `/restaurant/dockets` | Dockets | Restaurant | `/restaurant` | `screens/restaurant-22-dockets.png` | 390 × 1028 | 10 |
| 23 | `/restaurant/reports` | Restaurant Report | Restaurant | `/restaurant` | `screens/restaurant-23-restaurant-report.png` | 390 × 891 | 7 |
| 24 | `/bar` | Bar Hub | Bar | — | `screens/bar-24-bar-hub.png` | 390 × 790 | 11 |
| 25 | `/bar/stock` | Stock Audit | Bar | `/bar` | `screens/bar-25-stock-audit.png` | 390 × 889 | 17 |
| 26 | `/bar/sales` | Record Sales | Bar | `/bar` | `screens/bar-26-record-sales.png` | 390 × 823 | 10 |
| 27 | `/bar/reports` | Bar Daily Report | Bar | `/bar` | `screens/bar-27-bar-daily-report.png` | 390 × 858 | 7 |
| 28 | `/kitchen` | Kitchen Orders | Kitchen | — | `screens/kitchen-28-kitchen-orders.png` | 390 × 1044 | 9 |
| 29 | `/kitchen/stock` | Kitchen Stock | Kitchen | `/kitchen` | `screens/kitchen-29-kitchen-stock.png` | 390 × 849 | 13 |
| 30 | `/kitchen/recipes` | Recipe BOM Costing | Kitchen | `/kitchen` | `screens/kitchen-30-recipe-bom-costing.png` | 390 × 795 | 10 |
| 31 | `/kitchen/reports` | Kitchen Shift Report | Kitchen | `/kitchen` | `screens/kitchen-31-kitchen-shift-report.png` | 390 × 720 | 7 |
| 32 | `/manager` | Executive Desk | Manager | — | `screens/manager-32-executive-desk.png` | 390 × 1070 | 13 |
| 33 | `/manager/audits` | Shift Audits & Approvals | Manager | `/manager` | `screens/manager-33-shift-audits-approvals.png` | 390 × 887 | 8 |
| 34 | `/manager/reports` | Consolidated Report | Manager | `/manager` | `screens/manager-34-consolidated-report.png` | 390 × 874 | 6 |

## Popups (lightboxes)

| No | ID | Popup | Type | Opened from | Image | Size (px) |
|---|---|---|---|---|---|---|
| 1 | `book-room` | Book Room 104 | Bottom sheet | `/book` | `popups/popup-01-book-room.png` | 390 × 844 |
| 2 | `booking-confirmed` | Booking Confirmed | Center dialog | `/book` | `popups/popup-02-booking-confirmed.png` | 390 × 844 |
| 3 | `room-status` | Room 104 | Bottom sheet | `/frontdesk/rooms` | `popups/popup-03-room-status.png` | 390 × 844 |
| 4 | `walkin-success` | Registration Complete | Center dialog | `/frontdesk/walkin` | `popups/popup-04-walkin-success.png` | 390 × 844 |
| 5 | `checkin-success` | Check-In Complete | Center dialog | `/frontdesk/checkin` | `popups/popup-05-checkin-success.png` | 390 × 844 |
| 6 | `payment-receipt` | Payment Receipt | Bottom sheet | `/frontdesk/payments` | `popups/popup-06-payment-receipt.png` | 390 × 844 |
| 7 | `checkout-success` | Check-Out Successful | Center dialog | `/frontdesk/checkout` | `popups/popup-07-checkout-success.png` | 390 × 844 |
| 8 | `extension-success` | Stay Extended | Center dialog | `/frontdesk/extension` | `popups/popup-08-extension-success.png` | 390 × 844 |
| 9 | `add-menu-item` | Add Menu Item | Bottom sheet | `/restaurant/menu` | `popups/popup-09-add-menu-item.png` | 390 × 844 |
| 10 | `new-docket` | New Docket | Bottom sheet | `/restaurant/dockets` | `popups/popup-10-new-docket.png` | 390 × 844 |
| 11 | `add-bar-item` | Update Bar Stock | Bottom sheet | `/bar/stock` | `popups/popup-11-add-bar-item.png` | 390 × 844 |
| 12 | `add-ingredient` | Add Raw Ingredient | Bottom sheet | `/kitchen/stock` | `popups/popup-12-add-ingredient.png` | 390 × 844 |
| 13 | `guest-tray` | In-Room Dining Tray | Bottom sheet | `/guest/dining` | `popups/popup-13-guest-tray.png` | 390 × 844 |
| 14 | `guest-request` | Concierge Request | Bottom sheet | `/guest/services` | `popups/popup-14-guest-request.png` | 390 × 844 |
| 15 | `wifi` | Hotel WiFi | Bottom sheet | `/guest` | `popups/popup-15-wifi.png` | 390 × 844 |
| 16 | `guest-request-success` | Request Sent | Center dialog | `/guest/services` | `popups/popup-16-guest-request-success.png` | 390 × 844 |
| 17 | `audit-review` | Review Bar Shift Close | Bottom sheet | `/manager/audits` | `popups/popup-17-audit-review.png` | 390 × 844 |
| 18 | `shift-close` | Close Shift | Bottom sheet | `/frontdesk/reports` | `popups/popup-18-shift-close.png` | 390 × 844 |
| 19 | `sale-recorded` | Sale Recorded | Center dialog | `/bar/sales` | `popups/popup-19-sale-recorded.png` | 390 × 844 |

## Popup reachability

Every popup can be opened directly in the browser prototype with `prototype.html?modal=<id>` (optionally add `&screen=/route` to choose the background). The same list is available in the Screen map drawer inside the prototype.

> Popups that are not linked from a screen (for example `booking-confirmed`) are reached from another popup — they are all listed in the drawer so nothing is hidden.