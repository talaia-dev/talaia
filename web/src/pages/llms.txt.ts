import type { APIRoute } from 'astro';
import { getCollection } from 'astro:content';

export const GET: APIRoute = async () => {
  const entries = (await getCollection('docs', ({ data }) => !data.draft))
    .sort((a, b) => a.data.order - b.data.order);

  const body = [
    '# talaia',
    '> Human-in-command accountability standard for work done with AI: a certificate a person signs, and a method that makes the signature honest.',
    '',
    ...entries.map((e) => `${e.body?.trim() ?? ''}\n`),
  ].join('\n');

  return new Response(body, {
    headers: { 'Content-Type': 'text/plain; charset=utf-8' },
  });
};
