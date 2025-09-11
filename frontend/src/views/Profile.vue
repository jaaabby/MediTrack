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

                    <!-- Botón de editar -->
                    <div class="flex-shrink-0">
                        <button @click="editMode = !editMode"
                            class="inline-flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
                            <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                            </svg>
                            {{ editMode ? 'Cancelar' : 'Editar' }}
                        </button>
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

        <!-- Acciones rápidas -->
        <div class="mt-6 bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">Acciones Rápidas</h3>
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                <!-- Recepción de Insumos (Solo para Pabellones) -->
                <router-link to="/qr-reception" 
                    class="inline-flex items-center px-4 py-3 border border-transparent text-sm font-medium rounded-md text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 transition-colors">
                    <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    Recepcionar Insumos
                    <span class="ml-2 text-xs bg-green-500 px-2 py-1 rounded-full">Pabellón</span>
                </router-link>

                <!-- Consumir Insumos -->
                <router-link to="/qr-consumer" 
                    class="inline-flex items-center px-4 py-3 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors">
                    <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                    </svg>
                    Consumir Insumos
                    <span class="ml-2 text-xs bg-blue-500 px-2 py-1 rounded-full">Pabellón</span>
                </router-link>

                <!-- Escáner QR -->
                <router-link to="/qr" 
                    class="inline-flex items-center px-4 py-3 border border-transparent text-sm font-medium rounded-md text-white bg-purple-600 hover:bg-purple-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-purple-500 transition-colors">
                    <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h2M4 4h5l2 3h3l2-3h5v5M4 4v5m16-5v5" />
                    </svg>
                    Escáner QR
                </router-link>
            </div>
        </div>

        <!-- Botones de acción cuando está en modo edición -->
        <div v-if="editMode" class="mt-6 flex justify-end space-x-3">
            <button @click="cancelEdit"
                class="inline-flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
                Cancelar
            </button>
            <button @click="saveProfile" :disabled="loading"
                class="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50">
                <svg v-if="loading" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor"
                        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z">
                    </path>
                </svg>
                {{ loading ? 'Guardando...' : 'Guardar cambios' }}
            </button>
        </div>

        <!-- Mensaje de éxito/error -->
        <div v-if="message" class="mt-4">
            <div :class="messageType === 'success' ? 'bg-green-50 border-green-200 text-green-800' : 'bg-red-50 border-red-200 text-red-800'"
                class="border rounded-md p-4">
                {{ message }}
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { userService } from '@/services/userService'

const authStore = useAuthStore()
const editMode = ref(false)
const loading = ref(false)
const message = ref('')
const messageType = ref('success')

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
        showMessage('Error al cargar la información del perfil', 'error')
    } finally {
        loading.value = false
    }
}

// Guardar cambios del perfil
const saveProfile = async () => {
    try {
        loading.value = true

        const response = await userService.updateProfile(editForm)

        if (response.success) {
            // Actualizar store de autenticación
            authStore.updateUser({
                ...authStore.user,
                name: editForm.name,
                email: editForm.email
            })

            editMode.value = false
            showMessage('Perfil actualizado correctamente', 'success')
        } else {
            showMessage(response.error || 'Error al actualizar el perfil', 'error')
        }

    } catch (error) {
        console.error('Error al actualizar el perfil:', error)
        showMessage('Error al actualizar el perfil', 'error')
    } finally {
        loading.value = false
    }
}

// Cancelar edición
const cancelEdit = () => {
    editForm.name = authStore.getUserName || ''
    editForm.email = authStore.getUserEmail || ''
    editMode.value = false
    message.value = ''
}

// Mostrar mensaje
const showMessage = (text, type = 'success') => {
    message.value = text
    messageType.value = type
    setTimeout(() => {
        message.value = ''
    }, 5000)
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
