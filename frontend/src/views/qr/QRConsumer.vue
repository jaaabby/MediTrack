<template>
  <div class="max-w-4xl mx-auto p-4 sm:p-6 lg:p-8 space-y-6">
    <!-- Header -->
    <div class="bg-white shadow-sm rounded-lg border">
      <div class="px-4 py-5 sm:p-6">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-900">Consumir Insumo</h1>
            <p class="mt-1 text-sm text-gray-600">
              Consumir insumos médicos con estado "recepcionado"
            </p>
          </div>
          <div class="flex space-x-2">
            <router-link to="/qr" class="btn-secondary">
              <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
              Volver al Scanner
            </router-link>
            <router-link to="/qr-transfer" class="btn-secondary">
              <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
              </svg>
              Transferir Insumo
            </router-link>
          </div>
        </div>
      </div>
    </div>

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

        <!-- Status Validation -->
        <div v-if="!canConsume(scannedProduct)" class="bg-red-50 border border-red-200 rounded-md p-4 mb-6">
          <div class="flex">
            <div class="flex-shrink-0">
              <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div class="ml-3">
              <h3 class="text-sm font-medium text-red-800">No se puede consumir</h3>
              <div class="mt-2 text-sm text-red-700">
                {{ getConsumptionErrorMessage(scannedProduct) }}
              </div>
              <div v-if="shouldShowTransferOption(scannedProduct)" class="mt-3">
                <router-link to="/qr-transfer" class="btn-primary text-sm">
                  <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
                  </svg>
                  Ir a Transferir
                </router-link>
              </div>
            </div>
          </div>
        </div>

        <!-- Consumption Form -->
        <div v-if="canConsume(scannedProduct)" class="space-y-6">
          <!-- Propósito del Consumo - COMENTADO PARA OCULTAR -->
          <!-- 
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-3">
              Propósito del Consumo <span class="text-red-500">*</span>
            </label>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <div v-for="purpose in consumptionPurposes" :key="purpose.value"
                @click="selectedConsumptionPurpose = purpose.value" :class="[
                  'relative rounded-lg border p-4 cursor-pointer focus:outline-none',
                  selectedConsumptionPurpose === purpose.value
                    ? 'border-blue-500 ring-2 ring-blue-500 bg-blue-50'
                    : 'border-gray-300 hover:border-gray-400'
                ]">
                <div class="flex items-center justify-between">
                  <div class="flex items-center">
                    <div :class="[purpose.iconClass, 'flex-shrink-0 h-6 w-6 mr-3']">
                      <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="purpose.icon" />
                      </svg>
                    </div>
                    <div>
                      <div class="text-sm font-medium text-gray-900">{{ purpose.label }}</div>
                      <div class="text-xs text-gray-500">{{ purpose.description }}</div>
                    </div>
                  </div>
                  <div v-if="selectedConsumptionPurpose === purpose.value" class="text-blue-600">
                    <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd"
                        d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                        clip-rule="evenodd" />
                    </svg>
                  </div>
                </div>
              </div>
            </div>
          </div>
          -->

          <!-- Información del Usuario de Pabellón -->
          <!-- 
            Sección informativa que muestra los datos del usuario de pabellón que está realizando el consumo.
            Los datos se obtienen automáticamente de la sesión del usuario (useAuthStore).
            No es editable ya que se usa para identificar quién está consumiendo el insumo.
          -->
          <div class="bg-gray-50 rounded-lg p-4">
            <h4 class="text-sm font-medium text-gray-900 mb-2">Usuario de Pabellón</h4>
            <!-- RUT del usuario obtenido de la sesión actual -->
            <p class="text-sm text-gray-600">RUT: {{ currentUser?.rut || 'No disponible' }}</p>
            <!-- Nombre del usuario obtenido de la sesión actual -->
            <p class="text-sm text-gray-600">Nombre: {{ currentUser?.name || 'No disponible' }}</p>
          </div>

          <!-- Notas -->
          <!-- 
            Campo opcional para agregar observaciones adicionales sobre el consumo del insumo.
            Permite al usuario de pabellón documentar detalles específicos del uso del insumo médico.
            Las notas se guardan en el historial del insumo para trazabilidad.
          -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Notas Adicionales
            </label>
            <!-- 
              Textarea para ingresar notas opcionales sobre el consumo.
              Se vincula con consumptionForm.notes para capturar el texto ingresado.
            -->
            <textarea v-model="consumptionForm.notes" rows="3" placeholder="Detalles del consumo, observaciones, etc."
              class="form-textarea w-full"></textarea>
          </div>

          <!-- Consumption Button -->
          <!-- 
            Sección de botones de acción para el formulario de consumo.
            Incluye botón de cancelar y botón principal para consumir el insumo.
          -->
          <div class="flex justify-end space-x-3">
            <!-- 
              Botón de cancelar que limpia el producto escaneado y resetea el formulario.
              Permite al usuario cancelar la operación de consumo.
            -->
            <button @click="clearScannedProduct" class="btn-secondary">
              Cancelar
            </button>
            <!-- 
              Botón principal para consumir el insumo médico.
              Se deshabilita si el formulario no es válido o si ya se está procesando el consumo.
              Muestra estado de carga durante el procesamiento.
            -->
            <button @click="consumeProduct" :disabled="!validateConsumptionForm() || consuming" class="btn-primary">
              <!-- Estado de carga: muestra spinner y texto "Consumiendo..." -->
              <div v-if="consuming" class="flex items-center">
                <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor"
                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z">
                  </path>
                </svg>
                Consumiendo...
              </div>
              <!-- Estado normal: muestra icono de check y texto "Consumir Insumo" -->
              <div v-else class="flex items-center">
                <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
                Consumir Insumo
              </div>
            </button>
          </div>
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
          <h3 class="text-sm font-medium text-green-800">Consumo Exitoso</h3>
          <div class="mt-2 text-sm text-green-700">
            El insumo ha sido consumido correctamente y marcado como usado.
            <div v-if="consumptionSuccess.status_change" class="mt-1 text-xs text-green-600">
              Estado cambiado de "{{ consumptionSuccess.status_change.from }}" a "{{ consumptionSuccess.status_change.to
              }}"
            </div>
          </div>
          <div class="mt-4 flex flex-wrap gap-3">
            <button @click="resetForm" class="btn-primary text-sm">
              Consumir Otro Insumo
            </button>
            <router-link to="/qr" class="btn-secondary text-sm">
              <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
              Volver al Scanner
            </router-link>
            <router-link v-if="consumptionSuccess.qr_code" :to="`/qr/${consumptionSuccess.qr_code}/traceability`"
              class="btn-secondary text-sm">
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
import { useRouter } from 'vue-router'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import qrService from '@/services/qr/qrService'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

// Estado reactivo
const qrInput = ref('')
const scanning = ref(false)
const error = ref(null)
const scannedProduct = ref(null)
const consuming = ref(false)
const consumptionSuccess = ref(null)
const selectedConsumptionPurpose = ref('')

// Datos de configuración (simplificado)

// Formulario de consumo
const consumptionForm = ref({
  notes: ''
})

// Usuario actual
const currentUser = computed(() => authStore.user)

// Propósitos de consumo (sin transferir, ya que eso tiene su propia vista)
const consumptionPurposes = [
  {
    value: 'routine',
    label: 'Rutina',
    description: 'Consumo programado normal',
    icon: 'M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2',
    iconClass: 'text-blue-600'
  },
  {
    value: 'emergency',
    label: 'Emergencia',
    description: 'Consumo de urgencia',
    icon: 'M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z',
    iconClass: 'text-red-600'
  },
  {
    value: 'maintenance',
    label: 'Mantenimiento',
    description: 'Uso en mantenimiento',
    icon: 'M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z M15 12a3 3 0 11-6 0 3 3 0 016 0z',
    iconClass: 'text-green-600'
  }
]

// Auto-completar RUT del usuario actual y procesar QR de la URL
onMounted(() => {
  // Verificar si hay un QR en los parámetros de la URL
  const urlParams = new URLSearchParams(window.location.search)
  const qrFromUrl = urlParams.get('qr')
  if (qrFromUrl) {
    qrInput.value = qrFromUrl
    // Escanear automáticamente el QR de la URL
    scanQR()
  }
})

// Escanear QR
const scanQR = async () => {
  if (!qrInput.value.trim()) return

  scanning.value = true
  error.value = null
  scannedProduct.value = null
  consumptionSuccess.value = null

  try {
    const result = await qrService.scanQRCode(qrInput.value.trim(), {
      scan_purpose: 'consumption_check',
      scan_source: 'consumption_view'
    })

    console.log('Resultado del escaneo QR:', result) // Debug

    if (result) {
      scannedProduct.value = result

      // Debug: Mostrar información del producto escaneado
      console.log('Información del producto:', {
        type: result.type,
        is_consumed: result.is_consumed,
        status: result.status,
        current_status: result.current_status,
        supply_info: result.supply_info,
        can_consume: canConsume(result)
      })

      // Verificar si es un insumo individual
      if (result.type !== 'medical_supply') {
        error.value = 'Solo se pueden consumir insumos individuales. Este código QR corresponde a un lote.'
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
}

// Verificar si se puede consumir
const canConsume = (product) => {
  if (!product) {
    console.log('canConsume: No hay producto')
    return false
  }

  // No se pueden consumir insumos ya consumidos
  if (product.is_consumed) {
    console.log('canConsume: Insumo ya consumido')
    return false
  }

  // Solo se pueden consumir insumos con estado "recepcionado"
  const status = product.supply_info?.status ||
    product.status ||
    product.current_status ||
    (product.supply_info?.Status)

  console.log('canConsume: Estado encontrado:', status)
  console.log('canConsume: Estructura supply_info:', product.supply_info)

  const canConsumeResult = status === 'recepcionado'
  console.log('canConsume: Resultado final:', canConsumeResult)

  return canConsumeResult
}

// Verificar si se debe mostrar opción de transferir
const shouldShowTransferOption = (product) => {
  if (!product) return false

  // Mostrar opción de transferir si tiene estado "disponible"
  const status = product.supply_info?.status || product.status
  return status === 'disponible' && !product.is_consumed
}

// Obtener mensaje de error para consumo
const getConsumptionErrorMessage = (product) => {
  if (!product) return ''

  if (product.is_consumed) {
    return 'Este insumo ya ha sido consumido anteriormente.'
  }

  const status = product.supply_info?.status || product.status
  if (status === 'disponible') {
    return 'Este insumo está disponible para ser transferido.'
  }

  if (status !== 'recepcionado') {
    return `El insumo tiene estado "${status}" y no está listo para consumo. Solo se pueden consumir insumos con estado "recepcionado".`
  }

  return 'Este insumo no se puede consumir.'
}

// Validar formulario de consumo
const validateConsumptionForm = () => {
  // Solo requiere que el usuario tenga RUT - las notas son opcionales
  if (!currentUser.value?.rut) return false
  return true
}

// Consumir producto
const consumeProduct = async () => {
  if (!scannedProduct.value || !validateConsumptionForm()) return

  consuming.value = true
  error.value = null

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
      'recepcionado'

    const consumptionData = {
      qr_code: scannedProduct.value.qr_code,
      user_rut: userRUT,
      user_name: currentUser.value?.name || 'Usuario',
      destination_type: 'pavilion',
      destination_id: pavilionId,
      medical_center_id: currentUser.value?.medical_center_id || 1,
      consumption_purpose: selectedConsumptionPurpose.value || 'rutina', // Por defecto 'rutina' si no se selecciona
      notes: consumptionForm.value.notes,
      consumption_timestamp: new Date().toISOString(),
      consumption_context: {
        scan_source: 'consumption_view',
        user_agent: navigator.userAgent
      }
    }

    const result = await qrService.consumeSupply(consumptionData)

    if (result.success) {
      consumptionSuccess.value = {
        ...result.data,
        qr_code: scannedProduct.value.qr_code,
        batch_id: scannedProduct.value.supply_info?.batch?.id,
        status_change: {
          from: getStatusLabel({ supply_info: { status: previousStatus } }),
          to: 'Consumido'
        },
        traceability_info: result.traceability_info,
        batch_history_updated: true
      }

      // Limpiar producto escaneado pero mantener el mensaje de éxito visible
      scannedProduct.value = null
      qrInput.value = ''
      selectedConsumptionPurpose.value = ''
    } else {
      error.value = result.error || 'Error al consumir el insumo'
    }
  } catch (err) {
    error.value = err.message || 'Error al consumir el insumo'
    console.error('Error al consumir:', err)
  } finally {
    consuming.value = false
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
    case 'recepcionado': return 'bg-blue-100 text-blue-800'
    default: return 'bg-gray-100 text-gray-800'
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
  consumptionSuccess.value = null
  qrInput.value = ''
  error.value = null
  selectedConsumptionPurpose.value = ''
  consumptionForm.value = {
    notes: ''
  }
}
</script>

<style scoped>
.form-input {
  @apply block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500;
}

.form-select {
  @apply block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500;
}

.form-textarea {
  @apply block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500;
}

/* Usar clases de botones de style.css global */
</style>