// llms.txt and llms-full.txt in one language. llms.txt is the index:
// what talaia is, a link to the full text, the pages, Companies, the
// other languages. llms-full.txt is the same head, that language's
// README in full and the Companies page as one more part.
// `/llms.txt`, `/llms-full.txt`, `/es/…` and `/ca/…` are each one line
// around this.
import { getCollection } from 'astro:content';
import { langs, companiesPath, t, type Lang } from '../i18n/ui';
import { companiesText } from './companies-text';
import companiesEn from '../components/companies/en.astro?raw';
import companiesEs from '../components/companies/es.astro?raw';
import companiesCa from '../components/companies/ca.astro?raw';

const companiesSrc: Record<Lang, string> = { en: companiesEn, es: companiesEs, ca: companiesCa };
const contact: Record<Lang, string> = { en: 'Contact: the form on', es: 'Contacto: el formulario en', ca: 'Contacte: el formulari a' };

const site = 'https://talaia.dev';
const llmsUrl = (l: Lang, name = 'llms.txt') => `${site}${l === 'en' ? '' : '/' + l}/${name}`;

const head: Record<Lang, { lede: string[]; standard: string; fullLink: string; pages: string; pageList: string[]; others: string; follows: string; home: string }> = {
  en: {
    lede: [
      '> Human-in-charge accountability standard for work done with AI: a certificate a person signs, and a method (Brief → Build → Watch → Sign) that makes the signature honest.',
      '',
      'The signature is one line added to the work, address and name: `talaia.dev/certificate — John Doe`. In software, `git commit -s` with the address in the message.',
    ],
    standard: 'The standard',
    fullLink: 'The standard in full: the README, as published at',
    pages: 'Pages',
    pageList: [
      `- [The Certificate of Understanding](${site}/certificate): the text a person certifies, in the shape of the DCO.`,
      `- [AGENTS.md](${site}/agentsmd): the instructions file that puts the framework in front of an AI agent.`,
      `- [How to install AGENTS.md](${site}/how-to-agentsmd): where each tool reads it (Claude Code, Codex, Gemini CLI, Cursor, Copilot, Windsurf and Aider).`,
      `- [The pull request check](${site}/pr-check): the workflow that fails a pull request with an unsigned commit.`,
      `- [How to install the pull request check](${site}/how-to-pr-check): GitHub, Gitea, GitLab, Bitbucket, and how to make it required.`,
      '- [Repository](https://github.com/talaia-dev/talaia): the standard, the CLI and this site.',
    ],
    others: 'Other languages',
    follows: 'The standard follows in full, as published at',
    home: `${site}`,
  },
  es: {
    lede: [
      '> Estándar de responsabilidad human-in-charge para el trabajo hecho con IA: un certificado que firma una persona, y un método (Brief → Build → Watch → Sign) que hace honesta la firma.',
      '',
      'La firma es una línea añadida al trabajo, dirección y nombre: `talaia.dev/certificate — John Doe`. En software, `git commit -s` con la dirección en el mensaje.',
    ],
    standard: 'El estándar',
    fullLink: 'El estándar completo: el README, tal como se publica en',
    pages: 'Páginas (en inglés)',
    pageList: [
      `- [The Certificate of Understanding](${site}/certificate): el texto que certifica una persona, con la forma del DCO.`,
      `- [AGENTS.md](${site}/agentsmd): el archivo de instrucciones que pone el framework delante de un agente de IA.`,
      `- [How to install AGENTS.md](${site}/how-to-agentsmd): dónde lo lee cada herramienta (Claude Code, Codex, Gemini CLI, Cursor, Copilot, Windsurf y Aider).`,
      `- [The pull request check](${site}/pr-check): el workflow que hace fallar un pull request con un commit sin firmar.`,
      `- [How to install the pull request check](${site}/how-to-pr-check): GitHub, Gitea, GitLab, Bitbucket, y cómo hacerlo obligatorio.`,
      '- [Repositorio](https://github.com/talaia-dev/talaia): el estándar, la CLI y este sitio.',
    ],
    others: 'Otros idiomas',
    follows: 'Sigue el estándar completo, tal como se publica en',
    home: `${site}/es`,
  },
  ca: {
    lede: [
      "> Estàndard de responsabilitat human-in-charge per a la feina feta amb IA: un certificat que signa una persona, i un mètode (Brief → Build → Watch → Sign) que fa honesta la signatura.",
      '',
      "La signatura és una línia afegida a la feina, adreça i nom: `talaia.dev/certificate — John Doe`. En programari, `git commit -s` amb l'adreça al missatge.",
    ],
    standard: "L'estàndard",
    fullLink: "L'estàndard complet: el README, tal com es publica a",
    pages: 'Pàgines (en anglès)',
    pageList: [
      `- [The Certificate of Understanding](${site}/certificate): el text que certifica una persona, amb la forma del DCO.`,
      `- [AGENTS.md](${site}/agentsmd): el fitxer d'instruccions que posa el framework davant d'un agent d'IA.`,
      `- [How to install AGENTS.md](${site}/how-to-agentsmd): on el llegeix cada eina (Claude Code, Codex, Gemini CLI, Cursor, Copilot, Windsurf i Aider).`,
      `- [The pull request check](${site}/pr-check): el workflow que fa fallar un pull request amb un commit sense signar.`,
      `- [How to install the pull request check](${site}/how-to-pr-check): GitHub, Gitea, GitLab, Bitbucket, i com fer-lo obligatori.`,
      '- [Repositori](https://github.com/talaia-dev/talaia): l\'estàndard, la CLI i aquest lloc.',
    ],
    others: 'Altres idiomes',
    follows: "Segueix l'estàndard complet, tal com es publica a",
    home: `${site}/ca`,
  },
};

const names: Record<Lang, string> = { en: 'English', es: 'Español', ca: 'Català' };

export async function llms(lang: Lang, full = false): Promise<Response> {
  const id = lang === 'en' ? 'readme' : `${lang}/readme`;
  const entries = await getCollection('docs', (e) => e.id === id && !e.data.draft);
  const h = head[lang];
  const companiesUrl = `${site}${companiesPath[lang]}`;
  const body = full ? [
    '# talaia',
    ...h.lede,
    '',
    `${h.follows} ${h.home}.`,
    '',
    ...entries.map((e) => `${e.body?.trim() ?? ''}\n`),
    companiesText(companiesSrc[lang], companiesUrl, contact[lang]),
  ].join('\n') : [
    '# talaia',
    ...h.lede,
    '',
    `## ${h.standard}`,
    '',
    `- [llms-full.txt](${llmsUrl(lang, 'llms-full.txt')}): ${h.fullLink} ${h.home}.`,
    '',
    `## ${h.pages}`,
    '',
    ...h.pageList,
    '',
    `## ${t(lang, 'tab.companies')}`,
    '',
    `- [${t(lang, 'companies.title')}](${companiesUrl}): ${t(lang, 'companies.description')}`,
    '',
    `## ${h.others}`,
    '',
    ...langs.filter((l) => l !== lang).map((l) => `- [${names[l]}](${llmsUrl(l)})`),
    '',
  ].join('\n');

  return new Response(body, {
    headers: { 'Content-Type': 'text/plain; charset=utf-8' },
  });
}
