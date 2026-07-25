import { api, ErrorAPI } from '../models/api';
import { actualizar, obtener, suscribir } from '../models/estado';
import { renderMensajes, renderEscribiendo, renderShellChat, renderHeaderChat } from '../views/chat';

const MAX_SEGUNDOS_AUDIO = 60;
const INTERVALO_POLL_MS = 1500;
const RETRASO_ESCRIBIENDO_MS = 300;

let pollTimer: ReturnType<typeof setInterval> | null = null;
let escribiendoTimer: ReturnType<typeof setTimeout> | null = null;
let scrollManual = false;

// ── Grabación de audio ────────────────────────────────────────
let mediaRecorder: MediaRecorder | null = null;
let audioChunks: Blob[] = [];
let micTimer: ReturnType<typeof setInterval> | null = null;
let micSegundos = 0;

/** Monta el chat en el panel y arranca el polling. */
export function iniciarChat(panel: HTMLElement): void {
  renderShellChat(panel);

  const contenedor = document.getElementById('contenedor-mensajes')!;
  const scrollEl = document.getElementById('mensajes-scroll')!;
  const indicador = document.getElementById('indicador-escribiendo')!;
  const input = document.getElementById('input-mensaje') as HTMLInputElement;
  const btnEnviar = document.getElementById('btn-enviar')!;
  const btnMic = document.getElementById('btn-mic')!;
  const micGrabando = document.getElementById('mic-grabando')!;
  const contadorMic = document.getElementById('contador-mic')!;
  const btnDetenerMic = document.getElementById('btn-detener-mic')!;
  const btnNuevos = document.getElementById('btn-nuevos')!;

  // ── Envío con Enter o botón ➤ ────────────────────────────
  input.addEventListener('keydown', (e: KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      enviar(input);
    }
  });

  btnEnviar.addEventListener('click', () => enviar(input));

  // ── Micrófono ────────────────────────────────────────────
  btnMic.addEventListener('click', () => iniciarGrabacion(btnMic, micGrabando, contadorMic, input, btnEnviar));
  btnDetenerMic.addEventListener('click', () => detenerGrabacion(btnMic, micGrabando, input, btnEnviar));

  // ── Auto-scroll con detección de scroll manual ───────────
  scrollEl.addEventListener('scroll', () => {
    const distanciaAlFondo = scrollEl.scrollHeight - scrollEl.scrollTop - scrollEl.clientHeight;
    scrollManual = distanciaAlFondo > 80;
    if (!scrollManual) {
      btnNuevos.classList.remove('visible');
    }
  });

  btnNuevos.addEventListener('click', () => {
    scrollEl.scrollTop = scrollEl.scrollHeight;
    btnNuevos.classList.remove('visible');
    scrollManual = false;
  });

  // ── Polling & Reactividad de Estado ─────────────────────
  pollTimer = setInterval(() => poll(contenedor, scrollEl, indicador, btnNuevos), INTERVALO_POLL_MS);

  // Actualizar inmediatamente al cambiar de lead activo (clic en '💬 Ver chat')
  suscribir(() => poll(contenedor, scrollEl, indicador, btnNuevos));

  // Primera carga inmediata
  poll(contenedor, scrollEl, indicador, btnNuevos);
}

/** Detiene el polling (para cleanup). */
export function detenerChat(): void {
  if (pollTimer) { clearInterval(pollTimer); pollTimer = null; }
  if (escribiendoTimer) { clearTimeout(escribiendoTimer); escribiendoTimer = null; }
}

// ── Polling ─────────────────────────────────────────────────
async function poll(
  contenedor: HTMLElement,
  scrollEl: HTMLElement,
  indicador: HTMLElement,
  btnNuevos: HTMLElement,
): Promise<void> {
  const estado = obtener();
  if (!estado.leadActivo) return;

  // Actualizar nombre del lead en la cabecera de WhatsApp
  const leadEnCola = estado.cola?.leads.find(l => l.lead_id === estado.leadActivo);
  if (leadEnCola) {
    renderHeaderChat(leadEnCola.nombre);
  }

  try {
    const conv = await api.conversacion(estado.leadActivo);
    actualizar({ conversacion: conv });

    renderMensajes(contenedor, conv.mensajes);

    // "escribiendo…" desde RETRASO_ESCRIBIENDO_MS
    if (conv.turno_en_proceso) {
      if (!escribiendoTimer) {
        escribiendoTimer = setTimeout(() => renderEscribiendo(indicador, true), RETRASO_ESCRIBIENDO_MS);
      }
    } else {
      if (escribiendoTimer) { clearTimeout(escribiendoTimer); escribiendoTimer = null; }
      renderEscribiendo(indicador, false);
    }

    // Auto-scroll o botón flotante
    if (!scrollManual) {
      scrollEl.scrollTop = scrollEl.scrollHeight;
    } else {
      btnNuevos.classList.add('visible');
    }
  } catch (err) {
    if (err instanceof ErrorAPI && err.estadoHTTP >= 500) {
      mostrarToast('Error de conexión. Reintentando…');
    }
  }
}

// ── Envío de texto ──────────────────────────────────────────
async function enviar(input: HTMLInputElement): Promise<void> {
  const texto = input.value.trim();
  if (!texto) return;

  const estado = obtener();
  if (!estado.leadActivo) return;

  input.value = '';
  try {
    await api.enviarTexto(estado.leadActivo, texto);
  } catch (err) {
    // Restaurar el texto sin perderlo
    input.value = texto;
    if (err instanceof ErrorAPI) {
      mostrarToast(err.message);
    }
  }
}

// ── Grabación de audio ──────────────────────────────────────
async function iniciarGrabacion(
  btnMic: HTMLElement,
  micGrabando: HTMLElement,
  contadorMic: HTMLElement,
  input: HTMLInputElement,
  btnEnviar: HTMLElement,
): Promise<void> {
  try {
    const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
    mediaRecorder = new MediaRecorder(stream);
    audioChunks = [];
    micSegundos = 0;

    mediaRecorder.ondataavailable = (e: BlobEvent) => { audioChunks.push(e.data); };
    mediaRecorder.onstop = () => {
      stream.getTracks().forEach(t => t.stop());
      procesarAudio(btnMic, micGrabando, input, btnEnviar);
    };

    mediaRecorder.start();

    // UI: ocultar input/enviar/mic, mostrar grabando
    btnMic.style.display = 'none';
    input.style.display = 'none';
    btnEnviar.style.display = 'none';
    micGrabando.classList.add('activo');

    contadorMic.textContent = '0:00';
    micTimer = setInterval(() => {
      micSegundos++;
      const m = Math.floor(micSegundos / 60);
      const s = micSegundos % 60;
      contadorMic.textContent = `${m}:${s.toString().padStart(2, '0')}`;

      // Corte a MAX_SEGUNDOS_AUDIO
      if (micSegundos >= MAX_SEGUNDOS_AUDIO) {
        detenerGrabacion(btnMic, micGrabando, input, btnEnviar);
      }
    }, 1000);
  } catch {
    mostrarToast('No se pudo acceder al micrófono.');
  }
}

function detenerGrabacion(
  btnMic: HTMLElement,
  micGrabando: HTMLElement,
  input: HTMLInputElement,
  btnEnviar: HTMLElement,
): void {
  if (micTimer) { clearInterval(micTimer); micTimer = null; }
  if (mediaRecorder && mediaRecorder.state !== 'inactive') {
    mediaRecorder.stop();
  }
  // UI: restaurar barra de entrada
  micGrabando.classList.remove('activo');
  btnMic.style.display = '';
  input.style.display = '';
  btnEnviar.style.display = '';
}

async function procesarAudio(
  btnMic: HTMLElement,
  micGrabando: HTMLElement,
  input: HTMLInputElement,
  btnEnviar: HTMLElement,
): Promise<void> {
  if (audioChunks.length === 0) return;

  const estado = obtener();
  if (!estado.leadActivo) return;

  const blob = new Blob(audioChunks, { type: audioChunks[0].type });
  const reader = new FileReader();
  reader.onloadend = async () => {
    const base64 = (reader.result as string).split(',')[1];
    const mime = blob.type;
    const duracionS = micSegundos;

    try {
      await api.enviarAudio(estado.leadActivo!, base64, mime, duracionS);
    } catch (err) {
      if (err instanceof ErrorAPI) {
        if (err.codigo === 'AUDIO_INVALIDO') {
          mostrarToast('No te escuché bien, ¿me lo repites o me lo escribes?');
        } else {
          mostrarToast(err.message);
        }
      }
    }
  };
  reader.readAsDataURL(blob);

  // Restaurar UI por si no se hizo ya
  micGrabando.classList.remove('activo');
  btnMic.style.display = '';
  input.style.display = '';
  btnEnviar.style.display = '';
}

// ── Toast ───────────────────────────────────────────────────
function mostrarToast(mensaje: string): void {
  // Remover toast previo si existe
  const previo = document.querySelector('.toast-error');
  if (previo) previo.remove();

  const toast = document.createElement('div');
  toast.className = 'toast-error';
  toast.textContent = mensaje;
  document.body.appendChild(toast);

  setTimeout(() => toast.remove(), 4000);
}
