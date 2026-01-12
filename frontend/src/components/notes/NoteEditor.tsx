import { useState, useEffect } from 'react'
import { Modal, Button, Input } from '../ui'
import type { Note, CreateNoteInput, UpdateNoteInput } from '../../types'

interface NoteEditorProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (input: CreateNoteInput | UpdateNoteInput) => Promise<void>
  projectId: number
  note?: Note | null
}

export function NoteEditor({
  isOpen,
  onClose,
  onSubmit,
  projectId,
  note,
}: NoteEditorProps) {
  const [title, setTitle] = useState('')
  const [content, setContent] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  const isEditing = !!note

  useEffect(() => {
    if (note) {
      setTitle(note.title)
      setContent(note.content)
    } else {
      setTitle('')
      setContent('')
    }
  }, [note, isOpen])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsSubmitting(true)

    try {
      if (isEditing && note) {
        await onSubmit({
          id: note.id,
          title,
          content,
        })
      } else {
        await onSubmit({
          projectId,
          title,
          content,
        })
      }
      onClose()
    } catch (err) {
      console.error('Failed to save note:', err)
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={isEditing ? 'Edit Note' : 'New Note'}
      size="xl"
    >
      <form onSubmit={handleSubmit} className="space-y-4">
        <Input
          label="Title"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          required
          placeholder="Note title"
        />

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Content (Markdown supported)
          </label>
          <textarea
            value={content}
            onChange={(e) => setContent(e.target.value)}
            rows={15}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent font-mono text-sm"
            placeholder="# Heading&#10;&#10;Write your notes here using **markdown** formatting.&#10;&#10;- List item 1&#10;- List item 2&#10;&#10;```code block```"
          />
          <p className="text-xs text-gray-500 mt-1">
            Supports headers, bold, italic, lists, code blocks, and more
          </p>
        </div>

        <div className="flex justify-end gap-3 pt-4">
          <Button type="button" variant="ghost" onClick={onClose}>
            Cancel
          </Button>
          <Button type="submit" disabled={isSubmitting}>
            {isSubmitting ? 'Saving...' : isEditing ? 'Update' : 'Create Note'}
          </Button>
        </div>
      </form>
    </Modal>
  )
}
