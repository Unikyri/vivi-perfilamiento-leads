import { api, ErrorAPI } from '../models/api';
import { obtener, actualizar, suscribir } from '../models/estado';
import type { ColaLeads } from '../models/tipos';
import { renderCola } from '../views/cola';
import { renderFicha } from '../views/ficha';
import { renderBotoneraDemo, fechaAvanceDemo } from '../views/demo';
import { actualizarCabeceraChat } from '../views/chat';

const INTERVALO_POLL_COLA_MS = 5000;
let pollColaTimer: ReturnType<typeof setInterval> | null = null;

// Evita refetchear la Ficha en cada poll de la cola (cada 5 s) cuando el
// lead activo no cambió — la consulta más cara del sistema no debe
// dispararse por datos que nadie pidió (regresión ya corregida en #89).
let ultimaFichaCargada: string | null = null;

// `suscribir` notifica en CADA actualizar(), incluido el poll del chat cada
// 1.5s (que sólo cambia `conversacion`). Sin este guard, la lista de leads
// se re-renderizaría 1.5s sí, 1.5s también — perdiendo el foco del buscador
// y recalculando filtros por datos que nadie tocó. Mismo defecto de la #89,
// esta vez en el panel de leads en vez de en Ficha/Gerencia.
let ultimaColaRenderizada: ColaLeads | null = null;
let ultimoLeadRenderizado: string | null = null;
let ultimaSeccionRenderizada: SeccionNav | null = null;

type SeccionNav = 'leads' | 'conversaciones' | 'nutricion' | 'proyectos';
let seccionActiva: SeccionNav = 'leads';

/**
 * Inicializa el panel de leads (izquierda) y la ficha (derecha). El chat
 * (centro) lo monta e independiza `controllers/chat.ts` — este módulo sólo
 * le avisa del cambio de lead activo vía `actualizarCabeceraChat`.
 */
export function iniciarDashboard(
  leadsPanelEl: HTMLElement,
  detailsEl: HTMLElement,
  botoneraDemoEl: HTMLElement | null,
  navEl: HTMLElement | null,
): void {
  if (navEl) {
    navEl.querySelectorAll<HTMLButtonElement>('.nav-item').forEach(btn => {
      btn.addEventListener('click', () => {
        seccionActiva = (btn.dataset.seccion as SeccionNav) ?? 'leads';
        navEl.querySelectorAll('.nav-item').forEach(item => item.classList.toggle('active', item === btn));
        renderPanelLeads(leadsPanelEl);
      });
    });
  }

  if (botoneraDemoEl) montarBotoneraDemo(botoneraDemoEl);

  suscribir(() => {
    renderPanelLeads(leadsPanelEl);
    renderPanelFicha(detailsEl);
  });

  pollColaTimer = setInterval(cargarCola, INTERVALO_POLL_COLA_MS);
  cargarCola();
  renderPanelFicha(detailsEl);
}

export function detenerDashboard(): void {
  if (pollColaTimer) {
    clearInterval(pollColaTimer);
    pollColaTimer = null;
  }
}

async function cargarCola(): Promise<void> {
  try {
    const cola = await api.cola();
    actualizar({ cola });
  } catch (err) {
    console.warn('[dashboard] Error cargando cola:', err);
  }
}

function ordenarParaSeccion(cola: ColaLeads, seccion: SeccionNav): ColaLeads {
  if (seccion !== 'conversaciones') return cola;
  // "Conversaciones" reordena por actividad reciente en vez de prioridad —
  // es la única diferenciación real posible con los campos que expone
  // LeadEnCola hoy (no hay conteo de no-leídos en el contrato).
  const leads = [...cola.leads].sort(
    (a, b) => new Date(b.actualizado_en).getTime() - new Date(a.actualizado_en).getTime(),
  );
  return { ...cola, leads };
}

function renderPanelLeads(contenedor: HTMLElement, forzar = false): void {
  const st = obtener();

  if (seccionActiva === 'nutricion' || seccionActiva === 'proyectos') {
    if (!forzar && seccionActiva === ultimaSeccionRenderizada) return;
    ultimaSeccionRenderizada = seccionActiva;
    contenedor.innerHTML = `
      <div class="leads-heading"><span>${seccionActiva === 'nutricion' ? 'Nutrición' : 'Proyectos'}</span></div>
      <div class="details-vacio">
        <p>Próximamente.</p>
      </div>
    `;
    return;
  }

  if (!st.cola) {
    if (ultimaSeccionRenderizada !== null) return; // ya se mostró "cargando"
    ultimaSeccionRenderizada = null;
    contenedor.innerHTML = '<div class="leads-heading"><span>Todos los leads</span></div><p style="padding:1rem;color:#647198">Cargando leads…</p>';
    return;
  }

  // Sólo re-renderizar si de verdad cambió la cola, el lead activo o la
  // sección — no en cada poll del chat, que no afecta a este panel.
  if (!forzar && st.cola === ultimaColaRenderizada && st.leadActivo === ultimoLeadRenderizado && seccionActiva === ultimaSeccionRenderizada) {
    return;
  }
  const cambioDeLead = st.leadActivo !== ultimoLeadRenderizado;
  ultimaColaRenderizada = st.cola;
  ultimoLeadRenderizado = st.leadActivo;
  ultimaSeccionRenderizada = seccionActiva;

  // Cubre el lead seleccionado automáticamente al arrancar (main.ts crea el
  // lead inicial y setea leadActivo directo, sin pasar por seleccionarLead):
  // sin esto, el header del chat se quedaría en "Selecciona un lead" hasta
  // que el usuario hiciera clic manual sobre una fila ya resaltada.
  if (cambioDeLead && st.leadActivo) {
    const lead = st.cola.leads.find(l => l.lead_id === st.leadActivo);
    if (lead) actualizarCabeceraChat(lead.lead_id, lead.nombre);
  }

  const cola = ordenarParaSeccion(st.cola, seccionActiva);
  renderCola(contenedor, cola, st.leadActivo, leadId => seleccionarLead(leadId));
}

function seleccionarLead(leadId: string): void {
  const st = obtener();
  const lead = st.cola?.leads.find(l => l.lead_id === leadId);
  if (lead) actualizarCabeceraChat(lead.lead_id, lead.nombre);
  actualizar({ leadActivo: leadId });
}

async function renderPanelFicha(contenedor: HTMLElement): Promise<void> {
  const st = obtener();

  if (!st.leadActivo) {
    ultimaFichaCargada = null;
    contenedor.innerHTML = '<div class="details-vacio">Seleccioná un lead de la lista para ver su ficha comercial.</div>';
    return;
  }
  if (st.leadActivo === ultimaFichaCargada) return;
  ultimaFichaCargada = st.leadActivo;

  const leadEnCola = st.cola?.leads.find(l => l.lead_id === st.leadActivo);
  const nombreFallback = leadEnCola?.nombre ?? 'Lead';
  const ruta = leadEnCola?.ruta ?? null;

  try {
    const ficha = await api.ficha(st.leadActivo);
    renderFicha(contenedor, ficha, nombreFallback, ruta);
  } catch (err) {
    if (err instanceof ErrorAPI && err.estadoHTTP === 404) {
      renderFicha(contenedor, null, nombreFallback, ruta);
    } else {
      contenedor.innerHTML = `<div class="details-vacio">⚠️ Error cargando la ficha comercial: ${escapar((err as Error).message)}</div>`;
    }
  }
}

function avisarDemo(texto: string, esError = false): void {
  const el = document.getElementById('demo-aviso');
  if (!el) return;
  el.textContent = texto;
  el.classList.toggle('es-error', esError);
  setTimeout(() => {
    el.textContent = '';
    el.classList.remove('es-error');
  }, 4000);
}

function montarBotoneraDemo(contenedor: HTMLElement): void {
  renderBotoneraDemo(
    contenedor,
    async () => {
      try {
        const r = await api.avanzarTiempo(fechaAvanceDemo());
        avisarDemo(`Tiempo avanzado a ${fechaAvanceDemo()} · ${r.hitos_disparados} hito(s) disparado(s).`);
        cargarCola();
      } catch (e) {
        avisarDemo(`Error al avanzar tiempo: ${(e as Error).message}`, true);
      }
    },
    async () => {
      try {
        await api.reiniciar();

        let nuevoLead: string | null = null;
        try {
          nuevoLead = (await api.crearLead('ana')).lead_id;
        } catch (e) {
          console.error('[dashboard] no se pudo recrear el lead tras reiniciar:', e);
        }

        ultimaFichaCargada = null;
        actualizar({ leadActivo: nuevoLead });
        cargarCola();
        avisarDemo('Demo reiniciada al estado inicial.');
      } catch (e) {
        avisarDemo(`Error al reiniciar demo: ${(e as Error).message}`, true);
      }
    },
  );
}

function escapar(s: string): string {
  const d = document.createElement('div');
  d.textContent = s;
  return d.innerHTML;
}
