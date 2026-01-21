export type Status = 'under' | 'over' | 'on-track' | 'pending'

/**
 * Get the color class for a status (text)
 */
export function getStatusColor(status: Status): string {
  switch (status) {
    case 'under':
      return 'text-[var(--ld-green)]'
    case 'over':
      return 'text-[var(--ld-pink)]'
    case 'on-track':
      return 'text-ld-primary'
    case 'pending':
      return 'text-ld-muted'
    default:
      return 'text-ld-muted'
  }
}

/**
 * Get the background color class for a status (kept minimal, rely on border/text for brand pop)
 */
export function getStatusBgColor(_status: Status): string {
  return 'bg-transparent'
}

/**
 * Get the badge classes for a status
 */
export function getStatusBadgeClasses(status: Status): string {
  const base = 'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border';
  switch (status) {
    case 'under':
      return `${base} border-[var(--ld-green)] text-[var(--ld-green)]`
    case 'over':
      return `${base} border-[var(--ld-pink)] text-[var(--ld-pink)]`
    case 'on-track':
      return `${base} border-ld-primary text-ld-primary`
    case 'pending':
      return `${base} border-ld-border text-ld-muted`
    default:
      return `${base} border-ld-border text-ld-muted`
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
 * Chart colors for planned vs actual using brand palette
 */
export const chartColors = {
  planned: '#7084FF', // LD Powder Blue
  actual: '#A9FF5E',  // LD Green
  over: '#FF35A2',    // LD Pink
}
