import { createRouter, createWebHistory } from 'vue-router'
import Home from '@/views/Home.vue'
import Inventory from '@/views/Inventory.vue'
import QRScanner from '@/views/QRScanner.vue'
import QRDetails from '@/views/QRDetails.vue'
import QRConsumer from '@/views/QRConsumer.vue'
import AddSupply from '@/views/AddSupply.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home,
    meta: {
      title: 'Inicio - MediTrack',
      description: 'Sistema de trazabilidad para dispositivos médicos',
      requiresAuth: false
    }
  },
  {
    path: '/inventory',
    name: 'Inventory',
    component: Inventory,
    meta: {
      title: 'Inventario - MediTrack',
      description: 'Gestión completa del inventario médico',
      requiresAuth: false
    }
  },
  {
    path: '/add-supply',
    name: 'AddSupply',
    component: AddSupply,
    meta: {
      title: 'Agregar Insumo - MediTrack',
      description: 'Crear nuevos lotes con códigos QR únicos',
      requiresAuth: false
    }
  },
  {
    path: '/qr',
    name: 'QRScanner',
    component: QRScanner,
    meta: {
      title: 'Escáner QR - MediTrack',
      description: 'Escanear códigos QR de productos y lotes',
      requiresAuth: false
    }
  },
  {
    path: '/qr/:qrcode',
    name: 'QRDetails',
    component: QRDetails,
    props: true,
    meta: {
      title: 'Detalles QR - MediTrack',
      description: 'Información detallada del código QR escaneado',
      requiresAuth: false
    }
  },
  {
    path: '/consume',
    name: 'QRConsumer',
    component: QRConsumer,
    meta: {
      title: 'Consumir Productos - MediTrack',
      description: 'Registrar consumo de insumos médicos',
      requiresAuth: false
    }
  },
  // Rutas de redirección para compatibilidad
  {
    path: '/scanner',
    redirect: '/qr'
  },
  {
    path: '/scan',
    redirect: '/qr'
  },
  {
    path: '/qr-scanner',
    redirect: '/qr'
  },
  {
    path: '/consumption',
    redirect: '/consume'
  },
  {
    path: '/consume-supply',
    redirect: '/consume'
  },
  
  // Ruta 404
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue'),
    meta: {
      title: 'Página no encontrada - MediTrack',
      description: 'La página solicitada no existe',
      requiresAuth: false
    }
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
  scrollBehavior(to, from, savedPosition) {
    // Si hay una posición guardada (navegación con botón atrás/adelante)
    if (savedPosition) {
      return savedPosition
    }
    
    // Si hay un hash en la URL, ir a ese elemento
    if (to.hash) {
      return {
        el: to.hash,
        behavior: 'smooth'
      }
    }
    
    // Por defecto, ir al inicio de la página
    return { 
      top: 0,
      behavior: 'smooth'
    }
  }
})

// Guard de navegación global - Antes de cada ruta
router.beforeEach((to, from, next) => {
  // Actualizar título de la página
  if (to.meta.title) {
    document.title = to.meta.title
  }
  
  // Actualizar meta description
  if (to.meta.description) {
    const metaDescription = document.querySelector('meta[name="description"]')
    if (metaDescription) {
      metaDescription.setAttribute('content', to.meta.description)
    } else {
      // Crear meta description si no existe
      const meta = document.createElement('meta')
      meta.name = 'description'
      meta.content = to.meta.description
      document.head.appendChild(meta)
    }
  }
  
  // Verificar autenticación cuando esté implementada
  // if (to.meta.requiresAuth && !isAuthenticated()) {
  //   next({ name: 'Login', query: { redirect: to.fullPath } })
  //   return
  // }
  
  // Logging para desarrollo (remover en producción)
  if (import.meta.env.DEV) {
    console.log(`Navegando de ${from.fullPath} a ${to.fullPath}`)
  }
  
  next()
})

// Guard de navegación global - Después de cada ruta
router.afterEach((to, from) => {
  // Analytics o tracking aquí cuando sea necesario
  // gtag('config', 'GA_MEASUREMENT_ID', {
  //   page_path: to.fullPath
  // })
  
  // Limpieza de estados globales si es necesario
  // Ej: cerrar modales, limpiar notificaciones temporales, etc.
})

// Manejo de errores de navegación
router.onError((error) => {
  console.error('Error de navegación:', error)
  
  // En producción, podrías enviar este error a un servicio de logging
  // logError('Navigation Error', error)
})

// Funciones auxiliares para guards (implementar cuando sea necesario)

// Verificar si el usuario está autenticado
// function isAuthenticated() {
//   const token = localStorage.getItem('auth_token')
//   return token && !isTokenExpired(token)
// }

// Verificar si el token ha expirado
// function isTokenExpired(token) {
//   try {
//     const payload = JSON.parse(atob(token.split('.')[1]))
//     return payload.exp * 1000 < Date.now()
//   } catch (error) {
//     return true
//   }
// }

// Verificar permisos del usuario
// function hasPermission(permission) {
//   const userPermissions = JSON.parse(localStorage.getItem('user_permissions') || '[]')
//   return userPermissions.includes(permission)
// }

// Funciones para validación de rutas QR
export function validateQRRoute(qrcode) {
  // Validar formato de código QR
  const qrPattern = /^(BATCH|SUPPLY)_\d+_[a-f0-9]+$/i
  return qrPattern.test(qrcode)
}

// Función para generar rutas programáticas
export function generateQRRoute(qrcode) {
  return {
    name: 'QRDetails',
    params: { qrcode }
  }
}

export function generateConsumeRoute(qrcode = null) {
  const route = { name: 'QRConsumer' }
  if (qrcode) {
    route.query = { qr: qrcode }
  }
  return route
}

// Función para obtener la ruta anterior válida
export function getPreviousRoute(currentRoute, fallback = { name: 'Home' }) {
  const history = router.getRoutes()
  // Implementar lógica para obtener ruta anterior si es necesario
  return fallback
}

export default router