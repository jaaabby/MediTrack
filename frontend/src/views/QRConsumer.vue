<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <div class="flex items-center justify-between">
        <div>
          <h2 class="text-xl font-semibold text-gray-900">Consumir Insumos Médicos</h2>
          <p class="text-gray-600 mt-1">Escanea o ingresa códigos QR para registrar el consumo de productos</p>
        </div>
        <div class="text-right">
          <p class="text-sm text-gray-500">Productos consumidos hoy</p>
          <p class="text-2xl font-bold text-blue-600">{{ todayConsumptions }}</p>
        </div>
      </div>
    </div>

    <!-- Scanner/Input Section -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <h3 class="text-lg font-medium text-gray-900 mb-4">
        <svg class="h-5 w-5 inline mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h2M4 4h5l2 3h3l2-3h5v5M4 4v5m16-5v5" />
        </svg>
        Escanear Código QR
      </h3>

      <div class="grid md:grid-cols-2 gap-6">
        <!-- QR Input -->
        <div>
          <label for="qrInput" class="block text-sm font-medium text-gray-700 mb-2">
            Código QR del Producto:
          </label>
          <div class="flex space-x-2">
            <input
              id="qrInput"
              v-model="qrInput"
              type="text"
              placeholder="Ej: SUPPLY_1755580808_abc123def"
              class="form-input flex-1"
              @keyup.enter="scanQRCode"
            />
            <button
              @click="scanQRCode"
              :disabled="!qrInput.trim() || scanning"
              class="btn-primary"
            >
              <svg v-if="scanning" class="animate-spin h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              {{ scanning ? 'Escaneando...' : 'Escanear' }}
            </button>
          </div>
        </div>

        <!-- Camera Scanner (placeholder) -->
        <div class="text-center">
          <button
            @click="startCameraScanner"
            class="btn-secondary w-full h-20 border-2 border-dashed border-gray-300 hover:border-gray-400"
            :disabled="cameraActive"
          >
            <svg class="h-8 w-8 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
            <span class="text-sm">{{ cameraActive ? 'Cámara Activa' : 'Usar Cámara' }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Consumption Form -->
    <div v-if="scannedProduct && !error" class="bg-white rounded-lg shadow-sm border p-6">
      <h3 class="text-lg font-medium text-gray-900 mb-4">
        <svg class="h-5 w-5 inline mr-2 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        Producto Encontrado - Registrar Consumo
      </h3>

      <!-- Product Info Display -->
      <div class="bg-gray-50 rounded-lg p-4 mb-6">
        <div class="grid md:grid-cols-2 gap-4">
          <div>
            <label class="text-sm font-medium text-gray-600">Nombre del Producto:</label>
            <p class="text-gray-900 font-medium">{{ scannedProduct.supply_info?.supply_code_name || 'N/A' }}</p>
          </div>
          <div>
            <label class="text-sm font-medium text-gray-600">Código QR:</label>
            <p class="text-sm font-mono text-gray-800 bg-white px-2 py-1 rounded border">{{ scannedProduct.qr_code }}</p>
          </div>
          <div>
            <label class="text-sm font-medium text-gray-600">Proveedor:</label>
            <p class="text-gray-900">{{ scannedProduct.supply_info?.supplier || 'N/A' }}</p>
          </div>
          <div>
            <label class="text-sm font-medium text-gray-600">Vencimiento:</label>
            <p class="text-gray-900">{{ formatDate(scannedProduct.supply_info?.expiration_date) }}</p>
          </div>
        </div>

        <!-- Status Alert -->
        <div v-if="scannedProduct.is_consumed" class="mt-4 bg-red-50 border border-red-200 rounded p-3">
          <div class="flex items-center">
            <svg class="h-5 w-5 text-red-400 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span class="text-red-800 font-medium">Este producto ya ha sido consumido anteriormente</span>
          </div>
        </div>

        <div v-else-if="!scannedProduct.can_consume" class="mt-4 bg-yellow-50 border border-yellow-200 rounded p-3">
          <div class="flex items-center">
            <svg class="h-5 w-5 text-yellow-400 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span class="text-yellow-800 font-medium">Este producto no está disponible para consumo</span>
          </div>
        </div>
      </div>

      <!-- Consumption Form -->
      <form v-if="!scannedProduct.is_consumed && scannedProduct.can_consume" @submit.prevent="consumeProduct">
        <div class="grid md:grid-cols-2 gap-6">
          <!-- User RUT -->
          <div>
            <label for="userRUT" class="block text-sm font-medium text-gray-700 mb-2">
              RUT del Usuario <span class="text-red-500">*</span>
            </label>
            <input
              id="userRUT"
              v-model="consumptionForm.userRUT"
              type="text"
              placeholder="12.345.678-9"
              class="form-input"
              required
            />
          </div>

          <!-- Destination Type -->
          <div>
            <label for="destinationType" class="block text-sm font-medium text-gray-700 mb-2">
              Tipo de Destino <span class="text-red-500">*</span>
            </label>
            <select
              id="destinationType"
              v-model="consumptionForm.destinationType"
              class="form-select"
              required
            >
              <option value="">Seleccionar tipo</option>
              <option value="pavilion">Pabellón</option>
              <option value="store">Almacén</option>
            </select>
          </div>

          <!-- Destination ID -->
          <div>
            <label for="destinationID" class="block text-sm font-medium text-gray-700 mb-2">
              ID de Destino <span class="text-red-500">*</span>
            </label>
            <input
              id="destinationID"
              v-model="consumptionForm.destinationID"
              type="number"
              placeholder="ID del pabellón o almacén"
              class="form-input"
              required
            />
          </div>

          <!-- Notes -->
          <div>
            <label for="notes" class="block text-sm font-medium text-gray-700 mb-2">
              Notas (Opcional)
            </label>
            <input
              id="notes"
              v-model="consumptionForm.notes"
              type="text"
              placeholder="Observaciones adicionales"
              class="form-input"
            />
          </div>
        </div>

        <!-- Submit Button -->
        <div class="mt-6 flex justify-end space-x-3">
          <button
            type="button"
            @click="clearScannedProduct"
            class="btn-secondary"
          >
            Cancelar
          </button>
          <button
            type="submit"
            :disabled="consuming"
            class="btn-danger"
          >
            <svg v-if="consuming" class="animate-spin h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <svg v-else class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
            {{ consuming ? 'Consumiendo...' : 'Confirmar Consumo' }}
          </button>
        </div>
      </form>
    </div>

    <!-- Success Message -->
    <div v-if="consumptionSuccess" class="bg-green-50 border border-green-200 rounded-lg p-6">
      <div class="flex items-start">
        <div class="flex-shrink-0">
          <svg class="h-6 w-6 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <div class="ml-3 flex-1">
          <h3 class="text-lg font-medium text-green-800">Producto Consumido Exitosamente</h3>
          <div class="mt-2 text-sm text-green-700">
            <p>{{ consumptionSuccess.message }}</p>
            <p class="mt-1">
              <strong>Cantidad restante en lote:</strong> {{ consumptionSuccess.remaining_amount }} unidades
            </p>
          </div>
          <div class="mt-4">
            <button @click="clearAll" class="btn-secondary text-sm">
              Consumir Otro Producto
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Error Message -->
    <div v-if="error" class="bg-red-50 border border-red-200 rounded-lg p-6">
      <div class="flex items-start">
        <div class="flex-shrink-0">
          <svg class="h-6 w-6 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <div class="ml-3 flex-1">
          <h3 class="text-lg font-medium text-red-800">Error</h3>
          <div class="mt-2 text-sm text-red-700">
            <p>{{ error }}</p>
          </div>
          <div class="mt-4">
            <button @click="clearError" class="btn-secondary text-sm">
              Intentar de Nuevo
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Recent Consumptions -->
    <div v-if="recentConsumptions.length > 0" class="bg-white rounded-lg shadow-sm border p-6">
      <h3 class="text-lg font-medium text-gray-900 mb-4">
        <svg class="h-5 w-5 inline mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        Consumos Recientes
      </h3>

      <div class="space-y-3">
        <div 
          v-for="(consumption, index) in recentConsumptions" 
          :key="index"
          class="flex items-center justify-between p-3 bg-gray-50 rounded"
        >
          <div>
            <p class="font-medium text-gray-900">{{ consumption.product_name }}</p>
            <p class="text-sm text-gray-600">QR: {{ consumption.qr_code }}</p>
          </div>
          <div class="text-right">
            <p class="text-sm font-medium text-gray-900">{{ consumption.user_rut }}</p>
            <p class="text-xs text-gray-500">{{ formatDate(consumption.consumed_at) }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="bg-gray-50 rounded-lg p-6">
      <h3 class="font-medium text-gray-900 mb-4 text-center">Acciones Rápidas</h3>
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <router-link to="/inventory" class="btn-secondary p-4 text-center block hover:bg-gray-200 transition-colors">
          <svg class="h-6 w-6 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
          </svg>
          <div class="text-sm font-medium">Ver Inventario</div>
        </router-link>
        
        <router-link to="/qr" class="btn-secondary p-4 text-center block hover:bg-gray-200 transition-colors">
          <svg class="h-6 w-6 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h2M4 4h5l2 3h3l2-3h5v5M4 4v5m16-5v5" />
          </svg>
          <div class="text-sm font-medium">Escáner QR</div>
        </router-link>

        <button @click="syncBatchAmounts" class="btn-secondary p-4 text-center hover:bg-gray-200 transition-colors">
          <svg class="h-6 w-6 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          <div class="text-sm font-medium">Sincronizar</div>
        </button>

        <router-link to="/" class="btn-secondary p-4 text-center block hover:bg-gray-200 transition-colors">
          <svg class="h-6 w-6 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
          </svg>
          <div class="text-sm font-medium">Ir al Inicio</div>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import qrService from '@/services/qrService'

// Estado reactivo
const qrInput = ref('')
const scanning = ref(false)
const consuming = ref(false)
const cameraActive = ref(false)
const scannedProduct = ref(null)
const consumptionSuccess = ref(null)
const error = ref(null)
const todayConsumptions = ref(0)
const recentConsumptions = ref([])

// Formulario de consumo
const consumptionForm = ref({
  userRUT: '',
  destinationType: '',
  destinationID: '',
  notes: ''
})

// Métodos
const scanQRCode = async () => {
  if (!qrInput.value.trim()) return
  
  scanning.value = true
  error.value = null
  scannedProduct.value = null
  
  try {
    const result = await qrService.scanQRCode(qrInput.value.trim())
    
    if (result.success) {
      if (result.data.type !== 'medical_supply') {
        error.value = 'Solo se pueden consumir productos individuales (no lotes)'
        return
      }
      
      scannedProduct.value = result.data
    } else {
      error.value = result.error || 'Error desconocido al escanear código QR'
    }
  } catch (err) {
    error.value = err.response?.data?.error || 'Error de conexión al escanear código QR'
  } finally {
    scanning.value = false
  }
}

const consumeProduct = async () => {
  if (!scannedProduct.value) return
  
  consuming.value = true
  error.value = null
  
  try {
    const consumptionData = {
      qr_code: scannedProduct.value.qr_code,
      user_rut: consumptionForm.value.userRUT,
      destination_type: consumptionForm.value.destinationType,
      destination_id: parseInt(consumptionForm.value.destinationID),
      notes: consumptionForm.value.notes
    }
    
    const result = await qrService.consumeSupply(consumptionData)
    
    if (result.success) {
      consumptionSuccess.value = result.data || result
      
      // Agregar a consumos recientes
      recentConsumptions.value.unshift({
        qr_code: scannedProduct.value.qr_code,
        product_name: scannedProduct.value.supply_info?.supply_code_name || 'N/A',
        user_rut: consumptionForm.value.userRUT,
        consumed_at: new Date().toISOString()
      })
      
      // Mantener solo los últimos 10
      if (recentConsumptions.value.length > 10) {
        recentConsumptions.value = recentConsumptions.value.slice(0, 10)
      }
      
      // Actualizar contador del día
      todayConsumptions.value += 1
      
      // Limpiar formulario
      clearScannedProduct()
      qrInput.value = ''
      
    } else {
      error.value = result.error || 'Error al consumir el producto'
    }
  } catch (err) {
    error.value = err.response?.data?.error || 'Error de conexión al consumir producto'
  } finally {
    consuming.value = false
  }
}

const startCameraScanner = () => {
  cameraActive.value = true
  // Implementar escáner con cámara aquí
  alert('Funcionalidad de cámara en desarrollo. Por favor usa el input manual.')
  cameraActive.value = false
}

const syncBatchAmounts = async () => {
  try {
    await qrService.syncBatchAmounts()
    alert('Cantidades de lotes sincronizadas correctamente')
  } catch (error) {
    console.error('Error al sincronizar:', error)
    alert('Error al sincronizar las cantidades')
  }
}

const clearScannedProduct = () => {
  scannedProduct.value = null
  consumptionForm.value = {
    userRUT: '',
    destinationType: '',
    destinationID: '',
    notes: ''
  }
}

const clearError = () => {
  error.value = null
}

const clearAll = () => {
  qrInput.value = ''
  scannedProduct.value = null
  consumptionSuccess.value = null
  error.value = null
  clearScannedProduct()
}

const formatDate = (dateString) => {
  if (!dateString) return 'No disponible'
  try {
    return format(new Date(dateString), 'dd/MM/yyyy HH:mm', { locale: es })
  } catch (error) {
    return dateString
  }
}

// Lifecycle
onMounted(() => {
  // Aquí podrías cargar estadísticas iniciales si es necesario
  // Por ejemplo, consumos del día actual
})
</script>

