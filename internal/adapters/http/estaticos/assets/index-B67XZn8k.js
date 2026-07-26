(function(){const a=document.createElement("link").relList;if(a&&a.supports&&a.supports("modulepreload"))return;for(const n of document.querySelectorAll('link[rel="modulepreload"]'))o(n);new MutationObserver(n=>{for(const s of n)if(s.type==="childList")for(const i of s.addedNodes)i.tagName==="LINK"&&i.rel==="modulepreload"&&o(i)}).observe(document,{childList:!0,subtree:!0});function t(n){const s={};return n.integrity&&(s.integrity=n.integrity),n.referrerPolicy&&(s.referrerPolicy=n.referrerPolicy),n.crossOrigin==="use-credentials"?s.credentials="include":n.crossOrigin==="anonymous"?s.credentials="omit":s.credentials="same-origin",s}function o(n){if(n.ep)return;n.ep=!0;const s=t(n);fetch(n.href,s)}})();let $=50,N=!1;const g={"mock-1":[{mensaje_id:"m1",autor:"VIVI",tipo_contenido:"TEXTO",texto:"¡Hola Ana! 👋 Como afiliada tienes un subsidio de hasta $52,5M. ¿Sueñas con comprar este año?",creado_en:new Date().toISOString(),adjunto:null}],"mock-2":[{mensaje_id:"m20",autor:"VIVI",tipo_contenido:"TEXTO",texto:"¡Hola Carlos! 👋 Vi tu interés en opciones de vivienda. ¿Tienes algún proyecto en mente o presupuesto inicial?",creado_en:new Date().toISOString(),adjunto:null},{mensaje_id:"m21",autor:"LEAD",tipo_contenido:"TEXTO",texto:"Hola, tengo ahorrados $30M y busco un apto con buena ubicación.",creado_en:new Date().toISOString(),adjunto:null},{mensaje_id:"m22",autor:"VIVI",tipo_contenido:"TEXTO",texto:"¡Excelente! Con tus recursos y capacidad crediticia, el proyecto Versalles aplica muy bien a tus metas.",creado_en:new Date().toISOString(),adjunto:null}],"mock-3":[{mensaje_id:"m30",autor:"VIVI",tipo_contenido:"TEXTO",texto:"¡Hola Luisa! 👋 Como afiliada Cat. B podemos armar un plan de ahorro progresivo para tu cuota inicial.",creado_en:new Date().toISOString(),adjunto:null},{mensaje_id:"m31",autor:"LEAD",tipo_contenido:"TEXTO",texto:"Me interesa mucho la zona de Soacha, ¿qué proyectos tienen?",creado_en:new Date().toISOString(),adjunto:null}]};function ie(e){const a=g[e]??g["mock-1"];return{lead_id:e,estado:"PERFILANDO",turno_en_proceso:N,mensajes:[...a]}}const se={cupo_10:{usados:1,porcentaje_ventana:10},leads:[{lead_id:"mock-1",nombre:"Ana Rodríguez",estado:"ENTREGADO",ruta:"ASESOR",afiliado:!0,semaforo:"VERDE",prioridad:.91,resumen:"Afiliada cat. A · presupuesto $166.8M · intención alta",actualizado_en:new Date().toISOString()},{lead_id:"mock-2",nombre:"Carlos Martínez",estado:"PERFILANDO",ruta:"ASESOR",afiliado:!1,semaforo:"AMBAR",prioridad:.78,resumen:"No afiliado · presupuesto $210M · intención media",actualizado_en:new Date().toISOString()},{lead_id:"mock-3",nombre:"Luisa Gómez",estado:"CALIFICADO",ruta:"NUTRICION",afiliado:!0,semaforo:"VERDE",prioridad:.85,resumen:"Afiliada cat. B · presupuesto $145M · requiere plan cuota inicial",actualizado_en:new Date().toISOString()}]},ce={"mock-1":{ficha_id:"f1",lead_id:"mock-1",generada_en:new Date().toISOString(),confianza_perfil:.94,banda_advertencia:null,identificacion:{nombre:"Ana Rodríguez",afiliada:!0,categoria:"A",telefono:"+57 300 123 4567"},capacidad:{presupuesto_max:1668e5,credito_max:1143e5,subsidio_aplicable:525e5,recursos_propios:12e6,ratio:.28,confianza:.94,desglose:[{concepto:"Subsidio Mi Casa Ya / Caja",monto:525e5,regla:"Afiliado Cat A",fuente:"VERIFICADO_BASE"},{concepto:"Preaprobado Bancolombia",monto:1143e5,regla:"Capacidad de endeudamiento 30%",fuente:"INFERIDO"},{concepto:"Ahorro Declarado",monto:12e6,regla:"Declarado en chat",fuente:"DECLARADO"}]},perfil:{ingreso:{valor:26e5,fuente:"VERIFICADO_BASE",confianza:.95,requiere_confirmacion:!1,actualizado_en:new Date().toISOString()}},intencion:{nivel:"ALTA",confianza:"ALTA",senales:["Busca comprar antes de 6 meses","Tiene ahorro inicial preparado","Responde inmediatamente a Vivi"]},recomendaciones:[{proyecto_id:"mongui",nombre:"Monguí",zona:"Ciudadela Maiporé - Soacha",precio_desde:15647e4,razon:"Tu presupuesto cubre el 100% de la cuota inicial",vecinos:622,tasa_desistimiento:.12,brochure_url:"https://heyzine.com/flip-book/866af8f6a6.html",recorrido_360_url:"https://storage.net-fs.com/hosting/7532170/19/"}],beneficios:["Subsidio de vivienda Colsubsidio hasta $52,5M","Tasa preferencial crédito hipotecario"],argumentos_venta:["Cuota estimada mensual ($650k) es MENOR al arriendo promedio de la zona ($850k)","Entrega inmediata en 2026"],alerta_desistimiento:{activa:!1,tasa_vecinos:.12,detalle:null},consume_cupo_10:!1},"mock-2":{ficha_id:"f2",lead_id:"mock-2",generada_en:new Date().toISOString(),confianza_perfil:.82,banda_advertencia:"No afiliado a Colsubsidio — consume cupo del 10% regulatorio",identificacion:{nombre:"Carlos Martínez",afiliada:!1,categoria:"N/A",telefono:"+57 311 987 6543"},capacidad:{presupuesto_max:21e7,credito_max:18e7,subsidio_aplicable:0,recursos_propios:3e7,ratio:.32,confianza:.82,desglose:[{concepto:"Crédito solicitado",monto:18e7,regla:"Ingresos independientes",fuente:"DECLARADO"},{concepto:"Cuota Inicial Propia",monto:3e7,regla:"Recursos declarados",fuente:"DECLARADO"}]},perfil:{},intencion:{nivel:"MEDIA",confianza:"MEDIA",senales:["Interesado en proyectos VIS y no VIS","Evaluando opciones de crédito"]},recomendaciones:[],beneficios:["Opción de crédito con aliados de la caja"],argumentos_venta:["Proyecto Versalles cuenta con certificación EDGE para ahorro energético"],alerta_desistimiento:{activa:!1,tasa_vecinos:.08,detalle:null},consume_cupo_10:!0},"mock-3":{ficha_id:"f3",lead_id:"mock-3",generada_en:new Date().toISOString(),confianza_perfil:.88,banda_advertencia:null,identificacion:{nombre:"Luisa Gómez",afiliada:!0,categoria:"B",telefono:"+57 320 456 7890"},capacidad:{presupuesto_max:145e6,credito_max:95e6,subsidio_aplicable:4e7,recursos_propios:1e7,ratio:.25,confianza:.88,desglose:[{concepto:"Subsidio Colsubsidio Cat B",monto:4e7,regla:"Afiliado Cat B",fuente:"VERIFICADO_BASE"},{concepto:"Preaprobado Hipotecario",monto:95e6,regla:"Capacidad 30%",fuente:"INFERIDO"},{concepto:"Ahorro Declarado",monto:1e7,regla:"Recursos declarados",fuente:"DECLARADO"}]},perfil:{},intencion:{nivel:"ALTA",confianza:"MEDIA",senales:["Busca vivienda en zona Soacha","Evaluando plan de ahorro de cuota inicial"]},recomendaciones:[],beneficios:["Subsidio de vivienda Colsubsidio Cat B"],argumentos_venta:["Proyecto La Macarena ofrece el valor por m² más competitivo de la zona"],alerta_desistimiento:{activa:!1,tasa_vecinos:.08,detalle:null},consume_cupo_10:!1}},V={mongui:{proyecto_id:"mongui",nombre:"Monguí",muestras:312,afiliacion:{afiliados:198,no_afiliados:114},categoria:{A:110,B:68,C:20,SIN_DATO:114},rango_edad:{"20-35":165,"36-45":82,"46-55":40,"55+":12,SIN_DATO:13},tasa_desistimiento:.11,actualizado_en:new Date().toISOString()},macarena:{proyecto_id:"macarena",nombre:"La Macarena",muestras:185,afiliacion:{afiliados:140,no_afiliados:45},categoria:{A:85,B:42,C:13,SIN_DATO:45},rango_edad:{"20-35":100,"36-45":50,"46-55":20,"55+":8,SIN_DATO:7},tasa_desistimiento:.08,actualizado_en:new Date().toISOString()},versalles:{proyecto_id:"versalles",nombre:"Versalles",muestras:142,afiliacion:{afiliados:85,no_afiliados:57},categoria:{A:40,B:32,C:13,SIN_DATO:57},rango_edad:{"20-35":70,"36-45":42,"46-55":15,"55+":8,SIN_DATO:7},tasa_desistimiento:.15,actualizado_en:new Date().toISOString()}};function re(e,a){N=!0,setTimeout(()=>{$++;const t=g[e]??g["mock-1"],o=/proyecto|comprar|vivienda|casa|apto/i.test(a),n={mensaje_id:`m${$}`,autor:"VIVI",tipo_contenido:o?"TARJETAS_PROYECTOS":"TEXTO",texto:o?"Basándome en tu perfil, estos proyectos te pueden interesar:":"¡Qué bueno saberlo! ¿Te interesaría que exploremos opciones de vivienda según tu presupuesto?",creado_en:new Date().toISOString(),adjunto:o?{recomendaciones:[{proyecto_id:"mongui",nombre:"Monguí",zona:"Ciudadela Maiporé - Soacha",precio_desde:15647e4,razon:"Tu presupuesto cubre el 100% de la cuota inicial",vecinos:622,tasa_desistimiento:.12,brochure_url:"https://heyzine.com/flip-book/866af8f6a6.html",recorrido_360_url:"https://storage.net-fs.com/hosting/7532170/19/"},{proyecto_id:"macarena",nombre:"La Macarena",zona:"Ciudadela Maiporé - Soacha",precio_desde:12834e4,razon:"El más económico de la zona, ideal para tu ingreso",vecinos:374,tasa_desistimiento:.08,brochure_url:"https://heyzine.com/flip-book/b168b2f5ba.html",recorrido_360_url:""},{proyecto_id:"versalles",nombre:"Versalles",zona:"Ciudadela Maiporé - Soacha",precio_desde:1952e5,razon:"Certificación EDGE, ahorro en servicios",vecinos:174,tasa_desistimiento:.15,brochure_url:"https://heyzine.com/flip-book/be784b0d5c.html",recorrido_360_url:"https://shape.com.co/360/COLSUBSIDIO-Versalles_APTOA"}]}:null};t.push(n),N=!1},2e3)}function de(){const e=window.fetch;window.fetch=async(a,t)=>{const o=a.toString(),n=((t==null?void 0:t.method)??"GET").toUpperCase(),s=(i,c=200)=>new Response(JSON.stringify(i),{status:c,headers:{"Content-Type":"application/json"}});if(o.includes("/conversacion")&&n==="GET"){const i=o.match(/\/leads\/([^/]+)\/conversacion/),c=i?i[1]:"mock-1";return s(ie(c))}if(o.includes("/mensajes")&&n==="POST"){const i=o.match(/\/leads\/([^/]+)\/mensajes/),c=i?i[1]:"mock-1";$++;let r={};try{r=JSON.parse((t==null?void 0:t.body)||"{}")}catch{r={}}const l=r.tipo==="AUDIO";return(g[c]??g["mock-1"]).push({mensaje_id:`m${$}`,autor:"LEAD",tipo_contenido:"TEXTO",texto:l?"🎙️ [Nota de voz]":r.texto||"",creado_en:new Date().toISOString(),adjunto:l?{audio_original:!0}:null}),re(c,l?"audio":r.texto||""),s({mensaje_id:`m${$}`,turno_en_proceso:!0},201)}if(o.includes("/ficha")&&n==="GET"){const i=o.match(/\/leads\/([^/]+)\/ficha/),c=i?i[1]:"mock-1",r=ce[c];return r?s(r):s({error:{codigo:"FICHA_NO_DISPONIBLE",mensaje:"Ficha aún no disponible"}},404)}if(o.includes("/gerencia/buyer-persona")&&n==="GET"){const c=new URL(o,"http://localhost").searchParams.get("proyecto_id")??"mongui",r=V[c]??V.mongui;return s(r)}if(o.includes("/demo/tiempo")&&n==="POST")return s({fecha_simulada:new Date().toISOString(),hitos_disparados:2});if(o.includes("/demo/reiniciar")&&n==="POST")return Object.values(g).forEach(i=>{i.length=1}),s({reiniciado:!0,fecha_simulada:new Date().toISOString()});if(o.endsWith("/api/leads")&&n==="POST"){let i={};try{i=JSON.parse((t==null?void 0:t.body)||"{}")}catch{i={}}const r={ana:"mock-1",carlos:"mock-2",luisa:"mock-3"}[i.precargado_id??"ana"]??"mock-1";return s({lead_id:r,estado:"PERFILANDO",afiliado_detectado:r!=="mock-2"},201)}return o.endsWith("/api/leads")&&n==="GET"?s(se):e(a,t)},console.info("[mock] servidor simulado completo activo (Chat multi-lead + Ficha + Cola + Gerencia + Demo)")}const P={leadActivo:null,conversacion:null,cola:null,tabActiva:"cola"},ee=[];function O(){return P}function y(e){Object.assign(P,e),ee.forEach(a=>a(P))}function ae(e){ee.push(e)}const le="/api";class E extends Error{constructor(a,t,o){super(t),this.codigo=a,this.estadoHTTP=o,this.name="ErrorAPI"}}async function b(e,a){var o,n;const t=await fetch(`${le}${e}`,{headers:{"Content-Type":"application/json"},...a});if(!t.ok){const s=await t.json().catch(()=>null);throw new E(((o=s==null?void 0:s.error)==null?void 0:o.codigo)??"ERROR_INTERNO",((n=s==null?void 0:s.error)==null?void 0:n.mensaje)??`HTTP ${t.status}`,t.status)}return await t.json()}const v={crearLead:e=>b("/leads",{method:"POST",body:JSON.stringify({precargado_id:e,fuente:"DEMO"})}),enviarTexto:(e,a)=>b(`/leads/${e}/mensajes`,{method:"POST",body:JSON.stringify({tipo:"TEXTO",texto:a})}),enviarAudio:(e,a,t,o)=>b(`/leads/${e}/mensajes`,{method:"POST",body:JSON.stringify({tipo:"AUDIO",audio_base64:a,mime:t,duracion_s:o})}),conversacion:e=>b(`/leads/${e}/conversacion`),cola:()=>b("/leads"),ficha:e=>b(`/leads/${e}/ficha`),buyerPersona:e=>b(`/gerencia/buyer-persona${e?`?proyecto_id=${e}`:""}`),avanzarTiempo:e=>b("/demo/tiempo",{method:"POST",body:JSON.stringify({avanzar_hasta:e})}),reiniciar:()=>b("/demo/reiniciar",{method:"POST"})};function ue(e,a){e.innerHTML=a.map(me).join("")}function me(e){var i,c;if(e.tipo_contenido==="SISTEMA")return`<div class="pildora-sistema">${C(e.texto)}</div>`;const a=e.autor==="LEAD"?"derecha":"izquierda",t=new Date(e.creado_en).toLocaleTimeString("es-CO",{hour:"2-digit",minute:"2-digit"}),o=e.autor==="LEAD"?'<span class="chulos">✓✓</span>':"",n=(i=e.adjunto)!=null&&i.audio_original?'<span class="icono-audio" aria-label="nota de voz">🎙️</span>':"",s=(c=e.adjunto)!=null&&c.recomendaciones?pe(e.adjunto.recomendaciones):"";return`
    <div class="burbuja ${a}">
      ${n}<p>${C(e.texto)}</p>
      <span class="hora">${t}${o}</span>
    </div>${s}`}function pe(e){return`<div class="carrusel" role="list">${e.slice(0,3).map(a=>`
    <article class="tarjeta-proyecto" role="listitem">
      <header class="franja-azul">${C(a.nombre)}</header>
      <p class="zona">${C(a.zona)}</p>
      <p class="precio">Desde $${(a.precio_desde/1e6).toFixed(0)}M</p>
      <p class="razon">${C(a.razon)}</p>
      <p class="evidencia">${a.vecinos} personas con tu perfil compraron aquí ·
         ${(a.tasa_desistimiento*100).toFixed(0)}% desistió</p>
      <a class="btn-primario" href="${encodeURI(a.brochure_url)}" target="_blank" rel="noopener">Ver brochure</a>
      <a class="btn-secundario" href="${encodeURI(a.recorrido_360_url)}" target="_blank" rel="noopener">Recorrido 360°</a>
    </article>`).join("")}</div>`}function q(e,a){e.classList.toggle("visible",a)}function fe(e){const a=document.querySelector(".nombre-vivi");a&&(a.textContent=`Vivi — ${e}`)}function ve(e){e.innerHTML=`
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
  `}function C(e){const a=document.createElement("div");return a.textContent=e,a.innerHTML}const be=60,he=1500,ge=300;let S=null,j=!1,h=null,I=[],z=null,A=0;function ye(e){ve(e);const a=document.getElementById("contenedor-mensajes"),t=document.getElementById("mensajes-scroll"),o=document.getElementById("indicador-escribiendo"),n=document.getElementById("input-mensaje"),s=document.getElementById("btn-enviar"),i=document.getElementById("btn-mic"),c=document.getElementById("mic-grabando"),r=document.getElementById("contador-mic"),l=document.getElementById("btn-detener-mic"),d=document.getElementById("btn-nuevos");n.addEventListener("keydown",m=>{m.key==="Enter"&&!m.shiftKey&&(m.preventDefault(),H(n))}),s.addEventListener("click",()=>H(n)),i.addEventListener("click",()=>_e(i,c,r,n,s)),l.addEventListener("click",()=>te(i,c,n,s)),t.addEventListener("scroll",()=>{j=t.scrollHeight-t.scrollTop-t.clientHeight>80,j||d.classList.remove("visible")}),d.addEventListener("click",()=>{t.scrollTop=t.scrollHeight,d.classList.remove("visible"),j=!1}),setInterval(()=>w(a,t,o,d),he),ae(()=>w(a,t,o,d)),w(a,t,o,d)}async function w(e,a,t,o){var i;const n=O();if(!n.leadActivo)return;const s=(i=n.cola)==null?void 0:i.leads.find(c=>c.lead_id===n.leadActivo);s&&fe(s.nombre);try{const c=await v.conversacion(n.leadActivo);y({conversacion:c}),ue(e,c.mensajes),c.turno_en_proceso?S||(S=setTimeout(()=>q(t,!0),ge)):(S&&(clearTimeout(S),S=null),q(t,!1)),j?o.classList.add("visible"):a.scrollTop=a.scrollHeight}catch(c){c instanceof E&&c.estadoHTTP>=500&&L("Error de conexión. Reintentando…")}}async function H(e){const a=e.value.trim();if(!a)return;const t=O();if(t.leadActivo){e.value="";try{await v.enviarTexto(t.leadActivo,a)}catch(o){e.value=a,o instanceof E&&L(o.message)}}}async function _e(e,a,t,o,n){try{const s=await navigator.mediaDevices.getUserMedia({audio:!0});h=new MediaRecorder(s),I=[],A=0,h.ondataavailable=i=>{I.push(i.data)},h.onstop=()=>{s.getTracks().forEach(i=>i.stop()),Ae(e,a,o,n)},h.start(),e.style.display="none",o.style.display="none",n.style.display="none",a.classList.add("activo"),t.textContent="0:00",z=setInterval(()=>{A++;const i=Math.floor(A/60),c=A%60;t.textContent=`${i}:${c.toString().padStart(2,"0")}`,A>=be&&te(e,a,o,n)},1e3)}catch{L("No se pudo acceder al micrófono.")}}function te(e,a,t,o){z&&(clearInterval(z),z=null),h&&h.state!=="inactive"&&h.stop(),a.classList.remove("activo"),e.style.display="",t.style.display="",o.style.display=""}async function Ae(e,a,t,o){if(I.length===0)return;const n=O();if(!n.leadActivo)return;const s=new Blob(I,{type:I[0].type}),i=new FileReader;i.onloadend=async()=>{const c=i.result.split(",")[1],r=s.type,l=A;try{await v.enviarAudio(n.leadActivo,c,r,l)}catch(d){d instanceof E&&(d.codigo==="AUDIO_INVALIDO"?L("No te escuché bien, ¿me lo repites o me lo escribes?"):L(d.message))}},i.readAsDataURL(s),a.classList.remove("activo"),e.style.display="",t.style.display="",o.style.display=""}function L(e){const a=document.querySelector(".toast-error");a&&a.remove();const t=document.createElement("div");t.className="toast-error",t.textContent=e,document.body.appendChild(t),setTimeout(()=>t.remove(),4e3)}const G={VERDE:{icono:"🟢",etiqueta:"Alta Prio / Apto",clase:"badge-semaforo-verde"},AMBAR:{icono:"🟡",etiqueta:"En Validación",clase:"badge-semaforo-ambar"},GRIS:{icono:"⚪",etiqueta:"Sin Datos",clase:"badge-semaforo-gris"}};function Ee(e,a,t,o){const n=a.cupo_10.usados,s=a.cupo_10.porcentaje_ventana,i=Math.min(100,Math.round(n/(s||1)*100)),c=i>=80;e.innerHTML=`
    <div class="cola-container">
      <header class="cola-header-card">
        <div class="cola-header-titulo">
          <div class="header-icon-badge">🎯</div>
          <div>
            <h2>Cola Priorizada de Leads</h2>
            <p class="cola-header-sub">Asignación inteligente en tiempo real según modelo de scoring</p>
          </div>
          <span class="badge-total-leads">${a.leads.length} Activos</span>
        </div>

        <div class="cupo-regulado-card ${c?"alerta":""}" 
             title="${c?"¡Alerta! Cupo del 10% de no afiliados casi agotado (≥80%)":"Uso del cupo regulatorio del 10% para no afiliados"}">
          <div class="cupo-info-top">
            <span class="cupo-label">⚡ Cupo Regulado 10%</span>
            <span class="cupo-metrica">${n} / ${s} (${i}%)</span>
          </div>
          <div class="cupo-track-custom">
            <div class="cupo-fill-custom" style="transform: scaleX(${i/100}); transform-origin: left;"></div>
          </div>
        </div>
      </header>

      <div class="tabla-cola-wrapper">
        <table class="tabla-cola-leads">
          <thead>
            <tr>
              <th scope="col" class="th-pos"># Pos</th>
              <th scope="col" class="th-estado">Estado</th>
              <th scope="col" class="th-nombre">Lead / Prospecto</th>
              <th scope="col" class="th-ruta">Canal / Ruta</th>
              <th scope="col" class="th-prio">Score Prio</th>
              <th scope="col" class="th-resumen">Resumen de Interacción</th>
              <th scope="col" class="th-acciones">Acciones</th>
            </tr>
          </thead>
          <tbody>
            ${a.leads.map((r,l)=>Se(r,l+1)).join("")}
          </tbody>
        </table>
      </div>
    </div>
  `,e.querySelectorAll("[data-lead-id]").forEach(r=>{const l=r.dataset.leadId,d=r.querySelector("[data-btn-chat]"),m=r.querySelector("[data-btn-ficha]");d&&d.addEventListener("click",u=>{u.stopPropagation(),o(l)}),m&&m.addEventListener("click",u=>{u.stopPropagation(),t(l)}),r.addEventListener("click",()=>t(l)),r.addEventListener("keydown",u=>{(u.key==="Enter"||u.key===" ")&&(u.preventDefault(),t(l))})})}function Se(e,a){const t=G[e.semaforo]??G.GRIS,o=e.nombre.trim().charAt(0).toUpperCase()||"L",n=e.afiliado!==void 0?`<span class="badge-afiliado-pill ${e.afiliado?"es-afiliado":"no-afiliado"}">${e.afiliado?"Afiliado":"No afiliado"}</span>`:"";return`
    <tr class="fila-lead-tabla" data-lead-id="${e.lead_id}" tabindex="0" role="row">
      <td class="td-pos">
        <span class="pos-badge pos-${a}">#${a}</span>
      </td>
      <td class="td-estado">
        <span class="badge-semaforo ${t.clase}" title="Estado: ${t.etiqueta}">
          <span class="dot-icon">${t.icono}</span>
          <span class="txt-estado">${t.etiqueta}</span>
        </span>
      </td>
      <td class="td-nombre">
        <div class="lead-avatar-wrap">
          <div class="lead-avatar">${o}</div>
          <div class="lead-meta">
            <span class="lead-nombre-txt">${_(e.nombre)} ${n}</span>
            <span class="lead-id-sub">ID: ${_(e.lead_id)}</span>
          </div>
        </div>
      </td>
      <td class="td-ruta">
        <span class="badge-ruta" title="Origen del lead">${_(e.ruta)}</span>
      </td>
      <td class="td-prio">
        <span class="badge-prio-score cifra" title="Puntuación de Priorización Vivi">
          ${e.prioridad.toFixed(1)} <small>pts</small>
        </span>
      </td>
      <td class="td-resumen">
        <p class="resumen-corte">${_(e.resumen)}</p>
      </td>
      <td class="td-acciones">
        <div class="acciones-btn-group">
          <button class="btn-tabla-chat" data-btn-chat="true" title="Abrir chat de WhatsApp con ${_(e.nombre)}" type="button">
            💬 Chat
          </button>
          <button class="btn-tabla-ficha" data-btn-ficha="true" title="Ver Ficha Comercial de ${_(e.nombre)}" type="button">
            👁️ Ficha
          </button>
        </div>
      </td>
    </tr>
  `}function _(e){const a=document.createElement("div");return a.textContent=e,a.innerHTML}const U={VERIFICADO_BASE:{txt:"VERIFICADO",clase:"sello-verificado"},DECLARADO:{txt:"DECLARADO",clase:"sello-declarado"},INFERIDO:{txt:"INFERIDO",clase:"sello-inferido"}},X={ALTA:"zona-verde",MEDIA:"zona-ambar",BAJA:"zona-gris"};function $e(e){const a=U[e]??U.DECLARADO;return`<span class="sello-fuente ${a.clase}" title="Fuente: ${a.txt}">${a.txt}</span>`}function J(e,a,t="Lead"){if(!a){Ie(e,t);return}const o=a.banda_advertencia?`<div class="banda-advertencia" role="alert">
         <span class="banda-advertencia-etiqueta">Alerta</span>
         <span>${p(a.banda_advertencia)}</span>
       </div>`:"",n=Math.round(a.confianza_perfil*100),s=a.identificacion,i=a.capacidad,c=a.intencion,r=X[c.nivel]??X.BAJA;e.innerHTML=`
    <div class="hoja-informe informe-ficha">
      ${o}

      <!-- Masthead: identidad a la izquierda, medidor de confianza a la derecha -->
      <header class="masthead-ficha">
        <div class="masthead-identidad">
          <span class="masthead-kicker">Ficha comercial</span>
          <h3 class="masthead-nombre">${p(s.nombre||t)}</h3>
          <dl class="registro-identidad">
            <div class="registro-identidad-item">
              <dt>Afiliación</dt>
              <dd>${s.afiliada?`Afiliado · Categoría ${p(s.categoria)}`:"No afiliado"}</dd>
            </div>
            <div class="registro-identidad-item">
              <dt>Teléfono</dt>
              <dd class="cifra">${p(s.telefono||"No registrado")}</dd>
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

      <!-- Cartera de capacidad financiera -->
      <section class="seccion-informe seccion-capacidad">
        <div class="seccion-header-title">
          <span class="sec-icon">💳</span>
          <h4 class="seccion-titulo">Capacidad financiera</h4>
        </div>
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
              ${i.desglose.map(Ce).join("")}
            </tbody>
            <tfoot>
              <tr class="ledger-total">
                <td colspan="2">Presupuesto máximo</td>
                <td class="ledger-num cifra">$${B(i.presupuesto_max)}</td>
                <td></td>
              </tr>
              <tr class="ledger-ratio">
                <td colspan="2">Ratio de endeudamiento</td>
                <td class="ledger-num cifra">${(i.ratio*100).toFixed(0)}%</td>
                <td class="cifra">conf. ${(i.confianza*100).toFixed(0)}%</td>
              </tr>
            </tfoot>
          </table>
        </div>
      </section>

      <div class="informe-columnas">
        <!-- Zona de intención: banda con el veredicto de nivel -->
        <section class="seccion-informe zona-veredicto ${r}">
          <div class="seccion-header-title">
            <span class="sec-icon">🎯</span>
            <h4 class="seccion-titulo">Intención de compra</h4>
          </div>
          <p class="veredicto-nivel">Nivel <span class="cifra">${p(c.nivel)}</span></p>
          <p class="veredicto-confianza">Confianza: ${p(c.confianza)}</p>
          <ul class="lista-observaciones">
            ${c.senales.map(d=>`<li>${p(d)}</li>`).join("")}
          </ul>
        </section>

        <!-- Argumentos y beneficios -->
        <section class="seccion-informe">
          <div class="seccion-header-title">
            <span class="sec-icon">💡</span>
            <h4 class="seccion-titulo">Argumentos de venta</h4>
          </div>
          <ul class="lista-observaciones">
            ${a.argumentos_venta.map(d=>`<li>${p(d)}</li>`).join("")}
          </ul>
          <div class="seccion-header-title" style="margin-top:0.85rem">
            <span class="sec-icon">✨</span>
            <h5 class="subseccion-titulo" style="margin:0">Beneficios Colsubsidio</h5>
          </div>
          <ul class="lista-observaciones">
            ${a.beneficios.map(d=>`<li>${p(d)}</li>`).join("")}
          </ul>
        </section>
      </div>

      <!-- Banner de instrucción única -->
      <div class="banner-siguiente-paso">
        <span><strong class="banner-etiqueta">Siguiente paso</strong> Agendar visita a sala de ventas</span>
        <button class="btn-copiar-resumen" id="btn-copiar-resumen" type="button">
          📋 Copiar resumen
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
  `;const l=e.querySelector("#btn-copiar-resumen");l&&l.addEventListener("click",()=>{const d=`RESUMEN FICHA - ${s.nombre}
Presupuesto: $${B(i.presupuesto_max)} | Afiliada: ${s.afiliada?"Sí":"No"}
Intención: ${c.nivel}
Siguiente paso: Agendar visita sala de ventas.`;navigator.clipboard.writeText(d).then(()=>{l.textContent="✓ Copiado",setTimeout(()=>{l.textContent="📋 Copiar resumen"},2e3)})})}function Ce(e){return`
    <tr>
      <td>${p(e.concepto)}</td>
      <td class="ledger-regla">${p(e.regla)}</td>
      <td class="ledger-num cifra">$${B(e.monto)}</td>
      <td>${$e(e.fuente)}</td>
    </tr>
  `}function B(e){return`${(e/1e6).toFixed(1)}M`}function Ie(e,a){e.innerHTML=`
    <div class="hoja-informe informe-vacio">
      <span class="masthead-kicker">Ficha comercial</span>
      <h3>Aún sin generar</h3>
      <p>
        La ficha comercial completa de <strong>${p(a)}</strong> se generará automáticamente cuando Vivi complete la calificación.
      </p>
      <details class="timeline-plegable">
        <summary>Ver avance actual</summary>
        <p class="timeline-texto">
          Estado: En proceso de perfilamiento en tiempo real desde el chat.
        </p>
      </details>
    </div>
  `}function p(e){const a=document.createElement("div");return a.textContent=e,a.innerHTML}const Y=[{id:"mongui",nombre:"Monguí"},{id:"macarena",nombre:"La Macarena"},{id:"versalles",nombre:"Versalles"},{id:"todos",nombre:"Todos los proyectos"}];function W(e,a,t="mongui",o){var m;const n=a?a.muestras:312,s=a?(a.tasa_desistimiento*100).toFixed(0):"11",i=((m=Y.find(u=>u.id===t))==null?void 0:m.nombre)??"Proyecto",c=(a==null?void 0:a.afiliacion)??{afiliados:180,no_afiliados:132},r=(a==null?void 0:a.categoria)??{"Cat A":90,"Cat B":65,"Cat C":25,"No Afiliado":132},l=(a==null?void 0:a.rango_edad)??{"18-25":20,"26-35":150,"36-45":90,"46+":52};e.innerHTML=`
    <div class="hoja-informe informe-gerencia">
      <header class="masthead-gerencia">
        <div class="masthead-gerencia-titulo">
          <span class="masthead-kicker">Buyer persona vivo</span>
          <h2 class="masthead-nombre">${F(i)}</h2>
        </div>
        <div class="selector-proyecto-wrap">
          <label for="select-proyecto-gerencia">Proyecto</label>
          <select id="select-proyecto-gerencia" class="select-proyecto">
            ${Y.map(u=>`
              <option value="${u.id}" ${u.id===t?"selected":""}>
                ${F(u.nombre)}
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
          <span class="franja-metrica-cifra cifra">${s}%</span>
        </div>
      </div>

      <!-- Columnas de distribución (estilo junta directiva) -->
      <section class="panel-graficos">
        <div class="columna-grafico">
          <h4 class="seccion-titulo">Distribución por afiliación</h4>
          <div class="barras-lista">
            ${M(c)}
          </div>
        </div>
        <div class="columna-grafico">
          <h4 class="seccion-titulo">Categoría de afiliación</h4>
          <div class="barras-lista">
            ${M(r)}
          </div>
        </div>
        <div class="columna-grafico">
          <h4 class="seccion-titulo">Rango de edad</h4>
          <div class="barras-lista">
            ${M(l)}
          </div>
        </div>
      </section>
    </div>
  `;const d=e.querySelector("#select-proyecto-gerencia");d&&d.addEventListener("change",()=>{o(d.value)})}const Le={afiliados:"Afiliados",no_afiliados:"No afiliados",SIN_DATO:"Sin dato",A:"Categoría A",B:"Categoría B",C:"Categoría C"};function Oe(e){return Le[e]??e}function M(e){const a=Math.max(...Object.values(e),1);return Object.entries(e).map(([t,o])=>{const n=(o/a*100).toFixed(0);return`
      <div class="barra-item">
        <span class="barra-label">${F(Oe(t))}</span>
        <div class="barra-track">
          <div class="barra-fill" style="transform: scaleX(${Number(n)/100})"></div>
        </div>
        <span class="barra-val cifra">${o}</span>
      </div>
    `}).join("")}function F(e){const a=document.createElement("div");return a.textContent=e,a.innerHTML}function Te(e,a,t){e.innerHTML=`
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
  `;const o=e.querySelector("#btn-avanzar-tiempo");o&&o.addEventListener("click",a);const n=e.querySelector("#btn-reiniciar-demo");n&&n.addEventListener("click",t)}function De(e){return new Promise(a=>{const t=document.createElement("div");t.className="modal-overlay",t.setAttribute("role","dialog"),t.setAttribute("aria-modal","true");const o=e.icono??"❓",n=e.textoConfirmar??"Aceptar",s=e.textoCancelar??"Cancelar",i="btn-modal-primario";t.innerHTML=`
      <div class="modal-card">
        <div class="modal-header">
          <span class="modal-icon">${o}</span>
          <h3 class="modal-titulo">${f(e.titulo)}</h3>
        </div>
        <div class="modal-body">
          <p class="modal-mensaje">${f(e.mensaje)}</p>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn-modal-cancelar" data-modal-cancelar>${f(s)}</button>
          <button type="button" class="${i}" data-modal-confirmar>${f(n)}</button>
        </div>
      </div>
    `,document.body.appendChild(t),requestAnimationFrame(()=>t.classList.add("activo"));const c=m=>{t.classList.remove("activo"),setTimeout(()=>{t.parentNode&&t.parentNode.removeChild(t),a(m)},200)},r=t.querySelector("[data-modal-confirmar]"),l=t.querySelector("[data-modal-cancelar]");r==null||r.focus(),r==null||r.addEventListener("click",()=>c(!0)),l==null||l.addEventListener("click",()=>c(!1)),t.addEventListener("click",m=>{m.target===t&&c(!1)});const d=m=>{m.key==="Escape"&&(document.removeEventListener("keydown",d),c(!1))};document.addEventListener("keydown",d)})}function je(e){return new Promise(a=>{const t=document.createElement("div");t.className="modal-overlay",t.setAttribute("role","dialog"),t.setAttribute("aria-modal","true");const o=e.icono??"✏️",n=e.textoConfirmar??"Aceptar",s=e.textoCancelar??"Cancelar",i=e.valorDefecto??"";t.innerHTML=`
      <div class="modal-card">
        <div class="modal-header">
          <span class="modal-icon">${o}</span>
          <h3 class="modal-titulo">${f(e.titulo)}</h3>
        </div>
        <div class="modal-body">
          <p class="modal-mensaje">${f(e.mensaje)}</p>
          <input type="text" class="input-modal-texto" data-modal-input value="${f(i)}" />
        </div>
        <div class="modal-footer">
          <button type="button" class="btn-modal-cancelar" data-modal-cancelar>${f(s)}</button>
          <button type="button" class="btn-modal-primario" data-modal-confirmar>${f(n)}</button>
        </div>
      </div>
    `,document.body.appendChild(t),requestAnimationFrame(()=>t.classList.add("activo"));const c=t.querySelector("[data-modal-input]"),r=u=>{t.classList.remove("activo"),setTimeout(()=>{t.parentNode&&t.parentNode.removeChild(t),a(u)},200)};c&&(c.focus(),c.select(),c.addEventListener("keydown",u=>{u.key==="Enter"&&(u.preventDefault(),r(c.value))}));const l=t.querySelector("[data-modal-confirmar]"),d=t.querySelector("[data-modal-cancelar]");l==null||l.addEventListener("click",()=>r((c==null?void 0:c.value)??"")),d==null||d.addEventListener("click",()=>r(null)),t.addEventListener("click",u=>{u.target===t&&r(null)});const m=u=>{u.key==="Escape"&&(document.removeEventListener("keydown",m),r(null))};document.addEventListener("keydown",m)})}function T(e){return new Promise(a=>{const t=document.createElement("div");t.className="modal-overlay",t.setAttribute("role","dialog"),t.setAttribute("aria-modal","true");const o=e.icono??"ℹ️",n=e.textoConfirmar??"Entendido";t.innerHTML=`
      <div class="modal-card">
        <div class="modal-header">
          <span class="modal-icon">${o}</span>
          <h3 class="modal-titulo">${f(e.titulo)}</h3>
        </div>
        <div class="modal-body">
          <p class="modal-mensaje">${f(e.mensaje)}</p>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn-modal-primario" data-modal-confirmar>${f(n)}</button>
        </div>
      </div>
    `,document.body.appendChild(t),requestAnimationFrame(()=>t.classList.add("activo"));const s=()=>{t.classList.remove("activo"),setTimeout(()=>{t.parentNode&&t.parentNode.removeChild(t),a()},200)},i=t.querySelector("[data-modal-confirmar]");i==null||i.focus(),i==null||i.addEventListener("click",()=>s()),t.addEventListener("click",r=>{r.target===t&&s()});const c=r=>{(r.key==="Escape"||r.key==="Enter")&&(document.removeEventListener("keydown",c),s())};document.addEventListener("keydown",c)})}function f(e){const a=document.createElement("div");return a.textContent=e,a.innerHTML}const ze=5e3;let k="mongui",D=null,x=null;function ke(e,a,t){if(t){const o=t.querySelectorAll("button[data-tab]");o.forEach(n=>{n.addEventListener("click",()=>{const s=n.dataset.tab;xe(s,o)})})}a&&Ne(a),ae(()=>Z(e)),setInterval(()=>R(),ze),R(),Z(e)}async function R(){try{const e=await v.cola();y({cola:e})}catch(e){console.warn("[dashboard] Error cargando cola:",e)}}function xe(e,a){y({tabActiva:e}),a.forEach(t=>{const o=t.dataset.tab===e;t.setAttribute("aria-selected",o?"true":"false")})}let K=null,Q=null;function Z(e){const a=O(),t=a.tabActiva!==K;switch(a.tabActiva){case"cola":D=null,x=null,a.cola?(t||a.cola!==Q)&&(Q=a.cola,Ee(e,a.cola,o=>Re(o),o=>we(o))):e.innerHTML='<div style="padding:1rem; color:#6B7280">Cargando cola de leads…</div>';break;case"ficha":if(!a.leadActivo){D=null,e.innerHTML='<div style="padding:1.5rem; text-align:center; color:#6B7280">Selecciona un lead de la cola para ver su ficha comercial.</div>';break}(t||a.leadActivo!==D)&&(D=a.leadActivo,Me(e,a.leadActivo));break;case"gerencia":(t||k!==x)&&(x=k,oe(e,k));break}K=a.tabActiva}function Re(e){y({leadActivo:e,tabActiva:"ficha"});const a=document.querySelector(".tabs");a&&a.querySelectorAll("button[data-tab]").forEach(o=>o.setAttribute("aria-selected",o.dataset.tab==="ficha"?"true":"false"))}function we(e){y({leadActivo:e});const a=document.getElementById("panel-chat");a&&a.scrollIntoView({behavior:"smooth"})}async function Me(e,a){var s;const o=(s=O().cola)==null?void 0:s.leads.find(i=>i.lead_id===a),n=(o==null?void 0:o.nombre)??"Lead";try{const i=await v.ficha(a);J(e,i,n)}catch(i){i instanceof E&&i.estadoHTTP===404?J(e,null,n):e.innerHTML=`<div class="banda-advertencia">⚠️ Error cargando la ficha comercial: ${i.message}</div>`}}async function oe(e,a){const t=o=>{k=o,x=o,oe(e,o)};try{const o=await v.buyerPersona(a);W(e,o,a,t)}catch{W(e,null,a,t)}}function Ne(e){Te(e,async()=>{const a=await je({icono:"⏩",titulo:"Avanzar Tiempo Simulado",mensaje:"Ingresa la fecha a la cual deseas avanzar la simulación (Formato ISO AAAA-MM-DD):",valorDefecto:"2026-08-01",textoConfirmar:"Avanzar Tiempo",textoCancelar:"Cancelar"});if(a)try{const t=await v.avanzarTiempo(a);await T({icono:"✅",titulo:"Tiempo Avanzado",mensaje:`El reloj de simulación avanzó correctamente a ${a} (${t.hitos_disparados??0} hitos disparados).`}),R()}catch(t){await T({icono:"⚠️",titulo:"Error al Avanzar Tiempo",mensaje:t.message})}},async()=>{if(await De({icono:"🔄",titulo:"¿Reiniciar Demostración?",mensaje:"Esta acción restaurará la cola de leads y los datos de prueba a su estado inicial.",textoConfirmar:"Sí, Reiniciar Demo",textoCancelar:"Cancelar"}))try{await v.reiniciar();let t=null;try{t=(await v.crearLead("ana")).lead_id}catch{t="mock-1"}y({leadActivo:t,tabActiva:"cola"}),R(),await T({icono:"✅",titulo:"Demo Reiniciado",mensaje:"El estado de la demostración se ha restaurado con éxito."})}catch(t){await T({icono:"⚠️",titulo:"Error al Reiniciar Demo",mensaje:t.message})}})}const ne=new URLSearchParams(location.search);ne.get("mock")==="1"&&de();const Pe=["ana","carlos","luisa"];function Be(){const e=ne.get("precargado");return Pe.includes(e??"")?e:"ana"}async function Fe(){try{return(await v.crearLead(Be())).lead_id}catch(e){const a=e instanceof E?`${e.codigo}: ${e.message}`:String(e);return console.error("[main] no se pudo crear el lead inicial:",a),Ve(a),null}}function Ve(e){const a=document.getElementById("panel-chat");if(!a)return;const t=document.createElement("div");t.className="banda-advertencia",t.setAttribute("role","alert"),t.textContent=`No se pudo iniciar la conversación (${e}). Recargá la página o revisá que el backend esté arriba.`,a.prepend(t)}async function qe(){const e=document.getElementById("panel-chat");if(e){const n=await Fe();n&&(y({leadActivo:n}),ye(e))}const a=document.getElementById("contenido-tab"),t=document.querySelector(".tabs"),o=document.getElementById("botonera-demo");a&&ke(a,o,t),console.info("Vivi web iniciado completo (Chat + Dashboard)")}function He(){document.querySelectorAll("[data-colapsar]").forEach(e=>{const a=document.getElementById(e.dataset.colapsar);a&&e.addEventListener("click",()=>{const t=a.classList.toggle("panel-colapsado");e.setAttribute("aria-expanded",String(!t))})})}He();qe();
