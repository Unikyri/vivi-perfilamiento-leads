import './estilos/base.css';
import './estilos/whatsapp.css';
import './estilos/marca.css';
import { activarMock } from './mock/servidor-mock';
import { actualizar } from './models/estado';
import { iniciarChat } from './controllers/chat';
import { iniciarDashboard } from './controllers/dashboard';

if (new URLSearchParams(location.search).get('mock') === '1') {
  activarMock();
}

// Montar el chat en el panel izquierdo (#25).
const panelChat = document.getElementById('panel-chat');
if (panelChat) {
  actualizar({ leadActivo: 'mock-1' });
  iniciarChat(panelChat);
}

// Montar el panel del asesor en el panel derecho (#26).
const contenidoTab = document.getElementById('contenido-tab');
const navTabs = document.querySelector<HTMLElement>('.tabs');
const botoneraDemo = document.getElementById('botonera-demo');

if (contenidoTab) {
  iniciarDashboard(contenidoTab, botoneraDemo, navTabs);
}

console.info('Vivi web iniciado completo (Chat + Dashboard)');
