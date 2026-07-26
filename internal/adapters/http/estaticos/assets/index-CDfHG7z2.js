(function(){const a=document.createElement("link").relList;if(a&&a.supports&&a.supports("modulepreload"))return;for(const s of document.querySelectorAll('link[rel="modulepreload"]'))n(s);new MutationObserver(s=>{for(const i of s)if(i.type==="childList")for(const o of i.addedNodes)o.tagName==="LINK"&&o.rel==="modulepreload"&&n(o)}).observe(document,{childList:!0,subtree:!0});function t(s){const i={};return s.integrity&&(i.integrity=s.integrity),s.referrerPolicy&&(i.referrerPolicy=s.referrerPolicy),s.crossOrigin==="use-credentials"?i.credentials="include":s.crossOrigin==="anonymous"?i.credentials="omit":i.credentials="same-origin",i}function n(s){if(s.ep)return;s.ep=!0;const i=t(s);fetch(s.href,i)}})();const ie="/api";class A extends Error{constructor(a,t,n){super(t),this.codigo=a,this.estadoHTTP=n,this.name="ErrorAPI"}}async function b(e,a){var n,s;const t=await fetch(`${ie}${e}`,{headers:{"Content-Type":"application/json"},...a});if(!t.ok){const i=await t.json().catch(()=>null);throw new A(((n=i==null?void 0:i.error)==null?void 0:n.codigo)??"ERROR_INTERNO",((s=i==null?void 0:i.error)==null?void 0:s.mensaje)??`HTTP ${t.status}`,t.status)}return await t.json()}const g={crearLead:e=>b("/leads",{method:"POST",body:JSON.stringify({precargado_id:e,fuente:"DEMO"})}),enviarTexto:(e,a)=>b(`/leads/${e}/mensajes`,{method:"POST",body:JSON.stringify({tipo:"TEXTO",texto:a})}),enviarAudio:(e,a,t,n)=>b(`/leads/${e}/mensajes`,{method:"POST",body:JSON.stringify({tipo:"AUDIO",audio_base64:a,mime:t,duracion_s:n})}),conversacion:e=>b(`/leads/${e}/conversacion`),cola:()=>b("/leads"),ficha:e=>b(`/leads/${e}/ficha`),buyerPersona:e=>b(`/gerencia/buyer-persona${e?`?proyecto_id=${e}`:""}`),avanzarTiempo:e=>b("/demo/tiempo",{method:"POST",body:JSON.stringify({avanzar_hasta:e})}),reiniciar:()=>b("/demo/reiniciar",{method:"POST"})},F={leadActivo:null,conversacion:null,cola:null,tabActiva:"cola"},Y=[];function L(){return F}function O(e){Object.assign(F,e),Y.forEach(a=>a(F))}function oe(e){Y.push(e)}function re(e,a){e.innerHTML=a.map(ce).join("")}function ce(e){var o,c;if(e.tipo_contenido==="SISTEMA")return`<div class="pildora-sistema">${I(e.texto)}</div>`;const a=e.autor==="LEAD"?"derecha":"izquierda",t=new Date(e.creado_en).toLocaleTimeString("es-CO",{hour:"2-digit",minute:"2-digit"}),n=e.autor==="LEAD"?'<span class="chulos">✓✓</span>':"",s=(o=e.adjunto)!=null&&o.audio_original?'<span class="icono-audio" aria-label="nota de voz">🎙️</span>':"",i=(c=e.adjunto)!=null&&c.recomendaciones?le(e.adjunto.recomendaciones):"";return`
    <div class="burbuja ${a}">
      ${s}<p>${I(e.texto)}</p>
      <span class="hora">${t}${n}</span>
    </div>${i}`}function le(e){return`<div class="carrusel" role="list">${e.slice(0,3).map(a=>`
    <article class="tarjeta-proyecto" role="listitem">
      <header class="franja-azul">${I(a.nombre)}</header>
      <p class="zona">${I(a.zona)}</p>
      <p class="precio">Desde $${(a.precio_desde/1e6).toFixed(0)}M</p>
      <p class="razon">${I(a.razon)}</p>
      <p class="evidencia">${a.vecinos} personas con tu perfil compraron aquí ·
         ${(a.tasa_desistimiento*100).toFixed(0)}% desistió</p>
      <a class="btn-primario" href="${encodeURI(a.brochure_url)}" target="_blank" rel="noopener">Ver brochure</a>
      <a class="btn-secundario" href="${encodeURI(a.recorrido_360_url)}" target="_blank" rel="noopener">Recorrido 360°</a>
    </article>`).join("")}</div>`}function U(e,a){e.classList.toggle("visible",a)}function de(e){e.innerHTML=`
    <header class="chat-top">
      <button class="back" aria-label="Volver">‹</button>
      <span class="avatar chat-avatar" id="chat-avatar" aria-hidden="true">?</span>
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
  `}function ue(e){const a=document.getElementById("chat-nombre-lead"),t=document.getElementById("chat-avatar");if(a&&(a.textContent=e),t){const n=e.trim().split(/\s+/).map(s=>s[0]).slice(0,2).join("").toUpperCase()||"?";t.textContent=n}}function I(e){const a=document.createElement("div");return a.textContent=e,a.innerHTML}const pe=60,ve=1500,me=300;let T=null,w=!1,y=null,M=[],N=null,E=0;function be(e){de(e);const a=document.getElementById("contenedor-mensajes"),t=document.getElementById("mensajes-scroll"),n=document.getElementById("indicador-escribiendo"),s=document.getElementById("input-mensaje"),i=document.getElementById("btn-enviar"),o=document.getElementById("btn-mic"),c=document.getElementById("mic-grabando"),v=document.getElementById("contador-mic"),u=document.getElementById("btn-detener-mic"),l=document.getElementById("btn-nuevos");s.addEventListener("keydown",m=>{m.key==="Enter"&&!m.shiftKey&&(m.preventDefault(),G(s))}),i.addEventListener("click",()=>G(s)),o.addEventListener("click",()=>fe(o,c,v,s,i)),u.addEventListener("click",()=>ee(o,c,s,i)),t.addEventListener("scroll",()=>{w=t.scrollHeight-t.scrollTop-t.clientHeight>80,w||l.classList.remove("visible")}),l.addEventListener("click",()=>{t.scrollTop=t.scrollHeight,l.classList.remove("visible"),w=!1}),setInterval(()=>V(a,t,n,l),ve),V(a,t,n,l)}async function V(e,a,t,n){const s=L();if(s.leadActivo)try{const i=await g.conversacion(s.leadActivo);O({conversacion:i}),re(e,i.mensajes),i.turno_en_proceso?T||(T=setTimeout(()=>U(t,!0),me)):(T&&(clearTimeout(T),T=null),U(t,!1)),w?n.classList.add("visible"):a.scrollTop=a.scrollHeight}catch(i){i instanceof A&&i.estadoHTTP>=500&&C("Error de conexión. Reintentando…")}}async function G(e){const a=e.value.trim();if(!a)return;const t=L();if(t.leadActivo){e.value="";try{await g.enviarTexto(t.leadActivo,a)}catch(n){e.value=a,n instanceof A&&C(n.message)}}}async function fe(e,a,t,n,s){try{const i=await navigator.mediaDevices.getUserMedia({audio:!0});y=new MediaRecorder(i),M=[],E=0,y.ondataavailable=o=>{M.push(o.data)},y.onstop=()=>{i.getTracks().forEach(o=>o.stop()),ge(e,a,n,s)},y.start(),e.style.display="none",n.style.display="none",s.style.display="none",a.classList.add("activo"),t.textContent="0:00",N=setInterval(()=>{E++;const o=Math.floor(E/60),c=E%60;t.textContent=`${o}:${c.toString().padStart(2,"0")}`,E>=pe&&ee(e,a,n,s)},1e3)}catch{C("No se pudo acceder al micrófono.")}}function ee(e,a,t,n){N&&(clearInterval(N),N=null),y&&y.state!=="inactive"&&y.stop(),a.classList.remove("activo"),e.style.display="",t.style.display="",n.style.display=""}async function ge(e,a,t,n){if(M.length===0)return;const s=L();if(!s.leadActivo)return;const i=new Blob(M,{type:M[0].type}),o=new FileReader;o.onloadend=async()=>{const c=o.result.split(",")[1],v=i.type,u=E;try{await g.enviarAudio(s.leadActivo,c,v,u)}catch(l){l instanceof A&&(l.codigo==="AUDIO_INVALIDO"?C("No te escuché bien, ¿me lo repites o me lo escribes?"):C(l.message))}},o.readAsDataURL(i),a.classList.remove("activo"),e.style.display="",t.style.display="",n.style.display=""}function C(e){const a=document.querySelector(".toast-error");a&&a.remove();const t=document.createElement("div");t.className="toast-error",t.textContent=e,document.body.appendChild(t),setTimeout(()=>t.remove(),4e3)}const k=6;function H(e){return e>=.7?"Alta":e>=.4?"Media":"Baja"}let p=0,D="Todos",j="";function S(e,a,t,n){var l,m;const s={Todos:a.leads.length,Alta:0,Media:0,Baja:0};for(const r of a.leads)s[H(r.prioridad)]++;const i=a.leads.filter(r=>{const R=D==="Todos"||H(r.prioridad)===D,se=!j||r.nombre.toLocaleLowerCase("es").includes(j);return R&&se}),o=Math.max(1,Math.ceil(i.length/k));p>=o&&(p=o-1);const c=p*k,v=i.slice(c,c+k);e.innerHTML=`
    <div class="leads-heading">
      <span>Todos los leads</span>
      <button type="button" aria-label="Configurar leads">⚙</button>
    </div>
    <label class="search-wrap">
      <span>⌕</span>
      <input id="buscar-lead" type="search" placeholder="Buscar lead..." value="${_(j)}">
    </label>
    <nav class="filter-tabs" aria-label="Filtrar leads">
      ${["Todos","Alta","Media","Baja"].map(r=>`
        <button type="button" class="filter ${D===r?"active":""}" data-filtro="${r}">
          ${r} <strong>${s[r]}</strong>
        </button>`).join("")}
    </nav>
    <div class="lead-list" id="lead-list" role="list">
      ${v.length?v.map(r=>ye(r,r.lead_id===t)).join(""):'<p style="padding:1rem;color:#647198;font-size:13px">Ningún lead coincide.</p>'}
    </div>
    <div class="list-pager">
      <button type="button" class="pager-button" id="pager-prev" ${p===0?"disabled":""} aria-label="Anterior">‹</button>
      <span>Mostrando ${v.length?c+1:0}–${c+v.length} de ${i.length}</span>
      <button type="button" class="pager-button" id="pager-next" ${p>=o-1?"disabled":""} aria-label="Siguiente">›</button>
    </div>
  `,e.querySelectorAll("[data-lead-id]").forEach(r=>{r.addEventListener("click",()=>n(r.dataset.leadId))});const u=e.querySelector("#buscar-lead");if(u.addEventListener("input",()=>{j=u.value.trim().toLocaleLowerCase("es"),p=0,S(e,a,t,n)}),document.activeElement!==u&&j){const r=u.value.length;u.focus(),u.setSelectionRange(r,r)}e.querySelectorAll(".filter").forEach(r=>{r.addEventListener("click",()=>{D=r.dataset.filtro,p=0,S(e,a,t,n)})}),(l=e.querySelector("#pager-prev"))==null||l.addEventListener("click",()=>{p=Math.max(0,p-1),S(e,a,t,n)}),(m=e.querySelector("#pager-next"))==null||m.addEventListener("click",()=>{p=Math.min(o-1,p+1),S(e,a,t,n)})}function ye(e,a){const t=H(e.prioridad),n=t==="Alta"?"hot":t==="Media"?"medium":"",s=t==="Media"?"medium":t==="Baja"?"baja":"",i=e.nombre.trim().split(/\s+/).map(c=>c[0]).slice(0,2).join("").toUpperCase()||"?",o=e.afiliado?"Afiliado":"No afiliado";return`
    <button type="button" class="lead-row ${n} ${a?"selected":""}" data-lead-id="${e.lead_id}" role="listitem">
      <span class="avatar lead-avatar" aria-hidden="true">${_(i)}</span>
      <span class="lead-info">
        <span class="lead-name">${_(e.nombre)}</span>
        <span class="lead-meta">${_(o)}</span>
        <span class="lead-last">${_(e.resumen)}</span>
      </span>
      <span class="priority ${s}">${t}</span>
    </button>
  `}function _(e){const a=document.createElement("div");return a.textContent=e,a.innerHTML}const he={ASESOR:{icono:"👥",titulo:"ASESOR",detalle:"Lead listo para asesor comercial. Prioridad alta en cola."},NUTRICION:{icono:"🌱",titulo:"NUTRICIÓN",detalle:"Aún no alcanza el perfil objetivo. En plan de acompañamiento."},REMARKETING:{icono:"🔁",titulo:"REMARKETING",detalle:"Capacidad alta, intención baja. Reintentar contacto más adelante."},DESPEDIDA:{icono:"👋",titulo:"DESPEDIDA",detalle:"No cumple condiciones mínimas para continuar el proceso."}},q={VERIFICADO_BASE:{txt:"verificado",clase:"verified"},DECLARADO:{txt:"declarado",clase:"declared-badge"},INFERIDO:{txt:"inferido",clase:"declared-badge"}},$e={ingreso_hogar:"Ingreso hogar",recursos_propios:"Recursos propios",tiene_vivienda:"Tiene vivienda",recibio_subsidio:"Recibió subsidio antes",edad:"Edad",personas_hogar:"Personas en el hogar",zona_deseada:"Zona deseada"};let x=0;function J(e,a,t,n){if(x=0,!a){Se(e,t);return}ae(e,a,n)}function ae(e,a,t){var m;const n=a.identificacion,s=a.capacidad,i=t?he[t]:null,o=a.intencion.nivel,c=o==="MEDIA"?"media":o==="BAJA"?"baja":"",v=(n.nombre||"L").trim().split(/\s+/).map(r=>r[0]).slice(0,2).join("").toUpperCase(),u=je(a.perfil);e.innerHTML=`
    <header class="profile-head">
      <span class="avatar profile-avatar" aria-hidden="true">${d(v)}</span>
      <div>
        <div class="profile-name">${d(n.nombre||"Lead")}</div>
        <div class="profile-sub">${n.afiliada?`Afiliada · Categoría ${d(n.categoria||"N/A")}`:"No afiliado"}</div>
        <div class="lead-id">ID Lead: ${d(a.lead_id)} <button class="copy" id="btn-copiar-id" type="button" title="Copiar ID">▣</button></div>
      </div>
      <span class="attention ${c}">${o==="ALTA"?"Alta":o==="MEDIA"?"Media":"Baja"} intención</span>
    </header>

    ${a.banda_advertencia?`<div class="banda-advertencia" role="alert" style="margin-top:12px;padding:10px 12px;border-radius:8px;background:#fff1c9;color:#7a4b00;font-size:12px;font-weight:700">⚠️ ${d(a.banda_advertencia)}</div>`:""}

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
          <div class="stat"><small>Capacidad estimada</small><strong>${h(s.presupuesto_max)}</strong></div>
          <div class="stat"><small>Cuota mensual (40% del ingreso)</small><strong>${u!==null?h(u):"—"}</strong></div>
        </div>
        <div class="progress"><span style="width:40%"></span></div>
        <div class="rule">Regla cuota/ingreso: 40% · Ratio vs. proyecto: ${(s.ratio*100).toFixed(0)}% · Confianza: ${(s.confianza*100).toFixed(0)}%</div>
      </section>

      ${s.desglose.length?`
        <section class="breakdown">
          <h3>Desglose</h3>
          ${s.desglose.map(Ee).join("")}
        </section>
      `:""}

      ${Le(s.desglose)}

      ${Object.keys(a.perfil).length?`
        <section class="declared">
          <h3>Perfil declarado</h3>
          ${Object.entries(a.perfil).map(([r,R])=>Ae(r,R.valor,R.fuente)).join("")}
        </section>
      `:""}
    </div>

    ${Te(a)}
  `,(m=e.querySelector("#btn-copiar-id"))==null||m.addEventListener("click",()=>{navigator.clipboard.writeText(a.lead_id).catch(()=>{})});const l=e.querySelector("#btn-ver-alternativas");l&&a.recomendaciones.length>1&&l.addEventListener("click",()=>{x=(x+1)%a.recomendaciones.length,ae(e,a,t)})}function Ee(e){const a=e.fuente?q[e.fuente]:null;return`
    <div class="detail-line">
      <span>${d(e.concepto)}<br><small>regla: ${d(e.regla)}</small></span>
      <b>${h(e.monto)}${a?`<span class="${a.clase}">${a.txt}</span>`:""}</b>
    </div>
  `}function Ae(e,a,t){const n=q[t]??q.DECLARADO,s=$e[e]??e;let i;return typeof a=="boolean"?i=a?"Sí":"No":typeof a=="number"&&e.includes("ingreso")||e.includes("recursos")?i=h(a):i=String(a),`
    <div class="detail-line">
      <span>${d(s)}</span>
      <b>${d(i)} <span class="${n.clase}" style="display:inline">${n.txt}</span></b>
    </div>
  `}function Le(e){const a=e.filter(n=>n.concepto.toLowerCase().includes("subsidio"));if(!a.length)return"";const t=a.reduce((n,s)=>n+s.monto,0);return`
    <section class="subsidies">
      <div class="section-title"><span>Subsidios aplicables</span></div>
      ${a.map(n=>`<div class="subsidy-row"><span>${d(n.concepto)}</span><strong>${h(n.monto)}</strong></div>`).join("")}
      <div class="subsidy-total"><span>Total subsidios</span><span>${h(t)}</span></div>
    </section>
  `}function Te(e){if(!e.recomendaciones.length)return"";const a=e.recomendaciones[x]??e.recomendaciones[0];return`
    <section class="project-card">
      <div class="section-title">
        <span>Proyecto recomendado</span>
        ${e.recomendaciones.length>1?`<button class="link" id="btn-ver-alternativas" type="button">Ver alternativas (${x+1}/${e.recomendaciones.length}) ›</button>`:""}
      </div>
      <div class="project-info">
        <div class="project-image" aria-label="Imagen del proyecto ${d(a.nombre)}"></div>
        <div>
          <h3 class="project-name">${d(a.nombre)}</h3>
          <div class="project-place">${d(a.zona)}</div>
          <div class="project-meta">Desde ${h(a.precio_desde)}</div>
          <div class="project-foot">${a.vecinos} vecinos compraron aquí · <strong>${(a.tasa_desistimiento*100).toFixed(0)}% desistió</strong></div>
        </div>
      </div>
    </section>
  `}function je(e){var t;const a=(t=e.ingreso_hogar)==null?void 0:t.valor;return typeof a!="number"||a<=0?null:Math.round(a*.4)}function h(e){return Math.abs(e)>=1e6?`$${(e/1e6).toFixed(1)}M`:`$${e.toLocaleString("es-CO")}`}function Se(e,a){e.innerHTML=`
    <div class="details-vacio">
      <h3>Ficha aún sin generar</h3>
      <p>La ficha comercial completa de <strong>${d(a)}</strong> se generará automáticamente cuando Vivi complete la calificación.</p>
    </div>
  `}function d(e){const a=document.createElement("div");return a.textContent=e,a.innerHTML}const te="2026-08-10";function _e(e,a,t){var n,s;e.innerHTML=`
    <button id="btn-avanzar-tiempo" class="top-action primary" type="button" title="Avanzar la demo a ${te}">
      <span class="play">▶</span> Avanzar tiempo
    </button>
    <button id="btn-reiniciar-demo" class="top-action" type="button" title="Reiniciar demo al estado inicial">
      <span>↺</span> Reiniciar
    </button>
    <span id="demo-aviso" role="status" aria-live="polite" class="demo-aviso"></span>
  `,(n=e.querySelector("#btn-avanzar-tiempo"))==null||n.addEventListener("click",a),(s=e.querySelector("#btn-reiniciar-demo"))==null||s.addEventListener("click",t)}function K(){return te}const Ie=5e3;let B=null,X=null,Z=null,$=null,f="leads";function Me(e,a,t,n){n&&n.querySelectorAll(".nav-item").forEach(s=>{s.addEventListener("click",()=>{f=s.dataset.seccion??"leads",n.querySelectorAll(".nav-item").forEach(i=>i.classList.toggle("active",i===s)),Q(e)})}),t&&Oe(t),oe(()=>{Q(e),W(a)}),setInterval(z,Ie),z(),W(a)}async function z(){try{const e=await g.cola();O({cola:e})}catch(e){console.warn("[dashboard] Error cargando cola:",e)}}function Ce(e,a){if(a!=="conversaciones")return e;const t=[...e.leads].sort((n,s)=>new Date(s.actualizado_en).getTime()-new Date(n.actualizado_en).getTime());return{...e,leads:t}}function Q(e,a=!1){const t=L();if(f==="nutricion"||f==="proyectos"){if(!a&&f===$)return;$=f,e.innerHTML=`
      <div class="leads-heading"><span>${f==="nutricion"?"Nutrición":"Proyectos"}</span></div>
      <div class="details-vacio">
        <p>Próximamente.</p>
      </div>
    `;return}if(!t.cola){if($!==null)return;$=null,e.innerHTML='<div class="leads-heading"><span>Todos los leads</span></div><p style="padding:1rem;color:#647198">Cargando leads…</p>';return}if(!a&&t.cola===X&&t.leadActivo===Z&&f===$)return;X=t.cola,Z=t.leadActivo,$=f;const n=Ce(t.cola,f);S(e,n,t.leadActivo,s=>xe(s))}function xe(e){var n;const t=(n=L().cola)==null?void 0:n.leads.find(s=>s.lead_id===e);t&&ue(t.nombre),O({leadActivo:e})}async function W(e){var i;const a=L();if(!a.leadActivo){B=null,e.innerHTML='<div class="details-vacio">Seleccioná un lead de la lista para ver su ficha comercial.</div>';return}if(a.leadActivo===B)return;B=a.leadActivo;const t=(i=a.cola)==null?void 0:i.leads.find(o=>o.lead_id===a.leadActivo),n=(t==null?void 0:t.nombre)??"Lead",s=(t==null?void 0:t.ruta)??null;try{const o=await g.ficha(a.leadActivo);J(e,o,n,s)}catch(o){o instanceof A&&o.estadoHTTP===404?J(e,null,n,s):e.innerHTML=`<div class="details-vacio">⚠️ Error cargando la ficha comercial: ${Re(o.message)}</div>`}}function P(e,a=!1){const t=document.getElementById("demo-aviso");t&&(t.textContent=e,t.classList.toggle("es-error",a),setTimeout(()=>{t.textContent="",t.classList.remove("es-error")},4e3))}function Oe(e){_e(e,async()=>{try{const a=await g.avanzarTiempo(K());P(`Tiempo avanzado a ${K()} · ${a.hitos_disparados} hito(s) disparado(s).`),z()}catch(a){P(`Error al avanzar tiempo: ${a.message}`,!0)}},async()=>{try{await g.reiniciar();let a=null;try{a=(await g.crearLead("ana")).lead_id}catch(t){console.error("[dashboard] no se pudo recrear el lead tras reiniciar:",t)}B=null,O({leadActivo:a}),z(),P("Demo reiniciada al estado inicial.")}catch(a){P(`Error al reiniciar demo: ${a.message}`,!0)}})}function Re(e){const a=document.createElement("div");return a.textContent=e,a.innerHTML}const De=new URLSearchParams(location.search),Pe=["ana","carlos","luisa"];function we(){const e=De.get("precargado");return Pe.includes(e??"")?e:"ana"}const ne=document.querySelector("#app");if(!ne)throw new Error("No se encontró el contenedor de Vivi");ne.innerHTML=`
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
        <span class="avatar account-avatar" aria-hidden="true">AG</span>
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
`;async function Ne(){try{return(await g.crearLead(we())).lead_id}catch(e){const a=e instanceof A?`${e.codigo}: ${e.message}`:String(e);return console.error("[main] no se pudo crear el lead inicial:",a),Be(a),null}}function Be(e){const a=document.getElementById("panel-chat");if(!a)return;const t=document.createElement("div");t.setAttribute("role","alert"),t.style.cssText="margin:12px;padding:12px;border-radius:8px;background:#fff1c9;color:#7a4b00;font-size:13px;font-weight:600",t.textContent=`No se pudo iniciar la conversación (${e}). Recargá la página o revisá que el backend esté arriba.`,a.prepend(t)}async function ze(){const e=document.getElementById("leads-panel"),a=document.getElementById("panel-chat"),t=document.getElementById("details-panel"),n=document.getElementById("botonera-demo"),s=document.querySelector(".nav-main"),i=await Ne();i&&O({leadActivo:i}),be(a),Me(e,t,n,s),console.info("Vivi web iniciado (Leads + Chat + Ficha, conectados a la API real)")}ze();
