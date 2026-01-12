import { useEffect, useState } from 'react'
import ReactMarkdown from 'react-markdown'
import { Button, Card, CardHeader } from '../ui'
import { GoalModal } from './GoalModal'
import { useGoals } from '../../hooks'
import { formatDate } from '../../utils'
import type { Goal, GoalStatus, CreateGoalInput, UpdateGoalInput } from '../../types'

interface GoalListProps {
  projectId: number
}

const statusConfig: Record<GoalStatus, { label: string; color: string; bg: string }> = {
  pending: { label: 'Pending', color: 'text-gray-600', bg: 'bg-gray-100' },
  in_progress: { label: 'In Progress', color: 'text-blue-600', bg: 'bg-blue-100' },
  completed: { label: 'Completed', color: 'text-green-600', bg: 'bg-green-100' },
}

export function GoalList({ projectId }: GoalListProps) {
  const {
    goals,
    isLoading,
    loadGoals,
    createGoal,
    updateGoal,
    updateGoalStatus,
    deleteGoal,
  } = useGoals(projectId)

  const [isModalOpen, setIsModalOpen] = useState(false)
  const [editingGoal, setEditingGoal] = useState<Goal | null>(null)
  const [expandedId, setExpandedId] = useState<number | null>(null)

  useEffect(() => {
    loadGoals()
  }, [loadGoals])

  const handleSubmit = async (input: CreateGoalInput | UpdateGoalInput) => {
    if ('id' in input) {
      await updateGoal(input)
    } else {
      await createGoal(input)
    }
  }

  const handleEdit = (goal: Goal) => {
    setEditingGoal(goal)
    setIsModalOpen(true)
  }

  const handleDelete = async (id: number) => {
    if (window.confirm('Are you sure you want to delete this goal?')) {
      await deleteGoal(id)
    }
  }

  const handleStatusChange = async (id: number, status: GoalStatus) => {
    await updateGoalStatus(id, status)
  }

  const handleCloseModal = () => {
    setIsModalOpen(false)
    setEditingGoal(null)
  }

  const toggleExpand = (id: number) => {
    setExpandedId(expandedId === id ? null : id)
  }

  const completedCount = goals.filter((g) => g.status === 'completed').length
  const progressPercent = goals.length > 0 ? (completedCount / goals.length) * 100 : 0

  return (
    <Card>
      <CardHeader
        title="Goals"
        subtitle={`${completedCount}/${goals.length} completed`}
        action={
          <Button size="sm" onClick={() => setIsModalOpen(true)}>
            <svg className="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
            </svg>
            Add Goal
          </Button>
        }
      />

      {goals.length > 0 && (
        <div className="mb-4">
          <div className="flex justify-between text-sm text-gray-500 mb-1">
            <span>Progress</span>
            <span>{progressPercent.toFixed(0)}%</span>
          </div>
          <div className="w-full bg-gray-200 rounded-full h-2">
            <div
              className="bg-green-500 h-2 rounded-full transition-all"
              style={{ width: `${progressPercent}%` }}
            />
          </div>
        </div>
      )}

      {isLoading ? (
        <p className="text-gray-500">Loading goals...</p>
      ) : goals.length === 0 ? (
        <p className="text-gray-500 text-center py-8">No goals yet. Add your first goal!</p>
      ) : (
        <div className="space-y-3">
          {goals.map((goal) => {
            const config = statusConfig[goal.status]
            return (
              <div
                key={goal.id}
                className="border border-gray-200 rounded-lg overflow-hidden"
              >
                <div
                  className="flex items-center justify-between p-4 cursor-pointer hover:bg-gray-50"
                  onClick={() => toggleExpand(goal.id)}
                >
                  <div className="flex items-center gap-3">
                    <button
                      onClick={(e) => {
                        e.stopPropagation()
                        const nextStatus: GoalStatus =
                          goal.status === 'pending'
                            ? 'in_progress'
                            : goal.status === 'in_progress'
                            ? 'completed'
                            : 'pending'
                        handleStatusChange(goal.id, nextStatus)
                      }}
                      className={`w-6 h-6 rounded-full border-2 flex items-center justify-center transition-colors ${
                        goal.status === 'completed'
                          ? 'border-green-500 bg-green-500'
                          : goal.status === 'in_progress'
                          ? 'border-blue-500'
                          : 'border-gray-300'
                      }`}
                    >
                      {goal.status === 'completed' && (
                        <svg className="w-4 h-4 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
                        </svg>
                      )}
                      {goal.status === 'in_progress' && (
                        <div className="w-2 h-2 rounded-full bg-blue-500" />
                      )}
                    </button>
                    <div>
                      <h4 className={`font-medium ${goal.status === 'completed' ? 'text-gray-400 line-through' : 'text-gray-900'}`}>
                        {goal.title}
                      </h4>
                      <div className="flex items-center gap-2 mt-1">
                        <span className={`text-xs px-2 py-0.5 rounded-full ${config.bg} ${config.color}`}>
                          {config.label}
                        </span>
                        {goal.targetDate && (
                          <span className="text-xs text-gray-500">
                            Due: {formatDate(goal.targetDate)}
                          </span>
                        )}
                      </div>
                    </div>
                  </div>
                  <svg
                    className={`h-5 w-5 text-gray-400 transition-transform ${expandedId === goal.id ? 'rotate-180' : ''}`}
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                  </svg>
                </div>

                {expandedId === goal.id && (
                  <div className="px-4 pb-4 border-t border-gray-200 bg-gray-50">
                    {goal.description ? (
                      <div className="prose prose-sm max-w-none mt-3">
                        <ReactMarkdown>{goal.description}</ReactMarkdown>
                      </div>
                    ) : (
                      <p className="text-gray-500 text-sm mt-3 italic">No description</p>
                    )}
                    <div className="flex gap-2 mt-4">
                      <Button size="sm" variant="ghost" onClick={() => handleEdit(goal)}>
                        Edit
                      </Button>
                      <Button size="sm" variant="ghost" onClick={() => handleDelete(goal.id)}>
                        Delete
                      </Button>
                    </div>
                  </div>
                )}
              </div>
            )
          })}
        </div>
      )}

      <GoalModal
        isOpen={isModalOpen}
        onClose={handleCloseModal}
        onSubmit={handleSubmit}
        projectId={projectId}
        goal={editingGoal}
      />
    </Card>
  )
}
