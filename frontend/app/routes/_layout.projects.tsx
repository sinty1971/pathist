import ProjectGanttChartSimple from "../components/ProjectGanttChartSimple";

export default function GanttChart() {

  return (
    <div className="container mx-auto p-4">
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h2 className="card-title">工程表</h2>
          <ProjectGanttChartSimple />
        </div>
      </div>
    </div>
  );
}