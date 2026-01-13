import { useState } from "react";
import { AppProvider } from "./context/AppContext";
import { Layout } from "./components/layout";
import { Dashboard } from "./components/dashboard";
import { WeekView } from "./components/week";
import { ProjectList, ProjectDetail } from "./components/projects";
import { ReportsView } from "./components/reports";
import { SettingsView } from "./components/settings";
import type { ProjectWithStats } from "./types";

function AppContent() {
  const [activeTab, setActiveTab] = useState<"dashboard" | "week" | "projects" | "reports" | "settings">(
    "dashboard",
  );
  const [selectedProject, setSelectedProject] =
    useState<ProjectWithStats | null>(null);

  const handleProjectSelect = (project: ProjectWithStats) => {
    setSelectedProject(project);
  };

  const handleBackToList = () => {
    setSelectedProject(null);
  };

  return (
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
