<template>
  <div class="batch-history-container">
    <!-- Header -->
    <div class="bg-white rounded-lg shadow-sm border p-6 mb-6">
      <div class="flex items-center justify-between">
        <h2 class="text-2xl font-semibold text-gray-900 flex items-center">
          <svg class="h-8 w-8 mr-3 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          Historial de Movimientos - Lote {{ batchId }}
        </h2>
        <button @click="loadHistory" class="btn-primary">
          <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          Actualizar
        </button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="text-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
      <p class="mt-4 text-gray-600">Cargando historial del lote...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-lg p-6 text-center">
      <svg class="h-12 w-12 text-red-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.732-.833-2.5 0L4.268 15.5c-.77.833.192 2.5 1.732 2.5z" />
      </svg>
      <h3 class="text-lg font-medium text-red-900 mb-2">Error al cargar historial</h3>
      <p class="text-red-700">{{ error }}</p>
      <button @click="loadHistory" class="btn-danger mt-4">
        Intentar de nuevo
      </button>
    </div>

    <!-- History List -->
    <div v-else-if="historyData.length > 0" class="space-y-4">
      <!-- Filtros -->
      <div class="bg-white rounded-lg shadow-sm border p-4 mb-6">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Filtrar por usuario:</label>
            <input 
              v-model="userFilter" 
              type="text" 
              placeholder="RUT o nombre..."
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Tipo de cambio:</label>
            <select 
              v-model="changeTypeFilter"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="">Todos los cambios</option>
              <option value="Cantidad actualizada">Consumos</option>
              <option value="Lote creado">Creación</option>
              <option value="Lote actualizado">Actualizaciones</option>
            </select>
          </div>
          <div class="flex items-end">
            <button @click="clearFilters" class="px-4 py-2 bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-lg transition-colors">
              Limpiar filtros
            </button>
          </div>
        </div>
      </div>

      <!-- Historial -->
      <div class="bg-white rounded-lg shadow-sm border overflow-hidden">
        <div class="divide-y divide-gray-200">
          <div 
            v-for="(entry, index) in filteredHistory" 
            :key="index"
            class="p-6 hover:bg-gray-50 transition-colors"
          >
            <!-- Formato específico que pidió el usuario -->
            <div class="flex items-start justify-between">
              <div class="flex-1">
                <!-- Fecha -->
                <div class="text-lg font-semibold text-gray-900 mb-2">
                  {{ entry.display_format.date }}
                </div>
                
                <!-- Acción principal -->
                <div class="text-base font-medium text-blue-600 mb-2">
                  {{ entry.display_format.action }}
                </div>
                
                <!-- Cambio de cantidad (si aplica) -->
                <div v-if="entry.display_format.previous_amount && entry.display_format.new_amount" class="space-y-1 mb-3">
                  <div class="text-sm text-gray-600">
                    <span class="font-medium">Cant:</span> 
                    <span class="text-red-600">{{ entry.display_format.previous_amount }}</span>
                    <span class="mx-2">→</span>
                    <span class="text-green-600">{{ entry.display_format.new_amount }}</span>
                  </div>
                </div>
                
                <!-- Usuario responsable -->
                <div class="flex items-center space-x-4 text-sm text-gray-700">
                  <span class="font-mono">{{ entry.display_format.user_rut }}</span>
                  <span>{{ entry.display_format.user_name }}</span>
                </div>
                
                <!-- Destino (si aplica) -->
                <div v-if="entry.display_format.destination" class="mt-2 text-sm text-gray-600">
                  <span class="font-medium">Destino:</span> {{ entry.display_format.destination }}
                </div>
                
                <!-- Observaciones (si aplica) -->
                <div v-if="entry.display_format.observations" class="mt-2 text-sm text-gray-500 italic">
                  {{ entry.display_format.observations }}
                </div>
              </div>
              
              <!-- Indicador de tipo de cambio -->
              <div class="ml-4 flex-shrink-0">
                <span 
                  :class="getChangeTypeBadgeClass(entry.change_details)"
                  class="inline-flex px-2 py-1 text-xs font-medium rounded-full"
                >
                  {{ getChangeTypeLabel(entry.change_details) }}
                </span>
              </div>
            </div>
            
            <!-- Detalles técnicos (expandible) -->
            <div v-if="showTechnicalDetails[index]" class="mt-4 pt-4 border-t border-gray-100">
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-xs text-gray-500">
                <div>
                  <strong>Valores anteriores:</strong>
                  <pre class="mt-1 bg-gray-50 p-2 rounded">{{ formatJSON(entry.previous_values) }}</pre>
                </div>
                <div>
                  <strong>Valores nuevos:</strong>
                  <pre class="mt-1 bg-gray-50 p-2 rounded">{{ formatJSON(entry.new_values) }}</pre>
                </div>
              </div>
            </div>
            
            <!-- Toggle para detalles técnicos -->
            <button 
              @click="toggleTechnicalDetails(index)"
              class="mt-2 text-xs text-blue-600 hover:text-blue-800"
            >
              {{ showTechnicalDetails[index] ? 'Ocultar' : 'Mostrar' }} detalles técnicos
            </button>
          </div>
        </div>
      </div>
      
      <!-- Paginación -->
      <div v-if="totalPages > 1" class="flex justify-center space-x-2 mt-6">
        <button 
          @click="currentPage--" 
          :disabled="currentPage === 1"
          class="px-3 py-2 text-sm bg-gray-100 hover:bg-gray-200 disabled:opacity-50 disabled:cursor-not-allowed rounded-lg"
        >
          Anterior
        </button>
        <span class="px-3 py-2 text-sm bg-blue-100 text-blue-800 rounded-lg">
          {{ currentPage }} / {{ totalPages }}
        </span>
        <button 
          @click="currentPage++" 
          :disabled="currentPage === totalPages"
          class="px-3 py-2 text-sm bg-gray-100 hover:bg-gray-200 disabled:opacity-50 disabled:cursor-not-allowed rounded-lg"
        >
          Siguiente
        </button>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else class="text-center py-12">
      <svg class="h-12 w-12 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
      </svg>
      <h3 class="text-lg font-medium text-gray-900 mb-2">Sin historial disponible</h3>
      <p class="text-gray-600">Este lote aún no tiene movimientos registrados.</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useNotification } from '@/composables/useNotification'
import qrService from '@/services/qr/qrService'

// Props
const props = defineProps({
  batchId: {
    type: [String, Number],
    required: true
  }
})

const route = useRoute()
const { success: showSuccess, error: showError } = useNotification()

// Estado reactivo
const loading = ref(false)
const error = ref(null)
const historyData = ref([])
const showTechnicalDetails = ref({})

// Filtros
const userFilter = ref('')
const changeTypeFilter = ref('')
const currentPage = ref(1)
const itemsPerPage = 10

// Computed
const batchId = computed(() => props.batchId || route.params.batchId)

const filteredHistory = computed(() => {
  let filtered = [...historyData.value]
  
  // Filtro por usuario
  if (userFilter.value) {
    const filter = userFilter.value.toLowerCase()
    filtered = filtered.filter(entry => 
      entry.user_rut?.toLowerCase().includes(filter) ||
      entry.user_name?.toLowerCase().includes(filter)
    )
  }
  
  // Filtro por tipo de cambio
  if (changeTypeFilter.value) {
    filtered = filtered.filter(entry => 
      entry.change_details === changeTypeFilter.value
    )
  }
  
  // Paginación
  const start = (currentPage.value - 1) * itemsPerPage
  const end = start + itemsPerPage
  
  return filtered.slice(start, end)
})

const totalPages = computed(() => {
  let filtered = [...historyData.value]
  
  if (userFilter.value) {
    const filter = userFilter.value.toLowerCase()
    filtered = filtered.filter(entry => 
      entry.user_rut?.toLowerCase().includes(filter) ||
      entry.user_name?.toLowerCase().includes(filter)
    )
  }
  
  if (changeTypeFilter.value) {
    filtered = filtered.filter(entry => 
      entry.change_details === changeTypeFilter.value
    )
  }
  
  return Math.ceil(filtered.length / itemsPerPage)
})

// Métodos
const loadHistory = async () => {
  if (!batchId.value) {
    error.value = 'ID de lote no proporcionado'
    return
  }

  loading.value = true
  error.value = null

  try {
    const history = await qrService.getBatchHistoryFormatted(batchId.value)
    historyData.value = history.sort((a, b) => new Date(b.date_time) - new Date(a.date_time))
    if (history.length > 0) {
      showSuccess(`Historial cargado: ${history.length} movimientos encontrados`)
    }
  } catch (err) {
    console.error('Error al cargar historial:', err)
    const errorMessage = err.message || 'Error al cargar el historial del lote'
    error.value = errorMessage
    showError(errorMessage)
  } finally {
    loading.value = false
  }
}

const clearFilters = () => {
  userFilter.value = ''
  changeTypeFilter.value = ''
  currentPage.value = 1
}

const toggleTechnicalDetails = (index) => {
  showTechnicalDetails.value[index] = !showTechnicalDetails.value[index]
}

const getChangeTypeBadgeClass = (changeType) => {
  switch (changeType) {
    case 'Cantidad actualizada':
      return 'bg-orange-100 text-orange-800'
    case 'Lote creado':
      return 'bg-green-100 text-green-800'
    case 'Lote actualizado':
      return 'bg-blue-100 text-blue-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}

const getChangeTypeLabel = (changeType) => {
  switch (changeType) {
    case 'Cantidad actualizada':
      return 'Consumo'
    case 'Lote creado':
      return 'Creación'
    case 'Lote actualizado':
      return 'Actualización'
    default:
      return changeType
  }
}

const formatJSON = (jsonString) => {
  try {
    const parsed = JSON.parse(jsonString || '{}')
    return JSON.stringify(parsed, null, 2)
  } catch {
    return jsonString || '{}'
  }
}

// Watchers
watch([userFilter, changeTypeFilter], () => {
  currentPage.value = 1
})

// Lifecycle
onMounted(() => {
  loadHistory()
})
</script>

<style scoped>
.batch-history-container {
  max-width: 4xl;
  margin: 0 auto;
  padding: 1rem;
}

pre {
  white-space: pre-wrap;
  word-break: break-word;
  font-size: 0.75rem;
  line-height: 1.2;
}
</style>