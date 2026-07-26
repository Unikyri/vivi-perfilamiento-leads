import type { Mensaje, Recomendacion } from '../models/tipos';
import { avatarPersona, avatarProyecto } from '../util/avatares';

export function renderMensajes(contenedor: HTMLElement, mensajes: Mensaje[]): void {
  contenedor.innerHTML = mensajes.map(renderMensaje).join('');
}

function renderMensaje(m: Mensaje): string {
  if (m.tipo_contenido === 'SISTEMA') {
    return `<div class="pildora-sistema">${escapar(m.texto)}</div>`;
  }
  const lado = m.autor === 'LEAD' ? 'derecha' : 'izquierda';
  const hora = new Date(m.creado_en).toLocaleTimeString('es-CO', { hour: '2-digit', minute: '2-digit' });
  const chulos = m.autor === 'LEAD' ? '<span class="chulos">✓✓</span>' : '';
  const microfono = m.adjunto?.audio_original ? '<span class="icono-audio" aria-label="nota de voz">🎙️</span>' : '';
  const tarjetas = m.adjunto?.recomendaciones ? renderTarjetas(m.adjunto.recomendaciones) : '';
  return `
    <div class="burbuja ${lado}">
      ${microfono}<p>${escapar(m.texto)}</p>
      <span class="hora">${hora}${chulos}</span>
    </div>${tarjetas}`;
}

/** Máximo 3 tarjetas, carrusel horizontal (doc 12 §2.3, RF-M4-05). */
function renderTarjetas(recs: Recomendacion[]): string {
  return `<div class="carrusel" role="list">${recs.slice(0, 3).map(r => `
    <article class="tarjeta-proyecto" role="listitem">
      <div class="tarjeta-imagen" style="background-image:url('${avatarProyecto(r.proyecto_id)}')" aria-hidden="true"></div>
      <header class="franja-azul">${escapar(r.nombre)}</header>
      <p class="zona">${escapar(r.zona)}</p>
      <p class="precio">Desde $${(r.precio_desde / 1_000_000).toFixed(0)}M</p>
      <p class="razon">${escapar(r.razon)}</p>
      <p class="evidencia">${r.vecinos} personas con tu perfil compraron aquí ·
         ${(r.tasa_desistimiento * 100).toFixed(0)}% desistió</p>
      <a class="btn-primario" href="${encodeURI(r.brochure_url)}" target="_blank" rel="noopener">Ver brochure</a>
      <a class="btn-secundario" href="${encodeURI(r.recorrido_360_url)}" target="_blank" rel="noopener">Recorrido 360°</a>
    </article>`).join('')}</div>`;
}

export function renderEscribiendo(el: HTMLElement, activo: boolean): void {
  el.classList.toggle('visible', activo);
}

/** Renderiza la estructura HTML inicial del panel de chat (se monta una vez). */
export function renderShellChat(panel: HTMLElement): void {
  panel.innerHTML = `
    <header class="chat-top">
      <button class="back" aria-label="Volver">‹</button>
      <img class="avatar chat-avatar" id="chat-avatar" src="" alt="" aria-hidden="true" style="display:none">
      <div class="chat-name">
        <span id="chat-nombre-lead">Selecciona un lead</span>
        <small>en línea <span class="online">●</span></small>
      </div>
      <div class="chat-tools">
        <button type="button" aria-label="Llamar">☎</button>
        <button type="button" aria-label="Videollamar">▣</button>
        <button type="button" aria-label="Más acciones">⋮</button>
      </div>
    </header>
    <div class="mensajes-scroll" id="mensajes-scroll">
      <div class="mensajes" id="contenedor-mensajes"></div>
      <button class="btn-nuevos" id="btn-nuevos" type="button">↓ nuevos mensajes</button>
    </div>
    <span class="escribiendo" id="indicador-escribiendo" style="padding:0 16px">escribiendo…</span>
    <div class="barra-entrada">
      <button class="btn-mic" id="btn-mic" aria-label="Grabar nota de voz" type="button">🎤</button>
      <div class="mic-grabando" id="mic-grabando">
        <span class="punto-rojo"></span>
        <span class="contador-mic" id="contador-mic">0:00</span>
        <button class="btn-detener-mic" id="btn-detener-mic" type="button">■</button>
      </div>
      <input type="text" class="input-mensaje" id="input-mensaje"
             placeholder="Escribí un mensaje..." autocomplete="off" />
      <button class="btn-enviar" id="btn-enviar" aria-label="Enviar mensaje" type="button">➤</button>
    </div>
  `;
}

/** Actualiza sólo el nombre/avatar del encabezado — no remonta el panel
 * entero al cambiar de lead activo (el polling sigue corriendo). `leadId`
 * es la semilla del avatar generado: el mismo lead siempre muestra el
 * mismo avatar, tanto acá como en la lista de leads y en la ficha. */
export function actualizarCabeceraChat(leadId: string, nombre: string): void {
  const nombreEl = document.getElementById('chat-nombre-lead');
  const avatarEl = document.getElementById('chat-avatar') as HTMLImageElement | null;
  if (nombreEl) nombreEl.textContent = nombre;
  if (avatarEl) {
    avatarEl.src = avatarPersona(leadId);
    avatarEl.style.display = '';
  }
}

/** SIEMPRE escapar: el texto viene del LLM y del usuario (anti-XSS). */
function escapar(s: string): string {
  const d = document.createElement('div');
  d.textContent = s;
  return d.innerHTML;
}
