import { jsx, jsxs } from "react/jsx-runtime";
import { PassThrough } from "node:stream";
import { createReadableStreamFromReadable } from "@react-router/node";
import { ServerRouter, UNSAFE_withComponentProps, Outlet, Meta, Links, ScrollRestoration, Scripts, useLocation, Link } from "react-router";
import { isbot } from "isbot";
import { renderToPipeableStream } from "react-dom/server";
import React, { createContext, useState, useContext, useCallback, useEffect, useRef } from "react";
import { TreeItem, SimpleTreeView } from "@mui/x-tree-view";
import { Box, Typography, CircularProgress, Chip, Paper, Toolbar, IconButton, TextField, Breadcrumbs, Link as Link$1, Alert, Dialog, DialogTitle, DialogContent, Button } from "@mui/material";
import { FolderOpen, Folder, InsertDriveFile, Home, Refresh, Search } from "@mui/icons-material";
import CloseIcon from "@mui/icons-material/Close";
import WarningIcon from "@mui/icons-material/Warning";
import AttachFileIcon from "@mui/icons-material/AttachFile";
import LightbulbIcon from "@mui/icons-material/Lightbulb";
import { DatePicker } from "@mui/x-date-pickers/DatePicker";
import { LocalizationProvider } from "@mui/x-date-pickers/LocalizationProvider";
import { AdapterDateFns } from "@mui/x-date-pickers/AdapterDateFns";
import { ja } from "date-fns/locale";
const streamTimeout = 5e3;
function handleRequest(request, responseStatusCode, responseHeaders, routerContext, loadContext) {
  return new Promise((resolve, reject) => {
    let shellRendered = false;
    let userAgent = request.headers.get("user-agent");
    let readyOption = userAgent && isbot(userAgent) || routerContext.isSpaMode ? "onAllReady" : "onShellReady";
    const { pipe, abort } = renderToPipeableStream(
      /* @__PURE__ */ jsx(ServerRouter, { context: routerContext, url: request.url }),
      {
        [readyOption]() {
          shellRendered = true;
          const body = new PassThrough();
          const stream = createReadableStreamFromReadable(body);
          responseHeaders.set("Content-Type", "text/html");
          resolve(
            new Response(stream, {
              headers: responseHeaders,
              status: responseStatusCode
            })
          );
          pipe(body);
        },
        onShellError(error) {
          reject(error);
        },
        onError(error) {
          responseStatusCode = 500;
          if (shellRendered) {
            console.error(error);
          }
        }
      }
    );
    setTimeout(abort, streamTimeout + 1e3);
  });
}
const entryServer = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
  __proto__: null,
  default: handleRequest,
  streamTimeout
}, Symbol.toStringTag, { value: "Module" }));
const links$1 = () => [{
  rel: "preconnect",
  href: "https://fonts.googleapis.com"
}, {
  rel: "preconnect",
  href: "https://fonts.gstatic.com",
  crossOrigin: "anonymous"
}, {
  rel: "stylesheet",
  href: "https://fonts.googleapis.com/css2?family=Inter:ital,opsz,wght@0,14..32,100..900;1,14..32,100..900&display=swap"
}];
function Layout({
  children
}) {
  return /* @__PURE__ */ jsxs("html", {
    lang: "ja",
    children: [/* @__PURE__ */ jsxs("head", {
      children: [/* @__PURE__ */ jsx("meta", {
        charSet: "utf-8"
      }), /* @__PURE__ */ jsx("meta", {
        name: "viewport",
        content: "width=device-width, initial-scale=1"
      }), /* @__PURE__ */ jsx(Meta, {}), /* @__PURE__ */ jsx(Links, {})]
    }), /* @__PURE__ */ jsxs("body", {
      children: [children, /* @__PURE__ */ jsx(ScrollRestoration, {}), /* @__PURE__ */ jsx(Scripts, {})]
    })]
  });
}
const root = UNSAFE_withComponentProps(function App() {
  return /* @__PURE__ */ jsx(Outlet, {});
});
const route0 = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
  __proto__: null,
  Layout,
  default: root,
  links: links$1
}, Symbol.toStringTag, { value: "Module" }));
const FileInfoContext = createContext(void 0);
const FileInfoProvider = ({ children }) => {
  const [fileCount, setFileCount] = useState(0);
  const [currentPath, setCurrentPath] = useState("");
  return /* @__PURE__ */ jsx(FileInfoContext.Provider, { value: { fileCount, currentPath, setFileCount, setCurrentPath }, children });
};
const useFileInfo = () => {
  const context = useContext(FileInfoContext);
  if (!context) {
    throw new Error("useFileInfo must be used within a FileInfoProvider");
  }
  return context;
};
function Navigation({ projectCount } = {}) {
  const location = useLocation();
  const { fileCount, currentPath } = useFileInfo();
  const getPageTitle = () => {
    switch (location.pathname) {
      case "/":
        return "ホーム";
      case "/files":
        return `ファイル一覧${fileCount > 0 ? ` (${fileCount}項目)` : ""}${currentPath ? ` - ${currentPath}` : ""}`;
      case "/projects":
        return `工程表${projectCount !== void 0 ? ` (${projectCount}件)` : ""}`;
      case "/projects/gantt":
        return "ガントチャート";
      default:
        return "ファイル管理";
    }
  };
  return /* @__PURE__ */ jsx("nav", { className: "navigation", children: /* @__PURE__ */ jsxs("div", { className: "nav-container", children: [
    /* @__PURE__ */ jsxs("div", { className: "nav-logo", children: [
      /* @__PURE__ */ jsx("h1", { children: "Penguin フォルダー管理" }),
      /* @__PURE__ */ jsx("div", { className: "page-title", children: getPageTitle() })
    ] }),
    /* @__PURE__ */ jsxs("div", { className: "nav-links", children: [
      /* @__PURE__ */ jsx(
        Link,
        {
          to: "/",
          className: location.pathname === "/" ? "nav-link active" : "nav-link",
          children: "ホーム"
        }
      ),
      /* @__PURE__ */ jsx(
        Link,
        {
          to: "/files",
          className: location.pathname === "/files" ? "nav-link active" : "nav-link",
          children: "ファイル一覧"
        }
      ),
      /* @__PURE__ */ jsx(
        Link,
        {
          to: "/projects",
          className: location.pathname === "/projects" ? "nav-link active" : "nav-link",
          children: "工程表"
        }
      )
    ] })
  ] }) });
}
const ProjectContext = createContext(void 0);
function ProjectProvider({ children }) {
  const [projectCount, setProjectCount] = useState(0);
  return /* @__PURE__ */ jsx(ProjectContext.Provider, { value: { projectCount, setProjectCount }, children });
}
function useProject() {
  const context = useContext(ProjectContext);
  if (context === void 0) {
    throw new Error("useProject must be used within a ProjectProvider");
  }
  return context;
}
const links = () => [{
  rel: "stylesheet",
  href: "/app/styles/App.css"
}];
function LayoutContent() {
  const {
    projectCount
  } = useProject();
  return /* @__PURE__ */ jsxs("div", {
    className: "app",
    children: [/* @__PURE__ */ jsx(Navigation, {
      projectCount
    }), /* @__PURE__ */ jsx("main", {
      className: "main-content",
      children: /* @__PURE__ */ jsx(Outlet, {})
    })]
  });
}
const _layout = UNSAFE_withComponentProps(function Layout2() {
  return /* @__PURE__ */ jsx(ProjectProvider, {
    children: /* @__PURE__ */ jsx(FileInfoProvider, {
      children: /* @__PURE__ */ jsx(LayoutContent, {})
    })
  });
});
const route1 = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
  __proto__: null,
  default: _layout,
  links
}, Symbol.toStringTag, { value: "Module" }));
const _layout__index = UNSAFE_withComponentProps(function HomePage() {
  return /* @__PURE__ */ jsx("div", {
    style: {
      height: "100vh",
      overflowY: "auto"
    },
    children: /* @__PURE__ */ jsxs("div", {
      style: {
        padding: "3rem",
        maxWidth: "1200px",
        margin: "0 auto"
      },
      children: [/* @__PURE__ */ jsxs("div", {
        style: {
          textAlign: "center",
          marginBottom: "3rem"
        },
        children: [/* @__PURE__ */ jsx("h1", {
          style: {
            fontSize: "3rem",
            color: "#333",
            marginBottom: "1rem",
            fontWeight: "600"
          },
          children: "Penguin フォルダー管理システム"
        }), /* @__PURE__ */ jsx("p", {
          style: {
            fontSize: "1.2rem",
            color: "#666",
            fontWeight: "400"
          },
          children: "ファイル管理と工事プロジェクト管理を効率的に"
        })]
      }), /* @__PURE__ */ jsxs("div", {
        style: {
          display: "grid",
          gridTemplateColumns: "repeat(auto-fill, minmax(200px, 1fr))",
          gap: "1rem",
          marginTop: "2rem"
        },
        children: [/* @__PURE__ */ jsx(Link, {
          to: "/files",
          style: {
            textDecoration: "none"
          },
          children: /* @__PURE__ */ jsxs("div", {
            className: "feature-card",
            style: {
              background: "rgba(255, 255, 255, 0.95)",
              borderRadius: "12px",
              padding: "1.5rem",
              boxShadow: "0 4px 12px rgba(0, 0, 0, 0.1)",
              transition: "all 0.2s ease",
              cursor: "pointer",
              height: "180px",
              display: "flex",
              flexDirection: "column",
              alignItems: "center",
              justifyContent: "center"
            },
            onMouseEnter: (e) => {
              e.currentTarget.style.transform = "translateY(-3px)";
              e.currentTarget.style.boxShadow = "0 8px 20px rgba(0, 0, 0, 0.15)";
            },
            onMouseLeave: (e) => {
              e.currentTarget.style.transform = "translateY(0)";
              e.currentTarget.style.boxShadow = "0 4px 12px rgba(0, 0, 0, 0.1)";
            },
            children: [/* @__PURE__ */ jsx("div", {
              style: {
                fontSize: "2.5rem",
                marginBottom: "0.5rem"
              },
              children: "📁"
            }), /* @__PURE__ */ jsx("h3", {
              style: {
                fontSize: "1.1rem",
                color: "#333",
                marginBottom: "0.5rem",
                textAlign: "center"
              },
              children: "ファイル一覧"
            }), /* @__PURE__ */ jsx("p", {
              style: {
                color: "#666",
                fontSize: "0.85rem",
                textAlign: "center",
                lineHeight: "1.4"
              },
              children: "ファイルとフォルダーの閲覧・管理"
            })]
          })
        }), /* @__PURE__ */ jsxs("div", {
          className: "feature-card",
          style: {
            background: "rgba(255, 255, 255, 0.95)",
            borderRadius: "12px",
            padding: "1.5rem",
            boxShadow: "0 4px 12px rgba(0, 0, 0, 0.1)",
            transition: "all 0.2s ease",
            height: "220px",
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            justifyContent: "center",
            position: "relative"
          },
          onMouseEnter: (e) => {
            e.currentTarget.style.transform = "translateY(-3px)";
            e.currentTarget.style.boxShadow = "0 8px 20px rgba(0, 0, 0, 0.15)";
          },
          onMouseLeave: (e) => {
            e.currentTarget.style.transform = "translateY(0)";
            e.currentTarget.style.boxShadow = "0 4px 12px rgba(0, 0, 0, 0.1)";
          },
          children: [/* @__PURE__ */ jsxs(Link, {
            to: "/projects",
            style: {
              textDecoration: "none",
              display: "flex",
              flexDirection: "column",
              alignItems: "center",
              justifyContent: "center",
              flex: 1,
              width: "100%"
            },
            children: [/* @__PURE__ */ jsx("div", {
              style: {
                fontSize: "2.5rem",
                marginBottom: "0.5rem"
              },
              children: "📊"
            }), /* @__PURE__ */ jsx("h3", {
              style: {
                fontSize: "1.1rem",
                color: "#333",
                marginBottom: "0.5rem",
                textAlign: "center"
              },
              children: "工程表"
            }), /* @__PURE__ */ jsx("p", {
              style: {
                color: "#666",
                fontSize: "0.85rem",
                textAlign: "center",
                lineHeight: "1.4",
                marginBottom: "0.5rem"
              },
              children: "工事プロジェクトの一覧管理"
            })]
          }), /* @__PURE__ */ jsx(Link, {
            to: "/projects/gantt",
            style: {
              color: "#667eea",
              fontSize: "0.8rem",
              textDecoration: "none",
              cursor: "pointer",
              transition: "all 0.2s",
              position: "absolute",
              bottom: "12px",
              left: "50%",
              transform: "translateX(-50%)",
              display: "flex",
              alignItems: "center",
              gap: "4px",
              borderBottom: "1px solid transparent"
            },
            onMouseEnter: (e) => {
              e.currentTarget.style.color = "#4f46e5";
              e.currentTarget.style.borderBottomColor = "#4f46e5";
            },
            onMouseLeave: (e) => {
              e.currentTarget.style.color = "#667eea";
              e.currentTarget.style.borderBottomColor = "transparent";
            },
            onClick: (e) => {
              e.stopPropagation();
            },
            children: "📈 ガントチャートを見る"
          })]
        })]
      }), /* @__PURE__ */ jsxs("div", {
        style: {
          marginTop: "4rem",
          padding: "2rem",
          background: "#f8f9fa",
          borderRadius: "12px",
          textAlign: "center",
          border: "1px solid #e9ecef"
        },
        children: [/* @__PURE__ */ jsx("h3", {
          style: {
            color: "#333",
            marginBottom: "1rem",
            fontWeight: "600"
          },
          children: "システム情報"
        }), /* @__PURE__ */ jsxs("p", {
          style: {
            color: "#666",
            lineHeight: "1.8"
          },
          children: ["データ保存場所: ~/penguin", /* @__PURE__ */ jsx("br", {}), "工事プロジェクト: ~/penguin/豊田築炉/2-工事"]
        })]
      })]
    })
  });
});
const route2 = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
  __proto__: null,
  default: _layout__index
}, Symbol.toStringTag, { value: "Module" }));
const jsonBodySerializer = {
  bodySerializer: (body) => JSON.stringify(
    body,
    (_key, value) => typeof value === "bigint" ? value.toString() : value
  )
};
const getAuthToken = async (auth, callback) => {
  const token = typeof callback === "function" ? await callback(auth) : callback;
  if (!token) {
    return;
  }
  if (auth.scheme === "bearer") {
    return `Bearer ${token}`;
  }
  if (auth.scheme === "basic") {
    return `Basic ${btoa(token)}`;
  }
  return token;
};
const separatorArrayExplode = (style) => {
  switch (style) {
    case "label":
      return ".";
    case "matrix":
      return ";";
    case "simple":
      return ",";
    default:
      return "&";
  }
};
const separatorArrayNoExplode = (style) => {
  switch (style) {
    case "form":
      return ",";
    case "pipeDelimited":
      return "|";
    case "spaceDelimited":
      return "%20";
    default:
      return ",";
  }
};
const separatorObjectExplode = (style) => {
  switch (style) {
    case "label":
      return ".";
    case "matrix":
      return ";";
    case "simple":
      return ",";
    default:
      return "&";
  }
};
const serializeArrayParam = ({
  allowReserved,
  explode,
  name,
  style,
  value
}) => {
  if (!explode) {
    const joinedValues2 = (allowReserved ? value : value.map((v) => encodeURIComponent(v))).join(separatorArrayNoExplode(style));
    switch (style) {
      case "label":
        return `.${joinedValues2}`;
      case "matrix":
        return `;${name}=${joinedValues2}`;
      case "simple":
        return joinedValues2;
      default:
        return `${name}=${joinedValues2}`;
    }
  }
  const separator = separatorArrayExplode(style);
  const joinedValues = value.map((v) => {
    if (style === "label" || style === "simple") {
      return allowReserved ? v : encodeURIComponent(v);
    }
    return serializePrimitiveParam({
      allowReserved,
      name,
      value: v
    });
  }).join(separator);
  return style === "label" || style === "matrix" ? separator + joinedValues : joinedValues;
};
const serializePrimitiveParam = ({
  allowReserved,
  name,
  value
}) => {
  if (value === void 0 || value === null) {
    return "";
  }
  if (typeof value === "object") {
    throw new Error(
      "Deeply-nested arrays/objects aren’t supported. Provide your own `querySerializer()` to handle these."
    );
  }
  return `${name}=${allowReserved ? value : encodeURIComponent(value)}`;
};
const serializeObjectParam = ({
  allowReserved,
  explode,
  name,
  style,
  value,
  valueOnly
}) => {
  if (value instanceof Date) {
    return valueOnly ? value.toISOString() : `${name}=${value.toISOString()}`;
  }
  if (style !== "deepObject" && !explode) {
    let values = [];
    Object.entries(value).forEach(([key, v]) => {
      values = [
        ...values,
        key,
        allowReserved ? v : encodeURIComponent(v)
      ];
    });
    const joinedValues2 = values.join(",");
    switch (style) {
      case "form":
        return `${name}=${joinedValues2}`;
      case "label":
        return `.${joinedValues2}`;
      case "matrix":
        return `;${name}=${joinedValues2}`;
      default:
        return joinedValues2;
    }
  }
  const separator = separatorObjectExplode(style);
  const joinedValues = Object.entries(value).map(
    ([key, v]) => serializePrimitiveParam({
      allowReserved,
      name: style === "deepObject" ? `${name}[${key}]` : key,
      value: v
    })
  ).join(separator);
  return style === "label" || style === "matrix" ? separator + joinedValues : joinedValues;
};
const PATH_PARAM_RE = /\{[^{}]+\}/g;
const defaultPathSerializer = ({ path, url: _url }) => {
  let url = _url;
  const matches = _url.match(PATH_PARAM_RE);
  if (matches) {
    for (const match of matches) {
      let explode = false;
      let name = match.substring(1, match.length - 1);
      let style = "simple";
      if (name.endsWith("*")) {
        explode = true;
        name = name.substring(0, name.length - 1);
      }
      if (name.startsWith(".")) {
        name = name.substring(1);
        style = "label";
      } else if (name.startsWith(";")) {
        name = name.substring(1);
        style = "matrix";
      }
      const value = path[name];
      if (value === void 0 || value === null) {
        continue;
      }
      if (Array.isArray(value)) {
        url = url.replace(
          match,
          serializeArrayParam({ explode, name, style, value })
        );
        continue;
      }
      if (typeof value === "object") {
        url = url.replace(
          match,
          serializeObjectParam({
            explode,
            name,
            style,
            value,
            valueOnly: true
          })
        );
        continue;
      }
      if (style === "matrix") {
        url = url.replace(
          match,
          `;${serializePrimitiveParam({
            name,
            value
          })}`
        );
        continue;
      }
      const replaceValue = encodeURIComponent(
        style === "label" ? `.${value}` : value
      );
      url = url.replace(match, replaceValue);
    }
  }
  return url;
};
const createQuerySerializer = ({
  allowReserved,
  array,
  object
} = {}) => {
  const querySerializer = (queryParams) => {
    const search = [];
    if (queryParams && typeof queryParams === "object") {
      for (const name in queryParams) {
        const value = queryParams[name];
        if (value === void 0 || value === null) {
          continue;
        }
        if (Array.isArray(value)) {
          const serializedArray = serializeArrayParam({
            allowReserved,
            explode: true,
            name,
            style: "form",
            value,
            ...array
          });
          if (serializedArray) search.push(serializedArray);
        } else if (typeof value === "object") {
          const serializedObject = serializeObjectParam({
            allowReserved,
            explode: true,
            name,
            style: "deepObject",
            value,
            ...object
          });
          if (serializedObject) search.push(serializedObject);
        } else {
          const serializedPrimitive = serializePrimitiveParam({
            allowReserved,
            name,
            value
          });
          if (serializedPrimitive) search.push(serializedPrimitive);
        }
      }
    }
    return search.join("&");
  };
  return querySerializer;
};
const getParseAs = (contentType) => {
  if (!contentType) {
    return "stream";
  }
  const cleanContent = contentType.split(";")[0]?.trim();
  if (!cleanContent) {
    return;
  }
  if (cleanContent.startsWith("application/json") || cleanContent.endsWith("+json")) {
    return "json";
  }
  if (cleanContent === "multipart/form-data") {
    return "formData";
  }
  if (["application/", "audio/", "image/", "video/"].some(
    (type) => cleanContent.startsWith(type)
  )) {
    return "blob";
  }
  if (cleanContent.startsWith("text/")) {
    return "text";
  }
  return;
};
const setAuthParams = async ({
  security,
  ...options
}) => {
  for (const auth of security) {
    const token = await getAuthToken(auth, options.auth);
    if (!token) {
      continue;
    }
    const name = auth.name ?? "Authorization";
    switch (auth.in) {
      case "query":
        if (!options.query) {
          options.query = {};
        }
        options.query[name] = token;
        break;
      case "cookie":
        options.headers.append("Cookie", `${name}=${token}`);
        break;
      case "header":
      default:
        options.headers.set(name, token);
        break;
    }
    return;
  }
};
const buildUrl = (options) => {
  const url = getUrl({
    baseUrl: options.baseUrl,
    path: options.path,
    query: options.query,
    querySerializer: typeof options.querySerializer === "function" ? options.querySerializer : createQuerySerializer(options.querySerializer),
    url: options.url
  });
  return url;
};
const getUrl = ({
  baseUrl,
  path,
  query,
  querySerializer,
  url: _url
}) => {
  const pathUrl = _url.startsWith("/") ? _url : `/${_url}`;
  let url = (baseUrl ?? "") + pathUrl;
  if (path) {
    url = defaultPathSerializer({ path, url });
  }
  let search = query ? querySerializer(query) : "";
  if (search.startsWith("?")) {
    search = search.substring(1);
  }
  if (search) {
    url += `?${search}`;
  }
  return url;
};
const mergeConfigs = (a, b) => {
  const config = { ...a, ...b };
  if (config.baseUrl?.endsWith("/")) {
    config.baseUrl = config.baseUrl.substring(0, config.baseUrl.length - 1);
  }
  config.headers = mergeHeaders(a.headers, b.headers);
  return config;
};
const mergeHeaders = (...headers) => {
  const mergedHeaders = new Headers();
  for (const header of headers) {
    if (!header || typeof header !== "object") {
      continue;
    }
    const iterator = header instanceof Headers ? header.entries() : Object.entries(header);
    for (const [key, value] of iterator) {
      if (value === null) {
        mergedHeaders.delete(key);
      } else if (Array.isArray(value)) {
        for (const v of value) {
          mergedHeaders.append(key, v);
        }
      } else if (value !== void 0) {
        mergedHeaders.set(
          key,
          typeof value === "object" ? JSON.stringify(value) : value
        );
      }
    }
  }
  return mergedHeaders;
};
class Interceptors {
  _fns;
  constructor() {
    this._fns = [];
  }
  clear() {
    this._fns = [];
  }
  getInterceptorIndex(id) {
    if (typeof id === "number") {
      return this._fns[id] ? id : -1;
    } else {
      return this._fns.indexOf(id);
    }
  }
  exists(id) {
    const index = this.getInterceptorIndex(id);
    return !!this._fns[index];
  }
  eject(id) {
    const index = this.getInterceptorIndex(id);
    if (this._fns[index]) {
      this._fns[index] = null;
    }
  }
  update(id, fn) {
    const index = this.getInterceptorIndex(id);
    if (this._fns[index]) {
      this._fns[index] = fn;
      return id;
    } else {
      return false;
    }
  }
  use(fn) {
    this._fns = [...this._fns, fn];
    return this._fns.length - 1;
  }
}
const createInterceptors = () => ({
  error: new Interceptors(),
  request: new Interceptors(),
  response: new Interceptors()
});
const defaultQuerySerializer = createQuerySerializer({
  allowReserved: false,
  array: {
    explode: true,
    style: "form"
  },
  object: {
    explode: true,
    style: "deepObject"
  }
});
const defaultHeaders = {
  "Content-Type": "application/json"
};
const createConfig = (override = {}) => ({
  ...jsonBodySerializer,
  headers: defaultHeaders,
  parseAs: "auto",
  querySerializer: defaultQuerySerializer,
  ...override
});
const createClient = (config = {}) => {
  let _config = mergeConfigs(createConfig(), config);
  const getConfig = () => ({ ..._config });
  const setConfig = (config2) => {
    _config = mergeConfigs(_config, config2);
    return getConfig();
  };
  const interceptors = createInterceptors();
  const request = async (options) => {
    const opts = {
      ..._config,
      ...options,
      fetch: options.fetch ?? _config.fetch ?? globalThis.fetch,
      headers: mergeHeaders(_config.headers, options.headers)
    };
    if (opts.security) {
      await setAuthParams({
        ...opts,
        security: opts.security
      });
    }
    if (opts.requestValidator) {
      await opts.requestValidator(opts);
    }
    if (opts.body && opts.bodySerializer) {
      opts.body = opts.bodySerializer(opts.body);
    }
    if (opts.body === void 0 || opts.body === "") {
      opts.headers.delete("Content-Type");
    }
    const url = buildUrl(opts);
    const requestInit = {
      redirect: "follow",
      ...opts
    };
    let request2 = new Request(url, requestInit);
    for (const fn of interceptors.request._fns) {
      if (fn) {
        request2 = await fn(request2, opts);
      }
    }
    const _fetch = opts.fetch;
    let response = await _fetch(request2);
    for (const fn of interceptors.response._fns) {
      if (fn) {
        response = await fn(response, request2, opts);
      }
    }
    const result = {
      request: request2,
      response
    };
    if (response.ok) {
      if (response.status === 204 || response.headers.get("Content-Length") === "0") {
        return opts.responseStyle === "data" ? {} : {
          data: {},
          ...result
        };
      }
      const parseAs = (opts.parseAs === "auto" ? getParseAs(response.headers.get("Content-Type")) : opts.parseAs) ?? "json";
      let data;
      switch (parseAs) {
        case "arrayBuffer":
        case "blob":
        case "formData":
        case "json":
        case "text":
          data = await response[parseAs]();
          break;
        case "stream":
          return opts.responseStyle === "data" ? response.body : {
            data: response.body,
            ...result
          };
      }
      if (parseAs === "json") {
        if (opts.responseValidator) {
          await opts.responseValidator(data);
        }
        if (opts.responseTransformer) {
          data = await opts.responseTransformer(data);
        }
      }
      return opts.responseStyle === "data" ? data : {
        data,
        ...result
      };
    }
    let error = await response.text();
    try {
      error = JSON.parse(error);
    } catch {
    }
    let finalError = error;
    for (const fn of interceptors.error._fns) {
      if (fn) {
        finalError = await fn(error, response, request2, opts);
      }
    }
    finalError = finalError || {};
    if (opts.throwOnError) {
      throw finalError;
    }
    return opts.responseStyle === "data" ? void 0 : {
      error: finalError,
      ...result
    };
  };
  return {
    buildUrl,
    connect: (options) => request({ ...options, method: "CONNECT" }),
    delete: (options) => request({ ...options, method: "DELETE" }),
    get: (options) => request({ ...options, method: "GET" }),
    getConfig,
    head: (options) => request({ ...options, method: "HEAD" }),
    interceptors,
    options: (options) => request({ ...options, method: "OPTIONS" }),
    patch: (options) => request({ ...options, method: "PATCH" }),
    post: (options) => request({ ...options, method: "POST" }),
    put: (options) => request({ ...options, method: "PUT" }),
    request,
    setConfig,
    trace: (options) => request({ ...options, method: "TRACE" })
  };
};
const client = createClient(createConfig({
  baseUrl: "http://localhost:8080/api"
}));
const getFileFileinfos = (options) => {
  return (options?.client ?? client).get({
    url: "/file/fileinfos",
    ...options
  });
};
const getProjectGetByPath = (options) => {
  return (options.client ?? client).get({
    url: "/project/get/{path}",
    ...options
  });
};
const getProjectRecent = (options) => {
  return (options?.client ?? client).get({
    url: "/project/recent",
    ...options
  });
};
const postProjectRenameManagedFile = (options) => {
  return (options.client ?? client).post({
    url: "/project/rename-managed-file",
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options.headers
    }
  });
};
function timestampToString(timestamp) {
  if (!timestamp) {
    return void 0;
  }
  if (typeof timestamp === "string") {
    return timestamp;
  }
  if (timestamp["time.Time"]) {
    return timestamp["time.Time"];
  }
  return void 0;
}
const FileDetailModal = ({ fileInfo, isOpen, onClose }) => {
  if (!isOpen || !fileInfo) return null;
  const formatSize = (bytes) => {
    if (!bytes) return "0 B";
    const units = ["B", "KB", "MB", "GB"];
    let size = bytes;
    let unitIndex = 0;
    while (size >= 1024 && unitIndex < units.length - 1) {
      size /= 1024;
      unitIndex++;
    }
    return `${size.toFixed(1)} ${units[unitIndex]}`;
  };
  const formatDate = (dateString) => {
    if (!dateString) return "-";
    return new Date(dateString).toLocaleString("ja-JP");
  };
  return /* @__PURE__ */ jsx("div", { className: "modal-overlay", onClick: onClose, children: /* @__PURE__ */ jsxs("div", { className: "modal-content", onClick: (e) => e.stopPropagation(), children: [
    /* @__PURE__ */ jsxs("div", { className: "modal-header", children: [
      /* @__PURE__ */ jsx("h2", { children: fileInfo.name }),
      /* @__PURE__ */ jsx("button", { className: "close-button", onClick: onClose, children: "×" })
    ] }),
    /* @__PURE__ */ jsxs("div", { className: "modal-body", children: [
      /* @__PURE__ */ jsxs("div", { className: "info-row", children: [
        /* @__PURE__ */ jsx("span", { className: "label", children: "種類:" }),
        /* @__PURE__ */ jsx("span", { className: "value", children: fileInfo.is_directory ? "フォルダー" : "ファイル" })
      ] }),
      /* @__PURE__ */ jsxs("div", { className: "info-row", children: [
        /* @__PURE__ */ jsx("span", { className: "label", children: "サイズ:" }),
        /* @__PURE__ */ jsx("span", { className: "value", children: formatSize(fileInfo.size) })
      ] }),
      /* @__PURE__ */ jsxs("div", { className: "info-row", children: [
        /* @__PURE__ */ jsx("span", { className: "label", children: "パス:" }),
        /* @__PURE__ */ jsx("span", { className: "value", children: fileInfo.path })
      ] }),
      /* @__PURE__ */ jsxs("div", { className: "info-row", children: [
        /* @__PURE__ */ jsx("span", { className: "label", children: "更新日時:" }),
        /* @__PURE__ */ jsx("span", { className: "value", children: formatDate(timestampToString(fileInfo.modified_time)) })
      ] })
    ] })
  ] }) });
};
const CustomTreeItem = React.memo(({
  itemId,
  node,
  onNodeClick,
  onNodeExpand,
  ...props
}) => {
  const getNodeIcon = (node2, expanded = false) => {
    if (node2.isDirectory) {
      return expanded ? /* @__PURE__ */ jsx(FolderOpen, { color: "primary" }) : /* @__PURE__ */ jsx(Folder, { color: "primary" });
    }
    if (node2.name === ".detail.yaml") {
      return /* @__PURE__ */ jsx("span", { style: { fontSize: "16px" }, children: "⚙️" });
    }
    const ext = node2.name?.split(".").pop()?.toLowerCase();
    switch (ext) {
      case "pdf":
        return /* @__PURE__ */ jsx("span", { style: { fontSize: "16px" }, children: "📄" });
      case "xlsx":
      case "xls":
      case "xlsm":
        return /* @__PURE__ */ jsx("span", { style: { fontSize: "16px" }, children: "📊" });
      case "jpg":
      case "jpeg":
      case "png":
      case "gif":
        return /* @__PURE__ */ jsx("span", { style: { fontSize: "16px" }, children: "🖼️" });
      case "mp4":
      case "avi":
      case "mov":
        return /* @__PURE__ */ jsx("span", { style: { fontSize: "16px" }, children: "🎬" });
      case "mp3":
      case "wav":
        return /* @__PURE__ */ jsx("span", { style: { fontSize: "16px" }, children: "🎵" });
      default:
        return /* @__PURE__ */ jsx(InsertDriveFile, { color: "action" });
    }
  };
  const formatSize = (bytes) => {
    if (!bytes || bytes === 0) return "";
    const units = ["B", "KB", "MB", "GB"];
    let size = bytes;
    let unitIndex = 0;
    while (size >= 1024 && unitIndex < units.length - 1) {
      size /= 1024;
      unitIndex++;
    }
    return ` (${size.toFixed(1)} ${units[unitIndex]})`;
  };
  const handleClick = React.useCallback((e) => {
    e.stopPropagation();
    onNodeClick(node);
  }, [node, onNodeClick]);
  return /* @__PURE__ */ jsx(
    TreeItem,
    {
      itemId,
      onClick: handleClick,
      label: /* @__PURE__ */ jsxs(Box, { sx: { display: "flex", alignItems: "center", py: 0.5 }, children: [
        getNodeIcon(node),
        /* @__PURE__ */ jsxs(Typography, { variant: "body2", sx: { flexGrow: 1, mr: 1, ml: 1 }, children: [
          node.name,
          !node.isDirectory && formatSize(node.size)
        ] }),
        node.isLoading && /* @__PURE__ */ jsx(CircularProgress, { size: 16 }),
        node.name === ".detail.yaml" && /* @__PURE__ */ jsx(
          Chip,
          {
            label: "詳細",
            size: "small",
            color: "info",
            variant: "outlined",
            sx: { ml: 1, fontSize: "10px", height: "20px" }
          }
        )
      ] }),
      ...props,
      children: node.children?.map((child) => /* @__PURE__ */ jsx(
        CustomTreeItem,
        {
          itemId: child.id,
          node: child,
          onNodeClick,
          onNodeExpand
        },
        child.id
      ))
    }
  );
});
const Files = () => {
  const [treeData, setTreeData] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [currentPath, setCurrentPath] = useState("~/penguin");
  const [pathInput, setPathInput] = useState("~/penguin");
  const [selectedNode, setSelectedNode] = useState(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [expanded, setExpanded] = useState([]);
  const { setFileCount, setCurrentPath: setContextPath } = useFileInfo();
  const convertToRelativePath = (frontendPath) => {
    if (frontendPath === "~/penguin" || frontendPath === "/home/shin/penguin") {
      return "";
    }
    if (frontendPath.startsWith("~/penguin/")) {
      return frontendPath.substring("~/penguin/".length);
    }
    if (frontendPath.startsWith("/home/shin/penguin/")) {
      return frontendPath.substring("/home/shin/penguin/".length);
    }
    return frontendPath;
  };
  const convertToTreeNode = (fileInfo, basePath) => {
    return {
      id: `${basePath}/${fileInfo.name}`,
      name: fileInfo.name,
      path: fileInfo.path || `${basePath}/${fileInfo.name}`,
      isDirectory: fileInfo.is_directory,
      size: fileInfo.size,
      modifiedTime: fileInfo.modified_time,
      children: fileInfo.is_directory ? [] : void 0,
      isLoaded: !fileInfo.is_directory,
      isLoading: false
    };
  };
  const loadFiles = useCallback(async (path, isRefresh = false) => {
    const frontendPath = path || "~/penguin";
    const relativePath = convertToRelativePath(frontendPath);
    setLoading(true);
    setError(null);
    try {
      const response = await getFileFileinfos({
        query: relativePath ? { path: relativePath } : {}
      });
      if (response.data) {
        const data = response.data;
        const nodes = data.map((fileInfo) => convertToTreeNode(fileInfo, frontendPath));
        if (path === "~/penguin" || !path || isRefresh) {
          setTreeData(nodes);
          setFileCount(data.length);
          setContextPath(frontendPath);
          setCurrentPath(frontendPath);
        }
        return nodes;
      } else if (response.error) {
        throw new Error("APIエラー: " + JSON.stringify(response.error));
      }
    } catch (err) {
      console.error("Error loading files:", err);
      setError(err instanceof Error ? err.message : "エラーが発生しました");
      return [];
    } finally {
      setLoading(false);
    }
  }, [setFileCount, setContextPath]);
  useEffect(() => {
    loadFiles();
  }, [loadFiles]);
  const handleNodeExpand = useCallback(async (nodeId, node) => {
    if (!node.isDirectory || node.isLoaded || node.isLoading) return;
    setTreeData((prevData) => updateNodeInTree(prevData, nodeId, { ...node, isLoading: true }));
    try {
      const children = await loadFiles(node.path, false);
      setTreeData((prevData) => updateNodeInTree(prevData, nodeId, {
        ...node,
        children,
        isLoaded: true,
        isLoading: false
      }));
    } catch (err) {
      setTreeData((prevData) => updateNodeInTree(prevData, nodeId, { ...node, isLoading: false }));
    }
  }, [loadFiles]);
  const updateNodeInTree = (nodes, targetId, updatedNode) => {
    return nodes.map((node) => {
      if (node.id === targetId) {
        return updatedNode;
      }
      if (node.children) {
        return {
          ...node,
          children: updateNodeInTree(node.children, targetId, updatedNode)
        };
      }
      return node;
    });
  };
  const handleNodeClick = useCallback((node) => {
    if (!node.isDirectory) {
      setSelectedNode(node);
      setIsModalOpen(true);
    } else {
      setExpanded((prev) => {
        const isExpanded = prev.includes(node.id);
        if (isExpanded) {
          return prev.filter((id) => id !== node.id);
        } else {
          const newExpanded = [...prev, node.id];
          if (!node.isLoaded && !node.isLoading) {
            setTimeout(() => handleNodeExpand(node.id, node), 0);
          }
          return newExpanded;
        }
      });
    }
  }, [handleNodeExpand]);
  const handlePathSubmit = (e) => {
    e.preventDefault();
    if (pathInput !== currentPath) {
      setTreeData([]);
      setExpanded([]);
      loadFiles(pathInput, true);
    }
  };
  const handleRefresh = () => {
    setTreeData([]);
    setExpanded([]);
    loadFiles(currentPath, true);
  };
  const handleGoHome = () => {
    setPathInput("~/penguin");
    setTreeData([]);
    setExpanded([]);
    loadFiles("~/penguin", true);
  };
  const getBreadcrumbs = () => {
    const parts = currentPath.replace("~/penguin", "").split("/").filter(Boolean);
    const breadcrumbs = [
      { label: "penguin", path: "~/penguin" }
    ];
    let accumulatedPath = "~/penguin";
    parts.forEach((part) => {
      accumulatedPath += `/${part}`;
      breadcrumbs.push({ label: part, path: accumulatedPath });
    });
    return breadcrumbs;
  };
  return /* @__PURE__ */ jsxs(Box, { sx: { height: "100vh", display: "flex", flexDirection: "column" }, children: [
    /* @__PURE__ */ jsxs(Paper, { elevation: 1, sx: { mb: 1 }, children: [
      /* @__PURE__ */ jsxs(Toolbar, { variant: "dense", children: [
        /* @__PURE__ */ jsx(IconButton, { onClick: handleGoHome, size: "small", title: "ホームに戻る", children: /* @__PURE__ */ jsx(Home, {}) }),
        /* @__PURE__ */ jsx(IconButton, { onClick: handleRefresh, size: "small", title: "リフレッシュ", children: /* @__PURE__ */ jsx(Refresh, {}) }),
        /* @__PURE__ */ jsx(Box, { component: "form", onSubmit: handlePathSubmit, sx: { flexGrow: 1, mx: 2 }, children: /* @__PURE__ */ jsx(
          TextField,
          {
            size: "small",
            fullWidth: true,
            value: pathInput,
            onChange: (e) => setPathInput(e.target.value),
            placeholder: "パスを入力",
            slotProps: {
              input: {
                startAdornment: /* @__PURE__ */ jsx(Search, { sx: { mr: 1, color: "action.active" } })
              }
            }
          }
        ) })
      ] }),
      /* @__PURE__ */ jsx(Box, { sx: { px: 2, pb: 1 }, children: /* @__PURE__ */ jsx(Breadcrumbs, { children: getBreadcrumbs().map((crumb, index) => /* @__PURE__ */ jsx(
        Link$1,
        {
          component: "button",
          variant: "body2",
          color: index === getBreadcrumbs().length - 1 ? "text.primary" : "inherit",
          onClick: () => {
            setPathInput(crumb.path);
            if (crumb.path !== currentPath) {
              setTreeData([]);
              setExpanded([]);
              loadFiles(crumb.path, true);
            }
          },
          underline: "hover",
          children: crumb.label
        },
        index
      )) }) })
    ] }),
    error && /* @__PURE__ */ jsx(Alert, { severity: "error", sx: { mb: 1 }, children: error }),
    /* @__PURE__ */ jsxs(Paper, { sx: { flex: 1, overflow: "auto", p: 1 }, children: [
      loading && treeData.length === 0 ? /* @__PURE__ */ jsx(Box, { sx: { display: "flex", justifyContent: "center", p: 3 }, children: /* @__PURE__ */ jsx(CircularProgress, {}) }) : /* @__PURE__ */ jsx(
        SimpleTreeView,
        {
          expandedItems: expanded,
          onExpandedItemsChange: () => {
          },
          sx: {
            flexGrow: 1,
            overflowY: "auto",
            "& .MuiTreeItem-content": {
              "&:hover": {
                backgroundColor: "action.hover"
              }
            }
          },
          children: treeData.map((node) => /* @__PURE__ */ jsx(
            CustomTreeItem,
            {
              itemId: node.id,
              node,
              onNodeClick: handleNodeClick,
              onNodeExpand: handleNodeExpand
            },
            node.id
          ))
        }
      ),
      treeData.length === 0 && !loading && /* @__PURE__ */ jsx(Box, { sx: { display: "flex", justifyContent: "center", p: 3 }, children: /* @__PURE__ */ jsx(Typography, { color: "text.secondary", children: "フォルダーが空です" }) })
    ] }),
    /* @__PURE__ */ jsx(
      FileDetailModal,
      {
        fileInfo: selectedNode,
        isOpen: isModalOpen,
        onClose: () => {
          setIsModalOpen(false);
          setSelectedNode(null);
        }
      }
    )
  ] });
};
const _layout_files = UNSAFE_withComponentProps(function FilesPage() {
  return /* @__PURE__ */ jsx(Files, {});
});
const route3 = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
  __proto__: null,
  default: _layout_files
}, Symbol.toStringTag, { value: "Module" }));
const CalendarPicker = ({
  value,
  onChange,
  placeholder,
  disabled = false,
  minDate
}) => {
  const selectedDate = value ? new Date(value) : null;
  const [internalDate, setInternalDate] = React.useState(selectedDate);
  React.useEffect(() => {
    setInternalDate(selectedDate);
  }, [value]);
  const handleDateChange = (date) => {
    setInternalDate(date);
  };
  const handleDateAccept = (date) => {
    if (date) {
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, "0");
      const day = String(date.getDate()).padStart(2, "0");
      onChange(`${year}-${month}-${day}`);
    } else {
      onChange("");
    }
  };
  const minDateObj = minDate ? new Date(minDate) : void 0;
  return /* @__PURE__ */ jsx(LocalizationProvider, { dateAdapter: AdapterDateFns, adapterLocale: ja, children: /* @__PURE__ */ jsx(
    DatePicker,
    {
      value: internalDate,
      onChange: handleDateChange,
      onAccept: handleDateAccept,
      disabled,
      minDate: minDateObj,
      format: "yyyy年MM月dd日",
      enableAccessibleFieldDOMStructure: false,
      views: ["year", "month", "day"],
      openTo: "day",
      reduceAnimations: true,
      closeOnSelect: false,
      slots: {
        textField: TextField
      },
      slotProps: {
        textField: {
          placeholder,
          size: "small",
          fullWidth: true,
          sx: {
            "& .MuiOutlinedInput-root": {
              fontSize: "0.875rem",
              borderRadius: "6px",
              "&:hover .MuiOutlinedInput-notchedOutline": {
                borderColor: "#3b82f6"
              },
              "&.Mui-focused .MuiOutlinedInput-notchedOutline": {
                borderColor: "#3b82f6",
                boxShadow: "0 0 0 3px rgba(59, 130, 246, 0.1)"
              }
            },
            "& .MuiInputBase-input": {
              padding: "8px 12px"
            }
          }
        }
      }
    }
  ) });
};
const ProjectDetailModal = ({ isOpen, onClose, project, onUpdate, onProjectUpdate }) => {
  const [formData, setFormData] = useState({
    id: "",
    company_name: "",
    location_name: "",
    description: "",
    tags: "",
    start_date: "",
    end_date: ""
  });
  const [currentProject, setCurrentProject] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);
  const [isRenaming, setIsRenaming] = useState(false);
  const [hasFilenameChanges, setHasFilenameChanges] = useState(false);
  const [initialFilenameData, setInitialFilenameData] = useState({ start_date: "", company_name: "", location_name: "" });
  const extractDateString = (timestamp) => {
    if (!timestamp) return "";
    if (typeof timestamp === "string") {
      const match = timestamp.match(/^(\d{4})-(\d{2})-(\d{2})/);
      if (match) {
        return `${match[1]}-${match[2]}-${match[3]}`;
      }
    }
    if (typeof timestamp === "object") {
      const timeString = timestamp["time.Time"];
      if (timeString) {
        const match = timeString.match(/^(\d{4})-(\d{2})-(\d{2})/);
        if (match) {
          return `${match[1]}-${match[2]}-${match[3]}`;
        }
      }
    }
    return "";
  };
  useEffect(() => {
    if (project) {
      const startDate = extractDateString(project.start_date);
      const endDate = extractDateString(project.end_date);
      const companyName = project.company_name || "";
      const locationName = project.location_name || "";
      setCurrentProject(project);
      setFormData({
        id: project.id || "",
        company_name: companyName,
        location_name: locationName,
        description: project.description || "",
        tags: Array.isArray(project.tags) ? project.tags.join(", ") : project.tags || "",
        start_date: startDate,
        end_date: endDate
      });
      setInitialFilenameData({
        start_date: startDate,
        company_name: companyName,
        location_name: locationName
      });
      setHasFilenameChanges(false);
      setError(null);
    }
  }, [project]);
  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value
    }));
  };
  const generateRecommendedFileName = (originalFileName, formData2) => {
    if (!formData2.start_date || !formData2.company_name || !formData2.location_name) {
      return originalFileName;
    }
    const datePart = formData2.start_date.replace(/-/g, "").substring(0, 8);
    const formattedDate = `${datePart.substring(0, 4)}-${datePart.substring(4, 8)}`;
    const newPrefix = `${formattedDate} ${formData2.company_name} ${formData2.location_name}`;
    const fileExtension = originalFileName.includes(".") ? "." + originalFileName.split(".").pop() : "";
    const datePattern = /^\d{4}-\d{4}\s+/;
    if (datePattern.test(originalFileName)) {
      const afterPrefix = originalFileName.replace(/^\d{4}-\d{4}\s+[^\s]+\s+[^\s]+\s*/, "");
      return afterPrefix ? `${newPrefix} ${afterPrefix}` : `${newPrefix}${fileExtension}`;
    } else {
      const nameWithoutExt = fileExtension ? originalFileName.substring(0, originalFileName.lastIndexOf(".")) : originalFileName;
      return `${newPrefix} ${nameWithoutExt}${fileExtension}`;
    }
  };
  const updateManagedFileRecommendations = (formData2) => {
    if (!currentProject?.managed_files) return;
    const updatedManagedFiles = currentProject.managed_files.map((file) => {
      if (file.current) {
        const recommendedName = generateRecommendedFileName(file.current, formData2);
        return {
          ...file,
          recommended: recommendedName
        };
      }
      return file;
    });
    setCurrentProject((prev) => prev ? {
      ...prev,
      managed_files: updatedManagedFiles
    } : null);
  };
  const checkFilenameChanges = (currentFormData) => {
    const hasChanges = currentFormData.start_date !== initialFilenameData.start_date || currentFormData.company_name !== initialFilenameData.company_name || currentFormData.location_name !== initialFilenameData.location_name;
    setHasFilenameChanges(hasChanges);
    if (hasChanges) {
      updateManagedFileRecommendations(currentFormData);
    }
  };
  const handleFieldUpdateWithData = async (useFormData) => {
    if (!project) return;
    setIsLoading(true);
    setError(null);
    try {
      const updatedProject = {
        ...project,
        company_name: useFormData.company_name,
        location_name: useFormData.location_name,
        description: useFormData.description,
        tags: useFormData.tags ? useFormData.tags.split(",").map((tag) => tag.trim()).filter((tag) => tag.length > 0) : [],
        start_date: useFormData.start_date ? { "time.Time": `${useFormData.start_date}T00:00:00+09:00` } : void 0,
        end_date: useFormData.end_date ? { "time.Time": `${useFormData.end_date}T23:59:59+09:00` } : void 0
      };
      const originalFolderName = project.name;
      const savedProject = await onUpdate(updatedProject);
      const folderNameChanged = originalFolderName && savedProject.name && originalFolderName !== savedProject.name;
      if (savedProject.name) {
        const latestProjectResponse = await getProjectGetByPath({
          path: {
            path: savedProject.name
          }
        });
        if (latestProjectResponse.data) {
          const latestProject = latestProjectResponse.data;
          if (folderNameChanged) {
            if (onProjectUpdate) {
              onProjectUpdate(latestProject);
            }
            setTimeout(() => {
              onClose();
            }, 100);
            return;
          }
          setCurrentProject(latestProject);
          if (onProjectUpdate) {
            onProjectUpdate(latestProject);
          }
          const startDate = extractDateString(latestProject.start_date);
          const endDate = extractDateString(latestProject.end_date);
          setFormData({
            id: latestProject.id || "",
            company_name: latestProject.company_name || "",
            location_name: latestProject.location_name || "",
            description: latestProject.description || "",
            tags: Array.isArray(latestProject.tags) ? latestProject.tags.join(", ") : latestProject.tags || "",
            start_date: startDate,
            end_date: endDate
          });
        }
      } else {
        if (onProjectUpdate) {
          onProjectUpdate(savedProject);
        }
        const startDate = extractDateString(savedProject.start_date);
        const endDate = extractDateString(savedProject.end_date);
        setFormData({
          id: savedProject.id || "",
          company_name: savedProject.company_name || "",
          location_name: savedProject.location_name || "",
          description: savedProject.description || "",
          tags: Array.isArray(savedProject.tags) ? savedProject.tags.join(", ") : savedProject.tags || "",
          start_date: startDate,
          end_date: endDate
        });
      }
    } catch (err) {
      console.error("Error updating field:", err);
      setError(err instanceof Error ? err.message : "フィールドの更新に失敗しました");
    } finally {
      setIsLoading(false);
    }
  };
  const handleFieldUpdate = async (fieldName) => {
    if (hasFilenameChanges) {
      return;
    }
    if (fieldName === "company_name" || fieldName === "location_name") {
      return;
    }
    handleFieldUpdateWithData(formData);
  };
  const handleNonFilenameBlur = (fieldName) => () => {
    if (hasFilenameChanges) {
      return;
    }
    if (fieldName === "description" || fieldName === "tags") {
      handleFieldUpdateWithData(formData);
    }
  };
  const handleFilenameBlur = (fieldName) => () => {
    if (fieldName === "company_name" || fieldName === "location_name") {
      checkFilenameChanges(formData);
    }
  };
  const handleKeyDown = (e, fieldName) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleFieldUpdate(fieldName);
    }
  };
  const handleFilenameKeyDown = (e, fieldName) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      if (fieldName === "company_name" || fieldName === "location_name") {
        checkFilenameChanges(formData);
      }
    }
  };
  const handleDaisyDateChange = (dateString, fieldName) => {
    const newFormData = {
      ...formData,
      [fieldName]: dateString
    };
    setFormData(newFormData);
    if (fieldName === "start_date") {
      checkFilenameChanges(newFormData);
      return;
    }
    if (hasFilenameChanges) {
      return;
    }
    if (dateString) {
      handleFieldUpdateWithData(newFormData);
    }
  };
  const handleFilenameUpdate = async () => {
    await handleFieldUpdateWithData(formData);
    setHasFilenameChanges(false);
    setInitialFilenameData({
      start_date: formData.start_date || "",
      company_name: formData.company_name || "",
      location_name: formData.location_name || ""
    });
  };
  const handleRenameFiles = async () => {
    if (!currentProject || !currentProject.managed_files) return;
    setIsRenaming(true);
    setError(null);
    try {
      const currentFiles = currentProject.managed_files.filter((file) => file.current && file.recommended && file.current !== file.recommended).map((file) => file.current);
      if (currentFiles.length === 0) {
        setError("変更対象のファイルがありません");
        return;
      }
      const response = await postProjectRenameManagedFile({
        body: {
          project: currentProject,
          currents: currentFiles
        }
      });
      if (response.data) {
        if (currentProject.name) {
          const updatedProjectResponse = await getProjectGetByPath({
            path: {
              path: currentProject.name
            }
          });
          if (updatedProjectResponse.data && onProjectUpdate) {
            setCurrentProject(updatedProjectResponse.data);
            onProjectUpdate(updatedProjectResponse.data);
            const startDate = extractDateString(updatedProjectResponse.data.start_date);
            const endDate = extractDateString(updatedProjectResponse.data.end_date);
            setFormData({
              id: updatedProjectResponse.data.id || "",
              company_name: updatedProjectResponse.data.company_name || "",
              location_name: updatedProjectResponse.data.location_name || "",
              description: updatedProjectResponse.data.description || "",
              tags: Array.isArray(updatedProjectResponse.data.tags) ? updatedProjectResponse.data.tags.join(", ") : updatedProjectResponse.data.tags || "",
              start_date: startDate,
              end_date: endDate
            });
          }
        }
      } else if (response.error) {
        throw new Error("ファイル名の変更に失敗しました");
      }
    } catch (err) {
      console.error("Error renaming files:", err);
      setError(err instanceof Error ? err.message : "ファイル名の変更に失敗しました");
    } finally {
      setIsRenaming(false);
    }
  };
  return /* @__PURE__ */ jsxs(
    Dialog,
    {
      open: isOpen,
      onClose,
      maxWidth: "md",
      fullWidth: true,
      slotProps: {
        paper: {
          sx: {
            borderRadius: 2,
            minHeight: "500px"
          }
        }
      },
      children: [
        /* @__PURE__ */ jsxs(
          DialogTitle,
          {
            sx: {
              display: "flex",
              justifyContent: "space-between",
              alignItems: "center",
              pb: 1
            },
            children: [
              /* @__PURE__ */ jsx(Typography, { variant: "h6", component: "h2", children: "工事詳細編集" }),
              /* @__PURE__ */ jsx(
                IconButton,
                {
                  onClick: onClose,
                  size: "small",
                  sx: {
                    color: "grey.500"
                  },
                  children: /* @__PURE__ */ jsx(CloseIcon, {})
                }
              )
            ]
          }
        ),
        /* @__PURE__ */ jsxs(DialogContent, { dividers: true, children: [
          error && /* @__PURE__ */ jsx(Alert, { severity: "error", sx: { mb: 2 }, children: error }),
          /* @__PURE__ */ jsxs(Box, { sx: { mt: 1 }, children: [
            /* @__PURE__ */ jsxs(Box, { sx: { display: "flex", gap: 2, mb: 3, flexDirection: { xs: "column", sm: "row" } }, children: [
              /* @__PURE__ */ jsx(Box, { sx: { flex: 1 }, children: /* @__PURE__ */ jsxs(
                Paper,
                {
                  elevation: 0,
                  sx: {
                    p: 2,
                    backgroundColor: "primary.50",
                    border: "2px solid",
                    borderColor: "primary.200",
                    borderRadius: 2
                  },
                  children: [
                    /* @__PURE__ */ jsx(
                      Typography,
                      {
                        variant: "subtitle2",
                        sx: {
                          color: "primary.main",
                          fontWeight: 600,
                          mb: 1
                        },
                        children: "会社名"
                      }
                    ),
                    /* @__PURE__ */ jsx(
                      TextField,
                      {
                        fullWidth: true,
                        size: "small",
                        name: "company_name",
                        value: formData.company_name,
                        onChange: handleInputChange,
                        onKeyDown: (e) => handleFilenameKeyDown(e, "company_name"),
                        onBlur: handleFilenameBlur("company_name"),
                        required: true,
                        disabled: isLoading,
                        sx: {
                          "& .MuiOutlinedInput-root": {
                            backgroundColor: "white"
                          }
                        }
                      }
                    )
                  ]
                }
              ) }),
              /* @__PURE__ */ jsx(Box, { sx: { flex: 1 }, children: /* @__PURE__ */ jsxs(
                Paper,
                {
                  elevation: 0,
                  sx: {
                    p: 2,
                    backgroundColor: "primary.50",
                    border: "2px solid",
                    borderColor: "primary.200",
                    borderRadius: 2
                  },
                  children: [
                    /* @__PURE__ */ jsx(
                      Typography,
                      {
                        variant: "subtitle2",
                        sx: {
                          color: "primary.main",
                          fontWeight: 600,
                          mb: 1
                        },
                        children: "現場名"
                      }
                    ),
                    /* @__PURE__ */ jsx(
                      TextField,
                      {
                        fullWidth: true,
                        size: "small",
                        name: "location_name",
                        value: formData.location_name,
                        onChange: handleInputChange,
                        onKeyDown: (e) => handleFilenameKeyDown(e, "location_name"),
                        onBlur: handleFilenameBlur("location_name"),
                        required: true,
                        disabled: isLoading,
                        sx: {
                          "& .MuiOutlinedInput-root": {
                            backgroundColor: "white"
                          }
                        }
                      }
                    )
                  ]
                }
              ) })
            ] }),
            /* @__PURE__ */ jsxs(Box, { sx: { display: "flex", gap: 2, mb: 3, flexDirection: { xs: "column", sm: "row" } }, children: [
              /* @__PURE__ */ jsx(Box, { sx: { flex: 1 }, children: /* @__PURE__ */ jsxs(
                Paper,
                {
                  elevation: 0,
                  sx: {
                    p: 2,
                    backgroundColor: "primary.50",
                    border: "2px solid",
                    borderColor: "primary.200",
                    borderRadius: 2
                  },
                  children: [
                    /* @__PURE__ */ jsx(
                      Typography,
                      {
                        variant: "subtitle2",
                        sx: {
                          color: "primary.main",
                          fontWeight: 600,
                          mb: 1
                        },
                        children: "開始日"
                      }
                    ),
                    /* @__PURE__ */ jsx(
                      CalendarPicker,
                      {
                        value: formData.start_date || "",
                        onChange: (dateString) => handleDaisyDateChange(dateString, "start_date"),
                        placeholder: "開始日を選択",
                        disabled: isLoading
                      }
                    )
                  ]
                }
              ) }),
              /* @__PURE__ */ jsx(Box, { sx: { flex: 1 }, children: /* @__PURE__ */ jsxs(
                Paper,
                {
                  elevation: 0,
                  sx: {
                    p: 2,
                    backgroundColor: "grey.50",
                    border: "1px solid",
                    borderColor: "grey.200",
                    borderRadius: 2
                  },
                  children: [
                    /* @__PURE__ */ jsx(
                      Typography,
                      {
                        variant: "subtitle2",
                        sx: {
                          color: "text.primary",
                          fontWeight: 500,
                          mb: 1
                        },
                        children: "終了日"
                      }
                    ),
                    /* @__PURE__ */ jsx(
                      CalendarPicker,
                      {
                        value: formData.end_date || "",
                        onChange: (dateString) => handleDaisyDateChange(dateString, "end_date"),
                        placeholder: "終了日を選択",
                        disabled: isLoading,
                        minDate: formData.start_date
                      }
                    )
                  ]
                }
              ) })
            ] }),
            /* @__PURE__ */ jsxs(Box, { sx: { mb: 3 }, children: [
              /* @__PURE__ */ jsxs(Box, { sx: { mb: 2 }, children: [
                /* @__PURE__ */ jsx(
                  Typography,
                  {
                    variant: "subtitle2",
                    sx: {
                      color: "text.primary",
                      fontWeight: 500,
                      mb: 1
                    },
                    children: "説明"
                  }
                ),
                /* @__PURE__ */ jsx(
                  TextField,
                  {
                    fullWidth: true,
                    multiline: true,
                    rows: 3,
                    name: "description",
                    value: formData.description,
                    onChange: handleInputChange,
                    onKeyDown: (e) => handleKeyDown(e, "description"),
                    onBlur: handleNonFilenameBlur("description"),
                    disabled: isLoading,
                    placeholder: "工事の詳細説明を入力してください"
                  }
                )
              ] }),
              /* @__PURE__ */ jsxs(Box, { children: [
                /* @__PURE__ */ jsx(
                  Typography,
                  {
                    variant: "subtitle2",
                    sx: {
                      color: "text.primary",
                      fontWeight: 500,
                      mb: 1
                    },
                    children: "タグ"
                  }
                ),
                /* @__PURE__ */ jsx(
                  TextField,
                  {
                    fullWidth: true,
                    name: "tags",
                    value: formData.tags,
                    onChange: handleInputChange,
                    onKeyDown: (e) => handleKeyDown(e, "tags"),
                    onBlur: handleNonFilenameBlur("tags"),
                    placeholder: "カンマ区切りで入力",
                    disabled: isLoading
                  }
                )
              ] })
            ] }),
            hasFilenameChanges && /* @__PURE__ */ jsx(
              Alert,
              {
                severity: "warning",
                sx: { mb: 3 },
                icon: /* @__PURE__ */ jsx(WarningIcon, {}),
                action: /* @__PURE__ */ jsx(
                  Button,
                  {
                    color: "warning",
                    variant: "contained",
                    size: "small",
                    onClick: handleFilenameUpdate,
                    disabled: isLoading,
                    startIcon: isLoading ? /* @__PURE__ */ jsx(CircularProgress, { size: 16 }) : null,
                    children: isLoading ? "更新中..." : "ファイル名を更新"
                  }
                ),
                children: /* @__PURE__ */ jsx(Typography, { variant: "body2", fontWeight: 500, children: "ファイル名の変更が必要です。このボタンを押下して全ての変更をまとめて更新してください。" })
              }
            ),
            currentProject?.managed_files && currentProject.managed_files.length > 0 && /* @__PURE__ */ jsxs(Box, { sx: { mb: 3 }, children: [
              /* @__PURE__ */ jsx(
                Typography,
                {
                  variant: "subtitle2",
                  sx: {
                    color: "text.primary",
                    fontWeight: 500,
                    mb: 2
                  },
                  children: "管理ファイル"
                }
              ),
              /* @__PURE__ */ jsx(
                Paper,
                {
                  elevation: 0,
                  sx: {
                    border: "1px solid",
                    borderColor: "grey.200",
                    borderRadius: 2,
                    backgroundColor: "grey.50",
                    p: 2
                  },
                  children: currentProject.managed_files.map((file, index) => {
                    const needsRename = file.current && file.recommended && file.current !== file.recommended;
                    return /* @__PURE__ */ jsxs(
                      Paper,
                      {
                        elevation: 0,
                        sx: {
                          mb: index < currentProject.managed_files.length - 1 ? 2 : 0,
                          p: 2,
                          backgroundColor: "white",
                          border: "1px solid",
                          borderColor: "grey.200",
                          borderRadius: 1
                        },
                        children: [
                          /* @__PURE__ */ jsxs(Box, { sx: { display: "flex", alignItems: "center", mb: file.recommended ? 1 : 0 }, children: [
                            /* @__PURE__ */ jsx(AttachFileIcon, { sx: { mr: 1, color: "grey.600", fontSize: "1rem" } }),
                            /* @__PURE__ */ jsx(Typography, { variant: "body2", fontWeight: 500, sx: { mr: 1 }, children: "現在:" }),
                            /* @__PURE__ */ jsx(Typography, { variant: "body2", sx: { fontFamily: "monospace", color: "grey.700" }, children: file.current })
                          ] }),
                          file.recommended && /* @__PURE__ */ jsxs(Box, { sx: { display: "flex", alignItems: "center", justifyContent: "space-between" }, children: [
                            /* @__PURE__ */ jsxs(Box, { sx: { display: "flex", alignItems: "center" }, children: [
                              /* @__PURE__ */ jsx(LightbulbIcon, { sx: { mr: 1, color: "warning.main", fontSize: "1rem" } }),
                              /* @__PURE__ */ jsx(Typography, { variant: "body2", fontWeight: 500, sx: { mr: 1 }, children: "推奨:" }),
                              /* @__PURE__ */ jsx(Typography, { variant: "body2", sx: { fontFamily: "monospace", color: "primary.main" }, children: file.recommended })
                            ] }),
                            needsRename && /* @__PURE__ */ jsx(
                              Button,
                              {
                                variant: "contained",
                                size: "small",
                                onClick: handleRenameFiles,
                                disabled: isRenaming || isLoading || hasFilenameChanges,
                                startIcon: isRenaming ? /* @__PURE__ */ jsx(CircularProgress, { size: 16 }) : null,
                                title: hasFilenameChanges ? "ファイル名の変更が必要です。上のボタンで更新してください。" : "",
                                children: isRenaming ? "変更中..." : "変更"
                              }
                            )
                          ] })
                        ]
                      },
                      index
                    );
                  })
                }
              )
            ] }),
            currentProject && (!currentProject.managed_files || currentProject.managed_files.length === 0) && /* @__PURE__ */ jsxs(Box, { sx: { mb: 3 }, children: [
              /* @__PURE__ */ jsx(
                Typography,
                {
                  variant: "subtitle2",
                  sx: {
                    color: "text.primary",
                    fontWeight: 500,
                    mb: 2
                  },
                  children: "管理ファイル"
                }
              ),
              /* @__PURE__ */ jsxs(
                Paper,
                {
                  elevation: 0,
                  sx: {
                    border: "1px solid",
                    borderColor: "grey.200",
                    borderRadius: 2,
                    backgroundColor: "grey.50",
                    p: 3,
                    textAlign: "center"
                  },
                  children: [
                    /* @__PURE__ */ jsx(AttachFileIcon, { sx: { fontSize: "2rem", color: "grey.400", mb: 1 } }),
                    /* @__PURE__ */ jsx(Typography, { variant: "body2", color: "text.secondary", children: "管理ファイルはありません" })
                  ]
                }
              )
            ] })
          ] })
        ] })
      ]
    }
  );
};
const Projects = () => {
  const [projects, setProjects] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [selectedProject, setSelectedProject] = useState(
    null
  );
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [showHelp, setShowHelp] = useState(false);
  const { setProjectCount } = useProject();
  const loadProjects = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await getProjectRecent();
      if (response.data) {
        setProjects(response.data);
        setProjectCount(response.data.length);
      } else {
        setProjects([]);
        setProjectCount(0);
      }
    } catch (err) {
      console.error("Error loading kouji entries:", err);
      setError(
        `工事データの読み込みに失敗しました: ${err instanceof Error ? err.message : String(err)}`
      );
    } finally {
      setLoading(false);
    }
  };
  useEffect(() => {
    loadProjects();
  }, []);
  const handleProjectClick = (project) => {
    setSelectedProject(project);
    setIsEditModalOpen(true);
  };
  const updateProject = async (updatedProject) => {
    try {
      const response = await fetch("http://localhost:8080/api/project/update", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify(updatedProject)
      });
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || "更新に失敗しました");
      }
      const savedProject = await response.json();
      setProjects(
        (prevProjects) => prevProjects.map((p) => p.id === savedProject.id ? savedProject : p)
      );
      return savedProject;
    } catch (err) {
      console.error("Error updating project:", err);
      throw err;
    }
  };
  const closeEditModal = () => {
    setIsEditModalOpen(false);
    setSelectedProject(null);
  };
  const needsFileRename = (project) => {
    if (!project.managed_files || project.managed_files.length === 0) {
      return false;
    }
    const needsRename = project.managed_files.some((file) => {
      return file.current && file.recommended && file.current !== file.recommended;
    });
    return needsRename;
  };
  const handleProjectUpdate = (updatedProject) => {
    setSelectedProject(updatedProject);
    setProjects((prevProjects) => {
      const existingIndex = prevProjects.findIndex((p) => p.id === updatedProject.id);
      if (existingIndex !== -1) {
        const updatedProjects = [...prevProjects];
        updatedProjects[existingIndex] = updatedProject;
        return updatedProjects;
      } else {
        const oldProjectIndex = prevProjects.findIndex(
          (p) => p.company_name === updatedProject.company_name && p.location_name === updatedProject.location_name && p.id !== updatedProject.id
        );
        if (oldProjectIndex !== -1) {
          const updatedProjects = [...prevProjects];
          updatedProjects.splice(oldProjectIndex, 1);
          updatedProjects.push(updatedProject);
          return updatedProjects.sort((a, b) => {
            const dateA = a.start_date ? new Date(typeof a.start_date === "string" ? a.start_date : a.start_date["time.Time"]).getTime() : 0;
            const dateB = b.start_date ? new Date(typeof b.start_date === "string" ? b.start_date : b.start_date["time.Time"]).getTime() : 0;
            if (dateA > 0 && dateB === 0) return -1;
            if (dateA === 0 && dateB > 0) return 1;
            if (dateA > 0 && dateB > 0) return dateB - dateA;
            return (b.name || "").localeCompare(a.name || "");
          });
        } else {
          return [...prevProjects, updatedProject];
        }
      }
    });
  };
  if (loading) {
    return /* @__PURE__ */ jsx("div", { style: { padding: "20px" }, children: /* @__PURE__ */ jsx("div", { children: "工事データを読み込み中..." }) });
  }
  if (error) {
    return /* @__PURE__ */ jsxs("div", { style: { padding: "20px" }, children: [
      /* @__PURE__ */ jsx(
        "div",
        {
          style: {
            color: "red",
            padding: "10px",
            backgroundColor: "#ffe6e6",
            borderRadius: "4px"
          },
          children: error
        }
      ),
      /* @__PURE__ */ jsx(
        "button",
        {
          onClick: loadProjects,
          style: { marginTop: "10px", padding: "10px 20px" },
          children: "再試行"
        }
      )
    ] });
  }
  return /* @__PURE__ */ jsxs("div", { style: {
    padding: "20px",
    paddingTop: "60px",
    flex: 1,
    overflow: "hidden",
    display: "flex",
    flexDirection: "column",
    boxSizing: "border-box"
  }, children: [
    showHelp && /* @__PURE__ */ jsxs(
      "div",
      {
        style: {
          marginBottom: "20px",
          padding: "15px",
          backgroundColor: "#f0f8ff",
          borderRadius: "4px",
          border: "1px solid #b3d9ff",
          position: "relative",
          flexShrink: 0
        },
        children: [
          /* @__PURE__ */ jsx(
            "button",
            {
              onClick: () => setShowHelp(false),
              style: {
                position: "absolute",
                top: "10px",
                right: "10px",
                background: "none",
                border: "none",
                fontSize: "16px",
                cursor: "pointer",
                color: "#666"
              },
              title: "閉じる",
              children: "×"
            }
          ),
          /* @__PURE__ */ jsx("h3", { style: { marginTop: 0 }, children: "使用方法" }),
          /* @__PURE__ */ jsxs("p", { children: [
            "📝 ",
            /* @__PURE__ */ jsx("strong", { children: "リストをクリック" }),
            "して工事情報を編集できます"
          ] }),
          /* @__PURE__ */ jsx("p", { children: "✅ 開始日・終了日・説明・タグ・会社名・現場名を編集可能" }),
          /* @__PURE__ */ jsx("p", { children: "💾 編集後は自動で保存されます" }),
          /* @__PURE__ */ jsx("h3", { style: { marginTop: "15px" }, children: "開発状況" }),
          /* @__PURE__ */ jsx("p", { children: "✅ 工事データの取得" }),
          /* @__PURE__ */ jsx("p", { children: "✅ 編集モーダル機能" }),
          /* @__PURE__ */ jsx("p", { children: "🔄 工程表機能（次のステップ）" })
        ]
      }
    ),
    /* @__PURE__ */ jsxs(
      "div",
      {
        style: {
          border: "1px solid #ddd",
          borderRadius: "8px",
          overflow: "hidden",
          position: "relative",
          height: "calc(100vh - 240px)",
          display: "flex",
          flexDirection: "column"
        },
        children: [
          /* @__PURE__ */ jsxs(
            "div",
            {
              style: {
                backgroundColor: "#f5f5f5",
                padding: "10px 15px",
                fontWeight: "bold",
                borderBottom: "1px solid #ddd",
                display: "flex",
                alignItems: "center",
                gap: "16px",
                flexShrink: 0
              },
              children: [
                /* @__PURE__ */ jsx("div", { style: { minWidth: "90px", textAlign: "center", fontSize: "14px" }, children: "開始日" }),
                /* @__PURE__ */ jsx("div", { style: { minWidth: "120px", fontSize: "14px" }, children: "会社名" }),
                /* @__PURE__ */ jsx("div", { style: { minWidth: "120px", fontSize: "14px" }, children: "現場名" }),
                /* @__PURE__ */ jsx("div", { style: { flex: 1 } }),
                /* @__PURE__ */ jsx("div", { style: { minWidth: "90px", textAlign: "center", fontSize: "14px", marginRight: "24px" }, children: "終了日" }),
                /* @__PURE__ */ jsx("div", { style: { minWidth: "80px", textAlign: "center", fontSize: "14px" }, children: "ステータス" }),
                /* @__PURE__ */ jsx(
                  Link,
                  {
                    to: "/projects/gantt",
                    style: {
                      background: "#4CAF50",
                      color: "white",
                      border: "none",
                      borderRadius: "4px",
                      padding: "6px 12px",
                      fontSize: "12px",
                      textDecoration: "none",
                      cursor: "pointer",
                      transition: "background 0.2s",
                      marginRight: "8px"
                    },
                    onMouseEnter: (e) => {
                      e.currentTarget.style.background = "#45a049";
                    },
                    onMouseLeave: (e) => {
                      e.currentTarget.style.background = "#4CAF50";
                    },
                    children: "📈 ガントチャート"
                  }
                ),
                /* @__PURE__ */ jsx(
                  "button",
                  {
                    onClick: () => setShowHelp(!showHelp),
                    style: {
                      background: "none",
                      border: "1px solid #ccc",
                      borderRadius: "50%",
                      width: "24px",
                      height: "24px",
                      cursor: "pointer",
                      fontSize: "12px",
                      display: "flex",
                      alignItems: "center",
                      justifyContent: "center",
                      marginLeft: "8px",
                      color: "#666"
                    },
                    title: "使用方法を表示",
                    children: "?"
                  }
                )
              ]
            }
          ),
          /* @__PURE__ */ jsx("div", { style: {
            flex: 1,
            overflowY: "auto",
            minHeight: 0
          }, children: projects.length === 0 ? /* @__PURE__ */ jsx("div", { style: { padding: "20px", textAlign: "center", color: "#666" }, children: "工事データが見つかりません" }) : /* @__PURE__ */ jsx("div", { children: projects.map((project, index) => /* @__PURE__ */ jsxs(
            "div",
            {
              style: {
                padding: "15px",
                borderBottom: index < projects.length - 1 ? "1px solid #eee" : "none",
                display: "flex",
                justifyContent: "space-between",
                alignItems: "center",
                cursor: "pointer",
                transition: "background-color 0.3s"
              },
              onClick: () => handleProjectClick(project),
              onMouseEnter: (e) => e.currentTarget.style.backgroundColor = "#f8f9fa",
              onMouseLeave: (e) => e.currentTarget.style.backgroundColor = "transparent",
              title: "クリックして編集",
              children: [
                /* @__PURE__ */ jsxs("div", { style: { display: "flex", alignItems: "center", gap: "16px", width: "100%" }, children: [
                  /* @__PURE__ */ jsx("div", { style: {
                    fontWeight: "600",
                    fontSize: "14px",
                    color: "#fff",
                    backgroundColor: "#1976d2",
                    padding: "3px 8px",
                    borderRadius: "4px",
                    minWidth: "90px",
                    textAlign: "center"
                  }, children: project.start_date ? new Date(
                    typeof project.start_date === "string" ? project.start_date : project.start_date["time.Time"]
                  ).toLocaleDateString("ja-JP") : "未設定" }),
                  /* @__PURE__ */ jsx("div", { style: {
                    fontWeight: "600",
                    fontSize: "16px",
                    minWidth: "120px"
                  }, children: project.company_name || "会社名未設定" }),
                  /* @__PURE__ */ jsx("div", { style: {
                    fontWeight: "600",
                    fontSize: "16px",
                    minWidth: "120px"
                  }, children: project.location_name || "現場名未設定" }),
                  /* @__PURE__ */ jsx("div", { style: { flex: 1 } }),
                  /* @__PURE__ */ jsxs("div", { style: {
                    fontSize: "14px",
                    color: "#fff",
                    backgroundColor: "#666",
                    padding: "3px 8px",
                    borderRadius: "4px",
                    minWidth: "90px",
                    textAlign: "center",
                    marginRight: "24px"
                  }, children: [
                    "～",
                    project.end_date ? new Date(
                      typeof project.end_date === "string" ? project.end_date : project.end_date["time.Time"]
                    ).toLocaleDateString("ja-JP") : "未設定"
                  ] })
                ] }),
                /* @__PURE__ */ jsxs("div", { style: { display: "flex", alignItems: "center", gap: "8px" }, children: [
                  needsFileRename(project) && /* @__PURE__ */ jsx(
                    "span",
                    {
                      style: {
                        fontSize: "16px",
                        filter: "drop-shadow(0 0 3px rgba(0, 0, 0, 0.5))"
                      },
                      title: "管理ファイルの名前変更が必要です",
                      children: "⚠️"
                    }
                  ),
                  /* @__PURE__ */ jsx(
                    "div",
                    {
                      style: {
                        padding: "4px 16px",
                        borderRadius: "4px",
                        backgroundColor: project.status === "進行中" ? "#4CAF50" : project.status === "完了" ? "#9E9E9E" : project.status === "予定" ? "#FF9800" : "#2196F3",
                        color: "white",
                        fontSize: "12px",
                        minWidth: "80px",
                        textAlign: "center"
                      },
                      children: project.status || "未設定"
                    }
                  )
                ] })
              ]
            },
            project.id || index
          )) }) })
        ]
      }
    ),
    /* @__PURE__ */ jsx(
      ProjectDetailModal,
      {
        isOpen: isEditModalOpen,
        onClose: closeEditModal,
        project: selectedProject,
        onUpdate: updateProject,
        onProjectUpdate: handleProjectUpdate
      }
    )
  ] });
};
const _layout_projects = UNSAFE_withComponentProps(function ProjectsPage() {
  return /* @__PURE__ */ jsx("div", {
    className: "container mx-auto p-4",
    children: /* @__PURE__ */ jsx("div", {
      className: "bg-white border border-gray-200 rounded-lg",
      children: /* @__PURE__ */ jsx("div", {
        className: "p-6",
        children: /* @__PURE__ */ jsx(Projects, {})
      })
    })
  });
});
const route4 = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
  __proto__: null,
  default: _layout_projects
}, Symbol.toStringTag, { value: "Module" }));
const ProjectGanttChart = () => {
  const [projects, setProjects] = useState([]);
  const [ganttItems, setGanttItems] = useState([]);
  const [selectedProject, setSelectedProject] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [viewStartDate, setViewStartDate] = useState(/* @__PURE__ */ new Date());
  const [viewEndDate, setViewEndDate] = useState(/* @__PURE__ */ new Date());
  const [isDetailModalOpen, setIsDetailModalOpen] = useState(false);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [visibleProjects, setVisibleProjects] = useState([]);
  const scrollContainerRef = useRef(null);
  const ITEMS_PER_PAGE = 10;
  const DAY_WIDTH = 10;
  const ROW_HEIGHT = 40;
  const loadProjects = async () => {
    try {
      setLoading(true);
      setError(null);
      console.log("Loading projects...");
      const response = await getProjectRecent();
      console.log("API response:", response);
      const projects2 = response.data || [];
      console.log("Projects:", projects2);
      setProjects(projects2);
    } catch (err) {
      console.error("Error loading projects:", err);
      setError(`工事データの読み込みに失敗しました: ${err instanceof Error ? err.message : "Unknown error"}`);
    } finally {
      setLoading(false);
    }
  };
  useEffect(() => {
    loadProjects();
    setHasInitialScrolled(false);
  }, []);
  const [hasInitialScrolled, setHasInitialScrolled] = useState(false);
  useEffect(() => {
    if (scrollContainerRef.current && viewStartDate && viewEndDate && ganttItems.length > 0 && !hasInitialScrolled) {
      setTimeout(() => {
        if (scrollContainerRef.current) {
          const todayX = ((/* @__PURE__ */ new Date()).getTime() - viewStartDate.getTime()) / (1e3 * 60 * 60 * 24) * DAY_WIDTH;
          const containerWidth = scrollContainerRef.current.clientWidth;
          const scrollPosition = Math.max(0, todayX - containerWidth / 2);
          scrollContainerRef.current.scrollLeft = scrollPosition;
          setHasInitialScrolled(true);
        }
      }, 300);
    }
  }, [viewStartDate, viewEndDate, ganttItems, hasInitialScrolled]);
  useEffect(() => {
    if (projects.length === 0) return;
    let minDate = /* @__PURE__ */ new Date();
    let maxDate = /* @__PURE__ */ new Date();
    let hasValidDate = false;
    projects.forEach((project) => {
      try {
        const startDate = project.start_date ? new Date(project.start_date) : null;
        const endDate = project.end_date ? new Date(project.end_date) : null;
        if (startDate && !isNaN(startDate.getTime())) {
          if (!hasValidDate || startDate < minDate) {
            minDate = startDate;
          }
          hasValidDate = true;
        }
        if (endDate && !isNaN(endDate.getTime())) {
          if (!hasValidDate || endDate > maxDate) {
            maxDate = endDate;
          }
          hasValidDate = true;
        }
      } catch (error2) {
      }
    });
    if (hasValidDate) {
      const start = new Date(minDate.getFullYear(), minDate.getMonth() - 1, 1);
      const end = new Date(maxDate.getFullYear(), maxDate.getMonth() + 1, 0);
      setViewStartDate(start);
      setViewEndDate(end);
    } else {
      const today = /* @__PURE__ */ new Date();
      const start = new Date(today.getFullYear(), today.getMonth() - 6, 1);
      const end = new Date(today.getFullYear(), today.getMonth() + 6, 0);
      setViewStartDate(start);
      setViewEndDate(end);
    }
  }, [projects]);
  const updateVisibleProjects = (scrollLeft = 0) => {
    if (projects.length === 0 || !scrollContainerRef.current) return;
    const containerWidth = scrollContainerRef.current.clientWidth;
    const visibleStartDays = scrollLeft / DAY_WIDTH;
    const visibleEndDays = (scrollLeft + containerWidth) / DAY_WIDTH;
    const visibleStartDate = new Date(viewStartDate.getTime() + visibleStartDays * 24 * 60 * 60 * 1e3);
    const visibleEndDate = new Date(viewStartDate.getTime() + visibleEndDays * 24 * 60 * 60 * 1e3);
    const relevantProjects = projects.filter((project) => {
      try {
        const projectStart = project.start_date ? new Date(project.start_date) : /* @__PURE__ */ new Date();
        const projectEnd = project.end_date ? new Date(project.end_date) : new Date(projectStart.getTime() + 90 * 24 * 60 * 60 * 1e3);
        return projectStart <= visibleEndDate && projectEnd >= visibleStartDate;
      } catch {
        return false;
      }
    });
    const sortedRelevantProjects = relevantProjects.sort((a, b) => {
      const dateA = a.start_date ? new Date(a.start_date).getTime() : 0;
      const dateB = b.start_date ? new Date(b.start_date).getTime() : 0;
      return dateA - dateB;
    });
    let baselineDate;
    if (sortedRelevantProjects.length > 0) {
      baselineDate = sortedRelevantProjects[0].start_date ? new Date(sortedRelevantProjects[0].start_date).getTime() : 0;
    } else {
      const visibleStartTime = visibleStartDate.getTime();
      const projectsBeforeVisible = projects.filter((project) => {
        const projectStartDate = project.start_date ? new Date(project.start_date).getTime() : 0;
        return projectStartDate <= visibleStartTime;
      });
      if (projectsBeforeVisible.length > 0) {
        const closestProject = projectsBeforeVisible.sort((a, b) => {
          const dateA = a.start_date ? new Date(a.start_date).getTime() : 0;
          const dateB = b.start_date ? new Date(b.start_date).getTime() : 0;
          return dateB - dateA;
        })[0];
        baselineDate = closestProject.start_date ? new Date(closestProject.start_date).getTime() : 0;
      } else {
        const allProjectsSorted2 = [...projects].sort((a, b) => {
          const dateA = a.start_date ? new Date(a.start_date).getTime() : 0;
          const dateB = b.start_date ? new Date(b.start_date).getTime() : 0;
          return dateA - dateB;
        });
        if (allProjectsSorted2.length === 0) return;
        baselineDate = allProjectsSorted2[0].start_date ? new Date(allProjectsSorted2[0].start_date).getTime() : 0;
      }
    }
    const allProjectsSorted = [...projects].sort((a, b) => {
      const dateA = a.start_date ? new Date(a.start_date).getTime() : 0;
      const dateB = b.start_date ? new Date(b.start_date).getTime() : 0;
      return dateA - dateB;
    });
    const projectsFromBaselineDate = allProjectsSorted.filter((project) => {
      const projectStartDate = project.start_date ? new Date(project.start_date).getTime() : 0;
      return projectStartDate >= baselineDate;
    });
    let finalProjects = projectsFromBaselineDate.slice(0, ITEMS_PER_PAGE);
    if (finalProjects.length < ITEMS_PER_PAGE) {
      const allProjectsDescending = [...projects].sort((a, b) => {
        const dateA = a.start_date ? new Date(a.start_date).getTime() : 0;
        const dateB = b.start_date ? new Date(b.start_date).getTime() : 0;
        return dateB - dateA;
      });
      finalProjects = allProjectsDescending.slice(0, ITEMS_PER_PAGE).sort((a, b) => {
        const dateA = a.start_date ? new Date(a.start_date).getTime() : 0;
        const dateB = b.start_date ? new Date(b.start_date).getTime() : 0;
        return dateA - dateB;
      });
    }
    setVisibleProjects(finalProjects);
  };
  const handleScroll = () => {
    if (scrollContainerRef.current) {
      updateVisibleProjects(scrollContainerRef.current.scrollLeft);
    }
  };
  useEffect(() => {
    if (projects.length > 0 && scrollContainerRef.current) {
      updateVisibleProjects(scrollContainerRef.current.scrollLeft);
    }
  }, [projects, viewStartDate]);
  useEffect(() => {
    const handleResize = () => {
      if (scrollContainerRef.current) {
        updateVisibleProjects(scrollContainerRef.current.scrollLeft);
      }
    };
    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, [projects, viewStartDate]);
  useEffect(() => {
    if (visibleProjects.length === 0) return;
    const items = visibleProjects.map((project, index) => {
      let startDate;
      let endDate;
      try {
        startDate = project.start_date ? new Date(project.start_date) : /* @__PURE__ */ new Date();
        if (isNaN(startDate.getTime())) {
          startDate = /* @__PURE__ */ new Date();
        }
      } catch {
        startDate = /* @__PURE__ */ new Date();
      }
      try {
        endDate = project.end_date ? new Date(project.end_date) : new Date(startDate.getTime() + 90 * 24 * 60 * 60 * 1e3);
        if (isNaN(endDate.getTime())) {
          endDate = new Date(startDate.getTime() + 90 * 24 * 60 * 60 * 1e3);
        }
      } catch {
        endDate = new Date(startDate.getTime() + 90 * 24 * 60 * 60 * 1e3);
      }
      const daysDiff = (startDate.getTime() - viewStartDate.getTime()) / (1e3 * 60 * 60 * 24);
      const startX = Math.max(0, daysDiff * DAY_WIDTH);
      const endDaysDiff = (endDate.getTime() - viewStartDate.getTime()) / (1e3 * 60 * 60 * 24);
      const endX = endDaysDiff * DAY_WIDTH;
      const width = Math.max(DAY_WIDTH, endX - startX);
      return {
        ...project,
        startX: isNaN(startX) ? 0 : startX,
        width: isNaN(width) ? DAY_WIDTH : width,
        row: index
      };
    });
    setGanttItems(items);
  }, [visibleProjects, viewStartDate]);
  const handleProjectEdit = (project) => {
    setSelectedProject(project);
    setIsEditModalOpen(true);
  };
  const handleProjectNameClick = (project) => {
    if (!scrollContainerRef.current) return;
    try {
      const projectStart = project.start_date ? new Date(project.start_date) : /* @__PURE__ */ new Date();
      const projectEnd = project.end_date ? new Date(project.end_date) : new Date(projectStart.getTime() + 90 * 24 * 60 * 60 * 1e3);
      const projectMiddle = new Date((projectStart.getTime() + projectEnd.getTime()) / 2);
      const middleX = (projectMiddle.getTime() - viewStartDate.getTime()) / (1e3 * 60 * 60 * 24) * DAY_WIDTH;
      const containerWidth = scrollContainerRef.current.clientWidth;
      const scrollPosition = Math.max(0, middleX - containerWidth / 2);
      scrollContainerRef.current.scrollLeft = scrollPosition;
    } catch (error2) {
      console.error("Error calculating project center:", error2);
    }
  };
  const scrollToToday = () => {
    if (!scrollContainerRef.current) return;
    const todayX = ((/* @__PURE__ */ new Date()).getTime() - viewStartDate.getTime()) / (1e3 * 60 * 60 * 24) * DAY_WIDTH;
    const containerWidth = scrollContainerRef.current.clientWidth;
    const scrollPosition = Math.max(0, todayX - containerWidth / 2);
    scrollContainerRef.current.scrollLeft = scrollPosition;
  };
  const updateProject = async (updatedProject) => {
    try {
      const response = await fetch("http://localhost:8080/api/project/update", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify(updatedProject)
      });
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || "更新に失敗しました");
      }
      const savedProject = await response.json();
      setProjects(
        (prevProjects) => prevProjects.map((p) => p.id === savedProject.id ? savedProject : p)
      );
      return savedProject;
    } catch (err) {
      console.error("Error updating project:", err);
      throw err;
    }
  };
  const handleProjectUpdate = (updatedProject) => {
    setProjects((prevProjects) => {
      const existingIndex = prevProjects.findIndex((p) => p.id === updatedProject.id);
      if (existingIndex !== -1) {
        const updatedProjects = [...prevProjects];
        updatedProjects[existingIndex] = updatedProject;
        return updatedProjects;
      }
      if (selectedProject && selectedProject.id !== updatedProject.id) {
        const oldProjectIndex = prevProjects.findIndex((p) => p.id === selectedProject.id);
        if (oldProjectIndex !== -1) {
          const updatedProjects = [...prevProjects];
          updatedProjects.splice(oldProjectIndex, 1);
          updatedProjects.push(updatedProject);
          return updatedProjects.sort((a, b) => {
            const dateA = a.start_date ? new Date(typeof a.start_date === "string" ? a.start_date : a.start_date["time.Time"]).getTime() : 0;
            const dateB = b.start_date ? new Date(typeof b.start_date === "string" ? b.start_date : b.start_date["time.Time"]).getTime() : 0;
            if (dateA > 0 && dateB === 0) return -1;
            if (dateA === 0 && dateB > 0) return 1;
            if (dateA > 0 && dateB > 0) return dateA - dateB;
            return (a.name || "").localeCompare(b.name || "");
          });
        }
      }
      return [...prevProjects, updatedProject];
    });
    setSelectedProject(updatedProject);
  };
  const getStatusColor = (status) => {
    switch (status) {
      case "進行中":
        return "#4CAF50";
      case "完了":
        return "#9E9E9E";
      case "予定":
        return "#FF9800";
      default:
        return "#2196F3";
    }
  };
  const needsFileRename = (project) => {
    if (!project.managed_files || project.managed_files.length === 0) {
      return false;
    }
    const needsRename = project.managed_files.some((file) => {
      return file.current && file.recommended && file.current !== file.recommended;
    });
    return needsRename;
  };
  const generateMonthHeaders = () => {
    const headers = [];
    const current = new Date(viewStartDate);
    while (current <= viewEndDate) {
      const monthStart = new Date(current.getFullYear(), current.getMonth(), 1);
      const monthEnd = new Date(current.getFullYear(), current.getMonth() + 1, 0);
      const daysInMonth = monthEnd.getDate();
      headers.push({
        year: monthStart.getFullYear(),
        month: monthStart.getMonth() + 1,
        width: daysInMonth * DAY_WIDTH,
        startX: (monthStart.getTime() - viewStartDate.getTime()) / (1e3 * 60 * 60 * 24) * DAY_WIDTH
      });
      current.setMonth(current.getMonth() + 1);
    }
    return headers;
  };
  const generateDayHeaders = () => {
    const headers = [];
    const current = new Date(viewStartDate);
    while (current <= viewEndDate) {
      const day = current.getDate();
      if (day === 1 || day % 3 === 1) {
        headers.push({
          date: day,
          month: current.getMonth() + 1,
          year: current.getFullYear(),
          startX: (current.getTime() - viewStartDate.getTime()) / (1e3 * 60 * 60 * 24) * DAY_WIDTH,
          width: DAY_WIDTH * 3
          // 3日分の幅
        });
      }
      current.setDate(current.getDate() + 1);
    }
    return headers;
  };
  const generateMonthBoundaries = () => {
    const boundaries = [];
    const current = new Date(viewStartDate.getFullYear(), viewStartDate.getMonth() + 1, 1);
    while (current <= viewEndDate) {
      const startX = (current.getTime() - viewStartDate.getTime()) / (1e3 * 60 * 60 * 24) * DAY_WIDTH;
      boundaries.push({
        startX,
        year: current.getFullYear(),
        month: current.getMonth() + 1
      });
      current.setMonth(current.getMonth() + 1);
    }
    return boundaries;
  };
  if (loading) {
    return /* @__PURE__ */ jsx("div", { className: "loading", children: "工事データを読み込み中..." });
  }
  if (error) {
    return /* @__PURE__ */ jsx("div", { className: "error", children: error });
  }
  const monthHeaders = generateMonthHeaders();
  const dayHeaders = generateDayHeaders();
  const monthBoundaries = generateMonthBoundaries();
  const totalWidth = (viewEndDate.getTime() - viewStartDate.getTime()) / (1e3 * 60 * 60 * 24) * DAY_WIDTH;
  return /* @__PURE__ */ jsxs("div", { className: "gantt-container", children: [
    /* @__PURE__ */ jsx("h1", { children: "工程表" }),
    /* @__PURE__ */ jsxs("div", { className: "gantt-controls", children: [
      /* @__PURE__ */ jsx(
        "button",
        {
          onClick: scrollToToday,
          style: {
            padding: "8px 16px",
            background: "#FF5252",
            color: "white",
            border: "none",
            borderRadius: "4px",
            cursor: "pointer",
            fontSize: "14px"
          },
          children: "今日へ移動"
        }
      ),
      /* @__PURE__ */ jsx("div", { className: "info", children: /* @__PURE__ */ jsxs("span", { children: [
        "表示中: ",
        ganttItems.length,
        "件 / 全",
        projects.length,
        "件"
      ] }) })
    ] }),
    /* @__PURE__ */ jsxs("div", { className: "gantt-wrapper", children: [
      /* @__PURE__ */ jsxs("div", { className: "gantt-sidebar", children: [
        /* @__PURE__ */ jsx("div", { className: "gantt-header-left", children: "工事名" }),
        ganttItems.map((item, index) => {
          const today = /* @__PURE__ */ new Date();
          let startDate;
          let endDate;
          try {
            startDate = item.start_date ? new Date(item.start_date) : /* @__PURE__ */ new Date();
            if (isNaN(startDate.getTime())) startDate = /* @__PURE__ */ new Date();
          } catch {
            startDate = /* @__PURE__ */ new Date();
          }
          try {
            endDate = item.end_date ? new Date(item.end_date) : new Date(startDate.getTime() + 90 * 24 * 60 * 60 * 1e3);
            if (isNaN(endDate.getTime())) endDate = new Date(startDate.getTime() + 90 * 24 * 60 * 60 * 1e3);
          } catch {
            endDate = new Date(startDate.getTime() + 90 * 24 * 60 * 60 * 1e3);
          }
          const isActiveProject = today >= startDate && today <= endDate;
          return /* @__PURE__ */ jsx(
            "div",
            {
              className: "gantt-row-label",
              style: {
                height: ROW_HEIGHT,
                backgroundColor: isActiveProject ? "#fff3cd" : "transparent",
                borderLeft: isActiveProject ? "4px solid #ffc107" : "none"
              },
              children: /* @__PURE__ */ jsxs(
                "div",
                {
                  className: "project-name",
                  style: {
                    fontWeight: isActiveProject ? "bold" : "normal",
                    color: isActiveProject ? "#856404" : "inherit",
                    cursor: "pointer"
                  },
                  onClick: () => handleProjectNameClick(item),
                  title: "クリックしてプロジェクト期間の中央に移動",
                  children: [
                    item.company_name,
                    " - ",
                    item.location_name
                  ]
                }
              )
            },
            `${item.id}-${index}`
          );
        })
      ] }),
      /* @__PURE__ */ jsx("div", { className: "gantt-chart-container", ref: scrollContainerRef, onScroll: handleScroll, style: { backgroundColor: "#f5f5f5" }, children: /* @__PURE__ */ jsxs("div", { className: "gantt-chart", style: { width: totalWidth, backgroundColor: "#f5f5f5" }, children: [
        /* @__PURE__ */ jsx("div", { className: "gantt-header month-header-row", style: { height: "30px" }, children: monthHeaders.map((header, index) => /* @__PURE__ */ jsxs(
          "div",
          {
            className: "month-header",
            style: {
              position: "absolute",
              left: header.startX,
              width: header.width,
              height: "100%",
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              backgroundColor: "#e8f4fd",
              borderRight: "1px solid #ddd",
              fontWeight: "bold",
              fontSize: "13px",
              color: "#0066cc"
            },
            children: [
              header.year,
              "年",
              header.month,
              "月"
            ]
          },
          index
        )) }),
        /* @__PURE__ */ jsx("div", { className: "gantt-header day-header-row", style: { height: "25px", borderBottom: "2px solid #333" }, children: dayHeaders.map((header, index) => /* @__PURE__ */ jsx(
          "div",
          {
            className: "day-header",
            style: {
              position: "absolute",
              left: header.startX,
              width: header.width,
              height: "100%",
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              backgroundColor: "#f5f5f5",
              borderRight: "1px solid #ddd",
              fontSize: "11px",
              fontWeight: "500"
            },
            children: header.date
          },
          index
        )) }),
        /* @__PURE__ */ jsxs("div", { className: "gantt-body", children: [
          /* @__PURE__ */ jsx("div", { className: "gantt-grid", children: dayHeaders.map((header, index) => /* @__PURE__ */ jsx(
            "div",
            {
              className: "grid-line",
              style: { left: header.startX }
            },
            index
          )) }),
          /* @__PURE__ */ jsx("div", { className: "gantt-month-boundaries", children: monthBoundaries.map((boundary, index) => /* @__PURE__ */ jsx(
            "div",
            {
              className: "month-boundary-line",
              style: { left: boundary.startX },
              title: `${boundary.year}年${boundary.month}月開始`
            },
            `month-boundary-${index}`
          )) }),
          /* @__PURE__ */ jsx("div", { className: "gantt-horizontal-grid", children: ganttItems.map((_, index) => /* @__PURE__ */ jsx(
            "div",
            {
              className: "horizontal-grid-line",
              style: {
                top: (index + 1) * ROW_HEIGHT,
                width: "100%"
              }
            },
            `horizontal-${index}`
          )) }),
          ganttItems.map((item, index) => /* @__PURE__ */ jsxs(
            "div",
            {
              className: "gantt-bar",
              style: {
                left: item.startX,
                width: item.width,
                top: index * ROW_HEIGHT + 5,
                height: ROW_HEIGHT - 10,
                backgroundColor: getStatusColor(item.status)
              },
              onClick: () => handleProjectEdit(item),
              title: `${item.company_name} - ${item.location_name} (クリックして編集)`,
              children: [
                /* @__PURE__ */ jsx("span", { className: "gantt-bar-text", children: item.location_name }),
                needsFileRename(item) && /* @__PURE__ */ jsx(
                  "span",
                  {
                    className: "gantt-bar-rename-indicator",
                    title: "管理ファイルの名前変更が必要です",
                    children: "⚠️"
                  }
                )
              ]
            },
            `${item.id}-${index}`
          )),
          /* @__PURE__ */ jsx(
            "div",
            {
              className: "today-area",
              style: {
                left: Math.floor(((/* @__PURE__ */ new Date()).setHours(0, 0, 0, 0) - viewStartDate.getTime()) / (1e3 * 60 * 60 * 24)) * DAY_WIDTH,
                width: DAY_WIDTH,
                height: "100%",
                backgroundColor: "rgba(255, 192, 203, 0.3)",
                // 薄いピンク
                position: "absolute",
                top: 0,
                pointerEvents: "none",
                zIndex: 1
              }
            }
          )
        ] })
      ] }) })
    ] }),
    /* @__PURE__ */ jsx(
      ProjectDetailModal,
      {
        isOpen: isEditModalOpen,
        onClose: () => setIsEditModalOpen(false),
        project: selectedProject,
        onUpdate: updateProject,
        onProjectUpdate: handleProjectUpdate
      }
    )
  ] });
};
const _layout_gantt = UNSAFE_withComponentProps(function GanttChartPage() {
  return /* @__PURE__ */ jsx("div", {
    className: "container mx-auto p-4",
    children: /* @__PURE__ */ jsx("div", {
      className: "bg-white border border-gray-200 rounded-lg",
      children: /* @__PURE__ */ jsx("div", {
        className: "p-6",
        children: /* @__PURE__ */ jsx(ProjectGanttChart, {})
      })
    })
  });
});
const route5 = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
  __proto__: null,
  default: _layout_gantt
}, Symbol.toStringTag, { value: "Module" }));
const serverManifest = { "entry": { "module": "/assets/entry.client-BR_5jt8X.js", "imports": ["/assets/chunk-QMGIS6GS-CRLmNMPd.js", "/assets/index-BXUbddt1.js"], "css": [] }, "routes": { "root": { "id": "root", "parentId": void 0, "path": "", "index": void 0, "caseSensitive": void 0, "hasAction": false, "hasLoader": false, "hasClientAction": false, "hasClientLoader": false, "hasClientMiddleware": false, "hasErrorBoundary": false, "module": "/assets/root-CmFF9zS-.js", "imports": ["/assets/chunk-QMGIS6GS-CRLmNMPd.js", "/assets/index-BXUbddt1.js"], "css": [], "clientActionModule": void 0, "clientLoaderModule": void 0, "clientMiddlewareModule": void 0, "hydrateFallbackModule": void 0 }, "routes/_layout": { "id": "routes/_layout", "parentId": "root", "path": void 0, "index": void 0, "caseSensitive": void 0, "hasAction": false, "hasLoader": false, "hasClientAction": false, "hasClientLoader": false, "hasClientMiddleware": false, "hasErrorBoundary": false, "module": "/assets/_layout-BGsERk8U.js", "imports": ["/assets/chunk-QMGIS6GS-CRLmNMPd.js", "/assets/FileInfoContext-NxQDI7MG.js", "/assets/ProjectContext-DQ4aDs3-.js"], "css": [], "clientActionModule": void 0, "clientLoaderModule": void 0, "clientMiddlewareModule": void 0, "hydrateFallbackModule": void 0 }, "routes/_layout._index": { "id": "routes/_layout._index", "parentId": "routes/_layout", "path": void 0, "index": true, "caseSensitive": void 0, "hasAction": false, "hasLoader": false, "hasClientAction": false, "hasClientLoader": false, "hasClientMiddleware": false, "hasErrorBoundary": false, "module": "/assets/_layout._index-BxvzRpYg.js", "imports": ["/assets/chunk-QMGIS6GS-CRLmNMPd.js"], "css": [], "clientActionModule": void 0, "clientLoaderModule": void 0, "clientMiddlewareModule": void 0, "hydrateFallbackModule": void 0 }, "routes/_layout.files": { "id": "routes/_layout.files", "parentId": "routes/_layout", "path": "files", "index": void 0, "caseSensitive": void 0, "hasAction": false, "hasLoader": false, "hasClientAction": false, "hasClientLoader": false, "hasClientMiddleware": false, "hasErrorBoundary": false, "module": "/assets/_layout.files-3x5L2HQM.js", "imports": ["/assets/chunk-QMGIS6GS-CRLmNMPd.js", "/assets/sdk.gen-C0JJsBdk.js", "/assets/FileInfoContext-NxQDI7MG.js", "/assets/index-BXUbddt1.js"], "css": [], "clientActionModule": void 0, "clientLoaderModule": void 0, "clientMiddlewareModule": void 0, "hydrateFallbackModule": void 0 }, "routes/_layout.projects": { "id": "routes/_layout.projects", "parentId": "routes/_layout", "path": "projects", "index": void 0, "caseSensitive": void 0, "hasAction": false, "hasLoader": false, "hasClientAction": false, "hasClientLoader": false, "hasClientMiddleware": false, "hasErrorBoundary": false, "module": "/assets/_layout.projects-BhJjnzM-.js", "imports": ["/assets/chunk-QMGIS6GS-CRLmNMPd.js", "/assets/sdk.gen-C0JJsBdk.js", "/assets/ProjectDetailModal-BfP-5M3c.js", "/assets/ProjectContext-DQ4aDs3-.js", "/assets/index-BXUbddt1.js"], "css": [], "clientActionModule": void 0, "clientLoaderModule": void 0, "clientMiddlewareModule": void 0, "hydrateFallbackModule": void 0 }, "routes/_layout.gantt": { "id": "routes/_layout.gantt", "parentId": "routes/_layout", "path": "projects/gantt", "index": void 0, "caseSensitive": void 0, "hasAction": false, "hasLoader": false, "hasClientAction": false, "hasClientLoader": false, "hasClientMiddleware": false, "hasErrorBoundary": false, "module": "/assets/_layout.gantt-yY0nxKIi.js", "imports": ["/assets/chunk-QMGIS6GS-CRLmNMPd.js", "/assets/sdk.gen-C0JJsBdk.js", "/assets/ProjectDetailModal-BfP-5M3c.js", "/assets/index-BXUbddt1.js"], "css": ["/assets/_layout-TDnedI-G.css"], "clientActionModule": void 0, "clientLoaderModule": void 0, "clientMiddlewareModule": void 0, "hydrateFallbackModule": void 0 } }, "url": "/assets/manifest-b0e28d63.js", "version": "b0e28d63", "sri": void 0 };
const assetsBuildDirectory = "build/client";
const basename = "/";
const future = { "unstable_middleware": false, "unstable_optimizeDeps": false, "unstable_splitRouteModules": false, "unstable_subResourceIntegrity": false, "unstable_viteEnvironmentApi": false };
const ssr = true;
const isSpaMode = false;
const prerender = [];
const routeDiscovery = { "mode": "lazy", "manifestPath": "/__manifest" };
const publicPath = "/";
const entry = { module: entryServer };
const routes = {
  "root": {
    id: "root",
    parentId: void 0,
    path: "",
    index: void 0,
    caseSensitive: void 0,
    module: route0
  },
  "routes/_layout": {
    id: "routes/_layout",
    parentId: "root",
    path: void 0,
    index: void 0,
    caseSensitive: void 0,
    module: route1
  },
  "routes/_layout._index": {
    id: "routes/_layout._index",
    parentId: "routes/_layout",
    path: void 0,
    index: true,
    caseSensitive: void 0,
    module: route2
  },
  "routes/_layout.files": {
    id: "routes/_layout.files",
    parentId: "routes/_layout",
    path: "files",
    index: void 0,
    caseSensitive: void 0,
    module: route3
  },
  "routes/_layout.projects": {
    id: "routes/_layout.projects",
    parentId: "routes/_layout",
    path: "projects",
    index: void 0,
    caseSensitive: void 0,
    module: route4
  },
  "routes/_layout.gantt": {
    id: "routes/_layout.gantt",
    parentId: "routes/_layout",
    path: "projects/gantt",
    index: void 0,
    caseSensitive: void 0,
    module: route5
  }
};
export {
  serverManifest as assets,
  assetsBuildDirectory,
  basename,
  entry,
  future,
  isSpaMode,
  prerender,
  publicPath,
  routeDiscovery,
  routes,
  ssr
};
