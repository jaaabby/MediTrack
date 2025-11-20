<template>
  <div class="max-w-4xl mx-auto p-6">
    <!-- Header -->
    <div class="mb-6">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Recepción de Insumos</h1>
          <p class="text-gray-600 mt-2">Recepcionar insumos que están en camino al pabellón</p>
        </div>
        <div class="flex space-x-2">
          <router-link to="/qr" class="btn-secondary">
            <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
            Volver al Scanner
          </router-link>
        </div>
      </div>
    </div>

    <!-- Mensaje de Error -->
    <div v-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
      <div class="flex items-start space-x-3">
        <svg class="h-5 w-5 text-red-600 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <div class="flex-1">
          <h4 class="text-sm font-medium text-red-800">Error al escanear</h4>
          <p class="text-sm text-red-700 mt-1">{{ error }}</p>
          <button @click="clearError" class="text-sm text-red-600 hover:text-red-800 mt-2 underline">
            Limpiar error
          </button>
        </div>
      </div>
    </div>

    <!-- Información del Insumo Escaneado -->
    <div v-if="scannedProduct && !error" class="bg-white rounded-lg shadow-sm border overflow-hidden mb-6">
      <div class="px-4 py-5 sm:p-6">
        <h3 class="text-lg font-medium text-gray-900 mb-4">Información del Insumo</h3>

        <!-- Product Details -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
          <div>
            <div class="space-y-3">
              <div>
                <label class="text-sm font-medium text-gray-500">Código QR</label>
                <p class="text-sm text-gray-900 font-mono">{{ scannedProduct.qr_code }}</p>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-500">Nombre</label>
                <p class="text-sm text-gray-900">{{ 
                  scannedProduct.supply_code?.name || 
                  scannedProduct.supply_info?.SupplyCode?.name || 
                  scannedProduct.supply_info?.name || 
                  'No disponible' 
                }}</p>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-500">Estado</label>
                <span :class="getStatusBadgeClass(scannedProduct)"
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium">
                  {{ getStatusLabel(scannedProduct) }}
                </span>
              </div>
            </div>
          </div>
          <div>
            <div class="space-y-3">
              <div v-if="scannedProduct.supply_info?.batch">
                <label class="text-sm font-medium text-gray-500">Lote</label>
                <p class="text-sm text-gray-900">{{ scannedProduct.supply_info.batch.id }}</p>
              </div>
              <div v-if="scannedProduct.supply_info?.batch?.expiration_date">
                <label class="text-sm font-medium text-gray-500">Fecha de Vencimiento</label>
                <p class="text-sm text-gray-900">{{ formatDate(scannedProduct.supply_info.batch.expiration_date) }}</p>
              </div>
              <div v-if="scannedProduct.supply_info?.batch?.supplier">
                <label class="text-sm font-medium text-gray-500">Proveedor</label>
                <p class="text-sm text-gray-900">{{ scannedProduct.supply_info.batch.supplier }}</p>
              </div>
            </div>
          </div>
        </div>
      
        <!-- Estado del insumo -->
        <div class="border-t border-gray-200 pt-4">
          <div class="flex items-center justify-between">
            <div>
              <h4 class="text-sm font-medium text-gray-900">Estado Actual</h4>
              <p class="text-sm text-gray-600">{{ getStatusText(scannedProduct.supply_info?.Status || scannedProduct.status || scannedProduct.current_status) }}</p>
            </div>
            <div v-if="canReceive(scannedProduct)" class="flex items-center space-x-2">
              <div class="w-3 h-3 bg-yellow-400 rounded-full"></div>
              <span class="text-sm font-medium text-yellow-800">Listo para Recepcionar</span>
            </div>
            <div v-else class="flex items-center space-x-2">
              <div class="w-3 h-3 bg-red-400 rounded-full"></div>
              <span class="text-sm font-medium text-red-800">No se puede Recepcionar</span>
            </div>
          </div>
        </div>

        <!-- Status Validation -->
        <div v-if="!canReceive(scannedProduct)" class="bg-red-50 border border-red-200 rounded-md p-4 mt-4">
          <div class="flex">
            <div class="flex-shrink-0">
              <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div class="ml-3">
              <h3 class="text-sm font-medium text-red-800">No se puede recepcionar</h3>
              <div class="mt-2 text-sm text-red-700">
                {{ getReceptionErrorMessage(scannedProduct) }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Formulario de Recepción Simplificado para Pabellón -->
    <div v-if="canReceive(scannedProduct)" class="bg-white rounded-lg shadow-sm border p-6">
      <h3 class="text-lg font-medium text-gray-900 mb-4 flex items-center">
        <svg class="h-5 w-5 mr-2 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        Formulario de Recepción
      </h3>

      <div class="space-y-6">
        <!-- Información del Usuario de Pabellón -->
        <div class="bg-gray-50 rounded-lg p-4">
          <h4 class="text-sm font-medium text-gray-900 mb-2">Usuario de Pabellón</h4>
          <p class="text-sm text-gray-600">RUT: {{ currentUser?.rut || 'No disponible' }}</p>
          <p class="text-sm text-gray-600">Nombre: {{ currentUser?.name || 'No disponible' }}</p>
        </div>

        <!-- Notas de Recepción -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Notas de Recepción
          </label>
          <textarea
            v-model="receptionForm.notes"
            rows="3"
            placeholder="Observaciones sobre la recepción del insumo..."
            class="form-textarea w-full"
          ></textarea>
        </div>

        <!-- Botón de Recepción -->
        <div class="flex justify-end">
          <button
            @click="receiveSupply"
            :disabled="receiving"
            class="btn-primary"
          >
            <svg v-if="receiving" class="animate-spin h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ receiving ? 'Recepcionando...' : 'Recepcionar Insumo' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Success Message -->
    <div v-if="consumptionSuccess" class="bg-green-50 border border-green-200 rounded-md p-4">
      <div class="flex">
        <div class="flex-shrink-0">
          <svg class="h-5 w-5 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
        </div>
        <div class="ml-3">
          <h3 class="text-sm font-medium text-green-800">Recepción Exitosa</h3>
          <div class="mt-2 text-sm text-green-700">
            El insumo ha sido recepcionado correctamente y está listo para uso.
            <div v-if="consumptionSuccess.status_change" class="mt-1 text-xs text-green-600">
              Estado cambiado de "{{ consumptionSuccess.status_change.from }}" a "{{ consumptionSuccess.status_change.to }}"
            </div>
          </div>
          <div class="mt-4 flex flex-wrap gap-3">
            <button @click="resetForm" class="btn-primary text-sm">
              Recepcionar Otro Insumo
            </button>
            <router-link to="/qr" class="btn-secondary text-sm">
              <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
              Volver al Scanner
            </router-link>
            <router-link v-if="consumptionSuccess.qr_code" :to="`/qr/${consumptionSuccess.qr_code}/traceability`" class="btn-secondary text-sm">
              Ver Trazabilidad
            </router-link>
            <router-link v-if="consumptionSuccess.batch_id" :to="`/batch/${consumptionSuccess.batch_id}/history`"
              class="btn-secondary text-sm">
              Ver Historial del Lote
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import qrService from '@/services/qr/qrService'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

// Estado del componente
const loading = ref(false)
const qrInput = ref('')
const error = ref(null)
const scannedProduct = ref(null)
const receiving = ref(false)
const consumptionSuccess = ref(null) // Mantenemos el nombre para conservar la alerta

// Formulario de recepción simplificado
const receptionForm = ref({
  notes: ''
})

// Usuario actual
const currentUser = computed(() => authStore.user)

onMounted(() => {
  // Si hay un QR en la URL, escanearlo automáticamente
  const qrFromUrl = route.query.qr
  if (qrFromUrl) {
    qrInput.value = qrFromUrl
    scanQR()
  }
})

// Escanear QR
const scanQR = async () => {
  if (!qrInput.value.trim()) return
  
  loading.value = true
  error.value = null
  scannedProduct.value = null
  consumptionSuccess.value = null
  
  try {
    // Validar formato
    if (!isValidQRFormat(qrInput.value.trim())) {
      throw new Error('Formato de código QR inválido. Use SUPPLY_... para insumos individuales.')
    }
    
    // Crear contexto de escaneo
    const scanContext = {
      scan_purpose: 'reception',
      scan_source: 'reception_view',
      user_agent: navigator.userAgent
    }
    
    // Escanear código QR
    const result = await qrService.scanQRCode(qrInput.value.trim(), scanContext)
    
    console.log('Resultado del escaneo QR:', result) // Debug

    if (result) {
      scannedProduct.value = result

      // Verificar si es un insumo individual
      if (result.type !== 'medical_supply') {
        error.value = 'Solo se pueden recepcionar insumos individuales. Este código QR corresponde a un lote.'
        scannedProduct.value = null
        return
      }
    }
    
  } catch (err) {
    console.error('Error scanning QR:', err)
    error.value = err.response?.data?.error || err.message || 'Error al escanear el código QR'
  } finally {
    loading.value = false
  }
}

const isValidQRFormat = (qrCode) => {
  if (!qrCode || typeof qrCode !== 'string') return false
  const qrPattern = /^SUPPLY_\d+_[a-f0-9]+$/i
  return qrPattern.test(qrCode)
}

// Verificar si necesita retiro primero
const needsPickup = (product) => {
  if (!product || product.type !== 'medical_supply') return false
  if (product.is_consumed) return false
  
  const status = product.supply_info?.Status || 
                 product.supply_info?.status || 
                 product.status || 
                 product.current_status
  
  return status === 'pendiente_retiro'
}

// Verificar si se puede recepcionar
const canReceive = (product) => {
  if (!product || product.type !== 'medical_supply') return false
  if (product.is_consumed) return false
  
  const status = product.supply_info?.Status || 
                 product.supply_info?.status || 
                 product.status || 
                 product.current_status
  
  console.log('canReceive: Estado encontrado:', status)
  return status === 'en_camino_a_pabellon'
}

// Obtener mensaje de error para recepción
const getReceptionErrorMessage = (product) => {
  if (!product) return ''

  if (product.is_consumed) {
    return 'Este insumo ya ha sido consumido anteriormente.'
  }

  const status = product.supply_info?.Status || 
                 product.supply_info?.status || 
                 product.status || 
                 product.current_status

  if (status === 'recepcionado') {
    return 'Este insumo ya ha sido recepcionado anteriormente.'
  }

  if (status === 'disponible') {
    return 'Este insumo está disponible en el centro médico.'
  }

  if (status !== 'en_camino_a_pabellon') {
    return `El insumo tiene estado "${status}" y no está en camino al pabellón. Solo se pueden recepcionar insumos con estado "en_camino_a_pabellon".`
  }

  return 'Este insumo no se puede recepcionar.'
}

// Obtener texto del estado
const getStatusText = (status) => {
  const statusMap = {
    'disponible': 'Disponible',
    'en_camino_a_pabellon': 'En Camino a Pabellón',
    'recepcionado': 'Recepcionado',
    'consumido': 'Consumido',
    'transferido': 'Transferido'
  }
  return statusMap[status] || status || 'Desconocido'
}

// Recepcionar insumo
const receiveSupply = async () => {
  if (!scannedProduct.value || !qrInput.value.trim()) return
  
  receiving.value = true
  error.value = null
  consumptionSuccess.value = null
  
  try {
    // Obtener datos del usuario actual
    const userRUT = currentUser.value?.rut
    if (!userRUT) {
      throw new Error('No se pudo obtener el RUT del usuario de pabellón')
    }

    // Obtener el pabellón del usuario (asumiendo que está en el contexto del pabellón)
    const pavilionId = currentUser.value?.pavilion_id || 1 // Fallback a pabellón 1

    // Guardar el estado anterior para mostrar el cambio
    const previousStatus = scannedProduct.value.supply_info?.Status || 
                          scannedProduct.value.supply_info?.status || 
                          scannedProduct.value.status || 
                          scannedProduct.value.current_status || 
                          'en_camino_a_pabellon'
    
    const result = await qrService.receiveSupply(
      qrInput.value.trim(),
      userRUT,
      'pavilion',
      pavilionId,
      receptionForm.value.notes
    )
    
    console.log('Resultado de receiveSupply:', result) // Debug
    
    if (result.success) {
      console.log('Recepción exitosa, configurando alerta...') // Debug
      
      consumptionSuccess.value = {
        ...result.data,
        qr_code: scannedProduct.value.qr_code,
        batch_id: scannedProduct.value.supply_info?.batch?.id,
        status_change: {
          from: getStatusLabel({ supply_info: { status: previousStatus } }),
          to: 'Recepcionado'
        }
      }
      
      console.log('consumptionSuccess configurado:', consumptionSuccess.value) // Debug
      
      // Limpiar solo los campos del formulario pero mantener la alerta de éxito
      receptionForm.value = {
        notes: ''
      }
      qrInput.value = ''
      scannedProduct.value = null
      error.value = null
      // NO limpiar consumptionSuccess para que se muestre la alerta
      
    } else {
      throw new Error(result.error || 'Error al recepcionar el insumo')
    }
    
  } catch (err) {
    console.error('Error receiving supply:', err)
    error.value = err.response?.data?.error || err.message || 'Error al recepcionar el insumo'
  } finally {
    receiving.value = false
  }
}

// Funciones de utilidad
const formatDate = (dateString) => {
  if (!dateString) return 'No disponible'
  try {
    return format(new Date(dateString), 'dd/MM/yyyy', { locale: es })
  } catch (error) {
    return dateString
  }
}

const getStatusLabel = (product) => {
  if (!product) return 'Desconocido'

  if (product.is_consumed) return 'Consumido'

  const status = product.supply_info?.status || product.status
  switch (status) {
    case 'disponible': return 'Disponible'
    case 'en_camino_a_pabellon': return 'En Camino a Pabellón'
    case 'recepcionado': return 'Recepcionado'
    default: return status || 'Desconocido'
  }
}

const getStatusBadgeClass = (product) => {
  if (!product) return 'bg-gray-100 text-gray-800'

  if (product.is_consumed) return 'bg-red-100 text-red-800'

  const status = product.supply_info?.status || product.status
  switch (status) {
    case 'disponible': return 'bg-green-100 text-green-800'
    case 'en_camino_a_pabellon': return 'bg-yellow-100 text-yellow-800'
    case 'recepcionado': return 'bg-blue-100 text-blue-800'
    default: return 'bg-gray-100 text-gray-800'
  }
}

// Limpiar formulario
const resetForm = () => {
  receptionForm.value = {
    notes: ''
  }
  qrInput.value = ''
  scannedProduct.value = null
  error.value = null
  consumptionSuccess.value = null
}

// Limpiar error
const clearError = () => {
  error.value = null
}

// Navegación
const viewDetails = (qrInfo) => {
  router.push({
    name: 'QRDetails',
    params: { qrCode: qrInfo.qr_code }
  })
}

const viewBatch = (batchId) => {
  router.push({
    name: 'Inventory',
    query: { batch: batchId }
  })
}
</script>

<style scoped>
.form-input {
  display: block;
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  font-size: 0.875rem;
}

.form-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.form-input::placeholder {
  color: #9ca3af;
}

.form-select {
  display: block;
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  font-size: 0.875rem;
}

.form-select:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.form-textarea {
  display: block;
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  font-size: 0.875rem;
  resize: vertical;
}

.form-textarea:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.form-textarea::placeholder {
  color: #9ca3af;
}

/* Usar clases de botones de style.css global */
</style>