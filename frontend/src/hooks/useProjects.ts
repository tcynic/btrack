import { useCallback } from 'react'
import { useAppContext } from '../context/AppContext'
import type { CreateProjectInput, UpdateProjectInput, ProjectWithStats } from '../types'
import {
  CreateProject,
  GetAllProjects,
  GetProject,
  UpdateProject,
  DeleteProject,
  ToggleProjectActive,
} from '../../wailsjs/go/main/App'

export function useProjects() {
  const { state, dispatch } = useAppContext()

  const loadProjects = useCallback(async (activeOnly: boolean = true) => {
    dispatch({ type: 'SET_LOADING', payload: true })
    dispatch({ type: 'SET_ERROR', payload: null })
    try {
      const projects = await GetAllProjects(activeOnly)
      dispatch({ type: 'SET_PROJECTS', payload: projects as ProjectWithStats[] })
    } catch (err) {
      dispatch({ type: 'SET_ERROR', payload: String(err) })
    } finally {
      dispatch({ type: 'SET_LOADING', payload: false })
    }
  }, [dispatch])

  const createProject = useCallback(async (input: CreateProjectInput) => {
    dispatch({ type: 'SET_LOADING', payload: true })
    dispatch({ type: 'SET_ERROR', payload: null })
    try {
      const project = await CreateProject(input)
      dispatch({ type: 'ADD_PROJECT', payload: project as ProjectWithStats })
      return project as ProjectWithStats
    } catch (err) {
      dispatch({ type: 'SET_ERROR', payload: String(err) })
      throw err
    } finally {
      dispatch({ type: 'SET_LOADING', payload: false })
    }
  }, [dispatch])

  const updateProject = useCallback(async (input: UpdateProjectInput) => {
    dispatch({ type: 'SET_LOADING', payload: true })
    dispatch({ type: 'SET_ERROR', payload: null })
    try {
      const project = await UpdateProject(input)
      dispatch({ type: 'UPDATE_PROJECT', payload: project as ProjectWithStats })
      return project as ProjectWithStats
    } catch (err) {
      dispatch({ type: 'SET_ERROR', payload: String(err) })
      throw err
    } finally {
      dispatch({ type: 'SET_LOADING', payload: false })
    }
  }, [dispatch])

  const deleteProject = useCallback(async (id: number) => {
    dispatch({ type: 'SET_LOADING', payload: true })
    dispatch({ type: 'SET_ERROR', payload: null })
    try {
      await DeleteProject(id)
      dispatch({ type: 'DELETE_PROJECT', payload: id })
    } catch (err) {
      dispatch({ type: 'SET_ERROR', payload: String(err) })
      throw err
    } finally {
      dispatch({ type: 'SET_LOADING', payload: false })
    }
  }, [dispatch])

  const toggleProjectActive = useCallback(async (id: number) => {
    dispatch({ type: 'SET_LOADING', payload: true })
    dispatch({ type: 'SET_ERROR', payload: null })
    try {
      const project = await ToggleProjectActive(id)
      dispatch({ type: 'UPDATE_PROJECT', payload: project as ProjectWithStats })
      return project as ProjectWithStats
    } catch (err) {
      dispatch({ type: 'SET_ERROR', payload: String(err) })
      throw err
    } finally {
      dispatch({ type: 'SET_LOADING', payload: false })
    }
  }, [dispatch])

  const getProject = useCallback(async (id: number) => {
    try {
      const project = await GetProject(id)
      return project as ProjectWithStats
    } catch (err) {
      dispatch({ type: 'SET_ERROR', payload: String(err) })
      throw err
    }
  }, [dispatch])

  const selectProject = useCallback((id: number | null) => {
    dispatch({ type: 'SELECT_PROJECT', payload: id })
  }, [dispatch])

  return {
    projects: state.projects,
    selectedProjectId: state.selectedProjectId,
    isLoading: state.isLoading,
    error: state.error,
    loadProjects,
    createProject,
    updateProject,
    deleteProject,
    toggleProjectActive,
    getProject,
    selectProject,
  }
}
