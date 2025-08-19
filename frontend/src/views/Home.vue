<template>
  <div style="padding: 24px; background-color: #f9fafb; min-height: 100vh;">
    <!-- Header con saludo -->
    <div style="background: linear-gradient(135deg, #2563eb 0%, #1d4ed8 100%); border-radius: 12px; padding: 24px; color: white; margin-bottom: 24px;">
      <div style="display: flex; align-items: center; justify-content: space-between;">
        <div>
          <h1 style="font-size: 24px; font-weight: bold; margin: 0;">¡Hola INNOVO!</h1>
          <p style="color: #dbeafe; margin: 4px 0 0 0;">Sistema de gestión de inventario médico</p>
        </div>
        <div style="text-align: right;">
          <p style="font-size: 14px; color: #dbeafe; margin: 0;">Fecha actual</p>
          <p style="font-size: 18px; font-weight: 600; margin: 0;">{{ currentDate }}</p>
        </div>
      </div>
    </div>

    <!-- Barra de búsqueda principal -->
    <div style="background: white; border-radius: 8px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); border: 1px solid #e5e7eb; padding: 24px; margin-bottom: 24px;">
      <div style="display: flex; align-items: center; gap: 16px;">
        <div style="flex: 1;">
          <label for="search" style="display: block; font-size: 14px; font-weight: 500; color: #374151; margin-bottom: 8px;">
            Buscar insumo médico
          </label>
          <div style="position: relative;">
            <div style="position: absolute; top: 50%; left: 12px; transform: translateY(-50%); pointer-events: none;">
              <svg style="width: 20px; height: 20px; color: #9ca3af;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            <input
              id="search"
              type="text"
              placeholder="Ingrese nombre, código o lote del insumo..."
              style="width: 100%; padding: 8px 12px 8px 40px; border: 1px solid #d1d5db; border-radius: 8px; font-size: 14px;"
              v-model="searchQuery"
              @input="handleSearch"
            />
          </div>
        </div>
        <div style="display: flex; gap: 8px;">
          <button style="background: #2563eb; color: white; padding: 8px 16px; border: none; border-radius: 8px; display: flex; align-items: center; cursor: pointer;">
            <svg style="width: 16px; height: 16px; margin-right: 8px;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
            Buscar
          </button>
          <button style="background: #f3f4f6; color: #374151; padding: 8px 16px; border: none; border-radius: 8px; display: flex; align-items: center; cursor: pointer;">
            <svg style="width: 16px; height: 16px; margin-right: 8px;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.207A1 1 0 013 6.5V4z" />
            </svg>
            Filtros
          </button>
        </div>
      </div>
    </div>

    <!-- Funcionalidades principales - Solo Inventario -->
    <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 24px; margin-bottom: 24px;">
      <!-- Gestión de inventario -->
      <div style="background: white; border-radius: 8px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); border: 1px solid #e5e7eb; padding: 24px; cursor: pointer;" @click="navigateTo('/inventory')">
        <div style="display: flex; align-items: center; gap: 16px;">
          <div style="background: #dbeafe; padding: 12px; border-radius: 8px;">
            <svg style="width: 32px; height: 32px; color: #2563eb;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
            </svg>
          </div>
          <div>
            <h3 style="font-size: 18px; font-weight: 600; color: #111827; margin: 0;">Gestión de Inventario</h3>
            <p style="color: #6b7280; font-size: 14px; margin: 4px 0 0 0;">Administrar insumos médicos</p>
          </div>
        </div>
      </div>

      <!-- Agregar nuevo insumo -->
      <div style="background: white; border-radius: 8px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); border: 1px solid #e5e7eb; padding: 24px; cursor: pointer;" @click="navigateTo('/inventory/add')">
        <div style="display: flex; align-items: center; gap: 16px;">
          <div style="background: #dcfce7; padding: 12px; border-radius: 8px;">
            <svg style="width: 32px; height: 32px; color: #16a34a;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
          </div>
          <div>
            <h3 style="font-size: 18px; font-weight: 600; color: #111827; margin: 0;">Agregar Insumo</h3>
            <p style="color: #6b7280; font-size: 14px; margin: 4px 0 0 0;">Registrar nuevo producto médico</p>
          </div>
        </div>
      </div>

      <!-- Escanear QR -->
      <div style="background: white; border-radius: 8px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); border: 1px solid #e5e7eb; padding: 24px; cursor: pointer;" @click="navigateTo('/qr')">
        <div style="display: flex; align-items: center; gap: 16px;">
          <div style="background: #fef3c7; padding: 12px; border-radius: 8px;">
            <svg style="width: 32px; height: 32px; color: #d97706;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v2a2 2 0 002 2zm0 0h2a2 2 0 012 2v2a2 2 0 01-2 2H5a2 2 0 01-2-2v-2a2 2 0 012-2z" />
            </svg>
          </div>
          <div>
            <h3 style="font-size: 18px; font-weight: 600; color: #111827; margin: 0;">Escáner QR</h3>
            <p style="color: #6b7280; font-size: 14px; margin: 4px 0 0 0;">Trazabilidad y verificación</p>
          </div>
        </div>
      </div>

      <!-- Reportes de inventario -->
      <div style="background: white; border-radius: 8px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); border: 1px solid #e5e7eb; padding: 24px; cursor: pointer;" @click="navigateTo('/inventory/reports')">
        <div style="display: flex; align-items: center; gap: 16px;">
          <div style="background: #dbeafe; padding: 12px; border-radius: 8px;">
            <svg style="width: 32px; height: 32px; color: #2563eb;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
          </div>
          <div>
            <h3 style="font-size: 18px; font-weight: 600; color: #111827; margin: 0;">Reportes</h3>
            <p style="color: #6b7280; font-size: 14px; margin: 4px 0 0 0;">Análisis y estadísticas</p>
          </div>
        </div>
      </div>

      <!-- Movimientos de stock -->
      <div style="background: white; border-radius: 8px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); border: 1px solid #e5e7eb; padding: 24px; cursor: pointer;" @click="navigateTo('/inventory/movements')">
        <div style="display: flex; align-items: center; gap: 16px;">
          <div style="background: #f3e8ff; padding: 12px; border-radius: 8px;">
            <svg style="width: 32px; height: 32px; color: #9333ea;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
          </div>
          <div>
            <h3 style="font-size: 18px; font-weight: 600; color: #111827; margin: 0;">Movimientos</h3>
            <p style="color: #6b7280; font-size: 14px; margin: 4px 0 0 0;">Entradas y salidas de stock</p>
          </div>
        </div>
      </div>

      <!-- Configuración -->
      <div style="background: white; border-radius: 8px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); border: 1px solid #e5e7eb; padding: 24px; cursor: pointer;" @click="navigateTo('/inventory/settings')">
        <div style="display: flex; align-items: center; gap: 16px;">
          <div style="background: #f3f4f6; padding: 12px; border-radius: 8px;">
            <svg style="width: 32px; height: 32px; color: #6b7280;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
          </div>
          <div>
            <h3 style="font-size: 18px; font-weight: 600; color: #111827; margin: 0;">Configuración</h3>
            <p style="color: #6b7280; font-size: 14px; margin: 4px 0 0 0;">Ajustes del sistema</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Acceso rápido al inventario -->
    <div style="background: white; border-radius: 8px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); border: 1px solid #e5e7eb; padding: 24px;">
      <div style="text-align: center;">
        <h3 style="font-size: 18px; font-weight: 600; color: #111827; margin: 0 0 16px 0;">Acceso rápido al inventario</h3>
        <div style="display: flex; align-items: center; justify-content: center; gap: 16px;">
          <button style="padding: 8px; border-radius: 8px; background: #f3f4f6; border: none; cursor: pointer;">
            <svg style="width: 24px; height: 24px; color: #6b7280;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 12H4" />
            </svg>
          </button>
          <div style="background: #dbeafe; padding: 24px; border-radius: 8px;">
            <svg style="width: 64px; height: 64px; color: #2563eb;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
            </svg>
          </div>
          <button style="padding: 8px; border-radius: 8px; background: #f3f4f6; border: none; cursor: pointer;">
            <svg style="width: 24px; height: 24px; color: #6b7280;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
          </button>
        </div>
        <p style="color: #6b7280; margin: 8px 0 16px 0;">Toca para acceder al inventario completo</p>
        <button 
          @click="navigateTo('/inventory')"
          style="background: #2563eb; color: white; padding: 8px 24px; border: none; border-radius: 8px; cursor: pointer; font-weight: 500;"
        >
          Ir al Inventario
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const searchQuery = ref('')

const currentDate = computed(() => {
  return new Date().toLocaleDateString('es-ES', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
})

const handleSearch = () => {
  // TODO: Implementar lógica de búsqueda en inventario
  console.log('Buscando en inventario:', searchQuery.value)
}

const navigateTo = (path) => {
  router.push(path)
}
</script>