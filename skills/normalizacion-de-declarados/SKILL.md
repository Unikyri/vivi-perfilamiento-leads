---
name: normalizacion-de-declarados
description: >
  Convierte declaraciones coloquiales del hogar en datos estructurados y
  detecta ambigüedad. Usa esta skill cuando el cliente entregue ingresos,
  composición familiar o cualquier dato informal que deba registrarse.
agente: asesora
fuente_de_verdad: doc 05 §3
---

# Normalización de declarados

## Por qué existe
Las personas hablan de su realidad con palabras cotidianas, pero el motor
necesita datos comparables. Normalizar sin inventar protege la decisión y la
dignidad del cliente.

## Instrucciones
- Interpreta expresiones de dinero en pesos colombianos y conserva la precisión disponible.
- Convierte «dos palos y medio» en $2.500.000.
- Convierte «gano el mínimo» en $1.750.905 usando la referencia vigente autorizada.
- Convierte «como tres millones» en $3.000.000 con confianza media y marca el dato como aproximado.
- Interpreta «vivo con mi mamá y la niña» como tres personas y hogar monoparental ampliado.
- Separa siempre el valor normalizado, la confianza y la fuente declarada.
- Si la declaración es ambigua, como «gano bien», NO registres un valor: marca `requiere_confirmacion` y pregunta por un rango ofreciendo dos opciones concretas.
- No completes silenciosamente los datos que el cliente no declaró.

## Ejemplos
Entrada: «Entre mi pareja y yo ganamos dos palos y medio».
Salida: ingreso_hogar=$2.500.000, confianza=alta, fuente=declarado.

Entrada: «Como tres millones al mes».
Salida: ingreso_mensual=$3.000.000, confianza=media, fuente=declarado.

Caso límite:
Entrada: «Gano bien».
Salida: no registrar ingreso; `requiere_confirmacion=true`; preguntar: «¿Tu ingreso mensual está más cerca de $2 millones o de $4 millones?»

## Criterios de aceptación
- Cada valor normalizado conserva su confianza y su origen.
- Las expresiones ambiguas nunca se convierten en cifras inventadas.
- La confirmación solicita un rango con exactamente dos opciones concretas.
- La composición familiar refleja las personas mencionadas y su estructura.
