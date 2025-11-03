<template>
  <div class="bg-white rounded-lg shadow-sm border p-6 mb-6">
    <h3 class="text-lg font-medium text-gray-900 mb-4 flex items-center">
      <svg class="h-5 w-5 mr-2 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
      </svg>
      Ubicación Actual
    </h3>
    
    <div class="space-y-4">
      <!-- Selector de Centro Médico -->
      <div class="relative">
        <label class="block text-sm font-medium text-gray-700 mb-2">
          Centro Médico
        </label>
        <div class="relative">
          <input
            v-model="medicalCenterSearch"
            @input="onMedicalCenterSearch"
            @focus="showMedicalCenterOptions = true"
            @blur="hideMedicalCenterOptions"
            type="text"
            placeholder="Escribir o seleccionar centro médico..."
            class="form-input w-full"
            :disabled="loadingCenters"
          />
          <div class="absolute inset-y-0 right-0 pr-3 flex items-center">
            <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </div>
        </div>
        
        <!-- Dropdown de opciones -->
        <div
          v-show="showMedicalCenterOptions && filteredMedicalCenters.length > 0"
          class="absolute z-10 mt-1 w-full bg-white shadow-lg max-h-60 rounded-md py-1 text-base ring-1 ring-black ring-opacity-5 overflow-auto focus:outline-none sm:text-sm"
        >
          <div
            v-for="center in filteredMedicalCenters"
            :key="center.id"
            @mousedown="selectMedicalCenter(center)"
            class="cursor-pointer select-none relative py-2 pl-3 pr-9 hover:bg-purple-50"
          >
            <div class="flex items-center">
              <span class="font-medium text-gray-900 block truncate">
                {{ center.name }}
              </span>
            </div>
            <div v-if="center.address" class="text-sm text-gray-600 truncate">
              {{ center.address }}
            </div>
          </div>
        </div>
      </div>
      
      <!-- Selector de Pabellón -->
      <div class="relative">
        <label class="block text-sm font-medium text-gray-700 mb-2">
          Pabellón
        </label>
        <div class="relative">
          <input
            v-model="pavilionSearch"
            @input="onPavilionSearch"
            @focus="showPavilionOptions = true"
            @blur="hidePavilionOptions"
            type="text"
            placeholder="Escribir o seleccionar pabellón..."
            class="form-input w-full"
            :disabled="!selectedMedicalCenter || loadingPavilions"
          />
          <div class="absolute inset-y-0 right-0 pr-3 flex items-center">
            <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </div>
        </div>
        
        <!-- Dropdown de opciones -->
        <div
          v-show="showPavilionOptions && filteredPavilions.length > 0"
          class="absolute z-10 mt-1 w-full bg-white shadow-lg max-h-60 rounded-md py-1 text-base ring-1 ring-black ring-opacity-5 overflow-auto focus:outline-none sm:text-sm"
        >
          <div
            v-for="pavilion in filteredPavilions"
            :key="pavilion.id"
            @mousedown="selectPavilion(pavilion)"
            class="cursor-pointer select-none relative py-2 pl-3 pr-9 hover:bg-purple-50"
          >
            <div class="flex items-center">
              <span class="font-medium text-gray-900 block truncate">
                {{ pavilion.name }}
              </span>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Información de ubicación actual -->
      <div v-if="currentLocation" class="bg-purple-50 border border-purple-200 rounded-lg p-4">
        <div class="flex items-center">
          <svg class="h-5 w-5 text-purple-600 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <div class="flex-1">
            <p class="text-sm font-medium text-purple-900">Ubicación Configurada</p>
            <p class="text-sm text-purple-700">
              {{ currentLocation.medical_center_name }} - {{ currentLocation.pavilion_name }}
            </p>
          </div>
        </div>
        
        <div class="mt-3 flex items-center justify-between">
          <div class="text-xs text-purple-600">
            Configurado: {{ formatDate(currentLocation.set_at) }}
          </div>
          <button
            @click="clearLocation"
            class="text-xs text-purple-600 hover:text-purple-800 underline"
          >
            Cambiar ubicación
          </button>
        </div>
      </div>
      
      <!-- Botones de acción -->
      <div class="flex space-x-3">
        <button
          @click="saveLocation"
          :disabled="!selectedMedicalCenter || !selectedPavilion || saving"
          class="btn-primary flex-1"
        >
          <svg v-if="saving" class="animate-spin -ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <svg v-else class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
          {{ saving ? 'Guardando...' : 'Confirmar Ubicación' }}
        </button>
        
        <button
          @click="useCurrentGPS"
          :disabled="gettingGPS"
          class="btn-secondary"
        >
          <svg v-if="gettingGPS" class="animate-spin -ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <svg v-else class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
          </svg>
          GPS
        </button>
      </div>
      
      <!-- Ubicaciones recientes -->
      <div v-if="recentLocations.length > 0" class="border-t border-gray-200 pt-4">
        <h4 class="text-sm font-medium text-gray-700 mb-2">Ubicaciones Recientes</h4>
        <div class="space-y-2">
          <button
            v-for="location in recentLocations"
            :key="`${location.medical_center_id}-${location.pavilion_id}`"
            @click="selectRecentLocation(location)"
            class="w-full text-left p-2 rounded-md border border-gray-200 hover:border-gray-300 hover:bg-gray-50 text-sm"
          >
            <div class="font-medium text-gray-900">{{ location.medical_center_name }}</div>
            <div class="text-gray-600">{{ location.pavilion_name }}</div>
            <div class="text-xs text-gray-500">{{ formatDate(location.last_used) }}</div>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import medicalCenterService from '@/services/medicalCenterService'
import pavilionService from '@/services/pavilionService'
import Swal from 'sweetalert2'

// Emits
const emit = defineEmits(['location-changed'])

// Estado reactivo
const medicalCenters = ref([])
const pavilions = ref([])
const selectedMedicalCenter = ref('')
const selectedPavilion = ref('')
const currentLocation = ref(null)
const recentLocations = ref([])

// Variables para búsqueda con autocompletado
const medicalCenterSearch = ref('')
const pavilionSearch = ref('')
const showMedicalCenterOptions = ref(false)
const showPavilionOptions = ref(false)

const loadingCenters = ref(false)
const loadingPavilions = ref(false)
const saving = ref(false)
const gettingGPS = ref(false)

// Computed
const availablePavilions = computed(() => {
  if (!selectedMedicalCenter.value) return []
  return pavilions.value.filter(p => p.medical_center_id === parseInt(selectedMedicalCenter.value))
})

const filteredMedicalCenters = computed(() => {
  if (!medicalCenterSearch.value.trim()) {
    return medicalCenters.value
  }
  const search = medicalCenterSearch.value.toLowerCase().trim()
  return medicalCenters.value.filter(center => 
    center.name.toLowerCase().includes(search) ||
    (center.address && center.address.toLowerCase().includes(search))
  )
})

const filteredPavilions = computed(() => {
  const available = availablePavilions.value
  if (!pavilionSearch.value.trim()) {
    return available
  }
  const search = pavilionSearch.value.toLowerCase().trim()
  return available.filter(pavilion => 
    pavilion.name.toLowerCase().includes(search)
  )
})

// Watchers
watch(selectedMedicalCenter, (newValue) => {
  selectedPavilion.value = ''
  if (newValue) {
    loadPavilions()
  }
})

// Funciones principales
const loadMedicalCenters = async () => {
  console.log('Loading medical centers...')
  loadingCenters.value = true
  try {
    const response = await medicalCenterService.getAll()
    console.log('Medical centers response:', response)
    medicalCenters.value = response.data || response || []
    console.log('Medical centers loaded:', medicalCenters.value)
  } catch (error) {
    console.error('Error loading medical centers:', error)
  } finally {
    loadingCenters.value = false
  }
}

const loadPavilions = async () => {
  if (!selectedMedicalCenter.value) return
  console.log('Loading pavilions for medical center:', selectedMedicalCenter.value)
  loadingPavilions.value = true
  try {
    let response = []
    // Si existe el método getByMedicalCenter, úsalo, si no, usa getAllPavilions
    if (typeof pavilionService.getByMedicalCenter === 'function') {
      console.log('Using getByMedicalCenter method')
      response = await pavilionService.getByMedicalCenter(selectedMedicalCenter.value)
      response = response.data || response.Data || response || []
    } else if (typeof pavilionService.getAllPavilions === 'function') {
      console.log('Using getAllPavilions method')
      response = await pavilionService.getAllPavilions()
      response = response.data || response.Data || response || []
    }
    console.log('Pavilions response processed:', response)
    pavilions.value = response
  } catch (error) {
    console.error('Error loading pavilions:', error)
    pavilions.value = []
  } finally {
    loadingPavilions.value = false
  }
}

const onMedicalCenterChange = () => {
  selectedPavilion.value = ''
  pavilionSearch.value = ''
  updateLocationInfo()
}

const onPavilionChange = () => {
  updateLocationInfo()
}

// Nuevas funciones para búsqueda con autocompletado
const onMedicalCenterSearch = () => {
  showMedicalCenterOptions.value = true
  // Limpiar selección si el texto no coincide exactamente
  const exactMatch = medicalCenters.value.find(center => 
    center.name.toLowerCase() === medicalCenterSearch.value.toLowerCase()
  )
  if (exactMatch) {
    selectedMedicalCenter.value = exactMatch.id.toString()
  } else {
    selectedMedicalCenter.value = ''
  }
  onMedicalCenterChange()
}

const onPavilionSearch = () => {
  showPavilionOptions.value = true
  // Limpiar selección si el texto no coincide exactamente
  const exactMatch = availablePavilions.value.find(pavilion => 
    pavilion.name.toLowerCase() === pavilionSearch.value.toLowerCase()
  )
  if (exactMatch) {
    selectedPavilion.value = exactMatch.id.toString()
  } else {
    selectedPavilion.value = ''
  }
  onPavilionChange()
}

const selectMedicalCenter = (center) => {
  medicalCenterSearch.value = center.name
  selectedMedicalCenter.value = center.id.toString()
  showMedicalCenterOptions.value = false
  onMedicalCenterChange()
}

const selectPavilion = (pavilion) => {
  pavilionSearch.value = pavilion.name
  selectedPavilion.value = pavilion.id.toString()
  showPavilionOptions.value = false
  onPavilionChange()
}

const hideMedicalCenterOptions = () => {
  setTimeout(() => {
    showMedicalCenterOptions.value = false
  }, 200)
}

const hidePavilionOptions = () => {
  setTimeout(() => {
    showPavilionOptions.value = false
  }, 200)
}

const updateLocationInfo = () => {
  if (!selectedMedicalCenter.value || !selectedPavilion.value) {
    emit('location-changed', null)
    return
  }
  
  const medicalCenter = medicalCenters.value.find(c => c.id === parseInt(selectedMedicalCenter.value))
  const pavilion = availablePavilions.value.find(p => p.id === parseInt(selectedPavilion.value))
  
  if (medicalCenter && pavilion) {
    const location = {
      medical_center_id: medicalCenter.id,
      medical_center_name: medicalCenter.name,
      pavilion_id: pavilion.id,
      pavilion_name: pavilion.name,
      name: `${medicalCenter.name} - ${pavilion.name}`
    }
    
    emit('location-changed', location)
  }
}

const saveLocation = () => {
  if (!selectedMedicalCenter.value || !selectedPavilion.value) return
  
  saving.value = true
  
  const medicalCenter = medicalCenters.value.find(c => c.id === parseInt(selectedMedicalCenter.value))
  const pavilion = availablePavilions.value.find(p => p.id === parseInt(selectedPavilion.value))
  
  if (medicalCenter && pavilion) {
    const location = {
      medical_center_id: medicalCenter.id,
      medical_center_name: medicalCenter.name,
      pavilion_id: pavilion.id,
      pavilion_name: pavilion.name,
      name: `${medicalCenter.name} - ${pavilion.name}`,
      set_at: new Date().toISOString()
    }
    
    currentLocation.value = location
    
    // Guardar en localStorage
    localStorage.setItem('current-location', JSON.stringify(location))
    
    // Agregar a ubicaciones recientes
    addToRecentLocations(location)
    
    // Emitir cambio
    emit('location-changed', location)
    
    setTimeout(() => {
      saving.value = false
    }, 500)
  }
}

const clearLocation = () => {
  currentLocation.value = null
  selectedMedicalCenter.value = ''
  selectedPavilion.value = ''
  medicalCenterSearch.value = ''
  pavilionSearch.value = ''
  localStorage.removeItem('current-location')
  emit('location-changed', null)
}

const useCurrentGPS = async () => {
  gettingGPS.value = true
  
  try {
    if (!navigator.geolocation) {
      throw new Error('Geolocalización no soportada por este navegador')
    }
    
    const position = await new Promise((resolve, reject) => {
      navigator.geolocation.getCurrentPosition(resolve, reject, {
        enableHighAccuracy: true,
        timeout: 10000,
        maximumAge: 300000 // 5 minutes
      })
    })
    
    const { latitude, longitude } = position.coords
    
    // Aquí podrías hacer una llamada al backend para encontrar la ubicación más cercana
    // Por ahora, mostraremos un mensaje informativo
    Swal.fire({
      icon: 'info',
      title: 'Ubicación GPS obtenida',
      html: `${latitude.toFixed(6)}, ${longitude.toFixed(6)}<br><br>Funcionalidad de auto-detección en desarrollo.`
    })
    
  } catch (error) {
    console.error('Error getting GPS location:', error)
    Swal.fire({
      icon: 'error',
      title: 'Error',
      text: 'No se pudo obtener la ubicación GPS. Seleccione manualmente.'
    })
  } finally {
    gettingGPS.value = false
  }
}

const selectRecentLocation = (location) => {
  selectedMedicalCenter.value = location.medical_center_id.toString()
  selectedPavilion.value = location.pavilion_id.toString()
  
  // También actualizar los campos de búsqueda
  medicalCenterSearch.value = location.medical_center_name
  pavilionSearch.value = location.pavilion_name
  
  // Actualizar fecha de último uso
  location.last_used = new Date().toISOString()
  addToRecentLocations(location)
  
  updateLocationInfo()
}

const addToRecentLocations = (location) => {
  const locationKey = `${location.medical_center_id}-${location.pavilion_id}`
  
  // Remover si ya existe
  recentLocations.value = recentLocations.value.filter(
    l => `${l.medical_center_id}-${l.pavilion_id}` !== locationKey
  )
  
  // Agregar al inicio
  recentLocations.value.unshift({
    ...location,
    last_used: location.last_used || new Date().toISOString()
  })
  
  // Mantener solo las últimas 5
  recentLocations.value = recentLocations.value.slice(0, 5)
  
  // Guardar en localStorage
  localStorage.setItem('recent-locations', JSON.stringify(recentLocations.value))
}

const loadSavedLocation = () => {
  try {
    // Cargar ubicación actual
    const saved = localStorage.getItem('current-location')
    if (saved) {
      const location = JSON.parse(saved)
      currentLocation.value = location
      selectedMedicalCenter.value = location.medical_center_id.toString()
      selectedPavilion.value = location.pavilion_id.toString()
      medicalCenterSearch.value = location.medical_center_name
      pavilionSearch.value = location.pavilion_name
      emit('location-changed', location)
    }
    
    // Cargar ubicaciones recientes
    const recent = localStorage.getItem('recent-locations')
    if (recent) {
      recentLocations.value = JSON.parse(recent)
    }
  } catch (error) {
    console.error('Error loading saved location:', error)
  }
}

const formatDate = (dateString) => {
  if (!dateString) return 'No disponible'
  try {
    return format(new Date(dateString), 'dd/MM/yyyy HH:mm', { locale: es })
  } catch (error) {
    return dateString
  }
}

// Lifecycle
onMounted(async () => {
  await loadMedicalCenters()
  loadSavedLocation()
})
</script>

<style scoped>
/* Usar clases de formularios y botones de style.css global */
</style>