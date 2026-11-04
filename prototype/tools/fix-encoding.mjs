import fs from 'node:fs';

const p = 'assets/js/data-screens-guest.js';
let s = fs.readFileSync(p, 'utf8').replace(/^\uFEFF/, '');

const map = new Map([
  ['\u00e2\u20ac\u201c', '\u2013'], // â€“ -> –
  ['\u00e2\u20ac\u201d', '\u2014'], // â€” -> —
  ['\u00e2\u20ac\u2122', '\u2019'], // â€™ -> ’
  ['\u00e2\u20ac\u02dc', '\u2018'], // â€˜ -> ‘
  ['\u00e2\u20ac\u0153', '\u201c'], // â€œ -> “
  ['\u00e2\u20ac\u009d', '\u201d'], // â€  -> ”
  ['\u00e2\u20ac\u00a2', '\u2022'], // â€¢ -> •
  ['\u00e2\u20ac\u00a6', '\u2026'], // â€¦ -> …
  ['\u00c2\u00b7', '\u00b7'],       // Â· -> ·
  ['\u00c3\u2014', '\u00d7'],       // Ã— -> ×
]);

for (const [bad, good] of map) s = s.split(bad).join(good);

fs.writeFileSync(p, s, 'utf8');

const left = s.match(/[\u00c2\u00c3\u00e2]/g) || [];
console.log('remaining mojibake chars:', left.length);
console.log('has BOM:', fs.readFileSync(p)[0] === 0xef);
console.log(s.split('\n').filter((l) => l.includes('07:00') || l.includes('Room 104 ·')).slice(0, 4).join('\n'));
