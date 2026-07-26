(function(){const a=document.createElement("link").relList;if(a&&a.supports&&a.supports("modulepreload"))return;for(const s of document.querySelectorAll('link[rel="modulepreload"]'))n(s);new MutationObserver(s=>{for(const i of s)if(i.type==="childList")for(const o of i.addedNodes)o.tagName==="LINK"&&o.rel==="modulepreload"&&n(o)}).observe(document,{childList:!0,subtree:!0});function t(s){const i={};return s.integrity&&(i.integrity=s.integrity),s.referrerPolicy&&(i.referrerPolicy=s.referrerPolicy),s.crossOrigin==="use-credentials"?i.credentials="include":s.crossOrigin==="anonymous"?i.credentials="omit":i.credentials="same-origin",i}function n(s){if(s.ep)return;s.ep=!0;const i=t(s);fetch(s.href,i)}})();const pe="/api";class E extends Error{constructor(a,t,n){super(t),this.codigo=a,this.estadoHTTP=n,this.name="ErrorAPI"}}async function b(e,a){var n,s;const t=await fetch(`${pe}${e}`,{headers:{"Content-Type":"application/json"},...a});if(!t.ok){const i=await t.json().catch(()=>null);throw new E(((n=i==null?void 0:i.error)==null?void 0:n.codigo)??"ERROR_INTERNO",((s=i==null?void 0:i.error)==null?void 0:s.mensaje)??`HTTP ${t.status}`,t.status)}return await t.json()}const g={crearLead:e=>b("/leads",{method:"POST",body:JSON.stringify({precargado_id:e,fuente:"DEMO"})}),enviarTexto:(e,a)=>b(`/leads/${e}/mensajes`,{method:"POST",body:JSON.stringify({tipo:"TEXTO",texto:a})}),enviarAudio:(e,a,t,n)=>b(`/leads/${e}/mensajes`,{method:"POST",body:JSON.stringify({tipo:"AUDIO",audio_base64:a,mime:t,duracion_s:n})}),conversacion:e=>b(`/leads/${e}/conversacion`),cola:()=>b("/leads"),ficha:e=>b(`/leads/${e}/ficha`),buyerPersona:e=>b(`/gerencia/buyer-persona${e?`?proyecto_id=${e}`:""}`),avanzarTiempo:e=>b("/demo/tiempo",{method:"POST",body:JSON.stringify({avanzar_hasta:e})}),reiniciar:()=>b("/demo/reiniciar",{method:"POST"})},q={leadActivo:null,conversacion:null,cola:null,tabActiva:"cola"},te=[];function L(){return q}function w(e){Object.assign(q,e),te.forEach(a=>a(q))}function ve(e){te.push(e)}const G=[["#f4a261","#7a3e1d"],["#e76f51","#5c2418"],["#2a9d8f","#0b3b36"],["#264653","#0f1e24"],["#e9c46a","#7a5a15"],["#8ecae6","#1d4e6b"],["#cdb4db","#5a3d6b"],["#a3b18a","#3f4a30"]],J=[["#1375c2","#9bd4ef"],["#2a9d8f","#8fd9cf"],["#e9c46a","#f4dfa0"],["#8ecae6","#c9e8f5"],["#e76f51","#f3ab92"],["#264653","#5c7d8a"]];function B(e){let a=0;for(let t=0;t<e.length;t++)a=a*31+e.charCodeAt(t)>>>0;return a}function ne(e){return`data:image/svg+xml;utf8,${encodeURIComponent(e)}`}function k(e){const[a,t]=G[B(e)%G.length],s=`
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64">
      <circle cx="32" cy="32" r="32" fill="${B(e+"x")%2===0?"#eef2fb":"#e7f0ff"}"/>
      <circle cx="32" cy="26" r="12" fill="${a}"/>
      <path d="M10 58c2-13 12-21 22-21s20 8 22 21z" fill="${a}"/>
      <path d="M10 58c2-13 12-21 22-21s20 8 22 21" fill="none" stroke="${t}" stroke-width="1.5" opacity="0.35"/>
    </svg>`.trim();return ne(s)}function se(e){const[a,t]=J[B(e)%J.length],n=B(e+"y"),s=[22,34,26,40,18].map((u,l)=>u+(n>>l*3)%10),i=100/s.length,o=s.map((u,l)=>{const d=l*i+2,v=i-4,r=60-u;return`<rect x="${d}" y="${r}" width="${v}" height="${u}" fill="${t}" rx="1.5"/>`}).join(""),c=`
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 60" preserveAspectRatio="xMidYMax slice">
      <rect width="100" height="60" fill="${a}"/>
      ${o}
      <rect y="58" width="100" height="2" fill="${t}" opacity="0.6"/>
    </svg>`.trim();return ne(c)}function me(e,a){e.innerHTML=a.map(be).join("")}function be(e){var o,c;if(e.tipo_contenido==="SISTEMA")return`<div class="pildora-sistema">${S(e.texto)}</div>`;const a=e.autor==="LEAD"?"derecha":"izquierda",t=new Date(e.creado_en).toLocaleTimeString("es-CO",{hour:"2-digit",minute:"2-digit"}),n=e.autor==="LEAD"?'<span class="chulos">✓✓</span>':"",s=(o=e.adjunto)!=null&&o.audio_original?'<span class="icono-audio" aria-label="nota de voz">🎙️</span>':"",i=(c=e.adjunto)!=null&&c.recomendaciones?fe(e.adjunto.recomendaciones):"";return`
    <div class="burbuja ${a}">
      ${s}<p>${S(e.texto)}</p>
      <span class="hora">${t}${n}</span>
    </div>${i}`}function fe(e){return`<div class="carrusel" role="list">${e.slice(0,3).map(a=>`
    <article class="tarjeta-proyecto" role="listitem">
      <div class="tarjeta-imagen" style="background-image:url('${se(a.proyecto_id)}')" aria-hidden="true"></div>
      <header class="franja-azul">${S(a.nombre)}</header>
      <p class="zona">${S(a.zona)}</p>
      <p class="precio">Desde $${(a.precio_desde/1e6).toFixed(0)}M</p>
      <p class="razon">${S(a.razon)}</p>
      <p class="evidencia">${a.vecinos} personas con tu perfil compraron aquí ·
         ${(a.tasa_desistimiento*100).toFixed(0)}% desistió</p>
      <a class="btn-primario" href="${encodeURI(a.brochure_url)}" target="_blank" rel="noopener">Ver brochure</a>
      <a class="btn-secundario" href="${encodeURI(a.recorrido_360_url)}" target="_blank" rel="noopener">Recorrido 360°</a>
    </article>`).join("")}</div>`}function K(e,a){e.classList.toggle("visible",a)}function ge(e){e.innerHTML=`
    <header class="chat-top">
      <button class="back" aria-label="Volver">‹</button>
      <img class="avatar chat-avatar" id="chat-avatar" src="" alt="" aria-hidden="true" style="display:none">
      <div class="chat-name">
        <span id="chat-nombre-lead">Selecciona un lead</span>
        <small>en línea <span class="online">●</span></small>
      </div>
      <div class="chat-tools">
        <button type="button" aria-label="Llamar">☎</button>
        <button type="button" aria-label="Videollamar">▣</button>
        <button type="button" aria-label="Más acciones">⋮</button>
      </div>
    </header>
    <div class="mensajes-scroll" id="mensajes-scroll">
      <div class="mensajes" id="contenedor-mensajes"></div>
      <button class="btn-nuevos" id="btn-nuevos" type="button">↓ nuevos mensajes</button>
    </div>
    <span class="escribiendo" id="indicador-escribiendo" style="padding:0 16px">escribiendo…</span>
    <div class="barra-entrada">
      <button class="btn-mic" id="btn-mic" aria-label="Grabar nota de voz" type="button">🎤</button>
      <div class="mic-grabando" id="mic-grabando">
        <span class="punto-rojo"></span>
        <span class="contador-mic" id="contador-mic">0:00</span>
        <button class="btn-detener-mic" id="btn-detener-mic" type="button">■</button>
      </div>
      <input type="text" class="input-mensaje" id="input-mensaje"
             placeholder="Escribí un mensaje..." autocomplete="off" />
      <button class="btn-enviar" id="btn-enviar" aria-label="Enviar mensaje" type="button">➤</button>
    </div>
  `}function ie(e,a){const t=document.getElementById("chat-nombre-lead"),n=document.getElementById("chat-avatar");t&&(t.textContent=a),n&&(n.src=k(e),n.style.display="")}function S(e){const a=document.createElement("div");return a.textContent=e,a.innerHTML}const he=60,ye=1500,$e=300;let T=null,R=!1,h=null,M=[],P=null,A=0;function Ae(e){ge(e);const a=document.getElementById("contenedor-mensajes"),t=document.getElementById("mensajes-scroll"),n=document.getElementById("indicador-escribiendo"),s=document.getElementById("input-mensaje"),i=document.getElementById("btn-enviar"),o=document.getElementById("btn-mic"),c=document.getElementById("mic-grabando"),u=document.getElementById("contador-mic"),l=document.getElementById("btn-detener-mic"),d=document.getElementById("btn-nuevos");s.addEventListener("keydown",v=>{v.key==="Enter"&&!v.shiftKey&&(v.preventDefault(),Y(s))}),i.addEventListener("click",()=>Y(s)),o.addEventListener("click",()=>Ee(o,c,u,s,i)),l.addEventListener("click",()=>oe(o,c,s,i)),t.addEventListener("scroll",()=>{R=t.scrollHeight-t.scrollTop-t.clientHeight>80,R||d.classList.remove("visible")}),d.addEventListener("click",()=>{t.scrollTop=t.scrollHeight,d.classList.remove("visible"),R=!1}),setInterval(()=>X(a,t,n,d),ye),X(a,t,n,d)}async function X(e,a,t,n){const s=L();if(s.leadActivo)try{const i=await g.conversacion(s.leadActivo);w({conversacion:i}),me(e,i.mensajes),i.turno_en_proceso?T||(T=setTimeout(()=>K(t,!0),$e)):(T&&(clearTimeout(T),T=null),K(t,!1)),R?n.classList.add("visible"):a.scrollTop=a.scrollHeight}catch(i){i instanceof E&&i.estadoHTTP>=500&&j("Error de conexión. Reintentando…")}}async function Y(e){const a=e.value.trim();if(!a)return;const t=L();if(t.leadActivo){e.value="";try{await g.enviarTexto(t.leadActivo,a)}catch(n){e.value=a,n instanceof E&&j(n.message)}}}async function Ee(e,a,t,n,s){try{const i=await navigator.mediaDevices.getUserMedia({audio:!0});h=new MediaRecorder(i),M=[],A=0,h.ondataavailable=o=>{M.push(o.data)},h.onstop=()=>{i.getTracks().forEach(o=>o.stop()),Le(e,a,n,s)},h.start(),e.style.display="none",n.style.display="none",s.style.display="none",a.classList.add("activo"),t.textContent="0:00",P=setInterval(()=>{A++;const o=Math.floor(A/60),c=A%60;t.textContent=`${o}:${c.toString().padStart(2,"0")}`,A>=he&&oe(e,a,n,s)},1e3)}catch{j("No se pudo acceder al micrófono.")}}function oe(e,a,t,n){P&&(clearInterval(P),P=null),h&&h.state!=="inactive"&&h.stop(),a.classList.remove("activo"),e.style.display="",t.style.display="",n.style.display=""}async function Le(e,a,t,n){if(M.length===0)return;const s=L();if(!s.leadActivo)return;const i=new Blob(M,{type:M[0].type}),o=new FileReader;o.onloadend=async()=>{const c=o.result.split(",")[1],u=i.type,l=A;try{await g.enviarAudio(s.leadActivo,c,u,l)}catch(d){d instanceof E&&(d.codigo==="AUDIO_INVALIDO"?j("No te escuché bien, ¿me lo repites o me lo escribes?"):j(d.message))}},o.readAsDataURL(i),a.classList.remove("activo"),e.style.display="",t.style.display="",n.style.display=""}function j(e){const a=document.querySelector(".toast-error");a&&a.remove();const t=document.createElement("div");t.className="toast-error",t.textContent=e,document.body.appendChild(t),setTimeout(()=>t.remove(),4e3)}const F=6;function U(e){return e>=.7?"Alta":e>=.4?"Media":"Baja"}let m=0,C="Todos",_="";function x(e,a,t,n){var d,v;const s={Todos:a.leads.length,Alta:0,Media:0,Baja:0};for(const r of a.leads)s[U(r.prioridad)]++;const i=a.leads.filter(r=>{const de=C==="Todos"||U(r.prioridad)===C,ue=!_||r.nombre.toLocaleLowerCase("es").includes(_);return de&&ue}),o=Math.max(1,Math.ceil(i.length/F));m>=o&&(m=o-1);const c=m*F,u=i.slice(c,c+F);e.innerHTML=`
    <div class="leads-heading">
      <span>Todos los leads</span>
      <button type="button" aria-label="Configurar leads">⚙</button>
    </div>
    <label class="search-wrap">
      <span>⌕</span>
      <input id="buscar-lead" type="search" placeholder="Buscar lead..." value="${D(_)}">
    </label>
    <nav class="filter-tabs" aria-label="Filtrar leads">
      ${["Todos","Alta","Media","Baja"].map(r=>`
        <button type="button" class="filter ${C===r?"active":""}" data-filtro="${r}">
          ${r} <strong>${s[r]}</strong>
        </button>`).join("")}
    </nav>
    <div class="lead-list" id="lead-list" role="list">
      ${u.length?u.map(r=>Te(r,r.lead_id===t)).join(""):'<p style="padding:1rem;color:#647198;font-size:13px">Ningún lead coincide.</p>'}
    </div>
    <div class="list-pager">
      <button type="button" class="pager-button" id="pager-prev" ${m===0?"disabled":""} aria-label="Anterior">‹</button>
      <span>Mostrando ${u.length?c+1:0}–${c+u.length} de ${i.length}</span>
      <button type="button" class="pager-button" id="pager-next" ${m>=o-1?"disabled":""} aria-label="Siguiente">›</button>
    </div>
  `,e.querySelectorAll("[data-lead-id]").forEach(r=>{r.addEventListener("click",()=>n(r.dataset.leadId))});const l=e.querySelector("#buscar-lead");if(l.addEventListener("input",()=>{_=l.value.trim().toLocaleLowerCase("es"),m=0,x(e,a,t,n)}),document.activeElement!==l&&_){const r=l.value.length;l.focus(),l.setSelectionRange(r,r)}e.querySelectorAll(".filter").forEach(r=>{r.addEventListener("click",()=>{C=r.dataset.filtro,m=0,x(e,a,t,n)})}),(d=e.querySelector("#pager-prev"))==null||d.addEventListener("click",()=>{m=Math.max(0,m-1),x(e,a,t,n)}),(v=e.querySelector("#pager-next"))==null||v.addEventListener("click",()=>{m=Math.min(o-1,m+1),x(e,a,t,n)})}function Te(e,a){const t=U(e.prioridad),n=t==="Alta"?"hot":t==="Media"?"medium":"",s=t==="Media"?"medium":t==="Baja"?"baja":"",i=e.afiliado?"Afiliado":"No afiliado";return`
    <button type="button" class="lead-row ${n} ${a?"selected":""}" data-lead-id="${e.lead_id}" role="listitem">
      <img class="avatar lead-avatar" src="${k(e.lead_id)}" alt="" aria-hidden="true">
      <span class="lead-info">
        <span class="lead-name">${D(e.nombre)}</span>
        <span class="lead-meta">${D(i)}</span>
        <span class="lead-last">${D(e.resumen)}</span>
      </span>
      <span class="priority ${s}">${t}</span>
    </button>
  `}function D(e){const a=document.createElement("div");return a.textContent=e,a.innerHTML}const _e={ASESOR:{icono:"👥",titulo:"ASESOR",detalle:"Lead listo para asesor comercial. Prioridad alta en cola."},NUTRICION:{icono:"🌱",titulo:"NUTRICIÓN",detalle:"Aún no alcanza el perfil objetivo. En plan de acompañamiento."},REMARKETING:{icono:"🔁",titulo:"REMARKETING",detalle:"Capacidad alta, intención baja. Reintentar contacto más adelante."},DESPEDIDA:{icono:"👋",titulo:"DESPEDIDA",detalle:"No cumple condiciones mínimas para continuar el proceso."}},V={VERIFICADO_BASE:{txt:"verificado",clase:"verified"},DECLARADO:{txt:"declarado",clase:"declared-badge"},INFERIDO:{txt:"inferido",clase:"declared-badge"}},xe={ingreso_hogar:"Ingreso hogar",recursos_propios:"Recursos propios",tiene_vivienda:"Tiene vivienda",recibio_subsidio:"Recibió subsidio antes",edad:"Edad",personas_hogar:"Personas en el hogar",zona_deseada:"Zona deseada"};let I=0;function Z(e,a,t,n){if(I=0,!a){Ce(e,t);return}re(e,a,n)}function re(e,a,t){var d;const n=a.identificacion,s=a.capacidad,i=t?_e[t]:null,o=a.intencion.nivel,c=o==="MEDIA"?"media":o==="BAJA"?"baja":"",u=we(a.perfil);e.innerHTML=`
    <header class="profile-head">
      <img class="avatar profile-avatar" src="${k(a.lead_id)}" alt="" aria-hidden="true">
      <div>
        <div class="profile-name">${p(n.nombre||"Lead")}</div>
        <div class="profile-sub">${n.afiliada?`Afiliada · Categoría ${p(n.categoria||"N/A")}`:"No afiliado"}</div>
        <div class="lead-id">ID Lead: ${p(a.lead_id)} <button class="copy" id="btn-copiar-id" type="button" title="Copiar ID">▣</button></div>
      </div>
      <span class="attention ${c}">${o==="ALTA"?"Alta":o==="MEDIA"?"Media":"Baja"} intención</span>
    </header>

    ${a.banda_advertencia?`<div class="banda-advertencia" role="alert" style="margin-top:12px;padding:10px 12px;border-radius:8px;background:#fff1c9;color:#7a4b00;font-size:12px;font-weight:700">⚠️ ${p(a.banda_advertencia)}</div>`:""}

    ${i?`
      <div class="route-title"><span>Ruta asignada</span><span class="updated">Actualizado: ahora</span></div>
      <div class="route-card">
        <span class="route-icon" aria-hidden="true">${i.icono}</span>
        <div><strong>${i.titulo}</strong><p>${i.detalle}</p></div>
      </div>
    `:""}

    <div class="finance-grid" style="display:grid;grid-template-columns:1fr;gap:0">
      <section class="finance-card">
        <div class="section-title"><span>Capacidad financiera</span></div>
        <div class="stats">
          <div class="stat"><small>Capacidad estimada</small><strong>${y(s.presupuesto_max)}</strong></div>
          <div class="stat"><small>Cuota mensual (40% del ingreso)</small><strong>${u!==null?y(u):"—"}</strong></div>
        </div>
        <div class="progress"><span style="width:40%"></span></div>
        <div class="rule">Regla cuota/ingreso: 40% · Ratio vs. proyecto: ${(s.ratio*100).toFixed(0)}% · Confianza: ${(s.confianza*100).toFixed(0)}%</div>
      </section>

      ${s.desglose.length?`
        <section class="breakdown">
          <h3>Desglose</h3>
          ${s.desglose.map(Se).join("")}
        </section>
      `:""}

      ${je(s.desglose)}

      ${Object.keys(a.perfil).length?`
        <section class="declared">
          <h3>Perfil declarado</h3>
          ${Object.entries(a.perfil).map(([v,r])=>Me(v,r.valor,r.fuente)).join("")}
        </section>
      `:""}
    </div>

    ${Ie(a)}
  `,(d=e.querySelector("#btn-copiar-id"))==null||d.addEventListener("click",()=>{navigator.clipboard.writeText(a.lead_id).catch(()=>{})});const l=e.querySelector("#btn-ver-alternativas");l&&a.recomendaciones.length>1&&l.addEventListener("click",()=>{I=(I+1)%a.recomendaciones.length,re(e,a,t)})}function Se(e){const a=e.fuente?V[e.fuente]:null;return`
    <div class="detail-line">
      <span>${p(e.concepto)}<br><small>regla: ${p(e.regla)}</small></span>
      <b>${y(e.monto)}${a?`<span class="${a.clase}">${a.txt}</span>`:""}</b>
    </div>
  `}function Me(e,a,t){const n=V[t]??V.DECLARADO,s=xe[e]??e;let i;return typeof a=="boolean"?i=a?"Sí":"No":typeof a=="number"&&e.includes("ingreso")||e.includes("recursos")?i=y(a):i=String(a),`
    <div class="detail-line">
      <span>${p(s)}</span>
      <b>${p(i)} <span class="${n.clase}" style="display:inline">${n.txt}</span></b>
    </div>
  `}function je(e){const a=e.filter(n=>n.concepto.toLowerCase().includes("subsidio"));if(!a.length)return"";const t=a.reduce((n,s)=>n+s.monto,0);return`
    <section class="subsidies">
      <div class="section-title"><span>Subsidios aplicables</span></div>
      ${a.map(n=>`<div class="subsidy-row"><span>${p(n.concepto)}</span><strong>${y(n.monto)}</strong></div>`).join("")}
      <div class="subsidy-total"><span>Total subsidios</span><span>${y(t)}</span></div>
    </section>
  `}function Ie(e){if(!e.recomendaciones.length)return"";const a=e.recomendaciones[I]??e.recomendaciones[0];return`
    <section class="project-card">
      <div class="section-title">
        <span>Proyecto recomendado</span>
        ${e.recomendaciones.length>1?`<button class="link" id="btn-ver-alternativas" type="button">Ver alternativas (${I+1}/${e.recomendaciones.length}) ›</button>`:""}
      </div>
      <div class="project-info">
        <div class="project-image" style="background-image:url('${se(a.proyecto_id)}')" aria-label="Imagen del proyecto ${p(a.nombre)}"></div>
        <div>
          <h3 class="project-name">${p(a.nombre)}</h3>
          <div class="project-place">${p(a.zona)}</div>
          <div class="project-meta">Desde ${y(a.precio_desde)}</div>
          <div class="project-foot">${a.vecinos} vecinos compraron aquí · <strong>${(a.tasa_desistimiento*100).toFixed(0)}% desistió</strong></div>
        </div>
      </div>
    </section>
  `}function we(e){var t;const a=(t=e.ingreso_hogar)==null?void 0:t.valor;return typeof a!="number"||a<=0?null:Math.round(a*.4)}function y(e){return Math.abs(e)>=1e6?`$${(e/1e6).toFixed(1)}M`:`$${e.toLocaleString("es-CO")}`}function Ce(e,a){e.innerHTML=`
    <div class="details-vacio">
      <h3>Ficha aún sin generar</h3>
      <p>La ficha comercial completa de <strong>${p(a)}</strong> se generará automáticamente cuando Vivi complete la calificación.</p>
    </div>
  `}function p(e){const a=document.createElement("div");return a.textContent=e,a.innerHTML}const ce="2026-08-10";function Oe(e,a,t){var n,s;e.innerHTML=`
    <button id="btn-avanzar-tiempo" class="top-action primary" type="button" title="Avanzar la demo a ${ce}">
      <span class="play">▶</span> Avanzar tiempo
    </button>
    <button id="btn-reiniciar-demo" class="top-action" type="button" title="Reiniciar demo al estado inicial">
      <span>↺</span> Reiniciar
    </button>
    <span id="demo-aviso" role="status" aria-live="polite" class="demo-aviso"></span>
  `,(n=e.querySelector("#btn-avanzar-tiempo"))==null||n.addEventListener("click",a),(s=e.querySelector("#btn-reiniciar-demo"))==null||s.addEventListener("click",t)}function Q(){return ce}const Re=5e3;let N=null,W=null,H=null,$=null,f="leads";function Pe(e,a,t,n){n&&n.querySelectorAll(".nav-item").forEach(s=>{s.addEventListener("click",()=>{f=s.dataset.seccion??"leads",n.querySelectorAll(".nav-item").forEach(i=>i.classList.toggle("active",i===s)),ee(e)})}),t&&Be(t),ve(()=>{ee(e),ae(a)}),setInterval(z,Re),z(),ae(a)}async function z(){try{const e=await g.cola();w({cola:e})}catch(e){console.warn("[dashboard] Error cargando cola:",e)}}function De(e,a){if(a!=="conversaciones")return e;const t=[...e.leads].sort((n,s)=>new Date(s.actualizado_en).getTime()-new Date(n.actualizado_en).getTime());return{...e,leads:t}}function ee(e,a=!1){const t=L();if(f==="nutricion"||f==="proyectos"){if(!a&&f===$)return;$=f,e.innerHTML=`
      <div class="leads-heading"><span>${f==="nutricion"?"Nutrición":"Proyectos"}</span></div>
      <div class="details-vacio">
        <p>Próximamente.</p>
      </div>
    `;return}if(!t.cola){if($!==null)return;$=null,e.innerHTML='<div class="leads-heading"><span>Todos los leads</span></div><p style="padding:1rem;color:#647198">Cargando leads…</p>';return}if(!a&&t.cola===W&&t.leadActivo===H&&f===$)return;const n=t.leadActivo!==H;if(W=t.cola,H=t.leadActivo,$=f,n&&t.leadActivo){const i=t.cola.leads.find(o=>o.lead_id===t.leadActivo);i&&ie(i.lead_id,i.nombre)}const s=De(t.cola,f);x(e,s,t.leadActivo,i=>Ne(i))}function Ne(e){var n;const t=(n=L().cola)==null?void 0:n.leads.find(s=>s.lead_id===e);t&&ie(t.lead_id,t.nombre),w({leadActivo:e})}async function ae(e){var i;const a=L();if(!a.leadActivo){N=null,e.innerHTML='<div class="details-vacio">Seleccioná un lead de la lista para ver su ficha comercial.</div>';return}if(a.leadActivo===N)return;N=a.leadActivo;const t=(i=a.cola)==null?void 0:i.leads.find(o=>o.lead_id===a.leadActivo),n=(t==null?void 0:t.nombre)??"Lead",s=(t==null?void 0:t.ruta)??null;try{const o=await g.ficha(a.leadActivo);Z(e,o,n,s)}catch(o){o instanceof E&&o.estadoHTTP===404?Z(e,null,n,s):e.innerHTML=`<div class="details-vacio">⚠️ Error cargando la ficha comercial: ${ze(o.message)}</div>`}}function O(e,a=!1){const t=document.getElementById("demo-aviso");t&&(t.textContent=e,t.classList.toggle("es-error",a),setTimeout(()=>{t.textContent="",t.classList.remove("es-error")},4e3))}function Be(e){Oe(e,async()=>{try{const a=await g.avanzarTiempo(Q());O(`Tiempo avanzado a ${Q()} · ${a.hitos_disparados} hito(s) disparado(s).`),z()}catch(a){O(`Error al avanzar tiempo: ${a.message}`,!0)}},async()=>{try{await g.reiniciar();let a=null;try{a=(await g.crearLead("ana")).lead_id}catch(t){console.error("[dashboard] no se pudo recrear el lead tras reiniciar:",t)}N=null,w({leadActivo:a}),z(),O("Demo reiniciada al estado inicial.")}catch(a){O(`Error al reiniciar demo: ${a.message}`,!0)}})}function ze(e){const a=document.createElement("div");return a.textContent=e,a.innerHTML}const ke=new URLSearchParams(location.search),Fe=["ana","carlos","luisa"];function He(){const e=ke.get("precargado");return Fe.includes(e??"")?e:"ana"}const le=document.querySelector("#app");if(!le)throw new Error("No se encontró el contenedor de Vivi");le.innerHTML=`
  <div class="shell">
    <header class="topbar">
      <div class="brand">
        <img src="/logo-colsubsidio-blanco.png" alt="Colsubsidio">
        <span class="v-divider"></span>
        <span class="vivi-title"><strong>Vivi <span class="sun">☀</span></strong>Asesora de vivienda</span>
      </div>
      <div class="demo-actions" id="botonera-demo"></div>
      <div class="account">
        <span class="bell">🔔<span class="notification">3</span></span>
        <img class="avatar account-avatar" src="${k("asesora-ana-gomez")}" alt="" aria-hidden="true">
        <span>Ana Gómez</span>
        <span class="chevron">⌄</span>
      </div>
    </header>
    <main class="workspace">
      <aside class="sidebar" aria-label="Navegación principal">
        <nav class="nav-main">
          <button class="nav-item active" data-seccion="leads"><span class="nav-icon">👤</span><span class="nav-label">Leads</span></button>
          <button class="nav-item" data-seccion="conversaciones"><span class="nav-icon">💬</span><span class="nav-label">Conversaciones</span></button>
          <button class="nav-item" data-seccion="nutricion"><span class="nav-icon">🌱</span><span class="nav-label">Nutrición</span></button>
          <button class="nav-item" data-seccion="proyectos"><span class="nav-icon">🏢</span><span class="nav-label">Proyectos</span></button>
        </nav>
        <div class="side-promo">
          <span class="promo-icon">🏡</span>
          <h2>Convertimos<br>leads en<br>vecinos</h2>
          <p class="promo-copy">Estamos para ayudarte<br>a crear hogares felices.</p>
        </div>
      </aside>
      <section class="leads-panel" id="leads-panel" aria-label="Lista de leads"></section>
      <section class="chat-panel" id="panel-chat" aria-label="Conversación"></section>
      <section class="details" id="details-panel" aria-label="Ficha comercial del lead"></section>
    </main>
  </div>
`;async function qe(){try{return(await g.crearLead(He())).lead_id}catch(e){const a=e instanceof E?`${e.codigo}: ${e.message}`:String(e);return console.error("[main] no se pudo crear el lead inicial:",a),Ue(a),null}}function Ue(e){const a=document.getElementById("panel-chat");if(!a)return;const t=document.createElement("div");t.setAttribute("role","alert"),t.style.cssText="margin:12px;padding:12px;border-radius:8px;background:#fff1c9;color:#7a4b00;font-size:13px;font-weight:600",t.textContent=`No se pudo iniciar la conversación (${e}). Recargá la página o revisá que el backend esté arriba.`,a.prepend(t)}async function Ve(){const e=document.getElementById("leads-panel"),a=document.getElementById("panel-chat"),t=document.getElementById("details-panel"),n=document.getElementById("botonera-demo"),s=document.querySelector(".nav-main"),i=await qe();i&&w({leadActivo:i}),Ae(a),Pe(e,t,n,s),console.info("Vivi web iniciado (Leads + Chat + Ficha, conectados a la API real)")}Ve();
