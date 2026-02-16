<template>
  <div class="max-w-4xl mx-auto p-4 sm:p-6 lg:p-8 space-y-6">
    <!-- Header -->
    <div class="bg-white shadow-sm rounded-lg border">
      <div class="px-4 py-5 sm:p-6">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-900">Transferir Insumo</h1>
            <p class="mt-1 text-sm text-gray-600">
              Transferir un insumo disponible a un pabellón
            </p>
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
    </div>

    <!-- Scanner Section -->
    <!--<div class="bg-white shadow-sm rounded-lg border">
      <div class="px-4 py-5 sm:p-6">
        <h2 class="text-lg font-medium text-gray-900 mb-4">Escanear Código QR</h2>
        
        QR Input
        <div class="flex space-x-3 mb-4">
          <div class="flex-1">
            <input
              type="text"
              v-model="qrInput"
              placeholder="Ingrese o escanee el código QR del insumo"
              class="form-input w-full"
              @keyup.enter="scanQR"
              :disabled="scanning"
            />
          </div>
          <button
            @click="scanQR"
            :disabled="!qrInput.trim() || scanning"
            class="btn-primary"
          >
            <div v-if="scanning" class="flex items-center">
              <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Escaneando...
            </div>
            <div v-else class="flex items-center">
              <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v2a2 2 0 002 2zm0 0h2a2 2 0 012 2v2a2 2 0 01-2 2H5a2 2 0 01-2-2v-2a2 2 0 012-2z" />
              </svg>
              Escanear
            </div>
          </button>
        </div>

        Error Display
        <div v-if="error" class="mb-4 bg-red-50 border border-red-200 rounded-md p-4">
          <div class="flex">
            <div class="flex-shrink-0">
              <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div class="ml-3">
              <h3 class="text-sm font-medium text-red-800">Error</h3>
              <div class="mt-2 text-sm text-red-700">{{ error }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>-->

    <!-- Scanned Product Info -->
    <div v-if="scannedProduct" class="bg-white shadow-sm rounded-lg border">
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
                <p class="text-sm text-gray-900">{{ scannedProduct.supply_info?.name || 'No disponible' }}</p>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-500">Estado</label>
                <span :class="getStatusBadgeClass(scannedProduct)" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium">
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

        <!-- Status Validation -->
        <div v-if="!canTransfer(scannedProduct)" class="bg-red-50 border border-red-200 rounded-md p-4 mb-6">
          <div class="flex">
            <div class="flex-shrink-0">
              <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div class="ml-3">
              <h3 class="text-sm font-medium text-red-800">No se puede transferir</h3>
              <div class="mt-2 text-sm text-red-700">
                {{ getTransferErrorMessage(scannedProduct) }}
              </div>
            </div>
          </div>
        </div>

        <!-- Transfer Form -->
        <div v-if="canTransfer(scannedProduct)" class="space-y-6">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Usuario Responsable -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                RUT Usuario Responsable <span class="text-red-500">*</span>
              </label>
              <input
                type="text"
                v-model="transferForm.userRUT"
                placeholder="12.345.678-9"
                class="form-input w-full"
                required
              />
            </div>

            <!-- RUT Persona que Recibe -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                RUT Persona que Recibe <span class="text-red-500">*</span>
              </label>
              <input
                type="text"
                v-model="transferForm.receiverRUT"
                placeholder="12.345.678-9"
                class="form-input w-full"
                required
              />
            </div>



            <!-- Pabellón -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Pabellón <span class="text-red-500">*</span>
              </label>
              <select v-model="transferForm.destinationID" class="form-select w-full" required
                :disabled="loadingPavilions">
                <option value="">
                  {{ loadingPavilions ? 'Cargando pabellones...' : 
                     pavilions.length === 0 ? 'No hay pabellones disponibles' : 'Seleccionar pabellón' }}
                </option>
                <option v-for="pavilion in pavilions" :key="pavilion.id" :value="pavilion.id">
                  {{ pavilion.name }}
                </option>
              </select>
            </div>
          </div>

          <!-- Notas -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Notas Adicionales
            </label>
            <textarea
              v-model="transferForm.notes"
              rows="3"
              placeholder="Motivo de la transferencia, instrucciones especiales, etc."
              class="form-textarea w-full"
            ></textarea>
          </div>

          <!-- Transfer Button -->
          <div class="flex justify-end space-x-3">
            <button @click="clearScannedProduct" class="btn-secondary">
              Cancelar
            </button>
            <button
              @click="transferProduct"
              :disabled="!validateTransferForm() || transferring"
              class="btn-primary"
            >
              <div v-if="transferring" class="flex items-center">
                <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Transfiriendo...
              </div>
              <div v-else class="flex items-center">
                <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
                </svg>
                Transferir Insumo
              </div>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Success Message -->
    <div v-if="transferSuccess" class="bg-green-50 border border-green-200 rounded-md p-4">
      <div class="flex">
        <div class="flex-shrink-0">
          <svg class="h-5 w-5 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
        </div>
        <div class="ml-3">
          <h3 class="text-sm font-medium text-green-800">Transferencia Exitosa</h3>
          <div class="mt-2 text-sm text-green-700">
            El insumo se encuentra en tránsito a pabellón.
            <div v-if="transferSuccess.status_change" class="mt-1 text-xs text-green-600">
              Estado cambiado de "{{ transferSuccess.status_change.from }}" a "{{ transferSuccess.status_change.to }}"
            </div>
          </div>
          <div class="mt-4 flex flex-wrap gap-3">
            <button @click="resetForm" class="btn-primary text-sm">
              Transferir Otro Insumo
            </button>
            <router-link to="/qr" class="btn-secondary text-sm">
              <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
              Volver al Scanner
            </router-link>
            <router-link :to="`/qr/${transferSuccess.qr_code}/traceability`" class="btn-secondary text-sm">
              Ver Trazabilidad
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
import pavilionService from '@/services/config/pavilionService'
import { useAuthStore } from '@/stores/auth'
import { useNotification } from '@/composables/useNotification'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const { success: showSuccess, error: showError } = useNotification()

// Estado reactivo
const qrInput = ref('')
const scanning = ref(false)
const error = ref(null)
const scannedProduct = ref(null)
const transferring = ref(false)
const transferSuccess = ref(null)

// Datos de configuración
const pavilions = ref([])
const loadingPavilions = ref(false)

// Formulario de transferencia
const transferForm = ref({
  userRUT: '',
  receiverRUT: '',
  destinationID: '',
  notes: ''
})

// Usuario actual
const currentUser = computed(() => authStore.user)

// Auto-completar RUT del usuario actual
onMounted(() => {
  if (currentUser.value?.rut) {
    transferForm.value.userRUT = currentUser.value.rut
  }
  loadPavilions()
  
  // Si hay un QR en la URL, escanearlo automáticamente
  const qrFromUrl = route.query.qr
  if (qrFromUrl) {
    qrInput.value = qrFromUrl
    scanQR()
  }
})


// Cargar todos los pabellones desde la base de datos
const loadPavilions = async () => {
  try {
    loadingPavilions.value = true
    console.log('🏢 Cargando pabellones desde BD...')
    
    const allPavilions = await pavilionService.getAllPavilions()
    console.log('🏢 Respuesta completa del servicio de pabellones:', allPavilions)
    
    if (allPavilions && Array.isArray(allPavilions) && allPavilions.length > 0) {
      pavilions.value = allPavilions
      console.log('✅ Pabellones cargados desde BD:', pavilions.value)
    } else {
      console.warn('⚠️ No se encontraron pabellones en la BD, usando fallback')
      throw new Error('No data found')
    }
  } catch (error) {
    console.error('❌ Error al cargar pabellones desde BD:', error)
    console.log('🔄 Usando datos fallback para pabellones')
    // Fallback con datos por defecto
    pavilions.value = [
      { id: 1, name: 'Pabellón A', medical_center_id: 1 },
      { id: 2, name: 'Pabellón B', medical_center_id: 2 }
    ]
  } finally {
    loadingPavilions.value = false
  }
}


// Escanear QR
const scanQR = async () => {
  if (!qrInput.value.trim()) return

  scanning.value = true
  error.value = null
  scannedProduct.value = null
  transferSuccess.value = null
  
  try {
    // Escanear con contexto específico para evitar duplicados innecesarios
    const result = await qrService.scanQRCode(qrInput.value.trim(), {
      scan_purpose: 'transfer_verification', // Más específico
      scan_source: 'web', // Usar valor válido permitido por la BD
      prevent_duplicate_logging: true // Flag para evitar logging duplicado
    })
    
    if (result) {
      scannedProduct.value = result
      
      // Debug: mostrar información del producto escaneado
      console.log('Producto escaneado para transferencia:', result)
      console.log('Tipo:', result.type)
      console.log('Estado:', result.supply_info?.status || result.status)
      console.log('Es consumido:', result.is_consumed)
      console.log('Puede transferir:', canTransfer(result))
      
      // Verificar si es un insumo individual
      if (result.type !== 'medical_supply') {
        error.value = 'Solo se pueden transferir insumos individuales. Este código QR corresponde a un lote.'
        scannedProduct.value = null
        return
      }
    }
  } catch (err) {
    error.value = err.message || 'Error al escanear código QR'
    console.error('Error al escanear:', err)
  } finally {
    scanning.value = false
  }
}// Verificar si se puede transferir
const canTransfer = (product) => {
  if (!product) return false
  
  // Verificar que no esté consumido
  if (product.is_consumed) return false
  
  // Verificar el estado específico - se puede transferir desde "disponible" o "recepcionado"
  // Buscar el estado en diferentes lugares de la respuesta
  const status = product.supply_info?.status || 
                 product.status || 
                 product.current_status ||
                 product.supply_info?.Status ||
                 product.Status
  
  console.log('Estado encontrado para transferencia:', status)
  
  // Permitir transferencia si es "disponible", "recepcionado" o si no hay estado definido (asumir disponible)
  return status === 'disponible' || 
         status === 'available' || 
         status === 'recepcionado' || 
         status === 'received' || 
         !status
}

// Obtener mensaje de error para transferencia
const getTransferErrorMessage = (product) => {
  if (!product) return ''
  
  if (product.is_consumed) {
    return 'Este insumo ya ha sido consumido y no puede ser transferido.'
  }
  
  const status = product.supply_info?.status || 
                 product.status || 
                 product.current_status ||
                 product.supply_info?.Status ||
                 product.Status
                 
  if (status === 'en_camino_a_pabellon' || status === 'en_camino_a_bodega') {
    return 'Este insumo ya está en tránsito y no puede ser transferido nuevamente.'
  }
  
  if (status === 'consumido') {
    return 'Este insumo ya ha sido consumido y no puede ser transferido.'
  }
  
  if (status && status !== 'disponible' && status !== 'available' && status !== 'recepcionado' && status !== 'received') {
    return `El insumo tiene estado "${status}" y no está disponible para transferencia. Solo se pueden transferir insumos con estado "disponible" o "recepcionado".`
  }
  
  return 'Este insumo no se puede transferir por razones desconocidas.'
}

// Validar formulario de transferencia
const validateTransferForm = () => {
  if (!transferForm.value.userRUT.trim()) return false
  if (!transferForm.value.receiverRUT.trim()) return false
  if (!transferForm.value.destinationID) return false
  
  const destinationIdNum = parseInt(transferForm.value.destinationID)
  if (isNaN(destinationIdNum) || destinationIdNum < 1) return false
  
  return true
}

// Transferir producto
const transferProduct = async () => {
  if (!scannedProduct.value || !validateTransferForm()) return
  
  // Prevenir dobles clics
  if (transferring.value) {
    console.log('DEBUG - Transferencia ya en progreso, ignorando doble clic')
    return
  }
  
  transferring.value = true
  error.value = null
  
  try {
    console.log('DEBUG - Iniciando transferencia para QR:', scannedProduct.value.qr_code)
    
    const transferData = {
      qr_code: scannedProduct.value.qr_code,
      user_rut: transferForm.value.userRUT,
      receiver_rut: transferForm.value.receiverRUT,
      user_name: currentUser.value?.name || 'Usuario',
      destination_type: 'pavilion', // Siempre pabellón
      destination_id: parseInt(transferForm.value.destinationID),
      notes: transferForm.value.notes,
      transfer_timestamp: new Date().toISOString(),
      transfer_context: {
        scan_source: 'transfer_view',
        user_agent: navigator.userAgent
      }
    }
    
    console.log('DEBUG - Datos de transferencia:', transferData)
    
    const result = await qrService.transferSupply(transferData)
    
    console.log('DEBUG - Resultado de transferencia:', result)
    
    if (result.success) {
      transferSuccess.value = {
        ...result.data,
        qr_code: scannedProduct.value.qr_code
      }
      
      // Mostrar notificación de éxito
      showSuccess('Insumo transferido correctamente')
      
      // Limpiar formulario
      scannedProduct.value = null
      qrInput.value = ''
    } else {
      error.value = result.error || 'Error al transferir el insumo'
      showError(result.error || 'Error al transferir el insumo')
    }
  } catch (err) {
    console.error('DEBUG - Error en transferencia:', err)
    error.value = err.message || 'Error al transferir el insumo'
    showError(err.message || 'Error al transferir el insumo')
  } finally {
    transferring.value = false
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
  
  const status = product.supply_info?.status || 
                 product.status || 
                 product.current_status ||
                 product.supply_info?.Status ||
                 product.Status
                 
  switch (status) {
    case 'disponible': 
    case 'available': 
      return 'Disponible'
    case 'recepcionado': 
    case 'received': 
      return 'Recepcionado'
    case 'en_camino_a_pabellon':
      return 'En camino a pabellón'
    case 'en_camino_a_bodega':
      return 'En camino a bodega'
    case 'consumido':
      return 'Consumido'
    case 'creado':
      return 'Creado'
    default: 
      return status || 'Sin estado'
  }
}

const getStatusBadgeClass = (product) => {
  if (!product) return 'bg-gray-100 text-gray-800'
  
  if (product.is_consumed) return 'bg-red-100 text-red-800'
  
  const status = product.supply_info?.status || 
                 product.status || 
                 product.current_status ||
                 product.supply_info?.Status ||
                 product.Status
                 
  switch (status) {
    case 'disponible': 
    case 'available': 
      return 'bg-green-100 text-green-800'
    case 'recepcionado': 
    case 'received': 
      return 'bg-blue-100 text-blue-800'
    case 'en_camino_a_pabellon':
    case 'en_camino_a_bodega':
      return 'bg-yellow-100 text-yellow-800'
    case 'consumido':
      return 'bg-red-100 text-red-800'
    case 'creado':
      return 'bg-purple-100 text-purple-800'
    default: 
      return 'bg-gray-100 text-gray-800'
  }
}

// Limpiar producto escaneado
const clearScannedProduct = () => {
  scannedProduct.value = null
  error.value = null
}

// Resetear formulario completo
const resetForm = () => {
  scannedProduct.value = null
  transferSuccess.value = null
  qrInput.value = ''
  error.value = null
  transferForm.value = {
    userRUT: currentUser.value?.rut || '',
    receiverRUT: '',
    destinationID: '',
    notes: ''
  }
}
</script>

<!-- Los estilos .form-input, .form-select y .form-textarea están definidos globalmente en style.css -->