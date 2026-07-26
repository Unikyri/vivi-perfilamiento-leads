import { api, ErrorAPI } from '../models/api';
import { obtener, actualizar, suscribir } from '../models/estado';
import { renderCola } from '../views/cola';
import { renderFicha } from '../views/ficha';
import { renderGerencia } from '../views/gerencia';
import { renderBotoneraDemo } from '../views/demo';
import { solicitarTexto, mostrarConfirmacion, mostrarNotificacion } from '../views/modal';

const INTERVALO_POLL_COLA_MS = 5000;
let pollColaTimer: ReturnType<typeof setInterval> | null = null;
let proyectoGerenciaSeleccionado = 'mongui';

/**
 * Inicializa el Panel del Asesor y sus tres pestañas (Cola, Ficha, Gerencia) + Botonera Demo.
 */
export function iniciarDashboard(
  contenedorTab: HTMLElement,
  botoneraDemoEl: HTMLElement | null,
  navTabsEl: HTMLElement | null,
): void {
  // Configurar cambio de pestañas en la navegación
  if (navTabsEl) {
    const btns = navTabsEl.querySelectorAll<HTMLButtonElement>('button[data-tab]');
    btns.forEach(btn => {
      btn.addEventListener('click', () => {
        const tab = btn.dataset.tab as 'cola' | 'ficha' | 'gerencia';
        cambiarTab(tab, btns);
      });
    });
  }

  // Montar botonera demo si existe el contenedor
  if (botoneraDemoEl) {
    montarBotoneraDemo(botoneraDemoEl);
  }

  // Suscribirse a cambios de estado global para re-renderizar la pestaña activa
  suscribir(() => renderTabActiva(contenedorTab));

  // Polling periódico de la cola cada 5.000 ms (Contrato §0)
  pollColaTimer = setInterval(() => cargarCola(), INTERVALO_POLL_COLA_MS);

  // Carga inicial
  cargarCola();
  renderTabActiva(contenedorTab);
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

function cambiarTab(nuevaTab: 'cola' | 'ficha' | 'gerencia', btns: NodeListOf<HTMLButtonElement>): void {
  actualizar({ tabActiva: nuevaTab });
  btns.forEach(b => {
    const esActiva = b.dataset.tab === nuevaTab;
    b.setAttribute('aria-selected', esActiva ? 'true' : 'false');
  });
}

let ultimoTabRenderizado: string | null = null;
let ultimoColaRenderizado: unknown = null;
let ultimoLeadFichaRenderizado: string | null = null;

function renderTabActiva(contenedor: HTMLElement): void {
  const st = obtener();
  const tabCambio = st.tabActiva !== ultimoTabRenderizado;

  switch (st.tabActiva) {
    case 'cola':
      if (st.cola) {
        if (tabCambio || st.cola !== ultimoColaRenderizado) {
          ultimoColaRenderizado = st.cola;
          renderCola(
            contenedor,
            st.cola,
            leadId => seleccionarLead(leadId),
            leadId => seleccionarChat(leadId),
          );
        }
      } else {
        contenedor.innerHTML = '<div style="padding:1rem; color:#6B7280">Cargando cola de leads…</div>';
      }
      break;

    case 'ficha':
      if (st.leadActivo) {
        if (tabCambio || st.leadActivo !== ultimoLeadFichaRenderizado) {
          ultimoLeadFichaRenderizado = st.leadActivo;
          cargarYRenderizarFicha(contenedor, st.leadActivo);
        }
      } else {
        contenedor.innerHTML = '<div style="padding:1.5rem; text-align:center; color:#6B7280">Selecciona un lead de la cola para ver su ficha comercial.</div>';
      }
      break;

    case 'gerencia':
      if (tabCambio) {
        cargarYRenderizarGerencia(contenedor, proyectoGerenciaSeleccionado);
      }
      break;
  }

  ultimoTabRenderizado = st.tabActiva;
}

function seleccionarLead(leadId: string): void {
  actualizar({ leadActivo: leadId, tabActiva: 'ficha' });

  // Actualizar el estado visual de la tab en el DOM
  const navTabs = document.querySelector('.tabs');
  if (navTabs) {
    const btns = navTabs.querySelectorAll<HTMLButtonElement>('button[data-tab]');
    btns.forEach(b => b.setAttribute('aria-selected', b.dataset.tab === 'ficha' ? 'true' : 'false'));
  }
}

function seleccionarChat(leadId: string): void {
  actualizar({ leadActivo: leadId });
  const panelChat = document.getElementById('panel-chat');
  if (panelChat) {
    panelChat.scrollIntoView({ behavior: 'smooth' });
  }
}

async function cargarYRenderizarFicha(contenedor: HTMLElement, leadId: string): Promise<void> {
  const st = obtener();
  const leadEnCola = st.cola?.leads.find(l => l.lead_id === leadId);
  const nombreFallback = leadEnCola?.nombre ?? 'Lead';

  try {
    const ficha = await api.ficha(leadId);
    renderFicha(contenedor, ficha, nombreFallback);
  } catch (err) {
    if (err instanceof ErrorAPI && err.estadoHTTP === 404) {
      renderFicha(contenedor, null, nombreFallback);
    } else {
      contenedor.innerHTML = `<div class="banda-advertencia">⚠️ Error cargando la ficha comercial: ${(err as Error).message}</div>`;
    }
  }
}

async function cargarYRenderizarGerencia(contenedor: HTMLElement, proyectoId: string): Promise<void> {
  try {
    const bp = await api.buyerPersona(proyectoId);
    renderGerencia(contenedor, bp, proyectoId, id => {
      proyectoGerenciaSeleccionado = id;
      cargarYRenderizarGerencia(contenedor, id);
    });
  } catch {
    renderGerencia(contenedor, null, proyectoId, id => {
      proyectoGerenciaSeleccionado = id;
      cargarYRenderizarGerencia(contenedor, id);
    });
  }
}


function montarBotoneraDemo(contenedor: HTMLElement): void {
  renderBotoneraDemo(
    contenedor,
    async () => {
      const hasta = await solicitarTexto({
        icono: '⏩',
        titulo: 'Avanzar Tiempo Simulado',
        mensaje: 'Ingresa la fecha a la cual deseas avanzar la simulación (Formato ISO AAAA-MM-DD):',
        valorDefecto: '2026-08-01',
        textoConfirmar: 'Avanzar Tiempo',
        textoCancelar: 'Cancelar',
      });

      if (hasta) {
        try {
          await api.avanzarTiempo(hasta);
          await mostrarNotificacion({
            icono: '✅',
            titulo: 'Tiempo Avanzado',
            mensaje: `El reloj de simulación avanzó correctamente a ${hasta}.`,
          });
          cargarCola();
        } catch (e) {
          await mostrarNotificacion({
            icono: '⚠️',
            titulo: 'Error al Avanzar Tiempo',
            mensaje: (e as Error).message,
          });
        }
      }
    },
    async () => {
      const confirmado = await mostrarConfirmacion({
        icono: '🔄',
        titulo: '¿Reiniciar Demostración?',
        mensaje: 'Esta acción restaurará la cola de leads y los datos de prueba a su estado inicial.',
        textoConfirmar: 'Sí, Reiniciar Demo',
        textoCancelar: 'Cancelar',
        tipoBoton: 'advertencia',
      });

      if (confirmado) {
        try {
          await api.reiniciar();
          actualizar({ leadActivo: 'mock-1', tabActiva: 'cola' });
          cargarCola();
          await mostrarNotificacion({
            icono: '✅',
            titulo: 'Demo Reiniciado',
            mensaje: 'El estado de la demostración se ha restaurado con éxito.',
          });
        } catch (e) {
          await mostrarNotificacion({
            icono: '⚠️',
            titulo: 'Error al Reiniciar Demo',
            mensaje: (e as Error).message,
          });
        }
      }
    },
  );
}
