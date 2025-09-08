import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  // Rutas de autenticación
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: {
      title: 'Iniciar Sesión - MediTrack',
      description: 'Acceso al sistema de trazabilidad',
      requiresAuth: false,
      hideForAuth: true
    }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: {
      title: 'Registro - MediTrack',
      description: 'Crear nueva cuenta en el sistema de trazabilidad',
      requiresAuth: false,
      hideForAuth: true
    }
  },

  // Ruta principal
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home.vue'),
    meta: {
      title: 'Inicio - MediTrack',
      description: 'Sistema de trazabilidad para dispositivos médicos',
      requiresAuth: true
    }
  },
  
  // Rutas de inventario
  {
    path: '/inventory',
    name: 'Inventory',
    component: () => import('@/views/Inventory.vue'),
    meta: {
      title: 'Inventario - MediTrack',
      description: 'Gestión completa del inventario médico',
      requiresAuth: true
    }
  },
  {
    path: '/inventory/add',
    name: 'AddSupply',
    component: () => import('@/views/AddSupply.vue'),
    meta: {
      title: 'Agregar Insumo - MediTrack',
      description: 'Crear nuevos lotes con códigos QR únicos',
      requiresAuth: true
    }
  },

  // Rutas de códigos QR
  {
    path: '/qr',
    name: 'QRScanner',
    component: () => import('@/views/QRScanner.vue'),
    meta: {
      title: 'Escáner QR - MediTrack',
      description: 'Escanear códigos QR de productos y lotes',
      requiresAuth: true
    }
  },
  {
    path: '/qr/consumer',
    name: 'QRConsumer',
    component: () => import('@/views/QRConsumer.vue'),
    meta: {
      title: 'Consumo QR - MediTrack',
      description: 'Consumir insumos médicos mediante códigos QR',
      requiresAuth: true
    }
  },

  // Rutas específicas de QR con historial y trazabilidad
  {
    path: '/qr/:qrCode/details',
    name: 'QRDetails',
    component: () => import('@/views/QRDetails.vue'),
    props: route => ({ qrCode: route.params.qrCode }),
    meta: {
      title: 'Detalles QR - MediTrack',
      description: 'Información detallada del código QR',
      requiresAuth: true
    }
  },
  {
    path: '/qr/:qrCode/traceability',
    name: 'QRTraceability',
    component: () => import('@/views/QRTraceability.vue'),
    props: route => ({ qrCode: route.params.qrCode }),
    meta: {
      title: 'Trazabilidad QR - MediTrack',
      description: 'Trazabilidad completa del código QR',
      requiresAuth: true
    }
  },

  // Ruta para historial de lotes
  {
    path: '/batch/:batchId/history',
    name: 'BatchHistory',
    component: () => import('@/components/BatchHistory.vue'),
    props: route => ({ batchId: parseInt(route.params.batchId) }),
    meta: {
      title: 'Historial del Lote - MediTrack',
      description: 'Historial completo de movimientos del lote',
      requiresAuth: true
    }
  },

  // Rutas de solicitudes de insumo
  {
    path: '/supply-requests',
    name: 'SupplyRequestList',
    component: () => import('@/views/SupplyRequestList.vue'),
    meta: {
      title: 'Solicitudes de Insumo - MediTrack',
      description: 'Gestión de solicitudes de insumo con trazabilidad QR',
      requiresAuth: true
    }
  },
  {
    path: '/supply-requests/new',
    name: 'SupplyRequestForm',
    component: () => import('@/views/SupplyRequestForm.vue'),
    meta: {
      title: 'Nueva Solicitud - MediTrack',
      description: 'Crear nueva solicitud de insumo',
      requiresAuth: true
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
      requiresAuth: true
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
      requiresAuth: true
    }
  },

  // Ruta de perfil de usuario
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/Profile.vue'),
    meta: {
      title: 'Mi Perfil - MediTrack',
      description: 'Información personal y configuración de cuenta',
      requiresAuth: true
    }
  },

  // Redirecciones inteligentes para QR codes
  {
    path: '/qr/:qrCode',
    name: 'QRRedirect',
    redirect: to => {
      const qrCode = to.params.qrCode
      
      // Todos los QR codes van a detalles por ahora
      return `/qr/${qrCode}/details`
    },
    meta: {
      title: 'Redirigiendo QR - MediTrack',
      description: 'Procesando código QR',
      requiresAuth: true
    }
  },

  // Rutas de error
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
    name: 'CatchAll',
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
    if (savedPosition) {
      return savedPosition
    }
    
    if (to.hash) {
      return {
        el: to.hash,
        behavior: 'smooth'
      }
    }
    
    return { 
      top: 0,
      behavior: 'smooth'
    }
  }
})

// Guard de navegación global - Antes de cada ruta
router.beforeEach(async (to, from, next) => {
  if (to.meta.title) {
    document.title = to.meta.title
  }
  
  if (to.meta.description) {
    const metaDescription = document.querySelector('meta[name="description"]')
    if (metaDescription) {
      metaDescription.setAttribute('content', to.meta.description)
    } else {
      const meta = document.createElement('meta')
      meta.name = 'description'
      meta.content = to.meta.description
      document.head.appendChild(meta)
    }
  }
  
  const { useAuthStore } = await import('@/stores/auth')
  const authStore = useAuthStore()
  
  if (!authStore.isAuthenticated && authStore.token === null) {
    authStore.initializeAuth()
  }
  
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
    return
  }
  
  if (to.meta.hideForAuth && authStore.isAuthenticated) {
    next({ name: 'Home' })
    return
  }
  
  if (to.meta.requiresAuth && authStore.isAuthenticated && authStore.canAccessRoute && !authStore.canAccessRoute(to.name)) {
    next({ name: 'Home' })
    return
  }
  
  if (to.name === 'BatchHistory' && to.params.batchId) {
    const batchId = parseInt(to.params.batchId)
    if (isNaN(batchId) || batchId <= 0) {
      next({ name: 'NotFound' })
      return
    }
  }
  
  if (import.meta.env.DEV) {
    console.log(`Navegando de ${from.fullPath} a ${to.fullPath}`)
  }
  
  next()
})

router.afterEach((to, from) => {
  const notifications = document.querySelectorAll('[data-notification]')
  notifications.forEach(notification => {
    if (notification.dataset.temporary === 'true') {
      notification.remove()
    }
  })
})

router.onError((error) => {
  console.error('Error de navegación:', error)
  
  if (error.type === 'NavigationDuplicated') {
    return
  }
  
  router.push({ name: 'NotFound' }).catch(() => {
    window.location.reload()
  })
})

export default router