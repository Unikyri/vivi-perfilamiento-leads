import type { Conversacion, ColaLeads, Mensaje } from '../models/tipos';

let msgCounter = 1;
let turnoSimulado = false;

const mensajesDemo: Mensaje[] = [
  {
    mensaje_id: 'm1', autor: 'VIVI', tipo_contenido: 'TEXTO',
    texto: '¡Hola Ana! 👋 Como afiliada tienes un subsidio de hasta $52,5M. ¿Sueñas con comprar este año?',
    creado_en: new Date().toISOString(), adjunto: null,
  },
];

function obtenerConversacion(): Conversacion {
  return {
    lead_id: 'mock-1', estado: 'PERFILANDO',
    turno_en_proceso: turnoSimulado,
    mensajes: [...mensajesDemo],
  };
}

const colaDemo: ColaLeads = {
  cupo_10: { usados: 1, porcentaje_ventana: 10 },
  leads: [{
    lead_id: 'mock-1', nombre: 'Ana Rodríguez', estado: 'ENTREGADO', ruta: 'ASESOR',
    afiliado: true, semaforo: 'VERDE', prioridad: 0.91,
    resumen: 'Afiliada cat. A · presupuesto $166.8M · intención alta',
    actualizado_en: new Date().toISOString(),
  }],
};

function simularRespuesta(textoLead: string): void {
  turnoSimulado = true;

  setTimeout(() => {
    msgCounter++;

    // Si el lead pregunta por proyectos, mostrar tarjetas
    const hablaDeProyectos = /proyecto|comprar|vivienda|casa|apto/i.test(textoLead);

    const respuesta: Mensaje = {
      mensaje_id: `m${msgCounter}`, autor: 'VIVI', tipo_contenido: hablaDeProyectos ? 'TARJETAS_PROYECTOS' : 'TEXTO',
      texto: hablaDeProyectos
        ? 'Basándome en tu perfil, estos proyectos te pueden interesar:'
        : '¡Qué bueno saberlo! ¿Te interesaría que exploremos opciones de vivienda según tu presupuesto?',
      creado_en: new Date().toISOString(),
      adjunto: hablaDeProyectos ? {
        recomendaciones: [
          {
            proyecto_id: 'mongui', nombre: 'Monguí', zona: 'Ciudadela Maiporé - Soacha',
            precio_desde: 156470000, razon: 'Tu presupuesto cubre el 100% de la cuota inicial',
            vecinos: 622, tasa_desistimiento: 0.12,
            brochure_url: 'https://heyzine.com/flip-book/866af8f6a6.html',
            recorrido_360_url: 'https://storage.net-fs.com/hosting/7532170/19/',
          },
          {
            proyecto_id: 'macarena', nombre: 'La Macarena', zona: 'Ciudadela Maiporé - Soacha',
            precio_desde: 128340000, razon: 'El más económico de la zona, ideal para tu ingreso',
            vecinos: 374, tasa_desistimiento: 0.08,
            brochure_url: 'https://heyzine.com/flip-book/b168b2f5ba.html',
            recorrido_360_url: '',
          },
          {
            proyecto_id: 'versalles', nombre: 'Versalles', zona: 'Ciudadela Maiporé - Soacha',
            precio_desde: 195200000, razon: 'Certificación EDGE, ahorro en servicios',
            vecinos: 174, tasa_desistimiento: 0.15,
            brochure_url: 'https://heyzine.com/flip-book/be784b0d5c.html',
            recorrido_360_url: 'https://shape.com.co/360/COLSUBSIDIO-Versalles_APTOA',
          },
        ],
      } : null,
    };

    mensajesDemo.push(respuesta);
    turnoSimulado = false;
  }, 2000);
}

/** Intercepta fetch y responde con datos de ejemplo del Contrato §3. */
export function activarMock(): void {
  const original = window.fetch;
  window.fetch = async (entrada: RequestInfo | URL, opciones?: RequestInit) => {
    const url = entrada.toString();
    const json = (v: unknown, status = 200) =>
      new Response(JSON.stringify(v), { status, headers: { 'Content-Type': 'application/json' } });

    // GET conversación
    if (url.includes('/conversacion') && (!opciones?.method || opciones.method === 'GET')) {
      return json(obtenerConversacion());
    }

    // POST mensaje (texto o audio)
    if (url.includes('/mensajes') && opciones?.method === 'POST') {
      msgCounter++;
      const body = JSON.parse(opciones.body as string);
      const esAudio = body.tipo === 'AUDIO';

      // Agregar mensaje del lead
      mensajesDemo.push({
        mensaje_id: `m${msgCounter}`, autor: 'LEAD',
        tipo_contenido: 'TEXTO',
        texto: esAudio ? '🎙️ [Nota de voz]' : body.texto,
        creado_en: new Date().toISOString(),
        adjunto: esAudio ? { audio_original: true } : null,
      });

      // Simular respuesta de Vivi
      simularRespuesta(esAudio ? 'audio' : body.texto);

      return json({ mensaje_id: `m${msgCounter}`, turno_en_proceso: true }, 201);
    }

    // POST crear lead
    if (url.endsWith('/api/leads') && opciones?.method === 'POST') {
      return json({ lead_id: 'mock-1', estado: 'PERFILANDO', afiliado_detectado: true }, 201);
    }

    // GET cola
    if (url.endsWith('/api/leads') && (!opciones?.method || opciones.method === 'GET')) {
      return json(colaDemo);
    }

    return original(entrada, opciones);
  };
  console.info('[mock] servidor simulado activo');
}
