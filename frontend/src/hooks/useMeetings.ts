import { useState, useCallback } from 'react'
import type { Meeting, CreateMeetingInput, UpdateMeetingInput } from '../types'
import {
  CreateMeeting,
  GetMeetings,
  UpdateMeeting,
  DeleteMeeting,
} from '../../wailsjs/go/main/App'

export function useMeetings(projectId: number) {
  const [meetings, setMeetings] = useState<Meeting[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const loadMeetings = useCallback(async () => {
    setIsLoading(true)
    setError(null)
    try {
      const data = await GetMeetings(projectId)
      setMeetings(data as Meeting[])
    } catch (err) {
      setError(String(err))
    } finally {
      setIsLoading(false)
    }
  }, [projectId])

  const createMeeting = useCallback(async (input: CreateMeetingInput) => {
    setIsLoading(true)
    setError(null)
    try {
      const meeting = await CreateMeeting(input)
      setMeetings((prev) => [meeting as Meeting, ...prev])
      return meeting as Meeting
    } catch (err) {
      setError(String(err))
      throw err
    } finally {
      setIsLoading(false)
    }
  }, [])

  const updateMeeting = useCallback(async (input: UpdateMeetingInput) => {
    setIsLoading(true)
    setError(null)
    try {
      const meeting = await UpdateMeeting(input)
      setMeetings((prev) =>
        prev.map((m) => (m.id === input.id ? (meeting as Meeting) : m))
      )
      return meeting as Meeting
    } catch (err) {
      setError(String(err))
      throw err
    } finally {
      setIsLoading(false)
    }
  }, [])

  const deleteMeeting = useCallback(async (id: number) => {
    setIsLoading(true)
    setError(null)
    try {
      await DeleteMeeting(id)
      setMeetings((prev) => prev.filter((m) => m.id !== id))
    } catch (err) {
      setError(String(err))
      throw err
    } finally {
      setIsLoading(false)
    }
  }, [])

  return {
    meetings,
    isLoading,
    error,
    loadMeetings,
    createMeeting,
    updateMeeting,
    deleteMeeting,
  }
}
