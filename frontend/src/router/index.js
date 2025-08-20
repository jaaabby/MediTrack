import { createRouter, createWebHistory } from 'vue-router'
import Home from '@/views/Home.vue'
import Inventory from '@/views/Inventory.vue'
import QRScanner from '@/views/QRScanner.vue'
import QRDetails from '@/views/QRDetails.vue'
import AddSupply from '@/views/AddSupply.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home,
    meta: {
      title: 'Inicio - MediTrack',
      requiresAuth: false
    }
  },
  {
    path: '/inventory',
    name: 'Inventory',
    component: Inventory,
    meta: {
      title: 'Inventario - MediTrack',
      requiresAuth: false
    }
  },
  {
    path: '/add-supply',
    name: 'AddSupply',
    component: AddSupply,
    meta: {
      title: 'Agregar Insumo - MediTrack',
      requiresAuth: false
    }
  },
  {
    path: '/qr',
    name: 'QRScanner',
    component: QRScanner,
    meta: {
      title: 'Escáner QR - MediTrack',
      requiresAuth: false
    }
  },
  {
    path: '/qr/:qrcode',
    name: 'QRDetails',
    component: QRDetails,
    meta: {
      title: 'Detalles QR - MediTrack',
      requiresAuth: false
    }
  },
  // Ruta 404
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue'),
    meta: {
      title: 'Página no encontrada - MediTrack',
      requiresAuth: false
    }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  }
})

// Guardia de navegación para actualizar el título de la página
router.beforeEach((to, from, next) => {
  // Actualizar título de la página
  if (to.meta.title) {
    document.title = to.meta.title
  }
  
  // TODO: Implementar verificación de autenticación cuando esté lista
  // if (to.meta.requiresAuth && !isAuthenticated()) {
  //   next('/login')
  //   return
  // }
  
  next()
})

export default router