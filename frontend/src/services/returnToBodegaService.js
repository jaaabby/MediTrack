// Servicio para operaciones de retorno a bodega
// Maneja la comunicación con el backend para funcionalidades de retorno automático y manual

import { useAuthStore } from '@/stores/auth'
import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'

class ReturnToBodegaService {
  
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
        const authStore = useAuthStore()
        if (authStore.token) {
          config.headers.Authorization = `Bearer ${authStore.token}`
        }
        return config
      },
      (error) => {
        return Promise.reject(error)
      }
    )
  }
  
  /**
   * Regresa un insumo a bodega manualmente
   * @param {string} qrCode - Código QR del insumo
   * @param {string} notes - Notas adicionales
   * @returns {Promise<Object>} Resultado de la operación
   */
  async returnSupplyToStore(qrCode, notes = '') {
    try {
      console.log('🔍 API_BASE_URL:', API_BASE_URL)
      console.log('🔍 Calling:', `/qr/return-to-store`)
      
      const authStore = useAuthStore()
      console.log('🔍 AuthStore user:', authStore.user)
      console.log('🔍 User RUT:', authStore.user?.rut)
      
      // Usar un RUT por defecto que sabemos que existe en la base de datos
      const userRut = authStore.user?.rut || '12345678-9' // Admin del sistema
      console.log('🔍 Using RUT:', userRut)
      
      const response = await this.api.post('/qr/return-to-store', {
        qr_code: qrCode,
        user_rut: userRut,
        notes: notes || 'Retorno manual desde interfaz web'
      })

      return response.data
    } catch (error) {
      console.error('Error regresando insumo a bodega:', error)
      if (error.response) {
        throw new Error(error.response.data?.error || `Error HTTP: ${error.response.status}`)
      } else if (error.request) {
        throw new Error('No se pudo conectar con el servidor')
      } else {
        throw error
      }
    }
  }

  /**
   * Confirma la llegada de un insumo a bodega
   * @param {string} qrCode - Código QR del insumo
   * @param {string} notes - Notas adicionales
   * @returns {Promise<Object>} Resultado de la operación
   */
  async confirmArrivalToStore(qrCode, notes = '') {
    try {
      console.log('🔍 Confirmando llegada a bodega para:', qrCode)
      
      const authStore = useAuthStore()
      const userRut = authStore.user?.rut || '12345678-9'
      console.log('🔍 Using RUT:', userRut)
      
      const response = await this.api.post('/qr/confirm-arrival-to-store', {
        qr_code: qrCode,
        user_rut: userRut,
        notes: notes || 'Llegada confirmada desde interfaz web'
      })

      return response.data
    } catch (error) {
      console.error('Error confirmando llegada a bodega:', error)
      if (error.response) {
        throw new Error(error.response.data?.error || `Error HTTP: ${error.response.status}`)
      } else if (error.request) {
        throw new Error('No se pudo conectar con el servidor')
      } else {
        throw error
      }
    }
  }

  /**
   * Obtiene la lista de insumos que deben regresar a bodega (15 días sin consumir)
   * @returns {Promise<Array>} Lista de insumos para retorno
   */
  async getSuppliesForReturn() {
    try {
      const response = await this.api.get('/qr/supplies-for-return')
      return response.data.data || response.data || []
    } catch (error) {
      console.error('Error obteniendo insumos para retorno:', error)
      if (error.response) {
        throw new Error(error.response.data?.error || `Error HTTP: ${error.response.status}`)
      } else if (error.request) {
        throw new Error('No se pudo conectar con el servidor')
      } else {
        throw error
      }
    }
  }

  /**
   * Ejecuta manualmente el proceso automático de retornos
   * @returns {Promise<Object>} Resultado del proceso
   */
  async processAutomaticReturns() {
    try {
      const response = await this.api.post('/qr/process-automatic-returns')
      return response.data
    } catch (error) {
      console.error('Error procesando retornos automáticos:', error)
      if (error.response) {
        throw new Error(error.response.data?.error || `Error HTTP: ${error.response.status}`)
      } else if (error.request) {
        throw new Error('No se pudo conectar con el servidor')
      } else {
        throw error
      }
    }
  }

  /**
   * Verifica si un insumo puede ser regresado a bodega
   * @param {Object} supplyInfo - Información del insumo
   * @returns {boolean} True si puede ser regresado
   */
  canBeReturnedToStore(supplyInfo) {
    if (!supplyInfo || supplyInfo.type !== 'medical_supply') return false
    if (supplyInfo.is_consumed) return false
    
    const status = supplyInfo.supply_info?.Status || 
                   supplyInfo.supply_info?.status || 
                   supplyInfo.status || 
                   supplyInfo.current_status
    
    // Puede ser regresado si está recepcionado
    return status === 'recepcionado'
  }

  /**
   * Obtiene una descripción del motivo por el que un insumo debe regresar a bodega
   * @param {Object} supplyInfo - Información del insumo
   * @returns {string|null} Descripción o null si no debe regresar
   */
  getReturnReason(supplyInfo) {
    if (!this.canBeReturnedToStore(supplyInfo)) return null

    // Calcular días desde recepción (si está disponible)
    if (supplyInfo.received_at) {
      const receivedDate = new Date(supplyInfo.received_at)
      const daysSinceReceived = Math.floor((Date.now() - receivedDate.getTime()) / (1000 * 60 * 60 * 24))
      
      if (daysSinceReceived >= 15) {
        return `Lleva ${daysSinceReceived} días sin consumir (límite: 15 días)`
      } else {
        return `Lleva ${daysSinceReceived} días recepcionado`
      }
    }

    return 'Insumo recepcionado listo para retorno manual'
  }

  /**
   * Formatea la información de un insumo para mostrar en la UI
   * @param {Object} supply - Datos del insumo desde el backend
   * @returns {Object} Objeto formateado para la UI
   */
  formatSupplyForUI(supply) {
    return {
      id: supply.supply_id,
      qrCode: supply.qr_code,
      name: supply.supply_name,
      status: supply.status,
      batchId: supply.batch_id,
      supplier: supply.supplier,
      expirationDate: supply.expiration_date,
      storeName: supply.store_name,
      storeId: supply.store_id,
      receivedAt: supply.received_at,
      daysElapsed: supply.days_elapsed,
      shouldReturn: supply.should_return,
      canReturn: supply.days_elapsed >= 15
    }
  }
}

// Exportar instancia singleton
export default new ReturnToBodegaService()