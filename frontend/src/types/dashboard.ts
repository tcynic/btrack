export interface DashboardWeekData {
  weekStartDate: string
  totalPlannedHours: number
  totalActualHours: number
  projectCount: number
}

export interface DashboardSummary {
  totalActiveProjects: number
  atRiskProjects: number
  totalPlannedThisWeek: number
  totalActualThisWeek: number
  totalPlannedNextWeek: number
}
