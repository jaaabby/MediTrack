import axios from 'axios'
import { getApiBaseUrl } from '@/config/api.js'

const API_BASE_URL = getApiBaseUrl()

class SupplyHistoryService {
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

  async getAllSupplyHistoryWithDetails() {
    try {
      const response = await this.api.get('/supply-history/with-details')
      return response.data.data || response.data || []
    } catch (error) {
      console.error('Error en getAllSupplyHistoryWithDetails:', error)
      throw error
    }
  }

  async getConsumptionStatsBySurgery() {
    try {
      const response = await this.api.get('/supply-history/consumption-stats')
      return response.data.data?.consumption_stats || []
    } catch (error) {
      console.error('Error en getConsumptionStatsBySurgery:', error)
      throw error
    }
  }
}

export default new SupplyHistoryService()
