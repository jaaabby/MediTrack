<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Panel de ayuda colapsible -->
    <div v-if="showHelp" class="bg-blue-50 border-b border-blue-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
        <div class="bg-blue-100 rounded-lg p-4">
          <h3 class="text-lg font-semibold text-blue-900 mb-2">💡 Guía para crear solicitudes</h3>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm text-blue-800">
            <div>
              <h4 class="font-semibold mb-1">Información básica:</h4>
              <ul class="list-disc list-inside space-y-1">
                <li>Selecciona el pabellón que solicita</li>
                <li>Indica la fecha y hora de la cirugía</li>
                <li>Agrega observaciones si es necesario</li>
              </ul>
            </div>
            <div>
              <h4 class="font-semibold mb-1">Insumos solicitados:</h4>
              <ul class="list-disc list-inside space-y-1">
                <li>Puedes agregar múltiples insumos</li>
                <li>Especifica cantidades y características</li>
                <li>Marca si es para uso pediátrico</li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Contenido principal -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
      <!-- Componente del formulario -->
      <SupplyRequestForm :id="id" :editMode="editMode" @success="handleSuccess" @cancel="goBack" />
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import SupplyRequestForm from '@/components/SupplyRequestForm.vue'

const router = useRouter()

// Props recibidos desde las rutas
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

// Estado reactivo
const showHelp = ref(false)

// Métodos
const goBack = () => {
  router.push('/supply-requests')
}

const handleSuccess = (requestData) => {
  // Redirigir a la página de éxito con los datos de la solicitud
  router.push({
    name: 'SupplyRequestSuccess',
    params: { requestData }
  })
}
</script>

<style scoped>
/* Animaciones y transiciones */
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from, .fade-leave-to {
  opacity: 0;
}

/* Mejoras de UI */
.bg-gray-50 {
  background-color: #f9fafb;
}

.shadow-lg {
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
}

/* Hover effects */
.hover\:text-blue-600:hover {
  color: #2563eb;
}

.hover\:text-gray-600:hover {
  color: #4b5563;
}

.hover\:text-gray-900:hover {
  color: #111827;
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .max-w-7xl {
    padding-left: 1rem;
    padding-right: 1rem;
  }
}
</style>