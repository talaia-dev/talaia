import type { APIRoute } from 'astro';
import { getCollection } from 'astro:content';

export const GET: APIRoute = async () => {
  const entries = await getCollection('docs', ({ data }) => !data.draft);

  const body = [
    '# talaia',
    '> Human-in-command accountability standard for work done with AI: a certificate a person signs, and a method — Brief → Build → Watch → Sign — that makes the signature honest.',
    '',
    'The signature is one line added to the work, address and name: `talaia.dev/certificate — John Doe`. In software, `git commit -s` with the address in the message.',
    '',
    '## Pages',
    '',
    '- [The Certificate of Understanding](https://talaia.dev/certificate): the text a person certifies, in the shape of the DCO.',
    '- [AGENTS.md](https://talaia.dev/agentsmd): the instructions file that puts the framework in front of an AI agent.',
    '- [How to install AGENTS.md](https://talaia.dev/how-to-agentsmd): where each tool reads it — Claude Code, Codex, Gemini CLI, Cursor, Copilot, Windsurf, Aider, and the chat apps.',
    '- [The pull request check](https://talaia.dev/pr-check): the workflow that fails a pull request with an unsigned commit.',
    '- [How to install the pull request check](https://talaia.dev/how-to-pr-check): GitHub, Gitea, GitLab, Bitbucket, and how to make it required.',
    '- [Repository](https://github.com/talaia-dev/talaia): the standard, the CLI and this site.',
    '',
    'The standard follows in full, as published at https://talaia.dev.',
    '',
    ...entries.map((e) => `${e.body?.trim() ?? ''}\n`),
  ].join('\n');

  return new Response(body, {
    headers: { 'Content-Type': 'text/plain; charset=utf-8' },
  });
};
