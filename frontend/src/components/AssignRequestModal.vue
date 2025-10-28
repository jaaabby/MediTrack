<template>
  <div v-if="show" class="fixed inset-0 z-50 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
    <!-- Overlay -->
    <div class="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:block sm:p-0">
      <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" aria-hidden="true" @click="closeModal"></div>

      <!-- Centrado del modal -->
      <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>

      <!-- Contenido del modal -->
      <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
        <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
          <div class="sm:flex sm:items-start">
            <div class="mx-auto flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full bg-blue-100 sm:mx-0 sm:h-10 sm:w-10">
              <svg class="h-6 w-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
              </svg>
            </div>
            <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left flex-1">
              <h3 class="text-lg leading-6 font-medium text-gray-900" id="modal-title">
                Asignar Solicitud
              </h3>
              <div class="mt-2">
                <p class="text-sm text-gray-500">
                  Solicitud: <span class="font-semibold text-gray-900">{{ request?.request_number }}</span>
                </p>
                <p class="text-sm text-gray-500 mt-1">
                  Solicitante: <span class="font-semibold text-gray-900">{{ request?.requested_by_name }}</span>
                </p>
              </div>

              <!-- Formulario de asignación -->
              <div class="mt-4 space-y-4">
                <!-- Selector de encargado de bodega -->
                <div>
                  <label for="warehouse-manager" class="block text-sm font-medium text-gray-700 mb-1">
                    Asignar a <span class="text-red-500">*</span>
                  </label>
                  <select
                    id="warehouse-manager"
                    v-model="selectedWarehouseManager"
                    :disabled="loading"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100 disabled:cursor-not-allowed"
                    :class="{ 'border-red-500': errors.warehouseManager }"
                  >
                    <option value="">{{ loading ? 'Cargando encargados...' : 'Seleccione un encargado...' }}</option>
                    <option 
                      v-for="manager in warehouseManagers" 
                      :key="manager.rut" 
                      :value="manager.rut"
                    >
                      {{ manager.name }} - {{ manager.rut }}
                    </option>
                  </select>
                  <p v-if="errors.warehouseManager" class="mt-1 text-sm text-red-600">
                    {{ errors.warehouseManager }}
                  </p>
                  <p v-if="warehouseManagers.length === 0 && !loading && !errorMessage" class="mt-1 text-sm text-gray-500">
                    No hay encargados de bodega disponibles
                  </p>
                </div>

                <!-- Notas de Pavedad -->
                <div>
                  <label for="pavedad-notes" class="block text-sm font-medium text-gray-700 mb-1">
                    Notas adicionales
                  </label>
                  <textarea
                    id="pavedad-notes"
                    v-model="pavedadNotes"
                    rows="3"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    placeholder="Instrucciones o comentarios para el encargado de bodega..."
                  ></textarea>
                </div>

                <!-- Mensaje de error -->
                <div v-if="errorMessage" class="bg-red-50 border border-red-200 rounded-md p-3">
                  <div class="flex">
                    <div class="flex-shrink-0">
                      <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                      </svg>
                    </div>
                    <div class="ml-3">
                      <p class="text-sm text-red-800">{{ errorMessage }}</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Botones de acción -->
        <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse gap-2">
          <button
            type="button"
            @click="handleAssign"
            :disabled="loading || !selectedWarehouseManager"
            class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-blue-600 text-base font-medium text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <svg v-if="loading" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ loading ? 'Asignando...' : 'Asignar Solicitud' }}
          </button>
          <button
            type="button"
            @click="closeModal"
            :disabled="loading"
            class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50"
          >
            Cancelar
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useAuthStore } from '@/stores/auth'
import userService from '@/services/userService'
import supplyRequestService from '@/services/supplyRequestService'
import Swal from 'sweetalert2'

const props = defineProps({
  show: {
    type: Boolean,
    required: true
  },
  request: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['close', 'assigned'])

const authStore = useAuthStore()

// Estado
const loading = ref(false)
const warehouseManagers = ref([])
const selectedWarehouseManager = ref('')
const pavedadNotes = ref('')
const errorMessage = ref('')
const errors = ref({
  warehouseManager: ''
})

// Computed
const selectedManagerInfo = computed(() => {
  if (!selectedWarehouseManager.value) return null
  return warehouseManagers.value.find(m => m.rut === selectedWarehouseManager.value)
})

// Métodos
const loadWarehouseManagers = async () => {
  loading.value = true
  errorMessage.value = ''
  
  try {
    console.log('Cargando encargados de bodega...')
    const result = await userService.getUsersByRole('encargado de bodega')
    console.log('Resultado de getUsersByRole:', result)
    
    if (result.success && result.data) {
      warehouseManagers.value = result.data
      console.log('Encargados cargados:', warehouseManagers.value)
      
      if (warehouseManagers.value.length === 0) {
        errorMessage.value = 'No hay encargados de bodega registrados en el sistema'
      }
    } else {
      warehouseManagers.value = []
      errorMessage.value = result.error || 'No se pudieron cargar los encargados de bodega'
      console.error('Error en respuesta:', result.error)
    }
  } catch (error) {
    console.error('Error cargando encargados de bodega:', error)
    errorMessage.value = 'Error al cargar la lista de encargados de bodega'
    warehouseManagers.value = []
  } finally {
    loading.value = false
  }
}

const validateForm = () => {
  errors.value = {
    warehouseManager: ''
  }

  if (!selectedWarehouseManager.value) {
    errors.value.warehouseManager = 'Debe seleccionar un encargado de bodega'
    return false
  }

  return true
}

const handleAssign = async () => {
  if (!validateForm()) {
    return
  }

  loading.value = true
  errorMessage.value = ''

  try {
    const selectedManager = selectedManagerInfo.value
    
    const assignmentData = {
      assigned_to: selectedManager.rut,
      assigned_to_name: selectedManager.name,
      assigned_by_pavedad: authStore.getUserRut,
      assigned_by_pavedad_name: authStore.getUserName,
      pavedad_notes: pavedadNotes.value || ''
    }

    await supplyRequestService.assignRequestToWarehouseManager(props.request.id, assignmentData)

    // Mostrar mensaje de éxito
    await Swal.fire({
      icon: 'success',
      title: 'Solicitud Asignada',
      text: `La solicitud ha sido asignada exitosamente a ${selectedManager.name}`,
      timer: 2000,
      showConfirmButton: false
    })

    emit('assigned')
    closeModal()
    
    // Recargar la página completa
    window.location.reload()
  } catch (error) {
    console.error('Error asignando solicitud:', error)
    errorMessage.value = error.response?.data?.error || error.message || 'Error al asignar la solicitud'
    
    Swal.fire({
      icon: 'error',
      title: 'Error',
      text: errorMessage.value
    })
  } finally {
    loading.value = false
  }
}

const closeModal = () => {
  if (!loading.value) {
    resetForm()
    emit('close')
  }
}

const resetForm = () => {
  selectedWarehouseManager.value = ''
  pavedadNotes.value = ''
  errorMessage.value = ''
  errors.value = {
    warehouseManager: ''
  }
}

// Watchers
watch(() => props.show, (newValue) => {
  if (newValue) {
    loadWarehouseManagers()
    resetForm()
  }
})
</script>

<style scoped>
/* Animación de entrada del modal */
@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.fixed.inset-0 {
  animation: fadeIn 0.2s ease-out;
}

.inline-block {
  animation: slideIn 0.3s ease-out;
}

/* Estilos para select */
select:focus {
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

/* Estilos para textarea */
textarea:focus {
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

/* Animación del spinner */
@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.animate-spin {
  animation: spin 1s linear infinite;
}
</style>
