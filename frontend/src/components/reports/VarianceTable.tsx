import type { VarianceReport } from '../../hooks/useReports'

interface VarianceTableProps {
  data: VarianceReport[]
}

export function VarianceTable({ data }: VarianceTableProps) {
  if (data.length === 0) {
    return (
      <div className="flex items-center justify-center h-64 text-gray-500">
        No variance data available
      </div>
    )
  }

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
              % Used
            </th>
          </tr>
        </thead>
        <tbody className="bg-white divide-y divide-gray-200">
          {data.map((row) => (
            <tr key={row.projectId} className="hover:bg-gray-50">
              <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                {row.projectName}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-500">
                {row.planned}h
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-500">
                {row.actual}h
              </td>
              <td
                className={`px-6 py-4 whitespace-nowrap text-sm text-right font-medium ${
                  row.variance > 0
                    ? 'text-red-600'
                    : row.variance < 0
                    ? 'text-green-600'
                    : 'text-gray-500'
                }`}
              >
                {row.variance > 0 ? '+' : ''}
                {row.variance}h
              </td>
              <td
                className={`px-6 py-4 whitespace-nowrap text-sm text-right font-medium ${
                  row.percentage > 100
                    ? 'text-red-600'
                    : row.percentage > 80
                    ? 'text-yellow-600'
                    : 'text-green-600'
                }`}
              >
                {row.percentage}%
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
