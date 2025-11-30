import axios from 'axios'
import { getApiBaseUrl } from '@/config/api.js'

const API_BASE_URL = getApiBaseUrl()

class PavilionService {
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

  async getAllPavilions() {
    try {
      const response = await this.api.get('/pavilions/')
      return response.data.data || response.data || []
    } catch (error) {
      console.error('Error en getAllPavilions:', error)
      throw error
    }
  }

  async getByMedicalCenter(medicalCenterId) {
    try {
      const allPavilions = await this.getAllPavilions()
      const filtered = allPavilions.filter(p => p.medical_center_id === parseInt(medicalCenterId))
      return { data: filtered }
    } catch (error) {
      console.error('Error en getByMedicalCenter:', error)
      throw error
    }
  }

  async createPavilion(pavilionData) {
    try {
      const response = await this.api.post('/pavilions/', pavilionData)
      return response.data
    } catch (error) {
      console.error('Error en createPavilion:', error)
      throw error
    }
  }

  async updatePavilion(id, pavilionData) {
    try {
      const response = await this.api.put(`/pavilions/${id}`, pavilionData)
      return response.data
    } catch (error) {
      console.error('Error en updatePavilion:', error)
      throw error
    }
  }

  async deletePavilion(id) {
    try {
      const response = await this.api.delete(`/pavilions/${id}`)
      return response.data
    } catch (error) {
      console.error('Error en deletePavilion:', error)
      throw error
    }
  }
}

export default new PavilionService()