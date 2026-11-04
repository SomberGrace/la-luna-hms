/* Final regression: export copy, deep links, capture mode, landing page. */
import { chromium } from 'playwright-core';
import { pathToFileURL } from 'node:url';
import path from 'node:path';
import fs from 'node:fs';

const root = path.resolve(process.cwd(), '..');
const exe = 'C:/Program Files/Google/Chrome/Application/chrome.exe';
const exportRoot = path.join(root, 'export', 'LaLuna-HMS-prototype');
const appUrl = pathToFileURL(path.join(exportRoot, 'prototype.html')).href;
const landingUrl = pathToFileURL(path.join(exportRoot, 'index.html')).href;

const browser = await chromium.launch({ executablePath: exe, headless: true });
const page = await browser.newPage({ viewport: { width: 1280, height: 960 } });
const errors = [];
page.on('pageerror', (e) => errors.push('page: ' + e.message));
page.on('console', (m) => { if (m.type() === 'error') errors.push('console: ' + m.text()); });
page.on('requestfailed', (r) => errors.push('request: ' + r.url().split('/').pop()));

await page.goto(appUrl);
await page.waitForFunction(() => window.__proto && window.__proto.ready);
await page.waitForTimeout(300);

const report = { exportRoutes: 0, exportBroken: [], mojibake: [], popupDeepLink: false, captureModal: false, appSize: null };

const routes = await page.evaluate(() => Object.keys(PROTO.screens));
for (const r of routes) {
  await page.evaluate((route) => { location.hash = route; }, r);
  await page.waitForTimeout(45);
  const ok = await page.evaluate((route) => {
    const s = document.querySelector('.screen');
    return !!s && s.dataset.screen === route;
  }, r);
  if (ok) report.exportRoutes++;
  else report.exportBroken.push(r);
}

const bodyText = await page.evaluate(() => document.body.innerText);
const bad = bodyText.match(/[\u00c2\u00c3\u00e2]/g);
if (bad) report.mojibake = [...new Set(bad)];

await page.goto(appUrl + '?screen=/guest&modal=wifi');
await page.waitForTimeout(500);
report.popupDeepLink = await page.locator('.modal').isVisible();
report.popupTitle = await page.locator('.modal-title h2').innerText();

await page.setViewportSize({ width: 390, height: 844 });
await page.goto(appUrl + '?capture=1&screen=/frontdesk/rooms&modal=room-status');
await page.waitForTimeout(600);
report.captureModal = await page.evaluate(() => document.body.classList.contains('capture') && document.body.classList.contains('has-modal'));
report.appSize = await page.evaluate(() => {
  const r = document.querySelector('#app').getBoundingClientRect();
  return { w: Math.round(r.width), h: Math.round(r.height) };
});

/* landing + drawer + a popup shot for visual review */
const landing = await browser.newPage({ viewport: { width: 1440, height: 1000 } });
await landing.goto(landingUrl);
await landing.waitForTimeout(800);
await landing.screenshot({ path: path.join(root, 'tools', 'shots', '_landing_final.png') });

await page.setViewportSize({ width: 1280, height: 960 });
await page.goto(appUrl + '?map=1');
await page.waitForTimeout(500);
await page.screenshot({ path: path.join(root, 'tools', 'shots', '_drawer_final.png') });

await page.goto(appUrl + '?screen=/guest/dining&modal=guest-tray');
await page.waitForTimeout(500);
await page.locator('#app').screenshot({ path: path.join(root, 'tools', 'shots', '_popup_final.png') });

report.errors = errors;
fs.writeFileSync(path.join(root, 'tools', 'final-check.json'), JSON.stringify(report, null, 2));
console.log(JSON.stringify(report, null, 2));
await browser.close();
