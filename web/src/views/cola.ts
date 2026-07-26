import type { ColaLeads, LeadEnCola } from '../models/tipos';

const CONFIG_SEMAFORO = {
  VERDE: { icono: '🟢', etiqueta: 'Alta Prio / Apto', clase: 'badge-semaforo-verde' },
  AMBAR: { icono: '🟡', etiqueta: 'En Validación', clase: 'badge-semaforo-ambar' },
  GRIS:  { icono: '⚪', etiqueta: 'Sin Datos', clase: 'badge-semaforo-gris' },
} as const;

/**
 * Renderiza la cola priorizada de leads con un diseño de tabla ejecutivo.
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
  const pct = Math.min(100, Math.round((usados / (ventana || 1)) * 100));
  const esAlerta = pct >= 80;

  contenedor.innerHTML = `
    <div class="cola-container">
      <header class="cola-header-card">
        <div class="cola-header-titulo">
          <div class="header-icon-badge">🎯</div>
          <div>
            <h2>Cola Priorizada de Leads</h2>
            <p class="cola-header-sub">Asignación inteligente en tiempo real según modelo de scoring</p>
          </div>
          <span class="badge-total-leads">${cola.leads.length} Activos</span>
        </div>

        <div class="cupo-regulado-card ${esAlerta ? 'alerta' : ''}" 
             title="${esAlerta ? '¡Alerta! Cupo del 10% de no afiliados casi agotado (≥80%)' : 'Uso del cupo regulatorio del 10% para no afiliados'}">
          <div class="cupo-info-top">
            <span class="cupo-label">⚡ Cupo Regulado 10%</span>
            <span class="cupo-metrica">${usados} / ${ventana} (${pct}%)</span>
          </div>
          <div class="cupo-track-custom">
            <div class="cupo-fill-custom" style="transform: scaleX(${pct / 100}); transform-origin: left;"></div>
          </div>
        </div>
      </header>

      <div class="tabla-cola-wrapper">
        <table class="tabla-cola-leads">
          <thead>
            <tr>
              <th scope="col" class="th-pos"># Pos</th>
              <th scope="col" class="th-estado">Estado</th>
              <th scope="col" class="th-nombre">Lead / Prospecto</th>
              <th scope="col" class="th-ruta">Canal / Ruta</th>
              <th scope="col" class="th-prio">Score Prio</th>
              <th scope="col" class="th-resumen">Resumen de Interacción</th>
              <th scope="col" class="th-acciones">Acciones</th>
            </tr>
          </thead>
          <tbody>
            ${cola.leads.map((l, index) => renderFilaTabla(l, index + 1)).join('')}
          </tbody>
        </table>
      </div>
    </div>
  `;

  // Asignación de event listeners por fila y botones de acción
  contenedor.querySelectorAll<HTMLTableRowElement>('[data-lead-id]').forEach(el => {
    const leadId = el.dataset.leadId!;
    const btnChat = el.querySelector<HTMLButtonElement>('[data-btn-chat]');
    const btnFicha = el.querySelector<HTMLButtonElement>('[data-btn-ficha]');

    if (btnChat) {
      btnChat.addEventListener('click', (e) => {
        e.stopPropagation();
        onVerChat(leadId);
      });
    }

    if (btnFicha) {
      btnFicha.addEventListener('click', (e) => {
        e.stopPropagation();
        onVerFicha(leadId);
      });
    }

    // Clic directo en la fila abre la ficha comercial
    el.addEventListener('click', () => onVerFicha(leadId));
    el.addEventListener('keydown', (e: KeyboardEvent) => {
      if (e.key === 'Enter' || e.key === ' ') {
        e.preventDefault();
        onVerFicha(leadId);
      }
    });
  });
}

function renderFilaTabla(l: LeadEnCola, posicion: number): string {
  const sem = CONFIG_SEMAFORO[l.semaforo] ?? CONFIG_SEMAFORO.GRIS;
  const inicial = l.nombre.trim().charAt(0).toUpperCase() || 'L';
  const afilBadge = l.afiliado !== undefined 
    ? `<span class="badge-afiliado-pill ${l.afiliado ? 'es-afiliado' : 'no-afiliado'}">${l.afiliado ? 'Afiliado' : 'No afiliado'}</span>` 
    : '';

  return `
    <tr class="fila-lead-tabla" data-lead-id="${l.lead_id}" tabindex="0" role="row">
      <td class="td-pos">
        <span class="pos-badge pos-${posicion}">#${posicion}</span>
      </td>
      <td class="td-estado">
        <span class="badge-semaforo ${sem.clase}" title="Estado: ${sem.etiqueta}">
          <span class="dot-icon">${sem.icono}</span>
          <span class="txt-estado">${sem.etiqueta}</span>
        </span>
      </td>
      <td class="td-nombre">
        <div class="lead-avatar-wrap">
          <div class="lead-avatar">${inicial}</div>
          <div class="lead-meta">
            <span class="lead-nombre-txt">${escapar(l.nombre)} ${afilBadge}</span>
            <span class="lead-id-sub">ID: ${escapar(l.lead_id)}</span>
          </div>
        </div>
      </td>
      <td class="td-ruta">
        <span class="badge-ruta" title="Origen del lead">${escapar(l.ruta)}</span>
      </td>
      <td class="td-prio">
        <span class="badge-prio-score cifra" title="Puntuación de Priorización Vivi">
          ${l.prioridad.toFixed(1)} <small>pts</small>
        </span>
      </td>
      <td class="td-resumen">
        <p class="resumen-corte">${escapar(l.resumen)}</p>
      </td>
      <td class="td-acciones">
        <div class="acciones-btn-group">
          <button class="btn-tabla-chat" data-btn-chat="true" title="Abrir chat de WhatsApp con ${escapar(l.nombre)}" type="button">
            💬 Chat
          </button>
          <button class="btn-tabla-ficha" data-btn-ficha="true" title="Ver Ficha Comercial de ${escapar(l.nombre)}" type="button">
            👁️ Ficha
          </button>
        </div>
      </td>
    </tr>
  `;
}

/** SIEMPRE escapar strings para evitar XSS */
function escapar(s: string): string {
  const d = document.createElement('div');
  d.textContent = s;
  return d.innerHTML;
}
