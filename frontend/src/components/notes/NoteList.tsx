import { useEffect, useState } from 'react'
import { EntityList, itemCountSubtitle } from '../ui'
import { NoteEditor } from './NoteEditor'
import { NoteViewer } from './NoteViewer'
import { useNotes } from '../../hooks'
import { formatDate } from '../../utils'
import type { Note, CreateNoteInput, UpdateNoteInput } from '../../types'

interface NoteListProps {
  projectId: number
}

export function NoteList({ projectId }: NoteListProps) {
  const { notes, isLoading, loadNotes, createNote, updateNote, deleteNote } =
    useNotes(projectId)

  const [isEditorOpen, setIsEditorOpen] = useState(false)
  const [isViewerOpen, setIsViewerOpen] = useState(false)
  const [selectedNote, setSelectedNote] = useState<Note | null>(null)

  useEffect(() => {
    loadNotes()
  }, [loadNotes])

  const handleSubmit = async (input: CreateNoteInput | UpdateNoteInput) => {
    if ('id' in input) {
      await updateNote(input)
    } else {
      await createNote(input)
    }
  }

  const handleView = (note: Note) => {
    setSelectedNote(note)
    setIsViewerOpen(true)
  }

  const handleEdit = (note: Note) => {
    setSelectedNote(note)
    setIsEditorOpen(true)
    setIsViewerOpen(false)
  }

  const handleDelete = async (id: number) => {
    if (window.confirm('Are you sure you want to delete this note?')) {
      await deleteNote(id)
    }
  }

  const handleCloseEditor = () => {
    setIsEditorOpen(false)
    setSelectedNote(null)
  }

  const handleCloseViewer = () => {
    setIsViewerOpen(false)
    setSelectedNote(null)
  }

  const handleNewNote = () => {
    setSelectedNote(null)
    setIsEditorOpen(true)
  }

  const renderNoteItem = (note: Note) => (
    <div
      className="border border-gray-200 rounded-lg p-4 hover:border-gray-300 cursor-pointer transition-colors"
      onClick={() => handleView(note)}
    >
      <div className="flex items-start justify-between">
        <div className="flex items-center gap-2">
          <svg
            className="h-5 w-5 text-yellow-500 flex-shrink-0"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
            />
          </svg>
          <h4 className="font-medium text-gray-900 truncate">
            {note.title}
          </h4>
        </div>
        <button
          onClick={(e) => {
            e.stopPropagation()
            handleDelete(note.id)
          }}
          className="text-gray-400 hover:text-red-500 p-1"
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
              d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
            />
          </svg>
        </button>
      </div>
      <p className="text-sm text-gray-500 mt-2 line-clamp-2">
        {note.content || 'No content'}
      </p>
      <p className="text-xs text-gray-400 mt-2">
        Updated {formatDate(note.updatedAt)}
      </p>
    </div>
  )

  return (
    <>
      <EntityList
        title="Notes"
        items={notes}
        isLoading={isLoading}
        emptyMessage="No notes yet. Create your first note!"
        addButtonLabel="New Note"
        onAdd={handleNewNote}
        renderItem={renderNoteItem}
        subtitle={(items) => itemCountSubtitle(items.length, "note")}
        layout="grid"
        gridCols={2}
      />

      <NoteEditor
        isOpen={isEditorOpen}
        onClose={handleCloseEditor}
        onSubmit={handleSubmit}
        projectId={projectId}
        note={selectedNote}
      />

      <NoteViewer
        isOpen={isViewerOpen}
        onClose={handleCloseViewer}
        onEdit={() => selectedNote && handleEdit(selectedNote)}
        note={selectedNote}
      />
    </>
  )
}
