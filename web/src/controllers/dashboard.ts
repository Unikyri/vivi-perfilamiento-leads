import { api, ErrorAPI } from '../models/api';
import { obtener, actualizar, suscribir } from '../models/estado';
import { renderCola } from '../views/cola';
import { renderFicha } from '../views/ficha';
import { renderGerencia } from '../views/gerencia';
import { renderBotoneraDemo } from '../views/demo';

const INTERVALO_POLL_COLA_MS = 5000;
let pollColaTimer: ReturnType<typeof setInterval> | null = null;
let proyectoGerenciaSeleccionado = 'mongui';

// Evita refetchear Ficha/Gerencia en cada poll de la cola (cada 5 s) cuando
// nada relevante cambió. Se invalidan explícitamente al salir de la pestaña
// o al cambiar de lead/proyecto — nunca queda un dato congelado.
let ultimaFichaCargada: string | null = null;
let ultimaGerenciaCargada: string | null = null;

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

function renderTabActiva(contenedor: HTMLElement): void {
  const st = obtener();

  switch (st.tabActiva) {
    case 'cola':
      // Al volver a Cola, invalidar: si el asesor vuelve a Ficha/Gerencia
      // después, el dato debe refrescarse, no quedar congelado.
      ultimaFichaCargada = null;
      ultimaGerenciaCargada = null;
      if (st.cola) {
        renderCola(
          contenedor,
          st.cola,
          leadId => seleccionarLead(leadId),
          leadId => seleccionarChat(leadId),
        );
      } else {
        contenedor.innerHTML = '<div style="padding:1rem; color:#6B7280">Cargando cola de leads…</div>';
      }
      break;

    case 'ficha':
      if (!st.leadActivo) {
        ultimaFichaCargada = null;
        contenedor.innerHTML = '<div style="padding:1.5rem; text-align:center; color:#6B7280">Selecciona un lead de la cola para ver su ficha comercial.</div>';
        break;
      }
      if (st.leadActivo !== ultimaFichaCargada) {
        ultimaFichaCargada = st.leadActivo;
        cargarYRenderizarFicha(contenedor, st.leadActivo);
      }
      break;

    case 'gerencia':
      if (proyectoGerenciaSeleccionado !== ultimaGerenciaCargada) {
        ultimaGerenciaCargada = proyectoGerenciaSeleccionado;
        cargarYRenderizarGerencia(contenedor, proyectoGerenciaSeleccionado);
      }
      break;
  }
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
  const alCambiarProyecto = (id: string) => {
    proyectoGerenciaSeleccionado = id;
    ultimaGerenciaCargada = id;
    cargarYRenderizarGerencia(contenedor, id);
  };

  try {
    const bp = await api.buyerPersona(proyectoId);
    renderGerencia(contenedor, bp, proyectoId, alCambiarProyecto);
  } catch {
    renderGerencia(contenedor, null, proyectoId, alCambiarProyecto);
  }
}

/** Aviso inline en la botonera de demo: nunca abre una ventana nativa
 * (prompt/alert), que en pantalla compartida puede renderizarse fuera del
 * área capturada o quedar detrás de la ventana. */
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
      const input = document.getElementById('demo-fecha') as HTMLInputElement | null;
      const hasta = input?.value;
      if (!hasta) {
        avisarDemo('Elegí una fecha primero.', true);
        return;
      }
      try {
        const r = await api.avanzarTiempo(hasta);
        avisarDemo(`Tiempo avanzado a ${hasta} · ${r.hitos_disparados} hito(s) disparado(s).`);
        cargarCola();
      } catch (e) {
        avisarDemo(`Error al avanzar tiempo: ${(e as Error).message}`, true);
      }
    },
    async () => {
      try {
        await api.reiniciar();

        // Tras el reinicio el lead anterior ya no existe: hay que crear uno nuevo.
        let nuevoLead: string | null = null;
        try {
          nuevoLead = (await api.crearLead('ana')).lead_id;
        } catch (e) {
          console.error('[dashboard] no se pudo recrear el lead tras reiniciar:', e);
        }

        actualizar({ leadActivo: nuevoLead, tabActiva: 'cola' });
        cargarCola();
        avisarDemo('Demo reiniciado al estado inicial.');
      } catch (e) {
        avisarDemo(`Error al reiniciar demo: ${(e as Error).message}`, true);
      }
    },
  );
}
