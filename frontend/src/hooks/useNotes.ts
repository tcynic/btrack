import { useCallback, useMemo } from 'react'
import type { Note, CreateNoteInput, UpdateNoteInput } from '../types'
import {
  CreateNote,
  GetNotes,
  UpdateNote,
  DeleteNote,
} from '../../wailsjs/go/main/App'
import { useCrud } from './useCrud'

export function useNotes(projectId: number) {
  const config = useMemo(() => ({
    loadFn: () => GetNotes(projectId) as Promise<Note[]>,
    createFn: (input: CreateNoteInput) => CreateNote(input) as Promise<Note>,
    updateFn: (input: UpdateNoteInput) => UpdateNote(input) as Promise<Note>,
    deleteFn: DeleteNote,
    getId: (n: Note) => n.id,
  }), [projectId])

  const crud = useCrud<Note, CreateNoteInput, UpdateNoteInput>(config)

  // Alias for backward compatibility
  const loadNotes = useCallback(() => crud.load(), [crud.load])
  const createNote = useCallback((input: CreateNoteInput) => crud.create(input), [crud.create])
  const updateNote = useCallback((input: UpdateNoteInput) => crud.update(input), [crud.update])
  const deleteNote = useCallback((id: number) => crud.remove(id), [crud.remove])

  return {
    notes: crud.items,
    isLoading: crud.isLoading,
    error: crud.error,
    loadNotes,
    createNote,
    updateNote,
    deleteNote,
  }
}
