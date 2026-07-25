// Traducción literal del Contrato v1.1 §1 y §2.
// snake_case en español porque así viaja el JSON.

export type EstadoLead = 'NUEVO' | 'PERFILANDO' | 'CALIFICADO' | 'ENTREGADO'
  | 'EN_NUTRICION' | 'PAUSADO' | 'REMARKETING' | 'DESPEDIDO' | 'CERRADO';
export type Ruta = 'ASESOR' | 'NUTRICION' | 'REMARKETING' | 'DESPEDIDA';
export type FuenteCampo = 'VERIFICADO_BASE' | 'DECLARADO' | 'INFERIDO';
export type Nivel = 'ALTA' | 'MEDIA' | 'BAJA';
export type Semaforo = 'VERDE' | 'AMBAR' | 'GRIS';
export type TipoContenido = 'TEXTO' | 'TARJETAS_PROYECTOS' | 'OFERTA_PLAN'
  | 'HITO_NUTRICION' | 'SISTEMA';
export type AutorMensaje = 'LEAD' | 'VIVI' | 'SISTEMA';
export type EstadoPlan = 'ACTIVO' | 'PAUSADO' | 'COMPLETADO' | 'CANCELADO';

export interface CampoPerfil {
  valor: string | number | boolean;
  fuente: FuenteCampo;
  confianza: number;
  requiere_confirmacion: boolean;
  actualizado_en: string;
}

export interface ItemDesglose {
  concepto: string; monto: number; regla: string; fuente: FuenteCampo;
}

export interface Capacidad {
  presupuesto_max: number; credito_max: number; subsidio_aplicable: number;
  recursos_propios: number; ratio: number; confianza: number;
  desglose: ItemDesglose[];
}

export interface Intencion { nivel: Nivel; confianza: Nivel; senales: string[]; }

export interface Recomendacion {
  proyecto_id: string; nombre: string; zona: string; precio_desde: number;
  razon: string; vecinos: number; tasa_desistimiento: number;
  brochure_url: string; recorrido_360_url: string;
}

export interface Hito {
  hito_id: string; tipo: string; fecha: string; monto: number | null;
  descripcion: string; estado: string;
}

export interface PlanNutricion {
  plan_id: string; lead_id: string; estado: EstadoPlan;
  consentimiento_en: string | null; frecuencia: string;
  meta_monto: number; meta_descripcion: string; hitos: Hito[];
}

export interface Mensaje {
  mensaje_id: string; autor: AutorMensaje; tipo_contenido: TipoContenido;
  texto: string; creado_en: string;
  adjunto: {
    recomendaciones?: Recomendacion[];
    plan_propuesto?: PlanNutricion;
    hito?: Hito;
    audio_original?: boolean;
  } | null;
}

export interface Conversacion {
  lead_id: string; estado: EstadoLead; turno_en_proceso: boolean; mensajes: Mensaje[];
}

export interface Ficha {
  ficha_id: string; lead_id: string; generada_en: string;
  confianza_perfil: number; banda_advertencia: string | null;
  identificacion: { nombre: string; afiliada: boolean; categoria: string; telefono: string };
  capacidad: Capacidad;
  perfil: Record<string, CampoPerfil>;
  intencion: Intencion;
  recomendaciones: Recomendacion[];
  beneficios: string[];
  argumentos_venta: string[];
  alerta_desistimiento: { activa: boolean; tasa_vecinos: number; detalle: string | null };
  consume_cupo_10: boolean;
}

export interface LeadEnCola {
  lead_id: string; nombre: string; estado: EstadoLead; ruta: Ruta;
  afiliado: boolean; semaforo: Semaforo; prioridad: number;
  resumen: string; actualizado_en: string;
}

export interface ColaLeads {
  cupo_10: { usados: number; porcentaje_ventana: number };
  leads: LeadEnCola[];
}

export interface BuyerPersona {
  proyecto_id: string; nombre: string; muestras: number;
  afiliacion: { afiliados: number; no_afiliados: number };
  categoria: Record<string, number>;
  rango_edad: Record<string, number>;
  tasa_desistimiento: number; actualizado_en: string;
}

export interface RespuestaError {
  error: { codigo: string; mensaje: string; detalles: Record<string, unknown> };
}
