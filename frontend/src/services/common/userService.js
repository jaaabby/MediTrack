import axios from 'axios'
import { getApiBaseUrl } from '@/config/api.js'

const API_BASE_URL = getApiBaseUrl().replace('/api/v1', '/api')

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
      const response = await apiClient.get(`/users/profile?email=${encodeURIComponent(email)}`)
      return {
        success: true,
        data: response.data
      }
    } catch (error) {
      console.error('Error al obtener perfil:', error)
      return {
        success: false,
        error: error.response?.data?.error || error.response?.data?.message || 'Error al obtener información del perfil'
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
        error: error.response?.data?.error || error.response?.data?.message || 'Error al actualizar el perfil'
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

  // Obtener todos los usuarios (solo para administradores)
  async getAllUsers() {
    try {
      const response = await apiClient.get('/v1/users/')
      return {
        success: true,
        data: response.data.data || response.data || []
      }
    } catch (error) {
      console.error('Error al obtener todos los usuarios:', error)
      return {
        success: false,
        error: error.response?.data?.error || error.response?.data?.message || 'Error al obtener lista de usuarios'
      }
    }
  },

  // Crear nuevo usuario (solo para administradores)
  async createUser(userData) {
    try {
      const response = await apiClient.post('/v1/users/', userData)
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
      const response = await apiClient.put(`/v1/users/${userId}`, userData)
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
      const response = await apiClient.delete(`/v1/users/${userId}`)
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
  },

  // Obtener usuarios por rol
  async getUsersByRole(role) {
    try {
      const response = await apiClient.get(`/v1/users/by-role?role=${encodeURIComponent(role)}`)
      return {
        success: true,
        data: response.data.data || []
      }
    } catch (error) {
      console.error('Error obteniendo usuarios por rol:', error)
      return {
        success: false,
        error: error.response?.data?.error || error.message
      }
    }
  },

  // Buscar usuarios (accesible para admin y encargado de bodega)
  async searchUsers(searchTerm) {
    try {
      const response = await apiClient.get(`/v1/users/search?q=${encodeURIComponent(searchTerm)}`)
      return {
        success: true,
        data: response.data.data || []
      }
    } catch (error) {
      console.error('Error buscando usuarios:', error)
      return {
        success: false,
        error: error.response?.data?.error || error.message
      }
    }
  }
}

export default userService
