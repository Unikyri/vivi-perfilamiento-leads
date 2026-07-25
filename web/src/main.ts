import './estilos/base.css';
import { activarMock } from './mock/servidor-mock';

if (new URLSearchParams(location.search).get('mock') === '1') {
  activarMock();
}

// Las vistas se montan en las issues #25 (chat) y #26 (dashboard).
console.info('Vivi web iniciado');
