# La Luna HMS · Interaction Matrix

Every clickable target captured from the running prototype, grouped by screen. Coordinates are in logical pixels relative to the top-left of the screen image (0,0). Multiply by 2 for the pixel position in the @2x PNG.

Total: **378 hotspots** across **34 screens**. The same data is available as CSV/JSON in this folder.

## Action summary

| Action | Count | Axure interaction |
|---|---|---|
| Open page | 211 | On Click → Open Page |
| Toast / confirm | 107 | On Click → Show toast master / Set Text |
| Show lightbox | 60 | On Click → Show Lightbox |

## 1. / — Station Sign In

**Department:** Portal · **Image:** `screens/portal-01-station-sign-in.png` · **Hotspots:** 7

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Open Front Desk | Open page | `/frontdesk` | 31 | 234 | 161 | 46 |
| 2 | Open Restaurant | Open page | `/restaurant` | 199 | 234 | 161 | 46 |
| 3 | Open Bar | Open page | `/bar` | 31 | 287 | 161 | 46 |
| 4 | Open Kitchen | Open page | `/kitchen` | 199 | 287 | 161 | 46 |
| 5 | Open Manager | Open page | `/manager` | 31 | 340 | 161 | 46 |
| 6 | Open Guest Booking | Open page | `/book` | 199 | 340 | 161 | 46 |
| 7 | Sign in to front desk | Open page | `/frontdesk` | 31 | 534 | 328 | 40 |

## 2. /book — Book Your Stay

**Department:** Booking · **Image:** `screens/book-02-book-your-stay.png` · **Hotspots:** 8

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Notifications | Toast / confirm | — | 342 | 11 | 34 | 34 |
| 2 | Check availability | Toast / confirm | — | 26 | 425 | 338 | 40 |
| 3 | Book room 101 | Show lightbox | `book-room` | 14 | 515 | 177 | 80 |
| 4 | Book room 104 | Show lightbox | `book-room` | 199 | 515 | 177 | 80 |
| 5 | Book room 201 | Show lightbox | `book-room` | 14 | 603 | 177 | 80 |
| 6 | Book room 301 | Show lightbox | `book-room` | 199 | 603 | 177 | 80 |
| 7 | Book room 401 | Show lightbox | `book-room` | 14 | 691 | 177 | 80 |
| 8 | Guest portal login | Open page | `/book/login` | 15 | 881 | 360 | 54 |

## 3. /book/login — Guest Account

**Department:** Booking · **Image:** `screens/book-03-guest-account.png` · **Hotspots:** 5

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Back to /book | Open page | `/book` | 14 | 11 | 34 | 34 |
| 2 | Notifications | Toast / confirm | — | 342 | 11 | 34 | 34 |
| 3 | Sign in | Toast / confirm | — | 14 | 230 | 51 | 25 |
| 4 | Create account | Toast / confirm | — | 71 | 230 | 89 | 25 |
| 5 | Enter guest portal | Open page | `/guest` | 28 | 412 | 334 | 40 |

## 4. /guest — Your Stay

**Department:** Guest · **Image:** `screens/guest-04-your-stay.png` · **Hotspots:** 11

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Notifications | Toast / confirm | — | 342 | 11 | 34 | 34 |
| 2 | In-room dining | Open page | `/guest/dining` | 14 | 421 | 177 | 92 |
| 3 | My bill | Open page | `/guest/folio` | 199 | 421 | 177 | 92 |
| 4 | Concierge | Open page | `/guest/services` | 14 | 521 | 177 | 92 |
| 5 | Resort guide | Open page | `/guest/guide` | 199 | 521 | 177 | 92 |
| 6 | Extra bath towels Requested 18:42 · Housekeeping On the way | Toast / confirm | — | 15 | 652 | 360 | 54 |
| 7 | View WiFi details | Show lightbox | `wifi` | 14 | 719 | 362 | 52 |
| 8 | Stay | Open page | `/guest` | 8 | 800 | 92 | 48 |
| 9 | Dining | Open page | `/guest/dining` | 102 | 800 | 92 | 48 |
| 10 | My Bill | Open page | `/guest/folio` | 196 | 800 | 92 | 48 |
| 11 | Services | Open page | `/guest/services` | 290 | 800 | 92 | 48 |

## 5. /guest/dining — In-Room Dining

**Department:** Guest · **Image:** `screens/guest-05-in-room-dining.png` · **Hotspots:** 19

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Back to /guest | Open page | `/guest` | 14 | 11 | 34 | 34 |
| 2 | Notifications | Toast / confirm | — | 342 | 11 | 34 | 34 |
| 3 | All | Toast / confirm | — | 14 | 230 | 33 | 25 |
| 4 | Mains | Toast / confirm | — | 53 | 230 | 48 | 25 |
| 5 | Breakfast | Toast / confirm | — | 108 | 230 | 64 | 25 |
| 6 | Snacks | Toast / confirm | — | 177 | 230 | 54 | 25 |
| 7 | Drinks | Toast / confirm | — | 237 | 230 | 50 | 25 |
| 8 | Add Nsima | Show lightbox | — | 331 | 297 | 32 | 32 |
| 9 | Add Tilapia | Show lightbox | — | 331 | 392 | 32 | 32 |
| 10 | Add Tenderloin | Show lightbox | — | 331 | 487 | 32 | 32 |
| 11 | Add Mandasi | Show lightbox | — | 331 | 576 | 32 | 32 |
| 12 | Add Juice | Show lightbox | — | 331 | 658 | 32 | 32 |
| 13 | Add Gin & Tonic | Show lightbox | — | 331 | 741 | 32 | 32 |
| 14 | Review dining tray | Show lightbox | `guest-tray` | 14 | 808 | 362 | 44 |
| 15 | Tilapia & Mango Juice Yesterday 19:20 · Delivered MWK 6,000 | Toast / confirm | — | 15 | 891 | 360 | 54 |
| 16 | Stay | Open page | `/guest` | 8 | 975 | 92 | 48 |
| 17 | Dining | Open page | `/guest/dining` | 102 | 975 | 92 | 48 |
| 18 | My Bill | Open page | `/guest/folio` | 196 | 975 | 92 | 48 |
| 19 | Services | Open page | `/guest/services` | 290 | 975 | 92 | 48 |

## 6. /guest/folio — My Room Statement

**Department:** Guest · **Image:** `screens/guest-06-my-room-statement.png` · **Hotspots:** 8

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Back to /guest | Open page | `/guest` | 14 | 11 | 34 | 34 |
| 2 | Notifications | Toast / confirm | — | 342 | 11 | 34 | 34 |
| 3 | Query a charge | Show lightbox | `guest-request` | 15 | 492 | 360 | 55 |
| 4 | Email a copy Sent to alice.banda@example.mw Send | Toast / confirm | — | 15 | 547 | 360 | 54 |
| 5 | Stay | Open page | `/guest` | 8 | 631 | 92 | 48 |
| 6 | Dining | Open page | `/guest/dining` | 102 | 631 | 92 | 48 |
| 7 | My Bill | Open page | `/guest/folio` | 196 | 631 | 92 | 48 |
| 8 | Services | Open page | `/guest/services` | 290 | 631 | 92 | 48 |

## 7. /guest/services — Concierge & Requests

**Department:** Guest · **Image:** `screens/guest-07-concierge-requests.png` · **Hotspots:** 15

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Back to /guest | Open page | `/guest` | 14 | 11 | 34 | 34 |
| 2 | Notifications | Toast / confirm | — | 342 | 11 | 34 | 34 |
| 3 | Request bath towels | Show lightbox | `guest-request` | 14 | 255 | 116 | 74 |
| 4 | Request room clean | Show lightbox | `guest-request` | 137 | 255 | 116 | 74 |
| 5 | Request drinking water | Show lightbox | `guest-request` | 260 | 255 | 116 | 74 |
| 6 | Request late check-out | Show lightbox | `guest-request` | 14 | 336 | 116 | 74 |
| 7 | Request airport transfer | Show lightbox | `guest-request` | 137 | 336 | 116 | 74 |
| 8 | Request first aid | Show lightbox | `guest-request` | 260 | 336 | 116 | 74 |
| 9 | Send custom request | Show lightbox | `guest-request` | 26 | 532 | 338 | 40 |
| 10 | Extra bath towels 18:42 · Housekeeping On the way | Toast / confirm | — | 15 | 623 | 360 | 55 |
| 11 | Two bottles of water Yesterday 16:10 · Delivered Done | Toast / confirm | — | 15 | 678 | 360 | 54 |
| 12 | Stay | Open page | `/guest` | 8 | 762 | 92 | 48 |
| 13 | Dining | Open page | `/guest/dining` | 102 | 762 | 92 | 48 |
| 14 | My Bill | Open page | `/guest/folio` | 196 | 762 | 92 | 48 |
| 15 | Services | Open page | `/guest/services` | 290 | 762 | 92 | 48 |

## 8. /guest/guide — Resort Guide & Directory

**Department:** Guest · **Image:** `screens/guest-08-resort-guide-directory.png` · **Hotspots:** 19

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Back to /guest | Open page | `/guest` | 14 | 11 | 34 | 34 |
| 2 | Notifications | Toast / confirm | — | 342 | 11 | 34 | 34 |
| 3 | In-room dining Daily · 07:00 – 22:00 Ext. 2 | Toast / confirm | — | 15 | 256 | 360 | 55 |
| 4 | Restaurant Breakfast 07:00 · Dinner 18:30 Ext. 2 | Toast / confirm | — | 15 | 311 | 360 | 55 |
| 5 | Poolside bar Daily · 11:00 – 23:00 Ext. 1 | Toast / confirm | — | 15 | 366 | 360 | 55 |
| 6 | Pool & beach Open 06:00 – 19:00 Ext. 0 | Toast / confirm | — | 15 | 421 | 360 | 54 |
| 7 | Front Desk Reception Housekeeping, billing and general help Ext. 0 | Toast / confirm | — | 15 | 515 | 360 | 55 |
| 8 | In-Room Dining Kitchen Orders and dietary questions Ext. 2 | Toast / confirm | — | 15 | 570 | 360 | 55 |
| 9 | Housekeeping & Linen Towels, laundry and turndown Ext. 3 | Toast / confirm | — | 15 | 625 | 360 | 55 |
| 10 | General Manager On Duty Feedback and escalations Ext. 4 | Toast / confirm | — | 15 | 680 | 360 | 55 |
| 11 | Emergency & First Aid Urgent medical assistance Ext. 9 | Toast / confirm | — | 15 | 735 | 360 | 54 |
| 12 | Sunset cruise 17:00 · 90 min · MWK 45,000 | Toast / confirm | — | 14 | 828 | 177 | 92 |
| 13 | Snorkel safari 09:00 · 2 hrs · MWK 38,000 | Toast / confirm | — | 199 | 828 | 177 | 92 |
| 14 | Village walk 14:00 · 2 hrs · MWK 20,000 | Toast / confirm | — | 14 | 928 | 177 | 92 |
| 15 | Birdwatching 06:30 · 2 hrs · MWK 25,000 | Toast / confirm | — | 199 | 928 | 177 | 92 |
| 16 | Stay | Open page | `/guest` | 8 | 1120 | 92 | 48 |
| 17 | Dining | Open page | `/guest/dining` | 102 | 1120 | 92 | 48 |
| 18 | My Bill | Open page | `/guest/folio` | 196 | 1120 | 92 | 48 |
| 19 | Services | Open page | `/guest/services` | 290 | 1120 | 92 | 48 |

## 9. /frontdesk — Front Desk Hub

**Department:** Front Desk · **Image:** `screens/frontdesk-09-front-desk-hub.png` · **Hotspots:** 22

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Notifications | Toast / confirm | — | 342 | 11 | 34 | 34 |
| 2 | Check in Alice Banda | Open page | `/frontdesk/checkin` | 15 | 407 | 360 | 55 |
| 3 | Check in Thandiwe Phiri | Open page | `/frontdesk/checkin` | 15 | 462 | 360 | 55 |
| 4 | Check in James Mwale | Open page | `/frontdesk/checkin` | 15 | 517 | 360 | 54 |
| 5 | View all rooms | Open page | `/frontdesk/rooms` | 300 | 589 | 74 | 15 |
| 6 | Open folio for room 102 | Open page | `/frontdesk/folio` | 15 | 615 | 360 | 55 |
| 7 | Open folio for room 205 | Open page | `/frontdesk/folio` | 15 | 670 | 360 | 54 |
| 8 | Rooms | Open page | `/frontdesk/rooms` | 14 | 763 | 116 | 74 |
| 9 | Walk-in | Open page | `/frontdesk/walkin` | 137 | 763 | 116 | 74 |
| 10 | Folio | Open page | `/frontdesk/folio` | 260 | 763 | 116 | 74 |
| 11 | Checkout | Open page | `/frontdesk/checkout` | 14 | 844 | 116 | 74 |
| 12 | Payments | Open page | `/frontdesk/payments` | 137 | 844 | 116 | 74 |
| 13 | Housekeeping | Open page | `/frontdesk/housekeeping` | 260 | 844 | 116 | 74 |
| 14 | Stay extension | Open page | `/frontdesk/extension` | 14 | 925 | 116 | 74 |
| 15 | Offline sync | Open page | `/frontdesk/offline` | 137 | 925 | 116 | 74 |
| 16 | Reception reports | Open page | `/frontdesk/reports` | 260 | 925 | 116 | 74 |
| 17 | Review room 102 folio | Open page | `/frontdesk/folio` | 14 | 1011 | 362 | 59 |
| 18 | Hub | Open page | `/frontdesk` | 8 | 1098 | 73 | 48 |
| 19 | Rooms | Open page | `/frontdesk/rooms` | 83 | 1098 | 73 | 48 |
| 20 | Walk-In | Open page | `/frontdesk/walkin` | 158 | 1098 | 73 | 48 |
| 21 | Folio | Open page | `/frontdesk/folio` | 234 | 1098 | 73 | 48 |
| 22 | Checkout | Open page | `/frontdesk/checkout` | 309 | 1098 | 73 | 48 |

## 10. /frontdesk/rooms — Room Inventory

**Department:** Front Desk · **Image:** `screens/frontdesk-10-room-inventory.png` · **Hotspots:** 20

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Back to /frontdesk | Open page | `/frontdesk` | 14 | 11 | 34 | 34 |
| 2 | Notifications | Toast / confirm | — | 342 | 11 | 34 | 34 |
| 3 | All 12 | Toast / confirm | — | 14 | 230 | 45 | 25 |
| 4 | Available 2 | Toast / confirm | — | 65 | 230 | 69 | 25 |
| 5 | Occupied 9 | Toast / confirm | — | 140 | 230 | 72 | 25 |
| 6 | Out of service 1 | Toast / confirm | — | 218 | 230 | 89 | 25 |
| 7 | Change status for room 101 | Show lightbox | `room-status` | 14 | 269 | 177 | 80 |
| 8 | Change status for room 102 | Show lightbox | `room-status` | 199 | 269 | 177 | 80 |
| 9 | Change status for room 103 | Show lightbox | `room-status` | 14 | 357 | 177 | 80 |
| 10 | Change status for room 104 | Show lightbox | `room-status` | 199 | 357 | 177 | 80 |
| 11 | Change status for room 201 | Show lightbox | `room-status` | 14 | 445 | 177 | 80 |
| 12 | Change status for room 205 | Show lightbox | `room-status` | 199 | 445 | 177 | 80 |
| 13 | Change status for room 301 | Show lightbox | `room-status` | 14 | 533 | 177 | 80 |
| 14 | Change status for room 304 | Show lightbox | `room-status` | 199 | 533 | 177 | 80 |
| 15 | Change status for room 401 | Show lightbox | `room-status` | 14 | 621 | 177 | 80 |
| 16 | Hub | Open page | `/frontdesk` | 8 | 801 | 73 | 48 |
| 17 | Rooms | Open page | `/frontdesk/rooms` | 83 | 801 | 73 | 48 |
| 18 | Walk-In | Open page | `/frontdesk/walkin` | 158 | 801 | 73 | 48 |
| 19 | Folio | Open page | `/frontdesk/folio` | 234 | 801 | 73 | 48 |
| 20 | Checkout | Open page | `/frontdesk/checkout` | 309 | 801 | 73 | 48 |

## 11. /frontdesk/walkin — Rapid Walk-In

**Department:** Front Desk · **Image:** `screens/frontdesk-11-rapid-walk-in.png` · **Hotspots:** 8

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Back to /frontdesk | Open page | `/frontdesk` | 14 | 11 | 34 | 34 |
| 2 | Notifications | Toast / confirm | — | 342 | 11 | 34 | 34 |
| 3 | Complete walk-in check-in | Show lightbox | `walkin-success` | 28 | 584 | 334 | 40 |
| 4 | Hub | Open page | `/frontdesk` | 8 | 758 | 73 | 48 |
| 5 | Rooms | Open page | `/frontdesk/rooms` | 83 | 758 | 73 | 48 |
| 6 | Walk-In | Open page | `/frontdesk/walkin` | 158 | 758 | 73 | 48 |
| 7 | Folio | Open page | `/frontdesk/folio` | 234 | 758 | 73 | 48 |
| 8 | Checkout | Open page | `/frontdesk/checkout` | 309 | 758 | 73 | 48 |

## 12. /frontdesk/checkin — Reservation Check-In

**Department:** Front Desk · **Image:** `screens/frontdesk-12-reservation-check-in.png` · **Hotspots:** 11

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Back to /frontdesk | Open page | `/frontdesk` | 14 | 11 | 34 | 34 |
| 2 | Notifications | Toast / confirm | — | 342 | 11 | 34 | 34 |
| 3 | Alice Banda | Toast / confirm | — | 14 | 230 | 73 | 25 |
| 4 | Thandiwe Phiri | Toast / confirm | — | 93 | 230 | 87 | 25 |
| 5 | James Mwale | Toast / confirm | — | 186 | 230 | 82 | 25 |
| 6 | Confirm check-in | Show lightbox | `checkin-success` | 28 | 571 | 334 | 40 |
| 7 | Hub | Open page | `/frontdesk` | 8 | 675 | 73 | 48 |
| 8 | Rooms | Open page | `/frontdesk/rooms` | 83 | 675 | 73 | 48 |
| 9 | Walk-In | Open page | `/frontdesk/walkin` | 158 | 675 | 73 | 48 |
| 10 | Folio | Open page | `/frontdesk/folio` | 234 | 675 | 73 | 48 |
| 11 | Checkout | Open page | `/frontdesk/checkout` | 309 | 675 | 73 | 48 |

## 13. /frontdesk/folio — Guest Billing Folio

**Department:** Front Desk · **Image:** `screens/frontdesk-13-guest-billing-folio.png` · **Hotspots:** 10

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Back to /frontdesk | Open page | `/frontdesk` | 14 | 11 | 34 | 34 |
| 2 | Notifications | Toast / confirm | — | 342 | 11 | 34 | 34 |
| 3 | Post payment | Open page | `/frontdesk/payments` | 14 | 471 | 116 | 74 |
| 4 | Extend stay | Open page | `/frontdesk/extension` | 137 | 471 | 116 | 74 |
| 5 | Check out room 104 | Open page | `/frontdesk/checkout` | 260 | 471 | 116 | 74 |
| 6 | Hub | Open page | `/frontdesk` | 8 | 761 | 73 | 48 |
| 7 | Rooms | Open page | `/frontdesk/rooms` | 83 | 761 | 73 | 48 |
| 8 | Walk-In | Open page | `/frontdesk/walkin` | 158 | 761 | 73 | 48 |
| 9 | Folio | Open page | `/frontdesk/folio` | 234 | 761 | 73 | 48 |
| 10 | Checkout | Open page | `/frontdesk/checkout` | 309 | 761 | 73 | 48 |

## 14. /frontdesk/payments — Payment Terminal

**Department:** Front Desk · **Image:** `screens/frontdesk-14-payment-terminal.png` · **Hotspots:** 10

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Back to /frontdesk | Open page | `/frontdesk` | 14 | 11 | 34 | 34 |
| 2 | Notifications | Toast / confirm | — | 342 | 11 | 34 | 34 |
| 3 | Post payment | Show lightbox | `payment-receipt` | 28 | 425 | 334 | 40 |
| 4 | View receipt | Show lightbox | `payment-receipt` | 15 | 518 | 360 | 55 |
| 5 | View receipt | Show lightbox | `payment-receipt` | 15 | 573 | 360 | 54 |
| 6 | Hub | Open page | `/frontdesk` | 8 | 657 | 73 | 48 |
| 7 | Rooms | Open page | `/frontdesk/rooms` | 83 | 657 | 73 | 48 |
| 8 | Walk-In | Open page | `/frontdesk/walkin` | 158 | 657 | 73 | 48 |
| 9 | Folio | Open page | `/frontdesk/folio` | 234 | 657 | 73 | 48 |
| 10 | Checkout | Open page | `/frontdesk/checkout` | 309 | 657 | 73 | 48 |

## 15. /frontdesk/checkout — Check-Out Terminal

**Department:** Front Desk · **Image:** `screens/frontdesk-15-check-out-terminal.png` · **Hotspots:** 8

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Back to /frontdesk | Open page | `/frontdesk` | 14 | 11 | 34 | 34 |
| 2 | Notifications | Toast / confirm | — | 342 | 11 | 34 | 34 |
| 3 | Complete check-out | Show lightbox | `checkout-success` | 26 | 653 | 338 | 40 |
| 4 | Hub | Open page | `/frontdesk` | 8 | 755 | 73 | 48 |
| 5 | Rooms | Open page | `/frontdesk/rooms` | 83 | 755 | 73 | 48 |
| 6 | Walk-In | Open page | `/frontdesk/walkin` | 158 | 755 | 73 | 48 |
| 7 | Folio | Open page | `/frontdesk/folio` | 234 | 755 | 73 | 48 |
| 8 | Checkout | Open page | `/frontdesk/checkout` | 309 | 755 | 73 | 48 |

## 16. /frontdesk/extension — Stay Extension

**Department:** Front Desk · **Image:** `screens/frontdesk-16-stay-extension.png` · **Hotspots:** 8

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Back to /frontdesk | Open page | `/frontdesk` | 14 | 11 | 34 | 34 |
| 2 | Notifications | Toast / confirm | — | 342 | 11 | 34 | 34 |
| 3 | Confirm extension | Show lightbox | `extension-success` | 28 | 374 | 334 | 40 |
| 4 | Hub | Open page | `/frontdesk` | 8 | 561 | 73 | 48 |
| 5 | Rooms | Open page | `/frontdesk/rooms` | 83 | 561 | 73 | 48 |
| 6 | Walk-In | Open page | `/frontdesk/walkin` | 158 | 561 | 73 | 48 |
| 7 | Folio | Open page | `/frontdesk/folio` | 234 | 561 | 73 | 48 |
| 8 | Checkout | Open page | `/frontdesk/checkout` | 309 | 561 | 73 | 48 |

## 17. /frontdesk/housekeeping — Housekeeping

**Department:** Front Desk · **Image:** `screens/frontdesk-17-housekeeping.png` · **Hotspots:** 10

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Back to /frontdesk | Open page | `/frontdesk` | 14 | 11 | 34 | 34 |
| 2 | Notifications | Toast / confirm | — | 342 | 11 | 34 | 34 |
| 3 | Room 102 Dirty · Guest departed 10:00 Mark clean | Toast / confirm | — | 15 | 422 | 360 | 55 |
| 4 | Room 205 Dirty · Stayover service Mark clean | Toast / confirm | — | 15 | 477 | 360 | 55 |
| 5 | Room 201 Inspect · Supervisor assigned Pass | Toast / confirm | — | 15 | 532 | 360 | 55 |
| 6 | Hub | Open page | `/frontdesk` | 8 | 726 | 73 | 48 |
| 7 | Rooms | Open page | `/frontdesk/rooms` | 83 | 726 | 73 | 48 |
| 8 | Walk-In | Open page | `/frontdesk/walkin` | 158 | 726 | 73 | 48 |
| 9 | Folio | Open page | `/frontdesk/folio` | 234 | 726 | 73 | 48 |
| 10 | Checkout | Open page | `/frontdesk/checkout` | 309 | 726 | 73 | 48 |

## 18. /frontdesk/offline — Offline Sync & Queue

**Department:** Front Desk · **Image:** `screens/frontdesk-18-offline-sync-queue.png` · **Hotspots:** 11

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Back to /frontdesk | Open page | `/frontdesk` | 14 | 11 | 34 | 34 |
| 2 | Notifications | Toast / confirm | — | 342 | 11 | 34 | 34 |
| 3 | Walk-in check-in · Memory Tembo Room 103 · 08:01 · queued Retry | Toast / confirm | — | 15 | 422 | 360 | 55 |
| 4 | Folio charge · Room 102 Minibar · MWK 3,500 · 07:58 Retry | Toast / confirm | — | 15 | 477 | 360 | 55 |
| 5 | Cash payment · Room 205 MWK 18,000 · 07:54 Retry | Toast / confirm | — | 15 | 532 | 360 | 54 |
| 6 | Queue replays automatically Keep this screen open while the connection | Toast / confirm | — | 14 | 786 | 362 | 59 |
| 7 | Hub | Open page | `/frontdesk` | 8 | 874 | 73 | 48 |
| 8 | Rooms | Open page | `/frontdesk/rooms` | 83 | 874 | 73 | 48 |
| 9 | Walk-In | Open page | `/frontdesk/walkin` | 158 | 874 | 73 | 48 |
| 10 | Folio | Open page | `/frontdesk/folio` | 234 | 874 | 73 | 48 |
| 11 | Checkout | Open page | `/frontdesk/checkout` | 309 | 874 | 73 | 48 |

## 19. /frontdesk/reports — Reception Reports

**Department:** Front Desk · **Image:** `screens/frontdesk-19-reception-reports.png` · **Hotspots:** 8

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Back to /frontdesk | Open page | `/frontdesk` | 14 | 11 | 34 | 34 |
| 2 | Notifications | Toast / confirm | — | 342 | 11 | 34 | 34 |
| 3 | Close shift | Show lightbox | `shift-close` | 14 | 761 | 362 | 59 |
| 4 | Hub | Open page | `/frontdesk` | 8 | 849 | 73 | 48 |
| 5 | Rooms | Open page | `/frontdesk/rooms` | 83 | 849 | 73 | 48 |
| 6 | Walk-In | Open page | `/frontdesk/walkin` | 158 | 849 | 73 | 48 |
| 7 | Folio | Open page | `/frontdesk/folio` | 234 | 849 | 73 | 48 |
| 8 | Checkout | Open page | `/frontdesk/checkout` | 309 | 849 | 73 | 48 |

## 20. /restaurant — Restaurant Hub

**Department:** Restaurant · **Image:** `screens/restaurant-20-restaurant-hub.png` · **Hotspots:** 14

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Notifications | Toast / confirm | — | 342 | 11 | 34 | 34 |
| 2 | View all dockets | Open page | `/restaurant/dockets` | 319 | 400 | 55 | 15 |
| 3 | Open docket 0184 | Open page | `/restaurant/dockets` | 15 | 426 | 360 | 55 |
| 4 | Open docket 0183 | Open page | `/restaurant/dockets` | 15 | 481 | 360 | 55 |
| 5 | Open docket 0182 | Open page | `/restaurant/dockets` | 15 | 536 | 360 | 54 |
| 6 | New docket | Show lightbox | `new-docket` | 14 | 629 | 116 | 74 |
| 7 | Menu | Open page | `/restaurant/menu` | 137 | 629 | 116 | 74 |
| 8 | Dockets | Open page | `/restaurant/dockets` | 260 | 629 | 116 | 74 |
| 9 | Restaurant reports | Open page | `/restaurant/reports` | 14 | 710 | 116 | 74 |
| 10 | Open kitchen stock | Open page | `/kitchen/stock` | 14 | 796 | 362 | 59 |
| 11 | Hub | Open page | `/restaurant` | 8 | 884 | 92 | 48 |
| 12 | Menu | Open page | `/restaurant/menu` | 102 | 884 | 92 | 48 |
| 13 | Dockets | Open page | `/restaurant/dockets` | 196 | 884 | 92 | 48 |
| 14 | Reports | Open page | `/restaurant/reports` | 290 | 884 | 92 | 48 |

## 21. /restaurant/menu — Menu Items

**Department:** Restaurant · **Image:** `screens/restaurant-21-menu-items.png` · **Hotspots:** 18

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Back to /restaurant | Open page | `/restaurant` | 14 | 11 | 34 | 34 |
| 2 | Notifications | Toast / confirm | — | 342 | 11 | 34 | 34 |
| 3 | All 9 | Toast / confirm | — | 14 | 230 | 41 | 25 |
| 4 | Mains | Toast / confirm | — | 61 | 230 | 48 | 25 |
| 5 | Breakfast | Toast / confirm | — | 116 | 230 | 64 | 25 |
| 6 | Snacks | Toast / confirm | — | 185 | 230 | 54 | 25 |
| 7 | Drinks | Toast / confirm | — | 245 | 230 | 50 | 25 |
| 8 | Add Nsima with Stewed Chicken | Toast / confirm | — | 331 | 290 | 32 | 32 |
| 9 | Add Grilled Lake Tilapia (Chambo) | Toast / confirm | — | 331 | 373 | 32 | 32 |
| 10 | Add Prime Beef Tenderloin | Toast / confirm | — | 331 | 455 | 32 | 32 |
| 11 | Add Full English Resort Breakfast | Toast / confirm | — | 331 | 538 | 32 | 32 |
| 12 | Add Warm Mandasi & Local Honey | Toast / confirm | — | 331 | 620 | 32 | 32 |
| 13 | Add Malawi Gin & Tonic | Toast / confirm | — | 331 | 703 | 32 | 32 |
| 14 | Add menu item | Show lightbox | `add-menu-item` | 26 | 835 | 338 | 40 |
| 15 | Hub | Open page | `/restaurant` | 8 | 916 | 92 | 48 |
| 16 | Menu | Open page | `/restaurant/menu` | 102 | 916 | 92 | 48 |
| 17 | Dockets | Open page | `/restaurant/dockets` | 196 | 916 | 92 | 48 |
| 18 | Reports | Open page | `/restaurant/reports` | 290 | 916 | 92 | 48 |

## 22. /restaurant/dockets — Dockets

**Department:** Restaurant · **Image:** `screens/restaurant-22-dockets.png` · **Hotspots:** 10

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Back to /restaurant | Open page | `/restaurant` | 14 | 11 | 34 | 34 |
| 2 | Notifications | Toast / confirm | — | 342 | 11 | 34 | 34 |
| 3 | New docket | Show lightbox | `new-docket` | 300 | 400 | 74 | 15 |
| 4 | Bump | Toast / confirm | — | 300 | 529 | 63 | 27 |
| 5 | Bump | Toast / confirm | — | 300 | 680 | 63 | 27 |
| 6 | Bump | Toast / confirm | — | 300 | 831 | 63 | 27 |
| 7 | Hub | Open page | `/restaurant` | 8 | 970 | 92 | 48 |
| 8 | Menu | Open page | `/restaurant/menu` | 102 | 970 | 92 | 48 |
| 9 | Dockets | Open page | `/restaurant/dockets` | 196 | 970 | 92 | 48 |
| 10 | Reports | Open page | `/restaurant/reports` | 290 | 970 | 92 | 48 |

## 23. /restaurant/reports — Restaurant Report

**Department:** Restaurant · **Image:** `screens/restaurant-23-restaurant-report.png` · **Hotspots:** 7

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Back to /restaurant | Open page | `/restaurant` | 14 | 11 | 34 | 34 |
| 2 | Notifications | Toast / confirm | — | 342 | 11 | 34 | 34 |
| 3 | Send to manager Restaurant totals feed the manager consolidation repor | Toast / confirm | — | 14 | 746 | 362 | 59 |
| 4 | Hub | Open page | `/restaurant` | 8 | 833 | 92 | 48 |
| 5 | Menu | Open page | `/restaurant/menu` | 102 | 833 | 92 | 48 |
| 6 | Dockets | Open page | `/restaurant/dockets` | 196 | 833 | 92 | 48 |
| 7 | Reports | Open page | `/restaurant/reports` | 290 | 833 | 92 | 48 |

## 24. /bar — Bar Hub

**Department:** Bar · **Image:** `screens/bar-24-bar-hub.png` · **Hotspots:** 11

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Notifications | Toast / confirm | — | 342 | 11 | 34 | 34 |
| 2 | Manage bar stock | Open page | `/bar/stock` | 289 | 385 | 85 | 15 |
| 3 | Carlsberg Green Lager 500ml 8 bottles left · reorder at 12 Restock | Toast / confirm | — | 15 | 411 | 360 | 55 |
| 4 | Malawi Gin 750ml 3 bottles left · reorder at 4 Restock | Toast / confirm | — | 15 | 466 | 360 | 54 |
| 5 | Record bar sale | Open page | `/bar/sales` | 14 | 559 | 116 | 74 |
| 6 | Stock audit | Open page | `/bar/stock` | 137 | 559 | 116 | 74 |
| 7 | Bar reports | Open page | `/bar/reports` | 260 | 559 | 116 | 74 |
| 8 | Hub | Open page | `/bar` | 8 | 732 | 92 | 48 |
| 9 | Stock | Open page | `/bar/stock` | 102 | 732 | 92 | 48 |
| 10 | Sales | Open page | `/bar/sales` | 196 | 732 | 92 | 48 |
| 11 | Reports | Open page | `/bar/reports` | 290 | 732 | 92 | 48 |

## 25. /bar/stock — Stock Audit

**Department:** Bar · **Image:** `screens/bar-25-stock-audit.png` · **Hotspots:** 17

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Back to /bar | Open page | `/bar` | 14 | 11 | 34 | 34 |
| 2 | Notifications | Toast / confirm | — | 342 | 11 | 34 | 34 |
| 3 | All items | Toast / confirm | — | 14 | 230 | 59 | 25 |
| 4 | Low stock | Toast / confirm | — | 79 | 230 | 66 | 25 |
| 5 | Spirits | Toast / confirm | — | 151 | 230 | 50 | 25 |
| 6 | Beer | Toast / confirm | — | 207 | 230 | 42 | 25 |
| 7 | Update lager stock | Show lightbox | `add-bar-item` | 15 | 270 | 360 | 55 |
| 8 | Update gin stock | Show lightbox | `add-bar-item` | 15 | 325 | 360 | 55 |
| 9 | Update wine stock | Show lightbox | `add-bar-item` | 15 | 380 | 360 | 55 |
| 10 | Update stout stock | Show lightbox | `add-bar-item` | 15 | 435 | 360 | 55 |
| 11 | Update mixer stock | Show lightbox | `add-bar-item` | 15 | 490 | 360 | 55 |
| 12 | Update tonic stock | Show lightbox | `add-bar-item` | 15 | 545 | 360 | 54 |
| 13 | Add bar item | Show lightbox | `add-bar-item` | 26 | 679 | 338 | 40 |
| 14 | Hub | Open page | `/bar` | 8 | 831 | 92 | 48 |
| 15 | Stock | Open page | `/bar/stock` | 102 | 831 | 92 | 48 |
| 16 | Sales | Open page | `/bar/sales` | 196 | 831 | 92 | 48 |
| 17 | Reports | Open page | `/bar/reports` | 290 | 831 | 92 | 48 |

## 26. /bar/sales — Record Sales

**Department:** Bar · **Image:** `screens/bar-26-record-sales.png` · **Hotspots:** 10

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Back to /bar | Open page | `/bar` | 14 | 11 | 34 | 34 |
| 2 | Notifications | Toast / confirm | — | 342 | 11 | 34 | 34 |
| 3 | Record sale | Show lightbox | `sale-recorded` | 28 | 478 | 334 | 40 |
| 4 | 2 × Guinness Stout 12:51 · Room 205 MWK 5,000 | Toast / confirm | — | 15 | 571 | 360 | 55 |
| 5 | 1 × Gin & Tonic 12:38 · Cash MWK 3,500 | Toast / confirm | — | 15 | 626 | 360 | 55 |
| 6 | 2 × Sauvignon Blanc 12:20 · Card MWK 8,000 | Toast / confirm | — | 15 | 681 | 360 | 54 |
| 7 | Hub | Open page | `/bar` | 8 | 765 | 92 | 48 |
| 8 | Stock | Open page | `/bar/stock` | 102 | 765 | 92 | 48 |
| 9 | Sales | Open page | `/bar/sales` | 196 | 765 | 92 | 48 |
| 10 | Reports | Open page | `/bar/reports` | 290 | 765 | 92 | 48 |

## 27. /bar/reports — Bar Daily Report

**Department:** Bar · **Image:** `screens/bar-27-bar-daily-report.png` · **Hotspots:** 7

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Back to /bar | Open page | `/bar` | 14 | 11 | 34 | 34 |
| 2 | Notifications | Toast / confirm | — | 342 | 11 | 34 | 34 |
| 3 | Manager sign-off Variance is within tolerance. Send the ledger to the  | Toast / confirm | — | 14 | 713 | 362 | 59 |
| 4 | Hub | Open page | `/bar` | 8 | 800 | 92 | 48 |
| 5 | Stock | Open page | `/bar/stock` | 102 | 800 | 92 | 48 |
| 6 | Sales | Open page | `/bar/sales` | 196 | 800 | 92 | 48 |
| 7 | Reports | Open page | `/bar/reports` | 290 | 800 | 92 | 48 |

## 28. /kitchen — Kitchen Orders

**Department:** Kitchen · **Image:** `screens/kitchen-28-kitchen-orders.png` · **Hotspots:** 9

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Notifications | Toast / confirm | — | 342 | 11 | 34 | 34 |
| 2 | Bump | Toast / confirm | — | 300 | 535 | 63 | 27 |
| 3 | Bump | Toast / confirm | — | 300 | 696 | 63 | 27 |
| 4 | Bump | Toast / confirm | — | 300 | 847 | 63 | 27 |
| 5 | Open kitchen stock | Open page | `/kitchen/stock` | 14 | 898 | 362 | 59 |
| 6 | Orders | Open page | `/kitchen` | 8 | 986 | 92 | 48 |
| 7 | Stock | Open page | `/kitchen/stock` | 102 | 986 | 92 | 48 |
| 8 | Recipes | Open page | `/kitchen/recipes` | 196 | 986 | 92 | 48 |
| 9 | Reports | Open page | `/kitchen/reports` | 290 | 986 | 92 | 48 |

## 29. /kitchen/stock — Kitchen Stock

**Department:** Kitchen · **Image:** `screens/kitchen-29-kitchen-stock.png` · **Hotspots:** 13

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Back to /kitchen | Open page | `/kitchen` | 14 | 11 | 34 | 34 |
| 2 | Notifications | Toast / confirm | — | 342 | 11 | 34 | 34 |
| 3 | Update maize flour | Show lightbox | `add-ingredient` | 15 | 230 | 360 | 55 |
| 4 | Update tilapia | Show lightbox | `add-ingredient` | 15 | 285 | 360 | 55 |
| 5 | Update chicken | Show lightbox | `add-ingredient` | 15 | 340 | 360 | 55 |
| 6 | Update cooking oil | Show lightbox | `add-ingredient` | 15 | 395 | 360 | 55 |
| 7 | Update potatoes | Show lightbox | `add-ingredient` | 15 | 450 | 360 | 55 |
| 8 | Update tomatoes | Show lightbox | `add-ingredient` | 15 | 505 | 360 | 54 |
| 9 | Add raw ingredient | Show lightbox | `add-ingredient` | 26 | 639 | 338 | 40 |
| 10 | Orders | Open page | `/kitchen` | 8 | 791 | 92 | 48 |
| 11 | Stock | Open page | `/kitchen/stock` | 102 | 791 | 92 | 48 |
| 12 | Recipes | Open page | `/kitchen/recipes` | 196 | 791 | 92 | 48 |
| 13 | Reports | Open page | `/kitchen/reports` | 290 | 791 | 92 | 48 |

## 30. /kitchen/recipes — Recipe BOM Costing

**Department:** Kitchen · **Image:** `screens/kitchen-30-recipe-bom-costing.png` · **Hotspots:** 10

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Back to /kitchen | Open page | `/kitchen` | 14 | 11 | 34 | 34 |
| 2 | Notifications | Toast / confirm | — | 342 | 11 | 34 | 34 |
| 3 | Nsima & Stewed Chicken Flour 0.35 kg · chicken 0.25 bird · tomato 0.15 | Toast / confirm | — | 15 | 230 | 360 | 55 |
| 4 | Grilled Chambo Tilapia 1 pc · oil & spices 0.08 L · salad 0.1 kg Cost  | Toast / confirm | — | 15 | 285 | 360 | 55 |
| 5 | Crispy Chips Potatoes 0.45 kg · oil 0.1 L · sauce 0.05 L Cost 780 | Toast / confirm | — | 15 | 340 | 360 | 55 |
| 6 | Beef Tenderloin Beef 0.35 kg · potatoes 0.3 kg · butter & herbs 0.05 k | Toast / confirm | — | 15 | 395 | 360 | 54 |
| 7 | Orders | Open page | `/kitchen` | 8 | 737 | 92 | 48 |
| 8 | Stock | Open page | `/kitchen/stock` | 102 | 737 | 92 | 48 |
| 9 | Recipes | Open page | `/kitchen/recipes` | 196 | 737 | 92 | 48 |
| 10 | Reports | Open page | `/kitchen/reports` | 290 | 737 | 92 | 48 |

## 31. /kitchen/reports — Kitchen Shift Report

**Department:** Kitchen · **Image:** `screens/kitchen-31-kitchen-shift-report.png` · **Hotspots:** 7

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Back to /kitchen | Open page | `/kitchen` | 14 | 11 | 34 | 34 |
| 2 | Notifications | Toast / confirm | — | 342 | 11 | 34 | 34 |
| 3 | Handover Send the closing count to the evening chef and the manager co | Toast / confirm | — | 14 | 574 | 362 | 59 |
| 4 | Orders | Open page | `/kitchen` | 8 | 662 | 92 | 48 |
| 5 | Stock | Open page | `/kitchen/stock` | 102 | 662 | 92 | 48 |
| 6 | Recipes | Open page | `/kitchen/recipes` | 196 | 662 | 92 | 48 |
| 7 | Reports | Open page | `/kitchen/reports` | 290 | 662 | 92 | 48 |

## 32. /manager — Executive Desk

**Department:** Manager · **Image:** `screens/manager-32-executive-desk.png` · **Hotspots:** 13

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Notifications | Toast / confirm | — | 342 | 11 | 34 | 34 |
| 2 | Open consolidated report | Open page | `/manager/reports` | 293 | 554 | 81 | 15 |
| 3 | Front desk report | Open page | `/manager/reports` | 15 | 580 | 360 | 55 |
| 4 | Restaurant report | Open page | `/manager/reports` | 15 | 635 | 360 | 55 |
| 5 | Bar report | Open page | `/manager/reports` | 15 | 690 | 360 | 55 |
| 6 | Kitchen report | Open page | `/manager/reports` | 15 | 745 | 360 | 54 |
| 7 | Shift audits | Open page | `/manager/audits` | 14 | 838 | 116 | 74 |
| 8 | Consolidated report | Open page | `/manager/reports` | 137 | 838 | 116 | 74 |
| 9 | Add shift note | Show lightbox | `shift-close` | 260 | 838 | 116 | 74 |
| 10 | Review pending approvals | Open page | `/manager/audits` | 14 | 924 | 362 | 59 |
| 11 | Overview | Open page | `/manager` | 8 | 1012 | 123 | 48 |
| 12 | Audits | Open page | `/manager/audits` | 133 | 1012 | 123 | 48 |
| 13 | Consolidated | Open page | `/manager/reports` | 259 | 1012 | 123 | 48 |

## 33. /manager/audits — Shift Audits & Approvals

**Department:** Manager · **Image:** `screens/manager-33-shift-audits-approvals.png` · **Hotspots:** 8

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Back to /manager | Open page | `/manager` | 14 | 11 | 34 | 34 |
| 2 | Notifications | Toast / confirm | — | 342 | 11 | 34 | 34 |
| 3 | Review bar shift close | Show lightbox | `audit-review` | 15 | 422 | 360 | 55 |
| 4 | Review folio adjustment | Show lightbox | `audit-review` | 15 | 477 | 360 | 55 |
| 5 | Review stock variance | Show lightbox | `audit-review` | 15 | 532 | 360 | 54 |
| 6 | Overview | Open page | `/manager` | 8 | 829 | 123 | 48 |
| 7 | Audits | Open page | `/manager/audits` | 133 | 829 | 123 | 48 |
| 8 | Consolidated | Open page | `/manager/reports` | 259 | 829 | 123 | 48 |

## 34. /manager/reports — Consolidated Report

**Department:** Manager · **Image:** `screens/manager-34-consolidated-report.png` · **Hotspots:** 6

| # | Target | Action | Destination | X | Y | W | H |
|---|---|---|---|---|---|---|---|
| 1 | Back to /manager | Open page | `/manager` | 14 | 11 | 34 | 34 |
| 2 | Notifications | Toast / confirm | — | 342 | 11 | 34 | 34 |
| 3 | Board pack Export the consolidated report as a signed PDF for the week | Toast / confirm | — | 14 | 728 | 362 | 59 |
| 4 | Overview | Open page | `/manager` | 8 | 816 | 123 | 48 |
| 5 | Audits | Open page | `/manager/audits` | 133 | 816 | 123 | 48 |
| 6 | Consolidated | Open page | `/manager/reports` | 259 | 816 | 123 | 48 |

## Popup interactions

Popup actions are data-driven and identical in the browser prototype. Use the table when wiring lightbox buttons in Axure.

| Popup | Button / target | Action | Destination |
|---|---|---|---|
| `book-room` | Confirm booking | undefined | — |
| `book-room` | Cancel | Close overlay | — |
| `booking-confirmed` | Open guest portal | Open page | `/book/login` |
| `booking-confirmed` | Book another room | Close overlay | — |
| `room-status` | Status set to available | Toast / confirm | — |
| `room-status` | Status set to occupied | Toast / confirm | — |
| `room-status` | Status set to dirty | Toast / confirm | — |
| `room-status` | Status set to out of service | Toast / confirm | — |
| `room-status` | Save status | undefined | — |
| `room-status` | Cancel | Close overlay | — |
| `walkin-success` | Open folio | Open page | `/frontdesk/folio` |
| `walkin-success` | New walk-in | Close overlay | — |
| `checkin-success` | Open folio | Open page | `/frontdesk/folio` |
| `checkin-success` | Back to hub | Open page | `/frontdesk` |
| `payment-receipt` | Receipt emailed to the guest | Toast / confirm | — |
| `payment-receipt` | Proceed to checkout | Open page | `/frontdesk/checkout` |
| `payment-receipt` | Done | Close overlay | — |
| `checkout-success` | Open housekeeping | Open page | `/frontdesk/housekeeping` |
| `checkout-success` | Back to hub | Open page | `/frontdesk` |
| `extension-success` | Back to folio | Open page | `/frontdesk/folio` |
| `extension-success` | Done | Close overlay | — |
| `add-menu-item` | Add to menu | undefined | — |
| `add-menu-item` | Cancel | Close overlay | — |
| `new-docket` | Fire to kitchen | undefined | — |
| `new-docket` | Cancel | Close overlay | — |
| `add-bar-item` | Save stock | undefined | — |
| `add-bar-item` | Cancel | Close overlay | — |
| `add-ingredient` | Save ingredient | undefined | — |
| `add-ingredient` | Cancel | Close overlay | — |
| `guest-tray` | Place order | undefined | — |
| `guest-tray` | Add more | Close overlay | — |
| `guest-request` | Request type: bath towels | Toast / confirm | — |
| `guest-request` | Request type: room clean | Toast / confirm | — |
| `guest-request` | Request type: drinking water | Toast / confirm | — |
| `guest-request` | Request type: late check-out | Toast / confirm | — |
| `guest-request` | Send request | undefined | — |
| `guest-request` | Cancel | Close overlay | — |
| `wifi` | Copy password | Toast / confirm | — |
| `wifi` | Done | Close overlay | — |
| `guest-request-success` | Done | Close overlay | — |
| `audit-review` | Approve & sign | undefined | — |
| `audit-review` | Reject | Close overlay | — |
| `audit-review` | Cancel | Close overlay | — |
| `shift-close` | Close and hand over | undefined | — |
| `shift-close` | Cancel | Close overlay | — |
| `sale-recorded` | Done | Close overlay | — |
