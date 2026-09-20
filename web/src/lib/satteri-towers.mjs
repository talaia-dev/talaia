// Sätteri hast plugin: the manifesto's "Two kinds of tower" paragraph is
// the identity statement. This frames it in one enclosed block together
// with its illustration — Babel (tall, stepped, no fire, fading) beside
// the talaies (low, a fire, the next tower down the coast). The text
// keeps the ordinary body voice; the frame does the marking, drawn by
// CSS across the two adjacent siblings (.towers is the frame's top
// half, p.towers-text the bottom half). The markdown stays plain: the
// paragraph is found by its opening words.
const OPENING = 'Two kinds of tower';

const FIGURE = `<figure class="towers" aria-hidden="true">
<svg viewBox="0 0 250 78" width="250" height="78">
  <g class="babel" fill="currentColor">
    <path d="M18 64h54v-9H18zM24 55h42v-9H24zM30 46h30v-9H30zM35 37h20v-9H35zM39 28h12v-8H39zM42 20h6l-1-6h-4z"/>
  </g>
  <text class="cap" x="45" y="76" text-anchor="middle">Babel</text>
  <path class="sea" fill="currentColor" d="M110 64h130v1.6H110z"/>
  <g class="talaia" transform="translate(112 10)" fill="currentColor">
    <path fill-rule="evenodd" d="M20 54L23.5 27H21.5V22h21v5h-2L44 54zM30 44v-8h4v8z"/>
    <path d="M86 54l1-7h4l1 7z"/>
  </g>
  <circle class="fire" cx="144" cy="23" r="6"/>
  <circle class="fire fire-far" cx="200.5" cy="53.5" r="2"/>
  <text class="cap" x="176" y="76" text-anchor="middle">talaies</text>
</svg>
</figure>`;

export default function towersPlugin() {
  return {
    name: 'talaia-towers',
    element: {
      filter: ['p'],
      visit(node, ctx) {
        if (!ctx.textContent(node).startsWith(OPENING)) return;
        ctx.insertBefore(node, { type: 'raw', value: FIGURE });
        ctx.setProperty(node, 'className', ['towers-text']);
      },
    },
  };
}
