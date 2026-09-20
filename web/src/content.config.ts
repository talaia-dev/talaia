import { defineCollection, z } from 'astro:content';
import { glob } from 'astro/loaders';

const docs = defineCollection({
  loader: glob({
    pattern: '**/*.md',
    base: './src/content/docs',
    retainBody: true,          // necesario para /llms.txt
  }),
  schema: z.object({
    title: z.string(),
    order: z.number(),
    anchor: z.string(),
    summary: z.string().optional(),
    draft: z.boolean().default(false),
  }),
});

export const collections = { docs };
