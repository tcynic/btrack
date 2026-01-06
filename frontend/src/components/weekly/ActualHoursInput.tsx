import { useState, useEffect } from 'react'

interface ActualHoursInputProps {
  value: number
  onChange: (value: number) => Promise<void>
  disabled?: boolean
}

export function ActualHoursInput({ value, onChange, disabled }: ActualHoursInputProps) {
  const [localValue, setLocalValue] = useState(value.toString())
  const [isEditing, setIsEditing] = useState(false)
  const [isSaving, setIsSaving] = useState(false)

  useEffect(() => {
    if (!isEditing) {
      setLocalValue(value.toString())
    }
  }, [value, isEditing])

  const handleBlur = async () => {
    setIsEditing(false)
    const newValue = parseInt(localValue) || 0

    if (newValue !== value) {
      setIsSaving(true)
      try {
        await onChange(newValue)
      } catch {
        setLocalValue(value.toString())
      } finally {
        setIsSaving(false)
      }
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      (e.target as HTMLInputElement).blur()
    } else if (e.key === 'Escape') {
      setLocalValue(value.toString())
      setIsEditing(false)
    }
  }

  if (disabled) {
    return (
      <span className="text-gray-400 cursor-not-allowed">
        {value}
      </span>
    )
  }

  return (
    <div className="relative inline-block">
      <input
        type="number"
        value={localValue}
        onChange={(e) => {
          setLocalValue(e.target.value)
          setIsEditing(true)
        }}
        onBlur={handleBlur}
        onKeyDown={handleKeyDown}
        onFocus={() => setIsEditing(true)}
        min={0}
        className={`
          w-16 px-2 py-1 text-right border rounded
          focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500
          ${isSaving ? 'bg-gray-100' : 'bg-white'}
        `}
        disabled={isSaving}
      />
      {isSaving && (
        <span className="absolute right-2 top-1/2 -translate-y-1/2">
          <svg className="animate-spin h-4 w-4 text-blue-500" viewBox="0 0 24 24">
            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
          </svg>
        </span>
      )}
    </div>
  )
}
