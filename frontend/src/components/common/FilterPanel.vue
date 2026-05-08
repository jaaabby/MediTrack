<template>
  <div class="card">
    <div :class="gridClass">
      <template v-for="filter in filters" :key="filter.key">
        <!-- Text input -->
        <div v-if="filter.type === 'text'" :class="filter.colSpan === 2 ? 'lg:col-span-2' : ''">
          <label class="block text-sm font-medium text-gray-700 mb-1">{{ filter.label }}</label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <svg class="h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            <input
              type="text"
              :value="modelValues[filter.key]"
              :placeholder="filter.placeholder || 'Buscar...'"
              class="form-input pl-9"
              @input="onChange(filter.key, $event.target.value)"
            />
          </div>
        </div>

        <!-- Select dropdown -->
        <div v-else-if="filter.type === 'select'">
          <label class="block text-sm font-medium text-gray-700 mb-1">{{ filter.label }}</label>
          <select
            :value="modelValues[filter.key]"
            class="form-input"
            @change="onChange(filter.key, $event.target.value)"
          >
            <option v-for="opt in filter.options" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
          </select>
        </div>

        <!-- Date input -->
        <div v-else-if="filter.type === 'date'">
          <label class="block text-sm font-medium text-gray-700 mb-1">{{ filter.label }}</label>
          <input
            type="date"
            :value="modelValues[filter.key]"
            class="form-input"
            @change="onChange(filter.key, $event.target.value)"
          />
        </div>

        <!-- Checkbox -->
        <div v-else-if="filter.type === 'checkbox'" :class="filter.colSpan === 2 ? 'lg:col-span-2' : ''">
          <label class="block text-sm font-medium text-gray-700 mb-1">{{ filter.label }}</label>
          <label class="inline-flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              :checked="Boolean(modelValues[filter.key])"
              class="h-4 w-4 rounded border-gray-300 text-blue-600 shadow-sm focus:border-blue-300 focus:ring focus:ring-blue-200 focus:ring-opacity-50"
              @change="onChange(filter.key, $event.target.checked)"
            />
            <span class="text-sm text-gray-700">{{ filter.checkboxLabel || 'Activar' }}</span>
          </label>
        </div>

        <!-- Custom filter slot -->
        <div v-else-if="filter.type === 'custom'" :class="filter.colSpan === 2 ? 'lg:col-span-2' : ''">
          <slot
            :name="`filter-${filter.key}`"
            :filter="filter"
            :value="modelValues[filter.key]"
            :setValue="(value) => onChange(filter.key, value)"
          />
        </div>
      </template>
    </div>

    <!-- Toggle filters + counter + clear -->
    <div v-if="toggleFilters.length > 0" class="mt-4 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
      <div v-for="filter in toggleFilters" :key="filter.key" class="flex items-center gap-3">
        <label class="text-sm font-medium text-gray-700">{{ filter.label }}</label>
        <div class="flex gap-2">
          <button
            v-for="opt in filter.options"
            :key="opt.value"
            @click="onChange(filter.key, opt.value)"
            :class="modelValues[filter.key] === opt.value
              ? (opt.activeClass || 'bg-blue-600 text-white')
              : 'bg-gray-100 text-gray-700 hover:bg-gray-200'"
            class="px-3 py-1.5 text-xs font-medium rounded-full transition-colors"
          >{{ opt.label }}</button>
        </div>
      </div>

      <div v-if="showClear || hasActionsSlot" class="flex items-center gap-3">
        <span v-if="resultCount !== null" class="text-sm text-gray-500">{{ resultCount }} resultado(s)</span>
        <button
          v-if="showClear"
          @click="clearAll"
          :disabled="!hasActiveFilters"
          class="btn-secondary text-sm disabled:opacity-40 disabled:cursor-not-allowed"
        >Limpiar filtros</button>
        <slot name="actions" :has-active-filters="hasActiveFilters" :clear-all="clearAll" />
      </div>
    </div>

    <!-- Simple clear row when no toggles -->
    <div v-else-if="showClear || hasActionsSlot" class="mt-3 flex items-center justify-end gap-3">
      <span v-if="resultCount !== null" class="text-sm text-gray-500">{{ resultCount }} resultado(s)</span>
      <button
        v-if="showClear"
        @click="clearAll"
        :disabled="!hasActiveFilters"
        class="btn-secondary text-sm disabled:opacity-40 disabled:cursor-not-allowed"
      >Limpiar filtros</button>
      <slot name="actions" :has-active-filters="hasActiveFilters" :clear-all="clearAll" />
    </div>
  </div>
</template>

<script setup>
import { reactive, computed, watch, useSlots } from 'vue'

const props = defineProps({
  filters: {
    type: Array,
    required: true
  },
  showClear: {
    type: Boolean,
    default: true
  },
  resultCount: {
    type: Number,
    default: null
  }
})

const emit = defineEmits(['filter-change'])
const slots = useSlots()

// Initialize internal state from filter defaults
// `initialValue` allows pre-populating a field (e.g. from URL query params) without
// changing the `default` (reset target), so hasActiveFilters detects it correctly.
const modelValues = reactive(
  Object.fromEntries(props.filters.map(f => [f.key, f.initialValue ?? f.default ?? '']))
)

// Re-initialize when filters prop changes (e.g. new keys added)
watch(() => props.filters.map(f => f.key).join(','), () => {
  props.filters.forEach(f => {
    if (!(f.key in modelValues)) {
      modelValues[f.key] = f.default ?? ''
    }
  })
})

const toggleFilters = computed(() => props.filters.filter(f => f.type === 'toggle'))
const nonToggleFilters = computed(() => props.filters.filter(f => f.type !== 'toggle'))
const hasActionsSlot = computed(() => Boolean(slots.actions))

const gridClass = computed(() => {
  const count = nonToggleFilters.value.length
  if (count <= 1) return 'grid grid-cols-1 gap-4'
  if (count <= 2) return 'grid grid-cols-1 sm:grid-cols-2 gap-4'
  return 'grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4'
})

const hasActiveFilters = computed(() =>
  props.filters.some(f => modelValues[f.key] !== (f.default ?? ''))
)

const onChange = (key, value) => {
  modelValues[key] = value
  emit('filter-change', key, value)
}

const clearAll = () => {
  props.filters.forEach(f => {
    modelValues[f.key] = f.default ?? ''
    emit('filter-change', f.key, f.default ?? '')
  })
}
</script>
