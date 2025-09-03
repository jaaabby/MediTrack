import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  // Ruta principal
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home.vue'),
    meta: {
      title: 'Inicio - MediTrack',
      description: 'Sistema de trazabilidad para dispositivos médicos',
      requiresAuth: false
    }
  },
  
  // Rutas de inventario
  {
    path: '/inventory',
    name: 'Inventory',
    component: () => import('@/views/Inventory.vue'),
    meta: {
      title: 'Inventario - MediTrack',
      description: 'Gestión de inventario médico',
      requiresAuth: false
    }
  },
  {
    path: '/inventory/add',
    name: 'AddSupply',
    component: () => import('@/views/AddSupply.vue'),
    meta: {
      title: 'Agregar Insumo - MediTrack',
      description: 'Agregar nuevo insumo médico',
      requiresAuth: false
    }
  },

  // Rutas de códigos QR
  {
    path: '/qr',
    name: 'QRScanner',
    component: () => import('@/views/QRScanner.vue'),
    meta: {
      title: 'Escáner QR - MediTrack',
      description: 'Escanear códigos QR de insumos médicos',
      requiresAuth: false
    }
  },

  // ========================================
  // NUEVAS RUTAS DE SOLICITUDES DE INSUMO
  // ========================================
  
  {
    path: '/supply-requests',
    name: 'SupplyRequestList',
    component: () => import('@/views/SupplyRequestList.vue'),
    meta: {
      title: 'Solicitudes de Insumo - MediTrack',
      description: 'Gestión de solicitudes de insumo con trazabilidad QR',
      requiresAuth: false
    }
  },
  {
    path: '/supply-requests/new',
    name: 'SupplyRequestForm',
    component: () => import('@/views/SupplyRequestForm.vue'),
    meta: {
      title: 'Nueva Solicitud - MediTrack',
      description: 'Crear nueva solicitud de insumo',
      requiresAuth: false
    }
  },
  {
    path: '/supply-requests/:id',
    name: 'SupplyRequestDetail',
    component: () => import('@/views/SupplyRequestDetail.vue'),
    props: route => ({ id: parseInt(route.params.id) }),
    meta: {
      title: 'Detalle de Solicitud - MediTrack',
      description: 'Ver detalles y trazabilidad de solicitud de insumo',
      requiresAuth: false
    }
  },
  {
    path: '/supply-requests/:id/edit',
    name: 'SupplyRequestEdit',
    component: () => import('@/views/SupplyRequestForm.vue'),
    props: route => ({ id: parseInt(route.params.id), editMode: true }),
    meta: {
      title: 'Editar Solicitud - MediTrack',
      description: 'Editar solicitud de insumo existente',
      requiresAuth: false
    }
  },

  // Ruta específica para trazabilidad QR avanzada
  {
    path: '/qr/:qrCode/traceability',
    name: 'QRTraceability',
    component: () => import('@/views/QRTraceability.vue'),
    props: route => ({ qrCode: route.params.qrCode }),
    meta: {
      title: 'Trazabilidad QR - MediTrack',
      description: 'Trazabilidad completa del código QR',
      requiresAuth: false
    }
  },

  // Redirecciones y rutas de error
  {
    path: '/404',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue'),
    meta: {
      title: 'Página no encontrada - MediTrack',
      description: 'La página solicitada no existe',
      requiresAuth: false
    }
  },
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
  
  // Ocultar páginas de auth si ya está autenticado
  // if (to.meta.hideForAuth && isAuthenticated()) {
  //   next({ name: 'Home' })
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
//     const currentTime = Date.now() / 1000
//     return payload.exp < currentTime
//   } catch {
//     return true
//   }
// }

export default router