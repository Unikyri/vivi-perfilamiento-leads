---
name: ficha-comercial
description: >
  Organiza la ficha comercial del asesor con secciones estables y cifras
  auditables. Usa esta skill cuando se genere o explique la ficha de un lead
  calificado para entrega al asesor.
agente: documentadora
fuente_de_verdad: Contrato §2.7
---

# Ficha comercial

## Por qué existe
El asesor necesita una lectura rápida que conserve el contexto y permita
actuar. Un orden fijo hace comparable cada ficha y evita que una cifra generada
por conversación parezca una decisión del motor.

## Instrucciones
- Presenta siempre las secciones en este orden fijo: Identificación, Capacidad, Perfil, Intención, Recomendaciones, Beneficios, Argumentos de venta y Alertas.
- Conserva los nombres de los campos del Contrato y no agregues datos fuera del contrato.
- Usa únicamente números calculados por el motor para capacidad, subsidio, cuota, presupuesto y recomendaciones.
- Señala cuando el perfil sea parcial o tenga baja confianza.
- Distingue datos declarados, datos verificados y resultados del motor.
- Incluye la alerta de desistimiento cuando exista y no la conviertas en un juicio sobre el lead.
- No edites ni redondees cifras normativas en la capa de presentación.

## Ejemplos
Entrada: ficha de un lead calificado con capacidad, perfil e intención disponibles.
Salida: una ficha con las secciones en el orden fijo, recomendaciones justificadas y números copiados del resultado del motor.

Entrada: perfil con confianza baja.
Salida: mantener la ficha y mostrar la advertencia de perfil parcial para que el asesor valide los campos marcados.

Caso límite:
Entrada: existe una recomendación, pero falta el número correspondiente del motor.
Salida: no fabricar el número ni completar la ficha con una estimación; marcar el dato pendiente para validación.

## Criterios de aceptación
- Las secciones aparecen siempre en el orden fijo del Contrato §2.7.
- Los números de la ficha provienen solo del motor.
- La confianza y las advertencias del perfil quedan visibles.
- Los datos faltantes se marcan para validación y nunca se inventan.
