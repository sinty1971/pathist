import { getFileEntries, postKoujiEntriesSave } from '../api';
import type { ModelsFileEntry, ModelsFileEntriesListResponse } from '../api/types.gen';

export type FileEntry = ModelsFileEntry;
export type FileEntriesListResponse = ModelsFileEntriesListResponse;

export const folderService = {
  getFileEntries: async (path?: string) => {
    const response = await getFileEntries({ query: { path } });
    console.log('Full API response:', response);
    console.log('Response data:', response.data);
    console.log('Response data type:', typeof response.data);
    
    // レスポンス構造を確認して適切にデータを返す
    if (response.data && typeof response.data === 'object') {
      return response.data as FileEntriesListResponse;
    } else {
      // response自体がデータの場合
      return response as unknown as FileEntriesListResponse;
    }
  },
  saveKoujiEntries: async (path?: string, outputPath?: string) => {
    const response = await postKoujiEntriesSave({ query: { path, output_path: outputPath } });
    return response.data as unknown as { message: string; output_path: string; count: number };
  },
};