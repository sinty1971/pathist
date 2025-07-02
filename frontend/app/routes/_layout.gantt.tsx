import ProjectGanttChart from "../components/ProjectGanttChart";

export default function GanttChartPage() {
  return (
    <div className="container mx-auto p-4">
      <div className="bg-white border border-gray-200 rounded-lg">
        <div className="p-6">
          <ProjectGanttChart />
        </div>
      </div>
    </div>
  );
}