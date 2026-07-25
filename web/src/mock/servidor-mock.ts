import type { Conversacion, ColaLeads } from '../models/tipos';

const conversacionDemo: Conversacion = {
  lead_id: 'mock-1', estado: 'PERFILANDO', turno_en_proceso: false,
  mensajes: [{
    mensaje_id: 'm1', autor: 'VIVI', tipo_contenido: 'TEXTO',
    texto: '¡Hola Ana! 👋 Como afiliada tienes un subsidio de hasta $52,5M. ¿Sueñas con comprar este año?',
    creado_en: new Date().toISOString(), adjunto: null,
  }],
};

const colaDemo: ColaLeads = {
  cupo_10: { usados: 1, porcentaje_ventana: 10 },
  leads: [{
    lead_id: 'mock-1', nombre: 'Ana Rodríguez', estado: 'ENTREGADO', ruta: 'ASESOR',
    afiliado: true, semaforo: 'VERDE', prioridad: 0.91,
    resumen: 'Afiliada cat. A · presupuesto $166.8M · intención alta',
    actualizado_en: new Date().toISOString(),
  }],
};

/** Intercepta fetch y responde con datos de ejemplo del Contrato §3. */
export function activarMock(): void {
  const original = window.fetch;
  window.fetch = async (entrada: RequestInfo | URL, opciones?: RequestInit) => {
    const url = entrada.toString();
    const json = (v: unknown) =>
      new Response(JSON.stringify(v), { status: 200, headers: { 'Content-Type': 'application/json' } });

    if (url.includes('/conversacion')) return json(conversacionDemo);
    if (url.endsWith('/api/leads')) return json(colaDemo);
    if (url.includes('/api/leads') && opciones?.method === 'POST')
      return new Response(JSON.stringify({ lead_id: 'mock-1', estado: 'PERFILANDO', afiliado_detectado: true }), { status: 201 });
    return original(entrada, opciones);
  };
  console.info('[mock] servidor simulado activo');
}
