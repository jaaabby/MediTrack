import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'

class MedicalSpecialtyService {
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

  async getAllSpecialties() {
    try {
      const response = await this.api.get('/medical-specialties/all')
      return response.data.data?.specialties || response.data.specialties || []
    } catch (error) {
      console.error('Error en getAllSpecialties:', error)
      throw error
    }
  }

  async getSpecialtiesPaginated(page = 1, pageSize = 20, search = '') {
    try {
      const params = {
        page: page.toString(),
        page_size: pageSize.toString(),
        ...(search && { search })
      }
      const response = await this.api.get('/medical-specialties', { params })
      return response.data.data || response.data
    } catch (error) {
      console.error('Error en getSpecialtiesPaginated:', error)
      throw error
    }
  }

  async getSpecialtyById(id) {
    try {
      const response = await this.api.get(`/medical-specialties/${id}`)
      return response.data.data || response.data
    } catch (error) {
      console.error('Error en getSpecialtyById:', error)
      throw error
    }
  }

  async createSpecialty(specialtyData) {
    try {
      const response = await this.api.post('/medical-specialties/', specialtyData)
      return response.data.data || response.data
    } catch (error) {
      console.error('Error en createSpecialty:', error)
      throw error
    }
  }

  async updateSpecialty(id, specialtyData) {
    try {
      const response = await this.api.put(`/medical-specialties/${id}`, specialtyData)
      return response.data.data || response.data
    } catch (error) {
      console.error('Error en updateSpecialty:', error)
      throw error
    }
  }

  async deleteSpecialty(id) {
    try {
      const response = await this.api.delete(`/medical-specialties/${id}`)
      return response.data
    } catch (error) {
      console.error('Error en deleteSpecialty:', error)
      throw error
    }
  }

  async searchSpecialties(name) {
    try {
      const response = await this.api.get(`/medical-specialties/search?name=${encodeURIComponent(name)}`)
      return response.data.data?.specialties || response.data.specialties || []
    } catch (error) {
      console.error('Error en searchSpecialties:', error)
      throw error
    }
  }
}

export default new MedicalSpecialtyService()

