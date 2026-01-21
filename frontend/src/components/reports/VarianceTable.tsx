import type { VarianceReport } from '../../hooks/useReports'

interface VarianceTableProps {
  data: VarianceReport[]
}

export function VarianceTable({ data }: VarianceTableProps) {
  if (data.length === 0) {
    return (
      <div className="flex items-center justify-center h-64 text-ld-muted">
        No variance data available
      </div>
    )
  }

  return (
    <div className="overflow-x-auto">
      <table className="min-w-full divide-y divide-ld-border">
        <thead className="bg-ld-surface2">
          <tr>
            <th className="px-6 py-3 text-left text-xs font-medium text-ld-muted uppercase tracking-wider">
              Project
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
            <th className="px-6 py-3 text-right text-xs font-medium text-ld-muted uppercase tracking-wider">
              % Used
            </th>
          </tr>
        </thead>
        <tbody className="bg-ld-surface divide-y divide-ld-border">
          {data.map((row) => (
            <tr key={row.projectId} className="hover:bg-ld-surface2">
              <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-ld-text">
                {row.projectName}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-ld-muted">
                {row.planned}h
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-ld-muted">
                {row.actual}h
              </td>
              <td
                className={`px-6 py-4 whitespace-nowrap text-sm text-right font-medium ${
                  row.variance > 0
                    ? 'text-[var(--ld-pink)]'
                    : row.variance < 0
                    ? 'text-[var(--ld-green)]'
                    : 'text-ld-muted'
                }`}
              >
                {row.variance > 0 ? '+' : ''}
                {row.variance}h
              </td>
              <td
                className={`px-6 py-4 whitespace-nowrap text-sm text-right font-medium ${
                  row.percentage > 100
                    ? 'text-[var(--ld-pink)]'
                    : row.percentage > 80
                    ? 'text-[var(--ld-orange)]'
                    : 'text-[var(--ld-green)]'
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
