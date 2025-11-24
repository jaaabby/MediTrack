<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h2 class="text-2xl font-semibold text-gray-900">Gestión de Códigos de Insumos</h2>
          <p class="text-gray-600 mt-1">Gestiona los códigos de insumos y sus niveles críticos de stock</p>
          <p v-if="!loading" class="text-sm text-gray-500 mt-1">Total: {{ sortedSupplyCodes.length }} códigos</p>
        </div>
        <button @click="openCreateModal" class="btn-primary flex items-center justify-center">
          <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Nuevo Código
        </button>
      </div>
    </div>

    <!-- Búsqueda -->
    <div class="card">
      <div class="flex flex-col sm:flex-row sm:items-end gap-4">
        <div class="flex-1">
          <label class="block text-sm font-medium text-gray-700 mb-2">Buscar código de insumo</label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            <input type="text" placeholder="Buscar por código o nombre..." 
              class="form-input pl-10 w-full" v-model="searchTerm" @input="handleSearch" />
          </div>
        </div>
        <button class="btn-secondary px-4 py-2 h-10" @click="clearSearch" :disabled="!searchTerm">
          Limpiar
        </button>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="card">
      <div class="flex justify-center items-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        <span class="ml-3 text-gray-600">Cargando códigos de insumos...</span>
      </div>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="card">
      <div class="bg-red-50 border border-red-200 rounded-md p-4">
        <div class="flex">
          <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-red-800">Error al cargar códigos de insumos</h3>
            <div class="mt-2 text-sm text-red-700">{{ error }}</div>
            <button @click="loadSupplyCodes" class="btn-secondary mt-4 text-sm">Reintentar</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Tabla de códigos -->
    <div v-else class="card overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th scope="col" @click="sortBy('code')"
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                <div class="flex items-center space-x-1">
                  <span>Código</span>
                  <span class="flex flex-col -space-y-1">
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'code' && sortOrder === 'asc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                    </svg>
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'code' && sortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z"/>
                    </svg>
                  </span>
                </div>
              </th>
              <th scope="col" @click="sortBy('name')"
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                <div class="flex items-center space-x-1">
                  <span>Nombre</span>
                  <span class="flex flex-col -space-y-1">
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'name' && sortOrder === 'asc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                    </svg>
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'name' && sortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z"/>
                    </svg>
                  </span>
                </div>
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Código Proveedor
              </th>
              <th scope="col" @click="sortBy('critical_stock')"
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                <div class="flex items-center space-x-1">
                  <span>Stock Crítico</span>
                  <span class="flex flex-col -space-y-1">
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'critical_stock' && sortOrder === 'asc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                    </svg>
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'critical_stock' && sortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z"/>
                    </svg>
                  </span>
                </div>
              </th>
              <th scope="col" class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                Acciones
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="supplyCode in paginatedSupplyCodes" :key="supplyCode.code" 
              class="hover:bg-gray-50 transition-colors duration-150">
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm font-medium text-gray-900">{{ supplyCode.code }}</div>
              </td>
              <td class="px-6 py-4">
                <div class="text-sm text-gray-900">{{ supplyCode.name }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm text-gray-900">{{ supplyCode.code_supplier }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm text-gray-900">
                  <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                    :class="getCriticalStockBadgeClass(supplyCode.critical_stock)">
                    {{ supplyCode.critical_stock }}
                  </span>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <div class="flex justify-end space-x-2">
                  <button @click="openEditModal(supplyCode)" 
                    class="btn-primary text-xs px-3 py-1.5"
                    title="Editar">
                    <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                  </button>
                  <button @click="confirmDelete(supplyCode)" 
                    class="btn-danger text-xs px-3 py-1.5"
                    title="Eliminar">
                    <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Paginación -->
    <div v-if="!loading && sortedSupplyCodes.length > 0" class="card">
      <div class="flex flex-col sm:flex-row items-center justify-between gap-4">
        <div class="text-sm text-gray-700 text-center sm:text-left">
          Mostrando {{ startIndex + 1 }} a {{ endIndex }} de {{ sortedSupplyCodes.length }} códigos
        </div>
        <div class="flex items-center gap-2">
          <button class="btn-secondary px-3 py-2 text-sm min-w-[70px]" :disabled="currentPage === 1"
            @click="currentPage--">
            <span class="hidden sm:inline">Anterior</span>
            <span class="sm:hidden">Ant.</span>
          </button>
          <span class="px-3 py-2 text-sm text-gray-700 bg-gray-100 rounded-md min-w-[90px] text-center">
            Página {{ currentPage }} de {{ totalPages }}
          </span>
          <button class="btn-secondary px-3 py-2 text-sm min-w-[70px]" :disabled="currentPage === totalPages"
            @click="currentPage++">
            <span class="hidden sm:inline">Siguiente</span>
            <span class="sm:hidden">Sig.</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Mensaje sin resultados -->
    <div v-if="!loading && supplyCodes.length === 0" class="card text-center py-12">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900">No hay códigos de insumos</h3>
      <p class="mt-1 text-sm text-gray-500">Comienza creando un nuevo código de insumo.</p>
      <div class="mt-6">
        <button @click="openCreateModal" class="btn-primary">
          <svg class="h-5 w-5 mr-2 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Crear Código
        </button>
      </div>
    </div>

    <!-- Modal para crear/editar -->
    <Teleport to="body">
      <div v-if="showModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50" @click.self="closeModal">
        <div class="relative top-20 mx-auto p-5 border w-full max-w-lg shadow-lg rounded-md bg-white">
          <div class="space-y-4">
            <div class="flex justify-between items-center border-b pb-3">
              <h3 class="text-xl font-semibold text-gray-900">
                {{ isEditing ? 'Editar Código de Insumo' : 'Crear Código de Insumo' }}
              </h3>
              <button @click="closeModal" class="text-gray-400 hover:text-gray-600">
                <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            
            <form @submit.prevent="saveSupplyCode" class="space-y-4">
              <div v-if="!isEditing">
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Código del Insumo <span class="text-red-500">*</span>
                </label>
                <input v-model.number="supplyCodeForm.code" type="number" min="1" class="form-input" 
                  placeholder="123456" required :disabled="isEditing" />
                <p class="mt-1 text-xs text-gray-500">Código único del insumo (no se puede modificar después)</p>
              </div>
              <div v-else>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Código del Insumo
                </label>
                <input type="number" class="form-input bg-gray-100" 
                  :value="supplyCodeForm.code" disabled />
                <p class="mt-1 text-xs text-gray-500">El código del insumo no se puede modificar</p>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Nombre del Insumo <span class="text-red-500">*</span>
                </label>
                <input v-model="supplyCodeForm.name" type="text" class="form-input" 
                  placeholder="Ej: Jeringa 10ml" required />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Código de Proveedor <span class="text-red-500">*</span>
                </label>
                <input v-model.number="supplyCodeForm.code_supplier" type="number" min="1" class="form-input" 
                  placeholder="789" required />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Stock Crítico <span class="text-red-500">*</span>
                </label>
                <input v-model.number="supplyCodeForm.critical_stock" type="number" min="1" class="form-input" 
                  placeholder="1" required />
                <p class="mt-1 text-xs text-gray-500">
                  Nivel mínimo de stock para generar alertas. Para insumos específicos, generalmente se usa 1.
                </p>
              </div>

              <div class="flex justify-end space-x-3 pt-4 border-t">
                <button type="button" @click="closeModal" class="btn-secondary">Cancelar</button>
                <button type="submit" :disabled="saving" class="btn-primary">
                  <span v-if="saving">Guardando...</span>
                  <span v-else>{{ isEditing ? 'Actualizar' : 'Crear' }}</span>
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import supplyCodeService from '@/services/config/supplyCodeService'
import inventoryService from '@/services/inventory/inventoryService'
import Swal from 'sweetalert2'

const supplyCodes = ref([])
const loading = ref(false)
const error = ref(null)
const searchTerm = ref('')
const showModal = ref(false)
const isEditing = ref(false)
const saving = ref(false)

// Estado de ordenamiento
const sortKey = ref('code')
const sortOrder = ref('asc')

// Estado de paginación
const currentPage = ref(1)
const itemsPerPage = 10

const supplyCodeForm = ref({
  code: null,
  name: '',
  code_supplier: null,
  critical_stock: 1
})

let searchTimeout = null

// Computed para obtener la lista ordenada
const sortedSupplyCodes = computed(() => {
  if (!supplyCodes.value || supplyCodes.value.length === 0) return []
  
  const sorted = [...supplyCodes.value].sort((a, b) => {
    let aVal = a[sortKey.value]
    let bVal = b[sortKey.value]
    
    // Manejo de strings (comparación case-insensitive)
    if (typeof aVal === 'string') {
      aVal = aVal.toLowerCase()
      bVal = bVal.toLowerCase()
    }
    
    // Comparación
    if (aVal < bVal) return sortOrder.value === 'asc' ? -1 : 1
    if (aVal > bVal) return sortOrder.value === 'asc' ? 1 : -1
    return 0
  })
  
  // Aplicar búsqueda si existe
  if (searchTerm.value.trim()) {
    const term = searchTerm.value.toLowerCase().trim()
    return sorted.filter(sc => 
      String(sc.code).includes(term) ||
      sc.name.toLowerCase().includes(term) ||
      String(sc.code_supplier).includes(term)
    )
  }
  
  return sorted
})

// Computed properties para paginación
const totalPages = computed(() => Math.ceil(sortedSupplyCodes.value.length / itemsPerPage))

const startIndex = computed(() => (currentPage.value - 1) * itemsPerPage)
const endIndex = computed(() => Math.min(startIndex.value + itemsPerPage, sortedSupplyCodes.value.length))

const paginatedSupplyCodes = computed(() => {
  return sortedSupplyCodes.value.slice(startIndex.value, endIndex.value)
})

// Función para ordenar por columna
const sortBy = (key) => {
  if (sortKey.value === key) {
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortKey.value = key
    sortOrder.value = 'asc'
  }
  currentPage.value = 1
}

const loadSupplyCodes = async () => {
  loading.value = true
  error.value = null
  try {
    const data = await supplyCodeService.getAllSupplyCodes()
    supplyCodes.value = data
  } catch (err) {
    error.value = err.response?.data?.error || err.message || 'Error al cargar códigos de insumos'
    console.error('Error loading supply codes:', err)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  
  searchTimeout = setTimeout(() => {
    currentPage.value = 1
  }, 300)
}

const clearSearch = () => {
  searchTerm.value = ''
  currentPage.value = 1
}

const openCreateModal = () => {
  isEditing.value = false
  supplyCodeForm.value = {
    code: null,
    name: '',
    code_supplier: null,
    critical_stock: 1
  }
  showModal.value = true
}

// Función para obtener el stock actual de un código de insumo
const getCurrentStock = async (supplyCode) => {
  try {
    const inventory = await inventoryService.getInventory()
    // Sumar las cantidades de todos los lotes para este código
    // El código puede estar en diferentes campos: code, supply_code
    const totalStock = inventory
      .filter(item => {
        const itemCode = item.code || item.supply_code
        return itemCode && Number(itemCode) === Number(supplyCode)
      })
      .reduce((sum, item) => {
        const amount = item.amount || item.current_amount || 0
        return sum + Number(amount)
      }, 0)
    return totalStock
  } catch (error) {
    console.error('Error al obtener stock actual:', error)
    // Si hay error, retornar null para que no bloquee la operación
    return null
  }
}

// Función para verificar si el stock está en nivel crítico
const isStockCritical = (currentStock, criticalStock) => {
  return currentStock !== null && currentStock <= criticalStock
}

const openEditModal = async (supplyCode) => {
  // Verificar stock crítico antes de abrir el modal
  const currentStock = await getCurrentStock(supplyCode.code)
  const isCritical = isStockCritical(currentStock, supplyCode.critical_stock)
  
  if (isCritical) {
    const result = await Swal.fire({
      title: '⚠️ Advertencia: Stock Crítico',
      html: `
        <div class="text-left">
          <p class="mb-3">Este código de insumo tiene <strong>stock crítico</strong>:</p>
          <div class="bg-red-50 border border-red-200 rounded-lg p-4 mb-3">
            <div class="flex items-center mb-2">
              <svg class="h-5 w-5 text-red-600 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
              <span class="font-semibold text-red-800">Stock Actual: ${currentStock}</span>
            </div>
            <div class="text-sm text-red-700">
              <p>Stock Crítico: ${supplyCode.critical_stock}</p>
              <p class="mt-1">¿Estás seguro de que deseas editar este código?</p>
            </div>
          </div>
        </div>
      `,
      icon: 'warning',
      showCancelButton: true,
      confirmButtonColor: '#d33',
      cancelButtonColor: '#3085d6',
      confirmButtonText: 'Sí, continuar',
      cancelButtonText: 'Cancelar',
      customClass: {
        popup: 'text-left'
      }
    })
    
    if (!result.isConfirmed) {
      return // Cancelar si el usuario no confirma
    }
  }
  
  isEditing.value = true
  supplyCodeForm.value = {
    code: supplyCode.code,
    name: supplyCode.name,
    code_supplier: supplyCode.code_supplier,
    critical_stock: supplyCode.critical_stock
  }
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  isEditing.value = false
  supplyCodeForm.value = {
    code: null,
    name: '',
    code_supplier: null,
    critical_stock: 1
  }
}

const saveSupplyCode = async () => {
  saving.value = true
  try {
    if (isEditing.value) {
      await supplyCodeService.updateSupplyCode(supplyCodeForm.value.code, {
        name: supplyCodeForm.value.name,
        code_supplier: supplyCodeForm.value.code_supplier,
        critical_stock: supplyCodeForm.value.critical_stock
      })
      await Swal.fire({
        icon: 'success',
        title: '¡Código actualizado!',
        text: 'El código de insumo ha sido actualizado exitosamente.',
        timer: 2000,
        showConfirmButton: false
      })
    } else {
      await supplyCodeService.createSupplyCode(supplyCodeForm.value)
      await Swal.fire({
        icon: 'success',
        title: '¡Código creado!',
        text: 'El código de insumo ha sido creado exitosamente.',
        timer: 2000,
        showConfirmButton: false
      })
    }
    closeModal()
    await loadSupplyCodes()
  } catch (err) {
    const errorMessage = err.response?.data?.error || err.message || 'Error al guardar el código de insumo'
    await Swal.fire({
      icon: 'error',
      title: 'Error',
      text: errorMessage
    })
  } finally {
    saving.value = false
  }
}

const confirmDelete = async (supplyCode) => {
  // Verificar stock crítico antes de mostrar el modal de confirmación
  const currentStock = await getCurrentStock(supplyCode.code)
  const isCritical = isStockCritical(currentStock, supplyCode.critical_stock)
  
  let result
  
  if (isCritical) {
    // Mostrar advertencia de stock crítico
    result = await Swal.fire({
      title: '⚠️ Advertencia: Stock Crítico',
      html: `
        <div class="text-left">
          <p class="mb-3">Este código de insumo tiene <strong>stock crítico</strong>:</p>
          <div class="bg-red-50 border border-red-200 rounded-lg p-4 mb-3">
            <div class="flex items-center mb-2">
              <svg class="h-5 w-5 text-red-600 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
              <span class="font-semibold text-red-800">Stock Actual: ${currentStock}</span>
            </div>
            <div class="text-sm text-red-700">
              <p>Stock Crítico: ${supplyCode.critical_stock}</p>
              <p class="mt-2 font-semibold">¿Estás seguro de que deseas eliminar este código?</p>
              <p class="mt-1 text-xs">Esta acción no se puede deshacer.</p>
            </div>
          </div>
        </div>
      `,
      icon: 'warning',
      showCancelButton: true,
      confirmButtonColor: '#d33',
      cancelButtonColor: '#3085d6',
      confirmButtonText: 'Sí, eliminar',
      cancelButtonText: 'Cancelar',
      customClass: {
        popup: 'text-left'
      }
    })
  } else {
    // Mostrar confirmación normal
    result = await Swal.fire({
      title: '¿Eliminar código de insumo?',
      text: `¿Estás seguro de que deseas eliminar el código "${supplyCode.code} - ${supplyCode.name}"?`,
      icon: 'warning',
      showCancelButton: true,
      confirmButtonColor: '#d33',
      cancelButtonColor: '#3085d6',
      confirmButtonText: 'Sí, eliminar',
      cancelButtonText: 'Cancelar'
    })
  }

  if (result.isConfirmed) {
    try {
      await supplyCodeService.deleteSupplyCode(supplyCode.code)
      await Swal.fire({
        icon: 'success',
        title: '¡Código eliminado!',
        text: 'El código de insumo ha sido eliminado exitosamente.',
        timer: 2000,
        showConfirmButton: false
      })
      await loadSupplyCodes()
    } catch (err) {
      const errorMessage = err.response?.data?.error || err.message || 'Error al eliminar el código de insumo'
      await Swal.fire({
        icon: 'error',
        title: 'Error',
        text: errorMessage
      })
    }
  }
}

// Helper functions
const getCriticalStockBadgeClass = (criticalStock) => {
  if (criticalStock === 1) return 'bg-red-100 text-red-800'
  if (criticalStock <= 5) return 'bg-orange-100 text-orange-800'
  return 'bg-yellow-100 text-yellow-800'
}

onMounted(() => {
  loadSupplyCodes()
})
</script>

