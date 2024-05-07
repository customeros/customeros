import { useEffect } from 'react';

import { autorun } from 'mobx';
import { observer } from 'mobx-react-lite';
import { GrowthBook } from '@growthbook/growthbook-react';
import { GrowthBookProvider } from '@growthbook/growthbook-react';

import { useStore } from '@shared/hooks/useStore';

export const growthbook = new GrowthBook({
  enableDevMode: true,
  subscribeToChanges: true,
  trackingCallback: (_experiment, _result) => {
    // TODO: Use real analytics tracking system
  },
});

export const GrowthbookProvider = observer(
  ({ children }: { children: React.ReactNode }) => {
    const store = useStore();

    const tenant = store.sessionStore.value.tenant;
    const id = store.sessionStore.value.profile.id;
    const email = store.sessionStore.value.profile.email;

    useEffect(() => {
      autorun(() => {
        if (store.settingsStore.features.isBootstrapped) {
          growthbook.setFeatures(store.settingsStore.features.values);
        }
      });
    }, []);

    useEffect(() => {
      growthbook.setAttributes({
        id,
        email,
        tenant,
      });
    }, [tenant, id, email]);

    return (
      <GrowthBookProvider growthbook={growthbook}>
        {children}
      </GrowthBookProvider>
    );
  },
);
