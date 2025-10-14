export type ModelsTimestamp = string | { "time.Time": string };

export interface ModelsFileInfo {
  targetPath: string;
  standardPath?: string;
  path?: string;
  isDirectory: boolean;
  size: number;
  modifiedTime?: string;
}

export interface ModelsManagedFile {
  current?: string;
  desired?: string;
  recommended?: string;
}

export interface Company {
  id: string;
  folderName: string;
  targetFolder: string;
  shortName?: string;
  category: string;
  legalName?: string;
  phone?: string;
  email?: string;
  website?: string;
  address?: string;
  postalCode?: string;
  tags: string[];
  requiredFiles: ModelsFileInfo[];
}

export interface CompanyCategoryInfo {
  code: string;
  label: string;
}

export interface ModelsKoji {
  id: string;
  status?: string;
  targetFolder: string;
  folderPath: string;
  folderName?: string;
  name: string;
  companyName: string;
  locationName: string;
  description?: string;
  startDate?: ModelsTimestamp | string;
  endDate?: ModelsTimestamp | string;
  tags: string[];
  requiredFiles: ModelsFileInfo[];
  managed_files?: ModelsManagedFile[];
  assists?: ModelsManagedFile[];
}
