<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h2 class="text-2xl font-semibold text-gray-900">Gestión de Códigos de Insumos</h2>
          <p class="text-gray-600 mt-1">Gestiona los códigos de insumos y sus niveles críticos de stock</p>
        </div>
        <button @click="openCreateModal" class="btn-primary flex items-center justify-center">
          <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Nuevo Código
        </button>
      </div>
    </div>

    <!-- Búsqueda -->
    <FilterPanel :filters="[{ type: 'text', key: 'search', label: 'Buscar código de insumo', placeholder: 'Buscar por código o nombre...' }]" :result-count="filteredSupplyCodes.length" @filter-change="onFilterChange" />

    <!-- Loading -->
    <div v-if="loading" class="card">
      <div class="flex justify-center items-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        <span class="ml-3 text-gray-600">Cargando códigos de insumos...</span>
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
            <h3 class="text-sm font-medium text-red-800">Error al cargar códigos de insumos</h3>
            <div class="mt-2 text-sm text-red-700">{{ error }}</div>
            <button @click="loadSupplyCodes" class="btn-secondary mt-4 text-sm">Reintentar</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Tabla de códigos -->
    <DataTable
      v-else
      :columns="tableColumns"
      :rows="filteredSupplyCodes"
      default-sort-key="code"
      empty-message="No hay códigos de insumos registrados"
      :table-actions="[
        { type: 'edit', label: 'Editar', onClick: (row) => openEditModal(row) },
        { type: 'delete', onClick: (row) => confirmDelete(row) },
      ]"
    >
      <template #cell-code="{ row }">
        <div class="text-sm font-medium text-gray-900">{{ row.code }}</div>
      </template>
      <template #cell-name="{ row }">
        <div class="text-sm text-gray-900">{{ row.name }}</div>
      </template>
      <template #cell-code_supplier="{ row }">
        <div class="text-sm text-gray-900">{{ row.code_supplier }}</div>
      </template>
      <template #cell-critical_stock="{ row }">
        <span class="text-sm text-gray-900">{{ row.critical_stock }}</span>
      </template>
    </DataTable>

    <!-- Modal para crear/editar -->
    <Teleport to="body">
      <div v-if="showModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50" @click.self="closeModal">
        <div class="relative top-20 mx-auto p-5 border w-full max-w-lg shadow-lg rounded-md bg-white">
          <div class="space-y-4">
            <div class="flex justify-between items-center border-b pb-3">
              <h3 class="text-xl font-semibold text-gray-900">
                {{ isEditing ? 'Editar Código de Insumo' : 'Crear Código de Insumo' }}
              </h3>
              <button @click="closeModal" class="text-gray-400 hover:text-gray-600">
                <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            
            <form @submit.prevent="saveSupplyCode" class="space-y-4">
              <div v-if="!isEditing">
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Código del Insumo <span class="text-red-500">*</span>
                </label>
                <input v-model="supplyCodeForm.code" type="text" inputmode="numeric" class="form-input"
                  placeholder="123456" :disabled="isEditing"
                  @keydown="onlyNumericKeydown"
                  @input="onModalNumericInput($event, 'code')" />
                <p v-if="modalErrors.code" class="mt-1 text-xs text-red-600">{{ modalErrors.code }}</p>
                <p v-else class="mt-1 text-xs text-gray-500">Código único del insumo (no se puede modificar después)</p>
              </div>
              <div v-else>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Código del Insumo
                </label>
                <input type="number" class="form-input bg-gray-100" 
                  :value="supplyCodeForm.code" disabled />
                <p class="mt-1 text-xs text-gray-500">El código del insumo no se puede modificar</p>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Nombre del Insumo <span class="text-red-500">*</span>
                </label>
                <input v-model="supplyCodeForm.name" type="text" class="form-input"
                  placeholder="Ej: Jeringa 10ml"
                  @input="modalErrors.name = ''" />
                <p v-if="modalErrors.name" class="mt-1 text-xs text-red-600">{{ modalErrors.name }}</p>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Código de Proveedor <span class="text-red-500">*</span>
                </label>
                <input v-model="supplyCodeForm.code_supplier" type="text" inputmode="numeric" class="form-input"
                  placeholder="789"
                  @keydown="onlyNumericKeydown"
                  @input="onModalNumericInput($event, 'code_supplier')" />
                <p v-if="modalErrors.code_supplier" class="mt-1 text-xs text-red-600">{{ modalErrors.code_supplier }}</p>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Stock Crítico <span class="text-red-500">*</span>
                </label>
                <input v-model="supplyCodeForm.critical_stock" type="text" inputmode="numeric" class="form-input"
                  placeholder="1"
                  @keydown="onlyNumericKeydown"
                  @input="onModalNumericInput($event, 'critical_stock')" />
                <p v-if="modalErrors.critical_stock" class="mt-1 text-xs text-red-600">{{ modalErrors.critical_stock }}</p>
                <p v-else class="mt-1 text-xs text-gray-500">
                  Nivel mínimo de stock para generar alertas. Para insumos específicos, generalmente se usa 1.
                </p>
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
import supplyCodeService from '@/services/config/supplyCodeService'
import inventoryService from '@/services/inventory/inventoryService'
import { useNotification } from '@/composables/useNotification'
import { useAlert } from '@/composables/useAlert'
import FilterPanel from '@/components/common/FilterPanel.vue'
import DataTable from '@/components/common/DataTable.vue'

const { success: showSuccess, error: showError, warning: showWarning } = useNotification()
const { confirm, confirmDanger } = useAlert()

const supplyCodes = ref([])
const loading = ref(false)
const error = ref(null)
const filterState = reactive({ search: '' })
const showModal = ref(false)
const isEditing = ref(false)
const saving = ref(false)

// Bloquear caracteres no numéricos
const onlyNumericKeydown = (event) => {
  const allowed = ['Backspace', 'Delete', 'ArrowLeft', 'ArrowRight', 'ArrowUp', 'ArrowDown', 'Tab', 'Enter', 'Home', 'End']
  if (allowed.includes(event.key)) return
  if ((event.ctrlKey || event.metaKey) && ['a', 'c', 'v', 'x'].includes(event.key.toLowerCase())) return
  if (!/^\d$/.test(event.key)) event.preventDefault()
}

// Limpiar caracteres no numéricos tras pegar/autocompletar y limpiar error
const onModalNumericInput = (event, field) => {
  modalErrors.value[field] = ''
  const numeric = event.target.value.replace(/[^0-9]/g, '')
  if (event.target.value !== numeric) event.target.value = numeric
  supplyCodeForm.value[field] = numeric
}

const modalErrors = ref({
  code: '',
  name: '',
  code_supplier: '',
  critical_stock: ''
})

const tableColumns = [
  { key: 'code', label: 'Código' },
  { key: 'name', label: 'Nombre' },
  { key: 'code_supplier', label: 'Código Proveedor', sortable: false },
  { key: 'critical_stock', label: 'Stock Crítico' }
]
const supplyCodeForm = ref({
  code: null,
  name: '',
  code_supplier: null,
  critical_stock: 1
})

const onFilterChange = (key, value) => { filterState[key] = value }

const filteredSupplyCodes = computed(() => {
  if (!supplyCodes.value || supplyCodes.value.length === 0) return []
  if (!filterState.search.trim()) return supplyCodes.value
  const term = filterState.search.toLowerCase().trim()
  return supplyCodes.value.filter(sc =>
    String(sc.code).includes(term) ||
    sc.name.toLowerCase().includes(term) ||
    String(sc.code_supplier).includes(term)
  )
})

const loadSupplyCodes = async () => {
  loading.value = true
  error.value = null
  try {
    const data = await supplyCodeService.getAllSupplyCodes()
    supplyCodes.value = data
  } catch (err) {
    error.value = err.response?.data?.error || err.message || 'Error al cargar códigos de insumos'
    console.error('Error loading supply codes:', err)
  } finally {
    loading.value = false
  }
}

const openCreateModal = () => {
  isEditing.value = false
  supplyCodeForm.value = {
    code: null,
    name: '',
    code_supplier: null,
    critical_stock: 1
  }
  modalErrors.value = { code: '', name: '', code_supplier: '', critical_stock: '' }
  showModal.value = true
}

// Función para obtener el stock actual de un código de insumo
const getCurrentStock = async (supplyCode) => {
  try {
    const inventory = await inventoryService.getInventory()
    // Sumar las cantidades de todos los lotes para este código
    // El código puede estar en diferentes campos: code, supply_code
    const totalStock = inventory
      .filter(item => {
        const itemCode = item.code || item.supply_code
        return itemCode && Number(itemCode) === Number(supplyCode)
      })
      .reduce((sum, item) => {
        const amount = item.amount || item.current_amount || 0
        return sum + Number(amount)
      }, 0)
    return totalStock
  } catch (error) {
    console.error('Error al obtener stock actual:', error)
    // Si hay error, retornar null para que no bloquee la operación
    return null
  }
}

// Función para verificar si el stock está en nivel crítico
const isStockCritical = (currentStock, criticalStock) => {
  return currentStock !== null && currentStock <= criticalStock
}

const openEditModal = async (supplyCode) => {
  // Verificar stock crítico antes de abrir el modal
  const currentStock = await getCurrentStock(supplyCode.code)
  const isCritical = isStockCritical(currentStock, supplyCode.critical_stock)
  
  if (isCritical) {
    const confirmed = await confirm(
      `Este código de insumo tiene stock crítico:\n\nStock Actual: ${currentStock}\nStock Crítico: ${supplyCode.critical_stock}\n\n¿Estás seguro de que deseas editar este código?`,
      'ADVERTENCIA: Stock Crítico',
      { icon: 'warning' }
    )
    if (!confirmed) {
      return
    }
  }
  
  isEditing.value = true
  supplyCodeForm.value = { ...supplyCode }
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  isEditing.value = false
  supplyCodeForm.value = {
    code: null,
    name: '',
    code_supplier: null,
    critical_stock: 1
  }
  modalErrors.value = { code: '', name: '', code_supplier: '', critical_stock: '' }
}

const saveSupplyCode = async () => {
  // Validar campos obligatorios
  modalErrors.value = { code: '', name: '', code_supplier: '', critical_stock: '' }
  let hasErrors = false

  const codeVal = parseInt(supplyCodeForm.value.code)
  const codeSupplierVal = parseInt(supplyCodeForm.value.code_supplier)
  const criticalStockVal = parseInt(supplyCodeForm.value.critical_stock)

  if (!isEditing.value && (!supplyCodeForm.value.code || isNaN(codeVal) || codeVal < 1)) {
    modalErrors.value.code = 'El código del insumo es obligatorio y debe ser mayor a 0.'
    hasErrors = true
  }
  if (!supplyCodeForm.value.name || supplyCodeForm.value.name.trim() === '') {
    modalErrors.value.name = 'El nombre del insumo es obligatorio.'
    hasErrors = true
  }
  if (!supplyCodeForm.value.code_supplier || isNaN(codeSupplierVal) || codeSupplierVal < 1) {
    modalErrors.value.code_supplier = 'El código de proveedor es obligatorio y debe ser mayor a 0.'
    hasErrors = true
  }
  if (!supplyCodeForm.value.critical_stock || isNaN(criticalStockVal) || criticalStockVal < 1) {
    modalErrors.value.critical_stock = 'El stock crítico es obligatorio y debe ser mayor a 0.'
    hasErrors = true
  }

  if (hasErrors) return

  saving.value = true
  try {
    if (isEditing.value) {
      await supplyCodeService.updateSupplyCode(codeVal, {
        name: supplyCodeForm.value.name,
        code_supplier: codeSupplierVal,
        critical_stock: criticalStockVal
      })
      showSuccess('El código de insumo ha sido actualizado exitosamente.')
    } else {
      await supplyCodeService.createSupplyCode({
        code: codeVal,
        name: supplyCodeForm.value.name,
        code_supplier: codeSupplierVal,
        critical_stock: criticalStockVal
      })
      showSuccess('El código de insumo ha sido creado exitosamente.')
    }
    closeModal()
    await loadSupplyCodes()
  } catch (err) {
    const errorMessage = err.response?.data?.error || err.message || 'Error al guardar el código de insumo'
    showError(errorMessage)
  } finally {
    saving.value = false
  }
}

const confirmDelete = async (supplyCode) => {
  // Verificar stock crítico antes de mostrar el modal de confirmación
  const currentStock = await getCurrentStock(supplyCode.code)
  const isCritical = isStockCritical(currentStock, supplyCode.critical_stock)
  
  let confirmMessage
  
  if (isCritical) {
    confirmMessage = `ADVERTENCIA: Stock Crítico\n\nEste código de insumo tiene stock crítico:\n\nStock Actual: ${currentStock}\nStock Crítico: ${supplyCode.critical_stock}\n\n¿Estás seguro de que deseas eliminar este código?\n\nEsta acción no se puede deshacer.`
  } else {
    confirmMessage = `¿Estás seguro de que deseas eliminar el código "${supplyCode.code} - ${supplyCode.name}"?\n\nEsta acción no se puede deshacer.`
  }

  const confirmed = await confirmDanger(confirmMessage, 'Confirmar eliminación')
  if (!confirmed) {
    return
  }

  try {
    await supplyCodeService.deleteSupplyCode(supplyCode.code)
    showSuccess('El código de insumo ha sido eliminado exitosamente.')
    await loadSupplyCodes()
  } catch (err) {
    const rawError = err.response?.data?.error || err.message || ''
    const isForeignKeyViolation = rawError.includes('23503') || rawError.toLowerCase().includes('foreign key') || rawError.includes('llave foránea')
    const errorMessage = isForeignKeyViolation
      ? `No se puede eliminar el código "${supplyCode.code} - ${supplyCode.name}" porque está siendo utilizado por insumos médicos registrados en el sistema. Elimina o reasigna esos insumos antes de continuar.`
      : rawError || 'Error al eliminar el código de insumo'
    showError(errorMessage)
  }
}

// Helper functions
const getCriticalStockBadgeClass = (criticalStock) => {
  if (criticalStock === 1) return 'bg-red-100 text-red-800'
  if (criticalStock <= 5) return 'bg-orange-100 text-orange-800'
  return 'bg-yellow-100 text-yellow-800'
}

onMounted(() => {
  loadSupplyCodes()
})
</script>

