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
            @click="confirmReception"
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

    <!-- Formulario de Acciones para Insumo Recepcionado -->
    <div v-if="isReceived(scannedProduct)" class="bg-white rounded-lg shadow-sm border p-6">
      <h3 class="text-lg font-medium text-gray-900 mb-4 flex items-center">
        <svg class="h-5 w-5 mr-2 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
        </svg>
        Insumo Recepcionado - Seleccionar Acción
      </h3>

      <div class="space-y-6">
        <!-- Información del Usuario -->
        <div class="bg-gray-50 rounded-lg p-4">
          <h4 class="text-sm font-medium text-gray-900 mb-2">Usuario</h4>
          <p class="text-sm text-gray-600">RUT: {{ currentUser?.rut || 'No disponible' }}</p>
          <p class="text-sm text-gray-600">Nombre: {{ currentUser?.name || 'No disponible' }}</p>
        </div>

        <!-- Selección de Acción -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            ¿Qué desea hacer con este insumo? <span class="text-red-500">*</span>
          </label>
          <div class="space-y-3">
            <div class="flex items-center p-3 border rounded-lg hover:bg-gray-50 cursor-pointer" 
                 :class="{ 'border-green-500 bg-green-50': actionForm.action === 'consume' }"
                 @click="actionForm.action = 'consume'">
              <input 
                type="radio" 
                v-model="actionForm.action" 
                value="consume"
                class="h-4 w-4 text-green-600 focus:ring-green-500 border-gray-300"
              />
              <div class="ml-3">
                <label class="text-sm font-medium text-gray-900 cursor-pointer">
                  ✅ Consumir insumo
                </label>
                <p class="text-xs text-gray-500">Marcar como utilizado en la cirugía</p>
              </div>
            </div>
            <div class="flex items-center p-3 border rounded-lg hover:bg-gray-50 cursor-pointer"
                 :class="{ 'border-orange-500 bg-orange-50': actionForm.action === 'return' }"
                 @click="actionForm.action = 'return'">
              <input 
                type="radio" 
                v-model="actionForm.action" 
                value="return"
                class="h-4 w-4 text-orange-600 focus:ring-orange-500 border-gray-300"
              />
              <div class="ml-3">
                <label class="text-sm font-medium text-gray-900 cursor-pointer">
                  📦 Devolver a bodega
                </label>
                <p class="text-xs text-gray-500">El insumo no fue utilizado y será devuelto</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Notas de la Acción -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Notas {{ actionForm.action === 'return' ? '(opcional)' : '' }}
          </label>
          <textarea
            v-model="actionForm.notes"
            rows="3"
            :placeholder="actionForm.action === 'consume' ? 'Observaciones sobre el consumo...' : 'Motivo de la devolución...'"
            class="form-textarea w-full"
          ></textarea>
        </div>

        <!-- Botones de Acción -->
        <div class="flex justify-end gap-3">
          <button
            @click="resetForm"
            class="btn-secondary"
          >
            Cancelar
          </button>
          <button
            v-if="actionForm.action === 'consume'"
            @click="consumeSupply"
            :disabled="processing"
            class="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded-lg font-medium disabled:opacity-50 disabled:cursor-not-allowed flex items-center"
          >
            <svg v-if="processing" class="animate-spin h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ processing ? 'Consumiendo...' : 'Consumir Insumo' }}
          </button>
          <button
            v-if="actionForm.action === 'return'"
            @click="returnSupply"
            :disabled="processing"
            class="bg-orange-600 hover:bg-orange-700 text-white px-4 py-2 rounded-lg font-medium disabled:opacity-50 disabled:cursor-not-allowed flex items-center"
          >
            <svg v-if="processing" class="animate-spin h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ processing ? 'Devolviendo...' : 'Devolver a Bodega' }}
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
            <!--<button @click="resetForm" class="btn-primary text-sm">
              Recepcionar Otro Insumo
            </button>-->
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

    <!-- Modal de Confirmación de Recepción -->
    <div v-if="showReceptionConfirmModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50 flex items-center justify-center">
      <div class="relative bg-white rounded-lg shadow-xl max-w-md w-full mx-4">
        <div class="p-6">
          <div class="flex items-start mb-4">
            <div class="flex-shrink-0">
              <svg class="h-6 w-6 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
            </div>
            <div class="ml-3 flex-1">
              <h3 class="text-lg font-medium text-gray-900">
                Confirmar Recepción de Insumo
              </h3>
              <div class="mt-2 text-sm text-gray-500">
                <p>¿Está seguro de que desea marcar este insumo como recepcionado?</p>
                <div class="mt-3 bg-gray-50 rounded-lg p-3">
                  <p class="font-medium text-gray-900">{{ 
                    scannedProduct?.supply_code?.name || 
                    scannedProduct?.supply_info?.SupplyCode?.name || 
                    scannedProduct?.supply_info?.name || 
                    'Insumo' 
                  }}</p>
                  <p class="text-xs text-gray-600 mt-1">QR: {{ scannedProduct?.qr_code }}</p>
                </div>
                <p class="mt-3 text-yellow-700 font-medium">
                  Una vez recepcionado, deberá decidir si consumirlo o devolverlo a bodega.
                </p>
              </div>
            </div>
          </div>
          <div class="flex gap-3 mt-6">
            <button
              @click="showReceptionConfirmModal = false"
              class="flex-1 btn-secondary"
            >
              Cancelar
            </button>
            <button
              @click="proceedWithReception"
              class="flex-1 btn-primary"
            >
              Sí, Recepcionar
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
import { useNotificationStore } from '@/stores/notification'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const notificationStore = useNotificationStore()

// Estado del componente
const loading = ref(false)
const qrInput = ref('')
const error = ref(null)
const scannedProduct = ref(null)
const receiving = ref(false)
const processing = ref(false)
const consumptionSuccess = ref(null) // Mantenemos el nombre para conservar la alerta
const showReceptionConfirmModal = ref(false)

// Formulario de recepción simplificado
const receptionForm = ref({
  notes: ''
})

// Formulario de acciones para insumos recepcionados
const actionForm = ref({
  action: 'consume', // 'consume' o 'return'
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
  
  if (status !== 'en_camino_a_pabellon') return false
  
  // Solo los usuarios con rol 'pabellón' tienen restricción de pabellón
  const userRole = currentUser.value?.role
  if (userRole !== 'pabellón') return true
  
  const userPavilionId = currentUser.value?.pavilion_id
  
  // Usuario con rol 'pabellón' sin pabellón asignado no puede recepcionar nada
  if (userPavilionId == null) return false
  
  // Obtener el pabellón destino: campo explícito > solicitud > location del insumo
  const destinationPavilionId = product.destination_pavilion_id
    ?? product.supply_request?.pavilion_id
    ?? product.supply_info?.location_id
  
  // Si no podemos determinar el pabellón destino, bloquear por seguridad
  if (destinationPavilionId == null || destinationPavilionId === 0) return false
  
  return Number(userPavilionId) === Number(destinationPavilionId)
}

// Verificar si está recepcionado (para mostrar opciones de consumo/devolución)
const isReceived = (product) => {
  if (!product || product.type !== 'medical_supply') return false
  if (product.is_consumed) return false
  
  const status = product.supply_info?.Status || 
                 product.supply_info?.status || 
                 product.status || 
                 product.current_status
  
  console.log('isReceived: Estado encontrado:', status)
  
  if (status !== 'recepcionado') return false
  
  // Solo los usuarios con rol 'pabellón' tienen restricción de pabellón
  const userRole = currentUser.value?.role
  if (userRole !== 'pabellón') return true
  
  const userPavilionId = currentUser.value?.pavilion_id
  
  // Usuario con rol 'pabellón' sin pabellón asignado no puede recepcionar nada
  if (userPavilionId == null) return false
  
  // Obtener el pabellón destino: campo explícito > solicitud > location del insumo
  const destinationPavilionId = product.destination_pavilion_id
    ?? product.supply_request?.pavilion_id
    ?? product.supply_info?.location_id
  
  // Si no podemos determinar el pabellón destino, bloquear por seguridad
  if (destinationPavilionId == null || destinationPavilionId === 0) return false
  
  return Number(userPavilionId) === Number(destinationPavilionId)
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

  if (status === 'en_camino_a_pabellon') {
    // El estado es correcto pero el pabellón no coincide
    if (currentUser.value?.role === 'pabellón') {
      const userPavilionId = currentUser.value?.pavilion_id
      const destinationPavilionId = product.destination_pavilion_id
        ?? product.supply_request?.pavilion_id
        ?? product.supply_info?.location_id
      if (userPavilionId == null) {
        return 'Su usuario no tiene un pabellón asignado. Contacte al administrador.'
      }
      if (destinationPavilionId != null && Number(userPavilionId) !== Number(destinationPavilionId)) {
        return 'Este insumo fue solicitado para otro pabellón. Solo puede ser recepcionado por el usuario del pabellón asignado en la solicitud.'
      }
    }
  }

  if (status !== 'en_camino_a_pabellon') {
    return `El insumo tiene estado "${status}" y no está en camino al pabellón. Solo se pueden recepcionar insumos con estado "en_camino_a_pabellon".`
  }

  // El insumo está en camino pero verificar si es para este pabellón
  const destinationPavilionId = product.supply_info?.LocationID
  const userPavilionId = currentUser.value?.pavilion_id

  if (!userPavilionId) {
    return 'Su usuario no tiene un pabellón asignado. Contacte al administrador.'
  }

  if (Number(destinationPavilionId) !== Number(userPavilionId)) {
    return `Este insumo está destinado a otro pabellón (ID: ${destinationPavilionId}). Solo puede ser recepcionado por un usuario de ese pabellón.`
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

// Mostrar modal de confirmación
const confirmReception = () => {
  showReceptionConfirmModal.value = true
}

// Proceder con la recepción después de confirmar
const proceedWithReception = () => {
  showReceptionConfirmModal.value = false
  receiveSupply()
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
    
    const result = await qrService.receiveSupply(
      qrInput.value.trim(),
      userRUT,
      'pavilion',
      pavilionId,
      receptionForm.value.notes,
      true // Siempre recepcionado (el usuario decidirá después si consumir o devolver)
    )
    
    if (result.success) {
      // Mostrar notificación de éxito
      notificationStore.success('El insumo ha sido recepcionado correctamente')
      
      // Redirigir a la vista de QR
      router.push('/qr')
    } else {
      throw new Error(result.error || 'Error al recepcionar el insumo')
    }
    
  } catch (err) {
    console.error('Error receiving supply:', err)
    const errorMessage = err.response?.data?.error || err.message || 'Error al recepcionar el insumo'
    notificationStore.error(errorMessage)
    error.value = errorMessage
  } finally {
    receiving.value = false
  }
}

// Consumir insumo (para insumos ya recepcionados)
const consumeSupply = async () => {
  if (!scannedProduct.value || !qrInput.value.trim()) return
  
  processing.value = true
  error.value = null
  consumptionSuccess.value = null
  
  try {
    const userRUT = currentUser.value?.rut
    if (!userRUT) {
      throw new Error('No se pudo obtener el RUT del usuario')
    }

    const pavilionId = currentUser.value?.pavilion_id || 1

    const result = await qrService.consumeSupply({
      qr_code: qrInput.value.trim(),
      user_rut: userRUT,
      user_name: currentUser.value?.name || 'Usuario',
      destination_type: 'pavilion',
      destination_id: pavilionId,
      notes: actionForm.value.notes
    })
    
    if (result.success) {
      consumptionSuccess.value = {
        ...result,
        qr_code: scannedProduct.value.qr_code,
        batch_id: scannedProduct.value.supply_info?.batch?.id,
        status_change: {
          from: 'Recepcionado',
          to: 'Consumido'
        }
      }
      
      resetForm()
    } else {
      throw new Error(result.error || 'Error al consumir el insumo')
    }
    
  } catch (err) {
    console.error('Error consuming supply:', err)
    error.value = err.response?.data?.error || err.message || 'Error al consumir el insumo'
  } finally {
    processing.value = false
  }
}

// Devolver insumo a bodega (para insumos ya recepcionados)
const returnSupply = async () => {
  if (!scannedProduct.value || !qrInput.value.trim()) return
  
  processing.value = true
  error.value = null
  consumptionSuccess.value = null
  
  try {
    const userRUT = currentUser.value?.rut
    if (!userRUT) {
      throw new Error('No se pudo obtener el RUT del usuario')
    }

    // Obtener el ID de la bodega destino
    const storeId = currentUser.value?.store_id || 1

    const result = await qrService.transferSupply({
      qr_code: qrInput.value.trim(),
      user_rut: userRUT,
      receiver_rut: userRUT, // Por ahora el mismo usuario
      destination_type: 'store',
      destination_id: storeId,
      notes: actionForm.value.notes || 'Devolución desde pabellón'
    })
    
    if (result.success || result.transfer_mode) {
      consumptionSuccess.value = {
        ...result,
        qr_code: scannedProduct.value.qr_code,
        batch_id: scannedProduct.value.supply_info?.batch?.id,
        status_change: {
          from: 'Recepcionado',
          to: 'Disponible en Bodega'
        }
      }
      
      resetForm()
    } else {
      throw new Error(result.error || 'Error al devolver el insumo')
    }
    
  } catch (err) {
    console.error('Error returning supply:', err)
    error.value = err.response?.data?.error || err.message || 'Error al devolver el insumo'
  } finally {
    processing.value = false
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
  actionForm.value = {
    action: 'consume',
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