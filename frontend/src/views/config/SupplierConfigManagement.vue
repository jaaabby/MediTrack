<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h2 class="text-2xl font-semibold text-gray-900">Configuración de Proveedores</h2>
          <p class="text-gray-600 mt-1">Gestiona los proveedores del sistema</p>
          <p v-if="!loading" class="text-sm text-gray-500 mt-1">Total: {{ sortedConfigs.length }} configuraciones</p>
        </div>
        <div class="flex gap-2">
          <button 
            @click="loadConfigs" 
            :disabled="loading"
            class="btn-secondary flex items-center justify-center"
            title="Recargar configuraciones">
            <svg class="h-5 w-5" :class="{ 'animate-spin': loading }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
          </button>
          <button @click="openCreateModal" class="btn-primary flex items-center justify-center">
            <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            Nueva Configuración
          </button>
        </div>
      </div>
    </div>

    <!-- Búsqueda -->
    <div class="card">
      <div class="flex flex-col sm:flex-row sm:items-end gap-4">
        <div class="flex-1">
          <label class="block text-sm font-medium text-gray-700 mb-2">Buscar proveedor</label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            <input type="text" placeholder="Buscar por nombre de proveedor..." 
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
        <span class="ml-3 text-gray-600">Cargando configuraciones...</span>
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
            <h3 class="text-sm font-medium text-red-800">Error al cargar configuraciones</h3>
            <div class="mt-2 text-sm text-red-700">{{ error }}</div>
            <button @click="loadConfigs" class="btn-secondary mt-4 text-sm">Reintentar</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Tabla de configuraciones -->
    <div v-else class="card overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th scope="col" @click="sortBy('supplier_name')"
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                <div class="flex items-center space-x-1">
                  <span>Proveedor</span>
                  <span class="flex flex-col -space-y-1">
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'supplier_name' && sortOrder === 'asc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                    </svg>
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'supplier_name' && sortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z"/>
                    </svg>
                  </span>
                </div>
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Notas
              </th>
              <th scope="col" class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider sticky right-0 bg-gray-50 z-10 shadow-[-2px_0_4px_rgba(0,0,0,0.05)]">
                Acciones
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="config in paginatedConfigs" :key="config.supplier_name" 
              class="hover:bg-gray-50 transition-colors duration-150">
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm font-medium text-gray-900">{{ config.supplier_name }}</div>
              </td>
              <td class="px-6 py-4 text-sm text-gray-500 max-w-xs truncate" :title="config.notes || '-'">
                {{ config.notes || '-' }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium sticky right-0 bg-white z-10 shadow-[-2px_0_4px_rgba(0,0,0,0.05)]">
                <div class="flex justify-end space-x-2">
                  <button @click="openEditModal(config)" 
                    class="text-warning-600 hover:text-warning-800 hover:bg-warning-50 p-1.5 rounded inline-flex items-center gap-1 transition-colors"
                    title="Editar">
                    <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                    <span class="font-medium text-xs">Editar</span>
                  </button>
                  <button @click="confirmDelete(config)" 
                    class="text-danger-600 hover:text-danger-800 hover:bg-danger-50 p-1.5 rounded transition-colors"
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
    <div v-if="!loading && sortedConfigs.length > 0" class="card">
      <div class="flex flex-col sm:flex-row items-center justify-between gap-4">
        <div class="text-sm text-gray-700 text-center sm:text-left">
          Mostrando {{ startIndex + 1 }} a {{ endIndex }} de {{ sortedConfigs.length }} configuraciones
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
    <div v-if="!loading && configs.length === 0" class="card text-center py-12">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900">No hay configuraciones de proveedores</h3>
      <p class="mt-1 text-sm text-gray-500">Comienza creando una nueva configuración de proveedor.</p>
      <div class="mt-6">
        <button @click="openCreateModal" class="btn-primary">
          <svg class="h-5 w-5 mr-2 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Crear Configuración
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
                {{ isEditing ? 'Editar Configuración de Proveedor' : 'Crear Configuración de Proveedor' }}
              </h3>
              <button @click="closeModal" class="text-gray-400 hover:text-gray-600">
                <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            
            <form @submit.prevent="saveConfig" class="space-y-4">
              <div v-if="!isEditing">
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Nombre del Proveedor <span class="text-red-500">*</span>
                </label>
                <input v-model="configForm.supplier_name" type="text" class="form-input" 
                  placeholder="Ej: Proveedor ABC" required :disabled="isEditing" />
                <p class="mt-1 text-xs text-gray-500">Nombre único del proveedor</p>
              </div>
              <div v-else>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Nombre del Proveedor
                </label>
                <input type="text" class="form-input bg-gray-100" 
                  :value="configForm.supplier_name" disabled />
                <p class="mt-1 text-xs text-gray-500">El nombre del proveedor no se puede modificar</p>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Notas (opcional)
                </label>
                <textarea v-model="configForm.notes" rows="3" class="form-input" 
                  placeholder="Notas adicionales sobre este proveedor"></textarea>
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
import supplierConfigService from '@/services/config/supplierConfigService'
import { useNotification } from '@/composables/useNotification'
import { useAlert } from '@/composables/useAlert'

const { success: showSuccess, error: showError, warning: showWarning } = useNotification()
const { confirmDanger } = useAlert()

const configs = ref([])
const loading = ref(false)
const error = ref(null)
const searchTerm = ref('')
const showModal = ref(false)
const isEditing = ref(false)
const saving = ref(false)

// Estado de ordenamiento
const sortKey = ref('supplier_name')
const sortOrder = ref('asc')

// Estado de paginación
const currentPage = ref(1)
const itemsPerPage = 10

const configForm = ref({
  supplier_name: '',
  notes: ''
})

let searchTimeout = null

// Computed para obtener la lista ordenada
const sortedConfigs = computed(() => {
  if (!configs.value || configs.value.length === 0) return []
  
  const sorted = [...configs.value].sort((a, b) => {
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
    return sorted.filter(c => 
      c.supplier_name.toLowerCase().includes(term) || 
      (c.notes && c.notes.toLowerCase().includes(term))
    )
  }
  
  return sorted
})

// Computed properties para paginación
const totalPages = computed(() => Math.ceil(sortedConfigs.value.length / itemsPerPage))

const startIndex = computed(() => (currentPage.value - 1) * itemsPerPage)
const endIndex = computed(() => Math.min(startIndex.value + itemsPerPage, sortedConfigs.value.length))

const paginatedConfigs = computed(() => {
  return sortedConfigs.value.slice(startIndex.value, endIndex.value)
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

const loadConfigs = async () => {
  loading.value = true
  error.value = null
  try {
    const data = await supplierConfigService.getAllSupplierConfigs()
    configs.value = data
  } catch (err) {
    error.value = err.response?.data?.error || err.message || 'Error al cargar configuraciones de proveedores'
    console.error('Error loading configs:', err)
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
  configForm.value = {
    supplier_name: '',
    notes: ''
  }
  showModal.value = true
}

const openEditModal = (config) => {
  isEditing.value = true
  configForm.value = {
    supplier_name: config.supplier_name,
    notes: config.notes || ''
  }
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  isEditing.value = false
  configForm.value = {
    supplier_name: '',
    notes: ''
  }
}

const saveConfig = async () => {
  saving.value = true
  try {
    if (isEditing.value) {
      await supplierConfigService.updateSupplierConfig(configForm.value.supplier_name, {
        notes: configForm.value.notes
      })
      showSuccess('La configuración del proveedor ha sido actualizada exitosamente.')
    } else {
      await supplierConfigService.createSupplierConfig(configForm.value)
      showSuccess('La configuración del proveedor ha sido creada exitosamente.')
    }
    closeModal()
    await loadConfigs()
  } catch (err) {
    const errorMessage = err.response?.data?.error || err.message || 'Error al guardar la configuración'
    showError(errorMessage)
  } finally {
    saving.value = false
  }
}

const confirmDelete = async (config) => {
  const confirmed = await confirmDanger(
    `¿Estás seguro de que deseas eliminar la configuración para "${config.supplier_name}"?`,
    'Confirmar eliminación'
  )
  if (!confirmed) {
    return
  }

  try {
    await supplierConfigService.deleteSupplierConfig(config.supplier_name)
    showSuccess('La configuración ha sido eliminada exitosamente.')
    await loadConfigs()
  } catch (err) {
    const errorMessage = err.response?.data?.error || err.message || 'Error al eliminar la configuración'
    showError(errorMessage)
  }
}

onMounted(() => {
  loadConfigs()
})
</script>

