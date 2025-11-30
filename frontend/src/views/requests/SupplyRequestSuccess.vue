<template>
  <div class="min-h-screen bg-gray-50 flex items-center justify-center">
    <div class="max-w-md w-full bg-white rounded-lg shadow-lg p-8">
      <!-- Ícono de éxito -->
      <div class="flex justify-center mb-6">
        <div class="bg-green-100 p-4 rounded-full">
          <svg class="h-16 w-16 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
      </div>

      <!-- Título y mensaje -->
      <div class="text-center mb-6">
        <h1 class="text-2xl font-bold text-gray-900 mb-2">
          ¡Solicitud Creada Exitosamente!
        </h1>
        <p class="text-gray-600">
          Tu solicitud de insumos ha sido registrada correctamente en el sistema.
        </p>
      </div>

      <!-- Número de Solicitud con botón para copiar -->
      <div v-if="requestNumber" class="bg-blue-50 border-2 border-blue-200 rounded-lg p-4 mb-8">
        <div class="text-center">
          <p class="text-sm font-medium text-blue-900 mb-2">Número de Solicitud</p>
          <div class="flex items-center justify-center gap-3">
            <span class="text-2xl font-bold text-blue-600 font-mono">{{ requestNumber }}</span>
            <button
              @click="copyRequestNumber"
              class="p-2 text-blue-600 hover:bg-blue-100 rounded-lg transition-colors"
              title="Copiar número de solicitud"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
            </button>
          </div>
        </div>
      </div>

      <!-- Información de la solicitud 
      <div v-if="requestData || requestId" class="bg-gray-50 rounded-lg p-4 mb-6">
        <h3 class="text-sm font-semibold text-gray-700 mb-3">Detalles de la solicitud:</h3>
        <div class="space-y-2 text-sm">
          <div class="flex justify-between">
            <span class="text-gray-600">Número de solicitud:</span>
            <span class="font-medium">{{ requestData?.request_number || requestData?.id || requestId }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-gray-600">Estado:</span>
            <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800">
              Pendiente
            </span>
          </div>
          <div v-if="requestData?.priority" class="flex justify-between">
            <span class="text-gray-600">Prioridad:</span>
            <span class="font-medium capitalize">{{ getPriorityLabel(requestData.priority) }}</span>
          </div>
          <div v-if="requestData?.pavilion_id || requestData?.pavilion_name" class="flex justify-between">
            <span class="text-gray-600">Pabellón:</span>
            <span class="font-medium">{{ requestData?.pavilion_name || getPavilionName(requestData?.pavilion_id) }}</span>
          </div>
        </div>
      </div>

      <!-- Próximos pasos 
      <div class="bg-blue-50 rounded-lg p-4 mb-6">
        <h3 class="text-sm font-semibold text-blue-900 mb-2">Próximos pasos:</h3>
        <ul class="text-sm text-blue-800 space-y-1">
          <li>• La solicitud será revisada por el personal autorizado</li>
          <li>• Recibirás notificaciones sobre el estado de tu solicitud</li>
          <li>• Una vez aprobada, se asignarán códigos QR a los insumos</li>
        </ul>
      </div>-->

      <!-- Botones de acción -->
      <div class="flex flex-col space-y-3">
        <button
          @click="viewRequestDetails"
          class="btn-primary w-full"
        >
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
          </svg>
          Ver Detalles de la Solicitud
        </button>

        <div class="flex space-x-3">
          <button
            v-if="authStore.canCreateRequests"
            @click="createAnotherRequest"
            class="flex-1 bg-gray-100 text-gray-700 py-2 px-4 rounded-lg font-medium hover:bg-gray-200 transition-colors flex items-center justify-center"
          >
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
            Nueva Solicitud
          </button>

          <button
            @click="goToRequestsList"
            :class="authStore.canCreateRequests ? 'flex-1' : 'w-full'"
            class="bg-gray-100 text-gray-700 py-2 px-4 rounded-lg font-medium hover:bg-gray-200 transition-colors flex items-center justify-center"
          >
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
            </svg>
            Ver Solicitudes
          </button>
        </div>

        <button
          @click="goToHome"
          class="w-full text-gray-500 py-2 px-4 rounded-lg font-medium hover:text-gray-700 transition-colors flex items-center justify-center"
        >
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
          </svg>
          Volver al Inicio
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import supplyRequestService from '@/services/requests/supplyRequestService'
import pavilionService from '@/services/config/pavilionService'
import Swal from 'sweetalert2'

// Props
const props = defineProps({
  requestId: {
    type: Number,
    default: null
  },
  requestData: {
    type: Object,
    default: null
  }
})

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const requestData = ref(null)
const requestId = ref(null)
const requestNumber = ref(null)
const pavilions = ref([])

// Cargar datos de la solicitud desde los parámetros de la ruta
onMounted(async () => {
  // Obtener el ID de la solicitud (desde props o desde params)
  requestId.value = props.requestId || parseInt(route.params.id)
  
  // Obtener datos de la solicitud desde el state de la navegación
  if (props.requestData || route.params.requestData) {
    requestData.value = props.requestData || route.params.requestData
    // Extraer el número de solicitud
    requestNumber.value = requestData.value?.request_number || 
                          requestData.value?.request?.request_number ||
                          requestData.value?.RequestNumber
  }
  
  // Si no tenemos el número de solicitud pero tenemos el ID, cargar los datos
  if (!requestNumber.value && requestId.value) {
    try {
      const response = await supplyRequestService.getSupplyRequestById(requestId.value)
      if (response.success && response.data) {
        requestData.value = response.data.request
        requestNumber.value = response.data.request?.request_number
      }
    } catch (error) {
      console.error('Error cargando datos de la solicitud:', error)
    }
  }
  
  // Cargar pabellones para mostrar nombres
  try {
    pavilions.value = await pavilionService.getAllPavilions()
  } catch (error) {
    console.error('Error cargando pabellones:', error)
  }
})

// Copiar número de solicitud al portapapeles
const copyRequestNumber = async () => {
  if (!requestNumber.value) return
  
  try {
    await navigator.clipboard.writeText(requestNumber.value)
    Swal.fire({
      icon: 'success',
      title: 'Copiado',
      text: `Número de solicitud ${requestNumber.value} copiado al portapapeles`,
      toast: true,
      position: 'top-end',
      showConfirmButton: false,
      timer: 2000,
      timerProgressBar: true
    })
  } catch (err) {
    console.error('Error copiando al portapapeles:', err)
    Swal.fire({
      icon: 'error',
      title: 'Error',
      text: 'No se pudo copiar al portapapeles',
      toast: true,
      position: 'top-end',
      showConfirmButton: false,
      timer: 2000
    })
  }
}

// Métodos auxiliares
const getPriorityLabel = (priority) => {
  return supplyRequestService.getPriorityLabel(priority)
}

const getPavilionName = (pavilionId) => {
  // Buscar el pavilion_id en diferentes ubicaciones posibles
  const id = pavilionId || 
             requestData.value?.pavilion_id || 
             requestData.value?.request?.pavilion_id
  
  if (!id) return 'No especificado'
  
  const pavilion = pavilions.value.find(p => p.id === id)
  return pavilion ? pavilion.name : `Pabellón ${id}`
}

// Métodos de navegación
const viewRequestDetails = () => {
  if (requestId.value) {
    router.push({
      name: 'SupplyRequestDetails',
      params: { id: requestId.value }
    })
  } else {
    console.warn('No se encontró el ID de la solicitud, navegando al listado')
    router.push('/supply-requests')
  }
}

const createAnotherRequest = () => {
  router.push('/supply-requests/new')
}

const goToRequestsList = () => {
  router.push('/supply-requests')
}

const goToHome = () => {
  router.push('/home')
}
</script>

<style scoped>
/* Animaciones */
.animate-pulse {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: .5;
  }
}

/* Efectos hover */
.hover\:bg-blue-700:hover {
  background-color: #1d4ed8;
}

.hover\:bg-gray-200:hover {
  background-color: #e5e7eb;
}

.hover\:text-gray-700:hover {
  color: #374151;
}
</style>