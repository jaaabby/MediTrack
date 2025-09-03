import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'

// Crear instancia de axios con configuración base
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Interceptor para agregar token de autenticación
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('auth_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Interceptor para manejar respuestas y errores
apiClient.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      // Token expirado o inválido
      localStorage.removeItem('auth_token')
      localStorage.removeItem('user_data')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export const userService = {
  // Obtener información completa del perfil del usuario
  async getProfile(email) {
    try {
      const response = await apiClient.get(`/users/profile/${email}`)
      return {
        success: true,
        data: response.data
      }
    } catch (error) {
      console.error('Error al obtener perfil:', error)
      return {
        success: false,
        error: error.response?.data?.message || 'Error al obtener información del perfil'
      }
    }
  },

  // Actualizar información del perfil
  async updateProfile(profileData) {
    try {
      const response = await apiClient.put('/users/profile', profileData)
      return {
        success: true,
        data: response.data
      }
    } catch (error) {
      console.error('Error al actualizar perfil:', error)
      return {
        success: false,
        error: error.response?.data?.message || 'Error al actualizar el perfil'
      }
    }
  },

  // Cambiar contraseña
  async changePassword(passwordData) {
    try {
      const response = await apiClient.put('/users/change-password', passwordData)
      return {
        success: true,
        data: response.data
      }
    } catch (error) {
      console.error('Error al cambiar contraseña:', error)
      return {
        success: false,
        error: error.response?.data?.message || 'Error al cambiar la contraseña'
      }
    }
  },

  // Obtener estadísticas del usuario
  async getUserStats() {
    try {
      const response = await apiClient.get('/users/stats')
      return {
        success: true,
        data: response.data
      }
    } catch (error) {
      console.error('Error al obtener estadísticas:', error)
      return {
        success: false,
        error: error.response?.data?.message || 'Error al obtener estadísticas'
      }
    }
  },

  // Obtener historial de actividad del usuario
  async getUserActivity(limit = 10, offset = 0) {
    try {
      const response = await apiClient.get(`/users/activity?limit=${limit}&offset=${offset}`)
      return {
        success: true,
        data: response.data
      }
    } catch (error) {
      console.error('Error al obtener actividad:', error)
      return {
        success: false,
        error: error.response?.data?.message || 'Error al obtener historial de actividad'
      }
    }
  },

  // Subir foto de perfil
  async uploadProfilePicture(file) {
    try {
      const formData = new FormData()
      formData.append('profile_picture', file)

      const response = await apiClient.post('/users/profile-picture', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      })
      return {
        success: true,
        data: response.data
      }
    } catch (error) {
      console.error('Error al subir foto de perfil:', error)
      return {
        success: false,
        error: error.response?.data?.message || 'Error al subir la foto de perfil'
      }
    }
  },

  // Eliminar foto de perfil
  async deleteProfilePicture() {
    try {
      const response = await apiClient.delete('/users/profile-picture')
      return {
        success: true,
        data: response.data
      }
    } catch (error) {
      console.error('Error al eliminar foto de perfil:', error)
      return {
        success: false,
        error: error.response?.data?.message || 'Error al eliminar la foto de perfil'
      }
    }
  },

  // Obtener configuración de notificaciones del usuario
  async getNotificationSettings() {
    try {
      const response = await apiClient.get('/users/notification-settings')
      return {
        success: true,
        data: response.data
      }
    } catch (error) {
      console.error('Error al obtener configuración de notificaciones:', error)
      return {
        success: false,
        error: error.response?.data?.message || 'Error al obtener configuración de notificaciones'
      }
    }
  },

  // Actualizar configuración de notificaciones
  async updateNotificationSettings(settings) {
    try {
      const response = await apiClient.put('/users/notification-settings', settings)
      return {
        success: true,
        data: response.data
      }
    } catch (error) {
      console.error('Error al actualizar configuración de notificaciones:', error)
      return {
        success: false,
        error: error.response?.data?.message || 'Error al actualizar configuración de notificaciones'
      }
    }
  },

  // Obtener lista de usuarios (solo para administradores)
  async getUsers(filters = {}) {
    try {
      const queryParams = new URLSearchParams(filters).toString()
      const response = await apiClient.get(`/users?${queryParams}`)
      return {
        success: true,
        data: response.data
      }
    } catch (error) {
      console.error('Error al obtener usuarios:', error)
      return {
        success: false,
        error: error.response?.data?.message || 'Error al obtener lista de usuarios'
      }
    }
  },

  // Crear nuevo usuario (solo para administradores)
  async createUser(userData) {
    try {
      const response = await apiClient.post('/users', userData)
      return {
        success: true,
        data: response.data
      }
    } catch (error) {
      console.error('Error al crear usuario:', error)
      return {
        success: false,
        error: error.response?.data?.message || 'Error al crear usuario'
      }
    }
  },

  // Actualizar usuario (solo para administradores)
  async updateUser(userId, userData) {
    try {
      const response = await apiClient.put(`/users/${userId}`, userData)
      return {
        success: true,
        data: response.data
      }
    } catch (error) {
      console.error('Error al actualizar usuario:', error)
      return {
        success: false,
        error: error.response?.data?.message || 'Error al actualizar usuario'
      }
    }
  },

  // Eliminar usuario (solo para administradores)
  async deleteUser(userId) {
    try {
      const response = await apiClient.delete(`/users/${userId}`)
      return {
        success: true,
        data: response.data
      }
    } catch (error) {
      console.error('Error al eliminar usuario:', error)
      return {
        success: false,
        error: error.response?.data?.message || 'Error al eliminar usuario'
      }
    }
  }
}

export default userService
