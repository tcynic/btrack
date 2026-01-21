interface HeaderProps {
  activeTab: 'dashboard' | 'week' | 'projects' | 'tasks' | 'reports' | 'settings'
  onTabChange: (tab: 'dashboard' | 'week' | 'projects' | 'tasks' | 'reports' | 'settings') => void
}

export function Header({ activeTab, onTabChange }: HeaderProps) {
  const baseTab = 'px-4 py-2 rounded-md text-sm font-medium transition-colors';
  const active = 'bg-ld-primary text-white';
  const inactive = 'text-ld-text hover:bg-ld-surface2';

  return (
    <header className="bg-ld-surface border-b border-ld-border">
      <div className="mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center py-4">
          <div className="flex items-center space-x-3">
            <svg
              className="h-8 w-8 text-ld-primary"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
              />
            </svg>
            <h1 className="text-xl font-bold text-ld-text">Bandwidth Tracker</h1>
            <span className="text-xs text-ld-muted ml-4">
              Press <kbd className="px-1 py-0.5 text-xs font-semibold text-ld-text bg-ld-surface2 border border-ld-border rounded">?</kbd> for shortcuts
            </span>
          </div>
          <nav className="flex space-x-1">
            <button onClick={() => onTabChange('dashboard')} className={`${baseTab} ${activeTab === 'dashboard' ? active : inactive}`}>Dashboard</button>
            <button onClick={() => onTabChange('week')} className={`${baseTab} ${activeTab === 'week' ? active : inactive}`}>Week</button>
            <button onClick={() => onTabChange('projects')} className={`${baseTab} ${activeTab === 'projects' ? active : inactive}`}>Projects</button>
            <button onClick={() => onTabChange('tasks')} className={`${baseTab} ${activeTab === 'tasks' ? active : inactive}`}>Tasks</button>
            <button onClick={() => onTabChange('reports')} className={`${baseTab} ${activeTab === 'reports' ? active : inactive}`}>Reports</button>
            <button onClick={() => onTabChange('settings')} className={`${baseTab} ${activeTab === 'settings' ? active : inactive}`}>Settings</button>
          </nav>
        </div>
      </div>
    </header>
  )
}
