<template>
  <div class="space-y-6">
    <!-- Información básica del QR con imagen -->
    <div class="bg-gray-50 rounded-lg p-6">
      <div class="flex flex-col lg:flex-row lg:items-start lg:space-x-6">
        
        <!-- Imagen QR -->
        <div class="flex-shrink-0 mb-4 lg:mb-0">
          <div class="bg-white rounded-lg p-4 shadow-sm border text-center">
            <img 
              v-if="qrImageUrl" 
              :src="qrImageUrl" 
              :alt="`Código QR: ${qrInfo.qr_code}`"
              class="w-32 h-32 mx-auto object-contain"
              @error="handleImageError"
              @load="imageLoaded = true"
            />
            <div 
              v-else-if="imageError" 
              class="w-32 h-32 mx-auto bg-gray-100 rounded flex items-center justify-center"
            >
              <svg class="h-8 w-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.732-.833-2.5 0L4.268 16.5c-.77.833.192 2.5 1.732 2.5z" />
              </svg>
            </div>
            <div 
              v-else 
              class="w-32 h-32 mx-auto bg-gray-100 rounded flex items-center justify-center animate-pulse"
            >
              <svg class="h-8 w-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v14a2 2 0 002 2z" />
              </svg>
            </div>
            
            <!-- Botones de acción para la imagen -->
            <div v-if="imageLoaded" class="mt-3 space-y-2">
              <button 
                @click="downloadQRImage('normal')"
                class="w-full text-xs bg-gray-600 hover:bg-gray-700 text-white px-2 py-1 rounded transition-colors"
              >
                Descargar
              </button>
              <button 
                @click="downloadQRImage('high')"
                class="w-full text-xs bg-blue-600 hover:bg-blue-700 text-white px-2 py-1 rounded transition-colors"
              >
                HD
              </button>
            </div>
          </div>
        </div>

        <!-- Información básica -->
        <div class="flex-1">
          <div class="flex items-start justify-between mb-4">
            <div class="flex items-center">
              <div :class="getTypeIconClass(qrInfo.type)" class="p-3 rounded-lg mr-4">
                <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path v-if="qrInfo.type === 'batch'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                  <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
                </svg>
              </div>
              <div>
                <h3 class="text-lg font-semibold text-gray-900">{{ getTypeLabel(qrInfo.type) }}</h3>
                <p class="text-sm text-gray-600">ID: {{ qrInfo.id }}</p>
                <div class="flex items-center mt-1">
                  <span :class="getStatusBadgeClass()" class="badge text-xs">
                    {{ getStatusLabel() }}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <!-- Código QR -->
          <div class="bg-white rounded border p-3 mb-4">
            <p class="text-sm font-medium text-gray-700 mb-2">Código QR:</p>
            <div class="flex items-center justify-between">
              <code class="text-sm text-gray-800 font-mono bg-gray-50 px-2 py-1 rounded flex-1 mr-2">
                {{ qrInfo.qr_code }}
              </code>
              <button 
                @click="copyToClipboard(qrInfo.qr_code)"
                class="text-xs bg-gray-500 hover:bg-gray-600 text-white px-2 py-1 rounded transition-colors"
                title="Copiar código"
              >
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Información específica según el tipo -->
    <div v-if="qrInfo.type === 'batch' && qrInfo.batch_info" class="bg-white rounded-lg border p-6">
      <h4 class="font-semibold text-gray-900 flex items-center mb-4">
        <svg class="h-5 w-5 mr-2 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
        </svg>
        Información del Lote
      </h4>

      <div class="grid md:grid-cols-2 gap-6">
        <div class="space-y-4">
          <div>
            <label class="text-sm font-medium text-gray-600">Proveedor:</label>
            <p class="text-gray-900">{{ qrInfo.batch_info.supplier }}</p>
          </div>
          <div>
            <label class="text-sm font-medium text-gray-600">Fecha de Vencimiento:</label>
            <p class="text-gray-900">{{ formatDate(qrInfo.batch_info.expiration_date) }}</p>
          </div>
          <div>
            <label class="text-sm font-medium text-gray-600">Cantidad Actual:</label>
            <p class="text-gray-900 font-semibold">{{ qrInfo.batch_info.amount }} unidades</p>
          </div>
        </div>

        <!-- Estadísticas del lote -->
        <div v-if="qrInfo.batch_status" class="space-y-4">
          <h5 class="font-medium text-gray-700">Estado del Lote:</h5>
          <div class="bg-gray-50 rounded p-4 space-y-2">
            <div class="flex justify-between">
              <span class="text-sm text-gray-600">Total de productos:</span>
              <span class="font-medium">{{ qrInfo.batch_status.total_individual_supplies }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-sm text-gray-600">Consumidos:</span>
              <span class="font-medium text-red-600">{{ qrInfo.batch_status.consumed_supplies }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-sm text-gray-600">Disponibles:</span>
              <span class="font-medium text-green-600">{{ qrInfo.batch_status.available_supplies }}</span>
            </div>
            <div class="flex justify-between pt-2 border-t">
              <span class="text-sm text-gray-600">Sincronización:</span>
              <span :class="qrInfo.batch_status.amounts_synchronized ? 'text-green-600' : 'text-yellow-600'">
                {{ qrInfo.batch_status.amounts_synchronized ? 'Correcta' : 'Pendiente' }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Información del insumo individual -->
    <div v-if="qrInfo.type === 'medical_supply' && qrInfo.supply_info" class="bg-white rounded-lg border p-6">
      <h4 class="font-semibold text-gray-900 flex items-center mb-4">
        <svg class="h-5 w-5 mr-2 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
        </svg>
        Información del Insumo
      </h4>

      <div class="grid md:grid-cols-2 gap-6">
        <div class="space-y-4">
          <div>
            <label class="text-sm font-medium text-gray-600">Nombre:</label>
            <p class="text-gray-900 font-medium">{{ qrInfo.supply_info.supply_code_name }}</p>
          </div>
          <div>
            <label class="text-sm font-medium text-gray-600">Código:</label>
            <p class="text-gray-900">{{ qrInfo.supply_info.code }}</p>
          </div>
          <div>
            <label class="text-sm font-medium text-gray-600">Proveedor:</label>
            <p class="text-gray-900">{{ qrInfo.supply_info.supplier }}</p>
          </div>
          <div>
            <label class="text-sm font-medium text-gray-600">Estado:</label>
            <span :class="qrInfo.supply_info.is_consumed ? 'bg-red-100 text-red-800' : 'bg-green-100 text-green-800'" class="inline-block px-2 py-1 rounded text-sm font-medium">
              {{ qrInfo.supply_info.is_consumed ? 'Consumido' : 'Disponible' }}
            </span>
          </div>
        </div>

        <div class="space-y-4">
          <div>
            <label class="text-sm font-medium text-gray-600">Fecha de Vencimiento:</label>
            <p class="text-gray-900">{{ formatDate(qrInfo.supply_info.expiration_date) }}</p>
          </div>
          <div>
            <label class="text-sm font-medium text-gray-600">Almacén:</label>
            <p class="text-gray-900">{{ qrInfo.supply_info.store_name }}</p>
          </div>
          <div>
            <label class="text-sm font-medium text-gray-600">ID del Lote:</label>
            <p class="text-gray-900">{{ qrInfo.supply_info.batch_id }}</p>
          </div>
        </div>
      </div>

      <!-- Acciones para insumos individuales -->
      <div v-if="!qrInfo.supply_info.is_consumed && qrInfo.can_consume" class="mt-6 pt-4 border-t">
        <h5 class="font-medium text-gray-700 mb-3">Acciones:</h5>
        <div class="flex space-x-3">
          <button 
            @click="$emit('consume-supply', qrInfo)"
            class="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded font-medium transition-colors"
          >
            <svg class="h-4 w-4 inline mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
            Consumir Producto
          </button>
        </div>
      </div>
    </div>

    <!-- Información del supply code -->
    <div v-if="qrInfo.supply_code" class="bg-white rounded-lg border p-6">
      <h4 class="font-semibold text-gray-900 flex items-center mb-4">
        <svg class="h-5 w-5 mr-2 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
        </svg>
        Información del Código de Insumo
      </h4>

      <div class="grid md:grid-cols-3 gap-4">
        <div>
          <label class="text-sm font-medium text-gray-600">Código:</label>
          <p class="text-gray-900 font-semibold">{{ qrInfo.supply_code.code }}</p>
        </div>
        <div>
          <label class="text-sm font-medium text-gray-600">Nombre:</label>
          <p class="text-gray-900">{{ qrInfo.supply_code.name }}</p>
        </div>
        <div>
          <label class="text-sm font-medium text-gray-600">Código de Proveedor:</label>
          <p class="text-gray-900">{{ qrInfo.supply_code.code_supplier }}</p>
        </div>
      </div>
    </div>

    <!-- Historial (si existe) -->
    <div v-if="qrInfo.history && qrInfo.history.length > 0" class="bg-white rounded-lg border p-6">
      <h4 class="font-semibold text-gray-900 flex items-center mb-4">
        <svg class="h-5 w-5 mr-2 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        Historial de Movimientos
      </h4>

      <div class="space-y-3">
        <div 
          v-for="(item, index) in qrInfo.history" 
          :key="index"
          class="flex items-center justify-between p-3 bg-gray-50 rounded"
        >
          <div class="flex items-center">
            <div :class="getHistoryIconClass(item.status)" class="p-2 rounded-full mr-3">
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path v-if="item.status === 'consumido'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
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

    <!-- Acciones generales -->
    <div class="flex flex-wrap gap-3 justify-center">
      <button 
        @click="$emit('view-details', qrInfo)"
        class="btn-primary"
      >
        <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
        </svg>
        Ver Detalles Completos
      </button>

      <button 
        v-if="qrInfo.type === 'medical_supply'"
        @click="$emit('view-batch', qrInfo.supply_info?.batch_id)"
        class="btn-secondary"
      >
        <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
        </svg>
        Ver Lote Relacionado
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import qrService from '@/services/qrService'

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
  'consume-supply'
])

// Estado reactivo
const imageLoaded = ref(false)
const imageError = ref(false)

// URLs de imagen
const qrImageUrl = computed(() => {
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
onMounted(() => {
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
})
</script>

