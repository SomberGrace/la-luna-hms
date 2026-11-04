/* Generate the Axure documentation from the live prototype data + capture manifest. */
import fs from 'node:fs';
import path from 'node:path';
import vm from 'node:vm';
import { fileURLToPath } from 'node:url';

const here = path.dirname(fileURLToPath(import.meta.url));
const root = path.resolve(here, '..');
const axure = path.join(root, 'axure');

const window = {};
const context = vm.createContext({ window, console });
for (const f of ['data.js', 'data-screens-guest.js', 'data-screens-frontdesk.js', 'data-screens-ops.js', 'data-modals.js']) {
  vm.runInContext(fs.readFileSync(path.join(root, 'assets', 'js', f), 'utf8'), context, { filename: f });
}
const P = window.PROTO;
const manifest = JSON.parse(fs.readFileSync(path.join(axure, 'manifest.json'), 'utf8'));
const hotspots = JSON.parse(fs.readFileSync(path.join(axure, 'hotspot-map.json'), 'utf8'));

const byRoute = new Map();
for (const h of hotspots) {
  if (!byRoute.has(h.route)) byRoute.set(h.route, []);
  byRoute.get(h.route).push(h);
}
const screenByRoute = new Map(manifest.screens.map((s) => [s.route, s]));
const popupById = new Map(manifest.popups.map((p) => [p.id, p]));

const deptOrder = ['portal', 'frontdesk', 'restaurant', 'bar', 'kitchen', 'manager', 'guest', 'book'];
const deptDigits = { portal: '01', frontdesk: '02', restaurant: '03', bar: '04', kitchen: '05', manager: '06', guest: '07', book: '08' };

const cell = (v) => String(v == null ? '' : v).replace(/\|/g, '\\|').replace(/\n/g, ' ');
const fmtAction = (a) => ({
  nav: 'Open page',
  modal: 'Show lightbox',
  close: 'Close overlay',
  toast: 'Toast / confirm',
  back: 'Back',
  tray: 'Show lightbox',
}[a] || a);
const dest = (h) => h.target || (h.action === 'close' ? '—' : '—');

/* ---------------- SITEMAP ---------------- */
const deptScreens = (id) => Object.entries(P.screens).filter(([, s]) => s.dept === id);
const sitemap = [];
sitemap.push('# La Luna HMS · Axure Sitemap', '');
sitemap.push('Generated from the live prototype. Every screen below has a matching high-resolution capture in `screens/` and a searchable hotspot list in `hotspot-map.csv`.', '');
sitemap.push(`- **Screens:** ${manifest.screens.length}`);
sitemap.push(`- **Popups:** ${manifest.popups.length}`);
sitemap.push(`- **Hotspots:** ${manifest.hotspots}`);
sitemap.push(`- **Capture scale:** 2× (390 px logical width → 780 px image width)`);
sitemap.push('', '## Navigation model', '');
sitemap.push('```');
sitemap.push('Portal (station sign in)');
sitemap.push('├── Booking web  /book ──► /book/login ──► Guest portal');
sitemap.push('├── Guest portal /guest ──► dining · folio · services · guide');
sitemap.push('├── Front desk   /frontdesk ──► rooms · walk-in · check-in · folio · payments');
sitemap.push('│                              checkout · extension · housekeeping · offline · reports');
sitemap.push('├── Restaurant   /restaurant ──► menu · dockets · reports');
sitemap.push('├── Bar          /bar ──► stock · sales · reports');
sitemap.push('├── Kitchen      /kitchen ──► stock · recipes · reports');
sitemap.push('└── Manager      /manager ──► audits · reports');
sitemap.push('```', '');

sitemap.push('## Departments', '');
sitemap.push('| # | Department | Tag | Screens | Bottom navigation |');
sitemap.push('|---|---|---|---|---|');
for (const id of deptOrder) {
  const d = P.depts[id];
  const screens = deptScreens(id);
  const nav = (d.nav || []).map((n) => `\`${n.route}\``).join(' · ') || '—';
  sitemap.push(`| ${deptDigits[id]} | ${cell(d.label)} | ${cell(d.tag)} | ${screens.length} | ${cell(nav)} |`);
}
sitemap.push('');

sitemap.push('## Screens', '');
sitemap.push('| No | Route | Screen | Dept | Back | Image | Size (px) | Hotspots |');
sitemap.push('|---|---|---|---|---|---|---|---|');
for (const s of manifest.screens) {
  const def = P.screens[s.route];
  sitemap.push(`| ${s.no} | \`${cell(s.route)}\` | ${cell(s.title)} | ${cell(P.depts[s.dept].label)} | ${def.back ? '`' + def.back + '`' : '—'} | \`screens/${s.file.split('/').pop()}\` | ${s.width / 2} × ${s.height / 2} | ${s.hotspots} |`);
}
sitemap.push('');

sitemap.push('## Popups (lightboxes)', '');
sitemap.push('| No | ID | Popup | Type | Opened from | Image | Size (px) |');
sitemap.push('|---|---|---|---|---|---|---|');
const popupBackdrop = {
  'book-room': '/book', 'booking-confirmed': '/book', 'room-status': '/frontdesk/rooms',
  'walkin-success': '/frontdesk/walkin', 'checkin-success': '/frontdesk/checkin',
  'payment-receipt': '/frontdesk/payments', 'checkout-success': '/frontdesk/checkout',
  'extension-success': '/frontdesk/extension', 'add-menu-item': '/restaurant/menu',
  'new-docket': '/restaurant/dockets', 'add-bar-item': '/bar/stock', 'add-ingredient': '/kitchen/stock',
  'guest-tray': '/guest/dining', 'guest-request': '/guest/services', wifi: '/guest',
  'guest-request-success': '/guest/services', 'audit-review': '/manager/audits',
  'shift-close': '/frontdesk/reports', 'sale-recorded': '/bar/sales',
};
for (const p of manifest.popups) {
  const m = P.modals[p.id];
  const type = m.size === 'center' ? 'Center dialog' : 'Bottom sheet';
  sitemap.push(`| ${p.no} | \`${p.id}\` | ${cell(m.title)} | ${type} | \`${popupBackdrop[p.id]}\` | \`popups/${p.file.split('/').pop()}\` | 390 × 844 |`);
}
sitemap.push('');
sitemap.push('## Popup reachability', '');
sitemap.push('Every popup can be opened directly in the browser prototype with `prototype.html?modal=<id>` (optionally add `&screen=/route` to choose the background). The same list is available in the Screen map drawer inside the prototype.', '');
sitemap.push('> Popups that are not linked from a screen (for example `booking-confirmed`) are reached from another popup — they are all listed in the drawer so nothing is hidden.');

/* ---------------- INTERACTION MATRIX ---------------- */
const matrix = [];
matrix.push('# La Luna HMS · Interaction Matrix', '');
matrix.push('Every clickable target captured from the running prototype, grouped by screen. Coordinates are in logical pixels relative to the top-left of the screen image (0,0). Multiply by 2 for the pixel position in the @2x PNG.', '');
matrix.push(`Total: **${manifest.hotspots} hotspots** across **${manifest.screens.length} screens**. The same data is available as CSV/JSON in this folder.`, '');

matrix.push('## Action summary', '');
const actionAxure = {
  'Open page': 'On Click → Open Page',
  'Show lightbox': 'On Click → Show Lightbox',
  'Close overlay': 'On Click → Hide Lightbox',
  'Back': 'On Click → Open Page (back target)',
  'Toast / confirm': 'On Click → Show toast master / Set Text',
};
const counts = {};
for (const h of hotspots) {
  const label = fmtAction(h.action);
  counts[label] = (counts[label] || 0) + 1;
}
matrix.push('| Action | Count | Axure interaction |');
matrix.push('|---|---|---|');
for (const [label, v] of Object.entries(counts)) {
  matrix.push(`| ${label} | ${v} | ${actionAxure[label] || '—'} |`);
}
matrix.push('');

for (const s of manifest.screens) {
  const def = P.screens[s.route];
  const list = byRoute.get(s.route) || [];
  matrix.push(`## ${s.no}. ${s.route} — ${s.title}`, '');
  matrix.push(`**Department:** ${P.depts[s.dept].label} · **Image:** \`screens/${s.file.split('/').pop()}\` · **Hotspots:** ${list.length}`, '');
  if (!list.length) {
    matrix.push('_No interactive targets._', '');
    continue;
  }
  matrix.push('| # | Target | Action | Destination | X | Y | W | H |');
  matrix.push('|---|---|---|---|---|---|---|---|');
  for (const h of list) {
    matrix.push(`| ${h.index} | ${cell(h.label)} | ${fmtAction(h.action)} | ${h.target ? '`' + cell(h.target) + '`' : dest(h)} | ${h.x} | ${h.y} | ${h.w} | ${h.h} |`);
  }
  matrix.push('');
}

matrix.push('## Popup interactions', '');
matrix.push('Popup actions are data-driven and identical in the browser prototype. Use the table when wiring lightbox buttons in Axure.', '');
matrix.push('| Popup | Button / target | Action | Destination |');
matrix.push('|---|---|---|---|');
const walk = (blocks, onAction) => {
  for (const b of blocks || []) {
    if (b.action) onAction(b.action);
    if (b.submit) onAction(b.submit);
    for (const item of b.items || []) if (item.action) onAction(item.action);
    if (b.action && b.action.action) onAction(b.action.action);
  }
};
for (const [id, m] of Object.entries(P.modals)) {
  const acts = [];
  walk(m.blocks, (a) => acts.push({ label: a.label || a.text || 'Block action', action: a }));
  for (const a of m.actions || []) acts.push({ label: a.label, action: a.action });
  for (const { label, action } of acts) {
    matrix.push(`| \`${id}\` | ${cell(label)} | ${fmtAction(action.do)} | ${action.to ? '`' + cell(action.to) + '`' : '—'} |`);
  }
}
matrix.push('');

/* ---------------- write ---------------- */
const csvEsc = (v) => {
  const s = String(v == null ? '' : v);
  return /[",\n]/.test(s) ? '"' + s.replace(/"/g, '""') + '"' : s;
};

const hotspotDir = path.join(axure, 'hotspots');
fs.mkdirSync(hotspotDir, { recursive: true });
for (const s of manifest.screens) {
  const list = byRoute.get(s.route) || [];
  const lines = ['label,action,target,x,y,w,h'];
  for (const h of list) {
    lines.push([h.label, h.action, h.target, h.x, h.y, h.w, h.h].map(csvEsc).join(','));
  }
  fs.writeFileSync(path.join(hotspotDir, s.file.split('/').pop().replace(/\.png$/, '.csv')), lines.join('\n'), 'utf8');
}

const pageLines = ['no,route,title,dept,back,image,width,height,hotspots'];
for (const s of manifest.screens) {
  const def = P.screens[s.route];
  pageLines.push([s.no, s.route, s.title, P.depts[s.dept].label, def.back || '', s.file, s.width / 2, s.height / 2, s.hotspots].map(csvEsc).join(','));
}
fs.writeFileSync(path.join(axure, 'pages.csv'), pageLines.join('\n'), 'utf8');

fs.writeFileSync(path.join(axure, 'SITEMAP.md'), sitemap.join('\n'), 'utf8');
fs.writeFileSync(path.join(axure, 'INTERACTION-MATRIX.md'), matrix.join('\n'), 'utf8');
console.log('wrote SITEMAP.md, INTERACTION-MATRIX.md, pages.csv and per-screen hotspot CSVs');
console.log('screens', manifest.screens.length, 'popups', manifest.popups.length, 'hotspots', manifest.hotspots);
