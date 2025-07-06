import { Outlet } from "react-router";
import type { LinksFunction } from "react-router";
import { Navigation } from "../components/Navigation";
import { KojiProvider, useKoji } from "../contexts/KojiContext";
import { FileInfoProvider } from "../contexts/FileInfoContext";

export const links: LinksFunction = () => [
  { rel: "stylesheet", href: "/app/styles/app.css" },
];

function LayoutContent() {
  const { kojiCount } = useKoji();
  
  return (
    <div className="app">
      <Navigation />
      <main className="main-content">
        <Outlet />
      </main>
    </div>
  );
}

export default function Layout() {
  return (
    <KojiProvider>
      <FileInfoProvider>
        <LayoutContent />
      </FileInfoProvider>
    </KojiProvider>
  );
}
