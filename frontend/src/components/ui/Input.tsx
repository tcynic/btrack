import { InputHTMLAttributes, forwardRef } from 'react'

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string
  error?: string
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ label, error, className = '', id, ...props }, ref) => {
    const inputId = id || label?.toLowerCase().replace(/\s+/g, '-')

    return (
      <div className="w-full">
        {label && (
          <label
            htmlFor={inputId}
            className="block text-sm font-medium text-ld-text mb-1"
          >
            {label}
          </label>
        )}
        <input
          ref={ref}
          id={inputId}
          className={`
            block w-full rounded-md border-ld-border shadow-sm
            focus:border-ld-primary focus:ring-ld-primary
            disabled:bg-ld-surface2 disabled:text-ld-text
            text-sm px-3 py-2 border bg-ld-surface text-ld-text
            ${error ? 'border-[var(--ld-pink)] focus:border-[var(--ld-pink)] focus:ring-[var(--ld-pink)]' : ''}
            ${className}
          `}
          {...props}
        />
        {error && (
          <p className="mt-1 text-sm text-[var(--ld-pink)]">{error}</p>
        )}
      </div>
    )
  }
)

Input.displayName = 'Input'
