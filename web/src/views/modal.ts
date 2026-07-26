/**
 * Componente de Modales Centrados y Elegantes para la aplicación Vivi.
 * Reemplaza las ventanas nativas del navegador (alert, prompt, confirm).
 */

export interface OpcionesModal {
  icono?: string;
  titulo: string;
  mensaje: string;
  textoConfirmar?: string;
  textoCancelar?: string;
  valorDefecto?: string;
  tipoBoton?: 'primario' | 'peligro' | 'advertencia';
}

/**
 * Muestra un modal de confirmación centrado (remplaza confirm).
 */
export function mostrarConfirmacion(opciones: OpcionesModal): Promise<boolean> {
  return new Promise<boolean>((resolve) => {
    const overlay = document.createElement('div');
    overlay.className = 'modal-overlay';
    overlay.setAttribute('role', 'dialog');
    overlay.setAttribute('aria-modal', 'true');

    const icono = opciones.icono ?? '❓';
    const txtOk = opciones.textoConfirmar ?? 'Aceptar';
    const txtCancelar = opciones.textoCancelar ?? 'Cancelar';
    const claseBtnOk = opciones.tipoBoton === 'peligro' ? 'btn-modal-peligro' : 'btn-modal-primario';

    overlay.innerHTML = `
      <div class="modal-card">
        <div class="modal-header">
          <span class="modal-icon">${icono}</span>
          <h3 class="modal-titulo">${escapar(opciones.titulo)}</h3>
        </div>
        <div class="modal-body">
          <p class="modal-mensaje">${escapar(opciones.mensaje)}</p>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn-modal-cancelar" data-modal-cancelar>${escapar(txtCancelar)}</button>
          <button type="button" class="${claseBtnOk}" data-modal-confirmar>${escapar(txtOk)}</button>
        </div>
      </div>
    `;

    document.body.appendChild(overlay);
    // Trigger animación de entrada
    requestAnimationFrame(() => overlay.classList.add('activo'));

    const cerrar = (resultado: boolean) => {
      overlay.classList.remove('activo');
      setTimeout(() => {
        if (overlay.parentNode) {
          overlay.parentNode.removeChild(overlay);
        }
        resolve(resultado);
      }, 200);
    };

    const btnConfirmar = overlay.querySelector<HTMLButtonElement>('[data-modal-confirmar]');
    const btnCancelar = overlay.querySelector<HTMLButtonElement>('[data-modal-cancelar]');

    btnConfirmar?.focus();

    btnConfirmar?.addEventListener('click', () => cerrar(true));
    btnCancelar?.addEventListener('click', () => cerrar(false));

    overlay.addEventListener('click', (e) => {
      if (e.target === overlay) cerrar(false);
    });

    const keyHandler = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        document.removeEventListener('keydown', keyHandler);
        cerrar(false);
      }
    };
    document.addEventListener('keydown', keyHandler);
  });
}

/**
 * Muestra un modal de entrada de texto centrado (remplaza prompt).
 */
export function solicitarTexto(opciones: OpcionesModal): Promise<string | null> {
  return new Promise<string | null>((resolve) => {
    const overlay = document.createElement('div');
    overlay.className = 'modal-overlay';
    overlay.setAttribute('role', 'dialog');
    overlay.setAttribute('aria-modal', 'true');

    const icono = opciones.icono ?? '✏️';
    const txtOk = opciones.textoConfirmar ?? 'Aceptar';
    const txtCancelar = opciones.textoCancelar ?? 'Cancelar';
    const valDefault = opciones.valorDefecto ?? '';

    overlay.innerHTML = `
      <div class="modal-card">
        <div class="modal-header">
          <span class="modal-icon">${icono}</span>
          <h3 class="modal-titulo">${escapar(opciones.titulo)}</h3>
        </div>
        <div class="modal-body">
          <p class="modal-mensaje">${escapar(opciones.mensaje)}</p>
          <input type="text" class="input-modal-texto" data-modal-input value="${escapar(valDefault)}" />
        </div>
        <div class="modal-footer">
          <button type="button" class="btn-modal-cancelar" data-modal-cancelar>${escapar(txtCancelar)}</button>
          <button type="button" class="btn-modal-primario" data-modal-confirmar>${escapar(txtOk)}</button>
        </div>
      </div>
    `;

    document.body.appendChild(overlay);
    requestAnimationFrame(() => overlay.classList.add('activo'));

    const input = overlay.querySelector<HTMLInputElement>('[data-modal-input]');

    const cerrar = (valor: string | null) => {
      overlay.classList.remove('activo');
      setTimeout(() => {
        if (overlay.parentNode) {
          overlay.parentNode.removeChild(overlay);
        }
        resolve(valor);
      }, 200);
    };

    if (input) {
      input.focus();
      input.select();
      input.addEventListener('keydown', (e: KeyboardEvent) => {
        if (e.key === 'Enter') {
          e.preventDefault();
          cerrar(input.value);
        }
      });
    }

    const btnConfirmar = overlay.querySelector<HTMLButtonElement>('[data-modal-confirmar]');
    const btnCancelar = overlay.querySelector<HTMLButtonElement>('[data-modal-cancelar]');

    btnConfirmar?.addEventListener('click', () => cerrar(input?.value ?? ''));
    btnCancelar?.addEventListener('click', () => cerrar(null));

    overlay.addEventListener('click', (e) => {
      if (e.target === overlay) cerrar(null);
    });

    const keyHandler = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        document.removeEventListener('keydown', keyHandler);
        cerrar(null);
      }
    };
    document.addEventListener('keydown', keyHandler);
  });
}

/**
 * Muestra un modal de notificación centrado (remplaza alert).
 */
export function mostrarNotificacion(opciones: OpcionesModal): Promise<void> {
  return new Promise<void>((resolve) => {
    const overlay = document.createElement('div');
    overlay.className = 'modal-overlay';
    overlay.setAttribute('role', 'dialog');
    overlay.setAttribute('aria-modal', 'true');

    const icono = opciones.icono ?? 'ℹ️';
    const txtOk = opciones.textoConfirmar ?? 'Entendido';

    overlay.innerHTML = `
      <div class="modal-card">
        <div class="modal-header">
          <span class="modal-icon">${icono}</span>
          <h3 class="modal-titulo">${escapar(opciones.titulo)}</h3>
        </div>
        <div class="modal-body">
          <p class="modal-mensaje">${escapar(opciones.mensaje)}</p>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn-modal-primario" data-modal-confirmar>${escapar(txtOk)}</button>
        </div>
      </div>
    `;

    document.body.appendChild(overlay);
    requestAnimationFrame(() => overlay.classList.add('activo'));

    const cerrar = () => {
      overlay.classList.remove('activo');
      setTimeout(() => {
        if (overlay.parentNode) {
          overlay.parentNode.removeChild(overlay);
        }
        resolve();
      }, 200);
    };

    const btnConfirmar = overlay.querySelector<HTMLButtonElement>('[data-modal-confirmar]');
    btnConfirmar?.focus();
    btnConfirmar?.addEventListener('click', () => cerrar());

    overlay.addEventListener('click', (e) => {
      if (e.target === overlay) cerrar();
    });

    const keyHandler = (e: KeyboardEvent) => {
      if (e.key === 'Escape' || e.key === 'Enter') {
        document.removeEventListener('keydown', keyHandler);
        cerrar();
      }
    };
    document.addEventListener('keydown', keyHandler);
  });
}

function escapar(s: string): string {
  const d = document.createElement('div');
  d.textContent = s;
  return d.innerHTML;
}
