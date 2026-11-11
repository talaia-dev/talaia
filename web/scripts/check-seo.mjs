// Checks dist/ after a build for what search engines and LLMs read:
// robots.txt, one canonical per page with no trailing slash, the
// sitemap listing exactly those canonicals, hreflang, og:locale and
// JSON-LD on the three home pages, and one llms.txt per language.
//   npm run check
import { existsSync, readFileSync } from 'node:fs';

const site = 'https://talaia.dev';
const dist = (p) => new URL(`../dist/${p}`, import.meta.url);
const read = (p) => (existsSync(dist(p)) ? readFileSync(dist(p), 'utf8') : null);
const fail = (m) => { console.error('check-seo: ' + m); process.exitCode = 1; };
const parts = (md) => (md.replace(/^```[\s\S]*?^```/gm, '').match(/^# /gm) || []).length;

const robots = read('robots.txt');
if (!robots) fail('dist/robots.txt missing');
else {
  if (!robots.includes(`Sitemap: ${site}/sitemap-index.xml`)) fail('robots.txt: no Sitemap line');
  if (/^Disallow:\s*\/\s*$/m.test(robots)) fail('robots.txt: disallows the whole site');
}

// page file → the URL it is served at (no trailing slash but the root)
const pages = {
  'index.html': `${site}/`,
  'es.html': `${site}/es`,
  'ca.html': `${site}/ca`,
  'certificate.html': `${site}/certificate`,
  'how-to-agentsmd.html': `${site}/how-to-agentsmd`,
  'how-to-pr-check.html': `${site}/how-to-pr-check`,
  'companies.html': `${site}/companies`,
  'es/empresas.html': `${site}/es/empresas`,
  'ca/empreses.html': `${site}/ca/empreses`,
};
const homes = { en: 'index.html', es: 'es.html', ca: 'ca.html' };
const locales = { en: 'en_US', es: 'es_ES', ca: 'ca_ES' };

for (const [file, url] of Object.entries(pages)) {
  const html = read(file);
  if (!html) { fail(`dist/${file} missing`); continue; }
  const canon = [...html.matchAll(/<link rel="canonical" href="([^"]+)"/g)].map((m) => m[1]);
  if (canon.length !== 1 || canon[0] !== url) fail(`dist/${file}: canonical ${JSON.stringify(canon)}, expected ${url}`);
  const og = html.match(/<meta property="og:url" content="([^"]+)"/)?.[1];
  if (og !== url) fail(`dist/${file}: og:url ${og}, expected ${url}`);
}

for (const [lang, file] of Object.entries(homes)) {
  const html = read(file);
  if (!html) continue;
  for (const [l, f] of Object.entries(homes))
    if (!html.includes(`<link rel="alternate" hreflang="${l}" href="${pages[f]}"`)) fail(`dist/${file}: hreflang ${l} → ${pages[f]} missing`);
  if (!html.includes(`<link rel="alternate" hreflang="x-default" href="${site}/"`)) fail(`dist/${file}: hreflang x-default missing`);
  if (!html.includes(`<meta property="og:locale" content="${locales[lang]}"`)) fail(`dist/${file}: og:locale ${locales[lang]} missing`);
  for (const [l, loc] of Object.entries(locales))
    if (l !== lang && !html.includes(`<meta property="og:locale:alternate" content="${loc}"`)) fail(`dist/${file}: og:locale:alternate ${loc} missing`);
  for (const name of ['llms.txt', 'llms-full.txt']) {
    const href = lang === 'en' ? `/${name}` : `/${lang}/${name}`;
    if (!html.includes(`<link rel="alternate" type="text/plain" href="${href}"`)) fail(`dist/${file}: alternate ${href} missing`);
    const foot = html.slice(html.indexOf('<footer'), html.indexOf('</footer>'));
    if (!foot.includes(`href="${href}"`)) fail(`dist/${file}: footer link to ${href} missing`);
  }
  const ld = html.match(/<script type="application\/ld\+json">([\s\S]*?)<\/script>/)?.[1];
  if (!ld) { fail(`dist/${file}: no JSON-LD`); continue; }
  let data;
  try { data = JSON.parse(ld); } catch (e) { fail(`dist/${file}: JSON-LD does not parse: ${e.message}`); continue; }
  // the two terms the standard defines (Human-in-charge, the Certificate):
  // name = a part's title, description = a sentence of that README, url = an anchor on the page
  const terms = (data['@graph'] ?? []).filter((n) => n['@type'] === 'DefinedTerm');
  const md = readFileSync(new URL(`../../docs/${lang === 'en' ? '' : lang + '/'}README.md`, import.meta.url), 'utf8');
  const flat = md.replace(/\*/g, '').replace(/\s+/g, ' ');
  if (terms.length !== 2) fail(`dist/${file}: ${terms.length} DefinedTerm in JSON-LD, expected 2`);
  for (const term of terms) {
    if (!new RegExp(`^# ${term.name.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')}$`, 'm').test(md)) fail(`dist/${file}: term "${term.name}" is not a part of its README`);
    if (!term.description || !flat.includes(term.description)) fail(`dist/${file}: term "${term.name}": description not found in its README`);
    const id = term.url?.startsWith(pages[file] + '#') ? term.url.slice(pages[file].length + 1) : null;
    if (!id || !html.includes(`id="${id}"`)) fail(`dist/${file}: term "${term.name}": url ${term.url} is not an anchor on the page`);
  }
  const web = (data['@graph'] ?? [data]).find((n) => n['@type'] === 'WebSite');
  if (!web) fail(`dist/${file}: JSON-LD has no WebSite`);
  else {
    if (web.inLanguage !== lang) fail(`dist/${file}: JSON-LD inLanguage ${web.inLanguage}, expected ${lang}`);
    if (web.url !== pages[file]) fail(`dist/${file}: JSON-LD url ${web.url}, expected ${pages[file]}`);
  }
}

// the certificate page describes itself: its text, its version, in English
{
  const html = read('certificate.html') ?? '';
  const ld = html.match(/<script type="application\/ld\+json">([\s\S]*?)<\/script>/)?.[1];
  let work;
  try { work = ld && JSON.parse(ld); } catch (e) { fail(`certificate: JSON-LD does not parse: ${e.message}`); }
  if (!work) fail('certificate: no JSON-LD');
  else {
    if (work['@type'] !== 'CreativeWork') fail(`certificate: JSON-LD @type ${work['@type']}, expected CreativeWork`);
    if (work.url !== pages['certificate.html']) fail(`certificate: JSON-LD url ${work.url}`);
    if (work.inLanguage !== 'en') fail('certificate: JSON-LD inLanguage is not en');
    if (!work.version || !work.text?.includes(`Version ${work.version}`)) fail('certificate: JSON-LD version missing or not in its text');
  }
}

// Companies: its own page per language, linked from the home page's tab
// and hero, and no longer inside the home page
const companies = { en: 'companies.html', es: 'es/empresas.html', ca: 'ca/empreses.html' };
for (const [lang, file] of Object.entries(companies)) {
  const html = read(file);
  if (!html) continue;
  const href = pages[file].slice(site.length);
  if (!html.includes('<section id="companies"')) fail(`dist/${file}: the Companies section is missing`);
  if (!html.includes(`<html lang="${lang}"`)) fail(`dist/${file}: <html lang="${lang}"> missing`);
  const title = html.match(/<title>([^<]*)<\/title>/)?.[1];
  if (!title || title === read(homes[lang])?.match(/<title>([^<]*)<\/title>/)?.[1]) fail(`dist/${file}: no title of its own`);
  for (const [l, f] of Object.entries(companies))
    if (!html.includes(`<link rel="alternate" hreflang="${l}" href="${pages[f]}"`)) fail(`dist/${file}: hreflang ${l} → ${pages[f]} missing`);
  if (!html.includes(`<link rel="alternate" hreflang="x-default" href="${site}/companies"`)) fail(`dist/${file}: hreflang x-default missing`);
  if (!html.includes(`<meta property="og:locale" content="${locales[lang]}"`)) fail(`dist/${file}: og:locale missing`);
  if (!new RegExp(`<option value="${lang}"[^>]*data-href="${href}"`).test(html)) fail(`dist/${file}: the language menu does not keep the page`);
  const ld = html.match(/<script type="application\/ld\+json">([\s\S]*?)<\/script>/)?.[1];
  const node = ld && JSON.parse(ld)['@graph']?.find((n) => n['@type'] === 'WebPage');
  if (!node || node.url !== pages[file] || node.inLanguage !== lang) fail(`dist/${file}: JSON-LD WebPage missing or wrong`);
  const home = read(homes[lang]) ?? '';
  if (home.includes('<section id="companies"')) fail(`dist/${homes[lang]}: still carries the Companies section`);
  if (!home.includes(`data-tab="companies" href="${href}"`) && !home.includes(`href="${href}" data-tab="companies"`)) fail(`dist/${homes[lang]}: the Companies tab does not link to ${href}`);
  if (!home.includes(`class="btn" href="${href}"`)) fail(`dist/${homes[lang]}: the hero button does not link to ${href}`);
}

const sitemap = read('sitemap-0.xml');
if (!sitemap) fail('dist/sitemap-0.xml missing');
else {
  // @astrojs/sitemap writes the root as `https://talaia.dev` whenever
  // trailingSlash is 'never' (write-sitemap.js); an empty path is `/`
  // (RFC 3986 §6.2.3), so that one is read as the canonical root.
  const locs = [...sitemap.matchAll(/<loc>([^<]+)<\/loc>/g)].map((m) => (m[1] === site ? `${site}/` : m[1])).sort();
  const want = Object.values(pages).sort();
  if (JSON.stringify(locs) !== JSON.stringify(want)) fail(`sitemap: ${JSON.stringify(locs)}, expected ${JSON.stringify(want)}`);
  for (const m of sitemap.matchAll(/hreflang="[^"]+" href="([^"]+)"/g))
    if (!want.includes(m[1])) fail(`sitemap: alternate ${m[1]} is not a canonical URL`);
}

// llms.txt is the index (no README in it), llms-full.txt the README in full
for (const lang of ['en', 'es', 'ca']) {
  const dir = lang === 'en' ? '' : `${lang}/`;
  const index = read(`${dir}llms.txt`);
  const full = read(`${dir}llms-full.txt`);
  if (!index) fail(`dist/${dir}llms.txt missing`);
  if (!full) fail(`dist/${dir}llms-full.txt missing`);
  if (!index || !full) continue;
  const md = readFileSync(new URL(`../../docs/${dir}README.md`, import.meta.url), 'utf8');
  if (parts(index) !== 1) fail(`${dir}llms.txt: ${parts(index) - 1} README parts in the index, expected none`);
  if (!index.includes(`${site}/${dir}llms-full.txt`)) fail(`${dir}llms.txt: no link to ${site}/${dir}llms-full.txt`);
  const co = pages[companies[lang]];
  if (!index.includes(`(${co})`)) fail(`${dir}llms.txt: no link to ${co}`);
  // the head, the README's parts, and Companies as one more part
  if (parts(full) !== parts(md) + 2) fail(`${dir}llms-full.txt: ${parts(full) - 2} README parts (+ head + Companies), its README has ${parts(md)}`);
  if (!full.includes(co)) fail(`${dir}llms-full.txt: the Companies part does not give ${co}`);
  if (/<|\{contact\}/.test(full.slice(full.lastIndexOf('\n# ')))) fail(`${dir}llms-full.txt: markup left in the Companies part`);
  for (const l of ['en', 'es', 'ca'])
    if (l !== lang && !index.includes(`${site}/${l === 'en' ? '' : l + '/'}llms.txt`)) fail(`${dir}llms.txt: no link to the ${l} llms.txt`);
}

if (!process.exitCode) console.log("check-seo: ok (robots, canonicals, sitemap, hreflang, og:locale, JSON-LD, llms.txt + llms-full.txt ×3)");
