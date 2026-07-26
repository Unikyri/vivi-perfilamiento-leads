(function(){const e=document.createElement("link").relList;if(e&&e.supports&&e.supports("modulepreload"))return;for(const n of document.querySelectorAll('link[rel="modulepreload"]'))o(n);new MutationObserver(n=>{for(const i of n)if(i.type==="childList")for(const s of i.addedNodes)s.tagName==="LINK"&&s.rel==="modulepreload"&&o(s)}).observe(document,{childList:!0,subtree:!0});function t(n){const i={};return n.integrity&&(i.integrity=n.integrity),n.referrerPolicy&&(i.referrerPolicy=n.referrerPolicy),n.crossOrigin==="use-credentials"?i.credentials="include":n.crossOrigin==="anonymous"?i.credentials="omit":i.credentials="same-origin",i}function o(n){if(n.ep)return;n.ep=!0;const i=t(n);fetch(n.href,i)}})();let A=1,N=!1;const M=[{mensaje_id:"m1",autor:"VIVI",tipo_contenido:"TEXTO",texto:"¡Hola Ana! 👋 Como afiliada tienes un subsidio de hasta $52,5M. ¿Sueñas con comprar este año?",creado_en:new Date().toISOString(),adjunto:null}];function K(a){return{lead_id:a,estado:"PERFILANDO",turno_en_proceso:N,mensajes:[...M]}}const Q={cupo_10:{usados:1,porcentaje_ventana:10},leads:[{lead_id:"mock-1",nombre:"Ana Rodríguez",estado:"ENTREGADO",ruta:"ASESOR",afiliado:!0,semaforo:"VERDE",prioridad:.91,resumen:"Afiliada cat. A · presupuesto $166.8M · intención alta",actualizado_en:new Date().toISOString()},{lead_id:"mock-2",nombre:"Carlos Martínez",estado:"PERFILANDO",ruta:"ASESOR",afiliado:!1,semaforo:"AMBAR",prioridad:.78,resumen:"No afiliado · presupuesto $210M · intención media",actualizado_en:new Date().toISOString()},{lead_id:"mock-3",nombre:"Luisa Gómez",estado:"CALIFICADO",ruta:"NUTRICION",afiliado:!0,semaforo:"VERDE",prioridad:.85,resumen:"Afiliada cat. B · presupuesto $145M · requiere plan cuota inicial",actualizado_en:new Date().toISOString()}]},Z={"mock-1":{ficha_id:"f1",lead_id:"mock-1",generada_en:new Date().toISOString(),confianza_perfil:.94,banda_advertencia:null,identificacion:{nombre:"Ana Rodríguez",afiliada:!0,categoria:"A",telefono:"+57 300 123 4567"},capacidad:{presupuesto_max:1668e5,credito_max:1143e5,subsidio_aplicable:525e5,recursos_propios:12e6,ratio:.28,confianza:.94,desglose:[{concepto:"Subsidio Mi Casa Ya / Caja",monto:525e5,regla:"Afiliado Cat A",fuente:"VERIFICADO_BASE"},{concepto:"Preaprobado Bancolombia",monto:1143e5,regla:"Capacidad de endeudamiento 30%",fuente:"INFERIDO"},{concepto:"Ahorro Declarado",monto:12e6,regla:"Declarado en chat",fuente:"DECLARADO"}]},perfil:{ingreso:{valor:26e5,fuente:"VERIFICADO_BASE",confianza:.95,requiere_confirmacion:!1,actualizado_en:new Date().toISOString()}},intencion:{nivel:"ALTA",confianza:"ALTA",senales:["Busca comprar antes de 6 meses","Tiene ahorro inicial preparado","Responde inmediatamente a Vivi"]},recomendaciones:[{proyecto_id:"mongui",nombre:"Monguí",zona:"Ciudadela Maiporé - Soacha",precio_desde:15647e4,razon:"Tu presupuesto cubre el 100% de la cuota inicial",vecinos:622,tasa_desistimiento:.12,brochure_url:"https://heyzine.com/flip-book/866af8f6a6.html",recorrido_360_url:"https://storage.net-fs.com/hosting/7532170/19/"}],beneficios:["Subsidio de vivienda Colsubsidio hasta $52,5M","Tasa preferencial crédito hipotecario"],argumentos_venta:["Cuota estimada mensual ($650k) es MENOR al arriendo promedio de la zona ($850k)","Entrega inmediata en 2026"],alerta_desistimiento:{activa:!1,tasa_vecinos:.12,detalle:null},consume_cupo_10:!1},"mock-2":{ficha_id:"f2",lead_id:"mock-2",generada_en:new Date().toISOString(),confianza_perfil:.82,banda_advertencia:"No afiliado a Colsubsidio — consume cupo del 10% regulatorio",identificacion:{nombre:"Carlos Martínez",afiliada:!1,categoria:"N/A",telefono:"+57 311 987 6543"},capacidad:{presupuesto_max:21e7,credito_max:18e7,subsidio_aplicable:0,recursos_propios:3e7,ratio:.32,confianza:.82,desglose:[{concepto:"Crédito solicitado",monto:18e7,regla:"Ingresos independientes",fuente:"DECLARADO"},{concepto:"Cuota Inicial Propia",monto:3e7,regla:"Recursos declarados",fuente:"DECLARADO"}]},perfil:{},intencion:{nivel:"MEDIA",confianza:"MEDIA",senales:["Interesado en proyectos VIS y no VIS","Evaluando opciones de crédito"]},recomendaciones:[],beneficios:["Opción de crédito con aliados de la caja"],argumentos_venta:["Proyecto Versalles cuenta con certificación EDGE para ahorro energético"],alerta_desistimiento:{activa:!1,tasa_vecinos:.08,detalle:null},consume_cupo_10:!0}},w={mongui:{proyecto_id:"mongui",nombre:"Monguí",muestras:312,afiliacion:{afiliados:198,no_afiliados:114},categoria:{A:110,B:68,C:20,SIN_DATO:114},rango_edad:{"20-35":165,"36-45":82,"46-55":40,"55+":12,SIN_DATO:13},tasa_desistimiento:.11,actualizado_en:new Date().toISOString()},macarena:{proyecto_id:"macarena",nombre:"La Macarena",muestras:185,afiliacion:{afiliados:140,no_afiliados:45},categoria:{A:85,B:42,C:13,SIN_DATO:45},rango_edad:{"20-35":100,"36-45":50,"46-55":20,"55+":8,SIN_DATO:7},tasa_desistimiento:.08,actualizado_en:new Date().toISOString()},versalles:{proyecto_id:"versalles",nombre:"Versalles",muestras:142,afiliacion:{afiliados:85,no_afiliados:57},categoria:{A:40,B:32,C:13,SIN_DATO:57},rango_edad:{"20-35":70,"36-45":42,"46-55":15,"55+":8,SIN_DATO:7},tasa_desistimiento:.15,actualizado_en:new Date().toISOString()}};function ee(a){N=!0,setTimeout(()=>{A++;const e=/proyecto|comprar|vivienda|casa|apto/i.test(a),t={mensaje_id:`m${A}`,autor:"VIVI",tipo_contenido:e?"TARJETAS_PROYECTOS":"TEXTO",texto:e?"Basándome en tu perfil, estos proyectos te pueden interesar:":"¡Qué bueno saberlo! ¿Te interesaría que exploremos opciones de vivienda según tu presupuesto?",creado_en:new Date().toISOString(),adjunto:e?{recomendaciones:[{proyecto_id:"mongui",nombre:"Monguí",zona:"Ciudadela Maiporé - Soacha",precio_desde:15647e4,razon:"Tu presupuesto cubre el 100% de la cuota inicial",vecinos:622,tasa_desistimiento:.12,brochure_url:"https://heyzine.com/flip-book/866af8f6a6.html",recorrido_360_url:"https://storage.net-fs.com/hosting/7532170/19/"},{proyecto_id:"macarena",nombre:"La Macarena",zona:"Ciudadela Maiporé - Soacha",precio_desde:12834e4,razon:"El más económico de la zona, ideal para tu ingreso",vecinos:374,tasa_desistimiento:.08,brochure_url:"https://heyzine.com/flip-book/b168b2f5ba.html",recorrido_360_url:""},{proyecto_id:"versalles",nombre:"Versalles",zona:"Ciudadela Maiporé - Soacha",precio_desde:1952e5,razon:"Certificación EDGE, ahorro en servicios",vecinos:174,tasa_desistimiento:.15,brochure_url:"https://heyzine.com/flip-book/be784b0d5c.html",recorrido_360_url:"https://shape.com.co/360/COLSUBSIDIO-Versalles_APTOA"}]}:null};M.push(t),N=!1},2e3)}function ae(){const a=window.fetch;window.fetch=async(e,t)=>{const o=e.toString(),n=((t==null?void 0:t.method)??"GET").toUpperCase(),i=(s,r=200)=>new Response(JSON.stringify(s),{status:r,headers:{"Content-Type":"application/json"}});if(o.includes("/conversacion")&&n==="GET"){const s=o.match(/\/leads\/([^/]+)\/conversacion/);return i(K(s?s[1]:"mock-1"))}if(o.includes("/mensajes")&&n==="POST"){A++;let s={};try{s=JSON.parse((t==null?void 0:t.body)||"{}")}catch{s={}}const r=s.tipo==="AUDIO";return M.push({mensaje_id:`m${A}`,autor:"LEAD",tipo_contenido:"TEXTO",texto:r?"🎙️ [Nota de voz]":s.texto||"",creado_en:new Date().toISOString(),adjunto:r?{audio_original:!0}:null}),ee(r?"audio":s.texto||""),i({mensaje_id:`m${A}`,turno_en_proceso:!0},201)}if(o.includes("/ficha")&&n==="GET"){const s=o.match(/\/leads\/([^/]+)\/ficha/),r=s?s[1]:"mock-1",c=Z[r];return c?i(c):i({error:{codigo:"FICHA_NO_DISPONIBLE",mensaje:"Ficha aún no disponible"}},404)}if(o.includes("/gerencia/buyer-persona")&&n==="GET"){const r=new URL(o,"http://localhost").searchParams.get("proyecto_id")??"mongui",c=w[r]??w.mongui;return i(c)}if(o.includes("/demo/tiempo")&&n==="POST")return i({fecha_simulada:new Date().toISOString(),hitos_disparados:2});if(o.includes("/demo/reiniciar")&&n==="POST")return M.length=1,i({reiniciado:!0,fecha_simulada:new Date().toISOString()});if(o.endsWith("/api/leads")&&n==="POST"){let s={};try{s=JSON.parse((t==null?void 0:t.body)||"{}")}catch{s={}}const c={ana:"mock-1",carlos:"mock-2",luisa:"mock-3"}[s.precargado_id??"ana"]??"mock-1";return i({lead_id:c,estado:"PERFILANDO",afiliado_detectado:c!=="mock-2"},201)}return o.endsWith("/api/leads")&&n==="GET"?i(Q):a(e,t)},console.info("[mock] servidor simulado completo activo (Chat + Ficha + Cola + Gerencia + Demo)")}const P={leadActivo:null,conversacion:null,cola:null,tabActiva:"cola"},U=[];function I(){return P}function b(a){Object.assign(P,a),U.forEach(e=>e(P))}function te(a){U.push(a)}const oe="/api";class h extends Error{constructor(e,t,o){super(t),this.codigo=e,this.estadoHTTP=o,this.name="ErrorAPI"}}async function f(a,e){var o,n;const t=await fetch(`${oe}${a}`,{headers:{"Content-Type":"application/json"},...e});if(!t.ok){const i=await t.json().catch(()=>null);throw new h(((o=i==null?void 0:i.error)==null?void 0:o.codigo)??"ERROR_INTERNO",((n=i==null?void 0:i.error)==null?void 0:n.mensaje)??`HTTP ${t.status}`,t.status)}return await t.json()}const p={crearLead:a=>f("/leads",{method:"POST",body:JSON.stringify({precargado_id:a,fuente:"DEMO"})}),enviarTexto:(a,e)=>f(`/leads/${a}/mensajes`,{method:"POST",body:JSON.stringify({tipo:"TEXTO",texto:e})}),enviarAudio:(a,e,t,o)=>f(`/leads/${a}/mensajes`,{method:"POST",body:JSON.stringify({tipo:"AUDIO",audio_base64:e,mime:t,duracion_s:o})}),conversacion:a=>f(`/leads/${a}/conversacion`),cola:()=>f("/leads"),ficha:a=>f(`/leads/${a}/ficha`),buyerPersona:a=>f(`/gerencia/buyer-persona${a?`?proyecto_id=${a}`:""}`),avanzarTiempo:a=>f("/demo/tiempo",{method:"POST",body:JSON.stringify({avanzar_hasta:a})}),reiniciar:()=>f("/demo/reiniciar",{method:"POST"})};function ne(a,e){a.innerHTML=e.map(ie).join("")}function ie(a){var s,r;if(a.tipo_contenido==="SISTEMA")return`<div class="pildora-sistema">${E(a.texto)}</div>`;const e=a.autor==="LEAD"?"derecha":"izquierda",t=new Date(a.creado_en).toLocaleTimeString("es-CO",{hour:"2-digit",minute:"2-digit"}),o=a.autor==="LEAD"?'<span class="chulos">✓✓</span>':"",n=(s=a.adjunto)!=null&&s.audio_original?'<span class="icono-audio" aria-label="nota de voz">🎙️</span>':"",i=(r=a.adjunto)!=null&&r.recomendaciones?se(a.adjunto.recomendaciones):"";return`
    <div class="burbuja ${e}">
      ${n}<p>${E(a.texto)}</p>
      <span class="hora">${t}${o}</span>
    </div>${i}`}function se(a){return`<div class="carrusel" role="list">${a.slice(0,3).map(e=>`
    <article class="tarjeta-proyecto" role="listitem">
      <header class="franja-azul">${E(e.nombre)}</header>
      <p class="zona">${E(e.zona)}</p>
      <p class="precio">Desde $${(e.precio_desde/1e6).toFixed(0)}M</p>
      <p class="razon">${E(e.razon)}</p>
      <p class="evidencia">${e.vecinos} personas con tu perfil compraron aquí ·
         ${(e.tasa_desistimiento*100).toFixed(0)}% desistió</p>
      <a class="btn-primario" href="${encodeURI(e.brochure_url)}" target="_blank" rel="noopener">Ver brochure</a>
      <a class="btn-secundario" href="${encodeURI(e.recorrido_360_url)}" target="_blank" rel="noopener">Recorrido 360°</a>
    </article>`).join("")}</div>`}function B(a,e){a.classList.toggle("visible",e)}function re(a){a.innerHTML=`
    <div class="chat-header">
      <div class="avatar-vivi">V</div>
      <div class="info-header">
        <span class="nombre-vivi">Vivi</span>
        <span class="escribiendo" id="indicador-escribiendo">escribiendo…</span>
      </div>
    </div>
    <div class="mensajes-scroll" id="mensajes-scroll">
      <div class="mensajes" id="contenedor-mensajes"></div>
    </div>
    <button class="btn-nuevos" id="btn-nuevos" aria-label="Ir a nuevos mensajes">↓ nuevos mensajes</button>
    <div class="barra-entrada">
      <button class="btn-mic" id="btn-mic" aria-label="Grabar nota de voz" type="button">🎤</button>
      <div class="mic-grabando" id="mic-grabando">
        <span class="punto-rojo"></span>
        <span class="contador-mic" id="contador-mic">0:00</span>
        <button class="btn-detener-mic" id="btn-detener-mic" type="button">■</button>
      </div>
      <input type="text" class="input-mensaje" id="input-mensaje"
             placeholder="Escribe un mensaje…" autocomplete="off" />
      <button class="btn-enviar" id="btn-enviar" aria-label="Enviar mensaje" type="button">➤</button>
    </div>
  `}function E(a){const e=document.createElement("div");return e.textContent=a,e.innerHTML}const ce=60,le=1500,de=300;let _=null,D=!1,v=null,S=[],L=null,g=0;function ue(a){re(a);const e=document.getElementById("contenedor-mensajes"),t=document.getElementById("mensajes-scroll"),o=document.getElementById("indicador-escribiendo"),n=document.getElementById("input-mensaje"),i=document.getElementById("btn-enviar"),s=document.getElementById("btn-mic"),r=document.getElementById("mic-grabando"),c=document.getElementById("contador-mic"),l=document.getElementById("btn-detener-mic"),d=document.getElementById("btn-nuevos");n.addEventListener("keydown",u=>{u.key==="Enter"&&!u.shiftKey&&(u.preventDefault(),F(n))}),i.addEventListener("click",()=>F(n)),s.addEventListener("click",()=>me(s,r,c,n,i)),l.addEventListener("click",()=>J(s,r,n,i)),t.addEventListener("scroll",()=>{D=t.scrollHeight-t.scrollTop-t.clientHeight>80,D||d.classList.remove("visible")}),d.addEventListener("click",()=>{t.scrollTop=t.scrollHeight,d.classList.remove("visible"),D=!1}),setInterval(()=>k(e,t,o,d),le),k(e,t,o,d)}async function k(a,e,t,o){const n=I();if(n.leadActivo)try{const i=await p.conversacion(n.leadActivo);b({conversacion:i}),ne(a,i.mensajes),i.turno_en_proceso?_||(_=setTimeout(()=>B(t,!0),de)):(_&&(clearTimeout(_),_=null),B(t,!1)),D?o.classList.add("visible"):e.scrollTop=e.scrollHeight}catch(i){i instanceof h&&i.estadoHTTP>=500&&$("Error de conexión. Reintentando…")}}async function F(a){const e=a.value.trim();if(!e)return;const t=I();if(t.leadActivo){a.value="";try{await p.enviarTexto(t.leadActivo,e)}catch(o){a.value=e,o instanceof h&&$(o.message)}}}async function me(a,e,t,o,n){try{const i=await navigator.mediaDevices.getUserMedia({audio:!0});v=new MediaRecorder(i),S=[],g=0,v.ondataavailable=s=>{S.push(s.data)},v.onstop=()=>{i.getTracks().forEach(s=>s.stop()),pe(a,e,o,n)},v.start(),a.style.display="none",o.style.display="none",n.style.display="none",e.classList.add("activo"),t.textContent="0:00",L=setInterval(()=>{g++;const s=Math.floor(g/60),r=g%60;t.textContent=`${s}:${r.toString().padStart(2,"0")}`,g>=ce&&J(a,e,o,n)},1e3)}catch{$("No se pudo acceder al micrófono.")}}function J(a,e,t,o){L&&(clearInterval(L),L=null),v&&v.state!=="inactive"&&v.stop(),e.classList.remove("activo"),a.style.display="",t.style.display="",o.style.display=""}async function pe(a,e,t,o){if(S.length===0)return;const n=I();if(!n.leadActivo)return;const i=new Blob(S,{type:S[0].type}),s=new FileReader;s.onloadend=async()=>{const r=s.result.split(",")[1],c=i.type,l=g;try{await p.enviarAudio(n.leadActivo,r,c,l)}catch(d){d instanceof h&&(d.codigo==="AUDIO_INVALIDO"?$("No te escuché bien, ¿me lo repites o me lo escribes?"):$(d.message))}},s.readAsDataURL(i),e.classList.remove("activo"),a.style.display="",t.style.display="",o.style.display=""}function $(a){const e=document.querySelector(".toast-error");e&&e.remove();const t=document.createElement("div");t.className="toast-error",t.textContent=a,document.body.appendChild(t),setTimeout(()=>t.remove(),4e3)}const fe={VERDE:"🟢",AMBAR:"🟡",GRIS:"⚪"};function ve(a,e,t,o){const n=e.cupo_10.usados,i=e.cupo_10.porcentaje_ventana,r=n/(i||1)*100>=80;a.innerHTML=`
    <div class="cola-container">
      <header class="cola-header">
        <h2>Leads en Cola (${e.leads.length})</h2>
        <div class="cupo-bar ${r?"alerta":""}"
             title="${r?"Cupo regulatorio de no afiliados casi lleno (≥80%)":"Uso del cupo regulatorio del 10%"}">
          <span>Cupo 10%:</span>
          <progress value="${n}" max="${i}"></progress>
          <span>${n}/${i}</span>
        </div>
      </header>
      <ul class="lista-leads" role="list">
        ${e.leads.map(be).join("")}
      </ul>
    </div>
  `,a.querySelectorAll("[data-lead-id]").forEach(c=>{const l=c.dataset.leadId,d=c.querySelector("[data-btn-chat]");d&&d.addEventListener("click",u=>{u.stopPropagation(),o(l)}),c.addEventListener("click",()=>t(l)),c.addEventListener("keydown",u=>{(u.key==="Enter"||u.key===" ")&&(u.preventDefault(),t(l))})})}function be(a){const e=fe[a.semaforo]??"⚪";return`
    <li class="fila-lead" data-lead-id="${a.lead_id}" tabindex="0" role="listitem">
      <span class="semaforo-dot" aria-label="Semáforo ${a.semaforo}">${e}</span>
      <span class="lead-nombre">${C(a.nombre)}</span>
      <span class="lead-ruta">${C(a.ruta)}</span>
      <span class="lead-prio-badge" title="Prioridad calculada">Prio ${a.prioridad.toFixed(2)}</span>
      <button class="btn-ver-chat" data-btn-chat="true" title="Ver chat en vivo con ${C(a.nombre)}" type="button">
        💬 Ver chat
      </button>
      <p class="lead-resumen">${C(a.resumen)}</p>
    </li>
  `}function C(a){const e=document.createElement("div");return e.textContent=a,e.innerHTML}const V={VERIFICADO_BASE:{txt:"VERIFICADO",icono:"✓",clase:"badge-verificado"},DECLARADO:{txt:"DECLARADO",icono:"✍",clase:"badge-declarado"},INFERIDO:{txt:"INFERIDO",icono:"~",clase:"badge-inferido"}};function O(a){const e=V[a]??V.DECLARADO;return`<span class="badge-fuente ${e.clase}" title="Fuente: ${e.txt}">${e.icono} ${e.txt}</span>`}function H(a,e,t="Lead"){if(!e){ge(a,t);return}const o=e.banda_advertencia?`<div class="banda-advertencia" role="alert">
         <span>⚠️</span> <span>${m(e.banda_advertencia)}</span>
       </div>`:"",n=(e.confianza_perfil*100).toFixed(0),i=e.identificacion,s=e.capacidad,r=e.intencion;a.innerHTML=`
    <div class="ficha-container">
      ${o}

      <!-- Identificación Principal (F-Layout Top) -->
      <header class="ficha-header-card">
        <div class="ficha-id-info">
          <h3>${m(i.nombre||t)}</h3>
          <div class="ficha-meta-grid">
            <div class="ficha-meta-item">
              <strong>Afiliada:</strong> ${i.afiliada?`Sí (Cat. ${m(i.categoria)})`:"No"}
            </div>
            <div class="ficha-meta-item">
              <strong>Teléfono:</strong> ${m(i.telefono||"No registrado")}
            </div>
            <div class="ficha-meta-item">
              <strong>Cupo 10%:</strong> ${e.consume_cupo_10?"Consume cupo":"No aplica"}
            </div>
          </div>
        </div>

        <div class="confianza-bar-container" title="Confianza general del perfilamiento">
          <span class="confianza-val">${n}% Confianza</span>
          <progress class="confianza-progress" value="${e.confianza_perfil}" max="1"></progress>
        </div>
      </header>

      <!-- Bloque 3 Columnas: Intención, Capacidad, Recomendación -->
      <div class="grid-tres-columnas">

        <!-- Columna 1: Capacidad Financiera -->
        <article class="card-seccion">
          <h4>Capacidad Financiera</h4>
          <div class="desglose-monto-item">
            <span>Presupuesto Máx:</span>
            <span class="monto-num">$${(s.presupuesto_max/1e6).toFixed(1)}M ${O("VERIFICADO_BASE")}</span>
          </div>
          <div class="desglose-monto-item">
            <span>Crédito Máx:</span>
            <span class="monto-num">$${(s.credito_max/1e6).toFixed(1)}M ${O("INFERIDO")}</span>
          </div>
          <div class="desglose-monto-item">
            <span>Subsidio Aplicable:</span>
            <span class="monto-num">$${(s.subsidio_aplicable/1e6).toFixed(1)}M ${O("VERIFICADO_BASE")}</span>
          </div>
          <div class="desglose-monto-item">
            <span>Recursos Propios:</span>
            <span class="monto-num">$${(s.recursos_propios/1e6).toFixed(1)}M ${O("DECLARADO")}</span>
          </div>
        </article>

        <!-- Columna 2: Intención de Compra (Nivel + Confianza + Señales) -->
        <article class="card-seccion">
          <h4>Intención de Compra</h4>
          <p class="nivel-destacado nivel-${r.nivel}">
            Nivel ${m(r.nivel)} <small style="font-size:0.75rem; color:#6B7280">(Confianza ${m(r.confianza)})</small>
          </p>
          <ul class="lista-puntos">
            ${r.senales.map(l=>`<li>${m(l)}</li>`).join("")}
          </ul>
        </article>

        <!-- Columna 3: Argumentos y Beneficios -->
        <article class="card-seccion">
          <h4>Argumentos de Venta</h4>
          <ul class="lista-puntos">
            ${e.argumentos_venta.map(l=>`<li>${m(l)}</li>`).join("")}
          </ul>
          <h4 style="margin-top:0.5rem">Beneficios Colsubsidio</h4>
          <ul class="lista-puntos">
            ${e.beneficios.map(l=>`<li>${m(l)}</li>`).join("")}
          </ul>
        </article>
      </div>

      <!-- Banner Azul Siguiente Paso (Doc 12 §3.4) -->
      <div class="banner-siguiente-paso">
        <span>▶ Siguiente paso: agendar visita sala de ventas</span>
        <button class="btn-copiar-resumen" id="btn-copiar-resumen" type="button">
          📋 Copiar Resumen
        </button>
      </div>

      <!-- Timeline Plegado de la Conversación -->
      <details class="timeline-plegable">
        <summary>💬 Ver historial de perfilamiento con Vivi</summary>
        <p style="font-size:0.83rem; margin-top:0.5rem; color:#4B5563">
          Conversación iniciada. Canal WhatsApp. Transcripción procesada por Vivi.
        </p>
      </details>
    </div>
  `;const c=a.querySelector("#btn-copiar-resumen");c&&c.addEventListener("click",()=>{const l=`RESUMEN FICHA - ${i.nombre}
Presupuesto: $${(s.presupuesto_max/1e6).toFixed(1)}M | Afiliada: ${i.afiliada?"Sí":"No"}
Intención: ${r.nivel}
Siguiente paso: Agendar visita sala de ventas.`;navigator.clipboard.writeText(l).then(()=>{c.textContent="✓ Copiado!",setTimeout(()=>{c.textContent="📋 Copiar Resumen"},2e3)})})}function ge(a,e){a.innerHTML=`
    <div class="estado-vacio-amable">
      <h3>💬 Este lead aún está conversando con Vivi</h3>
      <p style="margin-bottom: 1rem; font-size: 0.9rem;">
        La ficha comercial completa de <strong>${m(e)}</strong> se generará automáticamente cuando Vivi complete la calificación.
      </p>
      <details class="timeline-plegable" style="max-width: 400px; margin: 0 auto; text-align: left;">
        <summary>Ver avance actual</summary>
        <p style="font-size:0.82rem; margin-top:0.5rem; color:#4B5563">
          Estado: En proceso de perfilamiento en tiempo real desde el chat.
        </p>
      </details>
    </div>
  `}function m(a){const e=document.createElement("div");return e.textContent=a,e.innerHTML}const he=[{id:"mongui",nombre:"Monguí"},{id:"macarena",nombre:"La Macarena"},{id:"versalles",nombre:"Versalles"},{id:"todos",nombre:"Todos los proyectos"}];function G(a,e,t="mongui",o){const n=e?e.muestras:312,i=e?(e.tasa_desistimiento*100).toFixed(0):"11",s=(e==null?void 0:e.afiliacion)??{afiliados:180,no_afiliados:132},r=(e==null?void 0:e.categoria)??{"Cat A":90,"Cat B":65,"Cat C":25,"No Afiliado":132},c=(e==null?void 0:e.rango_edad)??{"18-25":20,"26-35":150,"36-45":90,"46+":52};a.innerHTML=`
    <div class="gerencia-container">
      <header class="gerencia-header">
        <div class="selector-proyecto-wrap">
          <label for="select-proyecto-gerencia">Proyecto:</label>
          <select id="select-proyecto-gerencia" class="select-proyecto">
            ${he.map(d=>`
              <option value="${d.id}" ${d.id===t?"selected":""}>
                ${X(d.nombre)}
              </option>
            `).join("")}
          </select>
        </div>

        <span class="nota-actualizacion">
          ℹ️ Se actualiza en tiempo real con cada lead perfilado
        </span>
      </header>

      <!-- Métricas Clave -->
      <div class="metricas-row">
        <div class="card-metrica">
          <div class="metrica-label">Personas en vivo interesadas</div>
          <div class="metrica-valor">${n}</div>
        </div>
        <div class="card-metrica">
          <div class="metrica-label">Tasa de Desistimiento Histórica</div>
          <div class="metrica-valor">${i}%</div>
        </div>
      </div>

      <!-- Gráficos de Barras Horizontales (Estilo Junta Directiva) -->
      <div class="grid-tres-columnas">

        <!-- Gráfico 1: Afiliación -->
        <article class="grafico-barras-card">
          <h4>Distribución por Afiliación</h4>
          <div class="barras-lista">
            ${j(s)}
          </div>
        </article>

        <!-- Gráfico 2: Categoría -->
        <article class="grafico-barras-card">
          <h4>Categoría de Afiliación</h4>
          <div class="barras-lista">
            ${j(r)}
          </div>
        </article>

        <!-- Gráfico 3: Rango de Edad -->
        <article class="grafico-barras-card">
          <h4>Rango de Edad</h4>
          <div class="barras-lista">
            ${j(c)}
          </div>
        </article>

      </div>
    </div>
  `;const l=a.querySelector("#select-proyecto-gerencia");l&&l.addEventListener("change",()=>{o(l.value)})}function j(a){const e=Math.max(...Object.values(a),1);return Object.entries(a).map(([t,o])=>{const n=(o/e*100).toFixed(0);return`
      <div class="barra-item">
        <span class="barra-label">${X(t)}</span>
        <div class="barra-track">
          <div class="barra-fill" style="width: ${n}%"></div>
        </div>
        <span class="barra-val">${o}</span>
      </div>
    `}).join("")}function X(a){const e=document.createElement("div");return e.textContent=a,e.innerHTML}function _e(a,e,t){a.innerHTML=`
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
  `;const o=a.querySelector("#btn-avanzar-tiempo");o&&o.addEventListener("click",e);const n=a.querySelector("#btn-reiniciar-demo");n&&n.addEventListener("click",t)}const ye=5e3;let R="mongui",T=null,z=null;function Ae(a,e,t){if(t){const o=t.querySelectorAll("button[data-tab]");o.forEach(n=>{n.addEventListener("click",()=>{const i=n.dataset.tab;Ee(i,o)})})}e&&Ce(e),te(()=>q(a)),setInterval(()=>x(),ye),x(),q(a)}async function x(){try{const a=await p.cola();b({cola:a})}catch(a){console.warn("[dashboard] Error cargando cola:",a)}}function Ee(a,e){b({tabActiva:a}),e.forEach(t=>{const o=t.dataset.tab===a;t.setAttribute("aria-selected",o?"true":"false")})}function q(a){const e=I();switch(e.tabActiva){case"cola":T=null,z=null,e.cola?ve(a,e.cola,t=>Se(t),t=>$e(t)):a.innerHTML='<div style="padding:1rem; color:#6B7280">Cargando cola de leads…</div>';break;case"ficha":if(!e.leadActivo){T=null,a.innerHTML='<div style="padding:1.5rem; text-align:center; color:#6B7280">Selecciona un lead de la cola para ver su ficha comercial.</div>';break}e.leadActivo!==T&&(T=e.leadActivo,Ie(a,e.leadActivo));break;case"gerencia":R!==z&&(z=R,Y(a,R));break}}function Se(a){b({leadActivo:a,tabActiva:"ficha"});const e=document.querySelector(".tabs");e&&e.querySelectorAll("button[data-tab]").forEach(o=>o.setAttribute("aria-selected",o.dataset.tab==="ficha"?"true":"false"))}function $e(a){b({leadActivo:a});const e=document.getElementById("panel-chat");e&&e.scrollIntoView({behavior:"smooth"})}async function Ie(a,e){var i;const o=(i=I().cola)==null?void 0:i.leads.find(s=>s.lead_id===e),n=(o==null?void 0:o.nombre)??"Lead";try{const s=await p.ficha(e);H(a,s,n)}catch(s){s instanceof h&&s.estadoHTTP===404?H(a,null,n):a.innerHTML=`<div class="banda-advertencia">⚠️ Error cargando la ficha comercial: ${s.message}</div>`}}async function Y(a,e){const t=o=>{R=o,z=o,Y(a,o)};try{const o=await p.buyerPersona(e);G(a,o,e,t)}catch{G(a,null,e,t)}}function y(a,e=!1){const t=document.getElementById("demo-aviso");t&&(t.textContent=a,t.classList.toggle("es-error",e),setTimeout(()=>{t.textContent="",t.classList.remove("es-error")},4e3))}function Ce(a){_e(a,async()=>{const e=document.getElementById("demo-fecha"),t=e==null?void 0:e.value;if(!t){y("Elegí una fecha primero.",!0);return}try{const o=await p.avanzarTiempo(t);y(`Tiempo avanzado a ${t} · ${o.hitos_disparados} hito(s) disparado(s).`),x()}catch(o){y(`Error al avanzar tiempo: ${o.message}`,!0)}},async()=>{try{await p.reiniciar();let e=null;try{e=(await p.crearLead("ana")).lead_id}catch(t){console.error("[dashboard] no se pudo recrear el lead tras reiniciar:",t)}b({leadActivo:e,tabActiva:"cola"}),x(),y("Demo reiniciado al estado inicial.")}catch(e){y(`Error al reiniciar demo: ${e.message}`,!0)}})}const W=new URLSearchParams(location.search);W.get("mock")==="1"&&ae();const Oe=["ana","carlos","luisa"];function Te(){const a=W.get("precargado");return Oe.includes(a??"")?a:"ana"}async function De(){try{return(await p.crearLead(Te())).lead_id}catch(a){const e=a instanceof h?`${a.codigo}: ${a.message}`:String(a);return console.error("[main] no se pudo crear el lead inicial:",e),Le(e),null}}function Le(a){const e=document.getElementById("panel-chat");if(!e)return;const t=document.createElement("div");t.className="banda-advertencia",t.setAttribute("role","alert"),t.textContent=`No se pudo iniciar la conversación (${a}). Recargá la página o revisá que el backend esté arriba.`,e.prepend(t)}async function Re(){const a=document.getElementById("panel-chat");if(a){const n=await De();n&&(b({leadActivo:n}),ue(a))}const e=document.getElementById("contenido-tab"),t=document.querySelector(".tabs"),o=document.getElementById("botonera-demo");e&&Ae(e,o,t),console.info("Vivi web iniciado completo (Chat + Dashboard)")}Re();
