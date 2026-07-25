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
      <button id="btn-avanzar-tiempo" class="btn-demo-action" type="button" title="Simular avance en la línea de tiempo">
        ⏩ Avanzar tiempo
      </button>

      <button id="btn-reiniciar-demo" class="btn-demo-action" type="button" title="Reiniciar demo al estado inicial (<3s)">
        ↺ Reiniciar demo
      </button>
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
