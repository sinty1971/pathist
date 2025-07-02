import { Outlet } from "react-router";
import type { LinksFunction } from "react-router";
import { Navigation } from "../components/Navigation";
import { ProjectProvider, useProject } from "../contexts/ProjectContext";
import { FileInfoProvider } from "../contexts/FileInfoContext";

export const links: LinksFunction = () => [
  { rel: "stylesheet", href: "/app/styles/App.css" },
];

function LayoutContent() {
  const { projectCount } = useProject();
  
  return (
    <div className="app">
      <Navigation projectCount={projectCount} />
      <main className="main-content">
        <Outlet />
      </main>
    </div>
  );
}

export default function Layout() {
  return (
    <ProjectProvider>
      <FileInfoProvider>
        <LayoutContent />
      </FileInfoProvider>
    </ProjectProvider>
  );
}
