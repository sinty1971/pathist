import type { RouteConfig } from "@react-router/dev/routes";
import { route, layout, index } from "@react-router/dev/routes";

export default [
  layout("routes/_layout.tsx", [
    index("routes/_layout._index.tsx"),
    route("files", "routes/_layout.files.tsx"),
    route("kojies", "routes/_layout.kojies.tsx"),
    route("kojies/gantt", "routes/_layout.gantt.tsx"),
    route("companies", "routes/_layout.companies.tsx"),
  ]),
] satisfies RouteConfig;