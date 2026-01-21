import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts'
import type { DashboardWeekData } from '../../types'
import { formatWeekRange, isCurrentWeek } from '../../utils'
import { chartColors } from '../../utils/colors'

interface WeeklyChartProps {
  data: DashboardWeekData[]
}

export function WeeklyChart({ data }: WeeklyChartProps) {
  if (data.length === 0) {
    return (
      <div className="h-80 flex items-center justify-center text-gray-500">
        No data available. Create a project to see your capacity chart.
      </div>
    )
  }

  const chartData = data.map((week) => ({
    ...week,
    weekLabel: formatWeekRange(week.weekStartDate),
    isCurrentWeek: isCurrentWeek(week.weekStartDate),
  }))

  return (
    <div className="h-80">
      <ResponsiveContainer width="100%" height="100%">
        <BarChart data={chartData} margin={{ top: 20, right: 30, left: 20, bottom: 5 }}>
          <CartesianGrid strokeDasharray="3 3" stroke={getComputedStyle(document.documentElement).getPropertyValue('--ld-border') || '#E5E7EB'} />
          <XAxis
            dataKey="weekLabel"
            tick={{ fontSize: 12 }}
            tickLine={false}
            axisLine={{ stroke: getComputedStyle(document.documentElement).getPropertyValue('--ld-border') || '#E5E7EB' }}
          />
          <YAxis
            tick={{ fontSize: 12 }}
            tickLine={false}
            axisLine={{ stroke: getComputedStyle(document.documentElement).getPropertyValue('--ld-border') || '#E5E7EB' }}
            label={{ value: 'Hours', angle: -90, position: 'insideLeft', style: { fontSize: 12 } }}
          />
          <Tooltip
            contentStyle={{
              backgroundColor: getComputedStyle(document.documentElement).getPropertyValue('--ld-surface') || 'white',
              border: `1px solid ${getComputedStyle(document.documentElement).getPropertyValue('--ld-border') || '#E5E7EB'}`,
              borderRadius: '8px',
              boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)',
              color: getComputedStyle(document.documentElement).getPropertyValue('--ld-text') || '#111111'
            }}
            formatter={(value: number, name: string) => [
              `${value} hrs`,
              name === 'totalPlannedHours' ? 'Planned' : 'Actual',
            ]}
            labelFormatter={(label) => `Week: ${label}`}
          />
          <Legend
            formatter={(value) => (value === 'totalPlannedHours' ? 'Planned Hours' : 'Actual Hours')}
          />
          <Bar
            dataKey="totalPlannedHours"
            fill={chartColors.planned}
            radius={[4, 4, 0, 0]}
            name="totalPlannedHours"
          />
          <Bar
            dataKey="totalActualHours"
            fill={chartColors.actual}
            radius={[4, 4, 0, 0]}
            name="totalActualHours"
          />
        </BarChart>
      </ResponsiveContainer>
    </div>
  )
}
