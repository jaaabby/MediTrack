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
          Recuperar Contraseña
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600">
          Ingresa tu correo electrónico y te enviaremos instrucciones para restablecer tu contraseña
        </p>
      </div>

      <!-- Mensaje de Éxito -->
      <div v-if="emailSent" class="rounded-md bg-green-50 p-4">
        <div class="flex">
          <div class="flex-shrink-0">
            <svg class="h-5 w-5 text-green-400" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
            </svg>
          </div>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-green-800">
              Correo enviado
            </h3>
            <div class="mt-2 text-sm text-green-700">
              Si el correo existe en nuestro sistema, recibirás un enlace para restablecer tu contraseña.
              Por favor, revisa tu bandeja de entrada y carpeta de spam.
            </div>
          </div>
        </div>
      </div>

      <!-- Formulario -->
      <form v-if="!emailSent" class="mt-8 space-y-6" @submit.prevent="handleSubmit">
        <div>
          <label for="email" class="block text-sm font-medium text-gray-700 mb-2">
            Correo electrónico
          </label>
          <input
            id="email"
            v-model="email"
            name="email"
            type="email"
            autocomplete="email"
            required
            class="appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-brand-blue-dark focus:border-brand-blue-dark focus:z-10 sm:text-sm"
            :class="{ 'border-red-500': errors.email }"
            placeholder="tu-email@ejemplo.com"
          />
          <p v-if="errors.email" class="mt-1 text-sm text-red-600">{{ errors.email }}</p>
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
            {{ isLoading ? 'Enviando...' : 'Enviar instrucciones' }}
          </button>

          <router-link
            to="/login"
            class="block text-center w-full py-2 px-4 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-brand-blue-dark"
          >
            Volver al inicio de sesión
          </router-link>
        </div>
      </form>

      <!-- Botón para volver después de enviar -->
      <div v-if="emailSent" class="mt-6">
        <router-link
          to="/login"
          class="block text-center w-full py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-brand-blue-dark hover:bg-brand-blue focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-brand-blue-dark"
        >
          Volver al inicio de sesión
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import authService from '@/services/auth/authService'

const router = useRouter()

const email = ref('')
const isLoading = ref(false)
const emailSent = ref(false)
const errorMessage = ref('')
const errors = reactive({
  email: ''
})

const validateEmail = () => {
  errors.email = ''
  
  // Regex para validar email
  const emailRegex = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
  
  if (!email.value) {
    errors.email = 'El correo electrónico es requerido'
    return false
  }
  
  if (!emailRegex.test(email.value)) {
    errors.email = 'Ingresa un correo electrónico válido'
    return false
  }
  
  return true
}

const handleSubmit = async () => {
  errorMessage.value = ''
  
  if (!validateEmail()) {
    return
  }
  
  isLoading.value = true
  
  try {
    await authService.requestPasswordReset(email.value)
    emailSent.value = true
  } catch (error) {
    errorMessage.value = error.message || 'Error al enviar el correo de recuperación'
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
/* Estilos adicionales si son necesarios */
</style>
