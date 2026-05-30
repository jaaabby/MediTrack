<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <div class="bg-white rounded-lg shadow-sm border p-4 sm:p-6 mb-4 sm:mb-6 mx-3 sm:mx-6 lg:mx-8 mt-4">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div class="flex items-center gap-3">
          <router-link to="/qr" class="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-full flex-shrink-0">
            <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
          </router-link>
          <div>
            <h1 class="text-xl sm:text-2xl font-semibold text-gray-900">Trazabilidad del Insumo</h1>
            <p class="text-gray-600 mt-1 text-sm">Registro completo de movimientos y ubicaciones</p>
          </div>
        </div>
        <div class="flex items-center gap-2 sm:gap-3 flex-shrink-0">
          <button @click="exportTraceability" class="btn-secondary text-xs sm:text-sm p-2 sm:px-4 sm:py-2">
            <svg class="h-4 w-4 sm:mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            <span class="hidden sm:inline">Exportar</span>
          </button>
          <button @click="printTraceability" class="btn-secondary text-xs sm:text-sm p-2 sm:px-4 sm:py-2">
            <svg class="h-4 w-4 sm:mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z" />
            </svg>
            <span class="hidden sm:inline">Imprimir</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex justify-center items-center py-8 sm:py-12">
      <div class="text-center">
        <div class="animate-spin rounded-full h-10 w-10 sm:h-12 sm:w-12 border-b-2 border-blue-600 mx-auto"></div>
        <p class="mt-4 text-sm sm:text-base text-gray-600">Cargando trazabilidad...</p>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="max-w-4xl mx-auto px-3 sm:px-4 py-4 sm:py-8">
      <div class="bg-red-50 border border-red-200 rounded-lg p-4 sm:p-6">
        <div class="flex">
          <div class="flex-shrink-0">
            <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div class="ml-3 flex-1 min-w-0">
            <h3 class="text-xs sm:text-sm font-medium text-red-800">Error al cargar trazabilidad</h3>
            <div class="mt-2 text-xs sm:text-sm text-red-700 break-words">{{ error }}</div>
            <div class="mt-4">
              <button @click="loadTraceability" class="btn-primary text-xs sm:text-sm">
                Reintentar
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div v-else-if="traceabilityData" class="max-w-7xl mx-auto px-3 sm:px-6 lg:px-8 pb-4 sm:pb-8">
      <!-- Product Summary Card -->
      <div class="bg-white rounded-lg shadow-sm border mb-4 sm:mb-8">
        <div class="p-4 sm:p-6">
          <!-- Product Info Section -->
          <div class="flex flex-col sm:flex-row sm:items-start gap-3 sm:gap-4 mb-4 sm:mb-4">
            <!-- QR Code Icon -->
            <div class="bg-blue-100 p-3 rounded-lg flex-shrink-0 self-start">
              <svg class="h-8 w-8 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v2a2 2 0 002 2zm0 0h2a2 2 0 012 2v2a2 2 0 01-2 2H5a2 2 0 01-2-2v-2a2 2 0 012-2z" />
              </svg>
            </div>
            
            <!-- Product Info -->
            <div class="flex-1 min-w-0">
              <h2 class="text-lg sm:text-2xl font-bold text-gray-900 break-words leading-snug mb-2">
                {{ qrInfo?.supply_info?.name || 'Insumo Medico' }}
              </h2>
              <p class="text-gray-600 font-mono text-xs sm:text-sm break-all mb-3">{{ qrCode }}</p>
              <div class="flex flex-wrap items-center gap-2">
                <span :class="getCurrentStatusBadgeClass()" class="inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium whitespace-nowrap">
                  {{ getCurrentStatusLabel() }}
                </span>
                <span class="text-xs text-gray-500">
                  En sistema desde {{ getTimeInSystem() }}
                </span>
              </div>
            </div>
          </div>

          <!-- Current Location - Mobile Optimized -->
          <div v-if="getCurrentLocation()" class="mb-4 p-3 bg-blue-50 border border-blue-100 rounded-lg">
            <div class="text-xs font-medium text-blue-700 mb-1">Ubicacion actual</div>
            <div class="text-sm sm:text-base font-semibold text-gray-900 break-words leading-relaxed">
              {{ getCurrentLocation() }}
            </div>
          </div>

          <!-- Stats Bar - Mobile Optimized -->
          <div class="grid grid-cols-3 gap-2 sm:gap-4">
            <div class="bg-gray-50 rounded-lg p-3 sm:p-4 text-center">
              <div class="text-xl sm:text-2xl font-bold text-green-600">{{ getEventCount('movement') }}</div>
              <div class="text-xs sm:text-sm text-gray-600 mt-1 leading-tight">Cambios</div>
            </div>
            <div class="bg-gray-50 rounded-lg p-3 sm:p-4 text-center">
              <div class="text-xl sm:text-2xl font-bold text-purple-600">{{ getLocationCount() }}</div>
              <div class="text-xs sm:text-sm text-gray-600 mt-1 leading-tight">Ubicaciones</div>
            </div>
            <div class="bg-gray-50 rounded-lg p-3 sm:p-4 text-center">
              <div class="text-xl sm:text-2xl font-bold text-orange-600">{{ getUserCount() }}</div>
              <div class="text-xs sm:text-sm text-gray-600 mt-1 leading-tight">Usuarios</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Journey Map -->
      <div class="bg-white rounded-lg shadow-sm border mb-4 sm:mb-8">
        <div class="p-4 sm:p-6">
          <h3 class="text-base sm:text-lg font-semibold text-gray-900 mb-4">Recorrido del Insumo</h3>
          
          <!-- Progress Bar -->
          <div class="mb-6 sm:mb-8">
            <!-- Labels for Mobile -->
            <div class="flex items-center justify-between text-xs text-gray-500 mb-2 px-0">
              <span class="text-center flex-1">En bodega</span>
              <span class="text-center flex-1">En transito</span>
              <span class="text-center flex-1">Recepcionado</span>
              <span class="text-center flex-1">Consumido</span>
            </div>
            
            <div class="w-full bg-gray-200 rounded-full h-2.5">
              <div 
                :class="[
                  'h-2.5 rounded-full transition-all duration-1000',
                  getProgressBarColor()
                ]"
                :style="`width: ${getProgressPercentage()}%`"
              ></div>
            </div>
          </div>

          <!-- Location Journey - Mobile Optimized -->
          <div v-if="getLocationJourney().length > 0" class="space-y-3">
            <div
              v-for="(location, index) in getLocationJourney()"
              :key="index"
              class="relative"
            >
              <div :class="[
                'p-3 sm:p-4 rounded-lg border-2 transition-all duration-300',
                location.is_current ? 'border-blue-500 bg-blue-50' : 'border-gray-200 bg-white'
              ]">
                <div class="flex items-center justify-between mb-2">
                  <div :class="[
                    'flex items-center justify-center w-8 h-8 rounded-full text-sm font-semibold',
                    location.is_current ? 'bg-blue-500 text-white' : 'bg-gray-200 text-gray-600'
                  ]">
                    {{ index + 1 }}
                  </div>
                  <div v-if="location.is_current" class="flex items-center text-blue-600 text-xs font-medium">
                    <div class="w-2 h-2 bg-blue-600 rounded-full mr-1.5 animate-pulse"></div>
                    Actual
                  </div>
                </div>
                <div class="text-sm sm:text-base font-semibold text-gray-900 break-words leading-relaxed mb-1">
                  {{ location.name }}
                </div>
                <div class="text-xs text-gray-500">{{ formatDateTime(location.date) }}</div>
                <div v-if="location.duration" class="text-xs text-gray-400 mt-1">
                  Duracion: {{ location.duration }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Timeline -->
      <div class="bg-white rounded-lg shadow-sm border">
        <div class="p-4 sm:p-6">
          <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 mb-6">
            <h3 class="text-base sm:text-lg font-semibold text-gray-900">Historial Completo</h3>
            <div class="w-full sm:w-auto">
              <!-- Filters -->
              <select v-model="selectedEventFilter" class="form-select text-sm w-full" @change="applyFilters">
                <option value="">Todos los eventos</option>
                <option value="movement">Solo movimientos</option>
                <option value="consumption">Solo consumos</option>
                <option value="transfer">Solo transferencias</option>
              </select>
            </div>
          </div>

          <!-- Timeline Events -->
          <div class="flow-root">
            <ul class="space-y-0">
              <li v-for="(event, index) in filteredEvents" :key="index">
                <div class="relative pb-8">
                  <!-- Connector Line -->
                  <span
                    v-if="index !== filteredEvents.length - 1"
                    class="absolute top-5 left-4 -ml-px h-full w-0.5 bg-gray-200"
                    aria-hidden="true"
                  ></span>
                  
                  <div class="relative flex space-x-3">
                    <!-- Event Icon -->
                    <div class="flex-shrink-0">
                      <span :class="[
                        'h-8 w-8 rounded-full flex items-center justify-center ring-4 ring-white',
                        getEventIconClass(event)
                      ]">
                        <svg class="h-4 w-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="getEventIcon(event)" />
                        </svg>
                      </span>
                    </div>
                    
                    <!-- Event Content -->
                    <div class="min-w-0 flex-1">
                      <div>
                        <div class="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-1 mb-1">
                          <div class="text-sm font-semibold text-gray-900 break-words pr-2">
                            {{ getEventTitle(event) }}
                          </div>
                          <div class="text-xs text-gray-500 whitespace-nowrap flex-shrink-0">
                            {{ formatRelativeTime(event.date_time || event.timestamp) }}
                          </div>
                        </div>
                        <div class="text-sm text-gray-600 break-words leading-relaxed">
                          {{ getEventDescription(event) }}
                        </div>
                        
                        <!-- Event Details -->
                        <div v-if="hasEventDetails(event)" class="mt-2">
                          <div class="bg-gray-50 rounded-lg p-3 text-xs space-y-2">
                            <div v-if="event.user_name" class="flex flex-col sm:flex-row sm:justify-between gap-1">
                              <span class="text-gray-500 font-medium">Usuario:</span>
                              <span class="text-gray-900 break-words">{{ event.user_name }}</span>
                            </div>
                            <div v-if="event.location" class="flex flex-col sm:flex-row sm:justify-between gap-1">
                              <span class="text-gray-500 font-medium">Ubicacion:</span>
                              <span class="text-gray-900 break-words">{{ event.location }}</span>
                            </div>
                            <div v-if="event.notes" class="flex flex-col sm:flex-row sm:justify-between gap-1">
                              <span class="text-gray-500 font-medium">Notas:</span>
                              <span class="text-gray-900 break-words">{{ event.notes }}</span>
                            </div>
                            <div class="flex flex-col sm:flex-row sm:justify-between gap-1">
                              <span class="text-gray-500 font-medium">Fecha exacta:</span>
                              <span class="text-gray-900 whitespace-nowrap">{{ formatDateTime(event.date_time || event.timestamp) }}</span>
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </li>
            </ul>
          </div>

          <!-- Empty State -->
          <div v-if="filteredEvents.length === 0" class="text-center py-12">
            <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
            </svg>
            <h3 class="mt-2 text-sm font-medium text-gray-900">No hay eventos</h3>
            <p class="mt-1 text-sm text-gray-500">No se encontraron eventos para los filtros seleccionados.</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { format, formatDistanceToNow, differenceInDays, differenceInHours, differenceInMinutes } from 'date-fns'
import { es } from 'date-fns/locale'
import { useNotification } from '@/composables/useNotification'
import qrService from '@/services/qr/qrService'
import ExcelJS from 'exceljs'
import { jsPDF } from 'jspdf'
import autoTable from 'jspdf-autotable'

const route = useRoute()
const router = useRouter()
const { success: showSuccess, error: showError } = useNotification()

const props = defineProps({
  qrCode: {
    type: String,
    required: false
  }
})

const loading = ref(false)
const error = ref(null)
const traceabilityData = ref(null)
const qrInfo = ref(null)
const selectedEventFilter = ref('')

const qrCode = computed(() => props.qrCode || route.params.qrCode)

const filteredEvents = computed(() => {
  if (!traceabilityData.value || !selectedEventFilter.value) {
    return getAllEvents()
  }

  return getAllEvents().filter(event => {
    switch (selectedEventFilter.value) {
      case 'scan':
        return event.event_type === 'scan'
      case 'movement':
        return event.event_type === 'movement'
      case 'consumption':
        return event.event_type === 'movement' && event.status === 'consumido'
      case 'transfer':
        return event.event_type === 'movement' && event.status === 'transferido'
      default:
        return true
    }
  })
})

const loadTraceability = async () => {
  loading.value = true
  error.value = null
  
  try {
    const traceData = await qrService.getCompleteTraceability(qrCode.value)
    traceabilityData.value = traceData
    
    const qrData = await qrService.scanQRCode(qrCode.value, {
      scan_purpose: 'traceability_view',
      scan_source: 'web'
    })
    qrInfo.value = qrData
    
  } catch (err) {
    const errorMessage = err.message || 'Error al cargar la trazabilidad'
    error.value = errorMessage
    showError(errorMessage)
    console.error('Error cargando trazabilidad:', err)
  } finally {
    loading.value = false
  }
}

const getAllEvents = () => {
  if (!traceabilityData.value) return []

  const events = []
  
  if (traceabilityData.value.supply_history) {
    const uniqueMovements = new Map()
    
    traceabilityData.value.supply_history.forEach(movement => {
      const timestamp = movement.timestamp || movement.date_time
      const key = `${movement.status}_${timestamp}_${movement.destination_id}_${movement.user_rut}`
      
      if (!uniqueMovements.has(key)) {
        uniqueMovements.set(key, {
          ...movement,
          event_type: 'movement',
          date_time: timestamp,
          title: getMovementTitle(movement),
          location: getMovementLocation(movement)
        })
      }
    })
    
    events.push(...Array.from(uniqueMovements.values()))
  }

  return events.sort((a, b) => new Date(b.date_time) - new Date(a.date_time))
}

const getMovementTitle = (movement) => {
  if (movement.status === 'consumido') return 'Insumo Consumido'
  if (movement.status === 'transferido') return 'Insumo Transferido'
  if (movement.status === 'creado') return 'Insumo Creado'
  if (movement.status === 'en_camino_a_pabellon') return 'En camino a pabellon'
  if (movement.status === 'recepcionado') return 'Insumo Recepcionado'
  if (movement.status === 'disponible') return 'Insumo Disponible'
  return `Movimiento: ${movement.status || movement.movement_type || 'Cambio'}`
}

const getMovementLocation = (movement) => {
  // Si el estado es pendiente_retiro, el insumo físicamente está en bodega, no en el destino
  if (movement.status?.toLowerCase() === 'pendiente_retiro') {
    // Buscar la ubicación de bodega anterior en el historial
    // O mostrar que está en bodega
    if (movement.location && (movement.location.toLowerCase().includes('bodega') || 
        movement.location.toLowerCase().includes('almacen') || 
        movement.location.toLowerCase().includes('store'))) {
      return movement.location
    }
    // Si no hay ubicación previa, mostrar que está en bodega
    return 'Almacen: Bodega Principal (Pendiente de Retiro)'
  }
  
  if (movement.destination_name) {
    const type = movement.destination_type === 'pavilion' ? 'Pabellon' : 'Almacen'
    let location = `${type}: ${movement.destination_name}`
    if (movement.medical_center_name) {
      location += ` (${movement.medical_center_name})`
    }
    return location
  }
  return movement.location || 'Ubicacion no especificada'
}

const getCurrentStatusLabel = () => {
  if (!qrInfo.value) return 'Desconocido'
  
  if (qrInfo.value.is_consumed) return 'Consumido'
  
  const status = qrInfo.value.supply_info?.status || qrInfo.value.status
  switch (status?.toLowerCase()) {
    case 'disponible': return 'Disponible'
    case 'pendiente_retiro': return 'Pendiente de Retiro'
    case 'en_camino_a_pabellon': return 'En Camino a Pabellón'
    case 'recepcionado': return 'Recepcionado'
    case 'transferido': return 'Transferido'
    case 'en_camino_a_bodega': return 'En Camino a Bodega'
    case 'consumido': return 'Consumido'
    default: return status || 'Activo'
  }
}

const getCurrentStatusBadgeClass = () => {
  if (!qrInfo.value) return 'bg-gray-100 text-gray-800'
  
  if (qrInfo.value.is_consumed) return 'bg-red-100 text-red-800'
  
  const status = qrInfo.value.supply_info?.status || qrInfo.value.status
  switch (status?.toLowerCase()) {
    case 'disponible': return 'bg-green-100 text-green-800'
    case 'pendiente_retiro': return 'bg-yellow-100 text-yellow-800'
    case 'en_camino_a_pabellon': return 'bg-orange-100 text-orange-800'
    case 'recepcionado': return 'bg-blue-100 text-blue-800'
    case 'transferido': return 'bg-purple-100 text-purple-800'
    case 'en_camino_a_bodega': return 'bg-orange-100 text-orange-800'
    case 'consumido': return 'bg-red-100 text-red-800'
    default: return 'bg-gray-100 text-gray-800'
  }
}

const getCurrentLocation = () => {
  // Si el insumo está consumido, mostrar dónde fue consumido
  const status = qrInfo.value?.supply_info?.status || qrInfo.value?.status
  
  if (status?.toLowerCase() === 'consumido' || qrInfo.value?.is_consumed) {
    const events = getAllEvents()
    
    // Buscar el evento de consumo
    const consumptionEvent = events.find(e => {
      const isConsumption = e.event_type === 'movement' && 
                           (e.status === 'consumido' || e.status === 'consumed')
      return isConsumption
    })
    
    if (consumptionEvent && consumptionEvent.location) {
      return `Fue usado en ${consumptionEvent.location}`
    }
    
    // Buscar el último evento de recepción que debe tener la ubicación
    const receptionEvent = events.find(e => {
      const isReception = e.event_type === 'movement' && 
                         (e.status === 'recepcionado' || e.status === 'received')
      return isReception
    })
    
    if (receptionEvent) {
      // Verificar múltiples fuentes de ubicación
      const location = receptionEvent.location || 
                      (receptionEvent.destination_name ? 
                        `${receptionEvent.destination_type === 'pavilion' ? 'Pabellón' : 'Almacén'}: ${receptionEvent.destination_name}` : 
                        null)
    }
  }
  
  // Si el estado es pendiente_retiro, el insumo físicamente está en bodega
  if (status?.toLowerCase() === 'pendiente_retiro') {
    // Buscar la última ubicación de bodega antes del movimiento a pabellón
    const events = getAllEvents()
    const storeLocation = events
      .filter(e => e.event_type === 'movement' && e.location && 
              (e.location.toLowerCase().includes('bodega') || 
               e.location.toLowerCase().includes('almacen') ||
               e.location.toLowerCase().includes('store')))
      .reverse()[0]
    
    if (storeLocation) {
      return storeLocation.location
    }
    // Si no hay ubicación de bodega previa, buscar en supply_info
    const supplyInfo = qrInfo.value?.supply_info
    if (supplyInfo?.location_type === 'store' || supplyInfo?.store_name) {
      return supplyInfo.store_name || 'Bodega'
    }
    return 'Bodega (Pendiente de Retiro)'
  }
  
  // Para otros estados, usar la lógica normal
  const events = getAllEvents()
  const lastMovement = events.find(e => e.event_type === 'movement' && e.location)
  return lastMovement?.location || 'Sin ubicacion registrada'
}

const getTimeInSystem = () => {
  if (!traceabilityData.value?.created_date && !qrInfo.value?.created_at) return 'N/A'
  
  try {
    const createdDate = new Date(traceabilityData.value?.created_date || qrInfo.value?.created_at)
    const now = new Date()
    
    const days = differenceInDays(now, createdDate)
    if (days > 0) {
      return `${days} dia${days > 1 ? 's' : ''}`
    }
    
    const hours = differenceInHours(now, createdDate)
    if (hours > 0) {
      return `${hours} hora${hours > 1 ? 's' : ''}`
    }
    
    const minutes = differenceInMinutes(now, createdDate)
    return `${minutes} minuto${minutes > 1 ? 's' : ''}`
  } catch {
    return 'N/A'
  }
}

const getProgressPercentage = () => {
  if (!qrInfo.value) return 0
  
  const status = qrInfo.value.supply_info?.status || qrInfo.value.status
  
  switch (status) {
    case 'disponible':
      return 25
    case 'en_camino_a_pabellon':
      return 50
    case 'recepcionado':
      return 75
    case 'consumido':
      return 100
    default:
      return 0
  }
}

const getProgressBarColor = () => {
  if (!qrInfo.value) return 'bg-gray-500'
  
  const status = qrInfo.value.supply_info?.status || qrInfo.value.status
  
  switch (status) {
    case 'disponible':
      return 'bg-green-500'
    case 'en_camino_a_pabellon':
      return 'bg-yellow-500'
    case 'recepcionado':
      return 'bg-blue-500'
    case 'consumido':
      return 'bg-red-500'
    default:
      return 'bg-gray-500'
  }
}

const getEventCount = (type) => {
  if (type === 'scan') return 0
  return getAllEvents().filter(e => e.event_type === type).length
}

const getLocationCount = () => {
  const locations = new Set()
  getAllEvents().forEach(event => {
    if (event.location) locations.add(event.location)
  })
  return locations.size
}

const getUserCount = () => {
  const users = new Set()
  getAllEvents().forEach(event => {
    if (event.user_name) users.add(event.user_name)
  })
  return users.size
}

const getLocationJourney = () => {
  const movements = getAllEvents()
    .filter(e => e.event_type === 'movement' && e.location)
    .reverse()

  const journey = []
  const seenLocations = new Set()

  movements.forEach((movement, index) => {
    if (!seenLocations.has(movement.location)) {
      seenLocations.add(movement.location)
      
      const nextMovement = movements[index + 1]
      let duration = null
      
      if (nextMovement) {
        try {
          const start = new Date(movement.date_time)
          const end = new Date(nextMovement.date_time)
          const hours = differenceInHours(end, start)
          const days = differenceInDays(end, start)
          
          if (days > 0) {
            duration = `${days} dia${days > 1 ? 's' : ''}`
          } else if (hours > 0) {
            duration = `${hours} hora${hours > 1 ? 's' : ''}`
          } else {
            const minutes = differenceInMinutes(end, start)
            duration = `${minutes} min`
          }
        } catch {
          duration = null
        }
      }

      journey.push({
        name: movement.location,
        date: movement.date_time,
        duration: duration,
        is_current: index === movements.length - 1
      })
    }
  })

  return journey
}

const getEventTitle = (event) => {
  switch (event.event_type) {
    case 'scan':
      return 'Codigo QR Escaneado'
    case 'movement':
      return getMovementTitle(event)
    default:
      return 'Evento'
  }
}

const getEventDescription = (event) => {
  switch (event.event_type) {
    case 'scan':
      return `Escaneado por ${event.user_name || 'Usuario desconocido'} desde ${event.scan_source || 'aplicacion'}`
    case 'movement':
      const statusText = getStatusText(event.status)
      let description = event.observations || `Estado: ${statusText}`
      if (event.location) {
        description += ` - ${event.location}`
      }
      return description
    default:
      return 'Informacion no disponible'
  }
}

const getStatusText = (status) => {
  switch (status) {
    case 'disponible': return 'Disponible'
    case 'en_camino_a_pabellon': return 'En camino a pabellon'
    case 'recepcionado': return 'Recepcionado'
    case 'consumido': return 'Consumido'
    case 'transferido': return 'Transferido'
    case 'creado': return 'Creado'
    default: return status || 'Procesado'
  }
}

const getEventIcon = (event) => {
  switch (event.event_type) {
    case 'scan':
      return 'M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v2a2 2 0 002 2zm0 0h2a2 2 0 012 2v2a2 2 0 01-2 2H5a2 2 0 01-2-2v-2a2 2 0 012-2z'
    case 'movement':
      if (event.status === 'consumido') {
        return 'M5 13l4 4L19 7'
      } else if (event.status === 'transferido') {
        return 'M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4'
      }
      return 'M13 10V3L4 14h7v7l9-11h-7z'
    default:
      return 'M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2'
  }
}

const getEventIconClass = (event) => {
  switch (event.event_type) {
    case 'scan':
      return 'bg-blue-500'
    case 'movement':
      if (event.status === 'consumido') return 'bg-red-500'
      if (event.status === 'transferido') return 'bg-purple-500'
      if (event.status === 'creado') return 'bg-green-500'
      return 'bg-orange-500'
    default:
      return 'bg-gray-500'
  }
}

const hasEventDetails = (event) => {
  return event.user_name || event.location || event.notes || true
}

const formatDateTime = (dateString) => {
  if (!dateString) return 'N/A'
  try {
    return format(new Date(dateString), 'dd/MM/yyyy HH:mm', { locale: es })
  } catch {
    return dateString
  }
}

const formatRelativeTime = (dateString) => {
  if (!dateString) return 'N/A'
  try {
    return formatDistanceToNow(new Date(dateString), { addSuffix: true, locale: es })
  } catch {
    return dateString
  }
}

const exportTraceability = async () => {
  try {
    // Crear libro de Excel
    const workbook = new ExcelJS.Workbook()

    // ==================== Hoja 1: Resumen ====================
    const summarySheet = workbook.addWorksheet('Resumen')
    
    // Agregar datos
    summarySheet.addRow(['REPORTE DE TRAZABILIDAD DEL INSUMO MEDICO'])
    summarySheet.addRow([])
    summarySheet.addRow(['INFORMACION GENERAL'])
    summarySheet.addRow(['Código QR:', qrCode.value])
    summarySheet.addRow(['Nombre del Insumo:', qrInfo.value?.supply_info?.name || 'N/A'])
    summarySheet.addRow(['Código de Proveedor:', qrInfo.value?.supply_code?.code_supplier || 'N/A'])
    summarySheet.addRow(['Lote:', qrInfo.value?.supply_info?.batch?.id || qrInfo.value?.batch_id || 'N/A'])
    summarySheet.addRow(['Estado Actual:', getCurrentStatusLabel()])
    summarySheet.addRow(['Ubicación Actual:', getCurrentLocation()])
    summarySheet.addRow(['Tiempo en Sistema:', getTimeInSystem()])
    summarySheet.addRow(['Fecha de Exportación:', format(new Date(), 'dd/MM/yyyy HH:mm', { locale: es })])
    summarySheet.addRow([])
    summarySheet.addRow(['ESTADISTICAS DEL INSUMO'])
    summarySheet.addRow(['Total de Eventos Registrados:', getAllEvents().length])
    summarySheet.addRow(['Total de Movimientos:', getEventCount('movement')])
    summarySheet.addRow(['Ubicaciones Visitadas:', getLocationCount()])
    summarySheet.addRow(['Usuarios Involucrados:', getUserCount()])
    summarySheet.addRow([])
    summarySheet.addRow(['INFORMACION DEL LOTE'])
    summarySheet.addRow(['Proveedor:', qrInfo.value?.batch_info?.supplier || qrInfo.value?.supplier || 'N/A'])
    summarySheet.addRow(['Fecha de Vencimiento:', qrInfo.value?.batch_info?.expiration_date ? formatDateTime(qrInfo.value.batch_info.expiration_date) : 'N/A'])
    
    // Aplicar estilos y anchos de columna
    summarySheet.getColumn(1).width = 30
    summarySheet.getColumn(2).width = 50
    
    // Estilizar y mergear celdas para títulos
    summarySheet.getRow(1).font = { bold: true, size: 14 }
    summarySheet.mergeCells('A1:B1')
    summarySheet.getRow(3).font = { bold: true }
    summarySheet.mergeCells('A3:B3')
    summarySheet.getRow(13).font = { bold: true }
    summarySheet.mergeCells('A13:B13')
    summarySheet.getRow(19).font = { bold: true }
    summarySheet.mergeCells('A19:B19')

    // ==================== Hoja 2: Historial de Eventos ====================
    const events = getAllEvents()
    const eventsSheet = workbook.addWorksheet('Historial de Eventos')
    
    eventsSheet.addRow(['HISTORIAL COMPLETO DE EVENTOS'])
    eventsSheet.addRow([])
    eventsSheet.addRow(['Fecha y Hora', 'Tipo de Evento', 'Estado', 'Descripción', 'Usuario', 'Ubicación', 'Notas'])
    
    // Estilizar encabezados
    const eventsHeaderRow = eventsSheet.getRow(3)
    eventsHeaderRow.font = { bold: true }
    eventsHeaderRow.fill = {
      type: 'pattern',
      pattern: 'solid',
      fgColor: { argb: 'FFE0E0E0' }
    }
    
    events.forEach(event => {
      eventsSheet.addRow([
        formatDateTime(event.date_time || event.timestamp),
        getEventTitle(event),
        event.status || 'N/A',
        getEventDescription(event),
        event.user_name || 'N/A',
        event.location || 'N/A',
        event.notes || event.observations || 'N/A'
      ])
    })
    
    // Configurar anchos de columna
    eventsSheet.getColumn(1).width = 18
    eventsSheet.getColumn(2).width = 25
    eventsSheet.getColumn(3).width = 20
    eventsSheet.getColumn(4).width = 40
    eventsSheet.getColumn(5).width = 25
    eventsSheet.getColumn(6).width = 35
    eventsSheet.getColumn(7).width = 40
    
    // Mergear título
    eventsSheet.getRow(1).font = { bold: true, size: 14 }
    eventsSheet.mergeCells('A1:G1')

    // ==================== Hoja 3: Recorrido del Insumo ====================
    const journey = getLocationJourney()
    const journeySheet = workbook.addWorksheet('Recorrido')
    
    journeySheet.addRow(['RECORRIDO DEL INSUMO POR UBICACIONES'])
    journeySheet.addRow([])
    journeySheet.addRow(['Secuencia', 'Ubicación', 'Fecha de Llegada', 'Duración en Ubicación', 'Estado Actual'])
    
    // Estilizar encabezados
    const journeyHeaderRow = journeySheet.getRow(3)
    journeyHeaderRow.font = { bold: true }
    journeyHeaderRow.fill = {
      type: 'pattern',
      pattern: 'solid',
      fgColor: { argb: 'FFE0E0E0' }
    }
    
    journey.forEach((location, index) => {
      journeySheet.addRow([
        index + 1,
        location.name,
        formatDateTime(location.date),
        location.duration || 'Ubicación actual',
        location.is_current ? '✓ UBICACION ACTUAL' : 'Ubicación anterior'
      ])
    })
    
    // Configurar anchos de columna
    journeySheet.getColumn(1).width = 12
    journeySheet.getColumn(2).width = 45
    journeySheet.getColumn(3).width = 20
    journeySheet.getColumn(4).width = 25
    journeySheet.getColumn(5).width = 22
    
    // Mergear título
    journeySheet.getRow(1).font = { bold: true, size: 14 }
    journeySheet.mergeCells('A1:E1')

    // ==================== Hoja 4: Resumen por Usuario ====================
    const userSummarySheet = workbook.addWorksheet('Resumen por Usuario')
    
    userSummarySheet.addRow(['RESUMEN DE ACTIVIDAD POR USUARIO'])
    userSummarySheet.addRow([])
    userSummarySheet.addRow(['Usuario', 'Número de Eventos', 'Último Evento', 'Fecha del Último Evento'])
    
    // Estilizar encabezados
    const userSummaryHeaderRow = userSummarySheet.getRow(3)
    userSummaryHeaderRow.font = { bold: true }
    userSummaryHeaderRow.fill = {
      type: 'pattern',
      pattern: 'solid',
      fgColor: { argb: 'FFE0E0E0' }
    }
    
    // Agrupar eventos por usuario
    const userEvents = new Map()
    events.forEach(event => {
      const userName = event.user_name || 'Usuario no identificado'
      if (!userEvents.has(userName)) {
        userEvents.set(userName, [])
      }
      userEvents.get(userName).push(event)
    })
    
    // Crear resumen por usuario
    Array.from(userEvents.entries()).forEach(([userName, userEventList]) => {
      const lastEvent = userEventList[0] // Ya están ordenados por fecha descendente
      userSummarySheet.addRow([
        userName,
        userEventList.length,
        getEventTitle(lastEvent),
        formatDateTime(lastEvent.date_time || lastEvent.timestamp)
      ])
    })
    
    // Configurar anchos de columna
    userSummarySheet.getColumn(1).width = 30
    userSummarySheet.getColumn(2).width = 20
    userSummarySheet.getColumn(3).width = 35
    userSummarySheet.getColumn(4).width = 20
    
    // Mergear título
    userSummarySheet.getRow(1).font = { bold: true, size: 14 }
    userSummarySheet.mergeCells('A1:D1')

    // Generar archivo Excel
    const fileName = `Trazabilidad_${qrCode.value}_${format(new Date(), 'yyyy-MM-dd_HHmm')}.xlsx`
    const buffer = await workbook.xlsx.writeBuffer()
    const blob = new Blob([buffer], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = fileName
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    
    showSuccess('Trazabilidad exportada a Excel exitosamente')
  } catch (err) {
    console.error('Error al exportar trazabilidad:', err)
    showError('Error al exportar la trazabilidad: ' + err.message)
  }
}

const printTraceability = () => {
  // Manejo de error explícito: el PDF necesita los datos cargados.
  if (!traceabilityData.value || !qrInfo.value) {
    showError('Todavía no se cargó la trazabilidad del insumo. Esperá a que termine de cargar e intentá de nuevo.')
    return
  }

  try {
    const doc = new jsPDF({ orientation: 'portrait', unit: 'mm', format: 'a4' })
    const pageWidth = doc.internal.pageSize.getWidth()
    const pageHeight = doc.internal.pageSize.getHeight()
    const margin = 15
    const contentWidth = pageWidth - margin * 2

    const BLUE = [37, 99, 235]
    const GRAY = [107, 114, 128]

    // Barra de título de sección reutilizable; devuelve la Y donde empieza el contenido.
    const sectionTitle = (title, y) => {
      doc.setFillColor(BLUE[0], BLUE[1], BLUE[2])
      doc.rect(margin, y, contentWidth, 7, 'F')
      doc.setTextColor(255, 255, 255)
      doc.setFontSize(9)
      doc.setFont('helvetica', 'bold')
      doc.text(title, margin + 3, y + 4.8)
      return y + 11
    }

    // ===== ENCABEZADO =====
    doc.setFillColor(BLUE[0], BLUE[1], BLUE[2])
    doc.rect(0, 0, pageWidth, 28, 'F')
    doc.setTextColor(255, 255, 255)
    doc.setFontSize(18)
    doc.setFont('helvetica', 'bold')
    doc.text('Reporte de Trazabilidad', margin, 13)
    doc.setFontSize(9)
    doc.setFont('helvetica', 'normal')
    doc.text('MediTrack · Sistema de Trazabilidad Médica', margin, 20)
    doc.text(`Generado: ${format(new Date(), 'dd/MM/yyyy HH:mm', { locale: es })}`,
      pageWidth - margin, 13, { align: 'right' })

    let y = 36

    // ===== INFORMACIÓN DEL INSUMO =====
    y = sectionTitle('Información del insumo', y)
    const infoRows = [
      ['Insumo', qrInfo.value?.supply_info?.name || 'Insumo médico'],
      ['Código QR', qrCode.value || '—'],
      ['Estado actual', getCurrentStatusLabel()],
      ['Ubicación actual', getCurrentLocation()],
      ['Tiempo en sistema', getTimeInSystem()],
      ['Código proveedor', String(qrInfo.value?.supply_code?.code_supplier ?? 'N/A')]
    ]
    const proveedor = qrInfo.value?.supply_info?.batch?.supplier || qrInfo.value?.supply_info?.supplier_name
    if (proveedor) infoRows.push(['Proveedor', proveedor])

    autoTable(doc, {
      startY: y,
      theme: 'plain',
      body: infoRows,
      styles: { fontSize: 9, cellPadding: 1.5, textColor: [31, 41, 55], overflow: 'linebreak' },
      columnStyles: {
        0: { fontStyle: 'bold', cellWidth: 45, textColor: GRAY },
        1: { cellWidth: contentWidth - 45 }
      },
      margin: { left: margin, right: margin }
    })
    y = doc.lastAutoTable.finalY + 6

    // ===== RESUMEN / ESTADÍSTICAS =====
    y = sectionTitle('Resumen', y)
    autoTable(doc, {
      startY: y,
      theme: 'grid',
      head: [['Eventos', 'Movimientos', 'Ubicaciones', 'Usuarios']],
      body: [[
        String(getAllEvents().length),
        String(getEventCount('movement')),
        String(getLocationCount()),
        String(getUserCount())
      ]],
      headStyles: { fillColor: BLUE, textColor: 255, halign: 'center', fontSize: 8 },
      bodyStyles: { halign: 'center', fontSize: 13, fontStyle: 'bold', textColor: [17, 24, 39] },
      margin: { left: margin, right: margin }
    })
    y = doc.lastAutoTable.finalY + 6

    // ===== RECORRIDO DEL INSUMO =====
    const journey = getLocationJourney()
    if (journey.length) {
      y = sectionTitle('Recorrido del insumo', y)
      autoTable(doc, {
        startY: y,
        theme: 'striped',
        head: [['#', 'Ubicación', 'Fecha', 'Duración']],
        body: journey.map((l, i) => [
          String(i + 1),
          (l.name || '—') + (l.is_current ? '  (actual)' : ''),
          formatDateTime(l.date),
          l.duration || '—'
        ]),
        headStyles: { fillColor: BLUE, textColor: 255, fontSize: 8 },
        bodyStyles: { fontSize: 8, overflow: 'linebreak' },
        columnStyles: { 0: { cellWidth: 10, halign: 'center' }, 3: { cellWidth: 28 } },
        margin: { left: margin, right: margin }
      })
      y = doc.lastAutoTable.finalY + 6
    }

    // ===== HISTORIAL DETALLADO =====
    const events = getAllEvents()
    y = sectionTitle('Historial detallado', y)
    autoTable(doc, {
      startY: y,
      theme: 'striped',
      head: [['Fecha', 'Evento', 'Estado', 'Ubicación', 'Usuario']],
      body: events.length
        ? events.map(e => [
            formatDateTime(e.date_time || e.timestamp),
            getEventTitle(e),
            getStatusText(e.status),
            e.location || '—',
            e.user_name || '—'
          ])
        : [['—', 'Sin eventos registrados', '—', '—', '—']],
      headStyles: { fillColor: BLUE, textColor: 255, fontSize: 8 },
      bodyStyles: { fontSize: 8, overflow: 'linebreak' },
      columnStyles: { 0: { cellWidth: 30 }, 2: { cellWidth: 26 } },
      margin: { left: margin, right: margin }
    })

    // ===== PIE DE PÁGINA (todas las páginas) =====
    const pageCount = doc.internal.getNumberOfPages()
    for (let i = 1; i <= pageCount; i++) {
      doc.setPage(i)
      doc.setFontSize(7)
      doc.setFont('helvetica', 'normal')
      doc.setTextColor(GRAY[0], GRAY[1], GRAY[2])
      doc.text('MediTrack — Sistema de Trazabilidad Médica', margin, pageHeight - 8)
      doc.text(`Página ${i} de ${pageCount}`, pageWidth - margin, pageHeight - 8, { align: 'right' })
    }

    const fileName = `Trazabilidad_${qrCode.value}_${format(new Date(), 'yyyy-MM-dd_HHmm')}.pdf`
    doc.save(fileName)
    showSuccess('PDF generado exitosamente')
  } catch (err) {
    console.error('Error al generar PDF de trazabilidad:', err)
    showError('Error al generar el PDF: ' + (err?.message || 'error desconocido'))
  }
}

const applyFilters = () => {
  // Los filtros se aplican automaticamente
}

onMounted(() => {
  if (qrCode.value) {
    loadTraceability()
  }
})
</script>

<style scoped>
.form-select {
  @apply block w-full px-3 py-2 text-sm border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-white;
  min-height: 44px;
}

.btn-primary {
  /* Usar .btn-primary de style.css global */
  min-height: 44px;
}

/* Usar clases de botones de style.css global */
/* min-height se mantiene local para este componente */
.btn-secondary {
  min-height: 44px;
}

.flow-root {
  overflow: visible;
}

* {
  -webkit-tap-highlight-color: transparent;
}

button,
select,
.form-select {
  transition: all 0.2s ease-in-out;
}

@media (max-width: 640px) {
  select,
  .form-select {
    font-size: 16px;
  }
}

@media print {
  .no-print {
    display: none;
  }
  
  .print-only {
    display: block;
  }
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

.animate-pulse {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}
</style>