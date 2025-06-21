import type { RouteConfig } from "@react-router/dev/routes";
import { route, layout, index } from "@react-router/dev/routes";

export default [
  layout("routes/_layout.tsx", [
    index("routes/_layout._index.tsx"),
    route("kouji", "routes/_layout.kouji.tsx"),
    route("gantt", "routes/_layout.gantt.tsx"),
  ]),
] satisfies RouteConfig;