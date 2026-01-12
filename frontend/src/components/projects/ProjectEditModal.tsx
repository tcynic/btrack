import { useState, useEffect, useCallback } from 'react'
import { Modal, Button, Input } from '../ui'
import { HourDistributionPreview } from './HourDistributionPreview'
import { useDistributionPreview } from '../../hooks'
import type { UpdateProjectInput, ProjectWithStats } from '../../types'

interface ProjectEditModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (input: UpdateProjectInput) => Promise<void>
  project: ProjectWithStats
}

export function ProjectEditModal({ isOpen, onClose, onSubmit, project }: ProjectEditModalProps) {
  const [formData, setFormData] = useState<UpdateProjectInput>({
    id: project.id,
    name: project.name,
    totalSoldHours: project.totalSoldHours,
    specialistHours: project.specialistHours,
    startDate: project.startDate,
    endDate: project.endDate,
    isActive: project.isActive,
  })
  const [specialistMode, setSpecialistMode] = useState<'flat' | 'percent'>('flat')
  const [specialistPercent, setSpecialistPercent] = useState(0)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const { preview, isLoading: previewLoading, calculatePreview, clearPreview } = useDistributionPreview()

  // Initialize form with project data when modal opens
  useEffect(() => {
    if (isOpen) {
      setFormData({
        id: project.id,
        name: project.name,
        totalSoldHours: project.totalSoldHours,
        specialistHours: project.specialistHours,
        startDate: project.startDate,
        endDate: project.endDate,
        isActive: project.isActive,
      })
      // Calculate initial percentage if specialist hours exist
      if (project.specialistHours > 0 && project.totalSoldHours > 0) {
        setSpecialistPercent(Math.round((project.specialistHours / project.totalSoldHours) * 100))
      }
      setError(null)
    }
  }, [isOpen, project])

  // Calculate preview when form data changes
  const updatePreview = useCallback(() => {
    if (formData.name && formData.startDate && formData.endDate && formData.totalSoldHours > 0) {
      calculatePreview({
        name: formData.name,
        totalSoldHours: formData.totalSoldHours,
        specialistHours: formData.specialistHours,
        startDate: formData.startDate,
        endDate: formData.endDate,
      })
    } else {
      clearPreview()
    }
  }, [formData, calculatePreview, clearPreview])

  useEffect(() => {
    const timer = setTimeout(updatePreview, 300) // Debounce
    return () => clearTimeout(timer)
  }, [updatePreview])

  // Update specialist hours when mode or percent changes
  useEffect(() => {
    if (specialistMode === 'percent') {
      const hours = Math.floor((formData.totalSoldHours * specialistPercent) / 100)
      setFormData((prev) => ({ ...prev, specialistHours: hours }))
    }
  }, [specialistMode, specialistPercent, formData.totalSoldHours])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setIsSubmitting(true)

    try {
      await onSubmit(formData)
      onClose()
    } catch (err) {
      setError(String(err))
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Edit Project" size="lg">
      <form onSubmit={handleSubmit}>
        <div className="space-y-4">
          <Input
            label="Project Name"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            placeholder="Enter project name"
            required
          />

          <div className="grid grid-cols-2 gap-4">
            <Input
              label="Start Date"
              type="date"
              value={formData.startDate}
              onChange={(e) => setFormData({ ...formData, startDate: e.target.value })}
              required
            />
            <Input
              label="End Date"
              type="date"
              value={formData.endDate}
              onChange={(e) => setFormData({ ...formData, endDate: e.target.value })}
              min={formData.startDate}
              required
            />
          </div>

          <Input
            label="Total Sold Hours"
            type="number"
            value={formData.totalSoldHours || ''}
            onChange={(e) => setFormData({ ...formData, totalSoldHours: parseInt(e.target.value) || 0 })}
            placeholder="Enter total hours sold"
            min={1}
            required
          />

          <div>
            <div className="flex items-center justify-between mb-2">
              <label className="block text-sm font-medium text-gray-700">
                Specialist Hours
              </label>
              <div className="flex rounded-md shadow-sm">
                <button
                  type="button"
                  onClick={() => setSpecialistMode('flat')}
                  className={`px-3 py-1 text-xs font-medium rounded-l-md border ${
                    specialistMode === 'flat'
                      ? 'bg-blue-600 text-white border-blue-600'
                      : 'bg-white text-gray-700 border-gray-300 hover:bg-gray-50'
                  }`}
                >
                  Flat
                </button>
                <button
                  type="button"
                  onClick={() => setSpecialistMode('percent')}
                  className={`px-3 py-1 text-xs font-medium rounded-r-md border-t border-r border-b ${
                    specialistMode === 'percent'
                      ? 'bg-blue-600 text-white border-blue-600'
                      : 'bg-white text-gray-700 border-gray-300 hover:bg-gray-50'
                  }`}
                >
                  Percent
                </button>
              </div>
            </div>

            {specialistMode === 'flat' ? (
              <Input
                type="number"
                value={formData.specialistHours || ''}
                onChange={(e) => setFormData({ ...formData, specialistHours: parseInt(e.target.value) || 0 })}
                placeholder="Enter specialist hours"
                min={0}
                max={formData.totalSoldHours - 1}
              />
            ) : (
              <div className="flex items-center space-x-2">
                <Input
                  type="number"
                  value={specialistPercent || ''}
                  onChange={(e) => setSpecialistPercent(parseInt(e.target.value) || 0)}
                  placeholder="Enter percentage"
                  min={0}
                  max={99}
                />
                <span className="text-gray-500">%</span>
                <span className="text-sm text-gray-500">
                  = {formData.specialistHours} hrs
                </span>
              </div>
            )}
            <p className="mt-1 text-xs text-gray-500">
              Your hours: {formData.totalSoldHours - formData.specialistHours} hrs
            </p>
          </div>

          <div className="flex items-center">
            <input
              type="checkbox"
              id="isActive"
              checked={formData.isActive}
              onChange={(e) => setFormData({ ...formData, isActive: e.target.checked })}
              className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
            />
            <label htmlFor="isActive" className="ml-2 block text-sm text-gray-700">
              Project is active
            </label>
          </div>

          <div className="p-3 bg-yellow-50 border border-yellow-200 rounded-md">
            <p className="text-sm text-yellow-800">
              <strong>Note:</strong> Editing hours or dates will recalculate future weeks only. Past weeks with logged hours will be preserved.
            </p>
          </div>

          <HourDistributionPreview entries={preview} isLoading={previewLoading} />

          {error && (
            <div className="p-3 bg-red-50 border border-red-200 rounded-md text-sm text-red-600">
              {error}
            </div>
          )}
        </div>

        <div className="mt-6 flex justify-end space-x-3">
          <Button type="button" variant="secondary" onClick={onClose}>
            Cancel
          </Button>
          <Button type="submit" disabled={isSubmitting || preview.length === 0}>
            {isSubmitting ? 'Saving...' : 'Save Changes'}
          </Button>
        </div>
      </form>
    </Modal>
  )
}
