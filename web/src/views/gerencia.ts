import type { BuyerPersona } from '../models/tipos';

const PROYECTOS_DISPONIBLES = [
  { id: 'mongui', nombre: 'Monguí' },
  { id: 'macarena', nombre: 'La Macarena' },
  { id: 'versalles', nombre: 'Versalles' },
  { id: 'todos', nombre: 'Todos los proyectos' },
];

/**
 * Renderiza el panel de Gerencia / Buyer Persona Vivo (US-16) como un
 * informe de junta directiva: cifras tabulares sobre papel, sin tarjetas.
 */
export function renderGerencia(
  contenedor: HTMLElement,
  bp: BuyerPersona | null,
  proyectoSeleccionadoId = 'mongui',
  onCambiarProyecto: (proyectoId: string) => void,
): void {
  const muestras = bp ? bp.muestras : 312;
  const desistimientoPct = bp ? (bp.tasa_desistimiento * 100).toFixed(0) : '11';
  const proyectoActual = PROYECTOS_DISPONIBLES.find(p => p.id === proyectoSeleccionadoId)?.nombre ?? 'Proyecto';

  // Datos fallback si el backend aún no retorna agregados
  const afilData = bp?.afiliacion ?? { afiliados: 180, no_afiliados: 132 };
  const catData = bp?.categoria ?? { 'Cat A': 90, 'Cat B': 65, 'Cat C': 25, 'No Afiliado': 132 };
  const edadData = bp?.rango_edad ?? { '18-25': 20, '26-35': 150, '36-45': 90, '46+': 52 };

  contenedor.innerHTML = `
    <div class="hoja-informe informe-gerencia">
      <header class="masthead-gerencia">
        <div class="masthead-gerencia-titulo">
          <span class="masthead-kicker">Buyer persona vivo</span>
          <h2 class="masthead-nombre">${escapar(proyectoActual)}</h2>
        </div>
        <div class="selector-proyecto-wrap">
          <label for="select-proyecto-gerencia">Proyecto</label>
          <select id="select-proyecto-gerencia" class="select-proyecto">
            ${PROYECTOS_DISPONIBLES.map(p => `
              <option value="${p.id}" ${p.id === proyectoSeleccionadoId ? 'selected' : ''}>
                ${escapar(p.nombre)}
              </option>
            `).join('')}
          </select>
        </div>
      </header>

      <p class="nota-actualizacion">Se actualiza en tiempo real con cada lead perfilado.</p>

      <!-- Franja de métricas clave -->
      <div class="franja-metricas">
        <div class="franja-metrica">
          <span class="franja-metrica-etiqueta">Personas en vivo interesadas</span>
          <span class="franja-metrica-cifra cifra">${muestras}</span>
        </div>
        <div class="franja-metrica">
          <span class="franja-metrica-etiqueta">Tasa de desistimiento histórica</span>
          <span class="franja-metrica-cifra cifra">${desistimientoPct}%</span>
        </div>
      </div>

      <!-- Columnas de distribución (estilo junta directiva) -->
      <section class="panel-graficos">
        <div class="columna-grafico">
          <h4 class="seccion-titulo">Distribución por afiliación</h4>
          <div class="barras-lista">
            ${renderBarras(afilData)}
          </div>
        </div>
        <div class="columna-grafico">
          <h4 class="seccion-titulo">Categoría de afiliación</h4>
          <div class="barras-lista">
            ${renderBarras(catData)}
          </div>
        </div>
        <div class="columna-grafico">
          <h4 class="seccion-titulo">Rango de edad</h4>
          <div class="barras-lista">
            ${renderBarras(edadData)}
          </div>
        </div>
      </section>
    </div>
  `;

  // Listener para el selector de proyecto
  const selectEl = contenedor.querySelector<HTMLSelectElement>('#select-proyecto-gerencia');
  if (selectEl) {
    selectEl.addEventListener('change', () => {
      onCambiarProyecto(selectEl.value);
    });
  }
}

// Claves crudas del Contrato (data/compradores.json) que no son ya un
// rótulo legible — los rangos de edad ('20-35', '55+', ...) sí lo son
// y se muestran tal cual.
const ETIQUETAS_CLAVE: Record<string, string> = {
  afiliados: 'Afiliados',
  no_afiliados: 'No afiliados',
  SIN_DATO: 'Sin dato',
  A: 'Categoría A',
  B: 'Categoría B',
  C: 'Categoría C',
};

function etiquetarClave(clave: string): string {
  return ETIQUETAS_CLAVE[clave] ?? clave;
}

function renderBarras(datos: Record<string, number>): string {
  const maxVal = Math.max(...Object.values(datos), 1);

  return Object.entries(datos).map(([lbl, val]) => {
    const pctBarra = ((val / maxVal) * 100).toFixed(0);
    return `
      <div class="barra-item">
        <span class="barra-label">${escapar(etiquetarClave(lbl))}</span>
        <div class="barra-track">
          <div class="barra-fill" style="transform: scaleX(${Number(pctBarra) / 100})"></div>
        </div>
        <span class="barra-val cifra">${val}</span>
      </div>
    `;
  }).join('');
}

function escapar(s: string): string {
  const d = document.createElement('div');
  d.textContent = s;
  return d.innerHTML;
}
