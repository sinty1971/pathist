var __defProp = Object.defineProperty;
var __defNormalProp = (obj, key, value) => key in obj ? __defProp(obj, key, { enumerable: true, configurable: true, writable: true, value }) : obj[key] = value;
var __publicField = (obj, key, value) => __defNormalProp(obj, typeof key !== "symbol" ? key + "" : key, value);
import { jsx, jsxs } from "react/jsx-runtime";
import { PassThrough } from "node:stream";
import { createReadableStreamFromReadable } from "@react-router/node";
import { ServerRouter, UNSAFE_withComponentProps, Outlet, Meta, Links, ScrollRestoration, Scripts, useLocation, Link, useNavigate } from "react-router";
import { isbot } from "isbot";
import { renderToPipeableStream } from "react-dom/server";
import React, { useState, useCallback, useEffect } from "react";
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
const links = () => [{
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
  links
}, Symbol.toStringTag, { value: "Module" }));
function Navigation() {
  const location = useLocation();
  const getPageTitle = () => {
    switch (location.pathname) {
      case "/":
        return "ファイル一覧";
      case "/gantt":
        return "工程表";
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
          children: "フォルダー一覧"
        }
      ),
      /* @__PURE__ */ jsx(
        Link,
        {
          to: "/gantt",
          className: location.pathname === "/gantt" ? "nav-link active" : "nav-link",
          children: "工程表"
        }
      )
    ] })
  ] }) });
}
const _layout = UNSAFE_withComponentProps(function Layout2() {
  return /* @__PURE__ */ jsxs("div", {
    className: "app",
    children: [/* @__PURE__ */ jsx(Navigation, {}), /* @__PURE__ */ jsx("main", {
      className: "main-content",
      children: /* @__PURE__ */ jsx(Outlet, {})
    })]
  });
});
const route1 = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
  __proto__: null,
  default: _layout
}, Symbol.toStringTag, { value: "Module" }));
const jsonBodySerializer = {
  bodySerializer: (body) => JSON.stringify(
    body,
    (key, value) => typeof value === "bigint" ? value.toString() : value
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
  var _a;
  if (!contentType) {
    return "stream";
  }
  const cleanContent = (_a = contentType.split(";")[0]) == null ? void 0 : _a.trim();
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
  var _a;
  const config = { ...a, ...b };
  if ((_a = config.baseUrl) == null ? void 0 : _a.endsWith("/")) {
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
  constructor() {
    __publicField(this, "_fns");
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
      if (parseAs === "stream") {
        return opts.responseStyle === "data" ? response.body : {
          data: response.body,
          ...result
        };
      }
      let data = await response[parseAs]();
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
  return ((options == null ? void 0 : options.client) ?? client).get({
    url: "/file/fileinfos",
    ...options
  });
};
const getProjectRecent = (options) => {
  return ((options == null ? void 0 : options.client) ?? client).get({
    url: "/project/recent",
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
const FileInfoModal = ({ fileInfo, isOpen, onClose }) => {
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
const FileInfoGrid = () => {
  const navigate = useNavigate();
  const [fileInfos, setFileEntries] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [currentPath, setCurrentPath] = useState("");
  const [pathInput, setPathInput] = useState("");
  const [selectedFileInfo, setSelectedFileInfo] = useState(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const convertToRelativePath = (frontendPath) => {
    if (!frontendPath || frontendPath === "~/penguin" || frontendPath === "/home/shin/penguin") {
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
  const convertToDisplayPath = (relativePath) => {
    if (!relativePath) {
      return "~/penguin";
    }
    return `~/penguin/${relativePath}`;
  };
  const loadFileEntries = useCallback(async (path) => {
    const frontendPath = path || "";
    const relativePath = convertToRelativePath(frontendPath);
    setLoading(true);
    setError(null);
    try {
      console.log("Loading file entries for frontend path:", frontendPath);
      console.log("Converted to relative path:", relativePath);
      console.log("Calling API with query:", relativePath ? { path: relativePath } : {});
      const response = await getFileFileinfos({
        query: relativePath ? { path: relativePath } : {}
      });
      console.log("API response:", response);
      if (response.data) {
        const data = response.data;
        console.log("Received data:", data);
        setFileEntries(Array.isArray(data) ? data : []);
        setCurrentPath(frontendPath);
      } else if (response.error) {
        console.error("API returned error:", response.error);
        throw new Error("APIエラー: " + JSON.stringify(response.error));
      }
    } catch (err) {
      console.error("Error loading file entries:", err);
      setError(err instanceof Error ? err.message : "エラーが発生しました");
    } finally {
      setLoading(false);
    }
  }, [navigate]);
  useEffect(() => {
    loadFileEntries();
  }, [loadFileEntries]);
  const handleFileInfoClick = (fileInfo) => {
    if (fileInfo.is_directory) {
      const displayPath = convertToDisplayPath(convertToRelativePath(fileInfo.path || ""));
      setPathInput(displayPath);
      loadFileEntries(displayPath);
    } else {
      setSelectedFileInfo(fileInfo);
      setIsModalOpen(true);
    }
  };
  const handlePathSubmit = (e) => {
    e.preventDefault();
    const minPath = "~/penguin";
    if (pathInput.startsWith(minPath) || pathInput === minPath || pathInput === "") {
      loadFileEntries(pathInput);
    } else {
      setPathInput(minPath);
      loadFileEntries(minPath);
    }
  };
  const handleGoBack = () => {
    const pathParts = currentPath.split("/");
    if (pathParts.length > 1) {
      const parentPath = pathParts.slice(0, -1).join("/");
      const newPath = parentPath || "~/penguin";
      const minPath = "~/penguin";
      if (newPath.startsWith(minPath) || newPath === minPath) {
        setPathInput(newPath);
        loadFileEntries(newPath);
      }
    }
  };
  const getFileInfoIcon = (fileInfo) => {
    var _a, _b;
    if (fileInfo.is_directory) {
      return "📁";
    }
    const ext = (_b = (_a = fileInfo.name) == null ? void 0 : _a.split(".").pop()) == null ? void 0 : _b.toLowerCase();
    switch (ext) {
      case "pdf":
        return "📄";
      case "jpg":
      case "jpeg":
      case "png":
      case "gif":
        return "🖼️";
      case "mp4":
      case "avi":
      case "mov":
        return "🎬";
      case "mp3":
      case "wav":
        return "🎵";
      default:
        return "📄";
    }
  };
  console.log("FileInfoGrid render:", {
    loading,
    error,
    fileInfosCount: fileInfos.length,
    pathInput,
    currentPath
  });
  return /* @__PURE__ */ jsxs("div", { className: "folder-container", children: [
    /* @__PURE__ */ jsx("div", { className: "header", children: /* @__PURE__ */ jsxs("form", { onSubmit: handlePathSubmit, className: "path-form", children: [
      /* @__PURE__ */ jsx("button", { type: "button", onClick: handleGoBack, className: "back-button", children: /* @__PURE__ */ jsx("span", { className: "back-arrow", children: "⮜" }) }),
      /* @__PURE__ */ jsx(
        "input",
        {
          type: "text",
          value: pathInput,
          onChange: (e) => setPathInput(e.target.value),
          placeholder: "フォルダーパスを入力",
          className: "path-input"
        }
      ),
      /* @__PURE__ */ jsx("button", { type: "submit", className: "load-button", children: "読み込み" })
    ] }) }),
    /* @__PURE__ */ jsxs("div", { className: "folder-info", children: [
      /* @__PURE__ */ jsxs("span", { className: "folder-count", children: [
        fileInfos.length,
        " 項目"
      ] }),
      /* @__PURE__ */ jsx("span", { className: "current-path", children: currentPath || "~/penguin" })
    ] }),
    loading && /* @__PURE__ */ jsx("div", { className: "loading", children: "読み込み中..." }),
    error && /* @__PURE__ */ jsx("div", { className: "error", children: error }),
    /* @__PURE__ */ jsx("div", { className: "folder-list", children: fileInfos.map((fileInfo, index) => {
      return /* @__PURE__ */ jsxs(
        "div",
        {
          className: "folder-item",
          onClick: () => handleFileInfoClick(fileInfo),
          children: [
            /* @__PURE__ */ jsx("div", { className: "folder-icon", children: getFileInfoIcon(fileInfo) }),
            /* @__PURE__ */ jsxs("div", { className: "folder-info", children: [
              /* @__PURE__ */ jsx("div", { className: "folder-name", children: fileInfo.name }),
              /* @__PURE__ */ jsxs("div", { className: "folder-meta", children: [
                /* @__PURE__ */ jsx("span", { children: fileInfo.is_directory ? "フォルダー" : "ファイル" }),
                /* @__PURE__ */ jsxs("span", { className: "folder-date", children: [
                  " · 更新: ",
                  timestampToString(fileInfo.modified_time) ? new Date(timestampToString(fileInfo.modified_time)).toLocaleDateString("ja-JP", {
                    year: "numeric",
                    month: "2-digit",
                    day: "2-digit",
                    hour: "2-digit",
                    minute: "2-digit"
                  }) : "-"
                ] })
              ] })
            ] })
          ]
        },
        index
      );
    }) }),
    /* @__PURE__ */ jsx(
      FileInfoModal,
      {
        fileInfo: selectedFileInfo,
        isOpen: isModalOpen,
        onClose: () => setIsModalOpen(false)
      }
    )
  ] });
};
const _layout__index = UNSAFE_withComponentProps(function Index() {
  return /* @__PURE__ */ jsx(FileInfoGrid, {});
});
const route2 = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
  __proto__: null,
  default: _layout__index
}, Symbol.toStringTag, { value: "Module" }));
const ProjectEditModal = ({ isOpen, onClose, project, onSave }) => {
  const [formData, setFormData] = useState({
    id: "",
    company_name: "",
    location_name: "",
    description: "",
    tags: "",
    start_date: "",
    end_date: ""
  });
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);
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
      console.log("Project dates extracted:", { startDate, endDate });
      setFormData({
        id: project.id || "",
        company_name: project.company_name || "",
        location_name: project.location_name || "",
        description: project.description || "",
        tags: Array.isArray(project.tags) ? project.tags.join(", ") : project.tags || "",
        start_date: startDate,
        end_date: endDate
      });
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
  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!project) return;
    setIsLoading(true);
    setError(null);
    try {
      const updatedProject = {
        ...project,
        company_name: formData.company_name,
        location_name: formData.location_name,
        description: formData.description,
        tags: formData.tags ? formData.tags.split(",").map((tag) => tag.trim()).filter((tag) => tag.length > 0) : [],
        start_date: formData.start_date ? { "time.Time": `${formData.start_date}T00:00:00+09:00` } : void 0,
        end_date: formData.end_date ? { "time.Time": `${formData.end_date}T23:59:59+09:00` } : void 0
      };
      console.log("Saving project with dates:", updatedProject);
      await onSave(updatedProject);
      onClose();
    } catch (err) {
      console.error("Error saving project:", err);
      setError(err instanceof Error ? err.message : "保存に失敗しました");
    } finally {
      setIsLoading(false);
    }
  };
  if (!isOpen) return null;
  return /* @__PURE__ */ jsx("div", { className: "modal-overlay", onClick: onClose, children: /* @__PURE__ */ jsxs("div", { className: "modal-content", onClick: (e) => e.stopPropagation(), children: [
    /* @__PURE__ */ jsxs("div", { className: "modal-header", children: [
      /* @__PURE__ */ jsx("h2", { children: "工事詳細編集" }),
      /* @__PURE__ */ jsx("button", { type: "button", className: "close-button", onClick: onClose, children: "×" })
    ] }),
    error && /* @__PURE__ */ jsx("div", { className: "error-message", children: error }),
    /* @__PURE__ */ jsxs("form", { onSubmit: handleSubmit, className: "modal-body", children: [
      /* @__PURE__ */ jsxs("div", { className: "form-row", children: [
        /* @__PURE__ */ jsxs("div", { className: "form-group form-group-half", children: [
          /* @__PURE__ */ jsx("label", { htmlFor: "company_name", children: "会社名" }),
          /* @__PURE__ */ jsx(
            "input",
            {
              type: "text",
              id: "company_name",
              name: "company_name",
              value: formData.company_name,
              onChange: handleInputChange,
              className: "path-input",
              required: true
            }
          )
        ] }),
        /* @__PURE__ */ jsxs("div", { className: "form-group form-group-half", children: [
          /* @__PURE__ */ jsx("label", { htmlFor: "location_name", children: "現場名" }),
          /* @__PURE__ */ jsx(
            "input",
            {
              type: "text",
              id: "location_name",
              name: "location_name",
              value: formData.location_name,
              onChange: handleInputChange,
              className: "path-input",
              required: true
            }
          )
        ] })
      ] }),
      /* @__PURE__ */ jsxs("div", { className: "form-row", children: [
        /* @__PURE__ */ jsxs("div", { className: "form-group form-group-half", children: [
          /* @__PURE__ */ jsx("label", { htmlFor: "start_date", children: "開始日" }),
          /* @__PURE__ */ jsx(
            "input",
            {
              type: "date",
              id: "start_date",
              name: "start_date",
              value: formData.start_date,
              onChange: handleInputChange,
              className: "path-input"
            }
          )
        ] }),
        /* @__PURE__ */ jsxs("div", { className: "form-group form-group-half", children: [
          /* @__PURE__ */ jsx("label", { htmlFor: "end_date", children: "終了日" }),
          /* @__PURE__ */ jsx(
            "input",
            {
              type: "date",
              id: "end_date",
              name: "end_date",
              value: formData.end_date,
              onChange: handleInputChange,
              className: "path-input"
            }
          )
        ] })
      ] }),
      /* @__PURE__ */ jsxs("div", { className: "form-group", children: [
        /* @__PURE__ */ jsx("label", { htmlFor: "description", children: "説明" }),
        /* @__PURE__ */ jsx(
          "textarea",
          {
            id: "description",
            name: "description",
            value: formData.description,
            onChange: handleInputChange,
            className: "path-input",
            rows: 2
          }
        )
      ] }),
      /* @__PURE__ */ jsxs("div", { className: "form-group", children: [
        /* @__PURE__ */ jsx("label", { htmlFor: "tags", children: "タグ" }),
        /* @__PURE__ */ jsx(
          "input",
          {
            type: "text",
            id: "tags",
            name: "tags",
            value: formData.tags,
            onChange: handleInputChange,
            className: "path-input",
            placeholder: "カンマ区切りで入力"
          }
        )
      ] }),
      (project == null ? void 0 : project.managed_files) && project.managed_files.length > 0 && /* @__PURE__ */ jsxs("div", { className: "form-group", children: [
        /* @__PURE__ */ jsxs("label", { children: [
          "管理ファイル (",
          project.managed_files.length,
          "件)"
        ] }),
        /* @__PURE__ */ jsxs("div", { className: "managed-files-table", children: [
          /* @__PURE__ */ jsxs("div", { className: "files-header", children: [
            /* @__PURE__ */ jsx("div", { className: "header-item", children: "種別" }),
            /* @__PURE__ */ jsx("div", { className: "header-item", children: "ファイルパス" }),
            /* @__PURE__ */ jsx("div", { className: "header-item", children: "状態" })
          ] }),
          /* @__PURE__ */ jsx("div", { className: "files-body", children: project.managed_files.map((file, index) => /* @__PURE__ */ jsxs(React.Fragment, { children: [
            file.current && /* @__PURE__ */ jsxs("div", { className: "file-row", children: [
              /* @__PURE__ */ jsxs("div", { className: "file-type", children: [
                /* @__PURE__ */ jsx("span", { className: "file-icon", children: "📎" }),
                "現在ファイル"
              ] }),
              /* @__PURE__ */ jsx("div", { className: "file-path", children: file.current }),
              /* @__PURE__ */ jsx("div", { className: "file-status", children: /* @__PURE__ */ jsx("span", { className: "status-current", children: "使用中" }) })
            ] }),
            file.recommended && /* @__PURE__ */ jsxs("div", { className: "file-row", children: [
              /* @__PURE__ */ jsxs("div", { className: "file-type", children: [
                /* @__PURE__ */ jsx("span", { className: "file-icon", children: "💡" }),
                "推奨ファイル"
              ] }),
              /* @__PURE__ */ jsx("div", { className: "file-path", children: file.recommended }),
              /* @__PURE__ */ jsx("div", { className: "file-status", children: /* @__PURE__ */ jsx("span", { className: "status-recommended", children: "推奨" }) })
            ] })
          ] }, index)) })
        ] })
      ] }),
      project && (!project.managed_files || project.managed_files.length === 0) && /* @__PURE__ */ jsxs("div", { className: "form-group", children: [
        /* @__PURE__ */ jsx("label", { children: "管理ファイル" }),
        /* @__PURE__ */ jsxs("div", { className: "no-managed-files", children: [
          /* @__PURE__ */ jsx("span", { className: "no-files-icon", children: "📄" }),
          /* @__PURE__ */ jsx("span", { children: "管理ファイルはありません" })
        ] })
      ] }),
      /* @__PURE__ */ jsxs("div", { className: "form-actions", children: [
        /* @__PURE__ */ jsx(
          "button",
          {
            type: "button",
            className: "cancel-button",
            onClick: onClose,
            disabled: isLoading,
            children: "キャンセル"
          }
        ),
        /* @__PURE__ */ jsx(
          "button",
          {
            type: "submit",
            className: "submit-button",
            disabled: isLoading,
            children: isLoading ? "保存中..." : "保存"
          }
        )
      ] })
    ] })
  ] }) });
};
const ProjectGanttChartSimple = () => {
  const [projects, setProjects] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [selectedProject, setSelectedProject] = useState(null);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const loadProjects = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await getProjectRecent();
      if (response.data) {
        setProjects(response.data);
      } else {
        setProjects([]);
      }
    } catch (err) {
      console.error("Error loading kouji entries:", err);
      setError(`工事データの読み込みに失敗しました: ${err instanceof Error ? err.message : String(err)}`);
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
  const handleProjectSave = async (updatedProject) => {
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
        throw new Error(errorData.message || "保存に失敗しました");
      }
      console.log("プロジェクトが正常に保存されました");
      await loadProjects();
    } catch (err) {
      console.error("Error saving project:", err);
      throw err;
    }
  };
  const closeEditModal = () => {
    setIsEditModalOpen(false);
    setSelectedProject(null);
  };
  if (loading) {
    return /* @__PURE__ */ jsx("div", { style: { padding: "20px" }, children: /* @__PURE__ */ jsx("div", { children: "工事データを読み込み中..." }) });
  }
  if (error) {
    return /* @__PURE__ */ jsxs("div", { style: { padding: "20px" }, children: [
      /* @__PURE__ */ jsx("div", { style: { color: "red", padding: "10px", backgroundColor: "#ffe6e6", borderRadius: "4px" }, children: error }),
      /* @__PURE__ */ jsx("button", { onClick: loadProjects, style: { marginTop: "10px", padding: "10px 20px" }, children: "再試行" })
    ] });
  }
  return /* @__PURE__ */ jsxs("div", { style: { padding: "20px" }, children: [
    /* @__PURE__ */ jsx("div", { style: { marginBottom: "20px" }, children: /* @__PURE__ */ jsxs("p", { children: [
      "取得した工事データ: ",
      projects.length,
      "件"
    ] }) }),
    /* @__PURE__ */ jsxs("div", { style: { border: "1px solid #ddd", borderRadius: "8px", overflow: "hidden" }, children: [
      /* @__PURE__ */ jsx("div", { style: {
        backgroundColor: "#f5f5f5",
        padding: "10px",
        fontWeight: "bold",
        borderBottom: "1px solid #ddd"
      }, children: "工事一覧" }),
      projects.length === 0 ? /* @__PURE__ */ jsx("div", { style: { padding: "20px", textAlign: "center", color: "#666" }, children: "工事データが見つかりません" }) : /* @__PURE__ */ jsx("div", { children: projects.map((project, index) => /* @__PURE__ */ jsxs(
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
            /* @__PURE__ */ jsxs("div", { children: [
              /* @__PURE__ */ jsxs("div", { style: { fontWeight: "bold" }, children: [
                project.company_name || "会社名未設定",
                " - ",
                project.location_name || "現場名未設定"
              ] }),
              /* @__PURE__ */ jsxs("div", { style: { fontSize: "14px", color: "#666", marginTop: "5px" }, children: [
                "開始: ",
                project.start_date ? new Date(project.start_date).toLocaleDateString("ja-JP") : "未設定",
                " | 終了: ",
                project.end_date ? new Date(project.end_date).toLocaleDateString("ja-JP") : "未設定"
              ] })
            ] }),
            /* @__PURE__ */ jsx("div", { style: {
              padding: "4px 12px",
              borderRadius: "4px",
              backgroundColor: project.status === "進行中" ? "#4CAF50" : project.status === "完了" ? "#9E9E9E" : project.status === "予定" ? "#FF9800" : "#2196F3",
              color: "white",
              fontSize: "12px"
            }, children: project.status || "未設定" })
          ]
        },
        project.id || index
      )) })
    ] }),
    /* @__PURE__ */ jsxs("div", { style: { marginTop: "20px", padding: "15px", backgroundColor: "#f0f8ff", borderRadius: "4px" }, children: [
      /* @__PURE__ */ jsx("h3", { children: "使用方法" }),
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
      /* @__PURE__ */ jsx("p", { children: "🔄 ガントチャート機能（次のステップ）" })
    ] }),
    /* @__PURE__ */ jsx(
      ProjectEditModal,
      {
        isOpen: isEditModalOpen,
        onClose: closeEditModal,
        project: selectedProject,
        onSave: handleProjectSave
      }
    )
  ] });
};
const _layout_gantt = UNSAFE_withComponentProps(function GanttChart() {
  return /* @__PURE__ */ jsx(ProjectGanttChartSimple, {});
});
const route3 = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
  __proto__: null,
  default: _layout_gantt
}, Symbol.toStringTag, { value: "Module" }));
const serverManifest = { "entry": { "module": "/assets/entry.client-B8MmGEjj.js", "imports": ["/assets/chunk-NL6KNZEE-DVU-DdLP.js"], "css": [] }, "routes": { "root": { "id": "root", "parentId": void 0, "path": "", "index": void 0, "caseSensitive": void 0, "hasAction": false, "hasLoader": false, "hasClientAction": false, "hasClientLoader": false, "hasClientMiddleware": false, "hasErrorBoundary": false, "module": "/assets/root-CpgVR1Az.js", "imports": ["/assets/chunk-NL6KNZEE-DVU-DdLP.js"], "css": [], "clientActionModule": void 0, "clientLoaderModule": void 0, "clientMiddlewareModule": void 0, "hydrateFallbackModule": void 0 }, "routes/_layout": { "id": "routes/_layout", "parentId": "root", "path": void 0, "index": void 0, "caseSensitive": void 0, "hasAction": false, "hasLoader": false, "hasClientAction": false, "hasClientLoader": false, "hasClientMiddleware": false, "hasErrorBoundary": false, "module": "/assets/_layout-gjhVmDm9.js", "imports": ["/assets/chunk-NL6KNZEE-DVU-DdLP.js"], "css": ["/assets/_layout-B_lvcXFH.css"], "clientActionModule": void 0, "clientLoaderModule": void 0, "clientMiddlewareModule": void 0, "hydrateFallbackModule": void 0 }, "routes/_layout._index": { "id": "routes/_layout._index", "parentId": "routes/_layout", "path": void 0, "index": true, "caseSensitive": void 0, "hasAction": false, "hasLoader": false, "hasClientAction": false, "hasClientLoader": false, "hasClientMiddleware": false, "hasErrorBoundary": false, "module": "/assets/_layout._index-CINbqJPU.js", "imports": ["/assets/chunk-NL6KNZEE-DVU-DdLP.js", "/assets/sdk.gen-ZuoAz7BR.js"], "css": [], "clientActionModule": void 0, "clientLoaderModule": void 0, "clientMiddlewareModule": void 0, "hydrateFallbackModule": void 0 }, "routes/_layout.gantt": { "id": "routes/_layout.gantt", "parentId": "routes/_layout", "path": "gantt", "index": void 0, "caseSensitive": void 0, "hasAction": false, "hasLoader": false, "hasClientAction": false, "hasClientLoader": false, "hasClientMiddleware": false, "hasErrorBoundary": false, "module": "/assets/_layout.gantt-7Iv8td-9.js", "imports": ["/assets/chunk-NL6KNZEE-DVU-DdLP.js", "/assets/sdk.gen-ZuoAz7BR.js"], "css": [], "clientActionModule": void 0, "clientLoaderModule": void 0, "clientMiddlewareModule": void 0, "hydrateFallbackModule": void 0 } }, "url": "/assets/manifest-a49c986e.js", "version": "a49c986e", "sri": void 0 };
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
  "routes/_layout.gantt": {
    id: "routes/_layout.gantt",
    parentId: "routes/_layout",
    path: "gantt",
    index: void 0,
    caseSensitive: void 0,
    module: route3
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
