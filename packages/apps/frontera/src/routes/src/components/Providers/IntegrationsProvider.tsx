import { observer } from 'mobx-react-lite';
import { IntegrationAppProvider } from '@integration-app/react';

import { useStore } from '@shared/hooks/useStore';

export const IntegrationsProvider = observer(
  ({ children }: { children: React.ReactNode }) => {
    const store = useStore();

    return (
      <IntegrationAppProvider token={store.session.value.integrations_token}>
        {children}
      </IntegrationAppProvider>
    );
  },
);
