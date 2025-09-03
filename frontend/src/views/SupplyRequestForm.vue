<template>
  <SupplyRequestForm :id="id" :editMode="editMode" />
</template>

<script setup>
import SupplyRequestForm from '@/components/SupplyRequestForm.vue'

// Props recibidos desde las rutas
const props = defineProps({
  id: {
    type: Number,
    default: null
  },
  editMode: {
    type: Boolean,
    default: false
  }
})
</script>