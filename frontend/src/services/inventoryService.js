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
      const response = await this.api.get('/medical-supplies/inventory')
      return response.data.data || []
    } catch (error) {
      console.error('Error al obtener inventario:', error)
      throw error
    }
  }

  // Obtener inventario por bodega
  async getInventoryByStore(storeId) {
    try {
      const response = await this.api.get(`/medical-supplies/inventory/store/?store_id=${storeId}`)
      return response.data.data || []
    } catch (error) {
      console.error('Error al obtener inventario por bodega:', error)
      throw error
    }
  }

  // Obtener inventario por proveedor
  async getInventoryBySupplier(supplier) {
    try {
      const response = await this.api.get(`/medical-supplies/inventory/supplier/?supplier=${encodeURIComponent(supplier)}`)
      return response.data.data || []
    } catch (error) {
      console.error('Error al obtener inventario por proveedor:', error)
      throw error
    }
  }

  // Obtener todos los insumos médicos (método original)
  async getAllMedicalSupplies() {
    try {
      const response = await this.api.get('/medical-supplies/')
      return response.data.data || []
    } catch (error) {
      console.error('Error al obtener insumos médicos:', error)
      throw error
    }
  }

  // Crear insumo médico básico
  async createMedicalSupply(supply) {
    try {
      const response = await this.api.post('/medical-supplies/', supply)
      return response.data.data
    } catch (error) {
      console.error('Error al crear insumo médico:', error)
      throw error
    }
  }

  // Crear lote (batch)
  async createBatch(batchData) {
    try {
      console.log('Creando lote:', batchData)

      // Convertir fecha al formato RFC3339 que espera Go
      const formattedData = {
        ...batchData,
        expiration_date: batchData.expiration_date + 'T00:00:00Z'
      }

      console.log('Datos formateados:', formattedData)
      const response = await this.api.post('/batches/', formattedData)
      return response.data
    } catch (error) {
      const backendError = error.response?.data?.error || error.message
      console.error('Error al crear lote:', backendError)
      throw new Error(backendError)
    }
  }

  // Crear código de insumo (supply code)
  async createSupplyCode(supplyCodeData) {
    try {
      const response = await this.api.post('/supply-codes/', supplyCodeData)
      return response.data.data
    } catch (error) {
      console.error('Error al crear código de insumo:', error)
      throw error
    }
  }

  // Crear insumo completo con lote y código
  async createCompleteSupply(supplyData) {
    try {
      // 1. Crear el lote primero
      const batchData = {
        expiration_date: supplyData.batch.expiration_date,
        amount: supplyData.batch.amount,
        supplier: supplyData.batch.supplier,
        store_id: supplyData.batch.store_id
      }

      const createdBatch = await this.createBatch(batchData)

      // 2. Crear el código de insumo asociado al lote
      const supplyCodeData = {
        code: supplyData.supply_code.code,
        name: supplyData.supply_code.name,
        code_supplier: supplyData.supply_code.code_supplier,
        batch_id: createdBatch.data.id
      }

      const createdSupplyCode = await this.createSupplyCode(supplyCodeData)

      // 3. Crear el insumo médico individual
      const medicalSupplyData = {
        code: createdSupplyCode.code,
        batch_id: createdBatch.data.id
      }

      const createdSupply = await this.createMedicalSupply(medicalSupplyData)

      // Retornar toda la información combinada
      return {
        supply: createdSupply,
        batch: createdBatch,
        supply_code: createdSupplyCode
      }

    } catch (error) {
      console.error('Error al crear insumo completo:', error)
      throw error
    }
  }

  // Actualizar insumo médico
  async updateMedicalSupply(id, supply) {
    try {
      const response = await this.api.put(`/medical-supplies/${id}/`, supply)
      return response.data.data
    } catch (error) {
      console.error('Error al actualizar insumo médico:', error)
      throw error
    }
  }

  // Actualizar batch (lote)
  async updateBatch(id, batchData) {
    try {
      const response = await this.api.put(`/batches/${id}/`, batchData)
      return response.data.data
    } catch (error) {
      console.error('Error al actualizar batch:', error)
      throw error
    }
  }

  // Eliminar insumo médico
  async deleteBatch(id) {
    try {
      const response = await this.api.delete(`/batches/${id}/`)
      return response.data.data
    } catch (error) {
      console.error('Error al eliminar batch:', error)
      throw error
    }
  }

  async getBatchHistoryWithDetails() {
    try {
      const response = await this.api.get('/batch-history/details/')
      return response.data.data || []
    } catch (error) {
      console.error('Error al obtener historial de lotes:', error)
      throw error
    }
  }

  // Obtener historial de un lote específico por batch_id
  async getBatchHistory(batchId) {
    try {
      const response = await this.api.get(`/batch-history/search/${batchId}`)
      return response.data.data || []
    } catch (error) {
      console.error('Error al obtener historial del lote:', error)
      throw error
    }
  }

  // Buscar historial por número de lote
  async searchBatchHistoryByBatchNumber(batchNumber) {
    try {
      const response = await this.api.get(`/batch-history/search/${batchNumber}`)
      return response.data.data || []
    } catch (error) {
      console.error('Error al buscar historial por número de lote:', error)
      throw error
    }
  }

  // Obtener todas las bodegas
  async getAllStores() {
    try {
      const response = await this.api.get('/stores/')
      return response.data.data || []
    } catch (error) {
      console.error('Error al obtener bodegas:', error)
      throw error
    }
  }

  // Obtener todos los lotes
  async getAllBatches() {
    try {
      const response = await this.api.get('/batches/')
      return response.data.data || []
    } catch (error) {
      console.error('Error al obtener lotes:', error)
      throw error
    }
  }

  // Obtener todos los códigos de insumo
  async getAllSupplyCodes() {
    try {
      const response = await this.api.get('/supply-codes/')
      return response.data.data || []
    } catch (error) {
      console.error('Error al obtener códigos de insumo:', error)
      throw error
    }
  }

  // Buscar insumos por término
  async searchSupplies(searchTerm) {
    try {
      const inventory = await this.getInventory()

      if (!searchTerm) return inventory

      const filtered = inventory.filter(item =>
        item.name?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        item.code?.toString().includes(searchTerm) ||
        item.supplier?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        item.batch_id?.toString().includes(searchTerm)
      )

      return filtered
    } catch (error) {
      console.error('Error al buscar insumos:', error)
      throw error
    }
  }
}

export default new InventoryService()