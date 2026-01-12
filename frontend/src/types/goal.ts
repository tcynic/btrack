export type GoalStatus = 'pending' | 'in_progress' | 'completed'

export interface Goal {
  id: number
  projectId: number
  title: string
  description: string
  status: GoalStatus
  targetDate: string
  createdAt: string
  updatedAt: string
}

export interface CreateGoalInput {
  projectId: number
  title: string
  description: string
  targetDate: string
}

export interface UpdateGoalInput {
  id: number
  title: string
  description: string
  status: GoalStatus
  targetDate: string
}
