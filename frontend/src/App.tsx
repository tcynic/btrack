import { useState } from "react";
import { AppProvider } from "./context/AppContext";
import { Layout } from "./components/layout";
import { Dashboard } from "./components/dashboard";
import { ProjectList, ProjectDetail } from "./components/projects";
import type { ProjectWithStats } from "./types";

function AppContent() {
  const [activeTab, setActiveTab] = useState<"dashboard" | "projects">(
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
