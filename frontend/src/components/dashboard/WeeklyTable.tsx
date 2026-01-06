import type { DashboardWeekData } from '../../types'
import { formatWeekRange, isCurrentWeek, isWeekInPast } from '../../utils'

interface WeeklyTableProps {
  data: DashboardWeekData[]
}

export function WeeklyTable({ data }: WeeklyTableProps) {
  if (data.length === 0) {
    return (
      <div className="text-center py-8 text-gray-500">
        No weekly data available.
      </div>
    )
  }

  return (
    <div className="overflow-x-auto">
      <table className="min-w-full divide-y divide-gray-200">
        <thead className="bg-gray-50">
          <tr>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Week
            </th>
            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
              Projects
            </th>
            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
              Planned
            </th>
            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
              Actual
            </th>
            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
              Variance
            </th>
          </tr>
        </thead>
        <tbody className="bg-white divide-y divide-gray-200">
          {data.map((week) => {
            const variance = week.totalPlannedHours - week.totalActualHours
            const isCurrent = isCurrentWeek(week.weekStartDate)
            const isPast = isWeekInPast(week.weekStartDate)

            return (
              <tr
                key={week.weekStartDate}
                className={`
                  ${isCurrent ? 'bg-blue-50' : ''}
                  ${!isPast && !isCurrent ? 'text-gray-400' : ''}
                `}
              >
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                  <div className="flex items-center">
                    {formatWeekRange(week.weekStartDate)}
                    {isCurrent && (
                      <span className="ml-2 px-2 py-0.5 text-xs bg-blue-100 text-blue-800 rounded-full">
                        Current
                      </span>
                    )}
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-500">
                  {week.projectCount}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">
                  {week.totalPlannedHours} hrs
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">
                  {week.totalActualHours} hrs
                </td>
                <td className={`px-6 py-4 whitespace-nowrap text-sm text-right font-medium ${
                  variance > 0 ? 'text-green-600' : variance < 0 ? 'text-red-600' : 'text-gray-600'
                }`}>
                  {variance > 0 ? '+' : ''}{variance} hrs
                </td>
              </tr>
            )
          })}
        </tbody>
      </table>
    </div>
  )
}
