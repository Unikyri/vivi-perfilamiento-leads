import type { Conversacion, ColaLeads, Ficha, BuyerPersona, Mensaje } from '../models/tipos';

let msgCounter = 1;
let turnoSimulado = false;

const mensajesDemo: Mensaje[] = [
  {
    mensaje_id: 'm1', autor: 'VIVI', tipo_contenido: 'TEXTO',
    texto: '¡Hola Ana! 👋 Como afiliada tienes un subsidio de hasta $52,5M. ¿Sueñas con comprar este año?',
    creado_en: new Date().toISOString(), adjunto: null,
  },
];

function obtenerConversacion(leadId: string): Conversacion {
  return {
    lead_id: leadId, estado: 'PERFILANDO',
    turno_en_proceso: turnoSimulado,
    mensajes: [...mensajesDemo],
  };
}

const colaDemo: ColaLeads = {
  cupo_10: { usados: 1, porcentaje_ventana: 10 },
  leads: [
    {
      lead_id: 'mock-1', nombre: 'Ana Rodríguez', estado: 'ENTREGADO', ruta: 'ASESOR',
      afiliado: true, semaforo: 'VERDE', prioridad: 0.91,
      resumen: 'Afiliada cat. A · presupuesto $166.8M · intención alta',
      actualizado_en: new Date().toISOString(),
    },
    {
      lead_id: 'mock-2', nombre: 'Carlos Martínez', estado: 'PERFILANDO', ruta: 'ASESOR',
      afiliado: false, semaforo: 'AMBAR', prioridad: 0.78,
      resumen: 'No afiliado · presupuesto $210M · intención media',
      actualizado_en: new Date().toISOString(),
    },
    {
      lead_id: 'mock-3', nombre: 'Luisa Gómez', estado: 'CALIFICADO', ruta: 'NUTRICION',
      afiliado: true, semaforo: 'VERDE', prioridad: 0.85,
      resumen: 'Afiliada cat. B · presupuesto $145M · requiere plan cuota inicial',
      actualizado_en: new Date().toISOString(),
    },
  ],
};

const fichasDemo: Record<string, Ficha> = {
  'mock-1': {
    ficha_id: 'f1', lead_id: 'mock-1', generada_en: new Date().toISOString(),
    confianza_perfil: 0.94, banda_advertencia: null,
    identificacion: { nombre: 'Ana Rodríguez', afiliada: true, categoria: 'A', telefono: '+57 300 123 4567' },
    capacidad: {
      presupuesto_max: 166800000, credito_max: 114300000, subsidio_aplicable: 52500000,
      recursos_propios: 12000000, ratio: 0.28, confianza: 0.94,
      desglose: [
        { concepto: 'Subsidio Mi Casa Ya / Caja', monto: 52500000, regla: 'Afiliado Cat A', fuente: 'VERIFICADO_BASE' },
        { concepto: 'Preaprobado Bancolombia', monto: 114300000, regla: 'Capacidad de endeudamiento 30%', fuente: 'INFERIDO' },
        { concepto: 'Ahorro Declarado', monto: 12000000, regla: 'Declarado en chat', fuente: 'DECLARADO' },
      ],
    },
    perfil: {
      ingreso: { valor: 2600000, fuente: 'VERIFICADO_BASE', confianza: 0.95, requiere_confirmacion: false, actualizado_en: new Date().toISOString() },
    },
    intencion: { nivel: 'ALTA', confianza: 'ALTA', senales: ['Busca comprar antes de 6 meses', 'Tiene ahorro inicial preparado', 'Responde inmediatamente a Vivi'] },
    recomendaciones: [
      {
        proyecto_id: 'mongui', nombre: 'Monguí', zona: 'Ciudadela Maiporé - Soacha',
        precio_desde: 156470000, razon: 'Tu presupuesto cubre el 100% de la cuota inicial',
        vecinos: 622, tasa_desistimiento: 0.12,
        brochure_url: 'https://heyzine.com/flip-book/866af8f6a6.html',
        recorrido_360_url: 'https://storage.net-fs.com/hosting/7532170/19/',
      },
    ],
    beneficios: ['Subsidio de vivienda Colsubsidio hasta $52,5M', 'Tasa preferencial crédito hipotecario'],
    argumentos_venta: ['Cuota estimada mensual ($650k) es MENOR al arriendo promedio de la zona ($850k)', 'Entrega inmediata en 2026'],
    alerta_desistimiento: { activa: false, tasa_vecinos: 0.12, detalle: null },
    consume_cupo_10: false,
  },
  'mock-2': {
    ficha_id: 'f2', lead_id: 'mock-2', generada_en: new Date().toISOString(),
    confianza_perfil: 0.82, banda_advertencia: 'No afiliado a Colsubsidio — consume cupo del 10% regulatorio',
    identificacion: { nombre: 'Carlos Martínez', afiliada: false, categoria: 'N/A', telefono: '+57 311 987 6543' },
    capacidad: {
      presupuesto_max: 210000000, credito_max: 180000000, subsidio_aplicable: 0,
      recursos_propios: 30000000, ratio: 0.32, confianza: 0.82,
      desglose: [
        { concepto: 'Crédito solicitado', monto: 180000000, regla: 'Ingresos independientes', fuente: 'DECLARADO' },
        { concepto: 'Cuota Inicial Propia', monto: 30000000, regla: 'Recursos declarados', fuente: 'DECLARADO' },
      ],
    },
    perfil: {},
    intencion: { nivel: 'MEDIA', confianza: 'MEDIA', senales: ['Interesado en proyectos VIS y no VIS', 'Evaluando opciones de crédito'] },
    recomendaciones: [],
    beneficios: ['Opción de crédito con aliados de la caja'],
    argumentos_venta: ['Proyecto Versalles cuenta con certificación EDGE para ahorro energético'],
    alerta_desistimiento: { activa: false, tasa_vecinos: 0.08, detalle: null },
    consume_cupo_10: true,
  },
};

const buyerPersonaDemo: Record<string, BuyerPersona> = {
  mongui: {
    proyecto_id: 'mongui', nombre: 'Monguí', muestras: 312,
    afiliacion: { afiliados: 198, no_afiliados: 114 },
    categoria: { 'Cat A': 110, 'Cat B': 68, 'Cat C': 20, 'No Afiliado': 114 },
    rango_edad: { '18-25': 25, '26-35': 165, '36-45': 82, '46+': 40 },
    tasa_desistimiento: 0.11, actualizado_en: new Date().toISOString(),
  },
  macarena: {
    proyecto_id: 'macarena', nombre: 'La Macarena', muestras: 185,
    afiliacion: { afiliados: 140, no_afiliados: 45 },
    categoria: { 'Cat A': 85, 'Cat B': 42, 'Cat C': 13, 'No Afiliado': 45 },
    rango_edad: { '18-25': 15, '26-35': 95, '36-45': 50, '46+': 25 },
    tasa_desistimiento: 0.08, actualizado_en: new Date().toISOString(),
  },
  versalles: {
    proyecto_id: 'versalles', nombre: 'Versalles', muestras: 142,
    afiliacion: { afiliados: 85, no_afiliados: 57 },
    categoria: { 'Cat A': 40, 'Cat B': 32, 'Cat C': 13, 'No Afiliado': 57 },
    rango_edad: { '18-25': 10, '26-35': 70, '36-45': 42, '46+': 20 },
    tasa_desistimiento: 0.15, actualizado_en: new Date().toISOString(),
  },
};

function simularRespuesta(textoLead: string): void {
  turnoSimulado = true;

  setTimeout(() => {
    msgCounter++;

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
    const metodo = (opciones?.method ?? 'GET').toUpperCase();
    const json = (v: unknown, status = 200) =>
      new Response(JSON.stringify(v), { status, headers: { 'Content-Type': 'application/json' } });

    // GET conversación
    if (url.includes('/conversacion') && metodo === 'GET') {
      const m = url.match(/\/leads\/([^/]+)\/conversacion/);
      return json(obtenerConversacion(m ? m[1] : 'mock-1'));
    }

    // POST mensaje (texto o audio)
    if (url.includes('/mensajes') && metodo === 'POST') {
      msgCounter++;
      let body: { tipo?: string; texto?: string } = {};
      try { body = JSON.parse((opciones?.body as string) || '{}'); } catch { body = {}; }
      const esAudio = body.tipo === 'AUDIO';

      mensajesDemo.push({
        mensaje_id: `m${msgCounter}`, autor: 'LEAD',
        tipo_contenido: 'TEXTO',
        texto: esAudio ? '🎙️ [Nota de voz]' : (body.texto || ''),
        creado_en: new Date().toISOString(),
        adjunto: esAudio ? { audio_original: true } : null,
      });

      simularRespuesta(esAudio ? 'audio' : (body.texto || ''));
      return json({ mensaje_id: `m${msgCounter}`, turno_en_proceso: true }, 201);
    }

    // GET ficha comercial
    if (url.includes('/ficha') && metodo === 'GET') {
      const match = url.match(/\/leads\/([^/]+)\/ficha/);
      const leadId = match ? match[1] : 'mock-1';
      const ficha = fichasDemo[leadId];
      if (!ficha) {
        return json({ error: { codigo: 'FICHA_NO_DISPONIBLE', mensaje: 'Ficha aún no disponible' } }, 404);
      }
      return json(ficha);
    }

    // GET gerencia buyer persona
    if (url.includes('/gerencia/buyer-persona') && metodo === 'GET') {
      const u = new URL(url, 'http://localhost');
      const pId = u.searchParams.get('proyecto_id') ?? 'mongui';
      const bp = buyerPersonaDemo[pId] ?? buyerPersonaDemo['mongui'];
      return json(bp);
    }

    // POST demo avanzar tiempo
    if (url.includes('/demo/tiempo') && metodo === 'POST') {
      return json({ fecha_simulada: new Date().toISOString(), hitos_disparados: 2 });
    }

    // POST demo reiniciar
    if (url.includes('/demo/reiniciar') && metodo === 'POST') {
      mensajesDemo.length = 1;
      return json({ reiniciado: true, fecha_simulada: new Date().toISOString() });
    }

    // POST crear lead
    if (url.endsWith('/api/leads') && metodo === 'POST') {
      let body: { precargado_id?: string } = {};
      try { body = JSON.parse((opciones?.body as string) || '{}'); } catch { body = {}; }

      const porPrecargado: Record<string, string> = {
        ana: 'mock-1', carlos: 'mock-2', luisa: 'mock-3',
      };
      const leadId = porPrecargado[body.precargado_id ?? 'ana'] ?? 'mock-1';

      return json({
        lead_id: leadId,
        estado: 'PERFILANDO',
        afiliado_detectado: leadId !== 'mock-2',
      }, 201);
    }

    // GET cola
    if (url.endsWith('/api/leads') && metodo === 'GET') {
      return json(colaDemo);
    }

    return original(entrada, opciones);
  };
  console.info('[mock] servidor simulado completo activo (Chat + Ficha + Cola + Gerencia + Demo)');
}
