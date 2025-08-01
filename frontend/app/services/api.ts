import { getFileFileinfos } from '@/api/sdk.gen';
import type { ModelsFileInfo } from '@/api/types.gen';

export type FileInfo = ModelsFileInfo;

export const fileService = {
  getFileInfos: async (path?: string) => {
    const response = await getFileFileinfos({ 
      query: path ? { path } : {}
    });
    console.log('Full API response:', response);
    console.log('Response data:', response.data);
    console.log('Response data type:', typeof response.data);
    
    // レスポンス構造を確認して適切にデータを返す
    if (response.data && Array.isArray(response.data)) {
      return response.data as ModelsFileInfo[];
    } else {
      return [] as ModelsFileInfo[];
    }
  },
};