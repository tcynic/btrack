import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts'
import type { MonthlyTrend } from '../../hooks/useReports'

interface TrendChartProps {
  data: MonthlyTrend[]
}

export function TrendChart({ data }: TrendChartProps) {
  if (data.length === 0) {
    return (
      <div className="flex items-center justify-center h-64 text-gray-500">
        No trend data available
      </div>
    )
  }

  return (
    <ResponsiveContainer width="100%" height={400}>
      <LineChart data={data} margin={{ top: 5, right: 30, left: 20, bottom: 5 }}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="month" />
        <YAxis label={{ value: 'Hours', angle: -90, position: 'insideLeft' }} />
        <Tooltip />
        <Legend />
        <Line
          type="monotone"
          dataKey="plannedHours"
          stroke="#3b82f6"
          name="Planned Hours"
          strokeWidth={2}
        />
        <Line
          type="monotone"
          dataKey="actualHours"
          stroke="#10b981"
          name="Actual Hours"
          strokeWidth={2}
        />
      </LineChart>
    </ResponsiveContainer>
  )
}
