import type { ColaLeads, LeadEnCola } from '../models/tipos';

const COLOR_SEMAFORO = { VERDE: '🟢', AMBAR: '🟡', GRIS: '⚪' } as const;

/**
 * Renderiza la cola priorizada de leads.
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
    <div class="cola-container">
      <header class="cola-header">
        <h2>Leads en Cola (${cola.leads.length})</h2>
        <div class="cupo-bar ${esAlerta ? 'alerta' : ''}"
             title="${esAlerta ? 'Cupo regulatorio de no afiliados casi lleno (≥80%)' : 'Uso del cupo regulatorio del 10%'}">
          <span>Cupo 10%:</span>
          <progress value="${usados}" max="${ventana}"></progress>
          <span>${usados}/${ventana}</span>
        </div>
      </header>
      <ul class="lista-leads" role="list">
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
  const semaforoIcono = COLOR_SEMAFORO[l.semaforo] ?? '⚪';
  return `
    <li class="fila-lead" data-lead-id="${l.lead_id}" tabindex="0" role="listitem">
      <span class="semaforo-dot" aria-label="Semáforo ${l.semaforo}">${semaforoIcono}</span>
      <span class="lead-nombre">${escapar(l.nombre)}</span>
      <span class="lead-ruta">${escapar(l.ruta)}</span>
      <span class="lead-prio-badge" title="Prioridad calculada">Prio ${l.prioridad.toFixed(2)}</span>
      <button class="btn-ver-chat" data-btn-chat="true" title="Ver chat en vivo con ${escapar(l.nombre)}" type="button">
        💬 Ver chat
      </button>
      <p class="lead-resumen">${escapar(l.resumen)}</p>
    </li>
  `;
}

/** SIEMPRE escapar: anti-XSS */
function escapar(s: string): string {
  const d = document.createElement('div');
  d.textContent = s;
  return d.innerHTML;
}
