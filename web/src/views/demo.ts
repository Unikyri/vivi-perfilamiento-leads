/**
 * Renderiza la botonera del Demo (US-12, US-17).
 * Estilo "Control Room" (grafito #111827), visualmente separada del producto.
 */
export function renderBotoneraDemo(
  contenedor: HTMLElement,
  onAvanzarTiempo: () => void,
  onReiniciarDemo: () => void,
): void {
  contenedor.innerHTML = `
    <div class="botonera-demo">
      <label class="campo-fecha">
        Avanzar a
        <input type="date" id="demo-fecha" value="2026-08-01" min="2026-07-26" />
      </label>

      <button id="btn-avanzar-tiempo" class="btn-demo-action" type="button" title="Simular avance en la línea de tiempo">
        ⏩ Avanzar tiempo
      </button>

      <button id="btn-reiniciar-demo" class="btn-demo-action" type="button" title="Reiniciar demo al estado inicial (<3s)">
        ↺ Reiniciar demo
      </button>

      <span id="demo-aviso" role="status" aria-live="polite" class="demo-aviso"></span>
    </div>
  `;

  // Listeners
  const btnAvanzar = contenedor.querySelector<HTMLButtonElement>('#btn-avanzar-tiempo');
  if (btnAvanzar) {
    btnAvanzar.addEventListener('click', onAvanzarTiempo);
  }

  const btnReiniciar = contenedor.querySelector<HTMLButtonElement>('#btn-reiniciar-demo');
  if (btnReiniciar) {
    btnReiniciar.addEventListener('click', onReiniciarDemo);
  }
}
