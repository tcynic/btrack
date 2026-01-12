import { useState, useCallback } from 'react'
import type { Note, CreateNoteInput, UpdateNoteInput } from '../types'
import {
  CreateNote,
  GetNotes,
  UpdateNote,
  DeleteNote,
} from '../../wailsjs/go/main/App'

export function useNotes(projectId: number) {
  const [notes, setNotes] = useState<Note[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const loadNotes = useCallback(async () => {
    setIsLoading(true)
    setError(null)
    try {
      const data = await GetNotes(projectId)
      setNotes(data as Note[])
    } catch (err) {
      setError(String(err))
    } finally {
      setIsLoading(false)
    }
  }, [projectId])

  const createNote = useCallback(async (input: CreateNoteInput) => {
    setIsLoading(true)
    setError(null)
    try {
      const note = await CreateNote(input)
      setNotes((prev) => [note as Note, ...prev])
      return note as Note
    } catch (err) {
      setError(String(err))
      throw err
    } finally {
      setIsLoading(false)
    }
  }, [])

  const updateNote = useCallback(async (input: UpdateNoteInput) => {
    setIsLoading(true)
    setError(null)
    try {
      const note = await UpdateNote(input)
      setNotes((prev) =>
        prev.map((n) => (n.id === input.id ? (note as Note) : n))
      )
      return note as Note
    } catch (err) {
      setError(String(err))
      throw err
    } finally {
      setIsLoading(false)
    }
  }, [])

  const deleteNote = useCallback(async (id: number) => {
    setIsLoading(true)
    setError(null)
    try {
      await DeleteNote(id)
      setNotes((prev) => prev.filter((n) => n.id !== id))
    } catch (err) {
      setError(String(err))
      throw err
    } finally {
      setIsLoading(false)
    }
  }, [])

  return {
    notes,
    isLoading,
    error,
    loadNotes,
    createNote,
    updateNote,
    deleteNote,
  }
}
