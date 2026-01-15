import { useState, useEffect } from 'react'
import ReactMarkdown from 'react-markdown'
import { Modal, Button } from '../ui'
import { TaskList, TaskModal } from '../tasks'
import { GetTasksBySource, CreateTask, UpdateTask, DeleteTask, UpdateTaskStatus } from '../../../wailsjs/go/main/App'
import type { Note, Task, CreateTaskInput, UpdateTaskInput } from '../../types'

interface NoteViewerProps {
  isOpen: boolean
  onClose: () => void
  onEdit: () => void
  note: Note | null
}

export function NoteViewer({ isOpen, onClose, onEdit, note }: NoteViewerProps) {
  const [tasks, setTasks] = useState<Task[]>([])
  const [selectedTask, setSelectedTask] = useState<Task | null>(null)
  const [isTaskModalOpen, setIsTaskModalOpen] = useState(false)

  if (!note) return null

  const loadTasks = async () => {
    if (note) {
      try {
        const noteTasks = await GetTasksBySource('note', note.id)
        setTasks(noteTasks)
      } catch (err) {
        console.error('Failed to load tasks:', err)
      }
    }
  }

  useEffect(() => {
    if (isOpen && note) {
      loadTasks()
    }
  }, [isOpen, note])

  const handleAddTask = () => {
    setSelectedTask(null)
    setIsTaskModalOpen(true)
  }

  const handleEditTask = (task: Task) => {
    setSelectedTask(task)
    setIsTaskModalOpen(true)
  }

  const handleSubmitTask = async (input: CreateTaskInput | UpdateTaskInput) => {
    if ('id' in input) {
      await UpdateTask(input)
    } else {
      await CreateTask(input)
    }
    await loadTasks()
    setIsTaskModalOpen(false)
  }

  const handleDeleteTask = async (id: number) => {
    if (confirm('Are you sure you want to delete this task?')) {
      await DeleteTask(id)
      await loadTasks()
    }
  }

  const handleStatusChange = async (id: number, status: string) => {
    await UpdateTaskStatus(id, status)
    await loadTasks()
  }

  return (
    <Modal isOpen={isOpen} onClose={onClose} title={note.title} size="xl">
      <div className="prose prose-sm max-w-none mb-6">
        {note.content ? (
          <ReactMarkdown>{note.content}</ReactMarkdown>
        ) : (
          <p className="text-gray-500 italic">No content</p>
        )}
      </div>

      {/* Tasks Section */}
      <div className="border-t border-gray-200 pt-6">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-medium text-gray-900">Tasks</h3>
          <Button size="sm" onClick={handleAddTask}>Add Task</Button>
        </div>
        <TaskList
          tasks={tasks}
          onEdit={handleEditTask}
          onDelete={handleDeleteTask}
          onStatusChange={handleStatusChange}
          emptyMessage="No tasks for this note"
        />
      </div>

      <div className="flex justify-end gap-3 pt-6 mt-6 border-t border-gray-200">
        <Button variant="ghost" onClick={onClose}>
          Close
        </Button>
        <Button onClick={onEdit}>Edit Note</Button>
      </div>

      {/* Task Modal */}
      {isTaskModalOpen && (
        <TaskModal
          isOpen={isTaskModalOpen}
          onClose={() => setIsTaskModalOpen(false)}
          onSubmit={handleSubmitTask}
          projectId={note.projectId}
          sourceType="note"
          sourceId={note.id}
          task={selectedTask}
        />
      )}
    </Modal>
  )
}
