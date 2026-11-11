// Checks dist/ after a build: one page per language, each rendering
// only its own README (English when the translation is missing), and
// the English-only certificate and the English llms-full.txt reading only
// the English README (the translated llms.txt: check-seo.mjs).
//   npm run check
import { existsSync, readFileSync } from 'node:fs';

const read = (p) => readFileSync(new URL(p, import.meta.url), 'utf8');
const fail = (m) => { console.error('check-i18n: ' + m); process.exitCode = 1; };
const parts = (md) => (md.replace(/^```[\s\S]*?^```/gm, '').match(/^# /gm) || []).length;

const en = read('../../docs/README.md');
for (const [lang, page] of [['en', 'index.html'], ['es', 'es.html'], ['ca', 'ca.html']]) {
  const file = new URL(`../dist/${page}`, import.meta.url);
  if (!existsSync(file)) { fail(`dist/${page} missing`); continue; }
  const html = readFileSync(file, 'utf8');
  if (!html.includes(`<html lang="${lang}"`)) fail(`dist/${page}: <html lang="${lang}"> missing`);
  const src = new URL(`../../docs/${lang === 'en' ? '' : lang + '/'}README.md`, import.meta.url);
  const md = existsSync(src) ? readFileSync(src, 'utf8') : en;
  const article = html.slice(html.indexOf('<article class="readme">'), html.indexOf('</article>'));
  const h1 = (article.match(/<h1\b/g) || []).length;
  if (h1 !== parts(md)) fail(`dist/${page}: ${h1} parts rendered, its README has ${parts(md)}`);
  // the towers figure is found by the paragraph's opening words, which are translated
  if (!article.includes('class="towers"')) fail(`dist/${page}: the towers figure is missing`);
  // the social card: its own in each translated language, and it must exist
  const og = html.match(/<meta property="og:image" content="https:\/\/talaia\.dev\/([^"]+)"/)?.[1];
  const card = existsSync(src) ? (lang === 'en' ? 'og.png' : `og-${lang}.png`) : 'og.png';
  if (og !== card) fail(`dist/${page}: og:image is ${og}, expected ${card}`);
  else if (!existsSync(new URL(`../public/${og}`, import.meta.url))) fail(`public/${og} missing (sh scripts/og.sh)`);
  for (const l of ['en', 'es', 'ca'])
    if (!new RegExp(`<option value="${l}"[^>]*data-href=`).test(html)) fail(`dist/${page}: language option ${l} has no target`);
}

const llms = read('../dist/llms-full.txt');
if (parts(llms) !== parts(en) + 2) fail(`llms-full.txt: ${parts(llms) - 2} parts (+ head + Companies), the English README has ${parts(en)}`);
const cert = read('../dist/certificate.html');
if ((cert.match(/<pre\b/g) || []).length !== 1) fail('certificate: expected exactly one <pre>');
if (!cert.includes('<html lang="en"')) fail('certificate: not English');

if (!process.exitCode) console.log('check-i18n: ok (en, es, ca; certificate English only)');
