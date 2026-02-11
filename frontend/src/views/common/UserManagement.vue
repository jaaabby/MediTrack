<template>
  <div class="max-w-7xl mx-auto">
    <!-- Header -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-900">Gestión de Usuarios</h1>
      <p class="mt-1 text-sm text-gray-600">
        Administra las cuentas de usuarios del sistema
      </p>
    </div>

    <!-- Filtros y búsqueda -->
    <div class="card mb-6">
      <div class="flex flex-col gap-4">
        <!-- Primera fila: Buscador -->
        <div class="flex flex-col sm:flex-row sm:items-end gap-4">
          <!-- Buscador único -->
          <div class="flex-1">
            <label class="block text-sm font-medium text-gray-700 mb-2">Buscar usuario</label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
              </div>
              <input type="text" placeholder="Buscar por nombre, RUT o email..."
                class="form-input pl-10 w-full" v-model="searchTerm" />
            </div>
          </div>

          <!-- Botón de limpiar búsqueda -->
          <div class="w-full sm:w-auto">
            <button class="btn-secondary px-4 py-2 h-10 w-full sm:w-auto" @click="clearSearch" :disabled="!searchTerm">
              Limpiar
            </button>
          </div>
        </div>

        <!-- Segunda fila: Filtros -->
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <!-- Filtro por Rol -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Filtrar por Rol</label>
            <select v-model="roleFilter" class="form-select w-full">
              <option value="">Todos los roles</option>
              <option value="admin">Administrador</option>
              <option value="encargado de bodega">Encargado de Bodega</option>
              <option value="pabellón">Pabellón</option>
              <option value="pavedad">Pavedad</option>
              <option value="enfermera">Enfermera</option>
              <option value="doctor">Doctor</option>
            </select>
          </div>

          <!-- Filtro por Centro Médico -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Filtrar por Centro Médico</label>
            <select v-model="medicalCenterFilter" class="form-select w-full">
              <option value="">Todos los centros</option>
              <option v-for="center in medicalCenters" :key="center.id" :value="center.id">
                {{ center.name }}
              </option>
            </select>
          </div>

          <!-- Filtro por Estado -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Filtrar por Estado</label>
            <select v-model="statusFilter" class="form-select w-full">
              <option value="">Todos</option>
              <option value="true">Activo</option>
              <option value="false">Inactivo</option>
            </select>
          </div>

          <!-- Botón de crear usuario -->
          <div class="flex items-end">
            <button
              @click="openCreateModal"
              class="btn-primary w-full h-10"
            >
              <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
              </svg>
              <span class="hidden sm:inline">Crear Usuario</span>
              <span class="sm:hidden">Crear</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Lista de usuarios -->
    <div class="card">
      <div class="card-header">
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <h2 class="card-title">Lista de Usuarios</h2>
            <p class="text-sm text-gray-600">Total: {{ filteredUsers.length }} usuario(s)</p>
            <p v-if="hasActiveFilters" class="text-xs text-gray-500 mt-1">
              Mostrando {{ filteredUsers.length }} de {{ users.length }} usuario(s)
            </p>
          </div>
        </div>
      </div>

      <!-- Indicador de carga -->
      <div v-if="loading" class="flex justify-center items-center py-8">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        <span class="ml-2 text-gray-600">Cargando usuarios...</span>
      </div>

      <!-- Tabla de usuarios -->
      <div v-else class="table-container">
        <!-- Indicador de scroll horizontal para móviles -->
        <div class="md:hidden bg-blue-50 border-b border-blue-200 px-3 py-2 text-center sticky left-0 z-10">
          <div class="flex items-center justify-center text-blue-700 text-xs">
            <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 9l4-4 4 4m0 6l-4 4-4-4" />
            </svg>
            Desliza horizontalmente para ver todas las columnas
          </div>
        </div>

        <table class="min-w-full divide-y divide-gray-200" style="min-width: 1200px;">
          <thead class="table-header">
            <tr>
              <th class="table-header-cell">
                <span>Usuario</span>
              </th>
              <th class="table-header-cell">
                <span>Email</span>
              </th>
              <th class="table-header-cell">
                <span>Rol</span>
              </th>
              <th class="table-header-cell">
                <span>Centro Médico</span>
              </th>
              <th class="table-header-cell">
                <span>Pabellón</span>
              </th>
              <th class="table-header-cell">
                <span>Especialidad</span>
              </th>
              <th class="table-header-cell">
                <span>Estado</span>
              </th>
              <th class="table-header-cell sticky right-0 bg-gray-50 z-10 shadow-[-2px_0_4px_rgba(0,0,0,0.05)]">
                <span>Acciones</span>
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-if="filteredUsers.length === 0">
              <td colspan="8" class="px-6 py-8 text-center text-sm text-gray-500">
                <div v-if="hasActiveFilters || searchTerm" class="space-y-2">
                  <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M12 12h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <p>No se encontraron usuarios con los filtros aplicados</p>
                  <button @click="clearAllFilters" class="btn-secondary text-sm">
                    Limpiar filtros
                  </button>
                </div>
                <div v-else>
                  No hay usuarios registrados
                </div>
              </td>
            </tr>
            <tr v-for="user in paginatedUsers" :key="user.id" class="hover:bg-gray-50 transition-colors">
              <td class="table-cell">
                <div class="flex items-center">
                  <div class="flex-shrink-0 h-10 w-10 bg-blue-100 rounded-full flex items-center justify-center">
                    <svg class="h-6 w-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                    </svg>
                  </div>
                  <div class="ml-4">
                    <div class="text-sm font-medium text-gray-900">{{ user.name }}</div>
                    <div class="text-sm text-gray-500">{{ user.rut }}</div>
                  </div>
                </div>
              </td>
              <td class="table-cell">
                <div class="text-sm text-gray-900">{{ user.email }}</div>
              </td>
              <td class="table-cell">
                <span
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                  :class="getRoleBadgeClass(user.role)"
                >
                  {{ user.role }}
                </span>
              </td>
              <td class="table-cell">
                <div class="text-sm text-gray-900">{{ getMedicalCenterName(user.medical_center_id) }}</div>
              </td>
              <td class="table-cell">
                <div class="text-sm text-gray-900">{{ getPavilionName(user.pavilion_id) }}</div>
              </td>
              <td class="table-cell">
                <div class="text-sm text-gray-900">{{ getSpecialtyName(user.specialty_id) }}</div>
              </td>
              <td class="table-cell">
                <span
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                  :class="user.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'"
                >
                  {{ user.is_active ? 'Activo' : 'Inactivo' }}
                </span>
              </td>
              <td class="table-cell sticky right-0 bg-white z-10 shadow-[-2px_0_4px_rgba(0,0,0,0.05)]">
                <div class="flex justify-end space-x-2">
                  <button
                    @click.stop="openEditModal(user)"
                    class="text-warning-600 hover:text-warning-800 hover:bg-warning-50 p-1.5 rounded inline-flex items-center gap-1 transition-colors"
                    title="Editar usuario"
                  >
                    <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                    <span class="font-medium text-xs">Editar</span>
                  </button>
                  <button
                    @click.stop="openDeleteModal(user)"
                    class="text-danger-600 hover:text-danger-800 hover:bg-danger-50 p-1.5 rounded transition-colors"
                    title="Eliminar usuario"
                  >
                    <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Paginación -->
      <div v-if="filteredUsers.length > itemsPerPage"
        class="flex flex-col sm:flex-row items-center justify-between mt-4 gap-3 px-4 pb-4">
        <div class="text-xs sm:text-sm text-gray-700 text-center sm:text-left">
          Mostrando {{ startIndex + 1 }} a {{ endIndex }} de {{ filteredUsers.length }} usuario(s)
        </div>
        <div class="flex items-center gap-2">
          <button class="btn-secondary px-3 py-2 text-sm min-w-[80px]" :disabled="currentPage === 1"
            @click="currentPage--">
            <span class="hidden sm:inline">Anterior</span>
            <span class="sm:hidden">Ant.</span>
          </button>
          <span class="px-3 py-2 text-sm text-gray-700 bg-gray-100 rounded-md min-w-[100px] text-center">
            Página {{ currentPage }} de {{ totalPages }}
          </span>
          <button class="btn-secondary px-3 py-2 text-sm min-w-[80px]"
            :disabled="currentPage === totalPages" @click="currentPage++">
            <span class="hidden sm:inline">Siguiente</span>
            <span class="sm:hidden">Sig.</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Modal para crear/editar usuario -->
    <div
      v-if="showModal"
      class="fixed inset-0 z-50 overflow-y-auto"
      aria-labelledby="modal-title"
      role="dialog"
      aria-modal="true"
    >
      <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <!-- Overlay -->
        <div
          class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity"
          @click="closeModal"
        ></div>

        <!-- Modal panel -->
        <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
            <div class="sm:flex sm:items-start">
              <div class="mt-3 text-center sm:mt-0 sm:text-left w-full">
                <h3 class="text-lg leading-6 font-medium text-gray-900 mb-4" id="modal-title">
                  {{ isEditMode ? 'Editar Usuario' : 'Crear Nuevo Usuario' }}
                </h3>
                
                <form @submit.prevent="saveUser" class="space-y-4">
                  <!-- Nombre -->
                  <div>
                    <label for="name" class="block text-sm font-medium text-gray-700">
                      Nombre completo <span class="text-red-500">*</span>
                    </label>
                    <input
                      id="name"
                      v-model="userForm.name"
                      type="text"
                      required
                      class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                      :class="{ 'border-red-500': formErrors.name }"
                    />
                    <p v-if="formErrors.name" class="mt-1 text-sm text-red-600">{{ formErrors.name }}</p>
                  </div>

                  <!-- RUT -->
                  <div>
                    <label for="rut" class="block text-sm font-medium text-gray-700">
                      RUT <span class="text-red-500">*</span>
                    </label>
                    <input
                      id="rut"
                      v-model="userForm.rut"
                      type="text"
                      required
                      :disabled="isEditMode"
                      class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm disabled:bg-gray-100 disabled:cursor-not-allowed"
                      :class="{ 'border-red-500': formErrors.rut }"
                    />
                    <p v-if="formErrors.rut" class="mt-1 text-sm text-red-600">{{ formErrors.rut }}</p>
                  </div>

                  <!-- Email -->
                  <div>
                    <label for="email" class="block text-sm font-medium text-gray-700">
                      Email <span class="text-red-500">*</span>
                    </label>
                    <input
                      id="email"
                      v-model="userForm.email"
                      type="email"
                      required
                      class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                      :class="{ 'border-red-500': formErrors.email }"
                    />
                    <p v-if="formErrors.email" class="mt-1 text-sm text-red-600">{{ formErrors.email }}</p>
                  </div>

                  <!-- Mensaje informativo sobre contraseña automática (solo al crear) -->
                  <div v-if="!isEditMode" class="bg-blue-50 border border-blue-200 rounded-md p-4">
                    <div class="flex">
                      <svg class="h-5 w-5 text-blue-400 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                      </svg>
                      <div class="text-sm text-blue-700">
                        <p class="font-medium">Contraseña Automática</p>
                        <p class="mt-1">El sistema generará una contraseña segura y la enviará automáticamente al correo electrónico del usuario.</p>
                      </div>
                    </div>
                  </div>

                  <!-- Rol -->
                  <div>
                    <label for="role" class="block text-sm font-medium text-gray-700">
                      Rol <span class="text-red-500">*</span>
                    </label>
                    <select
                      id="role"
                      v-model="userForm.role"
                      required
                      class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                    >
                      <option value="">Seleccione un rol</option>
                      <option value="admin">Administrador</option>
                      <option value="encargado de bodega">Encargado de Bodega</option>
                      <option value="pabellón">Pabellón</option>
                      <option value="pavedad">Pavedad</option>
                      <option value="enfermera">Enfermera</option>
                      <option value="doctor">Doctor</option>
                    </select>
                    <p v-if="formErrors.role" class="mt-1 text-sm text-red-600">{{ formErrors.role }}</p>
                  </div>

                  <!-- Centro Médico -->
                  <div>
                    <label for="medical_center_id" class="block text-sm font-medium text-gray-700">
                      Centro Médico <span class="text-red-500">*</span>
                    </label>
                    <select
                      id="medical_center_id"
                      v-model="userForm.medical_center_id"
                      required
                      class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                    >
                      <option value="">Seleccione un centro médico</option>
                      <option v-for="center in medicalCenters" :key="center.id" :value="center.id">
                        {{ center.name }}
                      </option>
                    </select>
                    <p v-if="formErrors.medical_center_id" class="mt-1 text-sm text-red-600">{{ formErrors.medical_center_id }}</p>
                  </div>

                  <!-- Pabellón (solo para rol pabellón o pavedad) -->
                  <div v-if="userForm.role === 'pabellón' || userForm.role === 'pavedad'">
                    <label for="pavilion_id" class="block text-sm font-medium text-gray-700">
                      Pabellón <span class="text-red-500">*</span>
                    </label>
                    <select
                      id="pavilion_id"
                      v-model="userForm.pavilion_id"
                      :required="userForm.role === 'pabellón' || userForm.role === 'pavedad'"
                      class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                    >
                      <option value="">Seleccione un pabellón</option>
                      <option v-for="pavilion in pavilions" :key="pavilion.id" :value="pavilion.id">
                        {{ pavilion.name }}
                      </option>
                    </select>
                    <p v-if="formErrors.pavilion_id" class="mt-1 text-sm text-red-600">{{ formErrors.pavilion_id }}</p>
                  </div>

                  <!-- Especialidad Médica (solo para rol doctor) -->
                  <div v-if="userForm.role === 'doctor'">
                    <label for="specialty_id" class="block text-sm font-medium text-gray-700">
                      Especialidad Médica <span class="text-red-500">*</span>
                    </label>
                    <select
                      id="specialty_id"
                      v-model="userForm.specialty_id"
                      :required="userForm.role === 'doctor'"
                      class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                    >
                      <option value="">Seleccione una especialidad</option>
                      <option v-for="specialty in medicalSpecialties" :key="specialty.id" :value="specialty.id">
                        {{ specialty.name }}
                      </option>
                    </select>
                    <p v-if="formErrors.specialty_id" class="mt-1 text-sm text-red-600">{{ formErrors.specialty_id }}</p>
                  </div>

                  <!-- Estado (solo para editar) -->
                  <div v-if="isEditMode" class="flex items-center">
                    <input
                      id="is_active"
                      v-model="userForm.is_active"
                      type="checkbox"
                      class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                    />
                    <label for="is_active" class="ml-2 block text-sm text-gray-900">
                      Usuario activo
                    </label>
                  </div>
                </form>
              </div>
            </div>
          </div>
          <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
            <button
              @click="saveUser"
              :disabled="saving"
              type="button"
              class="btn-primary w-full sm:w-auto sm:ml-3 disabled:opacity-50"
            >
              {{ saving ? 'Guardando...' : (isEditMode ? 'Actualizar' : 'Crear Usuario') }}
            </button>
            <button
              @click="closeModal"
              type="button"
              class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 sm:mt-0 sm:w-auto sm:text-sm"
            >
              Cancelar
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- Modal de confirmación de eliminación -->
  <div v-if="showDeleteModal && userToDelete" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50 flex items-center justify-center p-4" @click.self="closeDeleteModal">
    <div class="w-full max-w-md p-5 border shadow-lg rounded-md bg-white">
      <div class="space-y-4">
        <!-- Header -->
        <div class="flex justify-between items-start border-b pb-3">
          <div class="flex items-center">
            <div class="flex-shrink-0 h-10 w-10 rounded-full bg-danger-100 flex items-center justify-center">
              <svg class="h-6 w-6 text-danger-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
            </div>
            <div class="ml-3">
              <h3 class="text-lg font-semibold text-gray-900">Confirmar Eliminación</h3>
              <p class="text-sm text-gray-500">{{ userToDelete.name }}</p>
            </div>
          </div>
          <button @click="closeDeleteModal" class="text-gray-400 hover:text-gray-600">
            <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Contenido -->
        <div class="space-y-4">
          <div class="bg-danger-50 border border-danger-200 rounded-md p-4">
            <div class="flex">
              <svg class="h-5 w-5 text-danger-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
              <div class="ml-3">
                <p class="text-sm text-danger-700">
                  ¿Está seguro de que desea eliminar este usuario?
                </p>
                <div class="mt-2 text-xs text-danger-600">
                  <p><strong>Usuario:</strong> {{ userToDelete.name }}</p>
                  <p><strong>RUT:</strong> {{ userToDelete.rut }}</p>
                  <p><strong>Email:</strong> {{ userToDelete.email }}</p>
                  <p><strong>Rol:</strong> {{ userToDelete.role }}</p>
                </div>
                <p class="mt-2 text-xs text-danger-700 font-semibold">
                  Esta acción no se puede deshacer.
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- Footer -->
        <div class="flex justify-end space-x-3 pt-4 border-t">
          <button @click="closeDeleteModal" class="btn-secondary" :disabled="isDeleting">
            Cancelar
          </button>
          <button @click="confirmDelete" class="btn-danger" :disabled="isDeleting">
            <span v-if="isDeleting">Eliminando...</span>
            <span v-else>Eliminar Usuario</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { userService } from '@/services/common/userService'
import medicalCenterService from '@/services/config/medicalCenterService'
import pavilionService from '@/services/config/pavilionService'
import medicalSpecialtyService from '@/services/config/medicalSpecialtyService'
import { useNotification } from '@/composables/useNotification'

const { success: showSuccess, error: showError } = useNotification()

// Estado
const users = ref([])
const medicalCenters = ref([])
const pavilions = ref([])
const medicalSpecialties = ref([])
const loading = ref(false)
const saving = ref(false)
const showModal = ref(false)
const isEditMode = ref(false)

// Refs para búsqueda y filtros
const searchTerm = ref('')
const roleFilter = ref('')
const medicalCenterFilter = ref('')
const statusFilter = ref('')

// Filtrado de usuarios
const filteredUsers = computed(() => {
  let filtered = users.value

  // Aplicar búsqueda por texto
  if (searchTerm.value) {
    const search = searchTerm.value.toLowerCase()
    filtered = filtered.filter(user => {
      return (
        user.name.toLowerCase().includes(search) ||
        user.rut.toLowerCase().includes(search) ||
        user.email.toLowerCase().includes(search)
      )
    })
  }

  // Aplicar filtro por rol
  if (roleFilter.value) {
    filtered = filtered.filter(user => user.role === roleFilter.value)
  }

  // Aplicar filtro por centro médico
  if (medicalCenterFilter.value) {
    filtered = filtered.filter(user => user.medical_center_id === parseInt(medicalCenterFilter.value))
  }

  // Aplicar filtro por estado
  if (statusFilter.value !== '') {
    const isActive = statusFilter.value === 'true'
    filtered = filtered.filter(user => user.is_active === isActive)
  }

  return filtered
})

// Verificar si hay filtros activos
const hasActiveFilters = computed(() => {
  return searchTerm.value || roleFilter.value || medicalCenterFilter.value || statusFilter.value !== ''
})

// Paginación
const currentPage = ref(1)
const itemsPerPage = ref(10)

// Computed para paginación
const totalPages = computed(() => Math.ceil(filteredUsers.value.length / itemsPerPage.value))
const startIndex = computed(() => (currentPage.value - 1) * itemsPerPage.value)
const endIndex = computed(() => Math.min(startIndex.value + itemsPerPage.value, filteredUsers.value.length))
const paginatedUsers = computed(() => {
  return filteredUsers.value.slice(startIndex.value, endIndex.value)
})

// Formulario de usuario
const userForm = reactive({
  id: null,
  name: '',
  rut: '',
  email: '',
  password: '',
  role: '',
  medical_center_id: '',
  pavilion_id: '',
  specialty_id: '',
  is_active: true
})

const formErrors = reactive({
  name: '',
  rut: '',
  email: '',
  password: '',
  role: '',
  medical_center_id: '',
  pavilion_id: '',
  specialty_id: ''
})

// Watcher para limpiar campos condicionales cuando cambia el rol
watch(() => userForm.role, (newRole, oldRole) => {
  if (newRole !== oldRole) {
    // Limpiar pabellón si no es rol pabellón o pavedad
    if (newRole !== 'pabellón' && newRole !== 'pavedad') {
      userForm.pavilion_id = ''
      formErrors.pavilion_id = ''
    }
    // Limpiar especialidad si no es doctor
    if (newRole !== 'doctor') {
      userForm.specialty_id = ''
      formErrors.specialty_id = ''
    }
  }
})

// Cargar usuarios
const loadUsers = async () => {
  try {
    loading.value = true
    const response = await userService.getAllUsers()
    if (response.success) {
      users.value = response.data || []
    } else {
      showError('Error al cargar usuarios')
    }
  } catch (error) {
    console.error('Error al cargar usuarios:', error)
    showError('Error al cargar usuarios')
  } finally {
    loading.value = false
  }
}

// Cargar centros médicos
const loadMedicalCenters = async () => {
  try {
    const response = await medicalCenterService.getAllMedicalCenters()
    if (response.success) {
      medicalCenters.value = response.data || []
    }
  } catch (error) {
    console.error('Error al cargar centros médicos:', error)
  }
}

// Cargar pabellones
const loadPavilions = async () => {
  try {
    const data = await pavilionService.getAllPavilions()
    pavilions.value = data || []
  } catch (error) {
    console.error('Error al cargar pabellones:', error)
  }
}

// Cargar especialidades médicas
const loadSpecialties = async () => {
  try {
    const data = await medicalSpecialtyService.getAllSpecialties()
    medicalSpecialties.value = data || []
  } catch (error) {
    console.error('Error al cargar especialidades:', error)
  }
}

// Obtener nombre del centro médico
const getMedicalCenterName = (medicalCenterId) => {
  if (!medicalCenterId) return 'Sin asignar'
  const center = medicalCenters.value.find(c => c.id === medicalCenterId)
  return center ? center.name : 'Sin asignar'
}

// Obtener nombre del pabellón
const getPavilionName = (pavilionId) => {
  if (!pavilionId) return 'N/A'
  const pavilion = pavilions.value.find(p => p.id === pavilionId)
  return pavilion ? pavilion.name : 'N/A'
}

// Obtener nombre de la especialidad
const getSpecialtyName = (specialtyId) => {
  if (!specialtyId) return 'N/A'
  const specialty = medicalSpecialties.value.find(s => s.id === specialtyId)
  return specialty ? specialty.name : 'N/A'
}

// Abrir modal para crear
const openCreateModal = () => {
  resetForm()
  isEditMode.value = false
  showModal.value = true
}

// Abrir modal para editar
const openEditModal = (user) => {
  resetForm()
  isEditMode.value = true
  userForm.id = user.id
  userForm.name = user.name
  userForm.rut = user.rut
  userForm.email = user.email
  userForm.role = user.role
  userForm.medical_center_id = user.medical_center_id
  userForm.pavilion_id = user.pavilion_id || ''
  userForm.specialty_id = user.specialty_id || ''
  userForm.is_active = user.is_active
  showModal.value = true
}

// Cerrar modal
const closeModal = () => {
  showModal.value = false
  resetForm()
}
// Resetear formulario
const resetForm = () => {
  userForm.id = null
  userForm.name = ''
  userForm.rut = ''
  userForm.email = ''
  userForm.password = ''
  userForm.role = ''
  userForm.medical_center_id = ''
  userForm.pavilion_id = ''
  userForm.specialty_id = ''
  userForm.is_active = true
  
  // Limpiar errores
  Object.keys(formErrors).forEach(key => {
    formErrors[key] = ''
  })
}

// Validar formulario
const validateForm = () => {
  let isValid = true
  
  // Limpiar errores
  Object.keys(formErrors).forEach(key => {
    formErrors[key] = ''
  })
  
  if (!userForm.name.trim()) {
    formErrors.name = 'El nombre es requerido'
    isValid = false
  }
  
  if (!userForm.rut.trim()) {
    formErrors.rut = 'El RUT es requerido'
    isValid = false
  }
  
  if (!userForm.email.trim()) {
    formErrors.email = 'El email es requerido'
    isValid = false
  } else if (!userForm.email.includes('@')) {
    formErrors.email = 'El email debe ser válido'
    isValid = false
  }
  
  if (!userForm.role) {
    formErrors.role = 'El rol es requerido'
    isValid = false
  }
  
  if (!userForm.medical_center_id) {
    formErrors.medical_center_id = 'El centro médico es requerido'
    isValid = false
  }

  // Validar pabellón para roles pabellón y pavedad
  if ((userForm.role === 'pabellón' || userForm.role === 'pavedad') && !userForm.pavilion_id) {
    formErrors.pavilion_id = 'El pabellón es requerido para este rol'
    isValid = false
  }

  // Validar especialidad para doctores
  if (userForm.role === 'doctor' && !userForm.specialty_id) {
    formErrors.specialty_id = 'La especialidad médica es requerida para doctores'
    isValid = false
  }
  
  return isValid
}

// Guardar usuario
const saveUser = async () => {
  if (!validateForm()) {
    return
  }
  
  try {
    saving.value = true
    let response
    
    // Preparar datos del usuario
    const userData = {
      name: userForm.name,
      rut: userForm.rut,
      email: userForm.email,
      role: userForm.role,
      medical_center_id: userForm.medical_center_id,
      is_active: userForm.is_active
    }
    
    // Incluir pabellón si aplica (convertir a número o null)
    if (userForm.role === 'pabellón' || userForm.role === 'pavedad') {
      userData.pavilion_id = userForm.pavilion_id ? parseInt(userForm.pavilion_id) : null
    }
    
    // Incluir especialidad si es doctor (convertir a número o null)
    if (userForm.role === 'doctor') {
      userData.specialty_id = userForm.specialty_id ? parseInt(userForm.specialty_id) : null
    }
    
    // NO se permite cambiar contraseña desde el formulario de edición
    // Los usuarios deben cambiarla mediante el sistema de recuperación de contraseña
    
    if (isEditMode.value) {
      // Usar RUT como identificador para actualizar
      response = await userService.updateUser(userForm.rut, userData)
    } else {
      response = await userService.createUser(userData)
    }
    
    if (response.success) {
      showSuccess(isEditMode.value ? 'Usuario actualizado correctamente' : 'Usuario creado correctamente. Se ha enviado un correo con la contraseña temporal.')
      closeModal()
      await loadUsers()
    } else {
      showError(response.error || 'Error al guardar el usuario')
    }
  } catch (error) {
    console.error('Error al guardar usuario:', error)
    showError('Error al guardar el usuario')
  } finally {
    saving.value = false
  }
}

// Confirmar eliminación
const showDeleteModal = ref(false)
const userToDelete = ref(null)
const isDeleting = ref(false)

const openDeleteModal = (user) => {
  userToDelete.value = user
  showDeleteModal.value = true
}

const closeDeleteModal = () => {
  showDeleteModal.value = false
  userToDelete.value = null
}

const confirmDelete = async () => {
  if (!userToDelete.value) return
  
  isDeleting.value = true
  try {
    await deleteUser(userToDelete.value.rut)
    closeDeleteModal()
  } catch (error) {
    console.error('Error al eliminar usuario:', error)
  } finally {
    isDeleting.value = false
  }
}

// Eliminar usuario
const deleteUser = async (userRut) => {
  try {
    const response = await userService.deleteUser(userRut)
    if (response.success) {
      showSuccess('Usuario eliminado correctamente')
      await loadUsers()
    } else {
      showError(response.error || 'Error al eliminar el usuario')
    }
  } catch (error) {
    console.error('Error al eliminar usuario:', error)
    showError('Error al eliminar el usuario')
  }
}

// Obtener clase CSS para el badge de rol
const getRoleBadgeClass = (role) => {
  const classes = {
    'admin': 'bg-purple-100 text-purple-800',
    'encargado de bodega': 'bg-blue-100 text-blue-800',
    'pabellón': 'bg-green-100 text-green-800',
    'pavedad': 'bg-teal-100 text-teal-800',
    'enfermera': 'bg-pink-100 text-pink-800',
    'doctor': 'bg-yellow-100 text-yellow-800'
  }
  return classes[role] || 'bg-gray-100 text-gray-800'
}

// Limpiar búsqueda
const clearSearch = () => {
  searchTerm.value = ''
  currentPage.value = 1
}

// Limpiar todos los filtros
const clearAllFilters = () => {
  searchTerm.value = ''
  roleFilter.value = ''
  medicalCenterFilter.value = ''
  statusFilter.value = ''
  currentPage.value = 1
}

// Cargar datos al montar
onMounted(() => {
  loadUsers()
  loadMedicalCenters()
  loadPavilions()
  loadSpecialties()
})

// Watchers para resetear paginación cuando cambian los filtros
watch([searchTerm, roleFilter, medicalCenterFilter, statusFilter], () => {
  currentPage.value = 1
})
</script>

<style scoped>
/* Estilos responsivos para la tabla */
.table-container {
  overflow-x: auto;
  border: 1px solid #e5e7eb;
  border-radius: 0.5rem;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
  -webkit-overflow-scrolling: touch;
  position: relative;
}

/* Scrollbar visible pero elegante */
.table-container::-webkit-scrollbar {
  height: 6px;
}

@media (min-width: 768px) {
  .table-container::-webkit-scrollbar {
    height: 8px;
  }
}

.table-container::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

.table-container::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 4px;
}

.table-container::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

.table-header {
  background-color: #f9fafb;
}

.table-header-cell {
  padding: 0.75rem 1.5rem;
  text-align: left;
  font-size: 0.75rem;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: #6b7280;
  white-space: nowrap;
}

.table-cell {
  padding: 1rem 1.5rem;
  white-space: nowrap;
  vertical-align: middle;
  font-size: 0.875rem;
}

/* Anchos específicos por columna */
.table-header th:nth-child(1),
.table-header-cell:first-child {
  min-width: 180px; /* Usuario */
}

.table-header th:nth-child(2) {
  min-width: 200px; /* Email */
}

.table-header th:nth-child(3) {
  min-width: 150px; /* Rol */
}

.table-header th:nth-child(4) {
  min-width: 160px; /* Centro Médico */
}

.table-header th:nth-child(5) {
  min-width: 160px; /* Pabellón */
}

.table-header th:nth-child(6) {
  min-width: 180px; /* Especialidad */
}

.table-header th:nth-child(7) {
  min-width: 100px; /* Estado */
}

.table-header th:nth-child(8) {
  min-width: 100px; /* Acciones */
}
</style>
