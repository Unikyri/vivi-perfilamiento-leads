export interface OpcionLeadDemo {
  id: string;
  nombre: string;
  descripcion: string;
}

const LEADS_PRECARGADOS: OpcionLeadDemo[] = [
  { id: 'mock-1', nombre: 'Ana Rodríguez', descripcion: 'Afiliada Cat. A · Presupuesto $166.8M' },
  { id: 'mock-2', nombre: 'Carlos Martínez', descripcion: 'No Afiliado · Presupuesto $210M' },
  { id: 'mock-3', nombre: 'Luisa Gómez', descripcion: 'Afiliada Cat. B · Presupuesto $145M' },
];

/**
 * Renderiza la botonera del Demo (US-12, US-17).
 * Estilo "Control Room" (grafito #111827), visualmente separada del producto.
 */
export function renderBotoneraDemo(
  contenedor: HTMLElement,
  leadActivoId: string | null,
  onSeleccionarLead: (leadId: string) => void,
  onAvanzarTiempo: () => void,
  onReiniciarDemo: () => void,
): void {
  contenedor.innerHTML = `
    <div class="botonera-demo">
      <select id="select-lead-demo" class="select-lead-demo" aria-label="Seleccionar lead precargado">
        ${LEADS_PRECARGADOS.map(l => `
          <option value="${l.id}" ${l.id === leadActivoId ? 'selected' : ''}>
            👤 ${escapar(l.nombre)}
          </option>
        `).join('')}
      </select>

      <button id="btn-avanzar-tiempo" class="btn-demo-action" type="button" title="Simular avance en la linea de tiempo">
        ⏩ Avanzar tiempo
      </button>

      <button id="btn-reiniciar-demo" class="btn-demo-action" type="button" title="Reiniciar demo al estado inicial (<3s)">
        ↺ Reiniciar demo
      </button>
    </div>
  `;

  // Listeners
  const selectLead = contenedor.querySelector<HTMLSelectElement>('#select-lead-demo');
  if (selectLead) {
    selectLead.addEventListener('change', () => onSeleccionarLead(selectLead.value));
  }

  const btnAvanzar = contenedor.querySelector<HTMLButtonElement>('#btn-avanzar-tiempo');
  if (btnAvanzar) {
    btnAvanzar.addEventListener('click', onAvanzarTiempo);
  }

  const btnReiniciar = contenedor.querySelector<HTMLButtonElement>('#btn-reiniciar-demo');
  if (btnReiniciar) {
    btnReiniciar.addEventListener('click', onReiniciarDemo);
  }
}

function escapar(s: string): string {
  const d = document.createElement('div');
  d.textContent = s;
  return d.innerHTML;
}
