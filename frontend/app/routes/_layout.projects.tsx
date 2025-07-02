import Projects from "../components/Projects";

export default function ProjectsPage() {

  return (
    <div className="container mx-auto p-4">
      <div className="bg-white border border-gray-200 rounded-lg">
        <div className="p-6">
          <Projects />
        </div>
      </div>
    </div>
  );
}