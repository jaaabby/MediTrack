import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'

class InventoryService {
  constructor() {
    this.api = axios.create({
      baseURL: API_BASE_URL,
      headers: {
        'Content-Type': 'application/json'
      }
    })
  }

  // Obtener todo el inventario
  async getInventory() {
    try {
      const response = await this.api.get('/medical-supplies/list')
      return response.data.data || []
    } catch (error) {
      console.error('Error al obtener inventario:', error)
      throw error
    }
  }

  // Obtener inventario por bodega
  async getInventoryByStore(storeId) {
    try {
      const response = await this.api.get(`/medical-supplies/inventory/store?store_id=${storeId}`)
      return response.data.data || []
    } catch (error) {
      console.error('Error al obtener inventario por bodega:', error)
      throw error
    }
  }

  // Obtener inventario por proveedor
  async getInventoryBySupplier(supplier) {
    try {
      const response = await this.api.get(`/medical-supplies/inventory/supplier?supplier=${encodeURIComponent(supplier)}`)
      return response.data.data || []
    } catch (error) {
      console.error('Error al obtener inventario por proveedor:', error)
      throw error
    }
  }

  // Obtener todos los insumos médicos (método original)
  async getAllMedicalSupplies() {
    try {
      const response = await this.api.get('/medical-supplies')
      return response.data.data || []
    } catch (error) {
      console.error('Error al obtener insumos médicos:', error)
      throw error
    }
  }

  // Crear insumo médico
  async createMedicalSupply(supply) {
    try {
      const response = await this.api.post('/medical-supplies', supply)
      return response.data.data
    } catch (error) {
      console.error('Error al crear insumo médico:', error)
      throw error
    }
  }

  // Actualizar insumo médico
  async updateMedicalSupply(id, supply) {
    try {
      const response = await this.api.put(`/medical-supplies/${id}`, supply)
      return response.data.data
    } catch (error) {
      console.error('Error al actualizar insumo médico:', error)
      throw error
    }
  }

  // Actualizar batch (lote)
  async updateBatch(id, batchData) {
    try {
      const response = await this.api.put(`/batches/${id}`, batchData)
      return response.data.data
    } catch (error) {
      console.error('Error al actualizar batch:', error)
      throw error
    }
  }

  // Eliminar insumo médico
  async deleteBatch(id) {
    try {
      const response = await this.api.delete(`/batches/${id}`)
      return response.data.data
    } catch (error) {
      console.error('Error al eliminar batch:', error)
      throw error
    }
  }

  async getBatchHistoryWithDetails() {
    try {
      const response = await this.api.get('/batch-histories/details')
      return response.data.data || []
    } catch (error) {
      console.error('Error al obtener historial de lotes:', error)
      throw error
    }
  }

  // Buscar historial por número de lote
  async searchBatchHistoryByBatchNumber(batchNumber) {
    try {
      const response = await this.api.get(`/batch-histories/search/${batchNumber}`)
      return response.data.data || []
    } catch (error) {
      console.error('Error al buscar historial por número de lote:', error)
      throw error
    }
  }
}

export default new InventoryService()
