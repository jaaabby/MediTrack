<template>
  <div class="space-y-6">
    <!-- Header del inventario -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Inventario de Insumos Médicos</h1>
        <p class="text-gray-600 mt-1">Gestión y control de stock médico</p>
      </div>
      <router-link to="/add-supply" class="btn-success">
        <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        Agregar Insumo
      </router-link>
    </div>

    <!-- Filtros y búsqueda -->
    <div class="card">
      <div class="flex items-end space-x-4">
        <!-- Buscador único -->
        <div class="flex-1">
          <label class="block text-sm font-medium text-gray-700 mb-2">Buscar insumo</label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            <input
              type="text"
              placeholder="Buscar por número de lote, nombre, código o proveedor..."
              class="form-input pl-10 w-full"
              v-model="searchTerm"
            />
          </div>
        </div>

        <!-- Botón de limpiar búsqueda -->
        <div>
          <button 
            class="btn-secondary px-4 py-2 h-10" 
            @click="clearSearch"
            :disabled="!searchTerm"
          >
            Limpiar
          </button>
        </div>
      </div>
    </div>

    <!-- Tabla de inventario -->
    <div class="card">
      <div class="card-header">
        <h2 class="card-title">Inventario de Insumos Médicos</h2>
        <p class="text-sm text-gray-600">Total: {{ filteredSupplies.length }} insumos</p>
      </div>

      <!-- Indicador de carga -->
      <div v-if="loading" class="flex justify-center items-center py-8">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
        <span class="ml-2 text-gray-600">Cargando inventario...</span>
      </div>

      <!-- Mensaje de error -->
      <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-md p-4 mx-4 mb-4">
        <div class="flex">
          <div class="flex-shrink-0">
            <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-red-800">Error al cargar inventario</h3>
            <div class="mt-2 text-sm text-red-700">{{ error }}</div>
            <div class="mt-4">
              <button @click="loadInventory" class="btn-secondary text-sm">
                Reintentar
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Tabla de datos -->
      <div v-else class="overflow-x-auto">
        <table class="table">
          <thead class="table-header">
            <tr>
              <th class="table-header-cell">
                <div class="flex items-center justify-between">
                  <span>N° de lote</span>
                  <div class="flex flex-col ml-2">
                    <button 
                      @click="sortBy('batch_id', 'asc')" 
                      class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-primary-600': sortField === 'batch_id' && sortDirection === 'asc' }"
                    >
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                      </svg>
                    </button>
                    <button 
                      @click="sortBy('batch_id', 'desc')" 
                      class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-primary-600': sortField === 'batch_id' && sortDirection === 'desc' }"
                    >
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                      </svg>
                    </button>
                  </div>
                </div>
              </th>
              <th class="table-header-cell">
                <div class="flex items-center justify-between">
                  <span>Nombre del Insumo</span>
                  <div class="flex flex-col ml-2">
                    <button 
                      @click="sortBy('name', 'asc')" 
                      class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-primary-600': sortField === 'name' && sortDirection === 'asc' }"
                    >
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                      </svg>
                    </button>
                    <button 
                      @click="sortBy('name', 'desc')" 
                      class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-primary-600': sortField === 'name' && sortDirection === 'desc' }"
                    >
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                      </svg>
                    </button>
                  </div>
                </div>
              </th>
              <th class="table-header-cell">
                <div class="flex items-center justify-between">
                  <span>Código Interno</span>
                  <div class="flex flex-col ml-2">
                    <button 
                      @click="sortBy('code', 'asc')" 
                      class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-primary-600': sortField === 'code' && sortDirection === 'asc' }"
                    >
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                      </svg>
                    </button>
                    <button 
                      @click="sortBy('code', 'desc')" 
                      class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-primary-600': sortField === 'code' && sortDirection === 'desc' }"
                    >
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                      </svg>
                    </button>
                  </div>
                </div>
              </th>
              <th class="table-header-cell">
                <div class="flex items-center justify-between">
                  <span>F. Vencimiento</span>
                  <div class="flex flex-col ml-2">
                    <button 
                      @click="sortBy('expiration_date', 'asc')" 
                      class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-primary-600': sortField === 'expiration_date' && sortDirection === 'asc' }"
                    >
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                      </svg>
                    </button>
                    <button 
                      @click="sortBy('expiration_date', 'desc')" 
                      class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-primary-600': sortField === 'expiration_date' && sortDirection === 'desc' }"
                    >
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                      </svg>
                    </button>
                  </div>
                </div>
              </th>
              <th class="table-header-cell">
                <div class="flex items-center justify-between">
                  <span>Cantidad</span>
                  <div class="flex flex-col ml-2">
                    <button 
                      @click="sortBy('amount', 'asc')" 
                      class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-primary-600': sortField === 'amount' && sortDirection === 'asc' }"
                    >
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                      </svg>
                    </button>
                    <button 
                      @click="sortBy('amount', 'desc')" 
                      class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-primary-600': sortField === 'amount' && sortDirection === 'desc' }"
                    >
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                      </svg>
                    </button>
                  </div>
                </div>
              </th>
              <th class="table-header-cell">
                <div class="flex items-center justify-between">
                  <span>Proveedor</span>
                  <div class="flex flex-col ml-2">
                    <button 
                      @click="sortBy('supplier', 'asc')" 
                      class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-primary-600': sortField === 'supplier' && sortDirection === 'asc' }"
                    >
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                      </svg>
                    </button>
                    <button 
                      @click="sortBy('supplier', 'desc')" 
                      class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-primary-600': sortField === 'supplier' && sortDirection === 'desc' }"
                    >
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                      </svg>
                    </button>
                  </div>
                </div>
              </th>
              <th class="table-header-cell">Acciones</th>
            </tr>
          </thead>
          <tbody class="table-body">
            <tr v-for="supply in paginatedSupplies" :key="supply.batch_id" class="table-row">
              <td class="table-cell font-mono text-sm">{{ supply.batch_id }}</td>
              <td class="table-cell">
                <div>
                  <div class="font-medium text-gray-900">{{ supply.name }}</div>
                </div>
              </td>
              <td class="table-cell">
                <span class="text-gray-700">{{ supply.code }}</span>
              </td>
              <td class="table-cell">
                <span :class="getExpirationClass(supply.expiration_date)">
                  {{ formatDate(supply.expiration_date) }}
                </span>
              </td>
              <td class="table-cell">
                <span class="font-medium">{{ supply.amount }}</span>
                <span class="text-gray-500 text-sm ml-1">unidades</span>
              </td>
              <td class="table-cell">
                <span class="text-gray-700">{{ supply.supplier }}</span>
              </td>
              <td class="table-cell">
                <div class="flex space-x-2">
                  <button class="text-primary-600 hover:text-primary-800" @click="viewSupply(supply)">
                    <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                    </svg>
                  </button>
                  <button class="text-warning-600 hover:text-warning-800" @click="editSupply(supply)">
                    <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2.5 2.5 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                  </button>
                  <button class="text-danger-600 hover:text-danger-800" @click="deleteSupply(supply)">
                    <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Paginación -->
      <div class="flex items-center justify-between mt-6">
        <div class="text-sm text-gray-700">
          Mostrando {{ startIndex + 1 }} a {{ endIndex }} de {{ filteredSupplies.length }} resultados
        </div>
        <div class="flex space-x-2">
          <button
            class="btn-secondary"
            :disabled="currentPage === 1"
            @click="currentPage--"
          >
            Anterior
          </button>
          <span class="px-3 py-2 text-sm text-gray-700">
            Página {{ currentPage }} de {{ totalPages }}
          </span>
          <button
            class="btn-secondary"
            :disabled="currentPage === totalPages"
            @click="currentPage++"
          >
            Siguiente
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import inventoryService from '@/services/inventoryService'

const route = useRoute()

// Estado reactivo
const supplies = ref([])
const loading = ref(false)
const error = ref(null)
const searchTerm = ref('')
const sortField = ref('name')
const sortDirection = ref('asc')
const currentPage = ref(1)
const itemsPerPage = 10

// Computed properties
const filteredSupplies = computed(() => {
  let filtered = [...supplies.value]
  
  if (searchTerm.value) {
    filtered = filtered.filter(supply => 
      supply.code.toString().includes(searchTerm.value) ||
      supply.name.toLowerCase().includes(searchTerm.value.toLowerCase()) ||
      supply.supplier.toLowerCase().includes(searchTerm.value.toLowerCase())
    )
  }
  
  // Ordenamiento
  filtered.sort((a, b) => {
    let result = 0
    
    switch (sortField.value) {
      case 'batch_id':
        result = a.batch_id - b.batch_id
        break
      case 'code':
        result = a.code - b.code
        break
      case 'expiration_date':
        result = new Date(a.expiration_date) - new Date(b.expiration_date)
        break
      case 'amount':
        result = a.amount - b.amount
        break
      case 'supplier':
        result = (a.supplier || '').localeCompare(b.supplier || '')
        break
      default:
        result = a.name.localeCompare(b.name)
        break
    }
    
    return sortDirection.value === 'asc' ? result : -result
  })
  
  return filtered
})

const totalPages = computed(() => Math.ceil(filteredSupplies.value.length / itemsPerPage))

const startIndex = computed(() => (currentPage.value - 1) * itemsPerPage)
const endIndex = computed(() => Math.min(startIndex.value + itemsPerPage, filteredSupplies.value.length))

const paginatedSupplies = computed(() => {
  return filteredSupplies.value.slice(startIndex.value, endIndex.value)
})

// Métodos
const clearSearch = () => {
  searchTerm.value = ''
  sortField.value = 'name'
  sortDirection.value = 'asc'
  currentPage.value = 1
}

const sortBy = (field, direction) => {
  sortField.value = field
  sortDirection.value = direction
  currentPage.value = 1
}

const formatDate = (dateString) => {
  try {
    return format(new Date(dateString), 'dd/MM/yyyy', { locale: es })
  } catch {
    return dateString
  }
}

const getExpirationClass = (expirationDate) => {
  const today = new Date()
  const expDate = new Date(expirationDate)
  const daysUntilExpiration = Math.ceil((expDate - today) / (1000 * 60 * 60 * 24))
  
  if (daysUntilExpiration < 0) return 'text-danger-600 font-semibold'
  if (daysUntilExpiration <= 30) return 'text-warning-600 font-semibold'
  return 'text-gray-900'
}

const viewSupply = (supply) => {
  console.log('Ver insumo:', supply)
  // TODO: Implementar vista detallada
}

const editSupply = (supply) => {
  console.log('Editar insumo:', supply)
  // TODO: Implementar edición
}

const deleteSupply = (supply) => {
  if (confirm(`¿Está seguro de que desea eliminar el insumo ${supply.name}?`)) {
    console.log('Eliminar insumo:', supply)
    // TODO: Implementar eliminación usando supply.batch_id
  }
}

// Métodos
const loadInventory = async () => {
  loading.value = true
  error.value = null
  
  try {
    const data = await inventoryService.getInventory()
    supplies.value = data
  } catch (err) {
    error.value = 'Error al cargar el inventario: ' + err.message
    console.error('Error al cargar inventario:', err)
  } finally {
    loading.value = false
  }
}

// Lifecycle
onMounted(() => {
  // Si viene con un término de búsqueda desde Home, aplicarlo
  if (route.query.search) {
    searchTerm.value = route.query.search
  }
  loadInventory()
})
</script>

<style scoped>
/* Transiciones suaves */
.table-row:hover {
  background-color: rgb(249 250 251);
  transition: background-color 0.2s ease-in-out;
}

/* Colores específicos para estados */
.text-danger-600 {
  color: #dc2626;
}

.text-warning-600 {
  color: #d97706;
}

.text-primary-600 {
  color: #2563eb;
}
</style>