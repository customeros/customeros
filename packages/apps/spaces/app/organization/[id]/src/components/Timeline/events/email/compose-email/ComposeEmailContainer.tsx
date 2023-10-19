import React, { useEffect, useState } from 'react';
import {
  GetGoogleSettings,
  OAuthUserSettingsInterface,
} from 'services/settings/settingsService';
import { useSession } from 'next-auth/react';
import {
  ComposeEmail,
  ComposeEmailProps,
} from '@organization/src/components/Timeline/events/email/compose-email/ComposeEmail';
import { KeymapperClose } from '@ui/form/RichTextEditor/components/keyboardShortcuts/KeymapperClose';
import { MissingPermissionsPrompt } from '@organization/src/components/Timeline/shared/EmailPermissionsPrompt/EmailPermissionsPrompt';

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
  >(true);

  useEffect(() => {
    if (!session?.user?.playerIdentityId) return;
    GetGoogleSettings(session.user.playerIdentityId)
      .then((res: OAuthUserSettingsInterface) => {
        setAllowSendingEmail(res.gmailSyncEnabled);
      })
      .catch((e) => {
        // throw toast
      });
  }, [session?.user?.playerIdentityId]);

  if (allowSendingEmail) {
    return (
      <ComposeEmail {...composeEmailProps}>
        <KeymapperClose onClose={onClose} />
      </ComposeEmail>
    );
  }

  if (!allowSendingEmail) {
    return (
      <MissingPermissionsPrompt
        modal={composeEmailProps.modal}
        onAllowSendingEmail={() => setAllowSendingEmail(true)}
      />
    );
  }

  return null;
};
