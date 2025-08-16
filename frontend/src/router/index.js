import { createRouter, createWebHistory } from 'vue-router'
import Home from '@/views/Home.vue'
import Inventory from '@/views/Inventory.vue'
//import Statistics from '@/views/Statistics.vue'
//import Movements from '@/views/Movements.vue'
//import Profile from '@/views/Profile.vue'
//import Tracking from '@/views/Tracking.vue'
//import Alerts from '@/views/Alerts.vue'

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
  /*{
    path: '/statistics',
    name: 'Statistics',
    component: Statistics,
    meta: {
      title: 'Estadísticas - MediTrack',
      requiresAuth: false
    }
  },
  {
    path: '/movements',
    name: 'Movements',
    component: Movements,
    meta: {
      title: 'Últimos Movimientos - MediTrack',
      requiresAuth: false
    }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: Profile,
    meta: {
      title: 'Perfil - MediTrack',
      requiresAuth: false
    }
  },
  {
    path: '/tracking',
    name: 'Tracking',
    component: Tracking,
    meta: {
      title: 'Trazabilidad - MediTrack',
      requiresAuth: false
    }
  },
  {
    path: '/alerts',
    name: 'Alerts',
    component: Alerts,
    meta: {
      title: 'Alertas - MediTrack',
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
  }*/
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