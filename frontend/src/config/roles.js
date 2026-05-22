/**
 * Definición centralizada de roles y sus etiquetas de presentación.
 * Los valores (value) son los almacenados en BD — NO modificar.
 * Los labels usan lenguaje inclusivo para la interfaz de usuario.
 */
export const ROLE_LABELS = {
  'admin': 'Administrador',
  'encargado de bodega': 'Encargado/a de Bodega',
  'pabellón': 'Pabellón',
  'pavedad': 'Pavedad',
  'enfermera': 'Enfermero/a',
  'doctor': 'Médico/a'
}

export const ROLE_OPTIONS = [
  { value: 'admin', label: 'Administrador' },
  { value: 'encargado de bodega', label: 'Encargado/a de Bodega' },
  { value: 'pabellón', label: 'Pabellón' },
  { value: 'pavedad', label: 'Pavedad' },
  { value: 'enfermera', label: 'Enfermero/a' },
  { value: 'doctor', label: 'Médico/a' }
]

export function getRoleLabel(role) {
  return ROLE_LABELS[role] ?? role
}
