import { useState, useCallback } from 'react'
import type { Goal, GoalStatus, CreateGoalInput, UpdateGoalInput } from '../types'
import {
  CreateGoal,
  GetGoals,
  UpdateGoal,
  UpdateGoalStatus,
  DeleteGoal,
} from '../../wailsjs/go/main/App'

export function useGoals(projectId: number) {
  const [goals, setGoals] = useState<Goal[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const loadGoals = useCallback(async () => {
    setIsLoading(true)
    setError(null)
    try {
      const data = await GetGoals(projectId)
      setGoals(data as Goal[])
    } catch (err) {
      setError(String(err))
    } finally {
      setIsLoading(false)
    }
  }, [projectId])

  const createGoal = useCallback(async (input: CreateGoalInput) => {
    setIsLoading(true)
    setError(null)
    try {
      const goal = await CreateGoal(input)
      setGoals((prev) => [goal as Goal, ...prev])
      return goal as Goal
    } catch (err) {
      setError(String(err))
      throw err
    } finally {
      setIsLoading(false)
    }
  }, [])

  const updateGoal = useCallback(async (input: UpdateGoalInput) => {
    setIsLoading(true)
    setError(null)
    try {
      const goal = await UpdateGoal(input)
      setGoals((prev) =>
        prev.map((g) => (g.id === input.id ? (goal as Goal) : g))
      )
      return goal as Goal
    } catch (err) {
      setError(String(err))
      throw err
    } finally {
      setIsLoading(false)
    }
  }, [])

  const updateGoalStatus = useCallback(async (id: number, status: GoalStatus) => {
    setIsLoading(true)
    setError(null)
    try {
      const goal = await UpdateGoalStatus(id, status)
      setGoals((prev) =>
        prev.map((g) => (g.id === id ? (goal as Goal) : g))
      )
      return goal as Goal
    } catch (err) {
      setError(String(err))
      throw err
    } finally {
      setIsLoading(false)
    }
  }, [])

  const deleteGoal = useCallback(async (id: number) => {
    setIsLoading(true)
    setError(null)
    try {
      await DeleteGoal(id)
      setGoals((prev) => prev.filter((g) => g.id !== id))
    } catch (err) {
      setError(String(err))
      throw err
    } finally {
      setIsLoading(false)
    }
  }, [])

  return {
    goals,
    isLoading,
    error,
    loadGoals,
    createGoal,
    updateGoal,
    updateGoalStatus,
    deleteGoal,
  }
}
