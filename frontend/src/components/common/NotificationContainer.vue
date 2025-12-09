<template>
  <Teleport to="body">
    <div class="fixed top-4 right-4 left-4 sm:left-auto z-50 pointer-events-none max-w-sm sm:w-full">
      <transition-group
        name="notification"
        tag="div"
        class="space-y-3"
      >
        <div
          v-for="notification in notifications"
          :key="notification.id"
          class="notification-toast pointer-events-auto rounded-2xl p-4 flex items-start space-x-3 transition-all duration-300"
          :class="getNotificationClass(notification.type)"
        >
          <!-- Icono con fondo circular -->
          <div class="flex-shrink-0">
            <div class="icon-circle" :class="getIconCircleClass(notification.type)">
              <svg v-if="notification.type === 'success'" 
                class="h-5 w-5" 
                :class="getIconColorClass(notification.type)"
                fill="none" 
                stroke="currentColor" 
                viewBox="0 0 24 24"
                stroke-width="2.5"
              >
                <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
              </svg>
              <svg v-else-if="notification.type === 'error'" 
                class="h-5 w-5" 
                :class="getIconColorClass(notification.type)"
                fill="none" 
                stroke="currentColor" 
                viewBox="0 0 24 24"
                stroke-width="2.5"
              >
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
              </svg>
              <svg v-else-if="notification.type === 'warning'" 
                class="h-5 w-5" 
                :class="getIconColorClass(notification.type)"
                fill="none" 
                stroke="currentColor" 
                viewBox="0 0 24 24"
                stroke-width="2.5"
              >
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01" />
              </svg>
              <svg v-else 
                class="h-5 w-5" 
                :class="getIconColorClass(notification.type)"
                fill="none" 
                stroke="currentColor" 
                viewBox="0 0 24 24"
                stroke-width="2.5"
              >
                <path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01" />
              </svg>
            </div>
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
            class="flex-shrink-0 rounded-lg p-1 transition-colors duration-200"
            :class="getCloseButtonClass(notification.type)"
            aria-label="Cerrar notificación"
          >
            <svg class="h-4 w-4 sm:h-5 sm:w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
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
    success: 'notification-success',
    error: 'notification-error',
    warning: 'notification-warning',
    info: 'notification-info'
  }
  return classes[type] || classes.info
}

const getIconCircleClass = (type) => {
  const classes = {
    success: 'bg-green-500',
    error: 'bg-red-500',
    warning: 'bg-yellow-500',
    info: 'bg-blue-500'
  }
  return classes[type] || classes.info
}

const getIconColorClass = (type) => {
  return 'text-white'
}

const getCloseButtonClass = (type) => {
  const classes = {
    success: 'text-green-600 hover:text-green-800 hover:bg-green-50',
    error: 'text-red-600 hover:text-red-800 hover:bg-red-50',
    warning: 'text-yellow-600 hover:text-yellow-800 hover:bg-yellow-50',
    info: 'text-blue-600 hover:text-blue-800 hover:bg-blue-50'
  }
  return classes[type] || classes.info
}
</script>

<style scoped>
/* Animaciones de entrada y salida */
.notification-enter-active {
  animation: slideIn 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.notification-leave-active {
  animation: slideOut 0.3s ease-in;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateX(100%) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateX(0) scale(1);
  }
}

@keyframes slideOut {
  from {
    opacity: 1;
    transform: translateX(0) scale(1);
  }
  to {
    opacity: 0;
    transform: translateX(100%) scale(0.95);
  }
}

.notification-move {
  transition: transform 0.4s ease;
}

/* Estilos del toast */
.notification-toast {
  background: white;
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.1), 
              0 8px 10px -6px rgba(0, 0, 0, 0.1),
              0 0 0 1px rgba(0, 0, 0, 0.05);
  position: relative;
  overflow: hidden;
}

.notification-toast::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: currentColor;
}

.notification-success {
  color: #10b981;
}

.notification-error {
  color: #ef4444;
}

.notification-warning {
  color: #f59e0b;
}

.notification-info {
  color: #3b82f6;
}

.notification-toast:hover {
  box-shadow: 0 20px 35px -5px rgba(0, 0, 0, 0.15), 
              0 10px 15px -6px rgba(0, 0, 0, 0.1),
              0 0 0 1px rgba(0, 0, 0, 0.05);
  transform: translateY(-2px);
}

/* Círculo del icono */
.icon-circle {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  animation: iconPulse 0.6s ease-out;
}

@keyframes iconPulse {
  0% {
    transform: scale(0);
    opacity: 0;
  }
  50% {
    transform: scale(1.1);
  }
  100% {
    transform: scale(1);
    opacity: 1;
  }
}

/* Responsive */
@media (max-width: 640px) {
  .notification-toast {
    margin: 0 0.5rem;
  }
}
</style>

