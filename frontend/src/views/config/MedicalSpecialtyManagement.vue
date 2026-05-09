<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h2 class="text-2xl font-semibold text-gray-900">Especialidades Médicas</h2>
          <p class="text-gray-600 mt-1">Gestiona las especialidades médicas del sistema</p>
        </div>
        <button @click="openCreateModal" class="btn-primary flex items-center justify-center">
          <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Nueva Especialidad
        </button>
      </div>
    </div>

    <!-- Búsqueda -->
    <FilterPanel
      :filters="[{ type: 'text', key: 'search', label: 'Buscar especialidad', placeholder: 'Buscar por nombre o código...' }]"
      :result-count="filteredSpecialties.length"
      @filter-change="onFilterChange"
    />

    <!-- Loading -->
    <div v-if="loading" class="card">
      <div class="flex justify-center items-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        <span class="ml-3 text-gray-600">Cargando especialidades...</span>
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
            <h3 class="text-sm font-medium text-red-800">Error al cargar especialidades</h3>
            <div class="mt-2 text-sm text-red-700">{{ error }}</div>
            <button @click="loadSpecialties" class="btn-secondary mt-4 text-sm">Reintentar</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Tabla de especialidades -->
    <DataTable
      v-else
      :columns="tableColumns"
      :rows="filteredSpecialties"
      default-sort-key="id"
      empty-message="No hay especialidades médicas registradas"
    >
      <template #cell-name="{ row }">
        <div class="text-sm font-medium text-gray-900">{{ row.name }}</div>
      </template>
      <template #cell-code="{ row }">{{ row.code || '-' }}</template>
      <template #cell-description="{ row }">
        <span class="block max-w-xs truncate" :title="row.description && row.description.trim() ? row.description : undefined">{{ row.description || '-' }}</span>
      </template>
      <template #cell-is_active="{ row }">
        <span v-if="row.is_active" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">Activa</span>
        <span v-else class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">Inactiva</span>
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
    <div v-if="!loading && specialties.length === 0" class="card text-center py-12">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900">No hay especialidades médicas</h3>
      <p class="mt-1 text-sm text-gray-500">Comienza creando una nueva especialidad médica.</p>
      <div class="mt-6">
        <button @click="openCreateModal" class="btn-primary">
          <svg class="h-5 w-5 mr-2 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Crear Especialidad
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
                {{ isEditing ? 'Editar Especialidad Médica' : 'Crear Especialidad Médica' }}
              </h3>
              <button @click="closeModal" class="text-gray-400 hover:text-gray-600">
                <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            
            <form @submit.prevent="saveSpecialty" class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Nombre <span class="text-red-500">*</span>
                </label>
                <input 
                  v-model="specialtyForm.name" 
                  type="text" 
                  class="form-input" 
                  :class="{ 'border-red-500 focus:border-red-500 focus:ring-red-500': formErrors.name }"
                  placeholder="Ej: Traumatología" 
                  @input="clearFieldError('name')" 
                />
                <p v-if="formErrors.name" class="text-sm text-red-600 mt-1">
                  {{ formErrors.name }}
                </p>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Código
                </label>
                <input 
                  v-model="specialtyForm.code" 
                  type="text" 
                  class="form-input" 
                  :class="{ 'border-red-500 focus:border-red-500 focus:ring-red-500': formErrors.code }"
                  placeholder="Ej: TRAUMA" 
                  @input="clearFieldError('code')" 
                />
                <p v-if="formErrors.code" class="text-sm text-red-600 mt-1">
                  {{ formErrors.code }}
                </p>
                <p v-else class="mt-1 text-xs text-gray-500">Código único para identificar la especialidad</p>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Descripción
                </label>
                <textarea 
                  v-model="specialtyForm.description" 
                  rows="3" 
                  class="form-input" 
                  :class="{ 'border-red-500 focus:border-red-500 focus:ring-red-500': formErrors.description }"
                  placeholder="Descripción de la especialidad médica"
                  @input="clearFieldError('description')" 
                ></textarea>
                <p v-if="formErrors.description" class="text-sm text-red-600 mt-1">
                  {{ formErrors.description }}
                </p>
              </div>

              <div>
                <label class="flex items-center space-x-2">
                  <input v-model="specialtyForm.is_active" type="checkbox" class="rounded border-gray-300 text-blue-600 focus:ring-blue-500" />
                  <span class="text-sm font-medium text-gray-700">Especialidad activa</span>
                </label>
                <p class="mt-1 text-xs text-gray-500">Las especialidades inactivas no aparecerán en listas de selección</p>
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
import medicalSpecialtyService from '@/services/config/medicalSpecialtyService'
import { useNotification } from '@/composables/useNotification'
import { useAlert } from '@/composables/useAlert'
import FilterPanel from '@/components/common/FilterPanel.vue'
import DataTable from '@/components/common/DataTable.vue'

const { success: showSuccess, error: showError, warning: showWarning } = useNotification()
const { confirmDanger } = useAlert()

const specialties = ref([])
const loading = ref(false)
const error = ref(null)
const filterState = reactive({ search: '' })
const onFilterChange = (key, value) => { filterState[key] = value }
const showModal = ref(false)
const isEditing = ref(false)
const saving = ref(false)

const tableColumns = [
  { key: 'name', label: 'Nombre' },
  { key: 'code', label: 'Código', sortable: false },
  { key: 'description', label: 'Descripción', sortable: false, wrap: true },
  { key: 'is_active', label: 'Estado' }
]

const specialtyForm = ref({
  name: '',
  code: '',
  description: '',
  is_active: true
})

// Estado de errores de validación
const formErrors = ref({
  name: '',
  code: '',
  description: ''
})

const filteredSpecialties = computed(() => {
  if (!specialties.value || specialties.value.length === 0) return []
  if (!filterState.search.trim()) return specialties.value
  const term = filterState.search.toLowerCase().trim()
  return specialties.value.filter(s =>
    s.name.toLowerCase().includes(term) ||
    (s.code && s.code.toLowerCase().includes(term)) ||
    (s.description && s.description.toLowerCase().includes(term))
  )
})

const loadSpecialties = async () => {
  loading.value = true
  error.value = null
  try {
    const data = await medicalSpecialtyService.getAllSpecialties()
    specialties.value = data
  } catch (err) {
    error.value = err.message || 'Error al cargar especialidades médicas'
    console.error('Error loading specialties:', err)
  } finally {
    loading.value = false
  }
}


const openCreateModal = () => {
  isEditing.value = false
  specialtyForm.value = {
    name: '',
    code: '',
    description: '',
    is_active: true
  }
  showModal.value = true
}

const openEditModal = (specialty) => {
  isEditing.value = true
  specialtyForm.value = {
    id: specialty.id,
    name: specialty.name || '',
    code: specialty.code || '',
    description: specialty.description || '',
    is_active: specialty.is_active !== undefined ? specialty.is_active : true
  }
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  specialtyForm.value = {
    name: '',
    code: '',
    description: '',
    is_active: true
  }
  // Resetear errores al cerrar
  formErrors.value = {
    name: '',
    code: '',
    description: ''
  }
}

// Validación de formulario
const validateForm = () => {
  // Resetear errores
  formErrors.value = {
    name: '',
    code: '',
    description: ''
  }
  
  let hasErrors = false
  
  // Validar nombre (obligatorio)
  if (!specialtyForm.value.name || !specialtyForm.value.name.trim()) {
    formErrors.value.name = 'El nombre de la especialidad es obligatorio.'
    hasErrors = true
  }
  if (specialtyForm.value.name && specialtyForm.value.name.trim().length > 0 && specialtyForm.value.name.trim().length < 3) {
    formErrors.value.name = 'El nombre debe tener al menos 3 caracteres.'
    hasErrors = true
  }
  if (specialtyForm.value.name && specialtyForm.value.name.trim().length > 100) {
    formErrors.value.name = 'El nombre no puede exceder los 100 caracteres.'
    hasErrors = true
  }
  
  // Validar código (opcional, pero si se ingresa debe ser válido)
  if (specialtyForm.value.code && specialtyForm.value.code.trim()) {
    if (specialtyForm.value.code.trim().length < 2) {
      formErrors.value.code = 'El código debe tener al menos 2 caracteres.'
      hasErrors = true
    }
    if (specialtyForm.value.code.trim().length > 20) {
      formErrors.value.code = 'El código no puede exceder los 20 caracteres.'
      hasErrors = true
    }
  }
  
  // Validar descripción (opcional, pero si se ingresa debe ser válido)
  if (specialtyForm.value.description && specialtyForm.value.description.trim()) {
    if (specialtyForm.value.description.trim().length > 500) {
      formErrors.value.description = 'La descripción no puede exceder los 500 caracteres.'
      hasErrors = true
    }
  }
  
  return !hasErrors
}

// Limpiar error individual cuando el usuario empiece a editar
const clearFieldError = (field) => {
  if (formErrors.value[field]) {
    formErrors.value[field] = ''
  }
}

const saveSpecialty = async () => {
  // Validar formulario
  if (!validateForm()) {
    showWarning('Por favor corrige los errores en el formulario')
    return
  }

  saving.value = true
  try {
    const specialtyData = {
      name: specialtyForm.value.name.trim(),
      code: specialtyForm.value.code.trim() || null,
      description: specialtyForm.value.description.trim() || null,
      is_active: Boolean(specialtyForm.value.is_active)
    }

    if (isEditing.value) {
      await medicalSpecialtyService.updateSpecialty(specialtyForm.value.id, specialtyData)
      await loadSpecialties()
      closeModal()
      showSuccess('Especialidad médica actualizada exitosamente')
    } else {
      await medicalSpecialtyService.createSpecialty(specialtyData)
      await loadSpecialties()
      closeModal()
      showSuccess('Especialidad médica creada exitosamente')
    }
  } catch (err) {
    console.error('Error al guardar:', err)
    let errorMessage = 'Error desconocido al guardar'
    
    if (err.response?.data?.error) {
      errorMessage = err.response.data.error
    } else if (err.response?.data?.message) {
      errorMessage = err.response.data.message
    } else if (err.message) {
      errorMessage = err.message
    }

    showError(errorMessage)
  } finally {
    saving.value = false
  }
}

const confirmDelete = async (specialty) => {
  const confirmed = await confirmDanger(
    `¿Deseas eliminar la especialidad "${specialty.name}"?\n\nEsta acción no se puede deshacer.`,
    'Confirmar eliminación'
  )
  if (!confirmed) {
    return
  }

  try {
    await medicalSpecialtyService.deleteSpecialty(specialty.id)
    await loadSpecialties()
    showSuccess('Especialidad médica eliminada exitosamente')
  } catch (err) {
    console.error('Error al eliminar:', err)
    showError('Error al eliminar: ' + (err.response?.data?.error || err.message))
  }
}

onMounted(() => {
  loadSpecialties()
})
</script>

