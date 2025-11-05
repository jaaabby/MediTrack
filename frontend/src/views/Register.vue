<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <!-- Header -->
      <div>
        <div class="mx-auto h-12 w-12 flex items-center justify-center rounded-full bg-green-100">
          <svg class="h-8 w-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
          </svg>
        </div>
        <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
          Registro de Usuario
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600">
          Crear nueva cuenta en MediTrack
        </p>
      </div>

      <!-- Formulario de Registro -->
      <form class="mt-8 space-y-6" @submit.prevent="handleRegister">
        <div class="space-y-4">
          <!-- Campo RUT -->
          <div>
            <label for="rut" class="block text-sm font-medium text-gray-700">RUT</label>
            <input
              id="rut"
              v-model="registerForm.rut"
              name="rut"
              type="text"
              required
              class="form-input"
              :class="{ 'border-red-500': errors.rut }"
              placeholder="12345678-9"
            />
            <p v-if="errors.rut" class="mt-1 text-sm text-red-600">{{ errors.rut }}</p>
          </div>

          <!-- Campo Nombre -->
          <div>
            <label for="name" class="block text-sm font-medium text-gray-700">Nombre Completo</label>
            <input
              id="name"
              v-model="registerForm.name"
              name="name"
              type="text"
              required
              class="form-input"
              :class="{ 'border-red-500': errors.name }"
              placeholder="Juan Pérez"
            />
            <p v-if="errors.name" class="mt-1 text-sm text-red-600">{{ errors.name }}</p>
          </div>

          <!-- Campo Email -->
          <div>
            <label for="email" class="block text-sm font-medium text-gray-700">Correo Electrónico</label>
            <input
              id="email"
              v-model="registerForm.email"
              name="email"
              type="email"
              autocomplete="email"
              required
              class="form-input"
              :class="{ 'border-red-500': errors.email }"
              placeholder="juan@example.com"
            />
            <p v-if="errors.email" class="mt-1 text-sm text-red-600">{{ errors.email }}</p>
          </div>

          <!-- Campo Contraseña -->
          <div>
            <label for="password" class="block text-sm font-medium text-gray-700">Contraseña</label>
            <input
              id="password"
              v-model="registerForm.password"
              name="password"
              type="password"
              autocomplete="new-password"
              required
              class="form-input"
              :class="{ 'border-red-500': errors.password }"
              placeholder="Mínimo 6 caracteres"
            />
            <p v-if="errors.password" class="mt-1 text-sm text-red-600">{{ errors.password }}</p>
          </div>

          <!-- Campo Confirmar Contraseña -->
          <div>
            <label for="confirmPassword" class="block text-sm font-medium text-gray-700">Confirmar Contraseña</label>
            <input
              id="confirmPassword"
              v-model="registerForm.confirmPassword"
              name="confirmPassword"
              type="password"
              autocomplete="new-password"
              required
              class="form-input"
              :class="{ 'border-red-500': errors.confirmPassword }"
              placeholder="Repetir contraseña"
            />
            <p v-if="errors.confirmPassword" class="mt-1 text-sm text-red-600">{{ errors.confirmPassword }}</p>
          </div>

          <!-- Campo Rol -->
          <div>
            <label for="role" class="block text-sm font-medium text-gray-700">Rol</label>
            <select
              id="role"
              v-model="registerForm.role"
              name="role"
              required
              class="form-select"
              :class="{ 'border-red-500': errors.role }"
            >
              <option value="">Seleccionar rol</option>
              <option value="admin">Administrador</option>
              <option value="pabellón">Pabellón</option>
              <option value="encargado de bodega">Encargado de Bodega</option>
              <option value="enfermera">Enfermera</option>
              <option value="doctor">Doctor</option>
            </select>
            <p v-if="errors.role" class="mt-1 text-sm text-red-600">{{ errors.role }}</p>
          </div>

          <!-- Campo Centro Médico -->
          <div>
            <label for="medicalCenterId" class="block text-sm font-medium text-gray-700">Centro Médico</label>
            <select
              id="medicalCenterId"
              v-model="registerForm.medicalCenterId"
              name="medicalCenterId"
              required
              :disabled="loadingCenters"
              class="form-select"
              :class="{ 'border-red-500': errors.medicalCenterId, 'opacity-50 cursor-not-allowed': loadingCenters }"
            >
              <option value="">{{ loadingCenters ? 'Cargando centros médicos...' : 'Seleccionar centro médico' }}</option>
              <option v-for="center in medicalCenters" :key="center.id" :value="center.id">
                {{ center.name }}
              </option>
            </select>
            <p v-if="errors.medicalCenterId" class="mt-1 text-sm text-red-600">{{ errors.medicalCenterId }}</p>
            <p v-if="loadingCenters" class="mt-1 text-sm text-gray-500">Cargando centros médicos...</p>
          </div>
        </div>

        <!-- Mensaje de Error General -->
        <div v-if="errorMessage" class="rounded-md bg-red-50 p-4">
          <div class="flex">
            <div class="flex-shrink-0">
              <svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
              </svg>
            </div>
            <div class="ml-3">
              <h3 class="text-sm font-medium text-red-800">
                Error en el registro
              </h3>
              <div class="mt-2 text-sm text-red-700">
                {{ errorMessage }}
              </div>
            </div>
          </div>
        </div>

        <!-- Mensaje de Éxito -->
        <div v-if="successMessage" class="rounded-md bg-green-50 p-4">
          <div class="flex">
            <div class="flex-shrink-0">
              <svg class="h-5 w-5 text-green-400" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
              </svg>
            </div>
            <div class="ml-3">
              <h3 class="text-sm font-medium text-green-800">
                Registro exitoso
              </h3>
              <div class="mt-2 text-sm text-green-700">
                {{ successMessage }}
              </div>
            </div>
          </div>
        </div>

        <!-- Botones -->
        <div class="space-y-3">
          <button
            type="submit"
            :disabled="isLoading"
            class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            <span v-if="isLoading" class="absolute left-0 inset-y-0 flex items-center pl-3">
              <svg class="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
            </span>
            {{ isLoading ? 'Registrando...' : 'Registrar Usuario' }}
          </button>

          <div class="text-center">
            <router-link
              to="/login"
              class="text-sm text-green-600 hover:text-green-500 font-medium"
            >
              ¿Ya tienes cuenta? Iniciar sesión
            </router-link>
          </div>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import authService from '@/services/authService'
import medicalCenterService from '@/services/medicalCenterService'

const router = useRouter()

// Estado del formulario
const registerForm = reactive({
  rut: '',
  name: '',
  email: '',
  password: '',
  confirmPassword: '',
  role: '',
  medicalCenterId: ''
})

// Estado de la UI
const isLoading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const medicalCenters = ref([])
const loadingCenters = ref(false)

// Errores de validación
const errors = reactive({
  rut: '',
  name: '',
  email: '',
  password: '',
  confirmPassword: '',
  role: '',
  medicalCenterId: ''
})

// Cargar centros médicos al montar el componente
onMounted(async () => {
  await loadMedicalCenters()
})

// Cargar lista de centros médicos desde el backend
const loadMedicalCenters = async () => {
  loadingCenters.value = true
  errorMessage.value = ''
  try {
    const response = await medicalCenterService.getAll()
    // Manejar diferentes formatos de respuesta
    medicalCenters.value = response.data || response || []
    console.log('Centros médicos cargados:', medicalCenters.value)
    
    if (medicalCenters.value.length === 0) {
      errorMessage.value = 'No hay centros médicos disponibles. Contacta al administrador.'
    }
  } catch (error) {
    console.error('Error al cargar centros médicos:', error)
    
    // Mensajes de error más específicos
    if (error.code === 'ERR_NETWORK' || error.message?.includes('ERR_CONNECTION_REFUSED')) {
      errorMessage.value = 'No se puede conectar al servidor. Verifica que el backend esté corriendo en http://localhost:8080'
    } else if (error.response?.status === 401) {
      errorMessage.value = 'Error de autenticación. Por favor, recarga la página.'
    } else if (error.response?.status === 404) {
      errorMessage.value = 'Endpoint no encontrado. Verifica la configuración del servidor.'
    } else {
      errorMessage.value = 'Error al cargar los centros médicos. Por favor, recarga la página o verifica que el servidor esté corriendo.'
    }
  } finally {
    loadingCenters.value = false
  }
}

// Validación del formulario
const validateForm = () => {
  // Limpiar errores anteriores
  Object.keys(errors).forEach(key => {
    errors[key] = ''
  })
  
  let isValid = true

  // Validar RUT
  if (!registerForm.rut) {
    errors.rut = 'El RUT es requerido'
    isValid = false
  } else if (!/^[0-9]+-[0-9kK]$/.test(registerForm.rut)) {
    errors.rut = 'El RUT debe tener el formato 12345678-9'
    isValid = false
  }

  // Validar nombre
  if (!registerForm.name) {
    errors.name = 'El nombre es requerido'
    isValid = false
  } else if (registerForm.name.length < 2) {
    errors.name = 'El nombre debe tener al menos 2 caracteres'
    isValid = false
  }

  // Validar email
  if (!registerForm.email) {
    errors.email = 'El email es requerido'
    isValid = false
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(registerForm.email)) {
    errors.email = 'El email debe ser válido'
    isValid = false
  }

  // Validar contraseña
  if (!registerForm.password) {
    errors.password = 'La contraseña es requerida'
    isValid = false
  } else if (registerForm.password.length < 6) {
    errors.password = 'La contraseña debe tener al menos 6 caracteres'
    isValid = false
  }

  // Validar confirmación de contraseña
  if (!registerForm.confirmPassword) {
    errors.confirmPassword = 'Debe confirmar la contraseña'
    isValid = false
  } else if (registerForm.password !== registerForm.confirmPassword) {
    errors.confirmPassword = 'Las contraseñas no coinciden'
    isValid = false
  }

  // Validar rol
  if (!registerForm.role) {
    errors.role = 'Debe seleccionar un rol'
    isValid = false
  }

  // Validar centro médico
  if (!registerForm.medicalCenterId) {
    errors.medicalCenterId = 'Debe seleccionar un centro médico'
    isValid = false
  }

  return isValid
}

// Manejo del registro
const handleRegister = async () => {
  if (!validateForm()) {
    return
  }
  
  isLoading.value = true
  errorMessage.value = ''
  successMessage.value = ''
  
  try {
    const response = await authService.register({
      rut: registerForm.rut,
      name: registerForm.name,
      email: registerForm.email,
      password: registerForm.password,
      role: registerForm.role,
      medical_center_id: parseInt(registerForm.medicalCenterId)
    })
    
    successMessage.value = 'Usuario registrado exitosamente. Serás redirigido al inicio de sesión.'
    
    // Limpiar formulario
    Object.keys(registerForm).forEach(key => {
      registerForm[key] = ''
    })
    
    // Redirigir al login después de 2 segundos
    setTimeout(() => {
      router.push('/login')
    }, 2000)
    
  } catch (error) {
    console.error('Error en registro:', error)
    // Manejar diferentes tipos de errores
    if (error.message) {
      errorMessage.value = error.message
    } else if (error.error) {
      errorMessage.value = error.error
    } else {
      errorMessage.value = 'Error al registrar usuario. Por favor, verifica los datos e intenta nuevamente.'
    }
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
/* Estilos adicionales si es necesario */
</style>
