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

        <!-- Sección de seguridad -->
        <div class="bg-white rounded-lg shadow-sm border border-gray-200 mt-6">
            <div class="px-6 py-4 border-b border-gray-200">
                <h2 class="text-lg font-semibold text-gray-900">Seguridad</h2>
            </div>
            <div class="px-6 py-4">
                <div class="flex items-center justify-between">
                    <div>
                        <p class="text-sm font-medium text-gray-900">Cerrar sesión en todos los dispositivos</p>
                        <p class="text-sm text-gray-500 mt-1">Invalida todos los tokens activos. Deberás iniciar sesión nuevamente en cada dispositivo.</p>
                    </div>
                    <button @click="showLogoutAllConfirm = true"
                        class="inline-flex items-center px-4 py-2 border border-red-300 rounded-md shadow-sm text-sm font-medium text-red-700 bg-white hover:bg-red-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500">
                        <svg class="mr-2 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                        </svg>
                        Cerrar todas las sesiones
                    </button>
                </div>
            </div>
        </div>

        <!-- Modal de confirmación -->
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

const authStore = useAuthStore()
const router = useRouter()
const { success: showSuccess, error: showError } = useNotification()
const editMode = ref(false)
const loading = ref(false)
const showLogoutAllConfirm = ref(false)
const logoutAllLoading = ref(false)

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

// Cargar datos al montar el componente
onMounted(() => {
    loadProfile()
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
