import type { RouteConfig } from "@react-router/dev/routes";
import { route, layout, index } from "@react-router/dev/routes";

export default [
  layout("routes/_layout.tsx", [
    index("routes/_layout._index.tsx"),
    route("files", "routes/_layout.files.tsx"),
    route("projects", "routes/_layout.projects.tsx"),
    route("projects/gantt", "routes/_layout.gantt.tsx"),
  ]),
] satisfies RouteConfig;