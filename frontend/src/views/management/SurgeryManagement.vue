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
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
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
              <th scope="col" @click="sortBy('specialty')"
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                <div class="flex items-center space-x-1">
                  <span>Especialidad</span>
                  <span class="flex flex-col -space-y-1">
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'specialty' && sortOrder === 'asc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                    </svg>
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'specialty' && sortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'" 
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
                {{ formatDuration(surgery.duration) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                <span v-if="surgery.specialty" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                  {{ surgery.specialty.name }}
                </span>
                <span v-else class="text-gray-400 text-xs">Sin especialidad</span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <div class="flex justify-end space-x-2">
                  <button @click="openEditModal(surgery)" 
                    class="btn-primary text-xs px-3 py-1.5"
                    title="Editar">
                    <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                  </button>
                  <button @click="confirmDelete(surgery)" 
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
                Duración <span class="text-red-500">*</span>
              </label>
              <div class="grid grid-cols-2 gap-3">
                <div>
                  <label class="block text-xs text-gray-600 mb-1">Horas</label>
                  <input 
                    v-model.number="surgeryForm.durationHours" 
                    type="number" 
                    min="0" 
                    max="24"
                    class="form-input" 
                    placeholder="0"
                    required 
                  />
                </div>
                <div>
                  <label class="block text-xs text-gray-600 mb-1">Minutos</label>
                  <select 
                    v-model.number="surgeryForm.durationMinutes" 
                    class="form-input"
                    required
                  >
                    <option :value="0">0 min</option>
                    <option :value="5">5 min</option>
                    <option :value="15">15 min</option>
                    <option :value="30">30 min</option>
                    <option :value="45">45 min</option>
                  </select>
                </div>
              </div>
              <p class="mt-1 text-xs text-gray-500">
                Duración total: {{ formatDurationFromInput(surgeryForm.durationHours, surgeryForm.durationMinutes) }}
              </p>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Especialidad Médica
              </label>
              <select v-model.number="surgeryForm.specialty_id" class="form-input">
                <option :value="null">Sin especialidad</option>
                <option v-for="specialty in specialties" :key="specialty.id" :value="specialty.id">
                  {{ specialty.name }}
                </option>
              </select>
              <p class="mt-1 text-xs text-gray-500">Seleccione la especialidad médica asociada (opcional)</p>
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
import surgeryService from '@/services/management/surgeryService'
import medicalSpecialtyService from '@/services/config/medicalSpecialtyService'
import { useNotification } from '@/composables/useNotification'
import { useAlert } from '@/composables/useAlert'

const { success: showSuccess, error: showError, warning: showWarning } = useNotification()
const { confirmDanger } = useAlert()

const surgeries = ref([])
const specialties = ref([])
const loading = ref(false)
const loadingSpecialties = ref(false)
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
  duration: 1,
  specialty_id: null
})

let searchTimeout = null

// Computed para obtener la lista ordenada
const sortedSurgeries = computed(() => {
  if (!surgeries.value || surgeries.value.length === 0) return []
  
  const sorted = [...surgeries.value].sort((a, b) => {
    let aVal = a[sortKey.value]
    let bVal = b[sortKey.value]
    
    // Manejo especial para ordenar por especialidad
    if (sortKey.value === 'specialty') {
      aVal = a.specialty?.name || ''
      bVal = b.specialty?.name || ''
    }
    
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

const loadSpecialties = async () => {
  loadingSpecialties.value = true
  try {
    const data = await medicalSpecialtyService.getAllSpecialties()
    specialties.value = data.filter(s => s.is_active)
  } catch (err) {
    console.error('Error loading specialties:', err)
  } finally {
    loadingSpecialties.value = false
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

// Funciones de conversión
const decimalToHoursMinutes = (decimalHours) => {
  if (!decimalHours || decimalHours <= 0) {
    return { hours: 0, minutes: 0 }
  }
  const hours = Math.floor(decimalHours)
  const minutes = Math.round((decimalHours - hours) * 60)
  // Redondear minutos a los valores permitidos (0, 5, 15, 30, 45)
  // Si está cerca de 5, usar 5; si no, redondear a múltiplos de 15
  let roundedMinutes
  if (minutes >= 2.5 && minutes < 10) {
    roundedMinutes = 5
  } else {
    roundedMinutes = Math.round(minutes / 15) * 15
  }
  return { hours, minutes: roundedMinutes >= 60 ? 0 : roundedMinutes }
}

const hoursMinutesToDecimal = (hours, minutes) => {
  return hours + (minutes / 60)
}

const formatDuration = (decimalHours) => {
  if (!decimalHours || decimalHours <= 0) {
    return '0h 0min'
  }
  const { hours, minutes } = decimalToHoursMinutes(decimalHours)
  if (minutes === 0) {
    return `${hours}h`
  }
  return `${hours}h ${minutes}min`
}

const formatDurationFromInput = (hours, minutes) => {
  if (!hours && !minutes) return '0h 0min'
  if (minutes === 0) {
    return `${hours || 0}h`
  }
  return `${hours || 0}h ${minutes}min`
}

const openCreateModal = () => {
  isEditing.value = false
  surgeryForm.value = {
    name: '',
    durationHours: 1,
    durationMinutes: 0,
    specialty_id: null
  }
  showModal.value = true
}

const openEditModal = (surgery) => {
  isEditing.value = true
  const { hours, minutes } = decimalToHoursMinutes(surgery.duration || 1)
  surgeryForm.value = {
    id: surgery.id,
    name: surgery.name,
    durationHours: hours,
    durationMinutes: minutes,
    specialty_id: surgery.specialty_id || null
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
    showWarning('El nombre del tipo de cirugía es obligatorio')
    return
  }

  // Validar duración
  const durationHours = surgeryForm.value.durationHours || 0
  const durationMinutes = surgeryForm.value.durationMinutes || 0
  const totalDuration = hoursMinutesToDecimal(durationHours, durationMinutes)
  
  if (totalDuration <= 0) {
    showWarning('La duración debe ser mayor a 0')
    return
  }

  saving.value = true
  try {
    // Convertir horas y minutos a decimal para guardar en BD
    const durationDecimal = hoursMinutesToDecimal(durationHours, durationMinutes)
    
    const surgeryData = {
      name: surgeryForm.value.name.trim(),
      duration: durationDecimal,
      specialty_id: surgeryForm.value.specialty_id || null
    }

    if (isEditing.value) {
      await surgeryService.updateSurgery(surgeryForm.value.id, surgeryData)
      await loadSurgeries()
      closeModal()
      showSuccess('Tipo de cirugía actualizado exitosamente')
    } else {
      await surgeryService.createSurgery(surgeryData)
      await loadSurgeries()
      closeModal()
      showSuccess('Tipo de cirugía creado exitosamente')
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

    showError('Error al guardar: ' + errorMessage)
  } finally {
    saving.value = false
  }
}

const confirmDelete = async (surgery) => {
  const confirmed = await confirmDanger(
    `¿Deseas eliminar el tipo de cirugía "${surgery.name}"?\n\nEsta acción no se puede deshacer.`,
    'Confirmar eliminación'
  )
  if (!confirmed) {
    return
  }

  try {
    await surgeryService.deleteSurgery(surgery.id)
    await loadSurgeries()
    showSuccess('Tipo de cirugía eliminado exitosamente')
  } catch (err) {
    console.error('Error al eliminar:', err)
    showError('Error al eliminar: ' + (err.response?.data?.error || err.message))
  }
}

const clearSearch = () => {
  searchTerm.value = ''
  currentPage.value = 1
  loadSurgeries()
}

onMounted(async () => {
  await Promise.all([
    loadSurgeries(),
    loadSpecialties()
  ])
})
</script>
