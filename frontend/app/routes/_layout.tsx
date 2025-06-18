import { Outlet } from "react-router";
import { Navigation } from "../components/Navigation";
import "../styles/App.css";

export default function Layout() {
  return (
    <div className="app">
      <Navigation />
      <main className="main-content">
        <Outlet />
      </main>
    </div>
  );
}
