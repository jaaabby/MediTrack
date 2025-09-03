<template>
  <div class="max-w-7xl mx-auto p-6">
    <!-- Header -->
    <div class="mb-6">
      <div class="flex items-center mb-4">
        <button
          @click="$router.go(-1)"
          class="mr-4 p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-full"
        >
          <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </button>
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Trazabilidad QR</h1>
          <p class="text-gray-600 mt-1">Seguimiento completo del código {{ qrCode }}</p>
        </div>
      </div>
      
      <!-- Código QR destacado -->
      <div class="bg-gradient-to-r from-blue-50 to-indigo-50 border border-blue-200 rounded-lg p-4">
        <div class="flex items-center justify-between">
          <div class="flex items-center">
            <div class="bg-white p-2 rounded border border-blue-300">
              <svg class="h-8 w-8 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v2a2 2 0 002 2zm0 0h2a2 2 0 012 2v2a2 2 0 01-2 2H5a2 2 0 01-2-2v-2a2 2 0 012-2z" />
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-700">Código QR</p>
              <code class="text-lg font-mono text-blue-900 bg-white px-2 py-1 rounded border">{{ qrCode }}</code>
            </div>
          </div>
          <button
            @click="loadTraceability"
            :disabled="loading"
            class="inline-flex items-center px-3 py-2 border border-blue-300 text-sm font-medium rounded-md text-blue-700 bg-white hover:bg-blue-50 disabled:opacity-50"
          >
            <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            Actualizar
          </button>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="text-center py-12">
      <div class="animate-spin rounded-full h-16 w-16 border-b-2 border-blue-600 mx-auto mb-4"></div>
      <p class="text-gray-600">Cargando trazabilidad...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="text-center py-12">
      <svg class="h-16 w-16 text-red-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <h3 class="text-lg font-medium text-gray-900 mb-2">Error al cargar trazabilidad</h3>
      <p class="text-gray-600 mb-4">{{ error }}</p>
      <button
        @click="loadTraceability"
        class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700"
      >
        Reintentar
      </button>
    </div>

    <!-- Contenido principal -->
    <div v-else-if="traceability" class="space-y-6">
      <!-- Información general del insumo -->
      <div v-if="qrInfo" class="bg-white rounded-lg shadow border p-6">
        <h2 class="text-lg font-semibold text-gray-900 mb-4">Información del Insumo</h2>
        
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <!-- Información básica -->
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700">Tipo</label>
              <span :class="getTypeIconClass(qrInfo.type)" class="inline-flex items-center px-2 py-1 text-sm font-medium rounded-full mt-1">
                <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path v-if="qrInfo.type === 'medical_supply'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                </svg>
                {{ getTypeLabel(qrInfo.type) }}
              </span>
            </div>
            
            <div v-if="qrInfo.supply_info">
              <label class="block text-sm font-medium text-gray-700">Estado Actual</label>
              <span :class="getStatusBadgeClass()" class="inline-flex px-2 py-1 text-sm font-semibold rounded-full mt-1">
                {{ getStatusLabel() }}
              </span>
            </div>
          </div>

          <!-- Información del código de insumo -->
          <div v-if="qrInfo.supply_code" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700">Código Insumo</label>
              <p class="text-sm text-gray-900 mt-1">{{ qrInfo.supply_code.code }}</p>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700">Nombre</label>
              <p class="text-sm text-gray-900 mt-1">{{ qrInfo.supply_code.name }}</p>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700">Código Proveedor</label>
              <p class="text-sm text-gray-900 mt-1">{{ qrInfo.supply_code.code_supplier }}</p>
            </div>
          </div>

          <!-- Información del lote -->
          <div v-if="qrInfo.supply_info?.batch" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700">Lote</label>
              <p class="text-sm text-gray-900 mt-1">ID: {{ qrInfo.supply_info.batch.id }}</p>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700">Proveedor</label>
              <p class="text-sm text-gray-900 mt-1">{{ qrInfo.supply_info.batch.supplier }}</p>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700">Fecha de Vencimiento</label>
              <p class="text-sm text-gray-900 mt-1">{{ formatDate(qrInfo.supply_info.batch.expiration_date) }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Estado actual y resumen -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Estado actual -->
        <div class="bg-white rounded-lg shadow border p-6">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Estado Actual</h3>
          
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700">Estado</label>
              <span :class="getCurrentStatusBadgeClass()" class="inline-flex px-3 py-1 text-sm font-semibold rounded-full mt-1">
                {{ traceability.current_status || 'Desconocido' }}
              </span>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700">Asignado a Solicitud</label>
              <p class="text-sm mt-1">
                {{ traceability.is_assigned_to_request ? 'Sí' : 'No' }}
              </p>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700">Última Actualización</label>
              <p class="text-sm text-gray-900 mt-1">{{ formatDate(traceability.last_updated) }}</p>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700">Fecha de Creación</label>
              <p class="text-sm text-gray-900 mt-1">{{ formatDate(traceability.created_date) }}</p>
            </div>
          </div>
        </div>

        <!-- Estadísticas -->
        <div class="bg-white rounded-lg shadow border p-6">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Estadísticas</h3>
          
          <div class="space-y-4">
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-600">Total de Asignaciones:</span>
              <span class="text-sm font-medium text-gray-900">{{ traceability.request_history?.length || 0 }}</span>
            </div>
            
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-600">Movimientos de Supply:</span>
              <span class="text-sm font-medium text-gray-900">{{ traceability.supply_history?.length || 0 }}</span>
            </div>
            
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-600">Tiempo en el Sistema:</span>
              <span class="text-sm font-medium text-gray-900">{{ calculateTimeInSystem() }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Historial de Solicitudes -->
      <div v-if="traceability.request_history && traceability.request_history.length > 0" class="bg-white rounded-lg shadow border">
        <div class="px-6 py-4 border-b border-gray-200">
          <h3 class="text-lg font-semibold text-gray-900">Historial de Solicitudes</h3>
        </div>
        
        <div class="p-6">
          <div class="flow-root">
            <ul class="-mb-8">
              <li
                v-for="(assignment, index) in traceability.request_history"
                :key="assignment.id"
                class="relative pb-8"
              >
                <span
                  v-if="index !== traceability.request_history.length - 1"
                  class="absolute top-4 left-4 -ml-px h-full w-0.5 bg-gray-200"
                  aria-hidden="true"
                ></span>
                
                <div class="relative flex space-x-3">
                  <!-- Icono del evento -->
                  <div>
                    <span :class="getAssignmentIconClass(assignment.status)" class="h-8 w-8 rounded-full flex items-center justify-center ring-8 ring-white">
                      <svg class="h-5 w-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path v-if="assignment.status === 'assigned'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
                        <path v-else-if="assignment.status === 'delivered'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                        <path v-else-if="assignment.status === 'consumed'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                        <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3a1 1 0 011-1h6a1 1 0 011 1v4m-9 4h12" />
                      </svg>
                    </span>
                  </div>
                  
                  <!-- Contenido del evento -->
                  <div class="min-w-0 flex-1 pt-1.5">
                    <div class="flex justify-between items-start">
                      <div>
                        <p class="text-sm text-gray-900">
                          <span class="font-medium">{{ getAssignmentStatusLabel(assignment.status) }}</span>
                          <span v-if="assignment.supply_request" class="ml-2">
                            - Solicitud {{ assignment.supply_request.request_number }}
                          </span>
                        </p>
                        <p v-if="assignment.assigned_by_name" class="text-sm text-gray-500 mt-1">
                          Por: {{ assignment.assigned_by_name }}
                        </p>
                        <p v-if="assignment.notes" class="text-sm text-gray-600 mt-2 bg-gray-50 p-2 rounded">
                          {{ assignment.notes }}
                        </p>
                      </div>
                      
                      <div class="text-right">
                        <span :class="getAssignmentStatusBadgeClass(assignment.status)" class="inline-flex px-2 py-1 text-xs font-semibold rounded-full">
                          {{ getAssignmentStatusLabel(assignment.status) }}
                        </span>
                        <p class="text-xs text-gray-500 mt-1">{{ formatDate(assignment.assigned_date) }}</p>
                      </div>
                    </div>
                  </div>
                </div>
              </li>
            </ul>
          </div>
        </div>
      </div>

      <!-- Historial de Movimientos del Insumo -->
      <div v-if="traceability.supply_history && traceability.supply_history.length > 0" class="bg-white rounded-lg shadow border">
        <div class="px-6 py-4 border-b border-gray-200">
          <h3 class="text-lg font-semibold text-gray-900">Historial de Movimientos</h3>
        </div>
        
        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Estado
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Fecha/Hora
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Destino
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Observaciones
                </th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr
                v-for="movement in traceability.supply_history"
                :key="movement.id"
                class="hover:bg-gray-50"
              >
                <td class="px-6 py-4 whitespace-nowrap">
                  <span :class="getMovementStatusBadgeClass(movement.status)" class="inline-flex px-2 py-1 text-xs font-semibold rounded-full">
                    {{ movement.status }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ formatDate(movement.date_time) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ movement.destination_type }} {{ movement.destination_id }}
                </td>
                <td class="px-6 py-4 text-sm text-gray-900">
                  {{ movement.observations || '-' }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Información adicional si no hay historial -->
      <div v-if="(!traceability.request_history || traceability.request_history.length === 0) && (!traceability.supply_history || traceability.supply_history.length === 0)" class="bg-yellow-50 border border-yellow-200 rounded-lg p-6">
        <div class="flex">
          <svg class="h-5 w-5 text-yellow-400 mr-3 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <div>
            <h3 class="text-sm font-medium text-yellow-800">Sin historial disponible</h3>
            <p class="text-sm text-yellow-700 mt-2">
              Este código QR no tiene movimientos registrados aún. Esto puede significar que es un insumo recién creado que no ha sido procesado en el sistema.
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import supplyRequestService from '../services/supplyRequestService'
import qrService from '../services/qrService'
import { format, differenceInDays, differenceInHours } from 'date-fns'
import { es } from 'date-fns/locale'

const route = useRoute()

// Props
const props = defineProps({
  qrCode: {
    type: String,
    required: true
  }
})

// Estado reactivo
const loading = ref(false)
const error = ref(null)
const traceability = ref(null)
const qrInfo = ref(null)

// Computed
const qrCode = computed(() => props.qrCode || route.params.qrCode)

// Métodos principales
const loadTraceability = async () => {
  loading.value = true
  error.value = null

  try {
    // Cargar información básica del QR
    const qrResult = await qrService.scanQRCode(qrCode.value)
    qrInfo.value = qrResult

    // Cargar trazabilidad específica
    const traceResult = await supplyRequestService.getQRTraceability(qrCode.value)
    if (traceResult.success && traceResult.data) {
      traceability.value = traceResult.data
    } else {
      // Si no hay trazabilidad específica, crear una básica
      traceability.value = {
        qr_code: qrCode.value,
        current_status: qrInfo.value.is_consumed ? 'consumed' : 'available',
        is_assigned_to_request: false,
        request_history: [],
        supply_history: qrInfo.value.history || [],
        created_date: new Date(),
        last_updated: new Date()
      }
    }

    console.log('Trazabilidad cargada:', traceability.value)
  } catch (err) {
    console.error('Error cargando trazabilidad:', err)
    error.value = 'Error al cargar la información de trazabilidad'
  } finally {
    loading.value = false
  }
}

// Métodos auxiliares
const calculateTimeInSystem = () => {
  if (!traceability.value?.created_date) return 'N/A'
  
  try {
    const created = new Date(traceability.value.created_date)
    const now = new Date()
    const days = differenceInDays(now, created)
    
    if (days === 0) {
      const hours = differenceInHours(now, created)
      return hours === 0 ? 'Menos de 1 hora' : `${hours} hora${hours > 1 ? 's' : ''}`
    } else {
      return `${days} día${days > 1 ? 's' : ''}`
    }
  } catch {
    return 'N/A'
  }
}

const formatDate = (dateString) => {
  if (!dateString) return 'N/A'
  try {
    return format(new Date(dateString), 'dd/MM/yyyy HH:mm', { locale: es })
  } catch {
    return dateString
  }
}

// Métodos de estilo
const getTypeLabel = (type) => {
  const labels = {
    'batch': 'Lote de Productos',
    'medical_supply': 'Insumo Individual'
  }
  return labels[type] || type
}

const getTypeIconClass = (type) => {
  const classes = {
    'batch': 'bg-blue-100 text-blue-600',
    'medical_supply': 'bg-green-100 text-green-600'
  }
  return classes[type] || 'bg-gray-100 text-gray-600'
}

const getStatusLabel = () => {
  if (!qrInfo.value) return 'Desconocido'
  return qrService.getStatusLabel(qrInfo.value)
}

const getStatusBadgeClass = () => {
  if (!qrInfo.value) return 'bg-gray-100 text-gray-800'
  const color = qrService.getStatusColor(qrInfo.value)
  const classes = {
    'green': 'bg-green-100 text-green-800',
    'red': 'bg-red-100 text-red-800',
    'yellow': 'bg-yellow-100 text-yellow-800',
    'gray': 'bg-gray-100 text-gray-800'
  }
  return classes[color] || classes.gray
}

const getCurrentStatusBadgeClass = () => {
  const status = traceability.value?.current_status || 'unknown'
  const classes = {
    'available': 'bg-green-100 text-green-800',
    'assigned': 'bg-blue-100 text-blue-800',
    'delivered': 'bg-purple-100 text-purple-800',
    'consumed': 'bg-gray-100 text-gray-800',
    'expired': 'bg-red-100 text-red-800',
    'unknown': 'bg-gray-100 text-gray-800'
  }
  return classes[status] || classes.unknown
}

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

const getAssignmentIconClass = (status) => {
  const classes = {
    'assigned': 'bg-blue-500',
    'delivered': 'bg-green-500',
    'consumed': 'bg-gray-500',
    'returned': 'bg-yellow-500'
  }
  return classes[status] || 'bg-gray-500'
}

const getMovementStatusBadgeClass = (status) => {
  const classes = {
    'creado': 'bg-green-100 text-green-800',
    'movido': 'bg-blue-100 text-blue-800',
    'consumido': 'bg-red-100 text-red-800',
    'vencido': 'bg-red-100 text-red-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

// Lifecycle
onMounted(() => {
  if (qrCode.value) {
    loadTraceability()
  } else {
    error.value = 'No se proporcionó código QR'
  }
})
</script>

<style scoped>
/* Estilos adicionales si son necesarios */
.flow-root {
  overflow: hidden;
}

/* Animaciones suaves */
.transition-all {
  transition: all 0.3s ease;
}

/* Gradientes */
.bg-gradient-to-r {
  background-image: linear-gradient(to right, var(--tw-gradient-stops));
}

.from-blue-50 {
  --tw-gradient-from: #eff6ff;
  --tw-gradient-stops: var(--tw-gradient-from), var(--tw-gradient-to, rgba(239, 246, 255, 0));
}

.to-indigo-50 {
  --tw-gradient-to: #eef2ff;
}
</style>