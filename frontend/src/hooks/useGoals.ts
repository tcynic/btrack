import { useState, useCallback, useMemo } from 'react'
import type { Goal, GoalStatus, CreateGoalInput, UpdateGoalInput } from '../types'
import {
  CreateGoal,
  GetGoals,
  UpdateGoal,
  UpdateGoalStatus,
  DeleteGoal,
} from '../../wailsjs/go/main/App'
import { useCrud } from './useCrud'

export function useGoals(projectId: number) {
  const config = useMemo(() => ({
    loadFn: () => GetGoals(projectId) as Promise<Goal[]>,
    createFn: (input: CreateGoalInput) => CreateGoal(input) as Promise<Goal>,
    updateFn: (input: UpdateGoalInput) => UpdateGoal(input) as Promise<Goal>,
    deleteFn: DeleteGoal,
    getId: (g: Goal) => g.id,
  }), [projectId])

  const crud = useCrud<Goal, CreateGoalInput, UpdateGoalInput>(config)
  const [statusLoading, setStatusLoading] = useState(false)

  // Alias for backward compatibility
  const loadGoals = useCallback(() => crud.load(), [crud.load])
  const createGoal = useCallback((input: CreateGoalInput) => crud.create(input), [crud.create])
  const updateGoal = useCallback((input: UpdateGoalInput) => crud.update(input), [crud.update])
  const deleteGoal = useCallback((id: number) => crud.remove(id), [crud.remove])

  // Status update is a special case not covered by generic CRUD
  const updateGoalStatus = useCallback(async (id: number, status: GoalStatus) => {
    setStatusLoading(true)
    try {
      const goal = await UpdateGoalStatus(id, status)
      crud.setItems((prev) =>
        prev.map((g) => (g.id === id ? (goal as Goal) : g))
      )
      return goal as Goal
    } catch (err) {
      throw err
    } finally {
      setStatusLoading(false)
    }
  }, [crud.setItems])

  return {
    goals: crud.items,
    isLoading: crud.isLoading || statusLoading,
    error: crud.error,
    loadGoals,
    createGoal,
    updateGoal,
    updateGoalStatus,
    deleteGoal,
  }
}
