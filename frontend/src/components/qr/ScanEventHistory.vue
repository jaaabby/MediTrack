<template>
  <div class="scan-event-history">
    <!-- Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h3 class="text-lg font-medium text-gray-900">Historial de Escaneos</h3>
        <p class="text-sm text-gray-600 mt-1">
          {{ filteredEvents.length }} eventos de escaneo
          <span v-if="hasActiveFilters">(filtrados)</span>
        </p>
      </div>
      
      <div class="flex space-x-2">
        <button
          v-if="showFilters"
          @click="toggleFiltersVisibility"
          class="btn-secondary text-sm"
        >
          <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.207A1 1 0 013 6.5V4z" />
          </svg>
          {{ filtersVisible ? 'Ocultar' : 'Mostrar' }} Filtros
        </button>
        
        <button
          @click="exportEvents"
          :disabled="filteredEvents.length === 0"
          class="btn-primary text-sm"
        >
          <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          Exportar
        </button>
      </div>
    </div>

    <!-- Filtros Avanzados -->
    <div v-if="showFilters && filtersVisible" class="bg-gray-50 rounded-lg p-4 mb-6">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Período</label>
          <select v-model="filters.period" class="form-select w-full text-sm">
            <option value="">Todo el tiempo</option>
            <option value="today">Hoy</option>
            <option value="yesterday">Ayer</option>
            <option value="week">Última semana</option>
            <option value="month">Último mes</option>
            <option value="custom">Personalizado</option>
          </select>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Usuario</label>
          <select v-model="filters.user" class="form-select w-full text-sm">
            <option value="">Todos los usuarios</option>
            <option
              v-for="user in uniqueUsers"
              :key="user"
              :value="user"
            >
              {{ user }}
            </option>
          </select>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Propósito</label>
          <select v-model="filters.purpose" class="form-select w-full text-sm">
            <option value="">Todos los propósitos</option>
            <option value="lookup">Consulta</option>
            <option value="consume">Consumo</option>
            <option value="verify">Verificación</option>
            <option value="inventory_check">Inventario</option>
            <option value="assign">Asignación</option>
          </select>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Resultado</label>
          <select v-model="filters.result" class="form-select w-full text-sm">
            <option value="">Todos los resultados</option>
            <option value="success">Exitoso</option>
            <option value="error">Error</option>
            <option value="not_found">No encontrado</option>
          </select>
        </div>
      </div>

      <!-- Filtros de fecha personalizada -->
      <div v-if="filters.period === 'custom'" class="mt-4 grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Fecha desde</label>
          <input
            v-model="filters.dateFrom"
            type="datetime-local"
            class="form-input w-full text-sm"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Fecha hasta</label>
          <input
            v-model="filters.dateTo"
            type="datetime-local"
            class="form-input w-full text-sm"
          />
        </div>
      </div>

      <!-- Botones de filtro -->
      <div class="mt-4 flex space-x-2">
        <button
          @click="clearFilters"
          class="btn-secondary text-sm"
        >
          Limpiar Filtros
        </button>
        <button
          @click="applyQuickFilter('today')"
          class="btn-secondary text-sm"
        >
          Solo Hoy
        </button>
        <button
          @click="applyQuickFilter('errors')"
          class="btn-secondary text-sm"
        >
          Solo Errores
        </button>
      </div>
    </div>

    <!-- Lista de Eventos -->
    <div class="space-y-4">
      <div
        v-for="(event, index) in paginatedEvents"
        :key="index"
        class="bg-white border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow"
      >
        <div class="flex items-start justify-between">
          <div class="flex-1">
            <!-- Header del evento -->
            <div class="flex items-center space-x-3 mb-2">
              <div
                :class="[
                  'w-2 h-2 rounded-full',
                  event.scan_result === 'success' ? 'bg-green-500' :
                  event.scan_result === 'error' ? 'bg-red-500' :
                  'bg-yellow-500'
                ]"
              ></div>
              
              <span class="text-sm font-medium text-gray-900">
                {{ formatDate(event.scanned_at) }}
              </span>
              
              <span
                :class="[
                  'px-2 py-1 rounded-full text-xs font-medium',
                  getScanPurposeBadgeClass(event.scan_purpose)
                ]"
              >
                {{ getScanPurposeLabel(event.scan_purpose) }}
              </span>
              
              <span
                :class="[
                  'px-2 py-1 rounded-full text-xs font-medium',
                  getScanSourceBadgeClass(event.scan_source)
                ]"
              >
                {{ event.scan_source }}
              </span>
            </div>

            <!-- Información del usuario y ubicación -->
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm text-gray-600">
              <div>
                <div class="flex items-center">
                  <svg class="h-4 w-4 mr-1 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                  </svg>
                  <span>{{ event.scanned_by_name || 'Usuario no identificado' }}</span>
                </div>
                <div v-if="event.scanned_by_rut" class="text-xs text-gray-500 ml-5">
                  RUT: {{ event.scanned_by_rut }}
                </div>
              </div>

              <div v-if="event.current_location || event.pavilion_name">
                <div class="flex items-center">
                  <svg class="h-4 w-4 mr-1 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                  </svg>
                  <span>{{ event.current_location || event.pavilion_name }}</span>
                </div>
                <div v-if="event.medical_center_name" class="text-xs text-gray-500 ml-5">
                  {{ event.medical_center_name }}
                </div>
              </div>
            </div>

            <!-- Información técnica -->
            <div class="mt-3 text-xs text-gray-500">
              <div class="grid grid-cols-1 md:grid-cols-3 gap-2">
                <div v-if="event.ip_address">
                  IP: {{ event.ip_address }}
                </div>
                <div v-if="event.device_info?.platform">
                  Plataforma: {{ event.device_info.platform }}
                </div>
                <div v-if="event.scan_sequence">
                  Secuencia: #{{ event.scan_sequence }}
                </div>
              </div>
            </div>

            <!-- Información del producto escaneado -->
            <div v-if="event.supply_name" class="mt-3 p-2 bg-gray-50 rounded text-sm">
              <div class="font-medium text-gray-900">{{ event.supply_name }}</div>
              <div v-if="event.current_status" class="text-gray-600">
                Estado: {{ event.current_status }}
              </div>
            </div>

            <!-- Error message si existe -->
            <div v-if="event.error_message" class="mt-3 p-2 bg-red-50 border border-red-200 rounded text-sm">
              <div class="text-red-800 font-medium">Error:</div>
              <div class="text-red-700">{{ event.error_message }}</div>
            </div>
          </div>

          <!-- Acciones -->
          <div class="flex space-x-2 ml-4">
            <button
              @click="viewEventDetails(event)"
              class="text-blue-600 hover:text-blue-800"
              title="Ver detalles"
            >
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
              </svg>
            </button>

            <button
              v-if="event.qr_code"
              @click="scanAgain(event.qr_code)"
              class="text-green-600 hover:text-green-800"
              title="Escanear nuevamente"
            >
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Paginación -->
    <div v-if="totalPages > 1" class="mt-6 flex items-center justify-between">
      <div class="text-sm text-gray-700">
        Mostrando {{ startIndex + 1 }} a {{ Math.min(endIndex, filteredEvents.length) }} de {{ filteredEvents.length }} eventos
      </div>
      
      <div class="flex space-x-2">
        <button
          @click="currentPage--"
          :disabled="currentPage === 1"
          class="btn-secondary text-sm"
        >
          Anterior
        </button>
        
        <span class="px-3 py-2 text-sm text-gray-700">
          Página {{ currentPage }} de {{ totalPages }}
        </span>
        
        <button
          @click="currentPage++"
          :disabled="currentPage === totalPages"
          class="btn-secondary text-sm"
        >
          Siguiente
        </button>
      </div>
    </div>

    <!-- Sin eventos -->
    <div v-if="filteredEvents.length === 0" class="text-center py-12">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900">No hay eventos de escaneo</h3>
      <p class="mt-1 text-sm text-gray-500">
        {{ hasActiveFilters ? 'No se encontraron eventos con los filtros aplicados.' : 'No se han registrado escaneos para este código QR.' }}
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { format, isToday, isYesterday, isThisWeek, isThisMonth, startOfDay, endOfDay } from 'date-fns'
import { es } from 'date-fns/locale'
import { useNotification } from '@/composables/useNotification'

// Props
const props = defineProps({
  events: {
    type: Array,
    default: () => []
  },
  showFilters: {
    type: Boolean,
    default: true
  },
  itemsPerPage: {
    type: Number,
    default: 10
  }
})

// Emits
const emit = defineEmits(['view-details', 'scan-again', 'export'])

// Notificaciones
const { success: showSuccess, error: showError } = useNotification()

// Estado reactivo
const filtersVisible = ref(false)
const currentPage = ref(1)

// Filtros
const filters = ref({
  period: '',
  user: '',
  purpose: '',
  result: '',
  dateFrom: '',
  dateTo: ''
})

// Computed properties
const filteredEvents = computed(() => {
  let events = props.events

  // Filtrar por período
  if (filters.value.period) {
    events = events.filter(event => {
      const eventDate = new Date(event.scanned_at)
      switch (filters.value.period) {
        case 'today':
          return isToday(eventDate)
        case 'yesterday':
          return isYesterday(eventDate)
        case 'week':
          return isThisWeek(eventDate)
        case 'month':
          return isThisMonth(eventDate)
        case 'custom':
          if (filters.value.dateFrom && filters.value.dateTo) {
            const from = new Date(filters.value.dateFrom)
            const to = new Date(filters.value.dateTo)
            return eventDate >= from && eventDate <= to
          }
          return true
        default:
          return true
      }
    })
  }

  // Filtrar por usuario
  if (filters.value.user) {
    events = events.filter(event => event.scanned_by_name === filters.value.user)
  }

  // Filtrar por propósito
  if (filters.value.purpose) {
    events = events.filter(event => event.scan_purpose === filters.value.purpose)
  }

  // Filtrar por resultado
  if (filters.value.result) {
    events = events.filter(event => event.scan_result === filters.value.result)
  }

  return events
})

const paginatedEvents = computed(() => {
  const start = (currentPage.value - 1) * props.itemsPerPage
  const end = start + props.itemsPerPage
  return filteredEvents.value.slice(start, end)
})

const totalPages = computed(() => {
  return Math.ceil(filteredEvents.value.length / props.itemsPerPage)
})

const startIndex = computed(() => {
  return (currentPage.value - 1) * props.itemsPerPage
})

const endIndex = computed(() => {
  return startIndex.value + props.itemsPerPage
})

const hasActiveFilters = computed(() => {
  return Object.values(filters.value).some(value => value !== '')
})

const uniqueUsers = computed(() => {
  const users = new Set()
  props.events.forEach(event => {
    if (event.scanned_by_name) users.add(event.scanned_by_name)
  })
  return Array.from(users).sort()
})

// Watchers
watch(() => filteredEvents.value.length, () => {
  currentPage.value = 1
})

// Funciones
const toggleFiltersVisibility = () => {
  filtersVisible.value = !filtersVisible.value
}

const clearFilters = () => {
  filters.value = {
    period: '',
    user: '',
    purpose: '',
    result: '',
    dateFrom: '',
    dateTo: ''
  }
}

const applyQuickFilter = (type) => {
  clearFilters()
  switch (type) {
    case 'today':
      filters.value.period = 'today'
      break
    case 'errors':
      filters.value.result = 'error'
      break
  }
}

const viewEventDetails = (event) => {
  emit('view-details', event)
}

const scanAgain = (qrCode) => {
  emit('scan-again', qrCode)
}

const exportEvents = () => {
  if (filteredEvents.value.length === 0) {
    showError('No hay eventos para exportar')
    return
  }
  
  try {
    const data = {
      export_date: new Date().toISOString(),
      total_events: filteredEvents.value.length,
      filters_applied: hasActiveFilters.value ? filters.value : null,
      events: filteredEvents.value
    }
    
    const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `scan-events-${format(new Date(), 'yyyy-MM-dd-HHmm')}.json`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
    
    showSuccess(`${filteredEvents.value.length} eventos exportados exitosamente`)
    emit('export', data)
  } catch (err) {
    console.error('Error al exportar eventos:', err)
    showError('Error al exportar los eventos')
  }
}

// Funciones de utilidad
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
  return labels[purpose] || purpose || 'Sin propósito'
}

const getScanSourceBadgeClass = (source) => {
  const classes = {
    'web': 'bg-blue-100 text-blue-800',
    'mobile': 'bg-green-100 text-green-800',
    'api': 'bg-purple-100 text-purple-800',
    'scanner': 'bg-orange-100 text-orange-800'
  }
  return classes[source] || 'bg-gray-100 text-gray-800'
}

const formatDate = (dateString) => {
  if (!dateString) return 'No disponible'
  try {
    return format(new Date(dateString), 'dd/MM/yyyy HH:mm:ss', { locale: es })
  } catch (error) {
    return dateString
  }
}
</script>

<style scoped>
.form-input, .form-select {
  @apply block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500;
}

/* Usar clases de botones de style.css global */
</style>