<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h2 class="text-2xl font-semibold text-gray-900">Insumos Típicos por Cirugía</h2>
          <p class="text-gray-600 mt-1">Gestiona los insumos típicos asociados a cada tipo de cirugía</p>
          <p v-if="!loading" class="text-sm text-gray-500 mt-1">Total: {{ typicalSupplies.length }} asociaciones</p>
        </div>
        <button @click="openCreateModal" class="btn-primary flex items-center justify-center">
          <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Nueva Asociación
        </button>
      </div>
    </div>

    <!-- Filtros -->
    <div class="card">
      <div class="flex flex-col sm:flex-row gap-4">
        <div class="flex-1">
          <label class="block text-sm font-medium text-gray-700 mb-2">Filtrar por Cirugía</label>
          <select v-model="selectedSurgeryId" @change="filterBySurgery" class="form-select">
            <option value="">Todas las cirugías</option>
            <option v-for="surgery in surgeries" :key="surgery.id" :value="surgery.id">
              {{ surgery.name }}
            </option>
          </select>
        </div>
        <div class="flex-1">
          <label class="block text-sm font-medium text-gray-700 mb-2">Buscar por Insumo</label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            <input type="text" placeholder="Buscar por código o nombre de insumo..." 
              class="form-input pl-10 w-full" v-model="searchTerm" @input="handleSearch" />
          </div>
        </div>
        <div class="flex items-end">
          <button class="btn-secondary px-4 py-2 h-10" @click="clearFilters" :disabled="!selectedSurgeryId && !searchTerm">
            Limpiar Filtros
          </button>
        </div>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="card">
      <div class="flex justify-center items-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        <span class="ml-3 text-gray-600">Cargando insumos típicos...</span>
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
            <h3 class="text-sm font-medium text-red-800">Error al cargar insumos típicos</h3>
            <div class="mt-2 text-sm text-red-700">{{ error }}</div>
            <button @click="loadTypicalSupplies" class="btn-secondary mt-4 text-sm">Reintentar</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Tabla de insumos típicos -->
    <div v-else class="card overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Cirugía
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Código Insumo
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Nombre Insumo
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Cantidad Típica
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Requerido
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Notas
              </th>
              <th scope="col" class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                Acciones
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="supply in filteredTypicalSupplies" :key="supply.id" 
              class="hover:bg-gray-50 transition-colors duration-150">
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm font-medium text-gray-900">{{ getSurgeryName(supply.surgery_id) }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {{ supply.supply_code }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ getSupplyName(supply.supply_code) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {{ supply.typical_quantity || 1 }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span v-if="supply.is_required" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">
                  Sí
                </span>
                <span v-else class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
                  No
                </span>
              </td>
              <td class="px-6 py-4 text-sm text-gray-500 max-w-xs truncate" :title="supply.notes">
                {{ supply.notes || '-' }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <div class="flex justify-end space-x-2">
                  <button @click="openEditModal(supply)" 
                    class="btn-primary text-xs px-3 py-1.5"
                    title="Editar">
                    <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                  </button>
                  <button @click="confirmDelete(supply)" 
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

    <!-- Mensaje sin resultados -->
    <div v-if="!loading && filteredTypicalSupplies.length === 0" class="card text-center py-12">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900">No hay insumos típicos</h3>
      <p class="mt-1 text-sm text-gray-500">{{ selectedSurgeryId || searchTerm ? 'No se encontraron resultados con los filtros aplicados.' : 'Comienza asociando insumos a cirugías.' }}</p>
      <div class="mt-6" v-if="!selectedSurgeryId && !searchTerm">
        <button @click="openCreateModal" class="btn-primary">
          <svg class="h-5 w-5 mr-2 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Crear Asociación
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
                {{ isEditing ? 'Editar Insumo Típico' : 'Crear Insumo Típico' }}
              </h3>
              <button @click="closeModal" class="text-gray-400 hover:text-gray-600">
                <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            
            <form @submit.prevent="saveTypicalSupply" class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Cirugía <span class="text-red-500">*</span>
                </label>
                <select v-model="typicalSupplyForm.surgery_id" class="form-select" required :disabled="isEditing">
                  <option value="">Seleccione una cirugía</option>
                  <option v-for="surgery in surgeries" :key="surgery.id" :value="surgery.id">
                    {{ surgery.name }}
                  </option>
                </select>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Insumo <span class="text-red-500">*</span>
                </label>
                <select v-model="typicalSupplyForm.supply_code" class="form-select" required :disabled="isEditing">
                  <option value="">Seleccione un insumo</option>
                  <option v-for="supply in supplyCodes" :key="supply.code" :value="supply.code">
                    {{ supply.code }} - {{ supply.name }}
                  </option>
                </select>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Cantidad Típica <span class="text-red-500">*</span>
                </label>
                <input v-model.number="typicalSupplyForm.typical_quantity" type="number" min="1" 
                  class="form-input" placeholder="Ej: 5" required />
                <p class="mt-1 text-xs text-gray-500">Cantidad típica necesaria para esta cirugía</p>
              </div>

              <div>
                <label class="flex items-center space-x-2">
                  <input v-model="typicalSupplyForm.is_required" type="checkbox" class="rounded border-gray-300 text-blue-600 focus:ring-blue-500" />
                  <span class="text-sm font-medium text-gray-700">Insumo requerido</span>
                </label>
                <p class="mt-1 text-xs text-gray-500">Los insumos requeridos son obligatorios para la cirugía</p>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Notas
                </label>
                <textarea v-model="typicalSupplyForm.notes" rows="3" class="form-input" 
                  placeholder="Notas adicionales sobre este insumo para la cirugía"></textarea>
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
import surgeryTypicalSupplyService from '@/services/surgeryTypicalSupplyService'
import surgeryService from '@/services/surgeryService'
import supplyCodeService from '@/services/supplyCodeService'
import Swal from 'sweetalert2'

const typicalSupplies = ref([])
const surgeries = ref([])
const supplyCodes = ref([])
const loading = ref(false)
const error = ref(null)
const searchTerm = ref('')
const selectedSurgeryId = ref('')
const showModal = ref(false)
const isEditing = ref(false)
const saving = ref(false)

const typicalSupplyForm = ref({
  surgery_id: '',
  supply_code: '',
  typical_quantity: 1,
  is_required: false,
  notes: ''
})

let searchTimeout = null

// Computed para filtrar insumos
const filteredTypicalSupplies = computed(() => {
  let filtered = [...typicalSupplies.value]

  // Filtrar por cirugía
  if (selectedSurgeryId.value) {
    filtered = filtered.filter(s => s.surgery_id === selectedSurgeryId.value)
  }

  // Filtrar por búsqueda
  if (searchTerm.value.trim()) {
    const term = searchTerm.value.toLowerCase().trim()
    filtered = filtered.filter(s => {
      const supplyName = getSupplyName(s.supply_code).toLowerCase()
      return supplyName.includes(term) || String(s.supply_code).includes(term)
    })
  }

  return filtered
})

// Funciones auxiliares
const getSurgeryName = (surgeryId) => {
  const surgery = surgeries.value.find(s => s.id === surgeryId)
  return surgery ? surgery.name : `Cirugía #${surgeryId}`
}

const getSupplyName = (supplyCode) => {
  const supply = supplyCodes.value.find(s => s.code === supplyCode)
  return supply ? supply.name : `Insumo #${supplyCode}`
}

const loadTypicalSupplies = async () => {
  loading.value = true
  error.value = null
  try {
    // Necesitamos cargar todas las asociaciones, pero el servicio actual solo tiene métodos por cirugía
    // Por ahora, cargaremos todas las cirugías y luego sus insumos típicos
    const allSurgeries = await surgeryService.getAllSurgeries()
    
    const allSupplies = []
    for (const surgery of allSurgeries) {
      try {
        const supplies = await surgeryTypicalSupplyService.getTypicalSuppliesBySurgeryId(surgery.id)
        allSupplies.push(...supplies.map(s => ({ ...s, surgery_id: surgery.id })))
      } catch (err) {
        // Ignorar errores si no hay insumos para esta cirugía
        console.log(`No se encontraron insumos para cirugía ${surgery.id}`)
      }
    }
    
    typicalSupplies.value = allSupplies
  } catch (err) {
    error.value = err.message || 'Error al cargar insumos típicos'
    console.error('Error loading typical supplies:', err)
  } finally {
    loading.value = false
  }
}

const loadSurgeries = async () => {
  try {
    const data = await surgeryService.getAllSurgeries()
    surgeries.value = data
  } catch (err) {
    console.error('Error loading surgeries:', err)
  }
}

const loadSupplyCodes = async () => {
  try {
    const data = await supplyCodeService.getAllSupplyCodes()
    supplyCodes.value = data
  } catch (err) {
    console.error('Error loading supply codes:', err)
  }
}

const handleSearch = () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    // La búsqueda se hace en el computed
  }, 300)
}

const filterBySurgery = () => {
  // El filtro se aplica en el computed
}

const clearFilters = () => {
  selectedSurgeryId.value = ''
  searchTerm.value = ''
}

const openCreateModal = () => {
  isEditing.value = false
  typicalSupplyForm.value = {
    surgery_id: '',
    supply_code: '',
    typical_quantity: 1,
    is_required: false,
    notes: ''
  }
  showModal.value = true
}

const openEditModal = (supply) => {
  isEditing.value = true
  typicalSupplyForm.value = {
    id: supply.id,
    surgery_id: supply.surgery_id,
    supply_code: supply.supply_code,
    typical_quantity: supply.typical_quantity || 1,
    is_required: supply.is_required || false,
    notes: supply.notes || ''
  }
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  typicalSupplyForm.value = {
    surgery_id: '',
    supply_code: '',
    typical_quantity: 1,
    is_required: false,
    notes: ''
  }
}

const saveTypicalSupply = async () => {
  // Validaciones
  if (!typicalSupplyForm.value.surgery_id) {
    await Swal.fire({
      icon: 'warning',
      title: 'Campo requerido',
      text: 'Debe seleccionar una cirugía',
      confirmButtonText: 'Aceptar'
    })
    return
  }

  if (!typicalSupplyForm.value.supply_code) {
    await Swal.fire({
      icon: 'warning',
      title: 'Campo requerido',
      text: 'Debe seleccionar un insumo',
      confirmButtonText: 'Aceptar'
    })
    return
  }

  if (!typicalSupplyForm.value.typical_quantity || typicalSupplyForm.value.typical_quantity < 1) {
    await Swal.fire({
      icon: 'warning',
      title: 'Cantidad inválida',
      text: 'La cantidad típica debe ser al menos 1',
      confirmButtonText: 'Aceptar'
    })
    return
  }

  saving.value = true
  try {
    const supplyData = {
      surgery_id: parseInt(typicalSupplyForm.value.surgery_id),
      supply_code: parseInt(typicalSupplyForm.value.supply_code),
      typical_quantity: parseInt(typicalSupplyForm.value.typical_quantity),
      is_required: typicalSupplyForm.value.is_required || false,
      notes: typicalSupplyForm.value.notes.trim() || null
    }

    if (isEditing.value) {
      await surgeryTypicalSupplyService.updateTypicalSupply(typicalSupplyForm.value.id, supplyData)
      await loadTypicalSupplies()
      closeModal()
      await Swal.fire({
        icon: 'success',
        title: 'Actualizado',
        text: 'Insumo típico actualizado exitosamente',
        timer: 2000,
        showConfirmButton: false
      })
    } else {
      await surgeryTypicalSupplyService.createTypicalSupply(supplyData)
      await loadTypicalSupplies()
      closeModal()
      await Swal.fire({
        icon: 'success',
        title: 'Creado',
        text: 'Insumo típico creado exitosamente',
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

const confirmDelete = async (supply) => {
  const surgeryName = getSurgeryName(supply.surgery_id)
  const supplyName = getSupplyName(supply.supply_code)
  
  const result = await Swal.fire({
    title: '¿Estás seguro?',
    html: `¿Deseas eliminar la asociación entre <strong>"${surgeryName}"</strong> y <strong>"${supplyName}"</strong>?<br><small class="text-gray-600">Esta acción no se puede deshacer.</small>`,
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
      await surgeryTypicalSupplyService.deleteTypicalSupply(supply.id)
      await loadTypicalSupplies()
      await Swal.fire({
        icon: 'success',
        title: 'Eliminado',
        text: 'Asociación eliminada exitosamente',
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
    loadSurgeries(),
    loadSupplyCodes(),
    loadTypicalSupplies()
  ])
})
</script>

