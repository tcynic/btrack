import { useState, useCallback, useEffect, useRef, useMemo } from 'react'

export interface UseQueryOptions<T> {
  queryFn: () => Promise<T>
  initialData?: T
  enabled?: boolean
  onSuccess?: (data: T) => void
  onError?: (error: Error) => void
}

export interface UseQueryResult<T> {
  data: T | undefined
  isLoading: boolean
  error: string | null
  refetch: () => Promise<void>
}

/**
 * Generic hook for read-only data fetching with loading and error states
 */
export function useQuery<T>(options: UseQueryOptions<T>): UseQueryResult<T> {
  const { queryFn, initialData, enabled = true, onSuccess, onError } = options

  const [data, setData] = useState<T | undefined>(initialData)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  
  // Use refs for callbacks to avoid infinite loops
  const onSuccessRef = useRef(onSuccess)
  const onErrorRef = useRef(onError)
  
  // Update refs when callbacks change
  useEffect(() => {
    onSuccessRef.current = onSuccess
    onErrorRef.current = onError
  }, [onSuccess, onError])

  const refetch = useCallback(async () => {
    if (!enabled) return

    setIsLoading(true)
    setError(null)
    try {
      const result = await queryFn()
      setData(result)
      onSuccessRef.current?.(result)
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : String(err)
      setError(errorMessage)
      if (err instanceof Error) {
        onErrorRef.current?.(err)
      } else {
        onErrorRef.current?.(new Error(errorMessage))
      }
    } finally {
      setIsLoading(false)
    }
  }, [queryFn, enabled])

  useEffect(() => {
    refetch()
  }, [refetch])

  return useMemo(
    () => ({
      data,
      isLoading,
      error,
      refetch,
    }),
    [data, isLoading, error, refetch]
  )
}
