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

    // Interceptor para agregar el token de autenticación
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

  // Obtener todo el inventario - Compatible con ambas versiones
  async getInventory() {
    try {
      // Intentar primero con el endpoint correcto
      const response = await this.api.get('/medical-supplies/inventory')
      return response.data.data || response.data || []
    } catch (error) {
      // Si falla, intentar con el endpoint de la versión avanzada
      try {
        const response = await this.api.get('/medical-supplies/inventory/advanced')
        return response.data.data || response.data || []
      } catch (fallbackError) {
        console.error('Error al obtener inventario:', error)
        throw error
      }
    }
  }

  // Obtener inventario por bodega
  async getInventoryByStore(storeId) {
    try {
      const response = await this.api.get(`/medical-supplies/inventory/store?store_id=${storeId}`)
      return response.data.data || response.data || []
    } catch (error) {
      // Fallback sin parámetros de query
      try {
        const response = await this.api.get(`/medical-supplies/inventory/store/?store_id=${storeId}`)
        return response.data.data || response.data || []
      } catch (fallbackError) {
        console.error('Error al obtener inventario por bodega:', error)
        throw error
      }
    }
  }

  // Obtener inventario por proveedor
  async getInventoryBySupplier(supplier) {
    try {
      const response = await this.api.get(`/medical-supplies/inventory/supplier?supplier=${encodeURIComponent(supplier)}`)
      return response.data.data || response.data || []
    } catch (error) {
      // Fallback con slash final
      try {
        const response = await this.api.get(`/medical-supplies/inventory/supplier/?supplier=${encodeURIComponent(supplier)}`)
        return response.data.data || response.data || []
      } catch (fallbackError) {
        console.error('Error al obtener inventario por proveedor:', error)
        throw error
      }
    }
  }

  // Obtener todos los insumos médicos (método original)
  async getAllMedicalSupplies() {
    try {
      const response = await this.api.get('/medical-supplies/')
      return response.data.data || response.data || []
    } catch (error) {
      console.error('Error al obtener insumos médicos:', error)
      throw error
    }
  }

  // Crear insumo médico básico
  async createMedicalSupply(supply) {
    try {
      const response = await this.api.post('/medical-supplies/', supply)
      return response.data.data || response.data
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
      return response.data.data || response.data
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
        batch_id: createdBatch.data?.id || createdBatch.id
      }

      const createdSupplyCode = await this.createSupplyCode(supplyCodeData)

      // 3. Crear el insumo médico individual
      const medicalSupplyData = {
        code: createdSupplyCode.code,
        batch_id: createdBatch.data?.id || createdBatch.id
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


  // Crear lote completo con múltiples insumos individuales (MÉTODO CORRECTO)
  async createBatchWithIndividualSupplies(batchData) {
    try {
      console.log('Creando lote con insumos individuales:', batchData)

      // Usar el endpoint correcto que ya existe en el backend
      const response = await this.api.post('/batches/create-with-supplies', {
        batch: {
          expiration_date: batchData.batch.expiration_date,
          amount: parseInt(batchData.batch.amount),
          supplier: batchData.batch.supplier,
          store_id: parseInt(batchData.batch.store_id),
          expiration_alert_days: batchData.batch.expiration_alert_days ? parseInt(batchData.batch.expiration_alert_days) : undefined
        },
        supply_code: {
          code: parseInt(batchData.supply_code.code),
          name: batchData.supply_code.name,
          code_supplier: parseInt(batchData.supply_code.code_supplier)
        },
        individual_count: parseInt(batchData.batch.amount) // Esta es la clave - cantidad de insumos individuales
      })

      console.log('Respuesta del backend:', response.data)
      return response.data
    } catch (error) {
      const backendError = error.response?.data?.error || error.message
      console.error('Error al crear lote con insumos individuales:', backendError)
      throw new Error(backendError)
    }
  }

  // Crear historial de lote
  async createBatchHistory(batchId, userRUT = "12345678-9", batchData) {
    try {
      const userName = 'Juan Pérez'; // Puedes reemplazar por el nombre real si lo tienes
      const now = new Date().toISOString();
      // Usar los datos reales del lote si están disponibles
      const realBatchData = batchData || {
        expiration_date: '',
        amount: '',
        supplier: '',
        store_id: ''
      };
      const response = await this.api.post('/batch-history/', {
        date_time: now,
        change_details: 'Lote creado',
        previous_values: '{}',
        new_values: JSON.stringify(realBatchData),
        user_name: userName,
        batch_id: batchId,
        user_rut: userRUT,
        batch_number: batchId
      });
      return response.data;
    } catch (error) {
      console.error('Error al crear historial de lote:', error);
      return null;
    }
  }

  // Actualizar insumo médico
  async updateMedicalSupply(id, supply) {
    try {
      const response = await this.api.put(`/medical-supplies/${id}/`, supply)
      return response.data.data || response.data
    } catch (error) {
      console.error('Error al actualizar insumo médico:', error)
      throw error
    }
  }

  // Actualizar batch (lote) - Compatible con ambas versiones
  async updateBatch(id, batchData) {
    try {
      const response = await this.api.put(`/batches/${id}/`, batchData)
      return response.data.data || response.data
    } catch (error) {
      console.error('Error al actualizar batch:', error)
      throw error
    }
  }

  // Eliminar batch (lote) - Compatible con ambas versiones
  async deleteBatch(id) {
    try {
      const response = await this.api.delete(`/batches/${id}/`)
      return response.data.data || response.data
    } catch (error) {
      console.error('Error al eliminar batch:', error)
      throw error
    }
  }

  // Obtener historial de lotes con detalles - Compatible con ambas versiones
  async getBatchHistoryWithDetails() {
    try {
      // Usar el endpoint correcto expuesto en el backend
      const response = await this.api.get('/batch-history/details')
      const data = response.data.data || response.data || []
      
      // Parsear los JSON strings de previous_values y new_values con verificación de tipo
      return data.map(entry => {
        let previousValues = entry.previous_values;
        let newValues = entry.new_values;
        
        // Solo parsear si es string y parece JSON
        if (typeof previousValues === 'string' && previousValues.trim().startsWith('{')) {
          try {
            previousValues = JSON.parse(previousValues);
          } catch (e) {
            console.warn('Could not parse previous_values:', previousValues, e);
            previousValues = {};
          }
        } else if (!previousValues) {
          previousValues = {};
        }
        
        if (typeof newValues === 'string' && newValues.trim().startsWith('{')) {
          try {
            newValues = JSON.parse(newValues);
          } catch (e) {
            console.warn('Could not parse new_values:', newValues, e);
            newValues = {};
          }
        } else if (!newValues) {
          newValues = {};
        }
        
        return {
          ...entry,
          previous_values: previousValues,
          new_values: newValues
        };
      })
    } catch (error) {
      console.error('Error al obtener historial de lotes:', error)
      throw error
    }
  }

  // Obtener historial de un lote específico por batch_id - Compatible con ambas versiones
  async getBatchHistory(batchId) {
    try {
      // Intentar primero con el endpoint de la versión original
      const response = await this.api.get(`/batch-histories/search/${batchId}`)
      const data = response.data.data || response.data || []
      
      // Parsear los JSON strings de previous_values y new_values con verificación de tipo
      return data.map(entry => {
        let previousValues = entry.previous_values;
        let newValues = entry.new_values;
        
        // Solo parsear si es string y parece JSON
        if (typeof previousValues === 'string' && previousValues.trim().startsWith('{')) {
          try {
            previousValues = JSON.parse(previousValues);
          } catch (e) {
            console.warn('Could not parse previous_values:', previousValues, e);
            previousValues = {};
          }
        } else if (!previousValues) {
          previousValues = {};
        }
        
        if (typeof newValues === 'string' && newValues.trim().startsWith('{')) {
          try {
            newValues = JSON.parse(newValues);
          } catch (e) {
            console.warn('Could not parse new_values:', newValues, e);
            newValues = {};
          }
        } else if (!newValues) {
          newValues = {};
        }
        
        return {
          ...entry,
          previous_values: previousValues,
          new_values: newValues
        };
      })
    } catch (error) {
      // Si falla, intentar con el endpoint de la nueva versión
      try {
        const response = await this.api.get(`/batch-history/search/${batchId}`)
        const data = response.data.data || response.data || []
        
        // Parsear los JSON strings de previous_values y new_values con verificación de tipo
        return data.map(entry => {
          let previousValues = entry.previous_values;
          let newValues = entry.new_values;
          
          // Solo parsear si es string y parece JSON
          if (typeof previousValues === 'string' && previousValues.trim().startsWith('{')) {
            try {
              previousValues = JSON.parse(previousValues);
            } catch (e) {
              console.warn('Could not parse previous_values:', previousValues, e);
              previousValues = {};
            }
          } else if (!previousValues) {
            previousValues = {};
          }
          
          if (typeof newValues === 'string' && newValues.trim().startsWith('{')) {
            try {
              newValues = JSON.parse(newValues);
            } catch (e) {
              console.warn('Could not parse new_values:', newValues, e);
              newValues = {};
            }
          } else if (!newValues) {
            newValues = {};
          }
          
          return {
            ...entry,
            previous_values: previousValues,
            new_values: newValues
          };
        })
      } catch (fallbackError) {
        console.error('Error al obtener historial del lote:', error)
        throw error
      }
    }
  }

  // Buscar historial por número de lote - Compatible con ambas versiones
  async searchBatchHistoryByBatchNumber(batchNumber) {
    try {
      // Intentar primero con el endpoint de la versión original
      const response = await this.api.get(`/batch-histories/search/${batchNumber}`)
      return response.data.data || response.data || []
    } catch (error) {
      // Si falla, intentar con el endpoint de la nueva versión
      try {
        const response = await this.api.get(`/batch-history/search/${batchNumber}`)
        return response.data.data || response.data || []
      } catch (fallbackError) {
        console.error('Error al buscar historial por número de lote:', error)
        throw error
      }
    }
  }

  // Obtener todas las bodegas
  async getAllStores() {
    try {
      const response = await this.api.get('/stores/')
      return response.data.data || response.data || []
    } catch (error) {
      console.error('Error al obtener bodegas:', error)
      throw error
    }
  }

  // Obtener todos los lotes
  async getAllBatches() {
    try {
      const response = await this.api.get('/batches/')
      return response.data.data || response.data || []
    } catch (error) {
      console.error('Error al obtener lotes:', error)
      throw error
    }
  }

  // Obtener todos los códigos de insumo
  async getAllSupplyCodes() {
    try {
      const response = await this.api.get('/supply-codes/')
      return response.data.data || response.data || []
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

  // Obtener insumos individuales disponibles por batch
  async getAvailableSuppliesByBatch(batchId) {
    try {
      const response = await this.api.get(`/medical-supplies/batch/${batchId}/available`)
      // Ajuste para usar el array correcto de la respuesta
      return response.data.data?.available_supplies || []
    } catch (error) {
      console.error('Error al obtener insumos individuales por lote:', error)
      throw error
    }
  }

  // ===== NUEVOS MÉTODOS PARA INVENTARIO POR UBICACIÓN =====

  // Obtener inventario de bodega con filtros
  async getStoreInventory(filters = {}) {
    try {
      const params = new URLSearchParams()
      if (filters.store_id) params.append('store_id', filters.store_id)
      if (filters.surgery_id && filters.surgery_id !== 'undefined') params.append('surgery_id', filters.surgery_id)
      if (filters.supply_code) params.append('supply_code', filters.supply_code)
      if (filters.supplier) params.append('supplier', filters.supplier)
      if (filters.near_expiration) params.append('near_expiration', 'true')
      if (filters.low_stock) params.append('low_stock', 'true')
      if (filters.page) params.append('page', filters.page)
      if (filters.page_size) params.append('page_size', filters.page_size)

      const response = await this.api.get(`/inventory/store?${params.toString()}`)
      // El backend devuelve {inventory: [...], total: ..., page: ...}
      // Extraer el array de inventory correctamente
      const data = response.data.data || response.data
      if (data && typeof data === 'object' && Array.isArray(data.inventory)) {
        return data.inventory
      }
      return Array.isArray(data) ? data : []
    } catch (error) {
      console.error('Error al obtener inventario de bodega:', error)
      throw error
    }
  }

  // Obtener inventario de pabellón
  async getPavilionInventory(pavilionId, includeInTransit = false) {
    try {
      const params = new URLSearchParams()
      if (includeInTransit) params.append('include_in_transit', 'true')

      const response = await this.api.get(`/inventory/pavilion/${pavilionId}?${params.toString()}`)
      // El backend devuelve {pavilion_id: ..., inventory: [...], count: ...}
      const data = response.data.data || response.data
      if (data && typeof data === 'object' && Array.isArray(data.inventory)) {
        return data.inventory
      }
      return Array.isArray(data) ? data : []
    } catch (error) {
      console.error('Error al obtener inventario de pabellón:', error)
      throw error
    }
  }

  // Obtener resumen general del inventario
  async getInventorySummary(medicalCenterId = null) {
    try {
      const params = medicalCenterId ? `?medical_center_id=${medicalCenterId}` : ''
      const response = await this.api.get(`/inventory/summary${params}`)
      return response.data.data || response.data || {}
    } catch (error) {
      console.error('Error al obtener resumen de inventario:', error)
      throw error
    }
  }

  // Obtener inventario agrupado por tipo de cirugía
  async getInventoryBySurgeryType(storeId = null) {
    try {
      const params = storeId ? `?store_id=${storeId}` : ''
      const response = await this.api.get(`/inventory/by-surgery${params}`)
      // El backend devuelve {inventory: [...], count: ...}
      const data = response.data.data || response.data
      if (data && typeof data === 'object' && Array.isArray(data.inventory)) {
        return data.inventory
      }
      return Array.isArray(data) ? data : []
    } catch (error) {
      console.error('Error al obtener inventario por tipo de cirugía:', error)
      throw error
    }
  }

  // Obtener todos los pabellones
  async getAllPavilions() {
    try {
      const response = await this.api.get('/pavilions/')
      return response.data.data || response.data || []
    } catch (error) {
      console.error('Error al obtener pabellones:', error)
      throw error
    }
  }

  // Obtener reporte de transferencias
  async getTransferReport(startDate, endDate, groupBy = 'date') {
    try {
      const params = new URLSearchParams()
      params.append('start_date', startDate)
      params.append('end_date', endDate)
      params.append('group_by', groupBy)

      const response = await this.api.get(`/inventory/reports/transfers?${params.toString()}`)
      return response.data.data || response.data || []
    } catch (error) {
      console.error('Error al obtener reporte de transferencias:', error)
      throw error
    }
  }

  // Sincronizar inventario (para admin)
  async syncInventory() {
    try {
      const response = await this.api.post('/inventory/sync')
      return response.data.data || response.data || {}
    } catch (error) {
      console.error('Error al sincronizar inventario:', error)
      throw error
    }
  }
}

export default new InventoryService()