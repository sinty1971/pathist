import { createContext, useContext, useState } from 'react';
import type { ReactNode } from 'react';

interface ProjectContextType {
  projectCount: number;
  setProjectCount: (count: number) => void;
}

const ProjectContext = createContext<ProjectContextType | undefined>(undefined);

export function ProjectProvider({ children }: { children: ReactNode }) {
  const [projectCount, setProjectCount] = useState<number>(0);

  return (
    <ProjectContext.Provider value={{ projectCount, setProjectCount }}>
      {children}
    </ProjectContext.Provider>
  );
}

export function useProject() {
  const context = useContext(ProjectContext);
  if (context === undefined) {
    throw new Error('useProject must be used within a ProjectProvider');
  }
  return context;
}