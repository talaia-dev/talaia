// Checks dist/index.html after a build: every plain-text "talaia" in the
// README body became the in-text wordmark, and nothing that must stay
// literal was touched (code, domains, the GitHub path, headings).
//   npm run check
import { readFileSync } from 'node:fs';

const html = readFileSync(new URL('../dist/index.html', import.meta.url), 'utf8');
const article = html.slice(html.indexOf('<article class="readme">'), html.indexOf('</article>'));
const fail = (m) => { console.error('check-wordmark: ' + m); process.exitCode = 1; };

const wm = (article.match(/<span class="wm"/g) || []).length;
if (wm < 10) fail(`expected the in-text wordmark throughout the body, found ${wm}`);

// nothing inside <code> or headings may carry the wordmark
for (const m of article.matchAll(/<(code|h[1-6])\b[^>]*>([\s\S]*?)<\/\1>/g)) {
  if (m[2].includes('class="wm"')) fail(`wordmark inside <${m[1]}>: ${m[2].slice(0, 60)}`);
}
// domains and the GitHub path stay literal
if (!/talaia\.dev/.test(article)) fail('talaia.dev no longer appears literally');
if (/<span class="wm"[^>]*>talaia<\/span>(\.dev|-dev)/.test(article)) fail('a domain was split by the wordmark');
if (!/github\.com\/talaia-dev\/talaia</.test(article)) fail('the GitHub path text was altered');
// no plain-text talaia left in paragraphs/list items (outside code)
const stripped = article.replace(/<code\b[^>]*>[\s\S]*?<\/code>/g, '').replace(/<h[1-6]\b[^>]*>[\s\S]*?<\/h[1-6]>/g, '')
  .replace(/<span class="wm"[^>]*>talaia<\/span>/g, '').replace(/talaia(\.dev|-dev)/g, '');
const left = stripped.match(/>[^<]*(?<![\/\w.-])[Tt]alaia\b[^<]*</g);   // path segments (/talaia) stay literal
if (left) fail(`plain "talaia" still in text: ${left.slice(0, 3).join(' | ')}`);

if (!process.exitCode) console.log(`check-wordmark: ok (${wm} wordmarks, code/domains/headings untouched)`);
