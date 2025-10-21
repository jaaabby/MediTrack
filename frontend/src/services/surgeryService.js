import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'

class SurgeryService {
  constructor() {
    this.api = axios.create({
      baseURL: API_BASE_URL,
      headers: {
        'Content-Type': 'application/json'
      }
    })

    // Interceptor para agregar token de autenticación
    this.api.interceptors.request.use(
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
  }

  // ========================
  // CRUD BÁSICO
  // ========================

  async createSurgery(surgeryData) {
    try {
      const response = await this.api.post('/surgeries/', surgeryData)
      return response.data
    } catch (error) {
      console.error('Error al crear tipo de cirugía:', error)
      throw error
    }
  }

  async getAllSurgeries() {
    try {
      const response = await this.api.get('/surgeries/all')
      return response.data.data?.surgeries || response.data.surgeries || []
    } catch (error) {
      console.error('Error al obtener tipos de cirugía:', error)
      throw error
    }
  }

  async getSurgeriesPaginated(page = 1, limit = 10) {
    try {
      const response = await this.api.get(`/surgeries/?page=${page}&limit=${limit}`)
      return response.data
    } catch (error) {
      console.error('Error al obtener tipos de cirugía paginados:', error)
      throw error
    }
  }

  async getSurgeryByID(id) {
    try {
      const response = await this.api.get(`/surgeries/${id}`)
      return response.data.data || response.data
    } catch (error) {
      console.error('Error al obtener tipo de cirugía:', error)
      throw error
    }
  }

  async updateSurgery(id, surgeryData) {
    try {
      const response = await this.api.put(`/surgeries/${id}`, surgeryData)
      return response.data
    } catch (error) {
      console.error('Error al actualizar tipo de cirugía:', error)
      throw error
    }
  }

  async deleteSurgery(id) {
    try {
      const response = await this.api.delete(`/surgeries/${id}`)
      return response.data
    } catch (error) {
      console.error('Error al eliminar tipo de cirugía:', error)
      throw error
    }
  }

  // ========================
  // BÚSQUEDA
  // ========================

  async searchSurgeries(searchTerm) {
    try {
      const response = await this.api.get(`/surgeries/search?name=${encodeURIComponent(searchTerm)}`)
      return response.data.data?.surgeries || response.data.surgeries || []
    } catch (error) {
      console.error('Error al buscar tipos de cirugía:', error)
      throw error
    }
  }
}

export default new SurgeryService()
