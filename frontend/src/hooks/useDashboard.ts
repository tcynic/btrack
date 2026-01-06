import { useCallback } from 'react'
import { useAppContext } from '../context/AppContext'
import type { DashboardSummary, DashboardWeekData } from '../types'
import { GetDashboardSummary, GetDashboardData } from '../../wailsjs/go/main/App'

export function useDashboard() {
  const { state, dispatch } = useAppContext()

  const loadDashboardSummary = useCallback(async () => {
    dispatch({ type: 'SET_LOADING', payload: true })
    dispatch({ type: 'SET_ERROR', payload: null })
    try {
      const summary = await GetDashboardSummary()
      dispatch({ type: 'SET_DASHBOARD_SUMMARY', payload: summary as DashboardSummary })
    } catch (err) {
      dispatch({ type: 'SET_ERROR', payload: String(err) })
    } finally {
      dispatch({ type: 'SET_LOADING', payload: false })
    }
  }, [dispatch])

  const loadDashboardData = useCallback(async (weeksBack: number = 4, weeksForward: number = 12) => {
    dispatch({ type: 'SET_LOADING', payload: true })
    dispatch({ type: 'SET_ERROR', payload: null })
    try {
      const data = await GetDashboardData(weeksBack, weeksForward)
      dispatch({ type: 'SET_DASHBOARD_WEEK_DATA', payload: data as DashboardWeekData[] })
    } catch (err) {
      dispatch({ type: 'SET_ERROR', payload: String(err) })
    } finally {
      dispatch({ type: 'SET_LOADING', payload: false })
    }
  }, [dispatch])

  const refreshDashboard = useCallback(async () => {
    await Promise.all([loadDashboardSummary(), loadDashboardData()])
  }, [loadDashboardSummary, loadDashboardData])

  return {
    summary: state.dashboardSummary,
    weekData: state.dashboardWeekData,
    isLoading: state.isLoading,
    error: state.error,
    loadDashboardSummary,
    loadDashboardData,
    refreshDashboard,
  }
}
