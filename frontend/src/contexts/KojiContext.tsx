'use client';

import { createContext, useContext, useState } from 'react';
import type { ReactNode } from 'react';

interface KojiContextType {
  kojiCount: number;
  setKojiCount: (count: number) => void;
}

const KojiContext = createContext<KojiContextType | undefined>(undefined);

export function KojiProvider({ children }: { children: ReactNode }) {
  const [kojiCount, setKojiCount] = useState<number>(0);

  return (
    <KojiContext.Provider value={{ kojiCount, setKojiCount }}>
      {children}
    </KojiContext.Provider>
  );
}

export function useKoji() {
  const context = useContext(KojiContext);
  if (context === undefined) {
    throw new Error('useKoji must be used within a KojiProvider');
  }
  return context;
}
