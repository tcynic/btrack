import type { WeeklyEntryWithProject } from '../../types'
import { ActualHoursInput } from '../weekly/ActualHoursInput'
import { Badge } from '../ui'
import type { Status } from '../../utils/colors'

interface WeekHoursTableProps {
  entries: WeeklyEntryWithProject[]
  onUpdateActualHours: (entryId: number, actualHours: number) => Promise<void>
  isCurrentWeek: boolean
}

export function WeekHoursTable({ 
  entries, 
  onUpdateActualHours,
  isCurrentWeek 
}: WeekHoursTableProps) {
  if (entries.length === 0) {
    return (
      <div className="text-center py-12 text-gray-500">
        <svg
          className="w-12 h-12 mx-auto mb-3 text-gray-300"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={1.5}
            d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"
          />
        </svg>
        <p>No active projects for this week</p>
      </div>
    )
  }

  const totalPlanned = entries.reduce((sum, e) => sum + e.plannedHours, 0)
  const totalActual = entries.reduce((sum, e) => sum + e.actualHours, 0)
  const totalVariance = totalPlanned - totalActual

  return (
    <div className="overflow-x-auto">
      <table className="min-w-full divide-y divide-gray-200">
        <thead className="bg-gray-50">
          <tr>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Project
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
            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
              Status
            </th>
          </tr>
        </thead>
        <tbody className="bg-white divide-y divide-gray-200">
          {entries.map((entry) => {
            const canEdit = entry.isPastWeek || isCurrentWeek
            
            return (
              <tr key={entry.id} className={isCurrentWeek ? 'bg-blue-50' : ''}>
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="text-sm font-medium text-gray-900">
                    {entry.projectName}
                  </div>
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
          })}
        </tbody>
        <tfoot className="bg-gray-50">
          <tr>
            <td className="px-6 py-3 text-sm font-medium text-gray-900">
              Total
            </td>
            <td className="px-6 py-3 text-sm font-medium text-right text-gray-900">
              {totalPlanned} hrs
            </td>
            <td className="px-6 py-3 text-sm font-medium text-right text-gray-900">
              {totalActual} hrs
            </td>
            <td className={`px-6 py-3 text-sm font-medium text-right ${
              totalVariance > 0 ? 'text-green-600' : totalVariance < 0 ? 'text-red-600' : 'text-gray-900'
            }`}>
              {totalVariance > 0 ? '+' : ''}{totalVariance} hrs
            </td>
            <td></td>
          </tr>
        </tfoot>
      </table>
    </div>
  )
}
