import type { ColaLeads, LeadEnCola } from '../models/tipos';

const ZONA_SEMAFORO = {
  VERDE: { etiqueta: 'VERDE', zona: 'zona-verde' },
  AMBAR: { etiqueta: 'ÁMBAR', zona: 'zona-ambar' },
  GRIS: { etiqueta: 'GRIS', zona: 'zona-gris' },
} as const;

/**
 * Renderiza la cola priorizada de leads como un registro tipo informe.
 * IMPORTANTE: El backend YA envía la lista ordenada por prioridad (US-14).
 * El frontend NO reordena la lista.
 */
export function renderCola(
  contenedor: HTMLElement,
  cola: ColaLeads,
  onVerFicha: (id: string) => void,
  onVerChat: (id: string) => void,
): void {
  const usados = cola.cupo_10.usados;
  const ventana = cola.cupo_10.porcentaje_ventana;
  const pct = (usados / (ventana || 1)) * 100;
  const esAlerta = pct >= 80;

  contenedor.innerHTML = `
    <div class="hoja-informe informe-cola">
      <header class="masthead-cola">
        <div class="masthead-cola-titulo">
          <span class="masthead-kicker">Registro de leads</span>
          <p class="masthead-cifra"><span class="cifra">${cola.leads.length}</span> en cola</p>
        </div>
        <div class="medidor-cupo ${esAlerta ? 'alerta' : ''}"
             title="${esAlerta ? 'Cupo regulatorio de no afiliados casi lleno (≥80%)' : 'Uso del cupo regulatorio del 10%'}">
          <span class="medidor-cupo-etiqueta">Cupo regulatorio 10%</span>
          <div class="medidor-cupo-escala" role="img" aria-label="Cupo usado: ${usados} de ${ventana}">
            ${renderTicks(ventana)}
            <div class="medidor-cupo-relleno" style="transform: scaleX(${pct / 100})"></div>
          </div>
          <span class="medidor-cupo-cifra cifra">${usados}/${ventana}</span>
        </div>
      </header>
      <ul class="registro-leads" role="list">
        ${cola.leads.map(renderFilaLead).join('')}
      </ul>
    </div>
  `;

  // Listeners para clic y botones dentro de cada fila
  contenedor.querySelectorAll<HTMLLIElement>('[data-lead-id]').forEach(el => {
    const leadId = el.dataset.leadId!;
    const btnChat = el.querySelector<HTMLButtonElement>('[data-btn-chat]');

    if (btnChat) {
      btnChat.addEventListener('click', (e) => {
        e.stopPropagation(); // Evitar que dispare la apertura de la ficha
        onVerChat(leadId);
      });
    }

    el.addEventListener('click', () => onVerFicha(leadId));
    el.addEventListener('keydown', (e: KeyboardEvent) => {
      if (e.key === 'Enter' || e.key === ' ') {
        e.preventDefault();
        onVerFicha(leadId);
      }
    });
  });
}

function renderFilaLead(l: LeadEnCola): string {
  const zona = ZONA_SEMAFORO[l.semaforo] ?? ZONA_SEMAFORO.GRIS;
  const inicial = l.nombre.trim().charAt(0).toUpperCase() || '?';
  return `
    <li class="fila-lead ${zona.zona}" data-lead-id="${l.lead_id}" tabindex="0" role="listitem">
      <span class="fila-lead-avatar" aria-hidden="true">${escapar(inicial)}</span>
      <div class="fila-lead-cuerpo">
        <div class="fila-lead-linea1">
          <span class="lead-nombre">${escapar(l.nombre)}</span>
          <span class="fila-lead-afiliado ${l.afiliado ? 'es-afiliado' : ''}">${l.afiliado ? 'Afiliado' : 'No afiliado'}</span>
          <span class="fila-lead-zona" aria-label="Semáforo ${l.semaforo}">${zona.etiqueta}</span>
          <span class="lead-ruta">${escapar(l.ruta)}</span>
        </div>
        <p class="lead-resumen">${escapar(l.resumen)}</p>
      </div>
      <span class="fila-lead-prioridad" title="Prioridad calculada">
        <span class="fila-lead-prioridad-etiqueta">Prioridad</span>
        <span class="fila-lead-prioridad-cifra cifra">${l.prioridad.toFixed(2)}</span>
      </span>
      <button class="btn-ver-chat" data-btn-chat="true" title="Ver chat en vivo con ${escapar(l.nombre)}" type="button">
        Ver chat
      </button>
    </li>
  `;
}

/** Marcas de calibración de la escala del cupo (una por unidad, sin el borde final). */
function renderTicks(max: number): string {
  const n = Math.max(1, Math.round(max));
  return Array.from({ length: n - 1 }, (_, i) =>
    `<span class="medidor-tick" style="left: ${((i + 1) / n) * 100}%"></span>`,
  ).join('');
}

/** SIEMPRE escapar: anti-XSS */
function escapar(s: string): string {
  const d = document.createElement('div');
  d.textContent = s;
  return d.innerHTML;
}
