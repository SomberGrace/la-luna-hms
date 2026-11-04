// Build an RP/1 spec that reproduces the La Luna HMS clickable prototype in Axure.
// Inputs: ../axure/pages.csv, ../axure/hotspot-map.csv, ../axure/screens, ../axure/popups
// Output: ../axure/rp-build/spec.json
//
// Generate the native file with zee-rp:
//   node <zee-rp>/scripts/proto.mjs generate \
//     --spec prototype/axure/rp-build/spec.json \
//     --out-dir prototype/axure/rp-build/out \
//     --out-rp prototype/axure/LaLuna-HMS.rp
// Set AXURE_BLANK_TEMPLATE_RP9 to a real Axure RP 9 blank project.

import { readFileSync, writeFileSync, mkdirSync } from 'node:fs';
import { dirname, join, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const HERE = dirname(fileURLToPath(import.meta.url));
const AXURE = resolve(HERE, '..', 'axure');
const OUT = join(AXURE, 'rp-build');

const parseCsv = text => {
  const rows = [];
  let row = [], field = '', quoted = false;
  for (let i = 0; i < text.length; i++) {
    const ch = text[i];
    if (quoted) {
      if (ch === '"' && text[i + 1] === '"') { field += '"'; i++; }
      else if (ch === '"') quoted = false;
      else field += ch;
    } else if (ch === '"') quoted = true;
    else if (ch === ',') { row.push(field); field = ''; }
    else if (ch === '\n') { row.push(field); rows.push(row); row = []; field = ''; }
    else if (ch !== '\r') field += ch;
  }
  if (field.length || row.length) { row.push(field); rows.push(row); }
  const [head, ...body] = rows.filter(r => r.some(Boolean));
  return body.map(r => Object.fromEntries(head.map((h, i) => [h, r[i] ?? ''])));
};

const pages = parseCsv(readFileSync(join(AXURE, 'pages.csv'), 'utf8'));
const hotspots = parseCsv(readFileSync(join(AXURE, 'hotspot-map.csv'), 'utf8'));

const routeToId = new Map(pages.map((p, i) => [p.route, `s${String(i + 1).padStart(2, '0')}`]));
const screenId = new Map(pages.map((p, i) => [Number(p.no), `s${String(i + 1).padStart(2, '0')}`]));

const POPUPS = [
  ['book-room', 'Book Room', 2],
  ['booking-confirmed', 'Booking Confirmed', 2],
  ['room-status', 'Room Status', 10],
  ['walkin-success', 'Walk-In Checked In', 11],
  ['checkin-success', 'Guest Checked In', 12],
  ['payment-receipt', 'Payment Receipt', 14],
  ['checkout-success', 'Check-Out Complete', 15],
  ['extension-success', 'Stay Extended', 16],
  ['add-menu-item', 'Add Menu Item', 21],
  ['new-docket', 'New Docket', 20],
  ['add-bar-item', 'Add Bar Item', 25],
  ['add-ingredient', 'Add Ingredient', 29],
  ['guest-tray', 'Dining Tray', 5],
  ['guest-request', 'Concierge Request', 7],
  ['wifi', 'WiFi Details', 4],
  ['guest-request-success', 'Request Sent', 7],
  ['audit-review', 'Audit Review', 33],
  ['shift-close', 'Close Shift', 32],
  ['sale-recorded', 'Sale Recorded', 26],
];
const slugToId = new Map(POPUPS.map(([slug], i) => [slug, `p${String(i + 1).padStart(2, '0')}`]));

const PHONE = {
  device: 'iphone-17',
  width: 402,
  height: 874,
  imageInset: 6,
  contentTop: 70,
  contentBottom: 24,
  contentHeight: 780,
};

const box = ({ id, name, x, y, w, h, onClick }) => ({
  type: 'box', id, name,
  x: Math.round(Number(x)), y: Math.round(Number(y)),
  w: Math.round(Number(w)), h: Math.round(Number(h)),
  fill: '#000000', opacity: 0.01, z: 10,
  ...(onClick ? { onClick } : {}),
});
const shot = (src, w, h) => ({ type: 'image', id: 'shot', name: 'Screen', src, x: 0, y: 0, w, h, z: 0 });
const page = (id, name, contentHeight, children) => ({
  id, name,
  device: PHONE.device,
  canvas: { w: PHONE.width, h: PHONE.height },
  root: {
    type: 'phone-shell',
    id: `${id}_phone`,
    device: PHONE.device,
    w: PHONE.width,
    h: PHONE.height,
    padding: [0, 0, PHONE.contentBottom, 0],
    gap: PHONE.contentTop - 28,
    statusBar: false,
    darkIsland: true,
    showHomeIndicator: true,
    children: [
      { type: 'image', id: `${id}_status`, name: 'Phone status bar', src: '../phone-status-bar.png', x: 0, y: 0, w: PHONE.width, h: 28 },
      {
        type: 'dynamic-panel',
        id: `${id}_scroll`,
        name: 'Scrollable screen content',
        w: PHONE.width,
        h: PHONE.contentHeight,
        states: [{
          id: 'main',
          root: {
            type: 'box',
            id: `${id}_content`,
            name: 'Screen content',
            layout: 'stack',
            w: PHONE.width,
            h: Math.max(PHONE.contentHeight, contentHeight),
            children,
          },
        }],
      },
    ],
  },
});

const specPages = pages.map(p => {
  const id = routeToId.get(p.route);
  const children = [{
    ...shot(`../${p.image}`, Number(p.width), Number(p.height)),
    x: PHONE.imageInset,
  }];
  let n = 0;
  for (const row of hotspots.filter(r => Number(r.screen_no) === Number(p.no))) {
    let target;
    if (row.action === 'nav') target = routeToId.get(row.target);
    else if ((row.action === 'modal' || row.action === 'tray') && row.target) target = slugToId.get(row.target);
    if (!target) continue;
    children.push(box({
      id: `${id}_hs${String(++n).padStart(2, '0')}`,
      name: row.label.slice(0, 80),
      x: Number(row.x) + PHONE.imageInset, y: row.y, w: row.w, h: row.h,
      onClick: `goto ${target}`,
    }));
  }
  return page(id, p.title, Number(p.height), children);
});

for (const [slug, label, parentNo] of POPUPS) {
  const id = slugToId.get(slug);
  const parent = pages.find(p => Number(p.no) === Number(parentNo));
  const parentId = parent ? routeToId.get(parent.route) : 's01';
  const children = [{
    ...shot(`../popups/popup-${String(POPUPS.findIndex(p => p[0] === slug) + 1).padStart(2, '0')}-${slug}.png`, 390, 844),
    x: PHONE.imageInset,
  }];
  // Close (X) in the modal header.
  children.push(box({ id: `${id}_close`, name: 'Close', x: 334 + PHONE.imageInset, y: 404, w: 48, h: 48, onClick: `goto ${parentId}` }));
  if (slug === 'book-room') {
    children.push(box({ id: `${id}_confirm`, name: 'Confirm booking', x: 30 + PHONE.imageInset, y: 722, w: 330, h: 46, onClick: 'goto p02' }));
    children.push(box({ id: `${id}_cancel`, name: 'Cancel', x: 20 + PHONE.imageInset, y: 800, w: 350, h: 44, onClick: `goto ${parentId}` }));
  }
  if (slug === 'guest-request') {
    children.push(box({ id: `${id}_send`, name: 'Send request', x: 30 + PHONE.imageInset, y: 678, w: 330, h: 44, onClick: 'goto p16' }));
    children.push(box({ id: `${id}_cancel`, name: 'Cancel', x: 20 + PHONE.imageInset, y: 800, w: 350, h: 44, onClick: `goto ${parentId}` }));
  }
  specPages.push(page(id, label, 844, children));
}

const spec = {
  language: 'rp/1',
  target: 9,
  name: 'La Luna HMS Prototype',
  theme: 'modern-dark',
  device: 'iphone-17',
  domain: 'mobile',
  shell: 'mobile-stack',
  static: false,
  pages: specPages,
};

mkdirSync(OUT, { recursive: true });
writeFileSync(join(OUT, 'spec.json'), JSON.stringify(spec, null, 2), 'utf8');
const countClicks = value => {
  if (Array.isArray(value)) return value.reduce((n, item) => n + countClicks(item), 0);
  if (!value || typeof value !== 'object') return 0;
  return Number(Boolean(value.onClick))
    + Object.values(value).reduce((n, item) => n + countClicks(item), 0);
};
console.log(JSON.stringify({ pages: specPages.length, screens: pages.length, popups: POPUPS.length,
  hotspots: countClicks(specPages),
  out: join(OUT, 'spec.json') }, null, 2));
