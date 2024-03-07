import { useContext } from 'react';

import { EnvContext } from '@shared/components/Providers/EnvProvider';

export const useEnv = () => {
  const env = useContext(EnvContext);

  if (!env) {
    throw new Error('useEnv must be used within an EnvProvider');
  }

  return env;
};
