# Decisión 0001: conservar la tabla `compradores`

Se conserva la tabla lógica `compradores` y su fuente JSON durante la separación entre tiempo operativo y tiempo simulado. La decisión mantiene estable el contrato existente y evita mezclar un cambio de reloj con una migración de datos.

## Alcance que se conserva

- `data/compradores.json` continúa siendo la fuente de carga del catálogo.
- El puerto de catálogo mantiene los compradores para el cálculo kNN del gemelo comprador.
- La ficha comercial y el endpoint de buyer-persona siguen derivando sus resultados de esos datos.
- El reinicio de demo restaura fechas y leads sintéticos, pero no elimina ni renombra compradores.

## Regla para una futura eliminación

Eliminar o cambiar `compradores` requiere un PR separado de modificación del Contrato v1.1 §9, con aprobación de ambos bloques (`feature/bloque-a` y `feature/bloque-b`). Ese PR debe incluir migración, compatibilidad y validación de los consumidores; no forma parte del cambio de reloj de issue #30.
