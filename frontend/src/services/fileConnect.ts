import { create } from "@bufbuild/protobuf";
import { timestampDate } from "@bufbuild/protobuf/wkt";
import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import type { ModelsFileInfo } from "@/types/models";
import {
  FileService,
  ListFilesRequestSchema,
  GetFileBasePathRequestSchema,
  type FileInfo as FileInfoMessage
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

const client = createClient(FileService, transport);

const timestampToIso = (value?: Timestamp): string | undefined => {
  if (!value) return undefined;
  const date = timestampDate(value);
  const time = date.getTime();
  if (Number.isNaN(time)) return undefined;
  return date.toISOString();
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

export const fileConnectClient = {
  list: async (path?: string): Promise<ModelsFileInfo[]> => {
    const request = create(ListFilesRequestSchema, path ? { path } : {});
    const response = await client.listFiles(request);
    return response.files?.map(normalizeFileInfo) ?? [];
  },

  getBasePath: async (): Promise<string> => {
    const response = await client.getFileBasePath(
      create(GetFileBasePathRequestSchema, {})
    );
    return response.basePath ?? "";
  }
};

export type FileConnectClient = typeof fileConnectClient;
