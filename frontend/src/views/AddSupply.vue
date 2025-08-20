<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Agregar Insumo Médico</h1>
        <p class="text-gray-600 mt-1">Registra un nuevo insumo médico con código QR</p>
      </div>
      <router-link to="/inventory" class="btn-secondary">
        <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16l-4-4m0 0l4-4m-4 4h18" />
        </svg>
        Volver al Inventario
      </router-link>
    </div>

    <div class="grid lg:grid-cols-2 gap-8">
      <!-- Formulario para crear insumo -->
      <div class="card">
        <div class="card-header">
          <h3 class="card-title text-blue-700">
            <svg class="h-6 w-6 inline mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
            Información del Insumo
          </h3>
        </div>

        <form @submit.prevent="createSupply" class="space-y-6">
          <!-- Código del insumo -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Código del Insumo *
            </label>
            <input
              type="number"
              required
              class="form-input w-full"
              v-model="supplyForm.code"
              placeholder="Ej: 1001"
            />
            <p class="text-xs text-gray-500 mt-1">Código numérico único para identificar el insumo</p>
          </div>

          <!-- Nombre del insumo -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Nombre del Insumo *
            </label>
            <input
              type="text"
              required
              class="form-input w-full"
              v-model="supplyForm.name"
              placeholder="Ej: Guantes de látex"
            />
          </div>

          <!-- Código del proveedor -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Código del Proveedor *
            </label>
            <input
              type="number"
              required
              class="form-input w-full"
              v-model="supplyForm.codeSupplier"
              placeholder="Ej: 5001"
            />
          </div>

          <!-- Información del lote -->
          <div class="border-t pt-4">
            <h4 class="text-md font-medium text-gray-900 mb-4">Información del Lote</h4>
            
            <div class="grid md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Fecha de Vencimiento *
                </label>
                <input
                  type="date"
                  required
                  class="form-input w-full"
                  v-model="batchForm.expirationDate"
                />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Cantidad del Lote *
                </label>
                <input
                  type="number"
                  required
                  min="1"
                  class="form-input w-full"
                  v-model="batchForm.amount"
                  placeholder="Ej: 100"
                />
              </div>
            </div>

            <div class="mt-4">
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Proveedor *
              </label>
              <input
                type="text"
                required
                class="form-input w-full"
                v-model="batchForm.supplier"
                placeholder="Ej: Proveedor Médico S.A."
              />
            </div>

            <div class="mt-4">
              <label class="block text-sm font-medium text-gray-700 mb-2">
                ID de Bodega *
              </label>
              <select v-model="batchForm.storeId" required class="form-select w-full">
                <option value="">Selecciona una bodega</option>
                <option value="1">Bodega Principal</option>
                <option value="2">Bodega Secundaria</option>
              </select>
            </div>
          </div>

          <!-- Botones -->
          <div class="flex space-x-4 pt-4">
            <button
              type="submit"
              :disabled="creating"
              class="btn-primary flex-1"
            >
              <svg v-if="!creating" class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
              <div v-else class="h-4 w-4 mr-2 spinner"></div>
              {{ creating ? 'Creando...' : 'Crear Insumo con QR' }}
            </button>
            
            <button
              type="button"
              @click="resetForm"
              class="btn-secondary px-6"
            >
              Limpiar
            </button>
          </div>
        </form>
      </div>

      <!-- Panel de QR generado -->
      <div class="space-y-6">
        <!-- QR Code generado -->
        <div v-if="generatedSupply" class="card">
          <div class="card-header">
            <h3 class="card-title text-green-700">
              <svg class="h-6 w-6 inline mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              Insumo Creado Exitosamente
            </h3>
          </div>

          <div class="text-center space-y-4">
            <!-- QR Code -->
            <div class="bg-white p-6 rounded-lg border-2 border-green-200 inline-block">
              <img 
                :src="`data:image/png;base64,${qrImageData}`"
                :alt="`QR del insumo: ${generatedSupply.qr_code}`"
                class="h-48 w-48 mx-auto"
              />
            </div>

            <!-- Información del insumo creado -->
            <div class="bg-green-50 rounded-lg p-4">
              <h4 class="font-medium text-green-900 mb-3">Detalles del Insumo</h4>
              <dl class="text-sm space-y-2">
                <div class="flex justify-between">
                  <dt class="text-green-700">ID:</dt>
                  <dd class="font-medium">#{{ generatedSupply.id }}</dd>
                </div>
                <div class="flex justify-between">
                  <dt class="text-green-700">Código:</dt>
                  <dd class="font-medium">{{ generatedSupply.code }}</dd>
                </div>
                <div class="flex justify-between">
                  <dt class="text-green-700">Código QR:</dt>
                  <dd class="font-mono text-xs">{{ generatedSupply.qr_code }}</dd>
                </div>
              </dl>
            </div>

            <!-- Acciones para el QR -->
            <div class="grid grid-cols-2 gap-3">
              <button @click="downloadQR" class="btn-secondary text-sm">
                <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
                Descargar QR
              </button>
              
              <button @click="testQR" class="btn-secondary text-sm">
                <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
                Probar Escáner
              </button>
            </div>

            <button @click="createAnother" class="btn-success w-full">
              <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
              Crear Otro Insumo
            </button>
          </div>
        </div>

        <!-- Instrucciones -->
        <div v-else class="card">
          <div class="card-header">
            <h3 class="card-title text-gray-700">
              <svg class="h-6 w-6 inline mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              Instrucciones
            </h3>
          </div>

          <div class="space-y-4 text-sm text-gray-600">
            <div class="flex items-start space-x-3">
              <div class="bg-blue-100 text-blue-600 rounded-full p-1 mt-0.5">
                <span class="text-xs font-bold">1</span>
              </div>
              <div>
                <p class="font-medium text-gray-900">Completa la información</p>
                <p>Llena todos los campos obligatorios del formulario</p>
              </div>
            </div>

            <div class="flex items-start space-x-3">
              <div class="bg-blue-100 text-blue-600 rounded-full p-1 mt-0.5">
                <span class="text-xs font-bold">2</span>
              </div>
              <div>
                <p class="font-medium text-gray-900">Genera el insumo</p>
                <p>Al crear el insumo, se generará automáticamente un código QR único</p>
              </div>
            </div>

            <div class="flex items-start space-x-3">
              <div class="bg-blue-100 text-blue-600 rounded-full p-1 mt-0.5">
                <span class="text-xs font-bold">3</span>
              </div>
              <div>
                <p class="font-medium text-gray-900">Descarga el QR</p>
                <p>Puedes descargar la imagen del QR para imprimirla o usarla</p>
              </div>
            </div>

            <div class="flex items-start space-x-3">
              <div class="bg-blue-100 text-blue-600 rounded-full p-1 mt-0.5">
                <span class="text-xs font-bold">4</span>
              </div>
              <div>
                <p class="font-medium text-gray-900">Prueba el escáner</p>
                <p>Verifica que el QR funcione correctamente en el escáner</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Mensaje de error -->
    <div v-if="error" class="bg-red-50 border border-red-200 rounded-md p-4">
      <div class="flex">
        <div class="flex-shrink-0">
          <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <div class="ml-3">
          <h3 class="text-sm font-medium text-red-800">Error al crear insumo</h3>
          <div class="mt-2 text-sm text-red-700">{{ error }}</div>
          <div class="mt-4">
            <button @click="clearError" class="btn-secondary text-sm">
              Continuar
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import inventoryService from '@/services/inventoryService'
import qrService from '@/services/qrService'

const router = useRouter()

// Estado reactivo
const creating = ref(false)
const error = ref(null)
const generatedSupply = ref(null)
const qrImageData = ref(null)

// Formularios
const supplyForm = ref({
  code: '',
  name: '',
  codeSupplier: ''
})

const batchForm = ref({
  expirationDate: '',
  amount: '',
  supplier: '',
  storeId: ''
})

// Método principal para crear insumo
const createSupply = async () => {
  creating.value = true
  error.value = null
  
  try {
    // Preparar datos completos para crear el insumo
    const completeSupplyData = {
      batch: {
        expiration_date: batchForm.value.expirationDate,
        amount: parseInt(batchForm.value.amount),
        supplier: batchForm.value.supplier,
        store_id: parseInt(batchForm.value.storeId)
      },
      supply_code: {
        code: parseInt(supplyForm.value.code),
        name: supplyForm.value.name,
        code_supplier: parseInt(supplyForm.value.codeSupplier)
      }
    }

    // Crear el insumo completo usando el servicio mejorado
    const result = await inventoryService.createCompleteSupply(completeSupplyData)
    
    // El resultado contiene: supply, batch, supply_code
    generatedSupply.value = result.supply
    
    // Generar la imagen del QR para mostrar
    const qrResult = await qrService.generateSupplyQR()
    if (qrResult.success && qrResult.data.image_data) {
      qrImageData.value = qrResult.data.image_data
      // Usar el QR code generado o el del resultado
      if (!generatedSupply.value.qr_code && qrResult.data.qr_code) {
        generatedSupply.value.qr_code = qrResult.data.qr_code
      }
    }

    // Si el backend no generó QR automáticamente, usar el del servicio QR
    if (!generatedSupply.value.qr_code && qrResult?.data?.qr_code) {
      generatedSupply.value.qr_code = qrResult.data.qr_code
    }

  } catch (err) {
    error.value = err.response?.data?.error || 'Error al crear el insumo médico'
    console.error('Error creating supply:', err)
  } finally {
    creating.value = false
  }
}

// Métodos de utilidad
const resetForm = () => {
  supplyForm.value = {
    code: '',
    name: '',
    codeSupplier: ''
  }
  batchForm.value = {
    expirationDate: '',
    amount: '',
    supplier: '',
    storeId: ''
  }
  generatedSupply.value = null
  qrImageData.value = null
  error.value = null
}

const createAnother = () => {
  resetForm()
}

const downloadQR = async () => {
  if (generatedSupply.value?.qr_code) {
    try {
      await qrService.downloadQRImage(
        generatedSupply.value.qr_code, 
        `insumo_${generatedSupply.value.code}_${generatedSupply.value.id}.png`
      )
    } catch (err) {
      console.error('Error downloading QR:', err)
    }
  }
}

const testQR = () => {
  if (generatedSupply.value?.qr_code) {
    router.push({
      name: 'QRScanner',
      query: { test: generatedSupply.value.qr_code }
    })
  }
}

const clearError = () => {
  error.value = null
}

// Inicializar fecha por defecto (30 días desde hoy)
const initDefaultDate = () => {
  const future = new Date()
  future.setDate(future.getDate() + 30)
  batchForm.value.expirationDate = future.toISOString().split('T')[0]
}

// Inicializar al montar el componente
initDefaultDate()
</script>

<style scoped>
/* Transiciones suaves */
.card {
  transition: all 0.2s ease-in-out;
}

.card:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
}

/* Colores específicos */
.text-blue-700 {
  color: #1d4ed8;
}

.text-green-700 {
  color: #15803d;
}

.bg-green-50 {
  background-color: #f0fdf4;
}

.border-green-200 {
  border-color: #bbf7d0;
}

/* Inputs focus */
.form-input:focus,
.form-select:focus {
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}
</style>