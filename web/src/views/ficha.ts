import type { Ficha, FuenteCampo, ItemDesglose, Nivel } from '../models/tipos';

const SELLO_FUENTE = {
  VERIFICADO_BASE: { txt: 'VERIFICADO', clase: 'sello-verificado' },
  DECLARADO: { txt: 'DECLARADO', clase: 'sello-declarado' },
  INFERIDO: { txt: 'INFERIDO', clase: 'sello-inferido' },
} as const;

const ZONA_NIVEL: Record<Nivel, string> = {
  ALTA: 'zona-verde',
  MEDIA: 'zona-ambar',
  BAJA: 'zona-gris',
};

function selloFuente(f: FuenteCampo): string {
  const b = SELLO_FUENTE[f] ?? SELLO_FUENTE.DECLARADO;
  return `<span class="sello-fuente ${b.clase}" title="Fuente: ${b.txt}">${b.txt}</span>`;
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
         <span class="banda-advertencia-etiqueta">Alerta</span>
         <span>${escapar(ficha.banda_advertencia)}</span>
       </div>`
    : '';

  const pctConfianza = Math.round(ficha.confianza_perfil * 100);
  const iden = ficha.identificacion;
  const cap = ficha.capacidad;
  const intenc = ficha.intencion;
  const zonaIntencion = ZONA_NIVEL[intenc.nivel] ?? ZONA_NIVEL.BAJA;

  contenedor.innerHTML = `
    <div class="hoja-informe informe-ficha">
      ${advertencia}

      <!-- Masthead: identidad a la izquierda, medidor de confianza a la derecha -->
      <header class="masthead-ficha">
        <div class="masthead-identidad">
          <span class="masthead-kicker">Ficha comercial</span>
          <h3 class="masthead-nombre">${escapar(iden.nombre || leadNombreFallback)}</h3>
          <dl class="registro-identidad">
            <div class="registro-identidad-item">
              <dt>Afiliación</dt>
              <dd>${iden.afiliada ? `Afiliado · Categoría ${escapar(iden.categoria)}` : 'No afiliado'}</dd>
            </div>
            <div class="registro-identidad-item">
              <dt>Teléfono</dt>
              <dd class="cifra">${escapar(iden.telefono || 'No registrado')}</dd>
            </div>
            <div class="registro-identidad-item">
              <dt>Cupo 10%</dt>
              <dd>${ficha.consume_cupo_10 ? 'Consume cupo' : 'No aplica'}</dd>
            </div>
          </dl>
        </div>

        <div class="medidor-confianza" role="img" aria-label="Confianza del perfil: ${pctConfianza} por ciento">
          <span class="medidor-confianza-etiqueta">Confianza del perfil</span>
          <div class="medidor-confianza-escala">
            <span class="medidor-confianza-tick" style="left: 25%"></span>
            <span class="medidor-confianza-tick" style="left: 50%"></span>
            <span class="medidor-confianza-tick" style="left: 75%"></span>
            <div class="medidor-confianza-relleno" style="transform: scaleX(${pctConfianza / 100})"></div>
            <div class="medidor-confianza-marcador" style="left: ${pctConfianza}%"></div>
          </div>
          <span class="medidor-confianza-cifra cifra">${pctConfianza}<span class="medidor-confianza-unidad">/100</span></span>
        </div>
      </header>

      <!-- Cartera de capacidad financiera -->
      <section class="seccion-informe seccion-capacidad">
        <div class="seccion-header-title">
          <span class="sec-icon">💳</span>
          <h4 class="seccion-titulo">Capacidad financiera</h4>
        </div>
        <div class="ledger-scroll">
          <table class="ledger">
            <thead>
              <tr>
                <th>Concepto</th>
                <th>Regla aplicada</th>
                <th class="ledger-num">Monto</th>
                <th>Fuente</th>
              </tr>
            </thead>
            <tbody>
              ${cap.desglose.map(renderFilaLedger).join('')}
            </tbody>
            <tfoot>
              <tr class="ledger-total">
                <td colspan="2">Presupuesto máximo</td>
                <td class="ledger-num cifra">$${formatoMonto(cap.presupuesto_max)}</td>
                <td></td>
              </tr>
              <tr class="ledger-ratio">
                <td colspan="2">Ratio de endeudamiento</td>
                <td class="ledger-num cifra">${(cap.ratio * 100).toFixed(0)}%</td>
                <td class="cifra">conf. ${(cap.confianza * 100).toFixed(0)}%</td>
              </tr>
            </tfoot>
          </table>
        </div>
      </section>

      <div class="informe-columnas">
        <!-- Zona de intención: banda con el veredicto de nivel -->
        <section class="seccion-informe zona-veredicto ${zonaIntencion}">
          <div class="seccion-header-title">
            <span class="sec-icon">🎯</span>
            <h4 class="seccion-titulo">Intención de compra</h4>
          </div>
          <p class="veredicto-nivel">Nivel <span class="cifra">${escapar(intenc.nivel)}</span></p>
          <p class="veredicto-confianza">Confianza: ${escapar(intenc.confianza)}</p>
          <ul class="lista-observaciones">
            ${intenc.senales.map(s => `<li>${escapar(s)}</li>`).join('')}
          </ul>
        </section>

        <!-- Argumentos y beneficios -->
        <section class="seccion-informe">
          <div class="seccion-header-title">
            <span class="sec-icon">💡</span>
            <h4 class="seccion-titulo">Argumentos de venta</h4>
          </div>
          <ul class="lista-observaciones">
            ${ficha.argumentos_venta.map(a => `<li>${escapar(a)}</li>`).join('')}
          </ul>
          <div class="seccion-header-title" style="margin-top:0.85rem">
            <span class="sec-icon">✨</span>
            <h5 class="subseccion-titulo" style="margin:0">Beneficios Colsubsidio</h5>
          </div>
          <ul class="lista-observaciones">
            ${ficha.beneficios.map(b => `<li>${escapar(b)}</li>`).join('')}
          </ul>
        </section>
      </div>

      <!-- Banner de instrucción única -->
      <div class="banner-siguiente-paso">
        <span><strong class="banner-etiqueta">Siguiente paso</strong> Agendar visita a sala de ventas</span>
        <button class="btn-copiar-resumen" id="btn-copiar-resumen" type="button">
          📋 Copiar resumen
        </button>
      </div>

      <!-- Nota al pie plegable -->
      <details class="timeline-plegable">
        <summary>Ver historial de perfilamiento con Vivi</summary>
        <p class="timeline-texto">
          Conversación iniciada. Canal WhatsApp. Transcripción procesada por Vivi.
        </p>
      </details>
    </div>
  `;

  // Copiar resumen al portapapeles
  const btnCopiar = contenedor.querySelector<HTMLButtonElement>('#btn-copiar-resumen');
  if (btnCopiar) {
    btnCopiar.addEventListener('click', () => {
      const texto = `RESUMEN FICHA - ${iden.nombre}\nPresupuesto: $${formatoMonto(cap.presupuesto_max)} | Afiliada: ${iden.afiliada ? 'Sí' : 'No'}\nIntención: ${intenc.nivel}\nSiguiente paso: Agendar visita sala de ventas.`;
      navigator.clipboard.writeText(texto).then(() => {
        btnCopiar.textContent = '✓ Copiado';
        setTimeout(() => { btnCopiar.textContent = '📋 Copiar resumen'; }, 2000);
      });
    });
  }
}

function renderFilaLedger(item: ItemDesglose): string {
  return `
    <tr>
      <td>${escapar(item.concepto)}</td>
      <td class="ledger-regla">${escapar(item.regla)}</td>
      <td class="ledger-num cifra">$${formatoMonto(item.monto)}</td>
      <td>${selloFuente(item.fuente)}</td>
    </tr>
  `;
}

function formatoMonto(monto: number): string {
  return `${(monto / 1_000_000).toFixed(1)}M`;
}

/** Estado vacío amable si la ficha aún no existe (404 FICHA_NO_DISPONIBLE) */
function renderEstadoVacio(contenedor: HTMLElement, nombreLead: string): void {
  contenedor.innerHTML = `
    <div class="hoja-informe informe-vacio">
      <span class="masthead-kicker">Ficha comercial</span>
      <h3>Aún sin generar</h3>
      <p>
        La ficha comercial completa de <strong>${escapar(nombreLead)}</strong> se generará automáticamente cuando Vivi complete la calificación.
      </p>
      <details class="timeline-plegable">
        <summary>Ver avance actual</summary>
        <p class="timeline-texto">
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
