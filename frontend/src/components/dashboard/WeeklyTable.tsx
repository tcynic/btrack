import type { DashboardWeekData } from '../../types'
import { formatWeekRange, isCurrentWeek, isWeekInPast } from '../../utils'

interface WeeklyTableProps {
  data: DashboardWeekData[]
}

export function WeeklyTable({ data }: WeeklyTableProps) {
  if (data.length === 0) {
    return (
      <div className="text-center py-8 text-ld-muted">
        No weekly data available.
      </div>
    )
  }

  return (
    <div className="overflow-x-auto">
      <table className="min-w-full divide-y divide-ld-border">
        <thead className="bg-ld-surface2">
          <tr>
            <th className="px-6 py-3 text-left text-xs font-medium text-ld-muted uppercase tracking-wider">
              Week
            </th>
            <th className="px-6 py-3 text-right text-xs font-medium text-ld-muted uppercase tracking-wider">
              Projects
            </th>
            <th className="px-6 py-3 text-right text-xs font-medium text-ld-muted uppercase tracking-wider">
              Planned
            </th>
            <th className="px-6 py-3 text-right text-xs font-medium text-ld-muted uppercase tracking-wider">
              Actual
            </th>
            <th className="px-6 py-3 text-right text-xs font-medium text-ld-muted uppercase tracking-wider">
              Variance
            </th>
          </tr>
        </thead>
        <tbody className="bg-ld-surface divide-y divide-ld-border">
          {data.map((week) => {
            const variance = week.totalPlannedHours - week.totalActualHours
            const isCurrent = isCurrentWeek(week.weekStartDate)
            const isPast = isWeekInPast(week.weekStartDate)

            return (
              <tr
                key={week.weekStartDate}
                className={`
                  ${isCurrent ? 'bg-ld-surface2' : ''}
                  ${!isPast && !isCurrent ? 'opacity-60' : ''}
                `}
              >
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-ld-text">
                  <div className="flex items-center">
                    {formatWeekRange(week.weekStartDate)}
                    {isCurrent && (
                      <span className="ml-2 px-2 py-0.5 text-xs border border-ld-primary text-ld-primary rounded-full">
                        Current
                      </span>
                    )}
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-ld-muted">
                  {week.projectCount}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-ld-text">
                  {week.totalPlannedHours} hrs
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-ld-text">
                  {week.totalActualHours} hrs
                </td>
                <td className={`px-6 py-4 whitespace-nowrap text-sm text-right font-medium ${
                  variance > 0 ? 'text-[var(--ld-green)]' : variance < 0 ? 'text-[var(--ld-pink)]' : 'text-ld-muted'
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
