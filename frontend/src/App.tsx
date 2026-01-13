import { useState } from "react";
import { AppProvider } from "./context/AppContext";
import { Layout } from "./components/layout";
import { Dashboard } from "./components/dashboard";
import { WeekView } from "./components/week";
import { ProjectList, ProjectDetail } from "./components/projects";
import { ReportsView } from "./components/reports";
import { SettingsView } from "./components/settings";
import { KeyboardShortcutsModal } from "./components/ui";
import { useKeyboardShortcuts } from "./hooks/useKeyboardShortcuts";
import type { ProjectWithStats } from "./types";

function AppContent() {
  const [activeTab, setActiveTab] = useState<"dashboard" | "week" | "projects" | "reports" | "settings">(
    "dashboard",
  );
  const [selectedProject, setSelectedProject] =
    useState<ProjectWithStats | null>(null);

  const { showHelp, setShowHelp } = useKeyboardShortcuts({
    onTabChange: setActiveTab,
  });

  const handleProjectSelect = (project: ProjectWithStats) => {
    setSelectedProject(project);
  };

  const handleBackToList = () => {
    setSelectedProject(null);
  };

  return (
    <>
      <Layout activeTab={activeTab} onTabChange={setActiveTab}>
        {activeTab === "dashboard" ? (
          <Dashboard />
        ) : activeTab === "week" ? (
          <WeekView />
        ) : activeTab === "reports" ? (
          <ReportsView />
        ) : activeTab === "settings" ? (
          <SettingsView />
        ) : selectedProject ? (
          <ProjectDetail project={selectedProject} onBack={handleBackToList} />
        ) : (
          <ProjectList onProjectSelect={handleProjectSelect} />
        )}
      </Layout>
      <KeyboardShortcutsModal isOpen={showHelp} onClose={() => setShowHelp(false)} />
    </>
  );
}

function App() {
  return (
    <AppProvider>
      <AppContent />
    </AppProvider>
  );
}

export default App;
