(function(){const a=document.createElement("link").relList;if(a&&a.supports&&a.supports("modulepreload"))return;for(const n of document.querySelectorAll('link[rel="modulepreload"]'))o(n);new MutationObserver(n=>{for(const i of n)if(i.type==="childList")for(const s of i.addedNodes)s.tagName==="LINK"&&s.rel==="modulepreload"&&o(s)}).observe(document,{childList:!0,subtree:!0});function t(n){const i={};return n.integrity&&(i.integrity=n.integrity),n.referrerPolicy&&(i.referrerPolicy=n.referrerPolicy),n.crossOrigin==="use-credentials"?i.credentials="include":n.crossOrigin==="anonymous"?i.credentials="omit":i.credentials="same-origin",i}function o(n){if(n.ep)return;n.ep=!0;const i=t(n);fetch(n.href,i)}})();let S=1,k=!1;const R=[{mensaje_id:"m1",autor:"VIVI",tipo_contenido:"TEXTO",texto:"¡Hola Ana! 👋 Como afiliada tienes un subsidio de hasta $52,5M. ¿Sueñas con comprar este año?",creado_en:new Date().toISOString(),adjunto:null}];function ae(e){return{lead_id:e,estado:"PERFILANDO",turno_en_proceso:k,mensajes:[...R]}}const te={cupo_10:{usados:1,porcentaje_ventana:10},leads:[{lead_id:"mock-1",nombre:"Ana Rodríguez",estado:"ENTREGADO",ruta:"ASESOR",afiliado:!0,semaforo:"VERDE",prioridad:.91,resumen:"Afiliada cat. A · presupuesto $166.8M · intención alta",actualizado_en:new Date().toISOString()},{lead_id:"mock-2",nombre:"Carlos Martínez",estado:"PERFILANDO",ruta:"ASESOR",afiliado:!1,semaforo:"AMBAR",prioridad:.78,resumen:"No afiliado · presupuesto $210M · intención media",actualizado_en:new Date().toISOString()},{lead_id:"mock-3",nombre:"Luisa Gómez",estado:"CALIFICADO",ruta:"NUTRICION",afiliado:!0,semaforo:"VERDE",prioridad:.85,resumen:"Afiliada cat. B · presupuesto $145M · requiere plan cuota inicial",actualizado_en:new Date().toISOString()}]},oe={"mock-1":{ficha_id:"f1",lead_id:"mock-1",generada_en:new Date().toISOString(),confianza_perfil:.94,banda_advertencia:null,identificacion:{nombre:"Ana Rodríguez",afiliada:!0,categoria:"A",telefono:"+57 300 123 4567"},capacidad:{presupuesto_max:1668e5,credito_max:1143e5,subsidio_aplicable:525e5,recursos_propios:12e6,ratio:.28,confianza:.94,desglose:[{concepto:"Subsidio Mi Casa Ya / Caja",monto:525e5,regla:"Afiliado Cat A",fuente:"VERIFICADO_BASE"},{concepto:"Preaprobado Bancolombia",monto:1143e5,regla:"Capacidad de endeudamiento 30%",fuente:"INFERIDO"},{concepto:"Ahorro Declarado",monto:12e6,regla:"Declarado en chat",fuente:"DECLARADO"}]},perfil:{ingreso:{valor:26e5,fuente:"VERIFICADO_BASE",confianza:.95,requiere_confirmacion:!1,actualizado_en:new Date().toISOString()}},intencion:{nivel:"ALTA",confianza:"ALTA",senales:["Busca comprar antes de 6 meses","Tiene ahorro inicial preparado","Responde inmediatamente a Vivi"]},recomendaciones:[{proyecto_id:"mongui",nombre:"Monguí",zona:"Ciudadela Maiporé - Soacha",precio_desde:15647e4,razon:"Tu presupuesto cubre el 100% de la cuota inicial",vecinos:622,tasa_desistimiento:.12,brochure_url:"https://heyzine.com/flip-book/866af8f6a6.html",recorrido_360_url:"https://storage.net-fs.com/hosting/7532170/19/"}],beneficios:["Subsidio de vivienda Colsubsidio hasta $52,5M","Tasa preferencial crédito hipotecario"],argumentos_venta:["Cuota estimada mensual ($650k) es MENOR al arriendo promedio de la zona ($850k)","Entrega inmediata en 2026"],alerta_desistimiento:{activa:!1,tasa_vecinos:.12,detalle:null},consume_cupo_10:!1},"mock-2":{ficha_id:"f2",lead_id:"mock-2",generada_en:new Date().toISOString(),confianza_perfil:.82,banda_advertencia:"No afiliado a Colsubsidio — consume cupo del 10% regulatorio",identificacion:{nombre:"Carlos Martínez",afiliada:!1,categoria:"N/A",telefono:"+57 311 987 6543"},capacidad:{presupuesto_max:21e7,credito_max:18e7,subsidio_aplicable:0,recursos_propios:3e7,ratio:.32,confianza:.82,desglose:[{concepto:"Crédito solicitado",monto:18e7,regla:"Ingresos independientes",fuente:"DECLARADO"},{concepto:"Cuota Inicial Propia",monto:3e7,regla:"Recursos declarados",fuente:"DECLARADO"}]},perfil:{},intencion:{nivel:"MEDIA",confianza:"MEDIA",senales:["Interesado en proyectos VIS y no VIS","Evaluando opciones de crédito"]},recomendaciones:[],beneficios:["Opción de crédito con aliados de la caja"],argumentos_venta:["Proyecto Versalles cuenta con certificación EDGE para ahorro energético"],alerta_desistimiento:{activa:!1,tasa_vecinos:.08,detalle:null},consume_cupo_10:!0}},B={mongui:{proyecto_id:"mongui",nombre:"Monguí",muestras:312,afiliacion:{afiliados:198,no_afiliados:114},categoria:{A:110,B:68,C:20,SIN_DATO:114},rango_edad:{"20-35":165,"36-45":82,"46-55":40,"55+":12,SIN_DATO:13},tasa_desistimiento:.11,actualizado_en:new Date().toISOString()},macarena:{proyecto_id:"macarena",nombre:"La Macarena",muestras:185,afiliacion:{afiliados:140,no_afiliados:45},categoria:{A:85,B:42,C:13,SIN_DATO:45},rango_edad:{"20-35":100,"36-45":50,"46-55":20,"55+":8,SIN_DATO:7},tasa_desistimiento:.08,actualizado_en:new Date().toISOString()},versalles:{proyecto_id:"versalles",nombre:"Versalles",muestras:142,afiliacion:{afiliados:85,no_afiliados:57},categoria:{A:40,B:32,C:13,SIN_DATO:57},rango_edad:{"20-35":70,"36-45":42,"46-55":15,"55+":8,SIN_DATO:7},tasa_desistimiento:.15,actualizado_en:new Date().toISOString()}};function ne(e){k=!0,setTimeout(()=>{S++;const a=/proyecto|comprar|vivienda|casa|apto/i.test(e),t={mensaje_id:`m${S}`,autor:"VIVI",tipo_contenido:a?"TARJETAS_PROYECTOS":"TEXTO",texto:a?"Basándome en tu perfil, estos proyectos te pueden interesar:":"¡Qué bueno saberlo! ¿Te interesaría que exploremos opciones de vivienda según tu presupuesto?",creado_en:new Date().toISOString(),adjunto:a?{recomendaciones:[{proyecto_id:"mongui",nombre:"Monguí",zona:"Ciudadela Maiporé - Soacha",precio_desde:15647e4,razon:"Tu presupuesto cubre el 100% de la cuota inicial",vecinos:622,tasa_desistimiento:.12,brochure_url:"https://heyzine.com/flip-book/866af8f6a6.html",recorrido_360_url:"https://storage.net-fs.com/hosting/7532170/19/"},{proyecto_id:"macarena",nombre:"La Macarena",zona:"Ciudadela Maiporé - Soacha",precio_desde:12834e4,razon:"El más económico de la zona, ideal para tu ingreso",vecinos:374,tasa_desistimiento:.08,brochure_url:"https://heyzine.com/flip-book/b168b2f5ba.html",recorrido_360_url:""},{proyecto_id:"versalles",nombre:"Versalles",zona:"Ciudadela Maiporé - Soacha",precio_desde:1952e5,razon:"Certificación EDGE, ahorro en servicios",vecinos:174,tasa_desistimiento:.15,brochure_url:"https://heyzine.com/flip-book/be784b0d5c.html",recorrido_360_url:"https://shape.com.co/360/COLSUBSIDIO-Versalles_APTOA"}]}:null};R.push(t),k=!1},2e3)}function ie(){const e=window.fetch;window.fetch=async(a,t)=>{const o=a.toString(),n=((t==null?void 0:t.method)??"GET").toUpperCase(),i=(s,r=200)=>new Response(JSON.stringify(s),{status:r,headers:{"Content-Type":"application/json"}});if(o.includes("/conversacion")&&n==="GET"){const s=o.match(/\/leads\/([^/]+)\/conversacion/);return i(ae(s?s[1]:"mock-1"))}if(o.includes("/mensajes")&&n==="POST"){S++;let s={};try{s=JSON.parse((t==null?void 0:t.body)||"{}")}catch{s={}}const r=s.tipo==="AUDIO";return R.push({mensaje_id:`m${S}`,autor:"LEAD",tipo_contenido:"TEXTO",texto:r?"🎙️ [Nota de voz]":s.texto||"",creado_en:new Date().toISOString(),adjunto:r?{audio_original:!0}:null}),ne(r?"audio":s.texto||""),i({mensaje_id:`m${S}`,turno_en_proceso:!0},201)}if(o.includes("/ficha")&&n==="GET"){const s=o.match(/\/leads\/([^/]+)\/ficha/),r=s?s[1]:"mock-1",d=oe[r];return d?i(d):i({error:{codigo:"FICHA_NO_DISPONIBLE",mensaje:"Ficha aún no disponible"}},404)}if(o.includes("/gerencia/buyer-persona")&&n==="GET"){const r=new URL(o,"http://localhost").searchParams.get("proyecto_id")??"mongui",d=B[r]??B.mongui;return i(d)}if(o.includes("/demo/tiempo")&&n==="POST")return i({fecha_simulada:new Date().toISOString(),hitos_disparados:2});if(o.includes("/demo/reiniciar")&&n==="POST")return R.length=1,i({reiniciado:!0,fecha_simulada:new Date().toISOString()});if(o.endsWith("/api/leads")&&n==="POST"){let s={};try{s=JSON.parse((t==null?void 0:t.body)||"{}")}catch{s={}}const d={ana:"mock-1",carlos:"mock-2",luisa:"mock-3"}[s.precargado_id??"ana"]??"mock-1";return i({lead_id:d,estado:"PERFILANDO",afiliado_detectado:d!=="mock-2"},201)}return o.endsWith("/api/leads")&&n==="GET"?i(te):e(a,t)},console.info("[mock] servidor simulado completo activo (Chat + Ficha + Cola + Gerencia + Demo)")}const x={leadActivo:null,conversacion:null,cola:null,tabActiva:"cola"},Z=[];function O(){return x}function b(e){Object.assign(x,e),Z.forEach(a=>a(x))}function se(e){Z.push(e)}const re="/api";class g extends Error{constructor(a,t,o){super(t),this.codigo=a,this.estadoHTTP=o,this.name="ErrorAPI"}}async function f(e,a){var o,n;const t=await fetch(`${re}${e}`,{headers:{"Content-Type":"application/json"},...a});if(!t.ok){const i=await t.json().catch(()=>null);throw new g(((o=i==null?void 0:i.error)==null?void 0:o.codigo)??"ERROR_INTERNO",((n=i==null?void 0:i.error)==null?void 0:n.mensaje)??`HTTP ${t.status}`,t.status)}return await t.json()}const p={crearLead:e=>f("/leads",{method:"POST",body:JSON.stringify({precargado_id:e,fuente:"DEMO"})}),enviarTexto:(e,a)=>f(`/leads/${e}/mensajes`,{method:"POST",body:JSON.stringify({tipo:"TEXTO",texto:a})}),enviarAudio:(e,a,t,o)=>f(`/leads/${e}/mensajes`,{method:"POST",body:JSON.stringify({tipo:"AUDIO",audio_base64:a,mime:t,duracion_s:o})}),conversacion:e=>f(`/leads/${e}/conversacion`),cola:()=>f("/leads"),ficha:e=>f(`/leads/${e}/ficha`),buyerPersona:e=>f(`/gerencia/buyer-persona${e?`?proyecto_id=${e}`:""}`),avanzarTiempo:e=>f("/demo/tiempo",{method:"POST",body:JSON.stringify({avanzar_hasta:e})}),reiniciar:()=>f("/demo/reiniciar",{method:"POST"})};function ce(e,a){e.innerHTML=a.map(de).join("")}function de(e){var s,r;if(e.tipo_contenido==="SISTEMA")return`<div class="pildora-sistema">${$(e.texto)}</div>`;const a=e.autor==="LEAD"?"derecha":"izquierda",t=new Date(e.creado_en).toLocaleTimeString("es-CO",{hour:"2-digit",minute:"2-digit"}),o=e.autor==="LEAD"?'<span class="chulos">✓✓</span>':"",n=(s=e.adjunto)!=null&&s.audio_original?'<span class="icono-audio" aria-label="nota de voz">🎙️</span>':"",i=(r=e.adjunto)!=null&&r.recomendaciones?le(e.adjunto.recomendaciones):"";return`
    <div class="burbuja ${a}">
      ${n}<p>${$(e.texto)}</p>
      <span class="hora">${t}${o}</span>
    </div>${i}`}function le(e){return`<div class="carrusel" role="list">${e.slice(0,3).map(a=>`
    <article class="tarjeta-proyecto" role="listitem">
      <header class="franja-azul">${$(a.nombre)}</header>
      <p class="zona">${$(a.zona)}</p>
      <p class="precio">Desde $${(a.precio_desde/1e6).toFixed(0)}M</p>
      <p class="razon">${$(a.razon)}</p>
      <p class="evidencia">${a.vecinos} personas con tu perfil compraron aquí ·
         ${(a.tasa_desistimiento*100).toFixed(0)}% desistió</p>
      <a class="btn-primario" href="${encodeURI(a.brochure_url)}" target="_blank" rel="noopener">Ver brochure</a>
      <a class="btn-secundario" href="${encodeURI(a.recorrido_360_url)}" target="_blank" rel="noopener">Recorrido 360°</a>
    </article>`).join("")}</div>`}function F(e,a){e.classList.toggle("visible",a)}function ue(e){e.innerHTML=`
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
  `}function $(e){const a=document.createElement("div");return a.textContent=e,a.innerHTML}const me=60,pe=1500,fe=300;let _=null,D=!1,v=null,I=[],L=null,h=0;function ve(e){ue(e);const a=document.getElementById("contenedor-mensajes"),t=document.getElementById("mensajes-scroll"),o=document.getElementById("indicador-escribiendo"),n=document.getElementById("input-mensaje"),i=document.getElementById("btn-enviar"),s=document.getElementById("btn-mic"),r=document.getElementById("mic-grabando"),d=document.getElementById("contador-mic"),l=document.getElementById("btn-detener-mic"),c=document.getElementById("btn-nuevos");n.addEventListener("keydown",u=>{u.key==="Enter"&&!u.shiftKey&&(u.preventDefault(),q(n))}),i.addEventListener("click",()=>q(n)),s.addEventListener("click",()=>be(s,r,d,n,i)),l.addEventListener("click",()=>K(s,r,n,i)),t.addEventListener("scroll",()=>{D=t.scrollHeight-t.scrollTop-t.clientHeight>80,D||c.classList.remove("visible")}),c.addEventListener("click",()=>{t.scrollTop=t.scrollHeight,c.classList.remove("visible"),D=!1}),setInterval(()=>V(a,t,o,c),pe),V(a,t,o,c)}async function V(e,a,t,o){const n=O();if(n.leadActivo)try{const i=await p.conversacion(n.leadActivo);b({conversacion:i}),ce(e,i.mensajes),i.turno_en_proceso?_||(_=setTimeout(()=>F(t,!0),fe)):(_&&(clearTimeout(_),_=null),F(t,!1)),D?o.classList.add("visible"):a.scrollTop=a.scrollHeight}catch(i){i instanceof g&&i.estadoHTTP>=500&&C("Error de conexión. Reintentando…")}}async function q(e){const a=e.value.trim();if(!a)return;const t=O();if(t.leadActivo){e.value="";try{await p.enviarTexto(t.leadActivo,a)}catch(o){e.value=a,o instanceof g&&C(o.message)}}}async function be(e,a,t,o,n){try{const i=await navigator.mediaDevices.getUserMedia({audio:!0});v=new MediaRecorder(i),I=[],h=0,v.ondataavailable=s=>{I.push(s.data)},v.onstop=()=>{i.getTracks().forEach(s=>s.stop()),he(e,a,o,n)},v.start(),e.style.display="none",o.style.display="none",n.style.display="none",a.classList.add("activo"),t.textContent="0:00",L=setInterval(()=>{h++;const s=Math.floor(h/60),r=h%60;t.textContent=`${s}:${r.toString().padStart(2,"0")}`,h>=me&&K(e,a,o,n)},1e3)}catch{C("No se pudo acceder al micrófono.")}}function K(e,a,t,o){L&&(clearInterval(L),L=null),v&&v.state!=="inactive"&&v.stop(),a.classList.remove("activo"),e.style.display="",t.style.display="",o.style.display=""}async function he(e,a,t,o){if(I.length===0)return;const n=O();if(!n.leadActivo)return;const i=new Blob(I,{type:I[0].type}),s=new FileReader;s.onloadend=async()=>{const r=s.result.split(",")[1],d=i.type,l=h;try{await p.enviarAudio(n.leadActivo,r,d,l)}catch(c){c instanceof g&&(c.codigo==="AUDIO_INVALIDO"?C("No te escuché bien, ¿me lo repites o me lo escribes?"):C(c.message))}},s.readAsDataURL(i),a.classList.remove("activo"),e.style.display="",t.style.display="",o.style.display=""}function C(e){const a=document.querySelector(".toast-error");a&&a.remove();const t=document.createElement("div");t.className="toast-error",t.textContent=e,document.body.appendChild(t),setTimeout(()=>t.remove(),4e3)}const H={VERDE:{etiqueta:"VERDE",zona:"zona-verde"},AMBAR:{etiqueta:"ÁMBAR",zona:"zona-ambar"},GRIS:{etiqueta:"GRIS",zona:"zona-gris"}};function ge(e,a,t,o){const n=a.cupo_10.usados,i=a.cupo_10.porcentaje_ventana,s=n/(i||1)*100,r=s>=80;e.innerHTML=`
    <div class="hoja-informe informe-cola">
      <header class="masthead-cola">
        <div class="masthead-cola-titulo">
          <span class="masthead-kicker">Registro de leads</span>
          <p class="masthead-cifra"><span class="cifra">${a.leads.length}</span> en cola</p>
        </div>
        <div class="medidor-cupo ${r?"alerta":""}"
             title="${r?"Cupo regulatorio de no afiliados casi lleno (≥80%)":"Uso del cupo regulatorio del 10%"}">
          <span class="medidor-cupo-etiqueta">Cupo regulatorio 10%</span>
          <div class="medidor-cupo-escala" role="img" aria-label="Cupo usado: ${n} de ${i}">
            ${_e(i)}
            <div class="medidor-cupo-relleno" style="transform: scaleX(${s/100})"></div>
          </div>
          <span class="medidor-cupo-cifra cifra">${n}/${i}</span>
        </div>
      </header>
      <ul class="registro-leads" role="list">
        ${a.leads.map(ye).join("")}
      </ul>
    </div>
  `,e.querySelectorAll("[data-lead-id]").forEach(d=>{const l=d.dataset.leadId,c=d.querySelector("[data-btn-chat]");c&&c.addEventListener("click",u=>{u.stopPropagation(),o(l)}),d.addEventListener("click",()=>t(l)),d.addEventListener("keydown",u=>{(u.key==="Enter"||u.key===" ")&&(u.preventDefault(),t(l))})})}function ye(e){const a=H[e.semaforo]??H.GRIS,t=e.nombre.trim().charAt(0).toUpperCase()||"?";return`
    <li class="fila-lead ${a.zona}" data-lead-id="${e.lead_id}" tabindex="0" role="listitem">
      <span class="fila-lead-avatar" aria-hidden="true">${A(t)}</span>
      <div class="fila-lead-cuerpo">
        <div class="fila-lead-linea1">
          <span class="lead-nombre">${A(e.nombre)}</span>
          <span class="fila-lead-afiliado ${e.afiliado?"es-afiliado":""}">${e.afiliado?"Afiliado":"No afiliado"}</span>
          <span class="fila-lead-zona" aria-label="Semáforo ${e.semaforo}">${a.etiqueta}</span>
          <span class="lead-ruta">${A(e.ruta)}</span>
        </div>
        <p class="lead-resumen">${A(e.resumen)}</p>
      </div>
      <span class="fila-lead-prioridad" title="Prioridad calculada">
        <span class="fila-lead-prioridad-etiqueta">Prioridad</span>
        <span class="fila-lead-prioridad-cifra cifra">${e.prioridad.toFixed(2)}</span>
      </span>
      <button class="btn-ver-chat" data-btn-chat="true" title="Ver chat en vivo con ${A(e.nombre)}" type="button">
        Ver chat
      </button>
    </li>
  `}function _e(e){const a=Math.max(1,Math.round(e));return Array.from({length:a-1},(t,o)=>`<span class="medidor-tick" style="left: ${(o+1)/a*100}%"></span>`).join("")}function A(e){const a=document.createElement("div");return a.textContent=e,a.innerHTML}const G={VERIFICADO_BASE:{txt:"VERIFICADO",clase:"sello-verificado"},DECLARADO:{txt:"DECLARADO",clase:"sello-declarado"},INFERIDO:{txt:"INFERIDO",clase:"sello-inferido"}},U={ALTA:"zona-verde",MEDIA:"zona-ambar",BAJA:"zona-gris"};function Ae(e){const a=G[e]??G.DECLARADO;return`<span class="sello-fuente ${a.clase}" title="Fuente: ${a.txt}">${a.txt}</span>`}function J(e,a,t="Lead"){if(!a){Se(e,t);return}const o=a.banda_advertencia?`<div class="banda-advertencia" role="alert">
         <span class="banda-advertencia-etiqueta">Alerta</span>
         <span>${m(a.banda_advertencia)}</span>
       </div>`:"",n=Math.round(a.confianza_perfil*100),i=a.identificacion,s=a.capacidad,r=a.intencion,d=U[r.nivel]??U.BAJA;e.innerHTML=`
    <div class="hoja-informe informe-ficha">
      ${o}

      <!-- Masthead: identidad a la izquierda, medidor de confianza a la derecha -->
      <header class="masthead-ficha">
        <div class="masthead-identidad">
          <span class="masthead-kicker">Ficha comercial</span>
          <h3 class="masthead-nombre">${m(i.nombre||t)}</h3>
          <dl class="registro-identidad">
            <div class="registro-identidad-item">
              <dt>Afiliación</dt>
              <dd>${i.afiliada?`Afiliado · Categoría ${m(i.categoria)}`:"No afiliado"}</dd>
            </div>
            <div class="registro-identidad-item">
              <dt>Teléfono</dt>
              <dd class="cifra">${m(i.telefono||"No registrado")}</dd>
            </div>
            <div class="registro-identidad-item">
              <dt>Cupo 10%</dt>
              <dd>${a.consume_cupo_10?"Consume cupo":"No aplica"}</dd>
            </div>
          </dl>
        </div>

        <div class="medidor-confianza" role="img" aria-label="Confianza del perfil: ${n} por ciento">
          <span class="medidor-confianza-etiqueta">Confianza del perfil</span>
          <div class="medidor-confianza-escala">
            <span class="medidor-confianza-tick" style="left: 25%"></span>
            <span class="medidor-confianza-tick" style="left: 50%"></span>
            <span class="medidor-confianza-tick" style="left: 75%"></span>
            <div class="medidor-confianza-relleno" style="transform: scaleX(${n/100})"></div>
            <div class="medidor-confianza-marcador" style="left: ${n}%"></div>
          </div>
          <span class="medidor-confianza-cifra cifra">${n}<span class="medidor-confianza-unidad">/100</span></span>
        </div>
      </header>

      <!-- Cartera de capacidad: renglones de informe, no tarjetas -->
      <section class="seccion-informe seccion-capacidad">
        <h4 class="seccion-titulo">Capacidad financiera</h4>
        <div class="ledger-scroll">
          <table class="ledger">
            <thead>
              <tr>
                <th>Concepto</th>
                <th>Regla aplicada</th>
                <th class="ledger-num">Monto</th>
                <th>Fuente</th>
              </tr>
            </thead>
            <tbody>
              ${s.desglose.map(Ee).join("")}
            </tbody>
            <tfoot>
              <tr class="ledger-total">
                <td colspan="2">Presupuesto máximo</td>
                <td class="ledger-num cifra">$${P(s.presupuesto_max)}</td>
                <td></td>
              </tr>
              <tr class="ledger-ratio">
                <td colspan="2">Ratio de endeudamiento</td>
                <td class="ledger-num cifra">${(s.ratio*100).toFixed(0)}%</td>
                <td class="cifra">conf. ${(s.confianza*100).toFixed(0)}%</td>
              </tr>
            </tfoot>
          </table>
        </div>
      </section>

      <div class="informe-columnas">
        <!-- Zona de intención: banda con el veredicto de nivel -->
        <section class="seccion-informe zona-veredicto ${d}">
          <h4 class="seccion-titulo">Intención de compra</h4>
          <p class="veredicto-nivel">Nivel <span class="cifra">${m(r.nivel)}</span></p>
          <p class="veredicto-confianza">Confianza: ${m(r.confianza)}</p>
          <ul class="lista-observaciones">
            ${r.senales.map(c=>`<li>${m(c)}</li>`).join("")}
          </ul>
        </section>

        <!-- Argumentos y beneficios -->
        <section class="seccion-informe">
          <h4 class="seccion-titulo">Argumentos de venta</h4>
          <ul class="lista-observaciones">
            ${a.argumentos_venta.map(c=>`<li>${m(c)}</li>`).join("")}
          </ul>
          <h5 class="subseccion-titulo">Beneficios Colsubsidio</h5>
          <ul class="lista-observaciones">
            ${a.beneficios.map(c=>`<li>${m(c)}</li>`).join("")}
          </ul>
        </section>
      </div>

      <!-- Banner de instrucción única (Doc 12 §3.4) -->
      <div class="banner-siguiente-paso">
        <span><strong class="banner-etiqueta">Siguiente paso</strong> Agendar visita a sala de ventas</span>
        <button class="btn-copiar-resumen" id="btn-copiar-resumen" type="button">
          Copiar resumen
        </button>
      </div>

      <!-- Nota al pie plegable -->
      <details class="timeline-plegable">
        <summary>Ver historial de perfilamiento con Vivi</summary>
        <p class="timeline-texto">
          Conversación iniciada. Canal WhatsApp. Transcripción procesada por Vivi.
        </p>
      </details>
    </div>
  `;const l=e.querySelector("#btn-copiar-resumen");l&&l.addEventListener("click",()=>{const c=`RESUMEN FICHA - ${i.nombre}
Presupuesto: $${P(s.presupuesto_max)} | Afiliada: ${i.afiliada?"Sí":"No"}
Intención: ${r.nivel}
Siguiente paso: Agendar visita sala de ventas.`;navigator.clipboard.writeText(c).then(()=>{l.textContent="Copiado",setTimeout(()=>{l.textContent="Copiar resumen"},2e3)})})}function Ee(e){return`
    <tr>
      <td>${m(e.concepto)}</td>
      <td class="ledger-regla">${m(e.regla)}</td>
      <td class="ledger-num cifra">$${P(e.monto)}</td>
      <td>${Ae(e.fuente)}</td>
    </tr>
  `}function P(e){return`${(e/1e6).toFixed(1)}M`}function Se(e,a){e.innerHTML=`
    <div class="hoja-informe informe-vacio">
      <span class="masthead-kicker">Ficha comercial</span>
      <h3>Aún sin generar</h3>
      <p>
        La ficha comercial completa de <strong>${m(a)}</strong> se generará automáticamente cuando Vivi complete la calificación.
      </p>
      <details class="timeline-plegable">
        <summary>Ver avance actual</summary>
        <p class="timeline-texto">
          Estado: En proceso de perfilamiento en tiempo real desde el chat.
        </p>
      </details>
    </div>
  `}function m(e){const a=document.createElement("div");return a.textContent=e,a.innerHTML}const X=[{id:"mongui",nombre:"Monguí"},{id:"macarena",nombre:"La Macarena"},{id:"versalles",nombre:"Versalles"},{id:"todos",nombre:"Todos los proyectos"}];function Y(e,a,t="mongui",o){var u;const n=a?a.muestras:312,i=a?(a.tasa_desistimiento*100).toFixed(0):"11",s=((u=X.find(y=>y.id===t))==null?void 0:u.nombre)??"Proyecto",r=(a==null?void 0:a.afiliacion)??{afiliados:180,no_afiliados:132},d=(a==null?void 0:a.categoria)??{"Cat A":90,"Cat B":65,"Cat C":25,"No Afiliado":132},l=(a==null?void 0:a.rango_edad)??{"18-25":20,"26-35":150,"36-45":90,"46+":52};e.innerHTML=`
    <div class="hoja-informe informe-gerencia">
      <header class="masthead-gerencia">
        <div class="masthead-gerencia-titulo">
          <span class="masthead-kicker">Buyer persona vivo</span>
          <h2 class="masthead-nombre">${w(s)}</h2>
        </div>
        <div class="selector-proyecto-wrap">
          <label for="select-proyecto-gerencia">Proyecto</label>
          <select id="select-proyecto-gerencia" class="select-proyecto">
            ${X.map(y=>`
              <option value="${y.id}" ${y.id===t?"selected":""}>
                ${w(y.nombre)}
              </option>
            `).join("")}
          </select>
        </div>
      </header>

      <p class="nota-actualizacion">Se actualiza en tiempo real con cada lead perfilado.</p>

      <!-- Franja de métricas clave -->
      <div class="franja-metricas">
        <div class="franja-metrica">
          <span class="franja-metrica-etiqueta">Personas en vivo interesadas</span>
          <span class="franja-metrica-cifra cifra">${n}</span>
        </div>
        <div class="franja-metrica">
          <span class="franja-metrica-etiqueta">Tasa de desistimiento histórica</span>
          <span class="franja-metrica-cifra cifra">${i}%</span>
        </div>
      </div>

      <!-- Columnas de distribución (estilo junta directiva) -->
      <section class="panel-graficos">
        <div class="columna-grafico">
          <h4 class="seccion-titulo">Distribución por afiliación</h4>
          <div class="barras-lista">
            ${N(r)}
          </div>
        </div>
        <div class="columna-grafico">
          <h4 class="seccion-titulo">Categoría de afiliación</h4>
          <div class="barras-lista">
            ${N(d)}
          </div>
        </div>
        <div class="columna-grafico">
          <h4 class="seccion-titulo">Rango de edad</h4>
          <div class="barras-lista">
            ${N(l)}
          </div>
        </div>
      </section>
    </div>
  `;const c=e.querySelector("#select-proyecto-gerencia");c&&c.addEventListener("change",()=>{o(c.value)})}const $e={afiliados:"Afiliados",no_afiliados:"No afiliados",SIN_DATO:"Sin dato",A:"Categoría A",B:"Categoría B",C:"Categoría C"};function Ie(e){return $e[e]??e}function N(e){const a=Math.max(...Object.values(e),1);return Object.entries(e).map(([t,o])=>{const n=(o/a*100).toFixed(0);return`
      <div class="barra-item">
        <span class="barra-label">${w(Ie(t))}</span>
        <div class="barra-track">
          <div class="barra-fill" style="transform: scaleX(${Number(n)/100})"></div>
        </div>
        <span class="barra-val cifra">${o}</span>
      </div>
    `}).join("")}function w(e){const a=document.createElement("div");return a.textContent=e,a.innerHTML}function Ce(e,a,t){e.innerHTML=`
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
  `;const o=e.querySelector("#btn-avanzar-tiempo");o&&o.addEventListener("click",a);const n=e.querySelector("#btn-reiniciar-demo");n&&n.addEventListener("click",t)}const Oe=5e3;let z="mongui",T=null,j=null;function Te(e,a,t){if(t){const o=t.querySelectorAll("button[data-tab]");o.forEach(n=>{n.addEventListener("click",()=>{const i=n.dataset.tab;De(i,o)})})}a&&Re(a),se(()=>W(e)),setInterval(()=>M(),Oe),M(),W(e)}async function M(){try{const e=await p.cola();b({cola:e})}catch(e){console.warn("[dashboard] Error cargando cola:",e)}}function De(e,a){b({tabActiva:e}),a.forEach(t=>{const o=t.dataset.tab===e;t.setAttribute("aria-selected",o?"true":"false")})}function W(e){const a=O();switch(a.tabActiva){case"cola":T=null,j=null,a.cola?ge(e,a.cola,t=>Le(t),t=>ze(t)):e.innerHTML='<div style="padding:1rem; color:#6B7280">Cargando cola de leads…</div>';break;case"ficha":if(!a.leadActivo){T=null,e.innerHTML='<div style="padding:1.5rem; text-align:center; color:#6B7280">Selecciona un lead de la cola para ver su ficha comercial.</div>';break}a.leadActivo!==T&&(T=a.leadActivo,je(e,a.leadActivo));break;case"gerencia":z!==j&&(j=z,Q(e,z));break}}function Le(e){b({leadActivo:e,tabActiva:"ficha"});const a=document.querySelector(".tabs");a&&a.querySelectorAll("button[data-tab]").forEach(o=>o.setAttribute("aria-selected",o.dataset.tab==="ficha"?"true":"false"))}function ze(e){b({leadActivo:e});const a=document.getElementById("panel-chat");a&&a.scrollIntoView({behavior:"smooth"})}async function je(e,a){var i;const o=(i=O().cola)==null?void 0:i.leads.find(s=>s.lead_id===a),n=(o==null?void 0:o.nombre)??"Lead";try{const s=await p.ficha(a);J(e,s,n)}catch(s){s instanceof g&&s.estadoHTTP===404?J(e,null,n):e.innerHTML=`<div class="banda-advertencia">⚠️ Error cargando la ficha comercial: ${s.message}</div>`}}async function Q(e,a){const t=o=>{z=o,j=o,Q(e,o)};try{const o=await p.buyerPersona(a);Y(e,o,a,t)}catch{Y(e,null,a,t)}}function E(e,a=!1){const t=document.getElementById("demo-aviso");t&&(t.textContent=e,t.classList.toggle("es-error",a),setTimeout(()=>{t.textContent="",t.classList.remove("es-error")},4e3))}function Re(e){Ce(e,async()=>{const a=document.getElementById("demo-fecha"),t=a==null?void 0:a.value;if(!t){E("Elegí una fecha primero.",!0);return}try{const o=await p.avanzarTiempo(t);E(`Tiempo avanzado a ${t} · ${o.hitos_disparados} hito(s) disparado(s).`),M()}catch(o){E(`Error al avanzar tiempo: ${o.message}`,!0)}},async()=>{try{await p.reiniciar();let a=null;try{a=(await p.crearLead("ana")).lead_id}catch(t){console.error("[dashboard] no se pudo recrear el lead tras reiniciar:",t)}b({leadActivo:a,tabActiva:"cola"}),M(),E("Demo reiniciado al estado inicial.")}catch(a){E(`Error al reiniciar demo: ${a.message}`,!0)}})}const ee=new URLSearchParams(location.search);ee.get("mock")==="1"&&ie();const Me=["ana","carlos","luisa"];function Ne(){const e=ee.get("precargado");return Me.includes(e??"")?e:"ana"}async function ke(){try{return(await p.crearLead(Ne())).lead_id}catch(e){const a=e instanceof g?`${e.codigo}: ${e.message}`:String(e);return console.error("[main] no se pudo crear el lead inicial:",a),xe(a),null}}function xe(e){const a=document.getElementById("panel-chat");if(!a)return;const t=document.createElement("div");t.className="banda-advertencia",t.setAttribute("role","alert"),t.textContent=`No se pudo iniciar la conversación (${e}). Recargá la página o revisá que el backend esté arriba.`,a.prepend(t)}async function Pe(){const e=document.getElementById("panel-chat");if(e){const n=await ke();n&&(b({leadActivo:n}),ve(e))}const a=document.getElementById("contenido-tab"),t=document.querySelector(".tabs"),o=document.getElementById("botonera-demo");a&&Te(a,o,t),console.info("Vivi web iniciado completo (Chat + Dashboard)")}function we(){document.querySelectorAll("[data-colapsar]").forEach(e=>{const a=document.getElementById(e.dataset.colapsar);a&&e.addEventListener("click",()=>{const t=a.classList.toggle("panel-colapsado");e.setAttribute("aria-expanded",String(!t))})})}we();Pe();
