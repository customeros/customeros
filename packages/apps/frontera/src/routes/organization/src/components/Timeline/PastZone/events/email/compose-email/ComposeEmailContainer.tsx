import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { MissingPermissionsPrompt } from '@organization/components/Timeline/shared/EmailPermissionsPrompt/EmailPermissionsPrompt';
import {
  ComposeEmail,
  ComposeEmailProps,
} from '@organization/components/Timeline/PastZone/events/email/compose-email/ComposeEmail';

interface ComposeEmailContainerProps extends ComposeEmailProps {
  onClose: () => void;
}

export const ComposeEmailContainer = observer(
  ({ onClose, ...composeEmailProps }: ComposeEmailContainerProps) => {
    const store = useStore();
    const allowSendingEmail =
      (store.globalCache?.value?.activeEmailTokens &&
        store.globalCache?.value?.activeEmailTokens?.length > 0) ||
      (store.globalCache?.value?.mailboxes &&
        store.globalCache?.value?.mailboxes?.length > 0);

    if (allowSendingEmail) {
      return <ComposeEmail {...composeEmailProps} />;
    }

    if (!allowSendingEmail) {
      return <MissingPermissionsPrompt modal={composeEmailProps.modal} />;
    }

    return null;
  },
);
