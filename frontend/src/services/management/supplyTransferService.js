import axios from 'axios'
import { getApiBaseUrl } from '@/config/api.js'

const API_BASE_URL = getApiBaseUrl()

class SupplyTransferService {
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
  // CREAR TRANSFERENCIAS
  // ========================

  async transferToPavilion(transferData) {
    try {
      const response = await this.api.post('/transfers/to-pavilion', transferData)
      return response.data
    } catch (error) {
      console.error('Error al dejar listo para retiro:', error)
      throw error
    }
  }

  async returnToStore(transferData) {
    try {
      const response = await this.api.post('/transfers/return-to-store', transferData)
      return response.data
    } catch (error) {
      console.error('Error al retornar a bodega:', error)
      throw error
    }
  }

  // ========================
  // CONFIRMAR Y GESTIONAR TRANSFERENCIAS
  // ========================

  async confirmReception(code, confirmationData = {}) {
    try {
      const response = await this.api.post(`/transfers/${code}/confirm`, confirmationData)
      return response.data
    } catch (error) {
      console.error('Error al confirmar recepción:', error)
      throw error
    }
  }

  async cancelTransfer(code, cancellationData = {}) {
    try {
      const response = await this.api.post(`/transfers/${code}/cancel`, cancellationData)
      return response.data
    } catch (error) {
      console.error('Error al cancelar transferencia:', error)
      throw error
    }
  }

  // ========================
  // CONSULTAR TRANSFERENCIAS
  // ========================

  async getTransferByCode(code) {
    try {
      const response = await this.api.get(`/transfers/${code}`)
      return response.data.data || response.data
    } catch (error) {
      console.error('Error al obtener transferencia por código:', error)
      throw error
    }
  }

  async getTransfers(filters = {}) {
    try {
      const queryParams = new URLSearchParams()
      
      if (filters.status) queryParams.append('status', filters.status)
      if (filters.from_date) queryParams.append('from_date', filters.from_date)
      if (filters.to_date) queryParams.append('to_date', filters.to_date)
      if (filters.pavilion_id) queryParams.append('pavilion_id', filters.pavilion_id)
      if (filters.store_id) queryParams.append('store_id', filters.store_id)
      // Solicitar todos los registros ya que la paginación se maneja en el cliente
      queryParams.append('page_size', filters.page_size || 9999)
      if (filters.page) queryParams.append('page', filters.page)
      
      const queryString = queryParams.toString()
      const url = queryString ? `/transfers?${queryString}` : '/transfers'
      
      const response = await this.api.get(url)
      return response.data.data || response.data || []
    } catch (error) {
      console.error('Error al obtener transferencias:', error)
      throw error
    }
  }
}

export default new SupplyTransferService()
