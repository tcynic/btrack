import { useCallback, useMemo } from 'react'
import type { Meeting, CreateMeetingInput, UpdateMeetingInput } from '../types'
import {
  CreateMeeting,
  GetMeetings,
  UpdateMeeting,
  DeleteMeeting,
} from '../../wailsjs/go/main/App'
import { useCrud } from './useCrud'

export function useMeetings(projectId: number) {
  const config = useMemo(() => ({
    loadFn: () => GetMeetings(projectId) as Promise<Meeting[]>,
    createFn: (input: CreateMeetingInput) => CreateMeeting(input) as Promise<Meeting>,
    updateFn: (input: UpdateMeetingInput) => UpdateMeeting(input) as Promise<Meeting>,
    deleteFn: DeleteMeeting,
    getId: (m: Meeting) => m.id,
  }), [projectId])

  const crud = useCrud<Meeting, CreateMeetingInput, UpdateMeetingInput>(config)

  // Alias for backward compatibility
  const loadMeetings = useCallback(() => crud.load(), [crud.load])
  const createMeeting = useCallback((input: CreateMeetingInput) => crud.create(input), [crud.create])
  const updateMeeting = useCallback((input: UpdateMeetingInput) => crud.update(input), [crud.update])
  const deleteMeeting = useCallback((id: number) => crud.remove(id), [crud.remove])

  return {
    meetings: crud.items,
    isLoading: crud.isLoading,
    error: crud.error,
    loadMeetings,
    createMeeting,
    updateMeeting,
    deleteMeeting,
  }
}
