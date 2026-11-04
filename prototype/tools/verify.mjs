import { chromium } from 'playwright-core';
import { pathToFileURL } from 'node:url';
import path from 'node:path';
import fs from 'node:fs';

const root = path.resolve(process.cwd(), '..');
const exe = 'C:/Program Files/Google/Chrome/Application/chrome.exe';
const url = pathToFileURL(path.join(root, 'prototype.html')).href;

const browser = await chromium.launch({ executablePath: exe, headless: true });
const page = await browser.newPage({ viewport: { width: 1320, height: 1020 }, deviceScaleFactor: 1 });
const errors = [];
page.on('console', (m) => { if (m.type() === 'error') errors.push('console: ' + m.text()); });
page.on('pageerror', (e) => errors.push('page: ' + e.message));
page.on('requestfailed', (r) => errors.push('request: ' + r.url().split('/').pop() + ' ' + (r.failure()?.errorText || '')));

await page.goto(url);
await page.waitForFunction(() => window.__proto && window.__proto.ready, null, { timeout: 10000 });
await page.waitForTimeout(400);

const report = { routes: [], brokenLinks: [], screens: {}, modals: 0, errors };

const routes = await page.evaluate(() => Object.keys(PROTO.screens));
const modals = await page.evaluate(() => Object.keys(PROTO.modals));
report.modals = modals.length;

for (const r of routes) {
  await page.evaluate((route) => { location.hash = route; }, r);
  await page.waitForTimeout(80);
  const info = await page.evaluate(() => {
    const screen = document.querySelector('.screen');
    return {
      rendered: !!screen,
      route: screen?.dataset.screen || '',
      rows: document.querySelectorAll('.screen .row, .screen .tile, .screen .room, .screen .dish, .screen .quick-item, .screen .ticket').length,
      height: document.querySelector('.screen')?.getBoundingClientRect().height || 0,
      title: document.title
    };
  });
  report.routes.push(r);
  report.screens[r] = info;
  if (!info.rendered || info.route !== r) report.brokenLinks.push('screen not rendered: ' + r);
}

// validate action targets across every screen + modal
const actionIssues = await page.evaluate(() => {
  const issues = [];
  const walk = (blocks, where) => (blocks || []).forEach((b) => {
    const acts = [];
    if (b.action) acts.push(b.action);
    if (b.submit) acts.push(b.submit);
    if (b.total && b.total.action) acts.push(b.total.action);
    (b.items || []).forEach((i) => { if (i.action) acts.push(i.action); });
    if (b.action && b.action.action) acts.push(b.action.action);
    acts.forEach((a) => {
      if (!a) return;
      if (a.do === 'nav' && !PROTO.screens[a.to]) issues.push(where + ' → missing route ' + a.to);
      if (a.do === 'modal' && !PROTO.modals[a.to]) issues.push(where + ' → missing modal ' + a.to);
    });
  });
  Object.entries(PROTO.screens).forEach(([r, s]) => walk(s.blocks, r));
  Object.entries(PROTO.modals).forEach(([id, m]) => walk(m.blocks, 'modal:' + id));
  return issues;
});
report.actionIssues = actionIssues;

// interactions: open a modal, close with Escape
await page.evaluate(() => { location.hash = '/frontdesk/rooms'; });
await page.waitForTimeout(150);
const modalTrigger = page.locator('.screen [data-do="modal"]').first();
await modalTrigger.click();
await page.waitForTimeout(350);
report.modalOpened = await page.locator('.modal').isVisible();
await page.keyboard.press('Escape');
await page.waitForTimeout(350);
report.modalClosed = !(await page.locator('.modal').isVisible());

// nav flow: click a nav action
await page.evaluate(() => { location.hash = '/frontdesk'; });
await page.waitForTimeout(200);
const navTrigger = page.locator('.screen [data-do="nav"]').first();
const before = await page.evaluate(() => location.hash);
await navTrigger.click();
await page.waitForTimeout(250);
report.navChanged = before !== await page.evaluate(() => location.hash);

// hotspot count + drawer
await page.evaluate(() => { location.hash = '/frontdesk'; });
await page.waitForTimeout(150);
report.hotspots = await page.evaluate(() => window.__proto.hotspots().length);
await page.click('#btn-screens');
await page.waitForTimeout(300);
report.drawerOpen = await page.locator('.drawer.is-open').isVisible();
await page.keyboard.press('Escape');
await page.waitForTimeout(300);

// screenshots for visual review
const shots = ['/frontdesk', '/frontdesk/rooms', '/guest/dining', '/restaurant/dockets', '/manager'];
fs.mkdirSync(path.join(root, 'tools', 'shots'), { recursive: true });
for (const r of shots) {
  await page.evaluate((route) => { location.hash = route; }, r);
  await page.waitForTimeout(200);
  await page.locator('.device').screenshot({ path: path.join(root, 'tools', 'shots', r.replace(/\//g, '_') + '.png') });
}
await page.evaluate(() => { location.hash = '/guest/dining'; });
await page.waitForTimeout(200);
await page.locator('.screen [data-do="modal"]').last().click();
await page.waitForTimeout(400);
await page.locator('.device').screenshot({ path: path.join(root, 'tools', 'shots', '_modal_tray.png') });

// landing page
const landing = await browser.newPage({ viewport: { width: 1440, height: 1000 } });
await landing.goto(pathToFileURL(path.join(root, 'index.html')).href);
await landing.waitForTimeout(700);
await landing.screenshot({ path: path.join(root, 'tools', 'shots', '_landing.png') });

console.log(JSON.stringify(report, null, 2));
await browser.close();
