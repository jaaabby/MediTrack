<template>
  <div class="card overflow-hidden">
    <div
      class="overflow-x-auto"
      :class="maxHeight ? 'overflow-y-auto' : ''"
      :style="maxHeight ? `max-height: ${maxHeight}` : ''"
    >
        <table class="min-w-full divide-y divide-gray-200">
          <thead :class="maxHeight ? 'sticky top-0 z-20' : ''">
            <tr>
              <th
                v-for="col in columns"
                :key="col.key"
                :class="[
                  'px-6 py-3 text-xs font-medium text-gray-500 uppercase tracking-wider bg-gray-50',
                  col.align === 'right' ? 'text-right' : col.align === 'center' ? 'text-center' : 'text-left',
                  col.sortable !== false ? 'cursor-pointer hover:bg-gray-100 select-none' : '',
                  col.sticky ? 'sticky right-0 z-30 shadow-[-2px_0_4px_rgba(0,0,0,0.05)]' : '',
                  col.headerClass || ''
                ]"
                @click="col.sortable !== false ? sortBy(col.key) : undefined"
              >
                <div
                  v-if="col.sortable !== false"
                  class="flex items-center space-x-1"
                  :class="col.align === 'right' ? 'justify-end' : col.align === 'center' ? 'justify-center' : ''"
                >
                  <span>{{ col.label }}</span>
                  <span class="flex flex-col -space-y-0.5">
                    <svg
                      class="h-3 w-3 transition-colors"
                      :class="sortKey === col.key && sortOrder === 'asc' ? 'text-blue-600' : 'text-gray-400'"
                      fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"
                      viewBox="0 0 24 24"
                    >
                      <path d="M18 15l-6-6-6 6" />
                    </svg>
                    <svg
                      class="h-3 w-3 transition-colors"
                      :class="sortKey === col.key && sortOrder === 'desc' ? 'text-blue-600' : 'text-gray-400'"
                      fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"
                      viewBox="0 0 24 24"
                    >
                      <path d="M6 9l6 6 6-6" />
                    </svg>
                  </span>
                </div>
                <span v-else>{{ col.label }}</span>
              </th>

              <!-- Actions column header -->
              <th
                v-if="$slots.actions || tableActions.length"
                class="px-4 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider bg-gray-50 sticky right-0 z-30 shadow-[-2px_0_4px_rgba(0,0,0,0.05)]"
              >
                Acciones
              </th>
            </tr>
          </thead>

          <tbody class="bg-white divide-y divide-gray-200">
            <!-- Empty state -->
            <tr v-if="paginatedRows.length === 0">
              <td
                :colspan="columns.length + ($slots.actions ? 1 : 0)"
                class="px-6 py-12 text-center"
              >
                <slot name="empty">
                  <p class="text-sm text-gray-500">{{ emptyMessage }}</p>
                </slot>
              </td>
            </tr>

            <!-- Data rows -->
            <tr
              v-for="(row, index) in paginatedRows"
              :key="row.id ?? index"
              class="transition-colors duration-150"
              :class="rowClass ? rowClass(row) : 'hover:bg-gray-50'"
              @click="emit('row-click', row)"
            >
              <td
                v-for="col in columns"
                :key="col.key"
                :class="[
                  'px-6 py-4 text-sm text-gray-900',
                  col.wrap ? '' : 'whitespace-nowrap',
                  col.align === 'right' ? 'text-right' : col.align === 'center' ? 'text-center' : '',
                  col.sticky ? 'sticky right-0 bg-white z-10 shadow-[-2px_0_4px_rgba(0,0,0,0.05)]' : '',
                  col.cellClass || ''
                ]"
              >
                <slot :name="`cell-${col.key}`" :row="row" :value="row[col.key]">
                  {{ row[col.key] ?? '—' }}
                </slot>
              </td>

              <!-- Actions cell -->
              <td
                v-if="$slots.actions || tableActions.length"
                class="px-4 py-4 text-sm font-medium sticky right-0 bg-white z-10 shadow-[-2px_0_4px_rgba(0,0,0,0.05)]"
              >
                <div class="flex flex-row items-center justify-center gap-2">
                  <template v-for="action in tableActions" :key="action.type">
                    <!-- View -->
                    <button
                      v-if="action.type === 'view'"
                      @click.stop="action.onClick(row)"
                      class="text-blue-600 hover:text-blue-800 hover:bg-blue-50 p-1.5 rounded inline-flex items-center gap-1 text-xs transition-colors"
                      :title="action.label || 'Ver detalles'"
                    >
                      <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                      </svg>
                    </button>
                    <!-- Edit -->
                    <button
                      v-else-if="action.type === 'edit'"
                      @click.stop="action.onClick(row)"
                      class="text-warning-600 hover:text-warning-800 hover:bg-warning-50 p-1.5 rounded inline-flex items-center gap-1 text-xs transition-colors"
                      :title="action.label || 'Editar'"
                    >
                      <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                      </svg>
                    </button>
                    <!-- Delete -->
                    <button
                      v-else-if="action.type === 'delete'"
                      @click.stop="action.onClick(row)"
                      class="text-danger-600 hover:text-danger-800 hover:bg-danger-50 p-1.5 rounded inline-flex items-center justify-center transition-colors"
                      :title="action.label || 'Eliminar'"
                    >
                      <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                      </svg>
                    </button>
                    <!-- Notify -->
                    <button
                      v-else-if="action.type === 'notify'"
                      @click.stop="action.onClick(row)"
                      :disabled="action.disabled ? action.disabled(row) : false"
                      class="text-amber-600 hover:text-amber-800 hover:bg-amber-50 p-1.5 rounded inline-flex items-center gap-1 text-xs transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                      :title="action.title ? action.title(row) : (action.label || 'Notificar')"
                    >
                      <svg v-if="action.loading && action.loading(row)" class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                      </svg>
                      <svg v-else class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                      </svg>
                    </button>
                  </template>
                  <slot name="actions" :row="row" />
                </div>
              </td>
            </tr>
          </tbody>
        </table>
    </div>

    <!-- Pagination -->
    <div
      v-if="totalPages > 1"
      class="flex flex-col sm:flex-row items-center justify-between gap-4 px-6 py-4 border-t border-gray-200"
    >
      <div class="text-sm text-gray-700 text-center sm:text-left">
        Mostrando {{ startIndex + 1 }} a {{ endIndex }} de {{ sortedRows.length }} registros
      </div>
      <div class="flex items-center gap-2">
        <button
          class="btn-secondary px-3 py-2 text-sm min-w-[70px]"
          :disabled="currentPage === 1"
          @click="currentPage--"
        >
          <span class="hidden sm:inline">Anterior</span>
          <span class="sm:hidden">Ant.</span>
        </button>
        <span class="px-3 py-2 text-sm text-gray-700 bg-gray-100 rounded-md min-w-[90px] text-center">
          Página {{ currentPage }} de {{ totalPages }}
        </span>
        <button
          class="btn-secondary px-3 py-2 text-sm min-w-[70px]"
          :disabled="currentPage === totalPages"
          @click="currentPage++"
        >
          <span class="hidden sm:inline">Siguiente</span>
          <span class="sm:hidden">Sig.</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { normalize } from '@/utils/normalize'

const props = defineProps({
  /**
   * Column definitions.
   * Each column: { key, label, sortable?, align?, wrap?, sticky?, headerClass?, cellClass? }
   * - sortable: defaults to true. Set to false to disable sorting for that column.
   * - align: 'left' (default) | 'center' | 'right'
   * - wrap: set to true to allow cell text wrapping (default is whitespace-nowrap)
   * - sticky: true to make the column sticky on the right (for actions-style columns)
   */
  columns: { type: Array, required: true },
  /**
   * Rows to display. The parent is responsible for filtering; this component
   * handles sorting and pagination.
   */
  rows: { type: Array, default: () => [] },
  /** Initial sort column key. Defaults to the first sortable column. */
  defaultSortKey: { type: String, default: '' },
  /** Initial sort direction. */
  defaultSortOrder: { type: String, default: 'asc' },
  /** Number of rows per page. */
  itemsPerPage: { type: Number, default: 10 },
  /**
   * Function (row) => string that returns Tailwind classes for a row.
   * When provided it replaces the default 'hover:bg-gray-50'.
   */
  rowClass: { type: Function, default: null },
  /** CSS max-height for the scrollable table body. Empty string disables scroll. */
  maxHeight: { type: String, default: '600px' },
  /** Message shown when there are no rows. */
  emptyMessage: { type: String, default: 'No hay registros para mostrar' },
  /**
   * Predefined action buttons rendered in the actions column.
   * Each action: { type: 'view' | 'edit' | 'delete', label?: string, onClick: (row) => void }
   * - view:   blue  — eye icon
   * - edit:   orange — pencil icon
   * - delete: red   — trash icon (no label)
   */
  tableActions: { type: Array, default: () => [] }
})

const emit = defineEmits(['sort-change', 'row-click'])

// Resolve initial sort key
const initialSortKey = props.defaultSortKey || props.columns.find(c => c.sortable !== false)?.key || ''

const sortKey = ref(initialSortKey)
const sortOrder = ref(props.defaultSortOrder)
const currentPage = ref(1)

// Reset to page 1 when rows change (e.g. after a filter)
watch(() => props.rows, () => { currentPage.value = 1 })

const sortBy = (key) => {
  if (sortKey.value === key) {
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortKey.value = key
    sortOrder.value = 'asc'
  }
  currentPage.value = 1
  emit('sort-change', { key: sortKey.value, order: sortOrder.value })
}

const sortedRows = computed(() => {
  if (!sortKey.value) return [...props.rows]
  return [...props.rows].sort((a, b) => {
    let aVal = a[sortKey.value] ?? ''
    let bVal = b[sortKey.value] ?? ''
    if (typeof aVal === 'string') {
      aVal = normalize(aVal)
      bVal = normalize((bVal || '').toString())
    }
    if (aVal < bVal) return sortOrder.value === 'asc' ? -1 : 1
    if (aVal > bVal) return sortOrder.value === 'asc' ? 1 : -1
    return 0
  })
})

const totalPages = computed(() => Math.max(1, Math.ceil(sortedRows.value.length / props.itemsPerPage)))
const startIndex = computed(() => (currentPage.value - 1) * props.itemsPerPage)
const endIndex = computed(() => Math.min(startIndex.value + props.itemsPerPage, sortedRows.value.length))
const paginatedRows = computed(() => sortedRows.value.slice(startIndex.value, endIndex.value))

// Expose row count so the parent can read it (e.g. for a badge)
defineExpose({ rowCount: computed(() => sortedRows.value.length) })
</script>
