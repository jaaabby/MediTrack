<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <!-- Header -->
      <div>
        <div class="mx-auto flex items-center justify-center mb-8">
          <img
            src="@/assets/images/MEDITRACK_LOGO.svg"
            alt="Meditrack"
            class="h-24 sm:h-32 md:h-40 w-auto"
          />
        </div>
        <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
          Verificación en dos pasos
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600">
          Ingresa el código de 6 dígitos de tu aplicación autenticadora
        </p>
      </div>

      <form class="mt-8 space-y-6" @submit.prevent="handleVerify">
        <div class="space-y-4">
          <!-- Campo código TOTP -->
          <div>
            <label for="code" class="block text-sm font-medium text-gray-700 mb-1">
              Código de verificación
            </label>
            <input
              id="code"
              v-model="code"
              name="code"
              type="text"
              inputmode="numeric"
              autocomplete="one-time-code"
              maxlength="6"
              placeholder="000000"
              class="appearance-none relative block w-full px-3 py-3 border placeholder-gray-400 text-gray-900 rounded-md focus:outline-none focus:ring-brand-blue-dark focus:border-brand-blue-dark text-center text-2xl tracking-widest font-mono"
              :class="errorMessage ? 'border-red-500' : 'border-gray-300'"
              @input="code = code.replace(/\D/g, '')"
            />
          </div>
        </div>

        <!-- Error -->
        <div v-if="errorMessage" class="rounded-md bg-red-50 p-4">
          <div class="flex">
            <div class="flex-shrink-0">
              <svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
              </svg>
            </div>
            <div class="ml-3">
              <p class="text-sm text-red-700">{{ errorMessage }}</p>
            </div>
          </div>
        </div>

        <!-- Botón verificar -->
        <div>
          <button
            type="submit"
            :disabled="isLoading || code.length !== 6"
            class="btn-primary w-full disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span v-if="isLoading" class="absolute left-0 inset-y-0 flex items-center pl-3">
              <svg class="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
            </span>
            {{ isLoading ? 'Verificando...' : 'Verificar' }}
          </button>
        </div>

        <!-- Volver al login -->
        <div class="text-center">
          <button
            type="button"
            class="text-sm font-medium text-brand-blue-dark hover:text-brand-blue"
            @click="handleCancel"
          >
            Volver al inicio de sesión
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const code = ref('')
const isLoading = ref(false)
const errorMessage = ref('')

const handleVerify = async () => {
  if (code.value.length !== 6) return

  isLoading.value = true
  errorMessage.value = ''

  try {
    await authStore.verifyTOTP(code.value)

    if (authStore.user?.must_change_password) {
      await router.replace('/first-time-password-change')
      return
    }

    const redirectTo = router.currentRoute.value.query.redirect || '/home'
    await router.replace(redirectTo)
  } catch (error) {
    errorMessage.value = error.message || 'Código inválido. Intenta de nuevo.'
    code.value = ''
  } finally {
    isLoading.value = false
  }
}

const handleCancel = () => {
  authStore.logout()
  router.replace('/login')
}
</script>
