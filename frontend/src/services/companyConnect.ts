import { create } from "@bufbuild/protobuf";
import { timestampDate } from "@bufbuild/protobuf/wkt";
import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import {
  CompanySchema,
  CompanyService,
  FileInfoSchema,
  ListCompaniesRequestSchema,
  ListCompanyCategoriesRequestSchema,
  UpdateCompanyRequestSchema,
  type Company as CompanyMessage,
  type CompanyCategoryInfo as CompanyCategoryMessage,
  type FileInfo as FileInfoMessage,
} from "@/gen/penguin/v1/penguin_pb";
import type {
  Company,
  CompanyCategoryInfo,
  ModelsFileInfo,
} from "@/types/models";
import type { Timestamp } from "@bufbuild/protobuf/wkt";

const DEFAULT_BASE_URL = "http://localhost:9090";

const baseUrl =
  process.env.NEXT_PUBLIC_CONNECT_BASE_URL ??
  process.env.CONNECT_BASE_URL ??
  DEFAULT_BASE_URL;

const transport = createConnectTransport({
  baseUrl,
  useBinaryFormat: true,
});

const client = createClient(CompanyService, transport);

const toPathLike = (value?: string): string => value ?? "";

const timestampToIso = (value?: Timestamp): string | undefined => {
  if (!value) return undefined;
  const date = timestampDate(value);
  const time = date.getTime();
  if (Number.isNaN(time)) return undefined;
  return date.toISOString();
};

const normalizeFileInfo = (source: FileInfoMessage): ModelsFileInfo => {
  const targetPath = toPathLike(source.targetPath);
  return {
    targetPath,
    standardPath: toPathLike(source.standardPath),
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

const normalizeCompany = (source: CompanyMessage): Company => {
  const targetFolder = toPathLike(source.targetFolder);
  return {
    id: toPathLike(source.id),
    folderName: extractFolderName(targetFolder),
    targetFolder,
    shortName: toPathLike(source.shortName) || undefined,
    category: toPathLike(source.category),
    legalName: toPathLike(source.legalName) || undefined,
    postalCode: toPathLike(source.postalCode) || undefined,
    address: toPathLike(source.address) || undefined,
    phone: toPathLike(source.phone) || undefined,
    email: toPathLike(source.email) || undefined,
    website: toPathLike(source.website) || undefined,
    tags: source.tags ? [...source.tags].filter(Boolean) : [],
    requiredFiles: source.requiredFiles
      ? source.requiredFiles.map(normalizeFileInfo)
      : [],
  };
};

const toCompanyMessage = (source: Company) =>
  create(CompanySchema, {
    id: toPathLike(source.id),
    targetFolder: toPathLike(source.targetFolder),
    shortName: source.shortName ?? "",
    category: toPathLike(source.category),
    legalName: source.legalName ?? "",
    postalCode: source.postalCode ?? "",
    address: source.address ?? "",
    phone: source.phone ?? "",
    email: source.email ?? "",
    website: source.website ?? "",
    tags: Array.isArray(source.tags) ? [...source.tags].filter(Boolean) : [],
    requiredFiles: Array.isArray(source.requiredFiles)
      ? source.requiredFiles.map((file) =>
          create(FileInfoSchema, {
            targetPath: file.targetPath,
            standardPath: file.standardPath ?? "",
            isDirectory: Boolean(file.isDirectory),
            size: BigInt(file.size ?? 0),
          })
        )
      : [],
  });

const normalizeCategory = (
  category: CompanyCategoryMessage
): CompanyCategoryInfo => ({
  code: toPathLike(category.code),
  label: toPathLike(category.label) || "業種未設定",
});

export const companyConnectClient = {
  list: async (): Promise<Company[]> => {
    const request = create(ListCompaniesRequestSchema, {});
    const response = await client.listCompanies(request);
    return response.companies?.map(normalizeCompany) ?? [];
  },

  update: async (company: Company): Promise<Company> => {
    const request = create(UpdateCompanyRequestSchema, {
      company: toCompanyMessage(company),
    });
    const response = await client.updateCompany(request);
    if (!response.company) {
      throw new Error("更新結果の Company データが空です");
    }
    return normalizeCompany(response.company);
  },

  listCategories: async (): Promise<CompanyCategoryInfo[]> => {
    const request = create(ListCompanyCategoriesRequestSchema, {});
    const response = await client.listCompanyCategories(request);
    return response.categories?.map(normalizeCategory) ?? [];
  },

  get: async (id: string): Promise<Company | null> => {
    const list = await companyConnectClient.list();
    return list.find((item) => item.id === id) ?? null;
  },
};

export type CompanyConnectClient = typeof companyConnectClient;
