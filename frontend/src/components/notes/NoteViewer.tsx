import ReactMarkdown from 'react-markdown'
import { Modal, Button } from '../ui'
import type { Note } from '../../types'

interface NoteViewerProps {
  isOpen: boolean
  onClose: () => void
  onEdit: () => void
  note: Note | null
}

export function NoteViewer({ isOpen, onClose, onEdit, note }: NoteViewerProps) {
  if (!note) return null

  return (
    <Modal isOpen={isOpen} onClose={onClose} title={note.title} size="xl">
      <div className="prose prose-sm max-w-none">
        {note.content ? (
          <ReactMarkdown>{note.content}</ReactMarkdown>
        ) : (
          <p className="text-gray-500 italic">No content</p>
        )}
      </div>
      <div className="flex justify-end gap-3 pt-6 mt-6 border-t border-gray-200">
        <Button variant="ghost" onClick={onClose}>
          Close
        </Button>
        <Button onClick={onEdit}>Edit Note</Button>
      </div>
    </Modal>
  )
}
