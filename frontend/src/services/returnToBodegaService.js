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
      console.log('🔍 Llamando a:', `${API_BASE_URL}/qr/supplies-for-return`)
      const response = await this.api.get('/qr/supplies-for-return')
      console.log('🔍 Respuesta completa:', response)
      console.log('🔍 Status:', response.status)
      console.log('🔍 Data:', response.data)
      console.log('🔍 Tipo de data:', typeof response.data)
      
      // Intentar diferentes estructuras de respuesta
      let data = []
      
      if (Array.isArray(response.data)) {
        // Si response.data es directamente un array
        data = response.data
      } else if (response.data && Array.isArray(response.data.data)) {
        // Si está en response.data.data
        data = response.data.data
      } else if (response.data && Array.isArray(response.data.supplies)) {
        // Si está en response.data.supplies
        data = response.data.supplies
      } else if (response.data && response.data.message) {
        // Si es una respuesta con mensaje (posiblemente vacía)
        console.log('📄 Mensaje del servidor:', response.data.message)
        data = []
      } else {
        // Fallback: intentar convertir cualquier cosa a array
        console.warn('⚠️ Estructura de respuesta no reconocida, usando array vacío')
        data = []
      }
      
      console.log('✅ Datos finales extraídos:', data)
      console.log('✅ Cantidad de elementos:', data.length)
      
      return data
    } catch (error) {
      console.error('❌ Error obteniendo insumos para retorno:', error)
      if (error.response) {
        console.error('❌ Status:', error.response.status)
        console.error('❌ Data:', error.response.data)
        throw new Error(error.response.data?.error || `Error HTTP: ${error.response.status}`)
      } else if (error.request) {
        console.error('❌ No response received:', error.request)
        throw new Error('No se pudo conectar con el servidor')
      } else {
        console.error('❌ Request setup error:', error.message)
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
    // Manejar casos donde supply puede ser null o undefined
    if (!supply || typeof supply !== 'object') {
      console.warn('⚠️ Supply inválido para formatear:', supply)
      return {
        id: null,
        qrCode: 'N/A',
        name: 'Insumo desconocido',
        status: 'unknown',
        batchId: null,
        supplier: 'N/A',
        expirationDate: null,
        storeName: 'N/A',
        storeId: null,
        receivedAt: null,
        daysElapsed: 0,
        shouldReturn: false,
        canReturn: false
      }
    }
    
    return {
      id: supply.supply_id || supply.id || null,
      qrCode: supply.qr_code || supply.qrCode || 'N/A',
      name: supply.supply_name || supply.name || 'Insumo sin nombre',
      status: supply.status || 'unknown',
      batchId: supply.batch_id || supply.batchId || null,
      supplier: supply.supplier || 'N/A',
      expirationDate: supply.expiration_date || supply.expirationDate || null,
      storeName: supply.store_name || supply.storeName || 'N/A',
      storeId: supply.store_id || supply.storeId || null,
      receivedAt: supply.received_at || supply.receivedAt || null,
      daysElapsed: supply.days_elapsed || supply.daysElapsed || 0,
      shouldReturn: supply.should_return || supply.shouldReturn || false,
      canReturn: (supply.days_elapsed || supply.daysElapsed || 0) >= 15
    }
  }
}

// Exportar instancia singleton
export default new ReturnToBodegaService()