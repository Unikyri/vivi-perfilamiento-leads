import type { Ficha, FuenteCampo, ItemDesglose, Ruta } from '../models/tipos';

const RUTA_INFO: Record<Ruta, { icono: string; titulo: string; detalle: string }> = {
  ASESOR: { icono: '👥', titulo: 'ASESOR', detalle: 'Lead listo para asesor comercial. Prioridad alta en cola.' },
  NUTRICION: { icono: '🌱', titulo: 'NUTRICIÓN', detalle: 'Aún no alcanza el perfil objetivo. En plan de acompañamiento.' },
  REMARKETING: { icono: '🔁', titulo: 'REMARKETING', detalle: 'Capacidad alta, intención baja. Reintentar contacto más adelante.' },
  DESPEDIDA: { icono: '👋', titulo: 'DESPEDIDA', detalle: 'No cumple condiciones mínimas para continuar el proceso.' },
};

const SELLO_FUENTE: Record<FuenteCampo, { txt: string; clase: string }> = {
  VERIFICADO_BASE: { txt: 'verificado', clase: 'verified' },
  DECLARADO: { txt: 'declarado', clase: 'declared-badge' },
  INFERIDO: { txt: 'inferido', clase: 'declared-badge' },
};

const NOMBRE_CAMPO: Record<string, string> = {
  ingreso_hogar: 'Ingreso hogar',
  recursos_propios: 'Recursos propios',
  tiene_vivienda: 'Tiene vivienda',
  recibio_subsidio: 'Recibió subsidio antes',
  edad: 'Edad',
  personas_hogar: 'Personas en el hogar',
  zona_deseada: 'Zona deseada',
};

let indiceAlternativa = 0;

/**
 * Renderiza la ficha comercial completa en una sola vista, sin pestañas
 * internas: Ruta asignada, Capacidad + Desglose, Subsidios, Perfil
 * declarado y Proyecto recomendado. `ruta` viene de la cola (LeadEnCola),
 * porque el objeto Ficha del Contrato no la incluye.
 */
export function renderFicha(contenedor: HTMLElement, ficha: Ficha | null, leadNombreFallback: string, ruta: Ruta | null): void {
  indiceAlternativa = 0;
  if (!ficha) {
    renderEstadoVacio(contenedor, leadNombreFallback);
    return;
  }
  pintar(contenedor, ficha, ruta);
}

function pintar(contenedor: HTMLElement, ficha: Ficha, ruta: Ruta | null): void {
  const iden = ficha.identificacion;
  const cap = ficha.capacidad;
  const infoRuta = ruta ? RUTA_INFO[ruta] : null;
  const nivelIntencion = ficha.intencion.nivel;
  const claseAtencion = nivelIntencion === 'MEDIA' ? 'media' : nivelIntencion === 'BAJA' ? 'baja' : '';
  const iniciales = (iden.nombre || 'L').trim().split(/\s+/).map(p => p[0]).slice(0, 2).join('').toUpperCase();
  const cuotaMensual = estimarCuotaMensual(ficha.perfil);

  contenedor.innerHTML = `
    <header class="profile-head">
      <span class="avatar profile-avatar" aria-hidden="true">${escapar(iniciales)}</span>
      <div>
        <div class="profile-name">${escapar(iden.nombre || 'Lead')}</div>
        <div class="profile-sub">${iden.afiliada ? `Afiliada · Categoría ${escapar(iden.categoria || 'N/A')}` : 'No afiliado'}</div>
        <div class="lead-id">ID Lead: ${escapar(ficha.lead_id)} <button class="copy" id="btn-copiar-id" type="button" title="Copiar ID">▣</button></div>
      </div>
      <span class="attention ${claseAtencion}">${nivelIntencion === 'ALTA' ? 'Alta' : nivelIntencion === 'MEDIA' ? 'Media' : 'Baja'} intención</span>
    </header>

    ${ficha.banda_advertencia ? `<div class="banda-advertencia" role="alert" style="margin-top:12px;padding:10px 12px;border-radius:8px;background:#fff1c9;color:#7a4b00;font-size:12px;font-weight:700">⚠️ ${escapar(ficha.banda_advertencia)}</div>` : ''}

    ${infoRuta ? `
      <div class="route-title"><span>Ruta asignada</span><span class="updated">Actualizado: ahora</span></div>
      <div class="route-card">
        <span class="route-icon" aria-hidden="true">${infoRuta.icono}</span>
        <div><strong>${infoRuta.titulo}</strong><p>${infoRuta.detalle}</p></div>
      </div>
    ` : ''}

    <div class="finance-grid" style="display:grid;grid-template-columns:1fr;gap:0">
      <section class="finance-card">
        <div class="section-title"><span>Capacidad financiera</span></div>
        <div class="stats">
          <div class="stat"><small>Capacidad estimada</small><strong>${formatoMonto(cap.presupuesto_max)}</strong></div>
          <div class="stat"><small>Cuota mensual (40% del ingreso)</small><strong>${cuotaMensual !== null ? formatoMonto(cuotaMensual) : '—'}</strong></div>
        </div>
        <div class="progress"><span style="width:40%"></span></div>
        <div class="rule">Regla cuota/ingreso: 40% · Ratio vs. proyecto: ${(cap.ratio * 100).toFixed(0)}% · Confianza: ${(cap.confianza * 100).toFixed(0)}%</div>
      </section>

      ${cap.desglose.length ? `
        <section class="breakdown">
          <h3>Desglose</h3>
          ${cap.desglose.map(renderLineaDesglose).join('')}
        </section>
      ` : ''}

      ${renderSubsidios(cap.desglose)}

      ${Object.keys(ficha.perfil).length ? `
        <section class="declared">
          <h3>Perfil declarado</h3>
          ${Object.entries(ficha.perfil).map(([campo, v]) => renderLineaPerfil(campo, v.valor, v.fuente)).join('')}
        </section>
      ` : ''}
    </div>

    ${renderProyecto(ficha)}
  `;

  contenedor.querySelector<HTMLButtonElement>('#btn-copiar-id')?.addEventListener('click', () => {
    navigator.clipboard.writeText(ficha.lead_id).catch(() => {});
  });

  const btnAlternativas = contenedor.querySelector<HTMLButtonElement>('#btn-ver-alternativas');
  if (btnAlternativas && ficha.recomendaciones.length > 1) {
    btnAlternativas.addEventListener('click', () => {
      indiceAlternativa = (indiceAlternativa + 1) % ficha.recomendaciones.length;
      pintar(contenedor, ficha, ruta);
    });
  }
}

function renderLineaDesglose(item: ItemDesglose): string {
  const sello = item.fuente ? SELLO_FUENTE[item.fuente] : null;
  return `
    <div class="detail-line">
      <span>${escapar(item.concepto)}<br><small>regla: ${escapar(item.regla)}</small></span>
      <b>${formatoMonto(item.monto)}${sello ? `<span class="${sello.clase}">${sello.txt}</span>` : ''}</b>
    </div>
  `;
}

function renderLineaPerfil(campo: string, valor: unknown, fuente: FuenteCampo): string {
  const sello = SELLO_FUENTE[fuente] ?? SELLO_FUENTE.DECLARADO;
  const nombre = NOMBRE_CAMPO[campo] ?? campo;
  let texto: string;
  if (typeof valor === 'boolean') texto = valor ? 'Sí' : 'No';
  else if (typeof valor === 'number' && campo.includes('ingreso') || campo.includes('recursos')) texto = formatoMonto(valor as number);
  else texto = String(valor);
  return `
    <div class="detail-line">
      <span>${escapar(nombre)}</span>
      <b>${escapar(texto)} <span class="${sello.clase}" style="display:inline">${sello.txt}</span></b>
    </div>
  `;
}

function renderSubsidios(desglose: ItemDesglose[]): string {
  const subsidios = desglose.filter(d => d.concepto.toLowerCase().includes('subsidio'));
  if (!subsidios.length) return '';
  const total = subsidios.reduce((acc, s) => acc + s.monto, 0);
  return `
    <section class="subsidies">
      <div class="section-title"><span>Subsidios aplicables</span></div>
      ${subsidios.map(s => `<div class="subsidy-row"><span>${escapar(s.concepto)}</span><strong>${formatoMonto(s.monto)}</strong></div>`).join('')}
      <div class="subsidy-total"><span>Total subsidios</span><span>${formatoMonto(total)}</span></div>
    </section>
  `;
}

function renderProyecto(ficha: Ficha): string {
  if (!ficha.recomendaciones.length) return '';
  const r = ficha.recomendaciones[indiceAlternativa] ?? ficha.recomendaciones[0];
  const hayMas = ficha.recomendaciones.length > 1;
  return `
    <section class="project-card">
      <div class="section-title">
        <span>Proyecto recomendado</span>
        ${hayMas ? `<button class="link" id="btn-ver-alternativas" type="button">Ver alternativas (${indiceAlternativa + 1}/${ficha.recomendaciones.length}) ›</button>` : ''}
      </div>
      <div class="project-info">
        <div class="project-image" aria-label="Imagen del proyecto ${escapar(r.nombre)}"></div>
        <div>
          <h3 class="project-name">${escapar(r.nombre)}</h3>
          <div class="project-place">${escapar(r.zona)}</div>
          <div class="project-meta">Desde ${formatoMonto(r.precio_desde)}</div>
          <div class="project-foot">${r.vecinos} vecinos compraron aquí · <strong>${(r.tasa_desistimiento * 100).toFixed(0)}% desistió</strong></div>
        </div>
      </div>
    </section>
  `;
}

/** No hay un campo `cuota_mensual` en el Contrato: se deriva con la MISMA
 * regla que ya usó el backend para el crédito (doc 13 §2.1, verificada en
 * issue #8): cuota = 0.40 × ingreso_hogar. Si el perfil no tiene el campo
 * (aún no declarado), no se muestra un número inventado. */
function estimarCuotaMensual(perfil: Ficha['perfil']): number | null {
  const ingreso = perfil['ingreso_hogar']?.valor;
  if (typeof ingreso !== 'number' || ingreso <= 0) return null;
  return Math.round(ingreso * 0.4);
}

function formatoMonto(monto: number): string {
  if (Math.abs(monto) >= 1_000_000) return `$${(monto / 1_000_000).toFixed(1)}M`;
  return `$${monto.toLocaleString('es-CO')}`;
}

function renderEstadoVacio(contenedor: HTMLElement, nombreLead: string): void {
  contenedor.innerHTML = `
    <div class="details-vacio">
      <h3>Ficha aún sin generar</h3>
      <p>La ficha comercial completa de <strong>${escapar(nombreLead)}</strong> se generará automáticamente cuando Vivi complete la calificación.</p>
    </div>
  `;
}

function escapar(s: string): string {
  const d = document.createElement('div');
  d.textContent = s;
  return d.innerHTML;
}
