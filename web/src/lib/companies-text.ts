// The Companies section as plain markdown for llms-full.txt, read from
// the component's own template (src/components/companies/<lang>.astro),
// so the page and the text never drift. The form, its script and the
// buttons are left out; the page's address is given for contact.

/** `src` is the .astro source; `url` the Companies page; `contact` its contact line. */
export function companiesText(src: string, url: string, contact: string): string {
  let h = src.split(/^---$/m).slice(2).join('---');
  h = h
    .replace(/<script[\s\S]*?<\/script>/g, '')
    .replace(/<form[\s\S]*?<\/form>/g, '')
    .replace(/<a class="btn[^"]*"[\s\S]*?<\/a>/g, '');
  // the first overline is the part's title; the others open its sections
  let first = true;
  h = h.replace(/<p class="co-over">([\s\S]*?)<\/p>/g, (_, x) => `\n\n${first ? ((first = false), '# ') : '## '}${x}\n\n`);
  h = h.replace(/<h([234])[^>]*>([\s\S]*?)<\/h\1>/g, (_, n, x) => `\n\n${'#'.repeat(Number(n))} ${x}\n\n`);
  h = h
    .replace(/<li>/g, '\n- ')
    .replace(/<\/li>/g, '\n')
    .replace(/<strong[^>]*>([\s\S]*?)<\/strong>/g, '**$1** ')
    .replace(/<br\s*\/?>/g, ' ')
    .replace(/<\/?(p|div|section|ul)\b[^>]*>/g, '\n\n')
    .replace(/<[^>]+>/g, ' ');
  // one line per heading, paragraph or list item
  const blocks = h.split(/\n\s*\n/).map((b) => {
    const out: string[] = [];
    for (const line of b.split('\n').map((l) => l.replace(/\s+/g, ' ').trim()).filter(Boolean)) {
      if (out.length && !/^(- |#)/.test(line)) out[out.length - 1] += ' ' + line;
      else out.push(line);
    }
    return out.map((l) => l.replace(/ ([.,:;)])/g, '$1').replace(/\*\* +\*\*/g, ' ')).join('\n');
  }).filter(Boolean);
  return [...blocks, `${contact} ${url}`].join('\n\n') + '\n';
}
