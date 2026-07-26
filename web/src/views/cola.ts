import type { ColaLeads, LeadEnCola } from '../models/tipos';
import { avatarPersona } from '../util/avatares';

const POR_PAGINA = 6;

/** Umbral de prioridad para derivar Alta/Media/Baja (LeadEnCola no trae un
 * campo de intención a nivel de cola — sólo `prioridad`, un float). Elegido
 * en frontend, sin pedir un cambio de contrato. */
type Bucket = 'Alta' | 'Media' | 'Baja';
function bucketDePrioridad(p: number): Bucket {
  if (p >= 0.7) return 'Alta';
  if (p >= 0.4) return 'Media';
  return 'Baja';
}

let paginaActual = 0;
let filtroActual: 'Todos' | Bucket = 'Todos';
let busquedaActual = '';

/**
 * Renderiza el panel de leads: encabezado, buscador, filtros, lista
 * paginada y pie. El backend YA envía la lista ordenada por prioridad
 * (US-14); el frontend nunca reordena, sólo filtra/pagina sobre ese orden.
 */
export function renderCola(
  contenedor: HTMLElement,
  cola: ColaLeads,
  leadActivoId: string | null,
  onSeleccionar: (id: string) => void,
): void {
  const conteos: Record<'Todos' | Bucket, number> = { Todos: cola.leads.length, Alta: 0, Media: 0, Baja: 0 };
  for (const l of cola.leads) conteos[bucketDePrioridad(l.prioridad)]++;

  const filtrados = cola.leads.filter(l => {
    const pasaFiltro = filtroActual === 'Todos' || bucketDePrioridad(l.prioridad) === filtroActual;
    const pasaBusqueda = !busquedaActual || l.nombre.toLocaleLowerCase('es').includes(busquedaActual);
    return pasaFiltro && pasaBusqueda;
  });

  const totalPaginas = Math.max(1, Math.ceil(filtrados.length / POR_PAGINA));
  if (paginaActual >= totalPaginas) paginaActual = totalPaginas - 1;
  const inicio = paginaActual * POR_PAGINA;
  const pagina = filtrados.slice(inicio, inicio + POR_PAGINA);

  contenedor.innerHTML = `
    <div class="leads-heading">
      <span>Todos los leads</span>
      <button type="button" aria-label="Configurar leads">⚙</button>
    </div>
    <label class="search-wrap">
      <span>⌕</span>
      <input id="buscar-lead" type="search" placeholder="Buscar lead..." value="${escapar(busquedaActual)}">
    </label>
    <nav class="filter-tabs" aria-label="Filtrar leads">
      ${(['Todos', 'Alta', 'Media', 'Baja'] as const).map(f => `
        <button type="button" class="filter ${filtroActual === f ? 'active' : ''}" data-filtro="${f}">
          ${f} <strong>${conteos[f]}</strong>
        </button>`).join('')}
    </nav>
    <div class="lead-list" id="lead-list" role="list">
      ${pagina.length ? pagina.map(l => renderFilaLead(l, l.lead_id === leadActivoId)).join('') : '<p style="padding:1rem;color:#647198;font-size:13px">Ningún lead coincide.</p>'}
    </div>
    <div class="list-pager">
      <button type="button" class="pager-button" id="pager-prev" ${paginaActual === 0 ? 'disabled' : ''} aria-label="Anterior">‹</button>
      <span>Mostrando ${pagina.length ? inicio + 1 : 0}–${inicio + pagina.length} de ${filtrados.length}</span>
      <button type="button" class="pager-button" id="pager-next" ${paginaActual >= totalPaginas - 1 ? 'disabled' : ''} aria-label="Siguiente">›</button>
    </div>
  `;

  contenedor.querySelectorAll<HTMLButtonElement>('[data-lead-id]').forEach(el => {
    el.addEventListener('click', () => onSeleccionar(el.dataset.leadId!));
  });

  const buscador = contenedor.querySelector<HTMLInputElement>('#buscar-lead')!;
  buscador.addEventListener('input', () => {
    busquedaActual = buscador.value.trim().toLocaleLowerCase('es');
    paginaActual = 0;
    renderCola(contenedor, cola, leadActivoId, onSeleccionar);
  });
  if (document.activeElement !== buscador && busquedaActual) {
    const cursor = buscador.value.length;
    buscador.focus();
    buscador.setSelectionRange(cursor, cursor);
  }

  contenedor.querySelectorAll<HTMLButtonElement>('.filter').forEach(btn => {
    btn.addEventListener('click', () => {
      filtroActual = btn.dataset.filtro as 'Todos' | Bucket;
      paginaActual = 0;
      renderCola(contenedor, cola, leadActivoId, onSeleccionar);
    });
  });

  contenedor.querySelector<HTMLButtonElement>('#pager-prev')?.addEventListener('click', () => {
    paginaActual = Math.max(0, paginaActual - 1);
    renderCola(contenedor, cola, leadActivoId, onSeleccionar);
  });
  contenedor.querySelector<HTMLButtonElement>('#pager-next')?.addEventListener('click', () => {
    paginaActual = Math.min(totalPaginas - 1, paginaActual + 1);
    renderCola(contenedor, cola, leadActivoId, onSeleccionar);
  });
}

function renderFilaLead(l: LeadEnCola, seleccionado: boolean): string {
  const bucket = bucketDePrioridad(l.prioridad);
  const claseBucket = bucket === 'Alta' ? 'hot' : bucket === 'Media' ? 'medium' : '';
  const clasePrioridad = bucket === 'Media' ? 'medium' : bucket === 'Baja' ? 'baja' : '';
  const afiliacion = l.afiliado ? 'Afiliado' : 'No afiliado';

  return `
    <button type="button" class="lead-row ${claseBucket} ${seleccionado ? 'selected' : ''}" data-lead-id="${l.lead_id}" role="listitem">
      <img class="avatar lead-avatar" src="${avatarPersona(l.lead_id)}" alt="" aria-hidden="true">
      <span class="lead-info">
        <span class="lead-name">${escapar(l.nombre)}</span>
        <span class="lead-meta">${escapar(afiliacion)}</span>
        <span class="lead-last">${escapar(l.resumen)}</span>
      </span>
      <span class="priority ${clasePrioridad}">${bucket}</span>
    </button>
  `;
}

function escapar(s: string): string {
  const d = document.createElement('div');
  d.textContent = s;
  return d.innerHTML;
}
