import { useEffect } from 'react'
import { WeeklyEntryRow } from './WeeklyEntryRow'
import { useWeeklyEntries } from '../../hooks'

interface WeeklyBreakdownProps {
  projectId: number
}

export function WeeklyBreakdown({ projectId }: WeeklyBreakdownProps) {
  const { entries, isLoading, loadEntries, updateActualHours } = useWeeklyEntries(projectId)

  useEffect(() => {
    loadEntries()
  }, [loadEntries])

  if (isLoading && entries.length === 0) {
    return (
      <div className="animate-pulse">
        <div className="h-10 bg-gray-200 rounded mb-2"></div>
        {[1, 2, 3, 4].map((i) => (
          <div key={i} className="h-16 bg-gray-100 rounded mb-1"></div>
        ))}
      </div>
    )
  }

  if (entries.length === 0) {
    return (
      <div className="text-center py-8 text-gray-500">
        No weekly entries found for this project.
      </div>
    )
  }

  const totalPlanned = entries.reduce((sum, e) => sum + e.plannedHours, 0)
  const totalActual = entries.reduce((sum, e) => sum + e.actualHours, 0)
  const totalVariance = totalPlanned - totalActual

  return (
    <div className="overflow-hidden">
      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Week
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
            {entries.map((entry) => (
              <WeeklyEntryRow
                key={entry.id}
                entry={entry}
                onUpdateActualHours={async (entryId, actualHours) => {
                  await updateActualHours({ entryId, actualHours })
                }}
              />
            ))}
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
    </div>
  )
}
