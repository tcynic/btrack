export interface Project {
  id: number
  name: string
  totalSoldHours: number
  specialistHours: number
  startDate: string
  endDate: string
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export interface ProjectWithStats extends Project {
  myHours: number
  totalWeeks: number
  totalPlannedHours: number
  totalActualHours: number
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
