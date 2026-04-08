<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <!-- Header -->
      <div class="text-center">
        <h2 class="text-3xl font-extrabold text-gray-900">
          Configurar verificación en dos pasos
        </h2>
        <p class="mt-2 text-sm text-gray-600">
          Escanea el código QR con Google Authenticator, Authy u otra aplicación TOTP
        </p>
      </div>

      <!-- Paso 1: QR Code -->
      <div v-if="step === 'scan'" class="space-y-6">
        <div v-if="isLoadingSetup" class="flex justify-center py-8">
          <svg class="animate-spin h-10 w-10 text-brand-blue-dark" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
        </div>

        <div v-else-if="setupData" class="space-y-4">
          <!-- QR Code -->
          <div class="flex justify-center">
            <div class="p-4 bg-white border border-gray-200 rounded-lg shadow-sm">
              <qrcode-vue :value="setupData.qr_url" :size="200" level="M" />
            </div>
          </div>

          <!-- Secreto manual -->
          <div class="rounded-md bg-gray-50 p-4">
            <p class="text-xs text-gray-500 mb-1">
              ¿No puedes escanear? Ingresa este código manualmente:
            </p>
            <div class="flex items-center gap-2">
              <code class="flex-1 text-sm font-mono bg-white border border-gray-300 rounded px-3 py-2 break-all">
                {{ setupData.secret }}
              </code>
              <button
                type="button"
                class="text-brand-blue-dark hover:text-brand-blue text-xs font-medium shrink-0"
                @click="copySecret"
              >
                {{ copied ? 'Copiado' : 'Copiar' }}
              </button>
            </div>
          </div>

          <button
            type="button"
            class="btn-primary w-full"
            @click="step = 'verify'"
          >
            Continuar
          </button>
        </div>

        <div v-else-if="loadError" class="space-y-3">
          <div class="rounded-md bg-red-50 p-4 text-sm text-red-700">{{ loadError }}</div>
          <button type="button" class="btn-primary w-full" @click="handleCancel">Volver</button>
        </div>
      </div>

      <!-- Paso 2: Verificar código -->
      <div v-if="step === 'verify'" class="space-y-6">
        <p class="text-sm text-gray-600 text-center">
          Ingresa el código de 6 dígitos que muestra tu aplicación para confirmar la activación
        </p>

        <form @submit.prevent="handleActivate" class="space-y-4">
          <div>
            <label for="verify-code" class="block text-sm font-medium text-gray-700 mb-1">
              Código de verificación
            </label>
            <input
              id="verify-code"
              v-model="verifyCode"
              type="text"
              inputmode="numeric"
              autocomplete="one-time-code"
              maxlength="6"
              placeholder="000000"
              class="appearance-none block w-full px-3 py-3 border placeholder-gray-400 text-gray-900 rounded-md focus:outline-none focus:ring-brand-blue-dark focus:border-brand-blue-dark text-center text-2xl tracking-widest font-mono"
              :class="activateError ? 'border-red-500' : 'border-gray-300'"
              @input="verifyCode = verifyCode.replace(/\D/g, '')"
            />
          </div>

          <div v-if="activateError" class="rounded-md bg-red-50 p-4">
            <p class="text-sm text-red-700">{{ activateError }}</p>
          </div>

          <button
            type="submit"
            :disabled="isActivating || verifyCode.length !== 6"
            class="btn-primary w-full disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ isActivating ? 'Activando...' : 'Activar TOTP' }}
          </button>

          <button
            type="button"
            class="w-full text-sm font-medium text-gray-500 hover:text-gray-700"
            @click="step = 'scan'"
          >
            Volver al QR
          </button>
        </form>
      </div>

      <!-- Paso 3: Éxito -->
      <div v-if="step === 'success'" class="space-y-6 text-center">
        <div class="flex justify-center">
          <div class="rounded-full bg-green-100 p-4">
            <svg class="h-12 w-12 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </div>
        </div>
        <h3 class="text-xl font-semibold text-gray-900">TOTP activado correctamente</h3>
        <p class="text-sm text-gray-600">
          La verificación en dos pasos está habilitada en tu cuenta. A partir de ahora necesitarás tu aplicación autenticadora para iniciar sesión.
        </p>
        <button type="button" class="btn-primary w-full" @click="handleDone">
          Volver al perfil
        </button>
      </div>

      <!-- Cancelar -->
      <div v-if="step !== 'success'" class="text-center">
        <button
          type="button"
          class="text-sm font-medium text-gray-500 hover:text-gray-700"
          @click="handleCancel"
        >
          Cancelar
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, getCurrentInstance } from 'vue'
import { useRouter } from 'vue-router'
import QrcodeVue from 'qrcode.vue'
import authService from '@/services/auth/authService'

const emit = defineEmits(['done', 'cancel'])
const router = useRouter()

// Si el componente no tiene padre que escuche los emits (usado como página standalone),
// navegar al perfil como fallback
const instance = getCurrentInstance()
const hasParentListeners = () => !!instance?.vnode?.props?.onDone

const handleDone = () => {
  emit('done')
  if (!hasParentListeners()) router.replace('/profile')
}

const handleCancel = () => {
  emit('cancel')
  if (!hasParentListeners()) router.replace('/profile')
}

const step = ref('scan')
const isLoadingSetup = ref(false)
const setupData = ref(null)
const loadError = ref('')
const copied = ref(false)

const verifyCode = ref('')
const isActivating = ref(false)
const activateError = ref('')

onMounted(async () => {
  isLoadingSetup.value = true
  try {
    setupData.value = await authService.getTOTPSetup()
  } catch (error) {
    loadError.value = error.message || 'Error al cargar la configuración TOTP'
  } finally {
    isLoadingSetup.value = false
  }
})

const copySecret = async () => {
  try {
    await navigator.clipboard.writeText(setupData.value.secret)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  } catch {
    // fallback silencioso
  }
}

const handleActivate = async () => {
  if (verifyCode.value.length !== 6) return

  isActivating.value = true
  activateError.value = ''

  try {
    await authService.activateTOTP(setupData.value.secret, verifyCode.value)
    step.value = 'success'
  } catch (error) {
    activateError.value = error.message || 'Código inválido. Intenta de nuevo.'
    verifyCode.value = ''
  } finally {
    isActivating.value = false
  }
}
</script>
