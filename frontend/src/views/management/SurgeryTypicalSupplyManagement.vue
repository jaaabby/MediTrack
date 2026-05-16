<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h2 class="text-2xl font-semibold text-gray-900">Insumos Típicos por Cirugía</h2>
          <p class="text-gray-600 mt-1">Gestiona los insumos típicos asociados a cada tipo de cirugía</p>
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
    <FilterPanel
      :key="filterPanelKey"
      :filters="filterConfig"
      :result-count="filteredTypicalSupplies.length"
      :show-clear="false"
      @filter-change="onFilterChange"
    >
      <template #filter-surgery_search="{ setValue }">
        <div class="relative">
          <label class="block text-sm font-medium text-gray-700 mb-2">Filtrar por Cirugía</label>
          <input
            type="text"
            :value="surgerySearch"
            placeholder="Buscar cirugía..."
            class="form-input w-full"
            @input="setValue($event.target.value)"
            @focus="showSurgeryOptions = true"
            @blur="hideSurgeryOptions"
            autocomplete="off"
          />

          <div
            v-if="showSurgeryOptions && filteredSurgeries.length > 0"
            class="absolute z-50 w-full mt-1 bg-white border border-gray-300 rounded-md shadow-lg max-h-60 overflow-auto"
          >
            <button
              v-for="surgery in filteredSurgeries"
              :key="surgery.id"
              @mousedown.prevent="selectSurgery(surgery)"
              class="w-full text-left px-4 py-2 hover:bg-blue-50 cursor-pointer text-sm border-b last:border-b-0"
            >
              {{ surgery.name }}
            </button>
          </div>

          <button
            v-if="surgerySearch"
            @click="clearSurgeryFilter"
            class="absolute right-2 top-9 text-gray-400 hover:text-gray-600"
            type="button"
          >
            <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </template>

      <template #actions>
        <button
          class="btn-secondary text-sm disabled:opacity-40 disabled:cursor-not-allowed"
          @click="clearFilters"
          :disabled="!selectedSurgeryId && !searchTerm && !surgerySearch"
        >
          Limpiar filtros
        </button>
      </template>
    </FilterPanel>

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
    <DataTable
      v-if="!loading && !error"
      :columns="tableColumns"
      :rows="filteredTypicalSupplies"
      default-sort-key="supply_code"
      empty-message="">
      <template #cell-surgery_id="{ row }">
        <div class="text-sm font-medium text-gray-900 truncate max-w-xs" :title="getSurgeryName(row.surgery_id)">{{ getSurgeryName(row.surgery_id) }}</div>
      </template>
      <template #cell-supply_name="{ row }">
        {{ getSupplyName(row.supply_code) }}
      </template>
      <template #cell-typical_quantity="{ row }">
        {{ row.typical_quantity || 1 }}
      </template>
      <template #cell-is_required="{ row }">
        <span v-if="row.is_required" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">Sí</span>
        <span v-else class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">No</span>
      </template>
      <template #cell-notes="{ row }">
        <span class="max-w-xs truncate block" :title="row.notes">{{ row.notes || '-' }}</span>
      </template>
      <template #empty>
        <div class="text-center py-12">
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
      </template>
      <template #actions="{ row }">
        <button @click="openEditModal(row)"
          class="text-warning-600 hover:text-warning-800 hover:bg-warning-50 p-1.5 rounded inline-flex items-center gap-1 transition-colors"
          title="Editar">
          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
          </svg>
        </button>
        <button @click="confirmDelete(row)"
          class="text-danger-600 hover:text-danger-800 hover:bg-danger-50 p-1.5 rounded transition-colors"
          title="Eliminar">
          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
        </button>
      </template>
    </DataTable>

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
            <div class="flex flex-col sm:flex-row gap-3">
              <div class="flex-1">
                <input
                  type="text"
                  v-model="supplySearchTerm"
                  placeholder="Buscar por nombre o código..."
                  class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                />
              </div>
              <div class="shrink-0">
                <select v-model="supplySortBy" class="px-4 py-2 pr-8 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500">
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
              <div
                v-for="supply in filteredAndSortedSupplies"
                :key="supply.code"
                class="p-4"
                :class="isTempSelected(supply.code) ? 'bg-blue-50' : 'hover:bg-gray-50'"
              >
                <div class="flex items-start">
                  <input
                    type="checkbox"
                    :checked="isTempSelected(supply.code)"
                    @change="toggleTempSupply(supply)"
                    class="mt-1 h-4 w-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500 cursor-pointer flex-shrink-0"
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
                    </div>

                    <!-- Cantidad y notas inline al seleccionar -->
                    <div v-if="isTempSelected(supply.code)" class="mt-3 grid grid-cols-1 sm:grid-cols-2 gap-3">
                      <div>
                        <label class="block text-xs font-medium text-gray-700 mb-1">
                          Cantidad Típica <span class="text-red-500">*</span>
                        </label>
                        <input
                          type="number"
                          min="1"
                          :value="getTempSupply(supply.code).typical_quantity"
                          @input="getTempSupply(supply.code).typical_quantity = parseInt($event.target.value) || 1"
                          class="w-full px-3 py-1.5 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                          placeholder="Ej: 5"
                          @click.stop
                        />
                      </div>
                      <div>
                        <label class="block text-xs font-medium text-gray-700 mb-1">Notas (opcional)</label>
                        <input
                          type="text"
                          :value="getTempSupply(supply.code).notes"
                          @input="getTempSupply(supply.code).notes = $event.target.value"
                          class="w-full px-3 py-1.5 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                          placeholder="Notas para este insumo..."
                          @click.stop
                        />
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Footer del modal -->
          <div class="flex justify-between items-center mt-4 pt-4 border-t">
            <div class="text-sm text-gray-600">
              <span v-if="tempSelectedSupplies.length > 0">
                {{ tempSelectedSupplies.length }} insumo(s) seleccionado(s)
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
                :disabled="tempSelectedSupplies.length === 0"
              >
                Confirmar Selección ({{ tempSelectedSupplies.length }})
              </button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Modal para crear/editar -->
    <Teleport to="body">
      <div v-if="showModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 flex items-start justify-center z-50 p-4 pt-10" @click.self="closeModal">
        <div class="w-full max-w-lg shadow-lg rounded-md bg-white border flex flex-col max-h-[90vh]"
             :class="showModalSurgeryOptions ? 'min-h-[26rem]' : ''"
        >

          <!-- Header fijo -->
          <div class="flex justify-between items-center border-b px-5 pt-5 pb-3 flex-shrink-0">
            <h3 class="text-xl font-semibold text-gray-900">
              {{ isEditing ? 'Editar Insumo Típico' : 'Crear Insumo Típico' }}
            </h3>
            <button @click="closeModal" class="text-gray-400 hover:text-gray-600">
              <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- Contenido scrolleable -->
          <div class="overflow-y-auto flex-1 px-5 py-4 space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Cirugía <span class="text-red-500">*</span>
              </label>
              <div class="relative">
                <input
                  type="text"
                  v-model="modalSurgerySearch"
                  :placeholder="isEditing ? getSurgeryName(typicalSupplyForm.surgery_id) : 'Buscar cirugía...' "
                  :disabled="isEditing"
                  class="form-input w-full"
                  :class="{ 'border-red-500': formErrors.surgery_id }"
                  autocomplete="off"
                  @input="onModalSurgeryInput"
                  @focus="showModalSurgeryOptions = true"
                  @blur="hideModalSurgeryOptions"
                />
                <div
                  v-if="showModalSurgeryOptions && modalFilteredSurgeries.length > 0"
                  class="absolute z-[70] w-full mt-1 bg-white border border-gray-300 rounded-md shadow-lg max-h-52 overflow-auto"
                >
                  <button
                    v-for="surgery in modalFilteredSurgeries"
                    :key="surgery.id"
                    type="button"
                    @mousedown.prevent="selectModalSurgery(surgery)"
                    class="w-full text-left px-4 py-2 hover:bg-blue-50 cursor-pointer text-sm border-b last:border-b-0"
                  >
                    {{ surgery.name }}
                  </button>
                </div>
              </div>
              <p v-if="formErrors.surgery_id" class="mt-1 text-sm text-red-600">{{ formErrors.surgery_id }}</p>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Insumos <span class="text-red-500">*</span>
              </label>
              <button
                type="button"
                @click="openSupplySelectionModal"
                class="w-full px-4 py-2 border rounded-md bg-white text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 mb-3"
                :class="formErrors.supplies ? 'border-red-500' : 'border-gray-300'"
                :disabled="isEditing || !typicalSupplyForm.surgery_id"
              >
                <span v-if="selectedSupplies.length === 0">Seleccionar Insumos</span>
                <span v-else>{{ selectedSupplies.length }} insumo(s) seleccionado(s)</span>
              </button>
              <p v-if="formErrors.supplies" class="text-sm text-red-600 mb-2">{{ formErrors.supplies }}</p>
              <p v-if="!typicalSupplyForm.surgery_id" class="mt-1 text-xs text-gray-500 mb-3">
                Primero seleccione una cirugía
              </p>
              <p v-else-if="selectedSupplies.length === 0 && !formErrors.supplies" class="mt-1 text-xs text-gray-500 mb-3">
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
                        :class="{ 'border-red-500': formErrors.supply_quantities[supply.code] }"
                        placeholder="Ej: 5"
                        @input="delete formErrors.supply_quantities[supply.code]"
                      />
                      <p v-if="formErrors.supply_quantities[supply.code]" class="mt-1 text-sm text-red-600">{{ formErrors.supply_quantities[supply.code] }}</p>
                      <p v-else class="mt-1 text-xs text-gray-500">Cantidad típica para este insumo</p>
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
          </div>

          <!-- Footer fijo -->
          <div class="flex justify-end space-x-3 px-5 py-4 border-t flex-shrink-0">
            <button type="button" @click="closeModal" class="btn-secondary">Cancelar</button>
            <button type="button" @click="saveTypicalSupply" :disabled="saving" class="btn-primary">
              <span v-if="saving">Guardando...</span>
              <span v-else>{{ isEditing ? 'Actualizar' : 'Crear' }}</span>
            </button>
          </div>

        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import FilterPanel from '@/components/common/FilterPanel.vue'
import DataTable from '@/components/common/DataTable.vue'
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
const filterPanelKey = ref(0)
const showModal = ref(false)
const isEditing = ref(false)
const saving = ref(false)

// Estado de errores de validación
const formErrors = ref({
  surgery_id: '',
  supplies: '',
  supply_quantities: {}
})

const tableColumns = [
  { key: 'surgery_id', label: 'Cirugía', sortable: false },
  { key: 'supply_code', label: 'Código Insumo' },
  { key: 'supply_name', label: 'Nombre Insumo', sortable: false },
  { key: 'typical_quantity', label: 'Cantidad Típica' },
  { key: 'is_required', label: 'Requerido', sortable: false },
  { key: 'notes', label: 'Notas', sortable: false }
]

// Estados para autocomplete de cirugías (filtro)
const surgerySearch = ref('')
const showSurgeryOptions = ref(false)

// Estados para combobox de cirugía en el modal
const modalSurgerySearch = ref('')
const showModalSurgeryOptions = ref(false)

// Estados para el modal de selección de insumos
const showSupplySelectionModal = ref(false)
const loadingSupplies = ref(false)
const supplySearchTerm = ref('')
const supplySortBy = ref('name')
const tempSelectedSupplies = ref([]) // Array of { ...supply, typical_quantity, is_required, notes }
const selectedSupplies = ref([]) // Insumos seleccionados para el formulario

const typicalSupplyForm = ref({
  surgery_id: '',
  supply_code: '',
  typical_quantity: 1,
  is_required: false,
  notes: ''
})

let searchTimeout = null

// Función auxiliar para normalizar texto (sin tildes ni mayúsculas)
const normalizeText = (text) => {
  if (!text) return ''
  return text
    .toLowerCase()
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
}

const filterConfig = computed(() => [
  {
    type: 'custom',
    key: 'surgery_search',
    label: 'Filtrar por Cirugía',
    default: surgerySearch.value
  },
  {
    type: 'text',
    key: 'search_term',
    label: 'Buscar por Insumo',
    placeholder: 'Buscar por código o nombre de insumo...',
    default: searchTerm.value
  }
])

// Computed para cirugías filtradas en el modal
const modalFilteredSurgeries = computed(() => {
  if (!modalSurgerySearch.value.trim()) {
    return surgeries.value
  }
  const term = normalizeText(modalSurgerySearch.value.trim())
  return surgeries.value.filter(surgery => {
    const name = normalizeText(surgery.name || '')
    return name.includes(term)
  })
})

// Computed para cirugías filtradas (autocomplete)
const filteredSurgeries = computed(() => {
  if (!surgerySearch.value.trim()) {
    return surgeries.value
  }

  const term = normalizeText(surgerySearch.value.trim())
  return surgeries.value
    .filter(surgery => {
      const name = normalizeText(surgery.name || '')
      const id = String(surgery.id || '')
      return name.includes(term) || id.includes(term)
    })
})

// Computed para filtrar insumos
const filteredTypicalSupplies = computed(() => {
  let filtered = [...typicalSupplies.value]

  // Filtrar por cirugía
  if (selectedSurgeryId.value) {
    filtered = filtered.filter(s => s.surgery_id === selectedSurgeryId.value)
  }

  // Filtrar por búsqueda (sin distinguir mayúsculas ni tildes)
  if (searchTerm.value.trim()) {
    const term = normalizeText(searchTerm.value.trim())
    filtered = filtered.filter(s => {
      const supplyName = normalizeText(getSupplyName(s.supply_code))
      const supplyCode = String(s.supply_code)
      return supplyName.includes(term) || supplyCode.includes(term)
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
    console.log('🔄 Cargando insumos típicos...')
    const allSupplies = await surgeryTypicalSupplyService.getAllTypicalSupplies()
    console.log('✅ Insumos típicos cargados:', allSupplies)
    console.log('📊 Cantidad de insumos:', allSupplies ? allSupplies.length : 0)
    
    // Asegurar que siempre sea un array
    typicalSupplies.value = Array.isArray(allSupplies) ? allSupplies : []
    
    if (typicalSupplies.value.length === 0) {
      console.warn('⚠️ No se encontraron insumos típicos en la base de datos')
    }
  } catch (err) {
    console.error('❌ Error detallado al cargar insumos típicos:', err)
    console.error('❌ Código de error:', err.response?.status)
    console.error('❌ Mensaje del servidor:', err.response?.data)
    
    // Mostrar mensaje de error más descriptivo
    if (err.response?.status === 404) {
      error.value = 'Endpoint no encontrado. Verifique que el backend esté ejecutándose en el puerto correcto.'
    } else if (err.response?.status === 401 || err.response?.status === 403) {
      error.value = 'No tiene permisos para acceder a este recurso. Por favor, inicie sesión nuevamente.'
    } else {
      error.value = err.message || 'Error al cargar insumos típicos'
    }
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
    // filter is reactive via filteredTypicalSupplies computed
  }, 300)
}

const onFilterChange = (key, value) => {
  if (key === 'search_term') {
    searchTerm.value = value
    handleSearch()
    return
  }

  if (key === 'surgery_search') {
    surgerySearch.value = value

    if (!value.trim()) {
      selectedSurgeryId.value = ''
      showSurgeryOptions.value = false
      return
    }

    const selectedSurgery = surgeries.value.find(s => s.id === selectedSurgeryId.value)
    if (!selectedSurgery || selectedSurgery.name !== value) {
      onSurgerySearch()
    }
  }
}

const filterBySurgery = () => {
  // filter is reactive via filteredTypicalSupplies computed
}

const clearFilters = () => {
  selectedSurgeryId.value = ''
  searchTerm.value = ''
  surgerySearch.value = ''
  showSurgeryOptions.value = false
  filterPanelKey.value += 1
}

// Funciones para autocomplete de cirugías
const onSurgerySearch = () => {
  showSurgeryOptions.value = true
  selectedSurgeryId.value = '' // Limpiar selección cuando se escribe
}

const selectSurgery = (surgery) => {
  selectedSurgeryId.value = surgery.id
  surgerySearch.value = surgery.name
  showSurgeryOptions.value = false
}

const hideSurgeryOptions = () => {
  setTimeout(() => {
    showSurgeryOptions.value = false
  }, 200)
}

const clearSurgeryFilter = () => {
  selectedSurgeryId.value = ''
  surgerySearch.value = ''
}

const openCreateModal = () => {
  isEditing.value = false
  resetForm()
  showModal.value = true
}

const onModalSurgeryInput = () => {
  showModalSurgeryOptions.value = true
  // Clear selection if user edits text
  typicalSupplyForm.value.surgery_id = ''
  formErrors.value.surgery_id = ''
  selectedSupplies.value = []
}

const selectModalSurgery = (surgery) => {
  typicalSupplyForm.value.surgery_id = surgery.id
  modalSurgerySearch.value = surgery.name
  showModalSurgeryOptions.value = false
  formErrors.value.surgery_id = ''
  selectedSupplies.value = []
}

const hideModalSurgeryOptions = () => {
  setTimeout(() => {
    showModalSurgeryOptions.value = false
    // If no surgery was selected, clear the text
    if (!typicalSupplyForm.value.surgery_id) {
      modalSurgerySearch.value = ''
    }
  }, 200)
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
  modalSurgerySearch.value = ''
  showModalSurgeryOptions.value = false
  
  // Limpiar errores
  formErrors.value = {
    surgery_id: '',
    supplies: '',
    supply_quantities: {}
  }
}

const openSupplySelectionModal = () => {
  if (!typicalSupplyForm.value.surgery_id) {
    showWarning('Debe seleccionar una cirugía antes de elegir insumos')
    return
  }
  
  // Inicializar con los insumos ya seleccionados (con copia para no mutar los originales)
  tempSelectedSupplies.value = selectedSupplies.value.map(s => ({ ...s }))
  supplySearchTerm.value = ''
  supplySortBy.value = 'name'
  showSupplySelectionModal.value = true
}

const closeSupplySelectionModal = () => {
  showSupplySelectionModal.value = false
  tempSelectedSupplies.value = []
  supplySearchTerm.value = ''
}

const isTempSelected = (code) => tempSelectedSupplies.value.some(s => s.code === code)

const getTempSupply = (code) => tempSelectedSupplies.value.find(s => s.code === code)

const toggleTempSupply = (supply) => {
  const idx = tempSelectedSupplies.value.findIndex(s => s.code === supply.code)
  if (idx >= 0) {
    tempSelectedSupplies.value.splice(idx, 1)
  } else {
    // Preserve existing config if supply was already in selectedSupplies
    const existing = selectedSupplies.value.find(s => s.code === supply.code)
    tempSelectedSupplies.value.push(existing ? { ...existing } : {
      ...supply,
      typical_quantity: 1,
      is_required: false,
      notes: ''
    })
  }
}

const confirmSupplySelection = () => {
  // Limpiar error de insumos al seleccionar
  formErrors.value.supplies = ''
  selectedSupplies.value = tempSelectedSupplies.value.map(s => ({ ...s }))
  closeSupplySelectionModal()
}

const removeSelectedSupply = (index) => {
  selectedSupplies.value.splice(index, 1)
}

const isSupplyAlreadySelected = (code) => {
  return selectedSupplies.value.some(s => s.code === code)
}

const selectAllFilteredSupplies = () => {
  filteredAndSortedSupplies.value.forEach(supply => {
    if (!isTempSelected(supply.code)) {
      const existing = selectedSupplies.value.find(s => s.code === supply.code)
      tempSelectedSupplies.value.push(existing ? { ...existing } : {
        ...supply,
        typical_quantity: 1,
        is_required: false,
        notes: ''
      })
    }
  })
}

const deselectAllSupplies = () => {
  tempSelectedSupplies.value = []
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
  modalSurgerySearch.value = getSurgeryName(supply.surgery_id)
  
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
  // Resetear errores
  formErrors.value = {
    surgery_id: '',
    supplies: '',
    supply_quantities: {}
  }
  
  let hasErrors = false
  const errorMessages = []
  
  // Validar cirugía
  if (!typicalSupplyForm.value.surgery_id) {
    formErrors.value.surgery_id = 'La cirugía es obligatoria.'
    errorMessages.push('Debe seleccionar una cirugía')
    hasErrors = true
  }
  
  // Validar que haya al menos un insumo seleccionado
  if (selectedSupplies.value.length === 0) {
    formErrors.value.supplies = 'Debe seleccionar al menos un insumo.'
    errorMessages.push('Debe seleccionar al menos un insumo')
    hasErrors = true
  }
  
  // Validar que cada insumo tenga una cantidad típica válida
  for (const supply of selectedSupplies.value) {
    if (!supply.typical_quantity || supply.typical_quantity < 1) {
      formErrors.value.supply_quantities[supply.code] = 'La cantidad debe ser mayor a 0.'
      errorMessages.push(`El insumo "${supply.name}" debe tener una cantidad típica de al menos 1`)
      hasErrors = true
    }
  }
  
  return { valid: !hasErrors, errors: errorMessages }
}

const saveTypicalSupply = async (closeAfterSave = true) => {
  // Validaciones
  const validation = validateForm()
  if (!validation.valid) {
    // Mostrar todos los errores en una notificación
    const errorMessage = validation.errors.join('\n• ')
    showError('Errores en el formulario:\n• ' + errorMessage)
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
  
  // Inicializar surgerySearch si hay una cirugía seleccionada
  if (selectedSurgeryId.value) {
    const surgery = surgeries.value.find(s => s.id === selectedSurgeryId.value)
    if (surgery) {
      surgerySearch.value = surgery.name
    }
  }
})
</script>
