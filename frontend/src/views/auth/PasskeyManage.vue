<template>
  <div class="max-w-2xl mx-auto py-8 px-4">
    <div class="bg-white shadow rounded-lg p-6 space-y-6">
      <!-- Encabezado -->
      <div>
        <h2 class="text-xl font-semibold text-gray-900">Llaves de acceso (Passkeys)</h2>
        <p class="mt-1 text-sm text-gray-500">
          Las passkeys permiten iniciar sesión con tu huella digital, Face ID o PIN del dispositivo,
          sin necesidad de contraseña.
        </p>
      </div>

      <!-- Error global -->
      <div v-if="error" class="rounded-md bg-red-50 p-4">
        <p class="text-sm text-red-700">{{ error }}</p>
      </div>

      <!-- Lista de passkeys -->
      <div>
        <h3 class="text-sm font-medium text-gray-700 mb-3">Tus passkeys registradas</h3>

        <div v-if="loading" class="flex justify-center py-6">
          <svg class="animate-spin h-6 w-6 text-gray-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
        </div>

        <ul v-else-if="passkeys.length" class="divide-y divide-gray-200 border border-gray-200 rounded-md">
          <li v-for="pk in passkeys" :key="pk.id" class="flex items-center justify-between px-4 py-3">
            <div class="flex items-center gap-3">
              <!-- Ícono de llave -->
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-400" viewBox="0 0 24 24" fill="none"
                stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/>
              </svg>
              <div>
                <p class="text-sm font-medium text-gray-900">{{ pk.name || 'Passkey' }}</p>
                <p class="text-xs text-gray-400">Registrada el {{ formatDate(pk.created_at) }}</p>
              </div>
            </div>
            <button
              @click="confirmDelete(pk)"
              :disabled="deleting === pk.id"
              class="text-red-500 hover:text-red-700 text-sm font-medium disabled:opacity-50"
            >
              {{ deleting === pk.id ? 'Eliminando...' : 'Eliminar' }}
            </button>
          </li>
        </ul>

        <p v-else class="text-sm text-gray-500 italic py-4 text-center">
          No tienes passkeys registradas.
        </p>
      </div>

      <!-- Agregar nueva passkey -->
      <div class="border-t pt-4 space-y-3">
        <h3 class="text-sm font-medium text-gray-700">Agregar nueva passkey</h3>
        <div class="flex gap-2">
          <input
            v-model="newName"
            type="text"
            placeholder="Nombre del dispositivo (ej: Mi iPhone)"
            class="flex-1 border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-brand-blue-dark"
            :disabled="registering"
            @keyup.enter="handleRegister"
          />
          <button
            @click="handleRegister"
            :disabled="registering || !passkeySupported"
            class="btn-primary px-4 py-2 text-sm disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ registering ? 'Registrando...' : 'Agregar' }}
          </button>
        </div>
        <p v-if="!passkeySupported" class="text-xs text-yellow-600">
          Tu navegador no soporta Passkeys.
        </p>
        <p v-if="successMessage" class="text-xs text-green-600">{{ successMessage }}</p>
      </div>
    </div>

    <!-- Modal de confirmación -->
    <div v-if="toDelete" class="fixed inset-0 bg-black bg-opacity-40 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg shadow-xl p-6 max-w-sm w-full mx-4 space-y-4">
        <h3 class="text-base font-semibold text-gray-900">¿Eliminar passkey?</h3>
        <p class="text-sm text-gray-600">
          Eliminarás la passkey <strong>"{{ toDelete.name }}"</strong>. Ya no podrás usarla para iniciar sesión.
        </p>
        <div class="flex justify-end gap-3">
          <button @click="toDelete = null" class="text-sm text-gray-600 hover:text-gray-800">Cancelar</button>
          <button @click="handleDelete" class="btn-danger text-sm px-4 py-2">Eliminar</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const passkeys = ref([])
const loading = ref(false)
const error = ref('')
const deleting = ref(null)
const registering = ref(false)
const newName = ref('')
const successMessage = ref('')
const toDelete = ref(null)
const passkeySupported = ref(false)

onMounted(async () => {
  passkeySupported.value = !!window.PublicKeyCredential
  await loadPasskeys()
})

async function loadPasskeys() {
  loading.value = true
  error.value = ''
  try {
    passkeys.value = await authStore.listPasskeys()
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

async function handleRegister() {
  if (!passkeySupported.value) return
  registering.value = true
  error.value = ''
  successMessage.value = ''
  try {
    await authStore.registerPasskey(newName.value || 'Mi Passkey')
    newName.value = ''
    successMessage.value = 'Passkey registrada exitosamente.'
    await loadPasskeys()
  } catch (e) {
    error.value = e.message
  } finally {
    registering.value = false
  }
}

function confirmDelete(pk) {
  toDelete.value = pk
}

async function handleDelete() {
  if (!toDelete.value) return
  const pk = toDelete.value
  toDelete.value = null
  deleting.value = pk.id
  error.value = ''
  try {
    await authStore.deletePasskey(pk.id)
    await loadPasskeys()
  } catch (e) {
    error.value = e.message
  } finally {
    deleting.value = null
  }
}

function formatDate(unixSeconds) {
  if (!unixSeconds) return '—'
  return new Date(unixSeconds * 1000).toLocaleDateString('es-CL', {
    year: 'numeric', month: 'short', day: 'numeric'
  })
}
</script>
