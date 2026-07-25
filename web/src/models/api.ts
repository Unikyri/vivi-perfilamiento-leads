import type { ColaLeads, Conversacion, Ficha, BuyerPersona, RespuestaError, EstadoLead } from './tipos';

const BASE = '/api';

/** ErrorAPI conserva el código del Contrato §3 para que las vistas decidan qué mostrar. */
export class ErrorAPI extends Error {
  constructor(public codigo: string, mensaje: string, public estadoHTTP: number) {
    super(mensaje);
    this.name = 'ErrorAPI';
  }
}

async function pedir<T>(ruta: string, opciones?: RequestInit): Promise<T> {
  const r = await fetch(`${BASE}${ruta}`, {
    headers: { 'Content-Type': 'application/json' },
    ...opciones,
  });
  if (!r.ok) {
    const cuerpo = (await r.json().catch(() => null)) as RespuestaError | null;
    throw new ErrorAPI(
      cuerpo?.error?.codigo ?? 'ERROR_INTERNO',
      cuerpo?.error?.mensaje ?? `HTTP ${r.status}`,
      r.status,
    );
  }
  return (await r.json()) as T;
}

export const api = {
  crearLead: (precargadoId: string) =>
    pedir<{ lead_id: string; estado: EstadoLead; afiliado_detectado: boolean }>('/leads', {
      method: 'POST',
      body: JSON.stringify({ precargado_id: precargadoId, fuente: 'DEMO' }),
    }),

  enviarTexto: (leadId: string, texto: string) =>
    pedir<{ mensaje_id: string; turno_en_proceso: boolean }>(`/leads/${leadId}/mensajes`, {
      method: 'POST',
      body: JSON.stringify({ tipo: 'TEXTO', texto }),
    }),

  enviarAudio: (leadId: string, audioBase64: string, mime: string, duracionS: number) =>
    pedir<{ mensaje_id: string }>(`/leads/${leadId}/mensajes`, {
      method: 'POST',
      body: JSON.stringify({ tipo: 'AUDIO', audio_base64: audioBase64, mime, duracion_s: duracionS }),
    }),

  conversacion: (leadId: string) => pedir<Conversacion>(`/leads/${leadId}/conversacion`),
  cola: () => pedir<ColaLeads>('/leads'),
  ficha: (leadId: string) => pedir<Ficha>(`/leads/${leadId}/ficha`),
  buyerPersona: (proyectoId?: string) =>
    pedir<BuyerPersona>(`/gerencia/buyer-persona${proyectoId ? `?proyecto_id=${proyectoId}` : ''}`),

  avanzarTiempo: (hasta: string) =>
    pedir<{ fecha_simulada: string; hitos_disparados: number }>('/demo/tiempo', {
      method: 'POST', body: JSON.stringify({ avanzar_hasta: hasta }),
    }),

  reiniciar: () =>
    pedir<{ reiniciado: boolean; fecha_simulada: string }>('/demo/reiniciar', { method: 'POST' }),
};
