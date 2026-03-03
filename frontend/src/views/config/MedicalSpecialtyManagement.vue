<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h2 class="text-2xl font-semibold text-gray-900">Especialidades Médicas</h2>
          <p class="text-gray-600 mt-1">Gestiona las especialidades médicas del sistema</p>
          <p v-if="!loading" class="text-sm text-gray-500 mt-1">Total: {{ sortedSpecialties.length }} especialidades</p>
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
    <div class="card">
      <div class="flex flex-col sm:flex-row sm:items-end gap-4">
        <div class="flex-1">
          <label class="block text-sm font-medium text-gray-700 mb-2">Buscar especialidad</label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            <input type="text" placeholder="Buscar por nombre o código..." 
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
    <div v-else class="card overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th scope="col" @click="sortBy('id')"
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                <div class="flex items-center space-x-1">
                  <span>ID</span>
                  <span class="flex flex-col -space-y-1">
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'id' && sortOrder === 'asc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                    </svg>
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'id' && sortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'" 
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
                Código
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Descripción
              </th>
              <th scope="col" @click="sortBy('is_active')"
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                <div class="flex items-center space-x-1">
                  <span>Estado</span>
                  <span class="flex flex-col -space-y-1">
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'is_active' && sortOrder === 'asc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                    </svg>
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'is_active' && sortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z"/>
                    </svg>
                  </span>
                </div>
              </th>
              <th scope="col" class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider sticky right-0 bg-gray-50 z-10 shadow-[-2px_0_4px_rgba(0,0,0,0.05)]">
                Acciones
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="specialty in paginatedSpecialties" :key="specialty.id" 
              class="hover:bg-gray-50 transition-colors duration-150">
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                #{{ specialty.id }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm font-medium text-gray-900">{{ specialty.name }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ specialty.code || '-' }}
              </td>
              <td class="px-6 py-4 text-sm text-gray-500 max-w-xs truncate" 
                :title="specialty.description && specialty.description.trim() ? specialty.description : undefined">
                {{ specialty.description || '-' }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span v-if="specialty.is_active" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                  Activa
                </span>
                <span v-else class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">
                  Inactiva
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium sticky right-0 bg-white z-10 shadow-[-2px_0_4px_rgba(0,0,0,0.05)]">
                <div class="flex justify-end space-x-2">
                  <button @click="openEditModal(specialty)" 
                    class="text-warning-600 hover:text-warning-800 hover:bg-warning-50 p-1.5 rounded inline-flex items-center gap-1 transition-colors"
                    title="Editar">
                    <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                    <span class="font-medium text-xs">Editar</span>
                  </button>
                  <button @click="confirmDelete(specialty)" 
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
    <div v-if="!loading && sortedSpecialties.length > 0" class="card">
      <div class="flex flex-col sm:flex-row items-center justify-between gap-4">
        <div class="text-sm text-gray-700 text-center sm:text-left">
          Mostrando {{ startIndex + 1 }} a {{ endIndex }} de {{ sortedSpecialties.length }} especialidades
        </div>
        <div class="flex items-center gap-2">
          <button class="btn-secondary px-3 py-2 text-sm min-w-[70px]" :disabled="currentPage === 1"
            @click="previousPage">
            <span class="hidden sm:inline">Anterior</span>
            <span class="sm:hidden">Ant.</span>
          </button>
          <span class="px-3 py-2 text-sm text-gray-700 bg-gray-100 rounded-md min-w-[90px] text-center">
            Página {{ currentPage }} de {{ totalPages }}
          </span>
          <button class="btn-secondary px-3 py-2 text-sm min-w-[70px]" :disabled="currentPage === totalPages"
            @click="nextPage">
            <span class="hidden sm:inline">Siguiente</span>
            <span class="sm:hidden">Sig.</span>
          </button>
        </div>
      </div>
    </div>

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
import { ref, computed, onMounted } from 'vue'
import medicalSpecialtyService from '@/services/config/medicalSpecialtyService'
import { useNotification } from '@/composables/useNotification'
import { useAlert } from '@/composables/useAlert'

const { success: showSuccess, error: showError, warning: showWarning } = useNotification()
const { confirmDanger } = useAlert()

const specialties = ref([])
const loading = ref(false)
const error = ref(null)
const searchTerm = ref('')
const showModal = ref(false)
const isEditing = ref(false)
const saving = ref(false)

// Estado de ordenamiento
const sortKey = ref('id')
const sortOrder = ref('asc')

// Estado de paginación
const currentPage = ref(1)
const itemsPerPage = 10

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

let searchTimeout = null

// Computed para obtener la lista ordenada
const sortedSpecialties = computed(() => {
  if (!specialties.value || specialties.value.length === 0) return []
  
  const sorted = [...specialties.value].sort((a, b) => {
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
    return sorted.filter(s => 
      s.name.toLowerCase().includes(term) || 
      (s.code && s.code.toLowerCase().includes(term)) ||
      (s.description && s.description.toLowerCase().includes(term))
    )
  }
  
  return sorted
})

// Computed properties para paginación
const totalPages = computed(() => Math.ceil(sortedSpecialties.value.length / itemsPerPage))

const startIndex = computed(() => (currentPage.value - 1) * itemsPerPage)
const endIndex = computed(() => Math.min(startIndex.value + itemsPerPage, sortedSpecialties.value.length))

const paginatedSpecialties = computed(() => {
  return sortedSpecialties.value.slice(startIndex.value, endIndex.value)
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

const handleSearch = () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  
  searchTimeout = setTimeout(() => {
    currentPage.value = 1
  }, 300)
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

const clearSearch = () => {
  searchTerm.value = ''
  currentPage.value = 1
}

const previousPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
  }
}

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
  }
}

onMounted(() => {
  loadSpecialties()
})
</script>

