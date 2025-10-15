<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h2 class="text-2xl font-semibold text-gray-900">Tipos de Cirugía</h2>
          <p class="text-gray-600 mt-1">Gestiona los tipos de procedimientos quirúrgicos</p>
          <p v-if="!loading" class="text-sm text-gray-500 mt-1">Total: {{ sortedSurgeries.length }} tipos</p>
        </div>
        <button @click="openCreateModal" class="btn-primary flex items-center justify-center">
          <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Nuevo Tipo
        </button>
      </div>
    </div>

    <!-- Búsqueda -->
    <div class="card">
      <div class="flex flex-col sm:flex-row sm:items-end gap-4">
        <div class="flex-1">
          <label class="block text-sm font-medium text-gray-700 mb-2">Buscar tipo de cirugía</label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            <input type="text" placeholder="Buscar por nombre o descripción..." 
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
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
        <span class="ml-3 text-gray-600">Cargando tipos de cirugía...</span>
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
            <h3 class="text-sm font-medium text-red-800">Error al cargar tipos de cirugía</h3>
            <div class="mt-2 text-sm text-red-700">{{ error }}</div>
            <button @click="loadSurgeries" class="btn-secondary mt-4 text-sm">Reintentar</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Tabla de tipos de cirugía -->
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
              <th scope="col" @click="sortBy('duration')"
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                <div class="flex items-center space-x-1">
                  <span>Duración (horas)</span>
                  <span class="flex flex-col -space-y-1">
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'duration' && sortOrder === 'asc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                    </svg>
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'duration' && sortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'" 
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
            <tr v-for="surgery in paginatedSurgeries" :key="surgery.id" 
              class="hover:bg-gray-50 transition-colors duration-150">
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                #{{ surgery.id }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm font-medium text-gray-900">{{ surgery.name }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {{ surgery.duration }}h
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <div class="flex justify-end space-x-2">
                  <button @click="openEditModal(surgery)" 
                    class="inline-flex items-center px-3 py-1.5 border border-transparent text-xs font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors"
                    title="Editar">
                    <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                  </button>
                  <button @click="confirmDelete(surgery)" 
                    class="inline-flex items-center px-3 py-1.5 border border-transparent text-xs font-medium rounded-md text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 transition-colors"
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
    <div v-if="!loading && sortedSurgeries.length > 0" class="card">
      <div class="flex flex-col sm:flex-row items-center justify-between gap-4">
        <div class="text-sm text-gray-700 text-center sm:text-left">
          Mostrando {{ startIndex + 1 }} a {{ endIndex }} de {{ sortedSurgeries.length }} tipos de cirugía
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
    <div v-if="!loading && surgeries.length === 0" class="card text-center py-12">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900">No hay tipos de cirugía</h3>
      <p class="mt-1 text-sm text-gray-500">Comienza creando un nuevo tipo de cirugía.</p>
      <div class="mt-6">
        <button @click="openCreateModal" class="btn-primary">
          <svg class="h-5 w-5 mr-2 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Crear Tipo de Cirugía
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
              {{ isEditing ? 'Editar Tipo de Cirugía' : 'Crear Tipo de Cirugía' }}
            </h3>
            <button @click="closeModal" class="text-gray-400 hover:text-gray-600">
              <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          
          <form @submit.prevent="saveSurgery" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Nombre <span class="text-red-500">*</span>
              </label>
              <input v-model="surgeryForm.name" type="text" class="form-input" 
                placeholder="Ej: Apendicectomía" required />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Duración (horas) <span class="text-red-500">*</span>
              </label>
              <input v-model.number="surgeryForm.duration" type="number" step="0.5" min="0.5" 
                class="form-input" placeholder="Ej: 2.5" required />
              <p class="mt-1 text-xs text-gray-500">Ingrese la duración estimada en horas (mínimo 0.5)</p>
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
import surgeryService from '@/services/surgeryService'
import Swal from 'sweetalert2'

const surgeries = ref([])
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

const surgeryForm = ref({
  name: '',
  duration: 1
})

let searchTimeout = null

// Computed para obtener la lista ordenada
const sortedSurgeries = computed(() => {
  if (!surgeries.value || surgeries.value.length === 0) return []
  
  const sorted = [...surgeries.value].sort((a, b) => {
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
  
  return sorted
})

// Computed properties para paginación
const totalPages = computed(() => Math.ceil(sortedSurgeries.value.length / itemsPerPage))

const startIndex = computed(() => (currentPage.value - 1) * itemsPerPage)
const endIndex = computed(() => Math.min(startIndex.value + itemsPerPage, sortedSurgeries.value.length))

const paginatedSurgeries = computed(() => {
  return sortedSurgeries.value.slice(startIndex.value, endIndex.value)
})

// Función para ordenar por columna
const sortBy = (key) => {
  if (sortKey.value === key) {
    // Si ya estamos ordenando por esta columna, cambiar dirección
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
  } else {
    // Nueva columna, ordenar ascendente por defecto
    sortKey.value = key
    sortOrder.value = 'asc'
  }
  currentPage.value = 1 // Resetear a la primera página al ordenar
}

const loadSurgeries = async () => {
  loading.value = true
  error.value = null
  try {
    const data = await surgeryService.getAllSurgeries()
    surgeries.value = data
  } catch (err) {
    error.value = err.message || 'Error al cargar tipos de cirugía'
    console.error('Error loading surgeries:', err)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  
  searchTimeout = setTimeout(async () => {
    currentPage.value = 1 // Resetear a la primera página al buscar
    if (searchTerm.value.trim()) {
      loading.value = true
      try {
        const data = await surgeryService.searchSurgeries(searchTerm.value)
        surgeries.value = data
      } catch (err) {
        error.value = err.message
      } finally {
        loading.value = false
      }
    } else {
      loadSurgeries()
    }
  }, 300)
}

const openCreateModal = () => {
  isEditing.value = false
  surgeryForm.value = {
    name: '',
    duration: 1
  }
  showModal.value = true
}

const openEditModal = (surgery) => {
  isEditing.value = true
  surgeryForm.value = {
    id: surgery.id,
    name: surgery.name,
    duration: surgery.duration || 1
  }
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  surgeryForm.value = {}
}

const saveSurgery = async () => {
  // Validaciones
  if (!surgeryForm.value.name || !surgeryForm.value.name.trim()) {
    await Swal.fire({
      icon: 'warning',
      title: 'Campo requerido',
      text: 'El nombre del tipo de cirugía es obligatorio',
      confirmButtonText: 'Aceptar'
    })
    return
  }

  if (!surgeryForm.value.duration || surgeryForm.value.duration < 0.5) {
    await Swal.fire({
      icon: 'warning',
      title: 'Duración inválida',
      text: 'La duración debe ser al menos 0.5 horas',
      confirmButtonText: 'Aceptar'
    })
    return
  }

  saving.value = true
  try {
    const surgeryData = {
      name: surgeryForm.value.name.trim(),
      duration: Number(surgeryForm.value.duration)
    }

    if (isEditing.value) {
      await surgeryService.updateSurgery(surgeryForm.value.id, surgeryData)
      await loadSurgeries()
      closeModal()
      await Swal.fire({
        icon: 'success',
        title: 'Actualizado',
        text: 'Tipo de cirugía actualizado exitosamente',
        timer: 2000,
        showConfirmButton: false
      })
    } else {
      await surgeryService.createSurgery(surgeryData)
      await loadSurgeries()
      closeModal()
      await Swal.fire({
        icon: 'success',
        title: 'Creado',
        text: 'Tipo de cirugía creado exitosamente',
        timer: 2000,
        showConfirmButton: false
      })
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

    await Swal.fire({
      icon: 'error',
      title: 'Error al guardar',
      text: errorMessage,
      confirmButtonText: 'Aceptar'
    })
  } finally {
    saving.value = false
  }
}

const confirmDelete = async (surgery) => {
  const result = await Swal.fire({
    title: '¿Estás seguro?',
    html: `¿Deseas eliminar el tipo de cirugía <strong>"${surgery.name}"</strong>?<br><small class="text-gray-600">Esta acción no se puede deshacer.</small>`,
    icon: 'warning',
    showCancelButton: true,
    confirmButtonColor: '#dc2626',
    cancelButtonColor: '#6b7280',
    confirmButtonText: 'Sí, eliminar',
    cancelButtonText: 'Cancelar',
    reverseButtons: true
  })

  if (result.isConfirmed) {
    try {
      await surgeryService.deleteSurgery(surgery.id)
      await loadSurgeries()
      await Swal.fire({
        icon: 'success',
        title: 'Eliminado',
        text: 'Tipo de cirugía eliminado exitosamente',
        timer: 2000,
        showConfirmButton: false
      })
    } catch (err) {
      console.error('Error al eliminar:', err)
      await Swal.fire({
        icon: 'error',
        title: 'Error',
        text: 'Error al eliminar: ' + (err.response?.data?.error || err.message),
        confirmButtonText: 'Aceptar'
      })
    }
  }
}

const clearSearch = () => {
  searchTerm.value = ''
  currentPage.value = 1
  loadSurgeries()
}

onMounted(() => {
  loadSurgeries()
})
</script>
