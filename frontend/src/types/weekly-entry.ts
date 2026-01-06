export interface WeeklyEntry {
  id: number
  projectId: number
  weekStartDate: string
  weekNumber: number
  plannedHours: number
  actualHours: number
  createdAt: string
  updatedAt: string
}

export interface WeeklyEntryWithStatus extends WeeklyEntry {
  variance: number
  status: 'under' | 'over' | 'on-track' | 'pending'
  isPastWeek: boolean
}

export interface UpdateActualHoursInput {
  entryId: number
  actualHours: number
}
