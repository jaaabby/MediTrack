<template>
  <div class="max-w-4xl mx-auto p-3 sm:p-6 bg-white rounded-lg shadow-lg">
    <!-- Título principal -->
    <div class="mb-4 sm:mb-6">
      <h2 class="text-xl sm:text-2xl font-bold text-gray-900 mb-2">Nueva Solicitud de Insumo</h2>
      <p class="text-sm sm:text-base text-gray-600">Crear solicitud con trazabilidad QR completa</p>
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
        
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 sm:gap-4">
          <!-- Selección de Pabellón -->
          <div>
            <label for="pavilion" class="block text-xs sm:text-sm font-medium text-gray-700 mb-1">
              Pabellón <span class="text-red-500">*</span>
            </label>
            <select
              id="pavilion"
              v-model="requestForm.pavilion_id"
              required
              class="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
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
            <p v-if="loadingPavilions" class="text-xs text-gray-500 mt-1">Cargando pabellones...</p>
          </div>

          <!-- Prioridad -->
          <div>
            <label for="priority" class="block text-xs sm:text-sm font-medium text-gray-700 mb-1">
              Prioridad <span class="text-red-500">*</span>
            </label>
            <select
              id="priority"
              v-model="requestForm.priority"
              required
              class="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="low">Baja</option>
              <option value="normal">Normal</option>
              <option value="high">Alta</option>
              <option value="critical">Crítica</option>
            </select>
          </div>

          <!-- Solicitante -->
          <div>
            <label for="requestedBy" class="block text-xs sm:text-sm font-medium text-gray-700 mb-1">
              Solicitante (RUT) <span class="text-red-500">*</span>
            </label>
            <input
              type="text"
              id="requestedBy"
              v-model="requestForm.requested_by"
              required
              placeholder="12.345.678-9"
              class="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            />
          </div>

          <!-- Nombre del solicitante -->
          <div>
            <label for="requestedByName" class="block text-xs sm:text-sm font-medium text-gray-700 mb-1">
              Nombre Completo <span class="text-red-500">*</span>
            </label>
            <input
              type="text"
              id="requestedByName"
              v-model="requestForm.requested_by_name"
              required
              placeholder="Nombre y apellidos"
              class="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            />
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
            placeholder="Observaciones adicionales sobre la solicitud..."
            class="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
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
            class="inline-flex items-center justify-center px-3 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 w-full sm:w-auto"
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
                  class="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                />
              </div>

              <!-- Nombre del insumo -->
              <div>
                <label class="block text-xs sm:text-sm font-medium text-gray-700 mb-1">
                  Nombre Insumo <span class="text-red-500">*</span>
                </label>
                <input
                  type="text"
                  v-model="item.supply_name"
                  required
                  placeholder="Nombre del insumo"
                  class="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                />
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
                  class="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                />
              </div>
            </div>

            <!-- Especificaciones técnicas -->
            <div class="mt-3 sm:mt-4 grid grid-cols-1 sm:grid-cols-2 gap-3 sm:gap-4">
              <!-- Medidas/Tamaño -->
              <div>
                <label class="block text-xs sm:text-sm font-medium text-gray-700 mb-1">
                  Medidas/Tamaño
                </label>
                <input
                  type="text"
                  v-model="item.size"
                  placeholder="Ej: Grande, Mediano, 20cm, etc."
                  class="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                />
              </div>

              <!-- Marca -->
              <div>
                <label class="block text-xs sm:text-sm font-medium text-gray-700 mb-1">
                  Marca
                </label>
                <input
                  type="text"
                  v-model="item.brand"
                  placeholder="Marca preferida"
                  class="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                />
              </div>
            </div>

            <!-- Características especiales -->
            <div class="mt-3 sm:mt-4 grid grid-cols-1 sm:grid-cols-2 gap-3 sm:gap-4">
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

              <!-- Nivel de urgencia -->
              <div>
                <label class="block text-xs sm:text-sm font-medium text-gray-700 mb-1">
                  Urgencia del Item
                </label>
                <select
                  v-model="item.urgency_level"
                  class="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                >
                  <option value="low">Baja</option>
                  <option value="normal">Normal</option>
                  <option value="high">Alta</option>
                  <option value="critical">Crítica</option>
                </select>
              </div>
            </div>

            <!-- Especificaciones y observaciones del insumo -->
            <div class="mt-3 sm:mt-4">
              <label class="block text-xs sm:text-sm font-medium text-gray-700 mb-1">
                Especificaciones Técnicas
              </label>
              <textarea
                v-model="item.specifications"
                rows="2"
                placeholder="Especificaciones técnicas del insumo (material, dimensiones exactas, características especiales, etc.)"
                class="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              ></textarea>
            </div>

            <!-- Solicitudes especiales -->
            <div class="mt-3 sm:mt-4">
              <label class="block text-xs sm:text-sm font-medium text-gray-700 mb-1">
                Solicitudes Especiales
              </label>
              <textarea
                v-model="item.special_requests"
                rows="2"
                placeholder="Solicitudes especiales para este insumo (entrega urgente, manipulación especial, etc.)"
                class="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              ></textarea>
            </div>
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
            class="inline-flex items-center justify-center px-6 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed w-full sm:w-auto order-1"
            :disabled="submitting"
          >
            <svg v-if="submitting" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <svg v-if="!submitting" class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
            {{ submitting ? 'Creando Solicitud...' : 'Crear Solicitud' }}
          </button>
        </div>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import supplyRequestService from '../services/supplyRequestService'
import pavilionService from '../services/pavilionService'

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

// Estado reactivo
const submitting = ref(false)
const errors = ref([])
const pavilions = ref([])
const loadingPavilions = ref(false)

// Formulario de solicitud
const requestForm = reactive({
  pavilion_id: '',
  requested_by: '',
  requested_by_name: '',
  priority: 'normal',
  notes: '',
  items: []
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
    console.log('Pabellones cargados:', pavilions.value)
  } catch (error) {
    console.error('Error cargando pabellones:', error)
    errors.value.push('Error al cargar la lista de pabellones')
    emit('error', error)
  } finally {
    loadingPavilions.value = false
  }
}

const addSupplyItem = () => {
  requestForm.items.push(newSupplyItem())
}

const removeSupplyItem = (index) => {
  requestForm.items.splice(index, 1)
}

const validateForm = () => {
  const validation = supplyRequestService.validateSupplyRequest(requestForm)
  errors.value = validation.errors
  return validation.isValid
}

const resetForm = () => {
  Object.assign(requestForm, {
    pavilion_id: '',
    requested_by: '',
    requested_by_name: '',
    priority: 'normal',
    notes: '',
    items: []
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
    const formattedData = supplyRequestService.formatSupplyRequestForAPI(requestForm)
    console.log('Enviando solicitud:', formattedData)
    
    const result = await supplyRequestService.createSupplyRequest(formattedData)
    
    if (result.success) {
      console.log('Solicitud creada exitosamente:', result)
      
      // Emitir evento de éxito con los datos de la solicitud
      emit('success', result.data?.request || result.data)
    } else {
      const errorMessage = 'Error al crear la solicitud: ' + (result.error || 'Error desconocido')
      errors.value.push(errorMessage)
      emit('error', new Error(errorMessage))
    }
  } catch (error) {
    console.error('Error al enviar solicitud:', error)
    
    let errorMessage = 'Error desconocido al crear la solicitud'
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

// Lifecycle
onMounted(() => {
  loadPavilions()
  // Agregar un insumo por defecto
  addSupplyItem()
})
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