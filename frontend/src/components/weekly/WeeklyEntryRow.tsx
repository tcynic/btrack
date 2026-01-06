import type { WeeklyEntryWithStatus } from '../../types'
import { Badge } from '../ui'
import { ActualHoursInput } from './ActualHoursInput'
import { formatWeekRange, isCurrentWeek } from '../../utils'
import type { Status } from '../../utils/colors'

interface WeeklyEntryRowProps {
  entry: WeeklyEntryWithStatus
  onUpdateActualHours: (entryId: number, actualHours: number) => Promise<void>
}

export function WeeklyEntryRow({ entry, onUpdateActualHours }: WeeklyEntryRowProps) {
  const isCurrent = isCurrentWeek(entry.weekStartDate)
  const canEdit = entry.isPastWeek || isCurrent

  return (
    <tr className={`${isCurrent ? 'bg-blue-50' : ''}`}>
      <td className="px-6 py-4 whitespace-nowrap">
        <div className="flex items-center">
          <span className="text-sm font-medium text-gray-900">
            Week {entry.weekNumber}
          </span>
          {isCurrent && (
            <span className="ml-2 px-2 py-0.5 text-xs bg-blue-100 text-blue-800 rounded-full">
              Current
            </span>
          )}
        </div>
        <span className="text-xs text-gray-500">{formatWeekRange(entry.weekStartDate)}</span>
      </td>
      <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">
        {entry.plannedHours} hrs
      </td>
      <td className="px-6 py-4 whitespace-nowrap text-sm text-right">
        <ActualHoursInput
          value={entry.actualHours}
          onChange={async (value) => {
            await onUpdateActualHours(entry.id, value)
          }}
          disabled={!canEdit}
        />
        <span className="ml-1 text-gray-500">hrs</span>
      </td>
      <td className="px-6 py-4 whitespace-nowrap text-sm text-right">
        <span className={`font-medium ${
          entry.variance > 0 ? 'text-green-600' : entry.variance < 0 ? 'text-red-600' : 'text-gray-600'
        }`}>
          {entry.variance > 0 ? '+' : ''}{entry.variance}
        </span>
      </td>
      <td className="px-6 py-4 whitespace-nowrap text-right">
        <Badge status={entry.status as Status} />
      </td>
    </tr>
  )
}
