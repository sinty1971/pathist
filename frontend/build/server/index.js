import { jsx, jsxs, Fragment } from "react/jsx-runtime";
import { PassThrough } from "node:stream";
import { createReadableStreamFromReadable } from "@react-router/node";
import { ServerRouter, UNSAFE_withComponentProps, Outlet, Meta, Links, ScrollRestoration, Scripts, useLocation, Link } from "react-router";
import { isbot } from "isbot";
import { renderToPipeableStream } from "react-dom/server";
import React, { createContext, useState, useContext, useCallback, useEffect, useRef } from "react";
import { TreeItem, SimpleTreeView } from "@mui/x-tree-view";
import { Box, Typography, Chip, Paper, Toolbar, IconButton, Breadcrumbs, Link as Link$1, Alert, CircularProgress, TextField, Dialog, DialogTitle, DialogContent, Button } from "@mui/material";
import { ExpandMore, ChevronRight, InsertDriveFile, Home, Refresh } from "@mui/icons-material";
import CloseIcon from "@mui/icons-material/Close";
import WarningIcon from "@mui/icons-material/Warning";
import AttachFileIcon from "@mui/icons-material/AttachFile";
import LightbulbIcon from "@mui/icons-material/Lightbulb";
import { DatePicker } from "@mui/x-date-pickers/DatePicker";
import { LocalizationProvider } from "@mui/x-date-pickers/LocalizationProvider";
import { AdapterDateFns } from "@mui/x-date-pickers/AdapterDateFns";
import { ja } from "date-fns/locale";
import BusinessIcon from "@mui/icons-material/Business";
import PhoneIcon from "@mui/icons-material/Phone";
import EmailIcon from "@mui/icons-material/Email";
import LanguageIcon from "@mui/icons-material/Language";
import TagIcon from "@mui/icons-material/Tag";
import EditIcon from "@mui/icons-material/Edit";
import SaveIcon from "@mui/icons-material/Save";
import CancelIcon from "@mui/icons-material/Cancel";
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
function Navigation() {
  const location = useLocation();
  return /* @__PURE__ */ jsx("nav", { className: "navigation", children: /* @__PURE__ */ jsxs("div", { className: "nav-container", children: [
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
          to: "/kojies",
          className: location.pathname === "/kojies" ? "nav-link active" : "nav-link",
          children: "工事一覧"
        }
      ),
      /* @__PURE__ */ jsx(
        Link,
        {
          to: "/companies",
          className: location.pathname === "/companies" ? "nav-link active" : "nav-link",
          children: "会社一覧"
        }
      )
    ] }),
    /* @__PURE__ */ jsx("div", { className: "nav-logo", children: /* @__PURE__ */ jsx("h1", { children: "Penguin フォルダー管理" }) })
  ] }) });
}
const KojiContext = createContext(void 0);
function KojiProvider({ children }) {
  const [kojiCount, setKojiCount] = useState(0);
  return /* @__PURE__ */ jsx(KojiContext.Provider, { value: { kojiCount, setKojiCount }, children });
}
function useKoji() {
  const context = useContext(KojiContext);
  if (context === void 0) {
    throw new Error("useKoji must be used within a KojiProvider");
  }
  return context;
}
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
const links = () => [{
  rel: "stylesheet",
  href: "/app/styles/app.css"
}];
function LayoutContent() {
  const {
    kojiCount
  } = useKoji();
  return /* @__PURE__ */ jsxs("div", {
    className: "app",
    children: [/* @__PURE__ */ jsx(Navigation, {}), /* @__PURE__ */ jsx("main", {
      className: "main-content",
      children: /* @__PURE__ */ jsx(Outlet, {})
    })]
  });
}
const _layout = UNSAFE_withComponentProps(function Layout2() {
  return /* @__PURE__ */ jsx(KojiProvider, {
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
              height: "220px",
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
            to: "/kojies",
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
              children: "工事一覧"
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
            to: "/kojies/gantt",
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
            children: "📈 工程表を見る"
          })]
        }), /* @__PURE__ */ jsx(Link, {
          to: "/companies",
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
              height: "220px",
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
              children: "🏢"
            }), /* @__PURE__ */ jsx("h3", {
              style: {
                fontSize: "1.1rem",
                color: "#333",
                marginBottom: "0.5rem",
                textAlign: "center"
              },
              children: "会社一覧"
            }), /* @__PURE__ */ jsx("p", {
              style: {
                color: "#666",
                fontSize: "0.85rem",
                textAlign: "center",
                lineHeight: "1.4"
              },
              children: "取引先会社の情報管理"
            })]
          })
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
          children: ["データ保存場所: ~/penguin", /* @__PURE__ */ jsx("br", {}), "工事プロジェクト: ~/penguin/豊田築炉/2-工事", /* @__PURE__ */ jsx("br", {}), "会社情報: ~/penguin/豊田築炉/1 会社"]
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
const getBusinessBasePath = (options) => {
  return (options?.client ?? client).get({
    url: "/business/base-path",
    ...options
  });
};
const getBusinessCompanies = (options) => {
  return (options?.client ?? client).get({
    url: "/business/companies",
    ...options
  });
};
const putBusinessCompanies = (options) => {
  return (options.client ?? client).put({
    url: "/business/companies",
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options.headers
    }
  });
};
const getBusinessFiles = (options) => {
  return (options?.client ?? client).get({
    url: "/business/files",
    ...options
  });
};
const getBusinessKojies = (options) => {
  return (options?.client ?? client).get({
    url: "/business/kojies",
    ...options
  });
};
const putBusinessKojies = (options) => {
  return (options.client ?? client).put({
    url: "/business/kojies",
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options.headers
    }
  });
};
const putBusinessKojiesManagedFiles = (options) => {
  return (options.client ?? client).put({
    url: "/business/kojies/managed-files",
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options.headers
    }
  });
};
const getKojiesByPath = (options) => {
  return (options.client ?? client).get({
    url: "/kojies/{path}",
    ...options
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
const getNodeIcon = (node) => {
  if (node.isDirectory) {
    return null;
  }
  if (node.name === ".detail.yaml") {
    return /* @__PURE__ */ jsx("span", { style: { fontSize: "16px" }, children: "⚙️" });
  }
  const ext = node.name?.split(".").pop()?.toLowerCase();
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
const formatFileSize = (bytes) => {
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
const CustomTreeItem = React.memo(
  ({
    itemId,
    node,
    onNodeClick,
    onNodeExpand,
    isExpanded,
    expanded,
    ...props
  }) => {
    const handleClick = React.useCallback(
      (e) => {
        e.stopPropagation();
        onNodeClick(node);
      },
      [node, onNodeClick]
    );
    return /* @__PURE__ */ jsx(
      TreeItem,
      {
        itemId,
        onClick: handleClick,
        sx: {
          // ローディング中のノードの視覚的フィードバックのみ
          ...node.isLoading && {
            "& .MuiTreeItem-content": {
              opacity: 0.7
            }
          }
        },
        label: /* @__PURE__ */ jsxs(Box, { sx: { display: "flex", alignItems: "center", py: 0.5, pr: 2 }, children: [
          node.isDirectory ? isExpanded ? /* @__PURE__ */ jsx(ExpandMore, { sx: { mr: 0.5 } }) : /* @__PURE__ */ jsx(ChevronRight, { sx: { mr: 0.5 } }) : getNodeIcon(node),
          /* @__PURE__ */ jsxs(Typography, { variant: "body2", sx: { flexGrow: 1, mr: 1, ml: 1 }, children: [
            node.name,
            !node.isDirectory && formatFileSize(node.size)
          ] }),
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
            onNodeExpand,
            isExpanded: expanded.includes(child.id),
            expanded
          },
          child.id
        ))
      }
    );
  }
);
const Files = () => {
  const [treeData, setTreeData] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [currentPath, setCurrentPath] = useState("");
  const [basePath, setBasePath] = useState("");
  const [basePathError, setBasePathError] = useState(false);
  const [selectedNode, setSelectedNode] = useState(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [expanded, setExpanded] = useState([]);
  const { setFileCount, setCurrentPath: setContextPath } = useFileInfo();
  const loadBasePath = useCallback(async () => {
    try {
      const response = await getBusinessBasePath();
      if (response.data && response.data.businessBasePath) {
        setBasePath(response.data.businessBasePath);
        setBasePathError(false);
        return response.data.businessBasePath;
      }
      throw new Error("ベースパスが取得できませんでした");
    } catch (err) {
      setBasePathError(true);
      setError("基準となるパスを取得できないため、一覧を表示できません");
      return null;
    }
  }, []);
  const convertToTreeNode = (fileInfo) => {
    const node = {
      id: fileInfo.path,
      // バックエンドからのpathをIDとして使用
      name: fileInfo.name,
      path: fileInfo.path,
      // バックエンドからのpathをそのまま使用
      relativePath: fileInfo.path.replace(basePath + "/", "") || "",
      // 相対パス部分を保存
      isDirectory: fileInfo.is_directory,
      size: fileInfo.size,
      modifiedTime: fileInfo.modified_time,
      children: fileInfo.is_directory ? [] : void 0,
      isLoaded: !fileInfo.is_directory,
      isLoading: false
    };
    return node;
  };
  const loadFiles = useCallback(
    async (relativePath, isRefresh = false) => {
      const requestPath = relativePath || "";
      setLoading(true);
      setError(null);
      try {
        const response = await getBusinessFiles({
          query: requestPath ? { path: requestPath } : {}
        });
        if (response.data) {
          const data = response.data;
          const nodes = data.map((fileInfo) => {
            return convertToTreeNode(fileInfo);
          });
          if (!relativePath || isRefresh) {
            setTreeData(nodes);
            setFileCount(data.length);
            setContextPath(basePath);
            setCurrentPath(basePath);
          }
          return nodes;
        } else if (response.error) {
          throw new Error("APIエラー: " + JSON.stringify(response.error));
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : "エラーが発生しました");
        return [];
      } finally {
        setLoading(false);
      }
    },
    [basePath, setFileCount, setContextPath]
  );
  useEffect(() => {
    const initialize = async () => {
      setLoading(true);
      const basePathResult = await loadBasePath();
      if (basePathResult) {
        await loadFiles();
      }
      setLoading(false);
    };
    initialize();
  }, [loadBasePath, loadFiles]);
  const handleNodeExpand = useCallback(
    async (nodeId, node) => {
      if (!node.isDirectory || node.isLoaded || node.isLoading) return;
      setTreeData(
        (prevData) => updateNodeInTree(prevData, nodeId, { ...node, isLoading: true })
      );
      try {
        const targetRelativePath = node.relativePath || "";
        const children = await loadFiles(targetRelativePath, false);
        setTreeData(
          (prevData) => updateNodeInTree(prevData, nodeId, {
            ...node,
            children,
            isLoaded: true,
            isLoading: false
          })
        );
      } catch (err) {
        console.error("handleNodeExpand error:", err);
        setTreeData(
          (prevData) => updateNodeInTree(prevData, nodeId, { ...node, isLoading: false })
        );
      }
    },
    [loadFiles]
  );
  const handleNodeClick = useCallback(
    (node) => {
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
    },
    [handleNodeExpand]
  );
  const handleRefresh = () => {
    setTreeData([]);
    setExpanded([]);
    loadFiles(currentPath, true);
  };
  const handleGoHome = () => {
    setTreeData([]);
    setExpanded([]);
    loadFiles("", true);
  };
  const getBreadcrumbs = () => {
    const currentBasePath = currentPath || basePath;
    const parts = currentBasePath.replace(basePath, "").split("/").filter(Boolean);
    const breadcrumbs = [
      { label: basePath.split("/").pop() || "ホーム", path: "" }
    ];
    let accumulatedPath = "";
    parts.forEach((part) => {
      if (accumulatedPath === "") {
        accumulatedPath = `${basePath}/${part}`;
      } else {
        accumulatedPath += `/${part}`;
      }
      breadcrumbs.push({ label: part, path: accumulatedPath });
    });
    return breadcrumbs;
  };
  return /* @__PURE__ */ jsxs(Box, { sx: { height: "100%", display: "flex", flexDirection: "column" }, children: [
    /* @__PURE__ */ jsxs(Paper, { elevation: 1, sx: { mb: 1 }, children: [
      /* @__PURE__ */ jsxs(Toolbar, { variant: "dense", children: [
        /* @__PURE__ */ jsx(IconButton, { onClick: handleGoHome, size: "small", title: "ホームに戻る", children: /* @__PURE__ */ jsx(Home, {}) }),
        /* @__PURE__ */ jsx(IconButton, { onClick: handleRefresh, size: "small", title: "リフレッシュ", children: /* @__PURE__ */ jsx(Refresh, {}) })
      ] }),
      /* @__PURE__ */ jsx(Box, { sx: { px: 2, pb: 1 }, children: /* @__PURE__ */ jsx(Breadcrumbs, { children: getBreadcrumbs().map((crumb, index) => /* @__PURE__ */ jsx(
        Link$1,
        {
          component: "button",
          variant: "body2",
          color: index === getBreadcrumbs().length - 1 ? "text.primary" : "inherit",
          onClick: () => {
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
    /* @__PURE__ */ jsxs(
      Paper,
      {
        sx: {
          flex: 1,
          overflow: "auto",
          p: 1,
          position: "relative",
          minHeight: 0
        },
        children: [
          loading && treeData.length === 0 ? /* @__PURE__ */ jsxs(
            Box,
            {
              sx: {
                display: "flex",
                flexDirection: "column",
                alignItems: "center",
                justifyContent: "center",
                height: "100%",
                p: 3
              },
              children: [
                /* @__PURE__ */ jsxs(Box, { sx: { position: "relative", display: "inline-flex" }, children: [
                  /* @__PURE__ */ jsx(
                    CircularProgress,
                    {
                      size: 60,
                      thickness: 1,
                      sx: {
                        color: "primary.main",
                        animationDuration: "2s"
                      }
                    }
                  ),
                  /* @__PURE__ */ jsx(
                    CircularProgress,
                    {
                      variant: "determinate",
                      size: 60,
                      thickness: 2,
                      value: 25,
                      sx: {
                        color: "grey.300",
                        position: "absolute",
                        top: 0,
                        left: 0,
                        zIndex: 0
                      }
                    }
                  )
                ] }),
                /* @__PURE__ */ jsx(Typography, { color: "text.secondary", sx: { mt: 2 }, children: "データ取得中..." }),
                /* @__PURE__ */ jsx(
                  Typography,
                  {
                    variant: "caption",
                    color: "text.secondary",
                    sx: { mt: 0.5 },
                    children: "大量のファイルがある場合は時間がかかることがあります"
                  }
                )
              ]
            }
          ) : basePathError ? /* @__PURE__ */ jsx(Box, { sx: { display: "flex", justifyContent: "center", p: 3 }, children: /* @__PURE__ */ jsx(Typography, { color: "error", children: "基準となるパスを取得できないため、一覧を表示できません" }) }) : /* @__PURE__ */ jsxs(Fragment, { children: [
            /* @__PURE__ */ jsx(
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
                  },
                  "& .MuiTreeItem-iconContainer": {
                    display: "none"
                    // MUIのデフォルト展開アイコンを非表示（カスタムアイコンを使用）
                  }
                },
                children: treeData.map((node) => /* @__PURE__ */ jsx(
                  CustomTreeItem,
                  {
                    itemId: node.id,
                    node,
                    onNodeClick: handleNodeClick,
                    onNodeExpand: handleNodeExpand,
                    isExpanded: expanded.includes(node.id),
                    expanded
                  },
                  node.id
                ))
              }
            ),
            treeData.some(
              (node) => node.isLoading || node.children && node.children.some((child) => child.isLoading)
            ) && /* @__PURE__ */ jsxs(
              Box,
              {
                sx: {
                  position: "absolute",
                  top: "50%",
                  left: "50%",
                  transform: "translate(-50%, -50%)",
                  backgroundColor: "rgba(255, 255, 255, 0.95)",
                  borderRadius: 2,
                  p: 3,
                  boxShadow: 3,
                  zIndex: 1e3,
                  display: "flex",
                  flexDirection: "column",
                  alignItems: "center",
                  minWidth: 200
                },
                children: [
                  /* @__PURE__ */ jsx(
                    CircularProgress,
                    {
                      size: 48,
                      thickness: 3,
                      sx: {
                        color: "primary.main"
                      }
                    }
                  ),
                  /* @__PURE__ */ jsx(
                    Typography,
                    {
                      color: "text.primary",
                      sx: { mt: 2, fontWeight: 500 },
                      children: "フォルダーを読み込んでいます..."
                    }
                  ),
                  /* @__PURE__ */ jsx(
                    Typography,
                    {
                      variant: "caption",
                      color: "text.secondary",
                      sx: { mt: 0.5 },
                      children: "大量のファイルがある場合は時間がかかります"
                    }
                  )
                ]
              }
            )
          ] }),
          treeData.length === 0 && !loading && !basePathError && /* @__PURE__ */ jsx(Box, { sx: { display: "flex", justifyContent: "center", p: 3 }, children: /* @__PURE__ */ jsx(Typography, { color: "text.secondary", children: "フォルダーが空です" }) })
        ]
      }
    ),
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
const KojiDetailModal = ({ isOpen, onClose, koji, onUpdate, onKojiUpdate }) => {
  const [formData, setFormData] = useState({
    id: "",
    company_name: "",
    location_name: "",
    description: "",
    tags: "",
    start_date: "",
    end_date: ""
  });
  const [currentKoji, setCurrentKoji] = useState(null);
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
    if (koji) {
      const startDate = extractDateString(koji.start_date);
      const endDate = extractDateString(koji.end_date);
      const companyName = koji.company_name || "";
      const locationName = koji.location_name || "";
      setCurrentKoji(koji);
      setFormData({
        id: koji.id || "",
        company_name: companyName,
        location_name: locationName,
        description: koji.description || "",
        tags: Array.isArray(koji.tags) ? koji.tags.join(", ") : koji.tags || "",
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
  }, [koji]);
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
    if (!currentKoji?.managed_files) return;
    const updatedManagedFiles = currentKoji.managed_files.map((file) => {
      if (file.current) {
        const recommendedName = generateRecommendedFileName(file.current, formData2);
        return {
          ...file,
          recommended: recommendedName
        };
      }
      return file;
    });
    setCurrentKoji((prev) => prev ? {
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
    if (!koji) return;
    setIsLoading(true);
    setError(null);
    try {
      const updatedKoji = {
        ...koji,
        company_name: useFormData.company_name,
        location_name: useFormData.location_name,
        description: useFormData.description,
        tags: useFormData.tags ? useFormData.tags.split(",").map((tag) => tag.trim()).filter((tag) => tag.length > 0) : [],
        start_date: useFormData.start_date ? { "time.Time": `${useFormData.start_date}T00:00:00+09:00` } : void 0,
        end_date: useFormData.end_date ? { "time.Time": `${useFormData.end_date}T23:59:59+09:00` } : void 0
      };
      const originalFolderName = koji.name;
      const savedKoji = await onUpdate(updatedKoji);
      const folderNameChanged = originalFolderName && savedKoji.name && originalFolderName !== savedKoji.name;
      setCurrentKoji(savedKoji);
      if (onKojiUpdate) {
        onKojiUpdate(savedKoji);
      }
      const startDate = extractDateString(savedKoji.start_date);
      const endDate = extractDateString(savedKoji.end_date);
      setFormData({
        id: savedKoji.id || "",
        company_name: savedKoji.company_name || "",
        location_name: savedKoji.location_name || "",
        description: savedKoji.description || "",
        tags: Array.isArray(savedKoji.tags) ? savedKoji.tags.join(", ") : savedKoji.tags || "",
        start_date: startDate,
        end_date: endDate
      });
      if (folderNameChanged) {
        setTimeout(() => {
          onClose();
        }, 100);
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
    if (!currentKoji || !currentKoji.managed_files) return;
    setIsRenaming(true);
    setError(null);
    try {
      const currentFiles = currentKoji.managed_files.filter((file) => file.current && file.recommended && file.current !== file.recommended).map((file) => file.current);
      if (currentFiles.length === 0) {
        setError("変更対象のファイルがありません");
        return;
      }
      const response = await putBusinessKojiesManagedFiles({
        body: {
          koji: currentKoji,
          currents: currentFiles
        }
      });
      if (response.data) {
        if (typeof response.data === "object" && response.data.id) {
          const updatedKoji = response.data;
          setCurrentKoji(updatedKoji);
          if (onKojiUpdate) {
            onKojiUpdate(updatedKoji);
          }
          const startDate = extractDateString(updatedKoji.start_date);
          const endDate = extractDateString(updatedKoji.end_date);
          setFormData({
            id: updatedKoji.id || "",
            company_name: updatedKoji.company_name || "",
            location_name: updatedKoji.location_name || "",
            description: updatedKoji.description || "",
            tags: Array.isArray(updatedKoji.tags) ? updatedKoji.tags.join(", ") : updatedKoji.tags || "",
            start_date: startDate,
            end_date: endDate
          });
        } else {
          if (currentKoji.name) {
            const updatedKojiResponse = await getKojiesByPath({
              query: {
                path: currentKoji.name
              }
            });
            if (updatedKojiResponse.data && onKojiUpdate) {
              setCurrentKoji(updatedKojiResponse.data);
              onKojiUpdate(updatedKojiResponse.data);
              const startDate = extractDateString(updatedKojiResponse.data.start_date);
              const endDate = extractDateString(updatedKojiResponse.data.end_date);
              setFormData({
                id: updatedKojiResponse.data.id || "",
                company_name: updatedKojiResponse.data.company_name || "",
                location_name: updatedKojiResponse.data.location_name || "",
                description: updatedKojiResponse.data.description || "",
                tags: Array.isArray(updatedKojiResponse.data.tags) ? updatedKojiResponse.data.tags.join(", ") : updatedKojiResponse.data.tags || "",
                start_date: startDate,
                end_date: endDate
              });
            }
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
            currentKoji?.managed_files && currentKoji.managed_files.length > 0 && /* @__PURE__ */ jsxs(Box, { sx: { mb: 3 }, children: [
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
                  children: currentKoji.managed_files.map((file, index) => {
                    const needsRename = file.current && file.recommended && file.current !== file.recommended;
                    return /* @__PURE__ */ jsxs(
                      Paper,
                      {
                        elevation: 0,
                        sx: {
                          mb: index < currentKoji.managed_files.length - 1 ? 2 : 0,
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
            currentKoji && (!currentKoji.managed_files || currentKoji.managed_files.length === 0) && /* @__PURE__ */ jsxs(Box, { sx: { mb: 3 }, children: [
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
const Kojies = () => {
  const [kojies, setKojies] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [selectedKoji, setSelectedKoji] = useState(
    null
  );
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [showHelp, setShowHelp] = useState(false);
  const [shouldReloadOnClose, setShouldReloadOnClose] = useState(false);
  const { setKojiCount } = useKoji();
  const loadKojies = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await getBusinessKojies();
      if (response.data) {
        setKojies(response.data);
        setKojiCount(response.data.length);
      } else {
        setKojies([]);
        setKojiCount(0);
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
    loadKojies();
  }, []);
  const handleKojiClick = (koji) => {
    setSelectedKoji(koji);
    setIsEditModalOpen(true);
  };
  const updateKoji = async (updatedKoji) => {
    try {
      const response = await putBusinessKojies({
        body: updatedKoji
      });
      if (response.data) {
        setKojies(
          (prevKojies) => prevKojies.map((k) => k.path === response.data.path ? response.data : k)
        );
        return response.data;
      } else {
        throw new Error("更新に失敗しました");
      }
    } catch (err) {
      console.error("Error updating koji:", err);
      throw err;
    }
  };
  const closeEditModal = () => {
    setIsEditModalOpen(false);
    setSelectedKoji(null);
  };
  const needsFileRename = (koji) => {
    if (!koji.managed_files || koji.managed_files.length === 0) {
      return false;
    }
    const needsRename = koji.managed_files.some((file) => {
      return file.current && file.recommended && file.current !== file.recommended;
    });
    return needsRename;
  };
  const handleKojiUpdate = (updatedKoji) => {
    setSelectedKoji(updatedKoji);
    if (selectedKoji && selectedKoji.path !== updatedKoji.path) {
      setShouldReloadOnClose(true);
    }
    setKojies((prevKojies) => {
      const existingIndex = prevKojies.findIndex((k) => k.path === updatedKoji.path);
      if (existingIndex !== -1) {
        const updatedKojies = [...prevKojies];
        updatedKojies[existingIndex] = updatedKoji;
        return updatedKojies;
      } else {
        const oldKojiIndex = prevKojies.findIndex(
          (k) => selectedKoji && k.path === selectedKoji.path
        );
        if (oldKojiIndex !== -1) {
          const updatedKojies = [...prevKojies];
          updatedKojies.splice(oldKojiIndex, 1);
          updatedKojies.push(updatedKoji);
          return updatedKojies.sort((a, b) => {
            const dateA = a.start_date ? new Date(typeof a.start_date === "string" ? a.start_date : a.start_date["time.Time"]).getTime() : 0;
            const dateB = b.start_date ? new Date(typeof b.start_date === "string" ? b.start_date : b.start_date["time.Time"]).getTime() : 0;
            if (dateA > 0 && dateB === 0) return -1;
            if (dateA === 0 && dateB > 0) return 1;
            if (dateA > 0 && dateB > 0) return dateB - dateA;
            return (b.name || "").localeCompare(a.name || "");
          });
        } else {
          return [...prevKojies, updatedKoji];
        }
      }
    });
  };
  if (loading) {
    return /* @__PURE__ */ jsx("div", { className: "business-entity-loading", children: /* @__PURE__ */ jsx("div", { children: "工事データを読み込み中..." }) });
  }
  if (error) {
    return /* @__PURE__ */ jsxs("div", { className: "business-entity-error", children: [
      /* @__PURE__ */ jsx("div", { className: "business-entity-error-message", children: error }),
      /* @__PURE__ */ jsx(
        "button",
        {
          onClick: loadKojies,
          className: "business-entity-retry-button",
          children: "再試行"
        }
      )
    ] });
  }
  return /* @__PURE__ */ jsxs("div", { className: "business-entity-container", children: [
    /* @__PURE__ */ jsxs("div", { className: "business-entity-controls", children: [
      /* @__PURE__ */ jsx(
        Link,
        {
          to: "/kojies/gantt",
          className: "business-entity-gantt-button",
          children: "📊 工程表を表示"
        }
      ),
      /* @__PURE__ */ jsxs("div", { className: "business-entity-count", children: [
        "全",
        kojies.length,
        "件"
      ] })
    ] }),
    showHelp && /* @__PURE__ */ jsxs("div", { className: "business-entity-help-box", children: [
      /* @__PURE__ */ jsx(
        "button",
        {
          onClick: () => setShowHelp(false),
          className: "business-entity-help-close",
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
      /* @__PURE__ */ jsx("p", { children: "✅ 開始日・終了日・説明・タグ・会社名・現場名を編雈可能" }),
      /* @__PURE__ */ jsx("p", { children: "💾 編集後は自動で保存されます" }),
      /* @__PURE__ */ jsx("h3", { style: { marginTop: "15px" }, children: "開発状況" }),
      /* @__PURE__ */ jsx("p", { children: "✅ 工事データの取得" }),
      /* @__PURE__ */ jsx("p", { children: "✅ 編集モーダル機能" }),
      /* @__PURE__ */ jsx("p", { children: "🔄 工程表機能（次のステップ）" })
    ] }),
    /* @__PURE__ */ jsxs("div", { className: "business-entity-list-container", children: [
      /* @__PURE__ */ jsxs("div", { className: "business-entity-list-header", children: [
        /* @__PURE__ */ jsx("div", { className: "business-entity-header-date", children: "開始日" }),
        /* @__PURE__ */ jsx("div", { className: "business-entity-header-company", children: "会社名" }),
        /* @__PURE__ */ jsx("div", { className: "business-entity-header-location", children: "現場名" }),
        /* @__PURE__ */ jsx("div", { className: "business-entity-header-spacer" }),
        /* @__PURE__ */ jsx("div", { className: "business-entity-header-date", style: { marginRight: "24px" }, children: "終了日" }),
        /* @__PURE__ */ jsx("div", { className: "business-entity-header-status", children: "ステータス" }),
        /* @__PURE__ */ jsx(
          "button",
          {
            onClick: () => setShowHelp(!showHelp),
            className: "business-entity-help-button",
            title: "使用方法を表示",
            children: "?"
          }
        )
      ] }),
      /* @__PURE__ */ jsx("div", { className: "business-entity-scroll-area", children: kojies.length === 0 ? /* @__PURE__ */ jsx("div", { className: "business-entity-empty", children: "工事データが見つかりません" }) : /* @__PURE__ */ jsx("div", { children: kojies.map((koji, index) => /* @__PURE__ */ jsxs(
        "div",
        {
          className: "business-entity-item-row",
          onClick: () => handleKojiClick(koji),
          title: "クリックして編集",
          children: [
            /* @__PURE__ */ jsxs("div", { className: "business-entity-item-info", children: [
              /* @__PURE__ */ jsx("div", { className: "business-entity-item-info-date", children: koji.start_date ? new Date(
                typeof koji.start_date === "string" ? koji.start_date : koji.start_date["time.Time"]
              ).toLocaleDateString("ja-JP") : "未設定" }),
              /* @__PURE__ */ jsx("div", { className: "business-entity-item-info-company", children: koji.company_name || "会社名未設定" }),
              /* @__PURE__ */ jsx("div", { className: "business-entity-item-info-location", children: koji.location_name || "現場名未設定" }),
              /* @__PURE__ */ jsx("div", { className: "business-entity-item-info-description", children: koji.description || "" }),
              /* @__PURE__ */ jsxs("div", { className: "business-entity-item-info-date end-date", style: {
                marginRight: "24px"
              }, children: [
                "～",
                koji.end_date ? new Date(
                  typeof koji.end_date === "string" ? koji.end_date : koji.end_date["time.Time"]
                ).toLocaleDateString("ja-JP") : "未設定"
              ] })
            ] }),
            /* @__PURE__ */ jsxs("div", { style: { display: "flex", alignItems: "center", gap: "8px" }, children: [
              needsFileRename(koji) && /* @__PURE__ */ jsx(
                "span",
                {
                  className: "business-entity-item-rename-indicator",
                  title: "管理ファイルの名前変更が必要です",
                  children: "⚠️"
                }
              ),
              /* @__PURE__ */ jsx(
                "div",
                {
                  className: `business-entity-item-status ${koji.status === "進行中" ? "business-entity-item-status-ongoing" : koji.status === "完了" ? "business-entity-item-status-completed" : koji.status === "予定" ? "business-entity-item-status-planned" : ""}`,
                  children: koji.status || "未設定"
                }
              )
            ] })
          ]
        },
        koji.id || index
      )) }) })
    ] }),
    /* @__PURE__ */ jsx(
      KojiDetailModal,
      {
        isOpen: isEditModalOpen,
        onClose: () => {
          closeEditModal();
          if (shouldReloadOnClose) {
            loadKojies();
            setShouldReloadOnClose(false);
          }
        },
        koji: selectedKoji,
        onUpdate: updateKoji,
        onKojiUpdate: handleKojiUpdate
      }
    )
  ] });
};
const _layout_kojies = UNSAFE_withComponentProps(function KojiesPage() {
  return /* @__PURE__ */ jsx(Kojies, {});
});
const route4 = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
  __proto__: null,
  default: _layout_kojies
}, Symbol.toStringTag, { value: "Module" }));
const KojiGanttChart = () => {
  const [kojies, setKojies] = useState([]);
  const [ganttItems, setGanttItems] = useState([]);
  const [selectedKoji, setSelectedKoji] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [viewStartDate, setViewStartDate] = useState(/* @__PURE__ */ new Date());
  const [viewEndDate, setViewEndDate] = useState(/* @__PURE__ */ new Date());
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [visibleKojies, setVisibleKojies] = useState([]);
  const [shouldReloadOnClose, setShouldReloadOnClose] = useState(false);
  const scrollContainerRef = useRef(null);
  const [itemsPerPage, setItemsPerPage] = useState(10);
  const MIN_ITEMS = 5;
  const DAY_WIDTH = 10;
  const ROW_HEIGHT = 40;
  const loadKojies = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await getBusinessKojies();
      const kojies2 = response.data || [];
      setKojies(kojies2);
    } catch (err) {
      console.error("Error loading kojies:", err);
      setError(`工事データの読み込みに失敗しました: ${err instanceof Error ? err.message : "Unknown error"}`);
    } finally {
      setLoading(false);
    }
  };
  useEffect(() => {
    loadKojies();
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
    if (kojies.length === 0) return;
    let minDate = /* @__PURE__ */ new Date();
    let maxDate = /* @__PURE__ */ new Date();
    let hasValidDate = false;
    kojies.forEach((koji) => {
      try {
        const startDate = koji.start_date ? new Date(koji.start_date) : null;
        const endDate = koji.end_date ? new Date(koji.end_date) : null;
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
  }, [kojies]);
  const updateVisibleKojies = (scrollLeft = 0) => {
    if (kojies.length === 0 || !scrollContainerRef.current) return;
    const containerWidth = scrollContainerRef.current.clientWidth;
    const visibleStartDays = scrollLeft / DAY_WIDTH;
    const visibleEndDays = (scrollLeft + containerWidth) / DAY_WIDTH;
    const visibleStartDate = new Date(viewStartDate.getTime() + visibleStartDays * 24 * 60 * 60 * 1e3);
    const visibleEndDate = new Date(viewStartDate.getTime() + visibleEndDays * 24 * 60 * 60 * 1e3);
    const relevantKojies = kojies.filter((koji) => {
      try {
        const kojiStart = koji.start_date ? new Date(koji.start_date) : /* @__PURE__ */ new Date();
        const kojiEnd = koji.end_date ? new Date(koji.end_date) : new Date(kojiStart.getTime() + 90 * 24 * 60 * 60 * 1e3);
        return kojiStart <= visibleEndDate && kojiEnd >= visibleStartDate;
      } catch {
        return false;
      }
    });
    const sortedRelevantKojies = relevantKojies.sort((a, b) => {
      const dateA = a.start_date ? new Date(a.start_date).getTime() : 0;
      const dateB = b.start_date ? new Date(b.start_date).getTime() : 0;
      return dateA - dateB;
    });
    let baselineDate;
    if (sortedRelevantKojies.length > 0) {
      baselineDate = sortedRelevantKojies[0].start_date ? new Date(sortedRelevantKojies[0].start_date).getTime() : 0;
    } else {
      const visibleStartTime = visibleStartDate.getTime();
      const kojiesBeforeVisible = kojies.filter((koji) => {
        const kojiStartDate = koji.start_date ? new Date(koji.start_date).getTime() : 0;
        return kojiStartDate <= visibleStartTime;
      });
      if (kojiesBeforeVisible.length > 0) {
        const closestKoji = kojiesBeforeVisible.sort((a, b) => {
          const dateA = a.start_date ? new Date(a.start_date).getTime() : 0;
          const dateB = b.start_date ? new Date(b.start_date).getTime() : 0;
          return dateB - dateA;
        })[0];
        baselineDate = closestKoji.start_date ? new Date(closestKoji.start_date).getTime() : 0;
      } else {
        const allKojiesSorted2 = [...kojies].sort((a, b) => {
          const dateA = a.start_date ? new Date(a.start_date).getTime() : 0;
          const dateB = b.start_date ? new Date(b.start_date).getTime() : 0;
          return dateA - dateB;
        });
        if (allKojiesSorted2.length === 0) return;
        baselineDate = allKojiesSorted2[0].start_date ? new Date(allKojiesSorted2[0].start_date).getTime() : 0;
      }
    }
    const allKojiesSorted = [...kojies].sort((a, b) => {
      const dateA = a.start_date ? new Date(a.start_date).getTime() : 0;
      const dateB = b.start_date ? new Date(b.start_date).getTime() : 0;
      return dateA - dateB;
    });
    const kojiesFromBaselineDate = allKojiesSorted.filter((koji) => {
      const kojiStartDate = koji.start_date ? new Date(koji.start_date).getTime() : 0;
      return kojiStartDate >= baselineDate;
    });
    let finalKojies = kojiesFromBaselineDate.slice(0, itemsPerPage);
    if (finalKojies.length < itemsPerPage) {
      const allKojiesDescending = [...kojies].sort((a, b) => {
        const dateA = a.start_date ? new Date(a.start_date).getTime() : 0;
        const dateB = b.start_date ? new Date(b.start_date).getTime() : 0;
        return dateB - dateA;
      });
      finalKojies = allKojiesDescending.slice(0, itemsPerPage).sort((a, b) => {
        const dateA = a.start_date ? new Date(a.start_date).getTime() : 0;
        const dateB = b.start_date ? new Date(b.start_date).getTime() : 0;
        return dateA - dateB;
      });
    }
    setVisibleKojies(finalKojies);
  };
  const handleScroll = () => {
    if (scrollContainerRef.current) {
      updateVisibleKojies(scrollContainerRef.current.scrollLeft);
    }
  };
  useEffect(() => {
    if (kojies.length > 0 && scrollContainerRef.current) {
      updateVisibleKojies(scrollContainerRef.current.scrollLeft);
    }
  }, [kojies, viewStartDate, itemsPerPage]);
  const calculateItemsPerPage = () => {
    if (!scrollContainerRef.current) return MIN_ITEMS;
    const containerHeight = scrollContainerRef.current.clientHeight;
    const availableHeight = containerHeight - 55;
    const maxItems = Math.floor(availableHeight / ROW_HEIGHT);
    return Math.max(MIN_ITEMS, maxItems);
  };
  useEffect(() => {
    const handleResize = () => {
      if (scrollContainerRef.current) {
        const newItemsPerPage = calculateItemsPerPage();
        setItemsPerPage(newItemsPerPage);
        updateVisibleKojies(scrollContainerRef.current.scrollLeft);
      }
    };
    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, [kojies, viewStartDate]);
  useEffect(() => {
    if (scrollContainerRef.current && ganttItems.length > 0) {
      const newItemsPerPage = calculateItemsPerPage();
      setItemsPerPage(newItemsPerPage);
    }
  }, [ganttItems.length]);
  useEffect(() => {
    if (visibleKojies.length === 0) return;
    const items = visibleKojies.map((koji, index) => {
      let startDate;
      let endDate;
      try {
        startDate = koji.start_date ? new Date(koji.start_date) : /* @__PURE__ */ new Date();
        if (isNaN(startDate.getTime())) {
          startDate = /* @__PURE__ */ new Date();
        }
      } catch {
        startDate = /* @__PURE__ */ new Date();
      }
      try {
        endDate = koji.end_date ? new Date(koji.end_date) : new Date(startDate.getTime() + 90 * 24 * 60 * 60 * 1e3);
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
        ...koji,
        startX: isNaN(startX) ? 0 : startX,
        width: isNaN(width) ? DAY_WIDTH : width,
        row: index
      };
    });
    setGanttItems(items);
  }, [visibleKojies, viewStartDate]);
  const handleKojiEdit = (koji) => {
    setSelectedKoji(koji);
    setIsEditModalOpen(true);
  };
  const handleKojiNameClick = (koji) => {
    if (!scrollContainerRef.current) return;
    try {
      const kojiStart = koji.start_date ? new Date(koji.start_date) : /* @__PURE__ */ new Date();
      const kojiEnd = koji.end_date ? new Date(koji.end_date) : new Date(kojiStart.getTime() + 90 * 24 * 60 * 60 * 1e3);
      const kojiMiddle = new Date((kojiStart.getTime() + kojiEnd.getTime()) / 2);
      const middleX = (kojiMiddle.getTime() - viewStartDate.getTime()) / (1e3 * 60 * 60 * 24) * DAY_WIDTH;
      const containerWidth = scrollContainerRef.current.clientWidth;
      const scrollPosition = Math.max(0, middleX - containerWidth / 2);
      scrollContainerRef.current.scrollLeft = scrollPosition;
    } catch (error2) {
      console.error("Error calculating koji center:", error2);
    }
  };
  const scrollToToday = () => {
    if (!scrollContainerRef.current) return;
    const todayX = ((/* @__PURE__ */ new Date()).getTime() - viewStartDate.getTime()) / (1e3 * 60 * 60 * 24) * DAY_WIDTH;
    const containerWidth = scrollContainerRef.current.clientWidth;
    const scrollPosition = Math.max(0, todayX - containerWidth / 2);
    scrollContainerRef.current.scrollLeft = scrollPosition;
  };
  const updateKoji = async (updatedKoji) => {
    try {
      const response = await putBusinessKojies({
        body: updatedKoji
      });
      if (response.data) {
        setKojies(
          (prevKojies) => prevKojies.map((k) => k.path === response.data.path ? response.data : k)
        );
        return response.data;
      } else {
        throw new Error("更新に失敗しました");
      }
    } catch (err) {
      console.error("Error updating koji:", err);
      throw err;
    }
  };
  const handleKojiUpdate = (updatedKoji) => {
    setSelectedKoji(updatedKoji);
    if (selectedKoji && selectedKoji.path !== updatedKoji.path) {
      setShouldReloadOnClose(true);
    }
    setKojies((prevKojies) => {
      const existingIndex = prevKojies.findIndex((k) => k.path === updatedKoji.path);
      if (existingIndex !== -1) {
        const updatedKojies = [...prevKojies];
        updatedKojies[existingIndex] = updatedKoji;
        return updatedKojies;
      } else {
        const oldKojiIndex = prevKojies.findIndex(
          (k) => selectedKoji && k.path === selectedKoji.path
        );
        if (oldKojiIndex !== -1) {
          const updatedKojies = [...prevKojies];
          updatedKojies.splice(oldKojiIndex, 1);
          updatedKojies.push(updatedKoji);
          return updatedKojies.sort((a, b) => {
            const dateA = a.start_date ? new Date(typeof a.start_date === "string" ? a.start_date : a.start_date["time.Time"]).getTime() : 0;
            const dateB = b.start_date ? new Date(typeof b.start_date === "string" ? b.start_date : b.start_date["time.Time"]).getTime() : 0;
            if (dateA > 0 && dateB === 0) return -1;
            if (dateA === 0 && dateB > 0) return 1;
            if (dateA > 0 && dateB > 0) return dateA - dateB;
            return (a.name || "").localeCompare(b.name || "");
          });
        } else {
          return [...prevKojies, updatedKoji];
        }
      }
    });
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
  const needsFileRename = (koji) => {
    if (!koji.managed_files || koji.managed_files.length === 0) {
      return false;
    }
    const needsRename = koji.managed_files.some((file) => {
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
      if ((day === 1 || day % 3 === 1) && day !== 31) {
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
    /* @__PURE__ */ jsxs("div", { style: {
      marginBottom: "20px",
      display: "flex",
      justifyContent: "space-between",
      alignItems: "center"
    }, children: [
      /* @__PURE__ */ jsx(
        "button",
        {
          onClick: scrollToToday,
          className: "gantt-today-button",
          children: "📅 今日へ移動"
        }
      ),
      /* @__PURE__ */ jsxs("div", { style: {
        fontSize: "16px",
        color: "#666",
        fontWeight: "500"
      }, children: [
        "表示中: ",
        ganttItems.length,
        "件 / 全",
        kojies.length,
        "件"
      ] })
    ] }),
    /* @__PURE__ */ jsxs("div", { className: "gantt-wrapper", children: [
      /* @__PURE__ */ jsxs("div", { className: "gantt-sidebar", children: [
        /* @__PURE__ */ jsx("div", { className: "gantt-header-left", children: "会社名" }),
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
          const isActiveKoji = today >= startDate && today <= endDate;
          return /* @__PURE__ */ jsx(
            "div",
            {
              className: `gantt-row-label ${isActiveKoji ? "gantt-row-label-active" : ""}`,
              style: {
                height: ROW_HEIGHT
              },
              children: /* @__PURE__ */ jsx(
                "div",
                {
                  className: `koji-name koji-name-clickable ${isActiveKoji ? "koji-name-active" : ""}`,
                  onClick: () => handleKojiNameClick(item),
                  title: "クリックして工事期間の中央に移動",
                  children: item.company_name
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
            className: "month-header month-header-content",
            style: {
              left: header.startX,
              width: header.width
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
        /* @__PURE__ */ jsx("div", { className: "gantt-header day-header-row", style: { height: "25px", borderBottom: "1px solid #333" }, children: dayHeaders.map((header, index) => /* @__PURE__ */ jsx(
          "div",
          {
            className: "day-header day-header-content",
            style: {
              left: header.startX,
              width: header.width
            },
            children: header.date
          },
          index
        )) }),
        /* @__PURE__ */ jsx("div", { className: "gantt-month-boundaries", style: { top: 0, height: "100%" }, children: monthBoundaries.map((boundary, index) => /* @__PURE__ */ jsx(
          "div",
          {
            className: "month-boundary-line",
            style: {
              left: boundary.startX,
              top: 0,
              height: Math.max(400, itemsPerPage * ROW_HEIGHT + 55)
              // ヘッダー分も含む
            },
            title: `${boundary.year}年${boundary.month}月開始`
          },
          `month-boundary-${index}`
        )) }),
        /* @__PURE__ */ jsxs("div", { className: "gantt-body", children: [
          /* @__PURE__ */ jsx("div", { className: "gantt-grid", children: dayHeaders.map((header, index) => /* @__PURE__ */ jsx(
            "div",
            {
              className: "grid-line",
              style: {
                left: header.startX,
                height: Math.max(400, itemsPerPage * ROW_HEIGHT)
              }
            },
            index
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
                top: index * ROW_HEIGHT + 10,
                height: ROW_HEIGHT - 15,
                backgroundColor: getStatusColor(item.status)
              },
              onClick: () => handleKojiEdit(item),
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
                height: Math.max(400, itemsPerPage * ROW_HEIGHT),
                // 動的な高さ
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
      KojiDetailModal,
      {
        isOpen: isEditModalOpen,
        onClose: () => {
          setIsEditModalOpen(false);
          if (shouldReloadOnClose) {
            loadKojies();
            setShouldReloadOnClose(false);
          }
        },
        koji: selectedKoji,
        onUpdate: updateKoji,
        onKojiUpdate: handleKojiUpdate
      }
    )
  ] });
};
const _layout_gantt = UNSAFE_withComponentProps(function GanttChartPage() {
  return /* @__PURE__ */ jsx(KojiGanttChart, {});
});
const route5 = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
  __proto__: null,
  default: _layout_gantt
}, Symbol.toStringTag, { value: "Module" }));
const CompanyDetailModal = ({
  isOpen,
  onClose,
  company,
  onCompanyUpdate
}) => {
  const [formData, setFormData] = useState({
    id: "",
    name: "",
    short_name: "",
    business_type: "",
    phone: "",
    email: "",
    website: "",
    address: "",
    description: "",
    tags: ""
  });
  const [currentCompany, setCurrentCompany] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);
  const [isEditing, setIsEditing] = useState(false);
  const [hasChanges, setHasChanges] = useState(false);
  useEffect(() => {
    if (company) {
      setCurrentCompany(company);
      const newFormData = {
        id: company.id || "",
        name: company.name || "",
        short_name: company.short_name || "",
        business_type: company.business_type || "",
        phone: company.phone || "",
        email: company.email || "",
        website: company.website || "",
        address: company.address || "",
        tags: Array.isArray(company.tags) ? company.tags.join(", ") : company.tags || ""
      };
      setFormData(newFormData);
      setError(null);
      setIsEditing(false);
      setHasChanges(false);
    }
  }, [company]);
  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value
    }));
    setHasChanges(true);
  };
  const handleEditToggle = () => {
    if (isEditing && hasChanges) {
      if (confirm("変更を破棄しますか？")) {
        if (company) {
          setFormData({
            id: company.id || "",
            name: company.name || "",
            short_name: company.short_name || "",
            business_type: company.business_type || "",
            phone: company.phone || "",
            email: company.email || "",
            website: company.website || "",
            address: company.address || "",
            tags: Array.isArray(company.tags) ? company.tags.join(", ") : company.tags || ""
          });
        }
        setIsEditing(false);
        setHasChanges(false);
      }
    } else {
      setIsEditing(!isEditing);
      setHasChanges(false);
    }
  };
  const handleUpdate = async () => {
    if (!company) {
      setError("会社データがありません");
      return;
    }
    setIsLoading(true);
    setError(null);
    try {
      const tagsArray = formData.tags ? formData.tags.split(",").map((tag) => tag.trim()).filter((tag) => tag.length > 0) : [];
      const updatedCompany = {
        ...company,
        name: formData.name,
        short_name: formData.short_name,
        business_type: formData.business_type,
        phone: formData.phone,
        email: formData.email,
        website: formData.website,
        address: formData.address,
        tags: tagsArray
      };
      const response = await putBusinessCompanies({
        body: updatedCompany
      });
      if (response.data) {
        setCurrentCompany(response.data);
        setIsEditing(false);
        setHasChanges(false);
        if (onCompanyUpdate) {
          onCompanyUpdate(response.data);
        }
        setError(null);
      } else {
        throw new Error("更新レスポンスが無効です");
      }
    } catch (err) {
      console.error("Error updating company:", err);
      setError(err instanceof Error ? err.message : "会社情報の更新に失敗しました");
    } finally {
      setIsLoading(false);
    }
  };
  const displayTags = currentCompany?.tags ? Array.isArray(currentCompany.tags) ? currentCompany.tags : [currentCompany.tags] : [];
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
              /* @__PURE__ */ jsxs(Box, { sx: { display: "flex", alignItems: "center", gap: 1 }, children: [
                /* @__PURE__ */ jsx(Typography, { variant: "h6", component: "h2", children: "会社詳細情報" }),
                isEditing && /* @__PURE__ */ jsx(
                  Chip,
                  {
                    label: "編集中",
                    color: "warning",
                    size: "small",
                    icon: /* @__PURE__ */ jsx(EditIcon, {})
                  }
                )
              ] }),
              /* @__PURE__ */ jsxs(Box, { sx: { display: "flex", gap: 1 }, children: [
                /* @__PURE__ */ jsx(
                  Button,
                  {
                    variant: isEditing ? "outlined" : "contained",
                    size: "small",
                    onClick: handleEditToggle,
                    startIcon: isEditing ? /* @__PURE__ */ jsx(CancelIcon, {}) : /* @__PURE__ */ jsx(EditIcon, {}),
                    color: isEditing ? "secondary" : "primary",
                    children: isEditing ? "キャンセル" : "編集"
                  }
                ),
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
              ] })
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
                    /* @__PURE__ */ jsxs(
                      Typography,
                      {
                        variant: "subtitle2",
                        sx: {
                          color: "primary.main",
                          fontWeight: 600,
                          mb: 1,
                          display: "flex",
                          alignItems: "center",
                          gap: 1
                        },
                        children: [
                          /* @__PURE__ */ jsx(BusinessIcon, { fontSize: "small" }),
                          "会社名"
                        ]
                      }
                    ),
                    /* @__PURE__ */ jsx(
                      TextField,
                      {
                        fullWidth: true,
                        size: "small",
                        name: "name",
                        value: formData.name,
                        onChange: handleInputChange,
                        disabled: isLoading || !isEditing,
                        sx: {
                          "& .MuiOutlinedInput-root": {
                            backgroundColor: isEditing ? "white" : "grey.50"
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
                        children: "略称"
                      }
                    ),
                    /* @__PURE__ */ jsx(
                      TextField,
                      {
                        fullWidth: true,
                        size: "small",
                        name: "short_name",
                        value: formData.short_name,
                        onChange: handleInputChange,
                        disabled: isLoading || !isEditing,
                        sx: {
                          "& .MuiOutlinedInput-root": {
                            backgroundColor: isEditing ? "white" : "grey.50"
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
                        children: "業種"
                      }
                    ),
                    /* @__PURE__ */ jsx(
                      TextField,
                      {
                        fullWidth: true,
                        size: "small",
                        name: "business_type",
                        value: formData.business_type,
                        onChange: handleInputChange,
                        disabled: isLoading || !isEditing,
                        sx: {
                          "& .MuiOutlinedInput-root": {
                            backgroundColor: isEditing ? "white" : "grey.50"
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
                        children: "ID"
                      }
                    ),
                    /* @__PURE__ */ jsx(
                      Typography,
                      {
                        variant: "body2",
                        sx: {
                          fontFamily: "monospace",
                          padding: "8px 12px",
                          backgroundColor: "white",
                          border: "1px solid",
                          borderColor: "grey.300",
                          borderRadius: 1,
                          color: "grey.700"
                        },
                        children: formData.id || "ID未設定"
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
                    backgroundColor: "grey.50",
                    border: "1px solid",
                    borderColor: "grey.200",
                    borderRadius: 2
                  },
                  children: [
                    /* @__PURE__ */ jsxs(
                      Typography,
                      {
                        variant: "subtitle2",
                        sx: {
                          color: "text.primary",
                          fontWeight: 500,
                          mb: 1,
                          display: "flex",
                          alignItems: "center",
                          gap: 1
                        },
                        children: [
                          /* @__PURE__ */ jsx(PhoneIcon, { fontSize: "small" }),
                          "電話番号"
                        ]
                      }
                    ),
                    /* @__PURE__ */ jsx(
                      TextField,
                      {
                        fullWidth: true,
                        size: "small",
                        name: "phone",
                        value: formData.phone,
                        onChange: handleInputChange,
                        disabled: isLoading || !isEditing,
                        sx: {
                          "& .MuiOutlinedInput-root": {
                            backgroundColor: isEditing ? "white" : "grey.50"
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
                    backgroundColor: "grey.50",
                    border: "1px solid",
                    borderColor: "grey.200",
                    borderRadius: 2
                  },
                  children: [
                    /* @__PURE__ */ jsxs(
                      Typography,
                      {
                        variant: "subtitle2",
                        sx: {
                          color: "text.primary",
                          fontWeight: 500,
                          mb: 1,
                          display: "flex",
                          alignItems: "center",
                          gap: 1
                        },
                        children: [
                          /* @__PURE__ */ jsx(EmailIcon, { fontSize: "small" }),
                          "メールアドレス"
                        ]
                      }
                    ),
                    /* @__PURE__ */ jsx(
                      TextField,
                      {
                        fullWidth: true,
                        size: "small",
                        name: "email",
                        type: "email",
                        value: formData.email,
                        onChange: handleInputChange,
                        disabled: isLoading || !isEditing,
                        sx: {
                          "& .MuiOutlinedInput-root": {
                            backgroundColor: isEditing ? "white" : "grey.50"
                          }
                        }
                      }
                    )
                  ]
                }
              ) })
            ] }),
            /* @__PURE__ */ jsx(Box, { sx: { display: "flex", gap: 2, mb: 3, flexDirection: { xs: "column", sm: "row" } }, children: /* @__PURE__ */ jsx(Box, { sx: { flex: 1 }, children: /* @__PURE__ */ jsxs(
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
                  /* @__PURE__ */ jsxs(
                    Typography,
                    {
                      variant: "subtitle2",
                      sx: {
                        color: "text.primary",
                        fontWeight: 500,
                        mb: 1,
                        display: "flex",
                        alignItems: "center",
                        gap: 1
                      },
                      children: [
                        /* @__PURE__ */ jsx(LanguageIcon, { fontSize: "small" }),
                        "ウェブサイト"
                      ]
                    }
                  ),
                  /* @__PURE__ */ jsx(
                    TextField,
                    {
                      fullWidth: true,
                      size: "small",
                      name: "website",
                      type: "url",
                      value: formData.website,
                      onChange: handleInputChange,
                      disabled: isLoading || !isEditing,
                      sx: {
                        "& .MuiOutlinedInput-root": {
                          backgroundColor: isEditing ? "white" : "grey.50"
                        }
                      }
                    }
                  )
                ]
              }
            ) }) }),
            /* @__PURE__ */ jsx(Box, { sx: { mb: 3 }, children: /* @__PURE__ */ jsxs(
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
                      children: "住所"
                    }
                  ),
                  /* @__PURE__ */ jsx(
                    TextField,
                    {
                      fullWidth: true,
                      size: "small",
                      name: "address",
                      value: formData.address,
                      onChange: handleInputChange,
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
                    disabled: isLoading || !isEditing,
                    placeholder: "会社の詳細説明を入力してください",
                    sx: {
                      "& .MuiOutlinedInput-root": {
                        backgroundColor: isEditing ? "white" : "grey.50"
                      }
                    }
                  }
                )
              ] }),
              /* @__PURE__ */ jsxs(Box, { children: [
                /* @__PURE__ */ jsxs(
                  Typography,
                  {
                    variant: "subtitle2",
                    sx: {
                      color: "text.primary",
                      fontWeight: 500,
                      mb: 1,
                      display: "flex",
                      alignItems: "center",
                      gap: 1
                    },
                    children: [
                      /* @__PURE__ */ jsx(TagIcon, { fontSize: "small" }),
                      "タグ"
                    ]
                  }
                ),
                /* @__PURE__ */ jsx(
                  TextField,
                  {
                    fullWidth: true,
                    name: "tags",
                    value: formData.tags,
                    onChange: handleInputChange,
                    placeholder: "カンマ区切りで入力",
                    disabled: isLoading || !isEditing,
                    sx: {
                      "& .MuiOutlinedInput-root": {
                        backgroundColor: isEditing ? "white" : "grey.50"
                      }
                    }
                  }
                ),
                displayTags.length > 0 && /* @__PURE__ */ jsxs(Box, { sx: { mt: 2 }, children: [
                  /* @__PURE__ */ jsx(Typography, { variant: "body2", color: "text.secondary", sx: { mb: 1 }, children: "現在のタグ:" }),
                  /* @__PURE__ */ jsx(Box, { sx: { display: "flex", flexWrap: "wrap", gap: 1 }, children: displayTags.map((tag, index) => /* @__PURE__ */ jsx(
                    Chip,
                    {
                      label: tag,
                      size: "small",
                      color: "primary",
                      variant: "outlined"
                    },
                    index
                  )) })
                ] })
              ] })
            ] }),
            isEditing && /* @__PURE__ */ jsxs(Box, { sx: { display: "flex", justifyContent: "flex-end", gap: 2, mt: 3 }, children: [
              /* @__PURE__ */ jsx(
                Button,
                {
                  variant: "outlined",
                  onClick: handleEditToggle,
                  disabled: isLoading,
                  startIcon: /* @__PURE__ */ jsx(CancelIcon, {}),
                  children: "キャンセル"
                }
              ),
              /* @__PURE__ */ jsx(
                Button,
                {
                  variant: "contained",
                  onClick: handleUpdate,
                  disabled: isLoading || !hasChanges,
                  startIcon: isLoading ? /* @__PURE__ */ jsx(CircularProgress, { size: 16 }) : /* @__PURE__ */ jsx(SaveIcon, {}),
                  color: "primary",
                  children: isLoading ? "保存中..." : "保存"
                }
              )
            ] }),
            !isEditing && hasChanges === false && currentCompany && /* @__PURE__ */ jsx(Alert, { severity: "success", sx: { mt: 2 }, children: "会社情報が正常に更新されました。" })
          ] })
        ] })
      ]
    }
  );
};
const Companies = () => {
  const [companies, setCompanies] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showHelp, setShowHelp] = useState(false);
  const [selectedCompany, setSelectedCompany] = useState(null);
  const [isDetailModalOpen, setIsDetailModalOpen] = useState(false);
  const loadCompanies = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await getBusinessCompanies();
      if (response.data) {
        setCompanies(response.data);
      } else {
        setCompanies([]);
      }
    } catch (err) {
      console.error("Error loading company entries:", err);
      setError(
        `会社データの読み込みに失敗しました: ${err instanceof Error ? err.message : String(err)}`
      );
    } finally {
      setLoading(false);
    }
  };
  useEffect(() => {
    loadCompanies();
  }, []);
  const handleCompanyClick = (company) => {
    setSelectedCompany(company);
    setIsDetailModalOpen(true);
  };
  const closeDetailModal = () => {
    setIsDetailModalOpen(false);
    setSelectedCompany(null);
  };
  const handleCompanyUpdate = (updatedCompany) => {
    setSelectedCompany(updatedCompany);
    setCompanies(
      (prevCompanies) => prevCompanies.map((c) => c.name === updatedCompany.name ? updatedCompany : c)
    );
  };
  if (loading) {
    return /* @__PURE__ */ jsx("div", { className: "business-entity-loading", children: /* @__PURE__ */ jsx("div", { children: "会社データを読み込み中..." }) });
  }
  if (error) {
    return /* @__PURE__ */ jsxs("div", { className: "business-entity-error", children: [
      /* @__PURE__ */ jsx("div", { className: "business-entity-error-message", children: error }),
      /* @__PURE__ */ jsx(
        "button",
        {
          onClick: loadCompanies,
          className: "business-entity-retry-button",
          children: "再試行"
        }
      )
    ] });
  }
  return /* @__PURE__ */ jsxs("div", { className: "business-entity-container", children: [
    /* @__PURE__ */ jsx("div", { className: "business-entity-controls", children: /* @__PURE__ */ jsxs("div", { className: "business-entity-count", children: [
      "全",
      companies.length,
      "件"
    ] }) }),
    showHelp && /* @__PURE__ */ jsxs("div", { className: "business-entity-help-box", children: [
      /* @__PURE__ */ jsx(
        "button",
        {
          onClick: () => setShowHelp(false),
          className: "business-entity-help-close",
          title: "閉じる",
          children: "×"
        }
      ),
      /* @__PURE__ */ jsx("h3", { style: { marginTop: 0 }, children: "使用方法" }),
      /* @__PURE__ */ jsxs("p", { children: [
        "📝 ",
        /* @__PURE__ */ jsx("strong", { children: "会社情報表示" })
      ] }),
      /* @__PURE__ */ jsxs("p", { children: [
        "🏢 ",
        /* @__PURE__ */ jsx("strong", { children: "業種別分類" })
      ] }),
      /* @__PURE__ */ jsxs("p", { children: [
        "📞 ",
        /* @__PURE__ */ jsx("strong", { children: "連絡先情報" })
      ] }),
      /* @__PURE__ */ jsx("p", { children: "✅ 会社データの取得" }),
      /* @__PURE__ */ jsx("p", { children: "✅ フォルダー名からの自動解析" })
    ] }),
    /* @__PURE__ */ jsxs("div", { className: "business-entity-list-container", children: [
      /* @__PURE__ */ jsxs("div", { className: "business-entity-list-header", children: [
        /* @__PURE__ */ jsx("div", { className: "business-entity-header-company", children: "会社名" }),
        /* @__PURE__ */ jsx("div", { className: "business-entity-header-company", children: "業種" }),
        /* @__PURE__ */ jsx("div", { className: "business-entity-header-spacer" }),
        /* @__PURE__ */ jsx(
          "button",
          {
            onClick: () => setShowHelp(!showHelp),
            className: "business-entity-help-button",
            title: "使用方法を表示",
            children: "?"
          }
        )
      ] }),
      /* @__PURE__ */ jsx("div", { className: "business-entity-scroll-area", children: companies.length === 0 ? /* @__PURE__ */ jsx("div", { className: "business-entity-empty", children: "会社データが見つかりません" }) : /* @__PURE__ */ jsx("div", { children: companies.map((company, index) => /* @__PURE__ */ jsxs(
        "div",
        {
          className: "business-entity-item-row",
          onClick: () => handleCompanyClick(company),
          title: "クリックして詳細表示",
          style: { cursor: "pointer" },
          children: [
            /* @__PURE__ */ jsxs("div", { className: "business-entity-item-info", children: [
              /* @__PURE__ */ jsx("div", { className: "business-entity-item-info-company", children: company.short_name || company.name || "会社名未設定" }),
              /* @__PURE__ */ jsx("div", { className: "business-entity-item-info-company", children: company.business_type || "業種未設定" }),
              /* @__PURE__ */ jsx("div", { className: "business-entity-item-info-spacer" })
            ] }),
            /* @__PURE__ */ jsxs("div", { style: { display: "flex", alignItems: "center", gap: "8px" }, children: [
              company.tags && company.tags.length > 0 && /* @__PURE__ */ jsxs("div", { style: { display: "flex", gap: "4px", flexWrap: "wrap" }, children: [
                company.tags.slice(0, 3).map((tag, tagIndex) => /* @__PURE__ */ jsx(
                  "span",
                  {
                    style: {
                      fontSize: "0.7em",
                      background: "#e3f2fd",
                      color: "#1976d2",
                      padding: "2px 6px",
                      borderRadius: "10px",
                      whiteSpace: "nowrap"
                    },
                    children: tag
                  },
                  tagIndex
                )),
                company.tags.length > 3 && /* @__PURE__ */ jsxs("span", { style: { fontSize: "0.7em", color: "#666" }, children: [
                  "+",
                  company.tags.length - 3
                ] })
              ] }),
              /* @__PURE__ */ jsx("div", { className: "business-entity-item-status business-entity-item-status-completed", children: company.website ? "🌐" : "📁" })
            ] })
          ]
        },
        company.id || index
      )) }) })
    ] }),
    /* @__PURE__ */ jsx(
      CompanyDetailModal,
      {
        isOpen: isDetailModalOpen,
        onClose: closeDetailModal,
        company: selectedCompany,
        onCompanyUpdate: handleCompanyUpdate
      }
    )
  ] });
};
const _layout_companies = UNSAFE_withComponentProps(function CompaniesPage() {
  return /* @__PURE__ */ jsx(Companies, {});
});
const route6 = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
  __proto__: null,
  default: _layout_companies
}, Symbol.toStringTag, { value: "Module" }));
const serverManifest = { "entry": { "module": "/assets/entry.client-BR_5jt8X.js", "imports": ["/assets/chunk-QMGIS6GS-CRLmNMPd.js", "/assets/index-BXUbddt1.js"], "css": [] }, "routes": { "root": { "id": "root", "parentId": void 0, "path": "", "index": void 0, "caseSensitive": void 0, "hasAction": false, "hasLoader": false, "hasClientAction": false, "hasClientLoader": false, "hasClientMiddleware": false, "hasErrorBoundary": false, "module": "/assets/root-CmFF9zS-.js", "imports": ["/assets/chunk-QMGIS6GS-CRLmNMPd.js", "/assets/index-BXUbddt1.js"], "css": [], "clientActionModule": void 0, "clientLoaderModule": void 0, "clientMiddlewareModule": void 0, "hydrateFallbackModule": void 0 }, "routes/_layout": { "id": "routes/_layout", "parentId": "root", "path": void 0, "index": void 0, "caseSensitive": void 0, "hasAction": false, "hasLoader": false, "hasClientAction": false, "hasClientLoader": false, "hasClientMiddleware": false, "hasErrorBoundary": false, "module": "/assets/_layout-CcL890Jn.js", "imports": ["/assets/chunk-QMGIS6GS-CRLmNMPd.js", "/assets/KojiContext-CrloO4lv.js", "/assets/FileInfoContext-NxQDI7MG.js"], "css": ["/assets/_layout-Dt8MPuhg.css"], "clientActionModule": void 0, "clientLoaderModule": void 0, "clientMiddlewareModule": void 0, "hydrateFallbackModule": void 0 }, "routes/_layout._index": { "id": "routes/_layout._index", "parentId": "routes/_layout", "path": void 0, "index": true, "caseSensitive": void 0, "hasAction": false, "hasLoader": false, "hasClientAction": false, "hasClientLoader": false, "hasClientMiddleware": false, "hasErrorBoundary": false, "module": "/assets/_layout._index-DfI2G5t0.js", "imports": ["/assets/chunk-QMGIS6GS-CRLmNMPd.js"], "css": [], "clientActionModule": void 0, "clientLoaderModule": void 0, "clientMiddlewareModule": void 0, "hydrateFallbackModule": void 0 }, "routes/_layout.files": { "id": "routes/_layout.files", "parentId": "routes/_layout", "path": "files", "index": void 0, "caseSensitive": void 0, "hasAction": false, "hasLoader": false, "hasClientAction": false, "hasClientLoader": false, "hasClientMiddleware": false, "hasErrorBoundary": false, "module": "/assets/_layout.files-D4tx_qCC.js", "imports": ["/assets/chunk-QMGIS6GS-CRLmNMPd.js", "/assets/sdk.gen-1y8RMMdH.js", "/assets/FileInfoContext-NxQDI7MG.js", "/assets/useThemeProps-Q2zdfQ03.js", "/assets/index-BXUbddt1.js"], "css": ["/assets/_layout-D3lNBCqO.css"], "clientActionModule": void 0, "clientLoaderModule": void 0, "clientMiddlewareModule": void 0, "hydrateFallbackModule": void 0 }, "routes/_layout.kojies": { "id": "routes/_layout.kojies", "parentId": "routes/_layout", "path": "kojies", "index": void 0, "caseSensitive": void 0, "hasAction": false, "hasLoader": false, "hasClientAction": false, "hasClientLoader": false, "hasClientMiddleware": false, "hasErrorBoundary": false, "module": "/assets/_layout.kojies-DDh5Z7KN.js", "imports": ["/assets/chunk-QMGIS6GS-CRLmNMPd.js", "/assets/sdk.gen-1y8RMMdH.js", "/assets/KojiDetailModal-BnI-BK5R.js", "/assets/KojiContext-CrloO4lv.js", "/assets/index-BXUbddt1.js", "/assets/Close-XszXWwct.js", "/assets/useThemeProps-Q2zdfQ03.js"], "css": ["/assets/business-entity-list-D92NJj2m.css"], "clientActionModule": void 0, "clientLoaderModule": void 0, "clientMiddlewareModule": void 0, "hydrateFallbackModule": void 0 }, "routes/_layout.gantt": { "id": "routes/_layout.gantt", "parentId": "routes/_layout", "path": "kojies/gantt", "index": void 0, "caseSensitive": void 0, "hasAction": false, "hasLoader": false, "hasClientAction": false, "hasClientLoader": false, "hasClientMiddleware": false, "hasErrorBoundary": false, "module": "/assets/_layout.gantt-LxSOLQpK.js", "imports": ["/assets/chunk-QMGIS6GS-CRLmNMPd.js", "/assets/sdk.gen-1y8RMMdH.js", "/assets/KojiDetailModal-BnI-BK5R.js", "/assets/index-BXUbddt1.js", "/assets/Close-XszXWwct.js", "/assets/useThemeProps-Q2zdfQ03.js"], "css": ["/assets/_layout-D9Nv6k78.css"], "clientActionModule": void 0, "clientLoaderModule": void 0, "clientMiddlewareModule": void 0, "hydrateFallbackModule": void 0 }, "routes/_layout.companies": { "id": "routes/_layout.companies", "parentId": "routes/_layout", "path": "companies", "index": void 0, "caseSensitive": void 0, "hasAction": false, "hasLoader": false, "hasClientAction": false, "hasClientLoader": false, "hasClientMiddleware": false, "hasErrorBoundary": false, "module": "/assets/_layout.companies-BOobA9GM.js", "imports": ["/assets/chunk-QMGIS6GS-CRLmNMPd.js", "/assets/sdk.gen-1y8RMMdH.js", "/assets/Close-XszXWwct.js", "/assets/index-BXUbddt1.js"], "css": ["/assets/_layout-CM-RDp9U.css", "/assets/business-entity-list-D92NJj2m.css"], "clientActionModule": void 0, "clientLoaderModule": void 0, "clientMiddlewareModule": void 0, "hydrateFallbackModule": void 0 } }, "url": "/assets/manifest-1d339d91.js", "version": "1d339d91", "sri": void 0 };
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
  "routes/_layout.kojies": {
    id: "routes/_layout.kojies",
    parentId: "routes/_layout",
    path: "kojies",
    index: void 0,
    caseSensitive: void 0,
    module: route4
  },
  "routes/_layout.gantt": {
    id: "routes/_layout.gantt",
    parentId: "routes/_layout",
    path: "kojies/gantt",
    index: void 0,
    caseSensitive: void 0,
    module: route5
  },
  "routes/_layout.companies": {
    id: "routes/_layout.companies",
    parentId: "routes/_layout",
    path: "companies",
    index: void 0,
    caseSensitive: void 0,
    module: route6
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
