import axios from 'axios'
import { getApiBaseUrl } from '@/config/api.js'

// Usar la misma función centralizada para obtener la URL de la API
const API_BASE_URL = getApiBaseUrl().replace('/api/v1', '')

class CartService {
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
  // OPERACIONES DE CARRITO
  // ========================

  /**
   * Obtiene el carrito asociado a una solicitud
   * @param {number} requestId - ID de la solicitud
   * @returns {Promise} Carrito con sus items
   */
  async getCartByRequestId(requestId) {
    try {
      const response = await this.api.get(`/api/carts/request/${requestId}`)
      return response.data
    } catch (error) {
      // No loguear error si es 404 (carrito no existe todavía, es un estado válido)
      if (error.response?.status !== 404) {
        console.error('Error al obtener carrito por solicitud:', error)
      }
      throw error
    }
  }

  /**
   * Obtiene un carrito por su ID
   * @param {number} cartId - ID del carrito
   * @returns {Promise} Carrito con sus items
   */
  async getCartById(cartId) {
    try {
      const response = await this.api.get(`/api/carts/${cartId}`)
      return response.data
    } catch (error) {
      console.error('Error al obtener carrito:', error)
      throw error
    }
  }

  /**
   * Obtiene el carrito asociado a un código QR
   * @param {string} qrCode - Código QR
   * @returns {Promise} Carrito con sus items
   */
  async getCartByQRCode(qrCode) {
    try {
      const response = await this.api.get(`/api/carts/qr/${qrCode}`)
      return response.data
    } catch (error) {
      console.error('Error al obtener carrito por QR:', error)
      throw error
    }
  }

  /**
   * Obtiene los detalles completos de un carrito
   * @param {number} cartId - ID del carrito
   * @returns {Promise} Detalles del carrito
   */
  async getCartDetails(cartId) {
    try {
      const response = await this.api.get(`/api/carts/${cartId}/details`)
      return response.data
    } catch (error) {
      console.error('Error al obtener detalles del carrito:', error)
      throw error
    }
  }

  /**
   * Obtiene todos los carritos con paginación
   * @param {number} page - Número de página
   * @param {number} pageSize - Tamaño de página
   * @param {string} status - Filtro de estado
   * @returns {Promise} Lista de carritos
   */
  async getAllCarts(page = 1, pageSize = 10, status = '') {
    try {
      const params = { page, pageSize }
      if (status) params.status = status
      
      const response = await this.api.get('/api/carts', { params })
      return response.data
    } catch (error) {
      console.error('Error al obtener carritos:', error)
      throw error
    }
  }

  /**
   * Crea un carrito para una solicitud (manual)
   * @param {number} requestId - ID de la solicitud
   * @returns {Promise} Carrito creado
   */
  async createCartForRequest(requestId) {
    try {
      const response = await this.api.post(`/api/carts/request/${requestId}`)
      return response.data
    } catch (error) {
      console.error('Error al crear carrito:', error)
      throw error
    }
  }

  /**
   * Agrega un item al carrito
   * @param {number} cartId - ID del carrito
   * @param {number} assignmentId - ID de la asignación QR
   * @returns {Promise} Item agregado
   */
  async addItemToCart(cartId, assignmentId) {
    try {
      const response = await this.api.post(`/api/carts/${cartId}/items`, {
        assignment_id: assignmentId
      })
      return response.data
    } catch (error) {
      console.error('Error al agregar item al carrito:', error)
      throw error
    }
  }

  /**
   * Remueve un item del carrito
   * @param {number} cartId - ID del carrito
   * @param {number} itemId - ID del item
   * @returns {Promise} Resultado de la operación
   */
  async removeItemFromCart(cartId, itemId) {
    try {
      const response = await this.api.delete(`/api/carts/${cartId}/items/${itemId}`)
      return response.data
    } catch (error) {
      console.error('Error al remover item del carrito:', error)
      throw error
    }
  }

  /**
   * Cierra un carrito
   * @param {number} cartId - ID del carrito
   * @returns {Promise} Respuesta del servidor
   */
  async closeCart(cartId) {
    try {
      const response = await this.api.post(`/api/carts/${cartId}/close`)
      return response.data
    } catch (error) {
      console.error('Error al cerrar carrito:', error)
      throw error
    }
  }

  /**
   * Marca un item del carrito como utilizado
   * @param {number} cartId - ID del carrito
   * @param {number} itemId - ID del item
   * @returns {Promise} Respuesta del servidor
   */
  async markItemAsUsed(cartId, itemId) {
    try {
      const response = await this.api.post(`/api/carts/${cartId}/items/${itemId}/use`)
      return response.data
    } catch (error) {
      console.error('Error al marcar item como utilizado:', error)
      throw error
    }
  }

  /**
   * Marca un item del carrito para devolución
   * @param {number} cartId - ID del carrito
   * @param {number} itemId - ID del item
   * @param {string} reason - Motivo de la devolución
   * @returns {Promise} Respuesta del servidor
   */
  async markItemForReturn(cartId, itemId, reason = '') {
    try {
      const response = await this.api.post(`/api/carts/${cartId}/items/${itemId}/return`, {
        reason: reason || 'Sin especificar'
      })
      return response.data
    } catch (error) {
      console.error('Error al marcar item para devolución:', error)
      throw error
    }
  }

  /**
   * Procesa múltiples items del carrito en una sola operación
   * Permite marcar algunos como usados y otros como devueltos en un solo paso
   * @param {number} cartId - ID del carrito
   * @param {Array} items - Array de objetos { item_id, action: 'use'|'return', reason?: string }
   * @returns {Promise} Resultado de la operación múltiple
   */
  async batchOperationItems(cartId, items) {
    try {
      const response = await this.api.post(`/api/carts/${cartId}/items/batch-operation`, {
        items: items.map(item => ({
          item_id: item.itemId || item.item_id,
          action: item.action,
          reason: item.reason || ''
        }))
      })
      return response.data
    } catch (error) {
      console.error('Error en operación múltiple de items:', error)
      throw error
    }
  }

  /**
   * Transfiere todos los items del carrito al pabellón
   * @param {number} cartId - ID del carrito
   * @returns {Promise} Respuesta del servidor
   */
  async transferCartToPavilion(cartId) {
    try {
      const response = await this.api.post(`/api/carts/${cartId}/transfer-to-pavilion`)
      return response.data
    } catch (error) {
      console.error('Error al transferir carrito al pabellón:', error)
      throw error
    }
  }

  // ========================
  // UTILIDADES
  // ========================

  /**
   * Formatea la fecha del carrito
   * @param {string} dateString - Fecha en formato ISO
   * @returns {string} Fecha formateada
   */
  formatDate(dateString) {
    if (!dateString) return '-'
    const date = new Date(dateString)
    return date.toLocaleString('es-CL', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  /**
   * Obtiene el label del estado del carrito
   * @param {string} status - Estado del carrito
   * @returns {string} Label del estado
   */
  getStatusLabel(status) {
    const labels = {
      active: 'Activo',
      closed: 'Cerrado',
      cancelled: 'Cancelado'
    }
    return labels[status] || status
  }

  /**
   * Obtiene la clase CSS para el estado del carrito
   * @param {string} status - Estado del carrito
   * @returns {string} Clase CSS
   */
  getStatusClass(status) {
    const classes = {
      active: 'bg-green-100 text-green-800',
      closed: 'bg-gray-100 text-gray-800',
      cancelled: 'bg-red-100 text-red-800'
    }
    return classes[status] || 'bg-gray-100 text-gray-800'
  }
}

export default new CartService()
