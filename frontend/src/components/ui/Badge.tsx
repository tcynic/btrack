import type { Status } from '../../utils/colors'
import { getStatusBadgeClasses, getStatusLabel } from '../../utils/colors'

interface BadgeProps {
  status: Status
  showLabel?: boolean
  className?: string
}

export function Badge({ status, showLabel = true, className = '' }: BadgeProps) {
  return (
    <span className={`${getStatusBadgeClasses(status)} ${className}`}>
      {showLabel ? getStatusLabel(status) : status}
    </span>
  )
}

interface StatBadgeProps {
  value: number
  label: string
  trend?: 'up' | 'down' | 'neutral'
}

export function StatBadge({ value, label, trend }: StatBadgeProps) {
  const trendColors = {
    up: 'text-green-600',
    down: 'text-red-600',
    neutral: 'text-gray-600',
  }

  return (
    <div className="text-center">
      <div className={`text-2xl font-bold ${trend ? trendColors[trend] : 'text-gray-900'}`}>
        {value}
      </div>
      <div className="text-sm text-gray-500">{label}</div>
    </div>
  )
}
