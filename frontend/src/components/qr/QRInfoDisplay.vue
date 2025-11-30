<template>
  <div class="space-y-4">
    <!-- Header principal con estado destacado -->
    <div class="bg-white rounded-lg shadow border overflow-hidden">
      <!-- Estado y nombre del insumo (Lo más importante arriba) -->
      <div class="px-4 sm:px-6 py-4 border-b border-gray-200">
        <div class="flex items-start justify-between gap-4">
          <!-- Información principal -->
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-3 mb-2">
              <h3 class="text-xl sm:text-2xl font-bold text-gray-900 truncate">
                {{ qrInfo.supply_code?.name || qrInfo.supply_info?.name || qrInfo.name || 'Insumo Médico' }}
              </h3>
            </div>
            
            <!-- Estado destacado -->
            <span :class="getStatusBadgeClass()" class="inline-flex px-4 py-2 text-base font-bold rounded-lg shadow-sm mb-3">
              {{ getStatusLabel() }}
            </span>
            
            <!-- Información clave en una línea -->
            <div class="flex flex-wrap items-center gap-3 text-sm text-gray-700">
              <div class="flex items-center gap-1 bg-white px-2 py-1 rounded">
                <svg class="h-4 w-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
                </svg>
                <span class="font-mono text-xs">{{ qrInfo.qr_code }}</span>
              </div>
              <div v-if="scanActivity" class="flex items-center gap-1 bg-white px-2 py-1 rounded">
                <svg class="h-4 w-4 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                </svg>
                <span>{{ scanActivity.total_scans }} escaneos</span>
              </div>
            </div>
          </div>
          
          <!-- QR Code miniatura -->
          <div class="flex-shrink-0 w-16 h-16 sm:w-20 sm:h-20 bg-white border-2 border-gray-300 rounded-lg p-1 shadow-sm">
            <img v-if="qrInfo.qr_image_url" :src="qrInfo.qr_image_url" alt="QR" class="w-full h-full object-contain" />
            <qrcode-vue v-else :value="qrInfo.qr_code" :size="64" :level="'M'" class="w-full h-full" />
          </div>
        </div>
      </div>

      <!-- Acciones principales (botones grandes y visibles) -->
      <div v-if="qrInfo.type !== 'batch'" class="px-4 sm:px-6 py-3 border-b border-gray-200">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-2">
          <router-link
            :to="`/qr/${qrInfo.qr_code}/traceability`"
            class="inline-flex items-center justify-center px-4 py-2.5 border-2 border-blue-600 text-base font-semibold rounded-lg text-blue-700 bg-white hover:bg-blue-50 transition-colors shadow-sm"
          >
            <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
            Ver Trazabilidad
          </router-link>
          <button
            @click="downloadQRAsPDF"
            :disabled="isGeneratingPDF"
            class="inline-flex items-center justify-center px-4 py-2.5 border-2 border-green-600 text-base font-semibold rounded-lg text-green-700 bg-white hover:bg-green-50 transition-colors disabled:opacity-50 disabled:cursor-not-allowed shadow-sm"
          >
            <svg v-if="isGeneratingPDF" class="animate-spin h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <svg v-else class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            {{ isGeneratingPDF ? 'Generando...' : 'Descargar PDF' }}
          </button>
        </div>
      </div>

      <!-- Contenido del QR -->
      <div class="px-4 sm:px-6 py-4">
        <!-- Información esencial del insumo (Grid compacto) -->
        <div v-if="qrInfo.supply_info || qrInfo.batch_info" class="mb-4 grid grid-cols-2 sm:grid-cols-4 gap-3 text-sm">
          <div v-if="qrInfo.supply_code?.code_supplier || qrInfo.supplier_code" class="border border-gray-200 p-3 rounded-lg">
            <div class="text-xs text-gray-500 mb-1">Código Proveedor</div>
            <div class="font-semibold text-gray-900">{{ qrInfo.supply_code?.code_supplier || qrInfo.supplier_code }}</div>
          </div>
          
          <div v-if="qrInfo.supply_info?.batch?.id || qrInfo.batch_info?.id || qrInfo.batch_id" class="border border-gray-200 p-3 rounded-lg">
            <div class="text-xs text-gray-500 mb-1">Lote</div>
            <div class="font-semibold text-gray-900">{{ qrInfo.supply_info?.batch?.id || qrInfo.batch_info?.id || qrInfo.batch_id }}</div>
          </div>
          
          <div v-if="qrInfo.batch_info?.supplier || qrInfo.supplier" class="border border-gray-200 p-3 rounded-lg">
            <div class="text-xs text-gray-500 mb-1">Proveedor</div>
            <div class="font-semibold text-gray-900">{{ qrInfo.batch_info?.supplier || qrInfo.supplier }}</div>
          </div>
          
          <div v-if="qrInfo.batch_info?.expiration_date || qrInfo.expiration_date" class="border border-gray-200 p-3 rounded-lg">
            <div class="text-xs text-gray-500 mb-1">Vencimiento</div>
            <div class="font-semibold text-gray-900">{{ formatDate(qrInfo.batch_info?.expiration_date || qrInfo.expiration_date) }}</div>
          </div>
        </div>

        <!-- Información de asignación a solicitud (más compacta) -->
        <div v-if="qrInfo.request_assignment" class="mb-4 bg-blue-50 border-l-4 border-blue-600 rounded-r-lg p-4">
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center">
              <svg class="h-5 w-5 text-blue-600 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
              </svg>
              <h4 class="font-semibold text-blue-900 text-base">Asignado a Solicitud</h4>
            </div>
            <span :class="getAssignmentStatusBadgeClass(qrInfo.request_assignment.status)" class="px-3 py-1 text-xs font-semibold rounded-full">
              {{ getAssignmentStatusLabel(qrInfo.request_assignment.status) }}
            </span>
          </div>
          
          <div class="space-y-2 text-sm">
            <div class="flex items-baseline">
              <span class="text-blue-700 font-medium mr-2 min-w-[120px]">Solicitud:</span>
              <router-link
                v-if="qrInfo.supply_request"
                :to="`/supply-requests/${qrInfo.supply_request.id}`"
                class="text-blue-900 hover:text-blue-700 underline font-semibold"
              >
                {{ qrInfo.supply_request.request_number }}
              </router-link>
              <span v-else class="text-blue-900">N/A</span>
            </div>
            
            <div class="flex items-baseline">
              <span class="text-blue-700 font-medium mr-2 min-w-[120px]">Asignado:</span>
              <span class="text-blue-900">{{ formatDate(qrInfo.request_assignment.assigned_date) }}</span>
            </div>
            
            <div v-if="qrInfo.request_assignment.delivered_date" class="flex items-baseline">
              <span class="text-blue-700 font-medium mr-2 min-w-[120px]">Entregado:</span>
              <span class="text-blue-900">{{ formatDate(qrInfo.request_assignment.delivered_date) }}</span>
            </div>
            
            <div v-if="qrInfo.request_assignment.notes" class="pt-2 border-t border-blue-200">
              <span class="text-blue-700 font-medium">Notas:</span>
              <p class="text-blue-900 mt-1 bg-white p-2 rounded text-xs">{{ qrInfo.request_assignment.notes }}</p>
            </div>
          </div>

          <!-- Carrito de Insumos (solo si está activo) - Colapsado por defecto -->
          <div v-if="qrInfo.request_assignment.cart && qrInfo.request_assignment.cart.status === 'active' && qrInfo.qr_code" class="mt-4 pt-3 border-t border-blue-200">
            <SupplyCart 
              :qr-code="qrInfo.qr_code"
              :can-close="false"
              :can-remove-items="false"
              :can-manage-items="true"
              :show-add-button="false"
              :default-collapsed="true"
              @cart-loaded="onCartLoaded"
              @cart-closed="onCartClosed"
              @error="onCartError"
            />
          </div>
        </div>

        <!-- Último escaneo (información compacta) -->
        <div v-if="lastScanInfo && qrInfo.type !== 'batch'" class="mb-4 p-3 border border-gray-200 rounded-lg text-sm">
          <div class="flex flex-wrap items-center gap-x-4 gap-y-2">
            <div class="flex items-center text-gray-600">
              <svg class="h-4 w-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
              </svg>
              <span class="font-medium">{{ lastScanInfo.user_name || 'Usuario no identificado' }}</span>
            </div>
            <div class="flex items-center text-gray-500">
              <svg class="h-4 w-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <span>{{ formatDate(lastScanInfo.scanned_at) }}</span>
            </div>
            <div v-if="lastScanInfo.location" class="flex items-center text-gray-500">
              <svg class="h-4 w-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
              </svg>
              <span>{{ lastScanInfo.location }}</span>
            </div>
          </div>
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
import qrService from '@/services/qr/qrService'
import { useQRPdfDownload } from '@/composables/useQRPdfDownload'
import SupplyCart from '@/components/requests/SupplyCart.vue'

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
  console.log('Carrito cargado en QRInfoDisplay:', cart)
}

const onCartClosed = (cart) => {
  console.log('Carrito cerrado en QRInfoDisplay:', cart)
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
      case 'pendiente_retiro':
        return 'Pendiente de Retiro'
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
      case 'pendiente_retiro':
        return 'text-yellow-600'
      case 'recepcionado':
        return 'text-blue-600'
      case 'en_camino_a_pabellon':
        return 'text-orange-600'
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