interface TaskStatusBadgeProps {
  status: string
  onClick?: () => void
}

export function TaskStatusBadge({ status, onClick }: TaskStatusBadgeProps) {
  const getStatusClasses = (status: string) => {
    switch (status) {
      case 'in_progress':
        return 'bg-blue-100 text-blue-800'
      case 'completed':
        return 'bg-green-100 text-green-800'
      case 'cancelled':
        return 'bg-gray-100 text-gray-800'
      default:
        return 'bg-yellow-100 text-yellow-800'
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
