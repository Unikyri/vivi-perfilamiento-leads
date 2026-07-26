import './estilos/base.css';
import './estilos/whatsapp.css';
import './estilos/marca.css';
import { activarMock } from './mock/servidor-mock';
import { actualizar } from './models/estado';
import { api, ErrorAPI } from './models/api';
import { iniciarChat } from './controllers/chat';
import { iniciarDashboard } from './controllers/dashboard';

const params = new URLSearchParams(location.search);

if (params.get('mock') === '1') {
  activarMock();
}

/** Precargado del Contrato §3.1: "ana" | "carlos" | "luisa". */
const PRECARGADOS = ['ana', 'carlos', 'luisa'] as const;
type Precargado = (typeof PRECARGADOS)[number];

function precargadoPedido(): Precargado {
  const p = params.get('precargado');
  return (PRECARGADOS as readonly string[]).includes(p ?? '') ? (p as Precargado) : 'ana';
}

/**
 * Crea el lead en el backend y devuelve su lead_id real.
 * El saludo (RF-M4-01) llega después por el polling de la conversación.
 */
async function crearLeadInicial(): Promise<string | null> {
  try {
    const r = await api.crearLead(precargadoPedido());
    return r.lead_id;
  } catch (err) {
    const detalle = err instanceof ErrorAPI ? `${err.codigo}: ${err.message}` : String(err);
    console.error('[main] no se pudo crear el lead inicial:', detalle);
    mostrarFalloArranque(detalle);
    return null;
  }
}

function mostrarFalloArranque(detalle: string): void {
  const panel = document.getElementById('panel-chat');
  if (!panel) return;
  const aviso = document.createElement('div');
  aviso.className = 'banda-advertencia';
  aviso.setAttribute('role', 'alert');
  aviso.textContent = `No se pudo iniciar la conversación (${detalle}). Recargá la página o revisá que el backend esté arriba.`;
  panel.prepend(aviso);
}

async function arrancar(): Promise<void> {
  // Montar el chat en el panel izquierdo (#25).
  const panelChat = document.getElementById('panel-chat');

  if (panelChat) {
    const leadId = await crearLeadInicial();
    if (leadId) {
      actualizar({ leadActivo: leadId });
      iniciarChat(panelChat);
    }
  }

  // Montar el panel del asesor en el panel derecho (#26).
  const contenidoTab = document.getElementById('contenido-tab');
  const navTabs = document.querySelector<HTMLElement>('.tabs');
  const botoneraDemo = document.getElementById('botonera-demo');

  if (contenidoTab) {
    iniciarDashboard(contenidoTab, botoneraDemo, navTabs);
  }

  console.info('Vivi web iniciado completo (Chat + Dashboard)');
}

/** Sidebars contraíbles: cada panel se colapsa a una franja angosta. */
function iniciarColapsoPaneles(): void {
  document.querySelectorAll<HTMLButtonElement>('[data-colapsar]').forEach(btn => {
    const panel = document.getElementById(btn.dataset.colapsar!);
    if (!panel) return;
    btn.addEventListener('click', () => {
      const colapsado = panel.classList.toggle('panel-colapsado');
      btn.setAttribute('aria-expanded', String(!colapsado));
    });
  });
}

iniciarColapsoPaneles();
void arrancar();
