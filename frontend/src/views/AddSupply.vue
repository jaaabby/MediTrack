<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <h2 class="text-2xl font-semibold text-gray-900">Agregar Nuevo Insumo Médico</h2>
      <p class="text-gray-600 mt-1">Crea un lote de productos con códigos QR únicos para cada unidad individual</p>
    </div>

    <!-- Formulario Principal -->
    <div v-if="!generatedBatch" class="bg-white rounded-lg shadow-sm border p-6">
      <form @submit.prevent="createSupply" class="space-y-8">
        
        <!-- Información del Insumo -->
        <div>
          <h3 class="text-lg font-medium text-gray-900 mb-4 flex items-center">
            <svg class="h-5 w-5 mr-2 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
            </svg>
            Información del Insumo
          </h3>
          
          <div class="grid md:grid-cols-3 gap-6">
            <div>
              <label for="supply-code" class="block text-sm font-medium text-gray-700 mb-2">
                Código del Insumo <span class="text-red-500">*</span>
              </label>
              <input
                id="supply-code"
                v-model="supplyForm.code"
                type="number"
                placeholder="123456"
                class="form-input"
                required
              />
            </div>
            
            <div>
              <label for="supply-name" class="block text-sm font-medium text-gray-700 mb-2">
                Nombre del Insumo <span class="text-red-500">*</span>
              </label>
              <input
                id="supply-name"
                v-model="supplyForm.name"
                type="text"
                placeholder="Ej: Jeringa 10ml"
                class="form-input"
                required
              />
            </div>
            
            <div>
              <label for="supplier-code" class="block text-sm font-medium text-gray-700 mb-2">
                Código de Proveedor <span class="text-red-500">*</span>
              </label>
              <input
                id="supplier-code"
                v-model="supplyForm.codeSupplier"
                type="number"
                placeholder="789"
                class="form-input"
                required
              />
            </div>
          </div>
        </div>

        <!-- Información del Lote -->
        <div>
          <h3 class="text-lg font-medium text-gray-900 mb-4 flex items-center">
            <svg class="h-5 w-5 mr-2 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
            </svg>
            Información del Lote
          </h3>
          
          <div class="grid md:grid-cols-2 gap-6">
            <div class="space-y-6">
              <div>
                <label for="expiration-date" class="block text-sm font-medium text-gray-700 mb-2">
                  Fecha de Vencimiento <span class="text-red-500">*</span>
                </label>
                <input
                  id="expiration-date"
                  v-model="batchForm.expirationDate"
                  type="date"
                  class="form-input"
                  required
                />
              </div>
              
              <div>
                <label for="supplier" class="block text-sm font-medium text-gray-700 mb-2">
                  Proveedor <span class="text-red-500">*</span>
                </label>
                <input
                  id="supplier"
                  v-model="batchForm.supplier"
                  type="text"
                  placeholder="Nombre del proveedor"
                  class="form-input"
                  required
                />
              </div>
            </div>
            
            <div class="space-y-6">
              <div>
                <label for="amount" class="block text-sm font-medium text-gray-700 mb-2">
                  Cantidad de Productos <span class="text-red-500">*</span>
                </label>
                <input
                  id="amount"
                  v-model="batchForm.amount"
                  type="number"
                  min="1"
                  max="1000"
                  placeholder="50"
                  class="form-input"
                  required
                />
                <p class="text-sm text-gray-500 mt-1">Se generará un QR único para cada producto individual</p>
              </div>
              
              <div>
                <label for="store-id" class="block text-sm font-medium text-gray-700 mb-2">
                  ID del Almacén <span class="text-red-500">*</span>
                </label>
                <input
                  id="store-id"
                  v-model="batchForm.storeId"
                  type="number"
                  placeholder="1"
                  class="form-input"
                  required
                />
              </div>
            </div>
          </div>
        </div>

        <!-- Submit Button -->
        <div class="flex justify-end">
          <button
            type="submit"
            :disabled="creating"
            class="btn-primary text-lg px-8 py-3"
          >
            <svg v-if="creating" class="animate-spin h-5 w-5 mr-3" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <svg v-else class="h-5 w-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
            {{ creating ? 'Creando Lote...' : 'Crear Lote con QR Codes' }}
          </button>
        </div>
      </form>
    </div>

    <!-- Resultado Exitoso -->
    <div v-if="generatedBatch && !error" class="space-y-6">
      
      <!-- Información del Lote Creado -->
      <div class="bg-green-50 border border-green-200 rounded-lg p-6">
        <div class="flex items-start">
          <div class="flex-shrink-0">
            <svg class="h-8 w-8 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div class="ml-4 flex-1">
            <h3 class="text-lg font-medium text-green-800">¡Lote Creado Exitosamente!</h3>
            <div class="mt-2 text-sm text-green-700">
              <p>Se ha creado el lote <strong>ID: {{ generatedBatch.id }}</strong> con <strong>{{ generatedSupplies?.length || batchForm.amount }} productos individuales</strong></p>
              <p class="mt-1">Cada producto tiene su propio código QR único para trazabilidad completa.</p>
            </div>
          </div>
        </div>
      </div>

      <!-- QR del Lote -->
      <div class="bg-white rounded-lg shadow-sm border p-6">
        <h3 class="text-lg font-medium text-gray-900 mb-4 flex items-center">
          <svg class="h-5 w-5 mr-2 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
        </svg>
        Código QR del Lote
        </h3>
        
        <div class="flex flex-col lg:flex-row lg:items-start lg:space-x-8">
          <!-- Imagen QR del Lote -->
          <div class="flex-shrink-0 text-center mb-4 lg:mb-0">
            <div class="bg-gray-50 rounded-lg p-6 border">
              <img 
                v-if="batchQRImage" 
                :src="batchQRImage" 
                :alt="`QR del Lote: ${generatedBatch.qr_code}`"
                class="w-48 h-48 mx-auto object-contain border rounded"
              />
              <div v-else class="w-48 h-48 mx-auto bg-gray-100 rounded flex items-center justify-center">
                <svg class="h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v14a2 2 0 002 2z" />
                </svg>
              </div>
              
              <div class="mt-4 space-y-2">
                <button 
                  @click="downloadBatchQR('normal')"
                  class="w-full btn-secondary text-sm"
                >
                  Descargar QR Lote
                </button>
                <button 
                  @click="downloadBatchQR('high')"
                  class="w-full btn-primary text-sm"
                >
                  Descargar HD
                </button>
              </div>
            </div>
          </div>
          
          <!-- Información del Lote -->
          <div class="flex-1">
            <div class="bg-gray-50 rounded-lg p-4 space-y-3">
              <div class="grid md:grid-cols-2 gap-4">
                <div>
                  <label class="text-sm font-medium text-gray-600">ID del Lote:</label>
                  <p class="text-gray-900 font-semibold">{{ generatedBatch.id }}</p>
                </div>
                <div>
                  <label class="text-sm font-medium text-gray-600">Código QR del Lote:</label>
                  <code class="block text-sm text-gray-800 font-mono bg-white px-2 py-1 rounded border">
                    {{ generatedBatch.qr_code }}
                  </code>
                </div>
                <div>
                  <label class="text-sm font-medium text-gray-600">Proveedor:</label>
                  <p class="text-gray-900">{{ generatedBatch.supplier }}</p>
                </div>
                <div>
                  <label class="text-sm font-medium text-gray-600">Fecha de Vencimiento:</label>
                  <p class="text-gray-900">{{ formatDate(generatedBatch.expiration_date) }}</p>
                </div>
                <div>
                  <label class="text-sm font-medium text-gray-600">Cantidad Total:</label>
                  <p class="text-gray-900 font-semibold">{{ generatedBatch.amount }} unidades</p>
                </div>
                <div>
                  <label class="text-sm font-medium text-gray-600">Almacén:</label>
                  <p class="text-gray-900">ID: {{ generatedBatch.store_id }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Lista de Productos Individuales -->
      <div v-if="generatedSupplies && generatedSupplies.length > 0" class="bg-white rounded-lg shadow-sm border p-6">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-lg font-medium text-gray-900 flex items-center">
            <svg class="h-5 w-5 mr-2 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
            </svg>
            Productos Individuales ({{ generatedSupplies.length }})
          </h3>
          
          <div class="flex space-x-2">
            <button @click="showAllQRs = !showAllQRs" class="btn-secondary text-sm">
              {{ showAllQRs ? 'Ocultar QRs' : 'Mostrar QRs' }}
            </button>
            <button @click="downloadAllSupplyQRs" class="btn-primary text-sm">
              <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              Descargar Todos
            </button>
          </div>
        </div>

        <!-- Grid de Productos -->
        <div class="grid sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
          <div 
            v-for="(supply, index) in generatedSupplies" 
            :key="supply.id"
            class="bg-gray-50 rounded-lg p-4 border hover:shadow-sm transition-shadow"
          >
            <!-- Imagen QR del producto individual -->
            <div v-if="showAllQRs" class="text-center mb-3">
              <img 
                :src="getSupplyQRImageUrl(supply.qr_code)" 
                :alt="`QR Producto ${index + 1}`"
                class="w-20 h-20 mx-auto object-contain border rounded bg-white"
                @error="handleSupplyQRError"
              />
            </div>
            
            <div class="space-y-2">
              <div class="flex justify-between items-center">
                <span class="text-sm font-medium text-gray-600">Producto:</span>
                <span class="text-sm font-semibold text-gray-900">#{{ index + 1 }}</span>
              </div>
              
              <div>
                <span class="text-xs text-gray-500">ID:</span>
                <span class="text-xs text-gray-700 font-mono ml-1">{{ supply.id }}</span>
              </div>
              
              <div>
                <span class="text-xs text-gray-500">QR:</span>
                <code class="block text-xs text-gray-800 font-mono bg-white px-2 py-1 rounded border mt-1 break-all">
                  {{ supply.qr_code }}
                </code>
              </div>
              
              <div class="flex space-x-1 pt-2">
                <button 
                  @click="downloadSupplyQR(supply.qr_code, 'normal')"
                  class="flex-1 text-xs bg-gray-600 hover:bg-gray-700 text-white px-2 py-1 rounded transition-colors"
                >
                  PNG
                </button>
                <button 
                  @click="downloadSupplyQR(supply.qr_code, 'high')"
                  class="flex-1 text-xs bg-blue-600 hover:bg-blue-700 text-white px-2 py-1 rounded transition-colors"
                >
                  HD
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Acciones Finales -->
      <div class="flex flex-col sm:flex-row gap-4 justify-center">
        <router-link to="/inventory" class="btn-primary text-center">
          <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
          </svg>
          Ver en Inventario
        </router-link>
        
        <button @click="createAnother" class="btn-secondary">
          <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          Crear Otro Lote
        </button>
        
        <router-link :to="`/qr/${generatedBatch.qr_code}`" class="btn-secondary text-center">
          <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h2M4 4h5l2 3h3l2-3h5v5M4 4v5m16-5v5" />
          </svg>
          Escanear QR del Lote
        </router-link>
      </div>
    </div>

    <!-- Mensaje de Error -->
    <div v-if="error" class="bg-red-50 border border-red-200 rounded-lg p-6">
      <div class="flex items-start">
        <div class="flex-shrink-0">
          <svg class="h-6 w-6 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <div class="ml-3 flex-1">
          <h3 class="text-lg font-medium text-red-800">Error al crear el lote</h3>
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
  </div>
</template>

<script setup>

import { ref } from 'vue'
import qrService from '@/services/qrService'
import inventoryService from '@/services/inventoryService'

// Validación de formularios
const validateForm = () => {
  // Validar supplyForm
  if (!supplyForm.value.code || isNaN(parseInt(supplyForm.value.code))) {
    error.value = 'El código del insumo es obligatorio y debe ser numérico.'
    return false
  }
  if (!supplyForm.value.name || supplyForm.value.name.trim() === '') {
    error.value = 'El nombre del insumo es obligatorio.'
    return false
  }
  if (!supplyForm.value.codeSupplier || isNaN(parseInt(supplyForm.value.codeSupplier))) {
    error.value = 'El código de proveedor es obligatorio y debe ser numérico.'
    return false
  }
  // Validar batchForm
  if (!batchForm.value.expirationDate) {
    error.value = 'La fecha de vencimiento es obligatoria.'
    return false
  }
  if (!batchForm.value.amount || isNaN(parseInt(batchForm.value.amount)) || parseInt(batchForm.value.amount) < 1) {
    error.value = 'La cantidad debe ser un número mayor a 0.'
    return false
  }
  if (!batchForm.value.supplier || batchForm.value.supplier.trim() === '') {
    error.value = 'El proveedor es obligatorio.'
    return false
  }
  if (!batchForm.value.storeId || isNaN(parseInt(batchForm.value.storeId))) {
    error.value = 'El ID de almacén es obligatorio y debe ser numérico.'
    return false
  }
  error.value = null
  return true
}

// Estado reactivo
const creating = ref(false)
const error = ref(null)
const generatedBatch = ref(null)
const generatedSupplies = ref(null)
const batchQRImage = ref(null)
const showAllQRs = ref(false)

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

// Métodos principales
const createSupply = async () => {
  if (!validateForm()) return
  
  creating.value = true
  error.value = null

  try {
    // Preparar datos para el método correcto
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

    console.log('Datos a enviar:', completeSupplyData)

    // USAR EL MÉTODO CORRECTO que crea múltiples insumos
    const result = await inventoryService.createBatchWithIndividualSupplies(completeSupplyData)

    console.log('Resultado completo:', result)

    // Verificar la estructura de la respuesta
    if (result.success && result.data) {
      generatedBatch.value = result.data.batch
      generatedSupplies.value = result.data.individual_supplies || []
      
      console.log(`✅ Lote creado exitosamente con ${generatedSupplies.value.length} insumos individuales`)
      
      error.value = null
      await loadBatchQRImage()
    } else {
      error.value = result.error || 'Error desconocido al crear el lote'
    }

  } catch (err) {
    console.error('Error creating supply:', err)
    error.value = err.message || 'Error de conexión al crear el lote'
  } finally {
    creating.value = false
  }
}

const loadBatchQRImage = async () => {
  if (!generatedBatch.value?.qr_code) return
  
  try {
    batchQRImage.value = qrService.getQRImageUrl(generatedBatch.value.qr_code)
  } catch (error) {
    console.error('Error loading batch QR image:', error)
  }
}

// Métodos de descarga
const downloadBatchQR = async (resolution = 'normal') => {
  if (!generatedBatch.value?.qr_code) return
  
  try {
    await qrService.downloadQRImage(generatedBatch.value.qr_code, resolution)
  } catch (error) {
    console.error('Error downloading batch QR:', error)
  }
}

const downloadSupplyQR = async (qrCode, resolution = 'normal') => {
  try {
    await qrService.downloadQRImage(qrCode, resolution)
  } catch (error) {
    console.error('Error downloading supply QR:', error)
  }
}

const downloadAllSupplyQRs = async () => {
  if (!generatedSupplies.value || generatedSupplies.value.length === 0) return
  
  try {
    // Descargar todos los QRs de productos individuales
    for (const supply of generatedSupplies.value) {
      await downloadSupplyQR(supply.qr_code, 'normal')
      // Pequeña pausa entre descargas para evitar saturar el navegador
      await new Promise(resolve => setTimeout(resolve, 100))
    }
  } catch (error) {
    console.error('Error downloading all supply QRs:', error)
  }
}

// Métodos auxiliares
const getSupplyQRImageUrl = (qrCode) => {
  return qrService.getQRImageUrl(qrCode)
}

const handleSupplyQRError = (event) => {
  // Manejar error de carga de imagen QR individual
  event.target.style.display = 'none'
}

const formatDate = (dateString) => {
  if (!dateString) return 'No disponible'
  try {
    return format(new Date(dateString), 'dd/MM/yyyy', { locale: es })
  } catch (error) {
    return dateString
  }
}

const createAnother = () => {
  // Reset form
  generatedBatch.value = null
  generatedSupplies.value = null
  batchQRImage.value = null
  showAllQRs.value = false
  error.value = null
  
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
}

const clearError = () => {
  error.value = null
}
</script>

<style scoped>
.form-input {
  @apply block w-full rounded-lg border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm;
}

.btn-primary {
  @apply bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-lg transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 inline-flex items-center justify-center;
}

.btn-secondary {
  @apply bg-gray-200 hover:bg-gray-300 text-gray-900 font-medium py-2 px-4 rounded-lg transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 inline-flex items-center justify-center;
}
</style>