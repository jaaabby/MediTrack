import axios from 'axios'
import { getApiBaseUrl } from '@/config/api.js'

const API_BASE_URL = getApiBaseUrl()

class MedicalCenterService {
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

  async getAll() {
    try {
      const response = await this.api.get('/medical-centers/')
      return { data: response.data.data || response.data || [] }
    } catch (error) {
      console.error('Error en getAll:', error)
      throw error
    }
  }

  async getAllMedicalCenters() {
    try {
      const response = await this.api.get('/medical-centers/')
      return {
        success: true,
        data: response.data.data || response.data || []
      }
    } catch (error) {
      console.error('Error al obtener centros médicos:', error)
      return {
        success: false,
        error: error.response?.data?.error || error.message
      }
    }
  }
}

export default new MedicalCenterService()
