/* Generate the Axure kit: screen captures, popup captures and the hotspot map. */
import { chromium } from 'playwright-core';
import { pathToFileURL } from 'node:url';
import path from 'node:path';
import fs from 'node:fs';

const root = path.resolve(process.cwd(), '..');
const exe = 'C:/Program Files/Google/Chrome/Application/chrome.exe';
const baseUrl = pathToFileURL(path.join(root, 'prototype.html')).href;
const outDir = path.join(root, 'axure');
const screensDir = path.join(outDir, 'screens');
const popupsDir = path.join(outDir, 'popups');
fs.mkdirSync(screensDir, { recursive: true });
fs.mkdirSync(popupsDir, { recursive: true });

const slug = (s) =>
  String(s).toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-+|-+$/g, '').slice(0, 54);

const pad = (n) => String(n).padStart(2, '0');

const modalBackdrop = {
  'book-room': '/book',
  'booking-confirmed': '/book',
  'room-status': '/frontdesk/rooms',
  'walkin-success': '/frontdesk/walkin',
  'checkin-success': '/frontdesk/checkin',
  'payment-receipt': '/frontdesk/payments',
  'checkout-success': '/frontdesk/checkout',
  'extension-success': '/frontdesk/extension',
  'add-menu-item': '/restaurant/menu',
  'new-docket': '/restaurant/dockets',
  'add-bar-item': '/bar/stock',
  'add-ingredient': '/kitchen/stock',
  'guest-tray': '/guest/dining',
  'guest-request': '/guest/services',
  wifi: '/guest',
  'guest-request-success': '/guest/services',
  'audit-review': '/manager/audits',
  'shift-close': '/frontdesk/reports',
  'sale-recorded': '/bar/sales',
};

const browser = await chromium.launch({ executablePath: exe, headless: true });
const page = await browser.newPage({ viewport: { width: 390, height: 1000 }, deviceScaleFactor: 2 });
const errors = [];
page.on('pageerror', (e) => errors.push('page: ' + e.message));
page.on('console', (m) => { if (m.type() === 'error') errors.push('console: ' + m.text()); });

async function load(url) {
  await page.goto(url, { waitUntil: 'load' });
  await page.waitForFunction(() => window.__proto && window.__proto.ready, null, { timeout: 15000 });
  await page.evaluate(async () => { await document.fonts.ready; });
  await page.waitForTimeout(220);
}

async function readyElement(selector) {
  await page.locator(selector).waitFor({ state: 'visible' });
}

/* -------- screens -------- */
await load(baseUrl + '?capture=1&screen=/');
const routes = await page.evaluate(() => Object.keys(PROTO.screens));
const modals = await page.evaluate(() => Object.keys(PROTO.modals));

const screenNames = new Map();
const rows = [];
const manifest = { generatedAt: new Date().toISOString(), screens: [], popups: [], hotspots: 0 };

let routeIndex = 0;
for (const route of routes) {
  const info = await page.evaluate((r) => {
    const s = PROTO.screens[r];
    return { dept: s.dept, title: s.title };
  }, route);
  const name = `${info.dept}-${pad(++routeIndex)}-${slug(info.title)}.png`;
  screenNames.set(route, name);

  const url = baseUrl + '?capture=1&screen=' + encodeURIComponent(route);
  await load(url);
  await readyElement('#app');

  const box = await page.locator('#app').boundingBox();
  const height = Math.ceil(box.height);
  await page.setViewportSize({ width: 390, height: Math.max(844, height) });
  await page.waitForTimeout(120);

  const shotPath = path.join(screensDir, name);
  await page.locator('#app').screenshot({ path: shotPath, animations: 'disabled' });

  const hotspots = await page.evaluate(() => window.__proto.hotspots());
  hotspots.forEach((h, i) => rows.push({
    screen_no: routeIndex,
    route,
    file: name,
    index: i + 1,
    label: h.label,
    action: h.action,
    target: h.target,
    x: h.x,
    y: h.y,
    w: h.w,
    h: h.h,
  }));

  const png = fs.readFileSync(shotPath);
  manifest.screens.push({
    no: routeIndex, route, title: info.title, dept: info.dept, file: `screens/${name}`,
    width: png.readUInt32BE(16), height: png.readUInt32BE(20), hotspots: hotspots.length,
  });
  console.log(`screen ${pad(routeIndex)} ${route} -> ${name} (${hotspots.length} hotspots)`);
}

/* -------- popups -------- */
let popupIndex = 0;
for (const id of modals) {
  const backdrop = modalBackdrop[id] || '/frontdesk';
  const url = baseUrl + '?capture=1&screen=' + encodeURIComponent(backdrop) + '&modal=' + encodeURIComponent(id);
  await page.setViewportSize({ width: 390, height: 844 });
  await load(url);
  await readyElement('#app');
  await page.waitForTimeout(300);

  const title = await page.evaluate((m) => PROTO.modals[m].title, id);
  const name = `popup-${pad(++popupIndex)}-${slug(id)}.png`;
  const shotPath = path.join(popupsDir, name);
  await page.locator('#app').screenshot({ path: shotPath, animations: 'disabled' });

  const png = fs.readFileSync(shotPath);
  manifest.popups.push({
    no: popupIndex, id, title, backdrop, file: `popups/${name}`,
    width: png.readUInt32BE(16), height: png.readUInt32BE(20),
  });
  console.log(`popup ${pad(popupIndex)} ${id} -> ${name}`);
}

/* -------- hotspot map -------- */
const csvEsc = (v) => {
  const s = String(v == null ? '' : v);
  return /[",\n]/.test(s) ? '"' + s.replace(/"/g, '""') + '"' : s;
};
const head = ['screen_no', 'route', 'file', 'index', 'label', 'action', 'target', 'x', 'y', 'w', 'h', 'x_2x', 'y_2x', 'w_2x', 'h_2x'];
const lines = [head.join(',')];
for (const r of rows) {
  lines.push([
    r.screen_no, r.route, r.file, r.index, r.label, r.action, r.target,
    r.x, r.y, r.w, r.h, r.x * 2, r.y * 2, r.w * 2, r.h * 2,
  ].map(csvEsc).join(','));
}
fs.writeFileSync(path.join(outDir, 'hotspot-map.csv'), lines.join('\n'), 'utf8');
fs.writeFileSync(path.join(outDir, 'hotspot-map.json'), JSON.stringify(rows, null, 2), 'utf8');

manifest.hotspots = rows.length;
manifest.errors = errors;
fs.writeFileSync(path.join(outDir, 'manifest.json'), JSON.stringify(manifest, null, 2), 'utf8');

console.log(`\n${manifest.screens.length} screens, ${manifest.popups.length} popups, ${manifest.hotspots} hotspots`);
if (errors.length) console.log('ERRORS:', errors);

await browser.close();
