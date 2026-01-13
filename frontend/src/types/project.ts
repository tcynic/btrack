export interface Project {
  id: number
  name: string
  totalSoldHours: number
  specialistHours: number
  startDate: string
  endDate: string
  isActive: boolean
  isPersistent: boolean
  createdAt: string
  updatedAt: string
}

export type HealthStatus = 'on_track' | 'at_risk' | 'over_budget' | 'completed'

export interface ProjectHealth {
  status: HealthStatus
  message: string
  severity: 'info' | 'warning' | 'error'
}

export interface ProjectWithStats extends Project {
  myHours: number
  totalWeeks: number
  totalPlannedHours: number
  totalActualHours: number
  health: ProjectHealth
}

export interface CreateProjectInput {
  name: string
  totalSoldHours: number
  specialistHours: number
  startDate: string
  endDate: string
}

export interface UpdateProjectInput {
  id: number
  name: string
  totalSoldHours: number
  specialistHours: number
  startDate: string
  endDate: string
  isActive: boolean
}
