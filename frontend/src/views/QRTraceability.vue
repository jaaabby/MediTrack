<template>
  <QRTraceability :qrCode="qrCode" />
</template>

<script setup>
import QRTraceability from '@/components/QRTraceability.vue'

// Props recibidos desde las rutas
const props = defineProps({
  qrCode: {
    type: String,
    required: true
  }
})
</script>