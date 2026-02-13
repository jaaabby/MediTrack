import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useNotificationStore = defineStore('notification', () => {
  const notifications = ref([])
  let notificationId = 0

  /**
   * Muestra una notificación
   * @param {string} message - Mensaje a mostrar
   * @param {string} type - Tipo: 'success', 'error', 'warning', 'info'
   * @param {number} duration - Duración en milisegundos (default: 5000)
   */
  const show = (message, type = 'info', duration = 5000) => {
    const id = notificationId++
    const notification = {
      id,
      message,
      type,
      timestamp: Date.now()
    }
    
    notifications.value.push(notification)
    
    // Auto-remover después de la duración especificada
    if (duration > 0) {
      setTimeout(() => {
        remove(id)
      }, duration)
    }
    
    return id
  }

  /**
   * Muestra una notificación de éxito
   */
  const success = (message, duration = 5000) => {
    return show(message, 'success', duration)
  }

  /**
   * Muestra una notificación de error
   */
  const error = (message, duration = 6000) => {
    return show(message, 'error', duration)
  }

  /**
   * Muestra una notificación de advertencia
   */
  const warning = (message, duration = 5000) => {
    return show(message, 'warning', duration)
  }

  /**
   * Muestra una notificación informativa
   */
  const info = (message, duration = 5000) => {
    return show(message, 'info', duration)
  }

  /**
   * Remueve una notificación por ID
   */
  const remove = (id) => {
    const index = notifications.value.findIndex(n => n.id === id)
    if (index > -1) {
      notifications.value.splice(index, 1)
    }
  }

  /**
   * Limpia todas las notificaciones
   */
  const clear = () => {
    notifications.value = []
  }

  return {
    notifications,
    show,
    success,
    error,
    warning,
    info,
    remove,
    clear
  }
})

