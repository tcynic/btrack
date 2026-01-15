export type TaskStatus = 'pending' | 'in_progress' | 'completed' | 'cancelled'
export type TaskPriority = 'low' | 'medium' | 'high'
export type TaskSourceType = 'meeting' | 'note' | 'standalone'

// These match the Go models exactly
export interface Task {
  id: number
  projectId: number
  sourceType: string
  sourceId?: number
  title: string
  description: string
  status: string
  priority: string
  dueDate: string
  createdAt: string
  updatedAt: string
}

export interface TaskWithContext {
  id: number
  projectId: number
  sourceType: string
  sourceId?: number
  title: string
  description: string
  status: string
  priority: string
  dueDate: string
  createdAt: string
  updatedAt: string
  projectName: string
  sourceTitle: string
}

export interface CreateTaskInput {
  projectId: number
  sourceType: string
  sourceId?: number
  title: string
  description: string
  priority: string
  dueDate: string
}

export interface UpdateTaskInput {
  id: number
  title: string
  description: string
  status: string
  priority: string
  dueDate: string
}
