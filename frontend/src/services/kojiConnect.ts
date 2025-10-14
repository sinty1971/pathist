import { create } from "@bufbuild/protobuf";
import { timestampDate, timestampFromDate } from "@bufbuild/protobuf/wkt";
import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import type { ConnectError } from "@connectrpc/connect";
import type {
  ModelsFileInfo,
  ModelsKoji,
  ModelsTimestamp,
} from "@/types/models";
import {
  FileInfoSchema,
  GetKojiRequestSchema,
  KojiSchema,
  ListKojiesRequestSchema,
  UpdateKojiRequestSchema,
  UpdateKojiStandardFilesRequestSchema,
  type FileInfo as FileInfoMessage,
  type Koji as KojiMessage,
  KojiService
} from "@/gen/penguin/v1/penguin_pb";
import type { Timestamp } from "@bufbuild/protobuf/wkt";

const DEFAULT_BASE_URL = "http://localhost:9090";

const baseUrl =
  process.env.NEXT_PUBLIC_CONNECT_BASE_URL ??
  process.env.CONNECT_BASE_URL ??
  DEFAULT_BASE_URL;

const transport = createConnectTransport({
  baseUrl,
  useBinaryFormat: true
});

const client = createClient(KojiService, transport);

const timestampToIso = (value?: Timestamp): string | undefined => {
  if (!value) return undefined;
  const date = timestampDate(value);
  const time = date.getTime();
  if (Number.isNaN(time)) return undefined;
  return date.toISOString();
};

const isoToTimestamp = (value?: string | null) => {
  if (!value) return undefined;
  try {
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) {
      return undefined;
    }
    return timestampFromDate(date);
  } catch {
    return undefined;
  }
};

const normalizeFileInfo = (source: FileInfoMessage): ModelsFileInfo => {
  const targetPath = source.targetPath ?? "";
  return {
    targetPath,
    standardPath: source.standardPath ?? "",
    path: targetPath,
    isDirectory: source.isDirectory,
    size: Number(source.size ?? 0n),
    modifiedTime: timestampToIso(source.modifiedTime),
  };
};

const extractFolderName = (fullPath: string): string => {
  if (!fullPath) return "";
  const parts = fullPath.split(/[/\\]/).filter(Boolean);
  return parts.length > 0 ? parts[parts.length - 1] : "";
};

const normalizeKoji = (source: KojiMessage): ModelsKoji => {
  const startDateIso = source.startDate
    ? timestampToIso(source.startDate)
    : undefined;
  const endDateIso = source.endDate ? timestampToIso(source.endDate) : undefined;

  const requiredFiles = source.requiredFiles?.map(normalizeFileInfo) ?? [];
  const tags = Array.isArray(source.tags)
    ? source.tags.filter((tag): tag is string => Boolean(tag))
    : [];
  const targetFolder = source.targetFolder ?? "";
  const folderName = extractFolderName(targetFolder);

  return {
    id: source.id ?? "",
    status: source.status || undefined,
    targetFolder,
    folderPath: targetFolder,
    folderName,
    name: folderName,
    companyName: source.companyName || "",
    locationName: source.locationName || "",
    description: source.description || "",
    startDate: startDateIso,
    endDate: endDateIso,
    tags,
    requiredFiles
  };
};

const normalizeModelsTimestamp = (
  value: ModelsTimestamp | string | undefined
): string | undefined => {
  if (!value) return undefined;
  if (typeof value === "string") {
    return value;
  }
  if (typeof value === "object" && value !== null && "time.Time" in value) {
    const inner = (value as Record<string, unknown>)["time.Time"];
    return typeof inner === "string" ? inner : undefined;
  }
  return undefined;
};

const toFileInfoMessage = (source: ModelsFileInfo) =>
  create(FileInfoSchema, {
    targetPath: source.targetPath,
    standardPath: source.standardPath ?? "",
    isDirectory: Boolean(source.isDirectory),
    size: BigInt(source.size ?? 0),
    modifiedTime: isoToTimestamp(
      typeof source.modifiedTime === "string" ? source.modifiedTime : undefined
    ),
  });

const toKojiMessage = (source: ModelsKoji) => {
  const startDateIso = normalizeModelsTimestamp(source.startDate);
  const endDateIso = normalizeModelsTimestamp(source.endDate);

  return create(KojiSchema, {
    id: source.id ?? "",
    status: source.status ?? "",
    targetFolder: source.targetFolder ?? "",
    companyName: source.companyName ?? "",
    locationName: source.locationName ?? "",
    description: source.description ?? "",
    startDate: isoToTimestamp(startDateIso),
    endDate: isoToTimestamp(endDateIso),
    tags: Array.isArray(source.tags)
      ? source.tags.filter((tag): tag is string => Boolean(tag))
      : [],
    requiredFiles: Array.isArray(source.requiredFiles)
      ? source.requiredFiles.map(toFileInfoMessage)
      : [],
  });
};

export const kojiConnectClient = {
  list: async (filter?: string): Promise<ModelsKoji[]> => {
    const request = create(ListKojiesRequestSchema, {
      filter: filter ?? ""
    });
    const response = await client.listKojies(request);
    return response.kojies?.map(normalizeKoji) ?? [];
  },

  get: async (folderName: string): Promise<ModelsKoji | null> => {
    try {
      const request = create(GetKojiRequestSchema, { folderName });
      const response = await client.getKoji(request);
      return response.koji ? normalizeKoji(response.koji) : null;
    } catch (error) {
      const connectError = error as ConnectError;
      if (connectError.code === "not_found") {
        return null;
      }
      throw error;
    }
  },

  update: async (koji: ModelsKoji): Promise<ModelsKoji> => {
    const request = create(UpdateKojiRequestSchema, {
      koji: toKojiMessage(koji)
    });
    const response = await client.updateKoji(request);
    if (!response.koji) {
      throw new Error("更新結果の Koji データが空です");
    }
    return normalizeKoji(response.koji);
  },

  updateStandardFiles: async (koji: ModelsKoji, currents: string[]): Promise<ModelsKoji> => {
    const request = create(UpdateKojiStandardFilesRequestSchema, {
      koji: toKojiMessage(koji),
      currents
    });
    const response = await client.updateKojiStandardFiles(request);
    if (!response.koji) {
      throw new Error("標準ファイル更新後の Koji データが空です");
    }
    return normalizeKoji(response.koji);
  }
};

export type KojiConnectClient = typeof kojiConnectClient;
