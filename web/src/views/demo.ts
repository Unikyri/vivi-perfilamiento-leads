/** Fecha fija de la demo: el reset del backend siempre vuelve al 2026-07-26
 * (`approvedDemoDate` en demo_repository.go), y el hito de afiliación más
 * rápido en dispararse cae a +8 días de creado el plan. Saltar a esta fecha
 * lo cubre con margen sin ir tan lejos como para pasarse la prima de
 * diciembre, si en algún momento se quiere mostrar esa también. */
const FECHA_AVANCE_DEMO = '2026-08-10';

/**
 * Renderiza los controles de demo en el topbar: un botón para avanzar
 * tiempo a una fecha fija (sin selector, sin ventanas nativas) y uno para
 * reiniciar. El aviso de resultado aparece inline, nunca en un alert().
 */
export function renderBotoneraDemo(
  contenedor: HTMLElement,
  onAvanzarTiempo: () => void,
  onReiniciarDemo: () => void,
): void {
  contenedor.innerHTML = `
    <button id="btn-avanzar-tiempo" class="top-action primary" type="button" title="Avanzar la demo a ${FECHA_AVANCE_DEMO}">
      <span class="play">▶</span> Avanzar tiempo
    </button>
    <button id="btn-reiniciar-demo" class="top-action" type="button" title="Reiniciar demo al estado inicial">
      <span>↺</span> Reiniciar
    </button>
    <span id="demo-aviso" role="status" aria-live="polite" class="demo-aviso"></span>
  `;

  contenedor.querySelector<HTMLButtonElement>('#btn-avanzar-tiempo')?.addEventListener('click', onAvanzarTiempo);
  contenedor.querySelector<HTMLButtonElement>('#btn-reiniciar-demo')?.addEventListener('click', onReiniciarDemo);
}

export function fechaAvanceDemo(): string {
  return FECHA_AVANCE_DEMO;
}
