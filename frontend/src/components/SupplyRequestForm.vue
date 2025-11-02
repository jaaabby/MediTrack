<template>
  <div class="max-w-4xl mx-auto p-3 sm:p-6 bg-white rounded-lg shadow-lg">
    <!-- Título principal -->
    <div class="mb-6">
      <h2 class="text-2xl font-bold text-gray-900 mb-2">
        {{ props.editMode ? 'Editar Solicitud de Insumo' : 'Nueva Solicitud de Insumo' }}
      </h2>
      <p class="text-gray-600">
        {{ props.editMode ? 'Modificar y reenviar solicitud devuelta' : 'Crear solicitud con trazabilidad completa' }}
      </p>
    </div>

    <!-- Mostrar errores generales -->
    <div v-if="errors.length > 0" class="mb-4 sm:mb-6 p-3 sm:p-4 bg-red-50 border border-red-200 rounded-md">
      <div class="flex">
        <svg class="h-5 w-5 text-red-400 mr-2 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
        </svg>
        <div>
          <h3 class="text-sm font-medium text-red-800">Se encontraron los siguientes errores:</h3>
          <ul class="mt-2 text-sm text-red-700 list-disc list-inside">
            <li v-for="error in errors" :key="error">{{ error }}</li>
          </ul>
        </div>
      </div>
    </div>

    <form @submit.prevent="submitRequest" class="space-y-4 sm:space-y-6">
      <!-- Información básica -->
      <div class="bg-gray-50 p-3 sm:p-4 rounded-lg">
        <h3 class="text-base sm:text-lg font-semibold text-gray-900 mb-3 sm:mb-4">Información Básica</h3>
        
        <!-- Información del solicitante (automática) -->
        <div class="mb-4">
          <label class="block text-sm font-medium text-gray-700 mb-1">
            Solicitante
          </label>
          <div class="w-full px-3 py-2 bg-gray-50 border border-gray-300 rounded-md text-gray-700">
            <div class="font-medium">{{ authStore.getUserName || 'Usuario' }}</div>
            <div class="text-sm text-gray-500">{{ authStore.getUserRut || 'RUT no disponible' }}</div>
            <div v-if="authStore.getUserSpecialty" class="font-medium">
              <span class="inline-flex items-center">
                {{ authStore.getUserSpecialty }}
              </span>
            </div>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <!-- Fecha y hora de cirugía -->
          <div>
            <label for="surgery_datetime" class="block text-sm font-medium text-gray-700 mb-1">
              Fecha y Hora de Cirugía <span class="text-red-500">*</span>
            </label>
            <!-- Modo edición: mostrar como solo lectura -->
            <div v-if="props.editMode" class="w-full px-3 py-2 bg-gray-200 border border-gray-400 rounded-md text-gray-600 font-medium cursor-not-allowed">
              {{ originalRequestData.surgery_datetime_display }}
            </div>
            <!-- Modo creación: campo editable -->
            <input
              v-else
              type="datetime-local"
              id="surgery_datetime"
              v-model="requestForm.surgery_datetime"
              required
              :min="minDateTime"
              class="form-input"
            />
            <p class="text-xs text-gray-500 mt-1">
              {{ props.editMode ? 'Fecha y hora programada para la cirugía' : 'Selecciona la fecha y hora programada para la cirugía' }}
            </p>
          </div>

          <!-- Selección de Pabellón -->
          <div>
            <label for="pavilion" class="block text-xs sm:text-sm font-medium text-gray-700 mb-1">
              Pabellón <span class="text-red-500">*</span>
            </label>
            <!-- Modo edición: mostrar como solo lectura -->
            <div v-if="props.editMode" class="w-full px-3 py-2 bg-gray-200 border border-gray-400 rounded-md text-gray-600 font-medium cursor-not-allowed">
              {{ originalRequestData.pavilion_name }}
            </div>
            <!-- Modo creación: campo editable -->
            <select
              v-else
              id="pavilion"
              v-model="requestForm.pavilion_id"
              required
              class="form-select text-sm"
              :disabled="loadingPavilions"
            >
              <option value="">Seleccionar pabellón</option>
              <option
                v-for="pavilion in pavilions"
                :key="pavilion.id"
                :value="pavilion.id"
              >
                {{ pavilion.name }}
              </option>
            </select>
            <p v-if="loadingPavilions && !props.editMode" class="text-xs text-gray-500 mt-1">Cargando pabellones...</p>
            <p v-else class="text-xs text-gray-500 mt-1">
              {{ props.editMode ? 'Pabellón donde se realizará la cirugía' : 'Selecciona el pabellón donde se realizará la cirugía' }}
            </p>
          </div>

          <!-- Tipo de Cirugía -->
          <div>
            <label for="surgery_id" class="block text-sm font-medium text-gray-700 mb-1">
              Tipo de Cirugía <span class="text-red-500">*</span>
            </label>
            <select
              id="surgery_id"
              v-model="requestForm.surgery_id"
              required
              class="form-select text-sm"
              :disabled="loadingSurgeries"
            >
              <option :value="null">Seleccionar cirugía</option>
              <option v-for="surgery in filteredSurgeries" :key="surgery.id" :value="surgery.id">
                {{ surgery.name }}
              </option>
            </select>
            <p v-if="loadingSurgeries" class="text-xs text-gray-500 mt-1">Cargando cirugías...</p>
            <p v-else-if="filteredSurgeries.length === 0" class="text-xs text-yellow-600 mt-1">
              No hay cirugías disponibles para tu especialidad
            </p>
            <p v-else class="text-xs text-gray-500 mt-1">
              Cirugías de {{ authStore.getUserSpecialty || 'tu especialidad' }}
            </p>
          </div>
        </div>

        <!-- Observaciones generales -->
        <div class="mt-3 sm:mt-4">
          <label for="notes" class="block text-xs sm:text-sm font-medium text-gray-700 mb-1">
            Observaciones Generales
          </label>
          <textarea
            id="notes"
            v-model="requestForm.notes"
            rows="3"
            :placeholder="props.editMode ? 'Observaciones sobre los cambios realizados en esta solicitud...' : 'Observaciones adicionales sobre la solicitud...'"
                  class="form-input text-sm"
          ></textarea>
        </div>
      </div>

      <!-- Insumos solicitados -->
      <div class="bg-blue-50 p-3 sm:p-4 rounded-lg">
        <div class="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-3 sm:gap-0 mb-3 sm:mb-4">
          <h3 class="text-base sm:text-lg font-semibold text-gray-900">Insumos Necesarios</h3>
          <button
            type="button"
            @click="addSupplyItem"
            class="btn-primary w-full sm:w-auto"
          >
            <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            Agregar Insumo
          </button>
        </div>

        <!-- Lista de insumos -->
        <div v-if="requestForm.items.length === 0" class="text-center py-6 sm:py-8 text-gray-500">
          <svg class="h-12 w-12 mx-auto mb-4 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
          </svg>
          <p>No hay insumos agregados</p>
          <p class="text-sm">Haz clic en "Agregar Insumo" para comenzar</p>
        </div>

        <div class="space-y-3 sm:space-y-4">
          <div
            v-for="(item, index) in requestForm.items"
            :key="index"
            class="bg-white p-3 sm:p-4 rounded-lg border border-gray-200 relative"
          >
            <!-- Botón eliminar -->
            <button
              type="button"
              @click="removeSupplyItem(index)"
              class="absolute top-2 right-2 p-1.5 text-red-500 hover:text-red-700 hover:bg-red-50 rounded-full z-10"
              title="Eliminar insumo"
            >
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>

            <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3 sm:gap-4 pr-8">
              <!-- Código del insumo -->
              <div>
                <label class="block text-xs sm:text-sm font-medium text-gray-700 mb-1">
                  Código Insumo <span class="text-red-500">*</span>
                </label>
                <input
                  type="number"
                  v-model.number="item.supply_code"
                  required
                  placeholder="1234"
                  class="form-select text-sm"
                />
              </div>

              <!-- Nombre del insumo con autocompletado -->
              <div class="relative">
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  Nombre Insumo <span class="text-red-500">*</span>
                </label>
                <input
                  type="text"
                  :value="supplySearchTerms[index] || item.supply_name"
                  @input="onSupplyInputChange(index, $event.target.value)"
                  @focus="onSupplyInputFocus(index)"
                  @blur="onSupplyInputBlur(index)"
                  required
                  placeholder="Buscar insumo..."
                  class="form-input"
                  autocomplete="off"
                />
                
                <!-- Dropdown de sugerencias -->
                <div 
                  v-if="showSupplyDropdowns[index] && medicalSupplies.length > 0"
                  class="absolute z-10 w-full mt-1 bg-white border border-gray-300 rounded-md shadow-lg max-h-60 overflow-y-auto"
                >
                  <div 
                    v-for="supply in getFilteredSupplies(index)" 
                    :key="supply.id"
                    @click="selectSupply(index, supply)"
                    class="px-3 py-2 hover:bg-blue-50 cursor-pointer border-b border-gray-100 last:border-b-0"
                  >
                    <div class="font-medium text-gray-900">{{ supply.name }}</div>
                    <div class="text-sm text-gray-500">Código: {{ supply.code }}</div>
                  </div>
                  
                  <div v-if="getFilteredSupplies(index).length === 0" class="px-3 py-2 text-gray-500 text-center">
                    No se encontraron insumos
                  </div>
                </div>
                
                <!-- Indicador de carga -->
                <div v-if="loadingSupplies" class="absolute right-3 top-9 text-gray-400">
                  <svg class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                </div>
              </div>

              <!-- Cantidad solicitada -->
              <div>
                <label class="block text-xs sm:text-sm font-medium text-gray-700 mb-1">
                  Cantidad <span class="text-red-500">*</span>
                </label>
                <input
                  type="number"
                  v-model.number="item.quantity_requested"
                  required
                  min="1"
                  placeholder="1"
                  class="form-select text-sm"
                />
              </div>
            </div>

            <!-- Observaciones del item devuelto (solo en modo edición cuando item fue devuelto) -->
            <div v-if="props.editMode && item.item_status === 'devuelto'" class="mt-4 space-y-3">
              <!-- Mostrar observaciones anteriores del encargado -->
              <div v-if="item.item_notes" class="bg-orange-50 border border-orange-200 rounded-md p-3">
                <label class="block text-xs font-semibold text-orange-800 mb-1">
                  Motivo de devolución del encargado:
                </label>
                <p class="text-xs text-orange-900 whitespace-pre-wrap">{{ item.item_notes }}</p>
              </div>
              
              <!-- Campo para nuevas observaciones del doctor -->
              <div>
                <label class="block text-xs sm:text-sm font-medium text-gray-700 mb-1">
                  Observaciones Insumo
                </label>
                <textarea
                  v-model="item.resubmit_notes"
                  rows="2"
                  placeholder="Agregue sus observaciones..."
                  class="w-full px-3 py-2 text-sm border border-orange-300 bg-white rounded-md focus:outline-none focus:ring-2 focus:ring-orange-500 focus:border-orange-500"
                ></textarea>
                <p class="mt-1 text-xs text-orange-600">
                  Opcional: Agregue observaciones si modificó la cantidad o desea aclarar algo sobre este insumo.
                </p>
              </div>
            </div>

            <!-- Especificaciones técnicas -->
            <div class="mt-4 grid grid-cols-1 md:grid-cols-2 gap-4">
              <!-- Medidas/Tamaño 
              <div>
                <label class="block text-xs sm:text-sm font-medium text-gray-700 mb-1">
                  Medidas/Tamaño
                </label>
                <input
                  type="text"
                  v-model="item.size"
                  placeholder="Ej: Grande, Mediano, 20cm, etc."
                  class="form-select text-sm"
                />
              </div>-->

              <!-- Marca 
              <div>
                <label class="block text-xs sm:text-sm font-medium text-gray-700 mb-1">
                  Marca
                </label>
                <input
                  type="text"
                  v-model="item.brand"
                  placeholder="Marca preferida"
                  class="form-select text-sm"
                />
              </div>
            </div>-->

            <!-- Características especiales
            <div class="mt-4 grid grid-cols-1 md:grid-cols-2 gap-4">-->
              <!-- Es pediátrico -->
              <div class="flex items-center">
                <input
                  type="checkbox"
                  :id="`pediatric-${index}`"
                  v-model="item.is_pediatric"
                  class="h-4 w-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                />
                <label :for="`pediatric-${index}`" class="ml-2 text-xs sm:text-sm font-medium text-gray-700">
                  Es insumo pediátrico
                </label>
              </div>
            </div>

            <!-- Especificaciones y observaciones del insumo 
            <div class="mt-4">
              <label class="block text-sm font-medium text-gray-700 mb-1">
                Especificaciones Técnicas
              </label>
              <textarea
                v-model="item.specifications"
                rows="2"
                placeholder="Especificaciones técnicas del insumo (material, dimensiones exactas, características especiales, etc.)"
                class="form-select text-sm"
              ></textarea>
            </div>-->

            <!-- Solicitudes especiales
            <div class="mt-4">
              <label class="block text-sm font-medium text-gray-700 mb-1">
                Solicitudes Especiales
              </label>
              <textarea
                v-model="item.special_requests"
                rows="2"
                placeholder="Solicitudes especiales para este insumo (entrega urgente, manipulación especial, etc.)"
                class="form-select text-sm"
              ></textarea>
            </div>-->
          </div>
        </div>
      </div>

      <!-- Botones de acción -->
      <div class="flex flex-col sm:flex-row sm:justify-between gap-3 sm:gap-0 pt-3 sm:pt-4 border-t border-gray-200">
        <button
          type="button"
          @click="resetForm"
          class="inline-flex items-center justify-center px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 w-full sm:w-auto order-3 sm:order-1"
          :disabled="submitting"
        >
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          Limpiar Formulario
        </button>
        
        <div class="flex flex-col sm:flex-row gap-2 sm:gap-3 sm:space-x-0 order-1 sm:order-2">
          <button
            type="button"
            @click="cancelForm"
            class="inline-flex items-center justify-center px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-gray-500 w-full sm:w-auto order-2"
            :disabled="submitting"
          >
            Cancelar
          </button>
          
          <button
            type="submit"
            class="btn-primary disabled:opacity-50 disabled:cursor-not-allowed w-full sm:w-auto order-1"
            :disabled="submitting"
          >
            <svg v-if="submitting" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <svg v-if="!submitting" class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
            <span v-if="props.editMode">
              {{ submitting ? 'Reenviando Solicitud...' : 'Reenviar Solicitud' }}
            </span>
            <span v-else>
              {{ submitting ? 'Creando Solicitud...' : 'Crear Solicitud' }}
            </span>
          </button>
        </div>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, nextTick } from 'vue'
import { useAuthStore } from '@/stores/auth'
import supplyRequestService from '../services/supplyRequestService'
import pavilionService from '../services/pavilionService'
import inventoryService from '../services/inventoryService'
import surgeryService from '../services/surgeryService'
import Swal from 'sweetalert2'

// Props
const props = defineProps({
  id: {
    type: Number,
    default: null
  },
  editMode: {
    type: Boolean,
    default: false
  }
})

// Emits
const emit = defineEmits(['success', 'cancel', 'error'])

// Stores
const authStore = useAuthStore()

// Estado reactivo
const submitting = ref(false)
const errors = ref([])
const pavilions = ref([])
const loadingPavilions = ref(false)
const medicalSupplies = ref([])
const loadingSupplies = ref(false)
const supplySearchTerms = ref([])
const showSupplyDropdowns = ref([])
const surgeries = ref([])
const loadingSurgeries = ref(false)

// Fecha mínima para la cirugía (hoy)
const minDateTime = computed(() => {
  const now = new Date()
  const year = now.getFullYear()
  const month = String(now.getMonth() + 1).padStart(2, '0')
  const day = String(now.getDate()).padStart(2, '0')
  const hours = String(now.getHours()).padStart(2, '0')
  const minutes = String(now.getMinutes()).padStart(2, '0')
  return `${year}-${month}-${day}T${hours}:${minutes}`
})

// Formatear fecha para mostrar (en modo edición)
const formatDateTimeForDisplay = (dateString) => {
  if (!dateString) return 'No especificada'
  
  try {
    const date = new Date(dateString)
    if (isNaN(date.getTime())) return 'Fecha inválida'
    
    const options = {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      hour12: false
    }
    return date.toLocaleDateString('es-CL', options)
  } catch (error) {
    console.error('Error formateando fecha:', error)
    return dateString
  }
}

// Formulario de solicitud
const requestForm = reactive({
  pavilion_id: '',
  surgery_datetime: '',
  notes: '',
  items: [],
  surgery_id: null
})

// Datos originales para mostrar en modo edición (solo lectura)
const originalRequestData = ref({
  pavilion_name: '',
  surgery_datetime_display: '',
  requester_name: '',
  requester_rut: ''
})

// Computed: Filtrar cirugías según la especialidad del doctor logueado
const filteredSurgeries = computed(() => {
  const userSpecialtyId = authStore.getUserSpecialtyId
  
  // Si el usuario no tiene especialidad, mostrar todas las cirugías
  if (!userSpecialtyId) {
    return surgeries.value
  }
  
  // Filtrar cirugías que coincidan con la especialidad del usuario
  return surgeries.value.filter(surgery => surgery.specialty_id === userSpecialtyId)
})

// Template para nuevo item
const newSupplyItem = () => ({
  supply_code: '',
  supply_name: '',
  quantity_requested: 1,
  specifications: '',
  is_pediatric: false,
  special_requests: '',
  urgency_level: 'normal',
  size: '',
  brand: ''
})

// Métodos
const loadPavilions = async () => {
  loadingPavilions.value = true
  try {
    const result = await pavilionService.getAllPavilions()
    pavilions.value = result
  } catch (error) {
    console.error('Error cargando pabellones:', error)
    errors.value.push('Error al cargar la lista de pabellones')
    emit('error', error)
  } finally {
    loadingPavilions.value = false
  }
}

const loadSurgeries = async () => {
  loadingSurgeries.value = true
  try {
    const result = await surgeryService.getAllSurgeries()
    surgeries.value = result
  } catch (error) {
    console.error('Error cargando cirugías:', error)
  } finally {
    loadingSurgeries.value = false
  }
}

const loadMedicalSupplies = async () => {
  loadingSupplies.value = true
  try {
    const result = await inventoryService.getAllSupplyCodes()
    medicalSupplies.value = result
  } catch (error) {
    console.error('Error cargando códigos de insumo:', error)
    errors.value.push('Error al cargar la lista de códigos de insumo')
    emit('error', error)
  } finally {
    loadingSupplies.value = false
  }
}

const addSupplyItem = () => {
  requestForm.items.unshift(newSupplyItem())
  // Inicializar estados de búsqueda para el nuevo item al principio
  supplySearchTerms.value.unshift('')
  showSupplyDropdowns.value.unshift(false)
}

const removeSupplyItem = (index) => {
  requestForm.items.splice(index, 1)
  // También remover los estados de búsqueda correspondientes
  supplySearchTerms.value.splice(index, 1)
  showSupplyDropdowns.value.splice(index, 1)
}

// Funciones para autocompletado de insumos
const getFilteredSupplies = (index) => {
  const searchTerm = supplySearchTerms.value[index] || ''
  if (!searchTerm) return medicalSupplies.value.slice(0, 10) // Mostrar primeros 10 si no hay búsqueda
  
  return medicalSupplies.value.filter(supply => 
    supply.name?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    supply.code?.toString().includes(searchTerm)
  ).slice(0, 10) // Limitar a 10 resultados
}

const selectSupply = (index, supply) => {
  requestForm.items[index].supply_name = supply.name
  requestForm.items[index].supply_code = supply.code
  supplySearchTerms.value[index] = supply.name
  showSupplyDropdowns.value[index] = false
}

const onSupplyInputFocus = (index) => {
  showSupplyDropdowns.value[index] = true
  // Inicializar el término de búsqueda si no existe
  if (!supplySearchTerms.value[index]) {
    supplySearchTerms.value[index] = requestForm.items[index].supply_name || ''
  }
}

const onSupplyInputBlur = (index) => {
  // Delay para permitir clicks en el dropdown
  setTimeout(() => {
    showSupplyDropdowns.value[index] = false
  }, 200)
}

const onSupplyInputChange = (index, value) => {
  supplySearchTerms.value[index] = value
  requestForm.items[index].supply_name = value
  showSupplyDropdowns.value[index] = true
}

const validateForm = () => {
  const validation = supplyRequestService.validateSupplyRequest(requestForm)
  errors.value = validation.errors
  return validation.isValid
}

const resetForm = () => {
  Object.assign(requestForm, {
    pavilion_id: '',
    surgery_datetime: '',
    notes: '',
    items: [],
    surgery_id: null
  })
  errors.value = []
}

const cancelForm = () => {
  emit('cancel')
}

const submitRequest = async () => {
  if (!validateForm()) {
    return
  }

  submitting.value = true
  errors.value = []

  try {
    // Si está en modo edición, reenviar la solicitud devuelta
    if (props.editMode && props.id) {
      await resubmitRequest()
    } else {
      // Crear nueva solicitud
      await createNewRequest()
    }
  } catch (error) {
    console.error('Error al enviar solicitud:', error)
    
    let errorMessage = 'Error desconocido al procesar la solicitud'
    if (error.response?.data?.error) {
      errorMessage = 'Error del servidor: ' + error.response.data.error
    } else if (error.message) {
      errorMessage = 'Error: ' + error.message
    }
    
    errors.value.push(errorMessage)
    emit('error', error)
  } finally {
    submitting.value = false
  }
}

// Crear nueva solicitud
const createNewRequest = async () => {
  const formattedData = supplyRequestService.formatSupplyRequestForAPI(requestForm)
  console.log('Enviando solicitud:', formattedData)
  
  const result = await supplyRequestService.createSupplyRequest(formattedData)
  
  if (result.success) {
    console.log('Solicitud creada exitosamente:', result)
    emit('success', result.data?.request || result.data)
  } else {
    const errorMessage = 'Error al crear la solicitud: ' + (result.error || 'Error desconocido')
    errors.value.push(errorMessage)
    emit('error', new Error(errorMessage))
  }
}

// Reenviar solicitud devuelta
const resubmitRequest = async () => {
  // Preparar solo los items que fueron devueltos con sus nuevas cantidades y observaciones
  const updatedItems = requestForm.items
    .filter(item => item.item_status === 'devuelto')
    .map(item => ({
      item_id: item.id,
      quantity: item.quantity_requested,
      notes: item.resubmit_notes || '' // Nuevas observaciones del doctor sobre este item
    }))
  
  console.log('Reenviando solicitud con items:', updatedItems)
  console.log('Notas del solicitante:', requestForm.notes)
  
  const result = await supplyRequestService.resubmitReturnedRequest(props.id, updatedItems, requestForm.notes)
  
  if (result.success) {
    console.log('Solicitud reenviada exitosamente')
    
    await Swal.fire({
      icon: 'success',
      title: 'Solicitud Reenviada',
      text: 'La solicitud ha sido reenviada al encargado de bodega',
      timer: 2000,
      showConfirmButton: false
    })
    
    emit('success', { id: props.id })
  } else {
    const errorMessage = 'Error al reenviar la solicitud: ' + (result.error || 'Error desconocido')
    errors.value.push(errorMessage)
    emit('error', new Error(errorMessage))
  }
}

// Lifecycle
onMounted(async () => {
  // Cargar listas primero
  await Promise.all([
    loadPavilions(),
    loadMedicalSupplies(),
    loadSurgeries()
  ])
  
  // Si está en modo edición, cargar la solicitud DESPUÉS de que las listas estén listas
  if (props.editMode && props.id) {
    await loadRequestForEdit()
  } else {
    // Agregar un insumo por defecto solo si no está en modo edición
    addSupplyItem()
  }
})

// Función auxiliar para extraer solo las notas originales del solicitante
const extractOriginalNotes = (fullNotes) => {
  if (!fullNotes) return ''
  
  // Buscar el marcador de devolución
  const devolucionIndex = fullNotes.indexOf('[Devolución por')
  
  // Si no hay marcador, devolver todas las notas
  if (devolucionIndex === -1) {
    return fullNotes.trim()
  }
  
  // Extraer solo la parte antes del marcador
  return fullNotes.substring(0, devolucionIndex).trim()
}

// Cargar solicitud para editar
const loadRequestForEdit = async () => {
  try {
    console.log('🔄 Cargando solicitud para editar, ID:', props.id)
    const response = await supplyRequestService.getSupplyRequestById(props.id)
    console.log('📦 Respuesta de la solicitud:', response)
    
    if (response.success && response.data) {
      // El backend devuelve data.request, data.items, data.assignments
      const request = response.data.request || response.data
      console.log('✅ Datos de la solicitud completa:', JSON.stringify(request, null, 2))
      console.log('🏥 Pabellones disponibles:', JSON.stringify(pavilions.value, null, 2))
      console.log('🔍 Datos directos - pavilion_id:', request.pavilion_id, 'surgery_datetime:', request.surgery_datetime)
      
      const pavilionId = request.pavilion_id
      const surgeryDatetime = request.surgery_datetime
      
      console.log('📊 Valores finales - pavilionId:', pavilionId, 'surgeryDatetime:', surgeryDatetime)
      
      // Cargar datos básicos
      requestForm.pavilion_id = pavilionId ? pavilionId.toString() : ''
      requestForm.surgery_datetime = surgeryDatetime ? formatDateTimeForInput(surgeryDatetime) : ''
      // En modo edición, dejar las observaciones vacías para que el usuario ingrese nuevas
      requestForm.notes = ''
      
      // Cargar tipo de cirugía
      requestForm.surgery_id = request.surgery_id || null
      
      // Guardar datos originales para mostrar en modo solo lectura
      // Buscar el pabellón por ID - probar con conversión de tipos
      const pavilionIdToFind = parseInt(pavilionId || 0)
      let pavilionName = 'No especificado'
      
      console.log('🔍 Buscando pabellón con ID:', pavilionIdToFind)
      console.log('📋 Pabellones disponibles:', pavilions.value)
      
      if (pavilionIdToFind > 0) {
        const foundPavilion = pavilions.value.find(p => {
          const pId = parseInt(p.id)
          console.log(`  Comparando: ${pId} === ${pavilionIdToFind}?`, pId === pavilionIdToFind)
          return pId === pavilionIdToFind
        })
        
        if (foundPavilion) {
          pavilionName = foundPavilion.name
          console.log('✅ Pabellón encontrado:', pavilionName)
        } else {
          console.warn('⚠️ Pabellón no encontrado para ID:', pavilionIdToFind)
          pavilionName = `Pabellón ID: ${pavilionIdToFind}`
        }
      }
      
      originalRequestData.value = {
        pavilion_name: pavilionName,
        surgery_datetime_display: surgeryDatetime ? formatDateTimeForDisplay(surgeryDatetime) : 'No especificada',
        requester_name: request.requested_by_name || 'No disponible',
        requester_rut: request.requested_by || 'No disponible'
      }
      
      console.log('📝 Datos originales asignados:', originalRequestData.value)
      
      console.log('📝 Datos originales guardados:', originalRequestData.value)
      
      // Forzar actualización del DOM
      await nextTick()
      
      console.log('📝 Formulario actualizado:', {
        pavilion_id: requestForm.pavilion_id,
        surgery_datetime: requestForm.surgery_datetime,
        notes: requestForm.notes
      })
      
      // Cargar items
      const itemsResponse = await supplyRequestService.getSupplyRequestItems(props.id)
      if (itemsResponse.success && itemsResponse.data) {
        requestForm.items = itemsResponse.data.map(item => ({
          id: item.id, // Incluir ID para tracking
          supply_code: item.supply_code,
          supply_name: item.supply_name,
          quantity_requested: item.quantity_requested,
          is_pediatric: item.is_pediatric,
          item_status: item.item_status, // Para saber cuáles son editables
          item_notes: item.item_notes, // Observaciones anteriores del encargado (solo lectura)
          resubmit_notes: '', // Nuevas observaciones del doctor (editable)
          specifications: '',
          special_requests: '',
          urgency_level: 'normal',
          size: '',
          brand: ''
        }))
        
        // Inicializar arrays de búsqueda
        supplySearchTerms.value = requestForm.items.map(() => '')
        showSupplyDropdowns.value = requestForm.items.map(() => false)
      }
    }
  } catch (error) {
    console.error('Error cargando solicitud:', error)
    Swal.fire({
      icon: 'error',
      title: 'Error',
      text: 'No se pudo cargar la solicitud para editar'
    })
  }
}

// Función auxiliar para formatear fecha para input datetime-local
const formatDateTimeForInput = (dateString) => {
  const date = new Date(dateString)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  return `${year}-${month}-${day}T${hours}:${minutes}`
}

</script>

<style scoped>
/* Estilos adicionales si son necesarios */
.bg-red-50 {
  background-color: #fef2f2;
}

.border-red-200 {
  border-color: #fecaca;
}

.text-red-400 {
  color: #f87171;
}

.text-red-800 {
  color: #991b1b;
}

.text-red-700 {
  color: #b91c1c;
}

.text-red-500 {
  color: #ef4444;
}

.hover\:bg-red-50:hover {
  background-color: #fef2f2;
}

.hover\:text-red-700:hover {
  color: #b91c1c;
}

/* Mejorar experiencia táctil en móviles */
@media (max-width: 640px) {
  /* Aumentar área táctil de botones */
  button {
    min-height: 44px;
  }
  
  /* Mejorar áreas de entrada */
  input,
  select,
  textarea {
    min-height: 44px;
    font-size: 16px; /* Prevenir zoom automático en iOS */
  }
  
  /* Mejorar área de checkboxes */
  input[type="checkbox"] {
    min-width: 20px;
    min-height: 20px;
  }

  /* Suavizar transiciones */
  * {
    -webkit-tap-highlight-color: transparent;
  }
}

/* Transiciones suaves */
button,
input,
select,
textarea {
  transition: all 0.2s ease-in-out;
}
</style>