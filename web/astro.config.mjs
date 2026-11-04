import { defineConfig } from 'astro/config';
import sitemap from '@astrojs/sitemap';
import { satteri } from '@astrojs/markdown-satteri';
import wordmark from './src/lib/satteri-wordmark.mjs';
import towers from './src/lib/satteri-towers.mjs';

export default defineConfig({
  site: 'https://talaia.dev',
  output: 'static',
  // No trailing slash, the Developer's choice: `about.html` files, which
  // Cloudflare Pages serves at `/about` (a directory would 308 to `/about/`).
  trailingSlash: 'never',
  build: { format: 'file' },
  // English at `/`, the translations at `/es/` and `/ca/`.
  i18n: { defaultLocale: 'en', locales: ['en', 'es', 'ca'], routing: { prefixDefaultLocale: false } },
  // The two files the standard hands out live in git, not on the site:
  // the addresses redirect to the scaffold templates in the repository.
  redirects: {
    '/agentsmd': 'https://github.com/talaia-dev/talaia/blob/main/internal/modules/scaffold/templates/AGENTS.md',
    '/pr-check': 'https://github.com/talaia-dev/talaia/blob/main/internal/modules/scaffold/templates/talaia-workflow.yml',
  },
  integrations: [sitemap({
    filter: (page) => !page.includes('/og-card'),
    i18n: { defaultLocale: 'en', locales: { en: 'en', es: 'es', ca: 'ca' } },
  })],
  markdown: {
    shikiConfig: { themes: { light: 'github-light', dark: 'github-dark' } },
    processor: satteri({ hastPlugins: [wordmark, towers] }),
  },
});
