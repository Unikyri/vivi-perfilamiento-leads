---
name: siguiente-mejor-pregunta
description: >
  Elige la pregunta de mayor impacto para completar el perfil sin cansar al
  cliente. Usa esta skill cuando falte información para decidir capacidad,
  ruta o recomendación.
agente: asesora
fuente_de_verdad: doc 06 RF-M4-02
---

# Siguiente mejor pregunta

## Por qué existe
Cada pregunta tiene un costo de atención. Vivi debe pedir primero el dato que
más cambia la decisión y avanzar con respeto, no completar un formulario.

## Instrucciones
- Calcula qué campo faltante tiene mayor impacto sobre la decisión actual.
- Considera capacidad de compra, subsidio, intención y preferencias antes de elegir.
- NUNCA preguntes un campo con fuente `VERIFICADO_BASE`; ese valor es contexto autoritativo.
- Formula una sola pregunta clara por turno.
- Explica brevemente para qué sirve el dato antes de pedirlo si es sensible.
- Si hay empate, elige la pregunta más fácil de responder y más cercana al contexto del cliente.
- No repitas una pregunta ya respondida ni pidas datos que no cambian la decisión.

## Ejemplos
Entrada: falta ingreso_hogar y ya existe ahorro, ciudad y tipo de vivienda.
Salida: «Para estimar qué opciones te quedan cómodas, ¿cuánto reciben al mes en tu hogar?»

Entrada: ingreso_hogar está en `VERIFICADO_BASE`, pero falta la ciudad.
Salida: preguntar por la ciudad; nunca volver a preguntar el ingreso verificado.

Caso límite:
Entrada: faltan varios campos, pero todos los datos de capacidad están completos.
Salida: no pedir un campo financiero redundante; elegir la preferencia que más separa las recomendaciones y formular una sola pregunta.

## Criterios de aceptación
- La pregunta elegida corresponde al campo faltante de mayor impacto.
- Nunca se pregunta un campo `VERIFICADO_BASE`.
- Cada salida contiene una sola pregunta.
- La pregunta es accionable y no solicita información innecesaria.
