<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <!-- Header -->
      <div>
        <div class="mx-auto h-12 w-12 flex items-center justify-center rounded-full bg-blue-100">
          <svg class="h-8 w-8 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
          </svg>
        </div>
        <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
          Iniciar Sesión
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600">
          Sistema de Trazabilidad MediTrack
        </p>
      </div>

      <!-- Formulario de Login -->
      <form class="mt-8 space-y-6" @submit.prevent="handleLogin">
        <div class="rounded-md shadow-sm -space-y-px">
          <!-- Campo Email -->
          <div>
            <label for="email" class="sr-only">Email</label>
            <input
              id="email"
              v-model="loginForm.email"
              name="email"
              type="email"
              autocomplete="email"
              required
              class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-t-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 focus:z-10 sm:text-sm"
              :class="{ 'border-red-500': errors.email }"
              placeholder="Correo electrónico"
            />
            <p v-if="errors.email" class="mt-1 text-sm text-red-600">{{ errors.email }}</p>
          </div>

          <!-- Campo Contraseña -->
          <div>
            <label for="password" class="sr-only">Contraseña</label>
            <input
              id="password"
              v-model="loginForm.password"
              name="password"
              type="password"
              autocomplete="current-password"
              required
              class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-b-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 focus:z-10 sm:text-sm"
              :class="{ 'border-red-500': errors.password }"
              placeholder="Contraseña"
            />
            <p v-if="errors.password" class="mt-1 text-sm text-red-600">{{ errors.password }}</p>
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
                Error de autenticación
              </h3>
              <div class="mt-2 text-sm text-red-700">
                {{ errorMessage }}
              </div>
            </div>
          </div>
        </div>

        <!-- Botón de Login -->
        <div>
          <button
            type="submit"
            :disabled="isLoading"
            class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            <span v-if="isLoading" class="absolute left-0 inset-y-0 flex items-center pl-3">
              <svg class="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
            </span>
            {{ isLoading ? 'Iniciando sesión...' : 'Iniciar Sesión' }}
          </button>
        </div>

        <!-- Enlace al registro -->
        <div class="text-center">
          <router-link
            to="/register"
            class="text-sm text-blue-600 hover:text-blue-500 font-medium"
          >
            ¿No tienes cuenta? Regístrate
          </router-link>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

// Estado del formulario
const loginForm = reactive({
  email: '',
  password: ''
})

// Estado de la UI
const isLoading = ref(false)
const errorMessage = ref('')
const errors = reactive({
  email: '',
  password: ''
})

// Validación del formulario
const validateForm = () => {
  errors.email = ''
  errors.password = ''
  
  if (!loginForm.email) {
    errors.email = 'El email es requerido'
    return false
  }
  
  if (!loginForm.email.includes('@')) {
    errors.email = 'El email debe ser válido'
    return false
  }
  
  if (!loginForm.password) {
    errors.password = 'La contraseña es requerida'
    return false
  }
  
  if (loginForm.password.length < 6) {
    errors.password = 'La contraseña debe tener al menos 6 caracteres'
    return false
  }
  
  return true
}

// Manejo del login
const handleLogin = async () => {
  if (!validateForm()) {
    return
  }
  
  isLoading.value = true
  errorMessage.value = ''
  
  try {
    await authStore.login(loginForm.email, loginForm.password)
    
    // Redirigir al usuario después del login exitoso
    const redirectTo = router.currentRoute.value.query.redirect || '/home'
    console.log('Redirigiendo a:', redirectTo)
    
    // Usar replace en lugar de push para evitar que el usuario vuelva al login con el botón atrás
    await router.replace(redirectTo)
    console.log('Redirección completada')
    
  } catch (error) {
    console.error('Error en login:', error)
    errorMessage.value = error.message || 'Error al iniciar sesión. Verifica tus credenciales.'
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
/* Estilos adicionales si es necesario */
</style>
