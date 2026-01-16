import { useState, useCallback } from 'react'

/**
 * Generic CRUD hook factory that provides standard async state management.
 * Reduces boilerplate across useMeetings, useNotes, useGoals, etc.
 */

export interface CrudState<T> {
  items: T[]
  isLoading: boolean
  error: string | null
}

export interface CrudActions<T, CreateInput, UpdateInput> {
  load: () => Promise<void>
  create: (input: CreateInput) => Promise<T>
  update: (input: UpdateInput) => Promise<T>
  remove: (id: number) => Promise<void>
  setItems: React.Dispatch<React.SetStateAction<T[]>>
}

export interface CrudHookResult<T, CreateInput, UpdateInput>
  extends CrudState<T>,
    CrudActions<T, CreateInput, UpdateInput> {}

interface CrudConfig<T, CreateInput, UpdateInput> {
  loadFn: () => Promise<T[]>
  createFn: (input: CreateInput) => Promise<T>
  updateFn: (input: UpdateInput) => Promise<T>
  deleteFn: (id: number) => Promise<void>
  getId: (item: T) => number
}

export function useCrud<T, CreateInput, UpdateInput>(
  config: CrudConfig<T, CreateInput, UpdateInput>
): CrudHookResult<T, CreateInput, UpdateInput> {
  const { loadFn, createFn, updateFn, deleteFn, getId } = config
  
  const [items, setItems] = useState<T[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const load = useCallback(async () => {
    setIsLoading(true)
    setError(null)
    try {
      const data = await loadFn()
      setItems(data as T[])
    } catch (err) {
      setError(String(err))
    } finally {
      setIsLoading(false)
    }
  }, [loadFn])

  const create = useCallback(async (input: CreateInput): Promise<T> => {
    setIsLoading(true)
    setError(null)
    try {
      const item = await createFn(input)
      setItems((prev) => [item as T, ...prev])
      return item as T
    } catch (err) {
      setError(String(err))
      throw err
    } finally {
      setIsLoading(false)
    }
  }, [createFn])

  const update = useCallback(async (input: UpdateInput): Promise<T> => {
    setIsLoading(true)
    setError(null)
    try {
      const item = await updateFn(input)
      setItems((prev) =>
        prev.map((i) => (getId(i) === getId(item as T) ? (item as T) : i))
      )
      return item as T
    } catch (err) {
      setError(String(err))
      throw err
    } finally {
      setIsLoading(false)
    }
  }, [updateFn, getId])

  const remove = useCallback(async (id: number): Promise<void> => {
    setIsLoading(true)
    setError(null)
    try {
      await deleteFn(id)
      setItems((prev) => prev.filter((i) => getId(i) !== id))
    } catch (err) {
      setError(String(err))
      throw err
    } finally {
      setIsLoading(false)
    }
  }, [deleteFn, getId])

  return {
    items,
    isLoading,
    error,
    load,
    create,
    update,
    remove,
    setItems,
  }
}
