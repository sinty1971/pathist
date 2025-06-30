import ProjectGanttChart from "../components/ProjectGanttChart";

export default function GanttChartPage() {
  return (
    <div className="container mx-auto p-4">
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <ProjectGanttChart />
        </div>
      </div>
    </div>
  );
}