import type { Ficha, FuenteCampo } from '../models/tipos';

const BADGE_FUENTE = {
  VERIFICADO_BASE: { txt: 'VERIFICADO', icono: '✓', clase: 'badge-verificado' },
  DECLARADO:       { txt: 'DECLARADO',  icono: '✍', clase: 'badge-declarado'  },
  INFERIDO:        { txt: 'INFERIDO',   icono: '~', clase: 'badge-inferido'   },
} as const;

function badgeFuente(f: FuenteCampo): string {
  const b = BADGE_FUENTE[f] ?? BADGE_FUENTE.DECLARADO;
  return `<span class="badge-fuente ${b.clase}" title="Fuente: ${b.txt}">${b.icono} ${b.txt}</span>`;
}

/**
 * Renderiza la Ficha Comercial del Lead (US-15).
 * Cero fetch — render puro a partir del objeto Ficha.
 */
export function renderFicha(contenedor: HTMLElement, ficha: Ficha | null, leadNombreFallback = 'Lead'): void {
  if (!ficha) {
    renderEstadoVacio(contenedor, leadNombreFallback);
    return;
  }

  const advertencia = ficha.banda_advertencia
    ? `<div class="banda-advertencia" role="alert">
         <span>⚠️</span> <span>${escapar(ficha.banda_advertencia)}</span>
       </div>`
    : '';

  const pctConfianza = (ficha.confianza_perfil * 100).toFixed(0);
  const iden = ficha.identificacion;
  const cap = ficha.capacidad;
  const intenc = ficha.intencion;

  contenedor.innerHTML = `
    <div class="ficha-container">
      ${advertencia}

      <!-- Identificación Principal (F-Layout Top) -->
      <header class="ficha-header-card">
        <div class="ficha-id-info">
          <h3>${escapar(iden.nombre || leadNombreFallback)}</h3>
          <div class="ficha-meta-grid">
            <div class="ficha-meta-item">
              <strong>Afiliada:</strong> ${iden.afiliada ? `Sí (Cat. ${escapar(iden.categoria)})` : 'No'}
            </div>
            <div class="ficha-meta-item">
              <strong>Teléfono:</strong> ${escapar(iden.telefono || 'No registrado')}
            </div>
            <div class="ficha-meta-item">
              <strong>Cupo 10%:</strong> ${ficha.consume_cupo_10 ? 'Consume cupo' : 'No aplica'}
            </div>
          </div>
        </div>

        <div class="confianza-bar-container" title="Confianza general del perfilamiento">
          <span class="confianza-val">${pctConfianza}% Confianza</span>
          <progress class="confianza-progress" value="${ficha.confianza_perfil}" max="1"></progress>
        </div>
      </header>

      <!-- Bloque 3 Columnas: Intención, Capacidad, Recomendación -->
      <div class="grid-tres-columnas">

        <!-- Columna 1: Capacidad Financiera -->
        <article class="card-seccion">
          <h4>Capacidad Financiera</h4>
          <div class="desglose-monto-item">
            <span>Presupuesto Máx:</span>
            <span class="monto-num">$${(cap.presupuesto_max / 1_000_000).toFixed(1)}M ${badgeFuente('VERIFICADO_BASE')}</span>
          </div>
          <div class="desglose-monto-item">
            <span>Crédito Máx:</span>
            <span class="monto-num">$${(cap.credito_max / 1_000_000).toFixed(1)}M ${badgeFuente('INFERIDO')}</span>
          </div>
          <div class="desglose-monto-item">
            <span>Subsidio Aplicable:</span>
            <span class="monto-num">$${(cap.subsidio_aplicable / 1_000_000).toFixed(1)}M ${badgeFuente('VERIFICADO_BASE')}</span>
          </div>
          <div class="desglose-monto-item">
            <span>Recursos Propios:</span>
            <span class="monto-num">$${(cap.recursos_propios / 1_000_000).toFixed(1)}M ${badgeFuente('DECLARADO')}</span>
          </div>
        </article>

        <!-- Columna 2: Intención de Compra (Nivel + Confianza + Señales) -->
        <article class="card-seccion">
          <h4>Intención de Compra</h4>
          <p class="nivel-destacado nivel-${intenc.nivel}">
            Nivel ${escapar(intenc.nivel)} <small style="font-size:0.75rem; color:#6B7280">(Confianza ${escapar(intenc.confianza)})</small>
          </p>
          <ul class="lista-puntos">
            ${intenc.senales.map(s => `<li>${escapar(s)}</li>`).join('')}
          </ul>
        </article>

        <!-- Columna 3: Argumentos y Beneficios -->
        <article class="card-seccion">
          <h4>Argumentos de Venta</h4>
          <ul class="lista-puntos">
            ${ficha.argumentos_venta.map(a => `<li>${escapar(a)}</li>`).join('')}
          </ul>
          <h4 style="margin-top:0.5rem">Beneficios Colsubsidio</h4>
          <ul class="lista-puntos">
            ${ficha.beneficios.map(b => `<li>${escapar(b)}</li>`).join('')}
          </ul>
        </article>
      </div>

      <!-- Banner Azul Siguiente Paso (Doc 12 §3.4) -->
      <div class="banner-siguiente-paso">
        <span>▶ Siguiente paso: agendar visita sala de ventas</span>
        <button class="btn-copiar-resumen" id="btn-copiar-resumen" type="button">
          📋 Copiar Resumen
        </button>
      </div>

      <!-- Timeline Plegado de la Conversación -->
      <details class="timeline-plegable">
        <summary>💬 Ver historial de perfilamiento con Vivi</summary>
        <p style="font-size:0.83rem; margin-top:0.5rem; color:#4B5563">
          Conversación iniciada. Canal WhatsApp. Transcripción procesada por Vivi.
        </p>
      </details>
    </div>
  `;

  // Copiar resumen al portapapeles
  const btnCopiar = contenedor.querySelector<HTMLButtonElement>('#btn-copiar-resumen');
  if (btnCopiar) {
    btnCopiar.addEventListener('click', () => {
      const texto = `RESUMEN FICHA - ${iden.nombre}\nPresupuesto: $${(cap.presupuesto_max/1e6).toFixed(1)}M | Afiliada: ${iden.afiliada ? 'Sí' : 'No'}\nIntención: ${intenc.nivel}\nSiguiente paso: Agendar visita sala de ventas.`;
      navigator.clipboard.writeText(texto).then(() => {
        btnCopiar.textContent = '✓ Copiado!';
        setTimeout(() => { btnCopiar.textContent = '📋 Copiar Resumen'; }, 2000);
      });
    });
  }
}

/** Estado vacío amable si la ficha aún no existe (404 FICHA_NO_DISPONIBLE) */
function renderEstadoVacio(contenedor: HTMLElement, nombreLead: string): void {
  contenedor.innerHTML = `
    <div class="estado-vacio-amable">
      <h3>💬 Este lead aún está conversando con Vivi</h3>
      <p style="margin-bottom: 1rem; font-size: 0.9rem;">
        La ficha comercial completa de <strong>${escapar(nombreLead)}</strong> se generará automáticamente cuando Vivi complete la calificación.
      </p>
      <details class="timeline-plegable" style="max-width: 400px; margin: 0 auto; text-align: left;">
        <summary>Ver avance actual</summary>
        <p style="font-size:0.82rem; margin-top:0.5rem; color:#4B5563">
          Estado: En proceso de perfilamiento en tiempo real desde el chat.
        </p>
      </details>
    </div>
  `;
}

function escapar(s: string): string {
  const d = document.createElement('div');
  d.textContent = s;
  return d.innerHTML;
}
