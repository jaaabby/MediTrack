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
              <th scope="col" class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider sticky right-0 bg-gray-50 z-10 shadow-[-2px_0_4px_rgba(0,0,0,0.05)]">
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
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium sticky right-0 bg-white z-10 shadow-[-2px_0_4px_rgba(0,0,0,0.05)]">
                <div class="flex justify-end space-x-2">
                  <button @click="openEditModal(supply)" 
                    class="text-warning-600 hover:text-warning-800 hover:bg-warning-50 p-1.5 rounded inline-flex items-center gap-1 transition-colors"
                    title="Editar">
                    <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                    <span class="font-medium text-xs">Editar</span>
                  </button>
                  <button @click="confirmDelete(supply)" 
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

    <!-- Modal para seleccionar insumos -->
    <Teleport to="body">
      <div v-if="showSupplySelectionModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-[60]" @click.self="closeSupplySelectionModal">
        <div class="relative top-10 mx-auto p-5 border w-full max-w-4xl shadow-lg rounded-md bg-white max-h-[90vh] flex flex-col">
          <div class="flex justify-between items-center border-b pb-3 mb-4">
            <h3 class="text-xl font-semibold text-gray-900">Seleccionar Insumos</h3>
            <button @click="closeSupplySelectionModal" class="text-gray-400 hover:text-gray-600">
              <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- Búsqueda y Filtros -->
          <div class="mb-4 space-y-3">
            <div class="flex gap-3">
              <div class="flex-1">
                <input
                  type="text"
                  v-model="supplySearchTerm"
                  placeholder="Buscar por nombre o código..."
                  class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                />
              </div>
              <div class="w-48">
                <select v-model="supplySortBy" class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500">
                  <option value="name">Ordenar por Nombre</option>
                  <option value="code">Ordenar por Código</option>
                </select>
              </div>
            </div>
            <div class="flex items-center justify-between text-sm">
              <span class="text-gray-600">
                {{ filteredAndSortedSupplies.length }} insumo(s) disponible(s)
              </span>
              <div class="flex gap-2">
                <button
                  type="button"
                  @click="selectAllFilteredSupplies"
                  class="px-3 py-1 text-sm text-blue-600 hover:text-blue-800"
                >
                  Seleccionar Todos
                </button>
                <button
                  type="button"
                  @click="deselectAllSupplies"
                  class="px-3 py-1 text-sm text-gray-600 hover:text-gray-800"
                >
                  Deseleccionar Todos
                </button>
              </div>
            </div>
          </div>

          <!-- Lista de insumos -->
          <div class="flex-1 overflow-y-auto border border-gray-200 rounded-md">
            <div v-if="loadingSupplies" class="text-center py-8">
              <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
              <p class="mt-2 text-sm text-gray-600">Cargando insumos...</p>
            </div>
            <div v-else-if="filteredAndSortedSupplies.length === 0" class="text-center py-8 text-gray-500">
              <p>No se encontraron insumos</p>
            </div>
            <div v-else class="divide-y divide-gray-200">
              <label
                v-for="supply in filteredAndSortedSupplies"
                :key="supply.code"
                class="flex items-start p-4 hover:bg-gray-50 cursor-pointer"
              >
                <input
                  type="checkbox"
                  :value="supply.code"
                  v-model="tempSelectedSupplyCodes"
                  class="mt-1 h-4 w-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                />
                <div class="ml-3 flex-1">
                  <div class="flex items-start justify-between">
                    <div class="flex-1">
                      <div class="font-medium text-gray-900 text-base">{{ supply.name }}</div>
                      <div class="mt-1 flex flex-wrap gap-3 text-sm">
                        <div class="text-gray-600">
                          <span class="font-medium">Código:</span> {{ supply.code }}
                        </div>
                        <div v-if="supply.code_supplier" class="text-gray-600">
                          <span class="font-medium">Código Proveedor:</span> {{ supply.code_supplier }}
                        </div>
                        <div v-if="supply.critical_stock !== undefined && supply.critical_stock !== null" class="text-gray-600">
                          <span class="font-medium">Stock Crítico:</span> {{ supply.critical_stock }}
                        </div>
                      </div>
                    </div>
                    <div v-if="isSupplyAlreadySelected(supply.code)" class="ml-2 flex-shrink-0">
                      <span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-yellow-100 text-yellow-800">
                        Ya seleccionado
                      </span>
                    </div>
                  </div>
                </div>
              </label>
            </div>
          </div>

          <!-- Footer del modal -->
          <div class="flex justify-between items-center mt-4 pt-4 border-t">
            <div class="text-sm text-gray-600">
              <span v-if="tempSelectedSupplyCodes.length > 0">
                {{ tempSelectedSupplyCodes.length }} insumo(s) seleccionado(s)
              </span>
              <span v-else>Ningún insumo seleccionado</span>
            </div>
            <div class="flex gap-3">
              <button
                type="button"
                @click="closeSupplySelectionModal"
                class="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50"
              >
                Cancelar
              </button>
              <button
                type="button"
                @click="confirmSupplySelection"
                class="px-4 py-2 bg-blue-600 text-white rounded-md text-sm font-medium hover:bg-blue-700"
                :disabled="tempSelectedSupplyCodes.length === 0"
              >
                Confirmar Selección ({{ tempSelectedSupplyCodes.length }})
              </button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

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
                  Insumos <span class="text-red-500">*</span>
                </label>
                <button
                  type="button"
                  @click="openSupplySelectionModal"
                  class="w-full px-4 py-2 border border-gray-300 rounded-md bg-white text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 mb-3"
                  :disabled="isEditing || !typicalSupplyForm.surgery_id"
                >
                  <span v-if="selectedSupplies.length === 0">Seleccionar Insumos</span>
                  <span v-else>{{ selectedSupplies.length }} insumo(s) seleccionado(s)</span>
                </button>
                <p v-if="!typicalSupplyForm.surgery_id" class="mt-1 text-xs text-gray-500 mb-3">
                  Primero seleccione una cirugía
                </p>
                <p v-else-if="selectedSupplies.length === 0" class="mt-1 text-xs text-gray-500 mb-3">
                  Haga clic para seleccionar uno o más insumos
                </p>

                <!-- Lista de insumos seleccionados con sus configuraciones -->
                <div v-if="selectedSupplies.length > 0" class="space-y-4 mt-4">
                  <div 
                    v-for="(supply, index) in selectedSupplies" 
                    :key="supply.code"
                    class="p-4 bg-blue-50 border border-blue-200 rounded-lg"
                  >
                    <div class="flex items-start justify-between mb-3">
                      <div class="flex-1">
                        <div class="font-medium text-base text-gray-900">{{ supply.name }}</div>
                        <div class="text-sm text-gray-600 mt-1">Código: {{ supply.code }}</div>
                      </div>
                      <button
                        v-if="!isEditing"
                        type="button"
                        @click="removeSelectedSupply(index)"
                        class="ml-2 text-red-500 hover:text-red-700 flex-shrink-0"
                        title="Eliminar insumo"
                      >
                        <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                        </svg>
                      </button>
                    </div>

                    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <!-- Cantidad Típica -->
                      <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">
                          Cantidad Típica <span class="text-red-500">*</span>
                        </label>
                        <input 
                          v-model.number="supply.typical_quantity" 
                          type="number" 
                          min="1" 
                          class="form-input" 
                          placeholder="Ej: 5" 
                          required
                        />
                        <p class="mt-1 text-xs text-gray-500">Cantidad típica para este insumo</p>
                      </div>

                      <!-- Insumo Requerido -->
                      <div class="flex items-end">
                        <label class="flex items-center space-x-2 cursor-pointer">
                          <input 
                            v-model="supply.is_required" 
                            type="checkbox" 
                            class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                          />
                          <span class="text-sm font-medium text-gray-700">Insumo requerido</span>
                        </label>
                      </div>
                    </div>

                    <!-- Notas por insumo -->
                    <div class="mt-3">
                      <label class="block text-sm font-medium text-gray-700 mb-1">
                        Notas (opcional)
                      </label>
                      <textarea 
                        v-model="supply.notes" 
                        rows="2" 
                        class="form-input text-sm" 
                        placeholder="Notas específicas para este insumo..."
                      ></textarea>
                    </div>
                  </div>
                </div>
              </div>

              <div class="flex justify-end space-x-3 pt-4 border-t">
                <button type="button" @click="closeModal" class="btn-secondary">Cancelar</button>
                <button 
                  v-if="!isEditing" 
                  type="button" 
                  @click="saveTypicalSupplyAndAddAnother" 
                  :disabled="saving" 
                  class="btn-secondary"
                >
                  <span v-if="saving">Guardando...</span>
                  <span v-else>Crear y Agregar Más</span>
                </button>
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
import surgeryTypicalSupplyService from '@/services/management/surgeryTypicalSupplyService'
import surgeryService from '@/services/management/surgeryService'
import supplyCodeService from '@/services/config/supplyCodeService'
import { useNotification } from '@/composables/useNotification'
import { useAlert } from '@/composables/useAlert'

const { success: showSuccess, error: showError, warning: showWarning } = useNotification()
const { confirmDanger } = useAlert()

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

// Estados para el modal de selección de insumos
const showSupplySelectionModal = ref(false)
const loadingSupplies = ref(false)
const supplySearchTerm = ref('')
const supplySortBy = ref('name')
const tempSelectedSupplyCodes = ref([])
const selectedSupplies = ref([]) // Insumos seleccionados para el formulario

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

// Computed para filtrar y ordenar insumos en el modal de selección
const filteredAndSortedSupplies = computed(() => {
  let filtered = [...supplyCodes.value]

  // Filtrar por búsqueda
  if (supplySearchTerm.value.trim()) {
    const term = supplySearchTerm.value.toLowerCase().trim()
    filtered = filtered.filter(supply => {
      const name = (supply.name || '').toLowerCase()
      const code = String(supply.code || '')
      return name.includes(term) || code.includes(term)
    })
  }

  // Ordenar
  if (supplySortBy.value === 'name') {
    filtered.sort((a, b) => {
      const nameA = (a.name || '').toLowerCase()
      const nameB = (b.name || '').toLowerCase()
      return nameA.localeCompare(nameB)
    })
  } else if (supplySortBy.value === 'code') {
    filtered.sort((a, b) => {
      return (a.code || 0) - (b.code || 0)
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
    // Usar el nuevo endpoint optimizado que obtiene todos los insumos típicos de una vez
    const allSupplies = await surgeryTypicalSupplyService.getAllTypicalSupplies()
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
  resetForm()
  showModal.value = true
}

const resetForm = () => {
  typicalSupplyForm.value = {
    surgery_id: '',
    supply_code: '',
    typical_quantity: 1,
    is_required: false,
    notes: ''
  }
  selectedSupplies.value = []
}

const openSupplySelectionModal = () => {
  if (!typicalSupplyForm.value.surgery_id) {
    showWarning('Debe seleccionar una cirugía antes de elegir insumos')
    return
  }
  
  // Inicializar con los insumos ya seleccionados
  tempSelectedSupplyCodes.value = selectedSupplies.value.map(s => s.code)
  supplySearchTerm.value = ''
  supplySortBy.value = 'name'
  showSupplySelectionModal.value = true
}

const closeSupplySelectionModal = () => {
  showSupplySelectionModal.value = false
  tempSelectedSupplyCodes.value = []
  supplySearchTerm.value = ''
}

const confirmSupplySelection = () => {
  // Convertir códigos seleccionados a objetos de insumos con configuración inicial
  const newSelectedCodes = tempSelectedSupplyCodes.value.filter(code => 
    !selectedSupplies.value.some(s => s.code === code)
  )
  
  const newSupplies = newSelectedCodes
    .map(code => {
      const supply = supplyCodes.value.find(s => s.code === code)
      if (supply) {
        return {
          ...supply,
          typical_quantity: 1, // Valor por defecto
          is_required: false,
          notes: ''
        }
      }
      return null
    })
    .filter(s => s !== null)
  
  // Mantener los insumos ya seleccionados y agregar los nuevos
  selectedSupplies.value = [...selectedSupplies.value, ...newSupplies]
  closeSupplySelectionModal()
}

const removeSelectedSupply = (index) => {
  selectedSupplies.value.splice(index, 1)
}

const isSupplyAlreadySelected = (code) => {
  return selectedSupplies.value.some(s => s.code === code)
}

const selectAllFilteredSupplies = () => {
  const filteredCodes = filteredAndSortedSupplies.value.map(s => s.code)
  // Agregar solo los que no están ya seleccionados
  filteredCodes.forEach(code => {
    if (!tempSelectedSupplyCodes.value.includes(code)) {
      tempSelectedSupplyCodes.value.push(code)
    }
  })
}

const deselectAllSupplies = () => {
  tempSelectedSupplyCodes.value = []
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
  
  // En modo edición, cargar el insumo seleccionado con su configuración
  const supplyInfo = supplyCodes.value.find(s => s.code === supply.supply_code)
  if (supplyInfo) {
    selectedSupplies.value = [{
      ...supplyInfo,
      typical_quantity: supply.typical_quantity || 1,
      is_required: supply.is_required || false,
      notes: supply.notes || ''
    }]
  } else {
    selectedSupplies.value = []
  }
  
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  resetForm()
}

const validateForm = () => {
  if (!typicalSupplyForm.value.surgery_id) {
    return { valid: false, message: 'Debe seleccionar una cirugía' }
  }

  if (selectedSupplies.value.length === 0) {
    return { valid: false, message: 'Debe seleccionar al menos un insumo' }
  }

  // Validar que cada insumo tenga una cantidad típica válida
  for (const supply of selectedSupplies.value) {
    if (!supply.typical_quantity || supply.typical_quantity < 1) {
      return { valid: false, message: `El insumo "${supply.name}" debe tener una cantidad típica de al menos 1` }
    }
  }

  return { valid: true }
}

const saveTypicalSupply = async (closeAfterSave = true) => {
  // Validaciones
  const validation = validateForm()
  if (!validation.valid) {
    showWarning(validation.message)
    return
  }

  saving.value = true
  try {
    if (isEditing.value) {
      // Modo edición: solo actualizar el insumo existente
      // Tomar los valores desde selectedSupplies que es donde el usuario edita
      const editedSupply = selectedSupplies.value[0]
      
      const supplyData = {
        surgery_id: parseInt(typicalSupplyForm.value.surgery_id),
        supply_code: parseInt(typicalSupplyForm.value.supply_code),
        typical_quantity: parseInt(editedSupply.typical_quantity),
        is_required: editedSupply.is_required || false,
        notes: (editedSupply.notes || '').trim() || null
      }
      
      await surgeryTypicalSupplyService.updateTypicalSupply(typicalSupplyForm.value.id, supplyData)
      await loadTypicalSupplies()
      if (closeAfterSave) {
        closeModal()
      }
      showSuccess('Insumo típico actualizado exitosamente')
    } else {
      // Modo creación: crear múltiples insumos si hay varios seleccionados
      if (selectedSupplies.value.length === 0) {
        showWarning('Debe seleccionar al menos un insumo')
        saving.value = false
        return
      }

      const surgeryId = parseInt(typicalSupplyForm.value.surgery_id)

      // Crear todas las asociaciones con sus configuraciones individuales
      let successCount = 0
      let errorCount = 0
      const errors = []

      for (const supply of selectedSupplies.value) {
        try {
          const supplyData = {
            surgery_id: surgeryId,
            supply_code: parseInt(supply.code),
            typical_quantity: parseInt(supply.typical_quantity) || 1,
            is_required: supply.is_required || false,
            notes: (supply.notes || '').trim() || null
          }
          
          await surgeryTypicalSupplyService.createTypicalSupply(supplyData)
          successCount++
        } catch (err) {
          errorCount++
          const errorMsg = err.response?.data?.error || err.message || 'Error desconocido'
          errors.push(`${supply.name} (${supply.code}): ${errorMsg}`)
        }
      }

      await loadTypicalSupplies()

      if (errorCount === 0) {
        showSuccess(`${successCount} insumo(s) asociado(s) exitosamente`)
      } else {
        showWarning(`Se crearon ${successCount} asociación(es), pero ${errorCount} fallaron. Errores: ${errors.join('; ')}`)
      }
      
      if (closeAfterSave) {
        closeModal()
      } else {
        // Limpiar el formulario pero mantener la cirugía seleccionada
        const savedSurgeryId = typicalSupplyForm.value.surgery_id
        resetForm()
        typicalSupplyForm.value.surgery_id = savedSurgeryId
      }
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

const saveTypicalSupplyAndAddAnother = async () => {
  await saveTypicalSupply(false) // No cerrar el modal
}

const confirmDelete = async (supply) => {
  const surgeryName = getSurgeryName(supply.surgery_id)
  const supplyName = getSupplyName(supply.supply_code)
  
  const confirmed = await confirmDanger(
    `¿Deseas eliminar la asociación entre "${surgeryName}" y "${supplyName}"?\n\nEsta acción no se puede deshacer.`,
    'Confirmar eliminación'
  )
  if (!confirmed) {
    return
  }

  try {
    await surgeryTypicalSupplyService.deleteTypicalSupply(supply.id)
    await loadTypicalSupplies()
    showSuccess('Asociación eliminada exitosamente')
  } catch (err) {
    console.error('Error al eliminar:', err)
    showError('Error al eliminar: ' + (err.response?.data?.error || err.message))
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

