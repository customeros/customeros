import { observer } from 'mobx-react-lite';
import { NovuProvider } from '@novu/notification-center';

import { useStore } from '@shared/hooks/useStore';

interface ProvidersProps {
  isProduction?: boolean;
  children: React.ReactNode;
}

export const NotificationsProvider = observer(
  ({ children, isProduction }: ProvidersProps) => {
    const store = useStore();

    const id = store.session.value.profile.id ?? 'temp-id';
    const applicationIdentifier = isProduction
      ? import.meta.env.VITE_NOTIFICATION_PROD_APP_IDENTIFIER
      : import.meta.env.VITE_NOTIFICATION_TEST_APP_IDENTIFIER;

    return (
      <NovuProvider
        subscriberId={id}
        applicationIdentifier={applicationIdentifier}
      >
        {children}
      </NovuProvider>
    );
  },
);
