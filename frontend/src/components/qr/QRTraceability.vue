<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <div class="bg-white shadow-sm border-b">
      <div class="max-w-7xl mx-auto px-3 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between h-14 sm:h-16">
          <div class="flex items-center space-x-2 sm:space-x-4 flex-1 min-w-0">
            <router-link to="/qr" class="text-gray-400 hover:text-gray-600 flex-shrink-0">
              <svg class="h-5 w-5 sm:h-6 sm:w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
            </router-link>
            <h1 class="text-sm sm:text-xl font-semibold text-gray-900 truncate">Trazabilidad del Insumo</h1>
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
    <div v-else-if="traceabilityData" class="max-w-7xl mx-auto px-3 sm:px-6 lg:px-8 py-4 sm:py-8">
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
import * as XLSX from 'xlsx'
import { jsPDF } from 'jspdf'
import 'jspdf-autotable'

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
  
  console.log('🔍 getCurrentLocation - Debug:', {
    status,
    is_consumed: qrInfo.value?.is_consumed,
    supply_info_status: qrInfo.value?.supply_info?.status
  })
  
  if (status?.toLowerCase() === 'consumido' || qrInfo.value?.is_consumed) {
    const events = getAllEvents()
    console.log('📋 Eventos totales:', events.length)
    
    // Buscar el evento de consumo
    const consumptionEvent = events.find(e => {
      const isConsumption = e.event_type === 'movement' && 
                           (e.status === 'consumido' || e.status === 'consumed')
      if (isConsumption) {
        console.log('✅ Encontrado evento de consumo:', e)
      }
      return isConsumption
    })
    
    if (consumptionEvent && consumptionEvent.location) {
      console.log('✓ Usando ubicación del evento de consumo:', consumptionEvent.location)
      return `Fue usado en ${consumptionEvent.location}`
    }
    
    console.log('⚠️ No hay ubicación en evento de consumo, buscando recepción...')
    
    // Buscar el último evento de recepción que debe tener la ubicación
    const receptionEvent = events.find(e => {
      const isReception = e.event_type === 'movement' && 
                         (e.status === 'recepcionado' || e.status === 'received')
      if (isReception) {
        console.log('📍 Encontrado evento de recepción:', {
          location: e.location,
          destination_name: e.destination_name,
          event: e
        })
      }
      return isReception
    })
    
    if (receptionEvent) {
      // Verificar múltiples fuentes de ubicación
      const location = receptionEvent.location || 
                      (receptionEvent.destination_name ? 
                        `${receptionEvent.destination_type === 'pavilion' ? 'Pabellón' : 'Almacén'}: ${receptionEvent.destination_name}` : 
                        null)
      
      if (location) {
        console.log('✓ Usando ubicación del evento de recepción:', location)
        return `Fue usado en ${location}`
      }
    }
    
    // Si no hay evento de consumo pero hay información de ubicación en supply_info
    if (qrInfo.value?.supply_info?.location_type === 'pavilion' && qrInfo.value?.supply_info?.pavilion_name) {
      console.log('✓ Usando ubicación de supply_info')
      return `Fue usado en Pabellón: ${qrInfo.value.supply_info.pavilion_name}`
    }
    
    console.log('❌ No se encontró ubicación en ninguna fuente')
    return 'Consumido (ubicación no especificada)'
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

const exportTraceability = () => {
  try {
    // Crear libro de Excel
    const workbook = XLSX.utils.book_new()

    // ==================== Hoja 1: Resumen ====================
    const summaryData = [
      ['REPORTE DE TRAZABILIDAD DEL INSUMO MEDICO'],
      [''],
      ['INFORMACION GENERAL'],
      ['Código QR:', qrCode.value],
      ['Nombre del Insumo:', qrInfo.value?.supply_info?.name || 'N/A'],
      ['Código de Proveedor:', qrInfo.value?.supply_code?.code_supplier || 'N/A'],
      ['Lote:', qrInfo.value?.supply_info?.batch?.id || qrInfo.value?.batch_id || 'N/A'],
      ['Estado Actual:', getCurrentStatusLabel()],
      ['Ubicación Actual:', getCurrentLocation()],
      ['Tiempo en Sistema:', getTimeInSystem()],
      ['Fecha de Exportación:', format(new Date(), 'dd/MM/yyyy HH:mm', { locale: es })],
      [''],
      ['ESTADISTICAS DEL INSUMO'],
      ['Total de Eventos Registrados:', getAllEvents().length],
      ['Total de Movimientos:', getEventCount('movement')],
      ['Ubicaciones Visitadas:', getLocationCount()],
      ['Usuarios Involucrados:', getUserCount()],
      [''],
      ['INFORMACION DEL LOTE'],
      ['Proveedor:', qrInfo.value?.batch_info?.supplier || qrInfo.value?.supplier || 'N/A'],
      ['Fecha de Vencimiento:', qrInfo.value?.batch_info?.expiration_date ? formatDateTime(qrInfo.value.batch_info.expiration_date) : 'N/A'],
    ]
    
    const summarySheet = XLSX.utils.aoa_to_sheet(summaryData)
    
    // Aplicar estilos y anchos de columna
    summarySheet['!cols'] = [
      { wch: 30 }, // Columna A (etiquetas)
      { wch: 50 }  // Columna B (valores)
    ]
    
    // Mergear celdas para títulos
    if (!summarySheet['!merges']) summarySheet['!merges'] = []
    summarySheet['!merges'].push(
      { s: { r: 0, c: 0 }, e: { r: 0, c: 1 } }, // Título principal
      { s: { r: 2, c: 0 }, e: { r: 2, c: 1 } }, // INFORMACION GENERAL
      { s: { r: 12, c: 0 }, e: { r: 12, c: 1 } }, // ESTADISTICAS
      { s: { r: 18, c: 0 }, e: { r: 18, c: 1 } }  // INFORMACION DEL LOTE
    )
    
    XLSX.utils.book_append_sheet(workbook, summarySheet, 'Resumen')

    // ==================== Hoja 2: Historial de Eventos ====================
    const events = getAllEvents()
    const eventsData = [
      ['HISTORIAL COMPLETO DE EVENTOS'],
      [''],
      ['Fecha y Hora', 'Tipo de Evento', 'Estado', 'Descripción', 'Usuario', 'Ubicación', 'Notas']
    ]
    
    events.forEach(event => {
      eventsData.push([
        formatDateTime(event.date_time || event.timestamp),
        getEventTitle(event),
        event.status || 'N/A',
        getEventDescription(event),
        event.user_name || 'N/A',
        event.location || 'N/A',
        event.notes || event.observations || 'N/A'
      ])
    })
    
    const eventsSheet = XLSX.utils.aoa_to_sheet(eventsData)
    
    // Configurar anchos de columna
    eventsSheet['!cols'] = [
      { wch: 18 }, // Fecha y Hora
      { wch: 25 }, // Tipo de Evento
      { wch: 20 }, // Estado
      { wch: 40 }, // Descripción
      { wch: 25 }, // Usuario
      { wch: 35 }, // Ubicación
      { wch: 40 }  // Notas
    ]
    
    // Mergear título
    if (!eventsSheet['!merges']) eventsSheet['!merges'] = []
    eventsSheet['!merges'].push({ s: { r: 0, c: 0 }, e: { r: 0, c: 6 } })
    
    XLSX.utils.book_append_sheet(workbook, eventsSheet, 'Historial de Eventos')

    // ==================== Hoja 3: Recorrido del Insumo ====================
    const journey = getLocationJourney()
    const journeyData = [
      ['RECORRIDO DEL INSUMO POR UBICACIONES'],
      [''],
      ['Secuencia', 'Ubicación', 'Fecha de Llegada', 'Duración en Ubicación', 'Estado Actual']
    ]
    
    journey.forEach((location, index) => {
      journeyData.push([
        index + 1,
        location.name,
        formatDateTime(location.date),
        location.duration || 'Ubicación actual',
        location.is_current ? '✓ UBICACION ACTUAL' : 'Ubicación anterior'
      ])
    })
    
    const journeySheet = XLSX.utils.aoa_to_sheet(journeyData)
    
    // Configurar anchos de columna
    journeySheet['!cols'] = [
      { wch: 12 }, // Secuencia
      { wch: 45 }, // Ubicación
      { wch: 20 }, // Fecha de Llegada
      { wch: 25 }, // Duración
      { wch: 22 }  // Estado
    ]
    
    // Mergear título
    if (!journeySheet['!merges']) journeySheet['!merges'] = []
    journeySheet['!merges'].push({ s: { r: 0, c: 0 }, e: { r: 0, c: 4 } })
    
    XLSX.utils.book_append_sheet(workbook, journeySheet, 'Recorrido')

    // ==================== Hoja 4: Resumen por Usuario ====================
    const userSummaryData = [
      ['RESUMEN DE ACTIVIDAD POR USUARIO'],
      [''],
      ['Usuario', 'Número de Eventos', 'Último Evento', 'Fecha del Último Evento']
    ]
    
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
      userSummaryData.push([
        userName,
        userEventList.length,
        getEventTitle(lastEvent),
        formatDateTime(lastEvent.date_time || lastEvent.timestamp)
      ])
    })
    
    const userSummarySheet = XLSX.utils.aoa_to_sheet(userSummaryData)
    
    // Configurar anchos de columna
    userSummarySheet['!cols'] = [
      { wch: 30 }, // Usuario
      { wch: 20 }, // Número de Eventos
      { wch: 35 }, // Último Evento
      { wch: 20 }  // Fecha
    ]
    
    // Mergear título
    if (!userSummarySheet['!merges']) userSummarySheet['!merges'] = []
    userSummarySheet['!merges'].push({ s: { r: 0, c: 0 }, e: { r: 0, c: 3 } })
    
    XLSX.utils.book_append_sheet(workbook, userSummarySheet, 'Resumen por Usuario')

    // Generar archivo Excel
    const fileName = `Trazabilidad_${qrCode.value}_${format(new Date(), 'yyyy-MM-dd_HHmm')}.xlsx`
    XLSX.writeFile(workbook, fileName)
    
    showSuccess('Trazabilidad exportada a Excel exitosamente')
  } catch (err) {
    console.error('Error al exportar trazabilidad:', err)
    showError('Error al exportar la trazabilidad: ' + err.message)
  }
}

const printTraceability = () => {
  try {
    const doc = new jsPDF()
    const pageWidth = doc.internal.pageSize.width
    const margin = 15
    const contentWidth = pageWidth - (margin * 2)
    let yPos = 15

    // === ENCABEZADO ===
    // Logo/Título
    doc.setFillColor(37, 99, 235) // bg-blue-600
    doc.rect(0, 0, pageWidth, 35, 'F')
    doc.setTextColor(255, 255, 255)
    doc.setFontSize(20)
    doc.setFont('helvetica', 'bold')
    doc.text('REPORTE DE TRAZABILIDAD', margin, 15)
    doc.setFontSize(10)
    doc.setFont('helvetica', 'normal')
    doc.text(`Generado: ${format(new Date(), 'dd/MM/yyyy HH:mm', { locale: es })}`, margin, 25)
    
    yPos = 45

    // === INFORMACIÓN DEL INSUMO ===
    doc.setTextColor(0, 0, 0)
    doc.setFillColor(243, 244, 246) // bg-gray-100
    doc.roundedRect(margin, yPos, contentWidth, 50, 2, 2, 'F')
    
    yPos += 8
    doc.setFontSize(12)
    doc.setFont('helvetica', 'bold')
    doc.text(qrInfo.value?.supply_info?.name || 'Insumo Médico', margin + 5, yPos)
    
    yPos += 7
    doc.setFontSize(8)
    doc.setFont('helvetica', 'normal')
    doc.setTextColor(75, 85, 99) // text-gray-600
    doc.text(`QR: ${qrCode.value}`, margin + 5, yPos)
    
    yPos += 10
    // Estado y ubicación en dos columnas
    const colWidth = contentWidth / 2
    
    doc.setFont('helvetica', 'bold')
    doc.setFontSize(7)
    doc.setTextColor(107, 114, 128)
    doc.text('ESTADO ACTUAL', margin + 5, yPos)
    
    yPos += 5
    doc.setFont('helvetica', 'normal')
    doc.setFontSize(9)
    doc.setTextColor(0, 0, 0)
    doc.text(getCurrentStatusLabel(), margin + 5, yPos)
    const locationText = doc.splitTextToSize(getCurrentLocation(), colWidth - 10)
    doc.text(locationText, margin + 5 + colWidth, yPos)
    
    yPos += Math.max(5, locationText.length * 4)
    
    doc.setFont('helvetica', 'bold')
    doc.setFontSize(7)
    doc.setTextColor(107, 114, 128)
    doc.text('TIEMPO EN SISTEMA', margin + 5, yPos)
    doc.text('CODIGO PROVEEDOR', margin + 5 + colWidth, yPos)
    
    yPos += 5
    doc.setFont('helvetica', 'normal')
    doc.setFontSize(9)
    doc.setTextColor(0, 0, 0)
    doc.text(getTimeInSystem(), margin + 5, yPos)
    doc.text(String(qrInfo.value?.supply_code?.code_supplier || 'N/A'), margin + 5 + colWidth, yPos)
    
    yPos += 15

    // === ESTADÍSTICAS ===
    doc.setFillColor(59, 130, 246) // bg-blue-600
    doc.rect(margin, yPos, contentWidth, 6, 'F')
    doc.setTextColor(255, 255, 255)
    doc.setFontSize(9)
    doc.setFont('helvetica', 'bold')
    doc.text('ESTADISTICAS', margin + 3, yPos + 4)
    
    yPos += 10
    doc.setTextColor(0, 0, 0)
    doc.setFontSize(8)
    doc.setFont('helvetica', 'normal')
    
    const statWidth = contentWidth / 4
    const stats = [
      { label: 'Eventos', value: getAllEvents().length },
      { label: 'Movimientos', value: getEventCount('movement') },
      { label: 'Ubicaciones', value: getLocationCount() },
      { label: 'Usuarios', value: getUserCount() }
    ]
    
    stats.forEach((stat, i) => {
      const xPos = margin + (i * statWidth)
      doc.setFont('helvetica', 'bold')
      doc.setFontSize(14)
      doc.text(stat.value.toString(), xPos + statWidth/2, yPos + 5, { align: 'center' })
      doc.setFont('helvetica', 'normal')
      doc.setFontSize(7)
      doc.setTextColor(107, 114, 128)
      doc.text(stat.label, xPos + statWidth/2, yPos + 10, { align: 'center' })
      doc.setTextColor(0, 0, 0)
    })
    
    yPos += 20

    // === RECORRIDO ===
    doc.setFillColor(59, 130, 246)
    doc.rect(margin, yPos, contentWidth, 6, 'F')
    doc.setTextColor(255, 255, 255)
    doc.setFontSize(9)
    doc.setFont('helvetica', 'bold')
    doc.text('RECORRIDO DEL INSUMO', margin + 3, yPos + 4)
    
    yPos += 12
    const journey = getLocationJourney()
    
    journey.forEach((location, index) => {
      if (yPos > 260) {
        doc.addPage()
        yPos = 20
      }
      
      // Círculo numerado
      const circleColor = location.is_current ? [34, 197, 94] : [156, 163, 175]
      doc.setFillColor(circleColor[0], circleColor[1], circleColor[2])
      doc.circle(margin + 4, yPos + 2, 3, 'F')
      doc.setTextColor(255, 255, 255)
      doc.setFontSize(7)
      doc.setFont('helvetica', 'bold')
      doc.text((index + 1).toString(), margin + 4, yPos + 3, { align: 'center' })
      
      // Línea conectora
      if (index < journey.length - 1) {
        doc.setDrawColor(209, 213, 219)
        doc.setLineWidth(0.5)
        doc.line(margin + 4, yPos + 5, margin + 4, yPos + 18)
      }
      
      // Contenido
      doc.setTextColor(0, 0, 0)
      doc.setFont('helvetica', 'bold')
      doc.setFontSize(8)
      const locText = doc.splitTextToSize(location.name, contentWidth - 15)
      doc.text(locText, margin + 10, yPos + 2)
      
      doc.setFont('helvetica', 'normal')
      doc.setFontSize(7)
      doc.setTextColor(107, 114, 128)
      doc.text(formatDateTime(location.date), margin + 10, yPos + 7)
      
      if (location.duration) {
        doc.text(`Duración: ${location.duration}`, margin + 10, yPos + 11)
      }
      
      if (location.is_current) {
        doc.setTextColor(34, 197, 94)
        doc.setFont('helvetica', 'bold')
        doc.text('• UBICACION ACTUAL', margin + 10, yPos + location.duration ? 15 : 11)
      }
      
      yPos += 18
    })
    
    yPos += 5

    // === HISTORIAL DE EVENTOS ===
    if (yPos > 200) {
      doc.addPage()
      yPos = 20
    }
    
    doc.setFillColor(59, 130, 246)
    doc.rect(margin, yPos, contentWidth, 6, 'F')
    doc.setTextColor(255, 255, 255)
    doc.setFontSize(9)
    doc.setFont('helvetica', 'bold')
    doc.text('HISTORIAL DETALLADO', margin + 3, yPos + 4)
    
    yPos += 12
    const events = getAllEvents()
    
    events.forEach((event, index) => {
      if (yPos > 260) {
        doc.addPage()
        yPos = 20
      }
      
      // Fondo alternado
      if (index % 2 === 0) {
        doc.setFillColor(249, 250, 251)
        doc.rect(margin, yPos - 2, contentWidth, 20, 'F')
      }
      
      // Fecha
      doc.setTextColor(107, 114, 128)
      doc.setFontSize(7)
      doc.setFont('helvetica', 'normal')
      doc.text(formatDateTime(event.date_time || event.timestamp), margin + 2, yPos + 1)
      
      // Título del evento
      doc.setTextColor(0, 0, 0)
      doc.setFontSize(8)
      doc.setFont('helvetica', 'bold')
      doc.text(getEventTitle(event), margin + 2, yPos + 6)
      
      // Descripción
      doc.setFont('helvetica', 'normal')
      doc.setFontSize(7)
      doc.setTextColor(75, 85, 99)
      const desc = doc.splitTextToSize(getEventDescription(event), contentWidth - 4)
      doc.text(desc.slice(0, 1), margin + 2, yPos + 10)
      
      // Usuario
      if (event.user_name) {
        doc.setTextColor(107, 114, 128)
        doc.setFontSize(6)
        doc.text(`Usuario: ${event.user_name}`, margin + 2, yPos + 14)
      }
      
      yPos += 20
    })

    // Pie de página en todas las páginas
    const pageCount = doc.internal.getNumberOfPages()
    for (let i = 1; i <= pageCount; i++) {
      doc.setPage(i)
      doc.setFontSize(7)
      doc.setTextColor(156, 163, 175)
      doc.text(`Página ${i} de ${pageCount}`, pageWidth / 2, 285, { align: 'center' })
      doc.text('MediTrack - Sistema de Trazabilidad', pageWidth / 2, 290, { align: 'center' })
    }

    // Guardar PDF
    const fileName = `Trazabilidad_${qrCode.value}_${format(new Date(), 'yyyy-MM-dd_HHmm')}.pdf`
    doc.save(fileName)
    
    showSuccess('PDF generado exitosamente')
  } catch (err) {
    console.error('Error al generar PDF:', err)
    showError('Error al generar el PDF: ' + err.message)
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