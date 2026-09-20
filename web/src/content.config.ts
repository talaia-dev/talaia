import { defineCollection, z } from 'astro:content';
import { glob } from 'astro/loaders';

const docs = defineCollection({
  loader: glob({
    pattern: '**/*.md',
    base: './src/content/docs',
    retainBody: true,          // necesario para /llms.txt
  }),
  // The README carries no frontmatter: the site derives everything from
  // its headings, and GitHub would render the block as a table.
  schema: z.object({
    draft: z.boolean().default(false),
  }),
});

export const collections = { docs };
