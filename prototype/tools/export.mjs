/* Create a clean standalone copy of the browser prototype in export/. */
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const here = path.dirname(fileURLToPath(import.meta.url));
const root = path.resolve(here, '..');
const target = path.join(root, 'export', 'LaLuna-HMS-prototype');

const files = ['index.html', 'prototype.html', 'README.md'];
const dirs = ['assets', 'axure'];

fs.rmSync(path.join(root, 'export'), { recursive: true, force: true });
fs.mkdirSync(target, { recursive: true });

for (const f of files) fs.copyFileSync(path.join(root, f), path.join(target, f));
for (const d of dirs) fs.cpSync(path.join(root, d), path.join(target, d), { recursive: true });

const walk = (dir) => fs.readdirSync(dir, { withFileTypes: true }).flatMap((e) => {
  const p = path.join(dir, e.name);
  return e.isDirectory() ? walk(p) : [p];
});
const all = walk(target);
const bytes = all.reduce((n, f) => n + fs.statSync(f).size, 0);
console.log(`copied ${all.length} files (${(bytes / 1024).toFixed(0)} KB) to ${path.relative(root, target)}`);
