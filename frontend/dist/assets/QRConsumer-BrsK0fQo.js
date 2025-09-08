import{q as ce,K as Ae,B as j,r as Ue,L as ze,o as Be,c as u,a as s,b as v,g as E,e as de,t as y,h as n,F as H,A as Q,j as w,x as D,y as ee,C as Ne,M as Le,v as Ve,z as F,u as $e,i as r,n as ue}from"./vendor-DfF3K0x-.js";import{_ as je,u as He}from"./index-nc5P6C25.js";import{q as te}from"./qrService-D6s-cMlP.js";import{L as Qe,m as Fe}from"./LocationSelector-CCfdcKH2.js";import{p as Oe}from"./pavilionService-bOFFpgyG.js";import{a as Ke}from"./index-NIGUFBhG.js";import{f as O,l as K}from"./index-DJqi8Qh4.js";const Je="http://localhost:8080/api/v1";class Ye{constructor(){this.api=Ke.create({baseURL:Je,headers:{"Content-Type":"application/json"}})}async getAllStores(){try{const p=await this.api.get("/stores");return p.data.data||p.data||[]}catch(p){try{const i=await this.api.get("/stores/");return i.data.data||i.data||[]}catch{throw console.error("Error al obtener almacenes:",p),p}}}async getStoreById(p){try{const i=await this.api.get(`/stores/${p}`);return i.data.data||i.data}catch(i){throw console.error("Error al obtener almacén:",i),i}}async getByMedicalCenter(p){try{return(await this.getAllStores()).filter(I=>I.medical_center_id===parseInt(p))}catch(i){throw console.error("Error al obtener almacenes por centro médico:",i),i}}async createStore(p){try{return(await this.api.post("/stores",p)).data}catch(i){throw console.error("Error al crear almacén:",i),i}}async updateStore(p,i){try{return(await this.api.put(`/stores/${p}`,i)).data}catch(I){throw console.error("Error al actualizar almacén:",I),I}}async deleteStore(p){try{return(await this.api.delete(`/stores/${p}`)).data}catch(i){throw console.error("Error al eliminar almacén:",i),i}}}const Ge=new Ye,We={class:"space-y-6"},Xe={class:"bg-white rounded-lg shadow-sm border p-6"},Ze={class:"flex items-center justify-between"},et={class:"grid grid-cols-2 gap-4 text-right"},tt={class:"text-2xl font-bold text-blue-600"},st={class:"text-2xl font-bold text-green-600"},at={class:"mt-4 flex items-center justify-between bg-gray-50 rounded-lg p-3"},ot={class:"flex items-center space-x-4"},nt={class:"text-sm text-gray-900 ml-1"},rt={class:"text-sm text-gray-900 ml-1"},it={class:"text-xs text-gray-500"},lt={key:0,class:"bg-white rounded-lg shadow-sm border p-6"},ct={class:"grid grid-cols-1 md:grid-cols-3 gap-3"},dt=["onClick"],ut={class:"flex flex-col items-center"},mt=["d"],pt={class:"font-medium"},vt={class:"text-xs text-gray-500 mt-1"},gt={key:1,class:"bg-white rounded-lg shadow-sm border p-6"},yt={class:"grid md:grid-cols-2 gap-6"},ft={class:"flex space-x-2"},bt=["disabled"],ht=["disabled"],xt={class:"text-center"},_t=["disabled"],wt={class:"text-sm"},kt={key:2,class:"bg-white rounded-lg shadow-sm border p-6"},Ct={class:"grid md:grid-cols-2 gap-6"},It={key:0},St=["disabled"],Mt=["value"],Dt={key:0,class:"text-xs text-gray-500 mt-1"},Pt={key:1},Tt=["disabled"],Et=["value"],qt={key:0,class:"text-xs text-gray-500 mt-1"},Rt={key:2},At=["disabled"],Ut=["value"],zt={key:0,class:"text-xs text-gray-500 mt-1"},Bt={class:"mt-6 flex justify-end space-x-3"},Nt=["disabled"],Lt={key:1,class:"h-4 w-4 mr-2",fill:"none",stroke:"currentColor",viewBox:"0 0 24 24"},Vt={__name:"QRConsumer",setup(me){const p={template:`
    <svg class="animate-spin h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24">
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
    </svg>
  `},i={props:["currentLocation"],template:`
    <div class="mb-6 bg-green-50 border border-green-200 rounded-lg p-4">
      <h4 class="text-green-800 font-medium mb-3 flex items-center">
        <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        Modo Consumo Rápido - Pasos Restantes
      </h4>
      <div class="space-y-2 text-sm">
        <ProgressStep :completed="true" text="Producto escaneado automáticamente" />
        <ProgressStep :completed="!!currentLocation" :text="currentLocation ? 'Ubicación seleccionada' : 'Selecciona tu ubicación actual abajo'" />
        <ProgressStep :completed="false" text="Completa los datos del consumo y confirma" />
      </div>
    </div>
  `,components:{ProgressStep:{props:["completed","text"],template:`
        <div class="flex items-center">
          <svg v-if="completed" class="h-4 w-4 text-green-600 mr-2" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
          </svg>
          <svg v-else class="h-4 w-4 text-yellow-500 mr-2" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
          </svg>
          <span :class="completed ? 'text-green-700' : 'text-yellow-700'">
            {{ completed ? '✓' : '⚠' }} {{ text }}
          </span>
        </div>
      `}}},I={props:["product","scanContext","currentUser","currentLocation"],methods:{formatDate(t){if(!t)return"No disponible";try{return O(new Date(t),"dd/MM/yyyy HH:mm",{locale:K})}catch{return t}}},template:`
    <div class="bg-gray-50 rounded-lg p-4 mb-6">
      <div class="grid md:grid-cols-2 gap-4">
        <div>
          <label class="text-sm font-medium text-gray-600">Nombre del Producto:</label>
          <p class="text-gray-900 font-medium">{{ product.supply_code?.name || 'N/A' }}</p>
        </div>
        <div>
          <label class="text-sm font-medium text-gray-600">Código QR:</label>
          <p class="text-sm font-mono text-gray-800 bg-white px-2 py-1 rounded border">{{ product.qr_code }}</p>
        </div>
        <div>
          <label class="text-sm font-medium text-gray-600">Código del Producto:</label>
          <p class="text-gray-900">{{ product.supply_code?.code || product.supply_info?.code || 'N/A' }}</p>
        </div>
        <div>
          <label class="text-sm font-medium text-gray-600">ID Proveedor:</label>
          <p class="text-gray-900">{{ product.supply_code?.code_supplier || 'N/A' }}</p>
        </div>
      </div>

      <!-- Información de lote -->
      <div v-if="batchInfo" class="mt-4 pt-4 border-t border-gray-200">
        <h5 class="text-sm font-semibold text-gray-700 mb-2">Información del Lote:</h5>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
          <div v-if="batchInfo.batch_number">
            <label class="font-medium text-gray-600">Número de Lote:</label>
            <p class="text-gray-900">{{ batchInfo.batch_number }}</p>
          </div>
          <div v-if="batchInfo.supplier">
            <label class="font-medium text-gray-600">Proveedor del Lote:</label>
            <p class="text-gray-900">{{ batchInfo.supplier }}</p>
          </div>
          <div v-if="batchInfo.expiration_date">
            <label class="font-medium text-gray-600">Fecha de Vencimiento:</label>
            <p class="text-gray-900">{{ formatDate(batchInfo.expiration_date) }}</p>
          </div>
        </div>
      </div>

      <!-- Información de trazabilidad -->
      <div v-if="scanContext" class="mt-4 pt-4 border-t border-gray-200">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
          <div>
            <label class="font-medium text-gray-600">Escaneado por:</label>
            <p class="text-gray-900">{{ scanContext.user_name || currentUser?.name || 'N/A' }}</p>
          </div>
          <div>
            <label class="font-medium text-gray-600">Ubicación de escaneo:</label>
            <p class="text-gray-900">{{ scanContext.location || currentLocation?.name || 'N/A' }}</p>
          </div>
          <div>
            <label class="font-medium text-gray-600">Hora de escaneo:</label>
            <p class="text-gray-900">{{ formatDate(scanContext.scanned_at) }}</p>
          </div>
        </div>
      </div>

      <!-- Status Alert -->
      <StatusAlert :product="product" />
    </div>
  `,computed:{batchInfo(){var t;return((t=this.product.supply_info)==null?void 0:t.batch_info)||this.product.batch_info}},components:{StatusAlert:{props:["product"],template:`
        <div v-if="product.is_consumed" class="mt-4 bg-red-50 border border-red-200 rounded p-3">
          <div class="flex items-center">
            <svg class="h-5 w-5 text-red-400 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span class="text-red-800 font-medium">Este producto ya ha sido consumido anteriormente</span>
          </div>
        </div>
        <div v-else-if="!product.can_consume" class="mt-4 bg-yellow-50 border border-yellow-200 rounded p-3">
          <div class="flex items-center">
            <svg class="h-5 w-5 text-yellow-400 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span class="text-yellow-800 font-medium">Este producto no está disponible para consumo</span>
          </div>
        </div>
      `}}},pe={props:["successData"],emits:["clear-all","view-traceability","view-batch-history"],methods:{formatDate(t){if(!t)return"No disponible";try{return O(new Date(t),"dd/MM/yyyy HH:mm",{locale:K})}catch{return t}}},template:`
    <div class="bg-green-50 border border-green-200 rounded-lg p-6">
      <div class="flex items-start">
        <div class="flex-shrink-0">
          <svg class="h-6 w-6 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <div class="ml-3 flex-1">
          <h3 class="text-lg font-medium text-green-800">Producto Consumido Exitosamente</h3>
          <div class="mt-2 text-sm text-green-700">
            <p>{{ successData.message }}</p>
            <p v-if="successData.remaining_amount !== undefined" class="mt-1">
              <strong>Cantidad restante en lote:</strong> {{ successData.remaining_amount }} unidades
            </p>
            <p v-if="successData.batch_history_updated" class="mt-1">
              <strong>Historial del lote actualizado automáticamente</strong>
            </p>
          </div>
          
          <!-- Información de trazabilidad -->
          <div v-if="successData.traceability_info" class="mt-4 bg-green-100 rounded p-3">
            <h4 class="text-sm font-medium text-green-800 mb-2">Información de Trazabilidad</h4>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-2 text-xs text-green-700">
              <div>
                <span class="font-medium">Evento registrado:</span> {{ successData.traceability_info.event_id }}
              </div>
              <div>
                <span class="font-medium">Timestamp:</span> {{ formatDate(successData.traceability_info.timestamp) }}
              </div>
              <div>
                <span class="font-medium">Usuario:</span> {{ successData.traceability_info.user || 'N/A' }}
              </div>
              <div>
                <span class="font-medium">Ubicación:</span> {{ successData.traceability_info.location || 'N/A' }}
              </div>
            </div>
          </div>
          
          <div class="mt-4 flex space-x-3">
            <button @click="$emit('clear-all')" class="btn-secondary text-sm">
              Consumir Otro Producto
            </button>
            <button 
              v-if="successData.qr_code"
              @click="$emit('view-traceability', successData.qr_code)" 
              class="btn-primary text-sm"
            >
              Ver Trazabilidad Completa
            </button>
            <button 
              v-if="successData.batch_id"
              @click="$emit('view-batch-history', successData.batch_id)" 
              class="btn-secondary text-sm"
            >
              Ver Historial del Lote
            </button>
          </div>
        </div>
      </div>
    </div>
  `},ve={props:["errorText","isQuickMode","qrInput"],emits:["clear-error","retry-scan"],template:`
    <div class="bg-red-50 border border-red-200 rounded-lg p-6">
      <div class="flex items-start">
        <div class="flex-shrink-0">
          <svg class="h-6 w-6 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <div class="ml-3 flex-1">
          <h3 class="text-lg font-medium text-red-800">Error</h3>
          <div class="mt-2 text-sm text-red-700">
            <p>{{ errorText }}</p>
            <p v-if="isQuickMode" class="mt-2 text-sm text-red-600">
              <strong>Modo rápido detectado:</strong> El código QR "{{ qrInput }}" no se pudo escanear automáticamente. 
              Puedes intentar escanearlo manualmente o verificar que el producto existe en el sistema.
            </p>
          </div>
          <div class="mt-4 flex space-x-2">
            <button @click="$emit('clear-error')" class="btn-secondary text-sm">
              Intentar de Nuevo
            </button>
            <button v-if="isQuickMode && qrInput" @click="$emit('retry-scan')" class="btn-primary text-sm">
              Reintentar Escaneo Automático
            </button>
          </div>
        </div>
      </div>
    </div>
  `},ge={props:["consumptions"],emits:["view-traceability","view-batch-history"],methods:{formatDate(t){if(!t)return"No disponible";try{return O(new Date(t),"dd/MM/yyyy HH:mm",{locale:K})}catch{return t}},getConsumptionPurposeLabel(t){return{routine:"Rutina",emergency:"Emergencia",maintenance:"Mantenimiento"}[t]||t},getPurposeBadgeClass(t){return t==="emergency"?"bg-red-100 text-red-800":t==="routine"?"bg-blue-100 text-blue-800":"bg-green-100 text-green-800"}},template:`
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <h3 class="text-lg font-medium text-gray-900 mb-4">
        <svg class="h-5 w-5 inline mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        Consumos Recientes
      </h3>

      <div class="space-y-3">
        <div 
          v-for="(consumption, index) in consumptions" 
          :key="index"
          class="flex items-center justify-between p-3 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
        >
          <div class="flex-1">
            <div class="flex items-center justify-between">
              <p class="font-medium text-gray-900">{{ consumption.product_name }}</p>
              <span 
                :class="['px-2 py-1 rounded-full text-xs font-medium', getPurposeBadgeClass(consumption.consumption_purpose)]"
              >
                {{ getConsumptionPurposeLabel(consumption.consumption_purpose) }}
              </span>
            </div>
            <p class="text-sm text-gray-600">QR: {{ consumption.qr_code }}</p>
            <div class="text-xs text-gray-500 mt-1">
              <span>{{ consumption.user_rut }}</span>
              <span class="mx-2">•</span>
              <span>{{ consumption.location || 'Sin ubicación' }}</span>
              <span class="mx-2">•</span>
              <span>{{ formatDate(consumption.consumed_at) }}</span>
            </div>
          </div>
          <div class="ml-4 flex space-x-2">
            <button 
              @click="$emit('view-traceability', consumption.qr_code)"
              class="text-blue-600 hover:text-blue-800 p-1"
              title="Ver trazabilidad"
            >
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
            </button>
            <button 
              v-if="consumption.batch_id"
              @click="$emit('view-batch-history', consumption.batch_id)"
              class="text-green-600 hover:text-green-800 p-1"
              title="Ver historial del lote"
            >
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>
  `},se=$e(),_=He(),ae=ce({qrInput:"",scanning:!1,consuming:!1,cameraActive:!1,scannedProduct:null,consumptionSuccess:null,error:null,recentConsumptions:[],selectedConsumptionPurpose:"routine",currentLocation:null,lastScanContext:null,sessionStart:new Date,isQuickMode:!1,medicalCenters:[],pavilions:[],stores:[],loadingMedicalCenters:!1,loadingPavilions:!1,loadingStores:!1}),{qrInput:b,scanning:P,consuming:q,cameraActive:A,scannedProduct:m,consumptionSuccess:J,error:f,recentConsumptions:k,selectedConsumptionPurpose:T,currentLocation:h,lastScanContext:Y,sessionStart:ye,isQuickMode:U,medicalCenters:R,pavilions:z,stores:B,loadingMedicalCenters:N,loadingPavilions:L,loadingStores:V}=Ae(ae),C=ce({today:0,week:0}),S=j(()=>_.user);j(()=>_.isAuthenticated);const fe=j(()=>m.value&&!m.value.is_consumed&&m.value.can_consume),oe=j(()=>o.value.medicalCenterId?z.value.filter(t=>t.medical_center_id===parseInt(o.value.medicalCenterId)):[]),be=[{value:"routine",label:"Rutina",description:"Consumo programado normal",icon:"M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2",iconClass:"text-blue-600"},{value:"emergency",label:"Emergencia",description:"Consumo de urgencia",icon:"M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z",iconClass:"text-red-600"},{value:"maintenance",label:"Mantenimiento",description:"Uso en mantenimiento",icon:"M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z M15 12a3 3 0 11-6 0 3 3 0 016 0z",iconClass:"text-green-600"}],o=Ue({userRUT:"",destinationType:"",destinationID:"",medicalCenterId:"",notes:""});ze(()=>{var t;(t=S.value)!=null&&t.rut&&!o.value.userRUT&&(o.value.userRUT=S.value.rut)});const M=(t,e="Error de conexión")=>{var l,c;console.error(e+":",t),f.value=((c=(l=t.response)==null?void 0:l.data)==null?void 0:c.error)||t.message||e},he=()=>{var t,e;return{scan_purpose:"consume",consumption_purpose:T.value,pavilion_id:(t=h.value)==null?void 0:t.pavilion_id,medical_center_id:(e=h.value)==null?void 0:e.medical_center_id,scan_source:"web",user_agent:navigator.userAgent,device_info:{platform:navigator.platform,language:navigator.language,screen_resolution:`${screen.width}x${screen.height}`}}},xe=()=>{var l;const t=[];(l=m.value)!=null&&l.qr_code||t.push("Código QR requerido"),o.value.userRUT||t.push("RUT del usuario requerido"),o.value.destinationType||t.push("Tipo de destino requerido"),o.value.destinationType==="pavilion"?(o.value.medicalCenterId||t.push("Centro médico requerido"),o.value.destinationID||t.push("Pabellón requerido")):o.value.destinationType==="store"&&(o.value.destinationID||t.push("Almacén requerido"));const e=parseInt(o.value.destinationID);return isNaN(e)&&t.push("Destino debe ser un número válido"),{isValid:t.length===0,errors:t,destinationIdNum:e}},_e=t=>{var l,c,a,d;let e=`${o.value.destinationType} ${o.value.destinationID}`;if(o.value.destinationType==="pavilion"){const g=oe.value.find(Z=>Z.id===parseInt(o.value.destinationID)),x=R.value.find(Z=>Z.id===parseInt(o.value.medicalCenterId));e=g?`${(x==null?void 0:x.name)||"Centro"} - ${g.name}`:e}else if(o.value.destinationType==="store"){const g=B.value.find(x=>x.id===parseInt(o.value.destinationID));e=g?`Almacén ${g.name}`:e}k.value.unshift({qr_code:m.value.qr_code,product_name:((l=m.value.supply_info)==null?void 0:l.supply_code_name)||"N/A",user_rut:o.value.userRUT,user_name:t.user_name,consumed_at:new Date().toISOString(),consumption_purpose:T.value,location:(c=h.value)==null?void 0:c.name,destination:e,batch_id:(d=(a=m.value.supply_info)==null?void 0:a.batch)==null?void 0:d.id}),k.value.length>15&&(k.value=k.value.slice(0,15)),C.today+=1,C.week+=1,Me()},G=async()=>{N.value=!0;try{const t=await Fe.getAll();R.value=t.data||[]}catch(t){console.error("Error cargando centros médicos:",t),M(t,"Error al cargar centros médicos")}finally{N.value=!1}},W=async()=>{L.value=!0;try{z.value=await Oe.getAllPavilions()}catch(t){console.error("Error cargando pabellones:",t),M(t,"Error al cargar pabellones")}finally{L.value=!1}},ne=async()=>{V.value=!0;try{B.value=await Ge.getAllStores()}catch(t){console.error("Error cargando almacenes:",t),M(t,"Error al cargar almacenes")}finally{V.value=!1}},we=()=>{o.value.destinationID="",o.value.medicalCenterId="",o.value.destinationType==="pavilion"?(R.value.length===0&&G(),z.value.length===0&&W()):o.value.destinationType==="store"&&B.value.length===0&&ne()},ke=()=>{o.value.destinationID=""},$=async()=>{var t,e;if(b.value.trim()){P.value=!0,f.value=null,m.value=null;try{const l=he(),c=await te.scanQRCode(b.value.trim(),l);if(!c){f.value="No se recibió respuesta del servicio de QR";return}if(c.success&&(c.data||c.type)){const a=c.data||c;if(a.type&&a.type!=="medical_supply"){f.value="Solo se pueden consumir productos individuales (no lotes)";return}m.value=a,Y.value={...l,scanned_at:new Date,user_name:(t=S.value)==null?void 0:t.name,location:(e=h.value)==null?void 0:e.name}}else f.value=c.error||c.message||"Error desconocido al escanear código QR"}catch(l){M(l,"Error al escanear código QR")}finally{P.value=!1}}},Ce=async()=>{var t,e,l,c,a;if(m.value){q.value=!0,f.value=null;try{const d=xe();if(!d.isValid){f.value="Faltan campos obligatorios: "+d.errors.join(", ");return}const g={qr_code:m.value.qr_code,user_rut:o.value.userRUT,user_name:((t=S.value)==null?void 0:t.name)||"Encargado Bodega",destination_type:o.value.destinationType,destination_id:d.destinationIdNum,notes:o.value.notes,consumption_purpose:T.value,consumption_context:{pavilion_id:(e=h.value)==null?void 0:e.pavilion_id,medical_center_id:(l=h.value)==null?void 0:l.medical_center_id,user_agent:navigator.userAgent,scan_source:"web"}},x=await te.consumeIndividualSupply(g);x.success?(J.value={...x.data,qr_code:m.value.qr_code,batch_id:(a=(c=m.value.supply_info)==null?void 0:c.batch)==null?void 0:a.id,traceability_info:x.traceability_info,batch_history_updated:!0},_e(g),X(),b.value=""):f.value=x.error||"Error al consumir el producto"}catch(d){M(d,"Error al consumir producto")}finally{q.value=!1}}},Ie=()=>{A.value=!0,alert("Funcionalidad de cámara en desarrollo. Por favor usa el input manual."),A.value=!1},re=async t=>{try{await se.push({name:"QRTraceability",params:{qrCode:t}})}catch(e){M(e,"Error al navegar a la trazabilidad")}},ie=async t=>{if(!t){console.warn("No se proporcionó ID de lote");return}try{await se.push({name:"BatchHistory",params:{batchId:t}})}catch{try{const c=(await te.getBatchHistoryFormatted(t)).map(a=>{const d=a.display_format;return`${d.date}
${d.action}
Cant: ${d.previous_amount} → ${d.new_amount}
${d.user_rut} ${d.user_name}`}).join(`

`);alert(`Historial del Lote:

${c}`)}catch(l){M(l,"Error al cargar historial del lote")}}},Se=t=>{var e;h.value=t,m.value&&t&&!o.value.destinationType&&t.pavilion_id&&(o.value.destinationType="pavilion",o.value.medicalCenterId=((e=t.medical_center_id)==null?void 0:e.toString())||"",o.value.destinationID=t.pavilion_id.toString(),R.value.length===0&&G(),z.value.length===0&&W())},Me=()=>{const t={today_consumptions:C.today,week_consumptions:C.week,recent_consumptions:k.value,last_updated:new Date().toISOString()};localStorage.setItem("consumption-stats",JSON.stringify(t))},De=()=>{try{const t=localStorage.getItem("consumption-stats");if(t){const e=JSON.parse(t);C.today=e.today_consumptions||0,C.week=e.week_consumptions||0,k.value=e.recent_consumptions||[]}}catch(t){console.error("Error loading consumption stats:",t)}},X=()=>{var t;m.value=null,Y.value=null,o.value={userRUT:((t=S.value)==null?void 0:t.rut)||"",destinationType:"",destinationID:"",medicalCenterId:"",notes:""}},le=()=>{f.value=null},Pe=async()=>{b.value&&(le(),await $())},Te=()=>{Object.assign(ae,{qrInput:"",scannedProduct:null,consumptionSuccess:null,error:null,lastScanContext:null,isQuickMode:!1}),X()},Ee=t=>{const e="p-4 rounded-lg border-2 text-sm font-medium transition-all";return T.value===t?`${e} border-red-500 bg-red-50 text-red-700`:`${e} border-gray-200 bg-white text-gray-700 hover:border-gray-300 hover:bg-gray-50`},qe=t=>{if(!t)return"No disponible";try{return O(new Date(t),"dd/MM/yyyy HH:mm",{locale:K})}catch{return t}},Re=async()=>{const t=new URLSearchParams(window.location.search),e=t.get("qr"),l=t.get("quick")==="true",c=t.get("purpose");if(e&&l&&c==="consume"){U.value=!0;const a=t.get("consumption_purpose");a&&["routine","emergency","maintenance"].includes(a)&&(T.value=a),b.value=e,await new Promise(g=>setTimeout(g,500));try{await $()}catch(g){console.error("Error en escaneo automático:",g)}const d=new URL(window.location);["qr","quick","purpose","consumption_purpose"].forEach(g=>d.searchParams.delete(g)),window.history.replaceState({},document.title,d.pathname+(d.searchParams.toString()?"?"+d.searchParams.toString():""))}};return Be(async()=>{if(De(),_.isAuthenticated||_.initializeAuth(),_.isAuthenticated&&_.token&&(!_.user||!_.user.name))try{await _.fetchProfile()}catch(t){console.error("Error obteniendo perfil:",t)}await Promise.all([G(),W(),ne()]),await Re()}),(t,e)=>{var l,c;return r(),u("div",We,[s("div",Xe,[s("div",Ze,[e[9]||(e[9]=s("div",null,[s("h2",{class:"text-xl font-semibold text-gray-900"},"Consumir Insumos Médicos"),s("p",{class:"text-gray-600 mt-1"},"Escanea o ingresa códigos QR para registrar el consumo de productos con trazabilidad completa")],-1)),s("div",et,[s("div",null,[e[7]||(e[7]=s("p",{class:"text-sm text-gray-500"},"Productos consumidos hoy",-1)),s("p",tt,y(C.today),1)]),s("div",null,[e[8]||(e[8]=s("p",{class:"text-sm text-gray-500"},"Total esta semana",-1)),s("p",st,y(C.week),1)])])]),s("div",at,[s("div",ot,[s("div",null,[e[10]||(e[10]=s("span",{class:"text-sm font-medium text-gray-700"},"Usuario:",-1)),s("span",nt,y(((l=S.value)==null?void 0:l.name)||"No identificado"),1)]),s("div",null,[e[11]||(e[11]=s("span",{class:"text-sm font-medium text-gray-700"},"Ubicación:",-1)),s("span",rt,y(((c=n(h))==null?void 0:c.name)||"No seleccionada"),1)])]),s("div",it," Sesión iniciada: "+y(qe(n(ye))),1)])]),!n(m)||n(U)?(r(),u("div",lt,[e[12]||(e[12]=s("h3",{class:"text-lg font-medium text-gray-900 mb-4"},"Tipo de Consumo",-1)),s("div",ct,[(r(),u(H,null,Q(be,a=>s("button",{key:a.value,onClick:d=>T.value=a.value,class:ue(Ee(a.value))},[s("div",ut,[(r(),u("svg",{class:ue(["h-6 w-6 mb-2",a.iconClass]),fill:"none",stroke:"currentColor",viewBox:"0 0 24 24"},[s("path",{"stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"2",d:a.icon},null,8,mt)],2)),s("span",pt,y(a.label),1),s("span",vt,y(a.description),1)])],10,dt)),64))])])):v("",!0),n(m)?v("",!0):(r(),u("div",gt,[e[15]||(e[15]=s("h3",{class:"text-lg font-medium text-gray-900 mb-4"},[s("svg",{class:"h-5 w-5 inline mr-2",fill:"none",stroke:"currentColor",viewBox:"0 0 24 24"},[s("path",{"stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"2",d:"M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h2M4 4h5l2 3h3l2-3h5v5M4 4v5m16-5v5"})]),w(" Escanear Código QR ")],-1)),s("div",yt,[s("div",null,[e[13]||(e[13]=s("label",{for:"qrInput",class:"block text-sm font-medium text-gray-700 mb-2"}," Código QR del Producto: ",-1)),s("div",ft,[D(s("input",{id:"qrInput","onUpdate:modelValue":e[0]||(e[0]=a=>Le(b)?b.value=a:null),type:"text",placeholder:"Ej: SUPPLY_1755580808_abc123def",class:"form-input flex-1",onKeyup:Ne($,["enter"]),disabled:n(P)},null,40,bt),[[ee,n(b)]]),s("button",{onClick:$,disabled:!n(b).trim()||n(P),class:"btn-primary"},[n(P)?(r(),E(p,{key:0})):v("",!0),w(" "+y(n(P)?"Escaneando...":"Escanear"),1)],8,ht)])]),s("div",xt,[s("button",{onClick:Ie,class:"btn-secondary w-full h-20 border-2 border-dashed border-gray-300 hover:border-gray-400",disabled:n(A)},[e[14]||(e[14]=s("svg",{class:"h-8 w-8 mx-auto mb-2",fill:"none",stroke:"currentColor",viewBox:"0 0 24 24"},[s("path",{"stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"2",d:"M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z"}),s("path",{"stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"2",d:"M15 13a3 3 0 11-6 0 3 3 0 016 0z"})],-1)),s("span",wt,y(n(A)?"Cámara Activa":"Usar Cámara"),1)],8,_t)])])])),n(m)&&!n(f)?(r(),u("div",kt,[n(U)?(r(),E(i,{key:0,"current-location":n(h)},null,8,["current-location"])):v("",!0),e[27]||(e[27]=s("h3",{class:"text-lg font-medium text-gray-900 mb-4"},[s("svg",{class:"h-5 w-5 inline mr-2 text-green-600",fill:"none",stroke:"currentColor",viewBox:"0 0 24 24"},[s("path",{"stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"2",d:"M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"})]),w(" Producto Encontrado - Registrar Consumo ")],-1)),de(I,{product:n(m),"scan-context":n(Y),"current-user":S.value,"current-location":n(h)},null,8,["product","scan-context","current-user","current-location"]),fe.value?(r(),u("form",{key:1,onSubmit:Ve(Ce,["prevent"])},[s("div",Ct,[s("div",null,[e[16]||(e[16]=s("label",{for:"userRUT",class:"block text-sm font-medium text-gray-700 mb-2"},[w(" RUT del Usuario "),s("span",{class:"text-red-500"},"*")],-1)),D(s("input",{id:"userRUT","onUpdate:modelValue":e[1]||(e[1]=a=>o.value.userRUT=a),type:"text",placeholder:"12.345.678-9",class:"form-input",required:""},null,512),[[ee,o.value.userRUT]])]),s("div",null,[e[18]||(e[18]=s("label",{for:"destinationType",class:"block text-sm font-medium text-gray-700 mb-2"},[w(" Tipo de Destino "),s("span",{class:"text-red-500"},"*")],-1)),D(s("select",{id:"destinationType","onUpdate:modelValue":e[2]||(e[2]=a=>o.value.destinationType=a),onChange:we,class:"form-select",required:""},[...e[17]||(e[17]=[s("option",{value:""},"Seleccionar tipo",-1),s("option",{value:"pavilion"},"Pabellón",-1),s("option",{value:"store"},"Almacén",-1)])],544),[[F,o.value.destinationType]])]),o.value.destinationType==="pavilion"?(r(),u("div",It,[e[20]||(e[20]=s("label",{for:"medicalCenter",class:"block text-sm font-medium text-gray-700 mb-2"},[w(" Centro Médico "),s("span",{class:"text-red-500"},"*")],-1)),D(s("select",{id:"medicalCenter","onUpdate:modelValue":e[3]||(e[3]=a=>o.value.medicalCenterId=a),onChange:ke,class:"form-select",disabled:n(N),required:""},[e[19]||(e[19]=s("option",{value:""},"Seleccionar centro médico",-1)),(r(!0),u(H,null,Q(n(R),a=>(r(),u("option",{key:a.id,value:a.id},y(a.name),9,Mt))),128))],40,St),[[F,o.value.medicalCenterId]]),n(N)?(r(),u("p",Dt,"Cargando centros médicos...")):v("",!0)])):v("",!0),o.value.destinationType==="pavilion"&&o.value.medicalCenterId?(r(),u("div",Pt,[e[22]||(e[22]=s("label",{for:"pavilion",class:"block text-sm font-medium text-gray-700 mb-2"},[w(" Pabellón "),s("span",{class:"text-red-500"},"*")],-1)),D(s("select",{id:"pavilion","onUpdate:modelValue":e[4]||(e[4]=a=>o.value.destinationID=a),class:"form-select",disabled:n(L),required:""},[e[21]||(e[21]=s("option",{value:""},"Seleccionar pabellón",-1)),(r(!0),u(H,null,Q(oe.value,a=>(r(),u("option",{key:a.id,value:a.id},y(a.name),9,Et))),128))],8,Tt),[[F,o.value.destinationID]]),n(L)?(r(),u("p",qt,"Cargando pabellones...")):v("",!0)])):v("",!0),o.value.destinationType==="store"?(r(),u("div",Rt,[e[24]||(e[24]=s("label",{for:"store",class:"block text-sm font-medium text-gray-700 mb-2"},[w(" Almacén "),s("span",{class:"text-red-500"},"*")],-1)),D(s("select",{id:"store","onUpdate:modelValue":e[5]||(e[5]=a=>o.value.destinationID=a),class:"form-select",disabled:n(V),required:""},[e[23]||(e[23]=s("option",{value:""},"Seleccionar almacén",-1)),(r(!0),u(H,null,Q(n(B),a=>(r(),u("option",{key:a.id,value:a.id},y(a.name)+" - "+y(a.location),9,Ut))),128))],8,At),[[F,o.value.destinationID]]),n(V)?(r(),u("p",zt,"Cargando almacenes...")):v("",!0)])):v("",!0),s("div",null,[e[25]||(e[25]=s("label",{for:"notes",class:"block text-sm font-medium text-gray-700 mb-2"}," Notas (Opcional) ",-1)),D(s("input",{id:"notes","onUpdate:modelValue":e[6]||(e[6]=a=>o.value.notes=a),type:"text",placeholder:"Observaciones adicionales",class:"form-input"},null,512),[[ee,o.value.notes]])])]),s("div",Bt,[s("button",{type:"button",onClick:X,class:"btn-secondary"}," Cancelar "),s("button",{type:"submit",disabled:n(q),class:"btn-danger"},[n(q)?(r(),E(p,{key:0})):(r(),u("svg",Lt,[...e[26]||(e[26]=[s("path",{"stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"2",d:"M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"},null,-1)])])),w(" "+y(n(q)?"Consumiendo...":"Confirmar Consumo"),1)],8,Nt)])],32)):v("",!0)])):v("",!0),n(J)?(r(),E(pe,{key:3,"success-data":n(J),onClearAll:Te,onViewTraceability:re,onViewBatchHistory:ie},null,8,["success-data"])):v("",!0),n(f)?(r(),E(ve,{key:4,"error-text":n(f),"is-quick-mode":n(U),"qr-input":n(b),onClearError:le,onRetryScan:Pe},null,8,["error-text","is-quick-mode","qr-input"])):v("",!0),n(k).length>0?(r(),E(ge,{key:5,consumptions:n(k),onViewTraceability:re,onViewBatchHistory:ie},null,8,["consumptions"])):v("",!0),de(Qe,{onLocationChanged:Se})])}}},Jt=je(Vt,[["__scopeId","data-v-23f17271"]]);export{Jt as default};
