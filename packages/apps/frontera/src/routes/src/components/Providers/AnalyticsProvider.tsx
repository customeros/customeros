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

// eslint-disable-next-line @typescript-eslint/no-explicit-any
declare const window: any;

if (import.meta.env.PROD) {
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
}

export const AnalyticsProvider = observer(
  ({ children }: { isProduction?: boolean; children: React.ReactNode }) => {
    const store = useStore();

    useEffect(() => {
      autorun(() => {
        if (!store.session.isAuthenticated) return;

        const id = store.session.value.profile.id;
        const email = store.session.value.profile.email;
        const name = store.session.value.profile.name;

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

        if (import.meta.env.PROD) {
          if (window.Atlas) {
            window.Atlas.call('identify', {
              userId: id,
              name,
              email,
            });
          }
        }
      });
    }, []);

    return <ErrorBoundary>{children}</ErrorBoundary>;
  },
);
