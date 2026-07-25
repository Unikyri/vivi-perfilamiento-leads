import type { Conversacion, ColaLeads } from './tipos';

interface EstadoApp {
  leadActivo: string | null;
  conversacion: Conversacion | null;
  cola: ColaLeads | null;
  tabActiva: 'cola' | 'ficha' | 'gerencia';
}

type Escucha = (e: EstadoApp) => void;

const estado: EstadoApp = { leadActivo: null, conversacion: null, cola: null, tabActiva: 'cola' };
const escuchas: Escucha[] = [];

export function obtener(): Readonly<EstadoApp> { return estado; }

export function actualizar(cambios: Partial<EstadoApp>): void {
  Object.assign(estado, cambios);
  escuchas.forEach(f => f(estado));
}

export function suscribir(f: Escucha): void { escuchas.push(f); }
