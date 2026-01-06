export type Status = 'under' | 'over' | 'on-track' | 'pending'

/**
 * Get the color class for a status
 */
export function getStatusColor(status: Status): string {
  switch (status) {
    case 'under':
      return 'text-green-600'
    case 'over':
      return 'text-red-600'
    case 'on-track':
      return 'text-blue-600'
    case 'pending':
      return 'text-gray-400'
    default:
      return 'text-gray-600'
  }
}

/**
 * Get the background color class for a status
 */
export function getStatusBgColor(status: Status): string {
  switch (status) {
    case 'under':
      return 'bg-green-100'
    case 'over':
      return 'bg-red-100'
    case 'on-track':
      return 'bg-blue-100'
    case 'pending':
      return 'bg-gray-100'
    default:
      return 'bg-gray-100'
  }
}

/**
 * Get the badge classes for a status
 */
export function getStatusBadgeClasses(status: Status): string {
  const base = 'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium'
  switch (status) {
    case 'under':
      return `${base} bg-green-100 text-green-800`
    case 'over':
      return `${base} bg-red-100 text-red-800`
    case 'on-track':
      return `${base} bg-blue-100 text-blue-800`
    case 'pending':
      return `${base} bg-gray-100 text-gray-800`
    default:
      return `${base} bg-gray-100 text-gray-800`
  }
}

/**
 * Get human-readable status label
 */
export function getStatusLabel(status: Status): string {
  switch (status) {
    case 'under':
      return 'Under Budget'
    case 'over':
      return 'Over Budget'
    case 'on-track':
      return 'On Track'
    case 'pending':
      return 'Pending'
    default:
      return 'Unknown'
  }
}

/**
 * Get chart colors for planned vs actual
 */
export const chartColors = {
  planned: '#3B82F6', // blue-500
  actual: '#10B981',  // emerald-500
  over: '#EF4444',    // red-500
}
