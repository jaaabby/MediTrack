<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <div>
        <h2 class="text-2xl font-semibold text-gray-900">Centros Médicos</h2>
        <p class="text-gray-600 mt-1">Configura el correo de alertas de cada centro médico</p>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="card">
      <div class="flex justify-center items-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        <span class="ml-3 text-gray-600">Cargando...</span>
      </div>
    </div>

    <!-- Error -->
    <div v-else-if="loadError" class="card">
      <div class="bg-red-50 border border-red-200 rounded-md p-4">
        <p class="text-red-700 text-sm">{{ loadError }}</p>
        <button @click="loadCenters" class="btn-secondary mt-3 text-sm">Reintentar</button>
      </div>
    </div>

    <!-- Lista de centros -->
    <div v-else class="space-y-4">
      <div
        v-for="center in centers"
        :key="center.id"
        class="bg-white rounded-lg shadow-sm border p-6"
      >
        <div class="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4">
          <!-- Info del centro -->
          <div class="flex-1 min-w-0">
            <h3 class="text-lg font-semibold text-gray-900">{{ center.name }}</h3>
            <p v-if="center.address" class="text-sm text-gray-500 mt-0.5">{{ center.address }}</p>
            <p v-if="center.phone" class="text-sm text-gray-500">{{ center.phone }}</p>
          </div>

          <!-- Formulario de correo -->
          <div class="w-full sm:w-96 flex-shrink-0">
            <label class="block text-sm font-medium text-gray-700 mb-1">
              Correo de alertas
            </label>
            <div class="flex gap-2">
              <input
                v-model="editValues[center.id]"
                type="email"
                placeholder="alertas@clinica.cl"
                class="form-input flex-1 min-w-0"
                @keyup.enter="save(center)"
              />
              <button
                @click="save(center)"
                :disabled="saving[center.id]"
                class="btn-primary px-4 flex-shrink-0 flex items-center gap-1"
              >
                <svg v-if="saving[center.id]" class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
                </svg>
                <span>{{ saving[center.id] ? 'Guardando...' : 'Guardar' }}</span>
              </button>
            </div>
            <!-- Feedback -->
            <p v-if="successMsg[center.id]" class="text-green-600 text-xs mt-1">{{ successMsg[center.id] }}</p>
            <p v-if="errorMsg[center.id]" class="text-red-600 text-xs mt-1">{{ errorMsg[center.id] }}</p>
            <p v-if="!center.alert_email" class="text-amber-600 text-xs mt-1">
              ⚠ Sin correo configurado — las alertas no se enviarán
            </p>
          </div>
        </div>
      </div>

      <div v-if="centers.length === 0" class="bg-white rounded-lg border p-8 text-center text-gray-500">
        No hay centros médicos registrados.
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import medicalCenterService from '@/services/config/medicalCenterService'

const centers = ref([])
const loading = ref(true)
const loadError = ref(null)

// Estado por centro
const editValues = reactive({})
const saving = reactive({})
const successMsg = reactive({})
const errorMsg = reactive({})

async function loadCenters() {
  loading.value = true
  loadError.value = null
  const result = await medicalCenterService.getAllMedicalCenters()
  if (result.success) {
    centers.value = result.data
    result.data.forEach(c => {
      editValues[c.id] = c.alert_email || ''
    })
  } else {
    loadError.value = result.error || 'Error al cargar los centros médicos'
  }
  loading.value = false
}

async function save(center) {
  saving[center.id] = true
  successMsg[center.id] = null
  errorMsg[center.id] = null

  const payload = {
    name: center.name,
    address: center.address,
    phone: center.phone,
    email: center.email,
    alert_email: editValues[center.id]?.trim() || ''
  }

  const result = await medicalCenterService.update(center.id, payload)
  if (result.success) {
    center.alert_email = payload.alert_email
    successMsg[center.id] = 'Correo guardado correctamente'
    setTimeout(() => { successMsg[center.id] = null }, 3000)
  } else {
    errorMsg[center.id] = result.error || 'Error al guardar'
  }
  saving[center.id] = false
}

onMounted(loadCenters)
</script>
