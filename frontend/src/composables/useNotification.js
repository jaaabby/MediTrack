import { useNotificationStore } from '@/stores/notification'

/**
 * Composable para usar el sistema de notificaciones
 * @returns {Object} Métodos para mostrar notificaciones
 */
export function useNotification() {
  const store = useNotificationStore()

  return {
    /**
     * Muestra una notificación genérica
     * @param {string} message - Mensaje a mostrar
     * @param {string} type - Tipo: 'success', 'error', 'warning', 'info'
     * @param {number} duration - Duración en milisegundos
     */
    show: (message, type = 'info', duration = 5000) => {
      return store.show(message, type, duration)
    },
    
    /**
     * Muestra una notificación de éxito
     */
    success: (message, duration = 5000) => {
      return store.success(message, duration)
    },
    
    /**
     * Muestra una notificación de error
     */
    error: (message, duration = 6000) => {
      return store.error(message, duration)
    },
    
    /**
     * Muestra una notificación de advertencia
     */
    warning: (message, duration = 5000) => {
      return store.warning(message, duration)
    },
    
    /**
     * Muestra una notificación informativa
     */
    info: (message, duration = 5000) => {
      return store.info(message, duration)
    },
    
    /**
     * Remueve una notificación por ID
     */
    remove: (id) => {
      store.remove(id)
    },
    
    /**
     * Limpia todas las notificaciones
     */
    clear: () => {
      store.clear()
    }
  }
}

