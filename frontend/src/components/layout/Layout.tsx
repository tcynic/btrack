import { ReactNode } from 'react'
import { Header } from './Header'

interface LayoutProps {
  children: ReactNode
  activeTab: 'dashboard' | 'week' | 'projects' | 'tasks' | 'reports' | 'settings'
  onTabChange: (tab: 'dashboard' | 'week' | 'projects' | 'tasks' | 'reports' | 'settings') => void
}

export function Layout({ children, activeTab, onTabChange }: LayoutProps) {
  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <Header activeTab={activeTab} onTabChange={onTabChange} />
      <main className="mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {children}
      </main>
    </div>
  )
}
