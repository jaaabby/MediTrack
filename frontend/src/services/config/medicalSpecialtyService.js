import axios from 'axios'
import { getApiBaseUrl } from '@/config/api.js'

const API_BASE_URL = getApiBaseUrl()

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

  async getActiveSpecialties() {
    try {
      const response = await this.api.get('/medical-specialties/active')
      return response.data.data?.specialties || response.data.specialties || []
    } catch (error) {
      console.error('Error en getActiveSpecialties:', error)
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
}

export default new MedicalSpecialtyService()
