<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <div>
        <h2 class="text-2xl font-semibold text-gray-900">Gestión de Doctores</h2>
        <p class="text-gray-600 mt-1">Gestiona los doctores del sistema directamente desde la tabla de usuarios</p>
      </div>
    </div>

    <!-- Filtros -->
    <FilterPanel :filters="filterConfig" :result-count="filteredDoctors.length" @filter-change="onFilterChange" />

    <!-- Loading -->
    <div v-if="loading" class="card">
      <div class="flex justify-center items-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        <span class="ml-3 text-gray-600">Cargando doctores...</span>
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
            <h3 class="text-sm font-medium text-red-800">Error al cargar doctores</h3>
            <div class="mt-2 text-sm text-red-700">{{ error }}</div>
            <button @click="loadDoctors" class="btn-secondary mt-4 text-sm">Reintentar</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Tabla de doctores -->
    <DataTable
      v-else
      :columns="tableColumns"
      :rows="filteredDoctors"
      default-sort-key="name"
      empty-message="No hay doctores registrados"
    >
      <template #cell-rut="{ row }"><span class="font-medium">{{ row.rut }}</span></template>
      <template #cell-name="{ row }">
        <div class="text-sm font-medium text-gray-900">{{ row.name || '-' }}</div>
      </template>
      <template #cell-email="{ row }">{{ row.email || '-' }}</template>
      <template #cell-specialty_id="{ row }">{{ getSpecialtyName(row.specialty_id) }}</template>
      <template #cell-medical_center="{ row }">{{ row.medical_center?.name || '-' }}</template>
      <template #cell-is_active="{ row }">
        <span v-if="row.is_active" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">Activo</span>
        <span v-else class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">Inactivo</span>
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
    <div v-if="!loading && filteredDoctors.length === 0" class="card text-center py-12">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900">No hay doctores registrados</h3>
      <p class="mt-1 text-sm text-gray-500">{{ filterState.specialty || filterState.search ? 'No se encontraron resultados con los filtros aplicados.' : 'No hay doctores registrados. Crea doctores desde la gestión de usuarios.' }}</p>
    </div>

    <!-- Modal para editar -->
    <Teleport to="body">
      <div v-if="showModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50" @click.self="closeModal">
        <div class="relative top-20 mx-auto p-5 border w-full max-w-2xl shadow-lg rounded-md bg-white max-h-[90vh] overflow-y-auto">
          <div class="space-y-4">
            <div class="flex justify-between items-center border-b pb-3">
              <h3 class="text-xl font-semibold text-gray-900">
                Editar Doctor
              </h3>
              <button @click="closeModal" class="text-gray-400 hover:text-gray-600">
                <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            
            <form @submit.prevent="saveDoctor" class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  RUT <span class="text-red-500">*</span>
                </label>
                <input v-model="doctorForm.rut" type="text" class="form-input" 
                  placeholder="Ej: 12345678-9" required disabled />
                <p class="mt-1 text-xs text-gray-500">El RUT no puede modificarse</p>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Nombre <span class="text-red-500">*</span>
                </label>
                <input v-model="doctorForm.name" type="text" class="form-input" 
                  placeholder="Ej: Dr. Juan Pérez" required />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Email <span class="text-red-500">*</span>
                </label>
                <input v-model="doctorForm.email" type="email" class="form-input" 
                  placeholder="Ej: doctor@meditrack.com" required />
              </div>

              <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    Centro Médico <span class="text-red-500">*</span>
                  </label>
                  <select v-model.number="doctorForm.medical_center_id" class="form-select" required>
                    <option :value="null">Seleccionar centro médico</option>
                    <option v-for="center in medicalCenters" :key="center.id" :value="center.id">
                      {{ center.name }}
                    </option>
                  </select>
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    Especialidad Médica
                  </label>
                  <select v-model.number="doctorForm.specialty_id" class="form-select">
                    <option :value="null">Sin especialidad</option>
                    <option v-for="specialty in specialties" :key="specialty.id" :value="specialty.id">
                      {{ specialty.name }}
                    </option>
                  </select>
                </div>
              </div>

              <div>
                <label class="flex items-center space-x-2">
                  <input v-model="doctorForm.is_active" type="checkbox" class="rounded border-gray-300 text-blue-600 focus:ring-blue-500" />
                  <span class="text-sm font-medium text-gray-700">Doctor activo</span>
                </label>
                <p class="mt-1 text-xs text-gray-500">Los doctores inactivos no podrán iniciar sesión</p>
              </div>

              <div class="flex justify-end space-x-3 pt-4 border-t">
                <button type="button" @click="closeModal" class="btn-secondary">Cancelar</button>
                <button type="submit" :disabled="saving" class="btn-primary">
                  <span v-if="saving">Guardando...</span>
                  <span v-else>Actualizar</span>
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
import doctorInfoService from '@/services/config/doctorInfoService'
import medicalSpecialtyService from '@/services/config/medicalSpecialtyService'
import medicalCenterService from '@/services/config/medicalCenterService'
import { useNotification } from '@/composables/useNotification'
import { useAlert } from '@/composables/useAlert'
import FilterPanel from '@/components/common/FilterPanel.vue'
import DataTable from '@/components/common/DataTable.vue'

const { success: showSuccess, error: showError, warning: showWarning } = useNotification()
const { confirmDanger } = useAlert()

const doctors = ref([])
const specialties = ref([])
const medicalCenters = ref([])
const loading = ref(false)
const error = ref(null)
const filterState = reactive({ search: '', specialty: '' })
const showModal = ref(false)
const isEditing = ref(false)
const saving = ref(false)

const tableColumns = [
  { key: 'rut', label: 'RUT', sortable: false },
  { key: 'name', label: 'Nombre' },
  { key: 'email', label: 'Email', sortable: false },
  { key: 'specialty_id', label: 'Especialidad', sortable: false },
  { key: 'medical_center', label: 'Centro Médico', sortable: false },
  { key: 'is_active', label: 'Estado', sortable: false }
]
const doctorForm = ref({
  rut: '',
  name: '',
  email: '',
  password: '',
  role: 'doctor',
  medical_center_id: null,
  specialty_id: null,
  is_active: true
})

const onFilterChange = (key, value) => { filterState[key] = value }

const filterConfig = computed(() => [
  {
    type: 'select',
    key: 'specialty',
    label: 'Filtrar por Especialidad',
    options: [
      { value: '', label: 'Todas las especialidades' },
      ...specialties.value.map(s => ({ value: String(s.id), label: s.name }))
    ]
  },
  { type: 'text', key: 'search', label: 'Buscar Doctor', placeholder: 'Buscar por nombre, RUT o email...' }
])

// Computed para obtener la lista filtrada
const filteredDoctors = computed(() => {
  let filtered = [...doctors.value]

  if (filterState.specialty) {
    filtered = filtered.filter(d => d.specialty_id === parseInt(filterState.specialty))
  }

  if (filterState.search.trim()) {
    const term = filterState.search.toLowerCase().trim()
    filtered = filtered.filter(d => {
      const name = (d.name || '').toLowerCase()
      const email = (d.email || '').toLowerCase()
      const rut = (d.rut || '').toLowerCase()
      return name.includes(term) || email.includes(term) || rut.includes(term)
    })
  }

  return filtered
})

// Funciones auxiliares
const getSpecialtyName = (specialtyId) => {
  if (!specialtyId) return '-'
  const specialty = specialties.value.find(s => s.id === specialtyId)
  return specialty ? specialty.name : `Especialidad #${specialtyId}`
}

const loadDoctors = async () => {
  loading.value = true
  error.value = null
  try {
    const data = await doctorInfoService.getAllDoctors()
    doctors.value = data
  } catch (err) {
    error.value = err.message || 'Error al cargar doctores'
    console.error('Error loading doctors:', err)
  } finally {
    loading.value = false
  }
}

const loadSpecialties = async () => {
  try {
    const data = await medicalSpecialtyService.getActiveSpecialties()
    specialties.value = data
  } catch (err) {
    console.error('Error loading specialties:', err)
  }
}

const loadMedicalCenters = async () => {
  try {
    const response = await medicalCenterService.getAll()
    medicalCenters.value = response.data || []
  } catch (err) {
    console.error('Error loading medical centers:', err)
  }
}


const openEditModal = (doctor) => {
  isEditing.value = true
  doctorForm.value = {
    rut: doctor.rut,
    name: doctor.name || '',
    email: doctor.email || '',
    password: '', // No se muestra la contraseña
    role: 'doctor',
    medical_center_id: doctor.medical_center_id || null,
    specialty_id: doctor.specialty_id || null,
    is_active: doctor.is_active !== undefined ? doctor.is_active : true
  }
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  doctorForm.value = {
    rut: '',
    name: '',
    email: '',
    password: '',
    role: 'doctor',
    medical_center_id: null,
    specialty_id: null,
    is_active: true
  }
}

const saveDoctor = async () => {
  // Validaciones
  if (!doctorForm.value.rut || !doctorForm.value.rut.trim()) {
    showWarning('El RUT es obligatorio')
    return
  }

  if (!doctorForm.value.name || !doctorForm.value.name.trim()) {
    showWarning('El nombre es obligatorio')
    return
  }

  if (!doctorForm.value.email || !doctorForm.value.email.trim()) {
    showWarning('El email es obligatorio')
    return
  }

  // Validar formato de email
  const emailRegex = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
  if (!emailRegex.test(doctorForm.value.email.trim())) {
    showWarning('El formato del email no es válido')
    return
  }

  if (!doctorForm.value.medical_center_id) {
    showWarning('Debe seleccionar un centro médico')
    return
  }

  // Validar formato de RUT (básico)
  const rutRegex = /^\d{7,8}-[\dkK]$/
  if (!rutRegex.test(doctorForm.value.rut.trim())) {
    showWarning('El RUT debe tener el formato 12345678-9')
    return
  }

  saving.value = true
  try {
    const doctorData = {
      rut: doctorForm.value.rut.trim(),
      name: doctorForm.value.name.trim(),
      email: doctorForm.value.email.trim(),
      role: 'doctor',
      medical_center_id: doctorForm.value.medical_center_id,
      specialty_id: doctorForm.value.specialty_id || null,
      is_active: doctorForm.value.is_active !== undefined ? doctorForm.value.is_active : true
    }

    await doctorInfoService.updateDoctor(doctorForm.value.rut, doctorData)
    await loadDoctors()
    closeModal()
    showSuccess('Doctor actualizado exitosamente')
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

const confirmDelete = async (doctor) => {
  const doctorName = doctor.name || doctor.rut
  
  const confirmed = await confirmDanger(
    `¿Deseas eliminar al doctor "${doctorName}"?\n\nEsta acción no se puede deshacer.`,
    'Confirmar eliminación'
  )
  if (!confirmed) {
    return
  }

  try {
    await doctorInfoService.deleteDoctor(doctor.rut)
    await loadDoctors()
    showSuccess('Doctor eliminado exitosamente')
  } catch (err) {
    console.error('Error al eliminar:', err)
    showError('Error al eliminar: ' + (err.response?.data?.error || err.message))
  }
}

onMounted(async () => {
  await Promise.all([
    loadSpecialties(),
    loadMedicalCenters(),
    loadDoctors()
  ])
})
</script>
