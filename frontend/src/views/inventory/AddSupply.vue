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
          
          <div class="grid md:grid-cols-2 gap-6">
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
            
            <div>
              <label for="critical-stock" class="block text-sm font-medium text-gray-700 mb-2">
                Stock Crítico <span class="text-red-500">*</span>
              </label>
              <input
                id="critical-stock"
                v-model.number="supplyForm.criticalStock"
                type="number"
                min="1"
                placeholder="1"
                class="form-input"
                required
              />
              <p class="text-sm text-gray-500 mt-1">
                Nivel mínimo de stock para generar alertas. Para insumos específicos, generalmente se usa 1.
              </p>
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
            
            <div class="relative">
              <label for="store-search" class="block text-sm font-medium text-gray-700 mb-2">
                Almacén <span class="text-red-500">*</span>
              </label>
              <div class="relative">
                <input
                  id="store-search"
                  v-model="storeSearch"
                  @input="onStoreSearch"
                  @focus="showStoreOptions = true"
                  @blur="hideStoreOptions"
                  type="text"
                  :placeholder="loadingStores ? 'Cargando almacenes...' : 'Escribir o seleccionar almacén...'"
                  class="form-input w-full pr-10"
                  :disabled="loadingStores"
                  autocomplete="off"
                  required
                />
                <div class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
                  <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                  </svg>
                </div>
              </div>
              
              <!-- Dropdown de opciones -->
              <div
                v-show="showStoreOptions && filteredStores.length > 0"
                class="absolute z-10 mt-1 w-full bg-white shadow-lg max-h-60 rounded-md py-1 text-base ring-1 ring-black ring-opacity-5 overflow-auto focus:outline-none sm:text-sm border border-gray-200"
              >
                <div
                  v-for="store in filteredStores"
                  :key="store.id"
                  @mousedown="selectStore(store)"
                  class="cursor-pointer select-none relative py-2 pl-3 pr-9 hover:bg-blue-50"
                >
                  <div class="flex items-center">
                    <span class="font-medium text-gray-900 block truncate">
                      {{ store.name }}
                    </span>
                    <span v-if="store.type" class="ml-2 text-sm text-gray-500">
                      ({{ store.type }})
                    </span>
                  </div>
                </div>
              </div>
              
              <!-- Mensaje cuando no hay resultados -->
              <div
                v-show="showStoreOptions && storeSearch && filteredStores.length === 0 && !loadingStores"
                class="absolute z-10 mt-1 w-full bg-white shadow-lg rounded-md py-2 px-3 text-sm text-gray-500 border border-gray-200"
              >
                No se encontraron almacenes que coincidan con "{{ storeSearch }}"
              </div>
              
              <!-- Mensaje cuando no hay almacenes -->
              <p v-if="stores.length === 0 && !loadingStores" class="text-sm text-red-500 mt-1">
                No hay almacenes disponibles
              </p>
            </div>
            
            <div class="md:col-span-2">
              <label for="expiration-alert-days" class="block text-sm font-medium text-gray-700 mb-2">
                Días de Alerta de Vencimiento <span class="text-red-500">*</span>
              </label>
              <input
                id="expiration-alert-days"
                v-model.number="batchForm.expirationAlertDays"
                type="number"
                min="1"
                max="365"
                placeholder="90"
                class="form-input"
                required
              />
              <p class="text-sm text-gray-500 mt-1">
                Días de anticipación para alertas de vencimiento (mínimo 90 días recomendado)
              </p>
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
                  class="btn-primary text-xs px-2 py-1"
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
  </div>
</template>

<script setup>

import { ref, onMounted, computed } from 'vue'
import qrService from '@/services/qr/qrService'
import inventoryService from '@/services/inventory/inventoryService'
import storeService from '@/services/inventory/storeService'
import { useNotification } from '@/composables/useNotification'

const { success: showSuccess, error: showError } = useNotification()

// Validación de formularios
const validateForm = () => {
  // Validar supplyForm
  if (!supplyForm.value.code || isNaN(parseInt(supplyForm.value.code))) {
    showError('El código del insumo es obligatorio y debe ser numérico.')
    return false
  }
  if (!supplyForm.value.name || supplyForm.value.name.trim() === '') {
    showError('El nombre del insumo es obligatorio.')
    return false
  }
  if (!supplyForm.value.codeSupplier || isNaN(parseInt(supplyForm.value.codeSupplier))) {
    showError('El código de proveedor es obligatorio y debe ser numérico.')
    return false
  }
  // Validar batchForm
  if (!batchForm.value.expirationDate) {
    showError('La fecha de vencimiento es obligatoria.')
    return false
  }
  if (!batchForm.value.amount || isNaN(parseInt(batchForm.value.amount)) || parseInt(batchForm.value.amount) < 1) {
    showError('La cantidad debe ser un número mayor a 0.')
    return false
  }
  if (!batchForm.value.supplier || batchForm.value.supplier.trim() === '') {
    showError('El proveedor es obligatorio.')
    return false
  }
  if (!batchForm.value.expirationAlertDays || isNaN(parseInt(batchForm.value.expirationAlertDays)) || parseInt(batchForm.value.expirationAlertDays) < 1) {
    showError('Los días de alerta deben ser un número mayor a 0.')
    return false
  }
  if (!batchForm.value.storeId || batchForm.value.storeId === '') {
    showError('Debe seleccionar un almacén.')
    return false
  }
  
  // Verificar que el almacén seleccionado exista en la lista
  const selectedStore = stores.value.find(s => s.id === parseInt(batchForm.value.storeId))
  if (!selectedStore) {
    showError('Debe seleccionar un almacén válido.')
    return false
  }
  
  return true
}

// Estado reactivo
const creating = ref(false)
const error = ref(null)
const generatedBatch = ref(null)
const generatedSupplies = ref(null)
const batchQRImage = ref(null)
const showAllQRs = ref(false)
const stores = ref([])
const loadingStores = ref(false)
const storeSearch = ref('')
const showStoreOptions = ref(false)

// Formularios
const supplyForm = ref({
  code: '',
  name: '',
  codeSupplier: '',
  criticalStock: 1
})

const batchForm = ref({
  expirationDate: '',
  amount: '',
  supplier: '',
  storeId: '', // Ahora será el ID del almacén seleccionado
  expirationAlertDays: 90 // Valor por defecto: 90 días (3 meses)
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
        store_id: parseInt(batchForm.value.storeId), // Ya viene como número del select
        expiration_alert_days: parseInt(batchForm.value.expirationAlertDays)
      },
      supply_code: {
        code: parseInt(supplyForm.value.code),
        name: supplyForm.value.name,
        code_supplier: parseInt(supplyForm.value.codeSupplier),
        critical_stock: parseInt(supplyForm.value.criticalStock) || 1
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
      
      showSuccess(`Lote creado exitosamente con ${generatedSupplies.value.length} insumos individuales`)
      await loadBatchQRImage()
    } else {
      showError(result.error || 'Error desconocido al crear el lote')
    }

  } catch (err) {
    console.error('Error creating supply:', err)
    showError(err.message || 'Error de conexión al crear el lote')
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
    const date = new Date(dateString)
    const day = String(date.getDate()).padStart(2, '0')
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const year = date.getFullYear()
    return `${day}/${month}/${year}`
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
  
  supplyForm.value = {
    code: '',
    name: '',
    codeSupplier: '',
    criticalStock: 1
  }
  
  batchForm.value = {
    expirationDate: '',
    amount: '',
    supplier: '',
    storeId: '',
    expirationAlertDays: 90 // Resetear al valor por defecto
  }
  storeSearch.value = ''
  showStoreOptions.value = false
}

// Computed para almacenes filtrados
const filteredStores = computed(() => {
  if (!storeSearch.value.trim()) {
    // Si no hay búsqueda, mostrar todos (limitado a 10 para no ser invasivo)
    return stores.value.slice(0, 10)
  }
  const search = storeSearch.value.toLowerCase().trim()
  return stores.value.filter(store => 
    store.name?.toLowerCase().includes(search) ||
    store.type?.toLowerCase().includes(search) ||
    store.id?.toString().includes(search)
  ).slice(0, 10) // Limitar a 10 resultados para no ser invasivo
})

// Cargar almacenes al montar el componente
const loadStores = async () => {
  loadingStores.value = true
  try {
    const storesList = await storeService.getAllStores()
    stores.value = storesList || []
  } catch (err) {
    console.error('Error al cargar almacenes:', err)
    showError('Error al cargar la lista de almacenes. Por favor, recarga la página.')
  } finally {
    loadingStores.value = false
  }
}

// Funciones para manejar la búsqueda y selección de almacenes
const onStoreSearch = () => {
  showStoreOptions.value = true
  // Si el texto coincide exactamente con un almacén, seleccionarlo automáticamente
  const exactMatch = stores.value.find(store => 
    store.name.toLowerCase() === storeSearch.value.toLowerCase() ||
    `${store.name} ${store.type ? `(${store.type})` : ''}`.toLowerCase() === storeSearch.value.toLowerCase()
  )
  if (exactMatch) {
    batchForm.value.storeId = exactMatch.id.toString()
  } else if (storeSearch.value.trim() === '') {
    batchForm.value.storeId = ''
  }
}

const selectStore = (store) => {
  batchForm.value.storeId = store.id.toString()
  storeSearch.value = store.type ? `${store.name} (${store.type})` : store.name
  showStoreOptions.value = false
}

const hideStoreOptions = () => {
  // Delay para permitir que el click en una opción se registre antes de ocultar
  setTimeout(() => {
    showStoreOptions.value = false
    // Si hay un almacén seleccionado, mostrar su nombre completo
    if (batchForm.value.storeId) {
      const selectedStore = stores.value.find(s => s.id === parseInt(batchForm.value.storeId))
      if (selectedStore) {
        storeSearch.value = selectedStore.type ? `${selectedStore.name} (${selectedStore.type})` : selectedStore.name
      }
    } else if (!storeSearch.value.trim()) {
      storeSearch.value = ''
    }
  }, 200)
}

// Cargar almacenes al montar el componente
onMounted(() => {
  loadStores()
})
</script>

<style scoped>
.form-input {
  @apply block w-full rounded-lg border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm px-3 py-2 border;
}

.form-input:disabled {
  @apply bg-gray-100 cursor-not-allowed;
}

/* Usar clases de botones de style.css global */
</style>