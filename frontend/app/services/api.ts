import { getFileEntries } from '../api/sdk.gen';
import type { ModelsFileEntry, ModelsFileEntriesListResponse } from '../api/types.gen';

export type FileEntry = ModelsFileEntry;
export type FileEntriesListResponse = ModelsFileEntriesListResponse;

export const folderService = {
  getFileEntries: async (path?: string) => {
    const response = await getFileEntries({ 
      query: path ? { path } : {}
    });
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
  saveKoujiEntries: async (koujiEntries?: any[]) => {
    const { postKoujiSave } = await import('../api/sdk.gen');
    const response = await postKoujiSave({ 
      body: koujiEntries 
    });
    return response.data as unknown as { message: string; output_path: string; count: number };
  },
};