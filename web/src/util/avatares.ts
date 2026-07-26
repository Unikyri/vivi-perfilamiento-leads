/**
 * Genera avatares e ilustraciones de proyecto como data-URI SVG, deterministas
 * por semilla (id o nombre) y sin ninguna llamada de red.
 *
 * Por qué no una foto real: ni LeadEnCola ni Ficha ni Proyecto traen un campo
 * de foto en el Contrato — no hay de dónde sacarla. La versión anterior
 * (prototipo.css) resolvía esto con background-image fijo a UNA foto
 * ("ana-avatar.png") aplicada a TODOS los avatares sin importar el lead — es
 * decir, mentía. Esta versión da a cada semilla su propia ilustración
 * consistente (mismo lead → mismo avatar siempre) sin inventar una identidad
 * que no existe.
 */

const PALETA_PERSONA = [
  ['#f4a261', '#7a3e1d'], ['#e76f51', '#5c2418'], ['#2a9d8f', '#0b3b36'],
  ['#264653', '#0f1e24'], ['#e9c46a', '#7a5a15'], ['#8ecae6', '#1d4e6b'],
  ['#cdb4db', '#5a3d6b'], ['#a3b18a', '#3f4a30'],
] as const;

const PALETA_PROYECTO = [
  ['#1375c2', '#9bd4ef'], ['#2a9d8f', '#8fd9cf'], ['#e9c46a', '#f4dfa0'],
  ['#8ecae6', '#c9e8f5'], ['#e76f51', '#f3ab92'], ['#264653', '#5c7d8a'],
] as const;

function hashSemilla(semilla: string): number {
  let h = 0;
  for (let i = 0; i < semilla.length; i++) h = (h * 31 + semilla.charCodeAt(i)) >>> 0;
  return h;
}

function aSvgDataUri(svg: string): string {
  return `data:image/svg+xml;utf8,${encodeURIComponent(svg)}`;
}

/** Avatar de persona: círculo de color + silueta genérica de cabeza/hombros.
 * Misma semilla (lead_id) → siempre el mismo color y siempre la misma forma. */
export function avatarPersona(semilla: string): string {
  const [piel, sombra] = PALETA_PERSONA[hashSemilla(semilla) % PALETA_PERSONA.length];
  const fondo = hashSemilla(semilla + 'x') % 2 === 0 ? '#eef2fb' : '#e7f0ff';
  const svg = `
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64">
      <circle cx="32" cy="32" r="32" fill="${fondo}"/>
      <circle cx="32" cy="26" r="12" fill="${piel}"/>
      <path d="M10 58c2-13 12-21 22-21s20 8 22 21z" fill="${piel}"/>
      <path d="M10 58c2-13 12-21 22-21s20 8 22 21" fill="none" stroke="${sombra}" stroke-width="1.5" opacity="0.35"/>
    </svg>`.trim();
  return aSvgDataUri(svg);
}

/** Ilustración de proyecto: skyline plano de color, distinto por proyecto_id
 * pero siempre el mismo para el mismo proyecto (no pretende ser una foto real
 * del edificio — no existe ese dato). */
export function avatarProyecto(semilla: string): string {
  const [cielo, torres] = PALETA_PROYECTO[hashSemilla(semilla) % PALETA_PROYECTO.length];
  const h = hashSemilla(semilla + 'y');
  const alturas = [22, 34, 26, 40, 18].map((base, i) => base + ((h >> (i * 3)) % 10));
  const anchoBarra = 100 / alturas.length;
  const barras = alturas.map((alto, i) => {
    const x = i * anchoBarra + 2;
    const ancho = anchoBarra - 4;
    const y = 60 - alto;
    return `<rect x="${x}" y="${y}" width="${ancho}" height="${alto}" fill="${torres}" rx="1.5"/>`;
  }).join('');
  const svg = `
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 60" preserveAspectRatio="xMidYMax slice">
      <rect width="100" height="60" fill="${cielo}"/>
      ${barras}
      <rect y="58" width="100" height="2" fill="${torres}" opacity="0.6"/>
    </svg>`.trim();
  return aSvgDataUri(svg);
}
