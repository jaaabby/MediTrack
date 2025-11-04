<template>
  <div class="space-y-6">
    <!-- Información principal del QR con indicadores de trazabilidad -->
    <div class="bg-white rounded-lg shadow border overflow-hidden">
      <!-- Header reorganizado con imagen QR y datos principales -->
      <div class="px-6 py-4 bg-gray-50 border-b border-gray-200 flex flex-col md:flex-row md:items-center md:justify-between">
        <div class="flex items-center space-x-4">
          <!-- Imagen del QR: usa la URL si existe, si no genera el QR dinámicamente -->
          <div class="w-20 h-20 flex items-center justify-center bg-white border rounded-lg mr-4">
            <img v-if="qrInfo.qr_image_url" :src="qrInfo.qr_image_url" alt="QR" class="w-16 h-16 object-contain" />
            <qrcode-vue v-else :value="qrInfo.qr_code" :size="64" :level="'M'" />
          </div>
          <div>
            <h3 class="text-lg font-semibold text-gray-900">{{ getTypeLabel(qrInfo.type) }}</h3>
            <p class="text-sm text-gray-600">Código: {{ qrInfo.qr_code }}</p>
            <span :class="getStatusBadgeClass()" class="inline-flex px-3 py-1 text-sm font-semibold rounded-full mt-2">
              {{ getStatusLabel() }}
            </span>
          </div>
        </div>
        <div v-if="qrInfo.type !== 'batch'" class="flex flex-col md:flex-row md:items-center md:space-x-2 mt-4 md:mt-0">
          <button
            @click="downloadQRAsPDF"
            :disabled="isGeneratingPDF"
            class="inline-flex items-center px-3 py-2 border border-green-300 text-sm font-medium rounded-md text-green-700 bg-green-50 hover:bg-green-100 transition-colors disabled:opacity-50 disabled:cursor-not-allowed mb-2 md:mb-0"
          >
            <svg v-if="isGeneratingPDF" class="animate-spin h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <svg v-else class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            {{ isGeneratingPDF ? 'Generando...' : 'Descargar PDF' }}
          </button>
          <router-link
            :to="`/qr/${qrInfo.qr_code}/traceability`"
            class="inline-flex items-center px-3 py-2 border border-blue-300 text-sm font-medium rounded-md text-blue-700 bg-blue-50 hover:bg-blue-100 transition-colors"
          >
            <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
            Ver Trazabilidad
          </router-link>
          <div v-if="scanActivity" class="flex items-center space-x-1 text-xs text-gray-500 bg-gray-100 px-2 py-1 rounded ml-2">
            <svg class="h-3 w-3 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
            <span>{{ scanActivity.total_scans }} escaneos</span>
          </div>
        </div>
      </div>
      <!-- Último escaneo -->
      <div v-if="lastScanInfo && qrInfo.type !== 'batch'" class="px-6 py-2 border-b border-gray-200 bg-white">
        <div class="flex flex-wrap items-center text-sm gap-4">
          <div class="flex items-center text-gray-600">
            <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
            </svg>
            <span>Último escaneo: {{ lastScanInfo.user_name || 'Usuario no identificado' }}</span>
          </div>
          <div class="flex items-center text-gray-600">
            <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span>{{ formatDate(lastScanInfo.scanned_at) }}</span>
          </div>
          <div v-if="lastScanInfo.location" class="flex items-center text-gray-600">
            <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
            </svg>
            <span>{{ lastScanInfo.location }}</span>
          </div>
          <span v-if="lastScanInfo.scan_purpose" :class="['px-2 py-1 rounded-full text-xs font-medium', getScanPurposeBadgeClass(lastScanInfo.scan_purpose)]">
            {{ getScanPurposeLabel(lastScanInfo.scan_purpose) }}
          </span>
        </div>
      </div>

      <!-- Contenido del QR -->
      <div class="p-6">
        <!-- Información de asignación a solicitud -->
        <div v-if="qrInfo.request_assignment" class="mb-6 bg-blue-50 border border-blue-200 rounded-lg p-4">
          <div class="flex items-center mb-3">
            <svg class="h-5 w-5 text-blue-600 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
            </svg>
            <h4 class="font-semibold text-blue-900">Asignado a Solicitud</h4>
          </div>
          
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="text-sm font-medium text-blue-700">Número de Solicitud:</label>
              <router-link
                v-if="qrInfo.supply_request"
                :to="`/supply-requests/${qrInfo.supply_request.id}`"
                class="block text-sm text-blue-900 hover:text-blue-700 underline"
              >
                {{ qrInfo.supply_request.request_number }}
              </router-link>
              <p v-else class="text-sm text-blue-900">N/A</p>
            </div>
            
            <div>
              <label class="text-sm font-medium text-blue-700">Estado de Asignación:</label>
              <span :class="getAssignmentStatusBadgeClass(qrInfo.request_assignment.status)" class="inline-flex px-2 py-1 text-xs font-semibold rounded-full ml-1">
                {{ getAssignmentStatusLabel(qrInfo.request_assignment.status) }}
              </span>
            </div>
            
            <div>
              <label class="text-sm font-medium text-blue-700">Fecha de Asignación:</label>
              <p class="text-sm text-blue-900">{{ formatDate(qrInfo.request_assignment.assigned_date) }}</p>
            </div>
            
            <div v-if="qrInfo.request_assignment.delivered_date">
              <label class="text-sm font-medium text-blue-700">Fecha de Entrega:</label>
              <p class="text-sm text-blue-900">{{ formatDate(qrInfo.request_assignment.delivered_date) }}</p>
            </div>
          </div>
          
          <div v-if="qrInfo.request_assignment.notes" class="mt-3">
            <label class="text-sm font-medium text-blue-700">Notas:</label>
            <p class="text-sm text-blue-900 bg-white p-2 rounded border mt-1">{{ qrInfo.request_assignment.notes }}</p>
          </div>

          <!-- Información del Carrito -->
          <div v-if="qrInfo.request_assignment.cart" class="mt-4 pt-4 border-t border-blue-200">
            <div class="flex items-center mb-2">
              <svg class="h-5 w-5 text-blue-600 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z" />
              </svg>
              <h5 class="font-semibold text-blue-900">Carrito Asociado</h5>
            </div>
            <div class="bg-white rounded p-3 border border-blue-100">
              <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                <div>
                  <label class="text-xs font-medium text-blue-700">Número de Carrito:</label>
                  <p class="text-sm text-blue-900 font-mono">{{ qrInfo.request_assignment.cart.cart_number }}</p>
                </div>
                <div>
                  <label class="text-xs font-medium text-blue-700">Estado:</label>
                  <span :class="getCartStatusClass(qrInfo.request_assignment.cart.status)" class="inline-flex px-2 py-1 text-xs font-semibold rounded-full ml-1">
                    {{ getCartStatusLabel(qrInfo.request_assignment.cart.status) }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Carrito de Insumos (si existe) -->
        <div v-if="qrInfo.request_assignment && qrInfo.qr_code">
          <SupplyCart 
            :qr-code="qrInfo.qr_code"
            :can-close="false"
            :can-remove-items="false"
            :can-manage-items="true"
            :show-add-button="false"
            @cart-loaded="onCartLoaded"
            @error="onCartError"
          />
        </div>

        <!-- Estadísticas de escaneo básicas 
        <div v-if="showTraceability && scanActivity" class="mb-6 bg-gradient-to-r from-purple-50 to-indigo-50 border border-purple-200 rounded-lg p-4">
          <div class="flex items-center mb-3">
            <svg class="h-5 w-5 text-purple-600 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
            <h4 class="font-semibold text-purple-900">Actividad de Escaneo</h4>
          </div>
          
          <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div class="text-center">
              <div class="text-2xl font-bold text-purple-900">{{ scanActivity.total_scans }}</div>
              <div class="text-xs text-purple-700">Total Escaneos</div>
            </div>
            <div class="text-center">
              <div class="text-2xl font-bold text-indigo-900">{{ scanActivity.unique_users || 0 }}</div>
              <div class="text-xs text-indigo-700">Usuarios Únicos</div>
            </div>
            <div class="text-center">
              <div class="text-2xl font-bold text-blue-900">{{ scanActivity.locations_count || 0 }}</div>
              <div class="text-xs text-blue-700">Ubicaciones</div>
            </div>
            <div class="text-center">
              <div class="text-2xl font-bold text-gray-900">{{ scanActivity.today_scans || 0 }}</div>
              <div class="text-xs text-gray-700">Hoy</div>
            </div>
          </div>
        </div>-->

        <!-- Información de contexto de escaneo 
        <div v-if="showTraceability && scanContext" class="mb-6 bg-yellow-50 border border-yellow-200 rounded-lg p-4">
          <div class="flex items-center mb-3">
            <svg class="h-5 w-5 text-yellow-600 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <h4 class="font-semibold text-yellow-900">Contexto del Escaneo Actual</h4>
          </div>
          
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
            <div>
              <label class="font-medium text-yellow-700">Propósito:</label>
              <span 
                :class="[
                  'inline-flex px-2 py-1 rounded-full text-xs font-medium ml-1',
                  getScanPurposeBadgeClass(scanContext.scan_purpose)
                ]"
              >
                {{ getScanPurposeLabel(scanContext.scan_purpose) }}
              </span>
            </div>
            <div v-if="scanContext.scan_source">
              <label class="font-medium text-yellow-700">Fuente:</label>
              <p class="text-yellow-900">{{ scanContext.scan_source }}</p>
            </div>
            <div v-if="scanContext.device_info">
              <label class="font-medium text-yellow-700">Dispositivo:</label>
              <p class="text-yellow-900">{{ scanContext.device_info.platform || 'Desconocido' }}</p>
            </div>
          </div>
        </div>-->

        <!-- Resto del contenido existente -->
        <!-- Información resumida y no redundante del insumo y lote -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div v-if="qrInfo.supply_info" class="space-y-2">
            <h4 class="font-semibold text-gray-900 flex items-center">
              <svg class="h-5 w-5 mr-2 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              Información del Insumo
            </h4>
            <div class="text-sm">
              <div><span class="font-medium text-gray-600">Código:</span> {{ qrInfo.supply_code?.code_supplier || qrInfo.supplier_code || 'N/A' }}</div>
              <div><span class="font-medium text-gray-600">Estado:</span> 
                <span :class="getStatusColor(qrInfo)" class="font-medium">{{ getStatusText(qrInfo) }}</span>
              </div>
              <div v-if="qrInfo.supply_info?.batch || qrInfo.batch_info"><span class="font-medium text-gray-600">Lote:</span> {{ qrInfo.supply_info?.batch?.id || qrInfo.batch_info?.id || qrInfo.batch_id || 'N/A' }}</div>
              <div v-if="qrInfo.supply_code || qrInfo.supply_info?.name"><span class="font-medium text-gray-600">Nombre:</span> {{ qrInfo.supply_code?.name || qrInfo.supply_info?.name || qrInfo.name || 'N/A' }}</div>
            </div>
          </div>
          <div v-if="qrInfo.batch_info" class="space-y-2">
            <h4 class="font-semibold text-gray-900 flex items-center">
              <svg class="h-5 w-5 mr-2 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
              </svg>
              Información del Lote
            </h4>
            <div class="text-sm">
              <div><span class="font-medium text-gray-600">ID:</span> {{ qrInfo.batch_info?.id || qrInfo.batch_id || 'N/A' }}</div>
              <div><span class="font-medium text-gray-600">Proveedor:</span> {{ qrInfo.batch_info?.supplier || qrInfo.supplier || 'N/A' }}</div>
              <div><span class="font-medium text-gray-600">Cantidad Total:</span> {{ qrInfo.batch_info?.amount || qrInfo.amount || 'N/A' }}</div>
              <div><span class="font-medium text-gray-600">Vencimiento:</span> {{ formatDate(qrInfo.batch_info?.expiration_date || qrInfo.expiration_date) }}</div>
            </div>
          </div>
        </div>

      </div>
    </div>

    <!-- Historial de movimientos (versión resumida) -->
    <div v-if="qrInfo.history && qrInfo.history.length > 0 && qrInfo.type !== 'batch'" class="bg-white rounded-lg shadow border">
      <div class="px-6 py-4 border-b border-gray-200 flex justify-between items-center">
        <h4 class="font-semibold text-gray-900">Historial Reciente</h4>
        <router-link
          :to="`/qr/${qrInfo.qr_code}/traceability`"
          class="text-sm text-blue-600 hover:text-blue-800"
        >
          Ver historial completo →
        </router-link>
      </div>
      
      <div class="p-6">
        <div class="flow-root">
          <ul class="-mb-8">
            <li
              v-for="(movement, index) in qrInfo.history.slice(0, 3)"
              :key="movement.id || index"
              class="relative pb-8"
            >
              <span
                v-if="index !== Math.min(qrInfo.history.length - 1, 2)"
                class="absolute top-4 left-4 -ml-px h-full w-0.5 bg-gray-200"
                aria-hidden="true"
              ></span>
              
              <div class="relative flex space-x-3">
                <div>
                  <span :class="getHistoryIconClass(movement.status || movement.Status)" class="h-8 w-8 rounded-full flex items-center justify-center ring-8 ring-white">
                    <svg class="h-4 w-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                  </span>
                </div>
                <div class="min-w-0 flex-1 pt-1.5">
                  <div>
                    <p class="text-sm text-gray-900">
                      <span class="font-medium">{{ movement.status || movement.Status }}</span>
                    </p>
                    <p class="text-sm text-gray-500">{{ formatDate(movement.date_time || movement.DateTime) }}</p>
                  </div>
                </div>
              </div>
            </li>
          </ul>
        </div>
        
        <div v-if="qrInfo.history.length > 3" class="mt-4 text-center">
          <p class="text-sm text-gray-500">{{ qrInfo.history.length - 3 }} movimientos adicionales</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import QrcodeVue from 'qrcode.vue'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import qrService from '@/services/qrService'
import { useQRPdfDownload } from '@/composables/useQRPdfDownload'
import SupplyCart from '@/components/SupplyCart.vue'

// Props
const props = defineProps({
  qrInfo: {
    type: Object,
    required: true
  },
  showTraceability: {
    type: Boolean,
    default: false
  },
  scanContext: {
    type: Object,
    default: null
  }
})

// Emits
const emit = defineEmits(['consume-supply', 'mark-as-delivered'])

// Registrar el componente QR
defineExpose({ QrcodeVue })

// Estado reactivo para información de trazabilidad
const scanActivity = ref(null)
const lastScanInfo = ref(null)

// Composable para descarga de PDF
const { downloadQRAsPDF: downloadPDF, isGenerating: isGeneratingPDF, error: pdfError } = useQRPdfDownload()

// Función para descargar QR como PDF
const downloadQRAsPDF = async () => {
  try {
    const success = await downloadPDF(props.qrInfo, {
      filename: `QR_${props.qrInfo.qr_code}_${new Date().getTime()}.pdf`,
      includeInfo: true
    })
    
    if (!success && pdfError.value) {
      console.error('Error al generar PDF:', pdfError.value)
      // Aquí podrías mostrar una notificación de error al usuario
    }
  } catch (error) {
    console.error('Error al descargar PDF:', error)
  }
}

// Métodos para el carrito
const onCartLoaded = (cart) => {
  console.log('Carrito cargado en vista QR:', cart)
}

const onCartError = (error) => {
  console.log('No se encontró carrito para este QR (puede ser normal):', error)
}

// Cargar información de trazabilidad si está habilitada
onMounted(async () => {
  if (props.showTraceability && props.qrInfo?.qr_code) {
    try {
      // Cargar estadísticas básicas de escaneo
      const stats = await qrService.getScanStatistics(props.qrInfo.qr_code)
      scanActivity.value = stats
      
      // Cargar información del último escaneo significativo
      const scanHistory = await qrService.getScanHistory(props.qrInfo.qr_code, 10)
      if (scanHistory && scanHistory.length > 0) {
        // Filtrar solo eventos que representen cambios reales, no verificaciones
        const significantEvents = scanHistory.filter(event => 
          !event.scan_purpose || 
          !['transfer_verification', 'transfer_check', 'lookup'].includes(event.scan_purpose)
        )
        lastScanInfo.value = significantEvents[0] || scanHistory[0]
      }
    } catch (error) {
      console.warn('No se pudo cargar información de trazabilidad:', error)
    }
  }
})

// Métodos de formato
const formatDate = (dateString) => {
  if (!dateString) return 'No disponible'
  try {
    return format(new Date(dateString), 'dd/MM/yyyy HH:mm', { locale: es })
  } catch (error) {
    return dateString
  }
}

const getTypeLabel = (type) => {
  const labels = {
    'batch': 'Lote de Productos',
    'medical_supply': 'Insumo Individual'
  }
  return labels[type] || type
}

const getStatusLabel = () => {
  return qrService.getStatusLabel(props.qrInfo)
}

const getStatusBadgeClass = () => {
  const color = qrService.getStatusColor(props.qrInfo)
  const classes = {
    'green': 'bg-green-100 text-green-800',
    'red': 'bg-red-100 text-red-800',
    'yellow': 'bg-yellow-100 text-yellow-800',
    'gray': 'bg-gray-100 text-gray-800'
  }
  return classes[color] || classes.gray
}

const getTypeIconClass = (type) => {
  const classes = {
    'batch': 'bg-blue-100 text-blue-600',
    'medical_supply': 'bg-green-100 text-green-600'
  }
  return classes[type] || 'bg-gray-100 text-gray-600'
}

const getHistoryIconClass = (status) => {
  const classes = {
    'consumido': 'bg-red-500',
    'creado': 'bg-green-500',
    'movido': 'bg-blue-500'
  }
  return classes[status] || 'bg-gray-500'
}

// Métodos para solicitudes (manteniendo los originales)
const getAssignmentStatusLabel = (status) => {
  const labels = {
    'assigned': 'Asignado',
    'delivered': 'Entregado',
    'consumed': 'Consumido',
    'returned': 'Devuelto'
  }
  return labels[status] || status
}

const getAssignmentStatusBadgeClass = (status) => {
  const classes = {
    'assigned': 'bg-blue-100 text-blue-800',
    'delivered': 'bg-green-100 text-green-800',
    'consumed': 'bg-gray-100 text-gray-800',
    'returned': 'bg-yellow-100 text-yellow-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

// Nuevos métodos para trazabilidad
const getScanPurposeBadgeClass = (purpose) => {
  const classes = {
    'lookup': 'bg-green-100 text-green-800',
    'consume': 'bg-red-100 text-red-800',
    'verify': 'bg-blue-100 text-blue-800',
    'inventory_check': 'bg-purple-100 text-purple-800',
    'assign': 'bg-yellow-100 text-yellow-800'
  }
  return classes[purpose] || 'bg-gray-100 text-gray-800'
}

const getScanPurposeLabel = (purpose) => {
  const labels = {
    'lookup': 'Consulta',
    'consume': 'Consumo',
    'verify': 'Verificación',
    'inventory_check': 'Inventario',
    'assign': 'Asignación'
  }
  return labels[purpose] || purpose
}

const getCartStatusLabel = (status) => {
  const labels = {
    'active': 'Activo',
    'closed': 'Cerrado',
    'cancelled': 'Cancelado'
  }
  return labels[status] || status
}

const getCartStatusClass = (status) => {
  const classes = {
    'active': 'bg-green-100 text-green-800',
    'closed': 'bg-gray-100 text-gray-800',
    'cancelled': 'bg-red-100 text-red-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

// Funciones para determinar el estado correcto del insumo
const getStatusText = (qrInfo) => {
  // Prioridad: consumido > estado del supply_info > estado directo > disponible por defecto
  if (qrInfo.supply_info?.is_consumed || qrInfo.is_consumed) {
    return 'Consumido'
  }
  
  // Revisar el estado desde supply_info (que viene del backend)
  const status = qrInfo.supply_info?.Status || qrInfo.supply_info?.status || qrInfo.status
  
  if (status) {
    switch (status.toLowerCase()) {
      case 'disponible':
        return 'Disponible'
      case 'recepcionado':
        return 'Recepcionado'
      case 'en_camino_a_pabellon':
        return 'En camino a pabellón'
      case 'en_camino_a_bodega':
        return 'En camino a bodega'
      case 'consumido':
        return 'Consumido'
      default:
        return status
    }
  }
  
  return 'Disponible' // fallback
}

const getStatusColor = (qrInfo) => {
  // Determinar color basado en el estado
  if (qrInfo.supply_info?.is_consumed || qrInfo.is_consumed) {
    return 'text-red-600'
  }
  
  const status = qrInfo.supply_info?.Status || qrInfo.supply_info?.status || qrInfo.status
  
  if (status) {
    switch (status.toLowerCase()) {
      case 'disponible':
        return 'text-green-600'
      case 'recepcionado':
        return 'text-blue-600'
      case 'en_camino_a_pabellon':
        return 'text-yellow-600'
      case 'en_camino_a_bodega':
        return 'text-orange-600'
      case 'consumido':
        return 'text-red-600'
      default:
        return 'text-gray-600'
    }
  }
  
  return 'text-green-600' // fallback para disponible
}
</script>

<style scoped>
/* Estilos específicos */
.flow-root {
  overflow: hidden;
}

.transition-colors {
  transition: color 0.2s ease-in-out, background-color 0.2s ease-in-out;
}

/* Ring styles for timeline */
.ring-8 {
  box-shadow: 0 0 0 8px #fff;
}
</style>