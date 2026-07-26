import './estilos/consola.css';
import { api, ErrorAPI } from './models/api';
import { actualizar } from './models/estado';
import { iniciarChat } from './controllers/chat';
import { iniciarDashboard } from './controllers/dashboard';

const params = new URLSearchParams(location.search);

/** Precargado del Contrato §3.1: "ana" | "carlos" | "luisa". */
const PRECARGADOS = ['ana', 'carlos', 'luisa'] as const;
type Precargado = (typeof PRECARGADOS)[number];

function precargadoPedido(): Precargado {
  const p = params.get('precargado');
  return (PRECARGADOS as readonly string[]).includes(p ?? '') ? (p as Precargado) : 'ana';
}

const app = document.querySelector<HTMLDivElement>('#app');
if (!app) throw new Error('No se encontró el contenedor de Vivi');

app.innerHTML = `
  <div class="shell">
    <header class="topbar">
      <div class="brand">
        <img src="/logo-colsubsidio-blanco.png" alt="Colsubsidio">
        <span class="v-divider"></span>
        <span class="vivi-title"><strong>Vivi <span class="sun">☀</span></strong>Asesora de vivienda</span>
      </div>
      <div class="demo-actions" id="botonera-demo"></div>
      <div class="account">
        <span class="bell">🔔<span class="notification">3</span></span>
        <span class="avatar account-avatar" aria-hidden="true">AG</span>
        <span>Ana Gómez</span>
        <span class="chevron">⌄</span>
      </div>
    </header>
    <main class="workspace">
      <aside class="sidebar" aria-label="Navegación principal">
        <nav class="nav-main">
          <button class="nav-item active" data-seccion="leads"><span class="nav-icon">👤</span><span class="nav-label">Leads</span></button>
          <button class="nav-item" data-seccion="conversaciones"><span class="nav-icon">💬</span><span class="nav-label">Conversaciones</span></button>
          <button class="nav-item" data-seccion="nutricion"><span class="nav-icon">🌱</span><span class="nav-label">Nutrición</span></button>
          <button class="nav-item" data-seccion="proyectos"><span class="nav-icon">🏢</span><span class="nav-label">Proyectos</span></button>
        </nav>
        <div class="side-promo">
          <span class="promo-icon">🏡</span>
          <h2>Convertimos<br>leads en<br>vecinos</h2>
          <p class="promo-copy">Estamos para ayudarte<br>a crear hogares felices.</p>
        </div>
      </aside>
      <section class="leads-panel" id="leads-panel" aria-label="Lista de leads"></section>
      <section class="chat-panel" id="panel-chat" aria-label="Conversación"></section>
      <section class="details" id="details-panel" aria-label="Ficha comercial del lead"></section>
    </main>
  </div>
`;

/** Crea el lead en el backend y devuelve su lead_id real. El saludo
 * (RF-M4-01) llega después por el polling de la conversación. */
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
  aviso.setAttribute('role', 'alert');
  aviso.style.cssText = 'margin:12px;padding:12px;border-radius:8px;background:#fff1c9;color:#7a4b00;font-size:13px;font-weight:600';
  aviso.textContent = `No se pudo iniciar la conversación (${detalle}). Recargá la página o revisá que el backend esté arriba.`;
  panel.prepend(aviso);
}

async function arrancar(): Promise<void> {
  const leadsPanel = document.getElementById('leads-panel')!;
  const chatPanel = document.getElementById('panel-chat')!;
  const detailsPanel = document.getElementById('details-panel')!;
  const botoneraDemo = document.getElementById('botonera-demo');
  const nav = document.querySelector<HTMLElement>('.nav-main');

  const leadId = await crearLeadInicial();
  if (leadId) actualizar({ leadActivo: leadId });

  iniciarChat(chatPanel);
  iniciarDashboard(leadsPanel, detailsPanel, botoneraDemo, nav);

  console.info('Vivi web iniciado (Leads + Chat + Ficha, conectados a la API real)');
}

void arrancar();
