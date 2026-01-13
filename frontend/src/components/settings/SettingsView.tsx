import { useEffect, useState } from 'react'
import { Card, Button } from '../ui'
import { GetBackupInfo, CreateBackup, RestoreBackup } from '../../../wailsjs/go/main/App'

interface BackupInfo {
  databasePath: string
  databaseSize: number
  lastModified: string
}

export function SettingsView() {
  const [backupInfo, setBackupInfo] = useState<BackupInfo | null>(null)
  const [isLoading, setIsLoading] = useState(false)
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null)

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
