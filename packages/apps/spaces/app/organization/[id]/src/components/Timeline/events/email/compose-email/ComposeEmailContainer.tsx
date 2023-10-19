import React from 'react';

import {
  ComposeEmail,
  ComposeEmailProps,
} from '@organization/src/components/Timeline/events/email/compose-email/ComposeEmail';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';
import { KeymapperClose } from '@ui/form/RichTextEditor/components/keyboardShortcuts/KeymapperClose';
import { MissingPermissionsPrompt } from '@organization/src/components/Timeline/shared/EmailPermissionsPrompt/EmailPermissionsPrompt';

interface ComposeEmailContainerProps extends ComposeEmailProps {
  onClose: () => void;
}

export const ComposeEmailContainer: React.FC<ComposeEmailContainerProps> = ({
  onClose,
  ...composeEmailProps
}) => {
  const client = getGraphQLClient();
  const { data: globalCache } = useGlobalCacheQuery(client);
  const allowSendingEmail = globalCache?.global_Cache?.isGoogleActive;

  if (allowSendingEmail) {
    return (
      <ComposeEmail {...composeEmailProps}>
        <KeymapperClose onClose={onClose} />
      </ComposeEmail>
    );
  }

  if (!allowSendingEmail) {
    return <MissingPermissionsPrompt modal={composeEmailProps.modal} />;
  }

  return null;
};
