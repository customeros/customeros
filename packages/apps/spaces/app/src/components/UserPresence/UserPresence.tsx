import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { HStack } from '@ui/layout/Stack';
import { useChannel } from '@shared/hooks/useChannel';

import { UserHexagon } from '../UserHexagon';

interface UserPresenceProps {
  channelName: string;
}

export const UserPresence = ({ channelName }: UserPresenceProps) => {
  const { presentUsers, username } = useChannel(channelName);
  const isPresenceEnabled = useFeatureIsOn('presence');

  if (!isPresenceEnabled) return null;

  return (
    <HStack>
      {presentUsers.map(([user, color]) => (
        <UserHexagon
          key={user}
          name={user}
          color={color}
          isCurrent={user === username}
        />
      ))}
    </HStack>
  );
};
