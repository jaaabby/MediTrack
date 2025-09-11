<template>
  <div class="max-w-4xl mx-auto p-6">
    <!-- Header -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-900">Recepción de Insumos</h1>
      <p class="text-gray-600 mt-2">Recepcionar insumos que están en camino al pabellón</p>
    </div>

    <!-- Sección de Escaneo -->
    <div class="bg-white rounded-lg shadow-sm border p-6 mb-6">
      <h3 class="text-lg font-medium text-gray-900 mb-4 flex items-center">
        <svg class="h-5 w-5 mr-2 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h2M4 4h5l2 3h3l2-3h5v5M4 4v5m16-5v5" />
        </svg>
        Escanear Código QR del Insumo
      </h3>

      <div class="space-y-4">
        <label for="qr-input" class="block text-sm font-medium text-gray-700">
          Código QR del Insumo:
        </label>
        <div class="flex space-x-3">
          <input
            id="qr-input"
            v-model="qrInput"
            type="text"
            placeholder="SUPPLY_1755580808_def456"
            class="form-input flex-1"
            @keyup.enter="scanQR"
            :disabled="loading"
          />
          <button
            @click="scanQR"
            :disabled="!qrInput.trim() || loading"
            class="btn-primary"
          >
            <svg v-if="loading" class="animate-spin h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ loading ? 'Escaneando...' : 'Escanear' }}
          </button>
        </div>
        
        <!-- Ayuda de formato -->
        <div class="p-3 bg-blue-50 rounded-lg">
          <p class="text-sm text-blue-800">
            <svg class="h-4 w-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <strong>Formato requerido:</strong> SUPPLY_... (insumos individuales)
          </p>
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
      <QRInfoDisplay 
        :qr-info="scannedProduct"
        :show-traceability="true"
        :scan-context="lastScanContext"
        @view-details="viewDetails"
        @view-batch="viewBatch"
      />
      
      <!-- Estado del insumo -->
      <div class="p-4 border-t border-gray-200">
        <div class="flex items-center justify-between">
          <div>
            <h4 class="text-sm font-medium text-gray-900">Estado Actual</h4>
            <p class="text-sm text-gray-600">{{ getStatusText(scannedProduct.supply_info?.Status) }}</p>
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

    <!-- Mensaje de Éxito -->
    <div v-if="receptionSuccess" class="bg-green-50 border border-green-200 rounded-lg p-4 mb-6">
      <div class="flex items-start space-x-3">
        <svg class="h-5 w-5 text-green-600 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
        </svg>
        <div class="flex-1">
          <h4 class="text-sm font-medium text-green-800">Recepción Exitosa</h4>
          <p class="text-sm text-green-700 mt-1">{{ receptionSuccess }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import qrService from '@/services/qrService'
import QRInfoDisplay from '@/components/QRInfoDisplay.vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

// Estado del componente
const loading = ref(false)
const error = ref(null)
const qrInput = ref('')
const scannedProduct = ref(null)
const lastScanContext = ref(null)
const receiving = ref(false)
const receptionSuccess = ref(null)

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
  receptionSuccess.value = null
  
  try {
    // Validar formato
    if (!isValidQRFormat(qrInput.value.trim())) {
      throw new Error('Formato de código QR inválido. Use SUPPLY_... para insumos individuales.')
    }
    
    // Crear contexto de escaneo
    const scanContext = buildScanContext()
    
    // Escanear código QR
    const result = await qrService.scanQRCode(qrInput.value.trim(), scanContext)
    
    scannedProduct.value = result
    lastScanContext.value = scanContext
    
  } catch (err) {
    console.error('Error scanning QR:', err)
    error.value = err.response?.data?.error || err.message || 'Error al escanear el código QR'
  } finally {
    loading.value = false
  }
}

const buildScanContext = () => {
  return {
    scan_purpose: 'reception',
    scan_source: 'web',
    user_agent: navigator.userAgent,
    device_info: {
      platform: navigator.platform,
      language: navigator.language,
      screen_resolution: `${screen.width}x${screen.height}`
    }
  }
}

const isValidQRFormat = (qrCode) => {
  if (!qrCode || typeof qrCode !== 'string') return false
  const qrPattern = /^SUPPLY_\d+_[a-f0-9]+$/i
  return qrPattern.test(qrCode)
}

// Verificar si se puede recepcionar
const canReceive = (product) => {
  if (!product || product.type !== 'medical_supply') return false
  if (product.is_consumed) return false
  
  const status = product.supply_info?.Status || product.supply_info?.status || product.status || product.current_status
  return status === 'en_camino_a_pabellon'
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
  receptionSuccess.value = null
  
  try {
    // Obtener datos del usuario actual
    const userRUT = currentUser.value?.rut
    if (!userRUT) {
      throw new Error('No se pudo obtener el RUT del usuario de pabellón')
    }

    // Obtener el pabellón del usuario (asumiendo que está en el contexto del pabellón)
    const pavilionId = currentUser.value?.pavilion_id || 1 // Fallback a pabellón 1
    
    const result = await qrService.receiveSupply(
      qrInput.value.trim(),
      userRUT,
      'pavilion',
      pavilionId,
      receptionForm.value.notes
    )
    
    if (result.success) {
      receptionSuccess.value = `Insumo recepcionado exitosamente. Estado cambiado a "recepcionado".`
      
      // Limpiar formulario
      resetForm()
      
      // Re-escanear para actualizar la información
      setTimeout(() => {
        scanQR()
      }, 1000)
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

// Limpiar formulario
const resetForm = () => {
  receptionForm.value = {
    notes: ''
  }
  qrInput.value = ''
  scannedProduct.value = null
  error.value = null
  receptionSuccess.value = null
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
  @apply block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm;
}

.form-select {
  @apply block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm;
}

.form-textarea {
  @apply block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm;
}

.btn-primary {
  @apply inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed;
}

.btn-secondary {
  @apply inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md shadow-sm text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed;
}
</style>
