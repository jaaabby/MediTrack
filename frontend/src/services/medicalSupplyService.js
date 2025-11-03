import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'

class MedicalSupplyService {
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

  async createMedicalSupply(supplyData) {
    try {
      const response = await this.api.post('/medical-supplies/', supplyData)
      return response.data
    } catch (error) {
      console.error('Error al crear insumo médico:', error)
      throw error
    }
  }

  async getAllMedicalSupplies() {
    try {
      const response = await this.api.get('/medical-supplies/')
      return response.data.data || response.data || []
    } catch (error) {
      console.error('Error al obtener insumos médicos:', error)
      throw error
    }
  }

  async getMedicalSupplyByID(id) {
    try {
      const response = await this.api.get(`/medical-supplies/${id}/`)
      return response.data.data || response.data
    } catch (error) {
      console.error('Error al obtener insumo médico:', error)
      throw error
    }
  }

  async updateMedicalSupply(id, supplyData) {
    try {
      const response = await this.api.put(`/medical-supplies/${id}/`, supplyData)
      return response.data
    } catch (error) {
      console.error('Error al actualizar insumo médico:', error)
      throw error
    }
  }

  async deleteMedicalSupply(id) {
    try {
      const response = await this.api.delete(`/medical-supplies/${id}/`)
      return response.data
    } catch (error) {
      console.error('Error al eliminar insumo médico:', error)
      throw error
    }
  }

  // ========================
  // INVENTARIO
  // ========================

  async getInventoryList() {
    try {
      const response = await this.api.get('/medical-supplies/inventory/')
      return response.data.data || response.data || []
    } catch (error) {
      console.error('Error al obtener lista de inventario:', error)
      throw error
    }
  }

  async getInventoryListAdvanced() {
    try {
      const response = await this.api.get('/medical-supplies/inventory/advanced/')
      return response.data.data || response.data || []
    } catch (error) {
      console.error('Error al obtener inventario avanzado:', error)
      throw error
    }
  }

  // ========================
  // RUTAS QR Y FUNCIONALIDADES AVANZADAS
  // ========================

  async getMedicalSupplyByQR(qrcode) {
    try {
      const response = await this.api.get(`/medical-supplies/qr/${qrcode}/`)
      return response.data.data || response.data
    } catch (error) {
      console.error('Error al obtener insumo por QR:', error)
      throw error
    }
  }

  async getSupplyWithBatchInfo(qrcode) {
    try {
      const response = await this.api.get(`/medical-supplies/details/${qrcode}/`)
      return response.data.data || response.data
    } catch (error) {
      console.error('Error al obtener insumo con información de lote:', error)
      throw error
    }
  }

  async getIndividualSuppliesByCode(code) {
    try {
      const response = await this.api.get(`/medical-supplies/code/${code}/`)
      return response.data.data || response.data || []
    } catch (error) {
      console.error('Error al obtener insumos individuales por código:', error)
      throw error
    }
  }

  async getAvailableSuppliesByBatch(batchId) {
    try {
      const response = await this.api.get(`/medical-supplies/batch/${batchId}/available/`)
      return response.data.data || response.data || []
    } catch (error) {
      console.error('Error al obtener insumos disponibles por lote:', error)
      throw error
    }
  }

  // ========================
  // CREAR MÚLTIPLES Y CONSUMO
  // ========================

  async createMultipleSupplies(suppliesData) {
    try {
      const response = await this.api.post('/medical-supplies/create-multiple/', suppliesData)
      return response.data
    } catch (error) {
      console.error('Error al crear múltiples insumos:', error)
      throw error
    }
  }

  async consumeSupply(consumeData) {
    try {
      const response = await this.api.post('/medical-supplies/consume/', consumeData)
      return response.data
    } catch (error) {
      console.error('Error al consumir insumo:', error)
      throw error
    }
  }

  async syncBatchAmounts(syncData) {
    try {
      const response = await this.api.post('/medical-supplies/sync-amounts/', syncData)
      return response.data
    } catch (error) {
      console.error('Error al sincronizar cantidades de lotes:', error)
      throw error
    }
  }

  // ========================
  // ALERTAS DE INSUMOS NO CONSUMIDOS
  // ========================

  async getUnconsumedSupplies(days = 15) {
    try {
      const response = await this.api.get(`/medical-supplies/unconsumed/?days=${days}`)
      return response.data.data || response.data || []
    } catch (error) {
      console.error('Error al obtener insumos no consumidos:', error)
      throw error
    }
  }

  async checkUnconsumedSupplies(checkData = {}) {
    try {
      const response = await this.api.post('/medical-supplies/check-unconsumed/', checkData)
      return response.data
    } catch (error) {
      console.error('Error al verificar insumos no consumidos:', error)
      throw error
    }
  }

  // ========================
  // ALERTAS DE STOCK BAJO PARA INSUMOS INDIVIDUALES
  // ========================

  async checkLowStockForIndividualSupply(supplyCode) {
    try {
      const response = await this.api.post(`/medical-supplies/${supplyCode}/check-low-stock/`)
      return response.data
    } catch (error) {
      console.error('Error al verificar alerta de stock bajo para insumo individual:', error)
      throw error
    }
  }

  async checkAllIndividualSuppliesLowStock() {
    try {
      const response = await this.api.post('/medical-supplies/check-all-low-stock/')
      return response.data
    } catch (error) {
      console.error('Error al verificar alertas de stock bajo para todos los insumos individuales:', error)
      throw error
    }
  }

  // ========================
  // RETORNOS A BODEGA
  // ========================

  async getSuppliesForReturn() {
    try {
      const response = await this.api.get('/qr/supplies-for-return')
      return response.data.data || response.data || []
    } catch (error) {
      console.error('Error al obtener insumos para retorno:', error)
      throw error
    }
  }
}

export default new MedicalSupplyService()
