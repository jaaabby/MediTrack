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

    <!-- Scanner Section -->
    <div class="bg-white shadow-sm rounded-lg border">
      <div class="px-4 py-5 sm:p-6">
        <h2 class="text-lg font-medium text-gray-900 mb-4">Escanear Código QR</h2>

        <!-- QR Input -->
        <div class="flex space-x-3 mb-4">
          <div class="flex-1">
            <input type="text" v-model="qrInput" placeholder="Ingrese o escanee el código QR del insumo"
              class="form-input w-full" @keyup.enter="scanQR" :disabled="scanning" />
          </div>
          <button @click="scanQR" :disabled="!qrInput.trim() || scanning" class="btn-primary">
            <div v-if="scanning" class="flex items-center">
              <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z">
                </path>
              </svg>
              Escaneando...
            </div>
            <div v-else class="flex items-center">
              <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v2a2 2 0 002 2zm0 0h2a2 2 0 012 2v2a2 2 0 01-2 2H5a2 2 0 01-2-2v-2a2 2 0 012-2z" />
              </svg>
              Escanear
            </div>
          </button>
        </div>

        <!-- Error Display -->
        <div v-if="error" class="mb-4 bg-red-50 border border-red-200 rounded-md p-4">
          <div class="flex">
            <div class="flex-shrink-0">
              <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div class="ml-3">
              <h3 class="text-sm font-medium text-red-800">Error</h3>
              <div class="mt-2 text-sm text-red-700">{{ error }}</div>
            </div>
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
          <!-- Propósito del Consumo -->
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

          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Usuario Responsable -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                RUT Usuario Responsable <span class="text-red-500">*</span>
              </label>
              <input type="text" v-model="consumptionForm.userRUT" placeholder="12.345.678-9" class="form-input w-full"
                required />
            </div>

            <!-- Tipo de Destino -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Tipo de Destino <span class="text-red-500">*</span>
              </label>
              <select v-model="consumptionForm.destinationType" class="form-select w-full" required>
                <option value="">Seleccionar tipo</option>
                <option value="pavilion">Pabellón</option>
                <option value="warehouse">Almacén</option>
              </select>
            </div>

            <!-- Centro Médico -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Centro Médico <span class="text-red-500">*</span>
              </label>
              <select v-model="consumptionForm.medicalCenterId" class="form-select w-full" required
                :disabled="loadingCenters">
                <option value="">
                  {{ loadingCenters ? 'Cargando centros médicos...' : 'Seleccionar centro médico' }}
                </option>
                <option v-for="center in medicalCenters" :key="center.id" :value="center.id">
                  {{ center.name }}
                </option>
              </select>
            </div>

            <!-- Pabellón (cuando tipo es pabellón) -->
            <div v-if="consumptionForm.destinationType === 'pavilion'">
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Pabellón <span class="text-red-500">*</span>
              </label>
              <select v-model="consumptionForm.destinationID" class="form-select w-full" required
                :disabled="!consumptionForm.medicalCenterId || loadingPavilions">
                <option value="">
                  {{ !consumptionForm.medicalCenterId
                    ? 'Primero seleccione un centro médico'
                    : (loadingPavilions ? 'Cargando pabellones...' : 'Seleccionar pabellón')
                  }}
                </option>
                <option v-for="pavilion in filteredPavilions" :key="pavilion.id" :value="pavilion.id">
                  {{ pavilion.name }}
                </option>
              </select>
            </div>

            <!-- ID de Destino manual (cuando tipo es almacén u otro) -->
            <div v-if="consumptionForm.destinationType && consumptionForm.destinationType !== 'pavilion'">
              <label class="block text-sm font-medium text-gray-700 mb-2">
                {{ consumptionForm.destinationType === 'warehouse' ? 'ID Almacén' : 'ID Destino' }}
                <span class="text-red-500">*</span>
              </label>
              <input type="number" v-model="consumptionForm.destinationID"
                :placeholder="consumptionForm.destinationType === 'warehouse' ? 'Ej: 1' : 'Ej: 1'"
                class="form-input w-full" required />
            </div>
          </div>
        </div>

        <!-- Notas -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Notas Adicionales
          </label>
          <textarea v-model="consumptionForm.notes" rows="3" placeholder="Detalles del consumo, observaciones, etc."
            class="form-textarea w-full"></textarea>
        </div>

        <!-- Consumption Button -->
        <div class="flex justify-end space-x-3">
          <button @click="clearScannedProduct" class="btn-secondary">
            Cancelar
          </button>
          <button @click="consumeProduct" :disabled="!validateConsumptionForm() || consuming" class="btn-primary">
            <div v-if="consuming" class="flex items-center">
              <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z">
                </path>
              </svg>
              Consumiendo...
            </div>
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
          El insumo ha sido consumido correctamente. El historial del lote y la trazabilidad han sido actualizados.
        </div>
        <div class="mt-4 flex space-x-3">
          <button @click="resetForm" class="btn-primary text-sm">
            Consumir Otro Insumo
          </button>
          <router-link :to="`/qr/${consumptionSuccess.qr_code}/traceability`" class="btn-secondary text-sm">
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
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import qrService from '@/services/qrService'
import medicalCenterService from '@/services/medicalCenterService'
import pavilionService from '@/services/pavilionService'
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

// Datos de configuración
const medicalCenters = ref([])
const pavilions = ref([])
const filteredPavilions = ref([])
const loadingCenters = ref(false)
const loadingPavilions = ref(false)

// Formulario de consumo
const consumptionForm = ref({
  userRUT: '',
  destinationType: 'pavilion', // Por defecto pavilion
  destinationID: '',
  medicalCenterId: '',
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
  if (currentUser.value?.rut) {
    consumptionForm.value.userRUT = currentUser.value.rut
  }
  loadMedicalCenters()
  loadPavilions()

  // Verificar si hay un QR en los parámetros de la URL
  const urlParams = new URLSearchParams(window.location.search)
  const qrFromUrl = urlParams.get('qr')
  if (qrFromUrl) {
    qrInput.value = qrFromUrl
    // Escanear automáticamente el QR de la URL
    scanQR()
  }
})

// Cargar centros médicos desde la base de datos
const loadMedicalCenters = async () => {
  try {
    loadingCenters.value = true
    console.log('🏥 Cargando centros médicos desde BD...')
    
    const response = await medicalCenterService.getAll()
    console.log('🏥 Respuesta completa del servicio:', response)
    console.log('🏥 Datos de centros médicos:', response.data)
    
    if (response.data && Array.isArray(response.data) && response.data.length > 0) {
      medicalCenters.value = response.data
      console.log('✅ Centros médicos cargados desde BD:', medicalCenters.value)
    } else {
      console.warn('⚠️ No se encontraron centros médicos en la BD, usando fallback')
      throw new Error('No data found')
    }
  } catch (error) {
    console.error('❌ Error al cargar centros médicos desde BD:', error)
    console.log('🔄 Usando datos fallback para centros médicos')
    // Fallback con datos por defecto
    medicalCenters.value = [
      { id: 1, name: 'Centro Médico Principal' },
      { id: 2, name: 'Clínica Norte' }
    ]
  } finally {
    loadingCenters.value = false
  }
}

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

// Watcher para filtrar pabellones cuando se selecciona un centro médico
watch(() => consumptionForm.value.medicalCenterId, (newMedicalCenterId) => {
  console.log('🔄 Centro médico seleccionado:', newMedicalCenterId)
  console.log('🔄 Todos los pabellones disponibles:', pavilions.value)
  
  if (newMedicalCenterId) {
    const centerId = parseInt(newMedicalCenterId)
    filteredPavilions.value = pavilions.value.filter(p => p.medical_center_id === centerId)
    console.log('✅ Pabellones filtrados para centro', centerId, ':', filteredPavilions.value)
    
    if (filteredPavilions.value.length === 0) {
      console.warn('⚠️ No se encontraron pabellones para el centro médico', centerId)
    }
    
    // Limpiar selección de pabellón al cambiar de centro médico
    consumptionForm.value.destinationID = ''
  } else {
    filteredPavilions.value = []
    console.log('🔄 Centro médico deseleccionado, limpiando lista de pabellones')
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
    return 'Este insumo tiene estado "disponible" y solo puede ser transferido, no consumido. Use la vista de transferencia.'
  }

  if (status !== 'recepcionado') {
    return `El insumo tiene estado "${status}" y no está listo para consumo. Solo se pueden consumir insumos con estado "recepcionado".`
  }

  return 'Este insumo no se puede consumir.'
}

// Validar formulario de consumo
const validateConsumptionForm = () => {
  if (!selectedConsumptionPurpose.value) return false
  if (!consumptionForm.value.userRUT.trim()) return false
  if (!consumptionForm.value.destinationType) return false
  if (!consumptionForm.value.destinationID) return false
  if (!consumptionForm.value.medicalCenterId) return false

  const destinationIdNum = parseInt(consumptionForm.value.destinationID)
  if (isNaN(destinationIdNum) || destinationIdNum < 1) return false

  const medicalCenterIdNum = parseInt(consumptionForm.value.medicalCenterId)
  if (isNaN(medicalCenterIdNum) || medicalCenterIdNum < 1) return false

  return true
}

// Consumir producto
const consumeProduct = async () => {
  if (!scannedProduct.value || !validateConsumptionForm()) return

  consuming.value = true
  error.value = null

  try {
    const consumptionData = {
      qr_code: scannedProduct.value.qr_code,
      user_rut: consumptionForm.value.userRUT,
      user_name: currentUser.value?.name || 'Usuario',
      destination_type: consumptionForm.value.destinationType,
      destination_id: parseInt(consumptionForm.value.destinationID),
      medical_center_id: parseInt(consumptionForm.value.medicalCenterId),
      consumption_purpose: selectedConsumptionPurpose.value,
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
        traceability_info: result.traceability_info,
        batch_history_updated: true
      }

      // Limpiar formulario
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
    userRUT: currentUser.value?.rut || '',
    destinationType: 'pavilion', // Por defecto pavilion
    destinationID: '',
    medicalCenterId: '',
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

.btn-primary {
  @apply inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed;
}

.btn-secondary {
  @apply inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md shadow-sm text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500;
}
</style>