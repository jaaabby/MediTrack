<template>
    <div class="max-w-4xl mx-auto">
        <!-- Header del perfil -->
        <div class="bg-white rounded-lg shadow-sm border border-gray-200 mb-6">
            <div class="px-6 py-8">
                <div class="flex items-center space-x-6">
                    <!-- Avatar del usuario -->
                    <div class="flex-shrink-0">
                        <div class="h-24 w-24 bg-blue-100 rounded-full flex items-center justify-center">
                            <svg class="h-12 w-12 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                            </svg>
                        </div>
                    </div>

                    <!-- Información básica -->
                    <div class="flex-1">
                        <h1 class="text-3xl font-bold text-gray-900">{{ authStore.getUserName || 'Usuario' }}</h1>
                        <p class="text-lg text-gray-600 mt-1">{{ authStore.getUserEmail || 'email@ejemplo.com' }}</p>
                        <div class="flex items-center mt-2">
                            <span
                                class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-blue-100 text-blue-800">
                                {{ authStore.getUserRole || 'Usuario' }}
                            </span>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Modal de configuración TOTP -->
        <div v-if="showTOTPSetup" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 p-4">
            <div class="bg-white rounded-xl shadow-2xl w-full max-w-md max-h-[90vh] overflow-y-auto">
                <TOTPSetup @done="onTOTPDone" @cancel="showTOTPSetup = false" />
            </div>
        </div>

        <!-- Información detallada -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <!-- Información personal -->
            <div class="bg-white rounded-lg shadow-sm border border-gray-200">
                <div class="px-6 py-4 border-b border-gray-200">
                    <h2 class="text-lg font-semibold text-gray-900">Información Personal</h2>
                </div>
                <div class="px-6 py-4 space-y-4">
                    <div>
                        <label class="block text-sm font-medium text-gray-700">Nombre completo</label>
                        <div v-if="!editMode" class="mt-1 text-sm text-gray-900">{{ authStore.getUserName || 'No especificado' }}</div>
                        <input v-else v-model="editForm.name" type="text"
                            class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm">
                    </div>

                    <div>
                        <label class="block text-sm font-medium text-gray-700">RUT</label>
                        <div class="mt-1 text-sm text-gray-900">{{ authStore.getUserRut || 'No especificado' }}</div>
                    </div>

                    <div>
                        <label class="block text-sm font-medium text-gray-700">Email</label>
                        <div v-if="!editMode" class="mt-1 text-sm text-gray-900">{{ authStore.getUserEmail || 'No especificado' }}</div>
                        <input v-else v-model="editForm.email" type="email"
                            class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm">
                    </div>

                    <div>
                        <label class="block text-sm font-medium text-gray-700">Rol</label>
                        <div class="mt-1 text-sm text-gray-900">
                            <span
                                class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                                {{ authStore.getUserRole || 'Usuario' }}
                            </span>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Información del centro médico -->
            <div class="bg-white rounded-lg shadow-sm border border-gray-200">
                <div class="px-6 py-4 border-b border-gray-200">
                    <h2 class="text-lg font-semibold text-gray-900">Centro Médico</h2>
                </div>
                <div class="px-6 py-4 space-y-4">
                    <div>
                        <label class="block text-sm font-medium text-gray-700">Centro médico</label>
                        <div class="mt-1 text-sm text-gray-900">{{ authStore.getUserMedicalCenterName || 'No asignado' }}
                        </div>
                    </div>

                    <div>
                        <label class="block text-sm font-medium text-gray-700">Fecha de registro</label>
                        <div class="mt-1 text-sm text-gray-900">{{ formatDate(authStore.getUserCreatedAt) }}</div>
                    </div>

                    <div>
                        <label class="block text-sm font-medium text-gray-700">Última actualización</label>
                        <div class="mt-1 text-sm text-gray-900">{{ formatDate(authStore.getUserUpdatedAt) }}</div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Sección de seguridad unificada -->
        <div class="bg-white rounded-lg shadow-sm border border-gray-200 mt-6">
            <div class="px-6 py-4 border-b border-gray-200">
                <h2 class="text-lg font-semibold text-gray-900">Seguridad</h2>
                <p class="mt-1 text-sm text-gray-500">Gestiona el acceso y la protección de tu cuenta.</p>
            </div>

            <!-- Sub-ítem: Cerrar sesión en todos los dispositivos -->
            <div class="px-6 py-5">
                <div class="flex items-start justify-between gap-4">
                    <div class="flex items-start gap-3">
                        <div class="mt-0.5 flex-shrink-0 h-9 w-9 rounded-full bg-red-50 flex items-center justify-center">
                            <svg class="h-5 w-5 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                            </svg>
                        </div>
                        <div>
                            <p class="text-sm font-medium text-gray-900">Cerrar sesión en todos los dispositivos</p>
                            <p class="text-sm text-gray-500 mt-0.5">Invalida todos los tokens activos. Deberás iniciar sesión nuevamente en cada dispositivo.</p>
                        </div>
                    </div>
                    <button @click="showLogoutAllConfirm = true"
                        class="shrink-0 inline-flex items-center px-4 py-2 border border-red-300 rounded-md shadow-sm text-sm font-medium text-red-700 bg-white hover:bg-red-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500">
                        Cerrar todas las sesiones
                    </button>
                </div>
            </div>

            <div class="border-t border-gray-100" />

            <!-- Sub-ítem: Verificación en dos pasos -->
            <div class="px-6 py-5">
                <div class="flex items-start justify-between gap-4">
                    <div class="flex items-start gap-3">
                        <div class="mt-0.5 flex-shrink-0 h-9 w-9 rounded-full bg-blue-50 flex items-center justify-center">
                            <svg class="h-5 w-5 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                            </svg>
                        </div>
                        <div>
                            <p class="text-sm font-medium text-gray-900">Verificación en dos pasos</p>
                            <p class="text-sm text-gray-500 mt-0.5">Añade una capa extra de seguridad con una aplicación autenticadora.</p>
                            <span
                                v-if="authStore.user?.totp_enabled"
                                class="inline-flex items-center mt-2 px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800"
                            >
                                Activo
                            </span>
                            <span
                                v-else
                                class="inline-flex items-center mt-2 px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-600"
                            >
                                Inactivo
                            </span>
                        </div>
                    </div>
                    <div class="shrink-0 ml-4">
                        <button
                            v-if="!authStore.user?.totp_enabled"
                            type="button"
                            class="btn-primary text-sm"
                            @click="showTOTPSetup = true"
                        >
                            Activar
                        </button>
                        <button
                            v-else
                            type="button"
                            class="inline-flex items-center px-4 py-2 border border-red-300 rounded-md shadow-sm text-sm font-medium text-red-700 bg-white hover:bg-red-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
                            @click="showDisableTOTPModal = true"
                        >
                            Desactivar
                        </button>
                    </div>
                </div>
            </div>

            <div class="border-t border-gray-100" />

            <!-- Sub-ítem: Acceso biométrico -->
            <div class="px-6 py-5">
                <div class="flex items-start gap-3 mb-4">
                    <div class="mt-0.5 flex-shrink-0 h-9 w-9 rounded-full bg-green-50 flex items-center justify-center">
                        <svg class="h-5 w-5 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
                        </svg>
                    </div>
                    <div>
                        <p class="text-sm font-medium text-gray-900">Acceso biométrico</p>
                        <p class="text-sm text-gray-500 mt-0.5">Ingresa al sistema usando tu huella digital, reconocimiento facial o PIN del dispositivo, sin necesidad de contraseña.</p>
                    </div>
                </div>

                <div v-if="passkeyError" class="rounded-md bg-red-50 p-3">
                    <p class="text-sm text-red-700">{{ passkeyError }}</p>
                </div>
                <div v-if="passkeySuccess" class="rounded-md bg-green-50 p-3">
                    <p class="text-sm text-green-700">{{ passkeySuccess }}</p>
                </div>

                <!-- Lista de dispositivos registrados -->
                <div v-if="passkeyLoading" class="flex justify-center py-4">
                    <svg class="animate-spin h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                </div>

                <div v-else-if="passkeys.length">
                    <p class="text-xs font-medium text-gray-500 uppercase tracking-wide mb-2">Dispositivos registrados</p>
                    <ul class="divide-y divide-gray-200 border border-gray-200 rounded-md">
                        <li v-for="pk in passkeys" :key="pk.id" class="flex items-center justify-between px-4 py-3">
                            <div class="flex items-center gap-3">
                                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-400 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 11c0-1.1.9-2 2-2s2 .9 2 2-.9 2-2 2-2-.9-2-2zm0 0V7m0 4v4m-6 4h12a2 2 0 002-2V5a2 2 0 00-2-2H6a2 2 0 00-2 2v14a2 2 0 002 2z"/>
                                </svg>
                                <div>
                                    <p class="text-sm font-medium text-gray-900">{{ pk.name || 'Mi dispositivo' }}</p>
                                    <p class="text-xs text-gray-400">Configurado el {{ formatDate(pk.created_at) }}</p>
                                </div>
                            </div>
                            <button
                                @click="passkeyToDelete = pk"
                                :disabled="deletingPasskey === pk.id"
                                class="text-sm text-red-500 hover:text-red-700 font-medium disabled:opacity-50"
                            >
                                {{ deletingPasskey === pk.id ? 'Eliminando...' : 'Eliminar' }}
                            </button>
                        </li>
                    </ul>
                </div>

                <p v-else-if="!passkeyLoading" class="text-sm text-gray-400 italic">
                    No tienes ningún dispositivo configurado para acceso biométrico.
                </p>

                <!-- Agregar dispositivo -->
                <div v-if="passkeySupported" class="pt-1">
                    <p class="text-xs font-medium text-gray-500 uppercase tracking-wide mb-2">Agregar dispositivo</p>
                    <div class="flex gap-2">
                        <input
                            v-model="newPasskeyName"
                            type="text"
                            placeholder="Nombre del dispositivo (ej: Mi iPhone, PC Trabajo)"
                            class="flex-1 border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-brand-blue-dark"
                            :disabled="registeringPasskey"
                            @keyup.enter="handleRegisterPasskey"
                        />
                        <button
                            @click="handleRegisterPasskey"
                            :disabled="registeringPasskey"
                            class="btn-primary text-sm px-4 disabled:opacity-50"
                        >
                            {{ registeringPasskey ? 'Configurando...' : 'Configurar' }}
                        </button>
                    </div>
                    <p class="mt-1.5 text-xs text-gray-400">
                        Al hacer clic en "Configurar", tu dispositivo te pedirá confirmar con tu huella, rostro o PIN.
                    </p>
                </div>
                <p v-else class="text-xs text-yellow-600">
                    Tu navegador no soporta autenticación biométrica.
                </p>
            </div>
        </div>

        <!-- Modal de confirmación: Cerrar todas las sesiones -->
        <div v-if="showLogoutAllConfirm" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-40">
            <div class="bg-white rounded-lg shadow-xl max-w-md w-full mx-4 p-6">
                <div class="flex items-center mb-4">
                    <div class="flex-shrink-0 h-10 w-10 rounded-full bg-red-100 flex items-center justify-center mr-3">
                        <svg class="h-5 w-5 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                        </svg>
                    </div>
                    <h3 class="text-lg font-semibold text-gray-900">Cerrar todas las sesiones</h3>
                </div>
                <p class="text-sm text-gray-600 mb-6">
                    ¿Estás seguro? Esto cerrará tu sesión en todos los dispositivos donde hayas iniciado sesión,
                    incluyendo este. Deberás volver a ingresar tus credenciales.
                </p>
                <div class="flex justify-end space-x-3">
                    <button @click="showLogoutAllConfirm = false" :disabled="logoutAllLoading"
                        class="inline-flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50">
                        Cancelar
                    </button>
                    <button @click="handleLogoutAllDevices" :disabled="logoutAllLoading"
                        class="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 disabled:opacity-50">
                        <svg v-if="logoutAllLoading" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                        </svg>
                        {{ logoutAllLoading ? 'Cerrando sesiones...' : 'Cerrar todas las sesiones' }}
                    </button>
                </div>
            </div>
        </div>

        <!-- Modal confirmación eliminar dispositivo -->
        <div v-if="passkeyToDelete" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-40">
            <div class="bg-white rounded-lg shadow-xl max-w-sm w-full mx-4 p-6 space-y-4">
                <h3 class="text-lg font-semibold text-gray-900">¿Eliminar dispositivo?</h3>
                <p class="text-sm text-gray-600">
                    Eliminarás el acceso biométrico de <strong>"{{ passkeyToDelete.name }}"</strong>.
                    Ya no podrás usar ese dispositivo para ingresar al sistema.
                </p>
                <div class="flex justify-end gap-3 pt-2">
                    <button @click="passkeyToDelete = null" class="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50">
                        Cancelar
                    </button>
                    <button @click="handleDeletePasskey" class="px-4 py-2 border border-transparent rounded-md text-sm font-medium text-white bg-red-600 hover:bg-red-700">
                        Eliminar
                    </button>
                </div>
            </div>
        </div>

        <!-- Modal para desactivar TOTP -->
        <div v-if="showDisableTOTPModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-40">
            <div class="bg-white rounded-lg shadow-xl max-w-sm w-full mx-4 p-6 space-y-4">
                <h3 class="text-lg font-semibold text-gray-900">Desactivar verificación en dos pasos</h3>
                <p class="text-sm text-gray-600">Confirma tu contraseña para desactivar la verificación en dos pasos.</p>
                <input
                    v-model="disableTOTPPassword"
                    type="password"
                    placeholder="Tu contraseña actual"
                    class="block w-full border-gray-300 rounded-md shadow-sm focus:ring-red-500 focus:border-red-500 sm:text-sm"
                />
                <p v-if="disableTOTPError" class="text-sm text-red-600">{{ disableTOTPError }}</p>
                <div class="flex justify-end gap-3 pt-2">
                    <button
                        type="button"
                        class="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50"
                        @click="closeDisableTOTPModal"
                    >
                        Cancelar
                    </button>
                    <button
                        type="button"
                        :disabled="disablingTOTP || !disableTOTPPassword"
                        class="px-4 py-2 border border-transparent rounded-md text-sm font-medium text-white bg-red-600 hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed"
                        @click="handleDisableTOTP"
                    >
                        {{ disablingTOTP ? 'Desactivando...' : 'Desactivar' }}
                    </button>
                </div>
            </div>
        </div>

        <!-- Botones de acción cuando está en modo edición -->
        <div v-if="editMode" class="mt-6 flex justify-end space-x-3">
            <button @click="cancelEdit"
                class="inline-flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
                Cancelar
            </button>
            <button @click="saveProfile" :disabled="loading"
                class="btn-primary disabled:opacity-50">
                <svg v-if="loading" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor"
                        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z">
                    </path>
                </svg>
                {{ loading ? 'Guardando...' : 'Guardar cambios' }}
            </button>
        </div>

    </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { userService } from '@/services/common/userService'
import { useNotification } from '@/composables/useNotification'
import authService from '@/services/auth/authService'
import TOTPSetup from '@/views/auth/TOTPSetup.vue'
import { registerPasskey, listPasskeys, deletePasskey } from '@/services/auth/passkeyService'

const authStore = useAuthStore()
const router = useRouter()
const { success: showSuccess, error: showError } = useNotification()
const editMode = ref(false)
const loading = ref(false)
const showLogoutAllConfirm = ref(false)
const logoutAllLoading = ref(false)

// ── Passkeys ────────────────────────────────────────────────────────────
const passkeys = ref([])
const passkeyLoading = ref(false)
const passkeyError = ref('')
const passkeySuccess = ref('')
const registeringPasskey = ref(false)
const deletingPasskey = ref(null)
const passkeyToDelete = ref(null)
const newPasskeyName = ref('')
const passkeySupported = ref(false)

const loadPasskeys = async () => {
    passkeyLoading.value = true
    passkeyError.value = ''
    try {
        passkeys.value = await listPasskeys(authStore.token)
    } catch (e) {
        passkeyError.value = e.message
    } finally {
        passkeyLoading.value = false
    }
}

const handleRegisterPasskey = async () => {
    registeringPasskey.value = true
    passkeyError.value = ''
    passkeySuccess.value = ''
    try {
        await registerPasskey(authStore.token, newPasskeyName.value || 'Mi Passkey')
        newPasskeyName.value = ''
        passkeySuccess.value = 'Passkey registrada exitosamente.'
        await loadPasskeys()
    } catch (e) {
        passkeyError.value = e.message
    } finally {
        registeringPasskey.value = false
    }
}

const handleDeletePasskey = async () => {
    if (!passkeyToDelete.value) return
    const pk = passkeyToDelete.value
    passkeyToDelete.value = null
    deletingPasskey.value = pk.id
    passkeyError.value = ''
    try {
        await deletePasskey(authStore.token, pk.id)
        await loadPasskeys()
    } catch (e) {
        passkeyError.value = e.message
    } finally {
        deletingPasskey.value = null
    }
}

// Estado TOTP
const showTOTPSetup = ref(false)
const showDisableTOTPModal = ref(false)
const disableTOTPPassword = ref('')
const disableTOTPError = ref('')
const disablingTOTP = ref(false)

// Formulario de edición
const editForm = reactive({
    name: '',
    email: ''
})

// Cargar información del perfil
const loadProfile = async () => {
    try {
        loading.value = true
        
        // Inicializar formulario de edición con datos del authStore
        editForm.name = authStore.getUserName || ''
        editForm.email = authStore.getUserEmail || ''

    } catch (error) {
        console.error('Error al cargar el perfil:', error)
        showError('Error al cargar la información del perfil')
    } finally {
        loading.value = false
    }
}

// Guardar cambios del perfil
const saveProfile = async () => {
    try {
        loading.value = true

        // Validar email
        const emailRegex = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
        if (!editForm.email || !editForm.email.trim()) {
            showError('El email es requerido')
            loading.value = false
            return
        }
        
        if (!emailRegex.test(editForm.email.trim())) {
            showError('El formato del email no es válido')
            loading.value = false
            return
        }

        const response = await userService.updateProfile(editForm)

        if (response.success) {
            // Actualizar store de autenticación
            authStore.updateUser({
                ...authStore.user,
                name: editForm.name,
                email: editForm.email
            })

            editMode.value = false
            showSuccess('Perfil actualizado correctamente')
        } else {
            showError(response.error || 'Error al actualizar el perfil')
        }

    } catch (error) {
        console.error('Error al actualizar el perfil:', error)
        showError('Error al actualizar el perfil')
    } finally {
        loading.value = false
    }
}

// Cancelar edición
const cancelEdit = () => {
    editForm.name = authStore.getUserName || ''
    editForm.email = authStore.getUserEmail || ''
    editMode.value = false
}

// Formatear fecha
const formatDate = (dateValue) => {
    if (!dateValue) return 'No disponible'

    let date

    // Manejar diferentes tipos de entrada
    if (typeof dateValue === 'string') {
        // Si es una cadena ISO o similar
        date = new Date(dateValue)
    } else if (typeof dateValue === 'number') {
        // Si es un timestamp de Unix
        date = new Date(dateValue < 1e12 ? dateValue * 1000 : dateValue)
    } else if (dateValue instanceof Date) {
        // Si ya es un objeto Date
        date = dateValue
    } else {
        // Intentar convertir cualquier otro tipo
        date = new Date(dateValue)
    }

    // Verificar si la fecha es válida
    if (isNaN(date.getTime())) {
        return 'Fecha inválida'
    }

    return date.toLocaleDateString('es-ES', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    })
}

// Cerrar sesión en todos los dispositivos
const handleLogoutAllDevices = async () => {
    logoutAllLoading.value = true
    try {
        await authStore.logoutAllDevices()
        router.push('/login')
    } catch (error) {
        showError(error.message || 'Error al cerrar sesión en todos los dispositivos')
        showLogoutAllConfirm.value = false
    } finally {
        logoutAllLoading.value = false
    }
}

// Callbacks TOTP
const onTOTPDone = () => {
    showTOTPSetup.value = false
    authStore.updateUser({ ...authStore.user, totp_enabled: true })
    showSuccess('Verificación en dos pasos activada')
}

const closeDisableTOTPModal = () => {
    showDisableTOTPModal.value = false
    disableTOTPPassword.value = ''
    disableTOTPError.value = ''
}

const handleDisableTOTP = async () => {
    if (!disableTOTPPassword.value) return
    disablingTOTP.value = true
    disableTOTPError.value = ''
    try {
        await authService.disableTOTP(disableTOTPPassword.value)
        authStore.updateUser({ ...authStore.user, totp_enabled: false })
        closeDisableTOTPModal()
        showSuccess('Verificación en dos pasos desactivada')
    } catch (error) {
        disableTOTPError.value = error.message || 'Error al desactivar la verificación en dos pasos'
    } finally {
        disablingTOTP.value = false
    }
}

// Cargar datos al montar el componente
onMounted(() => {
    loadProfile()
    passkeySupported.value = !!window.PublicKeyCredential
    loadPasskeys()
})
</script>

<style scoped>
/* Estilos específicos del componente */
.animate-spin {
    animation: spin 1s linear infinite;
}

@keyframes spin {
    from {
        transform: rotate(0deg);
    }

    to {
        transform: rotate(360deg);
    }
}
</style>
