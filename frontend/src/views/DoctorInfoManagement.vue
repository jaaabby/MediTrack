<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h2 class="text-2xl font-semibold text-gray-900">Información de Doctores</h2>
          <p class="text-gray-600 mt-1">Gestiona la información extendida de los doctores del sistema</p>
          <p v-if="!loading" class="text-sm text-gray-500 mt-1">Total: {{ sortedDoctors.length }} doctores</p>
        </div>
        <button @click="openCreateModal" class="btn-primary flex items-center justify-center">
          <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Nuevo Doctor
        </button>
      </div>
    </div>

    <!-- Filtros -->
    <div class="card">
      <div class="flex flex-col sm:flex-row gap-4">
        <div class="flex-1">
          <label class="block text-sm font-medium text-gray-700 mb-2">Filtrar por Especialidad</label>
          <select v-model="selectedSpecialtyId" @change="filterBySpecialty" class="form-select">
            <option value="">Todas las especialidades</option>
            <option v-for="specialty in specialties" :key="specialty.id" :value="specialty.id">
              {{ specialty.name }}
            </option>
          </select>
        </div>
        <div class="flex-1">
          <label class="block text-sm font-medium text-gray-700 mb-2">Buscar Doctor</label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            <input type="text" placeholder="Buscar por nombre, RUT o licencia..." 
              class="form-input pl-10 w-full" v-model="searchTerm" @input="handleSearch" />
          </div>
        </div>
        <div class="flex items-end">
          <button class="btn-secondary px-4 py-2 h-10" @click="clearFilters" :disabled="!selectedSpecialtyId && !searchTerm">
            Limpiar Filtros
          </button>
        </div>
      </div>
    </div>

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
    <div v-else class="card overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                RUT
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Nombre
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Licencia Médica
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Especialidad
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Años Exp.
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Teléfono
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Estado
              </th>
              <th scope="col" class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                Acciones
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="doctor in paginatedDoctors" :key="doctor.user_rut" 
              class="hover:bg-gray-50 transition-colors duration-150">
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                {{ doctor.user_rut }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm font-medium text-gray-900">{{ doctor.user?.name || '-' }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ doctor.medical_license || '-' }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ getSpecialtyName(doctor.specialty_id) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {{ doctor.years_of_experience || '-' }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ doctor.phone || '-' }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span v-if="doctor.is_available" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                  Disponible
                </span>
                <span v-else class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">
                  No Disponible
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <div class="flex justify-end space-x-2">
                  <button @click="openEditModal(doctor)" 
                    class="btn-primary text-xs px-3 py-1.5"
                    title="Editar">
                    <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                  </button>
                  <button @click="confirmDelete(doctor)" 
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
    <div v-if="!loading && sortedDoctors.length > 0" class="card">
      <div class="flex flex-col sm:flex-row items-center justify-between gap-4">
        <div class="text-sm text-gray-700 text-center sm:text-left">
          Mostrando {{ startIndex + 1 }} a {{ endIndex }} de {{ sortedDoctors.length }} doctores
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
    <div v-if="!loading && filteredDoctors.length === 0" class="card text-center py-12">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900">No hay doctores registrados</h3>
      <p class="mt-1 text-sm text-gray-500">{{ selectedSpecialtyId || searchTerm ? 'No se encontraron resultados con los filtros aplicados.' : 'Comienza registrando información de un doctor.' }}</p>
      <div class="mt-6" v-if="!selectedSpecialtyId && !searchTerm">
        <button @click="openCreateModal" class="btn-primary">
          <svg class="h-5 w-5 mr-2 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Registrar Doctor
        </button>
      </div>
    </div>

    <!-- Modal para crear/editar -->
    <Teleport to="body">
      <div v-if="showModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50" @click.self="closeModal">
        <div class="relative top-20 mx-auto p-5 border w-full max-w-2xl shadow-lg rounded-md bg-white max-h-[90vh] overflow-y-auto">
          <div class="space-y-4">
            <div class="flex justify-between items-center border-b pb-3">
              <h3 class="text-xl font-semibold text-gray-900">
                {{ isEditing ? 'Editar Información de Doctor' : 'Registrar Información de Doctor' }}
              </h3>
              <button @click="closeModal" class="text-gray-400 hover:text-gray-600">
                <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            
            <form @submit.prevent="saveDoctorInfo" class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  RUT del Usuario <span class="text-red-500">*</span>
                </label>
                <input v-model="doctorForm.user_rut" type="text" class="form-input" 
                  placeholder="Ej: 12345678-9" required :disabled="isEditing" />
                <p class="mt-1 text-xs text-gray-500">El RUT debe corresponder a un usuario existente con rol "doctor"</p>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Licencia Médica
                </label>
                <input v-model="doctorForm.medical_license" type="text" class="form-input" 
                  placeholder="Ej: LIC-12345678" />
              </div>

              <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    Fecha de Vencimiento de Licencia
                  </label>
                  <input v-model="doctorForm.license_expiration_date" type="date" class="form-input" />
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    Años de Experiencia
                  </label>
                  <input v-model.number="doctorForm.years_of_experience" type="number" min="0" 
                    class="form-input" placeholder="Ej: 10" />
                </div>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Especialización
                </label>
                <input v-model="doctorForm.specialization" type="text" class="form-input" 
                  placeholder="Ej: Cirugía General" />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Especialidad Médica
                </label>
                <select v-model="doctorForm.specialty_id" class="form-select">
                  <option :value="null">Sin especialidad</option>
                  <option v-for="specialty in specialties" :key="specialty.id" :value="specialty.id">
                    {{ specialty.name }}
                  </option>
                </select>
              </div>

              <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    Teléfono
                  </label>
                  <input v-model="doctorForm.phone" type="tel" class="form-input" 
                    placeholder="Ej: +56 9 1234 5678" />
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    Contacto de Emergencia
                  </label>
                  <input v-model="doctorForm.emergency_contact" type="text" class="form-input" 
                    placeholder="Nombre del contacto" />
                </div>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Teléfono de Emergencia
                </label>
                <input v-model="doctorForm.emergency_phone" type="tel" class="form-input" 
                  placeholder="Ej: +56 9 8765 4321" />
              </div>

              <div>
                <label class="flex items-center space-x-2">
                  <input v-model="doctorForm.is_available" type="checkbox" class="rounded border-gray-300 text-blue-600 focus:ring-blue-500" />
                  <span class="text-sm font-medium text-gray-700">Doctor disponible</span>
                </label>
                <p class="mt-1 text-xs text-gray-500">Los doctores no disponibles no aparecerán en listas de selección</p>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Notas
                </label>
                <textarea v-model="doctorForm.notes" rows="3" class="form-input" 
                  placeholder="Notas adicionales sobre el doctor"></textarea>
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
import doctorInfoService from '@/services/doctorInfoService'
import medicalSpecialtyService from '@/services/medicalSpecialtyService'
import Swal from 'sweetalert2'

const doctors = ref([])
const specialties = ref([])
const loading = ref(false)
const error = ref(null)
const searchTerm = ref('')
const selectedSpecialtyId = ref('')
const showModal = ref(false)
const isEditing = ref(false)
const saving = ref(false)

// Estado de paginación
const currentPage = ref(1)
const itemsPerPage = 10

const doctorForm = ref({
  user_rut: '',
  medical_license: '',
  license_expiration_date: '',
  specialization: '',
  specialty_id: null,
  years_of_experience: 0,
  phone: '',
  emergency_contact: '',
  emergency_phone: '',
  is_available: true,
  notes: ''
})

let searchTimeout = null

// Computed para obtener la lista filtrada
const filteredDoctors = computed(() => {
  let filtered = [...doctors.value]

  // Filtrar por especialidad
  if (selectedSpecialtyId.value) {
    filtered = filtered.filter(d => d.specialty_id === selectedSpecialtyId.value)
  }

  // Filtrar por búsqueda
  if (searchTerm.value.trim()) {
    const term = searchTerm.value.toLowerCase().trim()
    filtered = filtered.filter(d => {
      const userName = (d.user?.name || '').toLowerCase()
      const license = (d.medical_license || '').toLowerCase()
      const rut = d.user_rut.toLowerCase()
      return userName.includes(term) || license.includes(term) || rut.includes(term)
    })
  }

  return filtered
})

// Computed para ordenar
const sortedDoctors = computed(() => {
  return [...filteredDoctors.value].sort((a, b) => {
    const nameA = (a.user?.name || '').toLowerCase()
    const nameB = (b.user?.name || '').toLowerCase()
    return nameA.localeCompare(nameB)
  })
})

// Computed properties para paginación
const totalPages = computed(() => Math.ceil(sortedDoctors.value.length / itemsPerPage))

const startIndex = computed(() => (currentPage.value - 1) * itemsPerPage)
const endIndex = computed(() => Math.min(startIndex.value + itemsPerPage, sortedDoctors.value.length))

const paginatedDoctors = computed(() => {
  return sortedDoctors.value.slice(startIndex.value, endIndex.value)
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
    error.value = err.message || 'Error al cargar información de doctores'
    console.error('Error loading doctors:', err)
  } finally {
    loading.value = false
  }
}

const loadSpecialties = async () => {
  try {
    const data = await medicalSpecialtyService.getAllSpecialties()
    specialties.value = data.filter(s => s.is_active)
  } catch (err) {
    console.error('Error loading specialties:', err)
  }
}

const handleSearch = () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    currentPage.value = 1
  }, 300)
}

const filterBySpecialty = () => {
  currentPage.value = 1
}

const clearFilters = () => {
  selectedSpecialtyId.value = ''
  searchTerm.value = ''
  currentPage.value = 1
}

const openCreateModal = () => {
  isEditing.value = false
  doctorForm.value = {
    user_rut: '',
    medical_license: '',
    license_expiration_date: '',
    specialization: '',
    specialty_id: null,
    years_of_experience: 0,
    phone: '',
    emergency_contact: '',
    emergency_phone: '',
    is_available: true,
    notes: ''
  }
  showModal.value = true
}

const openEditModal = (doctor) => {
  isEditing.value = true
  doctorForm.value = {
    user_rut: doctor.user_rut,
    medical_license: doctor.medical_license || '',
    license_expiration_date: doctor.license_expiration_date ? doctor.license_expiration_date.split('T')[0] : '',
    specialization: doctor.specialization || '',
    specialty_id: doctor.specialty_id || null,
    years_of_experience: doctor.years_of_experience || 0,
    phone: doctor.phone || '',
    emergency_contact: doctor.emergency_contact || '',
    emergency_phone: doctor.emergency_phone || '',
    is_available: doctor.is_available !== undefined ? doctor.is_available : true,
    notes: doctor.notes || ''
  }
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  doctorForm.value = {
    user_rut: '',
    medical_license: '',
    license_expiration_date: '',
    specialization: '',
    specialty_id: null,
    years_of_experience: 0,
    phone: '',
    emergency_contact: '',
    emergency_phone: '',
    is_available: true,
    notes: ''
  }
}

const saveDoctorInfo = async () => {
  // Validaciones
  if (!doctorForm.value.user_rut || !doctorForm.value.user_rut.trim()) {
    await Swal.fire({
      icon: 'warning',
      title: 'Campo requerido',
      text: 'El RUT del usuario es obligatorio',
      confirmButtonText: 'Aceptar'
    })
    return
  }

  // Validar formato de RUT (básico)
  const rutRegex = /^\d{7,8}-[\dkK]$/
  if (!rutRegex.test(doctorForm.value.user_rut.trim())) {
    await Swal.fire({
      icon: 'warning',
      title: 'RUT inválido',
      text: 'El RUT debe tener el formato 12345678-9',
      confirmButtonText: 'Aceptar'
    })
    return
  }

  saving.value = true
  try {
    const doctorData = {
      user_rut: doctorForm.value.user_rut.trim(),
      medical_license: doctorForm.value.medical_license.trim() || null,
      license_expiration_date: doctorForm.value.license_expiration_date || null,
      specialization: doctorForm.value.specialization.trim() || null,
      specialty_id: doctorForm.value.specialty_id || null,
      years_of_experience: doctorForm.value.years_of_experience || 0,
      phone: doctorForm.value.phone.trim() || null,
      emergency_contact: doctorForm.value.emergency_contact.trim() || null,
      emergency_phone: doctorForm.value.emergency_phone.trim() || null,
      is_available: doctorForm.value.is_available !== undefined ? doctorForm.value.is_available : true,
      notes: doctorForm.value.notes.trim() || null
    }

    if (isEditing.value) {
      await doctorInfoService.updateDoctorInfo(doctorForm.value.user_rut, doctorData)
      await loadDoctors()
      closeModal()
      await Swal.fire({
        icon: 'success',
        title: 'Actualizado',
        text: 'Información del doctor actualizada exitosamente',
        timer: 2000,
        showConfirmButton: false
      })
    } else {
      await doctorInfoService.createDoctorInfo(doctorData)
      await loadDoctors()
      closeModal()
      await Swal.fire({
        icon: 'success',
        title: 'Creado',
        text: 'Información del doctor registrada exitosamente',
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

const confirmDelete = async (doctor) => {
  const doctorName = doctor.user?.name || doctor.user_rut
  
  const result = await Swal.fire({
    title: '¿Estás seguro?',
    html: `¿Deseas eliminar la información del doctor <strong>"${doctorName}"</strong>?<br><small class="text-gray-600">Esta acción no se puede deshacer.</small>`,
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
      await doctorInfoService.deleteDoctorInfo(doctor.user_rut)
      await loadDoctors()
      await Swal.fire({
        icon: 'success',
        title: 'Eliminado',
        text: 'Información del doctor eliminada exitosamente',
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

onMounted(async () => {
  await Promise.all([
    loadSpecialties(),
    loadDoctors()
  ])
})
</script>

