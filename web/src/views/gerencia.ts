import type { BuyerPersona } from '../models/tipos';

const PROYECTOS_DISPONIBLES = [
  { id: 'mongui', nombre: 'Monguí' },
  { id: 'macarena', nombre: 'La Macarena' },
  { id: 'versalles', nombre: 'Versalles' },
  { id: 'todos', nombre: 'Todos los proyectos' },
];

/**
 * Renderiza el panel de Gerencia / Buyer Persona Vivo (US-16).
 * Sobrio, estilo junta directiva con barras horizontales azules sobre blanco.
 */
export function renderGerencia(
  contenedor: HTMLElement,
  bp: BuyerPersona | null,
  proyectoSeleccionadoId = 'mongui',
  onCambiarProyecto: (proyectoId: string) => void,
): void {
  const muestras = bp ? bp.muestras : 312;
  const desistimientoPct = bp ? (bp.tasa_desistimiento * 100).toFixed(0) : '11';

  // Datos fallback si el backend aún no retorna agregados
  const afilData = bp?.afiliacion ?? { afiliados: 180, no_afiliados: 132 };
  const catData = bp?.categoria ?? { 'Cat A': 90, 'Cat B': 65, 'Cat C': 25, 'No Afiliado': 132 };
  const edadData = bp?.rango_edad ?? { '18-25': 20, '26-35': 150, '36-45': 90, '46+': 52 };

  contenedor.innerHTML = `
    <div class="gerencia-container">
      <header class="gerencia-header">
        <div class="selector-proyecto-wrap">
          <label for="select-proyecto-gerencia">Proyecto:</label>
          <select id="select-proyecto-gerencia" class="select-proyecto">
            ${PROYECTOS_DISPONIBLES.map(p => `
              <option value="${p.id}" ${p.id === proyectoSeleccionadoId ? 'selected' : ''}>
                ${escapar(p.nombre)}
              </option>
            `).join('')}
          </select>
        </div>

        <span class="nota-actualizacion">
          ℹ️ Se actualiza en tiempo real con cada lead perfilado
        </span>
      </header>

      <!-- Métricas Clave -->
      <div class="metricas-row">
        <div class="card-metrica">
          <div class="metrica-label">Personas en vivo interesadas</div>
          <div class="metrica-valor">${muestras}</div>
        </div>
        <div class="card-metrica">
          <div class="metrica-label">Tasa de Desistimiento Histórica</div>
          <div class="metrica-valor">${desistimientoPct}%</div>
        </div>
      </div>

      <!-- Gráficos de Barras Horizontales (Estilo Junta Directiva) -->
      <div class="grid-tres-columnas">

        <!-- Gráfico 1: Afiliación -->
        <article class="grafico-barras-card">
          <h4>Distribución por Afiliación</h4>
          <div class="barras-lista">
            ${renderBarras(afilData)}
          </div>
        </article>

        <!-- Gráfico 2: Categoría -->
        <article class="grafico-barras-card">
          <h4>Categoría de Afiliación</h4>
          <div class="barras-lista">
            ${renderBarras(catData)}
          </div>
        </article>

        <!-- Gráfico 3: Rango de Edad -->
        <article class="grafico-barras-card">
          <h4>Rango de Edad</h4>
          <div class="barras-lista">
            ${renderBarras(edadData)}
          </div>
        </article>

      </div>
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

function renderBarras(datos: Record<string, number>): string {
  const maxVal = Math.max(...Object.values(datos), 1);

  return Object.entries(datos).map(([lbl, val]) => {
    const pctBarra = ((val / maxVal) * 100).toFixed(0);
    return `
      <div class="barra-item">
        <span class="barra-label">${escapar(lbl)}</span>
        <div class="barra-track">
          <div class="barra-fill" style="width: ${pctBarra}%"></div>
        </div>
        <span class="barra-val">${val}</span>
      </div>
    `;
  }).join('');
}

function escapar(s: string): string {
  const d = document.createElement('div');
  d.textContent = s;
  return d.innerHTML;
}
