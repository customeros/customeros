import { useFeatureIsOn } from '@growthbook/growthbook-react';

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
    <div className='flex'>
      {presentUsers.map(([user, color]) => (
        <UserHexagon
          key={user}
          name={user}
          color={color}
          isCurrent={user === username}
        />
      ))}
    </div>
  );
};
