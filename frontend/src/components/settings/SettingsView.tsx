import { useEffect, useState } from 'react'
import { Card, Button } from '../ui'
import { GetBackupInfo, CreateBackup, RestoreBackup } from '../../../wailsjs/go/main/App'
import { useTheme } from '../../hooks/useTheme'

interface BackupInfo {
  databasePath: string
  databaseSize: number
  lastModified: string
}

export function SettingsView() {
  const [backupInfo, setBackupInfo] = useState<BackupInfo | null>(null)
  const [isLoading, setIsLoading] = useState(false)
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null)
  const { theme, setTheme } = useTheme()

  useEffect(() => {
    loadBackupInfo()
  }, [])

  const loadBackupInfo = async () => {
    try {
      const info = await GetBackupInfo()
      setBackupInfo(info)
    } catch (err) {
      console.error('Failed to load backup info:', err)
    }
  }

  const handleCreateBackup = async () => {
    setIsLoading(true)
    setMessage(null)
    try {
      const result = await CreateBackup()
      if (result) {
        setMessage({ type: 'success', text: `Backup created successfully at: ${result}` })
      } else {
        setMessage({ type: 'error', text: 'Backup cancelled' })
      }
    } catch (err) {
      console.error('Failed to create backup:', err)
      setMessage({ type: 'error', text: 'Failed to create backup' })
    } finally {
      setIsLoading(false)
    }
  }

  const handleRestoreBackup = async () => {
    if (!confirm('⚠️ Restoring a backup will replace your current data. Are you sure?')) {
      return
    }

    setIsLoading(true)
    setMessage(null)
    try {
      await RestoreBackup()
      setMessage({ type: 'success', text: 'Backup restored successfully. Please restart the application.' })
      await loadBackupInfo()
    } catch (err) {
      console.error('Failed to restore backup:', err)
      setMessage({ type: 'error', text: 'Failed to restore backup' })
    } finally {
      setIsLoading(false)
    }
  }

  const formatBytes = (bytes: number) => {
    if (bytes === 0) return '0 Bytes'
    const k = 1024
    const sizes = ['Bytes', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Settings</h1>
        <p className="text-gray-500 mt-1">Manage your application settings and data</p>
      </div>

      {message && (
        <div
          className={`p-4 rounded-lg ${
            message.type === 'success'
              ? 'bg-green-50 border border-green-200 text-green-700'
              : 'bg-red-50 border border-red-200 text-red-700'
          }`}
        >
          {message.text}
        </div>
      )}

      {/* Theme Section */}
      <Card>
        <div className="space-y-4">
          <div>
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-1">Appearance</h2>
            <p className="text-sm text-gray-500 dark:text-gray-400">
              Choose your preferred color scheme
            </p>
          </div>

          <div className="flex gap-3">
            <button
              onClick={() => setTheme('light')}
              className={`flex-1 px-4 py-3 rounded-lg border-2 transition-colors ${
                theme === 'light'
                  ? 'border-blue-600 bg-blue-50 dark:bg-blue-900/20'
                  : 'border-gray-300 dark:border-gray-600 hover:border-gray-400 dark:hover:border-gray-500'
              }`}
            >
              <div className="flex items-center justify-center gap-2">
                <svg className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"
                  />
                </svg>
                <span className="font-medium text-gray-900 dark:text-white">Light</span>
              </div>
            </button>

            <button
              onClick={() => setTheme('dark')}
              className={`flex-1 px-4 py-3 rounded-lg border-2 transition-colors ${
                theme === 'dark'
                  ? 'border-blue-600 bg-blue-50 dark:bg-blue-900/20'
                  : 'border-gray-300 dark:border-gray-600 hover:border-gray-400 dark:hover:border-gray-500'
              }`}
            >
              <div className="flex items-center justify-center gap-2">
                <svg className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"
                  />
                </svg>
                <span className="font-medium text-gray-900 dark:text-white">Dark</span>
              </div>
            </button>

            <button
              onClick={() => setTheme('system')}
              className={`flex-1 px-4 py-3 rounded-lg border-2 transition-colors ${
                theme === 'system'
                  ? 'border-blue-600 bg-blue-50 dark:bg-blue-900/20'
                  : 'border-gray-300 dark:border-gray-600 hover:border-gray-400 dark:hover:border-gray-500'
              }`}
            >
              <div className="flex items-center justify-center gap-2">
                <svg className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
                  />
                </svg>
                <span className="font-medium text-gray-900 dark:text-white">System</span>
              </div>
            </button>
          </div>
        </div>
      </Card>

      {/* Backup & Restore Section */}
      <Card>
        <div className="space-y-6">
          <div>
            <h2 className="text-lg font-semibold text-gray-900 mb-1">Backup & Restore</h2>
            <p className="text-sm text-gray-500">
              Create backups of your data or restore from a previous backup
            </p>
          </div>

          {backupInfo && (
            <div className="bg-gray-50 p-4 rounded-lg space-y-2">
              <div className="flex justify-between text-sm">
                <span className="text-gray-600">Database Location:</span>
                <span className="text-gray-900 font-mono text-xs">{backupInfo.databasePath}</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-gray-600">Database Size:</span>
                <span className="text-gray-900">{formatBytes(backupInfo.databaseSize)}</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-gray-600">Last Modified:</span>
                <span className="text-gray-900">{backupInfo.lastModified}</span>
              </div>
            </div>
          )}

          <div className="flex gap-4">
            <Button
              onClick={handleCreateBackup}
              disabled={isLoading}
              className="flex items-center gap-2"
            >
              <svg
                className="h-4 w-4"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M8 7H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-3m-1 4l-3 3m0 0l-3-3m3 3V4"
                />
              </svg>
              {isLoading ? 'Creating...' : 'Create Backup'}
            </Button>

            <Button
              onClick={handleRestoreBackup}
              disabled={isLoading}
              variant="secondary"
              className="flex items-center gap-2"
            >
              <svg
                className="h-4 w-4"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"
                />
              </svg>
              {isLoading ? 'Restoring...' : 'Restore from Backup'}
            </Button>
          </div>

          <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
            <div className="flex">
              <svg
                className="h-5 w-5 text-yellow-600 mr-3 flex-shrink-0"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                />
              </svg>
              <div>
                <p className="text-sm font-medium text-yellow-800">Important</p>
                <p className="text-sm text-yellow-700 mt-1">
                  Always create a backup before restoring from an old backup. Restoring will
                  replace all your current data.
                </p>
              </div>
            </div>
          </div>
        </div>
      </Card>
    </div>
  )
}
