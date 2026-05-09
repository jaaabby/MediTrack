<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h2 class="text-2xl font-semibold text-gray-900">Configuración de Proveedores</h2>
          <p class="text-gray-600 mt-1">Gestiona los proveedores del sistema</p>
          <p v-if="!loading" class="text-sm text-gray-500 mt-1">Total: {{ filteredConfigs.length }} configuraciones</p>
        </div>
        <div class="flex gap-2">
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
    <FilterPanel
      :filters="[{ type: 'text', key: 'search', label: 'Buscar proveedor', placeholder: 'Buscar por nombre de proveedor...' }]"
      @filter-change="onFilterChange"
    />

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
    <DataTable
      v-else
      :columns="tableColumns"
      :rows="filteredConfigs"
      default-sort-key="supplier_name"
      empty-message="No hay configuraciones de proveedores registradas"
    >
      <template #cell-supplier_name="{ row }">
        <div class="text-sm font-medium text-gray-900">{{ row.supplier_name }}</div>
      </template>
      <template #cell-notes="{ row }">
        <span class="block max-w-xs truncate" :title="row.notes || '-'">{{ row.notes || '-' }}</span>
      </template>
      <template #actions="{ row }">
        <div class="flex justify-end space-x-2">
          <button @click="openEditModal(row)" class="text-warning-600 hover:text-warning-800 hover:bg-warning-50 p-1.5 rounded inline-flex items-center transition-colors" title="Editar">
            <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" /></svg>
          </button>
          <button @click="confirmDelete(row)" class="text-danger-600 hover:text-danger-800 hover:bg-danger-50 p-1.5 rounded transition-colors" title="Eliminar">
            <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>
          </button>
        </div>
      </template>
    </DataTable>

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
import { ref, computed, onMounted, reactive } from 'vue'
import supplierConfigService from '@/services/config/supplierConfigService'
import { useNotification } from '@/composables/useNotification'
import { useAlert } from '@/composables/useAlert'
import FilterPanel from '@/components/common/FilterPanel.vue'
import DataTable from '@/components/common/DataTable.vue'

const { success: showSuccess, error: showError, warning: showWarning } = useNotification()
const { confirmDanger } = useAlert()

const configs = ref([])
const loading = ref(false)
const error = ref(null)
const filterState = reactive({ search: '' })
const onFilterChange = (key, value) => { filterState[key] = value }
const showModal = ref(false)
const isEditing = ref(false)
const saving = ref(false)

const tableColumns = [
  { key: 'supplier_name', label: 'Proveedor' },
  { key: 'notes', label: 'Notas', sortable: false, wrap: true }
]

const filteredConfigs = computed(() => {
  if (!configs.value || configs.value.length === 0) return []
  if (!filterState.search.trim()) return configs.value
  const term = filterState.search.toLowerCase().trim()
  return configs.value.filter(c =>
    c.supplier_name.toLowerCase().includes(term) ||
    (c.notes && c.notes.toLowerCase().includes(term))
  )
})

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

