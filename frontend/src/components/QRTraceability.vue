<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <div class="bg-white shadow-sm border-b">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between h-16">
          <div class="flex items-center space-x-4">
            <router-link to="/qr" class="text-gray-400 hover:text-gray-600">
              <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
            </router-link>
            <h1 class="text-xl font-semibold text-gray-900">Trazabilidad del Insumo</h1>
          </div>
          <div class="flex items-center space-x-3">
            <button @click="exportTraceability" class="btn-secondary text-sm">
              <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              Exportar
            </button>
            <button @click="printTraceability" class="btn-secondary text-sm">
              <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z" />
              </svg>
              Imprimir
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex justify-center items-center py-12">
      <div class="text-center">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
        <p class="mt-4 text-gray-600">Cargando trazabilidad...</p>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="max-w-4xl mx-auto px-4 py-8">
      <div class="bg-red-50 border border-red-200 rounded-lg p-6">
        <div class="flex">
          <div class="flex-shrink-0">
            <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-red-800">Error al cargar trazabilidad</h3>
            <div class="mt-2 text-sm text-red-700">{{ error }}</div>
            <div class="mt-4">
              <button @click="loadTraceability" class="btn-primary text-sm">
                Reintentar
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div v-else-if="traceabilityData" class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Product Summary Card -->
      <div class="bg-white rounded-lg shadow-sm border mb-8">
        <div class="p-6">
          <div class="flex items-start justify-between">
            <div class="flex items-center space-x-4">
              <!-- QR Code Icon -->
              <div class="bg-blue-100 p-3 rounded-lg">
                <svg class="h-8 w-8 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v2a2 2 0 002 2zm0 0h2a2 2 0 012 2v2a2 2 0 01-2 2H5a2 2 0 01-2-2v-2a2 2 0 012-2z" />
                </svg>
              </div>
              
              <!-- Product Info -->
              <div>
                <h2 class="text-2xl font-bold text-gray-900">
                  {{ qrInfo?.supply_info?.name || 'Insumo Médico' }}
                </h2>
                <p class="text-gray-600 font-mono text-sm mt-1">{{ qrCode }}</p>
                <div class="flex items-center mt-2 space-x-4">
                  <span :class="getCurrentStatusBadgeClass()" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium">
                    {{ getCurrentStatusLabel() }}
                  </span>
                  <span class="text-sm text-gray-500">
                    En sistema desde {{ getTimeInSystem() }}
                  </span>
                </div>
              </div>
            </div>

            <!-- Current Location -->
            <div v-if="getCurrentLocation()" class="text-right">
              <div class="text-sm text-gray-500">Ubicación actual</div>
              <div class="text-lg font-semibold text-gray-900">{{ getCurrentLocation() }}</div>
            </div>
          </div>

          <!-- Stats Bar -->
          <div class="mt-6 grid grid-cols-2 sm:grid-cols-4 gap-4">
            <div class="bg-gray-50 rounded-lg p-4">
              <div class="text-2xl font-bold text-blue-600">{{ getEventCount('scan') }}</div>
              <div class="text-sm text-gray-600">Escaneos</div>
            </div>
            <div class="bg-gray-50 rounded-lg p-4">
              <div class="text-2xl font-bold text-green-600">{{ getEventCount('movement') }}</div>
              <div class="text-sm text-gray-600">Movimientos</div>
            </div>
            <div class="bg-gray-50 rounded-lg p-4">
              <div class="text-2xl font-bold text-purple-600">{{ getLocationCount() }}</div>
              <div class="text-sm text-gray-600">Ubicaciones</div>
            </div>
            <div class="bg-gray-50 rounded-lg p-4">
              <div class="text-2xl font-bold text-orange-600">{{ getUserCount() }}</div>
              <div class="text-sm text-gray-600">Usuarios</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Journey Map -->
      <div class="bg-white rounded-lg shadow-sm border mb-8">
        <div class="p-6">
          <h3 class="text-lg font-medium text-gray-900 mb-6">Recorrido del Insumo</h3>
          
          <!-- Progress Bar -->
          <div class="mb-8">
            <div class="flex items-center justify-between text-sm text-gray-500 mb-2">
              <span>Creado</span>
              <span>En Tránsito</span>
              <span>{{ qrInfo?.is_consumed ? 'Consumido' : 'Activo' }}</span>
            </div>
            <div class="w-full bg-gray-200 rounded-full h-2">
              <div 
                :class="[
                  'h-2 rounded-full transition-all duration-1000',
                  qrInfo?.is_consumed ? 'bg-red-500' : 'bg-blue-500'
                ]"
                :style="`width: ${getProgressPercentage()}%`"
              ></div>
            </div>
          </div>

          <!-- Location Journey -->
          <div v-if="getLocationJourney().length > 0" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            <div
              v-for="(location, index) in getLocationJourney()"
              :key="index"
              class="relative"
            >
              <div :class="[
                'p-4 rounded-lg border-2 transition-all duration-300',
                location.is_current ? 'border-blue-500 bg-blue-50' : 'border-gray-200 bg-white'
              ]">
                <div class="flex items-center justify-between mb-2">
                  <div :class="[
                    'flex items-center justify-center w-8 h-8 rounded-full text-sm font-medium',
                    location.is_current ? 'bg-blue-500 text-white' : 'bg-gray-200 text-gray-600'
                  ]">
                    {{ index + 1 }}
                  </div>
                  <div v-if="location.is_current" class="flex items-center text-blue-600 text-xs font-medium">
                    <div class="w-2 h-2 bg-blue-600 rounded-full mr-1 animate-pulse"></div>
                    Actual
                  </div>
                </div>
                <div class="text-sm font-medium text-gray-900">{{ location.name }}</div>
                <div class="text-xs text-gray-500 mt-1">{{ formatDateTime(location.date) }}</div>
                <div v-if="location.duration" class="text-xs text-gray-400 mt-1">
                  Duración: {{ location.duration }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Timeline -->
      <div class="bg-white rounded-lg shadow-sm border">
        <div class="p-6">
          <div class="flex items-center justify-between mb-6">
            <h3 class="text-lg font-medium text-gray-900">Historial Completo</h3>
            <div class="flex items-center space-x-3">
              <!-- Filters -->
              <select v-model="selectedEventFilter" class="form-select text-sm" @change="applyFilters">
                <option value="">Todos los eventos</option>
                <option value="scan">Solo escaneos</option>
                <option value="movement">Solo movimientos</option>
                <option value="consumption">Solo consumos</option>
                <option value="transfer">Solo transferencias</option>
              </select>
            </div>
          </div>

          <!-- Timeline Events -->
          <div class="flow-root">
            <ul class="-mb-8">
              <li v-for="(event, index) in filteredEvents" :key="index">
                <div class="relative pb-8">
                  <!-- Connector Line -->
                  <span
                    v-if="index !== filteredEvents.length - 1"
                    class="absolute top-4 left-4 -ml-px h-full w-0.5 bg-gray-200"
                    aria-hidden="true"
                  ></span>
                  
                  <div class="relative flex space-x-3">
                    <!-- Event Icon -->
                    <div>
                      <span :class="[
                        'h-8 w-8 rounded-full flex items-center justify-center ring-8 ring-white',
                        getEventIconClass(event)
                      ]">
                        <svg class="h-4 w-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="getEventIcon(event)" />
                        </svg>
                      </span>
                    </div>
                    
                    <!-- Event Content -->
                    <div class="min-w-0 flex-1 pt-1.5">
                      <div>
                        <div class="flex items-center justify-between">
                          <div class="text-sm text-gray-900 font-medium">
                            {{ getEventTitle(event) }}
                          </div>
                          <div class="text-sm text-gray-500">
                            {{ formatRelativeTime(event.date_time || event.timestamp) }}
                          </div>
                        </div>
                        <div class="mt-1 text-sm text-gray-600">
                          {{ getEventDescription(event) }}
                        </div>
                        
                        <!-- Event Details -->
                        <div v-if="hasEventDetails(event)" class="mt-2">
                          <div class="bg-gray-50 rounded-md p-3 text-xs space-y-1">
                            <div v-if="event.user_name" class="flex justify-between">
                              <span class="text-gray-500">Usuario:</span>
                              <span class="text-gray-900">{{ event.user_name }}</span>
                            </div>
                            <div v-if="event.location" class="flex justify-between">
                              <span class="text-gray-500">Ubicación:</span>
                              <span class="text-gray-900">{{ event.location }}</span>
                            </div>
                            <div v-if="event.notes" class="flex justify-between">
                              <span class="text-gray-500">Notas:</span>
                              <span class="text-gray-900">{{ event.notes }}</span>
                            </div>
                            <div class="flex justify-between">
                              <span class="text-gray-500">Fecha exacta:</span>
                              <span class="text-gray-900">{{ formatDateTime(event.date_time || event.timestamp) }}</span>
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
import qrService from '@/services/qrService'

const route = useRoute()
const router = useRouter()

// Props
const props = defineProps({
  qrCode: {
    type: String,
    required: false
  }
})

// Estado del componente
const loading = ref(false)
const error = ref(null)
const traceabilityData = ref(null)
const qrInfo = ref(null)
const selectedEventFilter = ref('')

// Computed properties
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

// Métodos principales
const loadTraceability = async () => {
  loading.value = true
  error.value = null
  
  try {
    // Cargar trazabilidad completa
    const traceData = await qrService.getCompleteTraceability(qrCode.value)
    traceabilityData.value = traceData
    
    // Cargar información del QR
    const qrData = await qrService.scanQRCode(qrCode.value, {
      scan_purpose: 'traceability_view',
      scan_source: 'traceability_starken'
    })
    qrInfo.value = qrData
    
  } catch (err) {
    error.value = err.message || 'Error al cargar la trazabilidad'
    console.error('Error cargando trazabilidad:', err)
  } finally {
    loading.value = false
  }
}

// Funciones de datos
const getAllEvents = () => {
  if (!traceabilityData.value) return []

  const events = []

  // Agregar escaneos
  if (traceabilityData.value.scan_history) {
    events.push(...traceabilityData.value.scan_history.map(scan => ({
      ...scan,
      event_type: 'scan',
      date_time: scan.scanned_at || scan.timestamp,
      title: 'Código QR Escaneado',
      user_name: scan.scanned_by_name || scan.user_name
    })))
  }

  // Agregar movimientos
  if (traceabilityData.value.supply_history) {
    events.push(...traceabilityData.value.supply_history.map(movement => ({
      ...movement,
      event_type: 'movement',
      date_time: movement.timestamp || movement.date_time,
      title: getMovementTitle(movement),
      location: getMovementLocation(movement)
    })))
  }

  // Ordenar por fecha (más reciente primero)
  return events.sort((a, b) => new Date(b.date_time) - new Date(a.date_time))
}

const getMovementTitle = (movement) => {
  if (movement.status === 'consumido') return 'Insumo Consumido'
  if (movement.status === 'transferido') return 'Insumo Transferido'
  if (movement.status === 'creado') return 'Insumo Creado'
  return `Movimiento: ${movement.status || movement.movement_type || 'Cambio'}`
}

const getMovementLocation = (movement) => {
  if (movement.destination_name) {
    const type = movement.destination_type === 'pavilion' ? 'Pabellón' : 'Almacén'
    let location = `${type}: ${movement.destination_name}`
    if (movement.medical_center_name) {
      location += ` (${movement.medical_center_name})`
    }
    return location
  }
  return movement.location || 'Ubicación no especificada'
}

// Funciones de estado y estadísticas
const getCurrentStatusLabel = () => {
  if (!qrInfo.value) return 'Desconocido'
  
  if (qrInfo.value.is_consumed) return 'Consumido'
  
  const status = qrInfo.value.supply_info?.status || qrInfo.value.status
  switch (status) {
    case 'disponible': return 'Disponible'
    case 'recepcionado': return 'Recepcionado'
    case 'transferido': return 'Transferido'
    default: return status || 'Activo'
  }
}

const getCurrentStatusBadgeClass = () => {
  if (!qrInfo.value) return 'bg-gray-100 text-gray-800'
  
  if (qrInfo.value.is_consumed) return 'bg-red-100 text-red-800'
  
  const status = qrInfo.value.supply_info?.status || qrInfo.value.status
  switch (status) {
    case 'disponible': return 'bg-green-100 text-green-800'
    case 'recepcionado': return 'bg-blue-100 text-blue-800'
    case 'transferido': return 'bg-purple-100 text-purple-800'
    default: return 'bg-gray-100 text-gray-800'
  }
}

const getCurrentLocation = () => {
  const events = getAllEvents()
  const lastMovement = events.find(e => e.event_type === 'movement' && e.location)
  return lastMovement?.location || 'Sin ubicación registrada'
}

const getTimeInSystem = () => {
  if (!traceabilityData.value?.created_date && !qrInfo.value?.created_at) return 'N/A'
  
  try {
    const createdDate = new Date(traceabilityData.value?.created_date || qrInfo.value?.created_at)
    const now = new Date()
    
    const days = differenceInDays(now, createdDate)
    if (days > 0) {
      return `${days} día${days > 1 ? 's' : ''}`
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
  
  if (qrInfo.value.is_consumed) return 100
  
  // Calcular progreso basado en eventos
  const events = getAllEvents()
  const hasMovements = events.some(e => e.event_type === 'movement')
  
  if (hasMovements) return 75
  if (events.length > 0) return 50
  return 25
}

// Funciones de conteo
const getEventCount = (type) => {
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

// Funciones de recorrido
const getLocationJourney = () => {
  const movements = getAllEvents()
    .filter(e => e.event_type === 'movement' && e.location)
    .reverse() // Orden cronológico

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
            duration = `${days} día${days > 1 ? 's' : ''}`
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

// Funciones de eventos del timeline
const getEventTitle = (event) => {
  switch (event.event_type) {
    case 'scan':
      return 'Código QR Escaneado'
    case 'movement':
      return getMovementTitle(event)
    default:
      return 'Evento'
  }
}

const getEventDescription = (event) => {
  switch (event.event_type) {
    case 'scan':
      return `Escaneado por ${event.user_name || 'Usuario desconocido'} desde ${event.scan_source || 'aplicación'}`
    case 'movement':
      let description = event.observations || `Estado: ${event.status || 'Procesado'}`
      if (event.location) {
        description += ` - ${event.location}`
      }
      return description
    default:
      return 'Información no disponible'
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
  return event.user_name || event.location || event.notes || true // Siempre mostrar fecha exacta
}

// Funciones de formato
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

// Funciones de exportación
const exportTraceability = () => {
  const data = {
    qr_code: qrCode.value,
    export_date: new Date().toISOString(),
    product_info: qrInfo.value,
    traceability_data: traceabilityData.value,
    events: getAllEvents(),
    summary: {
      total_events: getAllEvents().length,
      total_scans: getEventCount('scan'),
      total_movements: getEventCount('movement'),
      locations_visited: getLocationCount(),
      users_involved: getUserCount(),
      current_status: getCurrentStatusLabel(),
      time_in_system: getTimeInSystem()
    }
  }

  const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `trazabilidad-${qrCode.value}-${format(new Date(), 'yyyy-MM-dd')}.json`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}

const printTraceability = () => {
  window.print()
}

const applyFilters = () => {
  // Los filtros se aplican automáticamente a través del computed
}

// Lifecycle hooks
onMounted(() => {
  if (qrCode.value) {
    loadTraceability()
  }
})
</script>

<style scoped>
.form-select {
  @apply block px-3 py-1 text-sm border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500;
}

.btn-primary {
  @apply inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500;
}

.btn-secondary {
  @apply inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md shadow-sm text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500;
}

.flow-root {
  overflow: hidden;
}

@media print {
  .no-print {
    display: none !important;
  }
  
  .print-only {
    display: block !important;
  }
}

/* Animaciones */
@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: .5;
  }
}

.animate-pulse {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}
</style>