import { useEffect } from 'react';

import { autorun } from 'mobx';
import { H } from 'highlight.run';
import { observer } from 'mobx-react-lite';
import { ErrorBoundary } from '@highlight-run/react';

import { useStore } from '@shared/hooks/useStore';

declare const heap: {
  identify: (email: string | null | undefined) => void;
  addUserProperties: (
    properties: Record<string, string | null | undefined>,
  ) => void;
};

H.init('ldwno7wd', {
  environment: import.meta.env.MODE,
  serviceName: 'customer-os',
  tracingOrigins: true,
  networkRecording: {
    enabled: import.meta.env.MODE === 'production',
    recordHeadersAndBody: true,
    urlBlocklist: [],
  },
});

export const AnalyticsProvider = observer(
  ({ children }: { isProduction?: boolean; children: React.ReactNode }) => {
    const { sessionStore } = useStore();

    useEffect(() => {
      autorun(() => {
        if (!sessionStore.isAuthenticated) return;

        const id = sessionStore.value.profile.id;
        const email = sessionStore.value.profile.email;
        const name = sessionStore.value.profile.name;

        if (import.meta.env.PROD && typeof heap !== 'undefined') {
          heap.identify(email);
          heap.addUserProperties({
            name,
          });
        }
        if (import.meta.env.PROD) {
          H.identify(email, {
            name,
            playerIdentityId: id,
          });
        }
      });
    }, []);

    return <ErrorBoundary>{children}</ErrorBoundary>;
  },
);
