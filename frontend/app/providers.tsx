'use client';

import type { ReactNode } from 'react';
import { KojiProvider } from '@/contexts/KojiContext';
import { FileInfoProvider } from '@/contexts/FileInfoContext';

export function Providers({ children }: { children: ReactNode }) {
  return (
    <KojiProvider>
      <FileInfoProvider>{children}</FileInfoProvider>
    </KojiProvider>
  );
}
