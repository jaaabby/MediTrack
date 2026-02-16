<template>
  <div class="max-w-4xl mx-auto p-6">
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
            <div v-if="canPickup(scannedProduct)" class="flex items-center space-x-2">
              <div class="w-3 h-3 bg-blue-400 rounded-full"></div>
              <span class="text-sm font-medium text-blue-800">Listo para Retirar</span>
            </div>
            <div v-else class="flex items-center space-x-2">
              <div class="w-3 h-3 bg-red-400 rounded-full"></div>
              <span class="text-sm font-medium text-red-800">No se puede Retirar</span>
            </div>
          </div>
        </div>

        <!-- Status Validation -->
        <div v-if="!canPickup(scannedProduct)" class="bg-red-50 border border-red-200 rounded-md p-4 mt-4">
          <div class="flex">
            <div class="flex-shrink-0">
              <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div class="ml-3">
              <h3 class="text-sm font-medium text-red-800">No se puede retirar</h3>
              <div class="mt-2 text-sm text-red-700">
                {{ getPickupErrorMessage(scannedProduct) }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Formulario de Retiro -->
    <div v-if="canPickup(scannedProduct)" class="bg-white rounded-lg shadow-sm border p-6">
      <h3 class="text-lg font-medium text-gray-900 mb-4 flex items-center">
        <svg class="h-5 w-5 mr-2 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        Formulario de Retiro
      </h3>

      <div class="space-y-6">
        <!-- Información del Usuario -->
        <div class="bg-gray-50 rounded-lg p-4">
          <h4 class="text-sm font-medium text-gray-900 mb-2">Usuario</h4>
          <p class="text-sm text-gray-600">RUT: {{ currentUser?.rut || 'No disponible' }}</p>
          <p class="text-sm text-gray-600">Nombre: {{ currentUser?.name || 'No disponible' }}</p>
        </div>

        <!-- Notas de Retiro -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Notas de Retiro (opcional)
          </label>
          <textarea
            v-model="pickupForm.notes"
            rows="3"
            placeholder="Observaciones sobre el retiro del insumo..."
            class="form-textarea w-full"
          ></textarea>
        </div>

        <!-- Botón de Retiro -->
        <div class="flex justify-end">
          <button
            @click="confirmPickup"
            :disabled="pickingUp"
            class="btn-primary"
          >
            <svg v-if="pickingUp" class="animate-spin h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ pickingUp ? 'Procesando retiro...' : 'Retirar de Bodega' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Modal de Confirmación de Retiro -->
    <div v-if="showPickupConfirmModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50 flex items-center justify-center">
      <div class="relative bg-white rounded-lg shadow-xl max-w-md w-full mx-4">
        <div class="p-6">
          <div class="flex items-start mb-4">
            <div class="flex-shrink-0">
              <svg class="h-6 w-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div class="ml-3 flex-1">
              <h3 class="text-lg font-medium text-gray-900">
                Confirmar Retiro de Insumo
              </h3>
              <div class="mt-2 text-sm text-gray-500">
                <p>¿Está seguro de que desea retirar este insumo de la bodega?</p>
                <div class="mt-3 bg-gray-50 rounded-lg p-3">
                  <p class="font-medium text-gray-900">{{ 
                    scannedProduct?.supply_code?.name || 
                    scannedProduct?.supply_info?.SupplyCode?.name || 
                    scannedProduct?.supply_info?.name || 
                    'Insumo' 
                  }}</p>
                  <p class="text-xs text-gray-600 mt-1">QR: {{ scannedProduct?.qr_code }}</p>
                </div>
                <p class="mt-3 text-blue-700 font-medium">
                  El insumo quedará marcado como "En camino a pabellón".
                </p>
              </div>
            </div>
          </div>
          <div class="flex gap-3 mt-6">
            <button
              @click="showPickupConfirmModal = false"
              class="flex-1 btn-secondary"
            >
              Cancelar
            </button>
            <button
              @click="proceedWithPickup"
              class="flex-1 btn-primary"
            >
              Sí, Retirar
            </button>
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
import { useNotification } from '@/composables/useNotification'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const { success: showSuccess, error: showError } = useNotification()

// Estado del componente
const loading = ref(false)
const qrInput = ref('')
const error = ref(null)
const scannedProduct = ref(null)
const pickingUp = ref(false)
const pickupSuccess = ref(null)
const showPickupConfirmModal = ref(false)

// Formulario de retiro
const pickupForm = ref({
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
  pickupSuccess.value = null
  
  try {
    // Validar formato
    if (!isValidQRFormat(qrInput.value.trim())) {
      throw new Error('Formato de código QR inválido. Use SUPPLY_... para insumos individuales.')
    }
    
    // Crear contexto de escaneo
    const scanContext = {
      scan_purpose: 'pickup',
      scan_source: 'pickup_view',
      user_agent: navigator.userAgent
    }
    
    // Escanear código QR
    const result = await qrService.scanQRCode(qrInput.value.trim(), scanContext)
    
    console.log('Resultado del escaneo QR:', result)

    if (result) {
      scannedProduct.value = result

      // Verificar si es un insumo individual
      if (result.type !== 'medical_supply') {
        error.value = 'Solo se pueden retirar insumos individuales. Este código QR corresponde a un lote.'
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

// Verificar si se puede retirar
const canPickup = (product) => {
  if (!product || product.type !== 'medical_supply') return false
  if (product.is_consumed) return false
  
  const status = product.supply_info?.Status || 
                 product.supply_info?.status || 
                 product.status || 
                 product.current_status
  
  console.log('canPickup: Estado encontrado:', status)
  return status === 'pendiente_retiro'
}

// Obtener mensaje de error para retiro
const getPickupErrorMessage = (product) => {
  if (!product) return ''

  if (product.is_consumed) {
    return 'Este insumo ya ha sido consumido anteriormente.'
  }

  const status = product.supply_info?.Status || 
                 product.supply_info?.status || 
                 product.status || 
                 product.current_status

  if (status === 'disponible') {
    return 'Este insumo está disponible pero no asignado para retiro.'
  }

  if (status === 'en_camino_a_pabellon') {
    return 'Este insumo ya fue retirado y está en camino al pabellón.'
  }

  if (status === 'recepcionado') {
    return 'Este insumo ya fue recepcionado en el pabellón.'
  }

  if (status !== 'pendiente_retiro') {
    return `El insumo tiene estado "${status}" y no está pendiente de retiro. Solo se pueden retirar insumos con estado "pendiente_retiro".`
  }

  return 'Este insumo no se puede retirar.'
}

// Obtener texto del estado
const getStatusText = (status) => {
  const statusMap = {
    'disponible': 'Disponible',
    'pendiente_retiro': 'Pendiente de Retiro',
    'en_camino_a_pabellon': 'En Camino a Pabellón',
    'recepcionado': 'Recepcionado',
    'consumido': 'Consumido'
  }
  return statusMap[status] || status || 'Desconocido'
}

// Mostrar modal de confirmación
const confirmPickup = () => {
  showPickupConfirmModal.value = true
}

// Proceder con el retiro después de confirmar
const proceedWithPickup = () => {
  showPickupConfirmModal.value = false
  pickupSupply()
}

// Retirar insumo
const pickupSupply = async () => {
  if (!scannedProduct.value || !qrInput.value.trim()) return
  
  pickingUp.value = true
  error.value = null
  pickupSuccess.value = null
  
  try {
    const userRUT = currentUser.value?.rut
    if (!userRUT) {
      throw new Error('No se pudo obtener el RUT del usuario')
    }

    // Guardar el estado anterior
    const previousStatus = scannedProduct.value.supply_info?.Status || 
                          scannedProduct.value.supply_info?.status || 
                          scannedProduct.value.status || 
                          scannedProduct.value.current_status || 
                          'pendiente_retiro'
    
    const result = await qrService.pickupSupplyFromStore(
      qrInput.value.trim(),
      userRUT,
      pickupForm.value.notes
    )
    
    console.log('Resultado de pickupSupplyFromStore:', result)
    
    if (result.success) {
      console.log('Retiro exitoso, redirigiendo...')
      
      // Mostrar notificación de éxito y redirigir
      showSuccess('El insumo ha sido retirado correctamente de la bodega')
      
      // Redirigir a la vista de QR
      setTimeout(() => {
        router.push('/qr')
      }, 1500)
      
    } else {
      throw new Error(result.error || 'Error al retirar el insumo')
    }
    
  } catch (err) {
    console.error('Error picking up supply:', err)
    const errorMessage = err.response?.data?.error || err.message || 'Error al retirar el insumo'
    showError(errorMessage)
    
    // Redirigir a la vista de QR después de mostrar el error
    setTimeout(() => {
      router.push('/qr')
    }, 2000)
  } finally {
    pickingUp.value = false
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
    case 'pendiente_retiro': return 'Pendiente de Retiro'
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
    case 'pendiente_retiro': return 'bg-blue-100 text-blue-800'
    case 'en_camino_a_pabellon': return 'bg-yellow-100 text-yellow-800'
    case 'recepcionado': return 'bg-purple-100 text-purple-800'
    default: return 'bg-gray-100 text-gray-800'
  }
}

// Limpiar formulario
const resetForm = () => {
  pickupForm.value = {
    notes: ''
  }
  qrInput.value = ''
  scannedProduct.value = null
  error.value = null
  pickupSuccess.value = null
}

// Limpiar error
const clearError = () => {
  error.value = null
}
</script>

<style scoped>
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
</style>
