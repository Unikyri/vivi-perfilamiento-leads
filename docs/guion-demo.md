# Vivi — guion de demo, promoción y operación

> **S4–S5 de Issue #27.** Guía para una demostración controlada con datos sintéticos y la promoción nativa ya configurada. No ejecuta despliegues, no provisiona credenciales y no sustituye la aprobación operativa del proveedor.

## Resultado esperado

El jurado debe poder abrir una sola URL, comprobar que el servicio está vivo y recorrer el flujo real: creación de un lead precargado, conversación, ficha del asesor y controles de demo. La aceptación de backend usa la API real, PostgreSQL y el proveedor LLM configurado; los datos de `ana`, `carlos` y `luisa` son sintéticos.

**Regla de aceptación:** `?mock=1` solo sirve para revisar la composición visual cuando no hay backend. No demuestra API, PostgreSQL, migraciones, semilla ni conversación LLM.

## Configuración operativa

Configurar las variables en el entorno de ejecución o en el gestor de configuración del proveedor. Nunca escribir valores de secretos en este archivo, commits, capturas, logs o comentarios.

| Variable | Desarrollo local | Demo controlada | Notas |
|---|---|---|---|
| `PORT` | `8080` | La asignada por la plataforma | El proceso debe escuchar el puerto del entorno. |
| `DATABASE_URL` | URL de PostgreSQL local | La URL entregada por PostgreSQL administrado | Obligatoria para iniciar; usar solo un valor real fuera del repositorio. |
| `DEMO_SEED` | `false` | `true` durante la preparación y la evaluación | Por defecto es `false`. Solo habilitarla en la aplicación controlada; permite semilla de arranque y reset. |
| `LLM_PROVIDER` | `gemini` | `gemini` o el proveedor aprobado | La conversación real requiere un proveedor configurado. |
| `GEMINI_API_KEY` | Vacía si solo se prueba salud/SPA | Referencia a secreto del entorno | Nunca pegar el valor en la terminal compartida ni en documentación. |
| `QWEN_API_KEY` | Vacía salvo uso de Qwen | Referencia a secreto si aplica fallback | No es necesaria para el flujo Gemini primario. |
| `QWEN_BASE_URL` | Vacía salvo uso de Qwen | Referencia de configuración aprobada | No incluir tokens en la URL. |
| `LLM_FALLBACK` | `qwen` | Según la configuración aprobada | Verificar que el fallback tenga su secreto antes de anunciarlo. |
| `TASA_EA` | `0.107` | Valor aprobado por la configuración del release | No cambiar criterios deterministas como parte de la demo. |
| `LOG_NIVEL` | `info` | `info` | Los logs no deben contener secretos, PII ni respuestas de proveedor. |

## Promoción nativa desde GitHub y verificación pública

La aplicación confirmada para la demo es `vivi-37863aed9d29`. La integración nativa de Heroku con GitHub está configurada para desplegar `main` después de que CI termine correctamente; el repositorio no necesita ni debe añadir un workflow GitHub adicional para desplegar. La aplicación usa un dyno web Basic y PostgreSQL administrado. Los valores de configuración permanecen en Heroku y no se copian al repositorio.

Flujo de promoción autorizado:

1. Integrar el cambio en `main` mediante el flujo normal de PR y esperar el CI requerido.
2. Dejar que la integración GitHub de Heroku promueva automáticamente el SHA de `main`; no ejecutar comandos de despliegue desde este repositorio ni imprimir variables de entorno.
3. Verificar el release público antes de la demo, empezando por salud y luego por la entrada SPA:

```bash
export BASE_URL="${BASE_URL:-https://vivi-37863aed9d29.herokuapp.com}"
curl --fail --silent --show-error "$BASE_URL/salud" | jq '{estado,bd,fecha_simulada}'
curl --fail --silent --show-error "$BASE_URL/" | grep -Eiq '<title>|<html'
```

La respuesta de `/salud` debe ser JSON válido y mostrar la disponibilidad de `bd`; después se ejecuta el Script 1 contra `$BASE_URL` sin `?mock=1`. Si salud falla, se detiene la demostración y se usa el procedimiento normal de rollback del release desde Heroku/GitHub, seguido de una nueva comprobación pública de `/salud`. Esta guía no ejecuta ese rollback ni ninguna operación contra el proveedor.

La aplicación es de un solo dyno/proceso: el tracker de turnos es local a la instancia. No escalar horizontalmente como parte de esta guía.

## Script 1 — recorrido real ante el jurado

**Precondiciones:** PostgreSQL está disponible, las migraciones de arranque terminan, `DEMO_SEED=true` solo en la aplicación controlada y el proveedor LLM fue configurado mediante secretos del entorno. No se requiere ADK.

### 1. Comprobar salud y entrada

```bash
export BASE_URL="${BASE_URL:-http://localhost:8080}"
curl --fail --silent --show-error "$BASE_URL/salud" | jq .
curl --fail --silent --show-error "$BASE_URL/" | grep -Eiq '<title>|<html'
```

En `/salud`, confirmar que la respuesta JSON se puede leer y que `bd` está disponible. La fecha `fecha_simulada` es la del reloj de demo, no una prueba de que el proveedor LLM esté respondiendo.

### 2. Abrir el producto y elegir un lead

Abrir `$BASE_URL/?precargado=ana` en el navegador. La SPA crea el lead `ana` contra la API real. También se puede preparar el mismo paso desde terminal:

```bash
LEAD_JSON=$(curl --fail --silent --show-error \
  -X POST "$BASE_URL/api/leads" \
  -H 'Content-Type: application/json' \
  --data '{"precargado_id":"ana","fuente":"DEMO"}')
printf '%s\n' "$LEAD_JSON" | jq .
LEAD_ID=$(printf '%s' "$LEAD_JSON" | jq -r '.lead_id')
test -n "$LEAD_ID" && test "$LEAD_ID" != "null"
```

Se aceptan los precargados `ana`, `carlos` y `luisa`. El backend debe devolver un `lead_id` real; no usar un ID inventado para continuar.

### 3. Mostrar conversación y ficha

Enviar un mensaje de texto y consultar la conversación hasta que termine el turno asíncrono:

```bash
curl --fail --silent --show-error \
  -X POST "$BASE_URL/api/leads/$LEAD_ID/mensajes" \
  -H 'Content-Type: application/json' \
  --data '{"tipo":"TEXTO","texto":"Quiero conocer opciones de vivienda para mi familia."}' | jq .

curl --fail --silent --show-error "$BASE_URL/api/leads/$LEAD_ID/conversacion" | jq .
curl --fail --silent --show-error "$BASE_URL/api/leads/$LEAD_ID/ficha" | jq .
curl --fail --silent --show-error "$BASE_URL/api/leads" | jq .
```

En la pantalla, mostrar el chat y luego las pestañas **Cola**, **Ficha** y **Gerencia**. La respuesta del proveedor puede ser asíncrona: esperar el polling de conversación antes de concluir que el turno falló. La ficha y las recomendaciones deben corresponder al `lead_id` creado, no a la maqueta.

### 4. Mostrar el control temporal y dejar el estado limpio

Desde la botonera **Avanzar tiempo**, o directamente:

```bash
curl --fail --silent --show-error \
  -X POST "$BASE_URL/api/demo/tiempo" \
  -H 'Content-Type: application/json' \
  --data '{"avanzar_hasta":"2026-08-01"}' | jq .
```

Si la evaluación necesita repetir el recorrido, ejecutar el reset solo con la semilla habilitada:

```bash
time curl --fail --silent --show-error \
  -X POST "$BASE_URL/api/demo/reiniciar" \
  -H 'Content-Type: application/json' | tee /tmp/vivi-reset.json
jq -e '.reiniciado == true and (.fecha_simulada | type == "string")' /tmp/vivi-reset.json
curl --fail --silent --show-error "$BASE_URL/salud" | jq .fecha_simulada
```

El reset debe terminar en menos de tres segundos, restaurar la fecha aprobada y dejar una sola semilla utilizable para cada identidad. Solo elimina `fichas`, `hitos`, `planes`, `mensajes` y `leads`; preserva `compradores` y no ejecuta DDL. Al terminar la evaluación, volver a `DEMO_SEED=false` siguiendo el procedimiento aprobado del entorno.

## Script 2 — recuperación y repetición controlada

Usar este flujo cuando la preparación quedó sucia, la conversación no se puede repetir o se necesita demostrar idempotencia. No usarlo sobre una base con datos reales.

1. Confirmar que la aplicación es la instancia controlada y que `DEMO_SEED=true` fue habilitada deliberadamente.
2. Ejecutar `POST /api/demo/reiniciar` una vez y verificar HTTP 200 con `reiniciado: true`.
3. Ejecutarlo una segunda vez y comparar que devuelve la misma fecha simulada y las tres identidades sin duplicados:

```bash
for intento in 1 2; do
  curl --fail --silent --show-error \
    -X POST "$BASE_URL/api/demo/reiniciar" \
    -H 'Content-Type: application/json' | jq .
done

for persona in ana carlos luisa; do
  curl --fail --silent --show-error \
    -X POST "$BASE_URL/api/leads" \
    -H 'Content-Type: application/json' \
    --data "{\"precargado_id\":\"$persona\",\"fuente\":\"DEMO\"}" | jq --arg p "$persona" '{persona:$p,lead_id,estado}'
done
curl --fail --silent --show-error "$BASE_URL/api/leads" | jq .
curl --fail --silent --show-error "$BASE_URL/salud" | jq '{estado,bd,fecha_simulada}'
```

4. Recargar la SPA sin `mock=1`, abrir el chat y volver a recorrer el Script 1. Si el reset está deshabilitado, esperar HTTP 500 `ERROR_INTERNO` genérico y confirmar que no hubo mutación; la respuesta no debe revelar `DEMO_SEED` ni configuración interna.
5. Deshabilitar la semilla después de la prueba y conservar evidencia únicamente de estados agregados y sintéticos.

## Fallback visual (no aceptación)

Para una revisión de layout sin PostgreSQL ni proveedor, abrir:

```text
$BASE_URL/?mock=1
```

Este modo intercepta las llamadas del navegador y muestra fixtures locales. Debe rotularse como **visual-only / mock** y no cuenta como recorrido real, prueba de `/salud`, prueba de reset, evidencia de migraciones ni aceptación de conversación.

## Disponibilidad y anti-sleep

| Opción | Uso | Decisión/limitación |
|---|---|---|
| Dyno no suspensible | Demo durante la ventana de evaluación | Preferida si el presupuesto y la política del entorno lo permiten; es una decisión del operador. |
| Ping externo a `/salud` | Mantener despierta una instancia con suspensión | Opcional. Configurar un monitor autorizado con intervalo moderado y sin enviar secretos. Verificar primero la política del proveedor. |
| Ping desde el navegador del jurado | No recomendado | No es evidencia operativa y deja la disponibilidad atada a una pestaña abierta. |

El ping debe ser solo `GET /salud`, sin `POST`, sin reset y sin datos de lead. No se configura ni ejecuta ningún monitor desde este cambio.

## Rollback y límites

- **S4:** revertir o eliminar `docs/guion-demo.md`; no cambia código, esquema, datos ni configuración.
- **Semilla:** si el comportamiento de demo es inesperado, detener el uso de reset y devolver `DEMO_SEED=false` en el entorno controlado; no aplicar el reset a datos reales.
- **SPA/build:** si faltan assets hasheados, detener la promoción y volver al release anterior mediante el procedimiento normal del operador; no ocultar el error con `?mock=1`.
- **Promoción:** la integración nativa de Heroku con GitHub es el único mecanismo documentado; si el release falla, detener la demo, volver al release anterior mediante el procedimiento autorizado del operador y repetir la verificación pública de `/salud`.
- **Servicio:** conservar el único dyno y verificar `/salud` después de cualquier rollback autorizado.

Este cambio crea únicamente `app.json` como metadata secret-free. No crea `deploy.yml`, no añade una GitHub Action de despliegue y no ejecuta ninguna operación contra Heroku, Gemini, Qwen o un proveedor externo.

## Verificación de documentación

La copia local de Wiki/Doc 12 no está presente en el checkout ni en el historial Git accesible de esta rama. Por ello, este documento se contrastó contra el contrato del issue #27, la especificación `despliegue-heroku` y las rutas implementadas (`/salud`, `/api/leads`, conversación, ficha, tiempo y reset). Antes de fusionar, un mantenedor debe revisar los nombres exactos de los dos scripts de Doc 12 §4 contra la fuente oficial disponible.

Revisión de secretos: el documento contiene nombres de variables y placeholders, pero no valores de API keys, tokens, contraseñas, URLs con credenciales ni datos personales reales.
