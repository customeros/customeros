'use client';

import { MenuConfig, MenuProvider } from 'kmenu';

const menuConfig: MenuConfig = {};

export const Providers = ({ children }: { children: React.ReactNode }) => {
  return <MenuProvider config={menuConfig}>{children}</MenuProvider>;
};
