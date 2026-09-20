// Sätteri hast plugin: every plain-text "talaia" / "Talaia" in the README
// body becomes the in-text wordmark — lowercase, in the serif voice, with a
// normal i (the fire-dotted i is for the hero and the top bar only; in
// running text it sat too high). Untouched: code, headings, domains
// (talaia.dev, talaia-dev) and path segments (/talaia).
const SKIP = new Set(['code', 'pre', 'h1', 'h2', 'h3', 'h4', 'h5', 'h6', 'script', 'style']);
const RE = /(?<![\w/.-])[Tt]alaia(?![\w.-])/g;

const text = (value) => ({ type: 'text', value });
const wordmark = () => ({ type: 'element', tagName: 'span', properties: { className: ['wm'] }, children: [text('talaia')] });

export default function wordmarkPlugin() {
  return {
    name: 'talaia-wordmark',
    text(node, ctx) {
      if (!RE.test(node.value)) { RE.lastIndex = 0; return; }
      RE.lastIndex = 0;
      for (let p = ctx.parent(node); p && p.type === 'element'; p = ctx.parent(p)) {
        if (SKIP.has(p.tagName)) return;
      }
      const pieces = [];
      let last = 0;
      for (const m of node.value.matchAll(RE)) {
        if (m.index > last) pieces.push(text(node.value.slice(last, m.index)));
        pieces.push(wordmark());
        last = m.index + m[0].length;
      }
      if (last < node.value.length) pieces.push(text(node.value.slice(last)));
      ctx.insertBefore(node, pieces);
      ctx.removeNode(node);
    },
  };
}
