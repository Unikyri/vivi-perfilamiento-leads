import { api, ErrorAPI } from '../models/api';
import { obtener, actualizar, suscribir } from '../models/estado';
import { renderCola } from '../views/cola';
import { renderFicha } from '../views/ficha';
import { renderGerencia } from '../views/gerencia';
import { renderBotoneraDemo } from '../views/demo';

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

function renderTabActiva(contenedor: HTMLElement): void {
  const st = obtener();

  switch (st.tabActiva) {
    case 'cola':
      if (st.cola) {
        renderCola(contenedor, st.cola, leadId => seleccionarLead(leadId));
      } else {
        contenedor.innerHTML = '<div style="padding:1rem; color:#6B7280">Cargando cola de leads…</div>';
      }
      break;

    case 'ficha':
      if (st.leadActivo) {
        cargarYRenderizarFicha(contenedor, st.leadActivo);
      } else {
        contenedor.innerHTML = '<div style="padding:1.5rem; text-align:center; color:#6B7280">Selecciona un lead de la cola para ver su ficha comercial.</div>';
      }
      break;

    case 'gerencia':
      cargarYRenderizarGerencia(contenedor, proyectoGerenciaSeleccionado);
      break;
  }
}

async function seleccionarLead(leadId: string): Promise<void> {
  actualizar({ leadActivo: leadId, tabActiva: 'ficha' });

  // Actualizar el estado visual de la tab en el DOM
  const navTabs = document.querySelector('.tabs');
  if (navTabs) {
    const btns = navTabs.querySelectorAll<HTMLButtonElement>('button[data-tab]');
    btns.forEach(b => b.setAttribute('aria-selected', b.dataset.tab === 'ficha' ? 'true' : 'false'));
  }
}

async function cargarYRenderizarFicha(contenedor: HTMLElement, leadId: string): Promise<void> {
  // Buscar nombre en la cola si existe para fallback
  const st = obtener();
  const leadEnCola = st.cola?.leads.find(l => l.lead_id === leadId);
  const nombreFallback = leadEnCola?.nombre ?? 'Lead';

  try {
    const ficha = await api.ficha(leadId);
    renderFicha(contenedor, ficha, nombreFallback);
  } catch (err) {
    if (err instanceof ErrorAPI && err.estadoHTTP === 404) {
      // 404 FICHA_NO_DISPONIBLE: estado vacío amable (sin error rojo)
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
    // Si falla el endpoint, mostrar gerencia con datos por defecto
    renderGerencia(contenedor, null, proyectoId, id => {
      proyectoGerenciaSeleccionado = id;
      cargarYRenderizarGerencia(contenedor, id);
    });
  }
}

function montarBotoneraDemo(contenedor: HTMLElement): void {
  const render = () => {
    const st = obtener();
    renderBotoneraDemo(
      contenedor,
      st.leadActivo,
      leadId => seleccionarLead(leadId),
      async () => {
        const hasta = prompt('Avanzar fecha/tiempo simulado (formato ISO, ej: 2026-08-01):', '2026-08-01');
        if (hasta) {
          try {
            await api.avanzarTiempo(hasta);
            alert(`Tiempo avanzado exitosamente a ${hasta}.`);
            cargarCola();
          } catch (e) {
            alert(`Error al avanzar tiempo: ${(e as Error).message}`);
          }
        }
      },
      async () => {
        try {
          await api.reiniciar();
          actualizar({ leadActivo: 'mock-1', tabActiva: 'cola' });
          cargarCola();
          alert('Demo reiniciado al estado inicial.');
        } catch (e) {
          alert(`Error al reiniciar demo: ${(e as Error).message}`);
        }
      },
    );
  };

  suscribir(render);
  render();
}
