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
          Cambio de Contraseña Obligatorio
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600">
          Por seguridad, debes establecer una nueva contraseña antes de continuar
        </p>
      </div>

      <!-- Alerta informativa -->
      <div class="rounded-md bg-blue-50 p-4">
        <div class="flex">
          <div class="flex-shrink-0">
            <svg class="h-5 w-5 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div class="ml-3 flex-1">
            <p class="text-sm text-blue-700">
              Has ingresado con tu contraseña temporal. Debes establecer una nueva contraseña para continuar usando el sistema.
            </p>
          </div>
        </div>
      </div>

      <!-- Formulario de cambio de contraseña -->
      <form class="mt-8 space-y-6" @submit.prevent="handlePasswordChange">
        <div class="space-y-4">
          <!-- Nueva contraseña -->
          <div>
            <label for="new-password" class="block text-sm font-medium text-gray-700 mb-2">
              Nueva Contraseña
            </label>
            <div class="relative">
              <input
                id="new-password"
                v-model="passwordForm.newPassword"
                :type="showNewPassword ? 'text' : 'password'"
                required
                class="appearance-none relative block w-full px-3 py-2 pr-10 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-brand-blue-dark focus:border-brand-blue-dark focus:z-10 sm:text-sm"
                :class="{ 'border-red-500': errors.newPassword }"
                placeholder="Mínimo 6 caracteres"
              />
              <button type="button" tabindex="-1" @click="showNewPassword = !showNewPassword" class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600">
                <svg v-if="!showNewPassword" class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" /><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" /></svg>
                <svg v-else class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" /></svg>
              </button>
            </div>
            <p v-if="errors.newPassword" class="mt-1 text-sm text-red-600">{{ errors.newPassword }}</p>
            
            <!-- Barra de fortaleza -->
            <div v-if="passwordForm.newPassword" class="mt-3">
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

          <!-- Confirmar nueva contraseña -->
          <div>
            <label for="confirm-password" class="block text-sm font-medium text-gray-700 mb-2">
              Confirmar Nueva Contraseña
            </label>
            <div class="relative">
              <input
                id="confirm-password"
                v-model="passwordForm.confirmPassword"
                :type="showConfirmPassword ? 'text' : 'password'"
                required
                class="appearance-none relative block w-full px-3 py-2 pr-10 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-brand-blue-dark focus:border-brand-blue-dark focus:z-10 sm:text-sm"
                :class="{ 'border-red-500': errors.confirmPassword }"
                placeholder="Repite la nueva contraseña"
              />
              <button type="button" tabindex="-1" @click="showConfirmPassword = !showConfirmPassword" class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600">
                <svg v-if="!showConfirmPassword" class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" /><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" /></svg>
                <svg v-else class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" /></svg>
              </button>
            </div>
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

        <!-- Mensaje de error general -->
        <div v-if="errorMessage" class="rounded-md bg-red-50 p-4">
          <div class="flex">
            <div class="flex-shrink-0">
              <svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
              </svg>
            </div>
            <div class="ml-3">
              <h3 class="text-sm font-medium text-red-800">
                Error al cambiar contraseña
              </h3>
              <div class="mt-2 text-sm text-red-700">
                {{ errorMessage }}
              </div>
            </div>
          </div>
        </div>

        <!-- Botón de envío -->
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
            {{ isLoading ? 'Cambiando contraseña...' : 'Cambiar Contraseña' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import authService from '@/services/auth/authService'
import { useNotification } from '@/composables/useNotification'

const router = useRouter()
const authStore = useAuthStore()
const { success: showSuccess, error: showError } = useNotification()

// Mostrar/ocultar contraseñas
const showNewPassword = ref(false)
const showConfirmPassword = ref(false)

// Estado del formulario
const passwordForm = reactive({
  newPassword: '',
  confirmPassword: ''
})

// Calcular fortaleza de la contraseña
const passwordStrength = computed(() => {
  const password = passwordForm.newPassword
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

// Estado de la UI
const isLoading = ref(false)
const errorMessage = ref('')
const errors = reactive({
  newPassword: '',
  confirmPassword: ''
})

// Validación del formulario
const validateForm = () => {
  errors.newPassword = ''
  errors.confirmPassword = ''
  
  if (!passwordForm.newPassword) {
    errors.newPassword = 'La nueva contraseña es requerida'
    return false
  }
  
  if (passwordForm.newPassword.length < 6) {
    errors.newPassword = 'La contraseña debe tener al menos 6 caracteres'
    return false
  }
  
  if (!passwordForm.confirmPassword) {
    errors.confirmPassword = 'Debes confirmar la nueva contraseña'
    return false
  }
  
  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    errors.confirmPassword = 'Las contraseñas no coinciden'
    return false
  }
  
  return true
}

// Manejo del cambio de contraseña - CORREGIDO
const handlePasswordChange = async () => {
  if (!validateForm()) {
    return
  }
  
  isLoading.value = true
  errorMessage.value = ''
  
  try {
    await authService.firstTimePasswordChange(
      passwordForm.newPassword
    )
    
    // Actualizar el estado del usuario en el store
    if (authStore.user) {
      authStore.user.must_change_password = false
      localStorage.setItem('user_full', JSON.stringify(authStore.user))
    }
    
    showSuccess('Contraseña cambiada exitosamente')
    
    // Redirigir al home después de un breve delay
    setTimeout(() => {
      router.replace('/home')
    }, 1000)
    
  } catch (error) {
    console.error('Error al cambiar contraseña:', error)
    errorMessage.value = error.message || 'Error al cambiar la contraseña. Verifica tu contraseña temporal.'
  } finally {
    isLoading.value = false
  }
}

// Verificar al montar que el usuario realmente necesite cambiar contraseña
onMounted(() => {
  if (!authStore.isAuthenticated) {
    router.replace('/login')
    return
  }
  
  if (!authStore.user?.must_change_password) {
    // Si no necesita cambiar contraseña, redirigir a home
    router.replace('/home')
  }
})
</script>

<style scoped>
/* Estilos adicionales si es necesario */
</style>