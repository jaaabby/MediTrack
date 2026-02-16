<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <!-- Header -->
      <div>
        <div class="mx-auto flex items-center justify-center mb-8">
          <img 
            src="@/assets/images/MEDITRACK_LOGO.svg" 
            alt="MediTrack" 
            class="h-24 sm:h-32 md:h-40 w-auto"
          />
        </div>
        <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
          Restablecer Contraseña
        </h2>
        <p v-if="userInfo" class="mt-2 text-center text-sm text-gray-600">
          {{ userInfo.name }} ({{ userInfo.email }})
        </p>
      </div>

      <!-- Estado de carga inicial -->
      <div v-if="isValidating" class="text-center">
        <svg class="animate-spin h-10 w-10 mx-auto text-brand-blue-dark" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <p class="mt-4 text-gray-600">Validando token...</p>
      </div>

      <!-- Token inválido -->
      <div v-if="!isValidating && !tokenValid" class="rounded-md bg-red-50 p-4">
        <div class="flex">
          <div class="flex-shrink-0">
            <svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
            </svg>
          </div>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-red-800">
              Token inválido o expirado
            </h3>
            <div class="mt-2 text-sm text-red-700">
              El enlace de recuperación ha expirado o no es válido. Por favor, solicita un nuevo enlace.
            </div>
          </div>
        </div>
        <div class="mt-4">
          <router-link
            to="/forgot-password"
            class="block text-center w-full py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-brand-blue-dark hover:bg-brand-blue focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-brand-blue-dark"
          >
            Solicitar nuevo enlace
          </router-link>
        </div>
      </div>

      <!-- Mensaje de Éxito -->
      <div v-if="passwordChanged" class="rounded-md bg-green-50 p-4">
        <div class="flex">
          <div class="flex-shrink-0">
            <svg class="h-5 w-5 text-green-400" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
            </svg>
          </div>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-green-800">
              Contraseña actualizada
            </h3>
            <div class="mt-2 text-sm text-green-700">
              Tu contraseña ha sido actualizada exitosamente. Ya puedes iniciar sesión con tu nueva contraseña.
            </div>
          </div>
        </div>
        <div class="mt-4">
          <router-link
            to="/login"
            class="block text-center w-full py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-brand-blue-dark hover:bg-brand-blue focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-brand-blue-dark"
          >
            Ir al inicio de sesión
          </router-link>
        </div>
      </div>

      <!-- Formulario -->
      <form v-if="!isValidating && tokenValid && !passwordChanged" class="mt-8 space-y-6" @submit.prevent="handleSubmit">
        <div class="space-y-4">
          <!-- Nueva Contraseña -->
          <div>
            <label for="password" class="block text-sm font-medium text-gray-700 mb-2">
              Nueva Contraseña
            </label>
            <input
              id="password"
              v-model="form.password"
              name="password"
              type="password"
              autocomplete="new-password"
              required
              class="appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-brand-blue-dark focus:border-brand-blue-dark focus:z-10 sm:text-sm"
              :class="{ 'border-red-500': errors.password }"
              placeholder="Mínimo 6 caracteres"
            />
            <p v-if="errors.password" class="mt-1 text-sm text-red-600">{{ errors.password }}</p>
            
            <!-- Barra de fortaleza -->
            <div v-if="form.password" class="mt-3">
              <div class="flex items-center justify-between mb-1">
                <span class="text-xs font-medium text-gray-700">Fortaleza de la contraseña:</span>
                <span class="text-xs font-semibold" :class="passwordStrength.colorClass">
                  {{ passwordStrength.label }}
                </span>
              </div>
              <div class="w-full bg-gray-200 rounded-full h-2 overflow-hidden">
                <div 
                  class="h-full rounded-full transition-all duration-300 ease-in-out"
                  :class="passwordStrength.bgClass"
                  :style="{ width: passwordStrength.width }"
                ></div>
              </div>
            </div>
          </div>

          <!-- Confirmar Contraseña -->
          <div>
            <label for="confirmPassword" class="block text-sm font-medium text-gray-700 mb-2">
              Confirmar Nueva Contraseña
            </label>
            <input
              id="confirmPassword"
              v-model="form.confirmPassword"
              name="confirmPassword"
              type="password"
              autocomplete="new-password"
              required
              class="appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-brand-blue-dark focus:border-brand-blue-dark focus:z-10 sm:text-sm"
              :class="{ 'border-red-500': errors.confirmPassword }"
              placeholder="Repite tu nueva contraseña"
            />
            <p v-if="errors.confirmPassword" class="mt-1 text-sm text-red-600">{{ errors.confirmPassword }}</p>
          </div>
        </div>

        <!-- Requisitos de contraseña -->
        <div class="rounded-md bg-gray-50 p-4">
          <h4 class="text-sm font-medium text-gray-700 mb-2">Requisitos de contraseña:</h4>
          <ul class="text-xs text-gray-600 space-y-1 list-disc list-inside">
            <li>Mínimo 6 caracteres</li>
            <li>Recomendamos usar letras, números y símbolos</li>
          </ul>
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
                Error
              </h3>
              <div class="mt-2 text-sm text-red-700">
                {{ errorMessage }}
              </div>
            </div>
          </div>
        </div>

        <!-- Botones -->
        <div class="space-y-3">
          <button
            type="submit"
            :disabled="isLoading"
            class="btn-primary w-full disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span v-if="isLoading" class="absolute left-0 inset-y-0 flex items-center pl-3">
              <svg class="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
            </span>
            {{ isLoading ? 'Actualizando...' : 'Restablecer contraseña' }}
          </button>

          <router-link
            to="/login"
            class="block text-center w-full py-2 px-4 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-brand-blue-dark"
          >
            Cancelar
          </router-link>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import authService from '@/services/auth/authService'

const route = useRoute()
const router = useRouter()

const token = ref('')
const isValidating = ref(true)
const tokenValid = ref(false)
const userInfo = ref(null)
const passwordChanged = ref(false)
const isLoading = ref(false)
const errorMessage = ref('')

const form = reactive({
  password: '',
  confirmPassword: ''
})

const errors = reactive({
  password: '',
  confirmPassword: ''
})

// Calcular fortaleza de la contraseña
const passwordStrength = computed(() => {
  const password = form.password
  if (!password) {
    return { width: '0%', label: '', colorClass: '', bgClass: '' }
  }
  
  let strength = 0
  
  // Longitud
  if (password.length >= 6) strength += 20
  if (password.length >= 8) strength += 10
  if (password.length >= 10) strength += 10
  if (password.length >= 12) strength += 10
  
  // Contiene minúsculas
  if (/[a-z]/.test(password)) strength += 10
  
  // Contiene mayúsculas
  if (/[A-Z]/.test(password)) strength += 15
  
  // Contiene números
  if (/[0-9]/.test(password)) strength += 15
  
  // Contiene símbolos
  if (/[^a-zA-Z0-9]/.test(password)) strength += 20
  
  // Determinar nivel y estilos
  if (strength < 40) {
    return {
      width: `${strength}%`,
      label: 'Muy débil',
      colorClass: 'text-red-600',
      bgClass: 'bg-red-500'
    }
  } else if (strength < 60) {
    return {
      width: `${strength}%`,
      label: 'Débil',
      colorClass: 'text-orange-600',
      bgClass: 'bg-orange-500'
    }
  } else if (strength < 80) {
    return {
      width: `${strength}%`,
      label: 'Media',
      colorClass: 'text-yellow-600',
      bgClass: 'bg-yellow-500'
    }
  } else if (strength < 100) {
    return {
      width: `${strength}%`,
      label: 'Fuerte',
      colorClass: 'text-green-600',
      bgClass: 'bg-green-500'
    }
  } else {
    return {
      width: '100%',
      label: 'Muy fuerte',
      colorClass: 'text-green-700',
      bgClass: 'bg-green-600'
    }
  }
})

// Validar token al cargar la página
onMounted(async () => {
  token.value = route.query.token || ''
  
  if (!token.value) {
    isValidating.value = false
    tokenValid.value = false
    return
  }
  
  try {
    const response = await authService.validateResetToken(token.value)
    tokenValid.value = true
    userInfo.value = response
  } catch (error) {
    tokenValid.value = false
  } finally {
    isValidating.value = false
  }
})

const validateForm = () => {
  errors.password = ''
  errors.confirmPassword = ''
  let isValid = true
  
  if (!form.password) {
    errors.password = 'La contraseña es requerida'
    isValid = false
  } else if (form.password.length < 6) {
    errors.password = 'La contraseña debe tener al menos 6 caracteres'
    isValid = false
  }
  
  if (!form.confirmPassword) {
    errors.confirmPassword = 'Debes confirmar tu contraseña'
    isValid = false
  } else if (form.password !== form.confirmPassword) {
    errors.confirmPassword = 'Las contraseñas no coinciden'
    isValid = false
  }
  
  return isValid
}

const handleSubmit = async () => {
  errorMessage.value = ''
  
  if (!validateForm()) {
    return
  }
  
  isLoading.value = true
  
  try {
    await authService.resetPassword(token.value, form.password)
    passwordChanged.value = true
    
    // Opcionalmente redirigir automáticamente después de 3 segundos
    setTimeout(() => {
      router.push('/login')
    }, 3000)
  } catch (error) {
    errorMessage.value = error.message || 'Error al restablecer la contraseña'
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
/* Estilos adicionales si son necesarios */
</style>
