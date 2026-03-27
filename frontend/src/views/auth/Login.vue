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
          {{ step === 'credentials' ? 'Iniciar Sesión' : 'Verificación SMS' }}
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600">
          {{ step === 'credentials'
            ? 'Sistema de Trazabilidad MediTrack'
            : `Ingresa el código enviado al número ${phoneMasked}` }}
        </p>
      </div>

      <!-- PASO 1: Formulario de Login -->
      <form v-if="step === 'credentials'" class="mt-8 space-y-6" @submit.prevent="handleLogin">
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
        <div v-if="errorMessage" class="rounded-md bg-red-50 p-4">
          <div class="flex">
            <div class="flex-shrink-0">
              <svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
              </svg>
            </div>
            <div class="ml-3">
              <h3 class="text-sm font-medium text-red-800">Error de autenticación</h3>
              <div class="mt-2 text-sm text-red-700">{{ errorMessage }}</div>
            </div>
          </div>
        </div>

        <!-- Botón de Login -->
        <div>
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
            {{ isLoading ? 'Verificando...' : 'Iniciar Sesión' }}
          </button>
        </div>
      </form>

      <!-- PASO 2: Formulario de verificación OTP -->
      <form v-else class="mt-8 space-y-6" @submit.prevent="handleVerifyOtp">
        <div class="space-y-4">
          <!-- Campos individuales para el código de 6 dígitos -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-3 text-center">
              Código de verificación (6 dígitos)
            </label>
            <div class="flex justify-center gap-2">
              <input
                v-for="(_, i) in otpDigits"
                :key="i"
                :ref="el => otpInputs[i] = el"
                v-model="otpDigits[i]"
                type="text"
                inputmode="numeric"
                maxlength="1"
                @input="onOtpInput(i)"
                @keydown="onOtpKeydown($event, i)"
                @paste="onOtpPaste($event)"
                class="w-11 h-12 text-center text-lg font-semibold border rounded-md focus:outline-none focus:ring-2 focus:ring-brand-blue-dark focus:border-brand-blue-dark"
                :class="errors.otp ? 'border-red-500' : 'border-gray-300'"
              />
            </div>
            <p v-if="errors.otp" class="mt-2 text-sm text-red-600 font-medium text-center">{{ errors.otp }}</p>
          </div>
        </div>

        <!-- Mensaje de Error -->
        <div v-if="errorMessage" class="rounded-md bg-red-50 p-4">
          <div class="flex">
            <div class="flex-shrink-0">
              <svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
              </svg>
            </div>
            <div class="ml-3">
              <h3 class="text-sm font-medium text-red-800">Error de verificación</h3>
              <div class="mt-2 text-sm text-red-700">{{ errorMessage }}</div>
            </div>
          </div>
        </div>

        <!-- Botón verificar -->
        <div>
          <button
            type="submit"
            :disabled="isLoading || otpCode.length < 6"
            class="btn-primary w-full disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span v-if="isLoading" class="absolute left-0 inset-y-0 flex items-center pl-3">
              <svg class="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
            </span>
            {{ isLoading ? 'Verificando...' : 'Verificar código' }}
          </button>
        </div>

        <!-- Volver al login -->
        <div class="text-center">
          <button
            type="button"
            class="text-sm font-medium text-brand-blue-dark hover:text-brand-blue"
            @click="backToCredentials"
          >
            Volver e ingresar de nuevo
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, watch, computed, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

// Estado del formulario de credenciales
const loginForm = reactive({ email: '', password: '' })

// Estado del flujo OTP
const step = ref('credentials') // 'credentials' | 'otp'
const otpSessionId = ref('')
const phoneMasked = ref('')
const otpDigits = reactive(['', '', '', '', '', ''])
const otpInputs = ref([])

const otpCode = computed(() => otpDigits.join(''))

// Estado de la UI
const isLoading = ref(false)
const errorMessage = ref('')
const errors = reactive({ email: '', password: '', otp: '' })

// Limpiar errores al modificar campos
watch(() => loginForm.email, () => { errors.email = ''; errorMessage.value = '' })
watch(() => loginForm.password, () => { errors.password = ''; errorMessage.value = '' })

const validateEmail = () => {
  if (!loginForm.email) return
  const emailRegex = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
  if (!emailRegex.test(loginForm.email)) {
    errors.email = 'El email debe ser válido'
  }
}

const validateForm = () => {
  errors.email = ''
  errors.password = ''
  const emailRegex = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
  if (!loginForm.email) { errors.email = 'El email es requerido'; return false }
  if (!emailRegex.test(loginForm.email)) { errors.email = 'El email debe ser válido'; return false }
  if (!loginForm.password) { errors.password = 'La contraseña es requerida'; return false }
  return true
}

// Manejo de inputs OTP individuales
const onOtpInput = (index) => {
  const val = otpDigits[index]
  // Solo permitir dígitos
  if (!/^\d*$/.test(val)) {
    otpDigits[index] = ''
    return
  }
  // Avanzar al siguiente campo
  if (val && index < 5) {
    nextTick(() => otpInputs.value[index + 1]?.focus())
  }
  errors.otp = ''
  errorMessage.value = ''
}

const onOtpKeydown = (event, index) => {
  if (event.key === 'Backspace' && !otpDigits[index] && index > 0) {
    nextTick(() => otpInputs.value[index - 1]?.focus())
  }
}

const onOtpPaste = (event) => {
  event.preventDefault()
  const pasted = event.clipboardData.getData('text').replace(/\D/g, '').slice(0, 6)
  for (let i = 0; i < 6; i++) {
    otpDigits[i] = pasted[i] || ''
  }
  nextTick(() => {
    const lastFilled = Math.min(pasted.length, 5)
    otpInputs.value[lastFilled]?.focus()
  })
}

// PASO 1: Login con credenciales
const handleLogin = async () => {
  if (!validateForm()) return

  isLoading.value = true
  errorMessage.value = ''

  try {
    const result = await authStore.login(loginForm.email, loginForm.password)

    // Si el backend retorna otp_session_id → requiere verificación SMS
    if (result?.otp_session_id) {
      otpSessionId.value = result.otp_session_id
      phoneMasked.value = result.phone_masked || ''
      step.value = 'otp'
      nextTick(() => otpInputs.value[0]?.focus())
      return
    }

    // Login directo (sin teléfono registrado)
    if (authStore.user?.must_change_password) {
      await router.replace('/first-time-password-change')
      return
    }

    const redirectTo = router.currentRoute.value.query.redirect || '/home'
    await router.replace(redirectTo)

  } catch (error) {
    errorMessage.value = error.message || 'Las credenciales ingresadas son inválidas. Verifica tu correo y contraseña.'
  } finally {
    isLoading.value = false
  }
}

// PASO 2: Verificar código OTP
const handleVerifyOtp = async () => {
  errors.otp = ''
  if (otpCode.value.length < 6) {
    errors.otp = 'Ingresa los 6 dígitos del código'
    return
  }

  isLoading.value = true
  errorMessage.value = ''

  try {
    await authStore.verifyOtp(otpSessionId.value, otpCode.value)

    if (authStore.user?.must_change_password) {
      await router.replace('/first-time-password-change')
      return
    }

    const redirectTo = router.currentRoute.value.query.redirect || '/home'
    await router.replace(redirectTo)

  } catch (error) {
    errorMessage.value = error.message || 'Código incorrecto. Inténtalo de nuevo.'
    // Limpiar campos OTP
    for (let i = 0; i < 6; i++) otpDigits[i] = ''
    nextTick(() => otpInputs.value[0]?.focus())
  } finally {
    isLoading.value = false
  }
}

const backToCredentials = () => {
  step.value = 'credentials'
  otpSessionId.value = ''
  phoneMasked.value = ''
  for (let i = 0; i < 6; i++) otpDigits[i] = ''
  errorMessage.value = ''
  errors.otp = ''
}
</script>

<style scoped>
/* Estilos adicionales si es necesario */
</style>
