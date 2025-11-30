import * as XLSX from 'xlsx'

/**
 * Exporta datos a un archivo Excel
 * @param {Array} data - Array de objetos con los datos a exportar
 * @param {Array} columns - Array de objetos con la configuración de columnas { key: string, label: string }
 * @param {string} filename - Nombre del archivo (sin extensión)
 */
export function exportToExcel(data, columns, filename = 'export') {
  try {
    // Preparar los datos para Excel
    const excelData = data.map(item => {
      const row = {}
      columns.forEach(col => {
        // Obtener el valor, manejando valores anidados con notación de punto
        let value = item
        const keys = col.key.split('.')
        for (const key of keys) {
          value = value?.[key]
        }
        
        // Formatear el valor según el tipo
        if (value === null || value === undefined) {
          row[col.label] = ''
        } else if (col.formatter && typeof col.formatter === 'function') {
          row[col.label] = col.formatter(value, item)
        } else if (value instanceof Date) {
          row[col.label] = value.toLocaleString('es-CL')
        } else {
          row[col.label] = String(value)
        }
      })
      return row
    })

    // Crear workbook y worksheet
    const worksheet = XLSX.utils.json_to_sheet(excelData)
    const workbook = XLSX.utils.book_new()
    XLSX.utils.book_append_sheet(workbook, worksheet, 'Datos')

    // Ajustar ancho de columnas
    const columnWidths = columns.map(col => ({
      wch: Math.max(col.label.length, 15) // Ancho mínimo de 15 caracteres
    }))
    worksheet['!cols'] = columnWidths

    // Generar nombre de archivo con fecha
    const date = new Date()
    const dateStr = date.toISOString().split('T')[0]
    const fullFilename = `${filename}_${dateStr}.xlsx`

    // Descargar el archivo
    XLSX.writeFile(workbook, fullFilename)
    
    return true
  } catch (error) {
    console.error('Error al exportar a Excel:', error)
    throw error
  }
}

/**
 * Formatea una fecha para Excel
 */
export function formatDateForExcel(date) {
  if (!date) return ''
  return new Date(date).toLocaleString('es-CL', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

/**
 * Formatea un estado para Excel
 */
export function formatStatusForExcel(status) {
  const statusMap = {
    'pendiente': 'Pendiente',
    'en_transit': 'En Tránsito',
    'in_transit': 'En Tránsito',
    'en_transito': 'En Tránsito',
    'completed': 'Completado',
    'completado': 'Completado',
    'recibido': 'Recibido',
    'cancelled': 'Cancelado',
    'cancelado': 'Cancelado',
    'rechazado': 'Rechazado',
    'disponible': 'Disponible',
    'en_camino_a_pabellon': 'En Camino a Pabellón',
    'en_camino_a_bodega': 'En Camino a Bodega',
    'recepcionado': 'Recepcionado',
    'consumido': 'Consumido',
    'vencido': 'Vencido',
    'reservado': 'Reservado',
    'pendiente_pavedad': 'Pendiente Pavedad',
    'asignado_bodega': 'Asignado a Bodega',
    'devuelto': 'Devuelto al Solicitante',
    'devuelto_al_encargado': 'Devuelto al Encargado',
    'aprobado': 'Aprobado',
    'parcialmente_aprobado': 'Parcialmente Aprobado',
    'rechazado': 'Rechazado'
  }
  return statusMap[status] || status
}

