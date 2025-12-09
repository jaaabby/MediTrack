<template>
  <Teleport to="body">
    <div class="fixed top-4 right-4 left-4 sm:left-auto z-50 pointer-events-none max-w-sm sm:w-full">
      <transition-group
        name="notification"
        tag="div"
        class="space-y-2"
      >
        <div
          v-for="notification in notifications"
          :key="notification.id"
          class="notification-toast pointer-events-auto shadow-lg rounded-lg p-4 flex items-start space-x-3 transition-all duration-300"
          :class="getNotificationClass(notification.type)"
        >
          <!-- Icono -->
          <div class="flex-shrink-0">
            <svg v-if="notification.type === 'success'" 
              class="h-5 w-5 sm:h-6 sm:w-6" 
              fill="none" 
              stroke="currentColor" 
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <svg v-else-if="notification.type === 'error'" 
              class="h-5 w-5 sm:h-6 sm:w-6" 
              fill="none" 
              stroke="currentColor" 
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <svg v-else-if="notification.type === 'warning'" 
              class="h-5 w-5 sm:h-6 sm:w-6" 
              fill="none" 
              stroke="currentColor" 
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            <svg v-else 
              class="h-5 w-5 sm:h-6 sm:w-6" 
              fill="none" 
              stroke="currentColor" 
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          
          <!-- Mensaje -->
          <div class="flex-1 min-w-0">
            <p class="text-sm sm:text-base font-medium break-words">
              {{ notification.message }}
            </p>
          </div>
          
          <!-- Botón cerrar -->
          <button 
            @click="removeNotification(notification.id)" 
            class="flex-shrink-0 opacity-70 hover:opacity-100 transition-opacity"
            :class="getCloseButtonClass(notification.type)"
            aria-label="Cerrar notificación"
          >
            <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </transition-group>
    </div>
  </Teleport>
</template>

<script setup>
import { computed } from 'vue'
import { useNotificationStore } from '@/stores/notification'

const notificationStore = useNotificationStore()

const notifications = computed(() => notificationStore.notifications)

const removeNotification = (id) => {
  notificationStore.remove(id)
}

const getNotificationClass = (type) => {
  const classes = {
    success: 'bg-green-50 text-green-800 border border-green-200',
    error: 'bg-red-50 text-red-800 border border-red-200',
    warning: 'bg-yellow-50 text-yellow-800 border border-yellow-200',
    info: 'bg-blue-50 text-blue-800 border border-blue-200'
  }
  return classes[type] || classes.info
}

const getCloseButtonClass = (type) => {
  const classes = {
    success: 'text-green-600 hover:text-green-800',
    error: 'text-red-600 hover:text-red-800',
    warning: 'text-yellow-600 hover:text-yellow-800',
    info: 'text-blue-600 hover:text-blue-800'
  }
  return classes[type] || classes.info
}
</script>

<style scoped>
/* Animaciones de entrada y salida */
.notification-enter-active {
  transition: all 0.3s ease-out;
}

.notification-leave-active {
  transition: all 0.3s ease-in;
}

.notification-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.notification-leave-to {
  opacity: 0;
  transform: translateX(100%);
}

.notification-move {
  transition: transform 0.3s ease;
}

/* Estilos del toast */
.notification-toast {
  backdrop-filter: blur(10px);
}
</style>

