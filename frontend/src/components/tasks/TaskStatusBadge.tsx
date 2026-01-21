interface TaskStatusBadgeProps {
  status: string
  onClick?: () => void
}

export function TaskStatusBadge({ status, onClick }: TaskStatusBadgeProps) {
  const getStatusClasses = (status: string) => {
    switch (status) {
      case 'in_progress':
        return 'border border-ld-primary text-ld-primary bg-transparent'
      case 'completed':
        return 'border border-[var(--ld-green)] text-[var(--ld-green)] bg-transparent'
      case 'cancelled':
        return 'border border-ld-border text-ld-muted bg-transparent'
      default:
        return 'border border-[var(--ld-orange)] text-[var(--ld-orange)] bg-transparent'
    }
  }

  const getStatusLabel = (status: string) => {
    switch (status) {
      case 'in_progress':
        return 'In Progress'
      case 'completed':
        return 'Completed'
      case 'cancelled':
        return 'Cancelled'
      default:
        return 'Pending'
    }
  }

  return (
    <span
      className={`px-2 py-1 text-xs font-medium rounded-full ${getStatusClasses(status)} ${onClick ? 'cursor-pointer hover:opacity-80' : ''}`}
      onClick={onClick}
    >
      {getStatusLabel(status)}
    </span>
  )
}
