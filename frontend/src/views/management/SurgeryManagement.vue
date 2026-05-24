<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h2 class="text-2xl font-semibold text-gray-900">Tipos de Cirugía</h2>
          <p class="text-gray-600 mt-1">Gestiona los tipos de procedimientos quirúrgicos</p>
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
    <FilterPanel :filters="[{ type: 'text', key: 'search', label: 'Buscar tipo de cirugía', placeholder: 'Buscar por nombre o descripción...' }]" :result-count="filteredSurgeries.length" @filter-change="onFilterChange" />

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
    <DataTable
      v-else
      :columns="tableColumns"
      :rows="filteredSurgeries"
      default-sort-key="id"
      empty-message="No hay tipos de cirugía"
      :table-actions="[
        { type: 'edit', label: 'Editar', onClick: (row) => openEditModal(row) },
        { type: 'delete', onClick: (row) => confirmDelete(row) },
      ]"
    >
      <template #cell-duration="{ row }">
        {{ formatDuration(row.duration) }}
      </template>
      <template #cell-specialty="{ row }">
        <span v-if="row.specialty" class="text-sm text-gray-900">
          {{ row.specialty.name }}
        </span>
        <span v-else class="text-gray-400 text-xs">Sin especialidad</span>
      </template>

    </DataTable>

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
              <p v-if="durationError" class="mt-1 text-sm text-red-600">{{ durationError }}</p>
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
import { ref, reactive, computed, onMounted } from 'vue'
import surgeryService from '@/services/management/surgeryService'
import medicalSpecialtyService from '@/services/config/medicalSpecialtyService'
import { useNotification } from '@/composables/useNotification'
import { useAlert } from '@/composables/useAlert'
import FilterPanel from '@/components/common/FilterPanel.vue'
import DataTable from '@/components/common/DataTable.vue'
import { normalize } from '@/utils/normalize'

const { success: showSuccess, error: showError, warning: showWarning } = useNotification()
const { confirmDanger } = useAlert()

const surgeries = ref([])
const specialties = ref([])
const loading = ref(false)
const loadingSpecialties = ref(false)
const error = ref(null)
const filterState = reactive({ search: '' })
const showModal = ref(false)
const isEditing = ref(false)
const saving = ref(false)
const durationError = ref('')

const tableColumns = [
  { key: 'name', label: 'Nombre' },
  { key: 'duration', label: 'Duración (horas)' },
  { key: 'specialty', label: 'Especialidad', sortable: false }
]

const filteredSurgeries = computed(() => {
  if (!surgeries.value || surgeries.value.length === 0) return []
  if (!filterState.search.trim()) return surgeries.value
  const term = normalize(filterState.search)
  return surgeries.value.filter(s =>
    normalize(s.name).includes(term) ||
    normalize(s.description).includes(term)
  )
})

const surgeryForm = ref({
  name: '',
  duration: 1,
  specialty_id: null
})

const onFilterChange = (key, value) => { filterState[key] = value }

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
    const data = await medicalSpecialtyService.getActiveSpecialties()
    specialties.value = data
  } catch (err) {
    console.error('Error loading specialties:', err)
  } finally {
    loadingSpecialties.value = false
  }
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
  durationError.value = ''
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
  durationError.value = ''
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
  durationError.value = ''
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
    durationError.value = 'La duración debe ser mayor a 0'
    return
  }
  durationError.value = ''

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

onMounted(async () => {
  await Promise.all([
    loadSurgeries(),
    loadSpecialties()
  ])
})
</script>
