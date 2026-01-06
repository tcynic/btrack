import type { WeeklyEntry } from '../../types'
import { formatWeekRange } from '../../utils'

interface HourDistributionPreviewProps {
  entries: WeeklyEntry[]
  isLoading?: boolean
}

export function HourDistributionPreview({ entries, isLoading }: HourDistributionPreviewProps) {
  if (isLoading) {
    return (
      <div className="p-4 bg-gray-50 rounded-lg">
        <div className="animate-pulse flex space-x-4">
          <div className="flex-1 space-y-2">
            <div className="h-4 bg-gray-200 rounded w-24"></div>
            <div className="h-8 bg-gray-200 rounded"></div>
          </div>
        </div>
      </div>
    )
  }

  if (entries.length === 0) {
    return (
      <div className="p-4 bg-gray-50 rounded-lg text-center text-gray-500 text-sm">
        Enter project details to preview hour distribution
      </div>
    )
  }

  const totalHours = entries.reduce((sum, e) => sum + e.plannedHours, 0)

  return (
    <div className="bg-gray-50 rounded-lg p-4">
      <div className="flex justify-between items-center mb-3">
        <h4 className="text-sm font-medium text-gray-700">Hour Distribution Preview</h4>
        <span className="text-sm text-gray-500">
          {totalHours} total hours across {entries.length} weeks
        </span>
      </div>

      <div className="max-h-48 overflow-y-auto">
        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-2">
          {entries.map((entry) => (
            <div
              key={entry.weekNumber}
              className="bg-white rounded-md p-2 border border-gray-200"
            >
              <div className="text-xs text-gray-500">Week {entry.weekNumber}</div>
              <div className="text-lg font-semibold text-blue-600">
                {entry.plannedHours} <span className="text-xs font-normal">hrs</span>
              </div>
              <div className="text-xs text-gray-400">{formatWeekRange(entry.weekStartDate)}</div>
            </div>
          ))}
        </div>
      </div>

      <div className="mt-3 text-xs text-gray-500">
        Hours are frontloaded: earlier weeks receive more hours when there's a remainder.
      </div>
    </div>
  )
}
