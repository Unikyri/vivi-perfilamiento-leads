import './estilos/base.css';
import './estilos/whatsapp.css';
import { activarMock } from './mock/servidor-mock';
import { actualizar } from './models/estado';
import { iniciarChat } from './controllers/chat';

if (new URLSearchParams(location.search).get('mock') === '1') {
  activarMock();
}

// Montar el chat en el panel izquierdo (#25).
const panelChat = document.getElementById('panel-chat');
if (panelChat) {
  // En modo demo/mock, crear un lead automáticamente
  actualizar({ leadActivo: 'mock-1' });
  iniciarChat(panelChat);
}

// Las vistas del dashboard se montan en la issue #26.
console.info('Vivi web iniciado');
