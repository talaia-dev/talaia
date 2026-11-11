// The README of a language, rendered, with what every page around it
// needs: the social card and the index (`toc`). The home page shows it;
// the Companies page needs its toc for the tabs and the footer.
import { getCollection, render } from 'astro:content';
import { readmeFor, type Lang } from '../i18n/ui';

export type Toc = { label: string; anchor: string; items: { anchor: string; title: string }[] }[];

export async function readmeOf(lang: Lang) {
  const entries = await getCollection('docs', ({ data }) => !data.draft);
  const readme = readmeFor(entries, lang);
  const { Content, headings } = await render(readme);
  // The social card: each translated language has its own (scripts/og.sh).
  const card = readme.id === 'readme' ? 'og.png' : `og-${lang}.png`;

  // Sidebar: one group per `#` part of the README (Manifesto, Framework,
  // Using talaia, Contributing), each listing its `##` sections. Derived
  // from the rendered headings, so renaming a section never breaks the
  // grouping; only a `##` before the first `#` would end up ungrouped.
  const toc: Toc = [];
  for (const h of headings) {
    if (h.slug === 'footnote-label') continue;   // the footnotes' own (hidden) h2, not a section
    if (h.depth === 1) toc.push({ label: h.text, anchor: h.slug, items: [] });
    else if (h.depth === 2 && toc.length) toc[toc.length - 1].items.push({ anchor: h.slug, title: h.text });
  }
  return { readme, Content, headings, card, toc };
}
