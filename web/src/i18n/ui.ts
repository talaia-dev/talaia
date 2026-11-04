// The site's three languages. English is the source; a key or a README
// missing in a translation falls back to English, so a language can
// ship before its translation is complete.
import type { CollectionEntry } from 'astro:content';

export const langs = ['en', 'es', 'ca'] as const;
export type Lang = (typeof langs)[number];
export const langNames: Record<Lang, string> = { en: 'English', es: 'Español', ca: 'Català' };
export const home = (lang: Lang) => (lang === 'en' ? '/' : `/${lang}/`);

const en = {
  title: 'talaia — human in charge',
  description: 'A human-in-charge accountability standard for work done with AI: a certificate a person signs, and a method (Brief → Build → Watch → Sign) that makes the signature honest.',
  'tab.manifesto': 'Manifesto',
  'tab.teams': 'Individuals and Teams',
  'tab.companies': 'Companies',
  'hero.tagline1': 'You direct. The AI builds.',
  'hero.tagline2': 'Your name answers for it.',
  'hero.sub': 'Talaia is a human-in-charge standard for work done with AI, and its whole claim fits in one line: the machine may produce the work; only a person can answer for it. It gives that person a certificate to sign (under a published post, at the foot of a report, in a commit) and a method that makes the signature honest.',
  'hero.cta.manifesto': 'Read the manifesto',
  'hero.cta.teams': 'Apply it, solo or as a team',
  'hero.cta.companies': 'Bring talaia to your company',
  'og.alt': 'talaia: You direct. The AI builds. Your name answers for it.',
  'og.sub': 'An accountability standard for work done with AI.',
  'og.foot': 'talaia.dev · open source · works with any AI',
};
type Key = keyof typeof en;

const es: Partial<Record<Key, string>> = {
  title: 'talaia — una persona al mando',
  description: 'Un estándar de responsabilidad human-in-charge para el trabajo hecho con IA: un certificado que firma una persona, y un método (Brief → Build → Watch → Sign) que hace honesta la firma.',
  'tab.manifesto': 'Manifiesto',
  'tab.teams': 'Personas y equipos',
  'tab.companies': 'Empresas',
  'hero.tagline1': 'Tú diriges. La IA construye.',
  'hero.tagline2': 'Tu nombre responde de ello.',
  'hero.sub': 'Talaia es un estándar human-in-charge para el trabajo hecho con IA, y todo lo que afirma cabe en una línea: la máquina puede producir el trabajo; solo una persona puede responder de él. Le da a esa persona un certificado que firmar (bajo una publicación, al pie de un informe, en un commit) y un método que hace honesta la firma.',
  'hero.cta.manifesto': 'Leer el manifiesto',
  'hero.cta.teams': 'Aplicarlo, solo o en equipo',
  'hero.cta.companies': 'Llevar talaia a tu empresa',
  'og.alt': 'talaia: Tú diriges. La IA construye. Tu nombre responde de ello.',
  'og.sub': 'Un estándar de responsabilidad para trabajar con IA.',
  'og.foot': 'talaia.dev · código abierto · funciona con cualquier IA',
};

const ca: Partial<Record<Key, string>> = {
  title: 'talaia — una persona al comandament',
  description: 'Un estàndard de responsabilitat human-in-charge per a la feina feta amb IA: un certificat que signa una persona, i un mètode (Brief → Build → Watch → Sign) que fa honesta la signatura.',
  'tab.manifesto': 'Manifest',
  'tab.teams': 'Persones i equips',
  'tab.companies': 'Empreses',
  'hero.tagline1': 'Tu dirigeixes. La IA construeix.',
  'hero.tagline2': 'El teu nom en respon.',
  'hero.sub': 'Talaia és un estàndard human-in-charge per a la feina feta amb IA, i tot el que afirma cap en una línia: la màquina pot produir la feina; només una persona en pot respondre. Dona a aquesta persona un certificat per signar (sota una publicació, al peu d\'un informe, en un commit) i un mètode que fa honesta la signatura.',
  'hero.cta.manifesto': 'Llegir el manifest',
  'hero.cta.teams': 'Aplicar-lo, sol o en equip',
  'hero.cta.companies': 'Portar talaia a la teva empresa',
  'og.alt': 'talaia: Tu dirigeixes. La IA construeix. El teu nom en respon.',
  'og.sub': 'Un estàndard de responsabilitat per treballar amb IA.',
  'og.foot': 'talaia.dev · codi obert · funciona amb qualsevol IA',
};

const ui: Record<Lang, Partial<Record<Key, string>>> = { en, es, ca };

export const t = (lang: Lang, key: Key) => ui[lang][key] ?? en[key];

/** The README of a language: `docs/<lang>/README.md`, or the English one. */
export const readmeFor = (entries: CollectionEntry<'docs'>[], lang: Lang) =>
  entries.find((e) => e.id === `${lang}/readme`) ?? entries.find((e) => e.id === 'readme')!;
