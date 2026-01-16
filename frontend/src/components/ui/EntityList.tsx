import { ReactNode } from 'react'
import { Card, CardHeader, Button } from './'

interface EntityListProps<T> {
  title: string
  items: T[]
  isLoading: boolean
  emptyMessage: string
  addButtonLabel: string
  onAdd: () => void
  renderItem: (item: T) => ReactNode
  subtitle?: string | ((items: T[]) => string)
  headerAction?: ReactNode
  layout?: 'list' | 'grid'
  gridCols?: 1 | 2 | 3
}

export function EntityList<T>({
  title,
  items,
  isLoading,
  emptyMessage,
  addButtonLabel,
  onAdd,
  renderItem,
  subtitle,
  headerAction,
  layout = 'list',
  gridCols = 1,
}: EntityListProps<T>) {
  const subtitleText = typeof subtitle === 'function' ? subtitle(items) : subtitle

  const defaultAction = (
    <Button size="sm" onClick={onAdd}>
      <svg className="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
      </svg>
      {addButtonLabel}
    </Button>
  )

  const gridClassName = layout === 'grid' 
    ? `grid gap-3 ${gridCols === 1 ? 'grid-cols-1' : gridCols === 2 ? 'grid-cols-1 md:grid-cols-2' : 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3'}`
    : 'space-y-3'

  return (
    <Card>
      <CardHeader
        title={title}
        subtitle={subtitleText}
        action={headerAction || defaultAction}
      />

      {isLoading ? (
        <p className="text-gray-500">Loading {title.toLowerCase()}...</p>
      ) : items.length === 0 ? (
        <p className="text-gray-500 text-center py-8">{emptyMessage}</p>
      ) : (
        <div className={gridClassName}>
          {items.map((item, index) => (
            <div key={index}>{renderItem(item)}</div>
          ))}
        </div>
      )}
    </Card>
  )
}

// Utility function to generate default subtitle
export function itemCountSubtitle(count: number, singular: string, plural?: string): string {
  const word = count === 1 ? singular : (plural || `${singular}s`)
  return `${count} ${word}`
}
