'use client';

import { createContext } from 'react';

export type Env = {
  PRODUCTION: string;
  REALTIME_WS_PATH: string;
  NOTIFICATION_PROD_APP_IDENTIFIER: string;
  NOTIFICATION_TEST_APP_IDENTIFIER: string;
};

export const EnvContext = createContext<Env>({
  PRODUCTION: '',
  REALTIME_WS_PATH: '',
  NOTIFICATION_PROD_APP_IDENTIFIER: '',
  NOTIFICATION_TEST_APP_IDENTIFIER: '',
});

interface EnvProviderProps {
  env: Env;
  children: React.ReactNode;
}

export const EnvProvider = ({ children, env }: EnvProviderProps) => {
  return <EnvContext.Provider value={env}>{children}</EnvContext.Provider>;
};
