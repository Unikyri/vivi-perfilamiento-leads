# DESIGN.md — Sistema de Diseño Vivi Colsubsidio

## 🎨 Paleta de Colores
- **Azul Institucional:** `#003DA6` (Header, botones primarios, banderas de acción)
- **Amarillo Acento:** `#FFC700` (Pill de prioridad, alertas, badges destacados)
- **Verde Éxito:** `#2E8540` / `#075E54` (Sustento de semáforo verde, tema WhatsApp)
- **Gris Neutral:** `#F4F6F9` (Fondos de tarjetas, paneles secundarios)
- **Texto:** `#1E293B` (Legibilidad estricta sobre fondo claro)

## 🔤 Tipografía
- **Fuente Principal:** Inter, System Sans-Serif
- **Jerarquía:**
  - Títulos principales: 1.25rem - 1.5rem bold
  - Subtítulos / Kicker: 0.85rem - 1rem semibold
  - Cuerpo: 0.9rem - 0.95rem regular
  - Badges / Captions: 0.75rem - 0.8rem bold (uppercase)

## 🧩 Componentes Principales
1. **Header Institucional:** Barra superior azul `#003DA6` con logo blanco recortado de Colsubsidio y tagline "Vivi - Asesora Digital".
2. **Panel de WhatsApp (Izquierda):** Ancho móvil estrecho (`380px`), burbujas beige/verde, micro-animación de grabación de micrófono, carrusel horizontal de máximo 3 tarjetas de proyectos.
3. **Panel del Asesor (Derecha):**
   - **Cola de Leads:** Tabla ordenada por prioridad backend, barra de cupo 10% con alerta al 80%, semáforos verde/ámbar/rojo, botón interactivo `💬 Ver chat`.
   - **Ficha Comercial:** Layout en F, banner de advertencia, desglose de capacidad con badges de fuente (`VERIFICADO_BASE`, `DECLARADO`, `INFERIDO`), banner azul de Siguiente Paso, botón de copiar resumen al portapapeles.
   - **Gerencia:** Gráficas de barras horizontales en tiempo real para Buyer Persona (afiliación, categoría A/B/C, rangos de edad) y selector de proyecto.
