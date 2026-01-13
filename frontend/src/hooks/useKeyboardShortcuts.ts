import { useEffect, useState } from 'react'

export interface KeyboardShortcut {
  key: string
  label: string
  description: string
  modifiers?: ('cmd' | 'ctrl' | 'shift' | 'alt')[]
}

export const shortcuts: KeyboardShortcut[] = [
  { key: '1', label: '⌘1', description: 'Go to Dashboard', modifiers: ['cmd'] },
  { key: '2', label: '⌘2', description: 'Go to Week View', modifiers: ['cmd'] },
  { key: '3', label: '⌘3', description: 'Go to Projects', modifiers: ['cmd'] },
  { key: '4', label: '⌘4', description: 'Go to Reports', modifiers: ['cmd'] },
  { key: '5', label: '⌘5', description: 'Go to Settings', modifiers: ['cmd'] },
  { key: '?', label: '?', description: 'Show keyboard shortcuts' },
]

interface UseKeyboardShortcutsOptions {
  onTabChange?: (tab: 'dashboard' | 'week' | 'projects' | 'reports' | 'settings') => void
  onShowHelp?: () => void
}

export function useKeyboardShortcuts(options: UseKeyboardShortcutsOptions = {}) {
  const [showHelp, setShowHelp] = useState(false)

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      // Check for modifier keys (Cmd on Mac, Ctrl on Windows/Linux)
      const isMod = event.metaKey || event.ctrlKey

      // Tab navigation shortcuts (Cmd/Ctrl + 1-5)
      if (isMod && options.onTabChange) {
        switch (event.key) {
          case '1':
            event.preventDefault()
            options.onTabChange('dashboard')
            break
          case '2':
            event.preventDefault()
            options.onTabChange('week')
            break
          case '3':
            event.preventDefault()
            options.onTabChange('projects')
            break
          case '4':
            event.preventDefault()
            options.onTabChange('reports')
            break
          case '5':
            event.preventDefault()
            options.onTabChange('settings')
            break
        }
      }

      // Help shortcut (? key)
      if (event.key === '?' && !event.metaKey && !event.ctrlKey && !event.shiftKey) {
        // Only trigger if not in an input field
        const target = event.target as HTMLElement
        if (target.tagName !== 'INPUT' && target.tagName !== 'TEXTAREA') {
          event.preventDefault()
          setShowHelp(true)
          if (options.onShowHelp) {
            options.onShowHelp()
          }
        }
      }

      // Close help with Escape
      if (event.key === 'Escape') {
        setShowHelp(false)
      }
    }

    window.addEventListener('keydown', handleKeyDown)
    return () => window.removeEventListener('keydown', handleKeyDown)
  }, [options])

  return {
    showHelp,
    setShowHelp,
  }
}
