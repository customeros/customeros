import React, { useEffect, useState } from 'react';
import {
  GetGoogleSettings,
  OAuthUserSettingsInterface,
} from '../../../../../../../../../services/settings/settingsService';
import { useSession } from 'next-auth/react';
import {
  ComposeEmail,
  ComposeEmailProps,
} from '@organization/src/components/Timeline/events/email/compose-email/ComposeEmail';
import { KeymapperClose } from '@ui/form/RichTextEditor/components/keyboardShortcuts/KeymapperClose';
import { EmptyIssueMessage } from '@organization/src/components/Timeline/events/email/compose-email/MissingPermissionsMessage';

interface ComposeEmailContainerProps extends ComposeEmailProps {
  onClose: () => void;
}

export const ComposeEmailContainer: React.FC<ComposeEmailContainerProps> = ({
  onClose,
  ...composeEmailProps
}) => {
  const { data: session } = useSession();
  const [allowSendingEmail, setAllowSendingEmail] = useState<
    boolean | undefined
  >(undefined);

  useEffect(() => {
    if (session) {
      // @ts-expect-error look into it
      GetGoogleSettings(session.user.playerIdentityId)
        .then((res: OAuthUserSettingsInterface) => {
          setAllowSendingEmail(res.gmailSyncEnabled);
        })
        .catch((e) => console.log(e));
    }
  }, [session]);

  if (allowSendingEmail) {
    return (
      <ComposeEmail {...composeEmailProps}>
        <KeymapperClose onClose={onClose} />
      </ComposeEmail>
    );
  }

  if (!allowSendingEmail) {
    return (
      <EmptyIssueMessage
        modal={composeEmailProps.modal}
        onAllowSendingEmail={() => setAllowSendingEmail(true)}
      />
    );
  }

  return null;
};
