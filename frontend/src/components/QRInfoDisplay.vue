<template>
  <div class="text-xs text-gray-500 mt-2">Valor seleccionado: {{ selectedPavilion }}</div>
  <div class="space-y-6">
    <!-- Información básica del QR con imagen -->
    <div class="bg-gray-50 rounded-lg p-6">
      <div class="flex flex-col lg:flex-row lg:items-start lg:space-x-6">

        <!-- Imagen QR -->
        <div class="flex-shrink-0 mb-4 lg:mb-0">
          <div class="bg-white rounded-lg p-4 shadow-sm border text-center">
            <img v-if="qrImageUrl" :src="qrImageUrl" :alt="`Código QR: ${qrInfo.qr_code}`"
              class="w-32 h-32 mx-auto object-contain" @error="handleImageError" @load="imageLoaded = true" />
            <div v-else-if="imageError" class="w-32 h-32 mx-auto bg-gray-100 rounded flex items-center justify-center">
              <svg class="h-8 w-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.732-.833-2.5 0L4.268 16.5c-.77.833.192 2.5 1.732 2.5z" />
              </svg>
            </div>
            <div v-else class="w-32 h-32 mx-auto bg-gray-100 rounded flex items-center justify-center animate-pulse">
              <svg class="h-8 w-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v14a2 2 0 002 2z" />
              </svg>
            </div>

            <!-- Solo botón de descarga -->
            <div v-if="imageLoaded" class="mt-3">
              <button @click="downloadQRImage('normal')"
                class="w-full text-xs bg-gray-600 hover:bg-gray-700 text-white px-2 py-1 rounded transition-colors">
                Descargar
              </button>
            </div>
          </div>
        </div>

        <!-- Información del Insumo y Acciones -->
        <div class="flex-1 space-y-4">
          <!-- Información del Insumo -->
          <div v-if="qrInfo.supply_code" class="bg-white rounded border p-4">
            <h4 class="font-semibold text-gray-900 flex items-center mb-3">
              <svg class="h-5 w-5 mr-2 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
              </svg>
              Información del Insumo
            </h4>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <div>
                <label class="text-sm font-medium text-gray-600">Código:</label>
                <p class="text-gray-900 font-semibold">{{ qrInfo.supply_code.code }}</p>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-600">Nombre:</label>
                <p class="text-gray-900">{{ qrInfo.supply_code.name }}</p>
              </div>
              <div class="sm:col-span-2">
                <label class="text-sm font-medium text-gray-600">Código de Proveedor:</label>
                <p class="text-gray-900">{{ qrInfo.supply_code.code_supplier }}</p>
              </div>
            </div>
          </div>

          <!-- Enviar a Pabellón -->
          <div class="bg-white rounded border p-4">
            <h4 class="font-semibold text-gray-900 flex items-center mb-3">
              <svg class="h-5 w-5 mr-2 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
              </svg>
              Enviar a Pabellón
            </h4>

            <div class="space-y-3">
              <!-- Búsqueda de Pabellón -->
              <!-- <div>
                <label class="text-sm font-medium text-gray-600 block mb-1">Buscar Pabellón:</label>
                <div class="relative">
                  <input 
                    type="text" 
                    v-model="pavilionSearch"
                    @input="filterPavilions"
                    placeholder="Buscar pabellón..."
                    class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  />
                  <svg class="absolute right-3 top-2.5 h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                  </svg>
                </div>
              </div>
            -->

              <!-- Selector de Pabellón -->
              <div>
                <label class="text-sm font-medium text-gray-600 block mb-1">Seleccionar Pabellón:</label>
                <div v-if="loadingPavilions" class="text-blue-500 text-sm py-2">Cargando pabellones...</div>
                <div v-else-if="pavilionsError" class="text-red-500 text-sm py-2">Error cargando pabellones: {{
                  pavilionsError }}</div>
                <select v-model="selectedPavilion"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
                  <option value="">Seleccione un pabellón</option>
                  <option v-for="pavilion in filteredPavilions" :key="pavilion.id" :value="String(pavilion.id)">
                    {{ pavilion.name }}
                  </option>
                  <option v-if="filteredPavilions.length === 0 && !loadingPavilions" disabled>No se encontraron
                    coincidencias</option>
                </select>
              </div>

              <!-- Botón de Enviar -->
              <button @click="sendToPavilion" :disabled="!selectedPavilion" :class="[
                'w-full py-2 px-4 rounded-md text-sm font-medium transition-colors',
                selectedPavilion
                  ? 'bg-blue-600 hover:bg-blue-700 text-white'
                  : 'bg-gray-300 text-gray-500 cursor-not-allowed'
              ]">
                <svg class="h-4 w-4 mr-2 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
                </svg>
                Enviar a Pabellón
              </button>
            </div>
          </div>

          <!-- Código QR y acciones -->
          <div class="bg-white rounded border p-4">
            <h4 class="font-semibold text-gray-900 flex items-center mb-3">
              <svg class="h-5 w-5 mr-2 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h2" />
              </svg>
              Código QR
            </h4>
            <div class="flex items-center justify-between">
              <code class="text-sm text-gray-800 font-mono bg-gray-50 px-2 py-1 rounded flex-1 mr-2">
                {{ qrInfo.qr_code }}
              </code>
              <button @click="copyToClipboard(qrInfo.qr_code)"
                class="text-xs bg-gray-500 hover:bg-gray-600 text-white px-2 py-1 rounded transition-colors"
                title="Copiar código">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Historial (si existe) -->
    <div v-if="qrInfo.history && qrInfo.history.length > 0" class="bg-white rounded-lg border p-6">
      <h4 class="font-semibold text-gray-900 flex items-center mb-4">
        <svg class="h-5 w-5 mr-2 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        Historial de Movimientos
      </h4>

      <div class="space-y-3">
        <div v-for="(item, index) in qrInfo.history" :key="index"
          class="flex items-center justify-between p-3 bg-gray-50 rounded">
          <div class="flex items-center">
            <div :class="getHistoryIconClass(item.status)" class="p-2 rounded-full mr-3">
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path v-if="item.status === 'consumido'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
            </div>
            <div>
              <p class="font-medium text-gray-900">{{ item.status.charAt(0).toUpperCase() + item.status.slice(1) }}</p>
              <p class="text-sm text-gray-600">{{ formatDate(item.date_time) }}</p>
            </div>
          </div>
          <div class="text-right">
            <p class="text-sm text-gray-600">Usuario: {{ item.user_rut }}</p>
            <p class="text-sm text-gray-600">{{ item.destination_type }}: {{ item.destination_id }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Acciones Disponibles -->
    <div class="bg-white rounded-lg shadow-sm border p-4 sm:p-6">
      <h3 class="text-base sm:text-lg font-medium text-gray-900 mb-4 text-center">
        Acciones Disponibles
      </h3>

      <!-- Contenedor de botones completamente responsivo -->
      <div class="qr-main-actions">
        <button @click="$emit('view-details', qrInfo)" class="qr-action-button qr-action-primary">
          <div class="flex items-center justify-center">
            <svg class="h-5 w-5 mr-3 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
            </svg>
            <span class="font-medium">Ver Detalles Completos</span>
          </div>
        </button>

        <button v-if="qrInfo.type === 'medical_supply'" @click="$emit('view-batch', qrInfo.supply_info?.batch_id)"
          class="qr-action-button qr-action-secondary">
          <div class="flex items-center justify-center">
            <svg class="h-5 w-5 mr-3 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
            </svg>
            <span class="font-medium">Ver Lote Relacionado</span>
          </div>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import qrService from '@/services/qrService'
import pavilionService from '@/services/pavilionService'

// Props
const props = defineProps({
  qrInfo: {
    type: Object,
    required: true
  }
})

// Emits
const emit = defineEmits([
  'view-details',
  'view-batch',
  'consume-supply',
  'send-to-pavilion'
])

// Estado reactivo
const imageLoaded = ref(false)
const imageError = ref(false)
const pavilionSearch = ref('')
const selectedPavilion = ref('') // Siempre string para coincidir con el value del select
const pavilions = ref([])
const loadingPavilions = ref(false)
const pavilionsError = ref("")

// Computed
const filteredPavilions = computed(() => {
  const search = pavilionSearch.value.trim().toLowerCase()
  if (!search) return pavilions.value.map(p => ({ ...p }))
  const result = pavilions.value
    .filter(pavilion => pavilion.name && pavilion.name.toLowerCase().includes(search))
    .map(p => ({ ...p }))
  console.log('Filtrando pabellones:', search, result)
  return result
})

// URLs de imagen
const qrImageUrl = computed(() => {

  // Watcher para limpiar selección si no hay coincidencias
  watch(filteredPavilions, (newVal) => {
    if (newVal.length === 0) {
      selectedPavilion.value = ''
    }
  })
  if (!props.qrInfo?.qr_code) return null
  return qrService.getQRImageUrl(props.qrInfo.qr_code)
})

// Métodos
const handleImageError = () => {
  imageError.value = true
  imageLoaded.value = false
}

const downloadQRImage = async (resolution = 'normal') => {
  try {
    await qrService.downloadQRImage(props.qrInfo.qr_code, resolution)
  } catch (error) {
    console.error('Error al descargar imagen QR:', error)
  }
}

const copyToClipboard = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    // Podrías mostrar una notificación aquí
  } catch (error) {
    console.error('Error al copiar al portapapeles:', error)
  }
}

const filterPavilions = () => {
  // La función computed ya maneja esto automáticamente
}

const sendToPavilion = async () => {
  if (!selectedPavilion.value) return

  try {
    const pavilionData = pavilions.value.find(p => String(p.id) === selectedPavilion.value)

    // Emitir evento al componente padre con los datos del envío
    emit('send-to-pavilion', {
      qrCode: props.qrInfo.qr_code,
      pavilionId: selectedPavilion.value,
      pavilionName: pavilionData?.name,
      supplyInfo: props.qrInfo.supply_code
    })

    // Limpiar selección
    selectedPavilion.value = ''
    pavilionSearch.value = ''

  } catch (error) {
    console.error('Error al enviar a pabellón:', error)
  }
}

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
    'consumido': 'bg-red-100 text-red-600',
    'creado': 'bg-green-100 text-green-600',
    'movido': 'bg-blue-100 text-blue-600'
  }
  return classes[status] || 'bg-gray-100 text-gray-600'
}

// Lifecycle
onMounted(async () => {
  // Verificar si la imagen existe al montar el componente
  if (props.qrInfo?.qr_code) {
    qrService.checkQRImageExists(props.qrInfo.qr_code)
      .then(exists => {
        if (!exists) {
          imageError.value = true
        }
      })
      .catch(() => {
        imageError.value = true
      })
  }
  // Obtener pabellones desde el backend
  loadingPavilions.value = true
  pavilionsError.value = ""
  try {
    const result = await pavilionService.getAllPavilions()
    console.log("Respuesta pabellones (raw):", result)
    // Si la respuesta tiene la estructura { data: { success: true, data: [...] } }
    if (result && result.data && Array.isArray(result.data.data)) {
      pavilions.value = result.data.data
      if (result.data.data.length === 0) {
        pavilionsError.value = "No se encontraron pabellones."
      }
    } else if (Array.isArray(result)) {
      pavilions.value = result
      if (result.length === 0) {
        pavilionsError.value = "No se encontraron pabellones."
      }
    } else {
      pavilions.value = []
      pavilionsError.value = "La respuesta no es un array."
      console.error('Estructura inesperada de pabellones:', result)
    }
  } catch (error) {
    pavilions.value = []
    pavilionsError.value = error?.message || "Error desconocido"
    console.error('Error cargando pabellones:', error)
  } finally {
    loadingPavilions.value = false
  }
})
</script>

<style scoped>
/* Estilos específicos para los botones de QR */
.qr-main-actions {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  align-items: stretch;
}

.qr-action-button {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem 1.5rem;
  border-radius: 0.75rem;
  font-size: 0.9rem;
  font-weight: 500;
  transition: all 0.2s ease-in-out;
  border: 2px solid transparent;
  min-height: 52px;
  text-align: center;
  cursor: pointer;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06);
}

.qr-action-primary {
  background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
  color: white;
}

.qr-action-primary:hover {
  background: linear-gradient(135deg, #1d4ed8 0%, #1e40af 100%);
  transform: translateY(-1px);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
}

.qr-action-secondary {
  background: linear-gradient(135deg, #6b7280 0%, #4b5563 100%);
  color: white;
}

.qr-action-secondary:hover {
  background: linear-gradient(135deg, #4b5563 0%, #374151 100%);
  transform: translateY(-1px);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
}

/* Responsive adjustments */
@media (min-width: 640px) {
  .qr-main-actions {
    flex-direction: row;
    justify-content: center;
  }

  .qr-action-button {
    flex: 1;
    max-width: 300px;
  }
}

/* Badge styles */
.badge {
  display: inline-block;
  padding: 0.25rem 0.5rem;
  border-radius: 0.375rem;
  font-size: 0.75rem;
  font-weight: 500;
}
</style>