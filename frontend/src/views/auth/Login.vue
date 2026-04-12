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
          Iniciar Sesión
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600">
          Sistema de Trazabilidad MediTrack
        </p>
      </div>

      <!-- Formulario de Login -->
      <form class="mt-8 space-y-6" @submit.prevent="handleLogin">
        <div class="space-y-4">
          <!-- Campo Email -->
          <div>
            <label for="email" class="sr-only">Email</label>
            <input
              id="email"
              v-model="loginForm.email"
              name="email"
              type="email"
              autocomplete="email"
              @blur="validateEmail"
              class="appearance-none relative block w-full px-3 py-2 border placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-brand-blue-dark focus:border-brand-blue-dark focus:z-10 sm:text-sm"
              :class="errors.email ? 'border-red-500' : 'border-gray-300'"
              placeholder="Correo electrónico"
            />
            <p v-if="errors.email" class="mt-2 text-sm text-red-600 font-medium">{{ errors.email }}</p>
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
              class="appearance-none relative block w-full px-3 py-2 border placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-brand-blue-dark focus:border-brand-blue-dark focus:z-10 sm:text-sm"
              :class="errors.password ? 'border-red-500' : 'border-gray-300'"
              placeholder="Contraseña"
            />
            <p v-if="errors.password" class="mt-2 text-sm text-red-600 font-medium">{{ errors.password }}</p>
          </div>
        </div>

        <!-- Enlace para recuperación de contraseña -->
        <div class="flex items-center justify-end">
          <router-link
            to="/forgot-password"
            class="text-sm font-medium text-brand-blue-dark hover:text-brand-blue"
          >
            ¿Olvidaste tu contraseña?
          </router-link>
        </div>

        <!-- Mensaje de Error General -->
        <div v-if="lockoutSeconds > 0" class="rounded-md bg-yellow-50 border border-yellow-200 p-4">
          <div class="flex">
            <div class="flex-shrink-0">
              <svg class="h-5 w-5 text-yellow-500" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z" clip-rule="evenodd" />
              </svg>
            </div>
            <div class="ml-3">
              <h3 class="text-sm font-medium text-yellow-800">Cuenta bloqueada</h3>
              <div class="mt-1 text-sm text-yellow-700">
                Demasiados intentos fallidos. Podrás volver a intentarlo en
                <span class="font-semibold">{{ lockoutDisplay }}</span>.
              </div>
            </div>
          </div>
        </div>
        <div v-else-if="errorMessage" class="rounded-md bg-red-50 p-4">
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
            :disabled="isLoading || lockoutSeconds > 0"
            class="btn-primary w-full disabled:opacity-50 disabled:cursor-not-allowed"
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

        <!-- Botón de Passkey (solo si el dispositivo soporta WebAuthn) -->
        <div v-if="webAuthnSupported">
          <div class="relative">
            <div class="absolute inset-0 flex items-center">
              <div class="w-full border-t border-gray-300" />
            </div>
            <div class="relative flex justify-center text-sm">
              <span class="px-2 bg-gray-50 text-gray-500">o</span>
            </div>
          </div>
          <button
            type="button"
            :disabled="isPasskeyLoading"
            class="mt-3 w-full flex items-center justify-center gap-2 py-2 px-4 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-brand-blue-dark disabled:opacity-50 disabled:cursor-not-allowed"
            @click="handlePasskeyLogin"
          >
            <svg class="h-5 w-5 text-gray-500" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
            </svg>
            {{ isPasskeyLoading ? 'Verificando...' : 'Iniciar sesión con llave de acceso' }}
          </button>
        </div>

      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, watch, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

// Detectar soporte de WebAuthn para mostrar el botón de passkey
const webAuthnSupported = typeof window !== 'undefined' && !!window.PublicKeyCredential
const isPasskeyLoading = ref(false)

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

// Estado de bloqueo de cuenta
const lockoutSeconds = ref(0)
let countdownTimer = null

const lockoutDisplay = computed(() => {
  const min = Math.floor(lockoutSeconds.value / 60)
  const sec = lockoutSeconds.value % 60
  return min > 0
    ? `${min}:${String(sec).padStart(2, '0')} min`
    : `${sec} seg`
})

const startCountdown = (seconds) => {
  lockoutSeconds.value = seconds
  if (countdownTimer) clearInterval(countdownTimer)
  countdownTimer = setInterval(() => {
    lockoutSeconds.value--
    if (lockoutSeconds.value <= 0) {
      clearInterval(countdownTimer)
      countdownTimer = null
    }
  }, 1000)
}

onUnmounted(() => {
  if (countdownTimer) clearInterval(countdownTimer)
})

// Limpiar errores cuando el usuario modifica los campos
watch(() => loginForm.email, () => {
  errors.email = ''
  errorMessage.value = ''
})

watch(() => loginForm.password, () => {
  errors.password = ''
  errorMessage.value = ''
})
// Validar email cuando el usuario sale del campo
const validateEmail = () => {
  if (!loginForm.email) {
    return // No mostrar error si está vacío al salir del campo
  }
  
  const emailRegex = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
  
  if (!emailRegex.test(loginForm.email)) {
    errors.email = 'El email debe ser válido'
  }
}

// Validación del formulario
const validateForm = () => {
  errors.email = ''
  errors.password = ''
  
  // Regex para validar email
  const emailRegex = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
  
  if (!loginForm.email) {
    errors.email = 'El email es requerido'
    return false
  }
  
  if (!emailRegex.test(loginForm.email)) {
    errors.email = 'El email debe ser válido'
    return false
  }
  
  if (!loginForm.password) {
    errors.password = 'La contraseña es requerida'
    return false
  }
  
  return true
}

// Login explícito con passkey (triggered por el usuario)
const handlePasskeyLogin = async () => {
  isPasskeyLoading.value = true
  errorMessage.value = ''
  try {
    await authStore.loginWithPasskey()

    if (authStore.user?.must_change_password) {
      await router.replace('/first-time-password-change')
      return
    }

    const redirectTo = router.currentRoute.value.query.redirect || '/home'
    await router.replace(redirectTo)
  } catch (error) {
    if (error?.message !== 'NO_CREDENTIALS' && error?.message !== 'NO_SUPPORT') {
      errorMessage.value = error?.message || 'No se pudo autenticar con la llave de acceso.'
    }
  } finally {
    isPasskeyLoading.value = false
  }
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

    // Si se requiere TOTP, redirigir a la verificación
    if (authStore.totpRequired) {
      await router.replace('/totp-verify')
      return
    }

    // Verificar si el usuario debe cambiar su contraseña
    if (authStore.user?.must_change_password) {
      await router.replace('/first-time-password-change')
      return
    }

    // Redirigir al usuario después del login exitoso
    const redirectTo = router.currentRoute.value.query.redirect || '/home'
    console.log('Redirigiendo a:', redirectTo)

    // Usar replace en lugar de push para evitar que el usuario vuelva al login con el botón atrás
    await router.replace(redirectTo)
    console.log('Redirección completada')
    
  } catch (error) {
    console.error('Error en login:', error)
    if (error.responseData?.remaining_seconds) {
      startCountdown(Number(error.responseData.remaining_seconds))
      errorMessage.value = ''
    } else {
      errorMessage.value = error.message || 'Las credenciales ingresadas son inválidas. Verifica tu correo y contraseña.'
    }
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
/* Estilos adicionales si es necesario */
</style>
